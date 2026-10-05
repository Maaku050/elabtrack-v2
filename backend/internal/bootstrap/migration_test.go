package bootstrap

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
)

// Exercise the real middleware chain and migrated bind/auth paths without a database.
func TestFiberV3HTTPContracts(t *testing.T) {
	cfg := &config.Config{App: config.AppConfig{BodyLimit: 1 << 20}, Security: config.SecurityConfig{AllowedOrigins: []string{"http://localhost:5173"}, RateLimitMax: 100, RateLimitWindow: time.Minute}}
	app := newServer(cfg, nil, nil)
	health := handlers.NewHealthHandler(nil)
	auth := handlers.NewAuthHandler(nil, validator.New())
	app.Get("/api/v1/health", health.Health)
	app.Post("/api/v1/auth/login", auth.Login)
	app.Get("/api/v1/users/me", middleware.Auth(nil, nil), func(c fiber.Ctx) error { return c.SendStatus(200) })
	app.Get("/missing-principal", auth.Me)
	app.Get("/panic", func(c fiber.Ctx) error { panic("test") })
	cases := []struct {
		name, method, path, body string
		status                   int
	}{
		{"health", "GET", "/api/v1/health", "", 200},
		{"malformed JSON", "POST", "/api/v1/auth/login", "{", 400},
		{"invalid credentials input", "POST", "/api/v1/auth/login", "{}", 422},
		{"missing bearer", "GET", "/api/v1/users/me", "", 401},
		{"missing principal helper", "GET", "/missing-principal", "", 401},
		{"panic recovery", "GET", "/panic", "", 500},
		{"unknown route", "GET", "/unknown", "", 404},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Origin", "http://localhost:5173")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tt.status {
				t.Fatalf("status %d, want %d", resp.StatusCode, tt.status)
			}
			if resp.Header.Get("X-Request-ID") == "" {
				t.Error("missing generated request ID")
			}
			if resp.Header.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
				t.Error("allowlisted origin missing")
			}
			var body map[string]any
			if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if _, ok := body["success"]; !ok && tt.name != "health" {
				t.Error("missing response envelope")
			}
		})
	}
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	req.Header.Set("Origin", "https://untrusted.example")
	req.Header.Set("Access-Control-Request-Method", "POST")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.Header.Get("Access-Control-Allow-Origin") != "" {
		t.Error("untrusted origin accepted")
	}
}
