package rule

import (
	"reflect"
	"regexp"

	"src/core/validator/_base"
)

type ReplaceRule struct {
	Pattern     *regexp.Regexp
	Replacement string
}

var _ _base.IRule = ReplaceRule{}

func (r ReplaceRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	if !valuePointer.IsValid() || r.Pattern == nil {
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
			Code:    "string.replace",
			Message: "string must be a string",
		}}
	}

	original := target.String()
	if original == "" {
		return nil
	}

	replaced := r.Pattern.ReplaceAllString(original, r.Replacement)
	if !target.CanSet() {
		return nil
	}

	target.SetString(replaced)
	return nil
}
