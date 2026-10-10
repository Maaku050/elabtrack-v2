package borrowing

import (
	"context"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

func (s *Service) History(c context.Context, actor, id uuid.UUID, kind string, page, size int) (out d.HistoryPage, e error) {
	if !d.HistoryKind(kind) || page < 1 || page > 100000 || size < 1 || size > 100 {
		return out, shared.ErrInvalidInput
	}
	r, ok := s.repo.(d.HistoryRepository)
	if !ok {
		return out, shared.ErrConflict
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
		v, e := r.Preview(c, id, true)
		if e != nil {
			return e
		}
		if u.Role == user.RoleBorrower && v.BorrowerID != actor {
			return shared.ErrNotFound
		}
		out, e = r.History(c, id, kind, page, size)
		return e
	})
	return
}
