package bootstrap

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/logger"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
)

// newServer constructs a Fiber app with all global middleware attached.
// Route registration is performed by RegisterRoutes.
func newServer(cfg *config.Config, c *Container, log *logger.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:        cfg.App.Name,
		ReadTimeout:    cfg.App.ReadTimeout,
		WriteTimeout:   cfg.App.WriteTimeout,
		IdleTimeout:    cfg.App.IdleTimeout,
		ReadBufferSize: 8192, // Explicit header ceiling; read timeout covers headers and body.
		BodyLimit:      cfg.App.BodyLimit,
		ServerHeader:   "",
		ErrorHandler:   middleware.ErrorHandler(log, cfg.App.Env),
		TrustProxy:     false, // Only middleware.ClientInfo may honor forwarding headers.
	})

	// Global middleware (order matters). The outer recovery is a fallback
	// for correlation/IP/completion failures; inner recovery lets Logger observe
	// the final status of panics from the remaining pipeline.
	app.Use(middleware.Recovery(log))
	app.Use(middleware.RequestID())
	app.Use(middleware.ClientInfo(cfg.Security))
	app.Use(middleware.Logger(log))
	app.Use(middleware.Recovery(log))
	app.Use(middleware.SecurityHeaders(cfg.App.Env))
	app.Use(middleware.CORS(cfg.Security))
	app.Use(middleware.RateLimit(cfg.Security))
	app.Use(middleware.RequestSafety())
	app.Use(compress.New())
	// Fiber initializes/replaces its underlying server at listener startup.
	// Install the wrapper after that initialization, before serving connections.
	app.Hooks().OnListen(func(fiber.ListenData) error {
		server := app.Server()
		server.ErrorHandler = middleware.ParserErrorLogging(app, server.ErrorHandler, log)
		return nil
	})

	return app
}

// RegisterRoutes attaches the API routes to the app and returns it.
func RegisterRoutes(app *fiber.App, c *Container) *fiber.App {
	routes.Register(app, c.routeDeps())
	return app
}
