package rule

import (
	"reflect"
	"regexp"

	"src/core/validator/_base"
)

type RegexRule struct {
	Code    string
	Message string
	Pattern *regexp.Regexp
}

var _ _base.IRule = RegexRule{}

func (r RegexRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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
			Code:    r.Code,
			Message: "string must be a string",
		}}
	}

	value := target.String()
	if value == "" {
		return nil
	}

	if !r.Pattern.MatchString(value) {
		message := r.Message
		if message == "" {
			message = "string does not match required pattern"
		}
		return []_base.ValidationError{{
			Code:    r.Code,
			Message: message,
		}}
	}

	return nil
}
