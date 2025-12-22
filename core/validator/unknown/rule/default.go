package rule

import (
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/unknown/util"
)

type DefaultRule struct {
	DefaultValue any
}

var _ _base.IRule = DefaultRule{}

func (r DefaultRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	if !valuePointer.IsValid() {
		return nil
	}

	targetValue := valuePointer.Elem()

	if !targetValue.CanSet() {
		return nil
	}

	if !util.IsZeroValue(targetValue) {
		return nil
	}

	if targetValue.Kind() == reflect.Pointer {
		if r.DefaultValue == nil {
			targetValue.Set(reflect.Zero(targetValue.Type()))
			return nil
		}

		defaultValue := reflect.ValueOf(r.DefaultValue)

		if defaultValue.Type().AssignableTo(targetValue.Type()) {
			targetValue.Set(defaultValue)
			return nil
		}

		elemType := targetValue.Type().Elem()

		if !defaultValue.Type().AssignableTo(elemType) {
			if defaultValue.Type().ConvertibleTo(elemType) {
				defaultValue = defaultValue.Convert(elemType)
			} else {
				return nil
			}
		}

		newPtr := reflect.New(elemType)
		newPtr.Elem().Set(defaultValue)
		targetValue.Set(newPtr)
		return nil
	}

	if r.DefaultValue == nil {
		return nil
	}

	defaultValue := reflect.ValueOf(r.DefaultValue)

	if !defaultValue.Type().AssignableTo(targetValue.Type()) {
		if defaultValue.Type().ConvertibleTo(targetValue.Type()) {
			defaultValue = defaultValue.Convert(targetValue.Type())
		} else {
			return nil
		}
	}

	targetValue.Set(defaultValue)
	return nil
}
