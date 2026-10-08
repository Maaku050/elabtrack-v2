package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Maaku050/elabtrack-v2/backend/internal/application"
	appterms "github.com/Maaku050/elabtrack-v2/backend/internal/application/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/terms"
	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/user"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/database"
	"github.com/Maaku050/elabtrack-v2/backend/internal/infrastructure/persistence/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"os"
	"sync"
	"testing"
	"time"
)

func termsDatabase(t *testing.T) (context.Context, *database.Postgres, *database.Postgres, *config.Config) {
	t.Helper()
	batch1 := os.Getenv("ELABTRACK_BATCH1") == "1"
	if os.Getenv("ELABTRACK_PHASE4B") != "1" && !batch1 {
		t.Skip("requires explicitly isolated Phase 4B database")
	}
	cfg, err := config.Parse(map[string]string{"APP_ENV": "test", "JWT_SECRET": os.Getenv("JWT_SECRET"), "DB_HOST": "127.0.0.1", "DB_PORT": "35432", "DB_NAME": "elabtrack_v2_phase4b_test", "DB_USER": "elabtrack_runtime", "DB_PASSWORD": os.Getenv("DB_PASSWORD"), "MIGRATION_DATABASE_URL": os.Getenv("MIGRATION_DATABASE_URL")})
	if batch1 {
		cfg, err = config.Load()
		require(t, err == nil && cfg.DB.Name == "elabtrack_v2_batch1_test" && cfg.DB.Port == "54832", "isolated batch terms guard")
	}
	require(t, err == nil, "isolated terms config")
	ctx := context.Background()
	runtime, err := database.New(ctx, cfg.DB)
	require(t, err == nil, "runtime terms connection")
	t.Cleanup(runtime.Close)
	ownerCfg, err := cfg.MigrationConnection()
	require(t, err == nil, "owner config")
	owner, err := database.New(ctx, ownerCfg)
	require(t, err == nil, "owner terms connection")
	t.Cleanup(owner.Close)
	for _, pair := range []struct {
		db   *database.Postgres
		role string
	}{{runtime, "elabtrack_runtime"}, {owner, "elabtrack_migrator"}} {
		var safe bool
		err = pair.db.Pool.QueryRow(ctx, `SELECT current_database()=$2 AND current_user=$1 AND NOT rolsuper AND NOT rolcreaterole AND NOT rolcreatedb FROM pg_roles WHERE rolname=current_user`, pair.role, cfg.DB.Name).Scan(&safe)
		require(t, err == nil && safe, "strict isolated identity")
	}
	return ctx, runtime, owner, cfg
}
func termsUser(t *testing.T, ctx context.Context, db *database.Postgres, role user.Role) *user.User {
	t.Helper()
	u := user.NewUser("phase4b-test-"+uuid.NewString()+"@example.invalid", "Synthetic terms fixture", "not-a-login")
	u.Role = role
	require(t, postgres.NewUserRepository(db.Pool).Create(ctx, u) == nil, "owned account fixture")
	// Immutable history intentionally retains these IDs only in the disposable DB.
	return u
}

type controlledTx struct {
	base  application.Transactions
	after func() error
}

