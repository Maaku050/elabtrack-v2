package handlers

import (
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/profile"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/profile"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"io"
	"strconv"
)

type ProfileHandler struct {
	svc     *app.Service
	browser browserSessionPolicy
}

func NewProfileHandler(s *app.Service, e config.Environment, p config.SecurityConfig) *ProfileHandler {
	return &ProfileHandler{s, newBrowserSessionPolicy(e, p)}
}
func (h *ProfileHandler) Metadata(c fiber.Ctx) error {
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	v, e := h.svc.Metadata(c.Context(), actor(c), id)
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Profile image metadata retrieved.", v)
}
func (h *ProfileHandler) Image(c fiber.Ctx) error {
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	image, e := uuid.Parse(c.Params("image"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Image(c.Context(), actor(c), id, image)
	if e != nil {
		return response.Error(c, e)
	}
	c.Set("Content-Type", "image/png")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("Content-Disposition", "inline; filename=profile.png")
	return c.Send(v.PNG)
}
func (h *ProfileHandler) mutation(c fiber.Ctx) (uuid.UUID, error) {
	if !h.browser.enabled {
		return uuid.Nil, shared.ErrForbidden
	}
	if _, ok := h.browser.origins[c.Get("Origin")]; !ok {
		return uuid.Nil, shared.ErrForbidden
	}
	id, e := accountID(c)
	if e != nil {
		return id, e
	}
	if id != actor(c) {
		return id, shared.ErrForbidden
	}
	return id, nil
}
func (h *ProfileHandler) Save(c fiber.Ctx) error {
	id, e := h.mutation(c)
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
	if f.Size <= 0 || f.Size > d.MaxBytes {
		return response.Error(c, shared.ErrInvalidInput)
	}
	reader, e := f.Open()
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	defer reader.Close()
	raw, e := io.ReadAll(io.LimitReader(reader, d.MaxBytes+1))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Save(c.Context(), actor(c), id, c.Get("Idempotency-Key"), version, raw, false)
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Profile image updated.", v)
}
func (h *ProfileHandler) Remove(c fiber.Ctx) error {
	id, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in struct {
		Version *int64 `json:"expected_version"`
		Confirm bool   `json:"confirm"`
	}
	if strictObject(c, &in) != nil || in.Version == nil || !in.Confirm {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Save(c.Context(), actor(c), id, c.Get("Idempotency-Key"), *in.Version, nil, true)
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Profile image removed; initials remain available.", v)
}

func (h *ProfileHandler) Own(c fiber.Ctx) error {
	v, e := h.svc.Own(c.Context(), actor(c))
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Your read-only account information.", v)
}
