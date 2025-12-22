package rule

import (
	"fmt"
	"reflect"
	"time"

	"src/core/validator/_base"
)

type LessRule struct {
	Limit time.Time
}

var _ _base.IRule = LessRule{}

func (r LessRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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

	if !value.Before(r.Limit) {
		return []_base.ValidationError{{
			Code:    "date.less",
			Message: fmt.Sprintf("date must be less than %s", r.Limit.Format(time.RFC3339)),
		}}
	}

	return nil
}
