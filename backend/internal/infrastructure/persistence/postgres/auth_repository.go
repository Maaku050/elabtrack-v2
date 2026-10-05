package postgres

import (
	"context"
	"errors"
	"fmt"

	domainauth "github.com/fullstacktemplate/backend/internal/domain/auth"
	"github.com/fullstacktemplate/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuthRepository implements domain/auth.Repository on top of PostgreSQL.
type AuthRepository struct {
	pool *pgxpool.Pool
}

// NewAuthRepository constructs an AuthRepository backed by the given pool.
func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

func (r *AuthRepository) executor(ctx context.Context) executor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}
	return r.pool
}

const refreshTokenColumns = `id, token, user_id, expires_at, revoked, created_at, updated_at`

// Create inserts a new refresh-token record.
func (r *AuthRepository) Create(ctx context.Context, t *domainauth.RefreshToken) error {
	const q = `
		INSERT INTO refresh_tokens (id, token, user_id, expires_at, revoked, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.executor(ctx).Exec(ctx, q,
		t.ID, t.Token, t.UserID, t.ExpiresAt, t.Revoked, t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	return nil
}

// FindByToken returns a refresh-token record by its opaque token string.
func (r *AuthRepository) FindByToken(ctx context.Context, token string) (*domainauth.RefreshToken, error) {
	const q = `SELECT ` + refreshTokenColumns + ` FROM refresh_tokens WHERE token = $1`
	row := r.executor(ctx).QueryRow(ctx, q, token)
	t, err := scanRefreshToken(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainauth.ErrTokenNotFound
		}
		return nil, err
	}
	return t, nil
}

// Revoke marks a single refresh token as revoked.
func (r *AuthRepository) Revoke(ctx context.Context, token string) error {
	const q = `UPDATE refresh_tokens SET revoked = TRUE, updated_at = NOW() WHERE token = $1`
	tag, err := r.executor(ctx).Exec(ctx, q, token)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainauth.ErrTokenNotFound
	}
	return nil
}

// RevokeAllForUser revokes every active refresh token for a user.
func (r *AuthRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	const q = `UPDATE refresh_tokens SET revoked = TRUE, updated_at = NOW() WHERE user_id = $1 AND revoked = FALSE`
	_, err := r.executor(ctx).Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("revoke all refresh tokens: %w", err)
	}
	return nil
}

func scanRefreshToken(s scanner) (*domainauth.RefreshToken, error) {
	t := &domainauth.RefreshToken{}
	err := s.Scan(
		&t.ID,
		&t.Token,
		&t.UserID,
		&t.ExpiresAt,
		&t.Revoked,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}
