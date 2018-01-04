package yaorm_test

import (
	"context"
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
