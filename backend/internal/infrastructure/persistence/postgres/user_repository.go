package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository implements domain/user.Repository on top of PostgreSQL.
type UserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository constructs a UserRepository backed by the given pool.
func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// executor returns the pgx.Tx from the context if one is present, otherwise
// the pool. This lets the repository participate in transactions managed by
// database.TxManager.
func (r *UserRepository) executor(ctx context.Context) executor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}
	return r.pool
}

// executor is the minimal subset of pgxpool.Pool / pgx.Tx we use.
type executor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

const (
	userColumns = `id, email, name, password, role, is_active, created_at, updated_at`
)

// Create inserts a new user. Email uniqueness is enforced by a unique index;
// a duplicate returns domain/user.ErrEmailAlreadyExists.
func (r *UserRepository) Create(ctx context.Context, u *domainuser.User) error {
	const q = `
		INSERT INTO users (id, email, name, password, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.executor(ctx).Exec(ctx, q,
		u.ID, u.Email, u.Name, u.Password, string(u.Role), u.IsActive, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err, "users_email_key") {
			return domainuser.ErrEmailAlreadyExists
		}
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

// FindByID returns a single user by id or domain/user.ErrUserNotFound.
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	const q = `SELECT ` + userColumns + ` FROM users WHERE id = $1`
	row := r.executor(ctx).QueryRow(ctx, q, id)
	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainuser.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

// accountColumns intentionally excludes password and session data.
const accountColumns = `id, email, name, role, is_active, created_at, updated_at`

// FindAccountByID reads current authorization state without loading secrets.
func (r *UserRepository) FindAccountByID(ctx context.Context, id uuid.UUID) (*domainuser.Account, error) {
	const q = `SELECT ` + accountColumns + ` FROM users WHERE id = $1`
	account, err := scanAccount(r.executor(ctx).QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainuser.ErrUserNotFound
	}
	return account, err
}

// LockAccountByID prevents concurrent status/role updates until rotation ends.
// FOR SHARE also conflicts with non-key updates; KEY SHARE would not suffice.
func (r *UserRepository) LockAccountByID(ctx context.Context, id uuid.UUID) (*domainuser.Account, error) {
	tx, ok := database.TxFromContext(ctx)
	if !ok {
		return nil, shared.ErrInternal
	}
	const q = `SELECT ` + accountColumns + ` FROM users WHERE id = $1 FOR SHARE`
	account, err := scanAccount(tx.QueryRow(ctx, q, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domainuser.ErrUserNotFound
	}
	if err != nil {
		return nil, shared.Internal("account_store.lock", err)
	}
	return account, nil
}

// UpdateProfile cannot overwrite a concurrent demotion/deactivation/password
// change. Permission is rechecked in the same statement as the name mutation.
func (r *UserRepository) UpdateProfile(ctx context.Context, id uuid.UUID, name string) (*domainuser.Account, error) {
	const q = `UPDATE users SET name = $2, updated_at = NOW()
  WHERE id = $1 AND is_active = TRUE AND role IN ('user','admin')
  RETURNING ` + accountColumns
	account, err := scanAccount(r.executor(ctx).QueryRow(ctx, q, id, name))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, shared.ErrForbidden
	}
	return account, err
}

func scanAccount(s scanner) (*domainuser.Account, error) {
	a := &domainuser.Account{}
	var role string
	if err := s.Scan(&a.ID, &a.Email, &a.Name, &role, &a.IsActive, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return nil, err
	}
	a.Role = domainuser.Role(role)
	return a, nil
}

// FindByEmail returns a single user by email or domain/user.ErrUserNotFound.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	const q = `SELECT ` + userColumns + ` FROM users WHERE email = $1`
	row := r.executor(ctx).QueryRow(ctx, q, strings.ToLower(email))
	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainuser.ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

// Update updates an existing user identified by ID.
func (r *UserRepository) Update(ctx context.Context, u *domainuser.User) error {
	const q = `
		UPDATE users
		SET email = $2, name = $3, password = $4, role = $5, is_active = $6, updated_at = $7
		WHERE id = $1`
	tag, err := r.executor(ctx).Exec(ctx, q,
		u.ID, u.Email, u.Name, u.Password, string(u.Role), u.IsActive, u.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err, "users_email_key") {
			return domainuser.ErrEmailAlreadyExists
		}
		return fmt.Errorf("update user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainuser.ErrUserNotFound
	}
	return nil
}

// List returns a paginated list of users and the total count.
func (r *UserRepository) List(ctx context.Context, p shared.Page) ([]*domainuser.User, int, error) {
	// Whitelist sort columns to prevent SQL injection via ORDER BY.
	if !isAllowedSort(p.Sort) {
		p.Sort = "created_at"
	}

	var (
		rows  []*domainuser.User
		total int
	)

	// Count total matches (without offset/limit).
	countQ := `SELECT COUNT(*) FROM users`
	listQ := `SELECT ` + userColumns + ` FROM users`
	args := []any{}
	if p.Search != "" {
		countQ += ` WHERE name ILIKE $1 OR email ILIKE $1`
		listQ += ` WHERE name ILIKE $1 OR email ILIKE $1`
		args = append(args, "%"+p.Search+"%")
	}

	if err := r.executor(ctx).QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	listQ += fmt.Sprintf(` ORDER BY %s %s LIMIT $%d OFFSET $%d`, p.Sort, p.Order, len(args)+1, len(args)+2)
	args = append(args, p.PerPage, p.Offset())

	rs, err := r.executor(ctx).Query(ctx, listQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rs.Close()

	for rs.Next() {
		u, err := scanUser(rs)
		if err != nil {
			return nil, 0, err
		}
		rows = append(rows, u)
	}
	if err := rs.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate users: %w", err)
	}
	return rows, total, nil
}

// Delete removes a user by id.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM users WHERE id = $1`
	tag, err := r.executor(ctx).Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domainuser.ErrUserNotFound
	}
	return nil
}

// scanner is satisfied by both pgx.Row and pgx.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanUser(s scanner) (*domainuser.User, error) {
	u := &domainuser.User{}
	var role string
	err := s.Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.Password,
		&role,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.Role = domainuser.Role(role)
	return u, nil
}

var allowedSorts = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"email":      true,
	"name":       true,
}

func isAllowedSort(s string) bool { return allowedSorts[s] }

// isUniqueViolation reports whether err is a Postgres unique_violation for
// the given constraint name.
func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" && (constraint == "" || pgErr.ConstraintName == constraint)
	}
	return false
}
