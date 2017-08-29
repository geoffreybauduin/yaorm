package yaormfilter

import (
	"reflect"

	"github.com/geoffreybauduin/yaorm/tools"
)

// Int64Filter is the filter used to filter on int64 fields. Implements ValueFilter
type Int64Filter struct {
	valuefilterimpl
}

// NewInt64Filter returns a new int64 filter
func NewInt64Filter() ValueFilter {
	return &Int64Filter{}
}

// Equals adds an equal filter
func (f *Int64Filter) Equals(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	// make sure we have an int64
	if underlyingValue.Kind() != reflect.Int64 {
		panic("Value in Int64Filter is not an int64")
	}
	f.equals(underlyingValue.Interface())
	return f
}

// Like is not applicable on int64
func (f *Int64Filter) Like(v interface{}) ValueFilter {
	return f
}

// Nil adds a nil filter
func (f *Int64Filter) Nil(v bool) ValueFilter {
	f.nil(v)
	return f
}

// In adds a IN filter
func (f *Int64Filter) In(values ...interface{}) ValueFilter {
	interfaceValues := []interface{}{}
	for _, v := range values {
		underlyingValue := tools.GetNonPtrValue(v)
		// make sure we have an int64
		if underlyingValue.Kind() != reflect.Int64 {
			panic("Value in Int64Filter is not an int64")
		}
		interfaceValues = append(interfaceValues, underlyingValue.Interface())
	}
	f.in(interfaceValues)
	return f
}
