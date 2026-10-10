package integration

import (
	"context"
	"errors"
	accountapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/accounts"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	appinventory "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	appterms "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	accountdomain "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	termsdomain "github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/security"
	"github.com/google/uuid"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// A repository clock override is test-only; persisted deadlines are untouched.
type borrowingClock struct {
	d.Repository
	at time.Time
}

func (r borrowingClock) Clock(context.Context) (time.Time, error) { return r.at, nil }
func TestRealBorrowingFoundation(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	tx := database.NewTxManager(db.Pool)
	ar := postgres.NewAccountsRepository(db.Pool)
	ir := postgres.NewInventoryRepository(db.Pool)
	br := postgres.NewBorrowingRepository(db.Pool)
	ts := appterms.NewService(postgres.NewTermsRepository(db.Pool), postgres.NewUserRepository(db.Pool), tx)
	s := app.NewService(br, ar, ir, ts, tx)
	is := appinventory.NewService(ir, ar, tx, catalogimage.Validator{})
	admin := termsUser(t, ctx, db, user.RoleAdmin)
	staff := termsUser(t, ctx, db, user.RoleStaff)
	var cur *uuid.UUID
	var current uuid.UUID
	require(t, owner.Pool.QueryRow(ctx, `SELECT current_version_id FROM terms_publication WHERE id=1`).Scan(&cur) == nil, "publication")
	version, e := ts.Publish(ctx, appterms.PublishCommand{ActorID: admin.ID, Version: "TEST-" + uuid.NewString()[:12], Title: "TEST ONLY", Body: "Synthetic isolated integration terms; not FSMO policy.", ExpectedCurrentVersionID: cur})
	require(t, e == nil, "test terms publication")
	current = version.ID
	borrower := func() *user.User {
		u := termsUser(t, ctx, db, user.RoleBorrower)
		_, e := owner.Pool.Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type) VALUES($1,'FACULTY')`, u.ID)
		require(t, e == nil, "fixture profile")
		_, e = ts.Accept(ctx, u.ID, current)
		require(t, e == nil, "test borrower consent")
		return u
	}
	b1, b2 := borrower(), borrower()
	equipment := func(q int64) inventory.Equipment {
		v, e := is.Create(ctx, staff.ID, uuid.NewString(), inventory.Create{Metadata: inventory.Metadata{Name: "TEST Phase7 " + uuid.NewString()}, Opening: q, Reason: "Synthetic verification acquisition"})
		require(t, e == nil, "stock fixture")
		return v
	}
	input := func(v inventory.Equipment, q int64) d.Input {
		return d.Input{Confirm: true, Items: []d.Line{{EquipmentID: v.ID, Quantity: q}}}
	}
	verify := func(v inventory.Equipment, A, R, C int64) {
		x, e := ir.Get(ctx, v.ID)
		require(t, e == nil && x.Stock.Available == A && x.Stock.Reserved == R && x.Stock.CheckedOut == C && x.Stock.DamagedHeld == 0 && x.Stock.Total == A+R+C, "exact physical buckets")
		var a, r, c, total int64
		require(t, owner.Pool.QueryRow(ctx, `SELECT sum(delta_available),sum(delta_reserved),sum(delta_checked_out),sum(delta_total) FROM inventory_movements WHERE equipment_id=$1`, v.ID).Scan(&a, &r, &c, &total) == nil && a == A && r == R && c == C && total == x.Stock.Total, "ledger conservation")
	}
	t.Run("last_unit_and_many_reservations", func(t *testing.T) {
		v := equipment(1)
		var wins atomic.Int32
		var wg sync.WaitGroup
		for _, u := range []*user.User{b1, b2} {
			wg.Add(1)
			go func(u *user.User) {
				defer wg.Done()
				_, e := s.Submit(ctx, u.ID, uuid.NewString(), input(v, 1))
				if e == nil {
					wins.Add(1)
				} else if !errors.Is(e, d.ErrStock) {
					t.Error(e)
				}
			}(u)
		}
		wg.Wait()
		require(t, wins.Load() == 1, "one last-unit winner")
		verify(v, 0, 1, 0)
		v = equipment(7)
		wins.Store(0)
		for i := 0; i < 20; i++ {
			u := b1
			if i%2 == 0 {
				u = b2
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, e := s.Submit(ctx, u.ID, uuid.NewString(), input(v, 1))
				if e == nil {
					wins.Add(1)
				} else if !errors.Is(e, d.ErrStock) {
					t.Error(e)
				}
			}()
		}
		wg.Wait()
		require(t, wins.Load() == 7, "seven reservation winners")
		verify(v, 0, 7, 0)
	})
	t.Run("multi_item_rollback_retries_and_invalid", func(t *testing.T) {
		v, w := equipment(4), equipment(0)
		in := input(v, 2)
		in.Items = append(in.Items, d.Line{EquipmentID: w.ID, Quantity: 1})
		_, e := s.Submit(ctx, b1.ID, uuid.NewString(), in)
		require(t, errors.Is(e, d.ErrStock), "partial conflict")
		verify(v, 4, 0, 0)
		in = input(v, 2)
		key := uuid.NewString()
		r, e := s.Submit(ctx, b1.ID, key, in)
		require(t, e == nil && r.Status == "PENDING" && r.ExpiresAt.Sub(r.CreatedAt) == 24*time.Hour, "atomic pending")
		again, e := s.Submit(ctx, b1.ID, key, in)
		require(t, e == nil && again.ID == r.ID, "replay identity")
		in.Items[0].Quantity = 1
		_, e = s.Submit(ctx, b1.ID, key, in)
		require(t, errors.Is(e, d.ErrKey), "changed payload conflict")
		verify(v, 2, 2, 0)
		abort := app.NewService(br, ar, ir, ts, controlledTx{base: tx, after: func() error { return shared.ErrConflict }})
		_, e = abort.Submit(ctx, b2.ID, uuid.NewString(), input(v, 1))
		require(t, e != nil, "interrupted transaction")
		verify(v, 2, 2, 0)
		for _, q := range []int64{0, -1, inventory.MaxQuantity + 1} {
			_, e = s.Submit(ctx, b1.ID, uuid.NewString(), input(v, q))
			require(t, e != nil, "invalid quantities")
		}
		dup := input(v, 1)
		dup.Items = append(dup.Items, dup.Items[0])
		_, e = s.Submit(ctx, b1.ID, uuid.NewString(), dup)
		require(t, e != nil, "duplicate lines")
		_, e = s.Submit(ctx, staff.ID, uuid.NewString(), input(v, 1))
		require(t, errors.Is(e, shared.ErrForbidden), "staff cannot submit own")
	})
	t.Run("eligibility_terms_and_inactive_equipment", func(t *testing.T) {
		v := equipment(3)
		u := borrower()
		_, e := owner.Pool.Exec(ctx, `UPDATE users SET activation_required=true WHERE id=$1`, u.ID)
		require(t, e == nil, "test activation")
		_, e = s.Submit(ctx, u.ID, uuid.NewString(), input(v, 1))
		require(t, e != nil, "activation blocked")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET activation_required=false,is_active=false WHERE id=$1`, u.ID)
		require(t, e == nil, "test inactive")
		_, e = s.Submit(ctx, u.ID, uuid.NewString(), input(v, 1))
		require(t, e != nil, "inactive blocked")
		v, e = is.Status(ctx, staff.ID, v.ID, uuid.NewString(), inventory.StatusInput{Status: "INACTIVE", Confirm: true, ExpectedVersion: v.Version})
		require(t, e == nil, "inactive equipment")
		_, e = s.Submit(ctx, b1.ID, uuid.NewString(), input(v, 1))
		require(t, errors.Is(e, d.ErrStock), "inactive equipment blocked")
		verify(v, 3, 0, 0)
		noConsent := termsUser(t, ctx, db, user.RoleBorrower)
		_, e = owner.Pool.Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type) VALUES($1,'FACULTY')`, noConsent.ID)
		require(t, e == nil, "no acceptance profile")
		_, e = s.Submit(ctx, noConsent.ID, uuid.NewString(), input(v, 1))
		require(t, e != nil, "missing acceptance blocked")
	})
	t.Run("cancel_and_expiry_release_once_after_restart", func(t *testing.T) {
		v := equipment(5)
		r, e := s.Submit(ctx, b1.ID, uuid.NewString(), input(v, 2))
		require(t, e == nil, "submit")
		_, e = s.Decide(ctx, b2.ID, r.ID, uuid.NewString(), "CANCELLED", d.Decision{Confirm: true})
		require(t, errors.Is(e, shared.ErrNotFound), "other owner hidden")
		key := uuid.NewString()
		r, e = s.Decide(ctx, b1.ID, r.ID, key, "CANCELLED", d.Decision{Confirm: true})
		require(t, e == nil && r.Status == "CANCELLED", "cancel")
		_, e = s.Decide(ctx, b1.ID, r.ID, key, "CANCELLED", d.Decision{Confirm: true})
		require(t, e == nil, "cancel replay")
		verify(v, 5, 0, 0)
		r, e = s.Submit(ctx, b1.ID, uuid.NewString(), input(v, 2))
		require(t, e == nil, "expiry pending")
		late := app.NewService(borrowingClock{br, *r.ExpiresAt}, ar, ir, ts, tx)
		r, e = late.Decide(ctx, b1.ID, r.ID, uuid.NewString(), "CANCELLED", d.Decision{Confirm: true})
		require(t, errors.Is(e, d.ErrExpired) && r.Status == "EXPIRED", "boundary expiry commits before409")
		verify(v, 5, 0, 0)
		fresh := app.NewService(br, ar, ir, ts, tx)
		n, e := fresh.Sweep(ctx, 100)
		require(t, e == nil && n == 0, "restart no duplicate release")
		stored, e := fresh.Read(ctx, b1.ID, r.ID)
		require(t, e == nil && len(stored.Events) == 2 && stored.Events[1].ActorID == nil, "durable system expiry audit")
	})
}

func TestRealBorrowingOperations(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	tx := database.NewTxManager(db.Pool)
	ar := postgres.NewAccountsRepository(db.Pool)
	ir := postgres.NewInventoryRepository(db.Pool)
	br := postgres.NewBorrowingRepository(db.Pool)
	ts := appterms.NewService(postgres.NewTermsRepository(db.Pool), postgres.NewUserRepository(db.Pool), tx)
	s := app.NewService(br, ar, ir, ts, tx)
	is := appinventory.NewService(ir, ar, tx, catalogimage.Validator{})
	admin, staff := termsUser(t, ctx, db, user.RoleAdmin), termsUser(t, ctx, db, user.RoleStaff)
	version, e := ts.Current(ctx, admin.ID)
	require(t, e == nil, "existing isolated terms")
	borrower := func(student bool) *user.User {
		u := termsUser(t, ctx, db, user.RoleBorrower)
		kind := "FACULTY"
		var sid *string
		if student {
			kind = "STUDENT"
			id := "00" + uuid.NewString()
			sid = &id
		}
		_, e := owner.Pool.Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type,student_id) VALUES($1,$2,$3)`, u.ID, kind, sid)
		require(t, e == nil, "category fixture")
		_, e = ts.Accept(ctx, u.ID, version.ID)
		require(t, e == nil, "fixture current acceptance")
		return u
	}
	b1, b2 := borrower(true), borrower(false)
	equipment := func(q int64) inventory.Equipment {
		v, e := is.Create(ctx, staff.ID, uuid.NewString(), inventory.Create{Metadata: inventory.Metadata{Name: "TEST operations " + uuid.NewString()}, Opening: q, Reason: "Isolated operation fixture"})
		require(t, e == nil, "stock fixture")
		return v
	}
	input := func(v inventory.Equipment, q int64) d.Input {
		return d.Input{Confirm: true, Items: []d.Line{{EquipmentID: v.ID, Quantity: q}}}
	}
	due := time.Now().Add(14 * 24 * time.Hour).UTC().Truncate(time.Microsecond)
	checkout := d.Decision{Confirm: true, Handover: true, DueAt: &due}
	verify := func(v inventory.Equipment, A, R, C int64) {
		x, e := ir.Get(ctx, v.ID)
		require(t, e == nil && x.Stock.Available == A && x.Stock.Reserved == R && x.Stock.CheckedOut == C && x.Stock.Total == A+R+C, "exact resulting quantities")
		var a, r, c, total int64
		require(t, owner.Pool.QueryRow(ctx, `SELECT sum(delta_available),sum(delta_reserved),sum(delta_checked_out),sum(delta_total) FROM inventory_movements WHERE equipment_id=$1`, v.ID).Scan(&a, &r, &c, &total) == nil && a == A && r == R && c == C && total == x.Stock.Total, "movement sum")
	}
	t.Run("deny_reason_checkout_due_and_replay", func(t *testing.T) {
		v := equipment(4)
		r, e := s.Submit(ctx, b1.ID, uuid.NewString(), input(v, 2))
		require(t, e == nil, "request")
		_, e = s.Decide(ctx, staff.ID, r.ID, uuid.NewString(), "DENIED", d.Decision{Confirm: true})
		require(t, e != nil, "explanation mandatory")
		r, e = s.Decide(ctx, staff.ID, r.ID, uuid.NewString(), "DENIED", d.Decision{Confirm: true, Reason: "FSMO operational conflict"})
		require(t, e == nil && r.Reason == "FSMO operational conflict" && r.Status == "DENIED", "visible denial")
		verify(v, 4, 0, 0)
		r, e = s.Submit(ctx, b2.ID, uuid.NewString(), input(v, 2))
		require(t, e == nil, "pending issue")
		for _, in := range []d.Decision{{Confirm: true}, {Confirm: true, DueAt: &due}, {Handover: true, DueAt: &due}} {
			_, e = s.Decide(ctx, staff.ID, r.ID, uuid.NewString(), "CHECKOUT", in)
			require(t, e != nil, "handover/due/confirmation required")
		}
		key := uuid.NewString()
		r, e = s.Decide(ctx, staff.ID, r.ID, key, "CHECKOUT", checkout)
		require(t, e == nil && r.Status == "CHECKED_OUT" && r.DueAt.Equal(due) && r.Items[0].Issued == 2 && len(r.Events) == 2, "single atomic checkout beyond7days")
		again, e := s.Decide(ctx, staff.ID, r.ID, key, "CHECKOUT", checkout)
		require(t, e == nil && again.ID == r.ID, "checkout replay")
		_, e = s.Decide(ctx, admin.ID, r.ID, uuid.NewString(), "CHECKOUT", checkout)
		require(t, errors.Is(e, d.ErrState), "no second checkout")
		verify(v, 2, 0, 2)
	})
	t.Run("detail_read_never_mixes_transition_headers_and_items", func(t *testing.T) {
		v := equipment(2)
		r, e := s.Submit(ctx, b1.ID, uuid.NewString(), input(v, 2))
		require(t, e == nil, "read consistency request")
		start := make(chan struct{})
		results := make(chan error, 21)
		for n := 0; n < 20; n++ {
			go func() {
				<-start
				read, err := s.Read(ctx, admin.ID, r.ID)
				if err == nil {
					pending := read.Status == "PENDING" && read.Items[0].Reserved == 2 && read.Items[0].Issued == 0 && len(read.Events) == 1
					issued := read.Status == "CHECKED_OUT" && read.Items[0].Reserved == 0 && read.Items[0].Issued == 2 && len(read.Events) == 2
					if !pending && !issued {
						err = errors.New("mixed transition snapshot")
					}
				}
				results <- err
			}()
		}
		go func() {
			<-start
			_, err := s.Decide(ctx, staff.ID, r.ID, uuid.NewString(), "CHECKOUT", checkout)
			results <- err
		}()
		close(start)
		for n := 0; n < 21; n++ {
			require(t, <-results == nil, "coherent concurrent detail")
		}
		verify(v, 0, 0, 2)
	})

	t.Run("cancel_checkout_and_deny_cancel_races", func(t *testing.T) {
		for _, staffAction := range []string{"CHECKOUT", "DENIED"} {
			v := equipment(2)
			r, e := s.Submit(ctx, b1.ID, uuid.NewString(), input(v, 2))
			require(t, e == nil, "race request")
			var wins atomic.Int32
			start := make(chan struct{})
			var wg sync.WaitGroup
			for _, action := range []string{"CANCELLED", staffAction} {
				wg.Add(1)
				go func(action string) {
					defer wg.Done()
					<-start
					actor := b1.ID
					in := d.Decision{Confirm: true}
					if action != "CANCELLED" {
						actor = staff.ID
						if action == "CHECKOUT" {
							in = checkout
						} else {
							in.Reason = "Operational denial"
						}
					}
					_, e := s.Decide(ctx, actor, r.ID, uuid.NewString(), action, in)
					if e == nil {
						wins.Add(1)
					} else if !errors.Is(e, d.ErrState) {
						t.Error("race failure", e)
					}
				}(action)
			}
			close(start)
			wg.Wait()
			require(t, wins.Load() == 1, "one winning lifecycle edge")
			final, e := s.Read(ctx, admin.ID, r.ID)
			require(t, e == nil && len(final.Events) == 2, "one outcome audit")
			if final.Status == "CHECKED_OUT" {
				verify(v, 0, 0, 2)
			} else {
				verify(v, 2, 0, 0)
			}
		}
	})
	t.Run("direct_checkout_competes_with_reservation", func(t *testing.T) {
		v := equipment(1)
		var wins atomic.Int32
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			in := input(v, 1)
			in.BorrowerID = b1.ID
			in.DueAt = &due
			in.Handover = true
			_, e := s.Direct(ctx, staff.ID, uuid.NewString(), in)
			if e == nil {
				wins.Add(1)
			} else if !errors.Is(e, d.ErrStock) {
				t.Error(e)
			}
		}()
		go func() {
			defer wg.Done()
			_, e := s.Submit(ctx, b2.ID, uuid.NewString(), input(v, 1))
			if e == nil {
				wins.Add(1)
			} else if !errors.Is(e, d.ErrStock) {
				t.Error(e)
			}
		}()
		wg.Wait()
		require(t, wins.Load() == 1, "one available-unit owner")
		stock, e := ir.Get(ctx, v.ID)
		require(t, e == nil && stock.Stock.Reserved+stock.Stock.CheckedOut == 1, "single custody path")
		verify(v, 0, stock.Stock.Reserved, stock.Stock.CheckedOut)
	})
	t.Run("expiry_checkout_workers_restart_and_inactive_target", func(t *testing.T) {
		v := equipment(3)
		past := app.NewService(borrowingClock{br, time.Now().Add(-25 * time.Hour)}, ar, ir, ts, tx)
		r, e := past.Submit(ctx, b1.ID, uuid.NewString(), input(v, 2))
		require(t, e == nil, "persisted past-deadline clock fixture")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET is_active=false WHERE id=$1`, b1.ID)
		require(t, e == nil, "fixture deactivation")
		var wg sync.WaitGroup
		wg.Add(3)
		for i := 0; i < 2; i++ {
			go func() {
				defer wg.Done()
				fresh := app.NewService(br, ar, ir, ts, tx)
				_, e := fresh.Sweep(ctx, 100)
				if e != nil {
					t.Error(e)
				}
			}()
		}
		go func() {
			defer wg.Done()
			_, e := s.Decide(ctx, staff.ID, r.ID, uuid.NewString(), "CHECKOUT", checkout)
			if !errors.Is(e, d.ErrExpired) && !errors.Is(e, d.ErrState) {
				t.Error(e)
			}
		}()
		wg.Wait()
		final, e := s.Read(ctx, admin.ID, r.ID)
		require(t, e == nil && final.Status == "EXPIRED" && len(final.Events) == 2 && final.Events[1].ActorID == nil, "expiry winner durable system audit")
		verify(v, 3, 0, 0)
		n, e := s.Sweep(ctx, 100)
		require(t, e == nil && n == 0, "retry no double release")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET is_active=true WHERE id=$1`, b1.ID)
		require(t, e == nil, "restore synthetic active")
	})
	t.Run("issue_rollbacks_deactivation_and_role_revocation", func(t *testing.T) {
		v := equipment(3)
		r, e := s.Submit(ctx, b1.ID, uuid.NewString(), input(v, 2))
		require(t, e == nil, "pending rollback fixture")
		abort := app.NewService(br, ar, ir, ts, controlledTx{base: tx, after: func() error { return shared.ErrConflict }})
		_, e = abort.Decide(ctx, staff.ID, r.ID, uuid.NewString(), "CHECKOUT", checkout)
		require(t, e != nil, "issue interrupted")
		verify(v, 1, 2, 0)
		r, e = s.Read(ctx, admin.ID, r.ID)
		require(t, e == nil && r.Status == "PENDING" && len(r.Events) == 1, "no fabricated checkout audit")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET is_active=false WHERE id=$1`, b1.ID)
		require(t, e == nil, "deactivate target")
		_, e = s.Decide(ctx, staff.ID, r.ID, uuid.NewString(), "CHECKOUT", checkout)
		require(t, errors.Is(e, d.ErrEligibility), "inactive issue blocked")
		verify(v, 1, 2, 0)
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET is_active=true WHERE id=$1`, b1.ID)
		require(t, e == nil, "restore synthetic account")
		key := uuid.NewString()
		r, e = s.Decide(ctx, staff.ID, r.ID, key, "CHECKOUT", checkout)
		require(t, e == nil, "checkout")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET role='BORROWER' WHERE id=$1`, staff.ID)
		require(t, e == nil, "fixture staff role revoked")
		_, e = s.Decide(ctx, staff.ID, r.ID, key, "CHECKOUT", checkout)
		require(t, errors.Is(e, shared.ErrForbidden), "replay reauthorizes role")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET role='STAFF' WHERE id=$1`, staff.ID)
		require(t, e == nil, "restore fixture role")
		verify(v, 1, 0, 2)
	})
	t.Run("terms_revision_existing_binding_and_future_gates", func(t *testing.T) {
		v := equipment(4)
		r, e := s.Submit(ctx, b2.ID, uuid.NewString(), input(v, 1))
		require(t, e == nil, "original receipt")
		old := r.AcceptanceID
		next, e := ts.Publish(ctx, appterms.PublishCommand{ActorID: admin.ID, Version: "TEST-new-" + uuid.NewString()[:10], Title: "TEST revision", Body: "Synthetic isolated revision; not institutional content.", ExpectedCurrentVersionID: &version.ID})
		require(t, e == nil, "revision")
		r, e = s.Decide(ctx, staff.ID, r.ID, uuid.NewString(), "CHECKOUT", checkout)
		require(t, e == nil && r.AcceptanceID == old, "original pending consent retained")
		_, e = s.Submit(ctx, b2.ID, uuid.NewString(), input(v, 1))
		require(t, e != nil, "new request needs latest consent")
		in := input(v, 1)
		in.BorrowerID = b2.ID
		in.DueAt = &due
		in.Handover = true
		_, e = s.Direct(ctx, staff.ID, uuid.NewString(), in)
		require(t, e != nil, "direct requires target current consent")
		_, e = ts.Accept(ctx, b2.ID, next.ID)
		require(t, e == nil, "fixture revision acceptance")
		r, e = s.Direct(ctx, admin.ID, uuid.NewString(), in)
		require(t, e == nil && r.EntryPath == "DIRECT" && r.ExpiresAt == nil && r.Items[0].Reserved == 0, "direct nohold")
		verify(v, 2, 0, 2)
	})
	t.Run("constraints_immutable_history_and_data_rollback_refusal", func(t *testing.T) {
		for _, q := range []string{`DELETE FROM borrowings`, `UPDATE borrowings SET due_at=due_at+interval '1 hour' WHERE status='CHECKED_OUT'`, `TRUNCATE borrowing_items`, `UPDATE borrowing_events SET kind='DENIED'`, `DELETE FROM borrowing_operation_receipts`, `UPDATE equipment SET damaged_held=damaged_held+1`} {
			_, e := db.Pool.Exec(ctx, q)
			require(t, e != nil, "runtime history/custody guard")
		}
		m := database.NewMigrator(owner.Pool, "../../migrations")
		refused := false
		for n := 0; n < 4; n++ {
			if m.Down(ctx) != nil {
				refused = true
				break
			}
		}
		require(t, refused, "consequential down refusal after empty newer presentation pairs")
		require(t, m.Up(ctx) == nil, "restore only empty newer schemas after refusal")
	})
}

