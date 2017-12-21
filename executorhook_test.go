package yaorm_test

import (
	"context"
	"os"
	"testing"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/stretchr/testify/assert"
)

type customExecutorHookForTesting struct {
	yaorm.DefaultExecutorHook
	query string
	args  []interface{}
}

var query_ string
var args_ []interface{}

func register(query string, args ...interface{}) {
	query_ = query
	args_ = args
}

func (h *customExecutorHookForTesting) BeforeSelectOne(ctx context.Context, query string, args ...interface{}) {
	register(query, args...)
}
func (h *customExecutorHookForTesting) AfterSelectOne(ctx context.Context, query string, args ...interface{}) {
	register(query, args...)
}
func (h *customExecutorHookForTesting) BeforeSelect(ctx context.Context, query string, args ...interface{}) {
	register(query, args...)
}
func (h *customExecutorHookForTesting) AfterSelect(ctx context.Context, query string, args ...interface{}) {
	register(query, args...)
}
func (h *customExecutorHookForTesting) BeforeInsert(ctx context.Context, query string, args ...interface{}) {
	register(query, args...)
}
func (h *customExecutorHookForTesting) BeforeUpdate(ctx context.Context, query string, args ...interface{}) {
	register(query, args...)
}
func (h *customExecutorHookForTesting) BeforeDelete(ctx context.Context, query string, args ...interface{}) {
	register(query, args...)
}

type fakeModel struct {
	yaorm.DatabaseModel
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type fakeModelFilter struct {
	yaormfilter.ModelFilter
}

func TestExecutorHook_BeforeSelectOne(t *testing.T) {
	defer func() {
		os.Remove("/tmp/test_test.sqlite")
		yaorm.UnregisterDB("test")
	}()
	yaorm.NewTable("test", "model", &fakeModel{}).WithFilter(&fakeModelFilter{})
	err := yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
		Name:             "test",
		DSN:              "/tmp/test_test.sqlite",
		System:           yaorm.DatabaseSqlite3,
		AutoCreateTables: true,
		ExecutorHook:     &customExecutorHookForTesting{},
	})
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)
	yaorm.GenericSelectOne(dbp, &fakeModelFilter{})
	assert.Equal(t, `SELECT "model"."id", "model"."name" FROM "model" AS "model"`, query_)
	assert.Len(t, args_, 0)
}

func TestExecutorHook_BeforeInsert(t *testing.T) {
	defer func() {
		os.Remove("/tmp/test_test.sqlite")
		yaorm.UnregisterDB("test")
	}()
	yaorm.NewTable("test", "model", &fakeModel{}).WithFilter(&fakeModelFilter{})
	err := yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
		Name:             "test",
		DSN:              "/tmp/test_test.sqlite",
		System:           yaorm.DatabaseSqlite3,
		AutoCreateTables: true,
		ExecutorHook:     &customExecutorHookForTesting{},
	})
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)
	m := &fakeModel{Name: "test"}
	m.SetDBP(dbp)
	err = yaorm.GenericSave(m)
	assert.Nil(t, err)
	assert.Equal(t, `INSERT INTO model (name) VALUES ($1)`, query_)
	assert.Len(t, args_, 1)
	assert.Equal(t, args_[0], "test")
}

func TestExecutorHook_BeforeUpdate(t *testing.T) {
	defer func() {
		os.Remove("/tmp/test_test.sqlite")
		yaorm.UnregisterDB("test")
	}()
	yaorm.NewTable("test", "model", &fakeModel{}).WithFilter(&fakeModelFilter{})
	err := yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
		Name:             "test",
		DSN:              "/tmp/test_test.sqlite",
		System:           yaorm.DatabaseSqlite3,
		AutoCreateTables: true,
		ExecutorHook:     &customExecutorHookForTesting{},
	})
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)
	m := &fakeModel{Name: "test"}
	m.SetDBP(dbp)
	err = yaorm.GenericSave(m)
	assert.Nil(t, err)
	err = yaorm.GenericSave(m)
	assert.Nil(t, err)
	assert.Equal(t, `UPDATE model SET name = $1 WHERE id = $2`, query_)
	assert.Len(t, args_, 2)
	assert.Equal(t, args_[0], "test")
	assert.Equal(t, args_[1], m.ID)
}
