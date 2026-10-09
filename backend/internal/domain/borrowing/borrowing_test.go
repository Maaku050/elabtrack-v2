package borrowing

import (
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/google/uuid"
	"testing"
)

func TestCustodyVectors(t *testing.T) {
	base := inventory.Stock{Available: 5, DamagedHeld: 2, Total: 7}
	r, e := Transfer(base, 3, "RESERVE")
	if e != nil || r.Available != 2 || r.Reserved != 3 || r.Total != 7 || r.DamagedHeld != 2 {
		t.Fatal(r, e)
	}
	for _, kind := range []string{"CANCELLED", "DENIED", "EXPIRED"} {
		back, e := Transfer(r, 3, kind)
		if e != nil || back != base {
			t.Fatal(kind, back, e)
		}
	}
	c, e := Transfer(r, 3, "CHECKOUT")
	if e != nil || c.CheckedOut != 3 || c.Reserved != 0 || c.Available != 2 || c.Total != 7 {
		t.Fatal(c, e)
	}
	direct, e := Transfer(base, 3, "DIRECT")
	if e != nil || direct != c {
		t.Fatal(direct, e)
	}
	for _, q := range []int64{-1, 0, 6, inventory.MaxQuantity + 1} {
		if _, e = Transfer(base, q, "RESERVE"); e == nil {
			t.Fatal("invalid quantity accepted", q)
		}
	}
	if _, e = Transfer(base, 1, "RETURN"); e == nil {
		t.Fatal("deferred effect")
	}
}
func TestRequestLines(t *testing.T) {
	id := uuid.New()
	for _, lines := range [][]Line{nil, {{id, 0}}, {{id, -1}}, {{uuid.Nil, 1}}, {{id, 1}, {id, 2}}, make([]Line, 101)} {
		if _, e := Normalize(lines); e == nil {
			t.Fatal("invalid request accepted")
		}
	}
	a, b := uuid.New(), uuid.New()
	x, e := Normalize([]Line{{a, 2}, {b, 1}})
	if e != nil || len(x) != 2 || x[0].EquipmentID.String() > x[1].EquipmentID.String() {
		t.Fatal(x, e)
	}
	if ValidReason(" ") || !ValidReason("Explain why") || State("APPROVED") {
		t.Fatal("state/reason")
	}
}
