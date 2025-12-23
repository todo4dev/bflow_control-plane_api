package meta

import (
	"fmt"
	"reflect"
	"sync"
)

type ObjectMetadata struct {
	Description string
	Throws      []ThrowsMetadata
	Fields      map[string]*FieldMetadata
	Example     any
	Type        reflect.Type
}

type FieldMetadata struct {
	Description string
	Example     any
	Type        reflect.Type
	Nullable    bool
}

type ThrowsMetadata struct {
	ErrorType   reflect.Type
	Description string
}

var (
	registryMutex  sync.RWMutex
	structRegistry = map[reflect.Type]*ObjectMetadata{}
)

func getOrCreateObjectMetadata(structType reflect.Type) *ObjectMetadata {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	if metadata, exists := structRegistry[structType]; exists {
		return metadata
	}

	metadata := &ObjectMetadata{
		Fields: map[string]*FieldMetadata{},
		Type:   structType,
	}
	structRegistry[structType] = metadata

	return metadata
}

func GetObjectMetadataAs[T any]() *ObjectMetadata {
	var zero T
	structType := reflect.TypeOf(zero)
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}

	registryMutex.RLock()
	defer registryMutex.RUnlock()

	return structRegistry[structType]
}

func GetObjectMetadataOf(structInstance any) *ObjectMetadata {
	if structInstance == nil {
		return nil
	}

	structType := reflect.TypeOf(structInstance)
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}

	if structType.Kind() != reflect.Struct {
		return nil
	}

	registryMutex.RLock()
	defer registryMutex.RUnlock()

	return structRegistry[structType]
}

func GetObjectMetadataByType(structType reflect.Type) *ObjectMetadata {
	if structType == nil {
		return nil
	}

	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}

	registryMutex.RLock()
	defer registryMutex.RUnlock()

	return structRegistry[structType]
}

type ObjectOption interface {
	applyToObject(structPointer any, metadata *ObjectMetadata)
}

type FieldOption interface {
	applyToField(fieldMetadata *FieldMetadata)
}

func Describe(target any, options ...ObjectOption) {
	if target == nil {
		panic("doc: target is nil in Describe")
	}

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Pointer || targetValue.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("doc: Describe target must be pointer to struct, got %T", target))
	}

	structType := targetValue.Elem().Type()
	structMetadata := getOrCreateObjectMetadata(structType)

	for _, option := range options {
		if option == nil {
			continue
		}
		option.applyToObject(target, structMetadata)
	}
}

type descriptionDecorator struct {
	Text string
}

func Description(text string) descriptionDecorator {
	return descriptionDecorator{Text: text}
}

func (option descriptionDecorator) applyToObject(_ any, metadata *ObjectMetadata) {
	metadata.Description = option.Text
}

func (option descriptionDecorator) applyToField(fieldMetadata *FieldMetadata) {
	fieldMetadata.Description = option.Text
}

type exampleDecorator struct {
	Value any
}

func Example[T any](value T) exampleDecorator {
	return exampleDecorator{Value: value}
}

func (option exampleDecorator) applyToObject(_ any, metadata *ObjectMetadata) {
	metadata.Example = option.Value
}

func (option exampleDecorator) applyToField(fieldMetadata *FieldMetadata) {
	fieldMetadata.Example = option.Value
}

type throwsDecorator struct {
	ErrorType   reflect.Type
	Description string
}

func Throws[T any](description string) throwsDecorator {
	var zeroPointer *T
	errorType := reflect.TypeOf(zeroPointer).Elem()

	return throwsDecorator{
		ErrorType:   errorType,
		Description: description,
	}
}

func (option throwsDecorator) applyToObject(_ any, metadata *ObjectMetadata) {
	metadata.Throws = append(metadata.Throws, ThrowsMetadata(option))
}

type fieldDecorator struct {
	FieldPointer any
	Options      []FieldOption
}

func Field(fieldPointer any, options ...FieldOption) fieldDecorator {
	if fieldPointer == nil {
		panic("doc: fieldPointer is nil in Field")
	}

	return fieldDecorator{
		FieldPointer: fieldPointer,
		Options:      options,
	}
}

func (decorator fieldDecorator) applyToObject(structPointer any, structMetadata *ObjectMetadata) {
	fieldName := resolveFieldName(structPointer, decorator.FieldPointer)
	if fieldName == "" {
		panic("doc: could not resolve field name for provided fieldPointer")
	}

	fieldMetadata, exists := structMetadata.Fields[fieldName]
	if !exists {
		fieldMetadata = &FieldMetadata{}
		structMetadata.Fields[fieldName] = fieldMetadata
	}

	fieldValue := reflect.ValueOf(decorator.FieldPointer)
	if fieldValue.Kind() == reflect.Pointer {
		fieldMetadata.Type = fieldValue.Elem().Type()
		fieldMetadata.Nullable = true
	} else {
		fieldMetadata.Type = fieldValue.Type()
		switch fieldMetadata.Type.Kind() {
		case reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
			fieldMetadata.Nullable = true
		default:
			fieldMetadata.Nullable = false
		}
	}

	for _, option := range decorator.Options {
		if option == nil {
			continue
		}
		option.applyToField(fieldMetadata)
	}
}

func resolveFieldName(structPointer any, fieldPointer any) string {
	if structPointer == nil || fieldPointer == nil {
		return ""
	}

	structValue := reflect.ValueOf(structPointer)
	if structValue.Kind() != reflect.Pointer || structValue.Elem().Kind() != reflect.Struct {
		return ""
	}

	fieldValue := reflect.ValueOf(fieldPointer)
	if fieldValue.Kind() != reflect.Pointer {
		return ""
	}

	targetAddr := fieldValue.Pointer()
	return resolveFieldNameRecursive(structValue.Elem(), targetAddr)
}

func resolveFieldNameRecursive(structValue reflect.Value, targetAddr uintptr) string {
	if structValue.Kind() != reflect.Struct {
		return ""
	}

	structType := structValue.Type()
	for i := 0; i < structValue.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		if fieldValue.CanAddr() && fieldValue.Addr().Pointer() == targetAddr {
			return field.Name
		}
		switch fieldValue.Kind() {
		case reflect.Struct:
			if subName := resolveFieldNameRecursive(fieldValue, targetAddr); subName != "" {
				if field.Anonymous {
					return subName
				}
				return field.Name + "." + subName
			}

		case reflect.Pointer:
			if !fieldValue.IsNil() && fieldValue.Elem().Kind() == reflect.Struct {
				if subName := resolveFieldNameRecursive(fieldValue.Elem(), targetAddr); subName != "" {
					if field.Anonymous {
						return subName
					}
					return field.Name + "." + subName
				}
			}
		}
	}

	return ""
}
