package yaorm

import (
	"context"

	"github.com/go-gorp/gorp"
)

type SqlExecutor struct {
	gorp.SqlExecutor
	db  DB
	ctx context.Context
	dbp DBProvider
}

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

func (e *SqlExecutor) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	hook := e.db.ExecutorHook()
	hook.BeforeSelect(e.ctx, query, args...)
	v, err := e.SqlExecutor.Select(i, query, args...)
	hook.AfterSelect(e.ctx, query, args...)
	return v, err
}

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

/*
func (e *SqlExecutor) Get(i interface{}, keys ...interface{}) (interface{}, error) {}
func (e *SqlExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {}
func (e *SqlExecutor) SelectInt(query string, args ...interface{}) (int64, error)                 {}
func (e *SqlExecutor) SelectNullInt(query string, args ...interface{}) (sql.NullInt64, error)     {}
func (e *SqlExecutor) SelectFloat(query string, args ...interface{}) (float64, error)             {}
func (e *SqlExecutor) SelectNullFloat(query string, args ...interface{}) (sql.NullFloat64, error) {}
func (e *SqlExecutor) SelectStr(query string, args ...interface{}) (string, error)                {}
func (e *SqlExecutor) SelectNullStr(query string, args ...interface{}) (sql.NullString, error)    {}
func (e *SqlExecutor) Query(query string, args ...interface{}) (*sql.Rows, error)                 {}
func (e *SqlExecutor) QueryRow(query string, args ...interface{}) *sql.Row                        {}
*/
