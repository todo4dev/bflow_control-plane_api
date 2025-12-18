package exception

import (
	"github.com/todo4dev/bflow/control-plane/api/core/common"
	"github.com/todo4dev/bflow/control-plane/api/core/doc"
)

const (
	DefaultNotAcceptableExceptionCode    = "NOT_ACCEPTABLE"
	DefaultNotAcceptableExceptionMessage = "Requested media type is not supported"
)

type NotAcceptableException struct {
	common.ErrorPayload
}

func NewNotAcceptableException() *NotAcceptableException {
	return &NotAcceptableException{ErrorPayload: common.ErrorPayload{
		Code:    DefaultNotAcceptableExceptionCode,
		Message: DefaultNotAcceptableExceptionMessage,
	}}
}

func init() {
	notAcceptableException := NewNotAcceptableException()
	doc.Describe(
		&notAcceptableException,
		doc.Description("Requested representation cannot be served (content negotiation failed)"),
		doc.Field(
			&notAcceptableException.Code,
			doc.Description("Machine-readable error code"),
			doc.Example(DefaultNotAcceptableExceptionCode),
		),
		doc.Field(
			&notAcceptableException.Message,
			doc.Description("Human-readable error message"),
			doc.Example(DefaultNotAcceptableExceptionMessage),
		),
	)
}
