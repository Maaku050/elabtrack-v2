package user

import (
	"context"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/google/uuid"
)

// Repository is the persistence port for the user domain.
// Implementations live in infrastructure/persistence/postgres.
type Repository interface {
	Create(ctx context.Context, u *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, u *User) error
	// UpdateProfile changes only the display name and rechecks account access.
	UpdateProfile(ctx context.Context, id uuid.UUID, name string) (*Account, error)
	List(ctx context.Context, p shared.Page) ([]*User, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
