package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/google/uuid"
	"io"
	"os"
	"testing"
)

// Uses the actual hardened runner and a separately named disposable database.
func TestRealProductRoleMigration(t *testing.T) {
	if os.Getenv("ELABTRACK_PHASE4A") != "1" {
		t.Skip("requires isolated Phase 4A PostgreSQL")
	}
	cfg, err := config.Parse(map[string]string{"APP_ENV": "test", "JWT_SECRET": os.Getenv("JWT_SECRET"), "DB_HOST": "127.0.0.1", "DB_PORT": "25432", "DB_NAME": "elabtrack_v2_phase4a_test", "DB_USER": "elabtrack_runtime", "DB_PASSWORD": os.Getenv("DB_PASSWORD"), "MIGRATION_DATABASE_URL": os.Getenv("MIGRATION_DATABASE_URL")})
	must(t, err == nil, "isolated config")
	ownerCfg, err := cfg.MigrationConnection()
	must(t, err == nil, "owner config")
	ctx := context.Background()
	owner, err := New(ctx, ownerCfg)
	must(t, err == nil, "owner connection")
	defer owner.Close()
	var safe bool
	err = owner.Pool.QueryRow(ctx, `SELECT current_database()='elabtrack_v2_phase4a_test' AND current_user='elabtrack_migrator' AND NOT rolsuper AND NOT rolcreaterole FROM pg_roles WHERE rolname=current_user`).Scan(&safe)
	must(t, err == nil && safe, "isolated least-privilege owner")
	full := NewMigrator(owner.Pool, "../../../migrations")
	full.output = io.Discard
	files, err := full.discover()
	must(t, err == nil && len(files) >= 4 && files[3].version == "000004_product_roles", "historical product-role migration prefix")
	// Exercise exactly the historical four pairs as new migrations are added.
	productDir := t.TempDir()
	for _, f := range files[:4] {
		fixturePair(t, productDir, f.version, string(f.up), string(f.down))
	}
	full = NewMigrator(owner.Pool, productDir)
	full.output = io.Discard
	legacyDir := t.TempDir()
	for _, f := range files[:3] {
		fixturePair(t, legacyDir, f.version, string(f.up), string(f.down))
	}
	legacy := NewMigrator(owner.Pool, legacyDir)
	legacy.output = io.Discard
	// Refuse a historical replay if later migrations are already tracked.
	var trackingExists bool
	err = owner.Pool.QueryRow(ctx, `SELECT to_regclass('schema_migrations') IS NOT NULL`).Scan(&trackingExists)
	must(t, err == nil, "tracking availability")
	if trackingExists {
		var newer int
		err = owner.Pool.QueryRow(ctx, `SELECT count(*) FROM schema_migrations WHERE version > '000004_product_roles'`).Scan(&newer)
		must(t, err == nil && newer == 0, "historical replay requires isolated migration prefix")
	}
	// Re-runs are allowed only on an empty disposable schema at version 4.
	var exists bool
	err = owner.Pool.QueryRow(ctx, `SELECT to_regclass('users') IS NOT NULL`).Scan(&exists)
	must(t, err == nil, "schema probe")
	if exists {
		var count int
		err = owner.Pool.QueryRow(ctx, `SELECT count(*) FROM users`).Scan(&count)
		must(t, err == nil && count == 0, "refuse migration replay with existing accounts")
		var v4 bool
		err = owner.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version='000004_product_roles')`).Scan(&v4)
		must(t, err == nil, "tracking probe")
		if v4 {
			must(t, full.Down(ctx) == nil, "empty product-role down")
		}
	}
	must(t, legacy.Up(ctx) == nil, "historical foundation up")
	grant, err := os.ReadFile("../../../database/runtime-grants.sql")
	must(t, err == nil, "explicit grant file")
	_, err = owner.Pool.Exec(ctx, string(grant))
	must(t, err == nil, "named runtime grants")
	ids := []uuid.UUID{uuid.New(), uuid.New()}
	staff := uuid.New()
	defer func() {
		_, e := owner.Pool.Exec(ctx, `DELETE FROM users WHERE id=ANY($1::uuid[])`, append(ids, staff))
		must(t, e == nil, "owned fixtures cleaned")
	}()
	for i, role := range []string{"user", "admin"} {
		_, err = owner.Pool.Exec(ctx, `INSERT INTO users(id,email,name,password,role,is_active) VALUES($1,$2,'Synthetic migration fixture','not-a-login',$3,$4)`, ids[i], "phase4a-migration-"+ids[i].String()+"@example.invalid", role, i == 1)
		must(t, err == nil, "legacy account fixture")
		_, err = owner.Pool.Exec(ctx, `INSERT INTO refresh_tokens(id,user_id,token_hash,expires_at)VALUES($1,$2,$3,now()+interval '1 day')`, uuid.New(), ids[i], digestFixture(ids[i].String()))
		must(t, err == nil, "legacy session fixture")
	}
	snapshot := func() string {
		var value string
		err := owner.Pool.QueryRow(ctx, `SELECT jsonb_build_object('users',(SELECT jsonb_agg(to_jsonb(u)-'role' ORDER BY id) FROM users u WHERE id=ANY($1::uuid[])),'sessions',(SELECT jsonb_agg(to_jsonb(r) ORDER BY id) FROM refresh_tokens r WHERE user_id=ANY($1::uuid[])))::text`, ids).Scan(&value)
		must(t, err == nil, "identity/session snapshot")
		return value
	}
	before := snapshot()
	must(t, full.Up(ctx) == nil, "product-role up")
	must(t, snapshot() == before, "identity/password/status/timestamps/sessions preserved")
	checkRoles := func(a, b string) {
		var ok bool
		err := owner.Pool.QueryRow(ctx, `SELECT (SELECT role=$2 FROM users WHERE id=$1) AND (SELECT role=$4 FROM users WHERE id=$3)`, ids[0], a, ids[1], b).Scan(&ok)
		must(t, err == nil && ok, "exact role mapping")
	}
	checkRoles("BORROWER", "ADMIN")
	runtime, err := New(ctx, cfg.DB)
	must(t, err == nil, "runtime connection")
	defer runtime.Close()
	for _, role := range []string{"user", "admin", "Student", "Faculty", "SUPER_ADMIN"} {
		_, err = runtime.Pool.Exec(ctx, `UPDATE users SET role=$2 WHERE id=$1`, ids[0], role)
		must(t, err != nil, "unapproved role denied")
	}
	for _, sql := range []string{`CREATE TABLE phase4a_forbidden(id int)`, `ALTER TABLE users ADD COLUMN phase4a_forbidden int`, `SELECT version FROM schema_migrations`, `CREATE ROLE phase4a_forbidden`} {
		_, err = runtime.Pool.Exec(ctx, sql)
		must(t, err != nil, "runtime DDL/tracking/cluster denied")
	}
	must(t, full.Down(ctx) == nil, "representable down")
	checkRoles("user", "admin")
	must(t, snapshot() == before, "rollback preserves all non-role data")
	must(t, full.Up(ctx) == nil, "product-role reapply")
	_, err = runtime.Pool.Exec(ctx, `INSERT INTO users(id,email,name,password,role)VALUES($1,$2,'Synthetic Staff fixture','not-a-login','STAFF')`, staff, "phase4a-migration-"+staff.String()+"@example.invalid")
	must(t, err == nil, "runtime Staff creation DML")
	must(t, full.Down(ctx) != nil, "Staff rollback must refuse ambiguous downgrade")
	checkRoles("BORROWER", "ADMIN")
	must(t, snapshot() == before && full.Status(ctx) == nil, "refused rollback leaves data and checksummed tracking intact")
	var intact bool
	err = owner.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1 AND role='STAFF') AND EXISTS(SELECT 1 FROM schema_migrations WHERE version='000004_product_roles')`, staff).Scan(&intact)
	must(t, err == nil && intact, "Staff and version retained")
	t.Log("user→BORROWER, admin→ADMIN; identity/session preservation; up/down/up; Staff rollback refusal; runtime DML/DDL separation verified")
}

func digestFixture(value string) string {
	h := sha256.Sum256([]byte(value))
	return hex.EncodeToString(h[:])
}
