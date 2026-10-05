package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository is the persistence port for refresh tokens.
type Repository interface {
	Create(ctx context.Context, t *RefreshToken) error
	// FindByHashForUpdate requires a transaction and locks the presented row.
	FindByHashForUpdate(ctx context.Context, hash string) (*RefreshToken, error)
	// Consume requires a transaction and succeeds only for a still-valid row.
	Consume(ctx context.Context, hash string, replacement uuid.UUID, now time.Time) error
	// Revoke is idempotent, including unknown hashes.
	Revoke(ctx context.Context, hash string) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
	// Cleanup deletes at most limit rows terminal before cutoff.
	Cleanup(ctx context.Context, cutoff time.Time, limit int) (int64, error)
}
