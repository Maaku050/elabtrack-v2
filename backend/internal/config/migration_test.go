package config

import (
	"strings"
	"testing"
)

func TestMigrationConfiguration(t *testing.T) {
	cfg, err := Parse(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	db, err := cfg.MigrationConnection()
	if err != nil || db.User != cfg.DB.User {
		t.Fatal("development fallback")
	}
	values := productionValues()
	cfg, err = Parse(values)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = cfg.MigrationConnection(); err == nil {
		t.Fatal("production fallback")
	}
	values["MIGRATION_DATABASE_URL"] = "postgresql://schema_owner:credential-sentinel@db.example.invalid:5432/" + values["DB_NAME"] + "?sslmode=verify-full"
	cfg, err = Parse(values)
	if err != nil {
		t.Fatal(err)
	}
	migration, err := cfg.MigrationConnection()
	if err != nil {
		t.Fatal(err)
	}
	pool, err := migration.PoolConfig()
	if err != nil || pool.ConnConfig.TLSConfig == nil || pool.ConnConfig.TLSConfig.InsecureSkipVerify || pool.ConnConfig.TLSConfig.ServerName != "db.example.invalid" || pool.ConnConfig.Fallbacks != nil {
		t.Fatal("migration TLS weakened")
	}
	for _, url := range []string{"", "postgresql://owner:credential-sentinel@db.example.invalid/" + values["DB_NAME"] + "?sslmode=disable", "postgresql://owner:credential-sentinel@different.invalid/" + values["DB_NAME"] + "?sslmode=verify-full", "postgresql://" + values["DB_USER"] + ":credential-sentinel@db.example.invalid/" + values["DB_NAME"] + "?sslmode=verify-full"} {
		values["MIGRATION_DATABASE_URL"] = url
		if _, err := Parse(values); err == nil || strings.Contains(err.Error(), "credential-sentinel") {
			t.Fatal("unsafe/leaked configuration")
		}
	}
	cfg, err = Parse(map[string]string{"APP_ENV": "test", "JWT_SECRET": "test-only"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cfg.MigrationConnection(); err != nil {
		t.Fatal("test fallback")
	}
}
