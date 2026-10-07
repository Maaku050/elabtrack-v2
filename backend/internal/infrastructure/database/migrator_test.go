package database

import (
	"bytes"
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixturePair(t *testing.T, dir, version, up, down string) {
	t.Helper()
	for _, p := range []struct{ suffix, sql string }{{"up", up}, {"down", down}} {
		if err := os.WriteFile(filepath.Join(dir, version+"."+p.suffix+".sql"), []byte(p.sql), 0600); err != nil {
			t.Fatal(err)
		}
	}
}
func TestMigrationDiscovery(t *testing.T) {
	for _, c := range []struct {
		name, version     string
		orphan, duplicate bool
	}{{"invalid", "1_bad", false, false}, {"zero", "000000_zero", false, false}, {"counterpart", "000001_valid", true, false}, {"duplicate", "000001_first", false, true}} {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			fixturePair(t, dir, c.version, "SELECT 1;", "SELECT 1;")
			if c.orphan {
				_ = os.Remove(filepath.Join(dir, c.version+".down.sql"))
			}
			if c.duplicate {
				fixturePair(t, dir, "000001_other", "SELECT 1;", "SELECT 1;")
			}
			if _, err := NewMigrator(nil, dir).discover(); err == nil {
				t.Fatal("invalid mapping accepted")
			}
		})
	}
	dir := t.TempDir()
	fixturePair(t, dir, "000010_later", "SELECT 10;", "SELECT 10;")
	fixturePair(t, dir, "000002_earlier", "SELECT 2;", "SELECT 2;")
	files, err := NewMigrator(nil, dir).discover()
	if err != nil || files[0].number != 2 || files[1].number != 10 {
		t.Fatal("numeric ordering/gaps")
	}
}
func TestMigrationHistory(t *testing.T) {
	files := []migrationFile{{version: "000001_first", checksum: "one"}, {version: "000010_second", checksum: "two"}}
	for _, c := range []struct {
		name    string
		records map[string]migrationRecord
		fail    bool
	}{{"empty", map[string]migrationRecord{}, false}, {"prefix", map[string]migrationRecord{"000001_first": {checksum: "one"}}, false}, {"mismatch", map[string]migrationRecord{"000001_first": {checksum: "wrong"}}, true}, {"gap", map[string]migrationRecord{"000010_second": {checksum: "two"}}, true}, {"unknown", map[string]migrationRecord{"unknown-secret-sentinel": {checksum: "whatever"}}, true}} {
		t.Run(c.name, func(t *testing.T) {
			err := validateHistory(files, c.records)
			if (err != nil) != c.fail {
				t.Fatal("history decision")
			}
			if err != nil && strings.Contains(err.Error(), "sentinel") {
				t.Fatal("database version leaked")
			}
		})
	}
}
func TestMigrationSQLTransactionBoundary(t *testing.T) {
	for _, sql := range []string{"COMMIT;", "CREATE TABLE probe(id int); /*nested /* comment */ */ END;", "--comment\nABORT;", "START TRANSACTION;", "PREPARE TRANSACTION 'x';", "SET standard_conforming_strings=off;", "ROLLBACK;", "SELECT 'backslash\\';COMMIT;"} {
		if validateTransactionalSQL(sql) == nil {
			t.Fatal("transaction/session control accepted")
		}
	}
	for _, sql := range []string{"-- COMMIT\n SELECT 'COMMIT;';", `DO $body$ BEGIN RAISE NOTICE 'END;'; END $body$;`, `SELECT E'escaped\'quote;'; SELECT 1;`, `SELECT "COMMIT" FROM example;`, "/* /* COMMIT */ */ SELECT 1;"} {
		if err := validateTransactionalSQL(sql); err != nil {
			t.Fatal(err)
		}
	}
}
func TestMigrationPairChecksum(t *testing.T) {
	if pairChecksum([]byte("ab"), []byte("c")) == pairChecksum([]byte("a"), []byte("bc")) {
		t.Fatal("ambiguous pair")
	}
	if pairChecksum([]byte("up"), []byte("down")) == pairChecksum([]byte("up"), []byte("down\n")) {
		t.Fatal("down unprotected")
	}
	files, err := NewMigrator(nil, "../../../migrations").discover()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if legacyFoundationChecksums[f.version] != f.checksum {
			t.Fatal("legacy baseline changed")
		}
	}
}
func TestMigrationGeneration(t *testing.T) {
	dir := t.TempDir()
	m := NewMigrator(nil, dir)
	m.output = &bytes.Buffer{}
	fixturePair(t, dir, "000010_before", "SELECT 1;", "SELECT 1;")
	v, err := m.Create("Next pair")
	if err != nil || v != "000011_next_pair" {
		t.Fatal("deterministic version")
	}
	if _, err := m.discover(); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Create("../"); err == nil {
		t.Fatal("unsafe name")
	}
	_ = os.WriteFile(filepath.Join(dir, ".scaffold.lock"), nil, 0600)
	if _, err := m.Create("another"); err == nil {
		t.Fatal("lock ignored")
	}
}
func TestMigrationSafeFailureLogging(t *testing.T) {
	dir := t.TempDir()
	fixturePair(t, dir, "000001_safe", "SELECT 'sql-credential-sentinel';", "SELECT 1;")
	core, logs := observer.New(zap.InfoLevel)
	err := NewMigrator(nil, dir, zap.New(core)).Up(context.Background())
	if err == nil {
		t.Fatal("no DB accepted")
	}
	entries := logs.All()
	if len(entries) != 1 {
		t.Fatal("missing log")
	}
	fields := entries[0].ContextMap()
	if fields["result"] != "failure" || fields["event"] != "migration.command" || fields["direction"] != "up" {
		t.Fatal("metadata")
	}
	if strings.Contains(err.Error(), "sentinel") {
		t.Fatal("SQL leaked")
	}
}
