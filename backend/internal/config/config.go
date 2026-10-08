package config

import (
	"errors"
	"fmt"
	"net/netip"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type Environment string

const (
	Development       Environment = "development"
	Test              Environment = "test"
	Production        Environment = "production"
	DevelopmentSecret             = "elabtrack-local-development-only-signing-key"
)

// Config is validated before any infrastructure or HTTP listener is created.
type Config struct {
	App         AppConfig
	DB          DBConfig
	MigrationDB *DBConfig
	JWT         JWTConfig
	Security    SecurityConfig
	Log         LogConfig
	Accounts    AccountsConfig
}
type AppConfig struct {
	Env                       Environment
	Name, Port                string
	ReadTimeout, WriteTimeout time.Duration
	IdleTimeout               time.Duration
	BodyLimit                 int
}
type JWTConfig struct {
	Secret                string
	AccessTTL, RefreshTTL time.Duration
	Issuer                string
}
type SecurityConfig struct {
	FrontendURL                                                  string
	AllowedOrigins                                               []string
	RateLimitMax                                                 int
	RateLimitWindow                                              time.Duration
	LoginRateLimitMax, RefreshRateLimitMax, RegisterRateLimitMax int
	TrustedProxies                                               []netip.Prefix
}
type LogConfig struct{ Level, Format string }

func (c *Config) IsProduction() bool  { return c.App.Env == Production }
func (c *Config) IsDevelopment() bool { return c.App.Env == Development }

// Parse is deterministic: only the supplied settings are used. Errors contain
// setting names and policy descriptions, never supplied values.
func Parse(values map[string]string) (*Config, error) {
	r := reader{values: values}
	env := Environment(r.value("APP_ENV", string(Development)))
	if env != Development && env != Test && env != Production {
		r.fail("APP_ENV", "must be development, test or production")
	}
	for key, value := range values {
		if strings.HasPrefix(key, "PG") && value != "" {
			r.fail(key, "unsupported PostgreSQL setting; use DATABASE_URL or DB_* instead of PG* settings")
			break
		}
	}
	port := r.value("APP_PORT", r.value("PORT", "8080"))
	if a, ok := values["APP_PORT"]; ok {
		if b, present := values["PORT"]; present && a != b {
			r.fail("APP_PORT / PORT", "conflicting port settings")
		}
	}
	validatePort(&r, "APP_PORT / PORT", port)
	body := strings.ToUpper(r.value("APP_BODY_LIMIT", "1MB"))
	multiplier := int64(1)
	for suffix, size := range map[string]int64{"KB": 1 << 10, "MB": 1 << 20, "GB": 1 << 30} {
		if strings.HasSuffix(body, suffix) {
			body = strings.TrimSuffix(body, suffix)
			multiplier = size
			break
		}
	}
	n, err := strconv.ParseInt(strings.TrimSpace(body), 10, 64)
	if err != nil || n <= 0 || n > (1<<30)/multiplier {
		r.fail("APP_BODY_LIMIT", "must be positive bytes, KB, MB or GB, at most 1GB")
		n = 1
	}
	defaultSecret := ""
	if env == Development {
		defaultSecret = DevelopmentSecret
	}
	secret := r.value("JWT_SECRET", defaultSecret)
	if env == Production {
		if !strongSecret(secret) {
			r.fail("JWT_SECRET", "production requires explicit non-placeholder signing material of at least 32 bytes and sufficient diversity")
		}
	} else if strings.TrimSpace(secret) == "" {
		r.fail("JWT_SECRET", "must be explicit in test and nonblank in development")
	}
	frontendDefault, originsDefault, logFormat := "http://localhost:5173", "http://localhost:5173,http://localhost:4173", "console"
	if env == Production {
		frontendDefault, originsDefault, logFormat = "", "", "json"
	}
	frontend := readOrigin(&r, "FRONTEND_URL", r.value("FRONTEND_URL", frontendDefault), env)
	origins := strings.Split(r.value("ALLOWED_ORIGINS", originsDefault), ",")
	for i := range origins {
		origins[i] = readOrigin(&r, "ALLOWED_ORIGINS", origins[i], env)
	}
	frontendAllowed := false
	for _, origin := range origins {
		if origin == frontend {
			frontendAllowed = true
		}
	}
	if !frontendAllowed {
		r.fail("ALLOWED_ORIGINS", "must include FRONTEND_URL")
	}
	cfg := &Config{
		App:      AppConfig{Env: env, Name: r.nonblank("APP_NAME", "eLabTrack V2"), Port: port, ReadTimeout: r.duration("APP_READ_TIMEOUT", 10*time.Second), WriteTimeout: r.duration("APP_WRITE_TIMEOUT", 15*time.Second), IdleTimeout: r.duration("APP_IDLE_TIMEOUT", 60*time.Second), BodyLimit: int(n * multiplier)},
		JWT:      JWTConfig{Secret: secret, AccessTTL: r.duration("JWT_ACCESS_TTL", 15*time.Minute), RefreshTTL: r.duration("JWT_REFRESH_TTL", 168*time.Hour), Issuer: r.nonblank("JWT_ISSUER", "elabtrack-v2")},
		Security: SecurityConfig{FrontendURL: frontend, AllowedOrigins: origins, RateLimitMax: r.integer("RATE_LIMIT_MAX", 120, 1, 1<<31-1), RateLimitWindow: r.duration("RATE_LIMIT_WINDOW", time.Minute), LoginRateLimitMax: r.integer("LOGIN_RATE_LIMIT_MAX", 10, 1, 1<<31-1), RefreshRateLimitMax: r.integer("REFRESH_RATE_LIMIT_MAX", 60, 1, 1<<31-1), RegisterRateLimitMax: r.integer("REGISTER_RATE_LIMIT_MAX", 5, 1, 1<<31-1), TrustedProxies: readTrustedProxies(&r)},
		Log:      LogConfig{Level: r.value("LOG_LEVEL", "info"), Format: r.value("LOG_FORMAT", logFormat)},
	}
	if cfg.Security.RateLimitWindow < time.Second || cfg.Security.RateLimitWindow > time.Hour {
		r.fail("RATE_LIMIT_WINDOW", "must be between one second and one hour")
	}
	if cfg.JWT.RefreshTTL <= cfg.JWT.AccessTTL {
		r.fail("JWT_REFRESH_TTL", "must exceed JWT_ACCESS_TTL")
	}
	switch cfg.Log.Level {
	case "debug", "info", "warn", "error":
	default:
		r.fail("LOG_LEVEL", "must be debug, info, warn or error")
	}
	switch cfg.Log.Format {
	case "console", "json":
	default:
		r.fail("LOG_FORMAT", "must be console or json")
	}
	cfg.Accounts = parseAccounts(&r, env, frontend)
	cfg.DB = parseDatabase(&r, env)
	cfg.MigrationDB = parseMigrationDatabase(&r, env, cfg.DB)
	if len(r.errs) > 0 {
		return nil, errors.Join(r.errs...)
	}
	return cfg, nil
}

func strongSecret(s string) bool {
	if len(s) < 32 || strings.TrimSpace(s) != s {
		return false
	}
	lower := strings.ToLower(s)
	for _, marker := range []string{"change-me", "changeme", "replace", "placeholder", "example", "development", "password", "secret", "your-", "your_"} {
		if strings.Contains(lower, marker) {
			return false
		}
	}
	// Reject common keyboard/alphabet runs and repeated patterns even when
	// their length and character count look adequate.
	for _, run := range []string{"abcdefghijklmnopqrstuvwxyz", "0123456789", "qwertyuiop", "asdfghjkl"} {
		if strings.Contains(lower, run) {
			return false
		}
	}
	for size := 1; size <= len(s)/2 && size <= 32; size++ {
		if len(s)%size == 0 && strings.Repeat(s[:size], len(s)/size) == s {
			return false
		}
	}
	distinct := map[rune]bool{}
	for _, c := range s {
		if unicode.IsSpace(c) || unicode.IsControl(c) {
			return false
		}
		distinct[c] = true
	}
	// This rejects obvious repetition; it cannot establish cryptographic entropy.
	return len(distinct) >= 12
}

func readOrigin(r *reader, key, s string, env Environment) string {
	origin, err := CanonicalOrigin(strings.TrimSpace(s))
	if err != nil {
		r.fail(key, "must be an exact HTTP(S) origin without credentials, path or wildcard")
		return ""
	}
	if env == Production && !strings.HasPrefix(origin, "https://") {
		r.fail(key, "production requires explicit HTTPS origins")
	}
	return origin
}
func validatePort(r *reader, key, value string) {
	n, err := strconv.Atoi(value)
	decimal := value != ""
	for _, c := range value {
		if c < '0' || c > '9' {
			decimal = false
		}
	}
	if err != nil || !decimal || n < 1 || n > 65535 {
		r.fail(key, "must be a port from 1 to 65535")
	}
}

type reader struct {
	values map[string]string
	errs   []error
}

func (r *reader) fail(key, policy string) { r.errs = append(r.errs, fmt.Errorf("%s: %s", key, policy)) }
func (r *reader) value(key, fallback string) string {
	if v, ok := r.values[key]; ok {
		return v
	}
	return fallback
}
func (r *reader) nonblank(key, fallback string) string {
	v := r.value(key, fallback)
	if strings.TrimSpace(v) == "" {
		r.fail(key, "must not be blank")
	}
	return v
}
func (r *reader) integer(key string, fallback, min, max int) int {
	v, err := strconv.Atoi(r.value(key, strconv.Itoa(fallback)))
	if err != nil || v < min || v > max {
		r.fail(key, "integer outside permitted range")
		return fallback
	}
	return v
}
func (r *reader) duration(key string, fallback time.Duration) time.Duration {
	v, err := time.ParseDuration(r.value(key, fallback.String()))
	if err != nil || v <= 0 {
		r.fail(key, "must be a positive Go duration")
		return fallback
	}
	return v
}
