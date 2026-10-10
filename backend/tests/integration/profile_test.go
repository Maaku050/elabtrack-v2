package integration

import (
	"bytes"
	"errors"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/profile"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/profileimage"
	"github.com/google/uuid"
	"image"
	"image/png"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRealProfileImages(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	tx := database.NewTxManager(db.Pool)
	ar := postgres.NewAccountsRepository(db.Pool)
	repo := postgres.NewProfileRepository(db.Pool)
	s := app.NewService(repo, ar, tx, profileimage.Validator{})
	admin, staff, borrower, other := termsUser(t, ctx, db, user.RoleAdmin), termsUser(t, ctx, db, user.RoleStaff), termsUser(t, ctx, db, user.RoleBorrower), termsUser(t, ctx, db, user.RoleBorrower)
	var buffer bytes.Buffer
	require(t, png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, 32, 32))) == nil, "synthetic valid PNG")
	raw := buffer.Bytes()
	t.Run("ownership_scope_and_real_png", func(t *testing.T) {
		v, e := s.Metadata(ctx, borrower.ID, borrower.ID)
		require(t, e == nil && v.ImageID == nil && v.Version == 0, "honest initial fallback")
		key := uuid.NewString()
		v, e = s.Save(ctx, borrower.ID, borrower.ID, key, 0, raw, false)
		require(t, e == nil && v.ImageID != nil && v.Version == 1, "own optional image")
		again, e := s.Save(ctx, borrower.ID, borrower.ID, key, 0, raw, false)
		require(t, e == nil && again.AccountID == v.AccountID && again.Version == v.Version && *again.ImageID == *v.ImageID, "same key exact replay")
		for _, actor := range []uuid.UUID{borrower.ID, staff.ID, admin.ID} {
			img, e := s.Image(ctx, actor, borrower.ID, *v.ImageID)
			require(t, e == nil && img.Width == 32 && img.Height == 32 && len(img.PNG) > 0, "authorized protected canonical image")
		}
		_, e = s.Image(ctx, other.ID, borrower.ID, *v.ImageID)
		require(t, errors.Is(e, shared.ErrForbidden), "cross borrower forbidden")
		_, e = s.Save(ctx, admin.ID, borrower.ID, uuid.NewString(), 1, raw, false)
		require(t, errors.Is(e, shared.ErrForbidden), "even Admin cannot overwrite another profile photo")
		_, e = s.Metadata(ctx, staff.ID, admin.ID)
		require(t, errors.Is(e, shared.ErrForbidden), "staff cannot read Admin photo")
		_, e = s.Image(ctx, borrower.ID, borrower.ID, uuid.New())
		require(t, errors.Is(e, shared.ErrNotFound), "mismatched image identity")
		_, e = s.Save(ctx, borrower.ID, borrower.ID, key, 1, raw, false)
		require(t, e != nil, "changed payload same key conflicts")
	})
	t.Run("validation_and_constraints", func(t *testing.T) {
		for _, bad := range [][]byte{nil, []byte("<svg onload='alert(1)'/>"), bytes.Repeat([]byte{1}, 524289), raw[:12]} {
			_, e := s.Save(ctx, other.ID, other.ID, uuid.NewString(), 0, bad, false)
			require(t, errors.Is(e, shared.ErrInvalidInput), "invalid image rejected")
		}
		var huge bytes.Buffer
		require(t, png.Encode(&huge, image.NewRGBA(image.Rect(0, 0, 2049, 1))) == nil, "oversize pixel PNG")
		_, e := s.Save(ctx, other.ID, other.ID, uuid.NewString(), 0, huge.Bytes(), false)
		require(t, errors.Is(e, shared.ErrInvalidInput), "dimension limit")
		v, e := s.Metadata(ctx, other.ID, other.ID)
		require(t, e == nil && v.Version == 0, "invalid image leaves no partial profile")
		_, e = db.Pool.Exec(ctx, `INSERT INTO profile_images(account_id,image_id,version,png,sha256,width,height) VALUES($1,gen_random_uuid(),1,NULL,NULL,NULL,NULL)`, other.ID)
		require(t, e != nil, "null bypass blocked by database")
		_, e = db.Pool.Exec(ctx, `DELETE FROM profile_images WHERE account_id=$1`, borrower.ID)
		require(t, e != nil, "runtime cannot hard-delete")
		_, e = db.Pool.Exec(ctx, `TRUNCATE profile_images`)
		require(t, e != nil, "runtime cannot truncate")
	})
	t.Run("atomic_rollback_and_racing_stale_reviews", func(t *testing.T) {
		abort := app.NewService(repo, ar, controlledTx{base: tx, after: func() error { return shared.ErrConflict }}, profileimage.Validator{})
		key := uuid.NewString()
		_, e := abort.Save(ctx, other.ID, other.ID, key, 0, raw, false)
		require(t, e != nil, "forced rollback")
		v, e := s.Metadata(ctx, other.ID, other.ID)
		require(t, e == nil && v.Version == 0, "photo rolled back")
		var audits, receipts int
		require(t, owner.Pool.QueryRow(ctx, `SELECT (SELECT count(*) FROM account_audit_events WHERE account_id=$1),(SELECT count(*) FROM account_operation_receipts WHERE actor_id=$1)`, other.ID).Scan(&audits, &receipts) == nil && audits == 0 && receipts == 0, "audit and receipt rolled back")
		var wins atomic.Int32
		var wg sync.WaitGroup
		for n := 0; n < 2; n++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, e := s.Save(ctx, other.ID, other.ID, uuid.NewString(), 0, raw, false)
				if e == nil {
					wins.Add(1)
				} else if !errors.Is(e, shared.ErrConflict) {
					t.Error(e)
				}
			}()
		}
		wg.Wait()
		require(t, wins.Load() == 1, "one concurrent version winner")
		v, e = s.Metadata(ctx, other.ID, other.ID)
		require(t, e == nil && v.Version == 1, "exact version")
		old := *v.ImageID
		key = uuid.NewString()
		v, e = s.Save(ctx, other.ID, other.ID, key, 1, nil, true)
		require(t, e == nil && v.ImageID == nil && v.Version == 2, "remove returns initials without deleting history")
		_, e = s.Image(ctx, admin.ID, other.ID, old)
		require(t, errors.Is(e, shared.ErrNotFound), "old bytes no longer retrievable")
		_, e = s.Save(ctx, other.ID, other.ID, key, 1, nil, true)
		require(t, e == nil, "remove replay")
	})
	t.Run("current_account_and_role_enforcement", func(t *testing.T) {
		_, e := owner.Pool.Exec(ctx, `UPDATE users SET is_active=false WHERE id=$1`, borrower.ID)
		require(t, e == nil, "fixture inactive")
		_, e = s.Save(ctx, borrower.ID, borrower.ID, uuid.NewString(), 1, raw, false)
		require(t, errors.Is(e, shared.ErrUnauthorized), "deactivated write forbidden")
		_, e = s.Metadata(ctx, borrower.ID, borrower.ID)
		require(t, errors.Is(e, shared.ErrUnauthorized), "deactivated read forbidden")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET role='BORROWER' WHERE id=$1`, staff.ID)
		require(t, e == nil, "fixture demotion")
		_, e = s.Metadata(ctx, staff.ID, borrower.ID)
		require(t, errors.Is(e, shared.ErrForbidden), "current demotion removes operational read")
	})
}
