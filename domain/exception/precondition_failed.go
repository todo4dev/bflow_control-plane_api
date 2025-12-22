package exception

import (
	"src/core/common"
	"src/core/doc"
)

const (
	DefaultPreconditionFailedExceptionCode    = "PRECONDITION_FAILED"
	DefaultPreconditionFailedExceptionMessage = "Preconditions for this operation have not been satisfied"
)

type PreconditionFailedException struct {
	common.ErrorPayload
}

func NewPreconditionFailedException() *PreconditionFailedException {
	return &PreconditionFailedException{ErrorPayload: common.ErrorPayload{
		Code:    DefaultPreconditionFailedExceptionCode,
		Message: DefaultPreconditionFailedExceptionMessage,
	}}
}

func init() {
	preconditionFailedException := NewPreconditionFailedException()
	doc.Describe(preconditionFailedException,
		doc.Description("One or more preconditions required for this operation were not met"),
		doc.Example(preconditionFailedException),
		doc.Field(
			&preconditionFailedException.Code,
			doc.Description("Machine-readable error code"),
			doc.Example(DefaultPreconditionFailedExceptionCode),
		),
		doc.Field(
			&preconditionFailedException.Message,
			doc.Description("Human-readable error message"),
			doc.Example(DefaultPreconditionFailedExceptionMessage),
		),
	)
}
