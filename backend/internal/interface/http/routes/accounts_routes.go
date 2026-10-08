package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterAccounts(v1 fiber.Router, h *handlers.AccountsHandler, protected fiber.Handler) {
	noStore := func(c fiber.Ctx) error { c.Set("Cache-Control", "no-store"); return c.Next() }
	admin := middleware.RequirePermission(user.ManageAccounts)
	v1.Post("/auth/activate", noStore, h.Activate)
	b := v1.Group("/borrowers", noStore, protected, middleware.RequirePermission(user.ReadBorrowers))
	b.Get("/", h.List)
	b.Get("/policy", admin, h.Policy)
	b.Get("/student-template", admin, h.Template)
	b.Post("/rosters", admin, h.Prepare)
	b.Get("/rosters/:id", admin, h.Preview)
	b.Post("/rosters/:id/confirm", admin, h.Confirm)
	b.Post("/", admin, h.Create)
	b.Get("/:id", h.Detail)
	b.Patch("/:id", admin, h.Profile)
	b.Patch("/:id/status", admin, h.Status)
	b.Post("/:id/activation", admin, h.Resend)
	b.Get("/:id/audit", admin, h.Audits)
	s := v1.Group("/staff-accounts", noStore, protected, admin)
	s.Get("/", h.List)
	s.Post("/", h.Create)
	s.Get("/:id", h.Detail)
	s.Patch("/:id/status", h.Status)
	s.Post("/:id/activation", h.Resend)
	s.Get("/:id/audit", h.Audits)
}
