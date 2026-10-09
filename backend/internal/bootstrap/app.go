package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appborrowing "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// App is the fully wired application. Run it with Run; shut it down with
// Shutdown (called automatically by Run on SIGINT/SIGTERM).
type App struct {
	cfg        *config.Config
	log        *logger.Logger
	db         *database.Postgres
	server     *fiber.App
	borrowing  *appborrowing.Service
	stopExpiry context.CancelFunc
	expiryDone chan struct{}
}

// New wires the entire application: config -> infrastructure -> services
// -> handlers -> routes. It returns an App ready to Run.
func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	infra, err := initInfrastructure(ctx, cfg)
	if err != nil {
		return nil, err
	}
	c := buildContainer(infra)
	server := newServer(cfg, c, infra.Logger)
	RegisterRoutes(server, c)

	return &App{
		cfg:       cfg,
		log:       infra.Logger,
		db:        infra.DB,
		server:    server,
		borrowing: c.BorrowingSvc,
	}, nil
}

// Run starts the HTTP server and blocks until SIGINT/SIGTERM is received,
// then performs a graceful shutdown.
func (a *App) Run() error {
	workerCtx, cancel := context.WithCancel(context.Background())
	a.stopExpiry = cancel
	a.expiryDone = make(chan struct{})
	go a.runExpiry(workerCtx)
	addr := ":" + a.cfg.App.Port
	go func() {
		logStartup(a.log, a.cfg)
		if err := a.server.Listen(addr, fiber.ListenConfig{DisableStartupMessage: true}); err != nil && !errors.Is(err, http.ErrServerClosed) {
			class, _ := shared.FailureDetails(err)
			a.log.Fatal("server stopped unexpectedly", zap.String("event", "server.listen_failed"), zap.String("error_class", class))
		}
	}()

	// Wait for interrupt signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	signal.Stop(quit)
	a.log.Info("shutdown signal received", zap.String("event", "server.shutdown_started"))

	return a.Shutdown(context.Background())
}

// Shutdown stops the server and closes the database pool with a timeout.
func (a *App) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := a.server.ShutdownWithContext(shutdownCtx); err != nil {
		class, _ := shared.FailureDetails(err)
		a.log.Error("server shutdown error", zap.String("event", "server.shutdown_failed"), zap.String("error_class", class))
	}
	if a.stopExpiry != nil {
		a.stopExpiry()
		<-a.expiryDone
	}
	a.db.Close()
	a.log.Info("database pool closed", zap.String("event", "server.shutdown_completed"))
	a.log.Sync()
	return nil
}

// Logger exposes the application logger.
func (a *App) Logger() *logger.Logger { return a.log }

// Config exposes the application config.
func (a *App) Config() *config.Config { return a.cfg }

func logStartup(log *logger.Logger, cfg *config.Config) {
	log.Info("starting http server", zap.String("event", "server.starting"), zap.String("addr", ":"+cfg.App.Port), zap.String("env", string(cfg.App.Env)), zap.String("migration_policy", "explicit_command"), zap.Bool("startup_seed_enabled", false), zap.Bool("trusted_proxy_enabled", len(cfg.Security.TrustedProxies) > 0), zap.Int("allowed_origin_count", len(cfg.Security.AllowedOrigins)))
}
