package dbx

import (
	"context"
	"database/sql"
	"fmt"
)

type DQLExecutor interface {
	DQLExecContext(ctx context.Context, query string, argument DQLArgument) (sql.Result, error)

	//DQLExec executes a query without returning any rows.
	DQLExec(query string, argument DQLArgument) (sql.Result, error)

	DQLMustExecContext(ctx context.Context, query string, argument DQLArgument) sql.Result
	DQLMustExec(query string, argument DQLArgument) sql.Result

	DQLFirstContext(ctx context.Context, dest interface{}, query string, argument DQLArgument) error
	DQLFirst(dest interface{}, query string, argument DQLArgument) error
	DQLMustFirstContext(ctx context.Context, dest interface{}, query string, argument DQLArgument)
	DQLMustFirst(dest interface{}, query string, argument DQLArgument)

	DQLFindContext(ctx context.Context, dest interface{}, query string, argument DQLArgument) error
	DQLFind(dest interface{}, query string, argument DQLArgument) error
	DQLMustFindContext(ctx context.Context, dest interface{}, query string, argument DQLArgument)
	DQLMustFind(dest interface{}, query string, argument DQLArgument)
}

type SQLExecutor interface {
	SQLPrepareContext(ctx context.Context, query string) (*SQLStmt, error)

	SQLPrepare(query string) (*SQLStmt, error)
	// SQLExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	SQLExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// SQLExec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	SQLExec(query string, args ...interface{}) (sql.Result, error)

	SQLMustExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	SQLMustExec(query string, args ...interface{}) sql.Result

	SQLFirstContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error)
	SQLFirst(dest interface{}, query string, args ...interface{}) (err error)
	SQLMustFirstContext(ctx context.Context, dest interface{}, query string, args ...interface{})
	SQLMustFirst(dest interface{}, query string, args ...interface{})

	SQLFindContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error)
	SQLFind(dest interface{}, query string, args ...interface{}) (err error)
	SQLMustFindContext(ctx context.Context, dest interface{}, query string, args ...interface{})
	SQLMustFind(dest interface{}, query string, args ...interface{})
}

type StructExecutor interface {
	StructInsertContext(ctx context.Context, value interface{}) (rs sql.Result, err error)
	StructMustInsertContext(ctx context.Context, value interface{}) (rs sql.Result)
	StructInsert(value interface{}) (sql.Result, error)
	StructMustInsert(value interface{}) sql.Result
	StructUpdateContext(ctx context.Context, value interface{}) (rs sql.Result, err error)
	StructUpdate(value interface{}) (sql.Result, error)
	StructMustUpdateContext(ctx context.Context, value interface{}) (rs sql.Result)
	StructMustUpdate(value interface{}) (rs sql.Result)
}

type ComplexExecutor interface {
	StructExecutor
	SQLExecutor
	DQLExecutor
}

// ComplexExec implemented SQLExecutor, DQLExecutor and StructExecutor
type ComplexExec struct {
	option   *Options
	preparer preparer
}

func newComplexExec(preparer preparer, option *Options) *ComplexExec {
	return &ComplexExec{preparer: preparer, option: option}
}

//StructInsertContext
func (e *ComplexExec) StructInsertContext(ctx context.Context, value interface{}) (rs sql.Result, err error) {
	atv, query, values, err := e.option.generator.InsertSQL(value)
	if err != nil {
		return
	}
	rs, err = e.SQLExecContext(ctx, query, values...)
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

//StructMustInsertContext
func (e *ComplexExec) StructMustInsertContext(ctx context.Context, value interface{}) (rs sql.Result) {
	rs, err := e.StructInsertContext(ctx, value)
	if err != nil {
		panic(err)
	}
	return
}

// StructInsert insert a struct to database
func (e *ComplexExec) StructInsert(value interface{}) (sql.Result, error) {
	return e.StructInsertContext(context.Background(), value)
}

// StructMustInsert is like StructInsert but panics if cannot insert.
func (e *ComplexExec) StructMustInsert(value interface{}) sql.Result {
	return e.StructMustInsertContext(context.Background(), value)
}

//StructUpdateContext
func (e *ComplexExec) StructUpdateContext(ctx context.Context, value interface{}) (rs sql.Result, err error) {
	query, values, err := e.option.generator.UpdateSQL(value)
	if err != nil {
		return
	}
	return e.SQLExecContext(ctx, query, values...)
}

//StructUpdate
func (e *ComplexExec) StructUpdate(value interface{}) (sql.Result, error) {
	return e.StructUpdateContext(context.Background(), value)
}

//StructMustUpdateContext
func (e *ComplexExec) StructMustUpdateContext(ctx context.Context, value interface{}) (rs sql.Result) {
	rs, err := e.StructUpdateContext(ctx, value)
	if err != nil {
		panic(err)
	}
	return
}

//StructMustUpdate
func (e *ComplexExec) StructMustUpdate(value interface{}) (rs sql.Result) {
	return e.StructMustUpdateContext(context.Background(), value)
}

//SQL Executor
func (e *ComplexExec) SQLPrepare(query string) (*SQLStmt, error) {
	return newStmt(e.preparer, query, e.option)
}

//SQLPrepareContext
func (e *ComplexExec) SQLPrepareContext(ctx context.Context, query string) (*SQLStmt, error) {
	return newStmtContext(ctx, e.preparer, query, e.option)
}

//SQLFirstContext
func (e *ComplexExec) SQLFirstContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	stmt, err := e.SQLPrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			fmt.Printf("close stmt err:%v", err)
		}
	}()
	return stmt.FirstContext(ctx, dest, args...)
}

//SQLFirst
func (e *ComplexExec) SQLFirst(dest interface{}, query string, args ...interface{}) (err error) {
	return e.SQLFirstContext(context.Background(), dest, query, args...)
}

