package core

import (
	"net/http"
	"reflect"

	"src/core/validator"
	"src/domain/exception"
)

var HTTPStatusMap = map[reflect.Type]int{
	reflect.TypeFor[*validator.ValidationError]():     http.StatusBadRequest,          // 400
	reflect.TypeFor[*exception.Validation]():          http.StatusBadRequest,          // 400
	reflect.TypeFor[*exception.Unauthorized]():        http.StatusUnauthorized,        // 401
	reflect.TypeFor[*exception.Forbidden]():           http.StatusForbidden,           // 403
	reflect.TypeFor[*exception.NotFound]():            http.StatusNotFound,            // 404
	reflect.TypeFor[*exception.Conflict]():            http.StatusConflict,            // 409
	reflect.TypeFor[*exception.UnprocessableEntity](): http.StatusUnprocessableEntity, // 422
	reflect.TypeFor[*exception.PreconditionFailed]():  http.StatusPreconditionFailed,  // 412
	reflect.TypeFor[*exception.MethodNotAllowed]():    http.StatusMethodNotAllowed,    // 405
	reflect.TypeFor[*exception.NotAcceptable]():       http.StatusNotAcceptable,       // 406
	reflect.TypeFor[*exception.Internal]():            http.StatusInternalServerError, // 500
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
