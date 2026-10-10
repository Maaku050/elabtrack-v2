package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/gofiber/fiber/v3"
)

func RegisterProfile(v1 fiber.Router, h *handlers.ProfileHandler, protected fiber.Handler) {
	g := v1.Group("/profile-images", protected, func(c fiber.Ctx) error { c.Set("Cache-Control", "no-store"); return c.Next() })
	g.Get("/account/me", h.Own)
	g.Get("/:id", h.Metadata)
	g.Get("/:id/:image", h.Image)
	g.Post("/:id", h.Save)
	g.Patch("/:id", h.Remove)
}
