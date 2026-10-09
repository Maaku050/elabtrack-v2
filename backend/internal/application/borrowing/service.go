package borrowing

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	a "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
	"strings"
)

type Acceptance interface {
	RequireCurrentAcceptance(context.Context, uuid.UUID) (*terms.Acceptance, error)
}
type Service struct {
	repo     d.Repository
	accounts a.Repository
	stock    inventory.Repository
	terms    Acceptance
	tx       application.Transactions
}

func NewService(r d.Repository, a a.Repository, i inventory.Repository, t Acceptance, tx application.Transactions) *Service {
	return &Service{r, a, i, t, tx}
}
func (s *Service) actor(c context.Context, id uuid.UUID, staff bool) (a.Record, error) {
	v, e := s.accounts.Get(c, id)
	if e != nil {
		return v, e
	}
	if !v.IsActive || v.ActivationRequired {
		return v, shared.ErrUnauthorized
	}
	if staff && (v.Role != user.RoleStaff && v.Role != user.RoleAdmin) || !staff && v.Role != user.RoleBorrower {
		return v, shared.ErrForbidden
	}
	return v, nil
}
func (s *Service) eligible(c context.Context, id uuid.UUID) (a.Record, error) {
	v, e := s.accounts.Get(c, id)
	if e != nil {
		return v, e
	}
	if v.Role != user.RoleBorrower || !v.IsActive || v.ActivationRequired || (v.BorrowerType != "STUDENT" && v.BorrowerType != "FACULTY") {
		return v, d.ErrEligibility
	}
	return v, nil
}
func (s *Service) Read(c context.Context, actor, id uuid.UUID) (v d.Record, e error) {
	e = s.tx.Within(c, func(c context.Context) error {
		if e := s.accounts.LockAccounts(c, []uuid.UUID{actor}); e != nil {
			return e
		}
		u, e := s.accounts.Get(c, actor)
		if e != nil {
			return e
		}
		if !u.IsActive || u.ActivationRequired {
			return shared.ErrUnauthorized
		}
		if !u.Role.Valid() {
			return shared.ErrForbidden
		}
		// Hold the aggregate through item/history reads so another operator
		// cannot mix pre-transition headers with post-transition custody.
		v, e = s.repo.Get(c, id, true)
		if e != nil {
			return e
		}
		if u.Role == user.RoleBorrower && v.BorrowerID != actor {
			return shared.ErrNotFound
		}
		now, e := s.repo.Clock(c)
		if e == nil {
			v.Overdue = v.Status == "CHECKED_OUT" && v.DueAt != nil && now.After(*v.DueAt)
		}
		return e
	})
	return
}
func (s *Service) List(c context.Context, actor uuid.UUID, f d.Filter) (p d.Page, e error) {
	if !f.Valid() {
		return p, shared.ErrInvalidInput
	}
	e = s.tx.Within(c, func(c context.Context) error {
		if e := s.accounts.LockAccounts(c, []uuid.UUID{actor}); e != nil {
			return e
		}
		u, e := s.accounts.Get(c, actor)
		if e != nil {
			return e
		}
		if !u.IsActive || u.ActivationRequired {
			return shared.ErrUnauthorized
		}
		if !u.Role.Valid() {
			return shared.ErrForbidden
		}
		if u.Role == user.RoleBorrower {
			f.Owner = &actor
		}
		p, e = s.repo.List(c, f)
		if e != nil {
			return e
		}
		now, e := s.repo.Clock(c)
		for i := range p.Items {
			v := &p.Items[i]
			v.Overdue = v.Status == "CHECKED_OUT" && v.DueAt != nil && now.After(*v.DueAt)
		}
		return e
	})
	return
}

