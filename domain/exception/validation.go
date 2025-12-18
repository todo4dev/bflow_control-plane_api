package exception

import (
	"github.com/todo4dev/bflow/control-plane/api/core/common"
	"github.com/todo4dev/bflow/control-plane/api/core/doc"
)

const (
	DefaultValidationExceptionCode    = "VALIDATION_ERROR"
	DefaultValidationExceptionMessage = "One or more fields are invalid"
)

type ValidationException struct {
	common.ErrorPayload
}

func NewValidationException() *ValidationException {
	return &ValidationException{ErrorPayload: common.ErrorPayload{
		Code:    DefaultValidationExceptionCode,
		Message: DefaultValidationExceptionMessage,
	}}
}

func init() {
	validationException := NewValidationException()
	doc.Describe(
		&validationException,
		doc.Description("The request contains invalid or malformed data"),
		doc.Example(validationException),
		doc.Field(
			&validationException.Code,
			doc.Description("Machine-readable error code"),
			doc.Example(DefaultValidationExceptionCode),
		),
		doc.Field(
			&validationException.Message,
			doc.Description("Human-readable error message"),
			doc.Example(DefaultValidationExceptionMessage),
		),
	)
}
