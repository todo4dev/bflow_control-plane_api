package rule

import (
	"fmt"
	"math"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/number/util"
)

type MultipleRule struct {
	Base float64
}

var _ _base.IRule = MultipleRule{}

func (r MultipleRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	value, present, typeErr := util.ExtractNumber(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	if r.Base == 0 {
		return nil
	}

	remainder := math.Mod(value, r.Base)
	if math.Abs(remainder) > util.NumberEpsilon && math.Abs(r.Base) > util.NumberEpsilon {
		return []_base.ValidationError{{
			Code:    "number.multiple",
			Message: fmt.Sprintf("number must be a multiple of %v", r.Base),
		}}
	}

	return nil
}
