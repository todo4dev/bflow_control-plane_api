package builder

import "encoding/json"

type UpdateFn[TEntity any] func(entity *TEntity, updateBuilder *UpdateBuilder[TEntity])

type UpdateBuilder[TEntity any] struct {
	builderBase
	Changes map[string]any
}

func NewUpdate[TEntity any]() *UpdateBuilder[TEntity] {
	return &UpdateBuilder[TEntity]{
		builderBase: builderBase{},
		Changes:     make(map[string]any),
	}
}

func (b *UpdateBuilder[TEntity]) Set(fieldPointer any, value any) *UpdateBuilder[TEntity] {
	fieldName := b.builderBase.fieldPointerJSONTag(fieldPointer)
	b.Changes[fieldName] = value
	return b
}

func (b *UpdateBuilder[TEntity]) ToJSON() *UpdateBuilder[json.RawMessage] {
	if b == nil {
		return nil
	}

	return &UpdateBuilder[json.RawMessage]{
		builderBase: b.builderBase,
		Changes:     b.Changes,
	}
}
