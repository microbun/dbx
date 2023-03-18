package dbx

import (
	"database/sql"
)

type Tx struct {
	*executor
	tx *sql.Tx
}

func newTx(tx *sql.Tx, option *Options) *Tx {
	return &Tx{executor: newDefaultExecutor(tx, option), tx: tx}
}

//Commit the transaction
func (t *Tx) Commit() error {
	return t.tx.Commit()
}

//Rollback the transaction
func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}
