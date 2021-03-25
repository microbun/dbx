package dbx

import (
	"context"
	"database/sql"
)

// SQLStmt is a preparer statement.
// it is like sql.Stmt.
type SQLStmt struct {
	rawQuery string
	stmt     *sql.Stmt
	option   *Options
}

func newStmt(preparer preparer, query string, option *Options) (stmt *SQLStmt, err error) {
	return newStmtContext(context.Background(), preparer, query, option)
}

func newStmtContext(ctx context.Context, preparer preparer, query string, option *Options) (stmt *SQLStmt, err error) {
	s, err := preparer.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &SQLStmt{
		stmt:     s,
		rawQuery: query,
		option:   option,
	}, nil
}

func (s *SQLStmt) FirstContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	printSQL(s.rawQuery, args, s.option.out)
	rows, err := s.stmt.QueryContext(ctx, args...)
	if err != nil {
		return err
	}
	err = toDest(rows, dest, one)
	if err != nil {
		return err
	}
	return err
}

func (s *SQLStmt) First(dest interface{}, args ...interface{}) error {
	return s.FirstContext(context.Background(), dest, args...)
}

func (s *SQLStmt) FindContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	printSQL(s.rawQuery, args, s.option.out)
	rows, err := s.stmt.QueryContext(ctx, args...)
	if err != nil {
		return err
	}
	err = toDest(rows, dest, array)
	if err != nil {
		return err
	}
	return err
}

func (s *SQLStmt) Find(dest interface{}, args ...interface{}) error {
	return s.FindContext(context.Background(), dest, args...)
}

func (s *SQLStmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
	printSQL(s.rawQuery, args, s.option.out)
	return s.stmt.QueryRowContext(ctx, args...)
}

func (s *SQLStmt) QueryRow(args ...interface{}) *sql.Row {
	return s.stmt.QueryRow(args...)
}

func (s *SQLStmt) QueryContext(ctx context.Context, args ...interface{}) (*sql.Rows, error) {
	printSQL(s.rawQuery, args, s.option.out)
	return s.stmt.QueryContext(ctx, args...)
}

func (s *SQLStmt) Query() (*sql.Rows, error) {
	return s.stmt.Query()
}

func (s *SQLStmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	printSQL(s.rawQuery, args, s.option.out)
	return s.stmt.ExecContext(ctx, args...)
}

func (s *SQLStmt) Exec() (sql.Result, error) {
	return s.stmt.Exec()
}

//RawQuery 原生的语句
func (s *SQLStmt) RawQuery() string {
	return s.rawQuery
}

func (s *SQLStmt) Close() error {
	return s.stmt.Close()
}

type DQLStmt struct {
}
