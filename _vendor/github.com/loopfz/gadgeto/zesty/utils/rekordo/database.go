package rekordo

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/loopfz/gadgeto/zesty"
	"github.com/go-gorp/gorp"
)

// Default database settings.
const (
	maxOpenConns = 5
	maxIdleConns = 3
)

// DatabaseConfig represents the configuration used to
// register a new database.
type DatabaseConfig struct {
	Name             string
	DSN              string
	System           DBMS
	MaxOpenConns     int
	MaxIdleConns     int
	AutoCreateTables bool
	ConnMaxLifetime  time.Duration
	ConnMaxIdleTime  time.Duration
}

// RegisterDatabase creates a gorp map with tables and tc and
// registers it with zesty.
func RegisterDatabase(db *DatabaseConfig, tc gorp.TypeConverter, dialect gorp.Dialect) error {
	dbConn, err := sql.Open(db.System.DriverName(), db.DSN)
	if err != nil {
		return err
	}
	// Make sure we have proper values for the database
	// settings, and replace them with default if necessary
	// before applying to the new connection.
	if db.MaxOpenConns == 0 {
		db.MaxOpenConns = maxOpenConns
	}
	dbConn.SetMaxOpenConns(db.MaxOpenConns)
	if db.MaxIdleConns == 0 {
		db.MaxIdleConns = maxIdleConns
	}
	dbConn.SetMaxIdleConns(db.MaxIdleConns)
	dbConn.SetConnMaxLifetime(db.ConnMaxLifetime)
	dbConn.SetConnMaxIdleTime(db.ConnMaxIdleTime)

	// Select the proper dialect used by gorp.
	if dialect == nil {
		switch db.System {
		case DatabaseMySQL:
			dialect = gorp.MySQLDialect{}
		case DatabasePostgreSQL:
			dialect = gorp.PostgresDialect{}
		case DatabaseSqlite3:
			dialect = gorp.SqliteDialect{}
		default:
			return errors.New("unknown database system")
		}
	}
	dbmap := &gorp.DbMap{
		Db:            dbConn,
		Dialect:       dialect,
		TypeConverter: tc,
	}
	modelsMu.Lock()
	tableModels := models[db.Name]
	for _, t := range tableModels {
		// using gorp dialect to know if the driver supports schemas
		if t.Schema == "" || !strings.Contains(dialect.QuotedTableForQuery("schema", "table"), "schema") {
			dbmap.AddTableWithName(t.Model, t.Name).SetKeys(t.AutoIncrement, t.Keys...)
		} else {
			dbmap.AddTableWithNameAndSchema(t.Model, t.Schema, t.Name).SetKeys(t.AutoIncrement, t.Keys...)
		}
	}
	modelsMu.Unlock()

	if db.AutoCreateTables {
		err = dbmap.CreateTablesIfNotExists()
		if err != nil {
			return err
		}
	}
	return zesty.RegisterDB(zesty.NewDB(dbmap), db.Name)
}

// DBMS represents a database management system.
type DBMS uint8

// Database management systems.
const (
	DatabasePostgreSQL DBMS = iota ^ 42
	DatabaseMySQL
	DatabaseSqlite3
)

// DriverName returns the name of the driver for ds.
func (d DBMS) DriverName() string {
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
