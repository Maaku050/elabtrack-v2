package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/gofiber/fiber/v3"
)

// RegisterAuth wires the auth routes under /api/v1/auth.
// All auth routes are public (no access token required) except /me.
func RegisterAuth(v1 fiber.Router, h *handlers.AuthHandler, issuer application.TokenIssuer) {
	g := v1.Group("/auth")
	g.Post("/register", h.Register)
	g.Post("/login", h.Login)
	g.Post("/refresh", h.Refresh)
	g.Post("/logout", h.Logout)

	// Authenticated convenience endpoint.
	g.Get("/me", middleware.Auth(issuer), h.Me)
}
