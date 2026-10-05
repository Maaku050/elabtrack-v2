package auth

import (
	"context"

	domainuser "github.com/fullstacktemplate/backend/internal/domain/user"
)

// Login verifies credentials and issues a token pair.
// On any failure path it returns ErrInvalidCredentials to avoid leaking
// whether the email exists in the system.
func (s *Service) Login(ctx context.Context, req LoginRequest) (TokenPairDTO, error) {
	email, err := domainuser.ParseEmail(req.Email)
	if err != nil {
		return TokenPairDTO{}, domainuser.ErrInvalidCredentials
	}

	u, err := s.users.FindByEmail(ctx, email.String())
	if err != nil {
		return TokenPairDTO{}, domainuser.ErrInvalidCredentials
	}
	if !u.IsActive {
		return TokenPairDTO{}, domainuser.ErrUserInactive
	}
	if err := s.hasher.Compare(u.Password, req.Password); err != nil {
		return TokenPairDTO{}, domainuser.ErrInvalidCredentials
	}

	return s.issueTokenPair(ctx, u)
}
