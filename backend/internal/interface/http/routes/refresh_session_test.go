package routes_test

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	appauth "github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	appuser "github.com/Maaku050/elabtrack-v2/backend/internal/application/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/routes"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type httpSessions struct {
	domainauth.Repository
	session              *domainauth.RefreshToken
	lookupErr, revokeErr error
	revokeHash           string
	created              []*domainauth.RefreshToken
}

func (r *httpSessions) FindByHashForUpdate(_ context.Context, hash string) (*domainauth.RefreshToken, error) {
	if r.lookupErr != nil {
		return nil, r.lookupErr
	}
	if r.session == nil || r.session.TokenHash != hash {
		return nil, domainauth.ErrTokenNotFound
	}
	return r.session, nil
}
func (r *httpSessions) Revoke(_ context.Context, hash string) error {
	r.revokeHash = hash
	if r.revokeErr != nil {
		return r.revokeErr
	}
	if r.session != nil && r.session.TokenHash == hash && r.session.RevokedAt == nil {
		now := time.Now().UTC()
		r.session.RevokedAt = &now
	}
	return nil
}

func (r *httpSessions) Create(_ context.Context, session *domainauth.RefreshToken) error {
	r.created = append(r.created, session)
	return nil
}
func (r *httpSessions) Consume(_ context.Context, hash string, replacement uuid.UUID, now time.Time) error {
	if r.session == nil || r.session.TokenHash != hash || !r.session.ValidAt(now) {
		return domainauth.ErrTokenInvalid
	}
	r.session.RevokedAt = &now
	r.session.ReplacedBy = &replacement
	return nil
}

type httpTransaction struct{}

func (httpTransaction) Within(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type httpAccounts struct{ users *users }

func (a httpAccounts) LockAccountByID(ctx context.Context, id uuid.UUID) (*domainuser.Account, error) {
	return a.users.FindAccountByID(ctx, id)
}

func refreshHTTPApp(r *users, sessions *httpSessions) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.ErrorHandler(nil, config.Development)})
	i := &issuer{}
	svc := appauth.NewService(r, sessions, security.NewBcryptHasher(4), i, security.SHA256RefreshHasher{}, httpTransaction{}, httpAccounts{r}, time.Minute, time.Hour)
	v := validator.New()
	routes.Register(app, &routes.Deps{Auth: handlers.NewAuthHandler(svc, v, config.Development, config.SecurityConfig{AllowedOrigins: []string{"https://app.example.invalid"}}), User: handlers.NewUserHandler(appuser.NewService(r), v), Health: handlers.NewHealthHandler(nil), TokenIssuer: i, Accounts: appauth.NewAccountResolver(r), Environment: config.Development})
	return app
}
func TestRefreshFailuresHaveSameUnauthorizedResponse(t *testing.T) {
	raw := strings.Repeat("a", 64)
	hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
	var canonical string
	for _, name := range []string{"expired", "revoked", "replayed", "unknown", "malformed", "inactive", "deleted", "unknown role"} {
		t.Run(name, func(t *testing.T) {
			now := time.Now().UTC()
			sessions := &httpSessions{session: domainauth.NewRefreshToken(hash, selfID, now.Add(-time.Hour), now.Add(time.Hour))}
			r := &users{account: fixtureUser()}
			presented := raw
			switch name {
			case "expired":
				sessions.session.ExpiresAt = now.Add(-time.Hour)
			case "revoked":
				sessions.session.RevokedAt = &now
			case "replayed":
				sessions.session.RevokedAt = &now
				replacement := otherID
				sessions.session.ReplacedBy = &replacement
			case "unknown":
				sessions.session = nil
			case "malformed":
				presented = "secret-sentinel"
			case "inactive":
				r.account.IsActive = false
			case "deleted":
				r.account = nil
			case "unknown role":
				r.account.Role = "superadmin"
			}
			status, body, _ := cookieRequest(t, refreshHTTPApp(r, sessions), "/api/v1/auth/refresh", presented, "https://app.example.invalid", "")
			if status != http.StatusUnauthorized {
				t.Fatalf("status=%d", status)
			}
			if canonical == "" {
				canonical = body
			}
			if body != canonical || strings.Contains(body, raw) || strings.Contains(body, hash) || strings.Contains(body, presented) {
				t.Fatal("token-existence detail/credential leaked")
			}
		})
	}
}
func TestLogoutHTTPIdempotentAndSecretSafe(t *testing.T) {
	raw := strings.Repeat("a", 64)
	hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
	now := time.Now().UTC()
	sessions := &httpSessions{session: domainauth.NewRefreshToken(hash, selfID, now, now.Add(time.Hour))}
	app := refreshHTTPApp(&users{account: fixtureUser()}, sessions)
	for _, presented := range []string{raw, raw, strings.Repeat("b", 64), "secret-sentinel"} {
		status, body, _ := cookieRequest(t, app, "/api/v1/auth/logout", presented, "https://app.example.invalid", "")
		if status != http.StatusNoContent || body != "" {
			t.Fatal("logout exposes token existence")
		}
	}
	if sessions.session.RevokedAt == nil {
		t.Fatal("logout did not revoke current session")
	}
	status, _, _ := cookieRequest(t, app, "/api/v1/auth/refresh", raw, "https://app.example.invalid", "")
	if status != http.StatusUnauthorized {
		t.Fatal("logged-out session refreshed")
	}
	sessions.revokeErr = errors.New(raw + hash)
	status, body, _ := cookieRequest(t, app, "/api/v1/auth/logout", raw, "https://app.example.invalid", "")
	if status != http.StatusInternalServerError || strings.Contains(body, raw) || strings.Contains(body, hash) {
		t.Fatal("logout storage failure not safe")
	}
}
func TestRefreshHTTPStorageErrorIsSafeInternalFailure(t *testing.T) {
	raw := strings.Repeat("a", 64)
	hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
	sessions := &httpSessions{lookupErr: errors.New(raw + hash + " SQL credential-sentinel")}
	status, body, _ := cookieRequest(t, refreshHTTPApp(&users{account: fixtureUser()}, sessions), "/api/v1/auth/refresh", raw, "https://app.example.invalid", "")
	if status != http.StatusInternalServerError || strings.Contains(body, raw) || strings.Contains(body, hash) || strings.Contains(body, "SQL") {
		t.Fatal("unsafe storage response")
	}
}
