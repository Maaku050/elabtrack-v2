package routes

import (
	"github.com/fullstacktemplate/backend/internal/application"
	"github.com/fullstacktemplate/backend/internal/interface/http/handlers"
	"github.com/fullstacktemplate/backend/internal/interface/http/middleware"
	"github.com/gofiber/fiber/v3"
)

// RegisterUser wires the user routes under /api/v1/users.
// /me endpoints require authentication; /users (list) requires admin role.
func RegisterUser(v1 fiber.Router, h *handlers.UserHandler, issuer application.TokenIssuer) {
	g := v1.Group("/users")

	// Authenticated routes for the current user.
	g.Get("/me", middleware.Auth(issuer), h.Me)
	g.Patch("/me", middleware.Auth(issuer), h.UpdateMe)

	// Admin-only routes. The RequireRole middleware is a foundation for
	// role/permission-based authorization; extend it as needed.
	g.Get("/", middleware.Auth(issuer), middleware.RequireRole("admin"), h.List)
}
