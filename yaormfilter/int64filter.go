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

func (f *Int64Filter) getValue(v interface{}) interface{} {
	underlyingValue := tools.GetNonPtrValue(v)
	// make sure we have an int64
	if underlyingValue.Kind() != reflect.Int64 {
		panic("Value in Int64Filter is not an int64")
	}
	return underlyingValue.Interface()
}

// Equals adds an equal filter
func (f *Int64Filter) Equals(v interface{}) ValueFilter {
	f.equals(f.getValue(v))
	return f
}

// NotEquals adds an notEqual filter
func (f *Int64Filter) NotEquals(v interface{}) ValueFilter {
	f.notEquals(f.getValue(v))
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
		interfaceValues = append(interfaceValues, f.getValue(v))
	}
	f.in(interfaceValues)
	return f
}

// NotIn adds a NOT IN filter
func (f *Int64Filter) NotIn(values ...interface{}) ValueFilter {
	interfaceValues := []interface{}{}
	for _, v := range values {
		interfaceValues = append(interfaceValues, f.getValue(v))
	}
	f.notIn(interfaceValues)
	return f
}

// Lt adds a < filter
func (f *Int64Filter) Lt(v interface{}) ValueFilter {
	f.lt(f.getValue(v))
	return f
}

// Lte adds a <= filter
func (f *Int64Filter) Lte(v interface{}) ValueFilter {
	f.lte(f.getValue(v))
	return f
}

// Gt adds a > filter
func (f *Int64Filter) Gt(v interface{}) ValueFilter {
	f.gt(f.getValue(v))
	return f
}

// Gte adds a > filter
func (f *Int64Filter) Gte(v interface{}) ValueFilter {
	f.gte(f.getValue(v))
	return f
}

// Raw performs a Raw filter
func (f *Int64Filter) Raw(s RawFilterFunc) ValueFilter {
	f.raw(s)
	return f
}
