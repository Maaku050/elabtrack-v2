package integration

import (
	"encoding/csv"
	"errors"
	bapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/borrowing"
	iapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/inventory"
	app "github.com/Maaku050/elabtrack-v2/backend/internal/application/reporting"
	tapp "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	b "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	i "github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/reporting"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/catalogimage"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
	"strings"
	"testing"
	"time"
)

func TestRealReporting(t *testing.T) {
	ctx, db, owner, _ := batchDatabase(t)
	tx := database.NewTxManager(db.Pool)
	ar := postgres.NewAccountsRepository(db.Pool)
	ir := postgres.NewInventoryRepository(db.Pool)
	br := postgres.NewBorrowingRepository(db.Pool)
	rr := postgres.NewReportingRepository(db.Pool)
	s := app.NewService(rr, ar, tx)
	admin, staff, borrower, other := termsUser(t, ctx, db, user.RoleAdmin), termsUser(t, ctx, db, user.RoleStaff), termsUser(t, ctx, db, user.RoleBorrower), termsUser(t, ctx, db, user.RoleBorrower)
	_, e := owner.Pool.Exec(ctx, `INSERT INTO borrower_profiles(user_id,borrower_type)VALUES($1,'FACULTY'),($2,'FACULTY')`, borrower.ID, other.ID)
	require(t, e == nil, "fictional profiles")
	ts := tapp.NewService(postgres.NewTermsRepository(db.Pool), postgres.NewUserRepository(db.Pool), tx)
	v, e := ts.Current(ctx, admin.ID)
	require(t, e == nil, "synthetic terms")
	_, e = ts.Accept(ctx, borrower.ID, v.ID)
	require(t, e == nil, "explicit fixture acceptance")
	inv := iapp.NewService(ir, ar, tx, catalogimage.Validator{})
	marker := uuid.NewString()
	eq, e := inv.Create(ctx, staff.ID, uuid.NewString(), i.Create{Metadata: i.Metadata{Name: "=TEST \"Reporting\", " + marker}, Opening: 20, Reason: "TEST reporting stock"})
	require(t, e == nil, "formula-name equipment")
	eq2, e := inv.Create(ctx, staff.ID, uuid.NewString(), i.Create{Metadata: i.Metadata{Name: "TEST Reporting second " + marker}, Opening: 10, Reason: "TEST second stock"})
	require(t, e == nil, "second equipment")
	past := time.Now().Add(-49 * time.Hour).UTC().Truncate(time.Microsecond)
	due := past.Add(time.Hour)
	bs := bapp.NewService(accountabilityClock{br, past}, ar, ir, ts, tx)
	loan, e := bs.Direct(ctx, staff.ID, uuid.NewString(), b.Input{BorrowerID: borrower.ID, Items: []b.Line{{EquipmentID: eq.ID, Quantity: 3}, {EquipmentID: eq2.ID, Quantity: 1}}, DueAt: &due, Confirm: true, Handover: true})
	require(t, e == nil, "historical multi-item loan")
	current := bapp.NewService(br, ar, ir, ts, tx)
	var firstItem uuid.UUID
	for _, item := range loan.Items {
		if item.EquipmentID == eq.ID {
			firstItem = item.ID
		}
	}
	loan, e = current.Return(ctx, staff.ID, loan.ID, uuid.NewString(), b.ReturnInput{ExpectedEvents: len(loan.Events), Confirm: true, Reason: "TEST mixed verification", Lines: []b.ReturnLine{{ItemID: firstItem, Good: 1, Damaged: 1, Lost: 1}}})
	require(t, e == nil, "mixed disposition")
	f := d.Filter{Page: 1, PerPage: 25, BorrowerID: &borrower.ID}
	t.Run("all_reports_and_exact_shapes", func(t *testing.T) {
		for _, def := range d.Definitions {
			p, e := s.Report(ctx, admin.ID, def.Key, d.Filter{Page: 1, PerPage: 2})
			require(t, e == nil, "actual report "+def.Key)
			require(t, p.Rows != nil && len(p.Rows) <= 2 && p.Total >= int64(len(p.Rows)), "bounded result "+def.Key)
			for _, row := range p.Rows {
				require(t, len(row) == len(def.Columns), "headers align "+def.Key)
			}
		}
	})
	t.Run("aggregation_and_scope_reconcile_without_join_multiplication", func(t *testing.T) {
		mine, e := s.Dashboard(ctx, borrower.ID)
		require(t, e == nil, "borrower dashboard")
		require(t, mine.Metrics["active_loans"] == 1 && mine.Metrics["overdue_loans"] == 1 && mine.Metrics["outstanding_replacements"] == 2, "own loan and obligation counts")
		var assessed int64
		e = owner.Pool.QueryRow(ctx, `SELECT borrowing_fine_amount(due_at,statement_timestamp()) FROM borrowings WHERE id=$1`, loan.ID).Scan(&assessed)
		require(t, e == nil && mine.Metrics["assessed_fines_minor"] == assessed, "one assessment per loan")
		none, e := s.Dashboard(ctx, other.ID)
		require(t, e == nil && none.Metrics["active_loans"] == 0 && len(none.Recent) == 0, "other borrower has no foreign history")
		ops, e := s.Dashboard(ctx, staff.ID)
		require(t, e == nil, "ops metrics")
		var a, r, c, dam, total int64
		e = owner.Pool.QueryRow(ctx, `SELECT sum(available),sum(reserved),sum(checked_out),sum(damaged_held),sum(total_tracked) FROM equipment`).Scan(&a, &r, &c, &dam, &total)
		require(t, e == nil && ops.Metrics["available_units"] == a && ops.Metrics["reserved_units"] == r && ops.Metrics["checked_out_units"] == c && ops.Metrics["damaged_held_units"] == dam && ops.Metrics["total_tracked_units"] == total && total == a+r+c+dam, "inventory buckets once")
		p, e := s.Report(ctx, staff.ID, "issued", f)
		require(t, e == nil && p.Total == 1 && len(p.Rows) == 1, "multi-item loan counted once")
		f.EquipmentID = &eq.ID
		p, e = s.Report(ctx, staff.ID, "issued", f)
		require(t, e == nil && p.Total == 1, "equipment EXISTS filter")
		searchByEquipment := f
		searchByEquipment.Search = eq.Name
		p, e = s.Report(ctx, staff.ID, "issued", searchByEquipment)
		require(t, e == nil && p.Total == 1, "equipment name search retains one loan")
		f.EquipmentID = nil
	})
	t.Run("search_dates_pages_and_filtered_csv", func(t *testing.T) {
		p, e := s.Report(ctx, staff.ID, "inventory", d.Filter{Page: 1, PerPage: 1, Search: marker})
		require(t, e == nil && p.Total == 2 && len(p.Rows) == 1, "search before count")
		second, e := s.Report(ctx, staff.ID, "inventory", d.Filter{Page: 2, PerPage: 1, Search: marker})
		require(t, e == nil && second.Total == 2 && second.Rows[0][0] != p.Rows[0][0], "stable second page")
		empty, e := s.Report(ctx, staff.ID, "inventory", d.Filter{Page: 3, PerPage: 1, Search: marker})
		require(t, e == nil && empty.Total == 2 && len(empty.Rows) == 0, "out-of-range retains total")
		end := past
		closed := f
		closed.To = &end
		p, e = s.Report(ctx, staff.ID, "issued", closed)
		require(t, e == nil && p.Total == 0, "exclusive end boundary")
		closed.To = nil
		closed.From = &past
		p, e = s.Report(ctx, staff.ID, "issued", closed)
		require(t, e == nil && p.Total == 1, "inclusive start")
		csvBytes, e := s.Export(ctx, staff.ID, "inventory", d.Filter{Page: 1, PerPage: 25, EquipmentID: &eq.ID})
		require(t, e == nil, "filtered CSV")
		rows, e := csv.NewReader(strings.NewReader(string(csvBytes))).ReadAll()
		require(t, e == nil && len(rows) == 2 && rows[1][1] == "'"+eq.Name && rows[1][4] == "18", "real formula-safe quantities")
		require(t, strings.Contains(string(csvBytes), "\r\n"), "deterministic RFC4180")
	})
	t.Run("fine_full_balance_read_model_and_history", func(t *testing.T) {
		p, e := s.Report(ctx, staff.ID, "fines", f)
		require(t, e == nil && p.Total == 1 && p.Rows[0][7] == "0", "before clearance")
		_, e = current.ClearFine(ctx, admin.ID, loan.ID, uuid.NewString(), b.FineInput{Confirm: true, Method: "WAIVED", ExpectedMinor: loan.Fine.Outstanding})
		require(t, e == nil, "actual Admin clearance")
		p, e = s.Report(ctx, staff.ID, "fines", f)
		require(t, e == nil && p.Rows[0][8] == "0" && p.Rows[0][10] == "WAIVED", "live resolved fine projection")
		history, e := s.Report(ctx, staff.ID, "fine-clearances", f)
		require(t, e == nil && history.Total == 1 && history.Rows[0][5] == "WAIVED" && history.Rows[0][7] == admin.ID.String(), "immutable Admin fine resolution history")
		p, e = s.Report(ctx, staff.ID, "replacements", f)
		require(t, e == nil && p.Total == 2, "clearance preserves replacements")
	})
	t.Run("actual_export_bound_no_silent_truncation", func(t *testing.T) {
		action := "TEST_REPORT_LIMIT_" + uuid.NewString()
		_, e := owner.Pool.Exec(ctx, `INSERT INTO account_audit_events(id,actor_id,account_id,action)SELECT gen_random_uuid(),$1,$2,$3 FROM generate_series(1,5001)`, admin.ID, borrower.ID, action)
		require(t, e == nil, "synthetic immutable audit volume")
		f := d.Filter{Page: 1, PerPage: 25, Search: action}
		p, e := s.Report(ctx, admin.ID, "account-audit", f)
		require(t, e == nil && p.Total == 5001 && len(p.Rows) == 25, "bounded exact volume")
		blob, e := s.Export(ctx, admin.ID, "account-audit", f)
		require(t, errors.Is(e, d.ErrExportLimit) && blob == nil, "oversized export refuses rather than truncating")
		_, e = s.ExportData(ctx, admin.ID, "account-audit", f)
		require(t, errors.Is(e, d.ErrExportLimit), "full data also refuses truncation")
	})

	t.Run("complete_export_data_and_truthful_bounded_trends", func(t *testing.T) {
		var dueToday int64
		e := owner.Pool.QueryRow(ctx, `SELECT count(*) FROM borrowings WHERE status='CHECKED_OUT' AND (due_at AT TIME ZONE 'Asia/Manila')::date=(statement_timestamp() AT TIME ZONE 'Asia/Manila')::date`).Scan(&dueToday)
		require(t, e == nil, "actual Manila due-today count")
		page, e := s.ExportData(ctx, staff.ID, "active", d.Filter{Page: 2, PerPage: 1, DueToday: true})
		require(t, e == nil && page.Total == dueToday && int64(len(page.Rows)) == dueToday, "due-today drilldown/export/count agree")
		for _, def := range d.Definitions {
			page, e := s.ExportData(ctx, admin.ID, def.Key, d.Filter{Page: 2, PerPage: 1, Search: marker})
			require(t, e == nil && page.Page == 1 && int64(len(page.Rows)) == page.Total, "full-filter export, not server page: "+def.Key)
		}
		_, e = s.ExportData(ctx, borrower.ID, "inventory", d.Filter{Page: 1, PerPage: 25})
		require(t, errors.Is(e, shared.ErrForbidden), "borrower cannot export full data")
		for _, days := range []int{7, 30} {
			dashboard, e := s.Dashboard(ctx, admin.ID, days)
			if e != nil {
				t.Fatalf("dashboard SQL: %v", e)
			}
			require(t, len(dashboard.Trends) == days, "exact bounded daily series")
			require(t, dashboard.Metrics["due_today"] == dueToday, "due-today KPI agrees with exact report")
			var issued, denied int64
			for _, day := range dashboard.Trends {
				issued += day.Issued
				denied += day.Denied
				require(t, day.Pending >= 0 && day.Issued >= 0 && day.Denied >= 0, "nonnegative real analytics")
			}
			var wantIssued, wantDenied int64
			e = owner.Pool.QueryRow(ctx, "SELECT count(*) FILTER(WHERE kind IN ('CHECKED_OUT','DIRECT')), count(*) FILTER(WHERE kind='DENIED') FROM borrowing_events WHERE (occurred_at AT TIME ZONE 'Asia/Manila')::date >= (statement_timestamp() AT TIME ZONE 'Asia/Manila')::date-($1-1)", days).Scan(&wantIssued, &wantDenied)
			require(t, e == nil && issued == wantIssued && denied == wantDenied, "exact event reconciliation")
			require(t, dashboard.Metrics["available_units"]+dashboard.Metrics["reserved_units"]+dashboard.Metrics["checked_out_units"]+dashboard.Metrics["damaged_held_units"] == dashboard.Metrics["total_tracked_units"], "ring conserves physical stock")
		}
		_, e = s.Dashboard(ctx, admin.ID, 365)
		require(t, errors.Is(e, shared.ErrInvalidInput), "unbounded range rejected")
	})
	t.Run("current_role_account_status_and_filters", func(t *testing.T) {
		_, e := s.Report(ctx, borrower.ID, "inventory", d.Filter{Page: 1, PerPage: 25})
		require(t, errors.Is(e, shared.ErrForbidden), "borrower report forbidden")
		_, e = s.Report(ctx, staff.ID, "account-audit", d.Filter{Page: 1, PerPage: 25})
		require(t, errors.Is(e, shared.ErrForbidden), "Admin audit only")
		_, e = s.Report(ctx, admin.ID, "account-audit", f)
		require(t, errors.Is(e, shared.ErrInvalidInput), "unsupported borrower filter rejected")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET role='BORROWER' WHERE id=$1`, staff.ID)
		require(t, e == nil, "fixture demotion")
		_, e = s.Report(ctx, staff.ID, "inventory", d.Filter{Page: 1, PerPage: 25})
		require(t, errors.Is(e, shared.ErrForbidden), "live demotion")
		_, e = owner.Pool.Exec(ctx, `UPDATE users SET is_active=false WHERE id=$1`, other.ID)
		require(t, e == nil, "fixture deactivation")
		_, e = s.Dashboard(ctx, other.ID)
		require(t, errors.Is(e, shared.ErrUnauthorized), "live inactive dashboard")
	})
}
