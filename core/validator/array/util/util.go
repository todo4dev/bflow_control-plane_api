package util

import (
	"reflect"

	"src/core/validator/_base"
)

func ExtractSlice(valuePointer reflect.Value) (reflect.Value, bool, []_base.ValidationError) {
	if !valuePointer.IsValid() {
		return reflect.Value{}, false, nil
	}

	target := valuePointer.Elem()
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return reflect.Value{}, false, nil
		}
		target = target.Elem()
	}

	switch target.Kind() {
	case reflect.Slice, reflect.Array:
		return target, true, nil
	default:
		return reflect.Value{}, false, []_base.ValidationError{{
			Code:    "array.type",
			Message: "value must be an array or slice",
		}}
	}
}
