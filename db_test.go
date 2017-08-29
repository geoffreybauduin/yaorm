package yaorm_test

import (
	"os"
	"testing"

	"github.com/geoffreybauduin/yaorm"
	_ "github.com/geoffreybauduin/yaorm/testdata"
	"github.com/stretchr/testify/assert"
)

func TestRegisterDB(t *testing.T) {
	defer func() {
		os.Remove("/tmp/test_test.sqlite")
		yaorm.UnregisterDB("test")
	}()
	err := yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
		Name:             "test",
		DSN:              "/tmp/test_test.sqlite",
		System:           yaorm.DatabaseSqlite3,
		AutoCreateTables: true,
	})
	assert.Nil(t, err)
}

func TestRegisterDB_Conflicts(t *testing.T) {
	defer func() {
		os.Remove("/tmp/test_test.sqlite")
		yaorm.UnregisterDB("test")
	}()
	err := yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
		Name:             "test",
		DSN:              "/tmp/test_test.sqlite",
		System:           yaorm.DatabaseSqlite3,
		AutoCreateTables: true,
	})
	assert.Nil(t, err)
	err = yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
		Name:             "test",
		DSN:              "/tmp/test_test.sqlite",
		System:           yaorm.DatabaseSqlite3,
		AutoCreateTables: true,
	})
	assert.Equal(t, yaorm.ErrDatabaseConflict, err)
}

func TestUnregisterDB(t *testing.T) {
	defer func() {
		os.Remove("/tmp/test_test.sqlite")
		yaorm.UnregisterDB("test")
	}()
	err := yaorm.RegisterDB(&yaorm.DatabaseConfiguration{
		Name:             "test",
		DSN:              "/tmp/test_test.sqlite",
		System:           yaorm.DatabaseSqlite3,
		AutoCreateTables: true,
	})
	assert.Nil(t, err)

	err = yaorm.UnregisterDB("test")
	assert.Nil(t, err)
}

func TestUnregisterDB_NotFound(t *testing.T) {
	err := yaorm.UnregisterDB("test")
	assert.Equal(t, yaorm.ErrDbNotFound, err)
}
