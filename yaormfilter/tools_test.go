package yaormfilter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/geoffreybauduin/yaorm/yaormfilter"
)

func TestEquals(t *testing.T) {
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.Equals("abcdef"))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.Equals(int64(12)))
	assert.IsType(t, &yaormfilter.NilFilter{}, yaormfilter.Equals(nil))
	assert.IsType(t, &yaormfilter.DateFilter{}, yaormfilter.Equals(time.Now()))
}

func TestIn(t *testing.T) {
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.In("abcdef", "bcderzzer"))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.In(int64(12), int64(15)))
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.In([]string{"abcdef", "bcderzzer"}))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.In([]int64{int64(12), int64(15)}))
}

func TestLike(t *testing.T) {
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.Like("abcdef%"))
}
