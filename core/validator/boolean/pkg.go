package boolean

import (
	"src/core/validator/boolean/rule"
	"src/core/validator/boolean/schema"
)

type (
	FalsyRule  = rule.FalsyRule
	TruthyRule = rule.TruthyRule
)

type BooleanSchema = schema.BooleanSchema

var Boolean = schema.Boolean
