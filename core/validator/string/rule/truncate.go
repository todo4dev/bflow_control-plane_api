package rule

import (
	"reflect"

	"src/core/validator/_base"
)

type TruncateRule struct {
	Max int
}

var _ _base.IRule = TruncateRule{}

func (r TruncateRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	if !valuePointer.IsValid() || r.Max <= 0 {
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
			Code:    "string.truncate",
			Message: "string must be a string",
		}}
	}

	value := target.String()
	if value == "" {
		return nil
	}

	runes := []rune(value)
	if len(runes) <= r.Max {
		return nil
	}

	truncated := string(runes[:r.Max])
	if !target.CanSet() {
		return nil
	}

	target.SetString(truncated)
	return nil
}
