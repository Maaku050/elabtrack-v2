package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	appauth "github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/gofiber/fiber/v3"
)

// Register wires every route group onto the app under /api/v1.
// Add new bounded contexts here without touching the foundation.
func Register(app *fiber.App, deps *Deps) {
	v1 := app.Group("/api/v1")

	// Public health check.
	v1.Get("/health", deps.Health.Health)
	v1.Get("/ready", deps.Health.Ready)

	// Feature route groups.
	protected := middleware.Auth(deps.TokenIssuer, deps.Accounts)
	RegisterAuth(v1, deps.Auth, protected, deps.Environment)
	RegisterUser(v1, deps.User, protected)
	if deps.Terms != nil {
		RegisterTerms(v1, deps.Terms, protected)
	}
}

// Deps bundles the handlers + token issuer needed to register routes.
// Keeping this in one place makes adding new contexts trivial.
type Deps struct {
	Health      *handlers.HealthHandler
	Auth        *handlers.AuthHandler
	User        *handlers.UserHandler
	Terms       *handlers.TermsHandler
	TokenIssuer application.TokenIssuer
	Accounts    appauth.AccountResolver
	Environment config.Environment
}
