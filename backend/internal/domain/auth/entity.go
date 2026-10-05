package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// RefreshToken is the persisted refresh-token aggregate.
// Tokens are opaque (random bytes) and rotated on each refresh.
type RefreshToken struct {
	ID        uuid.UUID
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewRefreshToken constructs a new refresh-token record.
func NewRefreshToken(token string, userID uuid.UUID, expiresAt time.Time) *RefreshToken {
	now := time.Now().UTC()
	return &RefreshToken{
		ID:        uuid.New(),
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
		Revoked:   false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsExpired reports whether the token has passed its expiry.
func (r *RefreshToken) IsExpired() bool {
	return time.Now().UTC().After(r.ExpiresAt)
}

// IsRevoked reports whether the token has been revoked.
func (r *RefreshToken) IsRevoked() bool { return r.Revoked }

// IsValid reports whether the token is usable (not revoked, not expired).
func (r *RefreshToken) IsValid() bool { return !r.Revoked && !r.IsExpired() }

// Revoke marks the token as revoked.
func (r *RefreshToken) Revoke() {
	r.Revoked = true
	r.UpdatedAt = time.Now().UTC()
}

// Domain errors.
var (
	ErrTokenNotFound = errors.New("refresh token not found")
	ErrTokenExpired  = errors.New("refresh token expired")
	ErrTokenRevoked  = errors.New("refresh token revoked")
	ErrTokenInvalid  = errors.New("refresh token invalid")
)
