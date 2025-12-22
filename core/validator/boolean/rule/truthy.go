package rule

import (
	"reflect"
	"slices"

	"src/core/validator/_base"
)

type TruthyRule struct {
	TruthyValues []bool
}

var _ _base.IRule = TruthyRule{}

func (r TruthyRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	if valuePointer.Kind() == reflect.Bool {
		b := valuePointer.Bool()
		if slices.Contains(r.TruthyValues, b) {
			return nil
		}
		return []_base.ValidationError{{
			Code:    "boolean.truthy",
			Message: "value is not considered truthy by this rule",
		}}
	}
	return nil
}
