package bootstrap

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPhase12ProfileHTTP(t *testing.T) {
	infra := batchInfra(t)
	c := buildContainer(infra)
	users, _, _ := phase7Fixtures(t, infra, c)
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	_, issuer := newSecurityAdapters(infra.Config)
	tokens := map[string]string{}
	for role, u := range users {
		v, e := issuer.IssueAccessToken(t.Context(), u)
		if e != nil {
			t.Fatal(e)
		}
		tokens[role] = v
	}
	var pngBytes bytes.Buffer
	_ = png.Encode(&pngBytes, image.NewRGBA(image.Rect(0, 0, 8, 8)))
	id := users["borrower"].ID.String()
	request := func(role, method, path, origin, content, key string, raw []byte) (int, []byte, string) {
		t.Helper()
		r := httptest.NewRequest(method, "/api/v1/profile-images/"+path, bytes.NewReader(raw))
		if role != "" {
			r.Header.Set("Authorization", "Bearer "+tokens[role])
		}
		r.Header.Set("Origin", origin)
		r.Header.Set("Content-Type", content)
		r.Header.Set("Idempotency-Key", key)
		res, e := server.Test(r, fiber.TestConfig{Timeout: 15 * time.Second})
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		data, _ := io.ReadAll(res.Body)
		return res.StatusCode, data, res.Header.Get("Content-Type")
	}
	upload := func(version string, raw []byte, extra bool) ([]byte, string) {
		var data bytes.Buffer
		w := multipart.NewWriter(&data)
		_ = w.WriteField("expected_version", version)
		if extra {
			_ = w.WriteField("role", "ADMIN")
		}
		f, _ := w.CreateFormFile("image", "../../private.png")
		_, _ = f.Write(raw)
		_ = w.Close()
		return data.Bytes(), w.FormDataContentType()
	}

	for _, role := range []string{"borrower", "staff", "admin"} {
		code, data, _ := request(role, "GET", "account/me", "", "", "", nil)
		if code != 200 || !json.Valid(data) || strings.Contains(string(data), "password") || strings.Contains(string(data), "token") {
			t.Fatal("safe own metadata", role, code)
		}
	}
	ownCode, _, _ := request("", "GET", "account/me", "", "", "", nil)
	if ownCode != 401 {
		t.Fatal("own metadata requires authentication")
	}
	code, raw, _ := request("borrower", "GET", id, "", "", "", nil)
	if code != 200 || !strings.Contains(string(raw), `"image_id":null`) {
		t.Fatal("own fallback metadata")
	}
	payload, content := upload("0", pngBytes.Bytes(), false)
	key := uuid.NewString()
	code, raw, _ = request("borrower", "POST", id, infra.Config.Security.AllowedOrigins[0], content, key, payload)
	if code != 200 {
		t.Fatalf("own upload: %d", code)
	}
	var envelope struct {
		Data struct {
			ImageID uuid.UUID `json:"image_id"`
			Version int       `json:"version"`
		}
	}
	if json.Unmarshal(raw, &envelope) != nil || envelope.Data.Version != 1 {
		t.Fatal("saved metadata")
	}
	imageID := envelope.Data.ImageID.String()
	for _, tc := range []struct {
		role, method, path, origin string
		want                       int
	}{{"", "GET", id, "", 401}, {"borrower", "GET", users["staff"].ID.String(), "", 403}, {"staff", "GET", users["admin"].ID.String(), "", 403}, {"staff", "GET", id + "/" + imageID, "", 200}, {"admin", "GET", id + "/" + imageID, "", 200}, {"borrower", "GET", id + "/" + uuid.NewString(), "", 404}, {"borrower", "POST", id, "https://untrusted.invalid", 403}, {"admin", "POST", id, infra.Config.Security.AllowedOrigins[0], 403}} {
		code, _, mime := request(tc.role, tc.method, tc.path, tc.origin, content, key, payload)
		if code != tc.want {
			t.Fatalf("%s %s %d expected %d", tc.role, tc.method, code, tc.want)
		}
		if code == 200 && strings.Contains(tc.path, "/") && mime != "image/png" {
			t.Fatal("canonical binary MIME")
		}
	}
	bad, ctype := upload("1", []byte("<svg/>"), false)
	code, _, _ = request("borrower", "POST", id, infra.Config.Security.AllowedOrigins[0], ctype, uuid.NewString(), bad)
	if code != 400 {
		t.Fatal("invalid format")
	}
	bad, ctype = upload("1", pngBytes.Bytes(), true)
	code, _, _ = request("borrower", "POST", id, infra.Config.Security.AllowedOrigins[0], ctype, uuid.NewString(), bad)
	if code != 400 {
		t.Fatal("strict multipart fields")
	}
	code, _, _ = request("borrower", "POST", id, infra.Config.Security.AllowedOrigins[0], content, uuid.NewString(), payload)
	if code != 409 {
		t.Fatal("stale version")
	}
	key = uuid.NewString()
	body := []byte(`{"expected_version":1,"confirm":true}`)
	code, _, _ = request("borrower", "PATCH", id, infra.Config.Security.AllowedOrigins[0], "application/json", key, body)
	if code != 200 {
		t.Fatalf("remove %d", code)
	}
	code, _, _ = request("borrower", "PATCH", id, infra.Config.Security.AllowedOrigins[0], "application/json", key, body)
	if code != 200 {
		t.Fatal("remove replay")
	}
	code, _, _ = request("borrower", "GET", id+"/"+imageID, "", "", "", nil)
	if code != 404 {
		t.Fatal("old profile bytes unavailable")
	}
	code, _, _ = request("borrower", "PATCH", id, infra.Config.Security.AllowedOrigins[0], "application/json", uuid.NewString(), []byte(`{"expected_version":2,"confirm":true,"role":"ADMIN"}`))
	if code != 400 {
		t.Fatal("strict remove body")
	}
}
