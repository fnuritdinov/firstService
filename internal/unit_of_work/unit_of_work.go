package unit_of_work

import (
	"context"
	"database/sql"
)

type TransactionManager struct {
	db *sql.DB
}

func New(db *sql.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (txn *TransactionManager) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := txn.db.BeginTx(ctx, nil)
	return tx, err
}

func (txn *TransactionManager) CommitTx(ctx context.Context, tx *sql.Tx) error {
	return tx.Commit()
}

func (txn *TransactionManager) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}
