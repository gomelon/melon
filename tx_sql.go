package melon

import (
	"context"
	"database/sql"
)

// SQLExecutor (SQL Go database connection) is a wrapper for SQL database handler ( can be *sql.DB or *sql.Tx)
// It should be able to work with all SQL data that follows SQL standard.
type SQLExecutor interface {
	// PrepareContext creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the
	// returned statement.
	// The caller must call the statement's Close method
	// when the statement is no longer needed.
	//
	// The provided context is used for the preparation of the statement, not for the
	// execution of the statement.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	// Prepare creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the
	// returned statement.
	// The caller must call the statement's Close method
	// when the statement is no longer needed.
	//
	// Prepare uses context.Background internally; to specify the context, use
	// PrepareContext.
	Prepare(query string) (*sql.Stmt, error)

	// ExecContext executes a query that doesn't return rows.
	// For example: an INSERT and UPDATE.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	// Exec executes a query that doesn't return rows.
	// For example: an INSERT and UPDATE.
	//
	// Exec uses context.Background internally; to specify the context, use
	// ExecContext.
	Exec(query string, args ...any) (sql.Result, error)

	// QueryContext executes a query that returns rows, typically a SELECT.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	// Query executes a query that returns rows, typically a SELECT.
	//
	// Query uses context.Background internally; to specify the context, use
	// QueryContext.
	Query(query string, args ...any) (*sql.Rows, error)

	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row

	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until
	// Row's Scan method is called.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	// Otherwise, the *Row's Scan scans the first selected row and discards
	// the rest.
	//
	// QueryRow uses context.Background internally; to specify the context, use
	// QueryRowContext.
	QueryRow(query string, args ...any) *sql.Row
}

type SQLTXManager struct {
	name TXName
	db   *sql.DB
}

func NewSqlTxManager(name string, db *sql.DB) *SQLTXManager {
	return &SQLTXManager{
		name: TXName(name),
		db:   db,
	}
}

func (tm *SQLTXManager) Begin(ctx context.Context, opts *sql.TxOptions) (newCtx context.Context, err error) {
	if ctx.Value(tm.name) != nil {
		newCtx = ctx
		return
	}
	tx, err := tm.db.BeginTx(ctx, opts)
	if err != nil {
		return ctx, err
	}
	newCtx = context.WithValue(ctx, tm.name, tx)
	return
}

func (tm *SQLTXManager) Commit(tx any) error {
	sqlTx := tx.(*sql.Tx)
	return sqlTx.Commit()
}

func (tm *SQLTXManager) Rollback(tx any) error {
	sqlTx := tx.(*sql.Tx)
	return sqlTx.Rollback()
}

func (tm *SQLTXManager) Name() TXName {
	return tm.name
}

func (tm *SQLTXManager) DB() interface{} {
	return tm.db
}

func (tm *SQLTXManager) TX(ctx context.Context) interface{} {
	return ctx.Value(tm.name)
}

func (tm *SQLTXManager) TXOrDB(ctx context.Context) interface{} {
	tx := tm.TX(ctx)
	if nil != tx {
		return tx
	}
	return tm.db
}

func (tm *SQLTXManager) OriginTXOrDB(ctx context.Context) SQLExecutor {
	tx := tm.TX(ctx)
	if nil != tx {
		return tx.(SQLExecutor)
	}
	return tm.db
}
