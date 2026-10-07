package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBConfig contains validated connection settings, never a manually escaped DSN.
// Keep it out of logs: User/Password may contain credentials.
type DBConfig struct {
	Host, Port, Name, User, Password, SSLMode string
	MaxConns, MinConns                        int32
	MaxConnLifetime, MaxConnIdleTime          time.Duration
	tlsConfig                                 *tls.Config
}

func parseDatabase(r *reader, env Environment) DBConfig {
	required := func(key, fallback string) string {
		if env == Production {
			fallback = ""
		}
		return r.nonblank(key, fallback)
	}
	var c DBConfig
	source := "DB_PASSWORD"
	root := r.value("DB_SSLROOTCERT", "")
	modeDefault := "disable"
	if env == Production {
		modeDefault = ""
	}
	c.SSLMode = r.value("DB_SSLMODE", modeDefault)
	if raw, present := r.values["DATABASE_URL"]; present {
		source = "DATABASE_URL"
		for _, key := range []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD"} {
			if _, ok := r.values[key]; ok {
				r.fail("DATABASE_URL / DB_*", "choose URL or discrete connection credentials")
				break
			}
		}
		u, err := url.Parse(raw)
		if err != nil || u == nil || (u.Scheme != "postgres" && u.Scheme != "postgresql") || u.User == nil || u.Hostname() == "" || u.Fragment != "" || u.Opaque != "" {
			r.fail("DATABASE_URL", "must be a PostgreSQL URL with explicit host, user and database")
		} else {
			c.Host = u.Hostname()
			c.Port = u.Port()
			if c.Port == "" {
				c.Port = "5432"
			}
			c.User = u.User.Username()
			c.Password, _ = u.User.Password()
			c.Name = strings.TrimPrefix(u.Path, "/")
			query, err := url.ParseQuery(u.RawQuery)
			if err != nil {
				r.fail("DATABASE_URL", "invalid query parameters")
			}
			for key, values := range query {
				if (key != "sslmode" && key != "sslrootcert") || len(values) != 1 {
					r.fail("DATABASE_URL", "only single sslmode and sslrootcert query parameters are supported")
					continue
				}
				setting := "DB_SSLMODE"
				target := &c.SSLMode
				if key == "sslrootcert" {
					setting = "DB_SSLROOTCERT"
					target = &root
				}
				if explicit, ok := r.values[setting]; ok && explicit != values[0] {
					r.fail("DATABASE_URL / "+setting, "conflicting TLS settings")
				}
				*target = values[0]
			}
		}
	} else {
		c.Host = required("DB_HOST", "127.0.0.1")
		c.Port = r.value("DB_PORT", "5432")
		c.Name = required("DB_NAME", "elabtrack_v2")
		c.User = required("DB_USER", "postgres")
		c.Password = r.value("DB_PASSWORD", "")
	}
	if c.Host == "" || strings.ContainsAny(c.Host, "/, \t\r\n\x00") {
		r.fail("DB_HOST / DATABASE_URL", "must identify one TCP host")
	}
	if strings.TrimSpace(c.User) == "" || strings.TrimSpace(c.Name) == "" || strings.ContainsRune(c.User+c.Name+c.Password, 0) {
		r.fail("DB_USER / DB_PASSWORD / DB_NAME / DATABASE_URL", "explicit user and database required; NUL is forbidden")
	}
	validatePort(r, "DB_PORT / DATABASE_URL", c.Port)
	if env == Production && strings.TrimSpace(c.Password) == "" {
		r.fail(source, "production requires an explicit nonblank database password")
	}
	if env == Production && c.SSLMode != "verify-full" {
		r.fail("DB_SSLMODE / DATABASE_URL", "production requires explicit verify-full TLS")
	}
	switch c.SSLMode {
	case "disable":
		if root != "" {
			r.fail("DB_SSLROOTCERT", "cannot be used with disabled TLS")
		}
	case "verify-full":
		c.tlsConfig = &tls.Config{MinVersion: tls.VersionTLS12, ServerName: c.Host}
		if root != "" && root != "system" {
			pem, err := os.ReadFile(root)
			if err != nil {
				r.fail("DB_SSLROOTCERT / DATABASE_URL", "cannot read CA bundle")
			} else {
				ca := x509.NewCertPool()
				if !ca.AppendCertsFromPEM(pem) {
					r.fail("DB_SSLROOTCERT / DATABASE_URL", "invalid CA bundle")
				} else {
					c.tlsConfig.RootCAs = ca
				}
			}
		}
	default:
		r.fail("DB_SSLMODE / DATABASE_URL", "must be disable (local only) or verify-full")
	}
	c.MaxConns = int32(r.integer("DB_MAX_CONNS", 20, 1, 1<<31-1))
	c.MinConns = int32(r.integer("DB_MIN_CONNS", 2, 0, int(c.MaxConns)))
	c.MaxConnLifetime = r.duration("DB_MAX_CONN_LIFETIME", time.Hour)
	c.MaxConnIdleTime = r.duration("DB_MAX_CONN_IDLE_TIME", 15*time.Minute)
	return c
}

