package integration

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	appinventory "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	appterms "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
)

// This opt-in launches the normal compiled API, never a test HTTP endpoint.
// batchDatabase refuses normal/local/production DB names and ports first.
func TestRealBorrowingProcessRestart(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	binary := os.Getenv("ELABTRACK_PHASE7_API_BINARY")
	if binary == "" {
		t.Skip("set ELABTRACK_PHASE7_API_BINARY to the isolated verification build")
	}
	require(t, os.Getenv("APP_ENV") == "test" && os.Getenv("APP_PORT") == "18085", "isolated child listener guard")
	tx := database.NewTxManager(db.Pool)
	ar, ir, br := postgres.NewAccountsRepository(db.Pool), postgres.NewInventoryRepository(db.Pool), postgres.NewBorrowingRepository(db.Pool)
	ts := appterms.NewService(postgres.NewTermsRepository(db.Pool), postgres.NewUserRepository(db.Pool), tx)
	admin, staff, borrower := termsUser(t, ctx, db, user.RoleAdmin), termsUser(t, ctx, db, user.RoleStaff), termsUser(t, ctx, db, user.RoleBorrower)
	_, e := owner.Pool.Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type) VALUES($1,'FACULTY')`, borrower.ID)
	require(t, e == nil, "synthetic target category")
	current, e := ts.Current(ctx, admin.ID)
	require(t, e == nil, "isolated current TEST policy")
	_, e = ts.Accept(ctx, borrower.ID, current.ID)
	require(t, e == nil, "synthetic acceptance")
	inv := appinventory.NewService(ir, ar, tx, catalogimage.Validator{})
	eq, e := inv.Create(ctx, staff.ID, uuid.NewString(), inventory.Create{Metadata: inventory.Metadata{Name: "TEST process restart " + uuid.NewString()}, Opening: 1, Reason: "Isolated restart verification"})
	require(t, e == nil, "synthetic one-unit stock")
	deadline := time.Now().Add(5 * time.Second).UTC().Truncate(time.Microsecond)
	// Only this fixture's submission clock is shortened; the real child worker
	// uses clock_timestamp() and the persisted 24-hour deadline unchanged.
	svc := app.NewService(borrowingClock{br, deadline.Add(-d.PendingTTL)}, ar, ir, ts, tx)
	pending, e := svc.Submit(ctx, borrower.ID, uuid.NewString(), d.Input{Confirm: true, Items: []d.Line{{EquipmentID: eq.ID, Quantity: 1}}})
	require(t, e == nil, "durable pending before process start")
	client := &http.Client{Timeout: time.Second}
	start := func() *exec.Cmd {
		t.Helper()
		cmd := exec.Command(binary)
		log, err := os.OpenFile(filepath.Join(t.TempDir(), "api.log"), os.O_CREATE|os.O_WRONLY, 0600)
		require(t, err == nil, "private isolated process diagnostic")
		cmd.Stdout, cmd.Stderr = log, log
		require(t, cmd.Start() == nil, "start isolated product API")
		t.Cleanup(func() {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			_ = log.Close()
			if t.Failed() {
				diagnostic, _ := os.ReadFile(log.Name())
				t.Logf("Isolated restart diagnostic: %s", diagnostic)
			}
		})
		ready := false
		for n := 0; n < 100; n++ {
			response, err := client.Get("http://127.0.0.1:18085/api/v1/ready")
			if err == nil {
				ready = response.StatusCode == 200
				response.Body.Close()
				if ready {
					break
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
		require(t, ready, "isolated API ready")
		return cmd
	}
	first := start()
	stored, e := br.Get(ctx, pending.ID, false)
	require(t, e == nil && stored.Status == "PENDING", "pending survived real first startup")
	require(t, first.Process.Kill() == nil, "abruptly terminate own API while pending")
	_ = first.Wait()
	for time.Now().Before(deadline.Add(50 * time.Millisecond)) {
		time.Sleep(50 * time.Millisecond)
	}
	second := start()
	for n := 0; n < 100; n++ {
		stored, e = br.Get(ctx, pending.ID, false)
		if e == nil && stored.Status == "EXPIRED" {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	require(t, e == nil && stored.Status == "EXPIRED" && stored.Items[0].Reserved == 0 && len(stored.Events) == 2, "actual restart catches durable overdue pending")
	stock, e := ir.Get(ctx, eq.ID)
	require(t, e == nil && stock.Stock.Available == 1 && stock.Stock.Reserved == 0 && stock.Stock.CheckedOut == 0 && stock.Stock.Total == 1, "exact one-unit release")
	require(t, second.Process.Signal(os.Interrupt) == nil && second.Wait() == nil, "graceful own API shutdown")
	third := start()
	var releases int
	require(t, owner.Pool.QueryRow(ctx, `SELECT count(*) FROM inventory_movements WHERE equipment_id=$1 AND kind='EXPIRED'`, eq.ID).Scan(&releases) == nil && releases == 1, "repeated real startup cannot release twice")
	require(t, third.Process.Signal(os.Interrupt) == nil && third.Wait() == nil, "final isolated API shutdown")
}
