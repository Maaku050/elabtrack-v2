package borrowing

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestFineStartedDays(t *testing.T) {
	due := time.Date(2026, 10, 10, 12, 0, 0, 123, time.UTC)
	for _, tt := range []struct {
		offset time.Duration
		want   int64
	}{{-1, 0}, {0, 0}, {1, 1000}, {time.Minute, 1000}, {24*time.Hour - 1, 1000}, {24 * time.Hour, 1000}, {24*time.Hour + 1, 2000}, {48 * time.Hour, 2000}, {48*time.Hour + 1, 3000}} {
		got, e := FineAmount(due, due.Add(tt.offset))
		if e != nil || got != tt.want {
			t.Fatalf("%s: %d %v", tt.offset, got, e)
		}
	}
	end := due.Add(25 * time.Hour)
	final := int64(2000)
	v := Record{DueAt: &due, CompletedAt: &end, FinalMinor: &final, Clearances: []Clearance{{Amount: 1000}}}
	if e := v.ProjectFine(end.Add(80 * time.Hour)); e != nil || v.Fine.Assessed != 2000 || v.Fine.Outstanding != 1000 || !v.Fine.Final {
		t.Fatal(v.Fine, e)
	}
}
func TestReturnNormalizationAndStock(t *testing.T) {
	id := uuid.New()
	l := ReturnLine{ItemID: id, Good: 2, Damaged: 1, Lost: 1}
	out, e := ReturnStock(inventory.Stock{CheckedOut: 5, Total: 5}, l)
	if e != nil || out.Available != 2 || out.CheckedOut != 1 || out.DamagedHeld != 1 || out.Total != 4 {
		t.Fatal(out, e)
	}
	for _, in := range [][]ReturnLine{nil, {l, l}, {{ItemID: id}}, {{ItemID: id, Good: -1}}, {{ItemID: id, Lost: inventory.MaxQuantity + 1}}} {
		if _, e := NormalizeReturns(in); e == nil {
			t.Fatal("invalid accepted", in)
		}
	}
	if _, e := ReturnStock(inventory.Stock{CheckedOut: 3, Total: 3}, l); e == nil {
		t.Fatal("over-return")
	}
}
