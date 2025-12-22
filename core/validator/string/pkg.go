package string

import (
	"src/core/validator/string/rule"
	"src/core/validator/string/schema"
)

type (
	Base64Rule     = rule.Base64Rule
	CaseRule       = rule.CaseRule
	CreditCardRule = rule.CreditCardRule
	IPRule         = rule.IPRule
	ISODateRule    = rule.ISODateRule
	LengthRule     = rule.LengthRule
	MaxRule        = rule.MaxRule
	MinRule        = rule.MinRule
	RegexRule      = rule.RegexRule
	ReplaceRule    = rule.ReplaceRule
	TransformRule  = rule.TransformRule
	TruncateRule   = rule.TruncateRule
	URIRule        = rule.URIRule
)

type StringSchema = schema.StringSchema

var String = schema.String
