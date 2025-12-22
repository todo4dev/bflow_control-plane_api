package schema

import (
	"src/core/validator/_base"
	"src/core/validator/array/rule"
	"src/core/validator/unknown"
)

type ArraySchema struct {
	unknown.UnknownSchema
}

func Array(fieldPointer any) *ArraySchema {
	return &ArraySchema{UnknownSchema: unknown.Unknown(fieldPointer)}
}
func (s *ArraySchema) Custom(param any, validatorFunc _base.ValidatorFunc, optionalField ...string) *ArraySchema {
	s.UnknownSchema.Custom(param, validatorFunc)
	return s
}
func (s *ArraySchema) Required() *ArraySchema {
	s.UnknownSchema.AddRule(unknown.RequiredRule{})
	return s
}
func (s *ArraySchema) Optional() *ArraySchema {
	return s
}
func (s *ArraySchema) Default(value any) *ArraySchema {
	s.UnknownSchema.AddRule(unknown.DefaultRule{DefaultValue: value})
	return s
}
func (s *ArraySchema) Allow(values ...any) *ArraySchema {
	if len(values) == 0 {
		return s
	}

	allowed := make([]any, len(values))
	copy(allowed, values)

	s.UnknownSchema.AddRule(unknown.AllowRule{Allowed: allowed})
	return s
}
func (s *ArraySchema) Validate() error {
	return s.UnknownSchema.Validate()
}
func (s *ArraySchema) Length(length int) *ArraySchema {
	if length < 0 {
		return s
	}
	s.UnknownSchema.AddRule(rule.LengthRule{Length: length})
	return s
}
func (s *ArraySchema) Min(minLength int) *ArraySchema {
	if minLength < 0 {
		return s
	}
	s.UnknownSchema.AddRule(rule.MinRule{Min: minLength})
	return s
}
func (s *ArraySchema) Max(max int) *ArraySchema {
	if max < 0 {
		return s
	}
	s.UnknownSchema.AddRule(rule.MaxRule{Max: max})
	return s
}
func (s *ArraySchema) Has(expected any) *ArraySchema {
	s.UnknownSchema.AddRule(rule.HasRule{Expected: expected})
	return s
}
func (s *ArraySchema) Unique() *ArraySchema {
	s.UnknownSchema.AddRule(rule.UniqueRule{})
	return s
}
func (s *ArraySchema) Items() *ArraySchema {
	s.UnknownSchema.AddRule(rule.ItemsRule{})
	return s
}
