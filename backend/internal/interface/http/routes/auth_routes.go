package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/gofiber/fiber/v3"
)

// RegisterAuth wires the auth routes under /api/v1/auth.
// Browser login issues a refresh cookie; refresh/logout read only that cookie.
// Cookie-changing handlers require a trusted Origin. Provisioning is never a
// public HTTP action, including development/test; fixtures use isolated SQL.
func RegisterAuth(v1 fiber.Router, h *handlers.AuthHandler, protected fiber.Handler, environment config.Environment) {
	g := v1.Group("/auth")
	g.Post("/login", h.Login)
	g.Post("/refresh", h.Refresh)
	g.Post("/logout", h.Logout)

	// Authenticated convenience endpoint.
	g.Get("/me", protected, h.Me)
}
