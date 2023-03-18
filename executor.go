package dbx

import (
	"context"
	"database/sql"
	"errors"
)

var _ Executor = &executor{}

type Executor interface {
	// NamedExecContext executes a Named query without returning any rows.
	// The arg are for any placeholder parameters in the query.
	NamedExecContext(ctx context.Context, query string, arg map[string]interface{}) (sql.Result, error)

	// NamedExec executes a Named query without returning any rows.
	// The arg are for any placeholder parameters in the query.
	NamedExec(query string, arg map[string]interface{}) (sql.Result, error)

	// MustNamedExec is the same as  NamedExec but panics if cannot insert.
	MustNamedExec(query string, arg map[string]interface{}) sql.Result

	// NamedGetContext execute the Named query and scan the first row to dest, dest must be a pointer.
	// if dest is a type supported by the database (string , int, []byte, time.Time, etc.)
	// and there is only one column, the row will be assigned to dest.
	// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
	// A sql.ErrNoRows is returned if the result set is empty.
	NamedGetContext(ctx context.Context, dest interface{}, query string, arg map[string]interface{}) error

	// NamedGet execute the Named query and scan the first row to dest, dest must be a pointer.
	// if dest is a type supported by the database (string , int, []byte, time.Time, etc.)
	// and there is only one column, the row will be assigned to dest.
	// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
	// A sql.ErrNoRows is returned if the result set is empty.
	NamedGet(dest interface{}, query string, arg map[string]interface{}) error

	// MustNamedGet is the same as NamedGet, if there is a error in query, it will panics,
	// but sql.ErrNoRows will return false
	MustNamedGet(dest interface{}, query string, arg map[string]interface{}) bool

	// NamedQueryContext execute the query and scan the rows to dest, dest must be a slice of pointer.
	// if dest is a type supported by the database ([]string , []int, []byte, []time.Time, etc.)
	// and there is only one column, the rows will be assigned to dest.
	// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
	NamedQueryContext(ctx context.Context, dest interface{}, query string, arg map[string]interface{}) error

	// NamedQuery execute the query and scan the rows to dest, dest must be a slice of pointer.
	// if dest is a type supported by the database ([]string , []int, []byte, []time.Time, etc.)
	// and there is only one column, the rows will be assigned to dest.
	// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
	NamedQuery(dest interface{}, query string, arg map[string]interface{}) error

	// MustNamedQuery like NamedQuery but panics if cannot query.
	MustNamedQuery(dest interface{}, query string, arg map[string]interface{})

	// PrepareContext creates a prepared statement
	PrepareContext(ctx context.Context, query string) (*Stmt, error)

	// Prepare creates a prepared statement
	Prepare(query string) (*Stmt, error)

	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	Exec(query string, args ...interface{}) (sql.Result, error)

	// MustExec is the same as  Exec but panics if cannot insert.
	MustExec(query string, args ...interface{}) sql.Result

	// GetContext execute the query and scan the first row to dest, dest must be a pointer.
	// if dest is a type supported by the database (string , int, []byte, time.Time, etc.)
	// and there is only one column, the row will be assigned to dest.
	// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
	// A sql.ErrNoRows is returned if the result set is empty.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error)

	// Get execute the query and scan the first row to dest, dest must be a pointer.
	// If dest is a type supported by the database (string , int, []byte, time.Time, etc.)
	// and there is only one column, the row will be assigned to dest.
	// If dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
	// A sql.ErrNoRows is returned, if the query selects no rows.
	Get(dest interface{}, query string, args ...interface{}) (err error)

	// MustGet is the same as Get, if there is a error in query, it will panics,
	// but sql.ErrNoRows will return false
	MustGet(dest interface{}, query string, args ...interface{}) bool

	// QueryContext execute the query and scan the rows to dest, dest must be a slice of pointer.
	// if dest is a type supported by the database ([]string , []int, []byte, []time.Time, etc.)
	// and there is only one column, the rows will be assigned to dest.
	// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
	QueryContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error)

	// Query execute the query and scan the rows to dest, dest must be a slice of pointer.
	// if dest is a type supported by the database ([]string , []int, []byte, []time.Time, etc.)
	// and there is only one column, the rows will be assigned to dest.
	// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
	Query(dest interface{}, query string, args ...interface{}) (err error)

	// MustQuery is the same as Query but panics if cannot query.
	MustQuery(dest interface{}, query string, args ...interface{})

	// InsertContext insert a struct to database
	InsertContext(ctx context.Context, value interface{}) (rs sql.Result, err error)

	// Insert a struct to database
	Insert(value interface{}) (sql.Result, error)

	// MustInsert is the same as Insert but panics if cannot insert.
	MustInsert(value interface{}) sql.Result

	// UpdateContext update the rows according to the value of structure, if the column name is specified,
	// only the specified column is updated.
	UpdateContext(ctx context.Context, value interface{}, columns ...string) (rs sql.Result, err error)

	// Update the rows according to the value of structure, if the column name is specified,
	// only the specified column is updated.
	Update(value interface{}, columns ...string) (sql.Result, error)

	// MustUpdate is the same as Update, but panics if cannot update.
	MustUpdate(value interface{}, columns ...string) (rs sql.Result)
}

