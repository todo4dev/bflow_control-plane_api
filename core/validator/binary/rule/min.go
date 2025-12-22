package rule

import (
	"fmt"
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

	if length < r.Min {
		return []_base.ValidationError{{
			Code:    "binary.min",
			Message: fmt.Sprintf("binary must have at least %d bytes", r.Min),
		}}
	}

	return nil
}
