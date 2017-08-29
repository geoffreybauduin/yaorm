package yaorm

import (
	"github.com/go-gorp/gorp"
	"github.com/geoffreybauduin/yaorm/_vendor/github.com/loopfz/gadgeto/zesty"
)

// DBProvider provides an abstracted way of accessing the database
type DBProvider interface {
	zesty.DBProvider
}

type dbprovider struct {
	zesty.DBProvider
}

// NewDBProvider creates a new db provider
func NewDBProvider(name string) (DBProvider, error) {
	dblock.RLock()
	defer dblock.RUnlock()
	dbp, err := zesty.NewDBProvider(name)
	if err != nil {
		return nil, err
	}
	return &dbprovider{dbp}, nil
}

// DB returns a SQL Executor interface
func (dbp *dbprovider) DB() gorp.SqlExecutor {
	return &SqlExecutor{dbp.DBProvider.DB().(DB)}
}
