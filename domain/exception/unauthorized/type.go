package unauthorized

import (
	"src/core"
)

const (
	DefaultCode    = "UNAUTHORIZED"
	DefaultMessage = "Invalid authentication credentials"
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
