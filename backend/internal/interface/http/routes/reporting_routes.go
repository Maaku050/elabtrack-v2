package routes

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterReporting(v1 fiber.Router, h *handlers.ReportingHandler, protected fiber.Handler) {
	g := v1.Group("/reporting", func(c fiber.Ctx) error { c.Set("Cache-Control", "no-store"); return c.Next() }, protected)
	g.Get("/dashboard", h.Dashboard)
	staff := middleware.RequirePermission(user.StaffWorkspace)
	g.Get("/kinds", staff, h.Definitions)
	g.Get("/:kind.csv", staff, h.Export)
	g.Get("/:kind/export-data", staff, h.ExportData)
	g.Get("/:kind", staff, h.Report)
}
