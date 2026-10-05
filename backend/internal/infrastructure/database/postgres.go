package database

import (
	"context"
	"errors"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres wraps a *pgxpool.Pool so the rest of the codebase depends on a
// small, stable type rather than the pool directly.
type Postgres struct {
	Pool *pgxpool.Pool
}

// New constructs a Postgres connection pool from config and pings it.
// One pool is created per process and injected via dependencies.
func New(ctx context.Context, cfg config.DBConfig) (*Postgres, error) {
	poolCfg, err := cfg.PoolConfig()
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, errors.New("database: failed to create connection pool")
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, errors.New("database: connection check failed (check connectivity, credentials and TLS)")
	}

	return &Postgres{Pool: pool}, nil
}

// Close releases all pool connections.
func (p *Postgres) Close() {
	if p == nil || p.Pool == nil {
		return
	}
	p.Pool.Close()
}
