package exception

import (
	"src/core/common"
	"src/core/doc"
)

const (
	DefaultUnprocessableEntityExceptionCode    = "UNPROCESSABLE_ENTITY"
	DefaultUnprocessableEntityExceptionMessage = "The request could not be processed due to semantic errors"
)

type UnprocessableEntityException struct {
	common.ErrorPayload
}

func NewUnprocessableEntityException() *UnprocessableEntityException {
	return &UnprocessableEntityException{ErrorPayload: common.ErrorPayload{
		Code:    DefaultUnprocessableEntityExceptionCode,
		Message: DefaultUnprocessableEntityExceptionMessage,
	}}
}

func init() {
	unprocessableEntityException := NewUnprocessableEntityException()
	doc.Describe(unprocessableEntityException,
		doc.Description("The request was well-formed but contains semantic errors and could not be processed"),
		doc.Example(unprocessableEntityException),
		doc.Field(
			&unprocessableEntityException.Code,
			doc.Description("Machine-readable error code"),
			doc.Example(DefaultUnprocessableEntityExceptionCode),
		),
		doc.Field(
			&unprocessableEntityException.Message,
			doc.Description("Human-readable error message"),
			doc.Example(DefaultUnprocessableEntityExceptionMessage),
		),
	)
}
