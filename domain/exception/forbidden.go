package exception

import (
	"github.com/todo4dev/bflow_control-plane_api/core/common"
	"github.com/todo4dev/bflow_control-plane_api/core/doc"
)

const (
	DefaultForbiddenExceptionCode    = "FORBIDDEN"
	DefaultForbiddenExceptionMessage = "You do not have permission to access this resource"
)

type ForbiddenException struct {
	common.ErrorPayload
}

func NewForbiddenException() *ForbiddenException {
	return &ForbiddenException{ErrorPayload: common.ErrorPayload{
		Code:    DefaultForbiddenExceptionCode,
		Message: DefaultForbiddenExceptionMessage,
	}}
}

func init() {
	forbiddenException := NewForbiddenException()
	doc.Describe(
		forbiddenException,
		doc.Description("The authenticated user does not have permission to perform this operation"),
		doc.Field(
			&forbiddenException.Code,
			doc.Description("Machine-readable error code"),
			doc.Example(DefaultForbiddenExceptionCode),
		),
		doc.Field(
			&forbiddenException.Message,
			doc.Description("Human-readable error message"),
			doc.Example(DefaultForbiddenExceptionMessage),
		),
	)
}