//SQLMustFirstContext
func (e *ComplexExec) SQLMustFirstContext(ctx context.Context, dest interface{}, query string, args ...interface{}) {
	err := e.SQLFirstContext(ctx, dest, query, args...)
	if err != nil {
		panic(err)
	}
}

//SQLMustFirst
func (e *ComplexExec) SQLMustFirst(dest interface{}, query string, args ...interface{}) {
	e.SQLMustFirstContext(context.Background(), dest, query, args...)
}

//SQLFindContext
func (e *ComplexExec) SQLFindContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	stmt, err := e.SQLPrepareContext(ctx, query)
	if err != nil {
		return
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			fmt.Printf("close stmt err:%v", err)
		}
	}()
	return stmt.FindContext(ctx, dest, args...)
}

//SQLFind
func (e *ComplexExec) SQLFind(dest interface{}, query string, args ...interface{}) (err error) {
	return e.SQLFindContext(context.Background(), dest, query, args...)
}

//SQLMustFindContext
func (e *ComplexExec) SQLMustFindContext(ctx context.Context, dest interface{}, query string, args ...interface{}) {
	err := e.SQLFindContext(ctx, dest, query, args...)
	if err != nil {
		panic(err)
	}
}

//SQLMustFind
func (e *ComplexExec) SQLMustFind(dest interface{}, query string, args ...interface{}) {
	e.SQLMustFindContext(context.Background(), dest, query, args...)
}

//SQLExecContext
func (e *ComplexExec) SQLExecContext(ctx context.Context, query string, args ...interface{}) (rs sql.Result, err error) {
	printSQL(query, args, e.option.out)
	return e.preparer.ExecContext(ctx, query, args...)
}

//SQLExec
func (e *ComplexExec) SQLExec(query string, args ...interface{}) (rs sql.Result, err error) {
	return e.SQLExecContext(context.Background(), query, args...)
}

// SQLMustExecContext is like SQLExecContext but panics if the query cannot be execute.
func (e *ComplexExec) SQLMustExecContext(ctx context.Context, query string, args ...interface{}) (rs sql.Result) {
	rs, err := e.SQLExecContext(ctx, query, args...)
	if err != nil {
		panic(err)
	}
	return
}

// SQLMustExec is like SQLExec but panics if the query cannot be execute.
func (e *ComplexExec) SQLMustExec(query string, args ...interface{}) (rs sql.Result) {
	return e.SQLMustExecContext(context.Background(), query, args...)
}

// DQLExecContext executes a DQL query without returning any rows.
// The argument are for any placeholder parameters in the query.
func (e *ComplexExec) DQLExecContext(ctx context.Context, query string, argument DQLArgument) (sql.Result, error) {
	query, args, err := DSLCompile(query, argument)
	if err != nil {
		return nil, err
	}
	return e.SQLExecContext(ctx, query, args...)
}

// DQLExec executes a DQL query without returning any rows.
// The argument are for any placeholder parameters in the query.
func (e *ComplexExec) DQLExec(query string, argument DQLArgument) (sql.Result, error) {
	return e.DQLExecContext(context.Background(), query, argument)
}

//DQLMustExecContext
func (e *ComplexExec) DQLMustExecContext(ctx context.Context, query string, argument DQLArgument) (r sql.Result) {
	r, err := e.DQLExecContext(ctx, query, argument)
	if err != nil {
		panic(err)
	}
	return
}

//DQLMustExec
func (e *ComplexExec) DQLMustExec(query string, argument DQLArgument) sql.Result {
	return e.DQLMustExecContext(context.Background(), query, argument)
}

//DQLFirstContext
func (e *ComplexExec) DQLFirstContext(ctx context.Context, dest interface{}, query string, argument DQLArgument) (err error) {
	query, args, err := DSLCompile(query, argument)
	err.Error()
	if err != nil {
		return
	}
	return e.SQLFirstContext(ctx, dest, query, args...)
}

//DQLFirst
func (e *ComplexExec) DQLFirst(dest interface{}, query string, argument DQLArgument) (err error) {
	return e.DQLFirstContext(context.Background(), dest, query, argument)
}

//DQLMustFirstContext
func (e *ComplexExec) DQLMustFirstContext(ctx context.Context, dest interface{}, query string, argument DQLArgument) {
	err := e.DQLFirstContext(ctx, dest, query, argument)
	if err != nil {
		panic(err)
	}
}

//DQLMustFirst
func (e *ComplexExec) DQLMustFirst(dest interface{}, query string, argument DQLArgument) {
	e.DQLMustFirstContext(context.Background(), dest, query, argument)
}

// DQLFindContext executes DQL that put rows into dest
func (e *ComplexExec) DQLFindContext(ctx context.Context, dest interface{}, query string, argument DQLArgument) (err error) {
	query, args, err := DSLCompile(query, argument)
	if err != nil {
		return
	}
	return e.SQLFindContext(ctx, dest, query, args...)
}

// DQLFind executes DQL that put rows into dest
func (e *ComplexExec) DQLFind(dest interface{}, query string, argument DQLArgument) (err error) {
	return e.DQLFindContext(context.Background(), dest, query, argument)
}

// DQLMustFindContext executes DQL that put rows into dest
func (e *ComplexExec) DQLMustFindContext(ctx context.Context, dest interface{}, query string, argument DQLArgument) {
	err := e.DQLFindContext(ctx, dest, query, argument)
	if err != nil {
		panic(err)
	}
}

// DQLMustFind executes a dbx query language that put rows into dest
func (e *ComplexExec) DQLMustFind(dest interface{}, query string, argument DQLArgument) {
	e.DQLMustFindContext(context.Background(), dest, query, argument)
}
