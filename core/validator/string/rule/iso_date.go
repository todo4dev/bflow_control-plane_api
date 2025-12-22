package rule

import (
	"reflect"
	"time"

	"src/core/validator/_base"
)

type ISODateRule struct{}

var _ _base.IRule = ISODateRule{}

func (r ISODateRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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
			Code:    "string.isoDate",
			Message: "string must be a string",
		}}
	}

	value := target.String()
	if value == "" {
		return nil
	}

	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return []_base.ValidationError{{
			Code:    "string.isoDate",
			Message: "string must be a valid RFC3339 timestamp",
		}}
	}

	return nil
}
