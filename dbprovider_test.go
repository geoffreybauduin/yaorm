package yaorm_test

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/testdata"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestDbprovider_Tx(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)
	err = dbp.Tx()
	assert.Nil(t, err)
	m := &testdata.Category{Name: "category"}
	m.SetDBP(dbp)
	assert.Equal(t, int64(0), m.ID)
	err = yaorm.GenericInsert(m)
	assert.Nil(t, err)
	dbp2, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)

	// does not exist outside of tx
	_, err = yaorm.GenericSelectOne(dbp2, testdata.NewCategoryFilter().ID(yaormfilter.Equals(m.ID)))
	assert.NotNil(t, err)
	assert.True(t, errors.IsNotFound(err))

	// exists in tx
	_, err = yaorm.GenericSelectOne(dbp, testdata.NewCategoryFilter().ID(yaormfilter.Equals(m.ID)))
	assert.Nil(t, err)

	err = dbp.Commit()
	assert.Nil(t, err)

	// now it exists outside of tx
	_, err = yaorm.GenericSelectOne(dbp2, testdata.NewCategoryFilter().ID(yaormfilter.Equals(m.ID)))
	assert.Nil(t, err)
}

func TestDbprovider_Rollback(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)
	err = dbp.Tx()
	assert.Nil(t, err)
	m := &testdata.Category{Name: "category"}
	m.SetDBP(dbp)
	assert.Equal(t, int64(0), m.ID)
	err = yaorm.GenericInsert(m)
	assert.Nil(t, err)
	dbp2, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)

	// does not exist outside of tx
	_, err = yaorm.GenericSelectOne(dbp2, testdata.NewCategoryFilter().ID(yaormfilter.Equals(m.ID)))
	assert.NotNil(t, err)
	assert.True(t, errors.IsNotFound(err))

	// exists in tx
	_, err = yaorm.GenericSelectOne(dbp, testdata.NewCategoryFilter().ID(yaormfilter.Equals(m.ID)))
	assert.Nil(t, err)

	err = dbp.Rollback()
	assert.Nil(t, err)

	// still does not exist outside of tx
	_, err = yaorm.GenericSelectOne(dbp2, testdata.NewCategoryFilter().ID(yaormfilter.Equals(m.ID)))
	assert.NotNil(t, err)
	assert.True(t, errors.IsNotFound(err))

	// inside old tx it does not exist either
	_, err = yaorm.GenericSelectOne(dbp, testdata.NewCategoryFilter().ID(yaormfilter.Equals(m.ID)))
	assert.NotNil(t, err)
	assert.True(t, errors.IsNotFound(err))
}

func TestDBProvider_RunInTransaction(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.NoError(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.NoError(t, err)
	errTx := dbp.RunInTransaction(func() error {
		return nil
	})
	assert.NoError(t, errTx)
	errTx2 := dbp.RunInTransaction(func() error {
		return fmt.Errorf("test error")
	})
	assert.Error(t, errTx2)
	assert.Equal(t, "test error", errTx2.Error())
}

func TestDBProvider_Postgres_OnSessionCreated(t *testing.T) {
	if os.Getenv("DB") != "postgres" {
		return
	}
	exec := &pgSpecificExecutor{
		queries: make([]logParams, 0),
	}
	killDb, err := testdata.SetupPostgres("test", yaorm.PostgresSpecific{IntervalStyle: "iso_8601"}, exec)
	defer killDb()
	assert.NoError(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.NoError(t, err)
	assert.Len(t, exec.queries, 1)
	assert.True(t, reflect.DeepEqual(exec.queries, []logParams{
		{
			query: "SET intervalstyle = 'iso_8601'",
			args:  []interface{}{},
		},
	}))
	_, err = dbp.DB().Exec("SELECT 1")
	assert.NoError(t, err)
	assert.Len(t, exec.queries, 2)
	assert.True(t, reflect.DeepEqual(exec.queries, []logParams{
		{
			query: "SET intervalstyle = 'iso_8601'",
			args:  []interface{}{},
		},
		{
			query: "SELECT 1",
			args:  nil,
		},
	}))
}

type pgSpecificExecutor struct {
	yaorm.DefaultExecutorHook
	queries []logParams
}

type logParams struct {
	query string
	args  []interface{}
}

func (exec *pgSpecificExecutor) BeforeExec(ctx context.Context, query string, args ...interface{}) {
	exec.queries = append(exec.queries, logParams{
		query: query,
		args:  args,
	})
}
