package not_found

import (
	"src/core"
)

const (
	DefaultCode    = "NOT_FOUND"
	DefaultMessage = "The requested resource was not found"
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
