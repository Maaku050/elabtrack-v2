package inventory

import (
	"context"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/google/uuid"
	"strings"
	"time"
	"unicode"
)

const MaxQuantity int64 = 2147483647
const MaxImageBytes = 512 << 10

type Stock struct {
	Available   int64 `json:"available"`
	Reserved    int64 `json:"reserved"`
	CheckedOut  int64 `json:"checked_out"`
	DamagedHeld int64 `json:"damaged_held"`
	Total       int64 `json:"total_tracked"`
}

func (s Stock) Valid() bool {
	return s.Available >= 0 && s.Reserved >= 0 && s.CheckedOut >= 0 && s.DamagedHeld >= 0 && s.Total >= 0 && s.Total <= MaxQuantity && s.Total == s.Available+s.Reserved+s.CheckedOut+s.DamagedHeld && s.Available <= MaxQuantity && s.Reserved <= MaxQuantity && s.CheckedOut <= MaxQuantity && s.DamagedHeld <= MaxQuantity
}
func (s Stock) ChangeAvailable(delta int64) (Stock, error) {
	if !s.Valid() || delta < -MaxQuantity || delta > MaxQuantity {
		return s, shared.ErrInvalidInput
	}
	s.Available += delta
	s.Total += delta
	if !s.Valid() {
		return s, shared.ErrConflict
	}
	return s, nil
}

type Equipment struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	CategoryID    *uuid.UUID `json:"category_id"`
	CategoryName  string     `json:"category_name"`
	Status        string     `json:"status"`
	Stock         Stock      `json:"stock"`
	Sequence      int64      `json:"stock_sequence"`
	Version       int64      `json:"metadata_version"`
	ImageID       *uuid.UUID `json:"image_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ArchiveSafety string     `json:"archive_safety,omitempty"`
}
type Category struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Active  bool      `json:"is_active"`
	Version int64     `json:"version"`
}
type Metadata struct {
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	CategoryID      *uuid.UUID `json:"category_id"`
	ExpectedVersion int64      `json:"expected_version"`
}
type Create struct {
	Metadata
	Opening int64  `json:"opening_quantity"`
	Reason  string `json:"reason"`
}
type Adjustment struct {
	Kind             string `json:"kind"`
	Quantity         int64  `json:"quantity"`
	Reason           string `json:"reason"`
	ExpectedSequence *int64 `json:"expected_sequence"`
	Confirm          bool   `json:"confirm"`
}
type StatusInput struct {
	Status          string `json:"status"`
	ExpectedVersion int64  `json:"expected_version"`
	Confirm         bool   `json:"confirm"`
}
type CategoryInput struct {
	Name            string `json:"name"`
	Active          *bool  `json:"is_active"`
	ExpectedVersion int64  `json:"expected_version"`
}

func Text(s string, max int) bool {
	return len(s) <= max && !strings.ContainsFunc(s, func(r rune) bool { return unicode.IsControl(r) && r != '\n' && r != '\t' })
}
func (m Metadata) Valid() bool {
	return strings.TrimSpace(m.Name) != "" && Text(m.Name, 160) && Text(m.Description, 2000) && (m.CategoryID == nil || *m.CategoryID != uuid.Nil)
}

type Filter struct {
	Search, Status, Sort string
	CategoryID           *uuid.UUID
	AvailableOnly        bool
	Page, PerPage        int
	Borrower             bool
}

func (f Filter) Valid() bool {
	return f.Page >= 1 && f.Page <= 100000 && f.PerPage >= 1 && f.PerPage <= 100 && Text(f.Search, 100) && (f.Status == "" || f.Status == "ACTIVE" || f.Status == "INACTIVE" || f.Status == "ARCHIVED") && (f.Sort == "" || f.Sort == "name" || f.Sort == "available") && (f.CategoryID == nil || *f.CategoryID != uuid.Nil)
}

type Page struct {
	Items   []Equipment `json:"items"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
	Totals  Stock       `json:"totals"`
}
type Movement struct {
	ID          uuid.UUID `json:"id"`
	EquipmentID uuid.UUID `json:"equipment_id"`
	ActorID     uuid.UUID `json:"actor_id"`
	Sequence    int64     `json:"sequence"`
	Kind        string    `json:"kind"`
	Delta       Stock     `json:"delta"`
	After       Stock     `json:"after"`
	Reason      string    `json:"reason"`
	At          time.Time `json:"created_at"`
}
type Image struct {
	ID            uuid.UUID
	EquipmentID   uuid.UUID
	ActorID       uuid.UUID
	PNG           []byte
	Hash          string
	Width, Height int
}
type Repository interface {
	LockEquipment(context.Context, uuid.UUID) error
	Get(context.Context, uuid.UUID) (Equipment, error)
	List(context.Context, Filter) (Page, error)
	Insert(context.Context, Equipment) error
	UpdateMetadata(context.Context, Equipment) error
	UpdateStock(context.Context, Equipment) error
	Categories(context.Context, bool, int) ([]Category, error)
	GetCategory(context.Context, uuid.UUID, bool) (Category, error)
	PutCategory(context.Context, Category, bool) error
	Movement(context.Context, Movement) error
	Movements(context.Context, uuid.UUID, int) ([]Movement, error)
	Audit(context.Context, uuid.UUID, *uuid.UUID, *uuid.UUID, string, any, any) error
	ReadReceipt(context.Context, uuid.UUID, string, string, string) ([]byte, error)
	WriteReceipt(context.Context, uuid.UUID, string, string, string, []byte) error
	ArchiveSafety(context.Context) (string, error)
	PutImage(context.Context, Image) error
	GetImage(context.Context, uuid.UUID, uuid.UUID) (Image, error)
}
type ImageValidator interface{ Validate([]byte) (Image, error) }
