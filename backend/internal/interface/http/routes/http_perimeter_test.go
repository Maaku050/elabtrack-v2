package routes_test

import (
	"errors"
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
)

func perimeterAuthApp(t *testing.T, r *users, sessions *httpSessions) *fiber.App {
	t.Helper()
	cfg, err := config.Parse(map[string]string{"FRONTEND_URL": trustedOrigin, "ALLOWED_ORIGINS": trustedOrigin})
	if err != nil {
		t.Fatal(err)
	}
	app := fiber.New(fiber.Config{BodyLimit: cfg.App.BodyLimit, ErrorHandler: middleware.ErrorHandler(nil, cfg.App.Env)})
	app.Use(middleware.Recovery(nil))
	app.Use(middleware.RequestID())
	app.Use(middleware.ClientInfo(cfg.Security))
	app.Use(middleware.SecurityHeaders(cfg.App.Env))
	app.Use(middleware.CORS(cfg.Security))
	app.Use(middleware.RateLimit(cfg.Security))
	app.Use(middleware.RequestSafety())
	i := security.NewJWTIssuer(cfg.JWT)
	v := validator.New()
	svc := appauth.NewService(r, sessions, security.NewBcryptHasher(4), i, security.SHA256RefreshHasher{}, httpTransaction{}, httpAccounts{r}, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	routes.Register(app, &routes.Deps{Health: handlers.NewHealthHandler(nil), Auth: handlers.NewAuthHandler(svc, v, cfg.App.Env, cfg.Security), User: handlers.NewUserHandler(appuser.NewService(r), v), TokenIssuer: i, Accounts: appauth.NewAccountResolver(r), Environment: cfg.App.Env})
	return app
}

func TestRefreshRotationFitsDefaultPerimeterBudget(t *testing.T) {
	raw := strings.Repeat("a", 64)
	hash, _ := (security.SHA256RefreshHasher{}).Hash(raw)
	sessions := &httpSessions{session: domainauth.NewRefreshToken(hash, selfID, time.Now().Add(-time.Minute), time.Now().Add(time.Hour))}
	app := perimeterAuthApp(t, &users{account: fixtureUser()}, sessions)
	original := raw
	for i := 0; i < 8; i++ {
		status, _, cookies := cookieRequest(t, app, "/api/v1/auth/refresh", raw, trustedOrigin, "")
		if status != 200 || len(cookies) != 1 || len(sessions.created) != i+1 || sessions.session.RevokedAt == nil {
			t.Fatal("normal bootstrap/rotation denied or single use weakened")
		}
		raw = cookies[0].Value
		sessions.session = sessions.created[i]
	}
	status, _, cookies := cookieRequest(t, app, "/api/v1/auth/refresh", original, trustedOrigin, "")
	if status != 401 || len(cookies) != 1 || sessions.session.RevokedAt != nil {
		t.Fatal("replay semantics weakened by perimeter")
	}
	status, _, _ = cookieRequest(t, app, "/api/v1/auth/refresh", raw, trustedOrigin, "")
	if status != 200 {
		t.Fatal("replacement cannot continue within ordinary refresh budget")
	}
	for _, path := range []string{"/api/v1/auth/login", "/api/v1/auth/register", "/api/v1/auth/refresh", "/api/v1/auth/logout"} {
		before := len(sessions.created)
		status, _, _ = request(t, app, "GET", path, "", "")
		if (status != 404 && status != 405) || len(sessions.created) != before {
			t.Fatal("state-changing GET auth behavior")
		}
	}
}

func TestLoginFailureEnvelopeDoesNotEnumerateAccounts(t *testing.T) {
	hashed, err := security.NewBcryptHasher(4).Hash("synthetic-password")
	if err != nil {
		t.Fatal(err)
	}
	var expected string
	for _, name := range []string{"missing", "wrong password", "inactive", "unknown role", "lookup failure"} {
		t.Run(name, func(t *testing.T) {
			r := &users{account: fixtureUser()}
			r.account.Password = hashed
			password := "synthetic-password"
			switch name {
			case "missing":
				r.account = nil
			case "wrong password":
				password = "incorrect-password"
			case "inactive":
				r.account.IsActive = false
			case "unknown role":
				r.account.Role = domainuser.Role("unaccepted")
			case "lookup failure":
				r.lookupErr = errors.New("SQL private lookup sentinel")
			}
			sessions := &httpSessions{}
			app := perimeterAuthApp(t, r, sessions)
			status, body, cookies := cookieRequest(t, app, "/api/v1/auth/login", "", trustedOrigin, `{"email":"current@example.invalid","password":"`+password+`"}`)
			if status != 401 || len(cookies) != 0 || len(sessions.created) != 0 || strings.Contains(body, password) || strings.Contains(body, "inactive") || strings.Contains(body, "SQL") {
				t.Fatal("login failure enumerated account")
			}
			comparable := errorWithoutRequestID(t, body)
			if expected == "" {
				expected = comparable
			} else if comparable != expected {
				t.Fatal("login failure envelopes differ")
			}
		})
	}
}
