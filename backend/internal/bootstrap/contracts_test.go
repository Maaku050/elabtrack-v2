package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/observability"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestEnvelopeCorrelationAndSafeLogs(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	log := &logger.Logger{Logger: zap.New(core)}
	cfg := httpSecurityConfig(t)
	cfg.Security.LoginRateLimitMax = 1
	cfg.Security.RateLimitMax = 100
	app := newServer(cfg, nil, log)
	auth := handlers.NewAuthHandler(nil, validator.New(), cfg.App.Env, cfg.Security)
	app.Post("/api/v1/auth/login", auth.Login)
	app.Get("/api/v1/forbidden", func(c fiber.Ctx) error { return response.Error(c, shared.ErrForbidden) })
	app.Get("/api/v1/direct-failure", func(c fiber.Ctx) error { return response.Error(c, errors.New("private-secret-SQL")) })
	app.Get("/api/v1/returned-failure", func(c fiber.Ctx) error { return errors.New("private-secret-SQL") })
	app.Get("/api/v1/panic", func(c fiber.Ctx) error { panic("private-secret-SQL") })
	app.Get("/api/v1/panic-framework", func(c fiber.Ctx) error { panic(fiber.ErrUnauthorized) })
	app.Get("/api/v1/validation", func(c fiber.Ctx) error { return response.Error(c, validator.FieldErrors{"email": "is required"}) })
	app.Get("/api/v1/denied", func(c fiber.Ctx) error { return response.Error(c, shared.ErrUnauthorized) })
	app.Get("/api/v1/correlation", func(c fiber.Ctx) error {
		id := observability.RequestID(c.Context())
		if id == "" {
			panic("missing application context")
		}
		return response.OK(c, "Success", map[string]string{"requestId": id})
	})
	cases := []struct {
		method, path, body string
		status             int
		code               string
	}{
		{"POST", "/api/v1/auth/login", "{", 400, "BAD_REQUEST"},
		{"POST", "/api/v1/auth/login", "{}", 429, "RATE_LIMITED"},
		{"GET", "/api/v1/validation", "", 400, "VALIDATION_ERROR"},
		{"GET", "/api/v1/denied", "", 401, "UNAUTHORIZED"},
		{"GET", "/api/v1/forbidden", "", 403, "FORBIDDEN"},
		{"GET", "/api/v1/missing/private-secret-SQL?token=private-secret-SQL", "", 404, "NOT_FOUND"},
		{"POST", "/api/v1/validation", "{}", 405, "METHOD_NOT_ALLOWED"},
		{"GET", "/api/v1/direct-failure", "", 500, "INTERNAL_ERROR"},
		{"GET", "/api/v1/returned-failure", "", 500, "INTERNAL_ERROR"},
		{"GET", "/api/v1/panic", "", 500, "INTERNAL_ERROR"},
		{"GET", "/api/v1/panic-framework", "", 500, "INTERNAL_ERROR"},
	}
	for _, tt := range cases {
		t.Run(tt.path+tt.method, func(t *testing.T) {
			before := logs.Len()
			status, body, headers := securityRequest(t, app, tt.method, tt.path, "application/json", tt.body)
			var envelope response.Body
			if json.Unmarshal([]byte(body), &envelope) != nil || status != tt.status || envelope.Success || envelope.Error == nil || envelope.Error.Code != tt.code || envelope.Error.Message == "" || envelope.Message != envelope.Error.Message {
				t.Fatalf("wrong standard error: %d %s", status, body)
			}
			id := headers["X-Request-Id"] // helper canonicalization varies; accept exact standard spelling.
			if id == "" {
				id = headers["X-Request-ID"]
			}
			parsed, err := uuid.Parse(id)
			if err != nil || parsed.Version() != 4 || envelope.Error.RequestID != id || strings.Contains(body, "private") || strings.Contains(body, "goroutine") {
				t.Fatalf("unsafe/error correlation absent: %s", body)
			}
			completions := 0
			for _, entry := range logs.All()[before:] {
				fields := entry.ContextMap()
				if fields["request_id"] != id || strings.Contains(fmt.Sprint(fields), "private-secret-SQL") {
					t.Fatal("log correlation/leakage")
				}
				if fields["event"] == "http.request_completed" {
					completions++
					if fields["status"] != int64(status) {
						t.Fatalf("logged status differs: %v", fields)
					}
					if _, ok := fields["duration_ms"].(float64); !ok {
						t.Fatal("latency must be numeric")
					}
					if tt.status == 429 && entry.Level != zapcore.WarnLevel || tt.status >= 500 && entry.Level != zapcore.ErrorLevel {
						t.Fatal("incorrect completion severity")
					}
				}
			}
			if completions != 1 {
				t.Fatalf("need one final completion: %d", completions)
			}
			if tt.code == "VALIDATION_ERROR" && envelope.Error.Fields["email"] != "is required" {
				t.Fatal("field validation missing")
			}
			if status == 429 && headers["Retry-After"] == "" {
				t.Fatal("retry delay lost")
			}
		})
	}
	if logs.FilterMessage("panic recovered").Len() != 2 || logs.FilterMessage("request error").Len() != 2 {
		t.Fatal("panic/direct/returned failure visibility inconsistent")
	}
	for _, entry := range logs.FilterMessage("panic recovered").All() {
		if entry.ContextMap()["stack"] == nil {
			t.Fatal("panic source stack missing")
		}
	}
	status, body, headers := securityRequest(t, app, "GET", "/api/v1/correlation", "", "")
	var envelope struct {
		Data map[string]string `json:"data"`
	}
	if status != 200 || json.Unmarshal([]byte(body), &envelope) != nil || envelope.Data["requestId"] != headers["X-Request-ID"] {
		t.Fatal("standard context propagation")
	}
}

