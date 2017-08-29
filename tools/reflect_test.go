package tools_test

import (
	"reflect"
	"testing"

	"github.com/geoffreybauduin/yaorm/tools"
	"github.com/stretchr/testify/assert"
)

func TestGetNonPtrValue(t *testing.T) {
	firstValue := "str"
	val := tools.GetNonPtrValue(&firstValue)
	assert.Equal(t, "str", val.Interface())
	secondValue := &firstValue
	val = tools.GetNonPtrValue(&secondValue)
	assert.Equal(t, "str", val.Interface())
}

func TestGetNonPtrInterface(t *testing.T) {
	firstValue := "str"
	val := tools.GetNonPtrInterface(&firstValue)
	assert.Equal(t, "str", val)
	secondValue := &firstValue
	val = tools.GetNonPtrInterface(&secondValue)
	assert.Equal(t, "str", val)
}

func TestIsZeroValue(t *testing.T) {
	assert.True(t, tools.IsZeroValue(reflect.ValueOf(nil)))
	assert.True(t, tools.IsZeroValue(reflect.ValueOf(0)))
	assert.True(t, tools.IsZeroValue(reflect.ValueOf(struct{ Test string }{""})))
	assert.False(t, tools.IsZeroValue(reflect.ValueOf([]int{})))
	assert.False(t, tools.IsZeroValue(reflect.ValueOf([]int{0})))
	assert.False(t, tools.IsZeroValue(reflect.ValueOf(1)))
	assert.False(t, tools.IsZeroValue(reflect.ValueOf([]int{1})))
	assert.False(t, tools.IsZeroValue(reflect.ValueOf(struct{ Test string }{"test"})))
}
