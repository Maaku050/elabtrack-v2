package auth

import (
	"context"
	"errors"

	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

// Refresh returns credentials only after a single-use rotation commits.
func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (TokenPairDTO, error) {
	hash, err := s.tokenHasher.Hash(req.RefreshToken)
	if err != nil {
		return TokenPairDTO{}, domainauth.ErrTokenInvalid
	}
	var pair TokenPairDTO
	err = s.tx.Within(ctx, func(txCtx context.Context) error {
		stored, err := s.tokens.FindByHashForUpdate(txCtx, hash)
		if errors.Is(err, domainauth.ErrTokenNotFound) {
			return domainauth.ErrTokenInvalid
		}
		if err != nil {
			return shared.ErrInternal
		}
		if stored == nil || stored.TokenHash != hash || stored.UserID == uuid.Nil || !stored.ValidAt(s.now()) {
			return domainauth.ErrTokenInvalid
		}
		account, err := s.accounts.LockAccountByID(txCtx, stored.UserID)
		if errors.Is(err, domainuser.ErrUserNotFound) {
			return domainauth.ErrTokenInvalid
		}
		if err != nil {
			return shared.ErrInternal
		}
		if account == nil || account.ID != stored.UserID || !account.IsActive || !account.Role.Valid() {
			return domainauth.ErrTokenInvalid
		}
		// Waiting for locks must not extend session validity.
		if !stored.ValidAt(s.now()) {
			return domainauth.ErrTokenInvalid
		}
		u := &domainuser.User{ID: account.ID, Email: account.Email, Name: account.Name, Role: account.Role, IsActive: account.IsActive}
		candidate, replacement, err := s.newTokenPair(txCtx, u)
		if err != nil {
			return shared.ErrInternal
		}
		// Insert first to satisfy the immediate replacement FK. Both writes are
		// invisible outside this transaction until consumption and commit succeed.
		if err := s.tokens.Create(txCtx, replacement); err != nil {
			return shared.ErrInternal
		}
		if err := s.tokens.Consume(txCtx, hash, replacement.ID, s.now()); err != nil {
			if errors.Is(err, domainauth.ErrTokenInvalid) {
				return domainauth.ErrTokenInvalid
			}
			return shared.ErrInternal
		}
		pair = candidate
		return nil
	})
	if err != nil {
		if errors.Is(err, domainauth.ErrTokenInvalid) {
			return TokenPairDTO{}, domainauth.ErrTokenInvalid
		}
		return TokenPairDTO{}, shared.ErrInternal
	}
	return pair, nil
}

// Logout revokes only the presented session, with no token-existence detail.
func (s *Service) Logout(ctx context.Context, req RefreshRequest) error {
	hash, err := s.tokenHasher.Hash(req.RefreshToken)
	if err != nil {
		return nil
	}
	if err := s.tokens.Revoke(ctx, hash); err != nil {
		return shared.ErrInternal
	}
	return nil
}

func (s *Service) issueTokenPair(ctx context.Context, u *domainuser.User) (TokenPairDTO, error) {
	pair, session, err := s.newTokenPair(ctx, u)
	if err != nil {
		return TokenPairDTO{}, shared.ErrInternal
	}
	if err := s.tokens.Create(ctx, session); err != nil {
		return TokenPairDTO{}, shared.ErrInternal
	}
	return pair, nil
}

// Raw credentials exist only in the transport result. The independent
// persistence entity holds the digest, never the raw secret.
func (s *Service) newTokenPair(ctx context.Context, u *domainuser.User) (TokenPairDTO, *domainauth.RefreshToken, error) {
	access, err := s.issuer.IssueAccessToken(ctx, u)
	if err != nil {
		return TokenPairDTO{}, nil, shared.ErrInternal
	}
	raw, err := s.issuer.GenerateRefreshToken()
	if err != nil {
		return TokenPairDTO{}, nil, shared.ErrInternal
	}
	hash, err := s.tokenHasher.Hash(raw)
	if err != nil {
		return TokenPairDTO{}, nil, shared.ErrInternal
	}
	now := s.now()
	session := domainauth.NewRefreshToken(hash, u.ID, now, now.Add(s.refreshTTL))
	return TokenPairDTO{AccessToken: access, RefreshToken: raw, ExpiresAt: now.Add(s.accessTTL), TokenType: "Bearer", RefreshExpiresAt: session.ExpiresAt, User: AuthUserDTO{ID: u.ID, Email: u.Email, Name: u.Name, Role: string(u.Role), IsActive: u.IsActive}}, session, nil
}
