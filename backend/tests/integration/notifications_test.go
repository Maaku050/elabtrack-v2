package integration

import (
	"errors"
	borrowapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	invapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/notifications"
	termsapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	borrow "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	inv "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/notifications"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
	"sync"
	"testing"
	"time"
)

func TestRealNotifications(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	tx := database.NewTxManager(db.Pool)
	ar := postgres.NewAccountsRepository(db.Pool)
	ir := postgres.NewInventoryRepository(db.Pool)
	br := postgres.NewBorrowingRepository(db.Pool)
	nr := postgres.NewNotificationsRepository(db.Pool)
	s := app.NewService(nr, ar, tx, 86400)
	admin, staff, b1, b2 := termsUser(t, ctx, db, user.RoleAdmin), termsUser(t, ctx, db, user.RoleStaff), termsUser(t, ctx, db, user.RoleBorrower), termsUser(t, ctx, db, user.RoleBorrower)
	ts := termsapp.NewService(postgres.NewTermsRepository(db.Pool), postgres.NewUserRepository(db.Pool), tx)
	v, e := ts.Current(ctx, admin.ID)
	require(t, e == nil, "test policy")
	for _, u := range []uuid.UUID{b1.ID, b2.ID} {
		_, e = owner.Pool.Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type)VALUES($1,'FACULTY')`, u)
		require(t, e == nil, "fictional profile")
		_, e = ts.Accept(ctx, u, v.ID)
		require(t, e == nil, "synthetic acceptance")
	}
	bs := borrowapp.NewService(br, ar, ir, ts, tx)
	is := invapp.NewService(ir, ar, tx, catalogimage.Validator{})
	eq, e := is.Create(ctx, staff.ID, uuid.NewString(), inv.Create{Metadata: inv.Metadata{Name: "TEST Notification Ladle " + uuid.NewString()}, Opening: 20, Reason: "Synthetic notification stock"})
	require(t, e == nil, "equipment")
	drain := func() {
		for n := 0; n < 100; n++ {
			count, e := s.Sweep(ctx, 100)
			require(t, e == nil, "durable reconciliation")
			if count == 0 {
				return
			}
		}
		t.Fatal("bounded catch-up did not finish")
	}
	drain()
	request, e := bs.Submit(ctx, b1.ID, uuid.NewString(), borrow.Input{Confirm: true, Items: []borrow.Line{{EquipmentID: eq.ID, Quantity: 2}}})
	require(t, e == nil, "request")
	t.Run("rollback_then_concurrent_dispatch_exactly_once", func(t *testing.T) {
		abort := app.NewService(nr, ar, controlledTx{base: tx, after: func() error { return shared.ErrConflict }}, 86400)
		_, e = abort.Sweep(ctx, 100)
		require(t, e != nil, "forced notification rollback")
		var n int
		e = owner.Pool.QueryRow(ctx, `SELECT count(*) FROM notification_dispatch WHERE borrowing_id=$1`, request.ID).Scan(&n)
		require(t, e == nil && n == 0, "dispatch intent and recipients both rolled back")
		var wg sync.WaitGroup
		for n := 0; n < 4; n++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := app.NewService(nr, ar, tx, 86400).Sweep(ctx, 100)
				if err != nil {
					t.Error(err)
				}
			}()
		}
		wg.Wait()
		drain()
		page, e := s.List(ctx, b1.ID, d.Filter{Page: 1, PerPage: 1})
		require(t, e == nil && page.Total == 1 && page.Unread == 1 && len(page.Items) == 1 && page.Items[0].Kind == "SUBMITTED", "one durable owner notification")
		none, e := s.List(ctx, b2.ID, d.Filter{Page: 1, PerPage: 25})
		require(t, e == nil && none.Total == 0, "other borrower receives no foreign event")
		var duplicates int
		e = owner.Pool.QueryRow(ctx, `SELECT count(*) FROM (SELECT event_key,recipient_id FROM notifications GROUP BY event_key,recipient_id HAVING count(*)>1)x`).Scan(&duplicates)
		require(t, e == nil && duplicates == 0, "independent recipient dedup")
		operations, e := s.List(ctx, staff.ID, d.Filter{Page: 1, PerPage: 100})
		require(t, e == nil && operations.Total > 0, "relevant staff action")
	})
	t.Run("ownership_read_unread_restart_and_live_role", func(t *testing.T) {
		page, e := s.List(ctx, b1.ID, d.Filter{Page: 1, PerPage: 25})
		require(t, e == nil, "own notification")
		id := page.Items[0].ID
		_, e = s.Mark(ctx, b2.ID, id, true)
		require(t, errors.Is(e, shared.ErrNotFound), "other user cannot mark notification")
		read, e := s.Mark(ctx, b1.ID, id, true)
		require(t, e == nil && read.ReadAt != nil, "persisted read")
		again, e := app.NewService(nr, ar, tx, 86400).Count(ctx, b1.ID)
		require(t, e == nil && again == 0, "reconstructed service keeps read state")
		read, e = s.Mark(ctx, b1.ID, id, false)
		require(t, e == nil && read.ReadAt == nil, "explicit unread")
		again, e = s.Count(ctx, b1.ID)
		require(t, e == nil && again == 1, "exact unread count")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET role='BORROWER' WHERE id=$1`, staff.ID)
		require(t, e == nil, "synthetic demotion")
		page, e = s.List(ctx, staff.ID, d.Filter{Page: 1, PerPage: 100})
		require(t, e == nil && page.Total == 0, "prior operation content hidden after role change")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET role='STAFF' WHERE id=$1`, staff.ID)
		require(t, e == nil, "restore role")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET is_active=false WHERE id=$1`, b2.ID)
		require(t, e == nil, "synthetic deactivation")
		_, e = s.Count(ctx, b2.ID)
		require(t, errors.Is(e, shared.ErrUnauthorized), "inactive account denied")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET is_active=true WHERE id=$1`, b2.ID)
		require(t, e == nil, "restore fixture")
	})
	t.Run("due_overdue_dispositions_fines_and_pagination", func(t *testing.T) {
		due := time.Now().Add(time.Hour).UTC().Truncate(time.Microsecond)
		loan, e := bs.Decide(ctx, staff.ID, request.ID, uuid.NewString(), "CHECKOUT", borrow.Decision{Confirm: true, Handover: true, DueAt: &due})
		require(t, e == nil, "checkout")
		drain()
		page, e := s.List(ctx, b1.ID, d.Filter{Page: 1, PerPage: 100})
		require(t, e == nil, "due reminders")
		kinds := map[string]int{}
		for _, n := range page.Items {
			kinds[n.Kind]++
		}
		require(t, kinds["DUE_SOON"] == 1 && kinds["CHECKED_OUT"] == 1, "one configured due reminder")
		past := time.Now().Add(-49 * time.Hour).UTC().Truncate(time.Microsecond)
		pastDue := past.Add(time.Hour)
		historical := borrowapp.NewService(accountabilityClock{br, past}, ar, ir, ts, tx)
		late, e := historical.Direct(ctx, staff.ID, uuid.NewString(), borrow.Input{BorrowerID: b1.ID, Items: []borrow.Line{{EquipmentID: eq.ID, Quantity: 2}}, DueAt: &pastDue, Confirm: true, Handover: true})
		require(t, e == nil, "synthetic actual overdue source")
		drain()
		late, e = bs.Return(ctx, staff.ID, late.ID, uuid.NewString(), borrow.ReturnInput{ExpectedEvents: len(late.Events), Lines: []borrow.ReturnLine{{ItemID: late.Items[0].ID, Damaged: 1, Lost: 1}}, Reason: "TEST damage and loss", Confirm: true})
		require(t, e == nil, "incident")
		drain()
		for _, o := range late.Obligations {
			late, e = bs.Replace(ctx, staff.ID, late.ID, uuid.NewString(), borrow.ReplacementInput{ExpectedEvents: len(late.Events), ObligationID: o.ID, Quantity: 1, Equivalent: true, Reason: "TEST physical equivalence", Confirm: true})
			require(t, e == nil, "replacement")
		}
		late, e = bs.ClearFine(ctx, admin.ID, late.ID, uuid.NewString(), borrow.FineInput{Method: "OTHER_RESOLUTION", ExpectedMinor: late.Fine.Outstanding, Confirm: true})
		require(t, e == nil, "fine resolution")
		drain()
		page, e = s.List(ctx, b1.ID, d.Filter{Page: 1, PerPage: 100})
		require(t, e == nil, "all lifecycle notifications")
		kinds = map[string]int{}
		for _, n := range page.Items {
			kinds[n.Kind]++
		}
		for _, kind := range []string{"OVERDUE", "FINE_ASSESSED", "RETURN", "DAMAGED", "LOST", "REPLACEMENT_REQUIRED", "REPLACEMENT", "COMPLETED", "FINE_CLEARED"} {
			require(t, kinds[kind] > 0, "required notification event "+kind)
		}
		count := page.Total
		first, e := s.List(ctx, b1.ID, d.Filter{Page: 1, PerPage: 2})
		require(t, e == nil && first.Total == count && len(first.Items) == 2, "server count/page")
		second, e := s.List(ctx, b1.ID, d.Filter{Page: 2, PerPage: 2})
		require(t, e == nil && second.Total == count && len(second.Items) == 2 && first.Items[0].ID != second.Items[0].ID, "stable nonoverlapping pages")
		var overdue int
		e = owner.Pool.QueryRow(ctx, `SELECT count(*) FROM notifications WHERE recipient_id=$1 AND borrowing_id=$2 AND kind='OVERDUE'`, b1.ID, late.ID).Scan(&overdue)
		require(t, e == nil && overdue == 1, "restart/catch-up produces no reminder repeats")
		_ = loan
	})

	t.Run("bulk_read_authorized_scope_retry_rollback_and_restart", func(t *testing.T) {
		mine, e := s.Count(ctx, b1.ID)
		require(t, e == nil && mine > 1, "multiple own unread rows")
		foreign, e := s.Count(ctx, staff.ID)
		require(t, e == nil, "foreign scope baseline")
		abort := app.NewService(nr, ar, controlledTx{base: tx, after: func() error { return shared.ErrConflict }}, 86400)
		n, e := abort.MarkAll(ctx, b1.ID)
		require(t, e != nil && n == 0, "bulk rollback")
		count, e := s.Count(ctx, b1.ID)
		require(t, e == nil && count == mine, "bulk rollback retains all unread")
		n, e = s.MarkAll(ctx, b1.ID)
		require(t, e == nil && n == int64(mine), "all own unread across pages marked")
		count, e = app.NewService(nr, ar, tx, 86400).Count(ctx, b1.ID)
		require(t, e == nil && count == 0, "durable after reconstructed service")
		n, e = s.MarkAll(ctx, b1.ID)
		require(t, e == nil && n == 0, "repeated bulk read is idempotent")
		count, e = s.Count(ctx, staff.ID)
		require(t, e == nil && count == foreign, "foreign rows untouched")
		_, e = owner.Pool.Exec(ctx, "UPDATE users SET role='BORROWER' WHERE id=$1", staff.ID)
		require(t, e == nil, "fixture demotion")
		n, e = s.MarkAll(ctx, staff.ID)
		require(t, e == nil && n == 0, "hidden operational notifications not mutated after demotion")
		_, e = owner.Pool.Exec(ctx, "UPDATE users SET role='STAFF' WHERE id=$1", staff.ID)
		require(t, e == nil, "restore role")
		count, e = s.Count(ctx, staff.ID)
		require(t, e == nil && count == foreign, "hidden rows still unread")
		_, e = owner.Pool.Exec(ctx, "UPDATE users SET is_active=false WHERE id=$1", b2.ID)
		require(t, e == nil, "explicit inactive fixture")
		_, e = s.MarkAll(ctx, b2.ID)
		require(t, errors.Is(e, shared.ErrUnauthorized), "inactive account cannot bulk read")
	})
	t.Run("least_privilege_immutable_content_and_history_safe_down", func(t *testing.T) {
		for _, sql := range []string{`DELETE FROM notifications`, `UPDATE notifications SET title='spoof'`, `TRUNCATE notification_dispatch`} {
			_, e = db.Pool.Exec(ctx, sql)
			require(t, e != nil, "runtime protected content/history")
		}
		m := database.NewMigrator(owner.Pool, "../../migrations")
		require(t, m.Down(ctx) == nil, "empty profile pair before notification history gate")
		require(t, m.Down(ctx) != nil, "history-bearing notification rollback refused")
		require(t, m.Up(ctx) == nil, "restore empty profile schema after notification gate")
	})
}
