package yaormfilter_test

import (
	"testing"
	"time"

	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/stretchr/testify/assert"
)

func TestNewDateFilter(t *testing.T) {
	filter := yaormfilter.NewDateFilter()
	assert.IsType(t, &yaormfilter.DateFilter{}, filter)
}

func TestDateFilter_Equals(t *testing.T) {
	filter := yaormfilter.NewDateFilter()
	now := time.Now()
	assert.Equal(t, filter, filter.Equals(now))
	assert.Equal(t, filter, filter.Equals(&now))
	assert.Panics(t, func() { filter.Equals(0) })
}

func TestDateFilter_NotEquals(t *testing.T) {
	filter := yaormfilter.NewDateFilter()
	now := time.Now()
	assert.Equal(t, filter, filter.NotEquals(now))
	assert.Equal(t, filter, filter.NotEquals(&now))
	assert.Panics(t, func() { filter.NotEquals(0) })
}

func TestDateFilter_Like(t *testing.T) {
	filter := yaormfilter.NewDateFilter()
	assert.Equal(t, filter, filter.Like(time.Now()))
}

func TestDateFilter_Nil(t *testing.T) {
	filter := yaormfilter.NewDateFilter()
	assert.Equal(t, filter, filter.Nil(true))
}

func TestDateFilter_In(t *testing.T) {
	filter := yaormfilter.NewDateFilter()
	assert.Equal(t, filter, filter.In(time.Now()))
}
