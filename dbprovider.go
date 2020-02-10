package yaorm

import (
	"context"

	"github.com/geoffreybauduin/yaorm/_vendor/github.com/lann/squirrel"
	"github.com/geoffreybauduin/yaorm/_vendor/github.com/loopfz/gadgeto/zesty"
	"github.com/geoffreybauduin/yaorm/_vendor/github.com/satori/go.uuid"
	"github.com/geoffreybauduin/yaorm/tools"
	"github.com/go-gorp/gorp"
)

// DBProvider provides an abstracted way of accessing the database
type DBProvider interface {
	zesty.DBProvider
	EscapeValue(value string) string
	CanSelectForUpdate() bool
	getStatementGenerator() squirrel.StatementBuilderType
	Context() context.Context
	UUID() string
	getDialect() gorp.Dialect
	HasCapacity(capacity DatabaseCapacity) bool
	RunInTransaction(func() error) error
}

type dbprovider struct {
	zesty.DBProvider
	name string
	ctx  context.Context
	uuid string
}

// NewDBProvider creates a new db provider
func NewDBProvider(ctx context.Context, name string) (DBProvider, error) {
	dblock.RLock()
	defer dblock.RUnlock()
	dbp, err := zesty.NewDBProvider(name)
	if err != nil {
		return nil, err
	}
	uuid4 := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &dbprovider{DBProvider: dbp, name: name, ctx: ctx, uuid: uuid4.String()}, nil
}

// DB returns a SQL Executor interface
func (dbp *dbprovider) DB() gorp.SqlExecutor {
	dbRetrieved := dbp.getDb()
	dbUsed := dbp.DBProvider.DB()
	return &SqlExecutor{
		SqlExecutor: dbUsed,
		db:          dbRetrieved,
		ctx:         dbp.Context(),
		dbp:         dbp,
	}
}

// EscapeValue escapes the value sent according to the dialect
func (dbp *dbprovider) EscapeValue(value string) string {
	return dbp.getDialect().QuoteField(value)
}

// CanSelectForUpdate returns true if the current dialect can perform select for update statements
func (dbp *dbprovider) CanSelectForUpdate() bool {
	db := registry[dbp.name]
	switch db.System() {
	case DatabasePostgreSQL:
		return true
	}
	return false
}

func (dbp *dbprovider) getDb() DB {
	return registry[dbp.name]
}

func (dbp *dbprovider) getDialect() gorp.Dialect {
	v := tools.GetNonPtrValue(dbp.getDb())
	dbField := tools.GetNonPtrValue(v.FieldByName("DB").Interface())
	field := dbField.FieldByName("DbMap")
	s := field.Interface().(*gorp.DbMap)
	return s.Dialect
}

func (dbp *dbprovider) getStatementGenerator() squirrel.StatementBuilderType {
	switch dbp.getDb().System() {
	case DatabaseMySQL:
		return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)
	}
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

// Context returns the context stored inside this DBProvider instance
func (dbp *dbprovider) Context() context.Context {
	return dbp.ctx
}

// UUID returns an unique identifier for this DBProvider instance
func (dbp *dbprovider) UUID() string {
	return dbp.uuid
}

// HasCapacity returns true if used database has the provided capacity
func (dbp *dbprovider) HasCapacity(capacity DatabaseCapacity) bool {
	system := dbp.getDb().System()
	if _, ok := databaseCapacities[system]; !ok {
		return false
	}
	return databaseCapacities[system][capacity]
}

// RunInTraction will run the provided function inside a transaction.
// if an error occurs, the transaction is automatically rolled back.
// at the end of the transaction, the transaction is commit inside the
// dbms
func (dbp *dbprovider) RunInTransaction(fn func() error) error {
	shouldRollback := true
	errTx := dbp.Tx()
	if errTx != nil {
		return errTx
	}
	defer func() {
		if !shouldRollback {
			return
		}
		dbp.Rollback() //nolint:errcheck
	}()
	errFn := fn()
	if errFn != nil {
		return errFn
	}
	errCommit := dbp.Commit()
	if errCommit != nil {
		return errCommit
	}
	shouldRollback = false
	return nil
}
