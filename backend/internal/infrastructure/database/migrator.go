package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrator applies simple paired *.up.sql / *.down.sql migration files
// tracked by a `schema_migrations` table. It is intentionally lightweight
// (no external dependencies) and sufficient for the template. For more
// advanced needs (e.g. embedded migrations, dry-run, checksums) consider
// golang-migrate.
type Migrator struct {
	pool    *pgxpool.Pool
	dir     string
	nowFunc func() time.Time
}

// NewMigrator constructs a Migrator pointing at the given migrations dir.
func NewMigrator(pool *pgxpool.Pool, dir string) *Migrator {
	return &Migrator{pool: pool, dir: dir, nowFunc: time.Now}
}

// Up applies all pending up migrations in order.
func (m *Migrator) Up(ctx context.Context) error {
	if err := m.ensureTable(ctx); err != nil {
		return err
	}
	applied, err := m.appliedVersions(ctx)
	if err != nil {
		return err
	}

	files, err := m.listUpMigrations()
	if err != nil {
		return err
	}

	for _, f := range files {
		if applied[f.version] {
			continue
		}
		if err := m.applyOne(ctx, f); err != nil {
			return err
		}
		fmt.Printf("applied migration %s\n", f.version)
	}
	return nil
}

// Down rolls back the most recently applied migration.
func (m *Migrator) Down(ctx context.Context) error {
	if err := m.ensureTable(ctx); err != nil {
		return err
	}
	last, err := m.lastAppliedVersion(ctx)
	if err != nil {
		return err
	}
	if last == "" {
		fmt.Println("no migrations to roll back")
		return nil
	}

	downPath := filepath.Join(m.dir, last+".down.sql")
	content, err := os.ReadFile(downPath)
	if err != nil {
		return fmt.Errorf("read down migration %s: %w", downPath, err)
	}
	if err := m.execTx(ctx, string(content)); err != nil {
		return fmt.Errorf("apply down migration %s: %w", last, err)
	}
	if _, err := m.pool.Exec(ctx, `DELETE FROM schema_migrations WHERE version = $1`, last); err != nil {
		return fmt.Errorf("delete migration record: %w", err)
	}
	fmt.Printf("rolled back migration %s\n", last)
	return nil
}

// Create writes a new paired migration scaffold with the given name.
func (m *Migrator) Create(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("migration name is required")
	}
	version := m.nextVersion()
	base := fmt.Sprintf("%06d_%s", version, sanitize(name))
	up := filepath.Join(m.dir, base+".up.sql")
	down := filepath.Join(m.dir, base+".down.sql")
	if err := os.MkdirAll(m.dir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(up, []byte("-- "+base+" up\n\n"), 0o644); err != nil {
		return "", err
	}
	if err := os.WriteFile(down, []byte("-- "+base+" down\n\n"), 0o644); err != nil {
		return "", err
	}
	fmt.Printf("created migration %s\n", base)
	return base, nil
}

type migrationFile struct {
	version string
	path    string
}

func (m *Migrator) ensureTable(ctx context.Context) error {
	const q = `CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`
	_, err := m.pool.Exec(ctx, q)
	return err
}

func (m *Migrator) appliedVersions(ctx context.Context) (map[string]bool, error) {
	rows, err := m.pool.Query(ctx, `SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, fmt.Errorf("list applied migrations: %w", err)
	}
	defer rows.Close()
	out := map[string]bool{}
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out[v] = true
	}
	return out, rows.Err()
}

func (m *Migrator) lastAppliedVersion(ctx context.Context) (string, error) {
	var v string
	err := m.pool.QueryRow(ctx, `SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1`).Scan(&v)
	if err != nil {
		// no rows -> nothing applied
		return "", nil
	}
	return v, nil
}

func (m *Migrator) listUpMigrations() ([]migrationFile, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}
	var out []migrationFile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}
		version := strings.TrimSuffix(name, ".up.sql")
		out = append(out, migrationFile{version: version, path: filepath.Join(m.dir, name)})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].version < out[j].version })
	return out, nil
}

func (m *Migrator) applyOne(ctx context.Context, f migrationFile) error {
	content, err := os.ReadFile(f.path)
	if err != nil {
		return fmt.Errorf("read migration %s: %w", f.path, err)
	}
	if err := m.execTx(ctx, string(content)); err != nil {
		return fmt.Errorf("apply migration %s: %w", f.version, err)
	}
	_, err = m.pool.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, f.version)
	return err
}

// execTx runs the migration SQL inside a transaction so a partially-applied
// migration is rolled back on error.
func (m *Migrator) execTx(ctx context.Context, sql string) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, sql); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (m *Migrator) nextVersion() int {
	files, err := m.listUpMigrations()
	if err != nil || len(files) == 0 {
		return 1
	}
	last := files[len(files)-1].version
	parts := strings.SplitN(last, "_", 2)
	var n int
	fmt.Sscanf(parts[0], "%d", &n)
	return n + 1
}

func sanitize(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "_")
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Status reads tracking state without creating schema_migrations or executing
// migration SQL. Schema mutation is reserved for explicit Up/Down commands.
func (m *Migrator) Status(ctx context.Context) error {
	var exists bool
	if err := m.pool.QueryRow(ctx, `SELECT to_regclass('schema_migrations') IS NOT NULL`).Scan(&exists); err != nil {
		return err
	}
	applied := map[string]bool{}
	if exists {
		var err error
		applied, err = m.appliedVersions(ctx)
		if err != nil {
			return err
		}
	}
	files, err := m.listUpMigrations()
	if err != nil {
		return err
	}
	for _, f := range files {
		state := "pending"
		if applied[f.version] {
			state = "applied"
		}
		fmt.Printf("%s %s\n", state, f.version)
	}
	return nil
}
