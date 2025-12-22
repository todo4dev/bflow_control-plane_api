package util

import "reflect"

func IsZeroValue(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	switch value.Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			return true
		}
		return IsZeroValue(value.Elem())
	case reflect.Slice, reflect.Map, reflect.Array, reflect.String:
		return value.Len() == 0
	default:
		return value.IsZero()
	}
}

func DerefAll(value reflect.Value) reflect.Value {
	for value.IsValid() && value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return value
		}
		value = value.Elem()
	}
	return value
}
