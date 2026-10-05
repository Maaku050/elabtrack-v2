package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/gofiber/fiber/v3"
)

// Register wires every route group onto the app under /api/v1.
// Add new bounded contexts here without touching the foundation.
func Register(app *fiber.App, deps *Deps) {
	v1 := app.Group("/api/v1")

	// Public health check.
	v1.Get("/health", deps.Health.Health)

	// Feature route groups.
	RegisterAuth(v1, deps.Auth, deps.TokenIssuer)
	RegisterUser(v1, deps.User, deps.TokenIssuer)
}

// Deps bundles the handlers + token issuer needed to register routes.
// Keeping this in one place makes adding new contexts trivial.
type Deps struct {
	Health      *handlers.HealthHandler
	Auth        *handlers.AuthHandler
	User        *handlers.UserHandler
	TokenIssuer application.TokenIssuer
}
