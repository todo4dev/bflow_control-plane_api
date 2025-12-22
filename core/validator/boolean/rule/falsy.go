package rule

import (
	"reflect"
	"slices"

	"src/core/validator/_base"
)

type FalsyRule struct {
	FalsyValues []bool
}

var _ _base.IRule = FalsyRule{}

func (r FalsyRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	if valuePointer.Kind() == reflect.Bool {
		b := valuePointer.Bool()
		if slices.Contains(r.FalsyValues, b) {
			return nil
		}
		return []_base.ValidationError{{
			Code:    "boolean.falsy",
			Message: "value is not considered falsy by this rule",
		}}
	}
	return nil
}
