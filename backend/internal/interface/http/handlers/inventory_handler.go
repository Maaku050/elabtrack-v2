package handlers

import (
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"io"
	"strconv"
)

type InventoryHandler struct {
	svc     *app.Service
	browser browserSessionPolicy
}

func NewInventoryHandler(s *app.Service, e config.Environment, p config.SecurityConfig) *InventoryHandler {
	return &InventoryHandler{s, newBrowserSessionPolicy(e, p)}
}
func (h *InventoryHandler) mutation(c fiber.Ctx) (string, error) {
	_, ok := h.browser.origins[c.Get("Origin")]
	if !h.browser.enabled || !ok || (c.Method() != fiber.MethodPost && c.Method() != fiber.MethodPatch) {
		return "", shared.ErrForbidden
	}
	return commandKey(c)
}
func inventoryResult(c fiber.Ctx, v any, e error) error {
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Inventory request completed.", v)
}
func (h *InventoryHandler) List(c fiber.Ctx) error {
	p, e := strconv.Atoi(c.Query("page", "1"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	n, e := strconv.Atoi(c.Query("per_page", "25"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	var cat *uuid.UUID
	if raw := c.Query("category_id"); raw != "" {
		id, e := uuid.Parse(raw)
		if e != nil {
			return response.Error(c, shared.ErrInvalidInput)
		}
		cat = &id
	}
	avail := c.Query("available_only", "false")
	if avail != "true" && avail != "false" {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.List(c.Context(), actor(c), d.Filter{Page: p, PerPage: n, Search: c.Query("search"), Status: c.Query("status"), Sort: c.Query("sort"), CategoryID: cat, AvailableOnly: avail == "true"})
	return inventoryResult(c, v, e)
}
func (h *InventoryHandler) Detail(c fiber.Ctx) error {
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	v, e := h.svc.Detail(c.Context(), actor(c), id)
	return inventoryResult(c, v, e)
}
func (h *InventoryHandler) Create(c fiber.Ctx) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in d.Create
	if e = strictObject(c, &in); e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Create(c.Context(), actor(c), key, in)
	if e != nil {
		return response.Error(c, e)
	}
	return response.Created(c, "Equipment created.", v)
}
func (h *InventoryHandler) Edit(c fiber.Ctx) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in d.Metadata
	if e = strictObject(c, &in); e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Edit(c.Context(), actor(c), id, key, in)
	return inventoryResult(c, v, e)
}
func (h *InventoryHandler) Status(c fiber.Ctx) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in d.StatusInput
	if e = strictObject(c, &in); e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Status(c.Context(), actor(c), id, key, in)
	return inventoryResult(c, v, e)
}
func (h *InventoryHandler) Adjust(c fiber.Ctx) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in struct {
		Kind             string `json:"kind"`
		Quantity         *int64 `json:"quantity"`
		Reason           string `json:"reason"`
		ExpectedSequence *int64 `json:"expected_sequence"`
		Confirm          bool   `json:"confirm"`
	}
	if e = strictObject(c, &in); e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	if in.Quantity == nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Adjust(c.Context(), actor(c), id, key, d.Adjustment{Kind: in.Kind, Quantity: *in.Quantity, Reason: in.Reason, ExpectedSequence: in.ExpectedSequence, Confirm: in.Confirm})
	return inventoryResult(c, v, e)
}
func (h *InventoryHandler) Movements(c fiber.Ctx) error {
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	p, e := strconv.Atoi(c.Query("page", "1"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Movements(c.Context(), actor(c), id, p)
	return inventoryResult(c, v, e)
}
func (h *InventoryHandler) Categories(c fiber.Ctx) error {
	p, e := strconv.Atoi(c.Query("page", "1"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Categories(c.Context(), actor(c), p)
	return inventoryResult(c, v, e)
}
func (h *InventoryHandler) Category(c fiber.Ctx) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	id := uuid.Nil
	if c.Params("id") != "" {
		id, e = accountID(c)
		if e != nil {
			return response.Error(c, e)
		}
	}
	var in d.CategoryInput
	if e = strictObject(c, &in); e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Category(c.Context(), actor(c), id, key, in)
	if e == nil && id == uuid.Nil {
		return response.Created(c, "Equipment category created.", v)
	}
	return inventoryResult(c, v, e)
}
func (h *InventoryHandler) Image(c fiber.Ctx) error {
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	im, e := uuid.Parse(c.Params("image"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Image(c.Context(), actor(c), id, im)
	if e != nil {
		return response.Error(c, e)
	}
	c.Set("Content-Type", "image/png")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Content-Disposition", "inline; filename=equipment.png")
	return c.Send(v.PNG)
}
func (h *InventoryHandler) SaveImage(c fiber.Ctx) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	form, e := c.MultipartForm()
	if e != nil || len(form.File) != 1 || len(form.File["image"]) != 1 || len(form.Value) != 1 || len(form.Value["expected_version"]) != 1 {
		return response.Error(c, shared.ErrInvalidInput)
	}
	version, e := strconv.ParseInt(form.Value["expected_version"][0], 10, 64)
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	f := form.File["image"][0]
	if f.Size <= 0 || f.Size > d.MaxImageBytes {
		return response.Error(c, shared.ErrInvalidInput)
	}
	reader, e := f.Open()
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	defer reader.Close()
	raw, e := io.ReadAll(io.LimitReader(reader, d.MaxImageBytes+1))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.SaveImage(c.Context(), actor(c), id, key, version, raw)
	return inventoryResult(c, v, e)
}
