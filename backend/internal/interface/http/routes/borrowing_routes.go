package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterBorrowing(v1 fiber.Router, h *handlers.BorrowingHandler, protected fiber.Handler) {
	g := v1.Group("/borrowings", func(c fiber.Ctx) error { c.Set("Cache-Control", "no-store"); return c.Next() }, protected)
	staff := middleware.RequirePermission(user.ManageBorrowings)
	borrower := middleware.RequirePermission(user.SubmitBorrowing)
	g.Get("/", h.List)
	g.Post("/", borrower, h.Create)
	g.Post("/direct-checkout", staff, h.Direct)
	g.Get("/:id", h.Detail)
	g.Post("/:id/cancel", borrower, h.Cancel)
	g.Post("/:id/deny", staff, h.Deny)
	g.Post("/:id/approve", staff, h.Approve)
}
