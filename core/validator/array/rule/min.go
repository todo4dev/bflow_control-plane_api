package rule

import (
	"fmt"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/array/util"
)

type MinRule struct {
	Min int
}

var _ _base.IRule = MinRule{}

func (r MinRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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

	if length < r.Min {
		return []_base.ValidationError{{
			Code:    "array.min",
			Message: fmt.Sprintf("array must have at least %d items", r.Min),
		}}
	}

	return nil
}
