package rule

import (
	"encoding/base64"
	"fmt"
	"reflect"

	"src/core/validator/_base"
)

type Base64Rule struct{}

var _ _base.IRule = Base64Rule{}

func (r Base64Rule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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
			Code:    "string.base64",
			Message: "string must be a string",
		}}
	}

	value := target.String()
	if value == "" {
		return nil
	}

	if _, err := base64.StdEncoding.DecodeString(value); err != nil {
		return []_base.ValidationError{{
			Code:    "string.base64",
			Message: fmt.Sprintf("string must be a valid base64 string"),
		}}
	}

	return nil
}
