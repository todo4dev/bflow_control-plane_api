package rule

import (
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/number/util"
)

type PositiveRule struct{}

var _ _base.IRule = PositiveRule{}

func (r PositiveRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	value, present, typeErr := util.ExtractNumber(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	if !(value > 0) {
		return []_base.ValidationError{{
			Code:    "number.positive",
			Message: "number must be a positive number",
		}}
	}

	return nil
}
