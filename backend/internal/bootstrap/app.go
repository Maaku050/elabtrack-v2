package bootstrap

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fullstacktemplate/backend/internal/config"
	"github.com/fullstacktemplate/backend/internal/infrastructure/database"
	"github.com/fullstacktemplate/backend/internal/infrastructure/logger"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

// App is the fully wired application. Run it with Run; shut it down with
// Shutdown (called automatically by Run on SIGINT/SIGTERM).
type App struct {
	cfg      *config.Config
	log      *logger.Logger
	db       *database.Postgres
	server   *fiber.App
	migrator *database.Migrator
	seeder   *database.Seeder
}

// New wires the entire application: config -> infrastructure -> services
// -> handlers -> routes. It returns an App ready to Run.
func New(ctx context.Context) (*App, error) {
	cfg := config.Load()
	infra, err := initInfrastructure(ctx, cfg)
	if err != nil {
		return nil, err
	}
	c := buildContainer(infra)
	server := newServer(cfg, c, infra.Logger)
	RegisterRoutes(server, c)

	return &App{
		cfg:      cfg,
		log:      infra.Logger,
		db:       infra.DB,
		server:   server,
		migrator: infra.Migrator,
		seeder:   infra.Seeder,
	}, nil
}

// Run starts the HTTP server and blocks until SIGINT/SIGTERM is received,
// then performs a graceful shutdown.
func (a *App) Run() error {
	addr := ":" + a.cfg.App.Port
	go func() {
		a.log.Info("starting http server", zap.String("addr", addr), zap.String("env", a.cfg.App.Env))
		if err := a.server.Listen(addr, fiber.ListenConfig{DisableStartupMessage: false}); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.log.Fatal("server stopped unexpectedly")
		}
	}()

	// Wait for interrupt signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.log.Info("shutdown signal received")

	return a.Shutdown(context.Background())
}

// Shutdown stops the server and closes the database pool with a timeout.
func (a *App) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := a.server.ShutdownWithContext(shutdownCtx); err != nil {
		a.log.Error("server shutdown error", zap.Error(err))
	}
	a.db.Close()
	a.log.Info("database pool closed")
	a.log.Sync()
	return nil
}

// Migrator exposes the migrator for CLI subcommands.
func (a *App) Migrator() *database.Migrator { return a.migrator }

// Seeder exposes the seeder for CLI subcommands.
func (a *App) Seeder() *database.Seeder { return a.seeder }

// Logger exposes the application logger.
func (a *App) Logger() *logger.Logger { return a.log }

// Config exposes the application config.
func (a *App) Config() *config.Config { return a.cfg }
