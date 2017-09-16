package yaorm_test

import (
	"testing"

	"github.com/geoffreybauduin/yaorm"
	"github.com/geoffreybauduin/yaorm/testdata"
	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestGenericSelectOne_NotFound(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	model, err := yaorm.GenericSelectOne(dbp, &testdata.CategoryFilter{})
	assert.NotNil(t, err)
	assert.True(t, errors.IsNotFound(err))
	assert.Nil(t, model)
}

type fakeFilter struct {
	*yaormfilter.ModelFilter
}

func TestGenericSelectOne_NotRegisteredFilter(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	model, err := yaorm.GenericSelectOne(dbp, &fakeFilter{})
	assert.NotNil(t, err)
	assert.Equal(t, yaorm.ErrTableNotFound, err)
	assert.Nil(t, model)
}

func TestGenericSelectOne(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	m := &testdata.Category{Name: "category"}
	m.SetDBP(dbp)
	err = yaorm.GenericSave(m)
	assert.Nil(t, err)
	modelFound, err := yaorm.GenericSelectOne(dbp, testdata.NewCategoryFilter().Name(yaormfilter.Equals("category")))
	assert.Nil(t, err)
	assert.Equal(t, m.ID, modelFound.(*testdata.Category).ID)
}

func TestGenericInsert(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	m := &testdata.Category{Name: "category"}
	m.SetDBP(dbp)
	beforeCreatedAt := m.CreatedAt
	assert.Equal(t, int64(0), m.ID)
	err = yaorm.GenericInsert(m)
	assert.Nil(t, err)
	assert.NotEqual(t, int64(0), m.ID)
	assert.NotEqual(t, beforeCreatedAt, m.CreatedAt)
}

func TestGenericUpdate(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	m := &testdata.Category{Name: "category"}
	m.SetDBP(dbp)
	assert.Equal(t, int64(0), m.ID)
	err = yaorm.GenericInsert(m)
	assert.Nil(t, err)
	assert.NotEqual(t, int64(0), m.ID)
	m.Name = "test2"
	beforeUpdatedAt := m.UpdatedAt
	beforeCreatedAt := m.CreatedAt
	err = yaorm.GenericUpdate(m)
	assert.Nil(t, err)
	assert.Equal(t, beforeCreatedAt, m.CreatedAt)
	assert.NotEqual(t, beforeUpdatedAt, m.UpdatedAt)
}

func TestGenericSave(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	m := &testdata.Category{Name: "category"}
	m.SetDBP(dbp)
	assert.Equal(t, int64(0), m.ID)
	err = yaorm.GenericSave(m)
	assert.Nil(t, err)
	assert.NotEqual(t, int64(0), m.ID)
	m.Name = "test2"
	err = yaorm.GenericSave(m)
	assert.Nil(t, err)
}

func TestGenericSelectAll_NotFound(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	models, err := yaorm.GenericSelectAll(dbp, &testdata.CategoryFilter{})
	assert.Nil(t, err)
	assert.Len(t, models, 0)
}

func TestGenericSelectAll_NotRegisteredFilter(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	models, err := yaorm.GenericSelectAll(dbp, &fakeFilter{})
	assert.NotNil(t, err)
	assert.Equal(t, yaorm.ErrTableNotFound, err)
	assert.Nil(t, models)
}

func TestGenericSelectAll(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	m := &testdata.Category{Name: "category"}
	m.SetDBP(dbp)
	err = yaorm.GenericSave(m)
	assert.Nil(t, err)
	modelsFound, err := yaorm.GenericSelectAll(dbp, testdata.NewCategoryFilter().Name(yaormfilter.Equals("category")))
	assert.Nil(t, err)
	assert.Len(t, modelsFound, 1)
	assert.Equal(t, m.ID, modelsFound[0].(*testdata.Category).ID)
}

func saveModel(t *testing.T, dbp yaorm.DBProvider, m yaorm.Model) {
	m.SetDBP(dbp)
	err := yaorm.GenericSave(m)
	assert.Nil(t, err)
}

func TestGenericSelectOne_WithJoinFilters(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)
	post := &testdata.Post{Subject: "subject", CategoryID: category2.ID}
	saveModel(t, dbp, post)
	modelFound, err := yaorm.GenericSelectOne(dbp, testdata.NewPostFilter().Category(
		testdata.NewCategoryFilter().Name(yaormfilter.Equals("category2")),
	))
	assert.Nil(t, err)
	assert.Equal(t, post.ID, modelFound.(*testdata.Post).ID)

	modelFound, err = yaorm.GenericSelectOne(dbp, testdata.NewPostFilter().Category(
		testdata.NewCategoryFilter().Name(yaormfilter.Equals("category")),
	))
	assert.NotNil(t, err)
	assert.True(t, errors.IsNotFound(err))
	assert.Nil(t, modelFound)
}

