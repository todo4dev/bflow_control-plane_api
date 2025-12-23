package unprocessable_entity

import (
	"src/core"
)

const (
	DefaultCode    = "UNPROCESSABLE_ENTITY"
	DefaultMessage = "The request could not be processed due to semantic errors"
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
