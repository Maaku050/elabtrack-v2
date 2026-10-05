package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthChecker reports the health of the database connection pool.
type HealthChecker struct {
	pool *pgxpool.Pool
}

// NewHealthChecker constructs a HealthChecker for the given pool.
func NewHealthChecker(pool *pgxpool.Pool) *HealthChecker {
	return &HealthChecker{pool: pool}
}

// Check pings the database with a short timeout and returns nil if healthy.
func (h *HealthChecker) Check(ctx context.Context) error {
	if h == nil || h.pool == nil {
		return errNoPool
	}
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return h.pool.Ping(pingCtx)
}

var errNoPool = healthErr("database pool is not initialized")

type healthErr string

func (e healthErr) Error() string { return string(e) }
