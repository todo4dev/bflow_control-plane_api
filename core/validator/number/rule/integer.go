package rule

import (
	"math"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/number/util"
)

type IntegerRule struct{}

var _ _base.IRule = IntegerRule{}

func (r IntegerRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	value, present, typeErr := util.ExtractNumber(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	if math.Abs(value-math.Trunc(value)) > util.NumberEpsilon {
		return []_base.ValidationError{{
			Code:    "number.integer",
			Message: "number must be an integer",
		}}
	}

	return nil
}
