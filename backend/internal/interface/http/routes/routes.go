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
	if deps.AccountManagement != nil {
		RegisterAccounts(v1, deps.AccountManagement, protected)
	}
	if deps.Inventory != nil {
		RegisterInventory(v1, deps.Inventory, protected)
	}
	if deps.Borrowing != nil {
		RegisterBorrowing(v1, deps.Borrowing, protected)
	}
	if deps.Profile != nil {
		RegisterProfile(v1, deps.Profile, protected)
	}
	if deps.Reporting != nil {
		RegisterReporting(v1, deps.Reporting, protected)
	}
	if deps.Notifications != nil {
		RegisterNotifications(v1, deps.Notifications, protected)
	}
	if deps.Terms != nil {
		RegisterTerms(v1, deps.Terms, protected)
	}
}

// Deps bundles the handlers + token issuer needed to register routes.
// Keeping this in one place makes adding new contexts trivial.
type Deps struct {
	Health            *handlers.HealthHandler
	Auth              *handlers.AuthHandler
	User              *handlers.UserHandler
	Terms             *handlers.TermsHandler
	AccountManagement *handlers.AccountsHandler
	Inventory         *handlers.InventoryHandler
	Borrowing         *handlers.BorrowingHandler
	Notifications     *handlers.NotificationsHandler
	Reporting         *handlers.ReportingHandler
	Profile           *handlers.ProfileHandler
	TokenIssuer       application.TokenIssuer
	Accounts          appauth.AccountResolver
	Environment       config.Environment
}
