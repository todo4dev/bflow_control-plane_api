package binary

import (
	"src/core/validator/binary/rule"
	"src/core/validator/binary/schema"
)

type (
	LengthRule = rule.LengthRule
	MaxRule    = rule.MaxRule
	MinRule    = rule.MinRule
)

type BinarySchema = schema.BinarySchema

var Binary = schema.Binary
