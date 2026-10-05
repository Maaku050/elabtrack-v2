package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
)

// Refresh rotates a refresh token and issues a new token pair.
// The old refresh token is revoked before the new one is issued (rotation).
func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (TokenPairDTO, error) {
	stored, err := s.tokens.FindByToken(ctx, req.RefreshToken)
	if err != nil {
		return TokenPairDTO{}, domainauth.ErrTokenInvalid
	}
	if !stored.IsValid() {
		if !stored.IsRevoked() {
			_ = s.tokens.Revoke(ctx, stored.Token)
		}
		return TokenPairDTO{}, domainauth.ErrTokenInvalid
	}

	if err := s.tokens.Revoke(ctx, stored.Token); err != nil {
		return TokenPairDTO{}, fmt.Errorf("revoke old refresh token: %w", err)
	}

	u, err := s.users.FindByID(ctx, stored.UserID)
	if err != nil {
		return TokenPairDTO{}, domainauth.ErrTokenInvalid
	}
	if !u.IsActive {
		return TokenPairDTO{}, domainuser.ErrUserInactive
	}

	return s.issueTokenPair(ctx, u)
}

// Logout revokes the supplied refresh token. It is idempotent.
func (s *Service) Logout(ctx context.Context, req RefreshRequest) error {
	if err := s.tokens.Revoke(ctx, req.RefreshToken); err != nil {
		if errors.Is(err, domainauth.ErrTokenNotFound) {
			return nil
		}
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

// issueTokenPair is shared by Register, Login, and Refresh.
func (s *Service) issueTokenPair(ctx context.Context, u *domainuser.User) (TokenPairDTO, error) {
	access, err := s.issuer.IssueAccessToken(ctx, u)
	if err != nil {
		return TokenPairDTO{}, fmt.Errorf("issue access token: %w", err)
	}
	refreshStr, err := s.issuer.GenerateRefreshToken()
	if err != nil {
		return TokenPairDTO{}, fmt.Errorf("generate refresh token: %w", err)
	}
	now := timeNowUTC()
	rt := domainauth.NewRefreshToken(refreshStr, u.ID, now.Add(s.refreshTTL))
	if err := s.tokens.Create(ctx, rt); err != nil {
		return TokenPairDTO{}, fmt.Errorf("persist refresh token: %w", err)
	}
	return TokenPairDTO{
		AccessToken:  access,
		RefreshToken: refreshStr,
		ExpiresAt:    now.Add(s.accessTTL),
		TokenType:    "Bearer",
	}, nil
}

// timeNowUTC is a small seam for tests; defaults to time.Now().UTC().
var timeNowUTC = func() time.Time { return time.Now().UTC() }

// Ensure application.Claims is referenced (used by Service.CurrentUser).
var _ application.Claims
