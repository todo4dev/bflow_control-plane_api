package cqrs

import (
	"fmt"
	"reflect"
)

func normalizeType(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.Pointer {
		return t.Elem()
	}
	return t
}

func normalizedTypeKeyOfValue(value any, valueName string) (reflect.Type, error) {
	if value == nil {
		return nil, fmt.Errorf("cqrs: nil %s", valueName)
	}

	valueType := reflect.TypeOf(value)

	if valueType.Kind() == reflect.Pointer && reflect.ValueOf(value).IsNil() {
		return nil, fmt.Errorf("cqrs: nil %s pointer", valueName)
	}

	return normalizeType(valueType), nil
}

func coerce[TExpected any](value any, valueName string) (TExpected, error) {
	var zero TExpected

	if value == nil {
		return zero, fmt.Errorf("cqrs: nil %s", valueName)
	}

	expectedType := reflect.TypeFor[TExpected]()
	gotValue := reflect.ValueOf(value)
	gotType := gotValue.Type()

	if gotType.AssignableTo(expectedType) {
		return gotValue.Interface().(TExpected), nil
	}
	if gotType.ConvertibleTo(expectedType) {
		converted := gotValue.Convert(expectedType)
		return converted.Interface().(TExpected), nil
	}

	if expectedType.Kind() == reflect.Pointer && gotType.Kind() != reflect.Pointer {
		expectedElem := expectedType.Elem()

		var source reflect.Value
		switch {
		case gotType.AssignableTo(expectedElem):
			source = gotValue
		case gotType.ConvertibleTo(expectedElem):
			source = gotValue.Convert(expectedElem)
		default:
			source = reflect.Value{}
		}

		if source.IsValid() {
			pointer := reflect.New(expectedElem)
			pointer.Elem().Set(source)
			return pointer.Interface().(TExpected), nil
		}
	}

	if expectedType.Kind() != reflect.Pointer && gotType.Kind() == reflect.Pointer {
		if gotValue.IsNil() {
			return zero, fmt.Errorf("cqrs: nil %s pointer", valueName)
		}

		elem := gotValue.Elem()
		elemType := elem.Type()

		if elemType.AssignableTo(expectedType) {
			return elem.Interface().(TExpected), nil
		}
		if elemType.ConvertibleTo(expectedType) {
			converted := elem.Convert(expectedType)
			return converted.Interface().(TExpected), nil
		}
	}

	return zero, fmt.Errorf("cqrs: expected %s compatible with %v, got %T", valueName, expectedType, value)
}
