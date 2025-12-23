package forbidden

import (
	"src/core"
)

const (
	DefaultCode    = "FORBIDDEN"
	DefaultMessage = "You do not have permission to access this resource"
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
