package not_acceptable

import (
	"src/core"
)

const (
	DefaultCode    = "NOT_ACCEPTABLE"
	DefaultMessage = "Requested media type is not supported"
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
