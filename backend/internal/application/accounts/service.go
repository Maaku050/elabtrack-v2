package accounts

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

type Service struct {
	repo        d.Repository
	tx          application.Transactions
	hasher      application.PasswordHasher
	sender      d.Sender
	policy      d.Policy
	obligations d.ObligationReader
	Now         func() time.Time
}

func NewService(repo d.Repository, tx application.Transactions, hasher application.PasswordHasher, sender d.Sender, policy d.Policy) *Service {
	return &Service{repo: repo, tx: tx, hasher: hasher, sender: sender, policy: policy, obligations: d.UnavailableObligations{}, Now: func() time.Time { return time.Now().UTC() }}
}
func (s *Service) SetObligationReader(reader d.ObligationReader) { s.obligations = reader }
func digest(raw string) string                                   { h := sha256.Sum256([]byte(raw)); return hex.EncodeToString(h[:]) }
func (s *Service) authorize(ctx context.Context, actor uuid.UUID, admin bool, targets ...uuid.UUID) error {
	ids := append([]uuid.UUID{actor}, targets...)
	if actor == uuid.Nil {
		return shared.ErrUnauthorized
	}
	if err := s.repo.LockAccounts(ctx, ids); err != nil {
		return err
	}
	a, err := s.repo.Get(ctx, actor)
	if err != nil {
		return err
	}
	if !a.IsActive || !a.Role.Valid() || (admin && a.Role != user.RoleAdmin) || (!admin && a.Role != user.RoleAdmin && a.Role != user.RoleStaff) {
		return shared.ErrForbidden
	}
	return nil
}
func (s *Service) command(ctx context.Context, actor uuid.UUID, operation, key string, payload any, targets []uuid.UUID, identity bool, fn func(context.Context) (any, error)) (raw []byte, first bool, err error) {
	if _, e := uuid.Parse(key); e != nil {
		return nil, false, shared.ErrInvalidInput
	}
	data, e := json.Marshal(payload)
	if e != nil {
		return nil, false, shared.ErrInvalidInput
	}
	hash := digest(string(data))
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if identity {
			if e := s.repo.LockIdentity(ctx); e != nil {
				return e
			}
		}
		if e := s.authorize(ctx, actor, true, targets...); e != nil {
			return e
		}
		cached, e := s.repo.ReadReceipt(ctx, actor, operation, key, hash)
		if e != nil {
			return e
		}
		if cached != nil {
			raw = cached
			return nil
		}
		result, e := fn(ctx)
		if e != nil {
			return e
		}
		raw, e = json.Marshal(result)
		if e != nil {
			return shared.ErrInternal
		}
		if e = s.repo.WriteReceipt(ctx, actor, operation, key, hash, raw); e != nil {
			return e
		}
		first = true
		return nil
	})
	return
}
func (s *Service) enrich(ctx context.Context, r d.Record) d.Record {
	o, e := s.obligations.Read(ctx, r.ID)
	if e != nil {
		o = d.Obligations{Availability: "UNAVAILABLE"}
	}
	r.Obligations = o
	return r
}
func (s *Service) List(ctx context.Context, actor uuid.UUID, f d.Filter) (out d.Page, err error) {
	if f.EligibleForIssuance && f.Staff {
		return out, shared.ErrInvalidInput
	}
	if f.Page < 1 || f.Page > 10000 || f.PerPage < 1 || f.PerPage > 100 || len(f.Search) > 100 || (f.BorrowerType != "" && f.BorrowerType != "STUDENT" && f.BorrowerType != "FACULTY") || (f.Status != "" && f.Status != "ACTIVE" && f.Status != "INACTIVE" && f.Status != "PENDING") {
		return out, shared.ErrInvalidInput
	}
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if e := s.authorize(ctx, actor, f.Staff); e != nil {
			return e
		}
		var e error
		out, e = s.repo.List(ctx, f)
		if e == nil {
			for i := range out.Items {
				out.Items[i] = s.enrich(ctx, out.Items[i])
			}
		}
		return e
	})
	return
}
func (s *Service) Detail(ctx context.Context, actor, id uuid.UUID, staff bool) (out d.Record, err error) {
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if e := s.authorize(ctx, actor, staff, id); e != nil {
			return e
		}
		var e error
		out, e = s.repo.Get(ctx, id)
		if e != nil {
			return e
		}
		if (!staff && out.Role != user.RoleBorrower) || (staff && out.Role == user.RoleBorrower) {
			return shared.ErrNotFound
		}
		out = s.enrich(ctx, out)
		return nil
	})
	return
}
func (s *Service) Audits(ctx context.Context, actor, id uuid.UUID, page int) (out []d.Audit, err error) {
	if page < 1 || page > 10000 {
		return nil, shared.ErrInvalidInput
	}
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if e := s.authorize(ctx, actor, true, id); e != nil {
			return e
		}
		var e error
		out, e = s.repo.Audits(ctx, id, page)
		return e
	})
	return
}
func (s *Service) token(ctx context.Context, id uuid.UUID) (raw string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", shared.ErrInternal
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	if err = s.repo.InvalidateTokens(ctx, id); err != nil {
		return "", err
	}
	now := s.Now()
	err = s.repo.PutToken(ctx, d.Token{Hash: digest(raw), AccountID: id, CreatedAt: now, ExpiresAt: now.Add(s.policy.ActivationTTL)})
	return
}
func (s *Service) insert(ctx context.Context, in d.Input, staff bool) (r d.Record, err error) {
	existing, e := s.repo.FindIdentity(ctx, in.Email, in.StudentID)
	if e != nil {
		return r, e
	}
	if len(existing) > 0 {
		for _, a := range existing {
			if a.Email == in.Email {
				return r, user.ErrEmailAlreadyExists
			}
			if in.StudentID != "" && a.StudentID == in.StudentID {
				return r, d.ErrStudentIDExists
			}
		}
		return r, shared.ErrConflict
	}
	now := s.Now()
	role := user.RoleBorrower
	if staff {
		role = user.RoleStaff
	}
	r = d.Record{ID: uuid.New(), Name: in.Name, Email: in.Email, Role: role, IsActive: true, BorrowerType: in.BorrowerType, StudentID: in.StudentID, Course: in.Course, ContactNumber: in.ContactNumber, ActivationRequired: true, DeliveryStatus: "PENDING", CreatedAt: now, UpdatedAt: now}
	if err = s.repo.Insert(ctx, r); err != nil {
		return r, err
	}
	if err = s.repo.AddAudit(ctx, actorFrom(ctx), r.ID, "ACCOUNT_CREATED"); err != nil {
		return r, err
	}
	return s.repo.Get(ctx, r.ID)
}

