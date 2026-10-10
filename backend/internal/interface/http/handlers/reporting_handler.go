package handlers

import (
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/reporting"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/reporting"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type ReportingHandler struct{ svc *app.Service }

func NewReportingHandler(s *app.Service) *ReportingHandler { return &ReportingHandler{s} }
func (h *ReportingHandler) Dashboard(c fiber.Ctx) error {
	days, e := strconv.Atoi(c.Query("days", "7"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Dashboard(c.Context(), actor(c), days)
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Current authorized operational metrics.", v)
}
func (h *ReportingHandler) Definitions(c fiber.Ctx) error {
	v, e := h.svc.Definitions(c.Context(), actor(c))
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Authorized report definitions.", v)
}
func reportFilter(c fiber.Ctx) (f d.Filter, e error) {
	f.Page, e = strconv.Atoi(c.Query("page", "1"))
	if e != nil {
		return f, shared.ErrInvalidInput
	}
	f.PerPage, e = strconv.Atoi(c.Query("per_page", "25"))
	if e != nil {
		return f, shared.ErrInvalidInput
	}
	if raw := c.Query("due_today"); raw != "" {
		if raw != "true" && raw != "false" {
			return f, shared.ErrInvalidInput
		}
		f.DueToday = raw == "true"
	}
	f.Search = c.Query("search")
	f.Status = c.Query("status")
	for name, dst := range map[string]**time.Time{"from": &f.From, "to": &f.To} {
		if raw := c.Query(name); raw != "" {
			v, err := time.Parse(time.RFC3339Nano, raw)
			if err != nil {
				return f, shared.ErrInvalidInput
			}
			*dst = &v
		}
	}
	for name, dst := range map[string]**uuid.UUID{"equipment_id": &f.EquipmentID, "category_id": &f.CategoryID, "borrower_id": &f.BorrowerID} {
		if raw := c.Query(name); raw != "" {
			v, err := uuid.Parse(raw)
			if err != nil || v == uuid.Nil {
				return f, shared.ErrInvalidInput
			}
			*dst = &v
		}
	}
	if !f.Valid() {
		return f, shared.ErrInvalidInput
	}
	return f, nil
}
func (h *ReportingHandler) Report(c fiber.Ctx) error {
	f, e := reportFilter(c)
	if e != nil {
		return response.Error(c, e)
	}
	v, e := h.svc.Report(c.Context(), actor(c), c.Params("kind"), f)
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Filtered report snapshot.", v)
}
func (h *ReportingHandler) Export(c fiber.Ctx) error {
	f, e := reportFilter(c)
	if e != nil {
		return response.Error(c, e)
	}
	v, e := h.svc.Export(c.Context(), actor(c), c.Params("kind"), f)
	if e != nil {
		return response.Error(c, e)
	}
	c.Set("Content-Type", "text/csv; charset=utf-8")
	c.Set("Content-Disposition", `attachment; filename="elabtrack-`+c.Params("kind")+`.csv"`)
	return c.Send(v)
}

func (h *ReportingHandler) ExportData(c fiber.Ctx) error {
	f, e := reportFilter(c)
	if e != nil {
		return response.Error(c, e)
	}
	v, e := h.svc.ExportData(c.Context(), actor(c), c.Params("kind"), f)
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Complete bounded filtered export snapshot.", v)
}
