package schema

import (
	"fmt"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/unknown/rule"
)

type UnknownSchema struct {
	Label        string
	fieldPointer any
	rules        []_base.IRule
}

func Unknown(fieldPointer any) UnknownSchema {
	if fieldPointer == nil {
		panic("validator: fieldPointer cannot be nil")
	}

	return UnknownSchema{
		Label:        "value",
		fieldPointer: fieldPointer,
		rules:        make([]_base.IRule, 0),
	}
}
func (s *UnknownSchema) TargetPointer() any {
	return s.fieldPointer
}
func (s *UnknownSchema) AddRule(rule _base.IRule) {
	if rule == nil {
		return
	}
	s.rules = append(s.rules, rule)
}
func (s *UnknownSchema) Validate() error {
	valuePointer := reflect.ValueOf(s.fieldPointer)
	if valuePointer.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("validator: schema expects a pointer, got %s", valuePointer.Kind().String()))
	}

	for _, currentRule := range s.rules {
		if currentRule == nil {
			continue
		}
		if _, ok := currentRule.(rule.DefaultRule); ok {
			_ = currentRule.Apply(valuePointer)
		}
	}

	var errors _base.Error

	for _, currentRule := range s.rules {
		if currentRule == nil {
			continue
		}
		if _, ok := currentRule.(rule.DefaultRule); ok {
			continue
		}

		appliedErrors := currentRule.Apply(valuePointer)
		if len(appliedErrors) > 0 {
			errors.Errors = append(errors.Errors, appliedErrors...)
		}
	}

	if len(errors.Errors) == 0 {
		return nil
	}

	return &errors
}
func (s *UnknownSchema) Custom(param any, validatorFunc _base.ValidatorFunc) *UnknownSchema {
	if validatorFunc == nil {
		return s
	}

	s.rules = append(s.rules, rule.CustomRule{
		Param:         param,
		ValidatorFunc: validatorFunc,
	})

	return s
}
func (s *UnknownSchema) Required() *UnknownSchema {
	s.AddRule(rule.RequiredRule{})
	return s
}
func (s *UnknownSchema) Default(defaultValue any) *UnknownSchema {
	s.AddRule(rule.DefaultRule{DefaultValue: defaultValue})
	return s
}
func (s *UnknownSchema) Allow(allowedValues ...any) *UnknownSchema {
	if len(allowedValues) == 0 {
		return s
	}
	s.AddRule(rule.AllowRule{Allowed: allowedValues})
	return s
}
