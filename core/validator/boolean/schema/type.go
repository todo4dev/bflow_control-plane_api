package schema

import (
	"src/core/validator/_base"
	"src/core/validator/boolean/rule"
	"src/core/validator/unknown"
)

type BooleanSchema struct {
	unknown.UnknownSchema
}

func Boolean(fieldPointer any) *BooleanSchema {
	return &BooleanSchema{UnknownSchema: unknown.Unknown(fieldPointer)}
}
func (s *BooleanSchema) Custom(param any, validatorFunc _base.ValidatorFunc) *BooleanSchema {
	s.UnknownSchema.Custom(param, validatorFunc)
	return s
}
func (s *BooleanSchema) Required() *BooleanSchema {
	s.UnknownSchema.Required()
	return s
}
func (s *BooleanSchema) Default(value bool) *BooleanSchema {
	s.UnknownSchema.AddRule(unknown.DefaultRule{DefaultValue: value})
	return s
}
func (s *BooleanSchema) Validate() error {
	return s.UnknownSchema.Validate()
}
func (s *BooleanSchema) Truthy(values ...bool) *BooleanSchema {
	s.UnknownSchema.AddRule(rule.TruthyRule{TruthyValues: values})
	return s
}
func (s *BooleanSchema) Falsy(values ...bool) *BooleanSchema {
	s.UnknownSchema.AddRule(rule.FalsyRule{FalsyValues: values})
	return s
}