func (c controlledTx) Within(ctx context.Context, fn func(context.Context) error) error {
	return c.base.Within(ctx, func(ctx context.Context) error {
		if err := fn(ctx); err != nil {
			return err
		}
		return c.after()
	})
}
func waitForDatabaseLock(t *testing.T, ctx context.Context, db *database.Postgres) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var waiting bool
		err := db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pg_stat_activity WHERE datname=current_database() AND usename=current_user AND wait_event_type='Lock')`).Scan(&waiting)
		require(t, err == nil, "lock introspection")
		if waiting {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("expected real publication/acceptance lock waiter")
}

func TestRealTerms(t *testing.T) {
	ctx, runtime, owner, _ := termsDatabase(t)
	var count int
	require(t, owner.Pool.QueryRow(ctx, `SELECT count(*) FROM terms_versions`).Scan(&count) == nil && count == 0, "run terms verification before any publications; refuse existing history")
	repo := postgres.NewTermsRepository(runtime.Pool)
	accounts := postgres.NewUserRepository(runtime.Pool)
	tx := database.NewTxManager(runtime.Pool)
	svc := appterms.NewService(repo, accounts, tx)
	admin := termsUser(t, ctx, runtime, user.RoleAdmin)
	borrower := termsUser(t, ctx, runtime, user.RoleBorrower)
	second := termsUser(t, ctx, runtime, user.RoleBorrower)
	staff := termsUser(t, ctx, runtime, user.RoleStaff)
	inactive := termsUser(t, ctx, runtime, user.RoleBorrower)
	require(t, func() bool {
		_, err := runtime.Pool.Exec(ctx, `UPDATE users SET is_active=false WHERE id=$1`, inactive.ID)
		return err == nil
	}(), "inactive fixture")
	var current *terms.Version
	publish := func(version string) *terms.Version {
		t.Helper()
		var expected *uuid.UUID
		if current != nil {
			expected = &current.ID
		}
		v, err := svc.Publish(ctx, appterms.PublishCommand{ActorID: admin.ID, Version: version, Title: "SYNTHETIC TEST TERMS — NOT OFFICIAL", Body: "TEST ONLY. This document is not approved FSMO policy.\nVersion " + version, ExpectedCurrentVersionID: expected})
		require(t, err == nil, "synthetic isolated publication")
		current = v
		return v
	}
	t.Run("empty_migration_down_up_preserves_accounts_sessions", func(t *testing.T) {
		sum := sha256.Sum256([]byte(uuid.NewString()))
		_, err := runtime.Pool.Exec(ctx, `INSERT INTO refresh_tokens(id,user_id,token_hash,expires_at)VALUES($1,$2,$3,now()+interval '1 day')`, uuid.New(), borrower.ID, hex.EncodeToString(sum[:]))
		require(t, err == nil, "session fixture")
		snapshot := func() string {
			var s string
			err := owner.Pool.QueryRow(ctx, `SELECT jsonb_build_object('accounts',(SELECT jsonb_agg(to_jsonb(u) ORDER BY id) FROM users u),'sessions',(SELECT jsonb_agg(to_jsonb(r) ORDER BY id) FROM refresh_tokens r))::text`).Scan(&s)
			require(t, err == nil, "account/session snapshot")
			return s
		}
		before := snapshot()
		m := database.NewMigrator(owner.Pool, "../../migrations")
		require(t, m.Down(ctx) == nil, "empty inventory rollback before account and terms pairs")
		require(t, m.Down(ctx) == nil, "empty account migration rollback before terms")
		require(t, m.Down(ctx) == nil, "empty terms rollback")
		require(t, m.Up(ctx) == nil, "terms reapplication")
		require(t, snapshot() == before, "exact unchanged accounts/sessions")
	})
	t.Run("missing_terms_and_no_automatic_consent", func(t *testing.T) {
		status, err := svc.Status(ctx, borrower.ID)
		require(t, err == nil && status.State == "unpublished" && !status.CanInitiateBorrowing && status.Acceptance == nil, "safe missing state")
		_, err = svc.Accept(ctx, borrower.ID, uuid.New())
		require(t, errors.Is(err, terms.ErrNotPublished), "missing publication acceptance denied")
		err = tx.Within(ctx, func(ctx context.Context) error { _, e := svc.RequireCurrentAcceptance(ctx, borrower.ID); return e })
		require(t, errors.Is(err, terms.ErrNotPublished), "missing publication command eligibility denied")
	})
	t.Run("publication_and_role_authority", func(t *testing.T) {
		for _, actor := range []*user.User{borrower, staff, inactive} {
			_, err := svc.Publish(ctx, appterms.PublishCommand{ActorID: actor.ID, Version: "TEST-denied", Title: "Synthetic", Body: "TEST ONLY"})
			require(t, errors.Is(err, shared.ErrForbidden), "current role/status publication denied")
		}
		v := publish("TEST-1")
		require(t, v.PublishedBy == admin.ID && !v.PublishedAt.IsZero(), "durable publisher/time")
		sum := sha256.Sum256([]byte(v.Body))
		require(t, v.ContentHash == hex.EncodeToString(sum[:]), "exact published content hash")
		for _, actor := range []*user.User{staff, admin, inactive} {
			_, err := svc.Accept(ctx, actor.ID, v.ID)
			require(t, errors.Is(err, shared.ErrForbidden), "only active Borrower may accept")
		}
		_, err := svc.Accept(ctx, uuid.Nil, v.ID)
		require(t, errors.Is(err, shared.ErrUnauthorized), "identity required")
		_, err = svc.Accept(ctx, uuid.New(), v.ID)
		require(t, errors.Is(err, shared.ErrUnauthorized), "missing account denied")
	})
	var first *terms.Acceptance
	t.Run("first_repeat_and_account_scope", func(t *testing.T) {
		var err error
		first, err = svc.Accept(ctx, borrower.ID, current.ID)
		require(t, err == nil && first.UserID == borrower.ID && !first.AcceptedAt.IsZero(), "first durable acceptance")
		repeat, err := svc.Accept(ctx, borrower.ID, current.ID)
		require(t, err == nil && repeat.ID == first.ID && repeat.AcceptedAt.Equal(first.AcceptedAt), "idempotent original evidence")
		status, err := svc.Status(ctx, second.ID)
		require(t, err == nil && status.State == "required" && status.Acceptance == nil, "other account cannot inherit acceptance")
		_, err = svc.Accept(ctx, borrower.ID, uuid.New())
		require(t, errors.Is(err, terms.ErrVersionNotFound), "unknown version denied")
		err = tx.Within(ctx, func(ctx context.Context) error {
			a, e := svc.RequireCurrentAcceptance(ctx, borrower.ID)
			if e == nil && a.ID != first.ID {
				return shared.ErrInternal
			}
			return e
		})
		require(t, err == nil, "future command exact evidence")
	})
	oldID := current.ID
	t.Run("replacement_history_and_stale_consent", func(t *testing.T) {
		publish("TEST-2")
		status, err := svc.Status(ctx, borrower.ID)
		require(t, err == nil && status.State == "updated" && status.Acceptance == nil && status.HasPreviousAcceptance && !status.CanInitiateBorrowing, "updated version requires consent")
		_, err = svc.Accept(ctx, borrower.ID, oldID)
		require(t, errors.Is(err, terms.ErrVersionChanged), "old reviewed version conflicts")
		err = tx.Within(ctx, func(ctx context.Context) error { _, e := svc.RequireCurrentAcceptance(ctx, borrower.ID); return e })
		require(t, errors.Is(err, terms.ErrAcceptanceRequired), "old acceptance cannot authorize new command")
		_, err = svc.Accept(ctx, borrower.ID, current.ID)
		require(t, err == nil, "new version acceptance")
		retained, err := repo.FindAcceptance(ctx, borrower.ID, oldID)
		require(t, err == nil && retained.ID == first.ID && retained.AcceptedAt.Equal(first.AcceptedAt), "old immutable history retained")
	})
	t.Run("parallel_acceptance_one_original_receipt", func(t *testing.T) {
		const n = 12
		results := make(chan *terms.Acceptance, n)
		failures := make(chan error, n)
		var group sync.WaitGroup
		for i := 0; i < n; i++ {
			group.Add(1)
			go func() {
				defer group.Done()
				a, e := svc.Accept(ctx, second.ID, current.ID)
				results <- a
				failures <- e
			}()
		}
		group.Wait()
		close(results)
		close(failures)
		for e := range failures {
			require(t, e == nil, "parallel acceptance succeeded")
		}
		var winner *terms.Acceptance
		for a := range results {
			if winner == nil {
				winner = a
			}
			require(t, a.ID == winner.ID && a.AcceptedAt.Equal(winner.AcceptedAt), "same original receipt")
		}
		require(t, runtime.Pool.QueryRow(ctx, `SELECT count(*) FROM terms_acceptances WHERE user_id=$1 AND terms_version_id=$2`, second.ID, current.ID).Scan(&count) == nil && count == 1, "unique parallel evidence")
	})
	t.Run("acceptance_and_publication_transaction_rollback", func(t *testing.T) {
		fresh := termsUser(t, ctx, runtime, user.RoleBorrower)
		rollback := errors.New("controlled rollback")
		broken := appterms.NewService(repo, accounts, controlledTx{tx, func() error { return rollback }})
		_, err := broken.Accept(ctx, fresh.ID, current.ID)
		require(t, errors.Is(err, rollback), "controlled post-insert failure")
		a, err := repo.FindAcceptance(ctx, fresh.ID, current.ID)
		require(t, err == nil && a == nil, "acceptance rolled back")
		_, err = broken.Publish(ctx, appterms.PublishCommand{ActorID: admin.ID, Version: "TEST-rollback", Title: "Synthetic", Body: "TEST ONLY", ExpectedCurrentVersionID: &current.ID})
		require(t, errors.Is(err, rollback), "controlled post-publication failure")
		now, err := svc.Current(ctx, admin.ID)
		require(t, err == nil && now.ID == current.ID, "pointer rolled back")
		require(t, runtime.Pool.QueryRow(ctx, `SELECT count(*) FROM terms_versions WHERE version='TEST-rollback'`).Scan(&count) == nil && count == 0, "version insertion rolled back")
	})
	t.Run("constraints_immutability_and_runtime_privileges", func(t *testing.T) {
		for index, c := range []struct {
			sql, code string
			args      []any
		}{
			{`INSERT INTO terms_acceptances(id,user_id,terms_version_id)VALUES($1,$2,$3)`, "23505", []any{uuid.New(), borrower.ID, oldID}},
			{`INSERT INTO terms_acceptances(id,user_id,terms_version_id)VALUES($1,$2,$3)`, "23503", []any{uuid.New(), uuid.New(), current.ID}},
			{`INSERT INTO terms_acceptances(id,user_id,terms_version_id)VALUES($1,$2,$3)`, "23503", []any{uuid.New(), borrower.ID, uuid.New()}},
			{`INSERT INTO terms_acceptances(id,user_id,terms_version_id,accepted_at)VALUES($1,$2,$3,NULL)`, "23502", []any{uuid.New(), uuid.New(), current.ID}},
			{`UPDATE terms_acceptances SET accepted_at=now()`, "42501", nil}, {`DELETE FROM terms_acceptances`, "42501", nil}, {`UPDATE terms_versions SET body='changed'`, "42501", nil}, {`TRUNCATE terms_versions`, "42501", nil},
			{`UPDATE terms_publication SET id=2`, "42501", nil}, {`DELETE FROM terms_publication`, "42501", nil}, {`SELECT * FROM schema_migrations`, "42501", nil}, {`CREATE TABLE forbidden_terms(id int)`, "42501", nil},
			{`DELETE FROM users WHERE id=$1`, "23001", []any{borrower.ID}},
		} {
			_, err := runtime.Pool.Exec(ctx, c.sql, c.args...)
			var pgErr *pgconn.PgError
			if !errors.As(err, &pgErr) {
				t.Fatalf("constraint %d: expected SQLSTATE %s, received safe type %T", index, c.code, err)
			}
			if pgErr.Code != c.code {
				t.Fatalf("constraint %d: SQLSTATE %s, expected %s", index, pgErr.Code, c.code)
			}
		}
		for _, sql := range []string{`UPDATE terms_versions SET title=title`, `DELETE FROM terms_acceptances`, `TRUNCATE terms_acceptances`} {
			_, err := owner.Pool.Exec(ctx, sql)
			var pgErr *pgconn.PgError
			require(t, errors.As(err, &pgErr) && pgErr.Code == "P0001", "owner cannot edit/delete/truncate history")
		}
		_, err := svc.Publish(ctx, appterms.PublishCommand{ActorID: admin.ID, Version: "TEST-2", Title: "Synthetic", Body: "TEST ONLY", ExpectedCurrentVersionID: &current.ID})
		require(t, errors.Is(err, terms.ErrVersionExists), "duplicate version identifier refused")
	})
	t.Run("publication_wins_old_view_conflicts", func(t *testing.T) {
		fresh := termsUser(t, ctx, runtime, user.RoleBorrower)
		old := current.ID
		locked, release := make(chan struct{}), make(chan struct{})
		var releaseOnce sync.Once
		unblock := func() { releaseOnce.Do(func() { close(release) }) }
		defer unblock()
		done := make(chan error, 1)
		v, _ := terms.NewVersion("TEST-3", "SYNTHETIC TEST TERMS — NOT OFFICIAL", "TEST ONLY. Publication-first race.", admin.ID)
		go func() {
			done <- tx.Within(ctx, func(ctx context.Context) error {
				if _, err := accounts.LockAccountByID(ctx, admin.ID); err != nil {
					return err
				}
				if _, err := repo.LockPublication(ctx, true); err != nil {
					return err
				}
				if err := repo.Publish(ctx, v); err != nil {
					return err
				}
				close(locked)
				<-release
				return nil
			})
		}()
		<-locked
		accepted := make(chan error, 1)
		go func() { _, e := svc.Accept(ctx, fresh.ID, old); accepted <- e }()
		waitForDatabaseLock(t, ctx, runtime)
		unblock()
		require(t, <-done == nil, "publication committed")
		require(t, errors.Is(<-accepted, terms.ErrVersionChanged), "reviewed old version was not replaced silently")
		a, err := repo.FindAcceptance(ctx, fresh.ID, v.ID)
		require(t, err == nil && a == nil, "no fabricated new-version consent")
		current = v
	})
	t.Run("acceptance_wins_preserves_reviewed_version", func(t *testing.T) {
		fresh := termsUser(t, ctx, runtime, user.RoleBorrower)
		locked, release := make(chan struct{}), make(chan struct{})
		var releaseOnce sync.Once
		unblock := func() { releaseOnce.Do(func() { close(release) }) }
		defer unblock()
		accepted := make(chan error, 1)
		held := appterms.NewService(repo, accounts, controlledTx{tx, func() error { close(locked); <-release; return nil }})
		old := current.ID
		go func() { _, err := held.Accept(ctx, fresh.ID, old); accepted <- err }()
		<-locked
		published := make(chan *terms.Version, 1)
		failed := make(chan error, 1)
		go func() {
			v, e := svc.Publish(ctx, appterms.PublishCommand{ActorID: admin.ID, Version: "TEST-4", Title: "SYNTHETIC TEST TERMS — NOT OFFICIAL", Body: "TEST ONLY. Acceptance-first race.", ExpectedCurrentVersionID: &old})
			published <- v
			failed <- e
		}()
		waitForDatabaseLock(t, ctx, runtime)
		unblock()
		require(t, <-accepted == nil && <-failed == nil, "both ordered transactions committed")
		current = <-published
		a, err := repo.FindAcceptance(ctx, fresh.ID, old)
		require(t, err == nil && a != nil && a.TermsVersionID == old, "reviewed evidence retained")
		status, err := svc.Status(ctx, fresh.ID)
		require(t, err == nil && status.State == "updated", "new current version remains unaccepted")
	})
	t.Run("parallel_publishers_no_blind_overwrite", func(t *testing.T) {
		old := current.ID
		type result struct {
			version *terms.Version
			err     error
		}
		results := make(chan result, 2)
		for _, name := range []string{"TEST-5a", "TEST-5b"} {
			go func(name string) {
				v, e := svc.Publish(ctx, appterms.PublishCommand{ActorID: admin.ID, Version: name, Title: "SYNTHETIC TEST TERMS — NOT OFFICIAL", Body: "TEST ONLY. Optimistic publication race.", ExpectedCurrentVersionID: &old})
				results <- result{v, e}
			}(name)
		}
		wins, conflicts := 0, 0
		for i := 0; i < 2; i++ {
			r := <-results
			if r.err == nil {
				wins++
				current = r.version
			} else if errors.Is(r.err, terms.ErrPublicationChanged) {
				conflicts++
			} else {
				t.Fatal("unexpected publication error")
			}
		}
		require(t, wins == 1 && conflicts == 1, "one publication winner, one explicit conflict")
	})
	t.Run("history_safe_rollback_and_policy_transaction_boundary", func(t *testing.T) {
		m := database.NewMigrator(owner.Pool, "../../migrations")
		require(t, m.Down(ctx) == nil, "empty inventory rollback before account and terms pairs")
		require(t, m.Down(ctx) == nil, "empty account migration rollback before terms history gate")
		require(t, m.Down(ctx) != nil, "history-bearing rollback refused")
		require(t, m.Up(ctx) == nil, "restore empty account schema without touching terms history")
		var intact bool
		err := owner.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version='000005_terms_acceptance') AND EXISTS(SELECT 1 FROM terms_acceptances WHERE id=$1)`, first.ID).Scan(&intact)
		require(t, err == nil && intact, "history/checksum bookkeeping retained")
		_, err = svc.RequireCurrentAcceptance(ctx, borrower.ID)
		require(t, errors.Is(err, shared.ErrInternal), "future commands must supply transaction")
	})
}
