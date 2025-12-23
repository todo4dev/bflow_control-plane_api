package precondition_failed

import (
	"src/core"
)

const (
	DefaultCode    = "PRECONDITION_FAILED"
	DefaultMessage = "Preconditions for this operation have not been satisfied"
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
