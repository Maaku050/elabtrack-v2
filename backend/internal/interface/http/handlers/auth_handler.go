package handlers

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	appuser "github.com/Maaku050/elabtrack-v2/backend/internal/application/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
)

// AuthHandler exposes the auth endpoints.
type AuthHandler struct {
	auth *auth.Service
	user *appuser.Service
	v    *validator.Validator
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(authSvc *auth.Service, userSvc *appuser.Service, v *validator.Validator) *AuthHandler {
	return &AuthHandler{auth: authSvc, user: userSvc, v: v}
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(c fiber.Ctx) error {
	var req auth.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, "Invalid request body.")
	}
	if fe := h.v.Validate(req); fe != nil {
		return response.Error(c, fe)
	}
	pair, err := h.auth.Register(c.Context(), req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Created(c, "User registered successfully.", pair)
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req auth.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, "Invalid request body.")
	}
	if fe := h.v.Validate(req); fe != nil {
		return response.Error(c, fe)
	}
	pair, err := h.auth.Login(c.Context(), req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, "Login successful.", pair)
}

// Refresh handles POST /api/v1/auth/refresh.
func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	var req auth.RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, "Invalid request body.")
	}
	if fe := h.v.Validate(req); fe != nil {
		return response.Error(c, fe)
	}
	pair, err := h.auth.Refresh(c.Context(), req)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, "Token refreshed successfully.", pair)
}

// Logout handles POST /api/v1/auth/logout.
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	var req auth.RefreshRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, "Invalid request body.")
	}
	if fe := h.v.Validate(req); fe != nil {
		return response.Error(c, fe)
	}
	if err := h.auth.Logout(c.Context(), req); err != nil {
		return response.Error(c, err)
	}
	return response.NoContent(c)
}

// Me handles GET /api/v1/auth/me (convenience alias for the current user).
func (h *AuthHandler) Me(c fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return response.Unauthorized(c, "Authentication required.")
	}
	dto, err := h.auth.CurrentUser(c.Context(), claims)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, "Current user retrieved successfully.", dto)
}
