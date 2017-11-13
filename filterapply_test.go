package yaorm_test

import (
	"testing"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/testdata"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/stretchr/testify/assert"
)

func TestFilterApplier_ApplyLt(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)

	models, err := yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().ID(yaormfilter.Lt(category2.ID)))
	assert.Nil(t, err)
	assert.Len(t, models, 1)
	assert.Equal(t, models[0].(*testdata.Category).ID, category.ID)
}

func TestFilterApplier_ApplyLte(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)

	models, err := yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().ID(yaormfilter.Lte(category2.ID)))
	assert.Nil(t, err)
	assert.Len(t, models, 2)
}

func TestFilterApplier_ApplyGt(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)

	models, err := yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().ID(yaormfilter.Gt(category.ID)))
	assert.Nil(t, err)
	assert.Len(t, models, 1)
	assert.Equal(t, models[0].(*testdata.Category).ID, category2.ID)
}

func TestFilterApplier_ApplyGte(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)

	models, err := yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().ID(yaormfilter.Gte(category.ID)))
	assert.Nil(t, err)
	assert.Len(t, models, 2)
}
