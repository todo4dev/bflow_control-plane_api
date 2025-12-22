package rule

import (
	"math"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/number/util"
)

type PortRule struct{}

var _ _base.IRule = PortRule{}

func (r PortRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	value, present, typeErr := util.ExtractNumber(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	if math.Abs(value-math.Trunc(value)) > util.NumberEpsilon {
		return []_base.ValidationError{{
			Code:    "number.port",
			Message: "number must be an integer within valid port range (0-65535)",
		}}
	}

	if value < 0 || value > 65535 {
		return []_base.ValidationError{{
			Code:    "number.port",
			Message: "number must be within valid port range (0-65535)",
		}}
	}

	return nil
}
