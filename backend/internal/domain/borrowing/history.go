package borrowing

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
)

// Public history is paged independently from the complete mutation aggregate.
type HistoryPage struct {
	Items   []json.RawMessage `json:"items"`
	Kind    string            `json:"kind"`
	Page    int               `json:"page"`
	PerPage int               `json:"per_page"`
	Total   int               `json:"total"`
}
type HistoryRepository interface {
	Preview(context.Context, uuid.UUID, bool) (Record, error)
	History(context.Context, uuid.UUID, string, int, int) (HistoryPage, error)
}

func HistoryKind(v string) bool {
	return v == "events" || v == "returns" || v == "replacements" || v == "clearances" || v == "obligations"
}
