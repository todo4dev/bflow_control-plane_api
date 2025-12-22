package rule

import (
	"fmt"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/array/util"
)

type MaxRule struct {
	Max int
}

var _ _base.IRule = MaxRule{}

func (r MaxRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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

	if length > r.Max {
		return []_base.ValidationError{{
			Code:    "array.max",
			Message: fmt.Sprintf("array must have at most %d items", r.Max),
		}}
	}

	return nil
}
