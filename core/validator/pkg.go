package validator

import (
	"src/core/validator/_base"
	"src/core/validator/array"
	"src/core/validator/binary"
	"src/core/validator/boolean"
	"src/core/validator/condition"
	"src/core/validator/date"
	"src/core/validator/number"
	"src/core/validator/object"
	"src/core/validator/string"
	"src/core/validator/unknown"
)

const Err_ValidationError = _base.Err_ValidationError

type (
	IValidable      = _base.IValidable
	IRule           = _base.IRule
	Error           = _base.Error
	ValidationError = _base.ValidationError
)

type (
	ArraySchema     = array.ArraySchema
	BinarySchema    = binary.BinarySchema
	BooleanSchema   = boolean.BooleanSchema
	ConditionSchema = condition.ConditionSchema
	DateSchema      = date.DateSchema
	NumberSchema    = number.NumberSchema
	ObjectSchema    = object.ObjectSchema
	StringSchema    = string.StringSchema
	UnknownSchema   = unknown.UnknownSchema
)

var (
	Array     = array.Array
	Binary    = binary.Binary
	Boolean   = boolean.Boolean
	Condition = condition.Condition
	Date      = date.Date
	Number    = number.Number
	Object    = object.Object
	String    = string.String
	Unknown   = unknown.Unknown
)
