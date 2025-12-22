package exception

import (
	"src/core/common"
	"src/core/doc"
)

const (
	DefaultMethodNotAllowedExceptionCode    = "METHOD_NOT_ALLOWED"
	DefaultMethodNotAllowedExceptionMessage = "Method POST is not allowed for this resource"
)

type MethodNotAllowedException struct {
	common.ErrorPayload
}

func NewMethodNotAllowedException() *MethodNotAllowedException {
	return &MethodNotAllowedException{ErrorPayload: common.ErrorPayload{
		Code:    DefaultMethodNotAllowedExceptionCode,
		Message: DefaultMethodNotAllowedExceptionMessage,
	}}
}

func init() {
	methodNotAllowedException := NewMethodNotAllowedException()
	doc.Describe(methodNotAllowedException,
		doc.Description("HTTP method is not allowed for the requested resource"),
		doc.Field(
			&methodNotAllowedException.Code,
			doc.Description("Machine-readable error code"),
			doc.Example(DefaultMethodNotAllowedExceptionCode),
		),
		doc.Field(
			&methodNotAllowedException.Message,
			doc.Description("Human-readable error message"),
			doc.Example(DefaultMethodNotAllowedExceptionMessage),
		),
	)
}