type actorKey struct{}

func actorFrom(ctx context.Context) uuid.UUID { id, _ := ctx.Value(actorKey{}).(uuid.UUID); return id }
func (s *Service) Create(ctx context.Context, actor uuid.UUID, key string, in d.Input, staff bool) (out d.Record, err error) {
	in, err = s.policy.Validate(in, staff)
	if err != nil {
		return
	}
	var rawToken string
	operation := "borrower.create"
	if staff {
		operation = "staff.create"
	}
	raw, first, e := s.command(ctx, actor, operation, key, in, nil, true, func(ctx context.Context) (any, error) {
		r, e := s.insert(context.WithValue(ctx, actorKey{}, actor), in, staff)
		if e != nil {
			return nil, e
		}
		rawToken, e = s.token(ctx, r.ID)
		return r, e
	})
	if e != nil {
		return out, e
	}
	if json.Unmarshal(raw, &out) != nil {
		return out, shared.ErrInternal
	}
	if first {
		s.deliver(ctx, actor, out.ID, rawToken)
	}
	return s.current(ctx, out.ID)
}
func (s *Service) current(ctx context.Context, id uuid.UUID) (d.Record, error) {
	r, e := s.repo.Get(ctx, id)
	if e != nil {
		return r, e
	}
	return s.enrich(ctx, r), nil
}
func (s *Service) deliver(ctx context.Context, actor, id uuid.UUID, raw string) {
	// Credentials live only in this call. Provider acceptance is not delivery evidence.
	deliveryCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 8*time.Second)
	defer cancel()
	_ = s.tx.Within(deliveryCtx, func(ctx context.Context) error {
		if e := s.repo.LockAccounts(ctx, []uuid.UUID{actor, id}); e != nil {
			return e
		}
		a, e := s.repo.Get(ctx, id)
		if e != nil {
			return e
		}
		t, e := s.repo.FindToken(ctx, digest(raw), true)
		if e != nil {
			return e
		}
		if !a.IsActive || !a.ActivationRequired || t.Invalid || !t.ExpiresAt.After(s.Now()) {
			return nil
		}
		link, e := url.Parse(s.policy.ActivationURL)
		if e != nil {
			return shared.ErrInternal
		}
		link.Fragment = "token=" + raw
		state, message := "UNCONFIGURED", ""
		sendErr := d.ErrMailUnavailable
		if s.sender != nil {
			message, sendErr = s.sender.SendActivation(ctx, d.Mail{To: a.Email, Name: a.Name, ActivationURL: link.String()})
		}
		switch {
		case sendErr == nil:
			state = "ACCEPTED"
		case errors.Is(sendErr, d.ErrMailUnavailable):
			state = "UNCONFIGURED"
		case errors.Is(sendErr, d.ErrMailUnknown):
			state = "UNKNOWN"
		default:
			state = "FAILED"
		}
		if e = s.repo.SetDelivery(ctx, id, t.Hash, state, message); e != nil {
			return e
		}
		return s.repo.AddAudit(ctx, actor, id, "ACTIVATION_SUBMISSION_"+state)
	})
}

