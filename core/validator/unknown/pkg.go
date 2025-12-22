package unknown

import (
	"src/core/validator/unknown/rule"
	"src/core/validator/unknown/schema"
)

type (
	AllowRule    = rule.AllowRule
	CustomRule   = rule.CustomRule
	DefaultRule  = rule.DefaultRule
	RequiredRule = rule.RequiredRule
)

type UnknownSchema = schema.UnknownSchema

var Unknown = schema.Unknown
