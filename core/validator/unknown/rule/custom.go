package rule

import (
	"reflect"

	"src/core/validator/_base"
)

type CustomRule struct {
	ValidatorFunc     _base.ValidatorFunc
	Param             any
	Code              string
	Message           string
	optionalFieldName string
}

var _ _base.IRule = CustomRule{}

func (r CustomRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	if r.ValidatorFunc == nil {
		return nil
	}

	if !valuePointer.IsValid() {
		return nil
	}

	target := valuePointer.Elem()

	for target.Kind() == reflect.Pointer {
		if target.IsNil() {
			if err := r.ValidatorFunc(nil, r.Param); err != nil {
				return []_base.ValidationError{{
					Code:    r.ErrorCode(),
					Message: r.ErrorMessage(err),
				}}
			}
			return nil
		}
		target = target.Elem()
	}

	if !target.CanInterface() {
		return nil
	}

	value := target.Interface()

	if err := r.ValidatorFunc(value, r.Param); err != nil {
		return []_base.ValidationError{{
			Code:    r.ErrorCode(),
			Message: r.ErrorMessage(err),
		}}
	}

	return nil
}

func (r CustomRule) ErrorCode() string {
	if r.Code != "" {
		return r.Code
	}
	return "common.custom"
}

func (r CustomRule) ErrorMessage(err error) string {
	if r.Message != "" {
		return r.Message
	}
	if err == nil {
		return ""
	}
	return err.Error()
}
