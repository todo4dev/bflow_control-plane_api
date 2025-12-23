package exception

import (
	"src/domain/exception/conflict"
	"src/domain/exception/forbidden"
	"src/domain/exception/internal"
	"src/domain/exception/method_not_allowed"
	"src/domain/exception/not_acceptable"
	"src/domain/exception/not_found"
	"src/domain/exception/precondition_failed"
	"src/domain/exception/unauthorized"
	"src/domain/exception/unprocessable_entity"
	"src/domain/exception/validation"
)

type (
	Conflict            = conflict.Exception
	Forbidden           = forbidden.Exception
	Internal            = internal.Exception
	MethodNotAllowed    = method_not_allowed.Exception
	NotAcceptable       = not_acceptable.Exception
	NotFound            = not_found.Exception
	PreconditionFailed  = precondition_failed.Exception
	Unauthorized        = unauthorized.Exception
	UnprocessableEntity = unprocessable_entity.Exception
	Validation          = validation.Exception
)

var (
	NewConflict            = conflict.New
	NewForbidden           = forbidden.New
	NewInternal            = internal.New
	NewMethodNotAllowed    = method_not_allowed.New
	NewNotAcceptable       = not_acceptable.New
	NewNotFound            = not_found.New
	NewPreconditionFailed  = precondition_failed.New
	NewUnauthorized        = unauthorized.New
	NewUnprocessableEntity = unprocessable_entity.New
	NewValidation          = validation.New
)

func Register() {
	conflict.Register()
	forbidden.Register()
	internal.Register()
	method_not_allowed.Register()
	not_acceptable.Register()
	not_found.Register()
	precondition_failed.Register()
	unauthorized.Register()
	unprocessable_entity.Register()
	validation.Register()
}
