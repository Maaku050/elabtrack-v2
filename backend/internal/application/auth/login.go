package auth

import (
	"context"
	"errors"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
)

// Login verifies credentials and issues a token pair.
// Failed credentials do not reveal whether an account exists or is inactive.
// Account status is checked only after the password has been verified.
func (s *Service) Login(ctx context.Context, req LoginRequest) (TokenPairDTO, error) {
	email, err := domainuser.ParseEmail(req.Email)
	if err != nil {
		return TokenPairDTO{}, domainuser.ErrInvalidCredentials
	}

	u, err := s.users.FindByEmail(ctx, email.String())
	if err != nil {
		if errors.Is(err, domainuser.ErrUserNotFound) || errors.Is(err, shared.ErrNotFound) {
			return TokenPairDTO{}, domainuser.ErrInvalidCredentials
		}
		return TokenPairDTO{}, shared.Internal("auth.login_lookup", err)
	}
	if u == nil {
		return TokenPairDTO{}, domainuser.ErrInvalidCredentials
	}
	if err := s.hasher.Compare(u.Password, req.Password); err != nil {
		return TokenPairDTO{}, domainuser.ErrInvalidCredentials
	}
	if !u.IsActive || !u.Role.Valid() {
		return TokenPairDTO{}, domainuser.ErrUserInactive
	}

	return s.issueTokenPair(ctx, u)
}
