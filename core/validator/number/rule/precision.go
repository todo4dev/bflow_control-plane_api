package rule

import (
	"fmt"
	"math"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/number/util"
)

type PrecisionRule struct {
	Precision int
}

var _ _base.IRule = PrecisionRule{}

func (r PrecisionRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	value, present, typeErr := util.ExtractNumber(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	if math.Abs(value-math.Trunc(value)) < util.NumberEpsilon {
		return nil
	}

	if r.Precision < 0 {
		return nil
	}

	scale := math.Pow10(r.Precision)
	rounded := math.Round(value*scale) / scale

	if math.Abs(value-rounded) > util.NumberEpsilon {
		return []_base.ValidationError{{
			Code:    "number.precision",
			Message: fmt.Sprintf("number must have at most %d decimal places", r.Precision),
		}}
	}

	return nil
}
