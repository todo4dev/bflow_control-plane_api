package schema

import (
	"regexp"
	"strings"

	"src/core/validator/_base"
	"src/core/validator/string/rule"
	"src/core/validator/unknown"
)

type StringSchema struct {
	unknown.UnknownSchema
}

func String(fieldPointer any) *StringSchema {
	return &StringSchema{UnknownSchema: unknown.Unknown(fieldPointer)}
}
func (builder *StringSchema) Custom(param any, validatorFunc _base.ValidatorFunc) *StringSchema {
	builder.UnknownSchema.Custom(param, validatorFunc)
	return builder
}
func (builder *StringSchema) Required() *StringSchema {
	builder.UnknownSchema.AddRule(unknown.RequiredRule{})
	return builder
}
func (builder *StringSchema) Default(value string) *StringSchema {
	builder.UnknownSchema.AddRule(unknown.DefaultRule{DefaultValue: value})
	return builder
}
func (builder *StringSchema) Allow(values ...string) *StringSchema {
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
func (builder *StringSchema) Alphanum() *StringSchema {
	builder.UnknownSchema.AddRule(rule.RegexRule{
		Code:    "string.alphanum",
		Message: "value must contain only alphanumeric characters",
		Pattern: RegexAlphanum,
	})
	return builder
}
func (builder *StringSchema) Base64() *StringSchema {
	builder.UnknownSchema.AddRule(rule.Base64Rule{})
	return builder
}
func (builder *StringSchema) Case(caseRuleMode rule.CaseRuleMode) *StringSchema {
	builder.UnknownSchema.AddRule(rule.CaseRule{CaseRuleMode: caseRuleMode})
	return builder
}
func (builder *StringSchema) CreditCard() *StringSchema {
	builder.UnknownSchema.AddRule(rule.CreditCardRule{})
	return builder
}
func (builder *StringSchema) DataURI() *StringSchema {
	builder.UnknownSchema.AddRule(rule.RegexRule{
		Code:    "string.dataUri",
		Message: "value must be a valid data URI",
		Pattern: RegexDataURI,
	})
	return builder
}
func (builder *StringSchema) Domain() *StringSchema {
	builder.UnknownSchema.AddRule(rule.RegexRule{
		Code:    "string.domain",
		Message: "value must be a valid domain name",
		Pattern: RegexDomain,
	})
	return builder
}
func (builder *StringSchema) Email() *StringSchema {
	builder.UnknownSchema.AddRule(rule.RegexRule{
		Code:    "string.email",
		Message: "value must be a valid email address",
		Pattern: RegexEmail,
	})
	return builder
}
func (builder *StringSchema) GUID() *StringSchema {
	builder.UnknownSchema.AddRule(rule.RegexRule{
		Code:    "string.guid",
		Message: "value must be a valid GUID/UUID",
		Pattern: RegexGUID,
	})
	return builder
}
func (builder *StringSchema) Hex() *StringSchema {
	builder.UnknownSchema.AddRule(rule.RegexRule{
		Code:    "string.hex",
		Message: "value must contain only hexadecimal characters",
		Pattern: RegexHex,
	})
	return builder
}
func (builder *StringSchema) Hostname() *StringSchema {
	builder.UnknownSchema.AddRule(rule.RegexRule{
		Code:    "string.hostname",
		Message: "value must be a valid hostname",
		Pattern: RegexHostname,
	})
	return builder
}
func (builder *StringSchema) IP() *StringSchema {
	builder.UnknownSchema.AddRule(rule.IPRule{})
	return builder
}
func (builder *StringSchema) IsoDate() *StringSchema {
	builder.UnknownSchema.AddRule(rule.ISODateRule{})
	return builder
}
func (builder *StringSchema) IsoDuration() *StringSchema {
	builder.UnknownSchema.AddRule(rule.RegexRule{
		Code:    "string.isoDuration",
		Message: "value must be a valid ISO-8601 duration",
		Pattern: RegexIsoDuration,
	})
	return builder
}
func (builder *StringSchema) Length(expected int) *StringSchema {
	builder.UnknownSchema.AddRule(rule.LengthRule{Length: expected})
	return builder
}
func (builder *StringSchema) Lowercase() *StringSchema {
	builder.UnknownSchema.AddRule(rule.TransformRule{
		Code: "string.lowercase",
		Transform: func(value string) string {
			return strings.ToLower(value)
		},
	})
	return builder
}
func (builder *StringSchema) Max(max int) *StringSchema {
	builder.UnknownSchema.AddRule(rule.MaxRule{Max: max})
	return builder
}
func (builder *StringSchema) Min(min int) *StringSchema {
	builder.UnknownSchema.AddRule(rule.MinRule{Min: min})
	return builder
}
func (builder *StringSchema) Pattern(pattern *regexp.Regexp) *StringSchema {
	if pattern == nil {
		return builder
	}

	builder.UnknownSchema.AddRule(rule.RegexRule{
		Code:    "string.pattern",
		Message: "value does not match required pattern",
		Pattern: pattern,
	})
	return builder
}
func (builder *StringSchema) Replace(pattern *regexp.Regexp, replacement string) *StringSchema {
	if pattern == nil {
		return builder
	}

	builder.UnknownSchema.AddRule(rule.ReplaceRule{
		Pattern:     pattern,
		Replacement: replacement,
	})
	return builder
}
func (builder *StringSchema) Token() *StringSchema {
	builder.UnknownSchema.AddRule(rule.RegexRule{
		Code:    "string.token",
		Message: "value must contain only word characters",
		Pattern: RegexToken,
	})
	return builder
}
func (builder *StringSchema) Trim() *StringSchema {
	builder.UnknownSchema.AddRule(rule.TransformRule{
		Code: "string.trim",
		Transform: func(value string) string {
			return strings.TrimSpace(value)
		},
	})
	return builder
}
func (builder *StringSchema) Truncate(max int) *StringSchema {
	if max <= 0 {
		return builder
	}

	builder.UnknownSchema.AddRule(rule.TruncateRule{Max: max})
	return builder
}
func (builder *StringSchema) Uppercase() *StringSchema {
	builder.UnknownSchema.AddRule(rule.TransformRule{
		Code: "string.uppercase",
		Transform: func(value string) string {
			return strings.ToUpper(value)
		},
	})
	return builder
}
func (builder *StringSchema) URI() *StringSchema {
	builder.UnknownSchema.AddRule(rule.URIRule{})
	return builder
}
func (builder *StringSchema) Validate() error {
	return builder.UnknownSchema.Validate()
}
