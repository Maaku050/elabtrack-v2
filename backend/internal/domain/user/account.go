package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Account is a client-safe snapshot. It deliberately has no password or
// session fields and is separate from the persistence aggregate.
type Account struct {
	ID        uuid.UUID
	Email     string
	Name      string
	Role      Role
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AccountRepository is the narrow read port used at the request boundary.
type AccountRepository interface {
	FindAccountByID(ctx context.Context, id uuid.UUID) (*Account, error)
}

// Valid fails closed for unknown, legacy and category values after reconciliation.
func (r Role) Valid() bool { return r == RoleBorrower || r == RoleStaff || r == RoleAdmin }
