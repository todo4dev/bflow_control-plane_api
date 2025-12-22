package rule

import (
	"fmt"
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

	if target.Kind() != reflect.Slice || target.Type().Elem().Kind() != reflect.Uint8 {
		return []_base.ValidationError{{
			Code:    "binary.type",
			Message: "binary must be a byte slice ([]byte)",
		}}
	}

	length := target.Len()

	if length == 0 {
		return nil
	}

	if length != r.Length {
		return []_base.ValidationError{{
			Code:    "binary.length",
			Message: fmt.Sprintf("binary must have exact length of %d bytes", r.Length),
		}}
	}

	return nil
}
