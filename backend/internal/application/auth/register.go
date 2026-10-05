package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/fullstacktemplate/backend/internal/domain/shared"
	domainuser "github.com/fullstacktemplate/backend/internal/domain/user"
)

// Register creates a new user and immediately issues a token pair.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (TokenPairDTO, error) {
	email, err := domainuser.ParseEmail(req.Email)
	if err != nil {
		return TokenPairDTO{}, err
	}
	name, err := domainuser.ParseName(req.Name)
	if err != nil {
		return TokenPairDTO{}, err
	}
	if err := (domainuser.PlainPassword(req.Password)).Validate(); err != nil {
		return TokenPairDTO{}, err
	}

	if existing, err := s.users.FindByEmail(ctx, email.String()); err == nil && existing != nil {
		return TokenPairDTO{}, domainuser.ErrEmailAlreadyExists
	} else if err != nil && !errors.Is(err, domainuser.ErrUserNotFound) && !errors.Is(err, shared.ErrNotFound) {
		return TokenPairDTO{}, fmt.Errorf("check existing email: %w", err)
	}

	hashed, err := s.hasher.Hash(req.Password)
	if err != nil {
		return TokenPairDTO{}, fmt.Errorf("hash password: %w", err)
	}

	u := domainuser.NewUser(email.String(), name.String(), hashed)
	if err := s.users.Create(ctx, u); err != nil {
		return TokenPairDTO{}, fmt.Errorf("create user: %w", err)
	}

	return s.issueTokenPair(ctx, u)
}
