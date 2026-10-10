package bootstrap

import (
	"context"
	"encoding/json"
	"errors"
	appborrow "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	appterms "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
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

func phase7Fixtures(t *testing.T, infra *Infrastructure, c *Container) (map[string]*user.User, map[string]string, inventory.Equipment) {
	t.Helper()
	ctx := context.Background()
	users := map[string]*user.User{}
	passwords := map[string]string{}
	for _, role := range []user.Role{user.RoleAdmin, user.RoleStaff, user.RoleBorrower} {
		u, p := fixtureAccount(t, infra, role)
		users[strings.ToLower(string(role))] = u
		passwords[strings.ToLower(string(role))] = p
	}
	u := users["borrower"]
	if _, e := infra.DB.Pool.Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type) VALUES($1,'FACULTY')`, u.ID); e != nil {
		t.Fatal("test-only profile")
	}
	var expected *uuid.UUID
	v, e := c.TermsSvc.Current(ctx, users["admin"].ID)
	if e == nil {
		expected = &v.ID
	} else if !errors.Is(e, terms.ErrNotPublished) {
		t.Fatal("test publication lookup")
	}
	v, e = c.TermsSvc.Publish(ctx, appterms.PublishCommand{ActorID: users["admin"].ID, ExpectedCurrentVersionID: expected, Version: "TEST-P7-" + uuid.NewString()[:12], Title: "TEST ONLY — NOT OFFICIAL FSMO TERMS", Body: "Synthetic Phase7 verification policy in an isolated database. Not approved institutional content."})
	if e != nil {
		t.Fatal("test publication")
	}
	if _, e = c.TermsSvc.Accept(ctx, u.ID, v.ID); e != nil {
		t.Fatal("test consent")
	}
	eq, e := c.InventorySvc.Create(ctx, users["staff"].ID, uuid.NewString(), inventory.Create{Metadata: inventory.Metadata{Name: "TEST Culinary Ladle " + uuid.NewString()[:8], Description: "Synthetic Phase7 borrowing acceptance fixture"}, Opening: 20, Reason: "Synthetic test opening stock"})
	if e != nil {
		t.Fatal("test equipment")
	}
	return users, passwords, eq
}
func TestPhase7HTTP(t *testing.T) {
	infra := batchInfra(t)
	c := buildContainer(infra)
	users, _, eq := phase7Fixtures(t, infra, c)
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	_, issuer := newSecurityAdapters(infra.Config)
	tokens := map[string]string{}
	for role, u := range users {
		token, e := issuer.IssueAccessToken(context.Background(), u)
		if e != nil {
			t.Fatal("token fixture")
		}
		tokens[role] = token
	}
	request := func(method, path, role, origin, key string, body any) (int, map[string]any) {
		t.Helper()
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest(method, "/api/v1"+path, strings.NewReader(string(raw)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", origin)
		req.Header.Set("Idempotency-Key", key)
		if role != "" {
			req.Header.Set("Authorization", "Bearer "+tokens[role])
		}
		res, e := server.Test(req, fiber.TestConfig{Timeout: 15 * time.Second})
		if e != nil {
			t.Fatal("HTTP request")
		}
		defer res.Body.Close()
		data, _ := io.ReadAll(res.Body)
		var out map[string]any
		if json.Unmarshal(data, &out) != nil {
			t.Fatal("safe envelope")
		}
		return res.StatusCode, out
	}
	origin := infra.Config.Security.FrontendURL
	in := d.Input{Items: []d.Line{{EquipmentID: eq.ID, Quantity: 2}}, Confirm: true}
	key := uuid.NewString()
	t.Run("current_issuance_eligibility_and_server_filtered_directory", func(t *testing.T) {
		path := "/borrowings/eligibility/" + users["borrower"].ID.String()
		for _, role := range []string{"staff", "admin", "borrower", ""} {
			want := 200
			if role == "borrower" {
				want = 403
			}
			if role == "" {
				want = 401
			}
			status, _ := request("GET", path, role, origin, "", nil)
			if status != want {
				t.Fatal("eligibility authority", role, status)
			}
		}
		query := "/borrowers?eligible_for_issuance=true&search=" + users["borrower"].Email + "&per_page=1"
		status, out := request("GET", query, "staff", origin, "", nil)
		if status != 200 || out["data"].(map[string]any)["total"].(float64) != 1 {
			t.Fatal("eligible exact server count", status)
		}
		if _, e := infra.DB.Pool.Exec(context.Background(), `UPDATE users SET activation_required=true WHERE id=$1`, users["borrower"].ID); e != nil {
			t.Fatal(e)
		}
		status, _ = request("GET", path, "staff", origin, "", nil)
		if status != 403 {
			t.Fatal("activation eligibility", status)
		}
		status, out = request("GET", query, "staff", origin, "", nil)
		if status != 200 || out["data"].(map[string]any)["total"].(float64) != 0 {
			t.Fatal("activation filtered before pagination", status)
		}
		if _, e := infra.DB.Pool.Exec(context.Background(), `UPDATE users SET activation_required=false WHERE id=$1`, users["borrower"].ID); e != nil {
			t.Fatal(e)
		}
		status, _ = request("GET", "/borrowers?eligible_for_issuance=invalid", "staff", origin, "", nil)
		if status != 400 {
			t.Fatal("strict eligibility filter", status)
		}
	})
	t.Run("submit_replay_strict_origin_and_roles", func(t *testing.T) {
		for _, role := range []string{"staff", "admin", ""} {
			status, _ := request("POST", "/borrowings", role, origin, uuid.NewString(), in)
			if role == "" && status != 401 || role != "" && status != 403 {
				t.Fatal("submit authority", status)
			}
		}
		for _, bad := range []any{map[string]any{"items": in.Items, "confirm": true, "status": "CHECKED_OUT"}, nil, map[string]any{}} {
			status, _ := request("POST", "/borrowings", "borrower", origin, uuid.NewString(), bad)
			if status != 400 {
				t.Fatal("strict input", status)
			}
		}
		status, _ := request("POST", "/borrowings", "borrower", "https://untrusted.example.invalid", key, in)
		if status != 403 {
			t.Fatal("origin")
		}
		status, _ = request("POST", "/borrowings", "borrower", origin, "bad", in)
		if status != 400 {
			t.Fatal("key")
		}
	})
	status, out := request("POST", "/borrowings", "borrower", origin, key, in)
	if status != 201 {
		t.Fatal("submit", status)
	}
	id := out["data"].(map[string]any)["id"].(string)
	t.Run("real_scoped_page_and_cancel", func(t *testing.T) {
		status, out = request("POST", "/borrowings", "borrower", origin, key, in)
		if status != 201 || out["data"].(map[string]any)["id"] != id {
			t.Fatal("replay")
		}
		status, out = request("GET", "/borrowings?status=PENDING&per_page=1", "borrower", origin, "", nil)
		if status != 200 || out["data"].(map[string]any)["total"].(float64) != 1 {
			t.Fatal("scoped count")
		}
		status, _ = request("POST", "/borrowings/"+id+"/cancel", "staff", origin, uuid.NewString(), d.Decision{Confirm: true})
		if status != 403 {
			t.Fatal("staff cancel")
		}
		status, out = request("POST", "/borrowings/"+id+"/cancel", "borrower", origin, uuid.NewString(), d.Decision{Confirm: true})
		if status != 200 || out["data"].(map[string]any)["status"] != "CANCELLED" {
			t.Fatal("cancel")
		}
		status, _ = request("GET", "/borrowings?status=APPROVED", "borrower", origin, "", nil)
		if status != 400 {
			t.Fatal("state allowlist")
		}
	})
	t.Run("staff_approval_denial_and_direct_contracts", func(t *testing.T) {
		status, out := request("POST", "/borrowings", "borrower", origin, uuid.NewString(), in)
		if status != 201 {
			t.Fatal("pending request")
		}
		pendingID := out["data"].(map[string]any)["id"].(string)
		future := time.Now().Add(14 * 24 * time.Hour).UTC().Truncate(time.Microsecond)
		decision := d.Decision{Confirm: true, Handover: true, DueAt: &future}
		status, _ = request("POST", "/borrowings/"+pendingID+"/approve", "borrower", origin, uuid.NewString(), decision)
		if status != 403 {
			t.Fatal("borrower cannot issue")
		}
		status, _ = request("POST", "/borrowings/"+pendingID+"/approve", "staff", origin, uuid.NewString(), d.Decision{Confirm: true, DueAt: &future})
		if status != 400 {
			t.Fatal("handover required")
		}
		status, out = request("POST", "/borrowings/"+pendingID+"/approve", "staff", origin, uuid.NewString(), decision)
		if status != 200 || out["data"].(map[string]any)["status"] != "CHECKED_OUT" {
			t.Fatal("atomic checkout", status)
		}
		direct := in
		direct.BorrowerID = users["borrower"].ID
		direct.DueAt = &future
		direct.Handover = true
		status, _ = request("POST", "/borrowings/direct-checkout", "borrower", origin, uuid.NewString(), direct)
		if status != 403 {
			t.Fatal("direct authority")
		}
		status, out = request("POST", "/borrowings/direct-checkout", "admin", origin, uuid.NewString(), direct)
		if status != 201 || out["data"].(map[string]any)["entry_path"] != "DIRECT" {
			t.Fatal("admin direct")
		}
		status, out = request("POST", "/borrowings", "borrower", origin, uuid.NewString(), in)
		if status != 201 {
			t.Fatal("denial request")
		}
		pendingID = out["data"].(map[string]any)["id"].(string)
		status, _ = request("POST", "/borrowings/"+pendingID+"/deny", "staff", origin, uuid.NewString(), d.Decision{Confirm: true})
		if status != 400 {
			t.Fatal("denial reason")
		}
		status, out = request("POST", "/borrowings/"+pendingID+"/deny", "staff", origin, uuid.NewString(), d.Decision{Confirm: true, Reason: "TEST visible reason"})
		if status != 200 || out["data"].(map[string]any)["denial_reason"] != "TEST visible reason" {
			t.Fatal("denial persistence")
		}
	})
	t.Run("current_status_revocation", func(t *testing.T) {
		if _, e := infra.DB.Pool.Exec(context.Background(), `UPDATE users SET is_active=false WHERE id=$1`, users["borrower"].ID); e != nil {
			t.Fatal("synthetic deactivation")
		}
		status, _ := request("GET", "/borrowings", "borrower", origin, "", nil)
		if status != 403 {
			t.Fatal("live revoked status", status)
		}
	})
}

type phase7FixtureClock struct {
	d.Repository
	at time.Time
}

func (r phase7FixtureClock) Clock(context.Context) (time.Time, error) { return r.at, nil }
func TestPhase7BrowserServer(t *testing.T) {
	if os.Getenv("ELABTRACK_PHASE7_BROWSER") != "1" {
		t.Skip("explicit isolated private browser harness")
	}
	infra := batchInfra(t)
	c := buildContainer(infra)
	users, passwords, equipment := phase7Fixtures(t, infra, c)
	fixtures := map[string]any{"equipment": equipment}
	for role, u := range users {
		fixtures[role] = map[string]any{"id": u.ID, "email": u.Email, "password": passwords[role]}
	}
	file, e := os.OpenFile("/tmp/elabtrack-phase7-fixtures.json", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		t.Fatal("private fixture file exists")
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
	workerCtx, cancel := context.WithCancel(context.Background())
	a := &App{borrowing: c.BorrowingSvc, log: infra.Logger, expiryDone: make(chan struct{})}
	go a.runExpiry(workerCtx)
	defer func() { cancel(); <-a.expiryDone }()
	expiryMade := false
	deadline := time.Now().Add(40 * time.Minute)
	for time.Now().Before(deadline) {
		if _, e := os.Stat("/tmp/elabtrack-phase7-expiry-trigger"); e == nil && !expiryMade {
			// Test-only clock fixture: ordinary submit/terms/stock services still enforce
			// all rules. The real persisted deadline/sweep performs expiration.
			repo := postgres.NewBorrowingRepository(infra.DB.Pool)
			past := appborrow.NewService(phase7FixtureClock{repo, time.Now().Add(-25 * time.Hour)}, postgres.NewAccountsRepository(infra.DB.Pool), postgres.NewInventoryRepository(infra.DB.Pool), c.TermsSvc, infra.Tx)
			record, e := past.Submit(context.Background(), users["borrower"].ID, uuid.NewString(), d.Input{Confirm: true, Items: []d.Line{{EquipmentID: equipment.ID, Quantity: 1}}})
			if e != nil {
				t.Fatal("isolated expiration fixture submit")
			}
			if _, e = c.BorrowingSvc.Sweep(context.Background(), 100); e != nil {
				t.Fatal("real expiration sweep")
			}
			raw, _ := json.Marshal(map[string]any{"id": record.ID})
			file, e := os.OpenFile("/tmp/elabtrack-phase7-expiry-result.json", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
			if e != nil {
				t.Fatal("test expiration result exists")
			}
			_, _ = file.Write(raw)
			file.Close()
			expiryMade = true
		}
		if _, e := os.Stat("/tmp/elabtrack-phase7-browser-stop"); e == nil {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("browser harness timeout")
}
