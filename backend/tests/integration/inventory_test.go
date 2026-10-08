package integration

import (
	"bytes"
	"context"
	"errors"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
	"image"
	"image/png"
	"sync"
	"sync/atomic"
	"testing"
)

type inventoryAbortTx struct{ base *database.TxManager }

func (t inventoryAbortTx) Within(c context.Context, fn func(context.Context) error) error {
	return t.base.Within(c, func(ctx context.Context) error {
		if e := fn(ctx); e != nil {
			return e
		}
		return shared.ErrConflict
	})
}
func TestRealInventory(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	r := postgres.NewInventoryRepository(db.Pool)
	a := postgres.NewAccountsRepository(db.Pool)
	tx := database.NewTxManager(db.Pool)
	s := app.NewService(r, a, tx, catalogimage.Validator{})
	admin := termsUser(t, ctx, db, user.RoleAdmin)
	staff := termsUser(t, ctx, db, user.RoleStaff)
	borrower := termsUser(t, ctx, db, user.RoleBorrower)
	yes := true
	t.Run("empty_paired_rollback_preserves_accounts", func(t *testing.T) {
		var count int
		require(t, owner.Pool.QueryRow(ctx, `SELECT count(*) FROM equipment`).Scan(&count) == nil && count == 0, "empty inventory before history")
		m := database.NewMigrator(owner.Pool, "../../migrations")
		require(t, m.Down(ctx) == nil && m.Up(ctx) == nil, "empty007 rollback/reapply")
	})
	cat, e := s.Category(ctx, staff.ID, uuid.Nil, uuid.NewString(), d.CategoryInput{Name: "TEST " + uuid.NewString(), Active: &yes})
	require(t, e == nil, "category create")
	input := d.Create{Metadata: d.Metadata{Name: "TEST cookware " + uuid.NewString(), Description: "Synthetic controlled pool", CategoryID: &cat.ID}, Opening: 5, Reason: "Verified test opening count"}
	v, e := s.Create(ctx, staff.ID, uuid.NewString(), input)
	require(t, e == nil && v.Stock.Available == 5 && v.Stock.Total == 5 && v.Sequence == 1, "opening movement")
	t.Run("authority_category_identity_and_metadata", func(t *testing.T) {
		_, e := s.Create(ctx, borrower.ID, uuid.NewString(), input)
		require(t, errors.Is(e, shared.ErrForbidden), "Borrower cannot write")
		_, e = s.Category(ctx, borrower.ID, uuid.Nil, uuid.NewString(), d.CategoryInput{Name: "X", Active: &yes})
		require(t, errors.Is(e, shared.ErrForbidden), "Borrower category write denied")
		_, e = s.Category(ctx, admin.ID, uuid.Nil, uuid.NewString(), d.CategoryInput{Name: cat.Name, Active: &yes})
		require(t, errors.Is(e, shared.ErrConflict), "unique category case insensitive")
		bad := input
		absent := uuid.New()
		bad.CategoryID = &absent
		_, e = s.Create(ctx, staff.ID, uuid.NewString(), bad)
		require(t, errors.Is(e, shared.ErrNotFound), "category ownership FK")
		v, e = s.Edit(ctx, staff.ID, v.ID, uuid.NewString(), d.Metadata{Name: v.Name + " edited", Description: "Updated metadata", CategoryID: &cat.ID, ExpectedVersion: v.Version})
		require(t, e == nil && v.Version == 2 && v.Stock.Available == 5, "metadata stock unchanged")
		_, e = s.Edit(ctx, staff.ID, v.ID, uuid.NewString(), d.Metadata{Name: "stale", ExpectedVersion: 1})
		require(t, errors.Is(e, shared.ErrConflict), "stale metadata conflicts")
	})
	t.Run("bounded_search_filters_visibility", func(t *testing.T) {
		p, e := s.List(ctx, borrower.ID, d.Filter{Page: 1, PerPage: 1, Search: v.Name, Sort: "available", AvailableOnly: true})
		require(t, e == nil && len(p.Items) == 1 && p.Total == 1, "real bounded query")
		p, e = s.List(ctx, borrower.ID, d.Filter{Page: 1, PerPage: 25, Search: "%' OR TRUE--"})
		require(t, e == nil && len(p.Items) == 0, "literal wildcard sql input")
		_, e = s.List(ctx, admin.ID, d.Filter{Page: 0, PerPage: 1000})
		require(t, errors.Is(e, shared.ErrInvalidInput), "bounded pagination")
		v, e = s.Status(ctx, staff.ID, v.ID, uuid.NewString(), d.StatusInput{Status: "INACTIVE", Confirm: true, ExpectedVersion: v.Version})
		require(t, e == nil, "inactive")
		_, e = s.Detail(ctx, borrower.ID, v.ID)
		require(t, errors.Is(e, shared.ErrNotFound), "hidden inactive ID")
		p, e = s.List(ctx, borrower.ID, d.Filter{Page: 1, PerPage: 25, Status: "INACTIVE"})
		require(t, errors.Is(e, shared.ErrForbidden), "cannot query hidden status")
		v, e = s.Status(ctx, staff.ID, v.ID, uuid.NewString(), d.StatusInput{Status: "ACTIVE", Confirm: true, ExpectedVersion: v.Version})
		require(t, e == nil, "active")
	})
	t.Run("adjustment_retry_negative_overflow_and_atomicity", func(t *testing.T) {
		key := uuid.NewString()
		in := d.Adjustment{Kind: "ADD", Quantity: 3, Reason: "Verified acquisition", Confirm: true}
		v, e = s.Adjust(ctx, staff.ID, v.ID, key, in)
		require(t, e == nil && v.Stock.Available == 8, "add")
		again, e := s.Adjust(ctx, staff.ID, v.ID, key, in)
		require(t, e == nil && again.Sequence == v.Sequence, "samekey singlemovement")
		in.Quantity = 4
		_, e = s.Adjust(ctx, staff.ID, v.ID, key, in)
		require(t, errors.Is(e, shared.ErrConflict), "changed key")
		for _, q := range []int64{-1, 9, d.MaxQuantity} {
			_, e = s.Adjust(ctx, staff.ID, v.ID, uuid.NewString(), d.Adjustment{Kind: "REMOVE", Quantity: q, Reason: "Invalid removal", Confirm: true})
			require(t, e != nil, "negative insufficient overflow rejected")
		}
		abort := app.NewService(r, a, inventoryAbortTx{tx}, catalogimage.Validator{})
		_, e = abort.Adjust(ctx, staff.ID, v.ID, uuid.NewString(), d.Adjustment{Kind: "ADD", Quantity: 2, Reason: "Forced rollback", Confirm: true})
		require(t, e != nil, "forced receipt/audit stock rollback")
		current, e := s.Detail(ctx, staff.ID, v.ID)
		require(t, e == nil && current.Sequence == v.Sequence && current.Stock == v.Stock, "atomic unchanged")
	})
	t.Run("admin_reconciliation_preserves_reserved_custody_damaged", func(t *testing.T) {
		_, e := owner.Pool.Exec(ctx, `UPDATE equipment SET reserved=2,checked_out=3,damaged_held=4,total_tracked=available+9 WHERE id=$1`, v.ID)
		require(t, e == nil, "synthetic future buckets fixture")
		seq := v.Sequence
		in := d.Adjustment{Kind: "RECONCILE", Quantity: 6, ExpectedSequence: &seq, Reason: "Verified physical available count", Confirm: true}
		_, e = s.Adjust(ctx, staff.ID, v.ID, uuid.NewString(), in)
		require(t, errors.Is(e, shared.ErrForbidden), "Admin-only correction")
		v, e = s.Adjust(ctx, admin.ID, v.ID, uuid.NewString(), in)
		require(t, e == nil && v.Stock.Available == 6 && v.Stock.Reserved == 2 && v.Stock.CheckedOut == 3 && v.Stock.DamagedHeld == 4 && v.Stock.Total == 15, "custody originals preserved")
		_, e = s.Adjust(ctx, admin.ID, v.ID, uuid.NewString(), in)
		require(t, errors.Is(e, shared.ErrConflict), "stale count basis")
		_, e = s.Status(ctx, admin.ID, v.ID, uuid.NewString(), d.StatusInput{Status: "ARCHIVED", ExpectedVersion: v.Version, Confirm: true})
		require(t, errors.Is(e, shared.ErrConflict), "archive blocks holds custody")
		_, e = owner.Pool.Exec(ctx, `UPDATE equipment SET reserved=0,checked_out=0,total_tracked=available+damaged_held WHERE id=$1`, v.ID)
		require(t, e == nil, "fixture cleanup only synthetic custody")
	})
	t.Run("concurrent_last_stock_and_duplicate_key", func(t *testing.T) {
		var wins atomic.Int32
		var wg sync.WaitGroup
		for i := 0; i < 12; i++ {
			actor := staff.ID
			if i%2 == 0 {
				actor = admin.ID
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, e := s.Adjust(ctx, actor, v.ID, uuid.NewString(), d.Adjustment{Kind: "REMOVE", Quantity: 1, Reason: "Concurrent usable removal", Confirm: true})
				if e == nil {
					wins.Add(1)
				} else if !errors.Is(e, shared.ErrConflict) {
					t.Error("unexpected concurrency error", e)
				}
			}()
		}
		wg.Wait()
		require(t, wins.Load() == 6, "six units six winners")
		key := uuid.NewString()
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, e := s.Adjust(ctx, staff.ID, v.ID, key, d.Adjustment{Kind: "ADD", Quantity: 2, Reason: "Concurrent same key", Confirm: true})
				if e != nil {
					t.Error(e)
				}
			}()
		}
		wg.Wait()
		v, e = s.Detail(ctx, staff.ID, v.ID)
		require(t, e == nil && v.Stock.Available == 2 && v.Stock.DamagedHeld == 4, "retry one acquisition")
		var sum int64
		require(t, owner.Pool.QueryRow(ctx, `SELECT sum(delta_available) FROM inventory_movements WHERE equipment_id=$1`, v.ID).Scan(&sum) == nil && sum == v.Stock.Available, "ledger reconciles")
	})
	t.Run("safe_images_and_ownership", func(t *testing.T) {
		var b bytes.Buffer
		_ = png.Encode(&b, image.NewRGBA(image.Rect(0, 0, 8, 8)))
		key := uuid.NewString()
		old := v.Version
		v, e = s.SaveImage(ctx, staff.ID, v.ID, key, old, b.Bytes())
		require(t, e == nil && v.ImageID != nil, "catalog raster")
		same, e := s.SaveImage(ctx, staff.ID, v.ID, key, old, b.Bytes())
		require(t, e == nil && same.ImageID != nil && *same.ImageID == *v.ImageID, "image retry")
		im, e := s.Image(ctx, borrower.ID, v.ID, *v.ImageID)
		require(t, e == nil && len(im.PNG) > 0, "authorized active raster")
		_, e = s.Image(ctx, borrower.ID, uuid.New(), *v.ImageID)
		require(t, errors.Is(e, shared.ErrNotFound), "image parent scoped")
		_, e = s.SaveImage(ctx, borrower.ID, v.ID, uuid.NewString(), v.Version, b.Bytes())
		require(t, errors.Is(e, shared.ErrForbidden), "Borrower cannot upload")
		_, e = s.SaveImage(ctx, staff.ID, v.ID, uuid.NewString(), v.Version, []byte("<svg/>"))
		require(t, e != nil, "unsupported image rejects")
	})
	t.Run("archive_future_boundary_and_history_guards", func(t *testing.T) {
		_, e := owner.Pool.Exec(ctx, `CREATE TABLE borrowings(batch1_probe bool)`)
		require(t, e == nil, "test-only future table probe")
		_, e = s.Status(ctx, staff.ID, v.ID, uuid.NewString(), d.StatusInput{Status: "ARCHIVED", ExpectedVersion: v.Version, Confirm: true})
		require(t, errors.Is(e, shared.ErrConflict), "future liability unconnected fails closed")
		_, e = owner.Pool.Exec(ctx, `DROP TABLE borrowings`)
		require(t, e == nil, "owned probe cleanup")
		v, e = s.Status(ctx, staff.ID, v.ID, uuid.NewString(), d.StatusInput{Status: "ARCHIVED", ExpectedVersion: v.Version, Confirm: true})
		require(t, e == nil, "safe current archive")
		_, e = s.Adjust(ctx, staff.ID, v.ID, uuid.NewString(), d.Adjustment{Kind: "ADD", Quantity: 1, Confirm: true, Reason: "No archived acquisition"})
		require(t, errors.Is(e, shared.ErrConflict), "archived immutable stock")
		_, e = s.Detail(ctx, borrower.ID, v.ID)
		require(t, errors.Is(e, shared.ErrNotFound), "archive hidden")
		for _, q := range []string{`DELETE FROM equipment`, `UPDATE equipment SET reserved=reserved`, `UPDATE equipment SET checked_out=checked_out`, `UPDATE equipment SET damaged_held=damaged_held`, `UPDATE inventory_movements SET reason='rewrite'`, `DELETE FROM inventory_audit_events`, `TRUNCATE inventory_operation_receipts`, `UPDATE equipment_images SET width=1`, `CREATE TABLE forbidden_inventory(id int)`} {
			_, e = db.Pool.Exec(ctx, q)
			require(t, e != nil, "runtime least privilege immutable evidence")
		}
		_, e = owner.Pool.Exec(ctx, `UPDATE inventory_movements SET reason='rewrite'`)
		require(t, e != nil, "owner history trigger")
		var narrow bool
		require(t, owner.Pool.QueryRow(ctx, `SELECT NOT has_column_privilege('elabtrack_runtime','equipment','reserved','UPDATE') AND NOT has_column_privilege('elabtrack_runtime','equipment','checked_out','UPDATE') AND NOT has_column_privilege('elabtrack_runtime','equipment','damaged_held','UPDATE')`).Scan(&narrow) == nil && narrow, "no premature custody write grants")
		_, e = owner.Pool.Exec(ctx, `UPDATE equipment SET available=-1 WHERE id=$1`, v.ID)
		require(t, e != nil, "DB negative constraint")
		_, e = owner.Pool.Exec(ctx, `UPDATE equipment SET total_tracked=99 WHERE id=$1`, v.ID)
		require(t, e != nil, "DB conservation constraint")
		require(t, database.NewMigrator(owner.Pool, "../../migrations").Down(ctx) != nil, "history rollback refused")
	})
}

