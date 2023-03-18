package dbx

import (
	"context"
	"database/sql"
)

// Stmt is a preparer statement.
// it is like sql.Stmt.
type Stmt struct {
	rawQuery string
	stmt     *sql.Stmt
	option   *Options
}

func newStmt(preparer preparer, query string, option *Options) (stmt *Stmt, err error) {
	return newStmtContext(context.Background(), preparer, query, option)
}

func newStmtContext(ctx context.Context, preparer preparer, query string, option *Options) (stmt *Stmt, err error) {
	s, err := preparer.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Stmt{
		stmt:     s,
		rawQuery: query,
		option:   option,
	}, nil
}

func (s *Stmt) GetContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	printSQL(s.rawQuery, args, s.option)
	rows, err := s.stmt.QueryContext(ctx, args...)
	if err != nil {
		return err
	}
	err = mapping(rows, dest, single)
	if err != nil {
		return err
	}
	return err
}

func (s *Stmt) Get(dest interface{}, args ...interface{}) error {
	return s.GetContext(context.Background(), dest, args...)
}

func (s *Stmt) QueryContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	printSQL(s.rawQuery, args, s.option)
	rows, err := s.stmt.QueryContext(ctx, args...)
	if err != nil {
		return err
	}
	err = mapping(rows, dest, slice)
	if err != nil {
		return err
	}
	return err
}

func (s *Stmt) Query(dest interface{}, args ...interface{}) error {
	return s.QueryContext(context.Background(), dest, args...)
}

//func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *sql.Row {
//	printSQL(s.rawQuery, args, s.option)
//	return s.stmt.QueryRowContext(ctx, args...)
//}
//
//func (s *Stmt) QueryRow(args ...interface{}) *sql.Row {
//	return s.QueryRowContext(context.Background(), args...)
//}

func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	printSQL(s.rawQuery, args, s.option)
	return s.stmt.ExecContext(ctx, args...)
}

func (s *Stmt) Exec() (sql.Result, error) {
	return s.ExecContext(context.Background())
}

func (s *Stmt) Close() error {
	return s.stmt.Close()
}
