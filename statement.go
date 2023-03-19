package dbx

import (
	"context"
	"database/sql"
	"time"
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

func timeFormat(fn func(t *time.Time) string, args ...interface{}) []interface{} {
	var _args = make([]interface{}, 0)
	for _, arg := range args {

		if t, ok := arg.(*time.Time); ok {
			_args = append(_args, fn(t))
			continue
		}
		if t, ok := arg.(time.Time); ok {
			_args = append(_args, fn(&t))
			continue
		}
		_args = append(_args, arg)
	}
	return _args
}

func (s *Stmt) format(ctx context.Context, args ...interface{}) []interface{} {
	if s.option.TimeFormat != nil {
		args = timeFormat(s.option.TimeFormat, args...)
	}
	printSQL(s.rawQuery, args, s.option)
	return args
}

func (s *Stmt) GetContext(ctx context.Context, dest interface{}, args ...interface{}) error {
	args = s.format(ctx, args...)
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
	args = s.format(ctx, args...)
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

func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error) {
	args = s.format(ctx, args...)
	return s.stmt.ExecContext(ctx, args...)
}

func (s *Stmt) Exec() (sql.Result, error) {
	return s.ExecContext(context.Background())
}

func (s *Stmt) Close() error {
	return s.stmt.Close()
}