// All command receipts follow account locks, so same-actor retries serialize before writes.
func (s *Service) command(c context.Context, actor, target uuid.UUID, op, key string, input any, staff, targetEligible bool, fn func(context.Context) (d.Record, error)) (out d.Record, err error) {
	k, e := uuid.Parse(key)
	if e != nil || k == uuid.Nil || actor == uuid.Nil {
		return out, shared.ErrInvalidInput
	}
	payload, e := json.Marshal(input)
	if e != nil {
		return out, shared.ErrInvalidInput
	}
	digest := sha256.Sum256(payload)
	hash := hex.EncodeToString(digest[:])
	err = s.tx.Within(c, func(c context.Context) error {
		if e := s.accounts.LockAccounts(c, []uuid.UUID{actor, target}); e != nil {
			return e
		}
		if _, e := s.actor(c, actor, staff); e != nil {
			return e
		}
		if targetEligible {
			if _, e := s.eligible(c, target); e != nil {
				return e
			}
		}
		data, e := s.repo.ReadReceipt(c, actor, op, k.String(), hash)
		if e == nil {
			return json.Unmarshal(data, &out)
		}
		if !errors.Is(e, shared.ErrNotFound) {
			return e
		}
		out, e = fn(c)
		if e != nil {
			return e
		}
		data, e = json.Marshal(out)
		if e != nil {
			return e
		}
		return s.repo.WriteReceipt(c, actor, op, k.String(), hash, data)
	})
	return
}
func (s *Service) Submit(c context.Context, actor uuid.UUID, key string, in d.Input) (d.Record, error) {
	if !in.Confirm || in.BorrowerID != uuid.Nil || in.DueAt != nil || in.Handover {
		return d.Record{}, shared.ErrInvalidInput
	}
	lines, e := d.Normalize(in.Items)
	if e != nil {
		return d.Record{}, e
	}
	in.Items = lines
	return s.command(c, actor, actor, "submit", key, in, false, true, func(c context.Context) (d.Record, error) { return s.create(c, actor, actor, in, false) })
}
func (s *Service) Direct(c context.Context, actor uuid.UUID, key string, in d.Input) (d.Record, error) {
	if in.BorrowerID == uuid.Nil || in.DueAt == nil || !in.Handover || !in.Confirm {
		return d.Record{}, shared.ErrInvalidInput
	}
	if in.DueAt != nil {
		utc := in.DueAt.UTC()
		in.DueAt = &utc
	}
	lines, e := d.Normalize(in.Items)
	if e != nil {
		return d.Record{}, e
	}
	in.Items = lines
	return s.command(c, actor, in.BorrowerID, "direct", key, in, true, true, func(c context.Context) (d.Record, error) { return s.create(c, actor, in.BorrowerID, in, true) })
}
func (s *Service) lockStock(c context.Context, lines []d.Line, usable bool) ([]inventory.Equipment, error) {
	out := make([]inventory.Equipment, 0, len(lines))
	for _, line := range lines {
		if e := s.stock.LockEquipment(c, line.EquipmentID); e != nil {
			return nil, e
		}
		v, e := s.stock.Get(c, line.EquipmentID)
		if e != nil {
			return nil, e
		}
		if usable && v.Status != "ACTIVE" {
			return nil, d.ErrStock
		}
		out = append(out, v)
	}
	return out, nil
}
func (s *Service) create(c context.Context, actor, target uuid.UUID, in d.Input, direct bool) (v d.Record, e error) {
	u, e := s.eligible(c, target)
	if e != nil {
		return v, e
	}
	accept, e := s.terms.RequireCurrentAcceptance(c, target)
	if e != nil {
		return v, e
	}
	stock, e := s.lockStock(c, in.Items, true)
	if e != nil {
		return v, e
	}
	now, e := s.repo.Clock(c)
	if e != nil {
		return v, e
	}
	if direct && (in.DueAt == nil || !in.DueAt.After(now)) {
		return v, shared.ErrInvalidInput
	}
	id := uuid.New()
	v = d.Record{ID: id, Reference: "BR-" + strings.ToUpper(id.String()), BorrowerID: target, BorrowerName: u.Name, BorrowerType: u.BorrowerType, StudentID: u.StudentID, AcceptanceID: accept.ID, EntryPath: "REQUEST", Status: "PENDING", CreatedAt: now, Items: []d.Item{}, Events: []d.Event{}}
	kind := "RESERVE"
	event := "SUBMITTED"
	expiry := now.Add(d.PendingTTL)
	v.ExpiresAt = &expiry
	if direct {
		kind = "DIRECT"
		event = "DIRECT"
		v.Status = "CHECKED_OUT"
		v.EntryPath = "DIRECT"
		v.ExpiresAt = nil
		v.CheckedOutAt = &now
		v.DueAt = in.DueAt
	}
	for n, l := range in.Items {
		i := d.Item{ID: uuid.New(), EquipmentID: l.EquipmentID, Name: stock[n].Name, Quantity: l.Quantity}
		if direct {
			i.Issued = l.Quantity
		} else {
			i.Reserved = l.Quantity
		}
		v.Items = append(v.Items, i)
		if _, e = d.Transfer(stock[n].Stock, l.Quantity, kind); e != nil {
			return v, e
		}
	}
	if e = s.repo.Insert(c, v); e != nil {
		return v, e
	}
	for n, i := range v.Items {
		before := stock[n].Stock
		stock[n].Stock, e = d.Transfer(before, i.Quantity, kind)
		if e != nil {
			return v, e
		}
		stock[n].Sequence++
		if e = s.repo.Effect(c, v, i, &actor, kind, stock[n], before); e != nil {
			return v, e
		}
	}
	if e = s.repo.Event(c, v.ID, d.Event{ID: uuid.New(), ActorID: &actor, Kind: event, At: now}); e != nil {
		return v, e
	}
	return s.repo.Get(c, v.ID, false)
}
func (s *Service) Decide(c context.Context, actor, id uuid.UUID, key, action string, in d.Decision) (v d.Record, err error) {
	if id == uuid.Nil || !in.Confirm || (action != "CANCELLED" && action != "DENIED" && action != "CHECKOUT") {
		return v, shared.ErrInvalidInput
	}
	if in.DueAt != nil {
		utc := in.DueAt.UTC()
		in.DueAt = &utc
	}
	in.Reason = strings.TrimSpace(in.Reason)
	if action == "DENIED" && !d.ValidReason(in.Reason) || action != "DENIED" && in.Reason != "" || action == "CHECKOUT" && (in.DueAt == nil || !in.Handover) || action != "CHECKOUT" && (in.DueAt != nil || in.Handover) {
		return v, shared.ErrInvalidInput
	}
	initial, e := s.repo.Get(c, id, false)
	if e != nil {
		return v, e
	}
	staff := action != "CANCELLED"
	// Target eligibility is checked after lazy expiry and before checkout; release is always permitted.
	v, err = s.command(c, actor, initial.BorrowerID, action+":"+id.String(), key, in, staff, false, func(c context.Context) (d.Record, error) {
		current, e := s.repo.Get(c, id, true)
		if e != nil {
			return current, e
		}
		if !staff && current.BorrowerID != actor {
			return d.Record{}, shared.ErrNotFound
		}
		if current.Status != "PENDING" {
			return current, d.ErrState
		}
		return s.transition(c, current, &actor, action, in)
	})
	if err == nil && v.Status == "EXPIRED" {
		err = d.ErrExpired
	}
	return
}
func (s *Service) transition(c context.Context, v d.Record, actor *uuid.UUID, action string, in d.Decision) (d.Record, error) {
	lines := make([]d.Line, len(v.Items))
	for n, i := range v.Items {
		lines[n] = d.Line{EquipmentID: i.EquipmentID, Quantity: i.Quantity}
	}
	stock, e := s.lockStock(c, lines, false)
	if e != nil {
		return v, e
	}
	now, e := s.repo.Clock(c)
	if e != nil {
		return v, e
	}
	if v.ExpiresAt != nil && !now.Before(*v.ExpiresAt) {
		action = "EXPIRED"
		actor = nil
	} else if action == "EXPIRED" {
		return v, nil
	}
	if action == "CHECKOUT" {
		if _, e = s.eligible(c, v.BorrowerID); e != nil {
			return v, e
		}
		if in.DueAt == nil || !in.DueAt.After(now) {
			return v, shared.ErrInvalidInput
		}
		for _, s := range stock {
			if s.Status != "ACTIVE" {
				return v, d.ErrStock
			}
		}
		v.Status = "CHECKED_OUT"
		v.CheckedOutAt = &now
		v.DueAt = in.DueAt
	} else {
		v.Status = action
		v.TerminalAt = &now
		if action == "DENIED" {
			v.Reason = in.Reason
		}
	}
	for n := range v.Items {
		i := &v.Items[n]
		if i.Reserved != i.Quantity || i.Issued != 0 {
			return v, d.ErrState
		}
		i.Reserved = 0
		if action == "CHECKOUT" {
			i.Issued = i.Quantity
		}
		before := stock[n].Stock
		stock[n].Stock, e = d.Transfer(before, i.Quantity, action)
		if e != nil {
			return v, e
		}
		stock[n].Sequence++
		if e = s.repo.Effect(c, v, *i, actor, action, stock[n], before); e != nil {
			return v, e
		}
	}
	if e = s.repo.Transition(c, v); e != nil {
		return v, e
	}
	event := action
	if action == "CHECKOUT" {
		event = "CHECKED_OUT"
	}
	if e = s.repo.Event(c, v.ID, d.Event{ID: uuid.New(), ActorID: actor, Kind: event, Reason: v.Reason, At: now}); e != nil {
		return v, e
	}
	return s.repo.Get(c, v.ID, false)
}

// Sweep uses persisted deadlines and the SAME account→borrowing→equipment lock order.
func (s *Service) Sweep(c context.Context, limit int) (int, error) {
	if limit < 1 || limit > 100 {
		return 0, shared.ErrInvalidInput
	}
	ids, e := s.repo.Due(c, limit)
	if e != nil {
		return 0, e
	}
	count := 0
	for _, id := range ids {
		initial, e := s.repo.Get(c, id, false)
		if e != nil {
			return count, e
		}
		expired := false
		e = s.tx.Within(c, func(c context.Context) error {
			if e := s.accounts.LockAccounts(c, []uuid.UUID{initial.BorrowerID}); e != nil {
				return e
			}
			v, e := s.repo.Get(c, id, true)
			if e != nil {
				return e
			}
			if v.Status != "PENDING" {
				return nil
			}
			v, e = s.transition(c, v, nil, "EXPIRED", d.Decision{})
			expired = e == nil && v.Status == "EXPIRED"
			return e
		})
		if e != nil {
			return count, e
		}
		if expired {
			count++
		}
	}
	return count, nil
}
