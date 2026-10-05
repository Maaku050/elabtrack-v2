package application

import (
	"context"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/google/uuid"
)

// Ports are outbound interfaces the application layer depends on.
// Infrastructure provides concrete adapters. This keeps the application
// layer decoupled from bcrypt, JWT libraries, and PostgreSQL.

// PasswordHasher hashes and verifies plaintext passwords.
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hashed, plain string) error
}

// TokenIssuer issues and verifies JWT access tokens and generates opaque
// refresh-token strings.
type TokenIssuer interface {
	IssueAccessToken(ctx context.Context, u *user.User) (string, error)
	VerifyAccessToken(ctx context.Context, token string) (Claims, error)
	GenerateRefreshToken() (string, error)
}

// Claims describes the data extracted from a verified access token.
type Claims struct {
	UserID uuid.UUID
	Email  string
	Role   string
}

// Ports bundle groups outbound dependencies for a use case.
type UserRepo interface {
	user.Repository
}

type AuthRepo interface {
	auth.Repository
}
