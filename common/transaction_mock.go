package common

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockTransactionExecutor allows testing without real transactions
type MockTransactionExecutor struct {
	mock.Mock
}

func (m *MockTransactionExecutor) WithTxContext(ctx context.Context, fn func(*TxContext) error) error {
	// Create a mock transaction context and call the function directly
	mockTx := &MockTx{}
	txCtx := NewTxContext(ctx, mockTx)
	return fn(txCtx)
}

// MockTx provides a minimal transaction implementation for testing
type MockTx struct{}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) { return m, nil }
func (m *MockTx) Commit(ctx context.Context) error          { return nil }
func (m *MockTx) Rollback(ctx context.Context) error        { return nil }
func (m *MockTx) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return nil, nil
}
func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row { return nil }
func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults          { return nil }
func (m *MockTx) LargeObjects() pgx.LargeObjects                                        { return pgx.LargeObjects{} }
func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	return nil, nil
}
func (m *MockTx) Conn() *pgx.Conn { return nil }
func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return 0, nil
}
