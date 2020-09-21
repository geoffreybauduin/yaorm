package yaormfilter_test

import (
	"testing"

	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/stretchr/testify/assert"
)

func TestNewStringFilter(t *testing.T) {
	filter := yaormfilter.NewStringFilter()
	assert.IsType(t, &yaormfilter.StringFilter{}, filter)
}

func TestStringFilter_Equals(t *testing.T) {
	filter := yaormfilter.NewStringFilter()
	str := "bla"
	assert.Equal(t, filter, filter.Equals(str))
	assert.Equal(t, filter, filter.Equals(&str))
	assert.Panics(t, func() { filter.Equals(0) })
}

func TestStringFilter_NotEquals(t *testing.T) {
	filter := yaormfilter.NewStringFilter()
	str := "bla"
	assert.Equal(t, filter, filter.NotEquals(str))
	assert.Equal(t, filter, filter.NotEquals(&str))
	assert.Panics(t, func() { filter.NotEquals(0) })
}

func TestStringFilter_Like(t *testing.T) {
	filter := yaormfilter.NewStringFilter()
	str := "bla"
	assert.Equal(t, filter, filter.Like(str))
}

func TestStringFilter_Nil(t *testing.T) {
	filter := yaormfilter.NewStringFilter()
	assert.Equal(t, filter, filter.Nil(true))
}

func TestStringFilter_In(t *testing.T) {
	filter := yaormfilter.NewStringFilter()
	str := "bla"
	assert.Equal(t, filter, filter.In(str))
}

func TestStringFilter_Raw(t *testing.T) {
	filter := yaormfilter.NewStringFilter()
	assert.Equal(t, filter, filter.Raw(func(s string) string { return s }))
}
