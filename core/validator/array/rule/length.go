package rule

import (
	"fmt"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/array/util"
)

type LengthRule struct {
	Length int
}

var _ _base.IRule = LengthRule{}

func (r LengthRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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

	if length != r.Length {
		return []_base.ValidationError{{
			Code:    "array.length",
			Message: fmt.Sprintf("array must have exact length of %d", r.Length),
		}}
	}

	return nil
}
