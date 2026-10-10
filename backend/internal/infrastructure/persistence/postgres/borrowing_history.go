package postgres

import (
	"context"
	"encoding/json"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/google/uuid"
)

func (r *BorrowingRepository) Preview(c context.Context, id uuid.UUID, lock bool) (d.Record, error) {
	return r.get(c, id, lock, 25)
}
func (r *BorrowingRepository) History(c context.Context, id uuid.UUID, kind string, page, size int) (p d.HistoryPage, e error) {
	sources := map[string]string{
		"events":       `SELECT id,actor_id,kind,reason,occurred_at FROM borrowing_events WHERE borrowing_id=$1`,
		"returns":      `SELECT id,item_id,good,damaged,lost,actor_id,reason,occurred_at FROM return_lines WHERE borrowing_id=$1`,
		"replacements": `SELECT id,obligation_id,equipment_id,equipment_name,quantity,actor_id,reason,equivalent_confirmed,occurred_at FROM replacement_acceptances WHERE borrowing_id=$1`,
		"clearances":   `SELECT id,actor_id,assessed_minor,cleared_minor,method,note,occurred_at FROM fine_clearances WHERE borrowing_id=$1`,
		"obligations":  `SELECT o.id,o.item_id,o.equipment_id,o.kind,o.required,(SELECT COALESCE(sum(quantity),0) FROM replacement_acceptances a WHERE a.obligation_id=o.id) AS accepted,l.occurred_at FROM replacement_obligations o JOIN return_lines l ON l.id=o.return_line_id WHERE o.borrowing_id=$1`,
	}
	source, ok := sources[kind]
	if !ok || page < 1 || page > 100000 || size < 1 || size > 100 {
		return p, shared.ErrInvalidInput
	}
	p = d.HistoryPage{Kind: kind, Page: page, PerPage: size, Items: []json.RawMessage{}}
	var raw []byte
	e = r.executor(c).QueryRow(c, `WITH source AS (`+source+`) SELECT (SELECT count(*) FROM source),COALESCE((SELECT jsonb_agg(to_jsonb(v) ORDER BY occurred_at DESC,id DESC) FROM(SELECT * FROM source ORDER BY occurred_at DESC,id DESC LIMIT $2 OFFSET $3)v),'[]'::jsonb)`, id, size, (page-1)*size).Scan(&p.Total, &raw)
	if e != nil {
		return p, accountError(e)
	}
	e = json.Unmarshal(raw, &p.Items)
	return
}
func (r *BorrowingRepository) previewAccountability(c context.Context, v *d.Record) error {
	for _, kind := range []string{"returns", "replacements", "clearances", "obligations"} {
		size := 25
		if kind == "obligations" {
			size = 100
		}
		p, e := r.History(c, v.ID, kind, 1, size)
		if e != nil {
			return e
		}
		raw, e := json.Marshal(map[string]any{kind: p.Items})
		if e != nil {
			return e
		}
		if e = json.Unmarshal(raw, v); e != nil {
			return e
		}
	}
	rows, e := r.executor(c).Query(c, `SELECT item_id,sum(good),sum(damaged),sum(lost) FROM return_lines WHERE borrowing_id=$1 GROUP BY item_id`, v.ID)
	if e != nil {
		return accountError(e)
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var good, damaged, lost int64
		if e = rows.Scan(&id, &good, &damaged, &lost); e != nil {
			return accountError(e)
		}
		for n := range v.Items {
			if v.Items[n].ID == id {
				v.Items[n].Good = good
				v.Items[n].Damaged = damaged
				v.Items[n].Lost = lost
			}
		}
	}
	if e = rows.Err(); e != nil {
		return accountError(e)
	}
	rows.Close()
	for n := range v.Items {
		v.Items[n].Outstanding = v.Items[n].Issued - v.Items[n].Good - v.Items[n].Damaged - v.Items[n].Lost
	}
	var cleared, obligations int64
	e = r.executor(c).QueryRow(c, `SELECT (SELECT COALESCE(sum(cleared_minor),0) FROM fine_clearances WHERE borrowing_id=$1),(SELECT COALESCE(sum(required-(SELECT COALESCE(sum(quantity),0) FROM replacement_acceptances a WHERE a.obligation_id=o.id)),0) FROM replacement_obligations o WHERE borrowing_id=$1)`, v.ID).Scan(&cleared, &obligations)
	v.FineClearedMinor = &cleared
	v.ReplacementOutstanding = &obligations
	v.PagedHistory = true
	return accountError(e)
}

var _ d.HistoryRepository = (*BorrowingRepository)(nil)
