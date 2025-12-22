package schema

import (
	"fmt"
	"reflect"
	"strings"

	"src/core/validator/_base"
	"src/core/validator/unknown"
)

type ObjectSchema struct {
	unknown.UnknownSchema
	targetPointer any
	fields        []_base.IValidable
}

func Object(targetPointer any, fields ..._base.IValidable) *ObjectSchema {
	if targetPointer == nil {
		panic("validator: targetPointer cannot be nil in Object()")
	}

	return &ObjectSchema{
		targetPointer: targetPointer,
		fields:        append([]_base.IValidable(nil), fields...),
	}
}

func (s *ObjectSchema) Fields(fields ..._base.IValidable) *ObjectSchema {
	if len(fields) == 0 {
		return s
	}

	s.fields = append(s.fields, fields...)
	return s
}

type targetPointerProvider interface {
	TargetPointer() any
}

func buildFieldPointerMap(structPointer any) map[uintptr]string {
	result := map[uintptr]string{}

	value := reflect.ValueOf(structPointer)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return result
	}

	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return result
	}

	typ := value.Type()

	for index := 0; index < value.NumField(); index++ {
		fieldValue := value.Field(index)
		fieldType := typ.Field(index)

		if !fieldValue.CanAddr() {
			continue
		}

		fieldName := fieldType.Name

		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
			tagName := strings.Split(jsonTag, ",")[0]
			if tagName != "" {
				fieldName = tagName
			}
		}

		result[fieldValue.Addr().Pointer()] = fieldName
	}

	return result
}

func pointerOf(value any) (uintptr, bool) {
	if value == nil {
		return 0, false
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return 0, false
	}

	return v.Pointer(), true
}

func extractSchemaTargetPointerAddress(schema _base.IValidable) (uintptr, bool) {
	if schema == nil {
		return 0, false
	}

	provider, ok := schema.(targetPointerProvider)
	if !ok {
		return 0, false
	}

	return pointerOf(provider.TargetPointer())
}

func (s *ObjectSchema) Validate() error {
	if s.targetPointer == nil {
		panic("validator: targetPointer cannot be nil in ObjectSchema.Validate()")
	}

	targetValue := reflect.ValueOf(s.targetPointer)
	if targetValue.Kind() != reflect.Pointer || targetValue.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("validator: Object() expects a pointer to struct, got %s", targetValue.Kind().String()))
	}

	if len(s.fields) == 0 {
		if validable, ok := s.targetPointer.(_base.IValidable); ok && validable != nil {
			err := validable.Validate()
			if err == nil {
				return nil
			}

			if validationError, ok := err.(*_base.Error); ok {
				if len(validationError.Errors) == 0 {
					return nil
				}
				return validationError
			}

			return &_base.Error{Errors: []_base.ValidationError{{
				Code:    "struct",
				Message: err.Error(),
			}}}
		}

		return nil
	}

	pointerToFieldName := buildFieldPointerMap(s.targetPointer)

	var allErrors []_base.ValidationError

	for _, fieldSchema := range s.fields {
		if fieldSchema == nil {
			continue
		}

		fieldName := ""
		if pointerAddress, ok := extractSchemaTargetPointerAddress(fieldSchema); ok {
			if name, exists := pointerToFieldName[pointerAddress]; exists && name != "" {
				fieldName = name
			}
		}

		err := fieldSchema.Validate()
		if err == nil {
			continue
		}

		if validationError, ok := err.(*_base.Error); ok {
			fieldErrors := make([]_base.ValidationError, len(validationError.Errors))
			copy(fieldErrors, validationError.Errors)

			if fieldName != "" {
				for index := range fieldErrors {
					if fieldErrors[index].Field == "" {
						fieldErrors[index].Field = fieldName
					}
				}
			}

			allErrors = append(allErrors, fieldErrors...)
			continue
		}

		allErrors = append(allErrors, _base.ValidationError{
			Field:   fieldName,
			Code:    "struct.field",
			Message: err.Error(),
		})
	}

	if len(allErrors) == 0 {
		return nil
	}

	return &_base.Error{Errors: allErrors}
}
