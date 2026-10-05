package auth

import (
	"context"

	"github.com/google/uuid"
)

// Repository is the persistence port for refresh tokens.
type Repository interface {
	Create(ctx context.Context, t *RefreshToken) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
	Revoke(ctx context.Context, token string) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}
