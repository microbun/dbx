package dbx

import (
	"context"
	"database/sql"
	"io"
	"os"
)

type Options struct {
	out       io.Writer
	generator SQLGenerator
}

type DB struct {
	ComplexExecutor
	option *Options
	rawDB  *sql.DB
}

func newDBX(db *sql.DB, options *Options) *DB {
	exec := newComplexExec(db, options)
	return &DB{ComplexExecutor: exec, option: options, rawDB: db}
}

//Connect a database to dbx
func Connect(db *sql.DB) *DB {
	return newDBX(db, nil)
}

//Open a database
func Open(driverName string, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	switch driverName {
	case "sqlite3":
		return newDBX(db, &Options{
			out:       os.Stdout,
			generator: NewSQLiteSQLiteGenerator(),
		}), nil
	case "mysql":
		return newDBX(db, &Options{
			out:       os.Stdout,
			generator: NewMySQLSQLiteGenerator(),
		}), nil
	default:
		return newDBX(db, &Options{
			out:       os.Stdout,
			generator: NewCommonSQLGenerator(),
		}), nil
	}
}

// SetSQLGenerator set a SQLGenerator for struct.
func (d *DB) SetSQLGenerator(generator SQLGenerator) {
	d.option.generator = generator
}

// SetSQLOutput write SQL statements to io, usually used to view SQL execution logs
func (d *DB) SetSQLOutput(io io.Writer) {
	d.option.out = io
}

//BeginTx begin a Transaction with opts
func (d *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := d.rawDB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return newTx(tx, d.option), nil
}

//Begin a Transaction
func (d *DB) Begin() (*Tx, error) {
	return d.BeginTx(context.Background(), nil)
}

//TxFunc is Transaction process func
type TxFunc func(*Tx) error

//ExecTx begin a transaction and commit automatically,automatically roll back when there is an error.
func (d *DB) ExecTx(txProc TxFunc) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}

	defer func() {
		txErr := tx.Rollback()
		if txErr != sql.ErrTxDone {
			err = txErr
		}
	}()
	err = txProc(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// Close the database and prevents newDBX queries from starting.
// Close then waits for all queries that have started processing on the server
// to finish.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (d *DB) Close() error {
	return d.rawDB.Close()
}

//RawDB return a raw sql.DB object
func (d *DB) RawDB() *sql.DB {
	return d.rawDB
}
