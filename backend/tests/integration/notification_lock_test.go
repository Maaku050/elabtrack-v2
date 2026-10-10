package integration

import (
	"context"
	borrowapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	invapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	termsapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	inv "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	notification "github.com/Maaku050/elabtrack-v2/backend/internal/domain/notifications"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
	"testing"
	"time"
)

// Reproduce the exact FK/account/lifecycle interleaving found by the real
// process-restart gate, while also proving competing transitions stay exclusive.
func TestRealNotificationLifecycleLockOrder(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	txm := database.NewTxManager(db.Pool)
	ar, ir, br := postgres.NewAccountsRepository(db.Pool), postgres.NewInventoryRepository(db.Pool), postgres.NewBorrowingRepository(db.Pool)
	terms := termsapp.NewService(postgres.NewTermsRepository(db.Pool), postgres.NewUserRepository(db.Pool), txm)
	staff, borrower := termsUser(t, ctx, db, user.RoleStaff), termsUser(t, ctx, db, user.RoleBorrower)
	_, e := owner.Pool.Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type)VALUES($1,'FACULTY')`, borrower.ID)
	require(t, e == nil, "fictional borrower category")
	current, e := terms.Current(ctx, staff.ID)
	require(t, e == nil, "current isolated test terms")
	_, e = terms.Accept(ctx, borrower.ID, current.ID)
	require(t, e == nil, "explicit test consent")
	inventory := invapp.NewService(ir, ar, txm, catalogimage.Validator{})
	eq, e := inventory.Create(ctx, staff.ID, uuid.NewString(), inv.Create{Metadata: inv.Metadata{Name: "TEST notification lock " + uuid.NewString()}, Opening: 1, Reason: "Isolated FK/lifecycle contention"})
	require(t, e == nil, "one-unit stock")
	service := borrowapp.NewService(br, ar, ir, terms, txm)
	loan, e := service.Submit(ctx, borrower.ID, uuid.NewString(), d.Input{Confirm: true, Items: []d.Line{{EquipmentID: eq.ID, Quantity: 1}}})
	require(t, e == nil, "reserved fixture")
	lifecycle, e := db.Pool.Begin(ctx)
	require(t, e == nil, "lifecycle transaction")
	defer lifecycle.Rollback(context.Background())
	lifecycleCtx := database.WithTx(ctx, lifecycle)
	require(t, ar.LockAccounts(lifecycleCtx, []uuid.UUID{borrower.ID}) == nil, "account lock precedes lifecycle lock")
	fanout, e := db.Pool.Begin(ctx)
	require(t, e == nil, "notification transaction")
	defer fanout.Rollback(context.Background())
	fanoutPID, lifecyclePID := int32(fanout.Conn().PgConn().PID()), int32(lifecycle.Conn().PgConn().PID())
	emitted := make(chan error, 1)
	go func() {
		nr := postgres.NewNotificationsRepository(db.Pool)
		event := notification.Event{Key: "LOCK:" + uuid.NewString(), BorrowingID: loan.ID, BorrowerID: borrower.ID, Kind: "SUBMITTED", At: time.Now()}
		_, err := nr.Emit(database.WithTx(ctx, fanout), event, event.Message())
		emitted <- err
	}()
	waitBlocked := func(pid, blocker int32) {
		t.Helper()
		for n := 0; n < 100; n++ {
			var blocked bool
			err := owner.Pool.QueryRow(ctx, `SELECT pg_blocking_pids($1) @> ARRAY[$2]::int[]`, pid, blocker).Scan(&blocked)
			require(t, err == nil, "observe real lock wait")
			if blocked {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
		t.Fatal("expected exact lock wait was not observed")
	}
	waitBlocked(fanoutPID, lifecyclePID)
	locked, e := br.Get(lifecycleCtx, loan.ID, true)
	require(t, e == nil && locked.Status == "PENDING", "notification FK must not deadlock lifecycle authority")
	competing, e := db.Pool.Begin(ctx)
	require(t, e == nil, "competing lifecycle transaction")
	defer competing.Rollback(context.Background())
	competingPID := int32(competing.Conn().PgConn().PID())
	changed := make(chan error, 1)
	go func() { _, err := br.Get(database.WithTx(ctx, competing), loan.ID, true); changed <- err }()
	waitBlocked(competingPID, lifecyclePID)
	require(t, lifecycle.Commit(ctx) == nil, "release authoritative lifecycle lock")
	require(t, <-emitted == nil && fanout.Commit(ctx) == nil, "durable notification completes without deadlock")
	require(t, <-changed == nil && competing.Commit(ctx) == nil, "competing transition resumes only after winner releases")
	_, e = service.Decide(ctx, borrower.ID, loan.ID, uuid.NewString(), "CANCELLED", d.Decision{Confirm: true})
	require(t, e == nil, "normal authorized release remains valid")
	stock, e := ir.Get(ctx, eq.ID)
	require(t, e == nil && stock.Stock.Available == 1 && stock.Stock.Reserved == 0 && stock.Stock.Total == 1, "exact conserved stock after lock contention")
}
