package postgres

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	domainauth "github.com/Maaku050/elabtrack-v2/backend/internal/domain/auth"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type sessionTx struct {
	pgx.Tx
	query    string
	args     []any
	values   []any
	err      error
	affected int64
}

func (tx *sessionTx) Exec(_ context.Context, q string, args ...any) (pgconn.CommandTag, error) {
	tx.query = q
	tx.args = args
	return pgconn.NewCommandTag("UPDATE " + strconv.FormatInt(tx.affected, 10)), tx.err
}

func (tx *sessionTx) QueryRow(_ context.Context, q string, args ...any) pgx.Row {
	tx.query = q
	tx.args = args
	return sessionRow{tx.values, tx.err}
}

type sessionRow struct {
	values []any
	err    error
}

func (row sessionRow) Scan(dest ...any) error {
	if row.err != nil {
		return row.err
	}
	if len(dest) != len(row.values) {
		return errors.New("wrong session projection")
	}
	for i, value := range row.values {
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(value))
	}
	return nil
}
func sessionFixture() (*domainauth.RefreshToken, *sessionTx, context.Context) {
	now := time.Date(2026, 10, 6, 0, 0, 0, 0, time.UTC)
	s := domainauth.NewRefreshToken(strings.Repeat("b", 64), uuid.New(), now, now.Add(time.Hour))
	tx := &sessionTx{affected: 1, values: []any{s.ID, s.TokenHash, s.UserID, s.ExpiresAt, (*time.Time)(nil), (*uuid.UUID)(nil), s.CreatedAt, s.UpdatedAt}}
	return s, tx, database.WithTx(context.Background(), tx)
}
func TestSessionRepositoryHashOnlyAndLocking(t *testing.T) {
	session, tx, ctx := sessionFixture()
	repo := NewAuthRepository(nil)
	if err := repo.Create(ctx, session); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(tx.query, "token_hash") || strings.Contains(tx.query, "(id, token,") || tx.args[1] != session.TokenHash || len(tx.args) != 8 {
		t.Fatal("create not hash-only")
	}
	got, err := repo.FindByHashForUpdate(ctx, session.TokenHash)
	if err != nil || got.ID != session.ID || got.TokenHash != session.TokenHash {
		t.Fatal("session scan incorrect")
	}
	if !strings.Contains(tx.query, "token_hash = $1 FOR UPDATE") || tx.args[0] != session.TokenHash || strings.Contains(tx.query, session.TokenHash) {
		t.Fatal("parameterized locking lookup missing")
	}
	if _, err := repo.FindByHashForUpdate(context.Background(), session.TokenHash); !errors.Is(err, shared.ErrInternal) {
		t.Fatal("autocommit lookup would release lock")
	}
	tx.err = pgx.ErrNoRows
	if _, err := repo.FindByHashForUpdate(ctx, session.TokenHash); !errors.Is(err, domainauth.ErrTokenNotFound) {
		t.Fatal("missing record mapping")
	}
}
func TestConsumeConditionalAndTransactional(t *testing.T) {
	session, tx, ctx := sessionFixture()
	repo := NewAuthRepository(nil)
	replacement := uuid.New()
	if err := repo.Consume(ctx, session.TokenHash, replacement, session.CreatedAt); err != nil {
		t.Fatal(err)
	}
	for _, predicate := range []string{"token_hash = $1", "revoked_at IS NULL", "expires_at > $3", "expires_at > clock_timestamp()", "replaced_by = $2"} {
		if !strings.Contains(tx.query, predicate) {
			t.Fatal("consume missing predicate", predicate)
		}
	}
	if tx.args[0] != session.TokenHash || tx.args[1] != replacement || tx.args[2] != session.CreatedAt {
		t.Fatal("wrong consume parameters")
	}
	tx.affected = 0
	if err := repo.Consume(ctx, session.TokenHash, replacement, session.CreatedAt); err != domainauth.ErrTokenInvalid {
		t.Fatal("losing claim must fail")
	}
	if err := repo.Consume(context.Background(), session.TokenHash, replacement, session.CreatedAt); !errors.Is(err, shared.ErrInternal) {
		t.Fatal("autocommit consume accepted")
	}
}
func TestRevocationAndCleanupRepositoryContracts(t *testing.T) {
	session, tx, ctx := sessionFixture()
	repo := NewAuthRepository(nil)
	tx.affected = 0
	if err := repo.Revoke(ctx, session.TokenHash); err != nil {
		t.Fatal("unknown/already-revoked should succeed")
	}
	if tx.args[0] != session.TokenHash || !strings.Contains(tx.query, "revoked_at IS NULL") {
		t.Fatal("idempotent hash revocation missing")
	}
	if err := repo.RevokeAllForUser(ctx, session.UserID); err != nil {
		t.Fatal(err)
	}
	if tx.args[0] != session.UserID || !strings.Contains(tx.query, "user_id = $1 AND revoked_at IS NULL") {
		t.Fatal("revoke-all capability weakened")
	}
	tx.affected = 2
	n, err := repo.Cleanup(ctx, session.CreatedAt, 1000)
	if err != nil || n != 2 {
		t.Fatal("cleanup count incorrect")
	}
	for _, part := range []string{"LEAST(expires_at, COALESCE(revoked_at, expires_at)) < $1", "LIMIT $2 FOR UPDATE SKIP LOCKED", "DELETE FROM"} {
		if !strings.Contains(tx.query, part) {
			t.Fatal("cleanup must retain active/recent rows and bound work", part)
		}
	}
	if tx.args[0] != session.CreatedAt || tx.args[1] != 1000 {
		t.Fatal("cleanup parameters incorrect")
	}
	for _, limit := range []int{-1, 0, 1001} {
		if _, err := repo.Cleanup(ctx, session.CreatedAt, limit); err != shared.ErrInvalidInput {
			t.Fatal("unbounded cleanup accepted")
		}
	}
}
func TestSessionDatabaseErrorsDiscardCredentialDetail(t *testing.T) {
	session, tx, ctx := sessionFixture()
	repo := NewAuthRepository(nil)
	tx.err = &pgconn.PgError{Code: "23505", Message: "raw-token-sentinel", Detail: session.TokenHash}
	checks := []func() error{
		func() error { return repo.Create(ctx, session) },
		func() error { _, err := repo.FindByHashForUpdate(ctx, session.TokenHash); return err },
		func() error { return repo.Consume(ctx, session.TokenHash, uuid.New(), session.CreatedAt) },
		func() error { return repo.Revoke(ctx, session.TokenHash) },
		func() error { return repo.RevokeAllForUser(ctx, session.UserID) },
		func() error { _, err := repo.Cleanup(ctx, session.CreatedAt, 10); return err },
	}
	for _, check := range checks {
		err := check()
		if !errors.Is(err, shared.ErrInternal) || strings.Contains(err.Error(), session.TokenHash) || strings.Contains(err.Error(), "raw-token-sentinel") {
			t.Fatal("database detail escaped repository")
		}
	}
}
func TestRefreshAccountLockUsesSafeProjectionAndTransaction(t *testing.T) {
	tx := &recordingTx{}
	ctx := database.WithTx(context.Background(), tx)
	repo := NewUserRepository(nil)
	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	account, err := repo.LockAccountByID(ctx, id)
	if err != nil || account.ID != id {
		t.Fatal("safe locked account scan failed")
	}
	if !strings.Contains(tx.query, "FOR SHARE") || strings.Contains(tx.query, "KEY SHARE") || strings.Contains(tx.query, "password") || tx.args[0] != id {
		t.Fatal("lock must conflict with non-key status/role updates and omit secrets")
	}
	if _, err := repo.LockAccountByID(context.Background(), id); !errors.Is(err, shared.ErrInternal) {
		t.Fatal("autocommit account lock accepted")
	}
	tx.err = pgx.ErrNoRows
	if _, err := repo.LockAccountByID(ctx, id); err == nil {
		t.Fatal("missing account accepted")
	}
	tx.err = errors.New("password-sentinel")
	if _, err := repo.LockAccountByID(ctx, id); !errors.Is(err, shared.ErrInternal) {
		t.Fatal("locked account leaked detail")
	}
}