func TestRealInventoryContention(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	repo := postgres.NewInventoryRepository(db.Pool)
	s := app.NewService(repo, postgres.NewAccountsRepository(db.Pool), database.NewTxManager(db.Pool), catalogimage.Validator{})
	staff := termsUser(t, ctx, db, user.RoleStaff)
	admin := termsUser(t, ctx, db, user.RoleAdmin)
	yes, no := true, false
	cat, e := s.Category(ctx, staff.ID, uuid.Nil, uuid.NewString(), d.CategoryInput{Name: "TEST independent actors " + uuid.NewString(), Active: &yes})
	require(t, e == nil, "category")
	v, e := s.Create(ctx, staff.ID, uuid.NewString(), d.Create{Metadata: d.Metadata{Name: "TEST independent actor pool", CategoryID: &cat.ID}, Opening: 10, Reason: "Verified test opening"})
	require(t, e == nil, "opening")
	cat, e = s.Category(ctx, admin.ID, cat.ID, uuid.NewString(), d.CategoryInput{Name: cat.Name, Active: &no, ExpectedVersion: cat.Version})
	require(t, e == nil, "category inactive")
	v, e = s.Edit(ctx, staff.ID, v.ID, uuid.NewString(), d.Metadata{Name: v.Name, Description: "Keep historical inactive category", CategoryID: &cat.ID, ExpectedVersion: v.Version})
	require(t, e == nil && v.CategoryID != nil && *v.CategoryID == cat.ID, "unchanged inactive reference preserved")
	_, e = s.Create(ctx, staff.ID, uuid.NewString(), d.Create{Metadata: d.Metadata{Name: "New disallowed category", CategoryID: &cat.ID}})
	require(t, errors.Is(e, shared.ErrConflict), "new inactive association rejected")
	t.Run("competing_actor_equipment_locks", func(t *testing.T) {
		var wg sync.WaitGroup
		var wins atomic.Int32
		for i := 0; i < 20; i++ {
			actor := staff.ID
			if i%2 == 0 {
				actor = admin.ID
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := s.Adjust(ctx, actor, v.ID, uuid.NewString(), d.Adjustment{Kind: "REMOVE", Quantity: 1, Reason: "Independent actor removal", Confirm: true})
				if err == nil {
					wins.Add(1)
				} else if !errors.Is(err, shared.ErrConflict) {
					t.Error(err)
				}
			}()
		}
		wg.Wait()
		require(t, wins.Load() == 10, "ten last-stock winners across actors")
		v, e = s.Detail(ctx, staff.ID, v.ID)
		require(t, e == nil && v.Stock.Total == 0 && v.Sequence == 11, "no underflow or lost movements")
	})
	t.Run("optimistic_metadata_single_winner", func(t *testing.T) {
		var wg sync.WaitGroup
		var wins atomic.Int32
		version := v.Version
		for _, actor := range []uuid.UUID{staff.ID, admin.ID} {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := s.Edit(ctx, actor, v.ID, uuid.NewString(), d.Metadata{Name: "TEST edited " + actor.String(), CategoryID: &cat.ID, ExpectedVersion: version})
				if err == nil {
					wins.Add(1)
				} else if !errors.Is(err, shared.ErrConflict) {
					t.Error(err)
				}
			}()
		}
		wg.Wait()
		require(t, wins.Load() == 1, "one metadata winner")
		v, e = s.Detail(ctx, staff.ID, v.ID)
		require(t, e == nil && v.Version == version+1 && v.Stock.Total == 0, "independent metadata and stock versions")
	})
	t.Run("ledger_vector_and_audit_receipt_conservation", func(t *testing.T) {
		var sumA, sumR, sumC, sumD, sumT int64
		require(t, owner.Pool.QueryRow(ctx, `SELECT sum(delta_available),sum(delta_reserved),sum(delta_checked_out),sum(delta_damaged_held),sum(delta_total) FROM inventory_movements WHERE equipment_id=$1`, v.ID).Scan(&sumA, &sumR, &sumC, &sumD, &sumT) == nil, "ledger sum")
		require(t, v.Stock == (d.Stock{Available: sumA, Reserved: sumR, CheckedOut: sumC, DamagedHeld: sumD, Total: sumT}), "all physical vectors reconcile")
		var count int
		require(t, owner.Pool.QueryRow(ctx, `SELECT count(*) FROM inventory_audit_events WHERE equipment_id=$1`, v.ID).Scan(&count) == nil && count == 13, "creation+10 movements+2 metadata audited")
		_, e = s.Adjust(ctx, admin.ID, v.ID, uuid.NewString(), d.Adjustment{Kind: "ADD", Quantity: d.MaxQuantity, Reason: "Bounded max acquisition", Confirm: true})
		require(t, e == nil, "maximum supported pool")
		_, e = s.Adjust(ctx, staff.ID, v.ID, uuid.NewString(), d.Adjustment{Kind: "ADD", Quantity: 1, Reason: "Overflow rejected", Confirm: true})
		require(t, errors.Is(e, shared.ErrConflict), "real overflow rejects without evidence")
	})
	current, e := s.Detail(ctx, staff.ID, v.ID)
	require(t, e == nil, "final metadata lookup")
	_, e = s.Status(ctx, staff.ID, v.ID, uuid.NewString(), d.StatusInput{Status: "ARCHIVED", ExpectedVersion: current.Version, Confirm: true})
	require(t, e == nil, "preserve fixture ledger through safe archive")

}
