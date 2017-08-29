package yaorm_test

import (
	"testing"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/testdata"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewTable(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{})
	assert.NotNil(t, table)
}

func TestTable_Fields(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{})
	assert.NotNil(t, table)
	assert.Len(t, table.Fields(), 4)
	assert.Equal(t, table.Fields()[0], "id")
	assert.Equal(t, table.Fields()[1], "name")
	assert.Equal(t, table.Fields()[2], "created_at")
	assert.Equal(t, table.Fields()[3], "updated_at")
}

func TestTable_Name(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{})
	assert.NotNil(t, table)
	assert.Equal(t, "category", table.Name())
}

func TestTable_NewModel(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{})
	assert.NotNil(t, table)
	m, err := table.NewModel()
	assert.Nil(t, err)
	assert.IsType(t, &testdata.Category{}, m)
}

func TestTable_NewFilter(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{}).WithFilter(testdata.NewCategoryFilter())
	assert.NotNil(t, table)
	f, err := table.NewFilter()
	assert.Nil(t, err)
	assert.IsType(t, &testdata.CategoryFilter{}, f)
}

func TestTable_NewSlice(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{})
	assert.NotNil(t, table)
	m, err := table.NewSlice()
	assert.Nil(t, err)
	assert.IsType(t, []*testdata.Category{}, m)
}

func TestTable_FieldIndex(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{})
	assert.NotNil(t, table)
	assert.Equal(t, 1, table.FieldIndex("id"))
	assert.Equal(t, -1, table.FieldIndex("unknown"))
}

func TestTable_FilterFieldIndex(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{}).WithFilter(testdata.NewCategoryFilter())
	assert.NotNil(t, table)
	assert.Equal(t, 1, table.FilterFieldIndex("id"))
	assert.Equal(t, -1, table.FilterFieldIndex("unknown"))
}

func TestGetTable(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{})
	assert.NotNil(t, table)

	retrievedTable, err := yaorm.GetTable("test", "category")
	assert.NotNil(t, retrievedTable)
	assert.Nil(t, err)
	assert.Equal(t, table, retrievedTable)

	notFoundTable, err := yaorm.GetTable("test", "category2")
	assert.Nil(t, notFoundTable)
	assert.NotNil(t, err)
	assert.True(t, errors.IsNotFound(err))

	notFoundDb, err := yaorm.GetTable("test2", "category2")
	assert.Nil(t, notFoundDb)
	assert.NotNil(t, err)
	assert.True(t, errors.IsNotFound(err))
}

func TestGetTableByFilter(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{}).WithFilter(&testdata.CategoryFilter{})
	assert.NotNil(t, table)

	foundTable, err := yaorm.GetTableByFilter(&testdata.CategoryFilter{})
	assert.Nil(t, err)
	assert.NotNil(t, foundTable)
	assert.Equal(t, table, foundTable)
}

func TestGetTableByModel(t *testing.T) {
	table := yaorm.NewTable("test", "category", &testdata.Category{}).WithFilter(&testdata.CategoryFilter{})
	assert.NotNil(t, table)

	foundTable, err := yaorm.GetTableByModel(&testdata.Category{})
	assert.Nil(t, err)
	assert.NotNil(t, foundTable)
	assert.Equal(t, table, foundTable)
}
