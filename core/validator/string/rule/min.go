package rule

import (
	"reflect"

	"src/core/validator/_base"
)

type MinRule struct {
	Min int
}

var _ _base.IRule = MinRule{}

func (r MinRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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
			Code:    "string.min",
			Message: "string must be a string",
		}}
	}

	value := target.String()
	if value == "" {
		return nil
	}

	if runeCount := len([]rune(value)); runeCount < r.Min {
		return []_base.ValidationError{{
			Code:    "string.min",
			Message: "string must have at least minimum length",
		}}
	}

	return nil
}
