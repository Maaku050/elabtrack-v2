package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

// The existing guard permits only our disposable PostgreSQL database, never
// the owner's development data. Exercise HTTP parsing and the real SQL query.
func TestBatch1InventoryBorrowabilityFilter(t *testing.T) {
	infra := batchInfra(t)
	c := buildContainer(infra)
	server := RegisterRoutes(newServer(infra.Config, c, infra.Logger), c)
	ctx := context.Background()
	admin, _ := fixtureAccount(t, infra, user.RoleAdmin)
	_, issuer := newSecurityAdapters(infra.Config)
	tokens := map[string]string{}
	for _, role := range []user.Role{user.RoleAdmin, user.RoleStaff, user.RoleBorrower} {
		u := admin
		if role != user.RoleAdmin {
			u, _ = fixtureAccount(t, infra, role)
		}
		token, err := issuer.IssueAccessToken(ctx, u)
		if err != nil {
			t.Fatal(err)
		}
		tokens[string(role)] = token
	}
	prefix := "TEST FILTER " + uuid.NewString()
	yes, no := true, false
	cat, err := c.InventorySvc.Category(ctx, admin.ID, uuid.Nil, uuid.NewString(), d.CategoryInput{Name: prefix, Active: &yes})
	if err != nil {
		t.Fatal(err)
	}
	other, err := c.InventorySvc.Category(ctx, admin.ID, uuid.Nil, uuid.NewString(), d.CategoryInput{Name: prefix + " other", Active: &yes})
	if err != nil {
		t.Fatal(err)
	}
	fixtures := []d.Equipment{}
	create := func(name, status string, quantity int64, category uuid.UUID) {
		t.Helper()
		v, e := c.InventorySvc.Create(ctx, admin.ID, uuid.NewString(), d.Create{Metadata: d.Metadata{Name: prefix + " " + name, CategoryID: &category}, Opening: quantity, Reason: "Synthetic verified opening"})
		if e != nil {
			t.Fatal(e)
		}
		if status != "ACTIVE" {
			v, e = c.InventorySvc.Status(ctx, admin.ID, v.ID, uuid.NewString(), d.StatusInput{Status: status, Confirm: true, ExpectedVersion: v.Version})
			if e != nil {
				t.Fatal(e)
			}
		}
		fixtures = append(fixtures, v)
	}
	create("Spoon active positive", "ACTIVE", 20, cat.ID)
	create("Spoon active zero", "ACTIVE", 0, cat.ID)
	create("Spoon inactive positive", "INACTIVE", 20, cat.ID)
	create("Spoon inactive zero", "INACTIVE", 0, cat.ID)
	create("Spoon archived positive", "ARCHIVED", 20, cat.ID)
	for i := 0; i < 26; i++ {
		create(fmt.Sprintf("Utensil %02d", i), "ACTIVE", 1, other.ID)
	}
	// Existing active equipment remains visible even when its category is inactive.
	_, err = c.InventorySvc.Category(ctx, admin.ID, cat.ID, uuid.NewString(), d.CategoryInput{Name: cat.Name, Active: &no, ExpectedVersion: cat.Version})
	if err != nil {
		t.Fatal(err)
	}
	history := func() [3]int64 {
		t.Helper()
		var counts [3]int64
		e := infra.DB.Pool.QueryRow(ctx, `SELECT (SELECT count(*) FROM inventory_movements), (SELECT count(*) FROM inventory_audit_events), (SELECT count(*) FROM inventory_operation_receipts)`).Scan(&counts[0], &counts[1], &counts[2])
		if e != nil {
			t.Fatal(e)
		}
		return counts
	}
	for i, v := range fixtures {
		fixtures[i], err = c.InventorySvc.Detail(ctx, admin.ID, v.ID)
		if err != nil {
			t.Fatal(err)
		}
	}
	before := history()
	get := func(query, role string, wantStatus int) d.Page {
		t.Helper()
		params, e := url.ParseQuery(query)
		if e != nil {
			t.Fatal(e)
		}
		if !params.Has("search") {
			params.Set("search", prefix)
		}
		req := httptest.NewRequest("GET", "/api/v1/equipment?"+params.Encode(), nil)
		if role != "" {
			req.Header.Set("Authorization", "Bearer "+tokens[role])
		}
		res, e := server.Test(req)
		if e != nil {
			t.Fatal(e)
		}
		defer res.Body.Close()
		raw, _ := io.ReadAll(res.Body)
		if res.StatusCode != wantStatus {
			t.Fatalf("query %s: got %d, want %d", query, res.StatusCode, wantStatus)
		}
		var envelope struct {
			Data d.Page `json:"data"`
		}
		if wantStatus == 200 && json.Unmarshal(raw, &envelope) != nil {
			t.Fatal("invalid equipment envelope")
		}
		return envelope.Data
	}
	for _, tc := range []struct {
		name, query      string
		total, available int64
		rows             int
	}{
		{"unfiltered", "page=1&per_page=100", 31, 86, 31},
		{"positive active first page", "available_only=true&page=1&per_page=25", 27, 46, 25},
		{"positive active second page", "available_only=true&page=2&per_page=25", 27, 46, 2},
		{"beyond final page", "available_only=true&page=3&per_page=25", 27, 46, 0},
		{"active intersection", "available_only=true&status=ACTIVE&per_page=100", 27, 46, 27},
		{"inactive intersection", "available_only=true&status=INACTIVE", 0, 0, 0},
		{"archived intersection", "available_only=true&status=ARCHIVED", 0, 0, 0},
		{"inactive without availability", "available_only=false&status=INACTIVE", 2, 20, 2},
		{"archived original visibility", "status=ARCHIVED", 1, 20, 1},
		{"inactive category intersection", "available_only=true&category_id=" + cat.ID.String(), 1, 20, 1},
		{"independent other category", "available_only=true&category_id=" + other.ID.String() + "&per_page=100", 26, 26, 26},
		{"search intersection", "available_only=true&search=" + url.QueryEscape(prefix+" Spoon"), 1, 20, 1},
		{"search inactive no match", "available_only=true&search=" + url.QueryEscape(prefix+" Spoon inactive"), 0, 0, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			q, _ := url.ParseQuery(tc.query)
			for _, role := range []string{"ADMIN", "STAFF"} {
				p := get(tc.query, role, 200)
				if p.Total != tc.total || len(p.Items) != tc.rows || p.Totals.Available != tc.available || p.Totals.Total != tc.available {
					t.Fatalf("%s: rows=%d total=%d totals=%+v", role, len(p.Items), p.Total, p.Totals)
				}
				if q.Get("available_only") == "true" {
					for _, v := range p.Items {
						if v.Status != "ACTIVE" || v.Stock.Available <= 0 {
							t.Fatal("ineligible item")
						}
					}
				}
			}
		})
	}
	t.Run("stable pages and sort", func(t *testing.T) {
		for _, sort := range []string{"name", "available"} {
			seen := map[uuid.UUID]bool{}
			for page := 1; page <= 2; page++ {
				p := get(fmt.Sprintf("available_only=true&page=%d&per_page=25&sort=%s", page, sort), "ADMIN", 200)
				if p.Page != page || p.PerPage != 25 {
					t.Fatal("pagination metadata")
				}
				for _, v := range p.Items {
					if seen[v.ID] {
						t.Fatal("duplicate across pages")
					}
					seen[v.ID] = true
				}
			}
			if len(seen) != 27 {
				t.Fatal("missing eligible equipment across pages")
			}
		}
	})
	t.Run("authority validation and borrower visibility", func(t *testing.T) {
		get("available_only=true", "", 401)
		get("available_only=invalid", "ADMIN", 400)
		get("available_only=true&status=INACTIVE", "BORROWER", 403)
		for _, filter := range []string{"true", "false"} {
			p := get("available_only="+filter+"&per_page=100", "BORROWER", 200)
			want := int64(27)
			if filter == "false" {
				want = 28
			}
			if p.Total != want {
				t.Fatal("borrower total")
			}
			for _, v := range p.Items {
				if v.Status != "ACTIVE" {
					t.Fatal("borrower hidden status leaked")
				}
			}
		}
	})
	for _, v := range fixtures {
		current, e := c.InventorySvc.Detail(ctx, admin.ID, v.ID)
		if e != nil || !reflect.DeepEqual(current.Stock, v.Stock) || current.Status != v.Status || current.Sequence != v.Sequence || current.Version != v.Version || !current.UpdatedAt.Equal(v.UpdatedAt) {
			t.Fatal("filter changed equipment")
		}
	}
	if history() != before {
		t.Fatal("filter changed history or receipts")
	}
}
