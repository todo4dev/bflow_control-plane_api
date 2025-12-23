package validation

import (
	"src/core"
)

const (
	DefaultCode    = "VALIDATION_ERROR"
	DefaultMessage = "One or more fields are invalid"
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
