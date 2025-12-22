package date

import (
	"src/core/validator/date/rule"
	"src/core/validator/date/schema"
)

type (
	GreaterRule = rule.GreaterRule
	LessRule    = rule.LessRule
	MaxRule     = rule.MaxRule
	MinRule     = rule.MinRule
)

type DateSchema = schema.DateSchema

var Date = schema.Date
