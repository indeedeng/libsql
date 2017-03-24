package libsql

import (
	"context"
	"database/sql"
)

func newDatabase(db sqlDB) Database {
	return &databaseImpl{
		Queryer:  newQueryerMixin(db),
		Preparer: newPreparerMixin(db),
		db:       db,
		newTX:    newTransaction,
	}
}

type databaseImpl struct {
	Queryer
	Preparer

	db    sqlDB
	newTX func(sqlTx) Transaction
}

var _ Database = (*databaseImpl)(nil)

// Transaction implements Database.Transaction
func (d databaseImpl) Transaction(ctx context.Context, work func(Transaction) error) error {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return err
	}

	rollbackIfNeeded := func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			// failed to rollback a transaction
		}
	}
	defer rollbackIfNeeded()

	if err := work(d.newTX(tx)); err != nil {
		return err
	}

	return tx.Commit()
}

// Close implements io.Close
func (d *databaseImpl) Close() error {
	return d.db.Close()
}
