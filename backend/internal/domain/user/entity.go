package user

import (
	"errors"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/google/uuid"
)

// Role is the current FSMO product capability, independent of borrower category.
type Role string

const (
	RoleBorrower Role = "BORROWER"
	RoleStaff    Role = "STAFF"
	RoleAdmin    Role = "ADMIN"
)

// User is the aggregate root for the user domain.
type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	Password  string // hashed
	Role      Role
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser constructs a User with sensible defaults. The caller is responsible
// for hashing the password before calling this constructor.
func NewUser(email, name, hashedPassword string) *User {
	now := time.Now().UTC()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Name:      name,
		Password:  hashedPassword,
		Role:      RoleBorrower,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsAdmin reports whether the user has the admin role.
func (u *User) IsAdmin() bool { return u.Role == RoleAdmin }

// Domain errors. Use errors.Is to check.
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidName        = errors.New("invalid name")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
)

// Ensure shared package is referenced so importers can compose errors.
var _ = shared.ErrNotFound
