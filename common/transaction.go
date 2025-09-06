package common

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

// TxFunc represents a function that can be executed within a transaction
type TxFunc func(tx pgx.Tx) error

// withTransaction executes a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// If the function succeeds, the transaction is committed.
func withTransaction(ctx context.Context, pool *pgxpool.Pool, fn TxFunc) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				// Log rollback error but don't override original error
				// In production, you might want to use structured logging here
			}
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// TxContext holds transaction context that can be passed to repositories
type TxContext struct {
	ctx context.Context
	tx  pgx.Tx
}

// NewTxContext creates a new transaction context
func NewTxContext(ctx context.Context, tx pgx.Tx) *TxContext {
	return &TxContext{
		ctx: ctx,
		tx:  tx,
	}
}

// Context returns the underlying context
func (tc *TxContext) Context() context.Context {
	return tc.ctx
}

// Tx returns the underlying transaction
func (tc *TxContext) Tx() pgx.Tx {
	return tc.tx
}

// TxContextFunc represents a function that can be executed within a transaction context
type TxContextFunc func(txCtx *TxContext) error

// withTxContext executes a function within a transaction, providing a TxContext
// that can be shared across multiple repositories.
func withTxContext(ctx context.Context, pool *pgxpool.Pool, fn TxContextFunc) error {
	return withTransaction(ctx, pool, func(tx pgx.Tx) error {
		txCtx := NewTxContext(ctx, tx)
		return fn(txCtx)
	})
}

// TransactionExecutor allows for dependency injection of transaction execution logic
type TransactionExecutor interface {
	WithTxContext(ctx context.Context, fn func(*TxContext) error) error
}

// PgxTransactionExecutor implements TransactionExecutor using a real pool
type PgxTransactionExecutor struct {
	pool *pgxpool.Pool
}

func NewPgxTransactionExecutor(pool *pgxpool.Pool) *PgxTransactionExecutor {
	return &PgxTransactionExecutor{pool: pool}
}

func (p *PgxTransactionExecutor) WithTxContext(ctx context.Context, fn func(*TxContext) error) error {
	return withTxContext(ctx, p.pool, fn)
}
