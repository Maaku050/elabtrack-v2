package borrowing

import (
	"context"
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/google/uuid"
	"sort"
	"strings"
	"time"
)

const PendingTTL = 24 * time.Hour
const MaxItems = 100 // transport bound, not a per-borrower quota
var (
	ErrStock       = errors.New("equipment unavailable")
	ErrState       = errors.New("borrowing state conflict")
	ErrExpired     = errors.New("borrowing expired")
	ErrKey         = errors.New("idempotency conflict")
	ErrEligibility = errors.New("borrower ineligible")
)

type Line struct {
	EquipmentID uuid.UUID `json:"equipment_id"`
	Quantity    int64     `json:"quantity"`
}
type Input struct {
	BorrowerID uuid.UUID  `json:"borrower_id,omitempty"`
	Items      []Line     `json:"items"`
	DueAt      *time.Time `json:"due_at,omitempty"`
	Confirm    bool       `json:"confirm"`
	Handover   bool       `json:"physical_handover_confirmed"`
}
type Decision struct {
	DueAt    *time.Time `json:"due_at,omitempty"`
	Reason   string     `json:"reason,omitempty"`
	Confirm  bool       `json:"confirm"`
	Handover bool       `json:"physical_handover_confirmed"`
}
type Item struct {
	ID          uuid.UUID `json:"id"`
	EquipmentID uuid.UUID `json:"equipment_id"`
	Name        string    `json:"name"`
	Quantity    int64     `json:"quantity"`
	Reserved    int64     `json:"reserved_quantity"`
	Issued      int64     `json:"issued_quantity"`
}
type Event struct {
	ID      uuid.UUID  `json:"id"`
	ActorID *uuid.UUID `json:"actor_id"`
	Kind    string     `json:"kind"`
	Reason  string     `json:"reason"`
	At      time.Time  `json:"occurred_at"`
}
type Record struct {
	ID           uuid.UUID  `json:"id"`
	Reference    string     `json:"reference"`
	BorrowerID   uuid.UUID  `json:"borrower_id"`
	BorrowerName string     `json:"borrower_name"`
	BorrowerType string     `json:"borrower_type"`
	StudentID    string     `json:"student_id"`
	Status       string     `json:"status"`
	EntryPath    string     `json:"entry_path"`
	AcceptanceID uuid.UUID  `json:"acceptance_id"`
	CreatedAt    time.Time  `json:"created_at"`
	ExpiresAt    *time.Time `json:"expires_at"`
	CheckedOutAt *time.Time `json:"checked_out_at"`
	DueAt        *time.Time `json:"due_at"`
	TerminalAt   *time.Time `json:"terminal_at"`
	Reason       string     `json:"denial_reason"`
	Items        []Item     `json:"items"`
	Events       []Event    `json:"events"`
	Overdue      bool       `json:"is_overdue"`
}
type Filter struct {
	Page, PerPage  int
	Status, Search string
	Owner          *uuid.UUID
}
type Page struct {
	Items   []Record `json:"items"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
	Total   int      `json:"total"`
}

func State(s string) bool {
	return s == "PENDING" || s == "CHECKED_OUT" || s == "DENIED" || s == "CANCELLED" || s == "EXPIRED"
}
func (f Filter) Valid() bool {
	return f.Page >= 1 && f.Page <= 100000 && f.PerPage >= 1 && f.PerPage <= 100 && (f.Status == "" || State(f.Status)) && inventory.Text(f.Search, 100)
}
func Normalize(lines []Line) ([]Line, error) {
	if len(lines) == 0 || len(lines) > MaxItems {
		return nil, shared.ErrInvalidInput
	}
	out := append([]Line(nil), lines...)
	sort.Slice(out, func(i, j int) bool { return out[i].EquipmentID.String() < out[j].EquipmentID.String() })
	for i, v := range out {
		if v.EquipmentID == uuid.Nil || v.Quantity <= 0 || v.Quantity > inventory.MaxQuantity || i > 0 && v.EquipmentID == out[i-1].EquipmentID {
			return nil, shared.ErrInvalidInput
		}
	}
	return out, nil
}
func ValidReason(reason string) bool {
	return strings.TrimSpace(reason) != "" && inventory.Text(reason, 1000)
}

// Transfer changes custody only; total and damaged stock are preserved.
func Transfer(s inventory.Stock, q int64, kind string) (inventory.Stock, error) {
	if !s.Valid() || q <= 0 || q > inventory.MaxQuantity {
		return s, shared.ErrInvalidInput
	}
	switch kind {
	case "RESERVE":
		s.Available -= q
		s.Reserved += q
	case "CANCELLED", "DENIED", "EXPIRED":
		s.Reserved -= q
		s.Available += q
	case "CHECKOUT":
		s.Reserved -= q
		s.CheckedOut += q
	case "DIRECT":
		s.Available -= q
		s.CheckedOut += q
	default:
		return s, shared.ErrInvalidInput
	}
	if !s.Valid() {
		return s, ErrStock
	}
	return s, nil
}

type Repository interface {
	Clock(context.Context) (time.Time, error)
	Get(context.Context, uuid.UUID, bool) (Record, error)
	List(context.Context, Filter) (Page, error)
	Insert(context.Context, Record) error
	Transition(context.Context, Record) error
	Effect(context.Context, Record, Item, *uuid.UUID, string, inventory.Equipment, inventory.Stock) error
	Event(context.Context, uuid.UUID, Event) error
	ReadReceipt(context.Context, uuid.UUID, string, string, string) ([]byte, error)
	WriteReceipt(context.Context, uuid.UUID, string, string, string, []byte) error
	Due(context.Context, int) ([]uuid.UUID, error)
}
