package melon

import (
	"context"
	"database/sql"
)

var DBNameDefault = "Default"

var dbProviders = make(map[string]*DBProvider)

func RegisterDBProvider(provider *DBProvider) {
	dbProviders[provider.Name()] = provider
}

type TxName string

func NewDBProvider(name string, db interface{}, txManager TxManager) *DBProvider {
	return &DBProvider{
		name:      name,
		db:        db,
		txManager: txManager,
	}
}

func GetDBProvider(name string) *DBProvider {
	return dbProviders[name]
}

type DBProvider struct {
	name      string
	db        interface{}
	txManager TxManager
}

func (d *DBProvider) Name() string {
	return d.name
}

func (d *DBProvider) DB() interface{} {
	return d.db
}

// Tx returns a tx stored inside a context, or nil if there isn't one.
func (d *DBProvider) GetTx(ctx context.Context) interface{} {
	txName := TxName(d.name)
	return ctx.Value(txName)
}

func (d *DBProvider) GetTxOrDB(ctx context.Context) interface{} {
	tx := d.GetTx(ctx)
	if nil != tx {
		return tx
	}
	return d.db
}

func (d *DBProvider) Begin(ctx context.Context, txOptions *sql.TxOptions) (context.Context, error) {
	if nil != d.GetTx(ctx) {
		return ctx, nil
	}
	tx, err := d.txManager.Begin(ctx, d.db, txOptions)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, TxName(d.name), tx), nil
}

func (d *DBProvider) Complete(ctx context.Context, err error, rollbackFor ...error) error {
	if err != nil && (len(rollbackFor) == 0 || IsErrorIn(err, rollbackFor...)) {
		return d.Rollback(ctx)
	} else {
		return d.Commit(ctx)
	}
}

func (d *DBProvider) Commit(ctx context.Context) error {
	tx := d.GetTx(ctx)
	if nil == tx {
		return nil
	}
	return d.txManager.Commit(tx)
}

func (d *DBProvider) Rollback(ctx context.Context) error {
	tx := d.GetTx(ctx)
	if nil == tx {
		return nil
	}
	return d.txManager.Rollback(tx)
}

type TxManager interface {
	Begin(ctx context.Context, db interface{}, opts *sql.TxOptions) (tx interface{}, err error)
	Commit(tx interface{}) error
	Rollback(tx interface{}) error
}
