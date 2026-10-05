package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager provides a transaction abstraction for workflows involving
// multiple related database operations. If any operation fails the
// transaction is rolled back; otherwise it is committed.
//
// Usage:
//
//	err := txm.Run(ctx, func(ctx context.Context, tx pgx.Tx) error {
//	    if err := repoA.WithTx(tx).Create(ctx, a); err != nil {
//	        return err
//	    }
//	    return repoB.WithTx(tx).Create(ctx, b)
//	})
type TxManager struct {
	pool *pgxpool.Pool
}

// NewTxManager constructs a TxManager backed by the given pool.
func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// Run executes fn inside a database transaction. The pgx.Tx is injected
// through the context so repositories can opt-in to participating in the
// transaction via TxFromContext.
func (m *TxManager) Run(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx) // safe to call after Commit; pgx treats it as no-op
	}()

	txCtx := WithTx(ctx, tx)
	if err := fn(txCtx, tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

type txKey struct{}

// WithTx stores a pgx.Tx in the context so repositories can pick it up.
func WithTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

// TxFromContext returns the pgx.Tx stored in the context, if any.
func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}
