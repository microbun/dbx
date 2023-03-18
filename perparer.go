package dbx

import (
	"context"
	"database/sql"
)

//prepare is a abstract interface for sql.DB and sql.Tx
type preparer interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (r sql.Result, err error)

	//Query(query string, args ...interface{}) (*sql.Rows, error)
	//QueryContext(ctx context.Context, query string, args ...interface{}) (r *sql.Rows, err error)
}
