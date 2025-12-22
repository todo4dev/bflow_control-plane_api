package rule

import (
	"fmt"
	"reflect"
	"time"

	"src/core/validator/_base"
)

type GreaterRule struct {
	Limit time.Time
}

var _ _base.IRule = GreaterRule{}

func (r GreaterRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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

	if target.Type() != reflect.TypeOf(time.Time{}) {
		return []_base.ValidationError{{
			Code:    "date.type",
			Message: "date must be a time.Time",
		}}
	}

	value := target.Interface().(time.Time)
	if value.IsZero() {
		return nil
	}

	if !value.After(r.Limit) {
		return []_base.ValidationError{{
			Code:    "date.greater",
			Message: fmt.Sprintf("date must be greater than %s", r.Limit.Format(time.RFC3339)),
		}}
	}

	return nil
}
