package handlers

import (
	"bytes"
	"encoding/json"
	appterms "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type TermsHandler struct {
	svc     *appterms.Service
	browser browserSessionPolicy
}

func NewTermsHandler(svc *appterms.Service, env config.Environment, cfg config.SecurityConfig) *TermsHandler {
	return &TermsHandler{svc, newBrowserSessionPolicy(env, cfg)}
}
func (h *TermsHandler) Current(c fiber.Ctx) error {
	principal, ok := middleware.PrincipalFromContext(c)
	if !ok {
		return response.Unauthorized(c, "Authentication required.")
	}
	v, err := h.svc.Current(c.Context(), principal.ID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, "Current terms retrieved.", v)
}
func (h *TermsHandler) Status(c fiber.Ctx) error {
	principal, ok := middleware.PrincipalFromContext(c)
	if !ok {
		return response.Unauthorized(c, "Authentication required.")
	}
	status, err := h.svc.Status(c.Context(), principal.ID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, "Terms acceptance status retrieved.", status)
}
func (h *TermsHandler) Accept(c fiber.Ctx) error {
	if !h.browser.allowed(c) {
		return originDenied(c)
	}
	principal, ok := middleware.PrincipalFromContext(c)
	if !ok {
		return response.Unauthorized(c, "Authentication required.")
	}
	id, err := uuid.Parse(c.Params("versionID"))
	if err != nil || id == uuid.Nil {
		return response.BadRequest(c, "Invalid terms version.")
	}
	var input struct{}
	if body := bytes.TrimSpace(c.Body()); len(body) == 0 || body[0] != '{' {
		return response.BadRequest(c, "Invalid request body.")
	}
	if err := bindStrictJSON(c, &input); err != nil {
		return response.BadRequest(c, "Invalid request body.")
	}
	receipt, err := h.svc.Accept(c.Context(), principal.ID, id)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, "Terms acceptance recorded.", receipt)
}
func (h *TermsHandler) Publish(c fiber.Ctx) error {
	if !h.browser.allowed(c) {
		return originDenied(c)
	}
	principal, ok := middleware.PrincipalFromContext(c)
	if !ok {
		return response.Unauthorized(c, "Authentication required.")
	}
	var input struct {
		Version  string          `json:"version"`
		Title    string          `json:"title"`
		Body     string          `json:"body"`
		Expected json.RawMessage `json:"expected_current_version_id"`
	}
	if err := bindStrictJSON(c, &input); err != nil || len(input.Expected) == 0 {
		return response.BadRequest(c, "Invalid request body.")
	}
	var expected *uuid.UUID
	if !bytes.Equal(bytes.TrimSpace(input.Expected), []byte("null")) {
		var raw string
		if json.Unmarshal(input.Expected, &raw) != nil {
			return response.BadRequest(c, "Invalid current terms version.")
		}
		id, err := uuid.Parse(raw)
		if err != nil || id == uuid.Nil {
			return response.BadRequest(c, "Invalid current terms version.")
		}
		expected = &id
	}
	v, err := h.svc.Publish(c.Context(), appterms.PublishCommand{ActorID: principal.ID, Version: input.Version, Title: input.Title, Body: input.Body, ExpectedCurrentVersionID: expected})
	if err != nil {
		return response.Error(c, err)
	}
	// The durable row records actor/time/content hash; existing request logging
	// supplies safe request correlation without logging legal content or identity.
	return response.Created(c, "Terms version published.", v)
}
