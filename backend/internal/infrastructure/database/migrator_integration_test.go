package database

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Explicit opt-in; only the named disposable integration database is accepted.
// Run before API/auth tests. All fault DDL lives in fixtures, never migrations/.
func TestRealMigrator(t *testing.T) {
	if os.Getenv("ELABTRACK_MIGRATION_INTEGRATION") != "1" {
		t.Skip("requires Phase 1H disposable PostgreSQL")
	}
	cfg, err := config.Parse(map[string]string{"APP_ENV": "test", "JWT_SECRET": "migration-tests-only", "DB_HOST": "127.0.0.1", "DB_PORT": "15432", "DB_NAME": "elabtrack_v2_integration", "DB_USER": "elabtrack_runtime", "DB_PASSWORD": os.Getenv("DB_PASSWORD"), "MIGRATION_DATABASE_URL": os.Getenv("MIGRATION_DATABASE_URL")})
	must(t, err == nil, "config")
	migration, err := cfg.MigrationConnection()
	must(t, err == nil && migration.User == "elabtrack_migrator", "migration identity config")
	ctx := context.Background()
	owner, err := New(ctx, migration)
	must(t, err == nil, "owner connection")
	t.Cleanup(owner.Close)
	runtime, err := New(ctx, cfg.DB)
	must(t, err == nil, "runtime connection")
	t.Cleanup(runtime.Close)
	var identity bool
	err = owner.Pool.QueryRow(ctx, `SELECT current_database()='elabtrack_v2_integration' AND current_user='elabtrack_migrator' AND NOT rolsuper AND NOT rolcreaterole AND NOT rolcreatedb FROM pg_roles WHERE rolname=current_user`).Scan(&identity)
	must(t, err == nil && identity, "non-superuser migrator")
	base, err := NewMigrator(nil, "../../../migrations").discover()
	must(t, err == nil && len(base) == 3, "immutable foundation files")
	exec := func(sql string, args ...any) {
		_, err := owner.Pool.Exec(ctx, sql, args...)
		must(t, err == nil, "fixture SQL")
	}
	queryBool := func(sql string, args ...any) bool {
		var value bool
		err := owner.Pool.QueryRow(ctx, sql, args...).Scan(&value)
		must(t, err == nil, "state query")
		return value
	}
	makeRunner := func(t *testing.T, up, down string) *Migrator {
		dir := t.TempDir()
		for _, f := range base {
			fixturePair(t, dir, f.version, string(f.up), string(f.down))
		}
		if up != "" {
			fixturePair(t, dir, "000099_probe", up, down)
		}
		m := NewMigrator(owner.Pool, dir)
		m.output = &bytes.Buffer{}
		return m
	}
	cleanup := func() {
		exec(`DROP TABLE IF EXISTS public.phase1h_probe; DELETE FROM public.schema_migrations WHERE version='000099_probe'`)
	}
	expectFail := func(t *testing.T, err error, reason string) {
		t.Helper()
		must(t, err != nil && strings.Contains(err.Error(), reason), "safe expected migration failure: "+reason)
	}
	absent := func() bool {
		return queryBool(`SELECT to_regclass('public.phase1h_probe') IS NULL AND NOT EXISTS(SELECT 1 FROM schema_migrations WHERE version='000099_probe')`)
	}
	t.Run("atomic_up_SQL_failure", func(t *testing.T) {
		m := makeRunner(t, "CREATE TABLE phase1h_probe(id int); SELECT 1/0;", "DROP TABLE phase1h_probe;")
		expectFail(t, m.Up(ctx), "SQL execution failed")
		must(t, absent(), "DDL and record rolled back")
	})
	t.Run("atomic_up_bookkeeping_failure", func(t *testing.T) {
		m := makeRunner(t, "CREATE TABLE phase1h_probe(id int);", "DROP TABLE phase1h_probe;")
		exec(`ALTER TABLE schema_migrations ADD CONSTRAINT phase1h_reject CHECK(version <> '000099_probe')`)
		defer exec(`ALTER TABLE schema_migrations DROP CONSTRAINT phase1h_reject`)
		expectFail(t, m.Up(ctx), "bookkeeping insert failed")
		must(t, absent(), "bookkeeping rejection rolls back DDL")
	})
	t.Run("atomic_down_SQL_failure", func(t *testing.T) {
		m := makeRunner(t, "CREATE TABLE phase1h_probe(id int);", "DROP TABLE phase1h_probe; SELECT 1/0;")
		must(t, m.Up(ctx) == nil, "probe up")
		defer cleanup()
		expectFail(t, m.Down(ctx), "SQL execution failed")
		must(t, queryBool(`SELECT to_regclass('phase1h_probe') IS NOT NULL AND EXISTS(SELECT 1 FROM schema_migrations WHERE version='000099_probe')`), "down preserves table and record")
	})
	t.Run("atomic_down_bookkeeping_failure", func(t *testing.T) {
		m := makeRunner(t, "CREATE TABLE phase1h_probe(id int);", "DROP TABLE phase1h_probe;")
		must(t, m.Up(ctx) == nil, "probe up")
		defer cleanup()
		exec(`CREATE FUNCTION phase1h_reject_delete() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'controlled failure'; END $$; CREATE TRIGGER phase1h_reject_delete BEFORE DELETE ON schema_migrations FOR EACH ROW EXECUTE FUNCTION phase1h_reject_delete()`)
		defer exec(`DROP TRIGGER phase1h_reject_delete ON schema_migrations; DROP FUNCTION phase1h_reject_delete()`)
		expectFail(t, m.Down(ctx), "bookkeeping removal failed")
		must(t, queryBool(`SELECT to_regclass('phase1h_probe') IS NOT NULL AND EXISTS(SELECT 1 FROM schema_migrations WHERE version='000099_probe')`), "down bookkeeping rollback")
	})
	t.Run("checksum_up_and_down_tampering", func(t *testing.T) {
		m := makeRunner(t, "CREATE TABLE phase1h_probe(id int);", "DROP TABLE phase1h_probe;")
		must(t, m.Up(ctx) == nil, "probe up")
		defer cleanup()
		for _, direction := range []string{"up", "down"} {
			path := filepath.Join(m.dir, "000099_probe."+direction+".sql")
			original, e := os.ReadFile(path)
			must(t, e == nil, "read fixture")
			must(t, os.WriteFile(path, append(original, []byte("\n-- changed\n")...), 0600) == nil, "tamper fixture")
			expectFail(t, m.Status(ctx), "checksum mismatch")
			expectFail(t, m.Up(ctx), "checksum mismatch")
			expectFail(t, m.Down(ctx), "checksum mismatch")
			must(t, os.WriteFile(path, original, 0600) == nil, "restore fixture")
			must(t, m.Status(ctx) == nil, "checksum was not overwritten")
		}
	})
	t.Run("missing_applied_file", func(t *testing.T) {
		m := makeRunner(t, "CREATE TABLE phase1h_probe(id int);", "DROP TABLE phase1h_probe;")
		must(t, m.Up(ctx) == nil, "probe up")
		defer cleanup()
		must(t, os.Remove(filepath.Join(m.dir, "000099_probe.up.sql")) == nil, "remove fixture")
		must(t, os.Remove(filepath.Join(m.dir, "000099_probe.down.sql")) == nil, "remove fixture pair")
		expectFail(t, m.Up(ctx), "applied version has no repository migration")
	})
	t.Run("historical_gap", func(t *testing.T) {
		m := makeRunner(t, "", "")
		var applied time.Time
		var hash string
		err := owner.Pool.QueryRow(ctx, `SELECT applied_at,checksum FROM schema_migrations WHERE version='000002_create_refresh_tokens'`).Scan(&applied, &hash)
		must(t, err == nil, "save history")
		exec(`DELETE FROM schema_migrations WHERE version='000002_create_refresh_tokens'`)
		defer exec(`INSERT INTO schema_migrations(version,applied_at,checksum)VALUES('000002_create_refresh_tokens',$1,$2)`, applied, hash)
		expectFail(t, m.Up(ctx), "inconsistent history")
	})
	t.Run("duplicate_numeric_version", func(t *testing.T) {
		m := makeRunner(t, "", "")
		fixturePair(t, m.dir, "000003_duplicate", "SELECT 1;", "SELECT 1;")
		expectFail(t, m.Up(ctx), "duplicate numeric version")
	})
	t.Run("database_history_permission_failure", func(t *testing.T) {
		m := makeRunner(t, "", "")
		m.pool = runtime.Pool
		expectFail(t, m.Down(ctx), "migration history")
		expectFail(t, m.Status(ctx), "migration history")
	})
	t.Run("unsupported_nontransactional_SQL", func(t *testing.T) {
		m := makeRunner(t, "CREATE TABLE phase1h_probe(id int); CREATE INDEX CONCURRENTLY phase1h_probe_idx ON phase1h_probe(id);", "DROP TABLE phase1h_probe;")
		expectFail(t, m.Up(ctx), "nontransactional SQL is unsupported")
		must(t, absent(), "unsupported statement leaves no DDL/record")
	})
	t.Run("transaction_control_cannot_commit_DDL", func(t *testing.T) {
		m := makeRunner(t, "CREATE TABLE phase1h_probe(id int); COMMIT; SELECT 1/0;", "DROP TABLE phase1h_probe;")
		expectFail(t, m.Up(ctx), "transaction control is unsupported")
		must(t, absent(), "file cannot escape transaction")
	})
	t.Run("real_concurrent_connections", func(t *testing.T) {
		m := makeRunner(t, "CREATE TABLE phase1h_probe(id int); SELECT pg_sleep(1);", "DROP TABLE phase1h_probe;")
		defer cleanup()
		finished := make(chan error, 1)
		go func() { finished <- m.Up(ctx) }()
		deadline := time.Now().Add(5 * time.Second)
		held := false
		for time.Now().Before(deadline) {
			held = queryBool(`SELECT EXISTS(SELECT 1 FROM pg_locks WHERE locktype='advisory' AND classid=$1::oid AND objid=$2::oid AND granted)`, migrationLockKey>>32, migrationLockKey&0xffffffff)
			if held {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		must(t, held, "native lock held")
		other := NewMigrator(owner.Pool, m.dir)
		other.output = &bytes.Buffer{}
		expectFail(t, other.Up(ctx), "another migrator holds the lock")
		must(t, <-finished == nil, "winning runner succeeds")
		must(t, m.Status(ctx) == nil, "lock released and coherent checksums")
		must(t, m.Down(ctx) == nil && absent(), "one coherent applied/down state")
	})
	t.Run("runtime_DML_and_DDL_denials", func(t *testing.T) {
		exec(`CREATE TABLE phase1h_probe(id int)`)
		defer cleanup()
		conn, err := runtime.Pool.Acquire(ctx)
		must(t, err == nil, "runtime session")
		defer conn.Release()
		tx, err := conn.Begin(ctx)
		must(t, err == nil, "runtime DML tx")
		_, err = tx.Exec(ctx, `INSERT INTO users(id,email,name,password)VALUES('ffffffff-ffff-4fff-8fff-ffffffffffff','phase1h-dml@example.invalid','Synthetic','not-a-login'); UPDATE users SET name='Updated' WHERE id='ffffffff-ffff-4fff-8fff-ffffffffffff'; SELECT id FROM users WHERE id='ffffffff-ffff-4fff-8fff-ffffffffffff'; DELETE FROM users WHERE id='ffffffff-ffff-4fff-8fff-ffffffffffff'`)
		must(t, err == nil, "runtime full application DML")
		must(t, tx.Rollback(ctx) == nil, "DML fixture rollback")
		for _, sql := range []string{`CREATE TABLE phase1h_runtime_forbidden(id int)`, `ALTER TABLE phase1h_probe ADD COLUMN phase1h_forbidden int`, `CREATE ROLE phase1h_forbidden`, `SELECT version FROM schema_migrations`, `UPDATE schema_migrations SET checksum=checksum`} {
			_, err := conn.Exec(ctx, sql)
			must(t, err != nil, "prohibited runtime operation denied")
		}
		// PostgreSQL may return only a warning for GRANT without grant option.
		// The effective ACL, rather than its command tag, is authoritative.
		_, err = conn.Exec(ctx, `GRANT SELECT ON users TO PUBLIC`)
		must(t, queryBool(`SELECT NOT has_table_privilege('public','users','SELECT') AND NOT has_table_privilege('elabtrack_runtime','users','SELECT WITH GRANT OPTION')`), "runtime cannot grant privileges")
		must(t, queryBool(`SELECT to_regclass('phase1h_runtime_forbidden') IS NULL`), "denied DDL absent")
	})
	t.Run("migrator_DDL_without_cluster_privileges", func(t *testing.T) {
		exec(`BEGIN; CREATE TABLE phase1h_probe(id int); ALTER TABLE phase1h_probe ADD COLUMN label text; DROP TABLE phase1h_probe; ROLLBACK`)
		must(t, absent(), "migrator DDL rollback")
	})
	t.Run("explicit_legacy_adoption", func(t *testing.T) {
		m := makeRunner(t, "", "")
		exec(`ALTER TABLE schema_migrations RENAME TO phase1h_saved_tracking; CREATE TABLE schema_migrations(version TEXT PRIMARY KEY,applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()); INSERT INTO schema_migrations(version,applied_at)SELECT version,applied_at FROM phase1h_saved_tracking`)
		defer exec(`DROP TABLE schema_migrations; ALTER TABLE phase1h_saved_tracking RENAME TO schema_migrations`)
		expectFail(t, m.Up(ctx), "legacy tracking requires explicit")
		exec(`INSERT INTO schema_migrations(version)VALUES('unknown_external')`)
		expectFail(t, m.AdoptLegacy(ctx), "unknown legacy version")
		exec(`DELETE FROM schema_migrations WHERE version='unknown_external'`)
		must(t, m.AdoptLegacy(ctx) == nil, "explicit known-prefix adoption")
		must(t, m.Status(ctx) == nil, "adopted checksum verification")
		must(t, queryBool(`SELECT NOT EXISTS(SELECT 1 FROM schema_migrations AS a JOIN phase1h_saved_tracking AS b USING(version) WHERE a.applied_at<>b.applied_at OR a.checksum<>b.checksum)`), "timestamps retained and verified mapping")
		expectFail(t, m.AdoptLegacy(ctx), "already has checksums")
	})
}
func must(t *testing.T, ok bool, message string) {
	t.Helper()
	if !ok {
		t.Fatal(message)
	}
}

// Keep a closed-pool error from becoming successful empty down/status.
func TestClosedMigrationPool(t *testing.T) {
	cfg, err := pgxpool.ParseConfig("postgresql://synthetic:synthetic@127.0.0.1/unused?sslmode=disable")
	must(t, err == nil, "closed pool config")
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	must(t, err == nil, "lazy pool")
	pool.Close()
	m := NewMigrator(pool, "../../../migrations")
	must(t, m.Down(context.Background()) != nil && m.Status(context.Background()) != nil, "closed DB propagates failure")
}
