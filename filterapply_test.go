package yaorm_test

import (
	"context"
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
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
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
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
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
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
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
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)

	models, err := yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().ID(yaormfilter.Gte(category.ID)))
	assert.Nil(t, err)
	assert.Len(t, models, 2)
}

func TestFilterApplier_ApplyWithOrderBy(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)

	models, err := yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().OrderBy("id", yaormfilter.OrderingWays.Desc))
	assert.Nil(t, err)
	assert.Len(t, models, 2)
	assert.Equal(t, models[0].(*testdata.Category).ID, category2.ID)
	assert.Equal(t, models[1].(*testdata.Category).ID, category.ID)

	models, err = yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().OrderBy("id", yaormfilter.OrderingWays.Asc))
	assert.Nil(t, err)
	assert.Len(t, models, 2)
	assert.Equal(t, models[0].(*testdata.Category).ID, category.ID)
	assert.Equal(t, models[1].(*testdata.Category).ID, category2.ID)
}

func TestFilterApplier_ApplyWithLimitAndOffset(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)

	models, err := yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().Limit(1))
	assert.Nil(t, err)
	assert.Len(t, models, 1)
	assert.Equal(t, models[0].(*testdata.Category).ID, category.ID)

	models, err = yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().Limit(1).Offset(1))
	assert.Nil(t, err)
	assert.Len(t, models, 1)
	assert.Equal(t, models[0].(*testdata.Category).ID, category2.ID)
}

func TestFilterApplier_LoadColumns(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()

	assert.NoError(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.NoError(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)

	models, err := yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter())
	assert.NoError(t, err)
	if assert.Len(t, models, 2) {
		assert.Equal(t, category.ID, models[0].(*testdata.Category).ID)
		assert.Equal(t, category.Name, models[0].(*testdata.Category).Name)
		assert.False(t, models[0].(*testdata.Category).CreatedAt.IsZero())
		assert.False(t, models[0].(*testdata.Category).UpdatedAt.IsZero())

		assert.Equal(t, category2.ID, models[1].(*testdata.Category).ID)
		assert.Equal(t, category2.Name, models[1].(*testdata.Category).Name)
		assert.False(t, models[1].(*testdata.Category).CreatedAt.IsZero())
		assert.False(t, models[1].(*testdata.Category).UpdatedAt.IsZero())
	}

	// Only Name, CreatedAt, UpdatedAt BUT without CreatedAt
	filter := testdata.NewCategoryFilter()
	filter.LoadColumns("name", "created_at", "updated_at")
	filter.DontLoadColumns("created_at")
	models, err = yaorm.GenericSelectAll(dbp, filter)
	assert.NoError(t, err)
	if assert.Len(t, models, 2) {
		assert.Zero(t, models[0].(*testdata.Category).ID) // not loaded
		assert.Equal(t, category.Name, models[0].(*testdata.Category).Name)
		assert.True(t, models[0].(*testdata.Category).CreatedAt.IsZero()) // not loaded
		assert.False(t, models[0].(*testdata.Category).UpdatedAt.IsZero())

		assert.Zero(t, models[1].(*testdata.Category).ID) // not loaded
		assert.Equal(t, category2.Name, models[1].(*testdata.Category).Name)
		assert.True(t, models[1].(*testdata.Category).CreatedAt.IsZero()) // not loaded
		assert.False(t, models[1].(*testdata.Category).UpdatedAt.IsZero())
	}
}

func TestFilterApplier_Distinct(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()

	assert.NoError(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.NoError(t, err)
	post := &testdata.Post{}
	saveModel(t, dbp, post)
	postChild := &testdata.Post{ParentPostID: post.ID}
	saveModel(t, dbp, postChild)
	postChild2 := &testdata.Post{ParentPostID: post.ID}
	saveModel(t, dbp, postChild2)

	f := testdata.NewPostFilter()
	f.ChildrenPosts(testdata.NewPostFilter().
		ParentPostID(yaormfilter.NewInt64Filter().Equals(post.ID)))
	models, err := yaorm.GenericSelectAll(dbp, f)
	assert.NoError(t, err)
	// We do not set the filter Distinct, so we should have one row per matching children
	assert.Equal(t, 2, len(models))

	f.Distinct()
	models, err = yaorm.GenericSelectAll(dbp, f)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(models))
}

func TestFilterApplier_ApplyNotIn(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider(context.TODO(), "test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)

	models, err := yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().ID(yaormfilter.NotIn(category2.ID)))
	assert.Nil(t, err)
	assert.Len(t, models, 1)
	assert.Equal(t, models[0].(*testdata.Category).ID, category.ID)
}
