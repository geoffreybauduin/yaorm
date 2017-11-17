package tools

import "reflect"

func GetNonPtrValue(v interface{}) reflect.Value {
	underlyingValue := reflect.ValueOf(v)
	for underlyingValue.Kind() == reflect.Ptr {
		underlyingValue = reflect.Indirect(underlyingValue)
	}
	return underlyingValue
}

func GetNonPtrInterface(v interface{}) interface{} {
	underlyingValue := GetNonPtrValue(v)
	return underlyingValue.Interface()
}

func IsZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice, reflect.Ptr:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && IsZeroValue(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && IsZeroValue(v.Field(i))
		}
		return z
	case reflect.Invalid:
		return !v.IsValid()
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	if v.CanInterface() && z.CanInterface() {
		return v.Interface() == z.Interface()
	} else {
		// Assuming this a zero value since unexported fields should not apply any filter
		return true
	}
}

// IsNil returns true if given value is nil
func IsNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice, reflect.Ptr:
		return v.IsNil()
	}
	return false
}