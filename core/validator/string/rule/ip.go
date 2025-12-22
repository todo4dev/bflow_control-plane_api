package rule

import (
	"net"
	"reflect"

	"src/core/validator/_base"
)

type IPRule struct{}

var _ _base.IRule = IPRule{}

func (r IPRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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
			Code:    "string.ip",
			Message: "string must be a string",
		}}
	}

	value := target.String()
	if value == "" {
		return nil
	}

	if net.ParseIP(value) == nil {
		return []_base.ValidationError{{
			Code:    "string.ip",
			Message: "string must be a valid IP address",
		}}
	}

	return nil
}
