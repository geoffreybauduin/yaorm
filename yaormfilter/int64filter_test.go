package yaormfilter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

func TestNewInt64Filter(t *testing.T) {
	filter := yaormfilter.NewInt64Filter()
	assert.IsType(t, &yaormfilter.Int64Filter{}, filter)
}

func TestInt64Filter_Equals(t *testing.T) {
	filter := yaormfilter.NewInt64Filter()
	v := int64(12)
	assert.Equal(t, filter, filter.Equals(v))
	assert.Equal(t, filter, filter.Equals(&v))
	assert.Panics(t, func() { filter.Equals("aazeae") })
}

func TestInt64Filter_Like(t *testing.T) {
	filter := yaormfilter.NewInt64Filter()
	assert.Equal(t, filter, filter.Like(int64(12)))
}

func TestInt64Filter_Nil(t *testing.T) {
	filter := yaormfilter.NewInt64Filter()
	assert.Equal(t, filter, filter.Nil(true))
}

func TestInt64Filter_In(t *testing.T) {
	filter := yaormfilter.NewInt64Filter()
	assert.Equal(t, filter, filter.In(int64(12), int64(13)))
}
