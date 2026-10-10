package postgres

import (
	"context"
	"encoding/json"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
	"time"
)

type BorrowingRepository struct{ pool *pgxpool.Pool }

func NewBorrowingRepository(p *pgxpool.Pool) *BorrowingRepository { return &BorrowingRepository{p} }
func (r *BorrowingRepository) executor(c context.Context) executor {
	if tx, ok := database.TxFromContext(c); ok {
		return tx
	}
	return r.pool
}
func (r *BorrowingRepository) Clock(c context.Context) (v time.Time, e error) {
	e = r.executor(c).QueryRow(c, `SELECT clock_timestamp()`).Scan(&v)
	return v, accountError(e)
}

const borrowingColumns = `id,reference,borrower_id,borrower_name,borrower_type,student_id,status,entry_path,acceptance_id,created_at,expires_at,checked_out_at,due_at,terminal_at,denial_reason,completed_at,fine_final_minor`

func scanBorrowing(s scanner) (v d.Record, e error) {
	e = s.Scan(&v.ID, &v.Reference, &v.BorrowerID, &v.BorrowerName, &v.BorrowerType, &v.StudentID, &v.Status, &v.EntryPath, &v.AcceptanceID, &v.CreatedAt, &v.ExpiresAt, &v.CheckedOutAt, &v.DueAt, &v.TerminalAt, &v.Reason, &v.CompletedAt, &v.FinalMinor)
	v.Items = []d.Item{}
	v.Events = []d.Event{}
	return v, accountError(e)
}
func (r *BorrowingRepository) Get(c context.Context, id uuid.UUID, lock bool) (d.Record, error) {
	return r.get(c, id, lock, 0)
}
func (r *BorrowingRepository) get(c context.Context, id uuid.UUID, lock bool, limit int) (v d.Record, e error) {
	suffix := ""
	if lock {
		// Lifecycle writes never change the borrowing primary key. This still
		// serializes competing transitions, but permits notification FK KEY SHARE
		// checks; FOR UPDATE creates a loan/account lock inversion with fan-out.
		suffix = " FOR NO KEY UPDATE"
	}
	v, e = scanBorrowing(r.executor(c).QueryRow(c, `SELECT `+borrowingColumns+` FROM borrowings WHERE id=$1`+suffix, id))
	if e != nil {
		return v, e
	}
	rows, e := r.executor(c).Query(c, `SELECT id,equipment_id,name,quantity,reserved_quantity,issued_quantity FROM borrowing_items WHERE borrowing_id=$1 ORDER BY equipment_id`, id)
	if e != nil {
		return v, accountError(e)
	}
	for rows.Next() {
		var i d.Item
		if e = rows.Scan(&i.ID, &i.EquipmentID, &i.Name, &i.Quantity, &i.Reserved, &i.Issued); e != nil {
			rows.Close()
			return v, accountError(e)
		}
		v.Items = append(v.Items, i)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return v, accountError(e)
	}
	if e = r.executor(c).QueryRow(c, `SELECT count(*) FROM borrowing_events WHERE borrowing_id=$1`, id).Scan(&v.EventCount); e != nil {
		return v, accountError(e)
	}
	eventLimit := 2147483647
	eventOrder := "occurred_at,id"
	if limit > 0 {
		eventLimit = limit
		eventOrder = "occurred_at DESC,id DESC"
	}
	rows, e = r.executor(c).Query(c, `SELECT id,actor_id,kind,reason,occurred_at FROM borrowing_events WHERE borrowing_id=$1 ORDER BY `+eventOrder+` LIMIT $2`, id, eventLimit)
	if e != nil {
		return v, accountError(e)
	}
	defer rows.Close()
	for rows.Next() {
		var i d.Event
		if e = rows.Scan(&i.ID, &i.ActorID, &i.Kind, &i.Reason, &i.At); e != nil {
			return v, accountError(e)
		}
		v.Events = append(v.Events, i)
	}
	if e = rows.Err(); e != nil {
		return v, accountError(e)
	}
	rows.Close()
	if limit > 0 {
		for left, right := 0, len(v.Events)-1; left < right; left, right = left+1, right-1 {
			v.Events[left], v.Events[right] = v.Events[right], v.Events[left]
		}
		e = r.previewAccountability(c, &v)
	} else {
		e = r.accountabilityDetails(c, &v)
	}
	return v, e
}
func (r *BorrowingRepository) List(c context.Context, f d.Filter) (p d.Page, e error) {
	p.Items = []d.Record{}
	p.Page = f.Page
	p.PerPage = f.PerPage
	search := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(f.Search)
	args := []any{f.Owner, f.Status, "%" + search + "%"}
	where := ` FROM borrowings WHERE ($1::uuid IS NULL OR borrower_id=$1) AND ($2='' OR status=$2) AND (reference ILIKE $3 OR borrower_name ILIKE $3 OR student_id ILIKE $3)`
	// One statement gives count and page a consistent snapshot under Read Committed.
	var data []byte
	e = r.executor(c).QueryRow(c, `WITH tick AS MATERIALIZED (SELECT clock_timestamp() AS at) SELECT (SELECT count(*)`+where+`),COALESCE((SELECT jsonb_agg(to_jsonb(b) ORDER BY b.created_at DESC,b.id DESC) FROM (SELECT `+borrowingColumns+`,jsonb_build_object('assessed_minor',COALESCE(fine_final_minor,borrowing_fine_amount(due_at,(SELECT at FROM tick))),'cleared_minor',(SELECT COALESCE(sum(cleared_minor),0) FROM fine_clearances f WHERE f.borrowing_id=borrowings.id),'outstanding_minor',COALESCE(fine_final_minor,borrowing_fine_amount(due_at,(SELECT at FROM tick)))-(SELECT COALESCE(sum(cleared_minor),0) FROM fine_clearances f WHERE f.borrowing_id=borrowings.id),'is_final',completed_at IS NOT NULL,'as_of',(SELECT at FROM tick)) AS fine,COALESCE((SELECT jsonb_agg(to_jsonb(i) ORDER BY i.equipment_id) FROM(SELECT id,equipment_id,name,quantity,reserved_quantity,issued_quantity FROM borrowing_items WHERE borrowing_id=borrowings.id ORDER BY equipment_id LIMIT 100)i),'[]'::jsonb) AS items`+where+` ORDER BY created_at DESC,id DESC LIMIT $4 OFFSET $5) b),'[]'::jsonb)`, append(args, f.PerPage, (f.Page-1)*f.PerPage)...).Scan(&p.Total, &data)
	if e != nil {
		return p, accountError(e)
	}
	e = json.Unmarshal(data, &p.Items)
	for n := range p.Items {
		if p.Items[n].Items == nil {
			p.Items[n].Items = []d.Item{}
		}
		p.Items[n].Events = []d.Event{}
	}
	return p, e
}
func (r *BorrowingRepository) Insert(c context.Context, v d.Record) error {
	_, e := r.executor(c).Exec(c, `INSERT INTO borrowings(id,reference,borrower_id,borrower_name,borrower_type,student_id,status,entry_path,acceptance_id,created_at,expires_at,checked_out_at,due_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`, v.ID, v.Reference, v.BorrowerID, v.BorrowerName, v.BorrowerType, v.StudentID, v.Status, v.EntryPath, v.AcceptanceID, v.CreatedAt, v.ExpiresAt, v.CheckedOutAt, v.DueAt)
	if e != nil {
		return accountError(e)
	}
	for _, i := range v.Items {
		_, e = r.executor(c).Exec(c, `INSERT INTO borrowing_items(id,borrowing_id,equipment_id,name,quantity,reserved_quantity,issued_quantity) VALUES($1,$2,$3,$4,$5,$6,$7)`, i.ID, v.ID, i.EquipmentID, i.Name, i.Quantity, i.Reserved, i.Issued)
		if e != nil {
			return accountError(e)
		}
	}
	return nil
}
func (r *BorrowingRepository) Transition(c context.Context, v d.Record) error {
	tag, e := r.executor(c).Exec(c, `UPDATE borrowings SET status=$2,checked_out_at=$3,due_at=$4,terminal_at=$5,denial_reason=$6 WHERE id=$1 AND status='PENDING'`, v.ID, v.Status, v.CheckedOutAt, v.DueAt, v.TerminalAt, v.Reason)
	if e != nil {
		return accountError(e)
	}
	if tag.RowsAffected() != 1 {
		return d.ErrState
	}
	for _, i := range v.Items {
		_, e = r.executor(c).Exec(c, `UPDATE borrowing_items SET reserved_quantity=$2,issued_quantity=$3 WHERE id=$1`, i.ID, i.Reserved, i.Issued)
		if e != nil {
			return accountError(e)
		}
	}
	return nil
}
func (r *BorrowingRepository) Effect(c context.Context, v d.Record, i d.Item, actor *uuid.UUID, kind string, equip inventory.Equipment, before inventory.Stock) error {
	tag, e := r.executor(c).Exec(c, `UPDATE equipment SET available=$2,reserved=$3,checked_out=$4,stock_sequence=$5,updated_at=clock_timestamp() WHERE id=$1 AND stock_sequence=$6 AND status<>'ARCHIVED'`, equip.ID, equip.Stock.Available, equip.Stock.Reserved, equip.Stock.CheckedOut, equip.Sequence, equip.Sequence-1)
	if e != nil {
		return accountError(e)
	}
	if tag.RowsAffected() != 1 {
		return d.ErrStock
	}
	at := v.CreatedAt
	if v.CheckedOutAt != nil {
		at = *v.CheckedOutAt
	}
	if v.TerminalAt != nil {
		at = *v.TerminalAt
	}
	_, e = r.executor(c).Exec(c, `INSERT INTO inventory_movements(id,equipment_id,sequence,actor_id,kind,borrowing_item_id,delta_available,delta_reserved,delta_checked_out,delta_damaged_held,delta_total,after_available,after_reserved,after_checked_out,after_damaged_held,after_total,reason,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,0,0,$10,$11,$12,$13,$14,$15,$16)`, uuid.New(), equip.ID, equip.Sequence, actor, kind, i.ID, equip.Stock.Available-before.Available, equip.Stock.Reserved-before.Reserved, equip.Stock.CheckedOut-before.CheckedOut, equip.Stock.Available, equip.Stock.Reserved, equip.Stock.CheckedOut, equip.Stock.DamagedHeld, equip.Stock.Total, "Borrowing "+v.Reference+" — "+kind, at)
	if e != nil {
		return accountError(e)
	}
	if actor != nil {
		// Immutable physical-operation audit retains issue-time catalog/borrower
		// labels separately from the original request snapshots.
		var borrowerName string
		if e = r.executor(c).QueryRow(c, `SELECT name FROM users WHERE id=$1`, v.BorrowerID).Scan(&borrowerName); e != nil {
			return accountError(e)
		}
		after, e := json.Marshal(map[string]any{"borrowing_id": v.ID, "borrowing_item_id": i.ID, "borrower_id": v.BorrowerID, "borrower_name": borrowerName, "equipment_name": equip.Name, "quantity": i.Quantity, "due_at": v.DueAt, "stock": equip.Stock})
		if e != nil {
			return e
		}
		_, e = r.executor(c).Exec(c, `INSERT INTO inventory_audit_events(id,actor_id,equipment_id,action,after_data,created_at) VALUES($1,$2,$3,$4,$5,$6)`, uuid.New(), actor, equip.ID, "BORROWING_"+kind, after, at)
		return accountError(e)
	}
	return nil
}
func (r *BorrowingRepository) Event(c context.Context, id uuid.UUID, v d.Event) error {
	_, e := r.executor(c).Exec(c, `INSERT INTO borrowing_events(id,borrowing_id,actor_id,kind,reason,occurred_at) VALUES($1,$2,$3,$4,$5,$6)`, v.ID, id, v.ActorID, v.Kind, v.Reason, v.At)
	return accountError(e)
}
func (r *BorrowingRepository) ReadReceipt(c context.Context, actor uuid.UUID, op, key, hash string) ([]byte, error) {
	var h string
	var data []byte
	e := r.executor(c).QueryRow(c, `SELECT payload_hash,result FROM borrowing_operation_receipts WHERE actor_id=$1 AND operation=$2 AND key=$3`, actor, op, key).Scan(&h, &data)
	if e != nil {
		return nil, accountError(e)
	}
	if h != hash {
		return nil, d.ErrKey
	}
	return data, nil
}
func (r *BorrowingRepository) WriteReceipt(c context.Context, actor uuid.UUID, op, key, hash string, data []byte) error {
	_, e := r.executor(c).Exec(c, `INSERT INTO borrowing_operation_receipts(actor_id,operation,key,payload_hash,result) VALUES($1,$2,$3,$4,$5)`, actor, op, key, hash, data)
	return accountError(e)
}
func (r *BorrowingRepository) Due(c context.Context, limit int) (ids []uuid.UUID, e error) {
	ids = []uuid.UUID{}
	rows, e := r.executor(c).Query(c, `SELECT id FROM borrowings WHERE status='PENDING' AND expires_at<=clock_timestamp() ORDER BY expires_at,id LIMIT $1`, limit)
	if e != nil {
		return nil, accountError(e)
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		if e = rows.Scan(&id); e != nil {
			return nil, accountError(e)
		}
		ids = append(ids, id)
	}
	return ids, accountError(rows.Err())
}

var _ d.Repository = (*BorrowingRepository)(nil)
