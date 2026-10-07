package routes_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	appauth "github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	appuser "github.com/Maaku050/elabtrack-v2/backend/internal/application/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/routes"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/observability"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func errorWithoutRequestID(t *testing.T, body string) string {
	t.Helper()
	var envelope map[string]any
	if json.Unmarshal([]byte(body), &envelope) != nil {
		t.Fatal("invalid error JSON")
	}
	detail, ok := envelope["error"].(map[string]any)
	if !ok {
		t.Fatal("missing error")
	}
	id, ok := detail["requestId"].(string)
	if !ok {
		t.Fatal("missing correlation")
	}
	parsed, err := uuid.Parse(id)
	if err != nil || parsed.Version() != 4 {
		t.Fatal("invalid request ID")
	}
	delete(detail, "requestId")
	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

type correlatedUsers struct {
	*users
	requestID string
}

func (r *correlatedUsers) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	r.requestID = observability.RequestID(ctx)
	return r.users.FindByEmail(ctx, email)
}
func contractRoutes(r *correlatedUsers, i application.TokenIssuer, log *logger.Logger) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler(log, config.Test)})
	app.Use(middleware.RequestID())
	app.Use(middleware.ClientInfo(config.SecurityConfig{}))
	app.Use(middleware.Logger(log))
	app.Use(middleware.Recovery(log))
	authSvc := appauth.NewService(r, &tokens{}, security.NewBcryptHasher(4), i, security.SHA256RefreshHasher{}, nil, nil, time.Minute, time.Hour)
	v := validator.New()
	routes.Register(app, &routes.Deps{Health: handlers.NewHealthHandler(nil), Auth: handlers.NewAuthHandler(authSvc, v, config.Test, config.SecurityConfig{AllowedOrigins: []string{trustedOrigin}}), User: handlers.NewUserHandler(appuser.NewService(r), v), TokenIssuer: i, Accounts: appauth.NewAccountResolver(r), Environment: config.Test})
	return app
}
func TestRetainedRoutesContractAndSecurityEvents(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	log := &logger.Logger{Logger: zap.New(core)}
	user := fixtureUser()
	user.Role = domainuser.RoleAdmin
	hash, err := security.NewBcryptHasher(4).Hash("synthetic-correct-password")
	if err != nil {
		t.Fatal(err)
	}
	user.Password = hash
	r := &correlatedUsers{users: &users{account: user}}
	app := contractRoutes(r, &issuer{claims: application.Claims{UserID: selfID}}, log)
	cases := []struct {
		method, path, body, auth string
		status                   int
		code                     string
	}{
		{"POST", "/api/v1/auth/login", `{"email":"current@example.invalid","password":"synthetic-wrong-password"}`, "", 401, "INVALID_CREDENTIALS"},
		{"POST", "/api/v1/auth/login", `{`, "", 400, "BAD_REQUEST"},
		{"POST", "/api/v1/auth/login", `{}`, "", 400, "VALIDATION_ERROR"},
		{"GET", "/api/v1/auth/me", "", "", 401, "UNAUTHORIZED"},
		{"PATCH", "/api/v1/users/me", `{"name":" "}`, "Bearer synthetic", 400, "VALIDATION_ERROR"},
		{"GET", "/api/v1/users/?page=abc", "", "Bearer synthetic", 400, "VALIDATION_ERROR"},
		{"GET", "/api/v1/users/?page=2147483647&per_page=100", "", "Bearer synthetic", 400, "VALIDATION_ERROR"},
		{"GET", "/api/v1/users/?per_page=0&sort=password&order=sideways", "", "Bearer synthetic", 400, "VALIDATION_ERROR"},
		{"GET", "/api/v1/users/?search=" + strings.Repeat("x", 257), "", "Bearer synthetic", 400, "VALIDATION_ERROR"},
		{"GET", "/api/v1/auth/me", "", "Bearer synthetic", 200, ""},
		{"GET", "/api/v1/users/me", "", "Bearer synthetic", 200, ""},
		{"GET", "/api/v1/users/", "", "Bearer synthetic", 200, ""},
		{"PATCH", "/api/v1/users/me", `{"name":"Updated"}`, "Bearer synthetic", 200, ""},
		{"POST", "/api/v1/auth/logout", "", "", 204, ""},
	}
	for _, tt := range cases {
		t.Run(tt.method+tt.path, func(t *testing.T) {
			status, envelope, raw := request(t, app, tt.method, tt.path, tt.body, tt.auth)
			if status != tt.status || strings.Contains(raw, "synthetic") || strings.Contains(raw, hash) {
				t.Fatalf("retained route contract invalid: %d %s", status, raw)
			}
			if status >= 400 {
				detail := envelope["error"].(map[string]any)
				if detail["code"] != tt.code || detail["message"] != envelope["message"] {
					t.Fatal("inconsistent error")
				}
				errorWithoutRequestID(t, raw)
			} else if status != 204 && envelope["success"] != true {
				t.Fatal("success envelope absent")
			}
		})
	}
	if r.requestID == "" {
		t.Fatal("application/repository correlation absent")
	}
	if logs.FilterField(zap.String("security_event", "auth.login_failed")).Len() != 3 || logs.FilterField(zap.String("security_event", "auth.logout_completed")).Len() != 1 {
		t.Fatal("safe auth security events absent")
	}
	status, _, raw := request(t, app, "POST", "/api/v1/auth/login", `{"email":"current@example.invalid","password":"synthetic-correct-password"}`, "")
	if status != 200 || !strings.Contains(raw, "synthetic-access") || logs.FilterField(zap.String("security_event", "auth.login_succeeded")).Len() != 1 || strings.Contains(fmt.Sprint(logs.All()), "synthetic-access") || strings.Contains(fmt.Sprint(logs.All()), strings.Repeat("a", 64)) {
		t.Fatal("successful auth event leaks access/refresh credential or is missing")
	}
	// Infrastructure lookup still has the identical public invalid-credentials
	// contract, but the operational failure is visible at ERROR with a safe class.
	r.lookupErr = errors.New("password-cookie-refresh-SQL-sentinel")
	var envelope map[string]any
	status, envelope, raw = request(t, app, "POST", "/api/v1/auth/login", `{"email":"current@example.invalid","password":"synthetic-wrong-password"}`, "")
	if status != 401 || envelope["error"].(map[string]any)["code"] != "INVALID_CREDENTIALS" || strings.Contains(raw, "sentinel") {
		t.Fatal("lookup failure enumerates account")
	}
	failures := logs.FilterMessage("request error").All()
	if len(failures) != 1 || failures[0].ContextMap()["request_id"] != r.requestID || failures[0].ContextMap()["operation"] != "auth.login_lookup" || strings.Contains(fmt.Sprint(logs.All()), "sentinel") {
		t.Fatal("masked failure unobservable/unsafe")
	}
	// Non-lookup internal failures continue to map to 500.
	if mapped := response.Map(shared.Internal("auth.issue_access", errors.New("private"))).Status; mapped != 500 {
		t.Fatal("unexpected failures must remain internal")
	}
}

