// core/http_errors.go
package core

import (
	"net/http"
	"reflect"

	"src/core/validator"
	"src/domain/exception"
)

var HTTPStatusMap = map[reflect.Type]int{
	reflect.TypeFor[*validator.ValidationError]():              http.StatusBadRequest,          // 400
	reflect.TypeFor[*exception.ValidationException]():          http.StatusBadRequest,          // 400
	reflect.TypeFor[*exception.UnauthorizedException]():        http.StatusUnauthorized,        // 401
	reflect.TypeFor[*exception.ForbiddenException]():           http.StatusForbidden,           // 403
	reflect.TypeFor[*exception.NotFoundException]():            http.StatusNotFound,            // 404
	reflect.TypeFor[*exception.ConflictException]():            http.StatusConflict,            // 409
	reflect.TypeFor[*exception.UnprocessableEntityException](): http.StatusUnprocessableEntity, // 422
	reflect.TypeFor[*exception.PreconditionFailedException]():  http.StatusPreconditionFailed,  // 412
	reflect.TypeFor[*exception.MethodNotAllowedException]():    http.StatusMethodNotAllowed,    // 405
	reflect.TypeFor[*exception.NotAcceptableException]():       http.StatusNotAcceptable,       // 406
	reflect.TypeFor[*exception.InternalException]():            http.StatusInternalServerError, // 500
}

func GetHTTPStatus(err any) int {
	if err == nil {
		return http.StatusOK
	}

	t := reflect.TypeOf(err)

	if t.Kind() == reflect.Pointer {
		if status, ok := HTTPStatusMap[t]; ok {
			return status
		}
	}

	if t.Kind() != reflect.Pointer {
		ptrType := reflect.PointerTo(t)
		if status, ok := HTTPStatusMap[ptrType]; ok {
			return status
		}
	}

	return http.StatusInternalServerError
}
