package shared

import (
	"time"

	"github.com/google/uuid"
)

// ID is a strongly-typed identifier wrapper around uuid.UUID.
type ID = uuid.UUID

// NewID generates a new random identifier.
func NewID() ID { return uuid.New() }

// ParseID parses a string into an identifier.
func ParseID(s string) (ID, error) { return uuid.Parse(s) }

// Timestamps bundles common audit columns.
type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
