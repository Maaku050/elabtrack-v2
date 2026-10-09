package borrowing

import (
	"context"
	"encoding/json"
	"errors"
	a "github.com/Maaku050/elabtrack-v2/backend/internal/domain/accounts"
	d "github.com/Maaku050/elabtrack-v2/backend/internal/domain/borrowing"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
	"testing"
)

type replayAccounts struct {
	a.Repository
	records map[uuid.UUID]a.Record
}

func (r replayAccounts) LockAccounts(context.Context, []uuid.UUID) error { return nil }
func (r replayAccounts) Get(_ context.Context, id uuid.UUID) (a.Record, error) {
	v, ok := r.records[id]
	if !ok {
		return v, shared.ErrNotFound
	}
	return v, nil
}

type replayRepository struct {
	d.Repository
	data []byte
}

func (r replayRepository) ReadReceipt(context.Context, uuid.UUID, string, string, string) ([]byte, error) {
	return r.data, nil
}

type unitTx struct{}

func (unitTx) Within(c context.Context, fn func(context.Context) error) error { return fn(c) }
func TestReplayReauthorizesBeforeReturningReceipt(t *testing.T) {
	actor, target := uuid.New(), uuid.New()
	record := d.Record{ID: uuid.New(), Status: "CHECKED_OUT"}
	raw, _ := json.Marshal(record)
	for _, tc := range []struct {
		name               string
		active, activation bool
		role               user.Role
		targetActive       bool
		want               error
	}{{"allowed", true, false, user.RoleStaff, true, nil}, {"deactivated", false, false, user.RoleStaff, true, shared.ErrUnauthorized}, {"activation", true, true, user.RoleStaff, true, shared.ErrUnauthorized}, {"role revoked", true, false, user.RoleBorrower, true, shared.ErrForbidden}, {"target inactive", true, false, user.RoleStaff, false, d.ErrEligibility}} {
		t.Run(tc.name, func(t *testing.T) {
			ar := replayAccounts{records: map[uuid.UUID]a.Record{actor: {ID: actor, Role: tc.role, IsActive: tc.active, ActivationRequired: tc.activation}, target: {ID: target, Role: user.RoleBorrower, IsActive: tc.targetActive, BorrowerType: "FACULTY"}}}
			s := NewService(replayRepository{data: raw}, ar, nil, nil, unitTx{})
			out, e := s.command(context.Background(), actor, target, "direct", uuid.NewString(), nil, true, true, func(context.Context) (d.Record, error) {
				t.Fatal("receipt replay invoked write path")
				return d.Record{}, nil
			})
			if !errors.Is(e, tc.want) || tc.want == nil && out.ID != record.ID {
				t.Fatal("replay authority", e)
			}
		})
	}
}
func TestInvalidCommandsNeverStartTransaction(t *testing.T) {
	s := NewService(nil, nil, nil, nil, nil)
	for _, in := range []d.Input{{}, {Confirm: true}, {Confirm: true, Items: []d.Line{{EquipmentID: uuid.New(), Quantity: -1}}}} {
		if _, e := s.Submit(context.Background(), uuid.New(), uuid.NewString(), in); !errors.Is(e, shared.ErrInvalidInput) {
			t.Fatal(e)
		}
	}
	if _, e := s.Direct(context.Background(), uuid.New(), uuid.NewString(), d.Input{}); !errors.Is(e, shared.ErrInvalidInput) {
		t.Fatal(e)
	}
}
