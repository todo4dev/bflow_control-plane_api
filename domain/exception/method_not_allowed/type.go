package method_not_allowed

import (
	"src/core"
)

const (
	DefaultCode    = "METHOD_NOT_ALLOWED"
	DefaultMessage = "Method POST is not allowed for this resource"
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
