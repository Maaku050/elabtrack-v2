package postgres

import (
	"context"
	"encoding/json"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/google/uuid"
	"time"
)

func (r *BorrowingRepository) accountabilityDetails(c context.Context, v *d.Record) error {
	v.Returns = []d.Disposition{}
	v.Obligations = []d.Obligation{}
	v.Replacements = []d.Replacement{}
	v.Clearances = []d.Clearance{}
	var data []byte
	e := r.executor(c).QueryRow(c, `SELECT jsonb_build_object('returns',COALESCE((SELECT jsonb_agg(to_jsonb(l) ORDER BY l.occurred_at,l.id) FROM return_lines l WHERE l.borrowing_id=$1),'[]'::jsonb),'obligations',COALESCE((SELECT jsonb_agg(to_jsonb(o)||jsonb_build_object('accepted',(SELECT COALESCE(sum(a.quantity),0) FROM replacement_acceptances a WHERE a.obligation_id=o.id)) ORDER BY o.id) FROM replacement_obligations o WHERE o.borrowing_id=$1),'[]'::jsonb),'replacements',COALESCE((SELECT jsonb_agg(to_jsonb(a) ORDER BY a.occurred_at,a.id) FROM replacement_acceptances a WHERE a.borrowing_id=$1),'[]'::jsonb),'clearances',COALESCE((SELECT jsonb_agg(to_jsonb(f) ORDER BY f.occurred_at,f.id) FROM fine_clearances f WHERE f.borrowing_id=$1),'[]'::jsonb))`, v.ID).Scan(&data)
	if e != nil {
		return accountError(e)
	}
	if e = json.Unmarshal(data, v); e != nil {
		return e
	}
	for n := range v.Items {
		i := &v.Items[n]
		for _, l := range v.Returns {
			if l.ItemID == i.ID {
				i.Good += l.Good
				i.Damaged += l.Damaged
				i.Lost += l.Lost
			}
		}
		i.Outstanding = i.Issued - i.Good - i.Damaged - i.Lost
	}
	return nil
}
func (r *BorrowingRepository) accountabilityEffect(c context.Context, v d.Record, itemID, sourceID uuid.UUID, actor uuid.UUID, kind, reason string, at time.Time, eq inventory.Equipment, before inventory.Stock) error {
	tag, e := r.executor(c).Exec(c, `UPDATE equipment SET available=$2,checked_out=$3,damaged_held=$4,total_tracked=$5,stock_sequence=$6,updated_at=$7 WHERE id=$1 AND stock_sequence=$8 AND status<>'ARCHIVED'`, eq.ID, eq.Stock.Available, eq.Stock.CheckedOut, eq.Stock.DamagedHeld, eq.Stock.Total, eq.Sequence, at, eq.Sequence-1)
	if e != nil {
		return accountError(e)
	}
	if tag.RowsAffected() != 1 {
		return d.ErrStock
	}
	var ret, rep *uuid.UUID
	if kind == "RETURN" {
		ret = &sourceID
	} else {
		rep = &sourceID
	}
	_, e = r.executor(c).Exec(c, `INSERT INTO inventory_movements(id,equipment_id,sequence,actor_id,kind,borrowing_item_id,return_line_id,replacement_acceptance_id,delta_available,delta_reserved,delta_checked_out,delta_damaged_held,delta_total,after_available,after_reserved,after_checked_out,after_damaged_held,after_total,reason,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,0,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19)`, uuid.New(), eq.ID, eq.Sequence, actor, kind, itemID, ret, rep, eq.Stock.Available-before.Available, eq.Stock.CheckedOut-before.CheckedOut, eq.Stock.DamagedHeld-before.DamagedHeld, eq.Stock.Total-before.Total, eq.Stock.Available, eq.Stock.Reserved, eq.Stock.CheckedOut, eq.Stock.DamagedHeld, eq.Stock.Total, reason, at)
	if e != nil {
		return accountError(e)
	}
	old, _ := json.Marshal(before)
	after, _ := json.Marshal(map[string]any{"stock": eq.Stock, "borrowing_id": v.ID, "source_id": sourceID, "equipment_name": eq.Name})
	_, e = r.executor(c).Exec(c, `INSERT INTO inventory_audit_events(id,actor_id,equipment_id,action,before_data,after_data,created_at) VALUES($1,$2,$3,$4,$5,$6,$7)`, uuid.New(), actor, eq.ID, kind, old, after, at)
	return accountError(e)
}
func (r *BorrowingRepository) RecordReturn(c context.Context, v d.Record, i d.Item, l d.Disposition, eq inventory.Equipment, before inventory.Stock) error {
	_, e := r.executor(c).Exec(c, `INSERT INTO return_lines(id,borrowing_id,equipment_id,item_id,good,damaged,lost,actor_id,reason,occurred_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, l.ID, v.ID, i.EquipmentID, i.ID, l.Good, l.Damaged, l.Lost, l.ActorID, l.Reason, l.At)
	if e != nil {
		return accountError(e)
	}
	for kind, q := range map[string]int64{"DAMAGED": l.Damaged, "LOST": l.Lost} {
		if q > 0 {
			_, e = r.executor(c).Exec(c, `INSERT INTO replacement_obligations(id,return_line_id,borrowing_id,item_id,equipment_id,kind,required) VALUES($1,$2,$3,$4,$5,$6,$7)`, uuid.New(), l.ID, v.ID, i.ID, i.EquipmentID, kind, q)
			if e != nil {
				return accountError(e)
			}
		}
	}
	return r.accountabilityEffect(c, v, i.ID, l.ID, l.ActorID, "RETURN", l.Reason, l.At, eq, before)
}
func (r *BorrowingRepository) RecordReplacement(c context.Context, v d.Record, o d.Obligation, a d.Replacement, eq inventory.Equipment, before inventory.Stock) error {
	_, e := r.executor(c).Exec(c, `INSERT INTO replacement_acceptances(id,obligation_id,borrowing_id,equipment_id,quantity,actor_id,reason,occurred_at,item_id,equivalent_confirmed,equipment_name) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,true,$10)`, a.ID, o.ID, v.ID, o.EquipmentID, a.Quantity, a.ActorID, a.Reason, a.At, o.ItemID, a.EquipmentName)
	if e != nil {
		return accountError(e)
	}
	return r.accountabilityEffect(c, v, o.ItemID, a.ID, a.ActorID, "REPLACEMENT", a.Reason, a.At, eq, before)
}
func (r *BorrowingRepository) Complete(c context.Context, v d.Record, at time.Time, amount int64) error {
	tag, e := r.executor(c).Exec(c, `UPDATE borrowings SET status='COMPLETED',completed_at=$2,fine_final_minor=$3 WHERE id=$1 AND status='CHECKED_OUT'`, v.ID, at, amount)
	if e != nil {
		return accountError(e)
	}
	if tag.RowsAffected() != 1 {
		return d.ErrState
	}
	return nil
}
func (r *BorrowingRepository) ClearFine(c context.Context, v d.Record, f d.Clearance) error {
	_, e := r.executor(c).Exec(c, `INSERT INTO fine_clearances(id,borrowing_id,actor_id,assessed_minor,cleared_minor,method,note,occurred_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, f.ID, v.ID, f.ActorID, f.Assessed, f.Amount, f.Method, f.Note, f.At)
	return accountError(e)
}

var _ d.AccountabilityRepository = (*BorrowingRepository)(nil)
