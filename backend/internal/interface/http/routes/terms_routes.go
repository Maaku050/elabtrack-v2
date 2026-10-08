package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
)

func RegisterTerms(v1 fiber.Router, h *handlers.TermsHandler, protected fiber.Handler) {
	g := v1.Group("/terms", func(c fiber.Ctx) error {
		c.Set("Cache-Control", "no-store")
		if c.Method() == fiber.MethodPost {
			limit := middleware.AuthBodyLimit
			if c.Path() == "/api/v1/terms/versions" {
				limit = 96 << 10
			}
			if len(c.Body()) > limit {
				return response.Fail(c, 413, "Request body is too large.", "PAYLOAD_TOO_LARGE")
			}
		}
		return c.Next()
	}, protected)
	g.Get("/current", h.Current)
	g.Get("/status", middleware.RequirePermission(user.AcceptBorrowerTerms), h.Status)
	g.Post("/versions", middleware.RequirePermission(user.PublishTerms), h.Publish)
	g.Post("/:versionID/accept", middleware.RequirePermission(user.AcceptBorrowerTerms), h.Accept)
}
