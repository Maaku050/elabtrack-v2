package notifications

import (
	"context"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Key                     string
	BorrowingID, BorrowerID uuid.UUID
	ActorID                 *uuid.UUID
	Kind                    string
	At                      time.Time
}
type Notification struct {
	ID          uuid.UUID  `json:"id"`
	BorrowingID uuid.UUID  `json:"borrowing_id"`
	Kind        string     `json:"kind"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Scope       string     `json:"scope"`
	CreatedAt   time.Time  `json:"created_at"`
	ReadAt      *time.Time `json:"read_at"`
}
type Filter struct {
	Page, PerPage int
	Read          string
}

func (f Filter) Valid() bool {
	return f.Page >= 1 && f.Page <= 100000 && f.PerPage >= 1 && f.PerPage <= 100 && (f.Read == "" || f.Read == "read" || f.Read == "unread")
}

type Page struct {
	Items   []Notification `json:"items"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
	Total   int            `json:"total"`
	Unread  int            `json:"unread"`
}
type Message struct {
	Title, Body string
	Operations  bool
}

func (e Event) Message() Message {
	switch e.Kind {
	case "SUBMITTED":
		return Message{"Request submitted", "The request is reserved and awaits Staff/Admin review.", true}
	case "CHECKED_OUT", "DIRECT":
		return Message{"Equipment handed over", "Physical checkout was recorded. Review the original due date.", false}
	case "DENIED":
		return Message{"Request denied", "Review the retained denial explanation. Reserved stock was released.", false}
	case "CANCELLED":
		return Message{"Request cancelled", "Reserved stock was released; request history is retained.", false}
	case "EXPIRED":
		return Message{"Reservation expired", "The24-hour reservation expired and reserved stock was released.", false}
	case "RETURN":
		return Message{"Return recorded", "Staff/Admin recorded a verified in-person disposition. Review physical and replacement accountability.", false}
	case "DAMAGED":
		return Message{"Damaged equipment received", "A damaged original was received and remains held. Review the replacement obligation.", true}
	case "LOST":
		return Message{"Equipment loss recorded", "Lost equipment was confirmed. Review the separate replacement obligation.", true}
	case "REPLACEMENT_REQUIRED":
		return Message{"Replacement required", "An exact damaged/lost replacement obligation was recorded. FSMO must verify its physical resolution.", true}
	case "REPLACEMENT":
		return Message{"Physical replacement accepted", "Equivalent usable replacement equipment was received. Review remaining accountability.", false}
	case "COMPLETED":
		return Message{"Borrowing completed", "Physical and replacement obligations are resolved. The final overdue assessment is frozen; money may remain outstanding.", false}
	case "FINE_CLEARED":
		return Message{"Fine resolved", "Admin recorded full current-balance clearance. Review the method and retained assessment.", false}
	case "DUE_SOON":
		return Message{"Borrowing due soon", "Check the due date and arrange face-to-face return processing with FSMO.", false}
	case "OVERDUE":
		return Message{"Borrowing overdue", "Physical or replacement obligations remain after the original due date. Contact FSMO for resolution.", true}
	case "FINE_ASSESSED":
		return Message{"Overdue fine assessed", "Review the authoritative current overdue assessment and outstanding balance.", false}
	default:
		return Message{}
	}
}

type Repository interface {
	TryWorkerLock(context.Context) (bool, error)
	Pending(context.Context, int, int64) ([]Event, error)
	Emit(context.Context, Event, Message) (bool, error)
	List(context.Context, uuid.UUID, user.Role, Filter) (Page, error)
	Count(context.Context, uuid.UUID, user.Role) (int, error)
	MarkAll(context.Context, uuid.UUID, user.Role) (int64, error)
	Mark(context.Context, uuid.UUID, user.Role, uuid.UUID, bool) (Notification, error)
}
