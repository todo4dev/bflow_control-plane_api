package rule

import (
	"fmt"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/number/util"
)

type LessRule struct {
	Limit float64
}

var _ _base.IRule = LessRule{}

func (r LessRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	value, present, typeErr := util.ExtractNumber(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	if !(value < r.Limit) {
		return []_base.ValidationError{{
			Code:    "number.less",
			Message: fmt.Sprintf("number must be less than %v", r.Limit),
		}}
	}

	return nil
}