func TestGenericSelectOne_WithSubqueryload(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)
	post := &testdata.Post{Subject: "subject", CategoryID: category2.ID}
	saveModel(t, dbp, post)
	modelFound, err := yaorm.GenericSelectOne(dbp, testdata.NewPostFilter().Category(
		testdata.NewCategoryFilter().Subqueryload(),
	))
	assert.Nil(t, err)
	assert.Equal(t, post.ID, modelFound.(*testdata.Post).ID)
	assert.NotNil(t, modelFound.(*testdata.Post).Category)
	assert.Equal(t, category2.ID, modelFound.(*testdata.Post).Category.ID)
}

func TestGenericSelectAll_WithJoinFilters(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)
	post := &testdata.Post{Subject: "subject", CategoryID: category2.ID}
	saveModel(t, dbp, post)
	modelsFound, err := yaorm.GenericSelectAll(dbp, testdata.NewPostFilter().Category(
		testdata.NewCategoryFilter().Name(yaormfilter.Equals("category2")),
	))
	assert.Nil(t, err)
	assert.Len(t, modelsFound, 1)
	assert.Equal(t, post.ID, modelsFound[0].(*testdata.Post).ID)

	modelsFound, err = yaorm.GenericSelectAll(dbp, testdata.NewPostFilter().Category(
		testdata.NewCategoryFilter().Name(yaormfilter.Equals("category")),
	))
	assert.Nil(t, err)
	assert.Len(t, modelsFound, 0)
}

func TestGenericSelectAll_WithSubqueryload(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)
	post := &testdata.Post{Subject: "subject", CategoryID: category2.ID}
	saveModel(t, dbp, post)
	modelsFound, err := yaorm.GenericSelectAll(dbp, testdata.NewPostFilter().Category(
		testdata.NewCategoryFilter().Subqueryload(),
	))
	assert.Nil(t, err)
	assert.Len(t, modelsFound, 1)
	assert.Equal(t, post.ID, modelsFound[0].(*testdata.Post).ID)
	assert.NotNil(t, modelsFound[0].(*testdata.Post).Category)
	assert.Equal(t, category2.ID, modelsFound[0].(*testdata.Post).Category.ID)
}

func TestGenericSave_MultiplePrimaryKeys(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	post := &testdata.Post{Subject: "subject", CategoryID: category.ID}
	saveModel(t, dbp, post)
	tag := &testdata.Tag{Tag: "tag"}
	saveModel(t, dbp, tag)
	posttag := &testdata.PostTag{TagID: tag.ID, PostID: post.ID}
	oldCreated := posttag.CreatedAt
	posttag.SetDBP(dbp)
	err = yaorm.GenericSave(posttag)
	assert.Nil(t, err)
	assert.NotEqual(t, oldCreated, posttag.CreatedAt)
	assert.Equal(t, posttag.CreatedAt, posttag.UpdatedAt)
	found, err := yaorm.GenericSelectOne(dbp, testdata.NewPostTagFilter())
	assert.Nil(t, err)
	assert.Equal(t, posttag.PostID, found.(*testdata.PostTag).PostID)
	assert.Equal(t, posttag.TagID, found.(*testdata.PostTag).TagID)
	err = yaorm.GenericSave(posttag)
	assert.Nil(t, err)
	assert.NotEqual(t, posttag.CreatedAt, posttag.UpdatedAt)
}

func TestDatabaseModel_Load(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	category := &testdata.Category{Name: "category"}
	saveModel(t, dbp, category)
	category2 := &testdata.Category{Name: "category2"}
	saveModel(t, dbp, category2)
	foundCategory := &testdata.Category{ID: category.ID}
	err = foundCategory.Load(dbp)
	assert.Nil(t, err)
	assert.Equal(t, category.Name, foundCategory.Name)
}

func TestGenericSelectOne_WithSelectForUpdate(t *testing.T) {
	killDb, err := testdata.SetupTestDatabase("test")
	defer killDb()
	assert.Nil(t, err)
	dbp, err := yaorm.NewDBProvider("test")
	assert.Nil(t, err)
	m := &testdata.Category{Name: "category"}
	m.SetDBP(dbp)
	err = yaorm.GenericSave(m)
	assert.Nil(t, err)
	modelFound, err := yaorm.GenericSelectOne(dbp, testdata.NewCategoryFilter().Name(
		yaormfilter.Equals("category"),
	).AddOption(
		yaormfilter.RequestOptions.SelectForUpdate,
	))
	assert.Nil(t, err)
	assert.Equal(t, m.ID, modelFound.(*testdata.Category).ID)
}