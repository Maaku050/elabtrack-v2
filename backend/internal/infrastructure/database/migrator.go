package database

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// One project-specific session lock, shared by up/down/status/adoption. The
// acquired connection is destroyed rather than returned locked to the pool.
const migrationLockKey int64 = 0x454c41424d494752 // ELABMIGR
var migrationName = regexp.MustCompile(`^([0-9]{6})_([a-z][a-z0-9]*(?:_[a-z0-9]+)*)\.(up|down)\.sql$`)

type Migrator struct {
	pool   *pgxpool.Pool
	dir    string
	log    *zap.Logger
	output io.Writer
}

func NewMigrator(pool *pgxpool.Pool, dir string, logs ...*zap.Logger) *Migrator {
	log := zap.NewNop()
	if len(logs) > 0 && logs[0] != nil {
		log = logs[0]
	}
	return &Migrator{pool: pool, dir: dir, log: log, output: os.Stdout}
}

// MigrationError contains only allowlisted reasons and canonical local versions;
// raw SQL/driver errors, paths and unknown database versions never escape.
type MigrationError struct{ Reason, Version string }

func (e *MigrationError) Error() string {
	if e.Version != "" {
		return "migration " + e.Version + ": " + e.Reason
	}
	return "migration: " + e.Reason
}
func migrationFailure(reason, version string) error {
	return &MigrationError{Reason: reason, Version: version}
}

type migrationFile struct {
	version  string
	number   int
	up, down []byte
	checksum string
}
type migrationRecord struct {
	checksum string
	applied  time.Time
}

func pairChecksum(up, down []byte) string {
	h := sha256.New()
	h.Write([]byte("elabtrack-v2-migration-pair-v1\x00"))
	for _, b := range [][]byte{up, down} {
		_ = binary.Write(h, binary.BigEndian, uint64(len(b)))
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))
}
func (m *Migrator) discover() ([]migrationFile, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return nil, migrationFailure("cannot read migration directory", "")
	}
	pairs := map[string]*migrationFile{}
	numbers := map[int]string{}
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		match := migrationName.FindStringSubmatch(name)
		if match == nil || !entry.Type().IsRegular() {
			return nil, migrationFailure("invalid migration filename or non-regular SQL file", "")
		}
		n, _ := strconv.Atoi(match[1])
		if n == 0 {
			return nil, migrationFailure("version zero is invalid", "")
		}
		version := match[1] + "_" + match[2]
		if old, ok := numbers[n]; ok && old != version {
			return nil, migrationFailure("duplicate numeric version", version)
		}
		numbers[n] = version
		f := pairs[version]
		if f == nil {
			f = &migrationFile{version: version, number: n}
			pairs[version] = f
		}
		content, err := os.ReadFile(filepath.Join(m.dir, name))
		if err != nil {
			return nil, migrationFailure("cannot read migration file", version)
		}
		if err := validateTransactionalSQL(string(content)); err != nil {
			return nil, migrationFailure("transaction control is unsupported in migration SQL", version)
		}
		if match[3] == "up" {
			f.up = content
		} else {
			f.down = content
		}
	}
	files := make([]migrationFile, 0, len(pairs))
	for _, f := range pairs {
		if f.up == nil || f.down == nil {
			return nil, migrationFailure("missing up/down counterpart", f.version)
		}
		f.checksum = pairChecksum(f.up, f.down)
		files = append(files, *f)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].number < files[j].number })
	return files, nil
}

func (m *Migrator) locked(ctx context.Context, action string, fn func(*pgx.Conn, []migrationFile) error) (err error) {
	start := time.Now()
	defer func() {
		result := "success"
		reason := ""
		version := ""
		if err != nil {
			result = "failure"
			var safe *MigrationError
			if errors.As(err, &safe) {
				reason = safe.Reason
				version = safe.Version
			}
		}
		m.log.Info("migration command completed", zap.String("event", "migration.command"), zap.String("direction", action), zap.String("result", result), zap.String("migration_version", version), zap.String("reason", reason), zap.Int64("duration_ms", time.Since(start).Milliseconds()))
	}()
	files, err := m.discover()
	if err != nil {
		return err
	}
	if m.pool == nil {
		return migrationFailure("database connection unavailable", "")
	}
	acquireCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pooled, err := m.pool.Acquire(acquireCtx)
	if err != nil {
		return migrationFailure("database connection failed", "")
	}
	conn := pooled.Hijack()
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = conn.Close(closeCtx)
	}()
	var owns bool
	if err := conn.QueryRow(acquireCtx, `SELECT pg_try_advisory_lock($1)`, migrationLockKey).Scan(&owns); err != nil {
		return migrationFailure("migration lock query failed", "")
	}
	if !owns {
		return migrationFailure("another migrator holds the lock; retry after it finishes", "")
	}
	return fn(conn, files)
}

