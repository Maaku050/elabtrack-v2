package accounts

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

var (
	ErrDomainsMissing    = errors.New("student domains not configured")
	ErrStudentDomain     = errors.New("student email domain is not allowed")
	ErrStudentIDExists   = errors.New("student ID already exists")
	ErrActivationInvalid = errors.New("activation is invalid or expired")
	ErrCooldown          = errors.New("activation resend cooldown")
)

type Input struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	BorrowerType  string `json:"borrower_type"`
	StudentID     string `json:"student_id"`
	Course        string `json:"course"`
	ContactNumber string `json:"contact_number"`
}

type Policy struct {
	StudentDomains []string
	ActivationTTL  time.Duration
	ActivationURL  string
	ResendCooldown time.Duration
}

func (p Policy) Validate(in Input, staff bool) (Input, error) {
	email, err := user.ParseEmail(in.Email)
	if err != nil || len(email.String()) > 254 || !strings.EqualFold(strings.TrimSpace(in.Email), email.String()) {
		return in, user.ErrInvalidEmail
	}
	name, err := user.ParseName(in.Name)
	if err != nil {
		return in, err
	}
	in.Email, in.Name = email.String(), name.String()
	in.StudentID, in.Course, in.ContactNumber = strings.TrimSpace(in.StudentID), strings.TrimSpace(in.Course), strings.TrimSpace(in.ContactNumber)
	if !bounded(in.Course, 128) || !bounded(in.ContactNumber, 32) {
		return in, shared.ErrInvalidInput
	}
	if staff {
		if in.BorrowerType != "" || in.StudentID != "" || in.Course != "" || in.ContactNumber != "" {
			return in, shared.ErrInvalidInput
		}
		return in, nil
	}
	switch in.BorrowerType {
	case "STUDENT":
		if !bounded(in.StudentID, 64) || in.StudentID == "" {
			return in, shared.ErrInvalidInput
		}
		if len(p.StudentDomains) == 0 {
			return in, ErrDomainsMissing
		}
		domain := in.Email[strings.LastIndex(in.Email, "@")+1:]
		allowed := false
		for _, d := range p.StudentDomains {
			if domain == d {
				allowed = true
			}
		}
		if !allowed {
			return in, ErrStudentDomain
		}
	case "FACULTY":
		if in.StudentID != "" {
			return in, shared.ErrInvalidInput
		}
	default:
		return in, shared.ErrInvalidInput
	}
	return in, nil
}
func bounded(s string, max int) bool {
	if len(s) > max {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

type Obligations struct {
	Availability      string `json:"availability"`
	PendingRequests   *int   `json:"pending_requests,omitempty"`
	ActiveBorrowings  *int   `json:"active_borrowings,omitempty"`
	OverdueBorrowings *int   `json:"overdue_borrowings,omitempty"`
	UnreturnedUnits   *int   `json:"unreturned_units,omitempty"`
	FineMinor         *int64 `json:"fine_minor,omitempty"`
	ReplacementUnits  *int   `json:"replacement_units,omitempty"`
}
type ObligationReader interface {
	Read(context.Context, uuid.UUID) (Obligations, error)
}
type UnavailableObligations struct{}

func (UnavailableObligations) Read(context.Context, uuid.UUID) (Obligations, error) {
	return Obligations{Availability: "UNAVAILABLE"}, nil
}

type Record struct {
	ID                 uuid.UUID   `json:"id"`
	Name               string      `json:"name"`
	Email              string      `json:"email"`
	Role               user.Role   `json:"role"`
	IsActive           bool        `json:"is_active"`
	BorrowerType       string      `json:"borrower_type"`
	StudentID          string      `json:"student_id"`
	Course             string      `json:"course"`
	ContactNumber      string      `json:"contact_number"`
	ActivationRequired bool        `json:"activation_required"`
	DeliveryStatus     string      `json:"delivery_status"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	Obligations        Obligations `json:"obligations"`
}
type Filter struct {
	Page, PerPage                int
	Search, BorrowerType, Status string
	Staff                        bool
}
type Page struct {
	Items   []Record `json:"items"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	PerPage int      `json:"per_page"`
}
type Token struct {
	Hash                 string
	AccountID            uuid.UUID
	CreatedAt, ExpiresAt time.Time
	Invalid              bool
}
type Mail struct{ To, Name, ActivationURL string }
type Sender interface {
	SendActivation(context.Context, Mail) (string, error)
}

// Delivery states describe provider submission, never verified recipient delivery.
var ErrMailUnavailable = errors.New("mail not configured")
var ErrMailUnknown = errors.New("mail submission outcome unknown")

type Row struct {
	Number        int         `json:"number"`
	Input         Input       `json:"input"`
	Error         string      `json:"error,omitempty"`
	TargetID      *uuid.UUID  `json:"target_id,omitempty"`
	TargetVersion *time.Time  `json:"target_version,omitempty"`
	Outcome       string      `json:"outcome,omitempty"`
	AccountID     *uuid.UUID  `json:"account_id,omitempty"`
	Obligations   Obligations `json:"obligations"`
}
type Batch struct {
	ID        uuid.UUID `json:"id"`
	ActorID   uuid.UUID `json:"-"`
	Operation string    `json:"operation"`
	Rows      []Row     `json:"rows"`
	ExpiresAt time.Time `json:"expires_at"`
	Selection []int     `json:"selection,omitempty"`
	Confirmed bool      `json:"confirmed"`
}
type Audit struct {
	ID        uuid.UUID `json:"id"`
	ActorID   uuid.UUID `json:"actor_id"`
	AccountID uuid.UUID `json:"account_id"`
	Action    string    `json:"action"`
	At        time.Time `json:"occurred_at"`
}

type Repository interface {
	LockIdentity(context.Context) error
	LockAccounts(context.Context, []uuid.UUID) error
	Get(context.Context, uuid.UUID) (Record, error)
	FindIdentity(context.Context, string, string) ([]Record, error)
	List(context.Context, Filter) (Page, error)
	Insert(context.Context, Record) error
	UpdateProfile(context.Context, Record) error
	SetActive(context.Context, uuid.UUID, bool) error
	SetPassword(context.Context, uuid.UUID, string) error
	PutToken(context.Context, Token) error
	FindToken(context.Context, string, bool) (Token, error)
	InvalidateTokens(context.Context, uuid.UUID) error
	SetDelivery(context.Context, uuid.UUID, string, string, string) error
	ReadReceipt(context.Context, uuid.UUID, string, string, string) ([]byte, error)
	WriteReceipt(context.Context, uuid.UUID, string, string, string, []byte) error
	AddAudit(context.Context, uuid.UUID, uuid.UUID, string) error
	Audits(context.Context, uuid.UUID, int) ([]Audit, error)
	PutBatch(context.Context, Batch) error
	GetBatch(context.Context, uuid.UUID, bool) (Batch, error)
	ConfirmBatch(context.Context, Batch) error
}

const MaxRosterUpload = 768 << 10

type Parser interface {
	Parse([]byte) ([]Row, error)
	Template() ([]byte, error)
}
