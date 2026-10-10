package bootstrap

import (
	"context"
	"encoding/json"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
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

type phase8Clock struct {
	*postgres.BorrowingRepository
	at time.Time
}

func (r phase8Clock) Clock(context.Context) (time.Time, error) { return r.at, nil }
func TestPhase8HTTP(t *testing.T) {
	infra := batchInfra(t)
	c := buildContainer(infra)
	users, _, eq := phase7Fixtures(t, infra, c)
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	_, issuer := newSecurityAdapters(infra.Config)
	tokens := map[string]string{}
	for role, u := range users {
		token, e := issuer.IssueAccessToken(context.Background(), u)
		if e != nil {
			t.Fatal("fixture token")
		}
		tokens[role] = token
	}
	due := time.Now().Add(time.Hour).UTC().Truncate(time.Microsecond)
	loan, e := c.BorrowingSvc.Direct(context.Background(), users["staff"].ID, uuid.NewString(), d.Input{BorrowerID: users["borrower"].ID, Items: []d.Line{{EquipmentID: eq.ID, Quantity: 3}}, Confirm: true, Handover: true, DueAt: &due})
	if e != nil {
		t.Fatal("fixture issue")
	}
	request := func(role, path, key string, input any) (int, map[string]any) {
		t.Helper()
		raw, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/api/v1/borrowings/"+loan.ID.String()+path, strings.NewReader(string(raw)))
		req.Header.Set("Authorization", "Bearer "+tokens[role])
		req.Header.Set("Origin", infra.Config.Security.FrontendURL)
		req.Header.Set("Idempotency-Key", key)
		req.Header.Set("Content-Type", "application/json")
		res, e := server.Test(req, fiber.TestConfig{Timeout: 15 * time.Second})
		if e != nil {
			t.Fatal("HTTP")
		}
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)
		var out map[string]any
		if json.Unmarshal(body, &out) != nil {
			t.Fatal("envelope")
		}
		return res.StatusCode, out
	}
	in := d.ReturnInput{ExpectedEvents: len(loan.Events), Confirm: true, Reason: "TEST physical inspection", Lines: []d.ReturnLine{{ItemID: loan.Items[0].ID, Good: 1, Damaged: 1, Lost: 1}}}
	t.Run("borrower_cannot_mutate_dispositions", func(t *testing.T) {
		for _, p := range []string{"/returns", "/replacements", "/fine-clearances"} {
			code, _ := request("borrower", p, uuid.NewString(), in)
			if code != 403 {
				t.Fatal("authority", code)
			}
		}
	})
	t.Run("strict_fields_and_review_token", func(t *testing.T) {
		code, _ := request("staff", "/returns", uuid.NewString(), map[string]any{"lines": in.Lines, "reason": in.Reason, "confirm": true, "expected_event_count": in.ExpectedEvents, "completed_at": "invented"})
		if code != 400 {
			t.Fatal("strict field", code)
		}
		bad := in
		bad.ExpectedEvents = 0
		code, _ = request("staff", "/returns", uuid.NewString(), bad)
		if code != 400 {
			t.Fatal("review required", code)
		}
	})
	t.Run("mixed_return_receipt_and_stale_review", func(t *testing.T) {
		key := uuid.NewString()
		code, out := request("staff", "/returns", key, in)
		if code != 200 {
			t.Fatal("mixed return", code)
		}
		data := out["data"].(map[string]any)
		if data["status"] != "CHECKED_OUT" || len(data["obligations"].([]any)) != 2 {
			t.Fatal("separate replacement accountability")
		}
		code, _ = request("staff", "/returns", key, in)
		if code != 200 {
			t.Fatal("safe retry", code)
		}
		code, _ = request("staff", "/returns", uuid.NewString(), in)
		if code != 409 {
			t.Fatal("stale review", code)
		}
	})
	t.Run("staff_cannot_clear_fine", func(t *testing.T) {
		code, _ := request("staff", "/fine-clearances", uuid.NewString(), d.FineInput{Method: "WAIVED", ExpectedMinor: 1000, Confirm: true})
		if code != 403 {
			t.Fatal("Admin only", code)
		}
	})
}
func TestPhase8BrowserServer(t *testing.T) {
	if os.Getenv("ELABTRACK_PHASE8_BROWSER") != "1" {
		t.Skip("explicit isolated browser gate")
	}
	infra := batchInfra(t)
	c := buildContainer(infra)
	users, passwords, eq := phase7Fixtures(t, infra, c)
	ctx := context.Background()
	start := time.Now().Add(-72 * time.Hour).UTC().Truncate(time.Microsecond)
	due := start.Add(24 * time.Hour)
	past := app.NewService(phase8Clock{postgres.NewBorrowingRepository(infra.DB.Pool), start}, postgres.NewAccountsRepository(infra.DB.Pool), postgres.NewInventoryRepository(infra.DB.Pool), c.TermsSvc, infra.Tx)
	loan, e := past.Direct(ctx, users["staff"].ID, uuid.NewString(), d.Input{BorrowerID: users["borrower"].ID, Confirm: true, Handover: true, DueAt: &due, Items: []d.Line{{EquipmentID: eq.ID, Quantity: 7}}})
	if e != nil {
		t.Fatal("synthetic historical issue")
	}
	loan, e = c.BorrowingSvc.Return(ctx, users["staff"].ID, loan.ID, uuid.NewString(), d.ReturnInput{ExpectedEvents: len(loan.Events), Confirm: true, Reason: "TEST synthetic mixed return history", Lines: []d.ReturnLine{{ItemID: loan.Items[0].ID, Good: 2, Damaged: 1, Lost: 1}}})
	if e != nil {
		t.Fatal("synthetic mixed return")
	}
	fixtures := map[string]any{"equipment": eq, "loan": loan}
	for role, u := range users {
		fixtures[role] = map[string]any{"id": u.ID, "email": u.Email, "password": passwords[role]}
	}
	file, e := os.OpenFile("/tmp/elabtrack-phase8-fixtures.json", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		t.Fatal("private fixture exists")
	}
	raw, _ := json.Marshal(fixtures)
	_, e = file.Write(raw)
	file.Close()
	if e != nil {
		t.Fatal("private fixture write")
	}
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	go func() { _ = server.Listen("127.0.0.1:18085", fiber.ListenConfig{DisableStartupMessage: true}) }()
	defer server.Shutdown()
	for deadline := time.Now().Add(40 * time.Minute); time.Now().Before(deadline); {
		if _, e := os.Stat("/tmp/elabtrack-phase8-browser-stop"); e == nil {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("browser gate timeout")
}
