package auth

import (
	"time"

	"github.com/google/uuid"
)

// TokenPairDTO is an internal result. Raw refresh material is never JSON.
type TokenPairDTO struct {
	AccessToken      string      `json:"access_token"`
	RefreshToken     string      `json:"-"`
	ExpiresAt        time.Time   `json:"expires_at"`
	TokenType        string      `json:"token_type"`
	RefreshExpiresAt time.Time   `json:"-"`
	User             AuthUserDTO `json:"user"`
}

// RegisterRequest is the input DTO for registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// LoginRequest is the input DTO for login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// RefreshRequest is an application credential command, never HTTP-bound.
type RefreshRequest struct {
	RefreshToken string `json:"-"`
}

// AuthUserDTO is the public representation of the authenticated user.
type AuthUserDTO struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active"`
}