// executor implemented SQLExecutor, NamedExecutor and StructExecutor
type executor struct {
	option   *Options
	preparer preparer
}

func newDefaultExecutor(preparer preparer, option *Options) *executor {
	return &executor{preparer: preparer, option: option}
}

// InsertContext insert a struct to database
func (e *executor) InsertContext(ctx context.Context, value interface{}) (rs sql.Result, err error) {
	atv, query, values, err := e.option.Generator.InsertSQL(value)
	if err != nil {
		return
	}
	rs, err = e.ExecContext(ctx, query, values...)
	if err != nil {
		return rs, err
	}
	id, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}
	if atv != nil {
		atv.SetInt(id)
	}
	return
}

// Insert a struct to database
func (e *executor) Insert(value interface{}) (sql.Result, error) {
	return e.InsertContext(context.Background(), value)
}

// MustInsert is the same as Insert but panics if cannot insert.
func (e *executor) MustInsert(value interface{}) sql.Result {
	rs, err := e.InsertContext(context.Background(), value)
	if err != nil {
		panic(err)
	}
	return rs
}

// UpdateContext update the rows according to the value of structure, if the column name is specified,
// only the specified column is updated.
func (e *executor) UpdateContext(ctx context.Context, value interface{}, columns ...string) (rs sql.Result, err error) {
	query, values, err := e.option.Generator.UpdateSQL(value, columns...)
	if err != nil {
		return
	}
	return e.ExecContext(ctx, query, values...)
}

// Update the rows according to the value of structure, if the column name is specified,
// only the specified column is updated.
func (e *executor) Update(value interface{}, columns ...string) (sql.Result, error) {
	return e.UpdateContext(context.Background(), value, columns...)
}

// MustUpdate is the same as Update, but panics if cannot update.
func (e *executor) MustUpdate(value interface{}, columns ...string) (rs sql.Result) {
	rs, err := e.UpdateContext(context.Background(), value, columns...)
	if err != nil {
		panic(err)
	}
	return
}

// Prepare creates a prepared statement
func (e *executor) Prepare(query string) (*Stmt, error) {
	return newStmt(e.preparer, query, e.option)
}

// PrepareContext creates a prepared statement
func (e *executor) PrepareContext(ctx context.Context, query string) (*Stmt, error) {
	return newStmtContext(ctx, e.preparer, query, e.option)
}

// GetContext execute the query and scan the first row to dest, dest must be a pointer.
// if dest is a type supported by the database (string , int, []byte, time.Time, etc.)
// and there is only one column, the row will be assigned to dest.
// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
// An sql.ErrNoRows is returned if the result set is empty.
func (e *executor) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	stmt, err := e.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer func() {
		_ = stmt.Close()
	}()
	return stmt.GetContext(ctx, dest, args...)
}

// Get execute the query and scan the first row to dest, dest must be a pointer.
// if dest is a type supported by the database (string , int, []byte, time.Time, etc.)
// and there is only one column, the row will be assigned to dest.
// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
// An sql.ErrNoRows is returned if the result set is empty.
func (e *executor) Get(dest interface{}, query string, args ...interface{}) (err error) {
	return e.GetContext(context.Background(), dest, query, args...)
}

// MustGet is the same as Get, if there is a error in query, it will panics,
// but sql.ErrNoRows will return false
func (e *executor) MustGet(dest interface{}, query string, args ...interface{}) bool {
	err := e.GetContext(context.Background(), dest, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
		panic(err)
	}
	return true
}

// QueryContext execute the query and scan the rows to dest, dest must be a slice of pointer.
// if dest is a type supported by the database ([]string , []int, []byte, []time.Time, etc.)
// and there is only one column, the rows will be assigned to dest.
// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
func (e *executor) QueryContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	stmt, err := e.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer func() {
		_ = stmt.Close()
	}()
	return stmt.QueryContext(ctx, dest, args...)
}

// Query execute the query and scan the rows to dest, dest must be a slice of pointer.
// if dest is a type supported by the database ([]string , []int, []byte, []time.Time, etc.)
// and there is only one column, the rows will be assigned to dest.
// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
func (e *executor) Query(dest interface{}, query string, args ...interface{}) (err error) {
	return e.QueryContext(context.Background(), dest, query, args...)
}

