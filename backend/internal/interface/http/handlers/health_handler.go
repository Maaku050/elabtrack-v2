package handlers

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/gofiber/fiber/v3"
)

// HealthHandler exposes the health-check endpoint.
type HealthHandler struct {
	db *database.HealthChecker
}

// NewHealthHandler constructs a HealthHandler.
func NewHealthHandler(db *database.HealthChecker) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health returns the aggregated health of the API and its dependencies.
//
// GET /api/v1/health
func (h *HealthHandler) Health(c fiber.Ctx) error {
	services := map[string]string{
		"api": "healthy",
	}
	status := "ok"
	httpStatus := 200

	if h.db != nil {
		if err := h.db.Check(c.Context()); err != nil {
			services["database"] = "unhealthy"
			status = "degraded"
			httpStatus = 503
		} else {
			services["database"] = "healthy"
		}
	} else {
		services["database"] = "unknown"
	}

	return c.Status(httpStatus).JSON(fiber.Map{
		"status":   status,
		"services": services,
	})
}
