package exception

import (
	"github.com/todo4dev/bflow/control-plane/api/core/common"
	"github.com/todo4dev/bflow/control-plane/api/core/doc"
)

const (
	DefaultInternalExceptionCode    = "INTERNAL_ERROR"
	DefaultInternalExceptionMessage = "An unexpected error occurred. Please try again later."
)

type InternalException struct {
	common.ErrorPayload
}

func NewInternalException() *InternalException {
	return &InternalException{ErrorPayload: common.ErrorPayload{
		Code:    DefaultInternalExceptionCode,
		Message: DefaultInternalExceptionMessage,
	}}
}

func init() {
	internalException := NewInternalException()
	doc.Describe(internalException,
		doc.Description("An unexpected internal error occurred while processing the request"),
		doc.Example(internalException),
		doc.Field(&internalException.Code,
			doc.Description("Machine-readable error code"),
			doc.Example(DefaultInternalExceptionCode),
		),
		doc.Field(&internalException.Message,
			doc.Description("Human-readable error message"),
			doc.Example(DefaultInternalExceptionMessage),
		),
	)
}
