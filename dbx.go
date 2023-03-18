package dbx

import (
	"context"
	"database/sql"
	"time"
)

type Options struct {
	Logger    Logger
	Generator SQLGenerator
	Location  *time.Location
}

type DB struct {
	*executor
	option *Options
	rawDB  *sql.DB
}

func newDBX(db *sql.DB, options *Options) *DB {
	exec := newDefaultExecutor(db, options)
	return &DB{executor: exec, option: options, rawDB: db}
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
	return newDBX(db, &Options{
		Logger:    logger,
		Generator: NewCommonSQLGenerator(),
		Location: time.Local,
	}), nil
}

func (d *DB) Options() *Options {
	return d.option
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

//Transaction begin a transaction and commit automatically,automatically roll back when there is an error.
func (d *DB) Transaction(fn func(*Tx) error) error {
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
	err = fn(tx)
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
