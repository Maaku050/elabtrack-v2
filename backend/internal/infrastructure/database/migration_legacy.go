package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

// These are the immutable repository foundation pairs verified in Phase 1G.
// Adoption is an operator attestation of a known local pre-production database,
// never proof that arbitrary historical DDL matches the current files.
var legacyFoundationChecksums = map[string]string{
	"000001_create_users":             "b5194dfa5b052b8af2c1453c00183f5048d6685d0ea6758b25b4f6bab8e81bae",
	"000002_create_refresh_tokens":    "6d55709989ae2e70960d4198e9ba7cbc0849bbeaef85db6f7b1d5782d30bde30",
	"000003_refresh_session_security": "baa3df61b17a57a2c4b3f97a09ce63967bd4059d5b5d7fd5b9e7e036ab42b0e0",
}

func (m *Migrator) AdoptLegacy(ctx context.Context) error {
	return m.locked(ctx, "adopt-legacy", func(conn *pgx.Conn, files []migrationFile) error {
		exists, err := trackingExists(ctx, conn)
		if err != nil {
			return err
		}
		if !exists {
			return migrationFailure("legacy tracking table missing", "")
		}
		var checksummed bool
		if err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pg_attribute WHERE attrelid='public.schema_migrations'::regclass AND attname='checksum' AND NOT attisdropped)`).Scan(&checksummed); err != nil {
			return migrationFailure("legacy metadata query failed", "")
		}
		if checksummed {
			return migrationFailure("tracking already has checksums; adoption refused", "")
		}
		known := map[string]migrationFile{}
		for _, f := range files {
			expected, ok := legacyFoundationChecksums[f.version]
			if ok {
				if expected != f.checksum {
					return migrationFailure("legacy foundation files differ from verified baseline", f.version)
				}
				known[f.version] = f
			}
		}
		rows, err := conn.Query(ctx, `SELECT version, applied_at FROM public.schema_migrations ORDER BY version`)
		if err != nil {
			return migrationFailure("legacy history query failed", "")
		}
		records := map[string]migrationRecord{}
		for rows.Next() {
			var version string
			var applied time.Time
			if err := rows.Scan(&version, &applied); err != nil {
				rows.Close()
				return migrationFailure("legacy history scan failed", "")
			}
			f, ok := known[version]
			if !ok {
				rows.Close()
				return migrationFailure("unknown legacy version; adoption refused", "")
			}
			if _, duplicate := records[version]; duplicate {
				rows.Close()
				return migrationFailure("duplicate legacy version; adoption refused", "")
			}
			records[version] = migrationRecord{checksum: f.checksum, applied: applied}
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return migrationFailure("legacy history iteration failed", "")
		}
		if err := validateHistory(files, records); err != nil {
			return err
		}
		tx, err := conn.Begin(ctx)
		if err != nil {
			return migrationFailure("legacy transaction begin failed", "")
		}
		defer func() {
			cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = tx.Rollback(cleanup)
		}()
		if _, err := tx.Exec(ctx, `ALTER TABLE public.schema_migrations ADD COLUMN checksum TEXT`); err != nil {
			return migrationFailure("legacy tracking upgrade failed", "")
		}
		for version, record := range records {
			if _, err := tx.Exec(ctx, `UPDATE public.schema_migrations SET checksum=$2 WHERE version=$1 AND checksum IS NULL`, version, record.checksum); err != nil {
				return migrationFailure("legacy checksum adoption failed", version)
			}
		}
		if _, err := tx.Exec(ctx, `ALTER TABLE public.schema_migrations ALTER COLUMN checksum SET NOT NULL, ADD CONSTRAINT schema_migrations_checksum_check CHECK(checksum ~ '^[0-9a-f]{64}$')`); err != nil {
			return migrationFailure("legacy tracking constraints failed", "")
		}
		if err := tx.Commit(ctx); err != nil {
			return migrationFailure("legacy commit failed; inspect status before retry", "")
		}
		return nil
	})
}
