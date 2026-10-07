package bootstrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	domainshared "github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func httpSecurityConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg, err := config.Parse(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	return cfg
}
func securityRequest(t *testing.T, app *fiber.App, method, path, contentType, body string, origins ...string) (int, string, map[string]string) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	origin := "http://localhost:5173"
	if len(origins) > 0 {
		origin = origins[0]
	}
	req.Header.Set("Origin", origin)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Fatal("HTTP request failed", err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	h := map[string]string{}
	for key := range res.Header {
		h[key] = res.Header.Get(key)
	}
	h["X-Request-ID"] = res.Header.Get("X-Request-ID")
	return res.StatusCode, string(raw), h
}

func TestHTTPBodyAndParserSafety(t *testing.T) {
	cfg := httpSecurityConfig(t)
	h := handlers.NewAuthHandler(nil, validator.New(), config.Development, cfg.Security)
	for _, tt := range []struct {
		name, path, media, body string
		status                  int
	}{
		{"auth oversized", "/api/v1/auth/login", "application/json", strings.Repeat(" ", 16<<10) + "{}", 413},
		{"auth alias oversized", "/API/V1/AUTH/LOGIN/", "application/json", strings.Repeat(" ", 16<<10) + "{}", 413},
		{"auth exact ceiling", "/api/v1/auth/login", "application/json", strings.Repeat(" ", (16<<10)-2) + "{}", 400},
		{"text rejected", "/api/v1/auth/login", "text/plain", "{}", 415},
		{"form rejected", "/api/v1/auth/login", "application/x-www-form-urlencoded", "email=a&password=b", 415},
		{"missing media rejected", "/api/v1/auth/login", "", "{}", 415},
		{"malformed JSON", "/api/v1/auth/login", "application/json", "{", 400},
		{"multiple documents", "/api/v1/auth/login", "application/json", "{} {}", 400},
		{"unknown login field", "/api/v1/auth/login", "application/json", `{"role":"admin"}`, 400},
		{"JSON charset supported", "/api/v1/auth/login", "application/json; charset=utf-8", "{}", 400},
		{"bodyless refresh", "/api/v1/auth/refresh", "", "", 204},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Each case gets its own server so parser tests do not consume the login budget.
			fresh := newServer(cfg, nil, nil)
			fresh.Post("/api/v1/auth/login", h.Login)
			fresh.Post("/api/v1/auth/refresh", func(c fiber.Ctx) error { return c.SendStatus(204) })
			status, body, headers := securityRequest(t, fresh, "POST", tt.path, tt.media, tt.body)
			if status != tt.status || strings.Contains(body, "runtime") || strings.Contains(body, "stack") || headers["X-Content-Type-Options"] != "nosniff" {
				t.Fatalf("unsafe parser response: status %d want %d", status, tt.status)
			}
		})
	}
}

func TestGlobalBodyParserCeiling(t *testing.T) {
	app := newServer(httpSecurityConfig(t), nil, nil)
	handlerCalled := false
	app.Post("/upload-not-implemented", func(c fiber.Ctx) error { handlerCalled = true; return c.SendStatus(204) })
	_, err := app.Test(httptest.NewRequest("POST", "/upload-not-implemented", strings.NewReader(strings.Repeat("x", (1<<20)+1))))
	if !errors.Is(err, fasthttp.ErrBodyTooLarge) || handlerCalled {
		t.Fatal("global parser ceiling did not reject before handler")
	}
}

func TestRoutingAliasesCannotBypassLoginBudget(t *testing.T) {
	cfg := httpSecurityConfig(t)
	cfg.Security.LoginRateLimitMax = 1
	app := newServer(cfg, nil, nil)
	app.Post("/api/v1/auth/login", func(c fiber.Ctx) error { return c.SendStatus(204) })
	status, _, _ := securityRequest(t, app, "POST", "/api/v1/auth/login", "application/json", "{}")
	if status != 204 {
		t.Fatal("ordinary login route denied")
	}
	for _, path := range []string{"/API/V1/AUTH/LOGIN/", "/api/v1/auth/login//", "/api/v1/auth/%6cogin", "/api/v1/auth//login"} {
		status, _, _ = securityRequest(t, app, "POST", path, "application/json", "{}")
		if status != 429 && status != 404 {
			t.Fatal("route alias bypassed sensitive budget")
		}
	}
}

