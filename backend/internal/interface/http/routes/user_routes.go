package routes

import (
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/gofiber/fiber/v3"
)

// RegisterUser wires the user routes under /api/v1/users.
// /me endpoints require authentication; /users (list) requires admin role.
func RegisterUser(v1 fiber.Router, h *handlers.UserHandler, protected fiber.Handler) {
	g := v1.Group("/users")

	// Authenticated routes for the current user.
	g.Get("/me", protected, h.Me)
	g.Patch("/me", protected, h.UpdateMe)

	// Temporary generic admin read access only; no account administration.
	g.Get("/", protected, middleware.RequireRole(domainuser.RoleAdmin), h.List)
}
