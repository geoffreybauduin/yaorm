package yaormfilter

import (
	"reflect"

	"github.com/geoffreybauduin/yaorm/tools"
)

type StringFilter struct {
	valuefilterimpl
}

// NewStringFilter returns a new string filter
func NewStringFilter() ValueFilter {
	return &StringFilter{}
}

// Equals adds an equal filter
func (f *StringFilter) Equals(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	// make sure we have a string
	if underlyingValue.Kind() != reflect.String {
		panic("Value in StringFilter is not a string")
	}
	f.equals(underlyingValue.Interface())
	return f
}

// NotEquals adds an notEqual filter
func (f *StringFilter) NotEquals(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	// make sure we have a string
	if underlyingValue.Kind() != reflect.String {
		panic("Value in StringFilter is not a string")
	}
	f.notEquals(underlyingValue.Interface())
	return f
}

// Like adds a Like filter
func (f *StringFilter) Like(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	// make sure we have a string
	if underlyingValue.Kind() != reflect.String {
		panic("Value in StringFilter is not a string")
	}
	f.like(underlyingValue.Interface())
	return f
}

// Nil adds a nil filter
func (f *StringFilter) Nil(v bool) ValueFilter {
	f.nil(v)
	return f
}

// In adds a IN filter
func (f *StringFilter) In(values ...interface{}) ValueFilter {
	interfaceValues := []interface{}{}
	for _, v := range values {
		underlyingValue := tools.GetNonPtrValue(v)
		// make sure we have a string
		if underlyingValue.Kind() != reflect.String {
			panic("Value in StringFilter is not an string")
		}
		interfaceValues = append(interfaceValues, underlyingValue.Interface())
	}
	f.in(interfaceValues)
	return f
}

// NotIn adds a NOT IN filter
func (f *StringFilter) NotIn(values ...interface{}) ValueFilter {
	interfaceValues := []interface{}{}
	for _, v := range values {
		underlyingValue := tools.GetNonPtrValue(v)
		// make sure we have a string
		if underlyingValue.Kind() != reflect.String {
			panic("Value in StringFilter is not an string")
		}
		interfaceValues = append(interfaceValues, underlyingValue.Interface())
	}
	f.notIn(interfaceValues)
	return f
}

// Lt adds a < filter
func (f *StringFilter) Lt(v interface{}) ValueFilter {
	return f
}

// Lte adds a <= filter
func (f *StringFilter) Lte(v interface{}) ValueFilter {
	return f
}

// Gt adds a > filter
func (f *StringFilter) Gt(v interface{}) ValueFilter {
	return f
}

// Gte adds a > filter
func (f *StringFilter) Gte(v interface{}) ValueFilter {
	return f
}

// Raw performs a Raw filter
func (f *StringFilter) Raw(s RawFilterFunc) ValueFilter {
	f.raw_ = s
	return f
}
