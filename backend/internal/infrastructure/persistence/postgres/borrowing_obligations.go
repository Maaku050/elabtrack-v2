package postgres

import (
	"context"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"github.com/google/uuid"
)

// Complete installed accountability facts. This is read-only; deactivation never resolves them.
type BorrowingObligations struct{ repo *BorrowingRepository }

func NewBorrowingObligations(r *BorrowingRepository) *BorrowingObligations {
	return &BorrowingObligations{r}
}
func (r *BorrowingObligations) Read(c context.Context, id uuid.UUID) (v d.Obligations, e error) {
	v.Availability = "UNAVAILABLE"
	var installed bool
	e = r.repo.executor(c).QueryRow(c, `SELECT to_regclass('public.return_lines') IS NOT NULL AND to_regclass('public.replacement_acceptances') IS NOT NULL AND to_regclass('public.fine_clearances') IS NOT NULL`).Scan(&installed)
	if e != nil {
		return v, accountError(e)
	}
	if !installed {
		return v, nil
	}
	var pending, active, overdue, units, replacements int
	var fine int64
	e = r.repo.executor(c).QueryRow(c, `WITH tick AS MATERIALIZED (SELECT clock_timestamp() at), loans AS (SELECT * FROM borrowings WHERE borrower_id=$1) SELECT
 count(*) FILTER(WHERE status='PENDING'),count(*) FILTER(WHERE status='CHECKED_OUT'),count(*) FILTER(WHERE status='CHECKED_OUT' AND due_at<(SELECT at FROM tick)),
 COALESCE((SELECT sum(i.issued_quantity-COALESCE((SELECT sum(good+damaged+lost) FROM return_lines l WHERE l.item_id=i.id),0)) FROM borrowing_items i JOIN loans b ON b.id=i.borrowing_id WHERE b.status='CHECKED_OUT'),0)::bigint,
 COALESCE((SELECT sum(o.required-COALESCE((SELECT sum(quantity) FROM replacement_acceptances a WHERE a.obligation_id=o.id),0)) FROM replacement_obligations o JOIN loans b ON b.id=o.borrowing_id),0)::bigint,
 COALESCE(sum(COALESCE(fine_final_minor,borrowing_fine_amount(due_at,(SELECT at FROM tick)))-COALESCE((SELECT sum(cleared_minor) FROM fine_clearances f WHERE f.borrowing_id=loans.id),0)),0)::bigint FROM loans`, id).Scan(&pending, &active, &overdue, &units, &replacements, &fine)
	if e != nil {
		return v, accountError(e)
	}
	v.Availability = "AVAILABLE"
	v.PendingRequests = &pending
	v.ActiveBorrowings = &active
	v.OverdueBorrowings = &overdue
	v.UnreturnedUnits = &units
	v.ReplacementUnits = &replacements
	v.FineMinor = &fine
	return v, nil
}
