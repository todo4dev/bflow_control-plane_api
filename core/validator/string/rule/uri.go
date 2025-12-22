package rule

import (
	"net/url"
	"reflect"

	"src/core/validator/_base"
)

type URIRule struct{}

var _ _base.IRule = URIRule{}

func (r URIRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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
			Code:    "string.uri",
			Message: "string must be a string",
		}}
	}

	value := target.String()
	if value == "" {
		return nil
	}

	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return []_base.ValidationError{{
			Code:    "string.uri",
			Message: "string must be a valid URI",
		}}
	}

	return nil
}
