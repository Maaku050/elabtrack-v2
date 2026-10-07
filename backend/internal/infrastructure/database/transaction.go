package database

import (
	"context"
	"time"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager provides a transaction abstraction for workflows involving
// multiple related database operations. If any operation fails the
// transaction is rolled back; otherwise it is committed.
//
// Usage:
//
//	err := txm.Within(ctx, func(ctx context.Context) error {
//	    if err := repoA.Create(ctx, a); err != nil {
//	        return err
//	    }
//	    return repoB.Create(ctx, b)
//	})
type TxManager struct {
	pool transactionBeginner
}

type transactionBeginner interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
}

// Within implements the inward application transaction port without pgx types.
func (m *TxManager) Within(ctx context.Context, fn func(context.Context) error) error {
	return m.Run(ctx, func(txCtx context.Context, _ pgx.Tx) error { return fn(txCtx) })
}

// NewTxManager constructs a TxManager backed by the given pool.
func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// Run executes fn inside a database transaction. The pgx.Tx is injected
// through the context so repositories can opt-in to participating in the
// transaction via TxFromContext.
func (m *TxManager) Run(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	// Explicit isolation: each locking read observes the committed winner's
	// revoked row after waiting, even when the server's default is different.
	tx, err := m.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return shared.Internal("transaction.begin", err)
	}
	defer func() {
		rollbackCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		_ = tx.Rollback(rollbackCtx)
	}()

	txCtx := WithTx(ctx, tx)
	if err := fn(txCtx, tx); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return shared.Internal("transaction.commit", err)
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
