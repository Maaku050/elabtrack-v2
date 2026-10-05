package config

import (
	"fmt"
	"strings"
)

// DSN builds a PostgreSQL connection string (URL form) suitable for pgxpool.
func (c *DBConfig) DSN() string {
	password := c.Password
	if password != "" {
		password = fmt.Sprintf(":%s", urlEscape(password))
	}
	return fmt.Sprintf(
		"postgres://%s%s@%s:%s/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d&pool_max_conn_lifetime=%s&pool_max_conn_idle_time=%s",
		c.User,
		password,
		c.Host,
		c.Port,
		c.Name,
		c.MaxConns,
		c.MinConns,
		c.MaxConnLifetime.String(),
		c.MaxConnIdleTime.String(),
	)
}

func urlEscape(s string) string {
	r := strings.NewReplacer(
		" ", "%20",
		"@", "%40",
		"/", "%2F",
		":", "%3A",
		"?", "%3F",
		"#", "%23",
		"&", "%26",
	)
	return r.Replace(s)
}
