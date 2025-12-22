package _base

import (
	"reflect"
)

type IValidable interface {
	Validate() error
}

type IRule interface {
	Apply(valuePointer reflect.Value) []ValidationError
}

type ValidatorFunc func(value any, param any) error
