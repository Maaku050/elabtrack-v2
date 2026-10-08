// Package terms owns immutable version and authenticated acceptance concepts.
package terms

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/google/uuid"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrNotPublished       = errors.New("terms not published")
	ErrVersionNotFound    = errors.New("terms version not found")
	ErrVersionChanged     = errors.New("terms version changed")
	ErrAcceptanceRequired = errors.New("terms acceptance required")
	ErrVersionExists      = errors.New("terms version exists")
	ErrPublicationChanged = errors.New("terms publication changed")
)

const MaxBodyBytes = 65536

var versionPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$`)

type Version struct {
	ID          uuid.UUID `json:"id"`
	Version     string    `json:"version"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	ContentHash string    `json:"content_hash"`
	PublishedBy uuid.UUID `json:"-"`
	CreatedAt   time.Time `json:"-"`
	PublishedAt time.Time `json:"published_at"`
}

// Content is plain UTF-8 text. It is never interpreted as HTML or Markdown.
func NewVersion(version, title, body string, publisher uuid.UUID) (*Version, error) {
	version, title = strings.TrimSpace(version), strings.TrimSpace(title)
	if !versionPattern.MatchString(version) || title == "" || len(title) > 200 || strings.ContainsAny(title, "\r\n\x00") ||
		strings.TrimSpace(body) == "" || len(body) > MaxBodyBytes || strings.ContainsRune(body, 0) || !utf8.ValidString(body) || !utf8.ValidString(title) || publisher == uuid.Nil {
		return nil, shared.ErrValidation
	}
	sum := sha256.Sum256([]byte(body))
	return &Version{ID: uuid.New(), Version: version, Title: title, Body: body, ContentHash: hex.EncodeToString(sum[:]), PublishedBy: publisher}, nil
}

type Acceptance struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"-"`
	TermsVersionID uuid.UUID `json:"terms_version_id"`
	AcceptedAt     time.Time `json:"accepted_at"`
}

type Status struct {
	State                 string      `json:"state"`
	CurrentTerms          *Version    `json:"current_terms"`
	Acceptance            *Acceptance `json:"acceptance"`
	HasPreviousAcceptance bool        `json:"has_previous_acceptance"`
	AcceptanceRequired    bool        `json:"acceptance_required"`
	CanInitiateBorrowing  bool        `json:"can_initiate_borrowing"`
}

func StatusFor(version *Version, acceptance *Acceptance, previous bool) Status {
	s := Status{State: "unpublished", CurrentTerms: version, Acceptance: acceptance, HasPreviousAcceptance: previous}
	if version == nil {
		s.Acceptance = nil
		return s
	}
	s.State = "required"
	s.AcceptanceRequired = true
	if previous {
		s.State = "updated"
	}
	if acceptance != nil && acceptance.TermsVersionID == version.ID {
		s.State = "accepted"
		s.AcceptanceRequired = false
		s.CanInitiateBorrowing = true
	}
	return s
}

// LockPublication requires the caller's transaction. Shared readers hold the
// pointer through acceptance/command commit; publication takes it exclusively.
// Call after sorted account locks, before future borrowing/equipment locks.
type Repository interface {
	LockPublication(context.Context, bool) (*Version, error)
	FindVersion(context.Context, uuid.UUID) (*Version, error)
	FindAcceptance(context.Context, uuid.UUID, uuid.UUID) (*Acceptance, error)
	HasAnyAcceptance(context.Context, uuid.UUID) (bool, error)
	InsertAcceptance(context.Context, uuid.UUID, uuid.UUID) (*Acceptance, error)
	Publish(context.Context, *Version) error
}
