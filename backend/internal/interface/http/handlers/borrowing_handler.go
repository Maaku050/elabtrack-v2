package handlers

import (
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"strconv"
)

type BorrowingHandler struct {
	svc     *app.Service
	browser browserSessionPolicy
}

func NewBorrowingHandler(s *app.Service, e config.Environment, p config.SecurityConfig) *BorrowingHandler {
	return &BorrowingHandler{s, newBrowserSessionPolicy(e, p)}
}
func (h *BorrowingHandler) mutation(c fiber.Ctx) (string, error) {
	if !h.browser.enabled {
		return "", shared.ErrForbidden
	}
	if _, ok := h.browser.origins[c.Get("Origin")]; !ok {
		return "", shared.ErrForbidden
	}
	return commandKey(c)
}
func borrowingResult(c fiber.Ctx, v any, e error) error {
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Borrowing request completed.", v)
}
func (h *BorrowingHandler) List(c fiber.Ctx) error {
	p, e := strconv.Atoi(c.Query("page", "1"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	n, e := strconv.Atoi(c.Query("per_page", "25"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	var owner *uuid.UUID
	if raw := c.Query("borrower_id"); raw != "" {
		id, e := uuid.Parse(raw)
		if e != nil || id == uuid.Nil {
			return response.Error(c, shared.ErrInvalidInput)
		}
		owner = &id
	}
	v, e := h.svc.List(c.Context(), actor(c), d.Filter{Owner: owner, Page: p, PerPage: n, Status: c.Query("status"), Search: c.Query("search")})
	return borrowingResult(c, v, e)
}
func (h *BorrowingHandler) Detail(c fiber.Ctx) error {
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	v, e := h.svc.Read(c.Context(), actor(c), id)
	return borrowingResult(c, v, e)
}
func (h *BorrowingHandler) Create(c fiber.Ctx) error { return h.create(c, false) }
func (h *BorrowingHandler) Direct(c fiber.Ctx) error { return h.create(c, true) }
func (h *BorrowingHandler) create(c fiber.Ctx, direct bool) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in d.Input
	if e = strictObject(c, &in); e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	var v d.Record
	if direct {
		v, e = h.svc.Direct(c.Context(), actor(c), key, in)
	} else {
		v, e = h.svc.Submit(c.Context(), actor(c), key, in)
	}
	if e != nil {
		return response.Error(c, e)
	}
	return response.Created(c, "Borrowing recorded.", v)
}
func (h *BorrowingHandler) Approve(c fiber.Ctx) error { return h.decide(c, "CHECKOUT") }
func (h *BorrowingHandler) Deny(c fiber.Ctx) error    { return h.decide(c, "DENIED") }
func (h *BorrowingHandler) Cancel(c fiber.Ctx) error  { return h.decide(c, "CANCELLED") }
func (h *BorrowingHandler) decide(c fiber.Ctx, action string) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in d.Decision
	if e = strictObject(c, &in); e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Decide(c.Context(), actor(c), id, key, action, in)
	return borrowingResult(c, v, e)
}
