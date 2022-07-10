package data

import (
	"context"
	"database/sql"
)

type TXName string

type TXManager interface {
	Begin(ctx context.Context, opts *sql.TxOptions) (newCtx context.Context, err error)
	Commit(tx any) error
	Rollback(tx any) error
	Name() TXName
	DB() interface{}
	TX(ctx context.Context) interface{}
	TXOrDB(ctx context.Context) interface{}
}
