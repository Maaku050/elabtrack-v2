package handlers

import (
	"errors"

	"github.com/Maaku050/elabtrack-v2/backend/internal/application/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
)

// AuthHandler exposes the auth endpoints.
type AuthHandler struct {
	auth    *auth.Service
	v       *validator.Validator
	browser browserSessionPolicy
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(authSvc *auth.Service, v *validator.Validator, env config.Environment, security config.SecurityConfig) *AuthHandler {
	return &AuthHandler{auth: authSvc, v: v, browser: newBrowserSessionPolicy(env, security)}
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(c fiber.Ctx) error {
	if !h.beginBrowserSession(c) {
		return originDenied(c)
	}
	var req auth.RegisterRequest
	if err := bindStrictJSON(c, &req); err != nil {
		return response.BadRequest(c, "Invalid request body.")
	}
	if fe := h.v.Validate(req); fe != nil {
		return response.Error(c, fe)
	}
	pair, err := h.auth.Register(c.Context(), req)
	if err != nil {
		return response.Error(c, err)
	}
	h.browser.issue(c, pair)
	return response.Created(c, "User registered successfully.", browserSession(pair))
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(c fiber.Ctx) error {
	if !h.beginBrowserSession(c) {
		return originDenied(c)
	}
	var req auth.LoginRequest
	if err := bindStrictJSON(c, &req); err != nil {
		return response.BadRequest(c, "Invalid request body.")
	}
	if fe := h.v.Validate(req); fe != nil {
		return response.Error(c, fe)
	}
	pair, err := h.auth.Login(c.Context(), req)
	if err != nil {
		if errors.Is(err, domainuser.ErrUserInactive) {
			err = domainuser.ErrInvalidCredentials
		}
		return response.Error(c, err)
	}
	h.browser.issue(c, pair)
	return response.OK(c, "Login successful.", browserSession(pair))
}

// Refresh handles POST /api/v1/auth/refresh.
func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	if !h.beginBrowserSession(c) {
		return originDenied(c)
	}
	// No body credential fallback: this web flow reads only the HttpOnly cookie.
	req := auth.RefreshRequest{RefreshToken: c.Cookies(refreshCookieName)}
	pair, err := h.auth.Refresh(c.Context(), req)
	if err != nil {
		if errors.Is(err, domainauth.ErrTokenInvalid) {
			h.browser.clear(c)
		}
		return response.Error(c, err)
	}
	h.browser.issue(c, pair)
	return response.OK(c, "Token refreshed successfully.", browserSession(pair))
}

// Logout handles POST /api/v1/auth/logout.
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	if !h.beginBrowserSession(c) {
		return originDenied(c)
	}
	req := auth.RefreshRequest{RefreshToken: c.Cookies(refreshCookieName)}
	err := h.auth.Logout(c.Context(), req)
	h.browser.clear(c)
	if err != nil {
		return response.Error(c, err)
	}
	return response.NoContent(c)
}

// Me returns the already resolved current-account snapshot, never JWT metadata.
func (h *AuthHandler) Me(c fiber.Ctx) error {
	principal, ok := middleware.PrincipalFromContext(c)
	if !ok {
		return response.Unauthorized(c, "Authentication required.")
	}
	return response.OK(c, "Current user retrieved successfully.", principal.UserDTO())
}