type ProfileInput struct {
	Name              string    `json:"name"`
	Course            string    `json:"course"`
	ContactNumber     string    `json:"contact_number"`
	ExpectedUpdatedAt time.Time `json:"expected_updated_at"`
}

func (s *Service) UpdateProfile(ctx context.Context, actor, id uuid.UUID, key string, in ProfileInput) (out d.Record, err error) {
	name, e := user.ParseName(in.Name)
	if e != nil {
		return out, e
	}
	in.Name = name.String()
	if len(in.Course) > 128 || len(in.ContactNumber) > 32 || in.ExpectedUpdatedAt.IsZero() {
		return out, shared.ErrInvalidInput
	}
	raw, _, e := s.command(ctx, actor, "borrower.profile:"+id.String(), key, in, []uuid.UUID{id}, false, func(ctx context.Context) (any, error) {
		r, e := s.repo.Get(ctx, id)
		if e != nil {
			return nil, e
		}
		if r.Role != user.RoleBorrower {
			return nil, shared.ErrForbidden
		}
		if !r.UpdatedAt.Equal(in.ExpectedUpdatedAt) {
			return nil, shared.ErrConflict
		}
		if r.BorrowerType == "" && (in.Course != "" || in.ContactNumber != "") {
			return nil, shared.ErrInvalidInput
		}
		r.Name, r.Course, r.ContactNumber = in.Name, in.Course, in.ContactNumber
		if e = s.repo.UpdateProfile(ctx, r); e != nil {
			return nil, e
		}
		if e = s.repo.AddAudit(ctx, actor, id, "PROFILE_UPDATED"); e != nil {
			return nil, e
		}
		return s.repo.Get(ctx, id)
	})
	if e != nil {
		return out, e
	}
	err = json.Unmarshal(raw, &out)
	return
}

type StatusInput struct {
	Active            bool      `json:"active"`
	Confirm           bool      `json:"confirm"`
	ExpectedUpdatedAt time.Time `json:"expected_updated_at"`
}

