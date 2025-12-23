package conflict

import (
	"src/core"
)

const (
	DefaultCode    = "CONFLICT"
	DefaultMessage = "Resource is already in the desired state"
)

type Exception struct {
	core.Error
}

func New() *Exception {
	return &Exception{
		Error: core.Error{
			Code:    DefaultCode,
			Message: DefaultMessage,
		},
	}
}
