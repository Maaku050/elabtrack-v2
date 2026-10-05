package postgres

import (
	"context"
	"errors"
	"time"

	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthRepository accepts only digests. Database detail is deliberately dropped
// because constraint errors can include the digest or other credential data.
type AuthRepository struct{ pool *pgxpool.Pool }

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository { return &AuthRepository{pool: pool} }
func (r *AuthRepository) executor(ctx context.Context) executor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}
	return r.pool
}

const refreshTokenColumns = `id, token_hash, user_id, expires_at, revoked_at, replaced_by, created_at, updated_at`

func (r *AuthRepository) Create(ctx context.Context, t *domainauth.RefreshToken) error {
	const q = `INSERT INTO refresh_tokens (id, token_hash, user_id, expires_at, revoked_at, replaced_by, created_at, updated_at)
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.executor(ctx).Exec(ctx, q, t.ID, t.TokenHash, t.UserID, t.ExpiresAt, t.RevokedAt, t.ReplacedBy, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return shared.ErrInternal
	}
	return nil
}

// A pool/autocommit SELECT would release the lock too early; reject it.
func (r *AuthRepository) FindByHashForUpdate(ctx context.Context, hash string) (*domainauth.RefreshToken, error) {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return nil, shared.ErrInternal
	}
	const q = `SELECT ` + refreshTokenColumns + ` FROM refresh_tokens WHERE token_hash = $1 FOR UPDATE`
	t := &domainauth.RefreshToken{}
	err := tx.QueryRow(ctx, q, hash).Scan(&t.ID, &t.TokenHash, &t.UserID, &t.ExpiresAt, &t.RevokedAt, &t.ReplacedBy, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainauth.ErrTokenNotFound
	}
	if err != nil {
		return nil, shared.ErrInternal
	}
	return t, nil
}

func (r *AuthRepository) Consume(ctx context.Context, hash string, replacement uuid.UUID, now time.Time) error {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return shared.ErrInternal
	}
	const q = `UPDATE refresh_tokens SET revoked_at = $3, replaced_by = $2, updated_at = $3
  WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > $3 AND expires_at > clock_timestamp()`
	tag, err := tx.Exec(ctx, q, hash, replacement, now)
	if err != nil {
		return shared.ErrInternal
	}
	if tag.RowsAffected() != 1 {
		return domainauth.ErrTokenInvalid
	}
	return nil
}

func (r *AuthRepository) Revoke(ctx context.Context, hash string) error {
	const q = `UPDATE refresh_tokens SET revoked_at = clock_timestamp(), updated_at = clock_timestamp()
  WHERE token_hash = $1 AND revoked_at IS NULL`
	_, err := r.executor(ctx).Exec(ctx, q, hash)
	if err != nil {
		return shared.ErrInternal
	}
	return nil
}

// This internal primitive affects sessions visible to the statement; it does
// not define a family-wide policy or serialize future logins/rotations.
func (r *AuthRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	const q = `UPDATE refresh_tokens SET revoked_at = clock_timestamp(), updated_at = clock_timestamp()
  WHERE user_id = $1 AND revoked_at IS NULL`
	_, err := r.executor(ctx).Exec(ctx, q, userID)
	if err != nil {
		return shared.ErrInternal
	}
	return nil
}

// Cleanup is one bounded statement; locked rows are left for the next pass.
func (r *AuthRepository) Cleanup(ctx context.Context, cutoff time.Time, limit int) (int64, error) {
	if limit < 1 || limit > 1000 {
		return 0, shared.ErrInvalidInput
	}
	const q = `WITH candidates AS (
  SELECT id FROM refresh_tokens
  WHERE LEAST(expires_at, COALESCE(revoked_at, expires_at)) < $1
  ORDER BY LEAST(expires_at, COALESCE(revoked_at, expires_at)), id
  LIMIT $2 FOR UPDATE SKIP LOCKED
 ) DELETE FROM refresh_tokens AS sessions USING candidates WHERE sessions.id = candidates.id`
	tag, err := r.executor(ctx).Exec(ctx, q, cutoff, limit)
	if err != nil {
		return 0, shared.ErrInternal
	}
	return tag.RowsAffected(), nil
}
