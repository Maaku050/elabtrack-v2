package postgres

import (
	"context"
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TermsRepository struct{ pool *pgxpool.Pool }

func NewTermsRepository(pool *pgxpool.Pool) *TermsRepository { return &TermsRepository{pool} }
func (r *TermsRepository) executor(ctx context.Context) executor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}
	return r.pool
}
func (r *TermsRepository) LockPublication(ctx context.Context, exclusive bool) (*terms.Version, error) {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return nil, shared.ErrInternal
	}
	q := `SELECT current_version_id FROM terms_publication WHERE id=1 FOR SHARE`
	if exclusive {
		q = `SELECT current_version_id FROM terms_publication WHERE id=1 FOR UPDATE`
	}
	var id *uuid.UUID
	if err := tx.QueryRow(ctx, q).Scan(&id); err != nil {
		return nil, err
	}
	if id == nil {
		return nil, nil
	}
	return r.FindVersion(ctx, *id)
}

const termsColumns = `id,version,title,body,content_hash,published_by,created_at,published_at`

func (r *TermsRepository) FindVersion(ctx context.Context, id uuid.UUID) (*terms.Version, error) {
	v := new(terms.Version)
	err := r.executor(ctx).QueryRow(ctx, `SELECT `+termsColumns+` FROM terms_versions WHERE id=$1`, id).Scan(&v.ID, &v.Version, &v.Title, &v.Body, &v.ContentHash, &v.PublishedBy, &v.CreatedAt, &v.PublishedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return v, err
}
func (r *TermsRepository) FindAcceptance(ctx context.Context, userID, versionID uuid.UUID) (*terms.Acceptance, error) {
	a := new(terms.Acceptance)
	err := r.executor(ctx).QueryRow(ctx, `SELECT id,user_id,terms_version_id,accepted_at FROM terms_acceptances WHERE user_id=$1 AND terms_version_id=$2`, userID, versionID).Scan(&a.ID, &a.UserID, &a.TermsVersionID, &a.AcceptedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return a, err
}
func (r *TermsRepository) HasAnyAcceptance(ctx context.Context, userID uuid.UUID) (bool, error) {
	var found bool
	err := r.executor(ctx).QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM terms_acceptances WHERE user_id=$1)`, userID).Scan(&found)
	return found, err
}
func (r *TermsRepository) InsertAcceptance(ctx context.Context, userID, versionID uuid.UUID) (*terms.Acceptance, error) {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return nil, shared.ErrInternal
	}
	_, err := tx.Exec(ctx, `INSERT INTO terms_acceptances(id,user_id,terms_version_id) VALUES($1,$2,$3) ON CONFLICT(user_id,terms_version_id) DO NOTHING`, uuid.New(), userID, versionID)
	if err != nil {
		return nil, err
	}
	return r.FindAcceptance(ctx, userID, versionID)
}
func (r *TermsRepository) Publish(ctx context.Context, v *terms.Version) error {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return shared.ErrInternal
	}
	err := tx.QueryRow(ctx, `INSERT INTO terms_versions(id,version,title,body,content_hash,published_by) VALUES($1,$2,$3,$4,$5,$6) RETURNING created_at,published_at`, v.ID, v.Version, v.Title, v.Body, v.ContentHash, v.PublishedBy).Scan(&v.CreatedAt, &v.PublishedAt)
	if isUniqueViolation(err, "terms_versions_version_key") {
		return terms.ErrVersionExists
	}
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE terms_publication SET current_version_id=$1 WHERE id=1`, v.ID)
	return err
}
