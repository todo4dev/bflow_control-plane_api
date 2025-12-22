package exception

import (
	"src/core/doc"

	"src/core/common"
)

const (
	DefaultConflitExceptionCode    = "CONFLICT"
	DefaultConflitExceptionMessage = "Resource is already in the desired state"
)

type ConflictException struct {
	common.ErrorPayload
}

func NewConflictException() *ConflictException {
	return &ConflictException{ErrorPayload: common.ErrorPayload{
		Code:    DefaultConflitExceptionCode,
		Message: DefaultConflitExceptionMessage,
	}}
}

func init() {
	conflictException := NewConflictException()
	doc.Describe(conflictException,
		doc.Description("Request could not be completed due to a conflict with the current state of the resource"),
		doc.Example(conflictException),
		doc.Field(
			&conflictException.Code,
			doc.Description("Machine-readable error code"),
			doc.Example(DefaultConflitExceptionCode),
		),
		doc.Field(
			&conflictException.Message,
			doc.Description("Human-readable error message"),
			doc.Example(DefaultConflitExceptionMessage),
		),
	)
}
