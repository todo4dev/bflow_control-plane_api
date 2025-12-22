package rule

import (
	"reflect"
	"strings"

	"src/core/validator/_base"
)

type CaseRuleMode int

const (
	CaseRuleMode_Lower CaseRuleMode = iota
	CaseRuleMode_Upper
)

type CaseRule struct {
	CaseRuleMode
}

var _ _base.IRule = CaseRule{}

func (r CaseRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	if !valuePointer.IsValid() {
		return nil
	}

	target := valuePointer.Elem()
	if target.Kind() == reflect.Pointer {
		if target.IsNil() {
			return nil
		}
		target = target.Elem()
	}

	if target.Kind() != reflect.String {
		return []_base.ValidationError{{
			Code:    "string.case",
			Message: "string must be a string",
		}}
	}

	value := target.String()
	if value == "" {
		return nil
	}

	switch r.CaseRuleMode {
	case CaseRuleMode_Lower:
		if value != strings.ToLower(value) {
			return []_base.ValidationError{{
				Code:    "string.case.lower",
				Message: "string must be lower case",
			}}
		}
	case CaseRuleMode_Upper:
		if value != strings.ToUpper(value) {
			return []_base.ValidationError{{
				Code:    "string.case.upper",
				Message: "string must be upper case",
			}}
		}
	default:
	}

	return nil
}
