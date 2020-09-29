package yaormfilter_test

import (
	"testing"
	"time"

	"github.com/geoffreybauduin/yaorm/yaormfilter"
	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.Equals("abcdef"))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.Equals(int64(12)))
	assert.IsType(t, &yaormfilter.NilFilter{}, yaormfilter.Equals(nil))
	assert.IsType(t, &yaormfilter.DateFilter{}, yaormfilter.Equals(time.Now()))
	assert.IsType(t, &yaormfilter.BoolFilter{}, yaormfilter.Equals(false))
}

func TestEquality(t *testing.T) {
	f := yaormfilter.Equals("abdef")
	assert.True(t, f.IsEquality())
	assert.Equal(t, "abdef", f.GetEquality())

	f = yaormfilter.NotEquals("bdef")
	assert.False(t, f.IsEquality())
}

func TestNotEquals(t *testing.T) {
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.NotEquals("abcdef"))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.NotEquals(int64(12)))
	assert.IsType(t, &yaormfilter.NilFilter{}, yaormfilter.NotEquals(nil))
	assert.IsType(t, &yaormfilter.DateFilter{}, yaormfilter.NotEquals(time.Now()))
	assert.IsType(t, &yaormfilter.BoolFilter{}, yaormfilter.NotEquals(false))
}

func TestIn(t *testing.T) {
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.In("abcdef", "bcderzzer"))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.In(int64(12), int64(15)))
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.In([]string{"abcdef", "bcderzzer"}))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.In([]int64{int64(12), int64(15)}))
	assert.IsType(t, &yaormfilter.BoolFilter{}, yaormfilter.In([]bool{true, false}))
}

func TestNotIn(t *testing.T) {
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.NotIn("abcdef", "bcderzzer"))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.NotIn(int64(12), int64(15)))
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.NotIn([]string{"abcdef", "bcderzzer"}))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.NotIn([]int64{int64(12), int64(15)}))
	assert.IsType(t, &yaormfilter.BoolFilter{}, yaormfilter.NotIn([]bool{true, false}))
}

func TestLike(t *testing.T) {
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.Like("abcdef%"))
}

func TestLt(t *testing.T) {
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.Lt(int64(12)))
	assert.IsType(t, &yaormfilter.DateFilter{}, yaormfilter.Lt(time.Now()))
}

func TestLte(t *testing.T) {
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.Lte(int64(12)))
	assert.IsType(t, &yaormfilter.DateFilter{}, yaormfilter.Lte(time.Now()))
}

func TestGt(t *testing.T) {
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.Gt(int64(12)))
	assert.IsType(t, &yaormfilter.DateFilter{}, yaormfilter.Gt(time.Now()))
}

func TestGte(t *testing.T) {
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.Gte(int64(12)))
	assert.IsType(t, &yaormfilter.DateFilter{}, yaormfilter.Gte(time.Now()))
}
