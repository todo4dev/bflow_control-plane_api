package rule

import (
	"fmt"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/number/util"
)

type MaxRule struct {
	Limit float64
}

var _ _base.IRule = MaxRule{}

func (r MaxRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	value, present, typeErr := util.ExtractNumber(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	if value > r.Limit {
		return []_base.ValidationError{{
			Code:    "number.max",
			Message: fmt.Sprintf("number must be less than or equal to %v", r.Limit),
		}}
	}

	return nil
}