func trackingExists(ctx context.Context, conn *pgx.Conn) (bool, error) {
	var exists bool
	if err := conn.QueryRow(ctx, `SELECT to_regclass('public.schema_migrations') IS NOT NULL`).Scan(&exists); err != nil {
		return false, migrationFailure("migration table lookup failed", "")
	}
	return exists, nil
}
func readRecords(ctx context.Context, conn *pgx.Conn) (map[string]migrationRecord, error) {
	var checksummed bool
	if err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pg_attribute WHERE attrelid='public.schema_migrations'::regclass AND attname='checksum' AND NOT attisdropped)`).Scan(&checksummed); err != nil {
		return nil, migrationFailure("migration metadata query failed", "")
	}
	if !checksummed {
		return nil, migrationFailure("legacy tracking requires explicit development-only --migrate-adopt-legacy after verifying provenance", "")
	}
	rows, err := conn.Query(ctx, `SELECT version, applied_at, checksum FROM public.schema_migrations ORDER BY version`)
	if err != nil {
		return nil, migrationFailure("migration history query failed", "")
	}
	defer rows.Close()
	records := map[string]migrationRecord{}
	for rows.Next() {
		var version string
		var record migrationRecord
		if err := rows.Scan(&version, &record.applied, &record.checksum); err != nil {
			return nil, migrationFailure("migration history scan failed", "")
		}
		if _, duplicate := records[version]; duplicate {
			return nil, migrationFailure("duplicate applied version in history", "")
		}
		records[version] = record
	}
	if rows.Err() != nil {
		return nil, migrationFailure("migration history iteration failed", "")
	}
	return records, nil
}
func validateHistory(files []migrationFile, records map[string]migrationRecord) error {
	known := map[string]bool{}
	missing := false
	for _, f := range files {
		known[f.version] = true
		record, applied := records[f.version]
		if !applied {
			missing = true
			continue
		}
		if missing {
			return migrationFailure("inconsistent history: lower repository version is pending", f.version)
		}
		if record.checksum != f.checksum {
			return migrationFailure("checksum mismatch; applied files are immutable", f.version)
		}
	}
	for version := range records {
		if !known[version] {
			return migrationFailure("applied version has no repository migration", "")
		}
	}
	return nil
}

const createTracking = `CREATE TABLE public.schema_migrations (
 version TEXT PRIMARY KEY CHECK (version ~ '^[0-9]{6}_[a-z][a-z0-9]*(_[a-z0-9]+)*$'),
 applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
 checksum TEXT NOT NULL CHECK (checksum ~ '^[0-9a-f]{64}$')
)`

func (m *Migrator) state(ctx context.Context, conn *pgx.Conn, files []migrationFile, create bool) (map[string]migrationRecord, error) {
	exists, err := trackingExists(ctx, conn)
	if err != nil {
		return nil, err
	}
	records := map[string]migrationRecord{}
	if exists {
		records, err = readRecords(ctx, conn)
		if err != nil {
			return nil, err
		}
	} else if create {
		if _, err := conn.Exec(ctx, createTracking); err != nil {
			return nil, migrationFailure("migration tracking creation failed", "")
		}
	}
	if err := validateHistory(files, records); err != nil {
		return nil, err
	}
	return records, nil
}
func (m *Migrator) Up(ctx context.Context) error {
	return m.locked(ctx, "up", func(conn *pgx.Conn, files []migrationFile) error {
		records, err := m.state(ctx, conn, files, true)
		if err != nil {
			return err
		}
		for _, f := range files {
			if _, applied := records[f.version]; applied {
				continue
			}
			if err := m.execute(ctx, conn, f, false); err != nil {
				return err
			}
		}
		return nil
	})
}
func (m *Migrator) Down(ctx context.Context) error {
	return m.locked(ctx, "down", func(conn *pgx.Conn, files []migrationFile) error {
		records, err := m.state(ctx, conn, files, false)
		if err != nil {
			return err
		}
		for i := len(files) - 1; i >= 0; i-- {
			if _, applied := records[files[i].version]; applied {
				return m.execute(ctx, conn, files[i], true)
			}
		}
		fmt.Fprintln(m.output, "no migrations applied")
		return nil
	})
}
func (m *Migrator) execute(ctx context.Context, conn *pgx.Conn, f migrationFile, down bool) (err error) {
	start := time.Now()
	direction := "up"
	sql := f.up
	if down {
		direction = "down"
		sql = f.down
	}
	defer func() {
		result := "success"
		if err != nil {
			result = "failure"
		}
		m.log.Info("migration completed", zap.String("event", "migration.apply"), zap.String("migration_version", f.version), zap.String("direction", direction), zap.String("result", result), zap.Int64("duration_ms", time.Since(start).Milliseconds()))
	}()
	tx, err := conn.Begin(ctx)
	if err != nil {
		return migrationFailure("transaction begin failed", f.version)
	}
	defer func() {
		rollbackCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tx.Rollback(rollbackCtx)
	}()
	if _, err := tx.Exec(ctx, string(sql)); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "25001" {
			return migrationFailure("nontransactional SQL is unsupported", f.version)
		}
		return migrationFailure("SQL execution failed; transaction rolled back", f.version)
	}
	if down {
		tag, e := tx.Exec(ctx, `DELETE FROM public.schema_migrations WHERE version=$1 AND checksum=$2`, f.version, f.checksum)
		if e != nil || tag.RowsAffected() != 1 {
			return migrationFailure("bookkeeping removal failed; transaction rolled back", f.version)
		}
	} else {
		if _, err := tx.Exec(ctx, `INSERT INTO public.schema_migrations(version,checksum) VALUES($1,$2)`, f.version, f.checksum); err != nil {
			return migrationFailure("bookkeeping insert failed; transaction rolled back", f.version)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return migrationFailure("commit failed; inspect status before retry (acknowledgement may be ambiguous)", f.version)
	}
	return nil
}
func (m *Migrator) Status(ctx context.Context) error {
	return m.locked(ctx, "status", func(conn *pgx.Conn, files []migrationFile) error {
		exists, err := trackingExists(ctx, conn)
		if err != nil {
			return err
		}
		records, err := m.state(ctx, conn, files, false)
		if err != nil {
			return err
		}
		if !exists {
			fmt.Fprintln(m.output, "migration table missing; no migrations applied")
		}
		latest := "none"
		for _, f := range files {
			state, checksum := "pending", "not-applied"
			if _, ok := records[f.version]; ok {
				state, checksum = "applied", "verified"
				latest = f.version
			}
			fmt.Fprintf(m.output, "%s %s checksum=%s\n", f.version, state, checksum)
		}
		fmt.Fprintf(m.output, "current version: %s\n", latest)
		return nil
	})
}

// Create is filesystem-only, deterministic and never overwrites a file. A local
// exclusive scaffold lock serializes paired generation; it is not the DB lock.
func (m *Migrator) Create(name string) (string, error) {
	name = sanitize(name)
	if name == "" || name[0] < 'a' || name[0] > 'z' {
		return "", migrationFailure("name must begin with a letter", "")
	}
	if err := os.MkdirAll(m.dir, 0755); err != nil {
		return "", migrationFailure("cannot create migration directory", "")
	}
	lockPath := filepath.Join(m.dir, ".scaffold.lock")
	lock, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return "", migrationFailure("scaffold lock unavailable; review any interrupted generation", "")
	}
	defer lock.Close()
	defer os.Remove(lockPath)
	files, err := m.discover()
	if err != nil {
		return "", err
	}
	next := 1
	if len(files) > 0 {
		next = files[len(files)-1].number + 1
	}
	if next > 999999 {
		return "", migrationFailure("six-digit version range exhausted", "")
	}
	base := fmt.Sprintf("%06d_%s", next, name)
	upPath := filepath.Join(m.dir, base+".up.sql")
	for _, direction := range []string{"up", "down"} {
		path := filepath.Join(m.dir, base+"."+direction+".sql")
		f, e := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if e != nil {
			if direction == "down" {
				_ = os.Remove(upPath)
			}
			return "", migrationFailure("cannot exclusively create migration pair", base)
		}
		_, e = f.WriteString("-- " + base + " " + direction + "\n\n")
		closeErr := f.Close()
		if e != nil || closeErr != nil {
			_ = os.Remove(path)
			_ = os.Remove(upPath)
			return "", migrationFailure("cannot write migration scaffold", base)
		}
	}
	fmt.Fprintf(m.output, "created migration %s\n", base)
	return base, nil
}
func sanitize(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if r == '_' || r == ' ' || r == '-' {
			b.WriteByte('_')
		}
	}
	return strings.Trim(regexp.MustCompile(`_+`).ReplaceAllString(b.String(), "_"), "_")
}
