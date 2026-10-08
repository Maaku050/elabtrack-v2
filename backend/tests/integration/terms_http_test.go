package integration

import (
	"bytes"
	"encoding/json"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/google/uuid"
	"io"
	"net/http"
	"testing"
)

func TestRealTermsHTTP(t *testing.T) {
	ctx, runtime, _, cfg := termsDatabase(t)
	borrower := termsUser(t, ctx, runtime, user.RoleBorrower)
	admin := termsUser(t, ctx, runtime, user.RoleAdmin)
	staff := termsUser(t, ctx, runtime, user.RoleStaff)
	issuer := security.NewJWTIssuer(cfg.JWT)
	token := func(u *user.User) string {
		t.Helper()
		raw, e := issuer.IssueAccessToken(ctx, u)
		require(t, e == nil, "test access issuance")
		return raw
	}
	bt, at, st := token(borrower), token(admin), token(staff)
	request := func(method, path, body, access, origin string) (int, map[string]any) {
		t.Helper()
		req, e := http.NewRequest(method, "http://localhost:18084/api/v1"+path, bytes.NewBufferString(body))
		require(t, e == nil, "request construction")
		req.Header.Set("Content-Type", "application/json")
		if access != "" {
			req.Header.Set("Authorization", "Bearer "+access)
		}
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		res, e := http.DefaultClient.Do(req)
		require(t, e == nil, "real terms HTTP")
		defer res.Body.Close()
		raw, e := io.ReadAll(res.Body)
		require(t, e == nil, "response read")
		var envelope map[string]any
		require(t, json.Unmarshal(raw, &envelope) == nil, "safe envelope")
		if origin == "" || origin == "http://localhost:15174" {
			require(t, res.Header.Get("Cache-Control") == "no-store", "private terms/consent responses never cache")
		}
		require(t, len(res.Cookies()) == 0, "terms never changes session cookies")
		require(t, res.Header.Get("X-Request-ID") != "", "server request correlation")
		if res.StatusCode >= 400 {
			fail, ok := envelope["error"].(map[string]any)
			require(t, ok && fail["requestId"] == res.Header.Get("X-Request-ID"), "safe correlated terms failure")
		}
		return res.StatusCode, envelope
	}
	var current string
	status, body := request("GET", "/terms/current", "", bt, "")
	require(t, status == 200, "published isolated current terms")
	current = body["data"].(map[string]any)["id"].(string)
	acceptPath := "/terms/" + current + "/accept"
	t.Run("authentication_and_role_matrix", func(t *testing.T) {
		for _, p := range []string{"/terms/current", "/terms/status"} {
			status, _ := request("GET", p, "", "", "")
			require(t, status == 401, "anonymous denial")
		}
		for _, access := range []string{st, at} {
			status, _ := request("POST", acceptPath, "{}", access, "http://localhost:15174")
			require(t, status == 403, "staff/admin cannot accept borrower terms")
		}
		for _, access := range []string{bt, st} {
			status, _ := request("POST", "/terms/versions", `{"version":"TEST-denied","title":"SYNTHETIC","body":"TEST ONLY","expected_current_version_id":null}`, access, "http://localhost:15174")
			require(t, status == 403, "publication Admin-only")
		}
	})
	t.Run("strict_identity_and_origin_boundary", func(t *testing.T) {
		for _, payload := range []string{"null", "[]", `{"user_id":"` + admin.ID.String() + `"}`, `{"accepted_at":"2000-01-01T00:00:00Z"}`, `{"accepted":true}`, `{} {}`} {
			status, _ := request("POST", acceptPath, payload, bt, "http://localhost:15174")
			require(t, status == 400, "cannot choose account/time or bypass binder")
		}
		for _, origin := range []string{"", "null", "https://untrusted.example.invalid"} {
			status, _ := request("POST", acceptPath, "{}", bt, origin)
			require(t, status == 403, "trusted Origin mandatory")
		}
		var count int
		require(t, runtime.Pool.QueryRow(ctx, `SELECT count(*) FROM terms_acceptances WHERE user_id=$1`, borrower.ID).Scan(&count) == nil && count == 0, "rejected bodies produced no consent")
		status, _ := request("POST", "/terms/"+uuid.NewString()+"/accept", "{}", bt, "http://localhost:15174")
		require(t, status == 404, "unknown version safe failure")
	})
	t.Run("current_acceptance_repeat_and_version_conflict", func(t *testing.T) {
		first, body := request("POST", acceptPath, "{}", bt, "http://localhost:15174")
		require(t, first == 200, "authenticated acceptance")
		old := body["data"].(map[string]any)
		second, body := request("POST", acceptPath, "{}", bt, "http://localhost:15174")
		now := body["data"].(map[string]any)
		require(t, second == 200 && old["id"] == now["id"] && old["accepted_at"] == now["accepted_at"], "repeat original receipt")
		status, body := request("GET", "/terms/status", "", bt, "")
		require(t, status == 200 && body["data"].(map[string]any)["state"] == "accepted", "server current acceptance")
		payload, _ := json.Marshal(map[string]any{"version": "TEST-HTTP-" + uuid.NewString(), "title": "SYNTHETIC TEST TERMS — NOT OFFICIAL", "body": "TEST ONLY. Isolated HTTP publication.", "expected_current_version_id": current})
		status, body = request("POST", "/terms/versions", string(payload), at, "http://localhost:15174")
		require(t, status == 201, "Admin HTTP publication")
		status, body = request("POST", acceptPath, "{}", bt, "http://localhost:15174")
		require(t, status == 409 && body["error"].(map[string]any)["code"] == "TERMS_VERSION_CHANGED", "old reviewed version conflicts")
		status, body = request("GET", "/terms/status", "", bt, "")
		require(t, status == 200 && body["data"].(map[string]any)["state"] == "updated", "outdated acceptance remains evidence")
		status, _ = request("POST", "/terms/versions", string(payload), at, "http://localhost:15174")
		require(t, status == 409, "optimistic publication conflict")
	})
	t.Run("current_database_role_and_inactive_after_issuance", func(t *testing.T) {
		_, err := runtime.Pool.Exec(ctx, `UPDATE users SET role='STAFF' WHERE id=$1`, admin.ID)
		require(t, err == nil, "owned demotion")
		status, _ := request("POST", "/terms/versions", `{"version":"TEST-stale","title":"SYNTHETIC","body":"TEST ONLY","expected_current_version_id":null}`, at, "http://localhost:15174")
		require(t, status == 403, "stale Admin JWT cannot publish")
		_, err = runtime.Pool.Exec(ctx, `UPDATE users SET is_active=false WHERE id=$1`, borrower.ID)
		require(t, err == nil, "owned disable")
		for _, p := range []string{"/terms/current", "/terms/status"} {
			status, _ := request("GET", p, "", bt, "")
			require(t, status == 403, "disabled account token denied")
		}
		status, _ = request("POST", acceptPath, "{}", bt, "http://localhost:15174")
		require(t, status == 403, "disabled account cannot accept")
	})
}
