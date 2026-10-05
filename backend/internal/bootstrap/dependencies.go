package bootstrap

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	"github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	appuser "github.com/Maaku050/elabtrack-v2/backend/internal/application/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/routes"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
)

// Container bundles all wired application services and HTTP handlers.
// It is constructed once during bootstrap and passed to the router.
type Container struct {
	// Application services
	AuthSvc *auth.Service
	UserSvc *appuser.Service

	// HTTP handlers
	Health *handlers.HealthHandler
	Auth   *handlers.AuthHandler
	User   *handlers.UserHandler

	// Outbound ports needed by route registration (auth middleware).
	TokenIssuer application.TokenIssuer
	Accounts    auth.AccountResolver
	Environment config.Environment
}

// buildContainer wires infrastructure -> repositories -> application services
// -> HTTP handlers. This is the single place where dependency injection
// happens, keeping wiring explicit and easy to follow.
func buildContainer(infra *Infrastructure) *Container {
	hasher, issuer := newSecurityAdapters(infra.Config)
	userRepo, authRepo := newRepositories(infra.DB)

	authSvc := auth.NewService(userRepo, authRepo, hasher, issuer, infra.Config.JWT.AccessTTL, infra.Config.JWT.RefreshTTL)
	userSvc := appuser.NewService(userRepo)

	v := validator.New()

	return &Container{
		AuthSvc:     authSvc,
		UserSvc:     userSvc,
		Health:      handlers.NewHealthHandler(infra.Health),
		Auth:        handlers.NewAuthHandler(authSvc, v),
		User:        handlers.NewUserHandler(userSvc, v),
		TokenIssuer: issuer,
		Accounts:    auth.NewAccountResolver(userRepo),
		Environment: infra.Config.App.Env,
	}
}

// routeDeps converts the container into the bundle expected by routes.Register.
func (c *Container) routeDeps() *routes.Deps {
	return &routes.Deps{
		Health:      c.Health,
		Auth:        c.Auth,
		User:        c.User,
		TokenIssuer: c.TokenIssuer,
		Accounts:    c.Accounts,
		Environment: c.Environment,
	}
}
