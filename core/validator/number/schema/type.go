package schema

import (
	"src/core/validator/_base"
	"src/core/validator/number/rule"
	"src/core/validator/unknown"
)

type NumberSchema struct {
	unknown.UnknownSchema
}

func Number(fieldPointer any) *NumberSchema {
	return &NumberSchema{UnknownSchema: unknown.Unknown(fieldPointer)}
}
func (builder *NumberSchema) Custom(param any, validatorFunc _base.ValidatorFunc) *NumberSchema {
	builder.UnknownSchema.Custom(param, validatorFunc)
	return builder
}
func (builder *NumberSchema) Required() *NumberSchema {
	builder.UnknownSchema.AddRule(unknown.RequiredRule{})
	return builder
}
func (builder *NumberSchema) Default(value float64) *NumberSchema {
	builder.UnknownSchema.AddRule(unknown.DefaultRule{DefaultValue: value})
	return builder
}
func (builder *NumberSchema) Allow(values ...float64) *NumberSchema {
	if len(values) == 0 {
		return builder
	}

	converted := make([]any, len(values))
	for index, value := range values {
		converted[index] = value
	}

	builder.UnknownSchema.AddRule(unknown.AllowRule{Allowed: converted})
	return builder
}
func (builder *NumberSchema) Min(limit float64) *NumberSchema {
	builder.UnknownSchema.AddRule(rule.MinRule{Limit: limit})
	return builder
}
func (builder *NumberSchema) Max(limit float64) *NumberSchema {
	builder.UnknownSchema.AddRule(rule.MaxRule{Limit: limit})
	return builder
}
func (builder *NumberSchema) Greater(limit float64) *NumberSchema {
	builder.UnknownSchema.AddRule(rule.GreaterRule{Limit: limit})
	return builder
}
func (builder *NumberSchema) Less(limit float64) *NumberSchema {
	builder.UnknownSchema.AddRule(rule.LessRule{Limit: limit})
	return builder
}
func (builder *NumberSchema) Integer() *NumberSchema {
	builder.UnknownSchema.AddRule(rule.IntegerRule{})
	return builder
}
func (builder *NumberSchema) Multiple(base float64) *NumberSchema {
	if base == 0 {
		return builder
	}
	builder.UnknownSchema.AddRule(rule.MultipleRule{Base: base})
	return builder
}
func (builder *NumberSchema) Negative() *NumberSchema {
	builder.UnknownSchema.AddRule(rule.NegativeRule{})
	return builder
}
func (builder *NumberSchema) Positive() *NumberSchema {
	builder.UnknownSchema.AddRule(rule.PositiveRule{})
	return builder
}
func (builder *NumberSchema) Port() *NumberSchema {
	builder.UnknownSchema.AddRule(rule.PortRule{})
	return builder
}
func (builder *NumberSchema) Precision(precision int) *NumberSchema {
	if precision < 0 {
		return builder
	}
	builder.UnknownSchema.AddRule(rule.PrecisionRule{Precision: precision})
	return builder
}
func (builder *NumberSchema) Validate() error {
	return builder.UnknownSchema.Validate()
}