// MustQuery is the same as Query but panics if cannot query.
func (e *executor) MustQuery(dest interface{}, query string, args ...interface{}) {
	err := e.QueryContext(context.Background(), dest, query, args...)
	if err != nil {
		panic(err)
	}
}

// ExecContext executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (e *executor) ExecContext(ctx context.Context, query string, args ...interface{}) (rs sql.Result, err error) {
	stmt, err := e.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer func() {
		_ = stmt.Close()
	}()
	return stmt.ExecContext(ctx, args...)
	//printSQL(query, args, e.option)
	//return e.preparer.ExecContext(ctx, query, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (e *executor) Exec(query string, args ...interface{}) (rs sql.Result, err error) {
	return e.ExecContext(context.Background(), query, args...)
}

// MustExec is the same as  Exec but panics if cannot insert.
func (e *executor) MustExec(query string, args ...interface{}) (rs sql.Result) {
	rs, err := e.ExecContext(context.Background(), query, args...)
	if err != nil {
		panic(err)
	}
	return
}

// NamedExecContext executes a Named query without returning any rows.
// The arg are for any placeholder parameters in the query.
func (e *executor) NamedExecContext(ctx context.Context, query string, arg map[string]interface{}) (sql.Result, error) {
	query, args, err := namedCompile(query, arg)
	if err != nil {
		return nil, err
	}
	return e.ExecContext(ctx, query, args...)
}

// NamedExec executes a Named query without returning any rows.
// The arg are for any placeholder parameters in the query.
func (e *executor) NamedExec(query string, arg map[string]interface{}) (sql.Result, error) {
	return e.NamedExecContext(context.Background(), query, arg)
}

// MustNamedExec is the same as  NamedExec but panics if cannot insert.
func (e *executor) MustNamedExec(query string, arg map[string]interface{}) sql.Result {
	r, err := e.NamedExecContext(context.Background(), query, arg)
	if err != nil {
		panic(err)
	}
	return r
}

// NamedGetContext execute the Named query and scan the first row to dest, dest must be a pointer.
// if dest is a type supported by the database (string , int, []byte, time.Time, etc.)
// and there is only one column, the row will be assigned to dest.
// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
// A sql.ErrNoRows is returned if the result set is empty.
func (e *executor) NamedGetContext(ctx context.Context, dest interface{}, query string, arg map[string]interface{}) (err error) {
	query, args, err := namedCompile(query, arg)
	if err != nil {
		return
	}
	return e.GetContext(ctx, dest, query, args...)
}

// NamedGet execute the Named query and scan the first row to dest, dest must be a pointer.
// if dest is a type supported by the database (string , int, []byte, time.Time, etc.)
// and there is only one column, the row will be assigned to dest.
// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
// A sql.ErrNoRows is returned if the result set is empty.
func (e *executor) NamedGet(dest interface{}, query string, arg map[string]interface{}) (err error) {
	return e.NamedGetContext(context.Background(), dest, query, arg)
}

// MustNamedGet is the same as NamedGet, if there is a error in query, it will panics,
// but sql.ErrNoRows will return false
func (e *executor) MustNamedGet(dest interface{}, query string, arg map[string]interface{}) bool {
	err := e.NamedGetContext(context.Background(), dest, query, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false
		}
		panic(err)
	}
	return true
}

// NamedQueryContext execute the query and scan the rows to dest, dest must be a slice of pointer.
// if dest is a type supported by the database ([]string , []int, []byte, []time.Time, etc.)
// and there is only one column, the rows will be assigned to dest.
// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
func (e *executor) NamedQueryContext(ctx context.Context, dest interface{}, query string, arg map[string]interface{}) (err error) {
	query, args, err := namedCompile(query, arg)
	if err != nil {
		return
	}
	return e.QueryContext(ctx, dest, query, args...)
}

// NamedQuery execute the query and scan the rows to dest, dest must be a slice of pointer.
// if dest is a type supported by the database ([]string , []int, []byte, []time.Time, etc.)
// and there is only one column, the rows will be assigned to dest.
// if dest is a struct, it will be mapped to the dbx tag field in the struct according to the name of each column.
func (e *executor) NamedQuery(dest interface{}, query string, arg map[string]interface{}) (err error) {
	return e.NamedQueryContext(context.Background(), dest, query, arg)
}

// MustNamedQuery like NamedQuery but panics if cannot query.
func (e *executor) MustNamedQuery(dest interface{}, query string, arg map[string]interface{}) {
	err := e.NamedQuery(dest, query, arg)
	if err != nil {
		panic(err)
	}
}
