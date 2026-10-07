package handlers

import (
	"context"
	"errors"

	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
)

type HealthChecker interface{ Check(context.Context) error }
type HealthHandler struct{ db HealthChecker }

func NewHealthHandler(db HealthChecker) *HealthHandler { return &HealthHandler{db: db} }

type HealthDTO struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// Health is process liveness only. It does not touch dependencies or claim DB readiness.
func (h *HealthHandler) Health(c fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")
	return response.OK(c, "Service is alive.", HealthDTO{Status: "ok", Service: "elabtrack-v2"})
}

// Ready uses the existing two-second DB checker; nil dependencies fail closed.
func (h *HealthHandler) Ready(c fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")
	if h.db == nil {
		return response.Unavailable(c, errors.New("readiness dependency absent"))
	}
	if err := h.db.Check(c.Context()); err != nil {
		return response.Unavailable(c, err)
	}
	return response.OK(c, "Service is ready.", HealthDTO{Status: "ready", Service: "elabtrack-v2"})
}
