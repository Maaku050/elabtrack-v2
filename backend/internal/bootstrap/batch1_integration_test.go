package bootstrap

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/accounts"
	appterms "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/spreadsheet"
	"github.com/Maaku050/elabtrack-v2/backend/internal/interface/http/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// Test-only adapter, compiled out of the product. Private inbox never enters reports.
type browserInbox struct {
	mu       sync.Mutex
	messages []d.Mail
	path     string
}

func (m *browserInbox) SendActivation(_ context.Context, v d.Mail) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.messages = append(m.messages, v)
	if m.path != "" {
		data, _ := json.Marshal(m.messages)
		if e := os.WriteFile(m.path, data, 0600); e != nil {
			return "", d.ErrMailUnknown
		}
	}
	return "SYNTHETIC-submission-not-live", nil
}
func batchInfra(t *testing.T) *Infrastructure {
	t.Helper()
	if os.Getenv("ELABTRACK_BATCH1") != "1" {
		t.Skip("requires owned isolated Batch 1 PG")
	}
	cfg, e := config.Load()
	if e != nil || cfg.DB.Name != "elabtrack_v2_batch1_test" || cfg.DB.Port != "54832" || cfg.DB.User != "elabtrack_runtime" {
		t.Fatal("refuse nonisolated target")
	}
	cfg.Security.RateLimitMax = 5000
	cfg.Security.LoginRateLimitMax = 100
	infra, e := initInfrastructure(context.Background(), cfg)
	if e != nil {
		t.Fatal("test infrastructure unavailable")
	}
	t.Cleanup(infra.DB.Close)
	return infra
}
func fixtureAccount(t *testing.T, infra *Infrastructure, role user.Role) (*user.User, string) {
	t.Helper()
	hasher, _ := newSecurityAdapters(infra.Config)
	b := make([]byte, 24)
	_, e := rand.Read(b)
	if e != nil {
		t.Fatal("fixture entropy")
	}
	password := hex.EncodeToString(b)
	hash, e := hasher.Hash(password)
	if e != nil {
		t.Fatal("fixture bcrypt")
	}
	u := user.NewUser("batch1-"+uuid.NewString()+"@example.invalid", "TEST "+string(role), hash)
	u.Role = role
	if postgres.NewUserRepository(infra.DB.Pool).Create(context.Background(), u) != nil {
		t.Fatal("owned fixture insert")
	}
	return u, password
}
func testAccountsContainer(infra *Infrastructure, inbox *browserInbox) *Container {
	c := buildContainer(infra)
	hasher, _ := newSecurityAdapters(infra.Config)
	c.AccountManagementSvc = app.NewService(postgres.NewAccountsRepository(infra.DB.Pool), infra.Tx, hasher, inbox, infra.Config.Accounts.Policy)
	c.AccountManagement = handlers.NewAccountsHandler(c.AccountManagementSvc, spreadsheet.StudentRoster{}, infra.Config.App.Env, infra.Config.Security)
	return c
}
func TestBatch1HTTP(t *testing.T) {
	infra := batchInfra(t)
	c := testAccountsContainer(infra, &browserInbox{})
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	admin, _ := fixtureAccount(t, infra, user.RoleAdmin)
	staff, _ := fixtureAccount(t, infra, user.RoleStaff)
	borrower, _ := fixtureAccount(t, infra, user.RoleBorrower)
	_, issuer := newSecurityAdapters(infra.Config)
	tokens := map[string]string{}
	for _, u := range []*user.User{admin, staff, borrower} {
		token, e := issuer.IssueAccessToken(context.Background(), u)
		if e != nil {
			t.Fatal("fixture token")
		}
		tokens[string(u.Role)] = token
	}
	request := func(method, path, body, role, origin string, headers map[string]string) (int, string) {
		t.Helper()
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if role != "" {
			req.Header.Set("Authorization", "Bearer "+tokens[role])
		}
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		res, e := server.Test(req, fiber.TestConfig{Timeout: 15 * time.Second})
		if e != nil {
			t.Fatal("test HTTP failed")
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		if strings.Contains(string(b), "token_hash") || strings.Contains(string(b), "password") || strings.Contains(string(b), "panic") {
			t.Fatal("sensitive response")
		}
		return res.StatusCode, string(b)
	}
	t.Run("role_origin_strict_payload_and_cors", func(t *testing.T) {
		body := `{"name":"Synthetic Faculty","email":"` + uuid.NewString() + `@example.invalid","borrower_type":"FACULTY"}`
		for _, role := range []string{"STAFF", "BORROWER", ""} {
			status, _ := request("POST", "/api/v1/borrowers", body, role, infra.Config.Security.FrontendURL, map[string]string{"Idempotency-Key": uuid.NewString()})
			want := 403
			if role == "" {
				want = 401
			}
			if status != want {
				t.Fatal("creation authority")
			}
		}
		for _, payload := range []string{body[:len(body)-1] + `,"role":"ADMIN"}`, body[:len(body)-1] + `,"password":"forbidden"}`, `null`, `{}`} {
			status, _ := request("POST", "/api/v1/borrowers", payload, "ADMIN", infra.Config.Security.FrontendURL, map[string]string{"Idempotency-Key": uuid.NewString()})
			if status != 400 {
				t.Fatal("strict payload boundary")
			}
		}
		status, _ := request("POST", "/api/v1/borrowers", body, "ADMIN", "", map[string]string{"Idempotency-Key": uuid.NewString()})
		if status != 403 {
			t.Fatal("missing Origin")
		}
		status, _ = request("OPTIONS", "/api/v1/borrowers", "", "", infra.Config.Security.FrontendURL, map[string]string{"Access-Control-Request-Method": "POST", "Access-Control-Request-Headers": "authorization,content-type,idempotency-key"})
		if status != 204 {
			t.Fatal("command header preflight")
		}
		status, response := request("POST", "/api/v1/borrowers", body, "ADMIN", infra.Config.Security.FrontendURL, map[string]string{"Idempotency-Key": uuid.NewString()})
		if status != 201 || !strings.Contains(response, "activation_required") {
			t.Fatal("authorized safe creation")
		}
		var created struct {
			Data d.Record `json:"data"`
		}
		if json.Unmarshal([]byte(response), &created) != nil {
			t.Fatal("created response")
		}
		profile, _ := json.Marshal(app.ProfileInput{Name: "TEST edited Faculty", ExpectedUpdatedAt: created.Data.UpdatedAt})
		status, _ = request("PATCH", "/api/v1/borrowers/"+created.Data.ID.String(), string(profile), "ADMIN", infra.Config.Security.FrontendURL, map[string]string{"Idempotency-Key": uuid.NewString()})
		if status != 200 {
			t.Fatal("trusted-Origin PATCH profile")
		}
		status, _ = request("PATCH", "/api/v1/borrowers/"+created.Data.ID.String(), string(profile), "ADMIN", "", map[string]string{"Idempotency-Key": uuid.NewString()})
		if status != 403 {
			t.Fatal("PATCH missing Origin denied")
		}
		for _, path := range []string{"/api/v1/staff-accounts", "/api/v1/borrowers/student-template", "/api/v1/borrowers/policy"} {
			status, _ := request("GET", path, "", "STAFF", infra.Config.Security.FrontendURL, nil)
			if status != 403 {
				t.Fatal("privileged route")
			}
		}
		status, _ = request("POST", "/api/v1/auth/register", "{}", "", infra.Config.Security.FrontendURL, nil)
		if status != 404 {
			t.Fatal("public registration unavailable")
		}
		status, _ = request("POST", "/api/v1/auth/activate", `{"token":"invalid","password":"valid-password"}`, "", infra.Config.Security.FrontendURL, nil)
		if status != 400 {
			t.Fatal("activation invalid safely")
		}
	})
}
func TestBatch1BrowserServer(t *testing.T) {
	if os.Getenv("BATCH1_BROWSER_SERVER") != "1" {
		t.Skip("explicit private browser harness")
	}
	infra := batchInfra(t)
	inboxPath := "/tmp/elabtrack-batch1-inbox.json"
	f, e := os.OpenFile(inboxPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		t.Fatal("private inbox exists; refuse overwrite")
	}
	_, _ = f.WriteString("[]")
	_ = f.Close()
	inbox := &browserInbox{path: inboxPath}
	c := testAccountsContainer(infra, inbox)
	if os.Getenv("BATCH1_BROWSER_MAIL_UNCONFIGURED") == "1" {
		// Exercise truthful missing-delivery UI through the real adapter, never
		// a live provider. This test-only switch cannot affect product binaries.
		if infra.Config.Accounts.BrevoKey != "" {
			t.Fatal("refuse live mail in unconfigured browser verification")
		}
		c = buildContainer(infra)
	}
	fixtures := map[string]any{}
	for _, role := range []user.Role{user.RoleAdmin, user.RoleStaff, user.RoleBorrower} {
		u, password := fixtureAccount(t, infra, role)
		fixtures[strings.ToLower(string(role))] = map[string]any{"id": u.ID, "email": u.Email, "name": u.Name, "password": password}
	}
	adminID := fixtures["admin"].(map[string]any)["id"].(uuid.UUID)
	var expected *uuid.UUID
	current, readErr := c.TermsSvc.Current(context.Background(), adminID)
	if readErr == nil {
		expected = &current.ID
	} else if !errors.Is(readErr, terms.ErrNotPublished) {
		t.Fatal("current synthetic terms lookup")
	}
	_, e = c.TermsSvc.Publish(context.Background(), appterms.PublishCommand{ExpectedCurrentVersionID: expected, ActorID: adminID, Version: "TEST-BATCH1-" + uuid.NewString(), Title: "TEST ONLY — NOT OFFICIAL FSMO TERMS", Body: "SYNTHETIC ISOLATED BROWSER TEST. This document is not institutional policy. Borrowing is not implemented in this batch."})
	if e != nil {
		t.Fatal("publish isolated synthetic document")
	}
	im := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			im.SetRGBA(x, y, color.RGBA{R: 52, G: 38, B: 210, A: 255})
		}
	}
	imageFile, err := os.OpenFile("/tmp/elabtrack-batch1-catalog.png", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal("private synthetic image exists")
	}
	if png.Encode(imageFile, im) != nil {
		t.Fatal("synthetic raster")
	}
	_ = imageFile.Close()
	parser := spreadsheet.StudentRoster{}
	raw, e := parser.Template()
	if e != nil {
		t.Fatal("template")
	}
	workbook, e := excelize.OpenReader(strings.NewReader(string(raw)))
	if e != nil {
		t.Fatal("test workbook")
	}
	for row := 2; row <= 5; row++ {
		values := []string{"00" + uuid.NewString(), "TEST roster student " + string(rune('A'+row-2)), uuid.NewString() + "@students.example.invalid", "TEST program", ""}
		if row == 4 {
			values[2] = "invalid@external.example.invalid"
		}
		if row == 5 {
			values[0], _ = workbook.GetCellValue("Sheet1", "A2")
			values[2], _ = workbook.GetCellValue("Sheet1", "C2")
		}
		for column, value := range values {
			cell, _ := excelize.CoordinatesToCellName(column+1, row)
			_ = workbook.SetCellStr("Sheet1", cell, value)
		}
	}
	if workbook.SaveAs("/tmp/elabtrack-batch1-students.xlsx") != nil {
		t.Fatal("synthetic roster")
	}
	_ = workbook.Close()
	data, _ := json.Marshal(fixtures)
	f, e = os.OpenFile("/tmp/elabtrack-batch1-fixtures.json", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if e != nil {
		t.Fatal("private fixtures exist; refuse overwrite")
	}
	_, _ = f.Write(data)
	_ = f.Close()
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	go func() {
		if e := server.Listen("127.0.0.1:18085", fiber.ListenConfig{DisableStartupMessage: true}); e != nil {
			t.Error("browser listener stopped")
		}
	}()
	defer server.Shutdown()
	deadline := time.Now().Add(45 * time.Minute)
	for time.Now().Before(deadline) {
		if _, e := os.Stat("/tmp/elabtrack-batch1-browser-stop"); e == nil {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatal("browser harness timeout")
}

// Regenerate only this owned synthetic workbook for browser retries, without DB writes.
func TestBatch1BrowserRoster(t *testing.T) {
	if os.Getenv("BATCH1_BROWSER_ROSTER") != "1" {
		t.Skip("explicit disposable roster generation")
	}
	_ = batchInfra(t)
	workbook := excelize.NewFile()
	defer workbook.Close()
	headers := []string{"studentId", "name", "email", "course", "contactNumber"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = workbook.SetCellStr("Sheet1", cell, h)
	}
	firstID, firstEmail := "00"+uuid.NewString(), uuid.NewString()+"@students.example.invalid"
	for row := 2; row <= 5; row++ {
		values := []string{"00" + uuid.NewString(), "TEST roster student " + string(rune('A'+row-2)), uuid.NewString() + "@students.example.invalid", "TEST program", ""}
		if row == 2 || row == 5 {
			values[0], values[2] = firstID, firstEmail
		}
		if row == 4 {
			values[2] = uuid.NewString() + "@external.example.invalid"
		}
		for column, value := range values {
			cell, _ := excelize.CoordinatesToCellName(column+1, row)
			_ = workbook.SetCellStr("Sheet1", cell, value)
		}
	}
	raw, e := workbook.WriteToBuffer()
	if e != nil {
		t.Fatal("synthetic workbook")
	}
	file, e := os.CreateTemp("/tmp", "elabtrack-batch1-roster-*.xlsx")
	if e != nil {
		t.Fatal("private roster file")
	}
	name := file.Name()
	defer os.Remove(name)
	if _, e = file.Write(raw.Bytes()); e != nil {
		t.Fatal("private roster write")
	}
	_ = file.Close()
	if os.Rename(name, "/tmp/elabtrack-batch1-students.xlsx") != nil {
		t.Fatal("private roster installation")
	}
}

func TestBatch1InventoryHTTP(t *testing.T) {
	infra := batchInfra(t)
	c := testAccountsContainer(infra, &browserInbox{})
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	users := map[string]*user.User{}
	tokens := map[string]string{}
	_, issuer := newSecurityAdapters(infra.Config)
	for _, role := range []user.Role{user.RoleAdmin, user.RoleStaff, user.RoleBorrower} {
		u, _ := fixtureAccount(t, infra, role)
		users[string(role)] = u
		tok, e := issuer.IssueAccessToken(context.Background(), u)
		if e != nil {
			t.Fatal("token")
		}
		tokens[string(role)] = tok
	}
	request := func(method, path, body, role, origin, key string) (int, []byte) {
		t.Helper()
		req := httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		if role != "" {
			req.Header.Set("Authorization", "Bearer "+tokens[role])
		}
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		if key != "" {
			req.Header.Set("Idempotency-Key", key)
		}
		res, e := server.Test(req)
		if e != nil {
			t.Fatal("HTTP request")
		}
		defer res.Body.Close()
		raw, _ := io.ReadAll(res.Body)
		return res.StatusCode, raw
	}
	origin := "http://localhost:15175"
	input := `{"name":"TEST HTTP inventory","description":"","category_id":null,"opening_quantity":2,"reason":"Synthetic opening"}`
	for _, role := range []string{"", "BORROWER"} {
		status, _ := request("POST", "/api/v1/equipment", input, role, origin, uuid.NewString())
		if status != 401 && status != 403 {
			t.Fatal("unauthorized inventory write")
		}
	}
	status, _ := request("POST", "/api/v1/equipment", input, "STAFF", "", uuid.NewString())
	if status != 403 {
		t.Fatal("Origin required")
	}
	for _, body := range []string{`null`, `{"name":"X","reserved":9}`, `{"name":"X","opening_quantity":-1}`, `{"name":"X","opening_quantity":1.5}`} {
		status, _ = request("POST", "/api/v1/equipment", body, "STAFF", origin, uuid.NewString())
		if status != 400 {
			t.Fatal("strict stock creation")
		}
	}
	status, _ = request("POST", "/api/v1/equipment", input, "STAFF", origin, "")
	if status != 400 {
		t.Fatal("key required")
	}
	status, raw := request("POST", "/api/v1/equipment", input, "STAFF", origin, uuid.NewString())
	if status != 201 {
		t.Fatal("Staff creates equipment")
	}
	var envelope struct {
		Data struct {
			ID       uuid.UUID `json:"id"`
			Version  int64     `json:"metadata_version"`
			Sequence int64     `json:"stock_sequence"`
		} `json:"data"`
	}
	if json.Unmarshal(raw, &envelope) != nil || envelope.Data.ID == uuid.Nil {
		t.Fatal("created ID")
	}
	id := envelope.Data.ID.String()
	status, _ = request("GET", "/api/v1/equipment/"+id, "", "BORROWER", "", "")
	if status != 200 {
		t.Fatal("active catalog")
	}
	for _, body := range []string{`{"kind":"RECONCILE","reason":"Count","expected_sequence":1,"confirm":true}`, `{"kind":"RECONCILE","quantity":null,"reason":"Count","expected_sequence":1,"confirm":true}`, `{"kind":"ADD","quantity":1,"reserved":1,"reason":"Count","confirm":true}`} {
		status, _ = request("POST", "/api/v1/equipment/"+id+"/adjustments", body, "ADMIN", origin, uuid.NewString())
		if status != 400 {
			t.Fatal("explicit quantity and no custody payload")
		}
	}
	reconcile := `{"kind":"RECONCILE","quantity":0,"reason":"Verified count","expected_sequence":1,"confirm":true}`
	status, _ = request("POST", "/api/v1/equipment/"+id+"/adjustments", reconcile, "STAFF", origin, uuid.NewString())
	if status != 403 {
		t.Fatal("Staff correction denied")
	}
	status, _ = request("POST", "/api/v1/equipment/"+id+"/adjustments", reconcile, "ADMIN", origin, uuid.NewString())
	if status != 200 {
		t.Fatal("Admin explicit zero reconciliation")
	}
	status, _ = request("GET", "/api/v1/equipment/"+id+"/movements", "", "BORROWER", "", "")
	if status != 403 {
		t.Fatal("operational history private")
	}
	status, _ = request("PATCH", "/api/v1/equipment/"+id+"/status", `{"status":"INACTIVE","expected_version":1,"confirm":true}`, "STAFF", origin, uuid.NewString())
	if status != 200 {
		t.Fatal("trusted PATCH")
	}
	status, _ = request("GET", "/api/v1/equipment/"+id, "", "BORROWER", "", "")
	if status != 404 {
		t.Fatal("inactive catalog scoped")
	}
	status, _ = request("GET", "/api/v1/equipment?status=ARCHIVED", "", "BORROWER", "", "")
	if status != 403 {
		t.Fatal("private status filter")
	}
	_, e := infra.DB.Pool.Exec(context.Background(), `UPDATE users SET is_active=false WHERE id=$1`, users["STAFF"].ID)
	if e != nil {
		t.Fatal("fixture status")
	}
	status, _ = request("GET", "/api/v1/equipment", "", "STAFF", "", "")
	if status != 403 {
		t.Fatal("current-account JWT authority")
	}
}
