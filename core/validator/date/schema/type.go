package schema

import (
	"time"

	"src/core/validator/_base"
	"src/core/validator/date/rule"
	"src/core/validator/unknown"
)

type DateSchema struct {
	unknown.UnknownSchema
}

func Date(fieldPointer any) *DateSchema {
	return &DateSchema{UnknownSchema: unknown.Unknown(fieldPointer)}
}
func (s *DateSchema) Custom(param any, validatorFunc _base.ValidatorFunc) *DateSchema {
	s.UnknownSchema.Custom(param, validatorFunc)
	return s
}
func (s *DateSchema) Required() *DateSchema {
	s.UnknownSchema.AddRule(unknown.RequiredRule{})
	return s
}
func (s *DateSchema) Default(value any) *DateSchema {
	s.UnknownSchema.AddRule(unknown.DefaultRule{DefaultValue: value})
	return s
}
func (s *DateSchema) Allow(values ...time.Time) *DateSchema {
	if len(values) == 0 {
		return s
	}

	converted := make([]any, len(values))
	for index, value := range values {
		converted[index] = value
	}

	s.UnknownSchema.AddRule(unknown.AllowRule{Allowed: converted})
	return s
}
func (s *DateSchema) Min(limit time.Time) *DateSchema {
	s.UnknownSchema.AddRule(rule.MinRule{Limit: limit})
	return s
}
func (s *DateSchema) Max(limit time.Time) *DateSchema {
	s.UnknownSchema.AddRule(rule.MaxRule{Limit: limit})
	return s
}
func (s *DateSchema) Greater(limit time.Time) *DateSchema {
	s.UnknownSchema.AddRule(rule.GreaterRule{Limit: limit})
	return s
}
func (s *DateSchema) Less(limit time.Time) *DateSchema {
	s.UnknownSchema.AddRule(rule.LessRule{Limit: limit})
	return s
}
func (s *DateSchema) Validate() error {
	return s.UnknownSchema.Validate()
}
