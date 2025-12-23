package internal

import (
	"src/core"
)

const (
	DefaultCode    = "INTERNAL_ERROR"
	DefaultMessage = "An unexpected error occurred. Please try again later."
)

type Exception struct {
	core.Error
}

func New() *Exception {
	return &Exception{Error: core.Error{
		Code:    DefaultCode,
		Message: DefaultMessage,
	}}
}
