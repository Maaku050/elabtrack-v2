package auth

import (
	"context"
	"errors"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

// Identity carries only the account locator from verified access claims.
// Token email/role are intentionally excluded from authorization evidence.
type Identity struct{ UserID uuid.UUID }

// Principal contains only the current account's safe fields, never a User
// aggregate, password hash, access token or refresh-token record.
type Principal domainuser.Account

type AccountResolver interface {
	ResolveCurrentAccount(ctx context.Context, identity Identity) (Principal, error)
}

type CurrentAccountService struct{ accounts domainuser.AccountRepository }

func NewAccountResolver(accounts domainuser.AccountRepository) *CurrentAccountService {
	return &CurrentAccountService{accounts: accounts}
}

func (s *CurrentAccountService) ResolveCurrentAccount(ctx context.Context, identity Identity) (Principal, error) {
	if identity.UserID == uuid.Nil {
		return Principal{}, shared.ErrUnauthorized
	}
	account, err := s.accounts.FindAccountByID(ctx, identity.UserID)
	if errors.Is(err, domainuser.ErrUserNotFound) || errors.Is(err, shared.ErrNotFound) {
		return Principal{}, shared.ErrUnauthorized
	}
	if err != nil {
		return Principal{}, shared.ErrInternal
	}
	if account == nil || account.ID != identity.UserID {
		return Principal{}, shared.ErrUnauthorized
	}
	if !account.IsActive || !account.Role.Valid() {
		return Principal{}, shared.ErrForbidden
	}
	return Principal(*account), nil
}

func (p Principal) UserDTO() AuthUserDTO {
	return AuthUserDTO{ID: p.ID, Email: p.Email, Name: p.Name, Role: string(p.Role), IsActive: p.IsActive}
}
