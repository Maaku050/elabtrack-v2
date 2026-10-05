package database

import (
	"context"
	"errors"
	"testing"

	"github.com/Maaku050/elabtrack-v2/backend/internal/domain/shared"
	"github.com/jackc/pgx/v5"
)

type transactionPool struct {
	tx       *managedTx
	beginErr error
	options  pgx.TxOptions
}

func (p *transactionPool) BeginTx(_ context.Context, opts pgx.TxOptions) (pgx.Tx, error) {
	p.options = opts
	return p.tx, p.beginErr
}

type managedTx struct {
	pgx.Tx
	commits, rollbacks int
	commitErr          error
	rollbackCanceled   bool
}

func (tx *managedTx) Commit(context.Context) error { tx.commits++; return tx.commitErr }
func (tx *managedTx) Rollback(ctx context.Context) error {
	tx.rollbacks++
	tx.rollbackCanceled = ctx.Err() != nil
	return nil
}
func TestTransactionPortCommitRollbackAndContext(t *testing.T) {
	for _, name := range []string{"success", "callback failure", "commit failure", "begin failure", "canceled request"} {
		t.Run(name, func(t *testing.T) {
			tx := &managedTx{}
			pool := &transactionPool{tx: tx}
			m := &TxManager{pool: pool}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			failure := errors.New("sensitive-database-detail")
			called := false
			if name == "commit failure" {
				tx.commitErr = failure
			}
			if name == "begin failure" {
				pool.beginErr = failure
			}
			err := m.Within(ctx, func(bound context.Context) error {
				called = true
				if actual, ok := TxFromContext(bound); !ok || actual != tx {
					t.Fatal("repositories do not share actual transaction")
				}
				if name == "canceled request" {
					cancel()
					return failure
				}
				if name == "callback failure" {
					return failure
				}
				return nil
			})
			if pool.options.IsoLevel != pgx.ReadCommitted {
				t.Fatal("isolation must be explicit")
			}
			switch name {
			case "success":
				if err != nil || tx.commits != 1 || tx.rollbacks != 1 {
					t.Fatal("commit lifecycle incorrect")
				}
			case "begin failure":
				if err != shared.ErrInternal || called || tx.commits != 0 || tx.rollbacks != 0 {
					t.Fatal("begin failure unsafe")
				}
			case "commit failure":
				if err != shared.ErrInternal || tx.commits != 1 || tx.rollbacks != 1 {
					t.Fatal("commit failure returned success/unsafe detail")
				}
			default:
				if err != failure || tx.commits != 0 || tx.rollbacks != 1 {
					t.Fatal("callback failure must rollback without commit")
				}
			}
			if tx.rollbackCanceled {
				t.Fatal("rollback inherited canceled request")
			}
		})
	}
}
