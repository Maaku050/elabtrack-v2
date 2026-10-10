package postgres

import (
	"context"
	"errors"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/profile"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct{ pool *pgxpool.Pool }

func NewProfileRepository(p *pgxpool.Pool) *ProfileRepository { return &ProfileRepository{p} }
func (r *ProfileRepository) executor(c context.Context) executor {
	if t, ok := database.TxFromContext(c); ok {
		return t
	}
	return r.pool
}
func (r *ProfileRepository) Metadata(c context.Context, id uuid.UUID) (v d.Metadata, e error) {
	v.AccountID = id
	e = r.executor(c).QueryRow(c, `SELECT image_id,version FROM profile_images WHERE account_id=$1`, id).Scan(&v.ImageID, &v.Version)
	if errors.Is(e, pgx.ErrNoRows) {
		e = nil
	}
	return v, accountError(e)
}
func (r *ProfileRepository) Image(c context.Context, id, image uuid.UUID) (v d.Image, e error) {
	v.AccountID = id
	e = r.executor(c).QueryRow(c, `SELECT image_id,version,png,sha256,width,height FROM profile_images WHERE account_id=$1 AND image_id=$2`, id, image).Scan(&v.ImageID, &v.Version, &v.PNG, &v.Hash, &v.Width, &v.Height)
	return v, accountError(e)
}
func (r *ProfileRepository) Save(c context.Context, v d.Image) error {
	var raw any = v.PNG
	var hash any = v.Hash
	var width any = v.Width
	var height any = v.Height
	if v.ImageID == nil {
		raw = nil
		hash = nil
		width = nil
		height = nil
	}
	tag, e := r.executor(c).Exec(c, `INSERT INTO profile_images(account_id,image_id,version,png,sha256,width,height)VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT(account_id) DO UPDATE SET image_id=EXCLUDED.image_id,version=EXCLUDED.version,png=EXCLUDED.png,sha256=EXCLUDED.sha256,width=EXCLUDED.width,height=EXCLUDED.height,updated_at=clock_timestamp() WHERE profile_images.version=EXCLUDED.version-1`, v.AccountID, v.ImageID, v.Version, raw, hash, width, height)
	if e != nil {
		return accountError(e)
	}
	if tag.RowsAffected() != 1 {
		return shared.ErrConflict
	}
	return nil
}
