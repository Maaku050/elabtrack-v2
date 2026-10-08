package terms

import (
	"context"
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainterms "github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
	"testing"
	"time"
)

type unitTxKey struct{}
type unitTx struct{}

func (unitTx) Within(ctx context.Context, fn func(context.Context) error) error {
	return fn(context.WithValue(ctx, unitTxKey{}, true))
}

type unitAccounts struct {
	account *user.Account
	err     error
}

func (a unitAccounts) LockAccountByID(ctx context.Context, id uuid.UUID) (*user.Account, error) {
	if ctx.Value(unitTxKey{}) != true {
		return nil, shared.ErrInternal
	}
	return a.account, a.err
}

type unitRepository struct {
	current   *domainterms.Version
	old       *domainterms.Version
	receipt   *domainterms.Acceptance
	previous  bool
	err       error
	inserts   int
	publishes int
}

func (r *unitRepository) LockPublication(ctx context.Context, exclusive bool) (*domainterms.Version, error) {
	if ctx.Value(unitTxKey{}) != true {
		return nil, shared.ErrInternal
	}
	return r.current, r.err
}
func (r *unitRepository) FindVersion(context.Context, uuid.UUID) (*domainterms.Version, error) {
	return r.old, r.err
}
func (r *unitRepository) FindAcceptance(context.Context, uuid.UUID, uuid.UUID) (*domainterms.Acceptance, error) {
	return r.receipt, r.err
}
func (r *unitRepository) HasAnyAcceptance(context.Context, uuid.UUID) (bool, error) {
	return r.previous, r.err
}
func (r *unitRepository) InsertAcceptance(_ context.Context, u, v uuid.UUID) (*domainterms.Acceptance, error) {
	r.inserts++
	if r.receipt == nil {
		r.receipt = &domainterms.Acceptance{ID: uuid.New(), UserID: u, TermsVersionID: v, AcceptedAt: time.Now().UTC()}
	}
	return r.receipt, r.err
}
func (r *unitRepository) Publish(_ context.Context, v *domainterms.Version) error {
	r.publishes++
	r.current = v
	return r.err
}
func unitService(role user.Role, active bool) (*Service, *unitRepository, uuid.UUID) {
	id := uuid.New()
	r := &unitRepository{current: &domainterms.Version{ID: uuid.New()}}
	return NewService(r, unitAccounts{account: &user.Account{ID: id, Role: role, IsActive: active}}, unitTx{}), r, id
}
func TestAcceptanceChecksIdentityStatusRoleAndExactVersion(t *testing.T) {
	ctx := context.Background()
	for _, role := range []user.Role{user.RoleStaff, user.RoleAdmin, "Student", "unknown"} {
		s, r, id := unitService(role, true)
		if _, err := s.Accept(ctx, id, r.current.ID); !errors.Is(err, shared.ErrForbidden) || r.inserts != 0 {
			t.Fatal("nonborrower acceptance")
		}
	}
	s, r, id := unitService(user.RoleBorrower, false)
	if _, err := s.Accept(ctx, id, r.current.ID); !errors.Is(err, shared.ErrForbidden) || r.inserts != 0 {
		t.Fatal("inactive acceptance")
	}
	s, r, id = unitService(user.RoleBorrower, true)
	if _, err := s.Accept(ctx, uuid.Nil, r.current.ID); !errors.Is(err, shared.ErrUnauthorized) {
		t.Fatal("anonymous acceptance")
	}
	if _, err := s.Accept(ctx, id, uuid.New()); !errors.Is(err, domainterms.ErrVersionNotFound) || r.inserts != 0 {
		t.Fatal("unknown version")
	}
	r.old = &domainterms.Version{ID: uuid.New()}
	if _, err := s.Accept(ctx, id, r.old.ID); !errors.Is(err, domainterms.ErrVersionChanged) || r.inserts != 0 {
		t.Fatal("stale version")
	}
	first, err := s.Accept(ctx, id, r.current.ID)
	if err != nil || first.UserID != id || first.TermsVersionID != r.current.ID {
		t.Fatal("wrong evidence identity")
	}
	again, err := s.Accept(ctx, id, r.current.ID)
	if err != nil || first.ID != again.ID || !first.AcceptedAt.Equal(again.AcceptedAt) {
		t.Fatal("repeat rewrote evidence")
	}
	r.current = nil
	if _, err := s.Accept(ctx, id, uuid.New()); !errors.Is(err, domainterms.ErrNotPublished) {
		t.Fatal("missing terms consent")
	}
}
func TestPolicyRequiresCallerTransactionAndDurableCurrentAcceptance(t *testing.T) {
	s, r, id := unitService(user.RoleBorrower, true)
	ctx := context.Background()
	if _, err := s.RequireCurrentAcceptance(ctx, id); !errors.Is(err, shared.ErrInternal) {
		t.Fatal("detached policy preflight allowed")
	}
	tx := unitTx{}
	if err := tx.Within(ctx, func(ctx context.Context) error { _, e := s.RequireCurrentAcceptance(ctx, id); return e }); !errors.Is(err, domainterms.ErrAcceptanceRequired) {
		t.Fatal("no evidence allowed")
	}
	r.receipt = &domainterms.Acceptance{ID: uuid.New(), UserID: id, TermsVersionID: r.current.ID, AcceptedAt: time.Now().UTC()}
	if err := tx.Within(ctx, func(ctx context.Context) error { _, e := s.RequireCurrentAcceptance(ctx, id); return e }); err != nil {
		t.Fatal("valid policy denied")
	}
	r.current = nil
	if err := tx.Within(ctx, func(ctx context.Context) error { _, e := s.RequireCurrentAcceptance(ctx, id); return e }); !errors.Is(err, domainterms.ErrNotPublished) {
		t.Fatal("missing policy allowed")
	}
}
func TestPublicationAuthorityAndOptimisticCurrentBoundary(t *testing.T) {
	for _, role := range []user.Role{user.RoleBorrower, user.RoleStaff} {
		s, r, id := unitService(role, true)
		_, err := s.Publish(context.Background(), PublishCommand{ActorID: id, Version: "test-only", Title: "Synthetic", Body: "TEST ONLY"})
		if !errors.Is(err, shared.ErrForbidden) || r.publishes != 0 {
			t.Fatal("unauthorized publication")
		}
	}
	s, r, id := unitService(user.RoleAdmin, true)
	cmd := PublishCommand{ActorID: id, Version: "test-only", Title: "Synthetic", Body: "TEST ONLY"}
	if _, err := s.Publish(context.Background(), cmd); !errors.Is(err, domainterms.ErrPublicationChanged) || r.publishes != 0 {
		t.Fatal("blind publication overwrite")
	}
	current := r.current.ID
	cmd.ExpectedCurrentVersionID = &current
	if _, err := s.Publish(context.Background(), cmd); err != nil || r.publishes != 1 {
		t.Fatal("authorized publication denied")
	}
	r.err = errors.New("private database sentinel")
	_, err := s.Status(context.Background(), id)
	if !errors.Is(err, shared.ErrForbidden) {
		t.Fatal("admin borrower status should be denied")
	}
	s, r, id = unitService(user.RoleBorrower, true)
	r.err = errors.New("private database sentinel")
	_, err = s.Status(context.Background(), id)
	if !errors.Is(err, shared.ErrInternal) || err.Error() == r.err.Error() {
		t.Fatal("raw database error escaped")
	}
}
