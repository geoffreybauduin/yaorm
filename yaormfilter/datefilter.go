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

// Equals applies an equal filter on Date
func (f *DateFilter) Equals(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	// make sure we have an int64
	if underlyingValue.Kind() != reflect.Struct || underlyingValue.Type() != timeType {
		panic("Value in DateFilter is not a time.Time object")
	}
	f.equals(underlyingValue.Interface())
	return f
}

// Like is not applicable on Date
func (f *DateFilter) Like(v interface{}) ValueFilter {
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
