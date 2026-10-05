package database

import (
	"context"
	"fmt"
	"time"

	"github.com/fullstacktemplate/backend/internal/config"
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
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	// Pool sizing is also expressed in the DSN; set here too for clarity.
	poolCfg.MaxConns = cfg.MaxConns
	poolCfg.MinConns = cfg.MinConns
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolCfg.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
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
