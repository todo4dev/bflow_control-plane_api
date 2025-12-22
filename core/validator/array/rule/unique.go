package rule

import (
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/array/util"
)

type UniqueRule struct{}

var _ _base.IRule = UniqueRule{}

func (r UniqueRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	slice, present, typeErr := util.ExtractSlice(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	length := slice.Len()
	if length <= 1 {
		return nil
	}

	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			a := slice.Index(i)
			b := slice.Index(j)

			if reflect.DeepEqual(a.Interface(), b.Interface()) {
				return []_base.ValidationError{{
					Code:    "array.unique",
					Message: "array must not contain duplicate values",
				}}
			}
		}
	}

	return nil
}
