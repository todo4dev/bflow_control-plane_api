package _base

import "strings"

const Err_ValidationError = "validation failed"

type Error struct {
	Errors []ValidationError
}

type ValidationError struct {
	Field   string
	Code    string
	Message string
}

func (e *Error) Error() string {
	if len(e.Errors) == 0 {
		return "no validation errors"
	}

	var builder strings.Builder
	builder.WriteString("validation failed: ")

	for index, fieldError := range e.Errors {
		if index > 0 {
			builder.WriteString("; ")
		}

		field := fieldError.Field
		if field != "" {
			builder.WriteString("[")
			builder.WriteString(field)
			builder.WriteString("] ")
		}

		builder.WriteString(fieldError.Code)
		builder.WriteString(": ")
		builder.WriteString(fieldError.Message)
	}

	return builder.String()
}

func (e *Error) HasErrors() bool {
	return len(e.Errors) > 0
}
