package exception

import (
	"github.com/todo4dev/bflow/control-plane/api/core/common"
	"github.com/todo4dev/bflow/control-plane/api/core/doc"
)

const (
	DefaultNotFoundExceptionCode    = "NOT_FOUND"
	DefaultNotFoundExceptionMessage = "The requested resource was not found"
)

type NotFoundException struct {
	common.ErrorPayload
}

func NewNotFoundException() *NotFoundException {
	return &NotFoundException{ErrorPayload: common.ErrorPayload{
		Code:    DefaultNotFoundExceptionCode,
		Message: DefaultNotFoundExceptionMessage,
	}}
}

func init() {
	notFoundException := NewNotFoundException()
	doc.Describe(
		&notFoundException,
		doc.Description("The requested resource could not be found"),
		doc.Example(notFoundException),
		doc.Field(
			&notFoundException.Code,
			doc.Description("Machine-readable error code"),
			doc.Example(DefaultNotFoundExceptionCode),
		),
		doc.Field(
			&notFoundException.Message,
			doc.Description("Human-readable error message"),
			doc.Example(DefaultNotFoundExceptionMessage),
		),
	)
}
