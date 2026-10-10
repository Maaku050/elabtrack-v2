package handlers

import (
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/notifications"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/notifications"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
	"strconv"
)

type NotificationsHandler struct {
	svc     *app.Service
	browser browserSessionPolicy
}

func NewNotificationsHandler(s *app.Service, e config.Environment, p config.SecurityConfig) *NotificationsHandler {
	return &NotificationsHandler{s, newBrowserSessionPolicy(e, p)}
}
func (h *NotificationsHandler) List(c fiber.Ctx) error {
	page, e := strconv.Atoi(c.Query("page", "1"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	per, e := strconv.Atoi(c.Query("per_page", "25"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.List(c.Context(), actor(c), d.Filter{Page: page, PerPage: per, Read: c.Query("read")})
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Notifications retrieved.", v)
}
func (h *NotificationsHandler) Count(c fiber.Ctx) error {
	n, e := h.svc.Count(c.Context(), actor(c))
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Unread notification count.", map[string]int{"unread": n})
}
func (h *NotificationsHandler) Mark(c fiber.Ctx) error {
	if !h.browser.enabled {
		return response.Error(c, shared.ErrForbidden)
	}
	if _, ok := h.browser.origins[c.Get("Origin")]; !ok {
		return response.Error(c, shared.ErrForbidden)
	}
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in struct {
		Read *bool `json:"read"`
	}
	if strictObject(c, &in) != nil || in.Read == nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Mark(c.Context(), actor(c), id, *in.Read)
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Notification read state updated.", v)
}

func (h *NotificationsHandler) MarkAll(c fiber.Ctx) error {
	if !h.browser.enabled {
		return response.Error(c, shared.ErrForbidden)
	}
	if _, ok := h.browser.origins[c.Get("Origin")]; !ok {
		return response.Error(c, shared.ErrForbidden)
	}
	var in struct {
		Confirm bool `json:"confirm"`
	}
	if strictObject(c, &in) != nil || !in.Confirm {
		return response.Error(c, shared.ErrInvalidInput)
	}
	n, e := h.svc.MarkAll(c.Context(), actor(c))
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Authorized unread notifications marked read.", map[string]int64{"updated": n})
}
