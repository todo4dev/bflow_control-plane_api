package rule

import (
	"reflect"
	"strings"

	"src/core/validator/_base"
	"src/core/validator/string/util"
)

type CreditCardRule struct{}

var _ _base.IRule = CreditCardRule{}

func (r CreditCardRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
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
			Code:    "string.creditCard",
			Message: "string must be a string",
		}}
	}

	value := strings.ReplaceAll(target.String(), " ", "")
	value = strings.ReplaceAll(value, "-", "")

	if value == "" {
		return nil
	}

	if len(value) < 12 || len(value) > 19 {
		return []_base.ValidationError{{
			Code:    "string.creditCard",
			Message: "string must be a valid credit card number",
		}}
	}

	for _, character := range value {
		if character < '0' || character > '9' {
			return []_base.ValidationError{{
				Code:    "string.creditCard",
				Message: "string must be a valid credit card number",
			}}
		}
	}

	if !util.LuhnValid(value) {
		return []_base.ValidationError{{
			Code:    "string.creditCard",
			Message: "string must be a valid credit card number",
		}}
	}

	return nil
}
