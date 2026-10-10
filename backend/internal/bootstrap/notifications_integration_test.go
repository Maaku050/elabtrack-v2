package bootstrap

import (
	"context"
	"encoding/json"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	n "github.com/Maaku050/elabtrack-v2/backend/internal/domain/notifications"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestPhase9HTTP(t *testing.T) {
	infra := batchInfra(t)
	c := buildContainer(infra)
	users, _, eq := phase7Fixtures(t, infra, c)
	ctx := context.Background()
	loan, e := c.BorrowingSvc.Submit(ctx, users["borrower"].ID, uuid.NewString(), d.Input{Confirm: true, Items: []d.Line{{EquipmentID: eq.ID, Quantity: 1}}})
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 100; i++ {
		count, e := c.NotificationsSvc.Sweep(ctx, 100)
		if e != nil {
			t.Fatal(e)
		}
		if count < 100 {
			break
		}
	}
	page, e := c.NotificationsSvc.List(ctx, users["borrower"].ID, n.Filter{Page: 1, PerPage: 100})
	if e != nil || len(page.Items) == 0 {
		t.Fatal("notification fixture", e)
	}
	var id string
	for _, v := range page.Items {
		if v.BorrowingID == loan.ID {
			id = v.ID.String()
		}
	}
	if id == "" {
		t.Fatal("missing source")
	}
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	_, issuer := newSecurityAdapters(infra.Config)
	token, e := issuer.IssueAccessToken(ctx, users["borrower"])
	if e != nil {
		t.Fatal(e)
	}
	request := func(method, path, body, origin string, authenticated bool) int {
		t.Helper()
		r := httptest.NewRequest(method, "/api/v1/notifications"+path, strings.NewReader(body))
		if authenticated {
			r.Header.Set("Authorization", "Bearer "+token)
		}
		if origin != "" {
			r.Header.Set("Origin", origin)
		}
		r.Header.Set("Content-Type", "application/json")
		res, e := server.Test(r, fiber.TestConfig{Timeout: 10 * time.Second})
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		raw, _ := io.ReadAll(res.Body)
		if !json.Valid(raw) {
			t.Fatal("envelope")
		}
		return res.StatusCode
	}
	for _, tc := range []struct {
		method, path, body, origin string
		auth                       bool
		want                       int
	}{
		{"GET", "/count", "", "", false, 401}, {"GET", "?page=0", "", "", true, 400}, {"GET", "?per_page=101", "", "", true, 400}, {"GET", "?read=invalid", "", "", true, 400},
		{"PATCH", "/" + id, `{"read":true}`, "", true, 403}, {"PATCH", "/" + id, `{"read":true}`, "https://untrusted.invalid", true, 403},
		{"PATCH", "/" + id, `{}`, infra.Config.Security.FrontendURL, true, 400}, {"PATCH", "/" + id, `{"read":true,"recipient_id":"forged"}`, infra.Config.Security.FrontendURL, true, 400},
		{"PATCH", "/" + uuid.NewString(), `{"read":true}`, infra.Config.Security.FrontendURL, true, 404},
		{"PATCH", "/" + id, `{"read":true}`, infra.Config.Security.FrontendURL, true, 200}, {"PATCH", "/" + id, `{"read":false}`, infra.Config.Security.FrontendURL, true, 200},
		{"GET", "?read=unread&page=1&per_page=1", "", "", true, 200}, {"GET", "/count", "", "", true, 200},
		{"POST", "/read-all", "{}", infra.Config.Security.FrontendURL, true, 400},
		{"POST", "/read-all", `{"confirm":true}`, "", true, 403},
		{"POST", "/read-all", `{"confirm":true}`, "https://untrusted.invalid", true, 403},
		{"POST", "/read-all", `{"confirm":true}`, infra.Config.Security.FrontendURL, false, 401},
		{"POST", "/read-all", `{"confirm":true}`, infra.Config.Security.FrontendURL, true, 200},
		{"POST", "/read-all", `{"confirm":true}`, infra.Config.Security.FrontendURL, true, 200},
	} {
		if got := request(tc.method, tc.path, tc.body, tc.origin, tc.auth); got != tc.want {
			t.Fatalf("%s %s: %d expected %d", tc.method, tc.path, got, tc.want)
		}
	}
	if count, e := c.NotificationsSvc.Count(ctx, users["borrower"].ID); e != nil || count != 0 {
		t.Fatal("durable bulk-read count", count, e)
	}
	if _, e := infra.DB.Pool.Exec(ctx, `UPDATE users SET is_active=false WHERE id=$1`, users["borrower"].ID); e != nil {
		t.Fatal(e)
	}
	if got := request("POST", "/read-all", `{"confirm":true}`, infra.Config.Security.FrontendURL, true); got != 403 {
		t.Fatal("inactive bulk-read token", got)
	}
	if got := request("GET", "/count", "", "", true); got != 403 {
		t.Fatal("deactivated token", got)
	}
}
func TestPhase9BrowserServer(t *testing.T) {
	if os.Getenv("ELABTRACK_PHASE9_BROWSER") != "1" {
		t.Skip("explicit isolated browser gate")
	}
	infra := batchInfra(t)
	c := buildContainer(infra)
	users, passwords, eq := phase7Fixtures(t, infra, c)
	ctx := context.Background()
	loan, e := c.BorrowingSvc.Submit(ctx, users["borrower"].ID, uuid.NewString(), d.Input{Confirm: true, Items: []d.Line{{EquipmentID: eq.ID, Quantity: 1}}})
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 100; i++ {
		n, e := c.NotificationsSvc.Sweep(ctx, 100)
		if e != nil {
			t.Fatal(e)
		}
		if n < 100 {
			break
		}
	}
	fixtures := map[string]any{"loan": loan}
	for role, u := range users {
		fixtures[role] = map[string]any{"email": u.Email, "password": passwords[role]}
	}
	f, e := os.OpenFile("/tmp/elabtrack-phase9-fixtures.json", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
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
		if _, e := os.Stat("/tmp/elabtrack-phase9-browser-stop"); e == nil {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("browser gate timeout")
}
