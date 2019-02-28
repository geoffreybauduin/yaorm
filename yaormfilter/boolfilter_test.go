package yaormfilter_test

import (
	"testing"

	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/stretchr/testify/assert"
)

func TestNewBoolFilter(t *testing.T) {
	filter := yaormfilter.NewBoolFilter()
	assert.IsType(t, &yaormfilter.BoolFilter{}, filter)
}

func TestBoolFilter_Equals(t *testing.T) {
	filter := yaormfilter.NewBoolFilter()
	v := true
	assert.Equal(t, filter, filter.Equals(v))
	assert.Equal(t, filter, filter.Equals(&v))
	assert.Panics(t, func() { filter.Equals("true") })
}

func TestBoolFilter_NotEquals(t *testing.T) {
	filter := yaormfilter.NewBoolFilter()
	v := true
	assert.Equal(t, filter, filter.NotEquals(v))
	assert.Equal(t, filter, filter.NotEquals(&v))
	assert.Panics(t, func() { filter.NotEquals("true") })
}

func TestBoolFilter_Like(t *testing.T) {
	filter := yaormfilter.NewBoolFilter()
	assert.Equal(t, filter, filter.Like(true))
}

func TestBoolFilter_Nil(t *testing.T) {
	filter := yaormfilter.NewBoolFilter()
	assert.Equal(t, filter, filter.Nil(true))
}

func TestBoolFilter_In(t *testing.T) {
	filter := yaormfilter.NewBoolFilter()
	assert.Equal(t, filter, filter.In(true, false))
}