type issuanceFailureIssuer struct{ *issuer }

func (i *issuanceFailureIssuer) IssueAccessToken(context.Context, *domainuser.User) (string, error) {
	return "", errors.New("private-signing-password-token-sentinel")
}
func TestLoginIssuanceFailureUsesInternalContract(t *testing.T) {
	core, logs := observer.New(zap.DebugLevel)
	user := fixtureUser()
	hash, err := security.NewBcryptHasher(4).Hash("synthetic-correct-password")
	if err != nil {
		t.Fatal(err)
	}
	user.Password = hash
	r := &correlatedUsers{users: &users{account: user}}
	app := contractRoutes(r, &issuanceFailureIssuer{issuer: &issuer{claims: application.Claims{UserID: selfID}}}, &logger.Logger{Logger: zap.New(core)})
	status, envelope, raw := request(t, app, "POST", "/api/v1/auth/login", `{"email":"current@example.invalid","password":"synthetic-correct-password"}`, "")
	if status != 500 || envelope["error"].(map[string]any)["code"] != "INTERNAL_ERROR" || strings.Contains(raw, "sentinel") || logs.FilterMessage("request error").Len() != 1 || strings.Contains(fmt.Sprint(logs.All()), "sentinel") {
		t.Fatal("unexpected issuance failure was masked, leaked or unobservable")
	}
}
