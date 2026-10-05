package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// RefreshToken is a hash-only session record. It never holds the bearer secret.
type RefreshToken struct {
	ID         uuid.UUID
	TokenHash  string
	UserID     uuid.UUID
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	ReplacedBy *uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewRefreshToken(hash string, userID uuid.UUID, createdAt, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{ID: uuid.New(), TokenHash: hash, UserID: userID,
		ExpiresAt: expiresAt, CreatedAt: createdAt, UpdatedAt: createdAt}
}

// ValidAt treats the exact expiry boundary as expired.
func (r *RefreshToken) ValidAt(now time.Time) bool {
	return r.RevokedAt == nil && now.Before(r.ExpiresAt)
}

var (
	ErrTokenNotFound = errors.New("refresh token not found")
	ErrTokenExpired  = errors.New("refresh token expired")
	ErrTokenRevoked  = errors.New("refresh token revoked")
	ErrTokenInvalid  = errors.New("refresh token invalid")
)
