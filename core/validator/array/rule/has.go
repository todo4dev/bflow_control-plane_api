package rule

import (
	"fmt"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/array/util"
)

type HasRule struct {
	Expected any
}

var _ _base.IRule = HasRule{}

func (r HasRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	slice, present, typeErr := util.ExtractSlice(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	length := slice.Len()
	if length == 0 {
		return nil
	}

	expectedValue := reflect.ValueOf(r.Expected)

	for index := 0; index < length; index++ {
		element := slice.Index(index)

		if reflect.DeepEqual(element.Interface(), r.Expected) {
			return nil
		}

		if expectedValue.IsValid() &&
			expectedValue.Type().ConvertibleTo(element.Type()) {
			converted := expectedValue.Convert(element.Type())
			if reflect.DeepEqual(element.Interface(), converted.Interface()) {
				return nil
			}
		}
	}

	return []_base.ValidationError{{
		Code:    "array.has",
		Message: fmt.Sprintf("array must contain at least one matching value"),
	}}
}
