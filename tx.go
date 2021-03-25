package dbx

import (
	"database/sql"
)

type Tx struct {
	*ComplexExec
	rawTx *sql.Tx
}

func newTx(tx *sql.Tx, option *Options) *Tx {
	return &Tx{newComplexExec(tx, option), tx}
}

//Commit the transaction
func (t *Tx) Commit() error {
	return t.rawTx.Commit()
}

//Rollback the transaction
func (t *Tx) Rollback() error {
	return t.rawTx.Rollback()
}
