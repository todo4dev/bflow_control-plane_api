package number

import (
	"src/core/validator/number/rule"
	"src/core/validator/number/schema"
)

type (
	GreaterRule   = rule.GreaterRule
	IntegerRule   = rule.IntegerRule
	LessRule      = rule.LessRule
	MaxRule       = rule.MaxRule
	MinRule       = rule.MinRule
	MultipleRule  = rule.MultipleRule
	NegativeRule  = rule.NegativeRule
	PortRule      = rule.PortRule
	PositiveRule  = rule.PositiveRule
	PrecisionRule = rule.PrecisionRule
)

type NumberSchema = schema.NumberSchema

var Number = schema.Number
