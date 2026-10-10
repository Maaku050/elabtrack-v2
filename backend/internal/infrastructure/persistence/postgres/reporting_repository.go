package postgres

import (
	"context"
	"encoding/json"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/reporting"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportingRepository struct{ pool *pgxpool.Pool }

func NewReportingRepository(p *pgxpool.Pool) *ReportingRepository { return &ReportingRepository{p} }
func (r *ReportingRepository) executor(c context.Context) executor {
	if t, ok := database.TxFromContext(c); ok {
		return t
	}
	return r.pool
}

var reportSources = map[string]string{
	"inventory":       `SELECT e.id id,e.created_at at,e.id equipment_id,e.category_id category_id,NULL::uuid borrower_id,e.status status,e.name||' '||COALESCE(c.name,'') search,ARRAY[COALESCE((e.id)::text,''),COALESCE((e.name)::text,''),COALESCE((COALESCE(c.name,'Uncategorized'))::text,''),COALESCE((e.status)::text,''),COALESCE((e.available)::text,''),COALESCE((e.reserved)::text,''),COALESCE((e.checked_out)::text,''),COALESCE((e.damaged_held)::text,''),COALESCE((e.total_tracked)::text,''),COALESCE((COALESCE(to_char(e.created_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,'')] cells FROM equipment e LEFT JOIN equipment_categories c ON c.id=e.category_id WHERE true`,
	"movements":       `SELECT m.id id,m.created_at at,e.id equipment_id,e.category_id category_id,b.borrower_id borrower_id,m.kind status,e.name||' '||m.reason||' '||COALESCE(b.reference,'') search,ARRAY[COALESCE((m.id)::text,''),COALESCE((e.name)::text,''),COALESCE((m.kind)::text,''),COALESCE((m.delta_available)::text,''),COALESCE((m.delta_reserved)::text,''),COALESCE((m.delta_checked_out)::text,''),COALESCE((m.delta_damaged_held)::text,''),COALESCE((m.delta_total)::text,''),COALESCE((m.after_available)::text,''),COALESCE((m.after_reserved)::text,''),COALESCE((m.after_checked_out)::text,''),COALESCE((m.after_damaged_held)::text,''),COALESCE((m.after_total)::text,''),COALESCE((m.actor_id)::text,''),COALESCE((m.reason)::text,''),COALESCE((COALESCE(to_char(m.created_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,'')] cells FROM inventory_movements m JOIN equipment e ON e.id=m.equipment_id LEFT JOIN borrowing_items i ON i.id=m.borrowing_item_id LEFT JOIN borrowings b ON b.id=i.borrowing_id WHERE true`,
	"requests":        `SELECT b.id id,b.created_at at,NULL::uuid equipment_id,NULL::uuid category_id,b.borrower_id borrower_id,b.status status,b.reference||' '||u.name||' '||COALESCE((SELECT string_agg(i.name,' ') FROM borrowing_items i WHERE i.borrowing_id=b.id),'') search,ARRAY[COALESCE((b.id)::text,''),COALESCE((b.reference)::text,''),COALESCE((u.name)::text,''),COALESCE((b.borrower_id)::text,''),COALESCE((b.status)::text,''),COALESCE((b.entry_path)::text,''),COALESCE((COALESCE(to_char(b.created_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.checked_out_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.due_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.completed_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,'')] cells FROM borrowings b JOIN users u ON u.id=b.borrower_id WHERE b.entry_path='REQUEST'`,
	"issued":          `SELECT b.id id,b.checked_out_at at,NULL::uuid equipment_id,NULL::uuid category_id,b.borrower_id borrower_id,b.status status,b.reference||' '||u.name||' '||COALESCE((SELECT string_agg(i.name,' ') FROM borrowing_items i WHERE i.borrowing_id=b.id),'') search,ARRAY[COALESCE((b.id)::text,''),COALESCE((b.reference)::text,''),COALESCE((u.name)::text,''),COALESCE((b.borrower_id)::text,''),COALESCE((b.status)::text,''),COALESCE((b.entry_path)::text,''),COALESCE((COALESCE(to_char(b.created_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.checked_out_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.due_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.completed_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,'')] cells FROM borrowings b JOIN users u ON u.id=b.borrower_id WHERE b.checked_out_at IS NOT NULL`,
	"active":          `SELECT b.id id,b.checked_out_at at,NULL::uuid equipment_id,NULL::uuid category_id,b.borrower_id borrower_id,b.status status,b.reference||' '||u.name||' '||COALESCE((SELECT string_agg(i.name,' ') FROM borrowing_items i WHERE i.borrowing_id=b.id),'') search,ARRAY[COALESCE((b.id)::text,''),COALESCE((b.reference)::text,''),COALESCE((u.name)::text,''),COALESCE((b.borrower_id)::text,''),COALESCE((b.status)::text,''),COALESCE((b.entry_path)::text,''),COALESCE((COALESCE(to_char(b.created_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.checked_out_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.due_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.completed_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,'')] cells FROM borrowings b JOIN users u ON u.id=b.borrower_id WHERE b.status='CHECKED_OUT'`,
	"overdue":         `SELECT b.id id,b.due_at at,NULL::uuid equipment_id,NULL::uuid category_id,b.borrower_id borrower_id,b.status status,b.reference||' '||u.name||' '||COALESCE((SELECT string_agg(i.name,' ') FROM borrowing_items i WHERE i.borrowing_id=b.id),'') search,ARRAY[COALESCE((b.id)::text,''),COALESCE((b.reference)::text,''),COALESCE((u.name)::text,''),COALESCE((b.borrower_id)::text,''),COALESCE((b.status)::text,''),COALESCE((b.entry_path)::text,''),COALESCE((COALESCE(to_char(b.created_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.checked_out_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.due_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.completed_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,'')] cells FROM borrowings b JOIN users u ON u.id=b.borrower_id WHERE b.status='CHECKED_OUT' AND b.due_at<statement_timestamp()`,
	"returns":         `SELECT l.id id,l.occurred_at at,e.id equipment_id,e.category_id category_id,b.borrower_id borrower_id,CASE WHEN l.lost>0 AND l.damaged>0 THEN 'MIXED' WHEN l.lost>0 THEN 'LOST' WHEN l.damaged>0 THEN 'DAMAGED' ELSE 'GOOD' END status,b.reference||' '||e.name||' '||l.reason search,ARRAY[COALESCE((l.id)::text,''),COALESCE((b.reference)::text,''),COALESCE((e.name)::text,''),COALESCE((l.good)::text,''),COALESCE((l.damaged)::text,''),COALESCE((l.lost)::text,''),COALESCE((l.actor_id)::text,''),COALESCE((l.reason)::text,''),COALESCE((COALESCE(to_char(l.occurred_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,'')] cells FROM return_lines l JOIN equipment e ON e.id=l.equipment_id JOIN borrowings b ON b.id=l.borrowing_id WHERE true`,
	"incidents":       `SELECT l.id id,l.occurred_at at,e.id equipment_id,e.category_id category_id,b.borrower_id borrower_id,CASE WHEN l.lost>0 AND l.damaged>0 THEN 'MIXED' WHEN l.lost>0 THEN 'LOST' WHEN l.damaged>0 THEN 'DAMAGED' ELSE 'GOOD' END status,b.reference||' '||e.name||' '||l.reason search,ARRAY[COALESCE((l.id)::text,''),COALESCE((b.reference)::text,''),COALESCE((e.name)::text,''),COALESCE((l.damaged)::text,''),COALESCE((l.lost)::text,''),COALESCE((l.actor_id)::text,''),COALESCE((l.reason)::text,''),COALESCE((COALESCE(to_char(l.occurred_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,'')] cells FROM return_lines l JOIN equipment e ON e.id=l.equipment_id JOIN borrowings b ON b.id=l.borrowing_id WHERE l.damaged+l.lost>0`,
	"replacements":    `SELECT o.id id,l.occurred_at at,e.id equipment_id,e.category_id category_id,b.borrower_id borrower_id,o.kind status,b.reference||' '||e.name search,ARRAY[COALESCE((o.id)::text,''),COALESCE((b.reference)::text,''),COALESCE((e.name)::text,''),COALESCE((o.kind)::text,''),COALESCE((o.required)::text,''),COALESCE((COALESCE((SELECT sum(quantity) FROM replacement_acceptances a WHERE a.obligation_id=o.id),0))::text,''),COALESCE((o.required-COALESCE((SELECT sum(quantity) FROM replacement_acceptances a WHERE a.obligation_id=o.id),0))::text,''),COALESCE((COALESCE(to_char(l.occurred_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,'')] cells FROM replacement_obligations o JOIN return_lines l ON l.id=o.return_line_id JOIN equipment e ON e.id=o.equipment_id JOIN borrowings b ON b.id=o.borrowing_id WHERE true`,
	"fines":           `SELECT b.id id,b.due_at at,NULL::uuid equipment_id,NULL::uuid category_id,b.borrower_id borrower_id,b.status status,b.reference||' '||u.name||' '||COALESCE((SELECT string_agg(i.name,' ') FROM borrowing_items i WHERE i.borrowing_id=b.id),'') search,ARRAY[COALESCE((b.id)::text,''),COALESCE((b.reference)::text,''),COALESCE((u.name)::text,''),COALESCE((b.status)::text,''),COALESCE((COALESCE(to_char(b.due_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(to_char(b.completed_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,''),COALESCE((COALESCE(b.fine_final_minor,borrowing_fine_amount(b.due_at,statement_timestamp()),0))::text,''),COALESCE((COALESCE((SELECT sum(cleared_minor) FROM fine_clearances fc WHERE fc.borrowing_id=b.id),0))::text,''),COALESCE((COALESCE(b.fine_final_minor,borrowing_fine_amount(b.due_at,statement_timestamp()),0)-COALESCE((SELECT sum(cleared_minor) FROM fine_clearances fc WHERE fc.borrowing_id=b.id),0))::text,''),COALESCE((b.status='COMPLETED')::text,''),COALESCE((COALESCE((SELECT string_agg(method,', ' ORDER BY occurred_at,id) FROM(SELECT method,occurred_at,id FROM fine_clearances fc WHERE fc.borrowing_id=b.id ORDER BY occurred_at DESC,id DESC LIMIT 25)fc),''))::text,'')] cells FROM borrowings b JOIN users u ON u.id=b.borrower_id WHERE b.checked_out_at IS NOT NULL`,
	"fine-clearances": `SELECT c.id id,c.occurred_at at,NULL::uuid equipment_id,NULL::uuid category_id,b.borrower_id borrower_id,c.method status,b.reference||' '||u.name||' '||c.note search,ARRAY[c.id::text,b.reference,u.name,c.assessed_minor::text,c.cleared_minor::text,c.method,c.note,c.actor_id::text,to_char(c.occurred_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"')] cells FROM fine_clearances c JOIN borrowings b ON b.id=c.borrowing_id JOIN users u ON u.id=b.borrower_id`,
	"account-audit":   `SELECT a.id id,a.occurred_at at,NULL::uuid equipment_id,NULL::uuid category_id,NULL::uuid borrower_id,a.action status,u.name||' '||a.action search,ARRAY[COALESCE((a.id)::text,''),COALESCE((a.account_id)::text,''),COALESCE((u.name)::text,''),COALESCE((a.actor_id)::text,''),COALESCE((a.action)::text,''),COALESCE((COALESCE(to_char(a.occurred_at AT TIME ZONE 'UTC','YYYY-MM-DD"T"HH24:MI:SS.US"Z"'),''))::text,'')] cells FROM account_audit_events a JOIN users u ON u.id=a.account_id WHERE true`,
}

func (r *ReportingRepository) Report(c context.Context, key string, f d.Filter, limit int) (p d.Page, e error) {
	def, ok := d.Find(key)
	if !ok {
		return p, fmtInvalidReport()
	}
	p.Kind = key
	p.Columns = def.Columns
	p.Page = f.Page
	p.PerPage = limit
	p.Rows = [][]string{}
	var raw []byte
	loan := key == "requests" || key == "issued" || key == "active" || key == "overdue" || key == "fines"
	custodyFilter := ` AND ($5::uuid IS NULL OR equipment_id=$5) AND ($6::uuid IS NULL OR category_id=$6)`
	if loan {
		custodyFilter = ` AND ($5::uuid IS NULL OR EXISTS(SELECT 1 FROM borrowing_items i WHERE i.borrowing_id=s.id AND i.equipment_id=$5)) AND ($6::uuid IS NULL OR EXISTS(SELECT 1 FROM borrowing_items i JOIN equipment e ON e.id=i.equipment_id WHERE i.borrowing_id=s.id AND e.category_id=$6))`
	}
	source := reportSources[key]
	if f.DueToday {
		source += ` AND (b.due_at AT TIME ZONE 'Asia/Manila')::date=(statement_timestamp() AT TIME ZONE 'Asia/Manila')::date`
	}
	query := `WITH source AS (` + source + `),filtered AS(SELECT * FROM source s WHERE ($1='' OR strpos(lower(search),lower($1))>0) AND ($2='' OR status=$2) AND ($3::timestamptz IS NULL OR at>=$3) AND ($4::timestamptz IS NULL OR at<$4) AND ($7::uuid IS NULL OR borrower_id=$7)` + custodyFilter + `) SELECT (SELECT count(*) FROM filtered),COALESCE((SELECT jsonb_agg(to_jsonb(v.cells) ORDER BY v.at DESC,v.id DESC) FROM(SELECT * FROM filtered ORDER BY at DESC,id DESC LIMIT $8 OFFSET $9)v),'[]'::jsonb),statement_timestamp()`
	e = r.executor(c).QueryRow(c, query, f.Search, f.Status, f.From, f.To, f.EquipmentID, f.CategoryID, f.BorrowerID, limit, (f.Page-1)*limit).Scan(&p.Total, &raw, &p.AsOf)
	if e != nil {
		return p, accountError(e)
	}
	e = json.Unmarshal(raw, &p.Rows)
	return
}
func fmtInvalidReport() error { return shared.ErrInvalidInput }
func (r *ReportingRepository) Dashboard(c context.Context, id uuid.UUID, role user.Role, days int) (v d.Dashboard, e error) {
	var metrics, recent, trends []byte
	e = r.executor(c).QueryRow(c, `WITH scoped AS(SELECT b.*,COALESCE(b.fine_final_minor,borrowing_fine_amount(b.due_at,statement_timestamp()),0) assessed,COALESCE((SELECT sum(cleared_minor) FROM fine_clearances f WHERE f.borrowing_id=b.id),0) cleared FROM borrowings b WHERE $2<>'BORROWER' OR b.borrower_id=$1) SELECT jsonb_build_object('equipment_types',(SELECT count(*) FROM equipment WHERE ($2<>'BORROWER' OR status='ACTIVE')),'available_units',(SELECT COALESCE(sum(available),0) FROM equipment WHERE ($2<>'BORROWER' OR status='ACTIVE')),'reserved_units',(SELECT COALESCE(sum(reserved),0) FROM equipment),'checked_out_units',(SELECT COALESCE(sum(checked_out),0) FROM equipment),'damaged_held_units',(SELECT COALESCE(sum(damaged_held),0) FROM equipment),'total_tracked_units',(SELECT COALESCE(sum(total_tracked),0) FROM equipment),'active_borrowers',(SELECT count(*) FROM users WHERE role='BORROWER' AND is_active AND NOT activation_required AND ($2<>'BORROWER' OR id=$1)),'pending_requests',(SELECT count(*) FROM scoped WHERE status='PENDING'),'active_loans',(SELECT count(*) FROM scoped WHERE status='CHECKED_OUT'),'overdue_loans',(SELECT count(*) FROM scoped WHERE status='CHECKED_OUT' AND due_at<statement_timestamp()),'due_today',(SELECT count(*) FROM scoped WHERE status='CHECKED_OUT' AND (due_at AT TIME ZONE 'Asia/Manila')::date=(statement_timestamp() AT TIME ZONE 'Asia/Manila')::date),'completed_loans',(SELECT count(*) FROM scoped WHERE status='COMPLETED'),'returned_good_units',(SELECT COALESCE(sum(l.good),0) FROM return_lines l JOIN scoped b ON b.id=l.borrowing_id),'outstanding_replacements',(SELECT COALESCE(sum(o.required-COALESCE(a.accepted,0)),0) FROM replacement_obligations o JOIN scoped b ON b.id=o.borrowing_id LEFT JOIN (SELECT obligation_id,sum(quantity) accepted FROM replacement_acceptances GROUP BY obligation_id)a ON a.obligation_id=o.id),'assessed_fines_minor',(SELECT COALESCE(sum(assessed),0) FROM scoped),'cleared_fines_minor',(SELECT COALESCE(sum(cleared),0) FROM scoped),'outstanding_fines_minor',(SELECT COALESCE(sum(assessed-cleared),0) FROM scoped),'outstanding_settlements',(SELECT count(*) FROM scoped WHERE assessed>cleared)), COALESCE((SELECT jsonb_agg(to_jsonb(v) ORDER BY v.at DESC,v.id DESC) FROM(SELECT e.id,e.borrowing_id,b.reference,e.kind,e.occurred_at at FROM borrowing_events e JOIN scoped b ON b.id=e.borrowing_id ORDER BY e.occurred_at DESC,e.id DESC LIMIT 10)v),'[]'::jsonb) , COALESCE((SELECT jsonb_agg(to_jsonb(t) ORDER BY t."day") FROM(
 SELECT to_char(bucket,'YYYY-MM-DD') AS "day",
 (SELECT count(*) FROM borrowing_events e JOIN scoped b ON b.id=e.borrowing_id WHERE e.kind IN ('CHECKED_OUT','DIRECT') AND (e.occurred_at AT TIME ZONE 'Asia/Manila')::date=bucket::date) issued,
 (SELECT count(*) FROM scoped b WHERE b.status='PENDING' AND (b.created_at AT TIME ZONE 'Asia/Manila')::date=bucket::date) pending,
 (SELECT count(*) FROM borrowing_events e JOIN scoped b ON b.id=e.borrowing_id WHERE e.kind='DENIED' AND (e.occurred_at AT TIME ZONE 'Asia/Manila')::date=bucket::date) denied
 FROM generate_series((statement_timestamp() AT TIME ZONE 'Asia/Manila')::date-($3::int-1),(statement_timestamp() AT TIME ZONE 'Asia/Manila')::date,'1 day'::interval) AS g(bucket))t),'[]'::jsonb),statement_timestamp()`, id, string(role), days).Scan(&metrics, &recent, &trends, &v.AsOf)
	if e != nil {
		return v, accountError(e)
	}
	if e = json.Unmarshal(metrics, &v.Metrics); e != nil {
		return v, e
	}
	if e = json.Unmarshal(trends, &v.Trends); e != nil {
		return v, e
	}
	e = json.Unmarshal(recent, &v.Recent)
	if role == user.RoleBorrower {
		for _, k := range []string{"reserved_units", "checked_out_units", "damaged_held_units", "total_tracked_units", "active_borrowers"} {
			delete(v.Metrics, k)
		}
	}
	return
}

var _ d.Repository = (*ReportingRepository)(nil)
