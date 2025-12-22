package builder

import "slices"

type Field struct {
	Select []string `json:"select,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

type FieldBuilder[TEntity any] struct {
	builderBase
	Field Field
}

type FieldFn[TEntity any] func(entity *TEntity, fieldBuilder *FieldBuilder[TEntity])

func NewField[TEntity any]() *FieldBuilder[TEntity] {
	return &FieldBuilder[TEntity]{builderBase: builderBase{}, Field: Field{}}
}

func (b *FieldBuilder[TEntity]) removeFieldName(list []string, fieldName string) []string {
	for i, v := range list {
		if v == fieldName {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func (b *FieldBuilder[TEntity]) containsFieldName(list []string, fieldName string) bool {
	return slices.Contains(list, fieldName)
}

func (b *FieldBuilder[TEntity]) Select(fieldPointer any) {
	fieldName := b.builderBase.fieldPointerJSONTag(fieldPointer)
	b.Field.Remove = b.removeFieldName(b.Field.Remove, fieldName)
	if !b.containsFieldName(b.Field.Select, fieldName) {
		b.Field.Select = append(b.Field.Select, fieldName)
	}
}

func (b *FieldBuilder[TEntity]) Remove(fieldPointer any) {
	fieldName := b.builderBase.fieldPointerJSONTag(fieldPointer)
	b.Field.Select = b.removeFieldName(b.Field.Select, fieldName)
	if !b.containsFieldName(b.Field.Remove, fieldName) {
		b.Field.Remove = append(b.Field.Remove, fieldName)
	}
}
