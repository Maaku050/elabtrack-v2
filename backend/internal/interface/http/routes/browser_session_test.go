package routes_test

import (
	"context"
	"encoding/json"
	appauth "github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const cookieName = "elabtrack_v2_refresh"
const trustedOrigin = "https://app.example.invalid"

func cookieRequest(t *testing.T, app *fiber.App, path, raw, origin, body string) (int, string, []*http.Cookie) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	if raw != "" {
		req.AddCookie(&http.Cookie{Name: cookieName, Value: raw})
	}
	res, err := app.Test(req)
	if err != nil {
		t.Fatal("HTTP request failed")
	}
	defer res.Body.Close()
	if res.Header.Get("Cache-Control") != "no-store" {
		t.Fatal("session response must not be cached")
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("response read failed")
	}
	return res.StatusCode, string(data), res.Cookies()
}
func assertCookie(t *testing.T, cookies []*http.Cookie, secure, cleared bool) *http.Cookie {
	t.Helper()
	if len(cookies) != 1 {
		t.Fatal("expected exactly one refresh cookie")
	}
	c := cookies[0]
	if c.Name != cookieName || c.Path != "/api/v1/auth" || c.Domain != "" || !c.HttpOnly || c.Secure != secure || c.SameSite != http.SameSiteLaxMode {
		t.Fatal("refresh cookie attribute policy incorrect")
	}
	if cleared {
		if c.Value != "" || c.MaxAge != -1 || !c.Expires.Before(time.Now()) {
			t.Fatal("cookie not consistently cleared")
		}
	} else {
		if c.Value == "" || !c.Expires.After(time.Now()) {
			t.Fatal("cookie credential/expiry missing")
		}
	}
	return c
}
func assertSafeSessionBody(t *testing.T, body string, raw, hash string) {
	t.Helper()
	var result struct{ Data map[string]json.RawMessage }
	if json.Unmarshal([]byte(body), &result) != nil {
		t.Fatal("invalid JSON")
	}
	for _, key := range []string{"refresh_token", "token_hash", "refresh_expires_at", "password", "password_hash", "replaced_by"} {
		if _, ok := result.Data[key]; ok {
			t.Fatal("unsafe response field")
		}
	}
	if len(result.Data) != 4 || result.Data["access_token"] == nil || result.Data["user"] == nil {
		t.Fatal("session response allowlist incorrect")
	}
	if strings.Contains(body, raw) || strings.Contains(body, hash) {
		t.Fatal("refresh material exposed in JSON")
	}
	var user map[string]any
	if json.Unmarshal(result.Data["user"], &user) != nil || len(user) != 5 || user["id"] != selfID.String() || user["role"] != "user" || user["is_active"] != true {
		t.Fatal("safe current account projection incorrect")
	}
}
func TestLoginCookiePolicyAndSafeResponse(t *testing.T) {
	for _, env := range []config.Environment{config.Development, config.Test, config.Production} {
		t.Run(string(env), func(t *testing.T) {
			r := &users{account: fixtureUser()}
			hasher := security.NewBcryptHasher(4)
			r.account.Password, _ = hasher.Hash("synthetic-password")
			persisted := &tokens{}
			app := newApp(env, r, &issuer{}, persisted)
			status, body, cookies := cookieRequest(t, app, "/api/v1/auth/login", "", trustedOrigin, `{"email":"current@example.invalid","password":"synthetic-password"}`)
			if status != http.StatusOK || len(persisted.created) != 1 {
				t.Fatal("login did not establish session")
			}
			cookie := assertCookie(t, cookies, env == config.Production, false)
			if !cookie.Expires.Equal(persisted.created[0].ExpiresAt.Truncate(time.Second)) {
				t.Fatal("cookie not aligned with persisted refresh expiry")
			}
			assertSafeSessionBody(t, body, cookie.Value, persisted.created[0].TokenHash)
			if strings.Contains(body, r.account.Password) {
				t.Fatal("password hash exposed")
			}
		})
	}
}
func TestRefreshCookieRotationAndLegacyBodyRemoval(t *testing.T) {
	raw := strings.Repeat("b", 64)
	hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
	now := time.Now().UTC()
	sessions := &httpSessions{session: domainauth.NewRefreshToken(hash, selfID, now, now.Add(time.Hour))}
	r := &users{account: fixtureUser()}
	app := refreshHTTPApp(r, sessions)
	status, body, cookies := cookieRequest(t, app, "/api/v1/auth/refresh", raw, trustedOrigin, "")
	if status != 200 || len(sessions.created) != 1 || sessions.session.RevokedAt == nil {
		t.Fatal("refresh did not rotate")
	}
	cookie := assertCookie(t, cookies, false, false)
	if cookie.Value == raw || !cookie.Expires.Equal(sessions.created[0].ExpiresAt.Truncate(time.Second)) {
		t.Fatal("rotation cookie lifetime/value incorrect")
	}
	assertSafeSessionBody(t, body, cookie.Value, sessions.created[0].TokenHash)
	status, _, cookies = cookieRequest(t, app, "/api/v1/auth/refresh", raw, trustedOrigin, "")
	if status != 401 || len(sessions.created) != 1 {
		t.Fatal("old cookie replay succeeded")
	}
	assertCookie(t, cookies, false, true)
	sessions.session = sessions.created[0]
	status, _, cookies = cookieRequest(t, app, "/api/v1/auth/refresh", "", trustedOrigin, `{"refresh_token":"`+cookie.Value+`"}`)
	if status != 401 || len(sessions.created) != 1 {
		t.Fatal("body-only credential accepted")
	}
	assertCookie(t, cookies, false, true)
	status, _, cookies = cookieRequest(t, app, "/api/v1/auth/logout", "", trustedOrigin, `{"refresh_token":"`+cookie.Value+`"}`)
	if status != 204 || sessions.session.RevokedAt != nil {
		t.Fatal("logout accepted body credential")
	}
	assertCookie(t, cookies, false, true)
}
func TestCookieClearingAndTransientFailure(t *testing.T) {
	for _, env := range []config.Environment{config.Development, config.Production} {
		t.Run(string(env), func(t *testing.T) {
			sessions := &httpSessions{}
			r := &users{account: fixtureUser()}
			svc := appauth.NewService(r, sessions, security.NewBcryptHasher(4), &issuer{}, security.SHA256RefreshHasher{}, httpTransaction{}, httpAccounts{r}, time.Minute, time.Hour)
			h := handlers.NewAuthHandler(svc, validator.New(), env, config.SecurityConfig{AllowedOrigins: []string{trustedOrigin}})
			app := fiber.New()
			app.Post("/api/v1/auth/refresh", h.Refresh)
			app.Post("/api/v1/auth/logout", h.Logout)
			for _, raw := range []string{"", strings.Repeat("a", 64), "invalid-sentinel"} {
				status, _, cookies := cookieRequest(t, app, "/api/v1/auth/refresh", raw, trustedOrigin, "")
				if status != 401 {
					t.Fatal("invalid refresh status")
				}
				assertCookie(t, cookies, env == config.Production, true)
				status, body, cookies := cookieRequest(t, app, "/api/v1/auth/logout", raw, trustedOrigin, "")
				if status != 204 || body != "" {
					t.Fatal("logout not idempotent")
				}
				assertCookie(t, cookies, env == config.Production, true)
			}
			sessions.lookupErr = context.DeadlineExceeded
			status, _, cookies := cookieRequest(t, app, "/api/v1/auth/refresh", strings.Repeat("a", 64), trustedOrigin, "")
			if status != 500 || len(cookies) != 0 {
				t.Fatal("transient storage failure should preserve cookie for recovery")
			}
			sessions.revokeErr = context.DeadlineExceeded
			status, _, cookies = cookieRequest(t, app, "/api/v1/auth/logout", strings.Repeat("a", 64), trustedOrigin, "")
			if status != 500 {
				t.Fatal("logout storage failure must fail safely")
			}
			assertCookie(t, cookies, env == config.Production, true)
		})
	}
}
func TestSessionCSRFOriginEnforcement(t *testing.T) {
	for _, path := range []string{"/api/v1/auth/login", "/api/v1/auth/register", "/api/v1/auth/refresh", "/api/v1/auth/logout"} {
		t.Run(path, func(t *testing.T) {
			for _, origin := range []string{"", "null", "https://attacker.example.invalid", trustedOrigin + ".attacker.invalid", trustedOrigin + "/", "http://app.example.invalid", "https://app.example.invalid:444"} {
				r := &users{account: fixtureUser()}
				persisted := &tokens{}
				app := newApp(config.Test, r, &issuer{}, persisted)
				status, body, cookies := cookieRequest(t, app, path, strings.Repeat("a", 64), origin, `{"email":"current@example.invalid","password":"synthetic-password"}`)
				if status != 403 || len(cookies) != 0 || r.creates != 0 || len(persisted.created) != 0 || strings.Contains(body, "attacker") {
					t.Fatal("untrusted origin reached session mutation")
				}
			}
		})
	}
}
func TestInvalidCookiePolicyFailsClosed(t *testing.T) {
	for _, test := range []struct {
		env     config.Environment
		origins []string
	}{
		{"unknown", []string{trustedOrigin}}, {"", []string{trustedOrigin}}, {config.Production, []string{"http://localhost:5173"}}, {config.Test, []string{"*"}}, {config.Development, nil},
	} {
		// Nil service proves invalid configuration fails before credential IO.
		h := handlers.NewAuthHandler(nil, validator.New(), test.env, config.SecurityConfig{AllowedOrigins: test.origins})
		app := fiber.New()
		app.Post("/api/v1/auth/refresh", h.Refresh)
		status, _, cookies := cookieRequest(t, app, "/api/v1/auth/refresh", strings.Repeat("a", 64), trustedOrigin, "")
		if status != 403 || len(cookies) != 0 {
			t.Fatal("invalid session cookie policy allowed IO")
		}
	}
}
func TestCredentialedCORSAndPostOnlySessionRoutes(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.CORS(config.SecurityConfig{AllowedOrigins: []string{trustedOrigin}}))
	app.Post("/api/v1/auth/refresh", func(c fiber.Ctx) error { return c.SendStatus(204) })
	for _, origin := range []string{trustedOrigin, "https://attacker.example.invalid"} {
		req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/refresh", nil)
		req.Header.Set("Origin", origin)
		req.Header.Set("Access-Control-Request-Method", "POST")
		res, err := app.Test(req)
		if err != nil {
			t.Fatal("preflight failed")
		}
		res.Body.Close()
		if origin == trustedOrigin {
			if res.Header.Get("Access-Control-Allow-Origin") != trustedOrigin || res.Header.Get("Access-Control-Allow-Credentials") != "true" {
				t.Fatal("credentialed preflight missing explicit origin")
			}
		} else if res.Header.Get("Access-Control-Allow-Origin") != "" {
			t.Fatal("arbitrary origin reflected")
		}
	}
	for _, path := range []string{"/api/v1/auth/refresh", "/api/v1/auth/logout"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		res, err := refreshHTTPApp(&users{account: fixtureUser()}, &httpSessions{}).Test(req)
		if err != nil {
			t.Fatal("GET check failed")
		}
		res.Body.Close()
		if res.StatusCode == 200 || len(res.Cookies()) != 0 {
			t.Fatal("state-changing GET auth route")
		}
	}
}

