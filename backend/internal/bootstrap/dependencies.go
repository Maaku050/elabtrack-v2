package bootstrap

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	appaccounts "github.com/Maaku050/elabtrack-v2/backend/internal/application/accounts"
	"github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	appborrowing "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	appinventory "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	appnotifications "github.com/Maaku050/elabtrack-v2/backend/internal/application/notifications"
	appprofile "github.com/Maaku050/elabtrack-v2/backend/internal/application/profile"
	appreporting "github.com/Maaku050/elabtrack-v2/backend/internal/application/reporting"
	appterms "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	appuser "github.com/Maaku050/elabtrack-v2/backend/internal/application/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/email"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/profileimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/spreadsheet"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/routes"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
)

// Container bundles all wired application services and HTTP handlers.
// It is constructed once during bootstrap and passed to the router.
type Container struct {
	// Application services
	AuthSvc              *auth.Service
	UserSvc              *appuser.Service
	TermsSvc             *appterms.Service
	AccountManagementSvc *appaccounts.Service
	InventorySvc         *appinventory.Service
	BorrowingSvc         *appborrowing.Service
	NotificationsSvc     *appnotifications.Service
	ReportingSvc         *appreporting.Service
	ProfileSvc           *appprofile.Service

	// HTTP handlers
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

	authSvc := auth.NewService(userRepo, authRepo, hasher, issuer, security.SHA256RefreshHasher{}, infra.Tx, userRepo, infra.Config.JWT.AccessTTL, infra.Config.JWT.RefreshTTL)
	userSvc := appuser.NewService(userRepo)
	termsSvc := appterms.NewService(postgres.NewTermsRepository(infra.DB.Pool), userRepo, infra.Tx)

	accountsSvc := appaccounts.NewService(postgres.NewAccountsRepository(infra.DB.Pool), infra.Tx, hasher, email.NewBrevo(infra.Config.Accounts.BrevoKey, infra.Config.Accounts.SenderEmail, infra.Config.Accounts.SenderName), infra.Config.Accounts.Policy)
	inventorySvc := appinventory.NewService(postgres.NewInventoryRepository(infra.DB.Pool), postgres.NewAccountsRepository(infra.DB.Pool), infra.Tx, catalogimage.Validator{})
	borrowingSvc := appborrowing.NewService(postgres.NewBorrowingRepository(infra.DB.Pool), postgres.NewAccountsRepository(infra.DB.Pool), postgres.NewInventoryRepository(infra.DB.Pool), termsSvc, infra.Tx)
	notificationsSvc := appnotifications.NewService(postgres.NewNotificationsRepository(infra.DB.Pool), postgres.NewAccountsRepository(infra.DB.Pool), infra.Tx, int64(infra.Config.NotificationDueSoonLeadSeconds))
	accountsSvc.SetObligationReader(postgres.NewBorrowingObligations(postgres.NewBorrowingRepository(infra.DB.Pool)))
	reportingSvc := appreporting.NewService(postgres.NewReportingRepository(infra.DB.Pool), postgres.NewAccountsRepository(infra.DB.Pool), infra.Tx)
	profileSvc := appprofile.NewService(postgres.NewProfileRepository(infra.DB.Pool), postgres.NewAccountsRepository(infra.DB.Pool), infra.Tx, profileimage.Validator{})
	v := validator.New()

	return &Container{
		ProfileSvc: profileSvc, Profile: handlers.NewProfileHandler(profileSvc, infra.Config.App.Env, infra.Config.Security),
		AuthSvc: authSvc, ReportingSvc: reportingSvc, Reporting: handlers.NewReportingHandler(reportingSvc),
		NotificationsSvc:     notificationsSvc,
		Notifications:        handlers.NewNotificationsHandler(notificationsSvc, infra.Config.App.Env, infra.Config.Security),
		BorrowingSvc:         borrowingSvc,
		Borrowing:            handlers.NewBorrowingHandler(borrowingSvc, infra.Config.App.Env, infra.Config.Security),
		InventorySvc:         inventorySvc,
		Inventory:            handlers.NewInventoryHandler(inventorySvc, infra.Config.App.Env, infra.Config.Security),
		AccountManagementSvc: accountsSvc,
		AccountManagement:    handlers.NewAccountsHandler(accountsSvc, spreadsheet.StudentRoster{}, infra.Config.App.Env, infra.Config.Security),
		UserSvc:              userSvc,
		TermsSvc:             termsSvc,
		Health:               handlers.NewHealthHandler(infra.Health),
		Auth:                 handlers.NewAuthHandler(authSvc, v, infra.Config.App.Env, infra.Config.Security),
		User:                 handlers.NewUserHandler(userSvc, v),
		Terms:                handlers.NewTermsHandler(termsSvc, infra.Config.App.Env, infra.Config.Security),
		TokenIssuer:          issuer,
		Accounts:             auth.NewAccountResolver(userRepo),
		Environment:          infra.Config.App.Env,
	}
}

// routeDeps converts the container into the bundle expected by routes.Register.
func (c *Container) routeDeps() *routes.Deps {
	return &routes.Deps{
		Health:            c.Health,
		Auth:              c.Auth,
		User:              c.User,
		Terms:             c.Terms,
		AccountManagement: c.AccountManagement,
		Inventory:         c.Inventory,
		Borrowing:         c.Borrowing,
		Notifications:     c.Notifications, Reporting: c.Reporting, Profile: c.Profile,
		TokenIssuer: c.TokenIssuer,
		Accounts:    c.Accounts,
		Environment: c.Environment,
	}
}
