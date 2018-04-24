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

func TestIn(t *testing.T) {
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.In("abcdef", "bcderzzer"))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.In(int64(12), int64(15)))
	assert.IsType(t, &yaormfilter.StringFilter{}, yaormfilter.In([]string{"abcdef", "bcderzzer"}))
	assert.IsType(t, &yaormfilter.Int64Filter{}, yaormfilter.In([]int64{int64(12), int64(15)}))
	assert.IsType(t, &yaormfilter.BoolFilter{}, yaormfilter.In([]bool{true, false}))
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
