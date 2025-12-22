package rule

import (
	"reflect"

	"src/core/validator/_base"
)

type TransformRule struct {
	Code      string
	Message   string
	Transform func(string) string
}

var _ _base.IRule = TransformRule{}

func (r TransformRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	if !valuePointer.IsValid() || r.Transform == nil {
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
			Code:    r.Code,
			Message: "string must be a string",
		}}
	}

	original := target.String()
	if original == "" {
		return nil
	}

	transformed := r.Transform(original)
	if !target.CanSet() {
		return nil
	}

	target.SetString(transformed)
	return nil
}
