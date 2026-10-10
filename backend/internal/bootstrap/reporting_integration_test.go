package bootstrap

import (
	"context"
	"encoding/json"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	b "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPhase10HTTP(t *testing.T) {
	infra := batchInfra(t)
	c := buildContainer(infra)
	users, _, eq := phase7Fixtures(t, infra, c)
	limitMarker := "TEST_HTTP_REPORT_" + uuid.NewString()
	if _, e := infra.DB.Pool.Exec(context.Background(), `INSERT INTO account_audit_events(id,actor_id,account_id,action) SELECT gen_random_uuid(),$1,$2,$3 FROM generate_series(1,5001)`, users["admin"].ID, users["borrower"].ID, limitMarker); e != nil {
		t.Fatal("synthetic audit export-limit fixture")
	}

	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	_, issuer := newSecurityAdapters(infra.Config)
	tokens := map[string]string{}
	for role, u := range users {
		v, e := issuer.IssueAccessToken(context.Background(), u)
		if e != nil {
			t.Fatal(e)
		}
		tokens[role] = v
	}
	request := func(role, path string) (int, string, string) {
		t.Helper()
		r := httptest.NewRequest("GET", "/api/v1/reporting/"+path, nil)
		if role != "" {
			r.Header.Set("Authorization", "Bearer "+tokens[role])
		}
		res, e := server.Test(r, fiber.TestConfig{Timeout: 15 * time.Second})
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		raw, _ := io.ReadAll(res.Body)
		return res.StatusCode, string(raw), res.Header.Get("Content-Type")
	}
	for _, tc := range []struct {
		role, path string
		want       int
	}{{"", "dashboard", 401}, {"borrower", "dashboard", 200}, {"borrower", "inventory", 403}, {"staff", "dashboard", 200}, {"admin", "dashboard", 200}, {"staff", "account-audit", 403}, {"admin", "account-audit", 200}, {"admin", "invented", 400}, {"admin", "inventory?page=0", 400}, {"admin", "inventory?per_page=101", 400}, {"admin", "inventory?from=invalid", 400}, {"admin", "inventory?equipment_id=invalid", 400}, {"admin", "inventory?from=2026-10-10T00:00:00Z&to=2026-10-09T00:00:00Z", 400}, {"admin", "account-audit?borrower_id=" + uuid.NewString(), 400}, {"staff", "fine-clearances", 200}, {"admin", "account-audit.csv?search=" + limitMarker, 409}, {"staff", "kinds", 200}, {"staff", "dashboard?days=30", 200}, {"staff", "dashboard?days=365", 400}, {"staff", "active?due_today=true", 200}, {"staff", "inventory?due_today=true", 400}, {"staff", "active?due_today=invalid", 400}, {"staff", "inventory/export-data", 200}, {"borrower", "inventory/export-data", 403}, {"staff", "account-audit/export-data", 403}, {"admin", "account-audit/export-data?search=" + limitMarker, 409}, {"", "inventory/export-data", 401}} {
		got, raw, _ := request(tc.role, tc.path)
		if got != tc.want {
			t.Fatalf("%s %s: %d expected %d", tc.role, tc.path, got, tc.want)
		}
		if !json.Valid([]byte(raw)) {
			t.Fatal("envelope")
		}
	}
	code, raw, content := request("staff", "inventory.csv?equipment_id="+eq.ID.String())
	if code != 200 || !strings.HasPrefix(content, "text/csv") || !strings.Contains(raw, eq.Name) || !strings.Contains(raw, "Available,Reserved") {
		t.Fatalf("filtered CSV route: %d %s", code, content)
	}
	code, _, _ = request("borrower", "inventory.csv")
	if code != 403 {
		t.Fatal("CSV scope")
	}
	code, _, _ = request("staff", "account-audit.csv")
	if code != 403 {
		t.Fatal("CSV admin audit")
	}
}
func TestPhase10BrowserServer(t *testing.T) {
	if os.Getenv("ELABTRACK_PHASE10_BROWSER") != "1" {
		t.Skip("explicit isolated browser gate")
	}
	infra := batchInfra(t)
	c := buildContainer(infra)
	users, passwords, eq := phase7Fixtures(t, infra, c)
	ctx := context.Background()
	past := time.Now().Add(-49 * time.Hour).UTC().Truncate(time.Microsecond)
	due := past.Add(time.Hour)
	bs := app.NewService(phase8Clock{postgres.NewBorrowingRepository(infra.DB.Pool), past}, postgres.NewAccountsRepository(infra.DB.Pool), postgres.NewInventoryRepository(infra.DB.Pool), c.TermsSvc, infra.Tx)
	loan, e := bs.Direct(ctx, users["staff"].ID, uuid.NewString(), b.Input{BorrowerID: users["borrower"].ID, Items: []b.Line{{EquipmentID: eq.ID, Quantity: 3}}, DueAt: &due, Confirm: true, Handover: true})
	if e != nil {
		t.Fatal(e)
	}
	closed, e := bs.Direct(ctx, users["staff"].ID, uuid.NewString(), b.Input{BorrowerID: users["borrower"].ID, Items: []b.Line{{EquipmentID: eq.ID, Quantity: 1}}, DueAt: &due, Confirm: true, Handover: true})
	if e != nil {
		t.Fatal(e)
	}
	closed, e = c.BorrowingSvc.Return(ctx, users["staff"].ID, closed.ID, uuid.NewString(), b.ReturnInput{ExpectedEvents: len(closed.Events), Confirm: true, Reason: "TEST physically returned complete example", Lines: []b.ReturnLine{{ItemID: closed.Items[0].ID, Good: 1}}})
	if e != nil {
		t.Fatal(e)
	}
	_, e = c.BorrowingSvc.ClearFine(ctx, users["admin"].ID, closed.ID, uuid.NewString(), b.FineInput{ExpectedMinor: closed.Fine.Outstanding, Confirm: true, Method: "WAIVED", Note: "TEST fictional full-balance Admin resolution"})
	if e != nil {
		t.Fatal(e)
	}
	fixtures := map[string]any{"equipment": eq, "loan": loan}
	for role, u := range users {
		fixtures[role] = map[string]any{"email": u.Email, "password": passwords[role]}
	}
	f, e := os.OpenFile("/tmp/elabtrack-phase10-fixtures.json", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		t.Fatal("private fixture exists")
	}
	raw, _ := json.Marshal(fixtures)
	_, e = f.Write(raw)
	f.Close()
	if e != nil {
		t.Fatal(e)
	}
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	go func() { _ = server.Listen("127.0.0.1:18085", fiber.ListenConfig{DisableStartupMessage: true}) }()
	defer server.Shutdown()
	for deadline := time.Now().Add(40 * time.Minute); time.Now().Before(deadline); {
		if _, e := os.Stat("/tmp/elabtrack-phase10-browser-stop"); e == nil {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("browser gate timeout")
}
