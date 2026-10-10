package borrowing

import (
	"context"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/inventory"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
	"strings"
)

func (s *Service) accountability() (d.AccountabilityRepository, error) {
	r, ok := s.repo.(d.AccountabilityRepository)
	if !ok {
		return nil, shared.ErrConflict
	}
	return r, nil
}
func (s *Service) Return(c context.Context, actor, id uuid.UUID, key string, in d.ReturnInput) (d.Record, error) {
	if !in.Confirm || in.ExpectedEvents < 1 || !d.ValidReason(in.Reason) {
		return d.Record{}, shared.ErrInvalidInput
	}
	lines, e := d.NormalizeReturns(in.Lines)
	if e != nil {
		return d.Record{}, e
	}
	in.Lines = lines
	in.Reason = strings.TrimSpace(in.Reason)
	initial, e := s.repo.Get(c, id, false)
	if e != nil {
		return initial, e
	}
	return s.command(c, actor, initial.BorrowerID, "return:"+id.String(), key, in, true, false, func(c context.Context) (d.Record, error) {
		r, e := s.accountability()
		if e != nil {
			return d.Record{}, e
		}
		v, e := s.repo.Get(c, id, true)
		if e != nil {
			return v, e
		}
		if v.Status != "CHECKED_OUT" || len(v.Events) != in.ExpectedEvents {
			return v, d.ErrState
		}
		items := map[uuid.UUID]d.Item{}
		stockLines := []d.Line{}
		for _, i := range v.Items {
			items[i.ID] = i
		}
		for _, l := range in.Lines {
			i, ok := items[l.ItemID]
			if !ok {
				return v, shared.ErrInvalidInput
			}
			if l.Good+l.Damaged+l.Lost > i.Outstanding {
				return v, d.ErrState
			}
			stockLines = append(stockLines, d.Line{EquipmentID: i.EquipmentID, Quantity: 1})
		}
		stockLines, e = d.Normalize(stockLines)
		if e != nil {
			return v, e
		}
		stock, e := s.lockStock(c, stockLines, false)
		if e != nil {
			return v, e
		}
		byID := map[uuid.UUID]inventory.Equipment{}
		for _, eq := range stock {
			byID[eq.ID] = eq
		}
		now, e := s.repo.Clock(c)
		if e != nil {
			return v, e
		}
		for _, l := range in.Lines {
			i := items[l.ItemID]
			eq := byID[i.EquipmentID]
			before := eq.Stock
			eq.Stock, e = d.ReturnStock(before, l)
			if e != nil {
				return v, e
			}
			eq.Sequence++
			disp := d.Disposition{ID: uuid.New(), ReturnLine: l, ActorID: actor, Reason: in.Reason, At: now}
			if e = r.RecordReturn(c, v, i, disp, eq, before); e != nil {
				return v, e
			}
		}
		if e = s.repo.Event(c, id, d.Event{ID: uuid.New(), ActorID: &actor, Kind: "RETURN", Reason: in.Reason, At: now}); e != nil {
			return v, e
		}
		return s.finish(c, id, actor)
	})
}
func (s *Service) Replace(c context.Context, actor, id uuid.UUID, key string, in d.ReplacementInput) (d.Record, error) {
	if !in.Confirm || in.ExpectedEvents < 1 || !in.Equivalent || in.ObligationID == uuid.Nil || in.Quantity <= 0 || in.Quantity > inventory.MaxQuantity || !d.ValidReason(in.Reason) {
		return d.Record{}, shared.ErrInvalidInput
	}
	in.Reason = strings.TrimSpace(in.Reason)
	initial, e := s.repo.Get(c, id, false)
	if e != nil {
		return initial, e
	}
	return s.command(c, actor, initial.BorrowerID, "replace:"+id.String(), key, in, true, false, func(c context.Context) (d.Record, error) {
		r, e := s.accountability()
		if e != nil {
			return d.Record{}, e
		}
		v, e := s.repo.Get(c, id, true)
		if e != nil {
			return v, e
		}
		if v.Status != "CHECKED_OUT" || len(v.Events) != in.ExpectedEvents {
			return v, d.ErrState
		}
		var o d.Obligation
		found := false
		for _, candidate := range v.Obligations {
			if candidate.ID == in.ObligationID {
				o = candidate
				found = true
				break
			}
		}
		if !found {
			return v, shared.ErrInvalidInput
		}
		if in.Quantity > o.Required-o.Accepted {
			return v, d.ErrState
		}
		stock, e := s.lockStock(c, []d.Line{{EquipmentID: o.EquipmentID, Quantity: in.Quantity}}, false)
		if e != nil {
			return v, e
		}
		eq := stock[0]
		before := eq.Stock
		eq.Stock.Available += in.Quantity
		eq.Stock.Total += in.Quantity
		eq.Sequence++
		if !eq.Stock.Valid() {
			return v, d.ErrStock
		}
		now, e := s.repo.Clock(c)
		if e != nil {
			return v, e
		}
		acceptance := d.Replacement{EquipmentID: eq.ID, EquipmentName: eq.Name, Equivalent: true, ID: uuid.New(), ObligationID: o.ID, Quantity: in.Quantity, ActorID: actor, Reason: in.Reason, At: now}
		if e = r.RecordReplacement(c, v, o, acceptance, eq, before); e != nil {
			return v, e
		}
		if e = s.repo.Event(c, id, d.Event{ID: uuid.New(), ActorID: &actor, Kind: "REPLACEMENT", Reason: in.Reason, At: now}); e != nil {
			return v, e
		}
		return s.finish(c, id, actor)
	})
}
func (s *Service) finish(c context.Context, id, actor uuid.UUID) (d.Record, error) {
	v, e := s.repo.Get(c, id, false)
	if e != nil {
		return v, e
	}
	now, e := s.repo.Clock(c)
	if e != nil {
		return v, e
	}
	complete := true
	for _, i := range v.Items {
		if i.Outstanding != 0 {
			complete = false
		}
	}
	for _, o := range v.Obligations {
		if o.Required != o.Accepted {
			complete = false
		}
	}
	if complete {
		r, e := s.accountability()
		if e != nil {
			return v, e
		}
		if e = v.ProjectFine(now); e != nil {
			return v, e
		}
		if e = r.Complete(c, v, now, v.Fine.Assessed); e != nil {
			return v, e
		}
		if e = s.repo.Event(c, id, d.Event{ID: uuid.New(), ActorID: &actor, Kind: "COMPLETED", At: now}); e != nil {
			return v, e
		}
		v, e = s.repo.Get(c, id, false)
		if e != nil {
			return v, e
		}
	}
	e = v.ProjectFine(now)
	return v, e
}
func (s *Service) ClearFine(c context.Context, actor, id uuid.UUID, key string, in d.FineInput) (d.Record, error) {
	if !in.Confirm || in.ExpectedMinor <= 0 || !inventory.Text(in.Note, 1000) || in.Method != "PAID" && in.Method != "WAIVED" && in.Method != "OTHER_RESOLUTION" {
		return d.Record{}, shared.ErrInvalidInput
	}
	initial, e := s.repo.Get(c, id, false)
	if e != nil {
		return initial, e
	}
	// Admin authority is rechecked inside the account-lock transaction, including replay.
	return s.command(c, actor, initial.BorrowerID, "clear:"+id.String(), key, in, true, false, func(c context.Context) (d.Record, error) {
		u, e := s.accounts.Get(c, actor)
		if e != nil {
			return d.Record{}, e
		}
		if u.Role != user.RoleAdmin {
			return d.Record{}, shared.ErrForbidden
		}
		r, e := s.accountability()
		if e != nil {
			return d.Record{}, e
		}
		v, e := s.repo.Get(c, id, true)
		if e != nil {
			return v, e
		}
		if v.Status != "CHECKED_OUT" && v.Status != "COMPLETED" {
			return v, d.ErrState
		}
		now, e := s.repo.Clock(c)
		if e != nil {
			return v, e
		}
		if e = v.ProjectFine(now); e != nil {
			return v, e
		}
		if v.Fine.Outstanding != in.ExpectedMinor {
			return v, d.ErrState
		}
		clear := d.Clearance{ID: uuid.New(), ActorID: actor, Assessed: v.Fine.Assessed, Amount: v.Fine.Outstanding, Method: in.Method, Note: in.Note, At: now}
		if e = r.ClearFine(c, v, clear); e != nil {
			return v, e
		}
		if e = s.repo.Event(c, id, d.Event{ID: uuid.New(), ActorID: &actor, Kind: "FINE_CLEARED", Reason: in.Method + ": " + in.Note, At: now}); e != nil {
			return v, e
		}
		v, e = s.repo.Get(c, id, false)
		if e != nil {
			return v, e
		}
		e = v.ProjectFine(now)
		return v, e
	})
}
