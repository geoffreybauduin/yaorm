package yaormfilter

import (
	"fmt"
	"reflect"
	"time"

	"github.com/geoffreybauduin/yaorm/tools"
)

// Equals returns the correct filter according the value sent
func Equals(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	switch underlyingValue.Kind() {
	case reflect.Int64:
		return NewInt64Filter().Equals(v)
	case reflect.String:
		return NewStringFilter().Equals(v)
	case reflect.Struct:
		if _, ok := underlyingValue.Interface().(time.Time); ok {
			return NewDateFilter().Equals(v)
		}
	}
	if v == nil {
		return NewNilFilter().Nil(true)
	}
	panic(fmt.Errorf("Unknown type: %+v for value %+v in Equals filter", underlyingValue.Kind(), v))
}

// In returns the correct filter according to the value sent
func In(values ...interface{}) ValueFilter {
	var t reflect.Type
	if tools.GetNonPtrValue(values).Len() == 0 {
		return nil
	}
	for idx, v := range values {
		underlyingValue := tools.GetNonPtrValue(v)
		if idx == 0 {
			t = underlyingValue.Type()
		} else {
			if underlyingValue.Type() != t {
				panic(fmt.Errorf("Inconsistent values sent, got types: %+v and %+v", t, underlyingValue.Type()))
			}
		}
	}
	switch t.Kind() {
	case reflect.Int64:
		return NewInt64Filter().In(values...)
	case reflect.String:
		return NewStringFilter().In(values...)
	case reflect.Slice:
		// if we receive a slice, we want to go through all the slices received an concat them inside one
		data := []interface{}{}
		for _, v := range values {
			underlyingValue := tools.GetNonPtrValue(v)
			for i := 0; i < underlyingValue.Len(); i++ {
				cell := tools.GetNonPtrValue(underlyingValue.Index(i).Interface())
				data = append(data, cell.Interface())
			}
		}
		return In(data...)
	}
	panic(fmt.Errorf("Unknown type: %v inside In filter", t.Kind()))
}

// Like returns the correct filter according to the value sent
func Like(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	switch underlyingValue.Kind() {
	case reflect.String:
		return NewStringFilter().Like(v)
	}
	panic(fmt.Errorf("Unknown type: %+v for value %+v in Like filter", underlyingValue.Kind(), v))
}

// Lt returns the correct filter according the value sent
func Lt(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	switch underlyingValue.Kind() {
	case reflect.Int64:
		return NewInt64Filter().Lt(v)
	case reflect.Struct:
		if _, ok := underlyingValue.Interface().(time.Time); ok {
			return NewDateFilter().Lt(v)
		}
	}
	panic(fmt.Errorf("Unknown type: %+v for value %+v in Lt filter", underlyingValue.Kind(), v))
}

// Lte returns the correct filter according the value sent
func Lte(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	switch underlyingValue.Kind() {
	case reflect.Int64:
		return NewInt64Filter().Lte(v)
	case reflect.Struct:
		if _, ok := underlyingValue.Interface().(time.Time); ok {
			return NewDateFilter().Lte(v)
		}
	}
	panic(fmt.Errorf("Unknown type: %+v for value %+v in Lte filter", underlyingValue.Kind(), v))
}

// Gt returns the correct filter according the value sent
func Gt(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	switch underlyingValue.Kind() {
	case reflect.Int64:
		return NewInt64Filter().Gt(v)
	case reflect.Struct:
		if _, ok := underlyingValue.Interface().(time.Time); ok {
			return NewDateFilter().Gt(v)
		}
	}
	panic(fmt.Errorf("Unknown type: %+v for value %+v in Gt filter", underlyingValue.Kind(), v))
}

// Gte returns the correct filter according the value sent
func Gte(v interface{}) ValueFilter {
	underlyingValue := tools.GetNonPtrValue(v)
	switch underlyingValue.Kind() {
	case reflect.Int64:
		return NewInt64Filter().Gte(v)
	case reflect.Struct:
		if _, ok := underlyingValue.Interface().(time.Time); ok {
			return NewDateFilter().Gte(v)
		}
	}
	panic(fmt.Errorf("Unknown type: %+v for value %+v in Gte filter", underlyingValue.Kind(), v))
}
