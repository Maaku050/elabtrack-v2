package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

type accountReader struct {
	account *domainuser.Account
	err     error
	calls   int
}

func (r *accountReader) FindAccountByID(_ context.Context, _ uuid.UUID) (*domainuser.Account, error) {
	r.calls++
	return r.account, r.err
}
func TestResolverRejectsMissingMismatchedAndUnavailableAccounts(t *testing.T) {
	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	for _, tt := range []struct {
		name          string
		account       *domainuser.Account
		err, expected error
	}{
		{"nil result", nil, nil, shared.ErrUnauthorized},
		{"wrong identity", &domainuser.Account{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), IsActive: true, Role: domainuser.RoleAdmin}, nil, shared.ErrUnauthorized},
		{"shared not found", nil, shared.ErrNotFound, shared.ErrUnauthorized},
		{"database error", nil, errors.New("database-private-sentinel"), shared.ErrInternal},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := &accountReader{account: tt.account, err: tt.err}
			principal, err := NewAccountResolver(r).ResolveCurrentAccount(context.Background(), Identity{UserID: id})
			if !errors.Is(err, tt.expected) || principal.ID != uuid.Nil || r.calls != 1 {
				t.Fatal("resolver did not fail closed")
			}
		})
	}
	r := &accountReader{}
	if _, err := NewAccountResolver(r).ResolveCurrentAccount(context.Background(), Identity{}); !errors.Is(err, shared.ErrUnauthorized) || r.calls != 0 {
		t.Fatal("zero identity must not reach account repository")
	}
}
