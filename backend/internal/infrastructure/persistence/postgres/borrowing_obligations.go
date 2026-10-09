package postgres

import (
	"context"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"github.com/google/uuid"
)

// Phase7 exposes only installed borrowing facts. Fine/replacement balances are
// unavailable, never fabricated zeroes; later disposition modules replace this adapter.
type BorrowingObligations struct{ repo *BorrowingRepository }

func NewBorrowingObligations(r *BorrowingRepository) *BorrowingObligations {
	return &BorrowingObligations{r}
}
func (r *BorrowingObligations) Read(c context.Context, id uuid.UUID) (v d.Obligations, e error) {
	v.Availability = "UNAVAILABLE"
	var installed, future bool
	e = r.repo.executor(c).QueryRow(c, `SELECT to_regclass('public.borrowing_items') IS NOT NULL,to_regclass('public.return_events') IS NOT NULL OR to_regclass('public.replacement_obligations') IS NOT NULL`).Scan(&installed, &future)
	if e != nil {
		return v, accountError(e)
	}
	if !installed || future {
		return v, nil
	}
	var pending, active, overdue, units int
	e = r.repo.executor(c).QueryRow(c, `SELECT count(*) FILTER(WHERE b.status='PENDING'),count(*) FILTER(WHERE b.status='CHECKED_OUT'),count(*) FILTER(WHERE b.status='CHECKED_OUT' AND b.due_at<clock_timestamp()),COALESCE(sum((SELECT COALESCE(sum(issued_quantity),0) FROM borrowing_items i WHERE i.borrowing_id=b.id)) FILTER(WHERE b.status='CHECKED_OUT'),0)::bigint FROM borrowings b WHERE b.borrower_id=$1`, id).Scan(&pending, &active, &overdue, &units)
	if e != nil {
		return v, accountError(e)
	}
	v.Availability = "PARTIAL"
	v.PendingRequests = &pending
	v.ActiveBorrowings = &active
	v.OverdueBorrowings = &overdue
	v.UnreturnedUnits = &units
	return v, nil
}
