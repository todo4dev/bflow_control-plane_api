package exception

import (
	"src/core/common"
	"src/core/doc"
)

const (
	DefaultUnauthorizedExceptionCode    = "UNAUTHORIZED"
	DefaultUnauthorizedExceptionMessage = "Invalid authentication credentials"
)

type UnauthorizedException struct {
	common.ErrorPayload
}

func NewUnauthorizedException() *UnauthorizedException {
	return &UnauthorizedException{ErrorPayload: common.ErrorPayload{
		Code:    DefaultUnauthorizedExceptionCode,
		Message: DefaultUnauthorizedExceptionMessage,
	}}
}

func init() {
	unauthorizedException := NewUnauthorizedException()
	doc.Describe(unauthorizedException,
		doc.Description("Authentication is required or the provided credentials are invalid"),
		doc.Field(
			&unauthorizedException.Code,
			doc.Description("Machine-readable error code"),
			doc.Example(DefaultUnauthorizedExceptionCode),
		),
		doc.Field(
			&unauthorizedException.Message,
			doc.Description("Human-readable error message"),
			doc.Example(DefaultUnauthorizedExceptionMessage),
		),
	)
}