type borrowingEventFailure struct{ d.Repository }

func (r borrowingEventFailure) Event(context.Context, uuid.UUID, d.Event) error {
	return shared.ErrInternal
}
func TestRealBorrowingCrossFeatureSafety(t *testing.T) {
	ctx, db, owner, cfg := batchDatabase(t)
	tx := database.NewTxManager(db.Pool)
	ar := postgres.NewAccountsRepository(db.Pool)
	ir := postgres.NewInventoryRepository(db.Pool)
	br := postgres.NewBorrowingRepository(db.Pool)
	ts := appterms.NewService(postgres.NewTermsRepository(db.Pool), postgres.NewUserRepository(db.Pool), tx)
	s := app.NewService(br, ar, ir, ts, tx)
	is := appinventory.NewService(ir, ar, tx, catalogimage.Validator{})
	admin := termsUser(t, ctx, db, user.RoleAdmin)
	staff := termsUser(t, ctx, db, user.RoleStaff)
	mail := &capturedMail{}
	as := accountapp.NewService(ar, tx, security.NewBcryptHasher(4), mail, cfg.Accounts.Policy)
	as.SetObligationReader(postgres.NewBorrowingObligations(br))
	current, err := ts.Current(ctx, admin.ID)
	if errors.Is(err, termsdomain.ErrNotPublished) {
		current, err = ts.Publish(ctx, appterms.PublishCommand{ActorID: admin.ID, Version: "TEST-safe-" + uuid.NewString()[:10], Title: "TEST ONLY cross-feature", Body: "Synthetic isolated verification terms; not FSMO content."})
	}
	require(t, err == nil && current != nil, "isolated current terms fixture")
	student, e := as.Create(ctx, admin.ID, uuid.NewString(), accountdomain.Input{Name: "TEST Phase7 Student", Email: uuid.NewString() + "@students.example.invalid", BorrowerType: "STUDENT", StudentID: "00" + uuid.NewString()}, false)
	require(t, e == nil && student.ActivationRequired, "real account provisioning fixture")
	v, e := is.Create(ctx, staff.ID, uuid.NewString(), inventory.Create{Metadata: inventory.Metadata{Name: "TEST cross-feature " + uuid.NewString()}, Opening: 5, Reason: "Synthetic integration opening"})
	require(t, e == nil, "equipment fixture")
	in := d.Input{Confirm: true, Items: []d.Line{{EquipmentID: v.ID, Quantity: 2}}}
	t.Run("activation_and_missing_terms_fail_closed", func(t *testing.T) {
		_, e := s.Submit(ctx, student.ID, uuid.NewString(), in)
		require(t, e != nil, "activation pending")
		password := uuid.NewString() + uuid.NewString()
		require(t, as.Activate(ctx, mail.lastToken(), password) == nil, "secure single-use test mailbox activation")
		var pointer *uuid.UUID
		require(t, owner.Pool.QueryRow(ctx, `SELECT current_version_id FROM terms_publication WHERE id=1`).Scan(&pointer) == nil, "pointer")
		_, e = s.Submit(ctx, student.ID, uuid.NewString(), in)
		require(t, errors.Is(e, termsdomain.ErrAcceptanceRequired), "current acceptance required")
		_, e = ts.Accept(ctx, student.ID, *pointer)
		require(t, e == nil, "synthetic approved fixture consent")
		_, e = owner.Pool.Exec(ctx, `UPDATE terms_publication SET current_version_id=NULL WHERE id=1`)
		require(t, e == nil, "isolated unpublished scenario")
		defer func() {
			_, e := owner.Pool.Exec(ctx, `UPDATE terms_publication SET current_version_id=$1 WHERE id=1`, pointer)
			require(t, e == nil, "restore isolated pointer")
		}()
		_, e = s.Submit(ctx, student.ID, uuid.NewString(), in)
		require(t, errors.Is(e, termsdomain.ErrNotPublished), "unpublished submit blocked")
		future := time.Now().Add(time.Hour)
		direct := in
		direct.BorrowerID = student.ID
		direct.DueAt = &future
		direct.Handover = true
		_, e = s.Direct(ctx, staff.ID, uuid.NewString(), direct)
		require(t, errors.Is(e, termsdomain.ErrNotPublished), "unpublished direct blocked")
		stock, e := ir.Get(ctx, v.ID)
		require(t, e == nil && stock.Stock.Available == 5 && stock.Stock.Reserved == 0 && stock.Stock.CheckedOut == 0, "eligibility failures no physical effects")
	})
	t.Run("required_audit_failure_rolls_back_every_effect", func(t *testing.T) {
		abort := app.NewService(borrowingEventFailure{br}, ar, ir, ts, tx)
		var before, after int
		require(t, owner.Pool.QueryRow(ctx, `SELECT count(*) FROM borrowing_items WHERE equipment_id=$1`, v.ID).Scan(&before) == nil, "item count")
		_, e := abort.Submit(ctx, student.ID, uuid.NewString(), in)
		require(t, errors.Is(e, shared.ErrInternal), "audit failure")
		require(t, owner.Pool.QueryRow(ctx, `SELECT count(*) FROM borrowing_items WHERE equipment_id=$1`, v.ID).Scan(&after) == nil && before == after, "no orphan item/header")
		stock, e := ir.Get(ctx, v.ID)
		require(t, e == nil && stock.Stock.Available == 5 && stock.Sequence == 1, "ledger/stock rolled back")
	})
	t.Run("concurrent_same_key_exactly_one_reservation", func(t *testing.T) {
		key := uuid.NewString()
		var wg sync.WaitGroup
		ids := make(chan uuid.UUID, 12)
		for i := 0; i < 12; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				r, e := s.Submit(ctx, student.ID, key, in)
				if e != nil {
					t.Error(e)
				} else {
					ids <- r.ID
				}
			}()
		}
		wg.Wait()
		close(ids)
		var id uuid.UUID
		for v := range ids {
			if id == uuid.Nil {
				id = v
			} else {
				require(t, id == v, "one replay identity")
			}
		}
		stock, e := ir.Get(ctx, v.ID)
		require(t, e == nil && stock.Stock.Available == 3 && stock.Stock.Reserved == 2 && stock.Sequence == 2, "exactly one hold and movement")
	})
	t.Run("deactivation_warns_without_resolving_stock_or_loans", func(t *testing.T) {
		list, e := s.List(ctx, staff.ID, d.Filter{Owner: &student.ID, Page: 1, PerPage: 25})
		require(t, e == nil && list.Total == 1, "scoped operational count")
		pending := list.Items[0]
		due := time.Now().Add(time.Hour)
		loan, e := s.Decide(ctx, staff.ID, pending.ID, uuid.NewString(), "CHECKOUT", d.Decision{Confirm: true, Handover: true, DueAt: &due})
		require(t, e == nil, "physical issue")
		record, e := as.Detail(ctx, admin.ID, student.ID, false)
		require(t, e == nil && record.Obligations.Availability == "AVAILABLE" && *record.Obligations.ActiveBorrowings == 1 && *record.Obligations.UnreturnedUnits == 2 && record.Obligations.FineMinor != nil && *record.Obligations.FineMinor == 0 && record.Obligations.ReplacementUnits != nil && *record.Obligations.ReplacementUnits == 0, "installed accountability provides authoritative known fine/replacement zeros")
		record, e = as.Status(ctx, admin.ID, student.ID, uuid.NewString(), accountapp.StatusInput{Active: false, Confirm: true, ExpectedUpdatedAt: record.UpdatedAt}, false)
		require(t, e == nil && !record.IsActive, "deactivation does not veto obligations")
		again, e := s.Read(ctx, admin.ID, loan.ID)
		require(t, e == nil && again.Status == "CHECKED_OUT" && again.DueAt.Equal(*loan.DueAt) && again.Items[0].Issued == 2, "loan/time preserved")
		stock, e := ir.Get(ctx, v.ID)
		require(t, e == nil && stock.Stock.CheckedOut == 2 && stock.Stock.Available == 3, "custody preserved")
		_, e = is.Status(ctx, admin.ID, v.ID, uuid.NewString(), inventory.StatusInput{Status: "ARCHIVED", Confirm: true, ExpectedVersion: v.Version})
		require(t, errors.Is(e, shared.ErrConflict), "borrowed equipment archive blocked")
	})
	t.Run("clear_equipment_archives_but_future_liability_stays_closed", func(t *testing.T) {
		clear, e := is.Create(ctx, staff.ID, uuid.NewString(), inventory.Create{Metadata: inventory.Metadata{Name: "TEST clear archive " + uuid.NewString()}, Opening: 1, Reason: "Synthetic clear stock"})
		require(t, e == nil, "unrelated equipment")
		detail, e := is.Detail(ctx, staff.ID, clear.ID)
		require(t, e == nil && detail.ArchiveSafety == "CLEAR", "actual equipment scoped archive check")
		clear, e = is.Status(ctx, staff.ID, clear.ID, uuid.NewString(), inventory.StatusInput{Status: "ARCHIVED", Confirm: true, ExpectedVersion: clear.Version})
		require(t, e == nil && clear.Status == "ARCHIVED", "unrelated safe archive preserved")
	})
}