// PoolConfig uses pgx's parser to initialize its internal defaults, then pins
// connection/TLS settings to the validated snapshot. No plaintext TLS fallback,
// password file, client certificate, service file or ambient PG* is supported.
func (c DBConfig) PoolConfig() (*pgxpool.Config, error) {
	u := &url.URL{Scheme: "postgresql", Host: net.JoinHostPort(c.Host, c.Port), Path: "/" + c.Name, User: url.UserPassword(c.User, c.Password)}
	q := url.Values{"sslmode": {c.SSLMode}, "passfile": {""}, "sslcert": {""}, "sslkey": {""}, "sslrootcert": {""}, "connect_timeout": {"5"}, "target_session_attrs": {"any"}}
	u.RawQuery = q.Encode()
	cfg, err := pgxpool.ParseConfig(u.String())
	if err != nil {
		return nil, errors.New("database: cannot initialize validated connection configuration")
	}
	cfg.ConnConfig.Host = c.Host
	cfg.ConnConfig.Database = c.Name
	cfg.ConnConfig.User = c.User
	cfg.ConnConfig.Password = c.Password
	cfg.ConnConfig.Fallbacks = nil
	cfg.ConnConfig.RuntimeParams = map[string]string{"application_name": "elabtrack-v2", "standard_conforming_strings": "on"}
	cfg.ConnConfig.TLSConfig = nil
	if c.tlsConfig != nil {
		cfg.ConnConfig.TLSConfig = c.tlsConfig.Clone()
	}
	cfg.MaxConns = c.MaxConns
	cfg.MinConns = c.MinConns
	cfg.MaxConnLifetime = c.MaxConnLifetime
	cfg.MaxConnIdleTime = c.MaxConnIdleTime
	cfg.HealthCheckPeriod = time.Minute
	return cfg, nil
}

// Migration credentials are optional for API startup. Production migration and
// seed actions cannot fall back; development/test may explicitly use their local
// runtime connection when no migration URL is configured.
func parseMigrationDatabase(r *reader, env Environment, runtime DBConfig) *DBConfig {
	raw, present := r.values["MIGRATION_DATABASE_URL"]
	if !present {
		return nil
	}
	values := map[string]string{"DATABASE_URL": raw}
	for _, key := range []string{"DB_MAX_CONNS", "DB_MIN_CONNS", "DB_MAX_CONN_LIFETIME", "DB_MAX_CONN_IDLE_TIME"} {
		if value, ok := r.values[key]; ok {
			values[key] = value
		}
	}
	migrationReader := reader{values: values}
	cfg := parseDatabase(&migrationReader, env)
	if len(migrationReader.errs) > 0 {
		r.fail("MIGRATION_DATABASE_URL", "invalid migration connection; requires the same URL/TLS policy as DATABASE_URL")
		return nil
	}
	if cfg.Host != runtime.Host || cfg.Port != runtime.Port || cfg.Name != runtime.Name {
		r.fail("MIGRATION_DATABASE_URL", "must target the runtime database host, port and name")
	}
	if env == Production && cfg.User == runtime.User {
		r.fail("MIGRATION_DATABASE_URL", "production migration and runtime users must differ")
	}
	return &cfg
}
func (c *Config) MigrationConnection() (DBConfig, error) {
	if c.MigrationDB != nil {
		return *c.MigrationDB, nil
	}
	if c.App.Env == Production {
		return DBConfig{}, errors.New("MIGRATION_DATABASE_URL: required for production migration commands")
	}
	return c.DB, nil
}
