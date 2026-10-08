package postgres

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	domainuser "github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Record the actual repository statement through its transaction seam, without
// a live server. These checks protect the sensitive projection and write scope.
type recordingTx struct {
	pgx.Tx
	query string
	args  []any
	err   error
}

func (tx *recordingTx) QueryRow(_ context.Context, q string, args ...any) pgx.Row {
	tx.query = q
	tx.args = args
	return accountRow{err: tx.err}
}

type accountRow struct{ err error }

func (row accountRow) Scan(dest ...any) error {
	if row.err != nil {
		return row.err
	}
	values := []any{uuid.MustParse("00000000-0000-0000-0000-000000000001"), "synthetic@example.invalid", "Safe Name", "BORROWER", true, time.Time{}, time.Time{}}
	if len(dest) != len(values) {
		return errors.New("unexpected sensitive projection")
	}
	for i, value := range values {
		reflect.ValueOf(dest[i]).Elem().Set(reflect.ValueOf(value))
	}
	return nil
}
func TestAccountLookupProjectsNoSecrets(t *testing.T) {
	tx := &recordingTx{}
	ctx := database.WithTx(context.Background(), tx)
	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	account, err := NewUserRepository(nil).FindAccountByID(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	if account.ID != id || account.Role != domainuser.RoleBorrower || !account.IsActive {
		t.Fatal("safe account not scanned")
	}
	if strings.Contains(tx.query, "password") || strings.Contains(tx.query, "refresh") || strings.Contains(tx.query, id.String()) || !strings.Contains(tx.query, "WHERE id = $1") || len(tx.args) != 1 || tx.args[0] != id {
		t.Fatal("lookup must project safe fields and parameterize identity")
	}
	tx.err = pgx.ErrNoRows
	if _, err := NewUserRepository(nil).FindAccountByID(ctx, id); !errors.Is(err, domainuser.ErrUserNotFound) {
		t.Fatal("missing account mapping incorrect")
	}
}
func TestSelfUpdateSQLPreservesSecurityFields(t *testing.T) {
	tx := &recordingTx{}
	ctx := database.WithTx(context.Background(), tx)
	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	name := "Robert'); UPDATE users SET role='admin'; --"
	if _, err := NewUserRepository(nil).UpdateProfile(ctx, id, name); err != nil {
		t.Fatal(err)
	}
	set := strings.Split(strings.Split(tx.query, "SET ")[1], "WHERE")[0]
	for _, sensitive := range []string{"role", "is_active", "password", "email", "created_at", "id ="} {
		if strings.Contains(set, sensitive) {
			t.Fatal("profile SQL writes security/identity field", sensitive)
		}
	}
	if !strings.Contains(tx.query, "is_active = TRUE") || !strings.Contains(tx.query, "role IN ('BORROWER','STAFF','ADMIN')") || strings.Contains(tx.query, name) || len(tx.args) != 2 || tx.args[0] != id || tx.args[1] != name {
		t.Fatal("mutation must recheck permission and parameterize name/identity")
	}
	if strings.Contains(tx.query, "password") {
		t.Fatal("profile return projection contains hash")
	}
	tx.err = pgx.ErrNoRows
	if _, err := NewUserRepository(nil).UpdateProfile(ctx, id, name); !errors.Is(err, shared.ErrForbidden) {
		t.Fatal("deleted/deactivated/unknown account must not be updated")
	}
}
