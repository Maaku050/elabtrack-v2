package user

import (
	"time"

	"github.com/google/uuid"
)

// DTOs are the wire-level data transfer objects for the user application
// layer. They are intentionally separate from domain entities so the
// domain never leaks persistence/transport concerns.

// UserDTO is the public representation of a user returned to clients.
type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateProfileRequest is the input DTO for PATCH /users/me.
// All fields are pointers so partial updates are supported.
type UpdateProfileRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// ListUsersRequest is the input DTO for GET /users.
type ListUsersRequest struct {
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
	Search  string `json:"search,omitempty"`
	Sort    string `json:"sort,omitempty"`
	Order   string `json:"order,omitempty"`
}
