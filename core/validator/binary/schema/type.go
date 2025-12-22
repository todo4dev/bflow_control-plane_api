package schema

import (
	"src/core/validator/_base"
	"src/core/validator/binary/rule"
	"src/core/validator/unknown"
)

type BinarySchema struct {
	unknown.UnknownSchema
}

func Binary(fieldPointer any) *BinarySchema {
	return &BinarySchema{UnknownSchema: unknown.Unknown(fieldPointer)}
}
func (builder *BinarySchema) Custom(param any, validatorFunc _base.ValidatorFunc) *BinarySchema {
	builder.UnknownSchema.Custom(param, validatorFunc)
	return builder
}
func (builder *BinarySchema) Required() *BinarySchema {
	builder.UnknownSchema.Required()
	return builder
}
func (builder *BinarySchema) Default(value []byte) *BinarySchema {
	builder.UnknownSchema.Default(value)
	return builder
}
func (builder *BinarySchema) Allow(values ...[]byte) *BinarySchema {
	if len(values) == 0 {
		return builder
	}

	converted := make([]any, len(values))
	for index, value := range values {
		converted[index] = value
	}

	builder.UnknownSchema.Allow(converted)
	return builder
}
func (builder *BinarySchema) Length(expectedLength int) *BinarySchema {
	if expectedLength < 0 {
		return builder
	}

	builder.UnknownSchema.AddRule(rule.LengthRule{Length: expectedLength})
	return builder
}
func (builder *BinarySchema) Min(minLength int) *BinarySchema {
	if minLength < 0 {
		return builder
	}

	builder.UnknownSchema.AddRule(rule.MinRule{Min: minLength})
	return builder
}
func (builder *BinarySchema) Max(maxLength int) *BinarySchema {
	if maxLength < 0 {
		return builder
	}

	builder.UnknownSchema.AddRule(rule.MaxRule{Max: maxLength})
	return builder
}
func (builder *BinarySchema) Validate() error {
	return builder.UnknownSchema.Validate()
}
