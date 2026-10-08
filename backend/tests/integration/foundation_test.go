package integration

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	appauth "github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/google/uuid"
)

// Explicit opt-in; the supplied Compose stack must already be migrated/running.
// This test never starts Docker, migrates, or reads an ambient DATABASE_URL.
func TestRealFoundation(t *testing.T) {
	batch1 := os.Getenv("ELABTRACK_BATCH1") == "1"
	phase4b := os.Getenv("ELABTRACK_PHASE4B") == "1"
	phase4a := os.Getenv("ELABTRACK_PHASE4A") == "1"
	if os.Getenv("ELABTRACK_INTEGRATION") != "1" && !phase4a && !phase4b && !batch1 {
		t.Skip("requires isolated integration Compose stack")
	}
	port, dbName, baseURL, origin := "15432", "elabtrack_v2_integration", "http://localhost:18080/api/v1", "http://localhost:15173"
	if phase4a {
		port, dbName, baseURL, origin = "25432", "elabtrack_v2_phase4a_test", "http://localhost:8080/api/v1", "http://localhost:5173"
	}
	if phase4b {
		port, dbName, baseURL, origin = "35432", "elabtrack_v2_phase4b_test", "http://localhost:18084/api/v1", "http://localhost:15174"
	}
	if batch1 {
		port, dbName, baseURL, origin = "54832", "elabtrack_v2_batch1_test", "http://localhost:18085/api/v1", "http://localhost:15175"
	}
	ctx := context.Background()
	cfg, err := config.Parse(map[string]string{
		"APP_ENV": "test", "JWT_SECRET": os.Getenv("JWT_SECRET"),
		"DB_HOST": "127.0.0.1", "DB_PORT": port, "DB_NAME": dbName,
		"DB_USER": "elabtrack_runtime", "DB_PASSWORD": os.Getenv("DB_PASSWORD"), "DB_SSLMODE": "disable",
	})
	require(t, err == nil, "integration configuration")
	db, err := database.New(ctx, cfg.DB)
	require(t, err == nil, "isolated database connection")
	t.Cleanup(db.Close)
	var identity string
	err = db.Pool.QueryRow(ctx, `SELECT current_database()`).Scan(&identity)
	require(t, err == nil && identity == dbName, "disposable database identity")
	var runtimeOnly bool
	err = db.Pool.QueryRow(ctx, `SELECT current_user='elabtrack_runtime' AND NOT rolsuper AND NOT rolcreaterole FROM pg_roles WHERE rolname=current_user`).Scan(&runtimeOnly)
	require(t, err == nil && runtimeOnly, "DML-only runtime identity")
	users := postgres.NewUserRepository(db.Pool)
	tokens := postgres.NewAuthRepository(db.Pool)
	hasher := security.NewBcryptHasher(0)
	password := uuid.NewString() + "-Synthetic-Only"
	hash, err := hasher.Hash(password)
	require(t, err == nil, "fixture hashing")
	fixture := func(role domainuser.Role, active bool) *domainuser.User {
		u := domainuser.NewUser("phase1g-"+uuid.NewString()+"@example.invalid", "Synthetic integration account", hash)
		u.Role, u.IsActive = role, active
		require(t, users.Create(ctx, u) == nil, "fixture creation")
		t.Cleanup(func() {
			_, e := db.Pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, u.ID)
			require(t, e == nil, "fixture cleanup")
		})
		return u
	}
	u := fixture(domainuser.RoleBorrower, true)
	staff := fixture(domainuser.RoleStaff, true)
	admin := fixture(domainuser.RoleAdmin, true)
	inactive := fixture(domainuser.RoleBorrower, false)
	client := &http.Client{Timeout: 10 * time.Second}
	var secretMu sync.Mutex
	secrets := []string{password, hash, os.Getenv("JWT_SECRET"), os.Getenv("DB_PASSWORD")}
	remember := func(s string) { secretMu.Lock(); secrets = append(secrets, s); secretMu.Unlock() }
	request := func(method, path string, body any, access, raw string) (*http.Response, map[string]any) {
		var encoded []byte
		if body != nil {
			encoded, _ = json.Marshal(body)
		}
		req, e := http.NewRequest(method, baseURL+path, bytes.NewReader(encoded))
		require(t, e == nil, "request creation")
		req.Header.Set("Origin", origin)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", "ignored-upstream-id")
		if access != "" {
			req.Header.Set("Authorization", "Bearer "+access)
		}
		if raw != "" {
			req.AddCookie(&http.Cookie{Name: "elabtrack_v2_refresh", Value: raw})
		}
		resp, e := client.Do(req)
		require(t, e == nil, "real HTTP request")
		defer resp.Body.Close()
		b, e := io.ReadAll(resp.Body)
		require(t, e == nil, "read response")
		var result map[string]any
		if len(b) > 0 {
			require(t, json.Unmarshal(b, &result) == nil, "JSON envelope")
		}
		id := resp.Header.Get("X-Request-ID")
		_, e = uuid.Parse(id)
		require(t, e == nil && id != "ignored-upstream-id", "server UUID")
		if resp.StatusCode >= 400 {
			er, ok := result["error"].(map[string]any)
			require(t, ok && er["requestId"] == id, "error/header correlation")
		}
		return resp, result
	}
	login := func(account *domainuser.User) (string, string) {
		r, b := request("POST", "/auth/login", map[string]string{"email": account.Email, "password": password}, "", "")
		require(t, r.StatusCode == 200, "login")
		data := b["data"].(map[string]any)
		access := data["access_token"].(string)
		_, present := data["refresh_token"]
		require(t, !present, "no refresh JSON")
		var raw string
		for _, c := range r.Cookies() {
			if c.Name == "elabtrack_v2_refresh" {
				raw = c.Value
				require(t, c.HttpOnly && c.SameSite == http.SameSiteLaxMode && c.Path == "/api/v1/auth" && !c.Secure, "local cookie policy")
			}
		}
		require(t, len(raw) == 64, "refresh cookie")
		remember(access)
		remember(raw)
		remember(digest(raw))
		return access, raw
	}
	access, raw := login(u)
	t.Run("health_and_safe_current_account", func(t *testing.T) {
		for _, path := range []string{"/health", "/ready", "/auth/me"} {
			r, b := request("GET", path, nil, access, "")
			require(t, r.StatusCode == 200 && b["success"] == true, "health/current account")
			if path == "/auth/me" {
				data := b["data"].(map[string]any)
				require(t, len(data) == 5 && data["id"] == u.ID.String() && data["role"] == "BORROWER", "safe current DB account")
			}
		}
	})
	t.Run("generic_invalid_login", func(t *testing.T) {
		var prior string
		for _, entry := range []struct{ email, password string }{{u.Email, "wrong-synthetic-password"}, {"absent@example.invalid", password}, {inactive.Email, password}} {
			r, b := request("POST", "/auth/login", map[string]string{"email": entry.email, "password": entry.password}, "", "")
			require(t, r.StatusCode == 401, "generic invalid login status")
			er := b["error"].(map[string]any)
			require(t, er["code"] == "INVALID_CREDENTIALS", "generic invalid login code")
			msg := er["message"].(string)
			if prior != "" {
				require(t, msg == prior, "no account enumeration detail")
			}
			prior = msg
		}
	})
	t.Run("product_permission_matrix", func(t *testing.T) {
		for _, account := range []*domainuser.User{fixture(domainuser.RoleBorrower, true), staff, admin} {
			a, _ := login(account)
			r, b := request("GET", "/auth/me", nil, a, "")
			require(t, r.StatusCode == 200 && b["data"].(map[string]any)["role"] == string(account.Role), "current product role")
			r, _ = request("GET", "/users/", nil, a, "")
			expected := 403
			if account.Role == domainuser.RoleAdmin {
				expected = 200
			}
			require(t, r.StatusCode == expected, "server role permission")
		}
		r, _ := request("POST", "/auth/register", map[string]string{"email": "unused@example.invalid", "password": password, "role": "ADMIN"}, "", "")
		require(t, r.StatusCode == 404, "no public registration")
	})
	t.Run("stale_authority", func(t *testing.T) {
		a, adminCookie := login(admin)
		r, _ := request("GET", "/users/", nil, a, "")
		require(t, r.StatusCode == 200, "temporary admin authorized")
		_, e := db.Pool.Exec(ctx, `UPDATE users SET role='BORROWER' WHERE id=$1`, admin.ID)
		require(t, e == nil, "demote synthetic account")
		r, _ = request("GET", "/users/", nil, a, "")
		require(t, r.StatusCode == 403, "old admin token loses privilege")
		_, e = db.Pool.Exec(ctx, `UPDATE users SET is_active=false WHERE id=$1`, admin.ID)
		require(t, e == nil, "disable synthetic account")
		r, _ = request("POST", "/auth/refresh", nil, "", adminCookie)
		require(t, r.StatusCode == 401, "disabled account cannot restore session")
		for _, path := range []string{"/users/", "/auth/me"} {
			r, _ = request("GET", path, nil, a, "")
			require(t, r.StatusCode == 403, "old token disabled")
		}
	})
	t.Run("hash_rotation_concurrent_http_logout", func(t *testing.T) {
		var n int
		e := db.Pool.QueryRow(ctx, `SELECT count(*) FROM refresh_tokens WHERE user_id=$1 AND token_hash=$2`, u.ID, digest(raw)).Scan(&n)
		require(t, e == nil && n == 1, "hash persistence")
		e = db.Pool.QueryRow(ctx, `SELECT count(*) FROM refresh_tokens WHERE token_hash=$1`, raw).Scan(&n)
		require(t, e == nil && n == 0, "raw absent")
		r, _ := request("POST", "/auth/refresh", nil, "", digest(raw))
		require(t, r.StatusCode == 401, "stored digest is not bearer")
		// Separate real HTTP connections race the exact same cookie.
		const competitors = 12
		start := make(chan struct{})
		var wg sync.WaitGroup
		var mu sync.Mutex
		wins, denials := 0, 0
		successor := ""
		for i := 0; i < competitors; i++ {
			wg.Go(func() {
				<-start
				r, b := request("POST", "/auth/refresh", nil, "", raw)
				mu.Lock()
				defer mu.Unlock()
				switch r.StatusCode {
				case 200:
					wins++
					require(t, b["data"].(map[string]any)["access_token"] != "", "renewed access")
					for _, c := range r.Cookies() {
						if c.Name == "elabtrack_v2_refresh" {
							successor = c.Value
						}
					}
				case 401:
					denials++
					require(t, len(r.Cookies()) == 0, "losing refresh cannot erase the winner cookie")
				default:
					t.Errorf("unexpected refresh status %d", r.StatusCode)
				}
			})
		}
		close(start)
		wg.Wait()
		require(t, wins == 1 && denials == competitors-1 && successor != "" && successor != raw, "single committed concurrent winner")
		t.Logf("real concurrent refresh: %d success / %d safe denials; no denial Set-Cookie", wins, denials)
		remember(successor)
		remember(digest(successor))
		e = db.Pool.QueryRow(ctx, `SELECT count(*) FROM refresh_tokens WHERE user_id=$1 AND revoked_at IS NULL AND expires_at>now()`, u.ID).Scan(&n)
		require(t, e == nil && n == 1, "exactly one usable replacement")
		e = db.Pool.QueryRow(ctx, `SELECT count(*) FROM refresh_tokens old JOIN refresh_tokens successor ON successor.id=old.replaced_by WHERE old.token_hash=$1 AND old.revoked_at IS NOT NULL AND successor.token_hash=$2`, digest(raw), digest(successor)).Scan(&n)
		require(t, e == nil && n == 1, "coherent consumption/replacement")
		r, _ = request("POST", "/auth/refresh", nil, "", raw)
		require(t, r.StatusCode == 401 && len(r.Cookies()) == 0, "old replay denied without cookie mutation")
		r, _ = request("POST", "/auth/refresh", nil, "", successor)
		require(t, r.StatusCode == 200, "winner remains usable after replay")
		for _, c := range r.Cookies() {
			if c.Name == "elabtrack_v2_refresh" {
				successor = c.Value
			}
		}
		remember(successor)
		remember(digest(successor))
		for i := 0; i < 2; i++ {
			r, _ = request("POST", "/auth/logout", nil, "", successor)
			require(t, r.StatusCode == 204, "idempotent logout")
			require(t, len(r.Cookies()) == 1 && r.Cookies()[0].Value == "", "cookie cleared")
		}
		r, _ = request("POST", "/auth/refresh", nil, "", successor)
		require(t, r.StatusCode == 401 && len(r.Cookies()) == 0, "logout revocation persisted without refresh cookie mutation")
	})
	t.Run("real_transaction_rollback", func(t *testing.T) {
		txm := database.NewTxManager(db.Pool)
		issuer := security.NewJWTIssuer(cfg.JWT)
		svc := appauth.NewService(users, tokens, hasher, issuer, security.SHA256RefreshHasher{}, txm, users, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
		pair, e := svc.Login(ctx, appauth.LoginRequest{Email: u.Email, Password: password})
		require(t, e == nil, "rollback fixture login")
		failing := appauth.NewService(users, failAfterConsume{tokens}, hasher, issuer, security.SHA256RefreshHasher{}, txm, users, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
		failed, e := failing.Refresh(ctx, appauth.RefreshRequest{RefreshToken: pair.RefreshToken})
		require(t, errors.Is(e, shared.ErrInternal) && failed.AccessToken == "" && failed.RefreshToken == "", "rollback returns no credentials")
		var n int
		e = db.Pool.QueryRow(ctx, `SELECT count(*) FROM refresh_tokens WHERE user_id=$1 AND revoked_at IS NULL`, u.ID).Scan(&n)
		require(t, e == nil && n == 1, "rollback removes successor and preserves old")
		_, e = svc.Refresh(ctx, appauth.RefreshRequest{RefreshToken: pair.RefreshToken})
		require(t, e == nil, "old session usable after rollback")
	})
	t.Run("missing_account_after_issuance", func(t *testing.T) {
		missing := fixture(domainuser.RoleBorrower, true)
		a, c := login(missing)
		_, e := db.Pool.Exec(ctx, `DELETE FROM users WHERE id=$1`, missing.ID)
		require(t, e == nil, "delete only owned fixture")
		r, _ := request("GET", "/auth/me", nil, a, "")
		require(t, r.StatusCode == 401, "missing identity denied")
		r, _ = request("POST", "/auth/refresh", nil, "", c)
		require(t, r.StatusCode == 401, "missing session denied")
	})
	t.Run("bounded_retention", func(t *testing.T) {
		old := fixture(domainuser.RoleBorrower, true)
		for i := 0; i < 3; i++ {
			_, e := db.Pool.Exec(ctx, `INSERT INTO refresh_tokens(id,token_hash,user_id,created_at,updated_at,expires_at) VALUES($1,$2,$3,now()-interval '10 days',now()-interval '10 days',now()-interval '9 days')`, uuid.New(), digest(uuid.NewString()), old.ID)
			require(t, e == nil, "expired fixture")
		}
		cutoff := time.Now().Add(-7 * 24 * time.Hour)
		n, e := tokens.Cleanup(ctx, cutoff, 2)
		require(t, e == nil && n == 2, "bounded batch")
		n, e = tokens.Cleanup(ctx, cutoff, 1000)
		require(t, e == nil && n == 1, "remaining terminal row")
		var active int
		e = db.Pool.QueryRow(ctx, `SELECT count(*) FROM refresh_tokens WHERE user_id=$1 AND revoked_at IS NULL AND expires_at>now()`, u.ID).Scan(&active)
		require(t, e == nil && active == 1, "active preserved")
	})
	if path := os.Getenv("INTEGRATION_SECRET_FILE"); path != "" {
		if previous, e := os.ReadFile(path); e == nil {
			var prior []string
			require(t, json.Unmarshal(previous, &prior) == nil, "private sentinel format")
			secrets = append(prior, secrets...)
		}
		b, _ := json.Marshal(secrets)
		require(t, os.WriteFile(path, b, 0600) == nil, "private redaction sentinels")
	}
}

type failAfterConsume struct{ domainauth.Repository }

func (r failAfterConsume) Consume(ctx context.Context, h string, id uuid.UUID, now time.Time) error {
	if e := r.Repository.Consume(ctx, h, id, now); e != nil {
		return e
	}
	return shared.ErrInternal
}
func digest(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func require(t *testing.T, ok bool, label string) {
	t.Helper()
	if !ok {
		t.Fatal(label)
	}
}
