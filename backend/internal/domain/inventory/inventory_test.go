package inventory

import (
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"testing"
)

func TestConservationBoundaries(t *testing.T) {
	for _, s := range []Stock{{Available: 2, Reserved: 3, CheckedOut: 4, DamagedHeld: 5, Total: 14}, {Available: MaxQuantity, Total: MaxQuantity}, {}} {
		if !s.Valid() {
			t.Fatal("valid physical vector")
		}
		v, e := s.ChangeAvailable(1)
		if s.Total == MaxQuantity {
			if !errors.Is(e, shared.ErrConflict) {
				t.Fatal("overflow")
			}
			continue
		}
		if e != nil || v.Reserved != s.Reserved || v.CheckedOut != s.CheckedOut || v.DamagedHeld != s.DamagedHeld || v.Total != s.Total+1 {
			t.Fatal("custody changed")
		}
		if _, e = s.ChangeAvailable(-s.Available - 1); e == nil {
			t.Fatal("negative accepted")
		}
	}
	for _, s := range []Stock{{Available: -1, Total: -1}, {Available: 1, Total: 2}, {Available: MaxQuantity + 1, Total: MaxQuantity + 1}, {CheckedOut: 1}} {
		if s.Valid() {
			t.Fatal("invalid vector accepted")
		}
	}
}
