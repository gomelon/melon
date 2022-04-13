package melon

import (
	"context"
	"database/sql"
)

func GetSqlExecutor(ctx context.Context, name string) SqlExecutor {
	dbProvider := GetDBProvider(name)
	dbObj := dbProvider.GetTxOrDB(ctx)
	return dbObj.(SqlExecutor)
}

// SqlExecutor (SQL Go database connection) is a wrapper for SQL database handler ( can be *sql.DB or *sql.Tx)
// It should be able to work with all SQL data that follows SQL standard.
type SqlExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
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
