package borrowing

import (
	"context"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/google/uuid"
	"math"
	"sort"
	"time"
)

type ReturnLine struct {
	ItemID  uuid.UUID `json:"item_id"`
	Good    int64     `json:"good"`
	Damaged int64     `json:"damaged"`
	Lost    int64     `json:"lost"`
}
type ReturnInput struct {
	ExpectedEvents int          `json:"expected_event_count"`
	Lines          []ReturnLine `json:"lines"`
	Reason         string       `json:"reason"`
	Confirm        bool         `json:"confirm"`
}
type ReplacementInput struct {
	ExpectedEvents int       `json:"expected_event_count"`
	ObligationID   uuid.UUID `json:"obligation_id"`
	Quantity       int64     `json:"quantity"`
	Reason         string    `json:"reason"`
	Equivalent     bool      `json:"equivalent_confirmed"`
	Confirm        bool      `json:"confirm"`
}
type FineInput struct {
	Method        string `json:"method"`
	Note          string `json:"note"`
	ExpectedMinor int64  `json:"expected_outstanding_minor"`
	Confirm       bool   `json:"confirm"`
}
type Obligation struct {
	ID          uuid.UUID `json:"id"`
	ItemID      uuid.UUID `json:"item_id"`
	EquipmentID uuid.UUID `json:"equipment_id"`
	Kind        string    `json:"kind"`
	Required    int64     `json:"required"`
	Accepted    int64     `json:"accepted"`
}
type Disposition struct {
	ID uuid.UUID `json:"id"`
	ReturnLine
	ActorID uuid.UUID `json:"actor_id"`
	Reason  string    `json:"reason"`
	At      time.Time `json:"occurred_at"`
}
type Replacement struct {
	EquipmentID   uuid.UUID `json:"equipment_id"`
	EquipmentName string    `json:"equipment_name"`
	Equivalent    bool      `json:"equivalent_confirmed"`
	ID            uuid.UUID `json:"id"`
	ObligationID  uuid.UUID `json:"obligation_id"`
	Quantity      int64     `json:"quantity"`
	ActorID       uuid.UUID `json:"actor_id"`
	Reason        string    `json:"reason"`
	At            time.Time `json:"occurred_at"`
}
type Clearance struct {
	ID       uuid.UUID `json:"id"`
	ActorID  uuid.UUID `json:"actor_id"`
	Assessed int64     `json:"assessed_minor"`
	Amount   int64     `json:"cleared_minor"`
	Method   string    `json:"method"`
	Note     string    `json:"note"`
	At       time.Time `json:"occurred_at"`
}
type Fine struct {
	Assessed    int64     `json:"assessed_minor"`
	Cleared     int64     `json:"cleared_minor"`
	Outstanding int64     `json:"outstanding_minor"`
	Final       bool      `json:"is_final"`
	AsOf        time.Time `json:"as_of"`
}

func NormalizeReturns(in []ReturnLine) ([]ReturnLine, error) {
	if len(in) == 0 || len(in) > MaxItems {
		return nil, shared.ErrInvalidInput
	}
	out := append([]ReturnLine(nil), in...)
	sort.Slice(out, func(i, j int) bool { return out[i].ItemID.String() < out[j].ItemID.String() })
	for n, l := range out {
		if l.ItemID == uuid.Nil || l.Good < 0 || l.Damaged < 0 || l.Lost < 0 || l.Good > inventory.MaxQuantity || l.Damaged > inventory.MaxQuantity || l.Lost > inventory.MaxQuantity || l.Good+l.Damaged+l.Lost <= 0 || l.Good+l.Damaged+l.Lost > inventory.MaxQuantity || n > 0 && l.ItemID == out[n-1].ItemID {
			return nil, shared.ErrInvalidInput
		}
	}
	return out, nil
}
func ReturnStock(s inventory.Stock, l ReturnLine) (inventory.Stock, error) {
	if !s.Valid() || l.Good < 0 || l.Damaged < 0 || l.Lost < 0 || l.Good > inventory.MaxQuantity || l.Damaged > inventory.MaxQuantity || l.Lost > inventory.MaxQuantity || l.Good+l.Damaged+l.Lost < 1 || l.Good+l.Damaged+l.Lost > inventory.MaxQuantity {
		return s, shared.ErrInvalidInput
	}
	s.Available += l.Good
	s.CheckedOut -= l.Good + l.Damaged + l.Lost
	s.DamagedHeld += l.Damaged
	s.Total -= l.Lost
	if !s.Valid() {
		return s, ErrStock
	}
	return s, nil
}

// FineAmount uses elapsed started 24-hour periods, bounded integer centavos.
// Unix seconds plus nanosecond comparison avoids time.Duration's 290-year limit.
func FineAmount(due, end time.Time) (int64, error) {
	if !end.After(due) {
		return 0, nil
	}
	seconds := end.Unix() - due.Unix()
	if seconds < 0 {
		return 0, shared.ErrInvalidInput
	}
	days := seconds / 86400
	if seconds%86400 > 0 || end.Nanosecond() > due.Nanosecond() {
		days++
	}
	if days > math.MaxInt64/1000 {
		return 0, shared.ErrInvalidInput
	}
	return days * 1000, nil
}
func (v *Record) ProjectFine(now time.Time) error {
	v.Fine = Fine{AsOf: now, Final: v.CompletedAt != nil}
	if v.CompletedAt != nil {
		if v.FinalMinor == nil || *v.FinalMinor < 0 || v.DueAt == nil {
			return ErrState
		}
		v.Fine.Assessed = *v.FinalMinor
	} else if v.DueAt != nil {
		end := now
		if v.CompletedAt != nil {
			end = *v.CompletedAt
		}
		n, e := FineAmount(*v.DueAt, end)
		if e != nil {
			return e
		}
		v.Fine.Assessed = n
	}
	for _, c := range v.Clearances {
		if c.Amount > math.MaxInt64-v.Fine.Cleared {
			return shared.ErrInvalidInput
		}
		v.Fine.Cleared += c.Amount
	}
	if v.FineClearedMinor != nil {
		if *v.FineClearedMinor < 0 {
			return ErrState
		}
		v.Fine.Cleared = *v.FineClearedMinor
	}
	v.Fine.Outstanding = v.Fine.Assessed - v.Fine.Cleared
	if v.Fine.Outstanding < 0 {
		return ErrState
	}
	return nil
}

type AccountabilityRepository interface {
	RecordReturn(context.Context, Record, Item, Disposition, inventory.Equipment, inventory.Stock) error
	RecordReplacement(context.Context, Record, Obligation, Replacement, inventory.Equipment, inventory.Stock) error
	Complete(context.Context, Record, time.Time, int64) error
	ClearFine(context.Context, Record, Clearance) error
}
