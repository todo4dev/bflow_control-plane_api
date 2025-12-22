package rule

import (
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/unknown/util"
)

type RequiredRule struct{}

var _ _base.IRule = RequiredRule{}

func (r RequiredRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	targetValue := valuePointer.Elem()

	if !valuePointer.IsValid() || util.IsZeroValue(targetValue) {
		return []_base.ValidationError{{
			Code:    "common.required",
			Message: "value is required",
		}}
	}

	return nil
}
