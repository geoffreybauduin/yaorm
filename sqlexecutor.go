package yaorm

import (
	"context"
	"database/sql"

	"github.com/go-gorp/gorp"
)

// SqlExecutor is a custom SQL Executor, on top of the one provided by gorp
// used to provide multiple hooks before executing statements
type SqlExecutor struct {
	gorp.SqlExecutor
	db  DB
	ctx context.Context
	dbp DBProvider
}

// SelectOne is a handler to select only 1 row from database and store it inside the first argument
func (e *SqlExecutor) SelectOne(holder interface{}, query string, args ...interface{}) error {
	hook := e.db.ExecutorHook()
	hook.BeforeSelectOne(e.ctx, query, args...)
	err := e.SqlExecutor.SelectOne(holder, query, args...)
	hook.AfterSelectOne(e.ctx, query, args...)
	return err
}

type queryArgs struct {
	query string
	args  []interface{}
}

// Select is a handler to select multiple rows from database and return them
func (e *SqlExecutor) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	hook := e.db.ExecutorHook()
	hook.BeforeSelect(e.ctx, query, args...)
	v, err := e.SqlExecutor.Select(i, query, args...)
	hook.AfterSelect(e.ctx, query, args...)
	return v, err
}

// Insert is a handler to insert a list of models inside the database
func (e *SqlExecutor) Insert(list ...interface{}) error {
	hook := e.db.ExecutorHook()
	for _, item := range list {
		var qArg queryArgs
		builder, err := buildInsert(e.dbp, item.(Model))
		if err != nil {
			return err
		}
		qArg.query, qArg.args, err = builder.ToSql()
		if err != nil {
			return err
		}
		hook.BeforeInsert(e.ctx, qArg.query, qArg.args...)
		err = e.SqlExecutor.Insert(item)
		if err != nil {
			return err
		}
		hook.AfterInsert(e.ctx, qArg.query, qArg.args...)
	}
	return nil
}

// Update is a handler to update a list of models inside the database
func (e *SqlExecutor) Update(list ...interface{}) (int64, error) {
	hook := e.db.ExecutorHook()
	var rv int64
	for _, item := range list {
		var qArg queryArgs
		builder, err := buildUpdate(e.dbp, item.(Model))
		if err != nil {
			return rv, err
		}
		qArg.query, qArg.args, err = builder.ToSql()
		if err != nil {
			return rv, err
		}
		hook.BeforeUpdate(e.ctx, qArg.query, qArg.args...)
		v, err := e.SqlExecutor.Update(item)
		if err != nil {
			return rv, err
		}
		rv += v
		hook.AfterUpdate(e.ctx, qArg.query, qArg.args...)
	}
	return rv, nil
}

// Delete is a handler to delete a list of models from the database
func (e *SqlExecutor) Delete(list ...interface{}) (int64, error) {
	hook := e.db.ExecutorHook()
	var rv int64
	for _, item := range list {
		var qArg queryArgs
		builder, err := buildDelete(e.dbp, item.(Model))
		if err != nil {
			return rv, err
		}
		qArg.query, qArg.args, err = builder.ToSql()
		if err != nil {
			return rv, err
		}
		hook.BeforeDelete(e.ctx, qArg.query, qArg.args...)
		v, err := e.SqlExecutor.Delete(item)
		if err != nil {
			return rv, err
		}
		rv += v
		hook.AfterDelete(e.ctx, qArg.query, qArg.args...)
	}
	return rv, nil
}

// Exec is a handler to execute a SQL query
func (e *SqlExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {
	hook := e.db.ExecutorHook()
	hook.BeforeExec(e.ctx, query, args...)
	v, err := e.SqlExecutor.Exec(query, args...)
	hook.AfterExec(e.ctx, query, args...)
	return v, err
}

/*
func (e *SqlExecutor) Get(i interface{}, keys ...interface{}) (interface{}, error) {}
func (e *SqlExecutor) SelectInt(query string, args ...interface{}) (int64, error)                 {}
func (e *SqlExecutor) SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error)     {}
func (e *SqlExecutor) SelectFloat(query string, args ...interface{}) (float64, error)             {}
func (e *SqlExecutor) SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error) {}
func (e *SqlExecutor) SelectStr(query string, args ...interface{}) (string, error)                {}
func (e *SqlExecutor) SelectNullStr(query string, args ...interface{}) (sql.NullString, error)    {}
func (e *SqlExecutor) Query(query string, args ...interface{}) (*sql.Rows, error)                 {}
func (e *SqlExecutor) QueryRow(query string, args ...interface{}) *sql.Row                        {}
*/
