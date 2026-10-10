package integration

import (
	borrowapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	invapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	termsapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	inv "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestRealNotificationProcessRestart(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	binary := os.Getenv("ELABTRACK_PHASE9_API_BINARY")
	if binary == "" {
		t.Skip("explicit isolated product binary")
	}
	require(t, os.Getenv("APP_ENV") == "test" && os.Getenv("APP_PORT") == "18085", "isolated child listener")
	tx := database.NewTxManager(db.Pool)
	ar := postgres.NewAccountsRepository(db.Pool)
	ir := postgres.NewInventoryRepository(db.Pool)
	br := postgres.NewBorrowingRepository(db.Pool)
	ts := termsapp.NewService(postgres.NewTermsRepository(db.Pool), postgres.NewUserRepository(db.Pool), tx)
	admin, staff, borrower := termsUser(t, ctx, db, user.RoleAdmin), termsUser(t, ctx, db, user.RoleStaff), termsUser(t, ctx, db, user.RoleBorrower)
	_, e := owner.Pool.Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type)VALUES($1,'FACULTY')`, borrower.ID)
	require(t, e == nil, "synthetic category")
	terms, e := ts.Current(ctx, admin.ID)
	require(t, e == nil, "synthetic terms")
	_, e = ts.Accept(ctx, borrower.ID, terms.ID)
	require(t, e == nil, "fixture consent")
	is := invapp.NewService(ir, ar, tx, catalogimage.Validator{})
	eq, e := is.Create(ctx, staff.ID, uuid.NewString(), inv.Create{Metadata: inv.Metadata{Name: "TEST notification restart " + uuid.NewString()}, Opening: 1, Reason: "Synthetic recovery stock"})
	require(t, e == nil, "equipment")
	past := time.Now().Add(-25 * time.Hour).UTC().Truncate(time.Microsecond)
	due := past.Add(time.Hour)
	bs := borrowapp.NewService(accountabilityClock{br, past}, ar, ir, ts, tx)
	loan, e := bs.Direct(ctx, staff.ID, uuid.NewString(), d.Input{BorrowerID: borrower.ID, Items: []d.Line{{EquipmentID: eq.ID, Quantity: 1}}, DueAt: &due, Confirm: true, Handover: true})
	require(t, e == nil, "persisted overdue fixture before process starts")
	start := func() *exec.Cmd {
		t.Helper()
		cmd := exec.Command(binary)
		cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
		require(t, cmd.Start() == nil, "start real isolated API")
		t.Cleanup(func() { _ = cmd.Process.Kill() })
		client := &http.Client{Timeout: time.Second}
		ready := false
		for n := 0; n < 100; n++ {
			res, err := client.Get("http://127.0.0.1:18085/api/v1/ready")
			if err == nil {
				ready = res.StatusCode == 200
				res.Body.Close()
				if ready {
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		require(t, ready, "API ready")
		return cmd
	}
	count := func() int {
		t.Helper()
		var n int
		require(t, owner.Pool.QueryRow(ctx, `SELECT count(*) FROM notifications WHERE recipient_id=$1 AND borrowing_id=$2`, borrower.ID, loan.ID).Scan(&n) == nil, "durable notification count")
		return n
	}
	first := start()
	for n := 0; n < 200 && count() < 3; n++ {
		time.Sleep(50 * time.Millisecond)
	}
	require(t, count() == 3, "actual worker catches direct/overdue/fine source after startup")
	_, e = owner.Pool.Exec(ctx, `UPDATE notifications SET read_at=clock_timestamp() WHERE recipient_id=$1 AND borrowing_id=$2 AND kind='OVERDUE'`, borrower.ID, loan.ID)
	require(t, e == nil, "fixture read state")
	require(t, first.Process.Kill() == nil, "abrupt own process stop")
	_ = first.Wait()
	second := start()
	time.Sleep(300 * time.Millisecond)
	require(t, count() == 3, "restart cannot duplicate visible events")
	var read int
	require(t, owner.Pool.QueryRow(ctx, `SELECT count(*) FROM notifications WHERE recipient_id=$1 AND borrowing_id=$2 AND read_at IS NOT NULL`, borrower.ID, loan.ID).Scan(&read) == nil && read == 1, "read state survives real process restart")
	require(t, second.Process.Signal(os.Interrupt) == nil && second.Wait() == nil, "graceful own shutdown")
}
