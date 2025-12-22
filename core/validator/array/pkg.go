package array

import (
	"src/core/validator/array/rule"
	"src/core/validator/array/schema"
)

type (
	HasRule    = rule.HasRule
	ItemsRule  = rule.ItemsRule
	LengthRule = rule.LengthRule
	MaxRule    = rule.MaxRule
	MinRule    = rule.MinRule
	UniqueRule = rule.UniqueRule
)

type ArraySchema = schema.ArraySchema

var Array = schema.Array
