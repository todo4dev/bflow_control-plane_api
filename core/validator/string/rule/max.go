package rule

import (
	"reflect"

	"src/core/validator/_base"
)

type MaxRule struct {
	Max int
}

var _ _base.IRule = MaxRule{}

func (r MaxRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	if !valuePointer.IsValid() {
		return nil
	}

	target := valuePointer.Elem()
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return nil
		}
		target = target.Elem()
	}

	if target.Kind() != reflect.String {
		return []_base.ValidationError{{
			Code:    "string.max",
			Message: "string must be a string",
		}}
	}

	value := target.String()
	if value == "" {
		return nil
	}

	if runeCount := len([]rune(value)); runeCount > r.Max {
		return []_base.ValidationError{{
			Code:    "string.max",
			Message: "string must have at most maximum length",
		}}
	}

	return nil
}