func TestCookieAloneDoesNotAuthenticateProtectedRoutes(t *testing.T) {
	r := &users{account: fixtureUser()}
	app := newApp(config.Development, r, &issuer{}, &tokens{})
	for _, path := range []string{"/api/v1/auth/me", "/api/v1/users/me", "/api/v1/users/"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.AddCookie(&http.Cookie{Name: cookieName, Value: strings.Repeat("a", 64)})
		res, err := app.Test(req)
		if err != nil {
			t.Fatal("protected request failed")
		}
		res.Body.Close()
		if res.StatusCode != 401 || r.lookups != 0 {
			t.Fatal("refresh cookie used as bearer authorization")
		}
	}
}
func TestDevelopmentLocalhostAndLocalRegistrationCookies(t *testing.T) {
	r := &users{account: fixtureUser()}
	r.account.Password, _ = security.NewBcryptHasher(4).Hash("synthetic-password")
	persisted := &tokens{}
	svc := appauth.NewService(r, persisted, security.NewBcryptHasher(4), &issuer{}, security.SHA256RefreshHasher{}, nil, nil, time.Minute, time.Hour)
	h := handlers.NewAuthHandler(svc, validator.New(), config.Development, config.SecurityConfig{AllowedOrigins: []string{"http://localhost:5173"}})
	app := fiber.New()
	app.Post("/api/v1/auth/login", h.Login)
	status, _, cookies := cookieRequest(t, app, "/api/v1/auth/login", "", "http://localhost:5173", `{"email":"current@example.invalid","password":"synthetic-password"}`)
	if status != 200 {
		t.Fatal("localhost development login failed")
	}
	assertCookie(t, cookies, false, false)
	r = &users{}
	persisted = &tokens{}
	app = newApp(config.Test, r, &issuer{}, persisted)
	status, body, cookies := cookieRequest(t, app, "/api/v1/auth/register", "", trustedOrigin, `{"email":"synthetic@example.invalid","name":"Local User","password":"synthetic-password"}`)
	if status != 201 || len(persisted.created) != 1 {
		t.Fatal("local registration failed")
	}
	cookie := assertCookie(t, cookies, false, false)
	if strings.Contains(body, cookie.Value) || strings.Contains(body, persisted.created[0].TokenHash) || strings.Contains(body, r.account.Password) {
		t.Fatal("local registration exposed refresh/password material")
	}
}
