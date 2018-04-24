package yaormfilter

import (
	"reflect"

	"github.com/geoffreybauduin/yaorm/tools"
)

// BoolFilter is the filter used to filter on bool fields. Implements ValueFilter
type BoolFilter struct {
	valuefilterimpl
}

// NewBoolFilter returns a new bool filter
func NewBoolFilter() ValueFilter {
	return &BoolFilter{}
}

func (f *BoolFilter) getValue(v interface{}) interface{} {
	underlyingValue := tools.GetNonPtrValue(v)
	// make sure we have a bool
	if underlyingValue.Kind() != reflect.Bool {
		panic("Value in BoolFilter is not a bool")
	}
	return underlyingValue.Interface()
}

// Equals adds an equal filter
func (f *BoolFilter) Equals(v interface{}) ValueFilter {
	f.equals(f.getValue(v))
	return f
}

// Like is not applicable on bool
func (f *BoolFilter) Like(v interface{}) ValueFilter {
	return f
}

// Nil adds a nil filter
func (f *BoolFilter) Nil(v bool) ValueFilter {
	f.nil(v)
	return f
}

// In adds a IN filter
func (f *BoolFilter) In(values ...interface{}) ValueFilter {
	interfaceValues := []interface{}{}
	for _, v := range values {
		interfaceValues = append(interfaceValues, f.getValue(v))
	}
	f.in(interfaceValues)
	return f
}

// Lt is not applicable on bool
func (f *BoolFilter) Lt(v interface{}) ValueFilter {
	return f
}

// Lte is not applicable on bool
func (f *BoolFilter) Lte(v interface{}) ValueFilter {
	return f
}

// Gt is not applicable on bool
func (f *BoolFilter) Gt(v interface{}) ValueFilter {
	return f
}

// Gte is not applicable on bool
func (f *BoolFilter) Gte(v interface{}) ValueFilter {
	return f
}
