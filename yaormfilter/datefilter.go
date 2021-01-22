package yaormfilter

import (
	"reflect"
	"time"

	"github.com/geoffreybauduin/yaorm/tools"
)

// DateFilter is a filter to operate on date fields
type DateFilter struct {
	valuefilterimpl
}

var (
	timeType = reflect.TypeOf(time.Time{})
)

// NewDateFilter returns a new DateFilter
func NewDateFilter() ValueFilter {
	return &DateFilter{}
}

func (f *DateFilter) getValue(v interface{}) interface{} {
	underlyingValue := tools.GetNonPtrValue(v)
	// make sure we have a time.Time
	if underlyingValue.Kind() != reflect.Struct || underlyingValue.Type() != timeType {
		panic("Value in DateFilter is not a time.Time object")
	}
	return underlyingValue.Interface()
}

// Equals applies an equal filter on Date
func (f *DateFilter) Equals(v interface{}) ValueFilter {
	f.equals(f.getValue(v))
	return f
}

// NotEquals applies an notEqual filter on Date
func (f *DateFilter) NotEquals(v interface{}) ValueFilter {
	f.notEquals(f.getValue(v))
	return f
}

// Like is not applicable on Date
func (f *DateFilter) Like(v interface{}) ValueFilter {
	return f
}

// ILike is not applicable on Date
func (f *DateFilter) ILike(v interface{}) ValueFilter {
	return f
}

// Nil applies a nil filter on Date
func (f *DateFilter) Nil(v bool) ValueFilter {
	f.nil(v)
	return f

}

// In adds a IN filter (not implemented)
func (f *DateFilter) In(values ...interface{}) ValueFilter {
	return f
}

// NotIn adds a NOT IN filter (not implemented)
func (f *DateFilter) NotIn(values ...interface{}) ValueFilter {
	return f
}

// Lt adds a < filter
func (f *DateFilter) Lt(v interface{}) ValueFilter {
	f.lt(f.getValue(v))
	return f
}

// Lte adds a <= filter
func (f *DateFilter) Lte(v interface{}) ValueFilter {
	f.lte(f.getValue(v))
	return f
}

// Gt adds a > filter
func (f *DateFilter) Gt(v interface{}) ValueFilter {
	f.gt(f.getValue(v))
	return f
}

// Gte adds a > filter
func (f *DateFilter) Gte(v interface{}) ValueFilter {
	f.gte(f.getValue(v))
	return f
}

// Raw performs a Raw filter
func (f *DateFilter) Raw(s RawFilterFunc) ValueFilter {
	f.raw(s)
	return f
}
