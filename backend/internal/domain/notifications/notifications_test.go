package notifications

import (
	"github.com/google/uuid"
	"testing"
)

func TestMessagesAndBounds(t *testing.T) {
	for _, kind := range []string{"SUBMITTED", "CHECKED_OUT", "DIRECT", "DENIED", "CANCELLED", "EXPIRED", "RETURN", "DAMAGED", "LOST", "REPLACEMENT_REQUIRED", "REPLACEMENT", "COMPLETED", "FINE_CLEARED", "DUE_SOON", "OVERDUE", "FINE_ASSESSED"} {
		m := (Event{Kind: kind, BorrowingID: uuid.New()}).Message()
		if m.Title == "" || m.Body == "" {
			t.Fatal("missing event", kind)
		}
	}
	if (Event{Kind: "UNKNOWN"}).Message().Title != "" {
		t.Fatal("unknown event accepted")
	}
	for _, f := range []Filter{{0, 25, ""}, {1, 101, ""}, {1, 25, "invalid"}} {
		if f.Valid() {
			t.Fatal("unbounded filter")
		}
	}
	if !(Filter{1, 25, "unread"}).Valid() {
		t.Fatal("valid unread")
	}
}
