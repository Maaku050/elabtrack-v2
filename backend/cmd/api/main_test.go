package main

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestCommandSelection(t *testing.T) {
	for _, args := range [][]string{{"--seed", "--migrate-up"}, {"--migrate-up", "--migrate-down"}, {"--migrate-create="}, {"--migrate-create", "!!!"}, {"--seed", "extra"}, {"--unknown=credential-sentinel"}} {
		_, err := parseCommand(args)
		if err == nil || strings.Contains(err.Error(), "credential-sentinel") {
			t.Fatal("invalid command must fail safely")
		}
	}
	for _, args := range [][]string{nil, {"--migrate-up"}, {"--migrate-down"}, {"--migrate-status"}, {"--migrate-create", "add_test"}, {"--seed"}} {
		if _, err := parseCommand(args); err != nil {
			t.Fatal(err)
		}
	}
}
func TestStartupValidationAndProductionSeedBeforeDatabase(t *testing.T) {
	// Isolate the real loader from optional settings in the developer's shell.
	for _, entry := range os.Environ() {
		key, value, _ := strings.Cut(entry, "=")
		controlled := key == "PORT" || key == "DATABASE_URL" || key == "FRONTEND_URL" || key == "ALLOWED_ORIGINS"
		for _, prefix := range []string{"APP_", "DB_", "JWT_", "RATE_LIMIT_", "LOG_", "PG"} {
			controlled = controlled || strings.HasPrefix(key, prefix)
		}
		if controlled {
			t.Cleanup(func() { _ = os.Setenv(key, value) })
			if err := os.Unsetenv(key); err != nil {
				t.Fatal(err)
			}
		}
	}
	for key, value := range map[string]string{"APP_ENV": "production", "JWT_SECRET": "change-me-in-production", "DB_HOST": "unreachable.example.invalid", "DB_PORT": "5432", "DB_NAME": "isolated_test", "DB_USER": "synthetic", "DB_PASSWORD": "synthetic", "DB_SSLMODE": "verify-full", "FRONTEND_URL": "https://app.example.invalid", "ALLOWED_ORIGINS": "https://app.example.invalid"} {
		t.Setenv(key, value)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := run(ctx, command{})
	if err == nil || !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Fatal("startup must validate before connecting")
	}
	t.Setenv("JWT_SECRET", "t9K2m7Q4v8R1z6N3p5W0x2C7b9H4f8L1")
	err = run(ctx, command{seed: true})
	if err == nil || err.Error() != "seed: permitted only in development" {
		t.Fatal("production seed must fail before connection")
	}
}
