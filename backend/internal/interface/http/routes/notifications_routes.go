package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/gofiber/fiber/v3"
)

func RegisterNotifications(v1 fiber.Router, h *handlers.NotificationsHandler, protected fiber.Handler) {
	g := v1.Group("/notifications", func(c fiber.Ctx) error { c.Set("Cache-Control", "no-store"); return c.Next() }, protected)
	g.Get("/", h.List)
	g.Get("/count", h.Count)
	g.Post("/read-all", h.MarkAll)
	g.Patch("/:id", h.Mark)
}
