package util

import (
	"reflect"

	"src/core/validator/_base"
)

const NumberEpsilon = 1e-9

func ExtractNumber(valuePointer reflect.Value) (float64, bool, []_base.ValidationError) {
	if !valuePointer.IsValid() {
		return 0, false, nil
	}

	target := valuePointer.Elem()
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return 0, false, nil
		}
		target = target.Elem()
	}

	switch target.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(target.Int()), true, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(target.Uint()), true, nil
	case reflect.Float32, reflect.Float64:
		return target.Float(), true, nil
	default:
		return 0, false, []_base.ValidationError{{
			Code:    "number.type",
			Message: "value must be a number",
		}}
	}
}
