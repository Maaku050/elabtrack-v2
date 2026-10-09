package handlers

import (
	"bytes"
	"errors"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/accounts"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"

	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"io"
	"strconv"
	"strings"
	"time"
)

type AccountsHandler struct {
	svc     *app.Service
	parser  d.Parser
	browser browserSessionPolicy
}

func NewAccountsHandler(s *app.Service, p d.Parser, env config.Environment, security config.SecurityConfig) *AccountsHandler {
	return &AccountsHandler{s, p, newBrowserSessionPolicy(env, security)}
}
func actor(c fiber.Ctx) uuid.UUID { p, _ := middleware.PrincipalFromContext(c); return p.ID }
func accountID(c fiber.Ctx) (uuid.UUID, error) {
	id, e := uuid.Parse(c.Params("id"))
	if e != nil || id == uuid.Nil {
		return id, shared.ErrInvalidInput
	}
	return id, nil
}
func commandKey(c fiber.Ctx) (string, error) {
	key := c.Get("Idempotency-Key")
	id, e := uuid.Parse(key)
	if e != nil || id == uuid.Nil {
		return "", shared.ErrInvalidInput
	}
	return id.String(), nil
}
func strictObject(c fiber.Ctx, target any) error {
	b := bytes.TrimSpace(c.Body())
	if len(b) == 0 || len(b) > 16<<10 || b[0] != '{' {
		return shared.ErrInvalidInput
	}
	return bindStrictJSON(c, target)
}
func staffPath(c fiber.Ctx) bool {
	return strings.Contains(strings.ToLower(c.Path()), "/staff-accounts")
}
func accountResult(c fiber.Ctx, value any, e error) error {
	if e != nil {
		return response.Error(c, e)
	}
	return response.OK(c, "Account management request completed.", value)
}
func (h *AccountsHandler) List(c fiber.Ctx) error {
	page, e := strconv.Atoi(c.Query("page", "1"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	per, e := strconv.Atoi(c.Query("per_page", "25"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.List(c.Context(), actor(c), d.Filter{Page: page, PerPage: per, Search: c.Query("search"), BorrowerType: c.Query("borrower_type"), Status: c.Query("status"), Staff: staffPath(c)})
	return accountResult(c, v, e)
}
func (h *AccountsHandler) Detail(c fiber.Ctx) error {
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	v, e := h.svc.Detail(c.Context(), actor(c), id, staffPath(c))
	return accountResult(c, v, e)
}
func (h *AccountsHandler) mutation(c fiber.Ctx) (string, error) {
	_, trusted := h.browser.origins[c.Get("Origin")]
	if !h.browser.enabled || !trusted || (c.Method() != fiber.MethodPost && c.Method() != fiber.MethodPatch) {
		return "", shared.ErrForbidden
	}
	return commandKey(c)
}
func (h *AccountsHandler) Create(c fiber.Ctx) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in d.Input
	if e = strictObject(c, &in); e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Create(c.Context(), actor(c), key, in, staffPath(c))
	if errors.Is(e, user.ErrEmailAlreadyExists) {
		return response.Fail(c, 409, "This email is already assigned to an account.", "CONFLICT")
	}
	if e != nil {
		return response.Error(c, e)
	}
	return response.Created(c, "Account created; activation remains required.", v)
}
func (h *AccountsHandler) Profile(c fiber.Ctx) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in app.ProfileInput
	if strictObject(c, &in) != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.UpdateProfile(c.Context(), actor(c), id, key, in)
	return accountResult(c, v, e)
}
func (h *AccountsHandler) Status(c fiber.Ctx) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	var body struct {
		Active            *bool     `json:"active"`
		Confirm           bool      `json:"confirm"`
		ExpectedUpdatedAt time.Time `json:"expected_updated_at"`
	}
	if strictObject(c, &body) != nil || body.Active == nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	in := app.StatusInput{Active: *body.Active, Confirm: body.Confirm, ExpectedUpdatedAt: body.ExpectedUpdatedAt}
	v, e := h.svc.Status(c.Context(), actor(c), id, key, in, staffPath(c))
	return accountResult(c, v, e)
}
func (h *AccountsHandler) Resend(c fiber.Ctx) error {
	key, e := h.mutation(c)
	if e != nil {
		return response.Error(c, e)
	}
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in struct{}
	if strictObject(c, &in) != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Resend(c.Context(), actor(c), id, key)
	return accountResult(c, v, e)
}
func (h *AccountsHandler) Audits(c fiber.Ctx) error {
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	page, e := strconv.Atoi(c.Query("page", "1"))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Audits(c.Context(), actor(c), id, page)
	return accountResult(c, v, e)
}
func (h *AccountsHandler) Activate(c fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")
	if !h.browser.allowed(c) {
		return originDenied(c)
	}
	var in struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if strictObject(c, &in) != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	return accountResult(c, nil, h.svc.Activate(c.Context(), in.Token, in.Password))
}
func (h *AccountsHandler) Policy(c fiber.Ctx) error {
	return response.OK(c, "Provisioning safeguards retrieved.", h.svc.ProvisioningPolicy())
}
func (h *AccountsHandler) Template(c fiber.Ctx) error {
	b, e := h.parser.Template()
	if e != nil {
		return response.Error(c, e)
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=elabtrack-students.xlsx")
	return c.Send(b)
}
func (h *AccountsHandler) Prepare(c fiber.Ctx) error {
	if !h.browser.allowed(c) {
		return originDenied(c)
	}
	form, e := c.MultipartForm()
	if e != nil || len(form.File) != 1 || len(form.Value) != 1 || len(form.File["roster"]) != 1 || len(form.Value["operation"]) != 1 {
		return response.Error(c, shared.ErrInvalidInput)
	}
	header := form.File["roster"][0]
	if header.Size > d.MaxRosterUpload || !strings.HasSuffix(strings.ToLower(header.Filename), ".xlsx") {
		return response.Error(c, shared.ErrInvalidInput)
	}
	f, e := header.Open()
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	defer f.Close()
	raw, e := io.ReadAll(io.LimitReader(f, d.MaxRosterUpload+1))
	if e != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	rows, e := h.parser.Parse(raw)
	if e != nil {
		return response.BadRequest(c, rosterErrorMessage(e))
	}
	v, e := h.svc.Prepare(c.Context(), actor(c), form.Value["operation"][0], rows)
	return accountResult(c, v, e)
}
func (h *AccountsHandler) Preview(c fiber.Ctx) error {
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	v, e := h.svc.Preview(c.Context(), actor(c), id)
	return accountResult(c, v, e)
}
func (h *AccountsHandler) Confirm(c fiber.Ctx) error {
	if !h.browser.allowed(c) {
		return originDenied(c)
	}
	id, e := accountID(c)
	if e != nil {
		return response.Error(c, e)
	}
	var in struct {
		Rows    []int `json:"rows"`
		Confirm bool  `json:"confirm"`
	}
	if strictObject(c, &in) != nil {
		return response.Error(c, shared.ErrInvalidInput)
	}
	v, e := h.svc.Confirm(c.Context(), actor(c), id, in.Rows, in.Confirm)
	return accountResult(c, v, e)
}