func TestPanicAndErrorRedaction(t *testing.T) {
	core, logs := observer.New(zap.InfoLevel)
	log := &logger.Logger{Logger: zap.New(core)}
	cfg, cfgErr := config.Parse(map[string]string{"APP_ENV": "production", "JWT_SECRET": "t9K2m7Q4v8R1z6N3p5W0x2C7b9H4f8L1", "DB_HOST": "db.example.invalid", "DB_NAME": "isolated_test", "DB_USER": "synthetic", "DB_PASSWORD": "synthetic-db-value", "DB_SSLMODE": "verify-full", "FRONTEND_URL": "https://app.example.invalid", "ALLOWED_ORIGINS": "https://app.example.invalid"})
	if cfgErr != nil {
		t.Fatal(cfgErr)
	}
	app := newServer(cfg, nil, log)
	const sensitive = "private-sentinel-password-cookie-jwt SQL SELECT /private/runtime.go"
	app.Get("/panic", func(c fiber.Ctx) error { panic(sensitive) })
	app.Get("/raw-error", func(c fiber.Ctx) error { return errors.New(sensitive) })
	app.Get("/framework-error", func(c fiber.Ctx) error { return fiber.NewError(400, sensitive) })
	app.Get("/wrapped-validation", func(c fiber.Ctx) error {
		return response.Error(c, fmt.Errorf("%s: %w", sensitive, domainshared.ErrInvalidInput))
	})
	for path, want := range map[string]int{"/panic": 500, "/raw-error": 500, "/framework-error": 400, "/wrapped-validation": 400, "/missing": 404} {
		status, body, h := securityRequest(t, app, "GET", path, "", "", "https://app.example.invalid")
		if status != want || strings.Contains(body, "private") || strings.Contains(body, "SQL") || strings.Contains(body, "goroutine") || h["X-Content-Type-Options"] != "nosniff" || h["Server"] != "" || h["Strict-Transport-Security"] != "" {
			t.Fatal("internal error/transport leakage")
		}
	}
	if logs.FilterMessage("panic recovered").Len() != 1 || logs.FilterMessage("request error").Len() != 1 {
		t.Fatal("safe diagnostics absent")
	}
	for _, entry := range logs.All() {
		fields := entry.ContextMap()
		if fields["error"] != nil || strings.Contains(fmt.Sprint(fields), sensitive) {
			t.Fatal("raw secret-bearing error logged")
		}
		if entry.Message == "panic recovered" && (fields["stack"] == nil || fields["panic_type"] != "string") {
			t.Fatal("safe panic source diagnostic absent")
		}
	}
}

func TestHealthAndServerBoundaries(t *testing.T) {
	cfg := httpSecurityConfig(t)
	app := newServer(cfg, nil, nil)
	app.Get("/api/v1/health", handlers.NewHealthHandler(nil).Health)
	status, body, headers := securityRequest(t, app, "GET", "/api/v1/health", "", "")
	var envelope struct {
		Success bool `json:"success"`
		Data    struct {
			Status  string `json:"status"`
			Service string `json:"service"`
		} `json:"data"`
	}
	var fields map[string]any
	if status != 200 || json.Unmarshal([]byte(body), &envelope) != nil || !envelope.Success || envelope.Data.Status != "ok" || envelope.Data.Service != "elabtrack-v2" || json.Unmarshal([]byte(body), &fields) != nil || len(fields) != 4 || headers["Cache-Control"] != "no-store" {
		t.Fatal("health envelope/liveness exposure incorrect")
	}

	actual := app.Config()
	if actual.ReadTimeout != 10*time.Second || actual.WriteTimeout != 15*time.Second || actual.IdleTimeout != time.Minute || actual.ReadBufferSize != 8192 || actual.BodyLimit != 1<<20 || actual.TrustProxy || actual.ProxyHeader != "" || actual.ServerHeader != "" {
		t.Fatal("server safety settings not wired")
	}
}
