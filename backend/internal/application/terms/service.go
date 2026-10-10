package terms

import (
	"context"
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainterms "github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

type Service struct {
	repo     domainterms.Repository
	accounts application.RefreshAccounts
	tx       application.Transactions
}

func NewService(repo domainterms.Repository, accounts application.RefreshAccounts, tx application.Transactions) *Service {
	return &Service{repo, accounts, tx}
}

func (s *Service) lockAccount(ctx context.Context, id uuid.UUID, permission user.Permission) (*user.Account, error) {
	if id == uuid.Nil {
		return nil, shared.ErrUnauthorized
	}
	account, err := s.accounts.LockAccountByID(ctx, id)
	if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, shared.ErrNotFound) {
		return nil, shared.ErrUnauthorized
	}
	if err != nil {
		return nil, shared.Internal("terms.account", err)
	}
	if account == nil || account.ID != id {
		return nil, shared.ErrUnauthorized
	}
	if !account.IsActive || !account.Role.Valid() || permission != "" && !account.Role.Allows(permission) {
		return nil, shared.ErrForbidden
	}
	return account, nil
}

func (s *Service) Current(ctx context.Context, id uuid.UUID) (v *domainterms.Version, err error) {
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if _, e := s.lockAccount(ctx, id, ""); e != nil {
			return e
		}
		var e error
		v, e = s.repo.LockPublication(ctx, false)
		if e != nil {
			return shared.Internal("terms.current", e)
		}
		if v == nil {
			return domainterms.ErrNotPublished
		}
		return nil
	})
	return
}
func (s *Service) Status(ctx context.Context, id uuid.UUID) (result domainterms.Status, err error) {
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if _, e := s.lockAccount(ctx, id, user.AcceptBorrowerTerms); e != nil {
			return e
		}
		v, e := s.repo.LockPublication(ctx, false)
		if e != nil {
			return shared.Internal("terms.status", e)
		}
		var a *domainterms.Acceptance
		if v != nil {
			a, e = s.repo.FindAcceptance(ctx, id, v.ID)
			if e != nil {
				return shared.Internal("terms.status_acceptance", e)
			}
		}
		previous, e := s.repo.HasAnyAcceptance(ctx, id)
		if e != nil {
			return shared.Internal("terms.status_history", e)
		}
		result = domainterms.StatusFor(v, a, previous)
		return nil
	})
	return
}
func (s *Service) Accept(ctx context.Context, id, versionID uuid.UUID) (receipt *domainterms.Acceptance, err error) {
	if versionID == uuid.Nil {
		return nil, shared.ErrValidation
	}
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if _, e := s.lockAccount(ctx, id, user.AcceptBorrowerTerms); e != nil {
			return e
		}
		v, e := s.repo.LockPublication(ctx, false)
		if e != nil {
			return shared.Internal("terms.accept_current", e)
		}
		if v == nil {
			return domainterms.ErrNotPublished
		}
		if v.ID != versionID {
			old, e := s.repo.FindVersion(ctx, versionID)
			if e != nil {
				return shared.Internal("terms.accept_version", e)
			}
			if old == nil {
				return domainterms.ErrVersionNotFound
			}
			return domainterms.ErrVersionChanged
		}
		receipt, e = s.repo.InsertAcceptance(ctx, id, versionID)
		if e != nil {
			return shared.Internal("terms.accept_insert", e)
		}
		return nil
	})
	return
}

type PublishCommand struct {
	ActorID                  uuid.UUID
	Version, Title, Body     string
	ExpectedCurrentVersionID *uuid.UUID
}

func (s *Service) Publish(ctx context.Context, cmd PublishCommand) (version *domainterms.Version, err error) {
	version, err = domainterms.NewVersion(cmd.Version, cmd.Title, cmd.Body, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	err = s.tx.Within(ctx, func(ctx context.Context) error {
		if _, e := s.lockAccount(ctx, cmd.ActorID, user.PublishTerms); e != nil {
			return e
		}
		current, e := s.repo.LockPublication(ctx, true)
		if e != nil {
			return shared.Internal("terms.publish_current", e)
		}
		if current == nil && cmd.ExpectedCurrentVersionID != nil || current != nil && (cmd.ExpectedCurrentVersionID == nil || *cmd.ExpectedCurrentVersionID != current.ID) {
			return domainterms.ErrPublicationChanged
		}
		if e = s.repo.Publish(ctx, version); errors.Is(e, domainterms.ErrVersionExists) {
			return e
		} else if e != nil {
			return shared.Internal("terms.publish", e)
		}
		return nil
	})
	return
}

// RequireCurrentAcceptance must run inside the borrowing command's transaction
// when authorizing a write. It retains account/publication locks until commit.
// A read-only eligibility preview may use it in its own transaction, but cannot
// authorize a later write: that command must repeat this check atomically.
// A direct checkout passes its separately authorized target Borrower ID.
func (s *Service) RequireCurrentAcceptance(ctx context.Context, borrowerID uuid.UUID) (*domainterms.Acceptance, error) {
	if _, err := s.lockAccount(ctx, borrowerID, user.AcceptBorrowerTerms); err != nil {
		return nil, err
	}
	v, err := s.repo.LockPublication(ctx, false)
	if err != nil {
		return nil, shared.Internal("terms.policy_current", err)
	}
	if v == nil {
		return nil, domainterms.ErrNotPublished
	}
	a, err := s.repo.FindAcceptance(ctx, borrowerID, v.ID)
	if err != nil {
		return nil, shared.Internal("terms.policy_acceptance", err)
	}
	if a == nil {
		return nil, domainterms.ErrAcceptanceRequired
	}
	return a, nil
}
