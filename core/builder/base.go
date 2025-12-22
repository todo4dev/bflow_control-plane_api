package builder

import (
	"fmt"
	"reflect"
	"strings"
)

type builderBase struct{}

func (b *builderBase) fieldPointerJSONTag(fieldPointer any) string {
	pointerValue := reflect.ValueOf(fieldPointer)
	if pointerValue.Kind() != reflect.Pointer {
		panic("fieldPointer must be a pointer to a struct field.")
	}

	parentObject := pointerValue.Elem().Addr()
	if parentObject.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("Invalid pointer type: expected addressable struct element, got %v", parentObject.Kind()))
	}

	parentObjectType := parentObject.Elem().Type()
	for i := 0; i < parentObjectType.NumField(); i++ {
		structField := parentObjectType.Field(i)
		if !structField.IsExported() {
			continue
		}

		field := parentObject.Elem().Field(i)
		if field.CanAddr() && field.Addr().Interface() == fieldPointer {
			tag := structField.Tag.Get("json")
			if tag == "" {
				panic("field must contain \"json\" tag")
			}

			if idx := strings.Index(tag, ","); idx != -1 {
				tag = tag[:idx]
			}

			if tag == "" || tag == "-" {
				panic("field must have a valid json name")
			}

			return tag
		}
	}

	panic("fieldPointer must point to a struct field with a \"json\" tag")
}
