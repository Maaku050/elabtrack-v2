package integration

import (
	"context"
	"encoding/json"
	"errors"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	invapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	termsapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	inv "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Preserve accountability ports while overriding only the synthetic fixture clock.
type accountabilityClock struct {
	*postgres.BorrowingRepository
	at time.Time
}

func (r accountabilityClock) Clock(context.Context) (time.Time, error) { return r.at, nil }

func TestRealAccountability(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	tx := database.NewTxManager(db.Pool)
	ar := postgres.NewAccountsRepository(db.Pool)
	ir := postgres.NewInventoryRepository(db.Pool)
	br := postgres.NewBorrowingRepository(db.Pool)
	ts := termsapp.NewService(postgres.NewTermsRepository(db.Pool), postgres.NewUserRepository(db.Pool), tx)
	admin, staff, borrower := termsUser(t, ctx, db, user.RoleAdmin), termsUser(t, ctx, db, user.RoleStaff), termsUser(t, ctx, db, user.RoleBorrower)
	_, e := owner.Pool.Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type) VALUES($1,'FACULTY')`, borrower.ID)
	require(t, e == nil, "fictional faculty")
	current, e := ts.Current(ctx, admin.ID)
	require(t, e == nil, "published synthetic terms")
	_, e = ts.Accept(ctx, borrower.ID, current.ID)
	require(t, e == nil, "explicit synthetic consent")
	inventory := invapp.NewService(ir, ar, tx, catalogimage.Validator{})
	start := time.Now().UTC().Truncate(time.Microsecond)
	due := start.Add(time.Hour)
	serviceAt := func(at time.Time) *app.Service { return app.NewService(accountabilityClock{br, at}, ar, ir, ts, tx) }
	issue := func(q int64) (d.Record, inv.Equipment) {
		eq, e := inventory.Create(ctx, staff.ID, uuid.NewString(), inv.Create{Metadata: inv.Metadata{Name: "TEST Accountability " + uuid.NewString()}, Opening: q + 5, Reason: "Synthetic acquisition"})
		require(t, e == nil, "equipment")
		loan, e := serviceAt(start).Direct(ctx, staff.ID, uuid.NewString(), d.Input{BorrowerID: borrower.ID, Confirm: true, Items: []d.Line{{EquipmentID: eq.ID, Quantity: q}}, Handover: true, DueAt: &due})
		require(t, e == nil, "physical issue")
		return loan, eq
	}
	verify := func(eq inv.Equipment, want inv.Stock) {
		got, e := ir.Get(ctx, eq.ID)
		require(t, e == nil && got.Stock == want, "exact custody buckets")
		var ledger inv.Stock
		e = owner.Pool.QueryRow(ctx, `SELECT sum(delta_available),sum(delta_reserved),sum(delta_checked_out),sum(delta_damaged_held),sum(delta_total) FROM inventory_movements WHERE equipment_id=$1`, eq.ID).Scan(&ledger.Available, &ledger.Reserved, &ledger.CheckedOut, &ledger.DamagedHeld, &ledger.Total)
		require(t, e == nil && ledger == want, "ledger reconciles every custody bucket")
	}
	ret := func(v d.Record, g, damage, lost int64) d.ReturnInput {
		return d.ReturnInput{ExpectedEvents: len(v.Events), Confirm: true, Reason: "In-person verified disposition", Lines: []d.ReturnLine{{ItemID: v.Items[0].ID, Good: g, Damaged: damage, Lost: lost}}}
	}
	t.Run("bounded_public_history_keeps_complete_custody_and_review_authority", func(t *testing.T) {
		loan, eq := issue(31)
		live := serviceAt(start.Add(30 * time.Minute))
		for n := 0; n < 30; n++ {
			var e error
			loan, e = live.Return(ctx, staff.ID, loan.ID, uuid.NewString(), ret(loan, 1, 0, 0))
			require(t, e == nil, "many actual partial returns")
		}
		view, e := live.Read(ctx, borrower.ID, loan.ID)
		require(t, e == nil && view.PagedHistory && len(view.Events) == 25 && len(view.Returns) == 25 && view.EventCount == 31 && view.Items[0].Good == 30 && view.Items[0].Outstanding == 1, "bounded previews preserve full aggregate")
		first, e := live.History(ctx, borrower.ID, loan.ID, "returns", 1, 25)
		require(t, e == nil && first.Total == 30 && len(first.Items) == 25, "history page and exact total")
		second, e := live.History(ctx, borrower.ID, loan.ID, "returns", 2, 25)
		require(t, e == nil && second.Total == 30 && len(second.Items) == 5, "remaining history retained")
		seen := map[string]bool{}
		for _, raw := range append(first.Items, second.Items...) {
			var row map[string]any
			require(t, json.Unmarshal(raw, &row) == nil, "history row")
			id := row["id"].(string)
			require(t, !seen[id], "no duplicate history rows")
			seen[id] = true
		}
		other := termsUser(t, ctx, db, user.RoleBorrower)
		_, e = live.History(ctx, other.ID, loan.ID, "events", 1, 25)
		require(t, errors.Is(e, shared.ErrNotFound), "history ownership")
		_, e = live.History(ctx, borrower.ID, loan.ID, "users;DROP TABLE users", 1, 25)
		require(t, errors.Is(e, shared.ErrInvalidInput), "history source whitelist")
		_, e = live.History(ctx, borrower.ID, loan.ID, "events", 1, 101)
		require(t, errors.Is(e, shared.ErrInvalidInput), "history size bound")
		in := ret(view, 1, 0, 0)
		_, e = live.Return(ctx, staff.ID, loan.ID, uuid.NewString(), in)
		require(t, errors.Is(e, d.ErrState), "truncated-array length cannot grant review authority")
		in.ExpectedEvents = view.EventCount
		loan, e = live.Return(ctx, staff.ID, loan.ID, uuid.NewString(), in)
		require(t, e == nil && loan.Status == "COMPLETED", "complete-count review succeeds")
		verify(eq, inv.Stock{Available: 36, Total: 36})
	})
	t.Run("bounded_clearance_preview_uses_full_fine_sum", func(t *testing.T) {
		eq, e := inventory.Create(ctx, staff.ID, uuid.NewString(), inv.Create{Metadata: inv.Metadata{Name: "TEST history fine " + uuid.NewString()}, Opening: 1, Reason: "Synthetic history aggregation"})
		require(t, e == nil, "fine equipment")
		issued := start.Add(-40 * 24 * time.Hour)
		deadline := issued.Add(time.Hour)
		loan, e := serviceAt(issued).Direct(ctx, staff.ID, uuid.NewString(), d.Input{BorrowerID: borrower.ID, Confirm: true, Handover: true, DueAt: &deadline, Items: []d.Line{{EquipmentID: eq.ID, Quantity: 1}}})
		require(t, e == nil, "historical fixture issue")
		for day := 1; day <= 30; day++ {
			v, e := serviceAt(deadline.Add(time.Duration(day)*24*time.Hour)).Read(ctx, admin.ID, loan.ID)
			require(t, e == nil && v.Fine.Outstanding == 1000, "next started-day balance")
			_, e = serviceAt(deadline.Add(time.Duration(day)*24*time.Hour)).ClearFine(ctx, admin.ID, loan.ID, uuid.NewString(), d.FineInput{Method: "OTHER_RESOLUTION", Note: "Fictional full-history aggregate check", ExpectedMinor: 1000, Confirm: true})
			require(t, e == nil, "actual historical Admin clearance")
		}
		v, e := serviceAt(deadline.Add(31*24*time.Hour)).Read(ctx, borrower.ID, loan.ID)
		require(t, e == nil && len(v.Clearances) == 25 && v.Fine.Cleared == 30000 && v.Fine.Assessed == 31000 && v.Fine.Outstanding == 1000, "full financial aggregate independent of preview")
		page, e := serviceAt(start).History(ctx, borrower.ID, loan.ID, "clearances", 2, 25)
		require(t, e == nil && page.Total == 30 && len(page.Items) == 5, "all fine classifications still browsable")
	})
	t.Run("mixed_partial_replacement_completion_and_frozen_fine", func(t *testing.T) {
		loan, eq := issue(5)
		late := serviceAt(due.Add(25 * time.Hour))
		key := uuid.NewString()
		in := ret(loan, 2, 1, 1)
		loan, e = late.Return(ctx, staff.ID, loan.ID, key, in)
		if e != nil {
			t.Fatalf("mixed return: %v", e)
		}
		t.Logf("mixed status=%s outstanding=%d obligations=%d fine=%d", loan.Status, loan.Items[0].Outstanding, len(loan.Obligations), loan.Fine.Outstanding)
		require(t, e == nil && loan.Status == "CHECKED_OUT" && loan.Items[0].Outstanding == 1 && len(loan.Obligations) == 2 && loan.Fine.Outstanding == 2000, "mixed partial return remains open")
		verify(eq, inv.Stock{Available: 7, CheckedOut: 1, DamagedHeld: 1, Total: 9})
		again, e := late.Return(ctx, staff.ID, loan.ID, key, in)
		require(t, e == nil && again.ID == loan.ID && len(again.Returns) == 1, "same key changes stock once")
		_, e = late.Return(ctx, staff.ID, loan.ID, uuid.NewString(), in)
		require(t, errors.Is(e, d.ErrState), "stale review rejected")
		_, e = late.Return(ctx, borrower.ID, loan.ID, uuid.NewString(), ret(loan, 1, 0, 0))
		require(t, errors.Is(e, shared.ErrForbidden), "borrower cannot record own return")
		loan, e = late.Return(ctx, staff.ID, loan.ID, uuid.NewString(), ret(loan, 1, 0, 0))
		require(t, e == nil && loan.Status == "CHECKED_OUT" && loan.Items[0].Outstanding == 0, "physical zero does not erase replacements")
		verify(eq, inv.Stock{Available: 8, DamagedHeld: 1, Total: 9})
		safety, e := ir.ArchiveSafety(ctx, eq.ID)
		require(t, e == nil && safety == "BLOCKED", "replacement liability blocks archive")
		for _, obligation := range loan.Obligations {
			in := d.ReplacementInput{ExpectedEvents: len(loan.Events), ObligationID: obligation.ID, Quantity: 1, Reason: "Equivalent replacement inspected and accepted in person", Equivalent: true, Confirm: true}
			loan, e = late.Replace(ctx, staff.ID, loan.ID, uuid.NewString(), in)
			require(t, e == nil, "physical replacement")
		}
		require(t, loan.Status == "COMPLETED" && loan.CompletedAt != nil && loan.FinalMinor != nil && *loan.FinalMinor == 2000, "all obligations complete and fine freezes")
		verify(eq, inv.Stock{Available: 10, DamagedHeld: 1, Total: 11})
		future, e := serviceAt(due.Add(500*time.Hour)).Read(ctx, borrower.ID, loan.ID)
		require(t, e == nil && future.Fine.Assessed == 2000 && future.Fine.Final, "read after restart cannot grow final fine")
		clear := d.FineInput{Method: "WAIVED", Note: "Synthetic authorized full resolution", ExpectedMinor: 2000, Confirm: true}
		key = uuid.NewString()
		_, e = late.ClearFine(ctx, staff.ID, loan.ID, key, clear)
		require(t, errors.Is(e, shared.ErrForbidden), "staff cannot clear fines")
		loan, e = late.ClearFine(ctx, admin.ID, loan.ID, key, clear)
		require(t, e == nil && loan.Fine.Outstanding == 0 && len(loan.Clearances) == 1, "Admin full balance immutable resolution")
		_, e = late.ClearFine(ctx, admin.ID, loan.ID, key, clear)
		require(t, e == nil, "clear replay")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET role='STAFF' WHERE id=$1`, admin.ID)
		require(t, e == nil, "synthetic demotion")
		_, e = late.ClearFine(ctx, admin.ID, loan.ID, key, clear)
		require(t, errors.Is(e, shared.ErrForbidden), "demoted Admin cannot replay privileged command")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET role='ADMIN' WHERE id=$1`, admin.ID)
		require(t, e == nil, "restore fixture authority")
		verify(eq, inv.Stock{Available: 10, DamagedHeld: 1, Total: 11})
	})
	t.Run("racing_review_once_and_transaction_rollback", func(t *testing.T) {
		loan, eq := issue(2)
		s := serviceAt(start.Add(time.Minute))
		in := ret(loan, 1, 0, 0)
		abort := app.NewService(accountabilityClock{br, start.Add(time.Minute)}, ar, ir, ts, controlledTx{base: tx, after: func() error { return shared.ErrConflict }})
		_, e = abort.Return(ctx, staff.ID, loan.ID, uuid.NewString(), in)
		require(t, e != nil, "forced transaction failure")
		verify(eq, inv.Stock{Available: 5, CheckedOut: 2, Total: 7})
		check, e := s.Read(ctx, staff.ID, loan.ID)
		require(t, e == nil && len(check.Returns) == 0 && len(check.Events) == len(loan.Events), "audit and disposition rolled back")
		var wins atomic.Int32
		var wg sync.WaitGroup
		for n := 0; n < 2; n++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := s.Return(ctx, staff.ID, loan.ID, uuid.NewString(), in)
				if err == nil {
					wins.Add(1)
				} else if !errors.Is(err, d.ErrState) {
					t.Error(err)
				}
			}()
		}
		wg.Wait()
		require(t, wins.Load() == 1, "one fresh review winner")
		verify(eq, inv.Stock{Available: 6, CheckedOut: 1, Total: 7})
		check, e = s.Read(ctx, staff.ID, loan.ID)
		require(t, e == nil, "reload")
		for _, bad := range []d.ReturnInput{ret(check, 2, 0, 0), ret(check, -1, 0, 0), ret(check, 0, 0, 0), {Confirm: true, ExpectedEvents: len(check.Events), Lines: []d.ReturnLine{{ItemID: check.Items[0].ID, Good: 1}}}} {
			_, e = s.Return(ctx, staff.ID, loan.ID, uuid.NewString(), bad)
			require(t, e != nil, "invalid or excess return rejected")
		}
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET is_active=false WHERE id=$1`, borrower.ID)
		require(t, e == nil, "deactivate synthetic borrower")
		check, e = s.Return(ctx, staff.ID, loan.ID, uuid.NewString(), ret(check, 1, 0, 0))
		require(t, e == nil && check.Status == "COMPLETED", "deactivation never prevents operational resolution")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET is_active=true WHERE id=$1`, borrower.ID)
		require(t, e == nil, "restore fixture")
		verify(eq, inv.Stock{Available: 7, Total: 7})
	})
	t.Run("open_fine_clearance_then_accrual_and_stale_balance", func(t *testing.T) {
		loan, eq := issue(1)
		first := serviceAt(due.Add(time.Minute))
		clear := d.FineInput{Method: "PAID", ExpectedMinor: 1000, Confirm: true}
		loan, e = first.ClearFine(ctx, admin.ID, loan.ID, uuid.NewString(), clear)
		require(t, e == nil && loan.Fine.Outstanding == 0 && loan.Status == "CHECKED_OUT", "clear does not return equipment or close loan")
		verify(eq, inv.Stock{Available: 5, CheckedOut: 1, Total: 6})
		next := serviceAt(due.Add(49 * time.Hour))
		read, e := next.Read(ctx, admin.ID, loan.ID)
		require(t, e == nil && read.Fine.Assessed == 3000 && read.Fine.Outstanding == 2000, "open loan continues started-day accrual")
		_, e = next.ClearFine(ctx, admin.ID, loan.ID, uuid.NewString(), clear)
		require(t, errors.Is(e, d.ErrState), "stale fine review fails")
		clear.ExpectedMinor = 2000
		clear.Method = "OTHER_RESOLUTION"
		loan, e = next.ClearFine(ctx, admin.ID, loan.ID, uuid.NewString(), clear)
		require(t, e == nil && loan.Fine.Outstanding == 0 && len(loan.Clearances) == 2, "full remaining balance separate audit")
		loan, e = next.Return(ctx, staff.ID, loan.ID, uuid.NewString(), ret(loan, 1, 0, 0))
		if e != nil {
			t.Fatalf("completion: %v", e)
		}
		require(t, e == nil && loan.Status == "COMPLETED" && loan.Fine.Outstanding == 0, "completion respects existing clearances")
	})
	t.Run("replacement_race_rollback_and_foreign_obligation", func(t *testing.T) {
		loan, eq := issue(2)
		s := serviceAt(start.Add(time.Minute))
		loan, e = s.Return(ctx, staff.ID, loan.ID, uuid.NewString(), ret(loan, 0, 2, 0))
		require(t, e == nil, "damage received")
		in := d.ReplacementInput{ExpectedEvents: len(loan.Events), ObligationID: loan.Obligations[0].ID, Quantity: 1, Reason: "Equivalent received", Equivalent: true, Confirm: true}
		bad := in
		bad.ObligationID = uuid.New()
		_, e = s.Replace(ctx, staff.ID, loan.ID, uuid.NewString(), bad)
		require(t, errors.Is(e, shared.ErrInvalidInput), "foreign obligation rejected")
		bad = in
		bad.Equivalent = false
		_, e = s.Replace(ctx, staff.ID, loan.ID, uuid.NewString(), bad)
		require(t, errors.Is(e, shared.ErrInvalidInput), "equivalence required")
		abort := app.NewService(accountabilityClock{br, start.Add(time.Minute)}, ar, ir, ts, controlledTx{base: tx, after: func() error { return shared.ErrConflict }})
		_, e = abort.Replace(ctx, staff.ID, loan.ID, uuid.NewString(), in)
		require(t, e != nil, "replacement rollback")
		verify(eq, inv.Stock{Available: 5, DamagedHeld: 2, Total: 7})
		var wins atomic.Int32
		var wg sync.WaitGroup
		for _, actor := range []uuid.UUID{staff.ID, admin.ID} {
			wg.Add(1)
			go func(actor uuid.UUID) {
				defer wg.Done()
				_, err := s.Replace(ctx, actor, loan.ID, uuid.NewString(), in)
				if err == nil {
					wins.Add(1)
				} else if !errors.Is(err, d.ErrState) {
					t.Error(err)
				}
			}(actor)
		}
		wg.Wait()
		require(t, wins.Load() == 1, "two actors one reviewed replacement winner")
		verify(eq, inv.Stock{Available: 6, DamagedHeld: 2, Total: 8})
		loan, e = s.Read(ctx, staff.ID, loan.ID)
		require(t, e == nil && len(loan.Replacements) == 1 && loan.Obligations[0].Accepted == 1, "exact acceptance history")
		in.ExpectedEvents = len(loan.Events)
		loan, e = s.Replace(ctx, admin.ID, loan.ID, uuid.NewString(), in)
		require(t, e == nil && loan.Status == "COMPLETED", "remaining replacement completes")
		verify(eq, inv.Stock{Available: 7, DamagedHeld: 2, Total: 9})
	})
	t.Run("final_return_replacement_race", func(t *testing.T) {
		loan, eq := issue(2)
		s := serviceAt(start.Add(time.Minute))
		loan, e = s.Return(ctx, staff.ID, loan.ID, uuid.NewString(), ret(loan, 0, 1, 0))
		require(t, e == nil, "partial damage")
		ri := ret(loan, 1, 0, 0)
		replacement := d.ReplacementInput{ExpectedEvents: len(loan.Events), ObligationID: loan.Obligations[0].ID, Quantity: 1, Reason: "Inspected equivalent", Equivalent: true, Confirm: true}
		results := make(chan error, 2)
		go func() { _, err := s.Return(ctx, staff.ID, loan.ID, uuid.NewString(), ri); results <- err }()
		go func() { _, err := s.Replace(ctx, admin.ID, loan.ID, uuid.NewString(), replacement); results <- err }()
		wins := 0
		for n := 0; n < 2; n++ {
			err := <-results
			if err == nil {
				wins++
			} else {
				require(t, errors.Is(err, d.ErrState), "other operation requires fresh review")
			}
		}
		require(t, wins == 1, "one review basis consumed")
		loan, e = s.Read(ctx, staff.ID, loan.ID)
		require(t, e == nil && loan.Status == "CHECKED_OUT", "partial combined resolution remains open")
		if loan.Items[0].Outstanding > 0 {
			loan, e = s.Return(ctx, staff.ID, loan.ID, uuid.NewString(), ret(loan, 1, 0, 0))
		} else {
			replacement.ExpectedEvents = len(loan.Events)
			loan, e = s.Replace(ctx, admin.ID, loan.ID, uuid.NewString(), replacement)
		}
		require(t, e == nil && loan.Status == "COMPLETED", "fresh remaining operation completes once")
		verify(eq, inv.Stock{Available: 7, DamagedHeld: 1, Total: 8})
	})
	t.Run("fine_clearance_race_and_rollback", func(t *testing.T) {
		loan, eq := issue(1)
		late := serviceAt(due.Add(time.Minute))
		in := d.FineInput{Method: "WAIVED", ExpectedMinor: 1000, Confirm: true}
		abort := app.NewService(accountabilityClock{br, due.Add(time.Minute)}, ar, ir, ts, controlledTx{base: tx, after: func() error { return shared.ErrConflict }})
		_, e = abort.ClearFine(ctx, admin.ID, loan.ID, uuid.NewString(), in)
		require(t, e != nil, "fine rollback")
		read, e := late.Read(ctx, admin.ID, loan.ID)
		require(t, e == nil && len(read.Clearances) == 0 && read.Fine.Outstanding == 1000, "fine audit rolled back")
		var wins atomic.Int32
		var wg sync.WaitGroup
		for n := 0; n < 2; n++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := late.ClearFine(ctx, admin.ID, loan.ID, uuid.NewString(), in)
				if err == nil {
					wins.Add(1)
				} else if !errors.Is(err, d.ErrState) {
					t.Error(err)
				}
			}()
		}
		wg.Wait()
		require(t, wins.Load() == 1, "one full-clear winner")
		verify(eq, inv.Stock{Available: 5, CheckedOut: 1, Total: 6})
	})
	t.Run("runtime_and_owner_history_immutable", func(t *testing.T) {
		for _, table := range []string{"return_lines", "replacement_obligations", "replacement_acceptances", "fine_clearances"} {
			for _, sql := range []string{"DELETE FROM " + table, "TRUNCATE " + table} {
				_, err := db.Pool.Exec(ctx, sql)
				require(t, err != nil, "runtime cannot erase accountability")
				_, err = owner.Pool.Exec(ctx, sql)
				require(t, err != nil, "owner trigger preserves history")
			}
		}
		var mismatch int
		e := owner.Pool.QueryRow(ctx, `SELECT count(*) FROM inventory_movements m LEFT JOIN return_lines l ON l.id=m.return_line_id LEFT JOIN replacement_acceptances a ON a.id=m.replacement_acceptance_id WHERE (m.kind='RETURN' AND (l.id IS NULL OR l.item_id<>m.borrowing_item_id OR l.equipment_id<>m.equipment_id)) OR (m.kind='REPLACEMENT' AND (a.id IS NULL OR a.item_id<>m.borrowing_item_id OR a.equipment_id<>m.equipment_id))`).Scan(&mismatch)
		require(t, e == nil && mismatch == 0, "ledger source identity matches immutable incident")
	})
}