func TestIncomingIDsAndSecretRedaction(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	log := &logger.Logger{Logger: zap.New(core)}
	app := newServer(httpSecurityConfig(t), nil, log)
	const secret = "synthetic-password-Authorization-Cookie-refresh-JWT-DB-secret"
	app.Post("/api/v1/safe", func(c fiber.Ctx) error {
		c.Set("Set-Cookie", "credential="+secret)
		return response.Error(c, errors.New(secret))
	})
	previous := ""
	for _, incoming := range []string{"", uuid.NewString(), "invalid/" + secret, strings.Repeat("x", 1024)} {
		// Explicit headers include a valid-looking incoming ID that must also be replaced.
		status, body, headers := securityRequestWithHeaders(t, app, "POST", "/api/v1/safe?password="+secret, "application/json", `{"password":"`+secret+`","refresh_token":"`+secret+`"}`, map[string]string{"X-Request-ID": incoming, "Authorization": "Bearer " + secret, "Cookie": "refresh=" + secret})
		id := headers["X-Request-ID"]
		if status != 500 || id == incoming || id == previous || strings.Contains(body, secret) {
			t.Fatal("server-owned ID/redaction failed")
		}
		previous = id
	}
	for _, entry := range logs.All() {
		if strings.Contains(fmt.Sprint(entry.ContextMap()), secret) {
			t.Fatal("log contains request/response credential")
		}
	}
	cfg := httpSecurityConfig(t)
	cfg.JWT.Secret = secret
	cfg.DB.Password = secret
	logStartup(log, cfg)
	if strings.Contains(fmt.Sprint(logs.All()), secret) {
		t.Fatal("startup config secret leak")
	}
}

// Capture ordinary HTTP responses with additional synthetic headers, no listener.
func securityRequestWithHeaders(t *testing.T, app *fiber.App, method, path, media, body string, headers map[string]string) (int, string, map[string]string) {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", media)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for key, values := range res.Header {
		out[key] = strings.Join(values, ",")
	}
	out["X-Request-ID"] = res.Header.Get("X-Request-ID")
	return res.StatusCode, string(raw), out
}

type checkDouble struct {
	err   error
	calls int
	id    string
}

func (d *checkDouble) Check(ctx context.Context) error {
	d.calls++
	d.id = observability.RequestID(ctx)
	return d.err
}
func TestHealthReadinessContracts(t *testing.T) {
	for _, mode := range []string{"healthy", "failed", "absent"} {
		t.Run(mode, func(t *testing.T) {
			d := &checkDouble{}
			var checker handlers.HealthChecker = d
			if mode == "failed" {
				d.err = errors.New("private-db-secret")
			}
			if mode == "absent" {
				checker = nil
			}
			h := handlers.NewHealthHandler(checker)
			app := newServer(httpSecurityConfig(t), nil, nil)
			app.Get("/api/v1/health", h.Health)
			app.Get("/api/v1/ready", h.Ready)
			status, body, headers := securityRequest(t, app, "GET", "/api/v1/health", "", "")
			if status != 200 || d.calls != 0 || headers["Cache-Control"] != "no-store" || !strings.Contains(body, `"service":"elabtrack-v2"`) || strings.Contains(body, "database") {
				t.Fatal("liveness touched/exposed dependency")
			}
			status, body, headers = securityRequest(t, app, "GET", "/api/v1/ready", "", "")
			want := 200
			if mode != "healthy" {
				want = 503
			}
			if status != want || strings.Contains(body, "private") || headers["Cache-Control"] != "no-store" || mode != "absent" && (d.calls != 1 || d.id != headers["X-Request-ID"]) {
				t.Fatal("readiness contract/correlation invalid")
			}
			if want == 503 && !strings.Contains(body, `"code":"SERVICE_UNAVAILABLE"`) {
				t.Fatal("readiness standard error absent")
			}
		})
	}
}
