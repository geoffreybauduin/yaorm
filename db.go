package yaorm

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/loopfz/gadgeto/zesty"
	"github.com/geoffreybauduin/yaorm/_vendor/github.com/loopfz/gadgeto/zesty/utils/rekordo"
	"github.com/geoffreybauduin/yaorm/tools"
	"github.com/go-gorp/gorp"
	"github.com/juju/errors"
)

var dblock sync.RWMutex

type DB interface {
	zesty.DB
	System() DMS
	ExecutorHook() ExecutorHook
}

type db struct {
	zesty.DB
	dbmap        *gorp.DbMap
	system       DMS
	executorHook ExecutorHook
}

func (d *db) System() DMS {
	return d.system
}

func (d *db) ExecutorHook() ExecutorHook {
	if d.executorHook == nil {
		return &DefaultExecutorHook{}
	}
	return reflect.New(tools.GetNonPtrValue(d.executorHook).Type()).Interface().(ExecutorHook)
}

var (
	ErrDatabaseConflict = errors.Errorf("Database name conflicts with existing")
	registry            = map[string]DB{}
)

// DatabaseConfiguration configures a database
type DatabaseConfiguration struct {
	Name             string
	System           DMS
	DSN              string
	MaxOpenConns     int
	MaxIdleConns     int
	ConnMaxLifetime  time.Duration
	AutoCreateTables bool
	// Dialect database dialect, leave empty for automatic guessing
	Dialect gorp.Dialect
	// ExecutorHook is a configurable hook to add logs, for example, to your sql requests
	ExecutorHook ExecutorHook
}

// GetDB returns a database object from its name
func GetDB(name string) (DB, error) {
	db, ok := registry[name]
	if !ok {
		return nil, errors.NotFoundf("database %s", name)
	}
	return db, nil
}

// RegisterDB creates a new database with configuration
func RegisterDB(config *DatabaseConfiguration) error {
	dblock.Lock()
	defer dblock.Unlock()

	if _, ok := registry[config.Name]; ok {
		return ErrDatabaseConflict
	}

	databaseConfiguration := &rekordo.DatabaseConfig{
		Name:             config.Name,
		DSN:              config.DSN,
		System:           config.System.RekordoValue(),
		AutoCreateTables: config.AutoCreateTables,
		MaxIdleConns:     config.MaxIdleConns,
		MaxOpenConns:     config.MaxOpenConns,
		ConnMaxLifetime:  config.ConnMaxLifetime,
	}

	// Register database.
	err := rekordo.RegisterDatabase(databaseConfiguration, TypeConverter{}, config.Dialect)
	if err != nil {
		return err
	}

	// Get database handler
	dbHandler, err := zesty.GetDB(config.Name)
	if err != nil {
		return err
	}

	registry[config.Name] = &db{
		DB:           dbHandler,
		system:       config.System,
		executorHook: config.ExecutorHook,
	}

	return nil
}

// UnregisterDB removes the database from the registry
func UnregisterDB(name string) error {
	dblock.Lock()
	defer dblock.Unlock()

	_, ok := registry[name]
	if !ok {
		return ErrDbNotFound
	}

	delete(registry, name)

	err := zesty.UnregisterDB(name)
	if err != nil {
		return err
	}

	return nil
}

// DMS represents a database management system
type DMS uint8

// Database management systems.
const (
	DatabasePostgreSQL DMS = iota ^ 42
	DatabaseMySQL
	DatabaseSqlite3
)

// DriverName returns the name of the driver for ds.
func (d DMS) DriverName() string {
	switch d {
	case DatabasePostgreSQL:
		return "postgres"
	case DatabaseMySQL:
		return "mysql"
	case DatabaseSqlite3:
		return "sqlite3"
	}
	return ""
}

// RekordoValue returns the rekordo value
func (d DMS) RekordoValue() rekordo.DBMS {
	switch d {
	case DatabasePostgreSQL:
		return rekordo.DatabasePostgreSQL
	case DatabaseMySQL:
		return rekordo.DatabaseMySQL
	case DatabaseSqlite3:
		return rekordo.DatabaseSqlite3
	}
	panic(fmt.Errorf("Invalid system"))
}
