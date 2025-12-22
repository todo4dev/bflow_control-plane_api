package rule

import (
	"reflect"

	"src/core/validator/_base"
)

type LengthRule struct {
	Length int
}

var _ _base.IRule = LengthRule{}

func (r LengthRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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
			Code:    "string.length",
			Message: "string must be a string",
		}}
	}

	value := target.String()
	if value == "" {
		return nil
	}

	if runeCount := len([]rune(value)); runeCount != r.Length {
		return []_base.ValidationError{{
			Code:    "string.length",
			Message: "string must have exact length",
		}}
	}

	return nil
}
