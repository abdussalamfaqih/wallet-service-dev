package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Executor defines the interface for executing SQL commands
// This can be implemented by *sql.DB, *sql.Tx, or any custom wrapper
type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type Transactor interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type TransactionFunc func(ctx context.Context, tx Executor) error

type TransactionManager struct {
	db      Transactor
	timeout time.Duration
	options *sql.TxOptions
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db Transactor) *TransactionManager {

	txOptions := &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}

	return &TransactionManager{
		db:      db,
		timeout: 30 * time.Second,
		options: txOptions,
	}
}

func (tm *TransactionManager) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	return tm.executeTransaction(ctx, fn)
}

// executeTransaction executes a single transaction attempt
func (tm *TransactionManager) executeTransaction(ctx context.Context, fn TransactionFunc) error {
	// Create context with timeout if not already set
	if tm.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, tm.timeout)
		defer cancel()
	}

	// Begin transaction
	tx, err := tm.db.BeginTx(ctx, tm.options)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is closed
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after cleanup
		}
	}()

	// Execute the transaction function
	if err := fn(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Repository provides a base repository with transaction support
type Repository struct {
	db    Executor
	txMgr *TransactionManager
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB, config *Config) *Repository {
	txMgr := NewTransactionManager(db)
	return &Repository{
		db:    db,
		txMgr: txMgr,
	}
}

// WithTransaction executes repository operations within a transaction
func (r *Repository) WithTransaction(ctx context.Context, fn func(ctx context.Context, repo *Repository) error) error {
	return r.txMgr.WithTransaction(ctx, func(ctx context.Context, tx Executor) error {
		txRepo := &Repository{
			db:    tx,
			txMgr: r.txMgr,
		}
		return fn(ctx, txRepo)
	})
}

// Exec executes a query without returning any rows
func (r *Repository) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.db.ExecContext(ctx, query, args...)
}

// Query executes a query that returns rows
func (r *Repository) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that is expected to return at most one row
func (r *Repository) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return r.db.QueryRowContext(ctx, query, args...)
}

// Prepare creates a prepared statement
func (r *Repository) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	return r.db.PrepareContext(ctx, query)
}
