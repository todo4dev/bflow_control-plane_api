package rule

import (
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/number/util"
)

type NegativeRule struct{}

var _ _base.IRule = NegativeRule{}

func (r NegativeRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	value, present, typeErr := util.ExtractNumber(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	if !(value < 0) {
		return []_base.ValidationError{{
			Code:    "number.negative",
			Message: "number must be a negative number",
		}}
	}

	return nil
}
