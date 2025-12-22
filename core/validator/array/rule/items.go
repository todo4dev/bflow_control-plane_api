package rule

import (
	"fmt"
	"reflect"

	"src/core/validator/_base"
	"src/core/validator/array/util"
)

type ItemsRule struct{}

var _ _base.IRule = ItemsRule{}

func (r ItemsRule) Apply(valuePointer reflect.Value) []_base.ValidationError {
	slice, present, typeErr := util.ExtractSlice(valuePointer)
	if typeErr != nil {
		return typeErr
	}
	if !present {
		return nil
	}

	length := slice.Len()
	if length == 0 {
		return nil
	}

	var allErrors []_base.ValidationError

	for index := 0; index < length; index++ {
		element := slice.Index(index)

		var elementPointer reflect.Value

		if element.Kind() == reflect.Pointer {
			if element.IsNil() {
				continue
			}
			elementPointer = element
		} else if element.CanAddr() {
			elementPointer = element.Addr()
		} else {
			continue
		}

		validable, ok := elementPointer.Interface().(_base.IValidable)
		if !ok || validable == nil {
			continue
		}

		err := validable.Validate()
		if err == nil {
			continue
		}

		if validationError, ok := err.(*_base.Error); ok {
			for _, inner := range validationError.Errors {
				allErrors = append(allErrors, _base.ValidationError{
					Code:    "array.items",
					Message: fmt.Sprintf("item %d: %s: %s", index, inner.Code, inner.Message),
				})
			}
		} else {
			allErrors = append(allErrors, _base.ValidationError{
				Code:    "array.items",
				Message: fmt.Sprintf("item %d: %s", index, err.Error()),
			})
		}
	}

	return allErrors
}
