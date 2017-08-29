package yaorm

import (
	"github.com/Sirupsen/logrus"
)

type SqlExecutor struct {
	DB
}

func (e *SqlExecutor) SelectOne(holder interface{}, query string, args ...interface{}) error {
	logrus.WithFields(logrus.Fields{
		"sql_query":  query,
		"sql_params": args,
	}).Infof("SQL Query")
	return e.DB.SelectOne(holder, query, args...)
}

func (e *SqlExecutor) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	logrus.WithFields(logrus.Fields{
		"sql_query":  query,
		"sql_params": args,
	}).Infof("SQL Query")
	return e.DB.Select(i, query, args...)
}

/*
func (e *SqlExecutor) Get(i interface{}, keys ...interface{}) (interface{}, error) {}
func (e *SqlExecutor) Insert(list ...interface{}) error {}
func (e *SqlExecutor) Update(list ...interface{}) (int64, error) {}
func (e *SqlExecutor) Delete(list ...interface{}) (int64, error) {}
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
