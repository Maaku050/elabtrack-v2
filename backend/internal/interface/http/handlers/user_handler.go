package handlers

import (
	appuser "github.com/Maaku050/elabtrack-v2/backend/internal/application/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/middleware"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/response"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/pagination"
	"github.com/Maaku050/elabtrack-v2/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v3"
)

// UserHandler exposes the user endpoints.
type UserHandler struct {
	svc *appuser.Service
	v   *validator.Validator
}

// NewUserHandler constructs a UserHandler.
func NewUserHandler(svc *appuser.Service, v *validator.Validator) *UserHandler {
	return &UserHandler{svc: svc, v: v}
}

// Me handles GET /api/v1/users/me.
func (h *UserHandler) Me(c fiber.Ctx) error {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		return response.Unauthorized(c, "Authentication required.")
	}
	dto, err := h.svc.GetByID(c.Context(), userID)
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, "User retrieved successfully.", dto)
}

// UpdateMe handles PATCH /api/v1/users/me.
func (h *UserHandler) UpdateMe(c fiber.Ctx) error {
	userID, ok := middleware.UserIDFromContext(c)
	if !ok {
		return response.Unauthorized(c, "Authentication required.")
	}

	var req appuser.UpdateProfileRequest
	if err := c.Bind().Body(&req); err != nil {
		return response.BadRequest(c, "Invalid request body.")
	}
	if fe := h.v.Validate(req); fe != nil {
		return response.Error(c, fe)
	}

	dto, err := h.svc.UpdateProfile(c.Context(), appuser.UpdateProfileCommand{
		UserID: userID,
		Name:   req.Name,
		Email:  req.Email,
	})
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, "Profile updated successfully.", dto)
}

// List handles GET /api/v1/users (admin-only).
func (h *UserHandler) List(c fiber.Ctx) error {
	page := pagination.FromContext(c)
	result, err := h.svc.List(c.Context(), appuser.ListUsersQuery{Page: page})
	if err != nil {
		return response.Error(c, err)
	}
	return response.Paginated(c, "Users retrieved successfully.", result.Items, result.Meta)
}
