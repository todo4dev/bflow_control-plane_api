package rule

import (
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/unknown/util"
)

type AllowRule struct {
	Allowed []any
}

var _ _base.IRule = AllowRule{}

func (r AllowRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	if !valuePointer.IsValid() {
		return nil
	}

	targetValue := valuePointer.Elem()

	if util.IsZeroValue(targetValue) {
		return nil
	}

	canonical := util.DerefAll(targetValue)

	for _, allowed := range r.Allowed {
		if allowed == nil {
			continue
		}

		allowedValue := reflect.ValueOf(allowed)
		allowedCanonical := util.DerefAll(allowedValue)

		if !allowedCanonical.IsValid() {
			continue
		}

		if canonical.Type() == allowedCanonical.Type() {
			if reflect.DeepEqual(canonical.Interface(), allowedCanonical.Interface()) {
				return nil
			}
			continue
		}

		if allowedCanonical.Type().ConvertibleTo(canonical.Type()) {
			converted := allowedCanonical.Convert(canonical.Type())
			if reflect.DeepEqual(canonical.Interface(), converted.Interface()) {
				return nil
			}
		}
	}

	return []_base.ValidationError{{
		Code:    "common.allow",
		Message: "value is not in the list of allowed values",
	}}
}
