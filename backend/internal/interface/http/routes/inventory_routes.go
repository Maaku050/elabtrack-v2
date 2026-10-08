package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterInventory(v1 fiber.Router, h *handlers.InventoryHandler, protected fiber.Handler) {
	noStore := func(c fiber.Ctx) error { c.Set("Cache-Control", "no-store"); return c.Next() }
	write := middleware.RequirePermission(user.ManageInventory)
	e := v1.Group("/equipment", noStore, protected)
	e.Get("/", h.List)
	e.Post("/", write, h.Create)
	e.Get("/:id", h.Detail)
	e.Patch("/:id", write, h.Edit)
	e.Patch("/:id/status", write, h.Status)
	e.Post("/:id/adjustments", write, h.Adjust)
	e.Get("/:id/movements", write, h.Movements)
	e.Post("/:id/image", write, h.SaveImage)
	e.Get("/:id/images/:image", h.Image)
	c := v1.Group("/equipment-categories", noStore, protected)
	c.Get("/", h.Categories)
	c.Post("/", write, h.Category)
	c.Patch("/:id", write, h.Category)
}
