package melon

import (
	"context"
	"database/sql"
)

func GetSqlExecutor(ctx context.Context, name string) TxDB {
	dbProvider := GetDBProvider(name)
	dbObj := dbProvider.GetTxOrDB(ctx)
	return dbObj.(TxDB)
}

// TxDB (SQL Go database connection) is a wrapper for SQL database handler ( can be *sql.DB or *sql.Tx)
// It should be able to work with all SQL data that follows SQL standard.
type TxDB interface {
	Exec(query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row

	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type SqlTxManager struct {
}

func (s *SqlTxManager) Begin(ctx context.Context, db interface{}, opts *sql.TxOptions) (tx interface{}, err error) {
	sqlDB := db.(*sql.DB)
	return sqlDB.BeginTx(ctx, opts)
}

func (s *SqlTxManager) Commit(tx interface{}) error {
	sqlTx := tx.(*sql.Tx)
	return sqlTx.Commit()
}

func (s *SqlTxManager) Rollback(tx interface{}) error {
	sqlTx := tx.(*sql.Tx)
	return sqlTx.Rollback()
}
