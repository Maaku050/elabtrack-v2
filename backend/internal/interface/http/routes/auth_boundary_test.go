package routes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	appauth "github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	appuser "github.com/Maaku050/elabtrack-v2/backend/internal/application/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/routes"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

var selfID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
var otherID = uuid.MustParse("00000000-0000-0000-0000-000000000002")

// A deliberately small stateful fake distinguishes safe-account lookup from
// aggregate/password lookup and records writes/privileged list execution.
type users struct {
	domainuser.Repository
	account                                          *domainuser.User
	lookupErr                                        error
	lookups, aggregateReads, updates, lists, creates int
	beforeUpdate                                     func()
}

func fixtureUser() *domainuser.User {
	return &domainuser.User{ID: selfID, Email: "current@example.invalid", Name: "Current User", Password: "private-hash-sentinel", Role: domainuser.RoleBorrower, IsActive: true, CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
}
func snapshot(u *domainuser.User) *domainuser.Account {
	if u == nil {
		return nil
	}
	return &domainuser.Account{ID: u.ID, Email: u.Email, Name: u.Name, Role: u.Role, IsActive: u.IsActive, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt}
}
func (r *users) FindAccountByID(_ context.Context, id uuid.UUID) (*domainuser.Account, error) {
	r.lookups++
	if r.lookupErr != nil {
		return nil, r.lookupErr
	}
	if r.account == nil || r.account.ID != id {
		return nil, domainuser.ErrUserNotFound
	}
	return snapshot(r.account), nil
}
func (r *users) FindByID(_ context.Context, id uuid.UUID) (*domainuser.User, error) {
	r.aggregateReads++
	return nil, errors.New("protected requests must not load password aggregates")
}
func (r *users) FindByEmail(_ context.Context, email string) (*domainuser.User, error) {
	if r.lookupErr != nil {
		return nil, r.lookupErr
	}
	if r.account == nil || r.account.Email != email {
		return nil, domainuser.ErrUserNotFound
	}
	return r.account, nil
}
func (r *users) Create(_ context.Context, u *domainuser.User) error {
	r.creates++
	r.account = u
	return nil
}
func (r *users) Update(_ context.Context, u *domainuser.User) error {
	return errors.New("self update must not perform aggregate write")
}
func (r *users) UpdateProfile(_ context.Context, id uuid.UUID, name string) (*domainuser.Account, error) {
	if r.beforeUpdate != nil {
		r.beforeUpdate()
	}
	if r.account == nil || r.account.ID != id || !r.account.IsActive || !r.account.Role.Valid() {
		return nil, shared.ErrForbidden
	}
	r.updates++
	r.account.Name = name
	return snapshot(r.account), nil
}
func (r *users) List(_ context.Context, _ shared.Page) ([]*domainuser.User, int, error) {
	r.lists++
	return []*domainuser.User{r.account}, 1, nil
}

type tokens struct {
	domainauth.Repository
	created []*domainauth.RefreshToken
}

func (r *tokens) Create(_ context.Context, t *domainauth.RefreshToken) error {
	r.created = append(r.created, t)
	return nil
}

type issuer struct {
	application.TokenIssuer
	claims application.Claims
	err    error
	calls  int
}

func (i *issuer) VerifyAccessToken(_ context.Context, _ string) (application.Claims, error) {
	i.calls++
	return i.claims, i.err
}
func (i *issuer) IssueAccessToken(_ context.Context, _ *domainuser.User) (string, error) {
	return "synthetic-access", nil
}
func (i *issuer) GenerateRefreshToken() (string, error) { return strings.Repeat("a", 64), nil }
func newApp(env config.Environment, r *users, i application.TokenIssuer, refresh *tokens) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler(nil, env)})
	accountService := appauth.NewAccountResolver(r)
	authService := appauth.NewService(r, refresh, security.NewBcryptHasher(4), i, security.SHA256RefreshHasher{}, nil, nil, time.Minute, time.Hour)
	v := validator.New()
	routes.Register(app, &routes.Deps{
		Health: handlers.NewHealthHandler(nil), Auth: handlers.NewAuthHandler(authService, v, env, config.SecurityConfig{AllowedOrigins: []string{"https://app.example.invalid"}}), User: handlers.NewUserHandler(appuser.NewService(r), v),
		TokenIssuer: i, Accounts: accountService, Environment: env,
	})
	return app
}
func request(t *testing.T, app *fiber.App, method, path, body, header string) (int, map[string]any, string) {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://app.example.invalid")
	if header != "" {
		req.Header.Set("Authorization", header)
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
	var payload map[string]any
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &payload); err != nil {
			t.Fatal("expected coherent JSON envelope", string(raw))
		}
	}
	return res.StatusCode, payload, string(raw)
}
func TestProtectedRouteAuthenticationMatrix(t *testing.T) {
	for _, route := range []struct{ method, path, body string }{{"GET", "/api/v1/auth/me", ""}, {"GET", "/api/v1/users/me", ""}, {"PATCH", "/api/v1/users/me", `{"name":"Renamed"}`}, {"GET", "/api/v1/users/", ""}} {
		for _, header := range []string{"", "Basic credentials", "Bearer", "Bearer one two", "Bearer invalid"} {
			t.Run(route.method+route.path+header, func(t *testing.T) {
				r := &users{account: fixtureUser()}
				i := &issuer{claims: application.Claims{UserID: selfID, Role: "admin"}, err: errors.New("token-secret-sentinel")}
				status, payload, raw := request(t, newApp(config.Production, r, i, &tokens{}), route.method, route.path, route.body, header)
				if status != 401 || payload["success"] != false {
					t.Fatal("authentication must deny with 401")
				}
				if r.lookups != 0 || r.updates != 0 || r.lists != 0 || strings.Contains(raw, "sentinel") {
					t.Fatal("invalid authentication reached database or leaked verifier data")
				}
			})
		}
	}
}
func TestCurrentAccountAndAuthorizationMatrix(t *testing.T) {
	for _, tt := range []struct {
		name               string
		role               domainuser.Role
		active, missing    bool
		tokenRole          string
		wantSelf, wantList int
	}{
		{"active user", domainuser.RoleBorrower, true, false, "user", 200, 403},
		{"active staff", domainuser.RoleStaff, true, false, "ADMIN", 200, 403},
		{"active admin", domainuser.RoleAdmin, true, false, "admin", 200, 200},
		{"demoted with stale admin token", domainuser.RoleBorrower, true, false, "admin", 200, 403},
		{"promoted with old user token", domainuser.RoleAdmin, true, false, "user", 200, 200},
		{"inactive admin", domainuser.RoleAdmin, false, false, "admin", 403, 403},
		{"unknown role", domainuser.Role("super_admin"), true, false, "admin", 403, 403},
		{"missing account", domainuser.RoleBorrower, true, true, "admin", 401, 401},
	} {
		for _, path := range []string{"/api/v1/auth/me", "/api/v1/users/me", "/api/v1/users/"} {
			t.Run(tt.name+path, func(t *testing.T) {
				r := &users{account: fixtureUser()}
				r.account.Role = tt.role
				r.account.IsActive = tt.active
				if tt.missing {
					r.account = nil
				}
				i := &issuer{claims: application.Claims{UserID: selfID, Role: tt.tokenRole, Email: "stale@example.invalid"}}
				status, payload, raw := request(t, newApp(config.Production, r, i, &tokens{}), "GET", path, `{"role":"admin","id":"`+otherID.String()+`"}`, "Bearer verified")
				want := tt.wantSelf
				if path == "/api/v1/users/" {
					want = tt.wantList
				}
				if status != want {
					t.Fatalf("status %d, expected %d: %s", status, want, raw)
				}
				if r.lookups != 1 || r.aggregateReads != 0 {
					t.Fatal("must resolve exactly one safe current-account snapshot")
				}
				if strings.Contains(raw, "private-hash-sentinel") || strings.Contains(raw, "stale@example.invalid") {
					t.Fatal("unsafe or stale fields returned")
				}
				if path == "/api/v1/users/" && (r.lists == 1) != (want == 200) {
					t.Fatal("unauthorized list handler ran")
				}
				if status == 200 && path != "/api/v1/users/" {
					data := payload["data"].(map[string]any)
					if data["id"] != selfID.String() || data["role"] != string(tt.role) || data["is_active"] != true {
						t.Fatal("body/token data displaced trusted principal")
					}
				}
			})
		}
	}
}
func TestPrincipalContextAndSafeMeResponse(t *testing.T) {
	r := &users{account: fixtureUser()}
	i := &issuer{claims: application.Claims{UserID: selfID, Role: "admin"}}
	app := newApp(config.Production, r, i, &tokens{})
	app.Get("/probe", middleware.Auth(i, appauth.NewAccountResolver(r)), func(c fiber.Ctx) error {
		principal, ok := middleware.PrincipalFromContext(c)
		identity, identityOK := middleware.IdentityFromContext(c)
		if !ok || !identityOK || principal.ID != selfID || identity.UserID != selfID || principal.Role != domainuser.RoleBorrower {
			return fmt.Errorf("typed context invalid")
		}
		return c.JSON(principal.UserDTO())
	})
	status, _, _ := request(t, app, "GET", "/probe", "", "Bearer verified")
	if status != 200 {
		t.Fatal("typed principal unavailable")
	}
	status, payload, raw := request(t, app, "GET", "/api/v1/auth/me", "", "Bearer verified")
	if status != 200 {
		t.Fatal(raw)
	}
	data := payload["data"].(map[string]any)
	allowed := map[string]bool{"id": true, "email": true, "name": true, "role": true, "is_active": true}
	if len(data) != len(allowed) {
		t.Fatal("unexpected current-user response fields")
	}
	for key := range data {
		if !allowed[key] {
			t.Fatalf("unsafe current-user field: %s", key)
		}
	}
	// Reused HTTP app/context must not leak a prior principal into an anonymous request.
	status, _, _ = request(t, app, "GET", "/api/v1/auth/me", "", "")
	if status != 401 {
		t.Fatal("prior request identity leaked")
	}
}
func TestResolutionErrorsAndZeroIdentityFailClosed(t *testing.T) {
	for _, tt := range []struct {
		name          string
		id            uuid.UUID
		repoErr       error
		want, lookups int
	}{
		{"zero identity", uuid.Nil, nil, 401, 0},
		{"database unavailable", selfID, errors.New("SQL-password-sentinel"), 500, 1},
		{"wrapped missing account", selfID, fmt.Errorf("lookup: %w", domainuser.ErrUserNotFound), 401, 1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := &users{account: fixtureUser(), lookupErr: tt.repoErr}
			i := &issuer{claims: application.Claims{UserID: tt.id, Role: "admin"}}
			status, _, raw := request(t, newApp(config.Production, r, i, &tokens{}), "GET", "/api/v1/auth/me", "", "Bearer verified")
			if status != tt.want || r.lookups != tt.lookups || strings.Contains(raw, "sentinel") {
				t.Fatal("resolution failure must deny safely")
			}
		})
	}
}
func TestSelfUpdateRejectsPrivilegedFieldsAndArbitraryUsers(t *testing.T) {
	for _, field := range []string{"id", "user_id", "email", "role", "is_active", "password", "password_hash", "refresh_token", "created_at", "updated_at"} {
		t.Run(field, func(t *testing.T) {
			r := &users{account: fixtureUser()}
			i := &issuer{claims: application.Claims{UserID: selfID}}
			body := `{"name":"Renamed","` + field + `":"injected-security-sentinel"}`
			status, _, raw := request(t, newApp(config.Production, r, i, &tokens{}), "PATCH", "/api/v1/users/me", body, "Bearer verified")
			if status != 400 || r.updates != 0 || r.account.Role != domainuser.RoleBorrower || !r.account.IsActive || strings.Contains(raw, "sentinel") {
				t.Fatal("mass assignment not contained")
			}
		})
	}
	r := &users{account: fixtureUser()}
	i := &issuer{claims: application.Claims{UserID: selfID, Role: "admin"}}
	app := newApp(config.Production, r, i, &tokens{})
	for _, method := range []string{"GET", "PATCH", "DELETE", "POST"} {
		status, _, _ := request(t, app, method, "/api/v1/users/"+otherID.String(), `{"role":"admin"}`, "Bearer verified")
		if status != 404 {
			t.Fatal("arbitrary-account endpoint must remain unavailable")
		}
	}
	status, payload, raw := request(t, app, "PATCH", "/api/v1/users/me", `{"name":"Renamed User"}`, "Bearer verified")
	if status != 200 || r.updates != 1 || payload["data"].(map[string]any)["id"] != selfID.String() {
		t.Fatal("self name update failed", raw)
	}
	if r.account.Password != "private-hash-sentinel" || r.account.Email != "current@example.invalid" || r.account.Role != domainuser.RoleBorrower || !r.account.IsActive {
		t.Fatal("self update modified security fields")
	}
}
func TestSelfUpdateDeniesConcurrentDeactivation(t *testing.T) {
	r := &users{account: fixtureUser()}
	r.beforeUpdate = func() { r.account.IsActive = false }
	i := &issuer{claims: application.Claims{UserID: selfID, Role: "admin"}}
	status, _, _ := request(t, newApp(config.Production, r, i, &tokens{}), "PATCH", "/api/v1/users/me", `{"name":"Renamed"}`, "Bearer verified")
	if status != 403 || r.updates != 0 || r.account.IsActive || r.account.Name != "Current User" {
		t.Fatal("deactivation must survive an in-flight self update")
	}
}
func TestRegistrationEnvironmentContainment(t *testing.T) {
	for _, env := range []config.Environment{config.Production, config.Development, config.Test, "unknown", ""} {
		for _, payload := range []string{`{"email":"synthetic@example.invalid","name":"Local User","password":"test-local-password"}`, `{"role":"ADMIN"}`} {
			r, refresh := &users{}, &tokens{}
			status, _, _ := request(t, newApp(env, r, &issuer{}, refresh), "POST", "/api/v1/auth/register", payload, "")
			if status != 404 || r.creates != 0 || len(refresh.created) != 0 {
				t.Fatal("public registration must be absent in every environment")
			}
		}
	}
}
func TestActualJWTVerifierFeedsOnlyIdentity(t *testing.T) {
	real := security.NewJWTIssuer(config.JWTConfig{Secret: "synthetic-signing-key-for-tests-only", Issuer: "test", AccessTTL: time.Minute})
	r := &users{account: fixtureUser()}
	claimed := *r.account
	claimed.Role = domainuser.RoleAdmin
	token, err := real.IssueAccessToken(context.Background(), &claimed)
	if err != nil {
		t.Fatal(err)
	}
	app := newApp(config.Production, r, real, &tokens{})
	status, _, _ := request(t, app, http.MethodGet, "/api/v1/users/", "", "Bearer "+token)
	if status != 403 {
		t.Fatal("real stale JWT granted privilege")
	}
	status, _, _ = request(t, app, http.MethodGet, "/api/v1/auth/me", "", "Bearer "+token)
	if status != 200 {
		t.Fatal("valid real token failed identity resolution")
	}
	status, _, _ = request(t, app, http.MethodGet, "/api/v1/auth/me", "", "Bearer malformed-token")
	if status != 401 {
		t.Fatal("malformed JWT accepted")
	}
}
func TestUnknownConfiguredRoleAndMissingBoundaryDenied(t *testing.T) {
	r := &users{account: fixtureUser()}
	i := &issuer{claims: application.Claims{UserID: selfID}}
	app := newApp(config.Production, r, i, &tokens{})
	app.Get("/unknown-required-role", middleware.Auth(i, appauth.NewAccountResolver(r)), middleware.RequireRole("unknown"), func(c fiber.Ctx) error { return c.SendStatus(200) })
	app.Get("/without-auth", middleware.RequireRole(domainuser.RoleAdmin), func(c fiber.Ctx) error { return c.SendStatus(200) })
	status, _, _ := request(t, app, "GET", "/unknown-required-role", "", "Bearer verified")
	if status != 403 {
		t.Fatal("unknown configured role accepted")
	}
	status, _, _ = request(t, app, "GET", "/without-auth", "", "Bearer verified")
	if status != 401 {
		t.Fatal("role middleware must not authenticate a request")
	}
}
func TestLoginDoesNotRevealInactiveAccountBeforePassword(t *testing.T) {
	hashed, err := security.NewBcryptHasher(4).Hash("correct-local-password")
	if err != nil {
		t.Fatal(err)
	}
	r := &users{account: fixtureUser()}
	r.account.IsActive = false
	r.account.Password = hashed
	app := newApp(config.Production, r, &issuer{}, &tokens{})
	for password, want := range map[string]int{"wrong-local-password": 401, "correct-local-password": 401} {
		body := `{"email":"current@example.invalid","password":"` + password + `"}`
		status, _, raw := request(t, app, "POST", "/api/v1/auth/login", body, "")
		if status != want || strings.Contains(raw, "inactive") {
			t.Fatal("login exposed account status")
		}
		if strings.Contains(raw, password) || strings.Contains(raw, hashed) {
			t.Fatal("password leaked")
		}
	}
	// Preserve the Phase 1B service ordering assertion independently of the now
	// uniform HTTP failures: correct password reaches status denial; wrong does not.
	svc := appauth.NewService(r, &tokens{}, security.NewBcryptHasher(4), &issuer{}, security.SHA256RefreshHasher{}, nil, nil, time.Minute, time.Hour)
	for password, want := range map[string]error{"wrong-local-password": domainuser.ErrInvalidCredentials, "correct-local-password": domainuser.ErrUserInactive} {
		if _, err := svc.Login(context.Background(), appauth.LoginRequest{Email: "current@example.invalid", Password: password}); !errors.Is(err, want) {
			t.Fatal("credential/status ordering changed")
		}
	}
}
