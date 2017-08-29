package yaorm

import (
	"fmt"
	"sync"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/loopfz/gadgeto/zesty"
	"github.com/geoffreybauduin/yaorm/_vendor/github.com/loopfz/gadgeto/zesty/utils/rekordo"
	"github.com/go-gorp/gorp"
	"github.com/juju/errors"
)

var dblock sync.RWMutex

type DB interface {
	zesty.DB
	System() DMS
}

type db struct {
	zesty.DB
	system DMS
}

func (d db) System() DMS {
	return d.system
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
	AutoCreateTables bool
	// Dialect database dialect, leave empty for automatic guessing
	Dialect gorp.Dialect
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

	registry[config.Name] = db{DB: dbHandler, system: config.System}

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

//TypeConverter defines convertion from db to golang
type TypeConverter struct{}

//ToDb converts to database
func (tc TypeConverter) ToDb(val interface{}) (interface{}, error) {
	return val, nil
}

//FromDb converts to golang
func (tc TypeConverter) FromDb(target interface{}) (gorp.CustomScanner, bool) {
	return gorp.CustomScanner{}, false
}
