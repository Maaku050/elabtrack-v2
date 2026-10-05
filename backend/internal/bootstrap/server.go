package bootstrap

import (
	"github.com/fullstacktemplate/backend/internal/config"
	"github.com/fullstacktemplate/backend/internal/infrastructure/logger"
	"github.com/fullstacktemplate/backend/internal/interface/http/middleware"
	"github.com/fullstacktemplate/backend/internal/interface/http/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
)

// newServer constructs a Fiber app with all global middleware attached.
// Route registration is performed by RegisterRoutes.
func newServer(cfg *config.Config, c *Container, log *logger.Logger) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
		IdleTimeout:  cfg.App.ReadTimeout + cfg.App.WriteTimeout,
		BodyLimit:    bodyLimitBytes(cfg.App.BodyLimit),
		ServerHeader: "",
		ErrorHandler: middleware.ErrorHandler(log),
	})

	// Global middleware (order matters).
	app.Use(middleware.RequestID())
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.CORS(cfg.Security))
	app.Use(middleware.Recovery(log))
	app.Use(middleware.Logger())
	app.Use(middleware.RateLimit(cfg.Security))
	app.Use(compress.New())

	return app
}

// RegisterRoutes attaches the API routes to the app and returns it.
func RegisterRoutes(app *fiber.App, c *Container) *fiber.App {
	routes.Register(app, c.routeDeps())
	return app
}
