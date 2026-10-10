package config

import (
	"crypto/tls"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Synthetic material only; never use this fixture as a deployment credential.
const testKey = "t9K2m7Q4v8R1z6N3p5W0x2C7b9H4f8L1"

// pgx reads PG* internally when initializing its defaults. Keep fixtures
// independent of the developer's shell; individual tests opt in explicitly.
func TestMain(m *testing.M) {
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "PG") {
			_ = os.Unsetenv(key)
		}
	}
	os.Exit(m.Run())
}

func productionValues() map[string]string {
	return map[string]string{"APP_ENV": "production", "JWT_SECRET": testKey, "DB_HOST": "db.example.invalid", "DB_NAME": "isolated_test", "DB_USER": "test_user", "DB_PASSWORD": "synthetic-db-value", "DB_SSLMODE": "verify-full", "FRONTEND_URL": "https://app.example.invalid", "ALLOWED_ORIGINS": "https://app.example.invalid"}
}
func TestEnvironmentConfiguration(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values map[string]string
		want   Environment
	}{
		{"development", map[string]string{}, Development},
		{"test", map[string]string{"APP_ENV": "test", "JWT_SECRET": testKey, "DB_NAME": "elabtrack_v2_test"}, Test},
		{"production", productionValues(), Production},
	} {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Parse(tt.values)
			if err != nil {
				t.Fatal(err)
			}
			if cfg.App.Env != tt.want || cfg.App.BodyLimit != 1<<20 {
				t.Fatal("incorrect typed config")
			}
			if tt.want == Development && cfg.JWT.Secret != DevelopmentSecret {
				t.Fatal("incorrect development default")
			}
		})
	}
}
func TestProductionSecretRejection(t *testing.T) {
	for name, value := range map[string]string{"blank": "", "whitespace": "   ", "old default": "change-me-in-production", "development": DevelopmentSecret, "short": "abcdefghijkLMNOP", "placeholder": "replace-with-your-production-secret-value", "repetition": strings.Repeat("a", 64), "low diversity": strings.Repeat("abcd", 16), "alphabet": "abcdefghijklmnopqrstuvwxyzABCDEF", "repeated diverse pattern": strings.Repeat("a1B2c3D4e5F6g7H8", 4)} {
		t.Run(name, func(t *testing.T) {
			values := productionValues()
			values["JWT_SECRET"] = value
			_, err := Parse(values)
			if err == nil || !strings.Contains(err.Error(), "JWT_SECRET") {
				t.Fatal("must reject secret")
			}
			if value != "" && strings.Contains(err.Error(), value) {
				t.Fatal("secret exposed")
			}
		})
	}
	values := productionValues()
	delete(values, "JWT_SECRET")
	if _, err := Parse(values); err == nil {
		t.Fatal("missing production signing material accepted")
	}
}
func TestInvalidSettingsFailClosed(t *testing.T) {
	for _, tt := range []struct{ key, value string }{
		{"APP_ENV", "staging"}, {"APP_ENV", ""}, {"APP_PORT", "not-a-port"}, {"APP_PORT", "0"}, {"APP_PORT", "65536"},
		{"APP_READ_TIMEOUT", "invalid"}, {"APP_WRITE_TIMEOUT", "0s"}, {"APP_BODY_LIMIT", "bogus"}, {"APP_BODY_LIMIT", "999999999999999GB"},
		{"DB_PORT", "-1"}, {"DB_PORT", "+5432"}, {"DB_MAX_CONNS", "0"}, {"DB_MIN_CONNS", "21"}, {"DB_MAX_CONN_IDLE_TIME", "-1m"},
		{"JWT_ACCESS_TTL", "0s"}, {"JWT_REFRESH_TTL", "1m"}, {"JWT_ISSUER", ""}, {"RATE_LIMIT_MAX", "bad"},
		{"RATE_LIMIT_WINDOW", "0s"}, {"LOG_LEVEL", "unknown"}, {"LOG_FORMAT", "unknown"}, {"ALLOWED_ORIGINS", "*"},
	} {
		t.Run(tt.key+tt.value, func(t *testing.T) {
			_, err := Parse(map[string]string{tt.key: tt.value})
			if err == nil || !strings.Contains(err.Error(), tt.key) {
				t.Fatalf("expected setting name %s", tt.key)
			}
		})
	}
	if _, err := Parse(map[string]string{"APP_PORT": "8080", "PORT": "8081"}); err == nil {
		t.Fatal("conflicting port accepted")
	}
	cfg, err := Parse(map[string]string{"PORT": "9090"})
	if err != nil || cfg.App.Port != "9090" {
		t.Fatal("PORT fallback failed")
	}
	if _, err := Parse(map[string]string{"APP_ENV": "test"}); err == nil {
		t.Fatal("test must supply explicit signing material")
	}
}
func TestProductionTransportPolicy(t *testing.T) {
	for _, mode := range []string{"", "disable", "allow", "prefer", "require", "verify-ca"} {
		t.Run(mode, func(t *testing.T) {
			values := productionValues()
			values["DB_SSLMODE"] = mode
			_, err := Parse(values)
			if err == nil || !strings.Contains(err.Error(), "DB_SSLMODE") {
				t.Fatal("unsafe TLS accepted")
			}
		})
	}
	for _, key := range []string{"DB_SSLMODE", "DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "FRONTEND_URL", "ALLOWED_ORIGINS"} {
		t.Run("missing "+key, func(t *testing.T) {
			values := productionValues()
			delete(values, key)
			if _, err := Parse(values); err == nil {
				t.Fatal("missing production setting accepted")
			}
		})
	}
	for _, key := range []string{"FRONTEND_URL", "ALLOWED_ORIGINS"} {
		values := productionValues()
		values[key] = "http://app.example.invalid"
		if _, err := Parse(values); err == nil {
			t.Fatal("insecure production origin accepted")
		}
	}
	values := productionValues()
	values["DB_PASSWORD"] = " "
	if _, err := Parse(values); err == nil {
		t.Fatal("blank production password accepted")
	}
	values = productionValues()
	values["DB_SSLROOTCERT"] = filepath.Join(t.TempDir(), "missing.pem")
	if _, err := Parse(values); err == nil {
		t.Fatal("missing certificate accepted")
	}
}
func TestCredentialEncodingAndTLS(t *testing.T) {
	values := productionValues()
	user, password, dbname := "user:@/% #&?+é", "p%:@/#&?='\\+ 空格", "database/% #?"
	values["DB_USER"] = user
	values["DB_PASSWORD"] = password
	values["DB_NAME"] = dbname
	for _, useURL := range []bool{false, true} {
		if useURL {
			u := url.URL{Scheme: "postgresql", Host: "db.example.invalid:5432", Path: "/" + dbname, User: url.UserPassword(user, password), RawQuery: "sslmode=verify-full"}
			for _, key := range []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_SSLMODE"} {
				delete(values, key)
			}
			values["DATABASE_URL"] = u.String()
		}
		cfg, err := Parse(values)
		if err != nil {
			t.Fatal(err)
		}
		pool, err := cfg.DB.PoolConfig()
		if err != nil {
			t.Fatal(err)
		}
		if pool.ConnConfig.User != user || pool.ConnConfig.Password != password || pool.ConnConfig.Database != dbname {
			t.Fatal("credentials altered")
		}
		tlsCfg := pool.ConnConfig.TLSConfig
		if tlsCfg == nil || tlsCfg.InsecureSkipVerify || tlsCfg.ServerName != "db.example.invalid" || tlsCfg.MinVersion < tls.VersionTLS12 || len(pool.ConnConfig.Fallbacks) > 0 {
			t.Fatal("TLS verification or downgrade policy broken")
		}
	}
	cfg, err := Parse(map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	pool, err := cfg.DB.PoolConfig()
	if err != nil {
		t.Fatal(err)
	}
	if pool.ConnConfig.TLSConfig != nil {
		t.Fatal("local TLS unexpectedly enabled")
	}
}
func TestURLRejectionAndRedaction(t *testing.T) {
	password := "synthetic-sensitive-sentinel"
	for _, tail := range []string{"?sslmode=disable", "?sslmode=verify-full&sslmode=disable", "?sslmode=verify-full&host=other", "?sslmode=verify-full&ssl=true", "?sslmode=verify-full&bad=%GG", ""} {
		raw := "postgresql://u:" + password + "@db.example.invalid/db" + tail
		values := productionValues()
		for _, key := range []string{"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_SSLMODE"} {
			delete(values, key)
		}
		values["DATABASE_URL"] = raw
		_, err := Parse(values)
		if err == nil {
			t.Fatal("invalid URL accepted")
		}
		if strings.Contains(err.Error(), password) || strings.Contains(err.Error(), raw) || strings.Contains(err.Error(), testKey) {
			t.Fatal("credential leaked")
		}
	}
	values := productionValues()
	values["DATABASE_URL"] = "postgresql://u:p@db.example.invalid/db?sslmode=verify-full"
	if _, err := Parse(values); err == nil {
		t.Fatal("mixed URL/discrete config accepted")
	}
	values = productionValues()
	for _, key := range []string{"DB_HOST", "DB_USER", "DB_NAME", "DB_PASSWORD"} {
		delete(values, key)
	}
	values["DATABASE_URL"] = "postgresql://u:p@db.example.invalid/db?sslmode=disable"
	if _, err := Parse(values); err == nil {
		t.Fatal("conflicting TLS accepted")
	}
	values = productionValues()
	for _, key := range []string{"DB_HOST", "DB_USER", "DB_NAME", "DB_PASSWORD", "DB_SSLMODE"} {
		delete(values, key)
	}
	values["DATABASE_URL"] = "postgresql://bad%GG:" + password + "@db/db"
	_, err := Parse(values)
	if err == nil || strings.Contains(err.Error(), password) {
		t.Fatal("malformed URL must fail safely")
	}
}
func TestAmbientPostgresSettingsRejected(t *testing.T) {
	if _, err := Parse(map[string]string{"PGPASSWORD": "sensitive-sentinel"}); err == nil || strings.Contains(err.Error(), "sensitive-sentinel") {
		t.Fatal("ambient settings must fail safely")
	}
	t.Setenv("PGSSLMODE", "disable")
	t.Setenv("PGHOST", "attacker.example.invalid")
	t.Setenv("PGPASSWORD", "attacker")
	t.Setenv("PGOPTIONS", "-c statement_timeout=1")
	cfg, err := Parse(productionValues())
	if err != nil {
		t.Fatal(err)
	}
	pool, err := cfg.DB.PoolConfig()
	if err != nil {
		t.Fatal(err)
	}
	if pool.ConnConfig.Host != "db.example.invalid" || pool.ConnConfig.Password != "synthetic-db-value" || pool.ConnConfig.TLSConfig == nil || len(pool.ConnConfig.RuntimeParams) != 2 || pool.ConnConfig.RuntimeParams["standard_conforming_strings"] != "on" {
		t.Fatal("ambient configuration won")
	}
}
func TestDotenvPrecedenceAndIsolation(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("APP_PORT=9090\nJWT_SECRET='local-example'\n"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadWithFile(map[string]string{"APP_PORT": "8081"}, path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.App.Port != "8081" || cfg.JWT.Secret != "local-example" {
		t.Fatal("dotenv precedence failed")
	}
	if err := os.WriteFile(path, []byte("invalid-sensitive-line\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadWithFile(map[string]string{}, path); err == nil || strings.Contains(err.Error(), "sensitive") {
		t.Fatal("malformed dotenv must fail without contents")
	}
	if _, err := loadWithFile(productionValues(), path); err != nil {
		t.Fatal("production must ignore dotenv", err)
	}
	if _, err := loadWithFile(map[string]string{"APP_ENV": "test", "JWT_SECRET": testKey}, path); err != nil {
		t.Fatal("test must ignore dotenv", err)
	}
	if _, err := loadWithFile(map[string]string{"APP_ENV": "unknown"}, path); err == nil || !strings.Contains(err.Error(), "APP_ENV") {
		t.Fatal("unknown env must fail")
	}
	if _, err := loadWithFile(map[string]string{"JWT_SECRET": ""}, filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatal("blank must not fall back")
	}
}

func TestAccountConfigurationGuards(t *testing.T) {
	for _, pair := range [][2]string{{"STUDENT_EMAIL_DOMAINS", "*.example.invalid"}, {"ACTIVATION_URL", "https://other.example.invalid/activate"}, {"ACTIVATION_URL", "http://localhost:5173/activate?token=bad"}, {"ACTIVATION_TTL", "1m"}, {"ACTIVATION_RESEND_COOLDOWN", "1s"}, {"BREVO_SENDER_EMAIL", "Display <sender@example.invalid>"}} {
		_, err := Parse(map[string]string{pair[0]: pair[1]})
		if err == nil {
			t.Fatal("unsafe account configuration accepted", pair[0])
		}
	}
	cfg, err := Parse(map[string]string{})
	if err != nil || len(cfg.Accounts.Policy.StudentDomains) != 0 || cfg.Accounts.BrevoKey != "" {
		t.Fatal("missing external dependencies must not stop independent startup")
	}
}

func TestExplicitListenerAddress(t *testing.T) {
	for _, v := range []struct{ host, address string }{{"", ":8080"}, {"127.0.0.1", "127.0.0.1:8080"}, {"::1", "[::1]:8080"}} {
		cfg, err := Parse(map[string]string{"APP_BIND_HOST": v.host})
		if err != nil || cfg.App.Address() != v.address {
			t.Fatalf("explicit listener configuration: %v", err)
		}
	}
	for _, host := range []string{"localhost", "127.0.0.1:8080", "https://example.invalid", "127.0.0.999"} {
		if _, err := Parse(map[string]string{"APP_BIND_HOST": host}); err == nil {
			t.Fatal("invalid bind address accepted")
		}
	}
}
