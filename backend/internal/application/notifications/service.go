package notifications

import (
	"context"
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	accounts "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/notifications"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

type Service struct {
	repo        d.Repository
	accounts    accounts.Repository
	tx          application.Transactions
	leadSeconds int64
}

func NewService(r d.Repository, a accounts.Repository, t application.Transactions, lead int64) *Service {
	return &Service{r, a, t, lead}
}
func (s *Service) actor(c context.Context, id uuid.UUID) (user.Role, error) {
	if e := s.accounts.LockAccounts(c, []uuid.UUID{id}); e != nil {
		return "", e
	}
	u, e := s.accounts.Get(c, id)
	if e != nil {
		return "", e
	}
	if !u.IsActive || u.ActivationRequired {
		return "", shared.ErrUnauthorized
	}
	if !u.Role.Valid() {
		return "", shared.ErrForbidden
	}
	return u.Role, nil
}
func (s *Service) List(c context.Context, id uuid.UUID, f d.Filter) (p d.Page, e error) {
	if !f.Valid() {
		return p, shared.ErrInvalidInput
	}
	e = s.tx.Within(c, func(c context.Context) error {
		role, err := s.actor(c, id)
		if err != nil {
			return err
		}
		p, err = s.repo.List(c, id, role, f)
		return err
	})
	return
}
func (s *Service) Count(c context.Context, id uuid.UUID) (n int, e error) {
	e = s.tx.Within(c, func(c context.Context) error {
		role, err := s.actor(c, id)
		if err != nil {
			return err
		}
		n, err = s.repo.Count(c, id, role)
		return err
	})
	return
}
func (s *Service) Mark(c context.Context, id, notification uuid.UUID, read bool) (v d.Notification, e error) {
	if notification == uuid.Nil {
		return v, shared.ErrInvalidInput
	}
	e = s.tx.Within(c, func(c context.Context) error {
		role, err := s.actor(c, id)
		if err != nil {
			return err
		}
		v, err = s.repo.Mark(c, id, role, notification, read)
		return err
	})
	return
}
func (s *Service) Sweep(c context.Context, limit int) (n int, e error) {
	if limit < 1 || limit > 100 || s.leadSeconds < 0 || s.leadSeconds > 604800 {
		return 0, shared.ErrInvalidInput
	}
	e = s.tx.Within(c, func(c context.Context) error {
		locked, err := s.repo.TryWorkerLock(c)
		if err != nil || !locked {
			return err
		}
		events, err := s.repo.Pending(c, limit, s.leadSeconds)
		if err != nil {
			return err
		}
		for _, event := range events {
			message := event.Message()
			if message.Title == "" {
				return shared.ErrInvalidInput
			}
			created, err := s.repo.Emit(c, event, message)
			if err != nil {
				return err
			}
			if created {
				n++
			}
		}
		return nil
	})
	if e != nil {
		n = 0
	}
	return
}

// MarkAll affects the same current, authorized account scope as List and Count.
func (s *Service) MarkAll(c context.Context, id uuid.UUID) (n int64, e error) {
	e = s.tx.Within(c, func(c context.Context) error {
		role, err := s.actor(c, id)
		if err != nil {
			return err
		}
		n, err = s.repo.MarkAll(c, id, role)
		return err
	})
	if e != nil {
		n = 0
	}
	return
}