func (s *Service) Status(ctx context.Context, actor, id uuid.UUID, key string, in StatusInput, staff bool) (out d.Record, err error) {
	if !in.Confirm || in.ExpectedUpdatedAt.IsZero() {
		return out, shared.ErrInvalidInput
	}
	op := "borrower.status:" + id.String()
	if staff {
		op = "staff.status:" + id.String()
	}
	raw, _, e := s.command(ctx, actor, op, key, in, []uuid.UUID{id}, false, func(ctx context.Context) (any, error) {
		r, e := s.repo.Get(ctx, id)
		if e != nil {
			return nil, e
		}
		if (!staff && r.Role != user.RoleBorrower) || (staff && r.Role != user.RoleStaff) {
			return nil, shared.ErrForbidden
		}
		if !r.UpdatedAt.Equal(in.ExpectedUpdatedAt) {
			return nil, shared.ErrConflict
		}
		if e = s.repo.SetActive(ctx, id, in.Active); e != nil {
			return nil, e
		}
		if !in.Active {
			if e = s.repo.InvalidateTokens(ctx, id); e != nil {
				return nil, e
			}
		}
		action := "ACCOUNT_DEACTIVATED"
		if in.Active {
			action = "ACCOUNT_REACTIVATED"
		}
		if e = s.repo.AddAudit(ctx, actor, id, action); e != nil {
			return nil, e
		}
		return s.repo.Get(ctx, id)
	})
	if e != nil {
		return out, e
	}
	if json.Unmarshal(raw, &out) != nil {
		return out, shared.ErrInternal
	}
	return s.current(ctx, id)
}
func (s *Service) Resend(ctx context.Context, actor, id uuid.UUID, key string) (out d.Record, err error) {
	var token string
	raw, first, e := s.command(ctx, actor, "activation.resend", key, id, []uuid.UUID{id}, false, func(ctx context.Context) (any, error) {
		r, e := s.repo.Get(ctx, id)
		if e != nil {
			return nil, e
		}
		if !r.IsActive || !r.ActivationRequired || r.Role == user.RoleAdmin {
			return nil, shared.ErrForbidden
		}
		old, e := s.repo.FindToken(ctx, "account:"+id.String(), false)
		if e != nil && !errors.Is(e, shared.ErrNotFound) {
			return nil, e
		}
		if e == nil && s.Now().Sub(old.CreatedAt) < s.policy.ResendCooldown {
			return nil, d.ErrCooldown
		}
		token, e = s.token(ctx, id)
		if e != nil {
			return nil, e
		}
		if e = s.repo.AddAudit(ctx, actor, id, "ACTIVATION_REISSUED"); e != nil {
			return nil, e
		}
		return r, nil
	})
	if e != nil {
		return out, e
	}
	if json.Unmarshal(raw, &out) != nil {
		return out, shared.ErrInternal
	}
	if first {
		s.deliver(ctx, actor, id, token)
	}
	return s.current(ctx, id)
}
func (s *Service) Activate(ctx context.Context, raw, password string) error {
	if len(raw) != 43 || strings.ContainsAny(raw, " \t\r\n") {
		return d.ErrActivationInvalid
	}
	if len(password) > 72 {
		return shared.ErrInvalidInput
	}
	if e := user.PlainPassword(password).Validate(); e != nil {
		return e
	}
	hash := digest(raw)
	initial, e := s.repo.FindToken(ctx, hash, false)
	if e != nil {
		return d.ErrActivationInvalid
	}
	hashed, e := s.hasher.Hash(password)
	if e != nil {
		return shared.ErrInternal
	}
	return s.tx.Within(ctx, func(ctx context.Context) error {
		if e := s.repo.LockAccounts(ctx, []uuid.UUID{initial.AccountID}); e != nil {
			return e
		}
		t, e := s.repo.FindToken(ctx, hash, true)
		if e != nil || t.Invalid || !t.ExpiresAt.After(s.Now()) {
			return d.ErrActivationInvalid
		}
		a, e := s.repo.Get(ctx, t.AccountID)
		if e != nil {
			return d.ErrActivationInvalid
		}
		if !a.IsActive || !a.ActivationRequired || !a.Role.Valid() || a.Role == user.RoleAdmin {
			return d.ErrActivationInvalid
		}
		if a.BorrowerType == "STUDENT" {
			_, e = s.policy.Validate(d.Input{Name: a.Name, Email: a.Email, BorrowerType: a.BorrowerType, StudentID: a.StudentID}, false)
			if e != nil {
				return e
			}
		}
		if e = s.repo.SetPassword(ctx, a.ID, hashed); e != nil {
			return e
		}
		if e = s.repo.InvalidateTokens(ctx, a.ID); e != nil {
			return e
		}
		return s.repo.AddAudit(ctx, a.ID, a.ID, "ACCOUNT_ACTIVATED")
	})
}
func (s *Service) Prepare(ctx context.Context, actor uuid.UUID, operation string, rows []d.Row) (out d.Batch, err error) {
	if operation != "CREATE" && operation != "DEACTIVATE" || len(rows) == 0 || len(rows) > 500 {
		return out, shared.ErrInvalidInput
	}
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if e := s.authorize(ctx, actor, true); e != nil {
			return e
		}
		out = d.Batch{ID: uuid.New(), ActorID: actor, Operation: operation, Rows: rows, ExpiresAt: s.Now().Add(30 * time.Minute)}
		emails, ids := map[string]int{}, map[string]int{}
		for i := range out.Rows {
			r := &out.Rows[i]
			if r.Error != "" {
				continue
			}
			validationPolicy := s.policy
			if operation == "DEACTIVATE" {
				if at := strings.LastIndex(r.Input.Email, "@"); at >= 0 {
					validationPolicy.StudentDomains = []string{strings.ToLower(strings.TrimSpace(r.Input.Email[at+1:]))}
				}
			}
			in, e := validationPolicy.Validate(r.Input, false)
			r.Input = in
			if e != nil {
				switch {
				case errors.Is(e, d.ErrDomainsMissing):
					r.Error = "Student email domains are not configured"
				case errors.Is(e, d.ErrStudentDomain):
					r.Error = "Email must use a configured institutional domain"
				case errors.Is(e, user.ErrInvalidEmail):
					r.Error = "Enter a valid email address"
				case errors.Is(e, user.ErrInvalidName):
					r.Error = "Enter a valid Student name"
				case strings.TrimSpace(r.Input.StudentID) == "":
					r.Error = "Student ID is required"
				default:
					r.Error = "Invalid required identity fields or field length"
				}
				continue
			}
			if in.BorrowerType != "STUDENT" {
				r.Error = "Student roster only"
				continue
			}
			emails[in.Email]++
			ids[in.StudentID]++
		}
		for i := range out.Rows {
			r := &out.Rows[i]
			r.Obligations = d.Obligations{Availability: "UNAVAILABLE"}
			if r.Error != "" {
				continue
			}
			if emails[r.Input.Email] > 1 || ids[r.Input.StudentID] > 1 {
				r.Error = "Duplicate email or Student ID in roster"
				continue
			}
			matches, e := s.repo.FindIdentity(ctx, r.Input.Email, r.Input.StudentID)
			if e != nil {
				return e
			}
			if operation == "CREATE" {
				if len(matches) > 0 {
					r.Error = "Existing account or conflicting identity"
				}
			} else {
				if len(matches) != 1 {
					r.Error = "Unmatched or conflicting identifiers"
					continue
				}
				m := matches[0]
				if m.Role != user.RoleBorrower || m.BorrowerType != "STUDENT" || m.Email != r.Input.Email || m.StudentID != r.Input.StudentID {
					r.Error = "Conflicting identifiers or non-Student account"
					continue
				}
				r.TargetID = &m.ID
				r.TargetVersion = &m.UpdatedAt
				r.Obligations = s.enrich(ctx, m).Obligations
				if !m.IsActive {
					r.Error = "Already inactive"
				}
			}
		}
		return s.repo.PutBatch(ctx, out)
	})
	return
}
func (s *Service) Preview(ctx context.Context, actor, id uuid.UUID) (out d.Batch, err error) {
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if e := s.authorize(ctx, actor, true); e != nil {
			return e
		}
		var e error
		out, e = s.repo.GetBatch(ctx, id, false)
		if e == nil && out.ActorID != actor {
			return shared.ErrNotFound
		}
		return e
	})
	return
}
func (s *Service) Confirm(ctx context.Context, actor, id uuid.UUID, selection []int, confirm bool) (out d.Batch, err error) {
	if !confirm || len(selection) == 0 || len(selection) > 500 {
		return out, shared.ErrInvalidInput
	}
	sort.Ints(selection)
	for i, n := range selection {
		if n < 2 || (i > 0 && n == selection[i-1]) {
			return out, shared.ErrInvalidInput
		}
	}
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if e := s.repo.LockIdentity(ctx); e != nil {
			return e
		}
		snapshot, e := s.repo.GetBatch(ctx, id, false)
		if e != nil {
			return e
		}
		if snapshot.ActorID != actor {
			return shared.ErrNotFound
		}
		targets := []uuid.UUID{}
		for _, r := range snapshot.Rows {
			if r.TargetID != nil {
				targets = append(targets, *r.TargetID)
			}
		}
		if e = s.authorize(ctx, actor, true, targets...); e != nil {
			return e
		}
		out, e = s.repo.GetBatch(ctx, id, true)
		if e != nil {
			return e
		}
		if out.ActorID != actor {
			return shared.ErrNotFound
		}
		if out.Confirmed {
			a, _ := json.Marshal(out.Selection)
			b, _ := json.Marshal(selection)
			if string(a) != string(b) {
				return shared.ErrConflict
			}
			return nil
		}
		if !out.ExpiresAt.After(s.Now()) {
			return shared.ErrConflict
		}
		selected := map[int]bool{}
		for _, n := range selection {
			selected[n] = true
		}
		seen := 0
		// Identity writers share a lock; account/status writers share sorted row locks.
		for i := range out.Rows {
			r := &out.Rows[i]
			if !selected[r.Number] {
				r.Outcome = "NOT_SELECTED"
				continue
			}
			seen++
			if r.Error != "" {
				return shared.ErrConflict
			}
			validationPolicy := s.policy
			if out.Operation == "DEACTIVATE" {
				validationPolicy.StudentDomains = []string{r.Input.Email[strings.LastIndex(r.Input.Email, "@")+1:]}
			}
			in, e := validationPolicy.Validate(r.Input, false)
			if e != nil {
				return e
			}
			if in.BorrowerType != "STUDENT" {
				return shared.ErrForbidden
			}
			if out.Operation == "CREATE" {
				m, e := s.insert(context.WithValue(ctx, actorKey{}, actor), in, false)
				if e != nil {
					return e
				}
				r.AccountID = &m.ID
				r.Outcome = "CREATED_ACTIVATION_PENDING"
			} else {
				if r.TargetID == nil || r.TargetVersion == nil {
					return shared.ErrConflict
				}
				m, e := s.repo.Get(ctx, *r.TargetID)
				if e != nil {
					return e
				}
				if m.Role != user.RoleBorrower || m.BorrowerType != "STUDENT" || m.Email != in.Email || m.StudentID != in.StudentID || !m.IsActive || !m.UpdatedAt.Equal(*r.TargetVersion) {
					return shared.ErrConflict
				}
				if e = s.repo.SetActive(ctx, m.ID, false); e != nil {
					return e
				}
				if e = s.repo.InvalidateTokens(ctx, m.ID); e != nil {
					return e
				}
				if e = s.repo.AddAudit(ctx, actor, m.ID, "ROSTER_STUDENT_DEACTIVATED"); e != nil {
					return e
				}
				r.AccountID = &m.ID
				r.Outcome = "DEACTIVATED"
			}
		}
		if seen != len(selection) {
			return shared.ErrInvalidInput
		}
		out.Selection = selection
		out.Confirmed = true
		return s.repo.ConfirmBatch(ctx, out)
	})
	return
}

// Configuration availability is separate from institutionally verified compliance.
func (s *Service) ProvisioningPolicy() any {
	return struct {
		StudentDomains        []string `json:"student_domains"`
		ActivationTTLSeconds  int64    `json:"activation_ttl_seconds"`
		ResendCooldownSeconds int64    `json:"resend_cooldown_seconds"`
	}{s.policy.StudentDomains, int64(s.policy.ActivationTTL.Seconds()), int64(s.policy.ResendCooldown.Seconds())}
}
