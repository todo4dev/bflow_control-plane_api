package builder

import (
	"encoding/json"
	"reflect"
)

type WhereEnum string

const (
	WhereEnum_Equal           WhereEnum = "eq"
	WhereEnum_NotEqual        WhereEnum = "neq"
	WhereEnum_Like            WhereEnum = "like"
	WhereEnum_NotLike         WhereEnum = "nlike"
	WhereEnum_Empty           WhereEnum = "empty"
	WhereEnum_NotEmpty        WhereEnum = "nempty"
	WhereEnum_In              WhereEnum = "in"
	WhereEnum_NotIn           WhereEnum = "nin"
	WhereEnum_GreaterThan     WhereEnum = "gt"
	WhereEnum_NotGreaterThan  WhereEnum = "ngt"
	WhereEnum_GreaterEqual    WhereEnum = "gte"
	WhereEnum_NotGreaterEqual WhereEnum = "ngte"
	WhereEnum_LowerThan       WhereEnum = "lt"
	WhereEnum_NotLowerThan    WhereEnum = "nlt"
	WhereEnum_LowerEqual      WhereEnum = "lte"
	WhereEnum_NotLowerEqual   WhereEnum = "nlte"
)

type WherePointerMap map[string]map[WhereEnum]any

type WhereBuilder[TEntity any] struct {
	builderBase
	PointerMap WherePointerMap
}

type WhereFn[TEntity any] func(entity *TEntity, whereBuilder *WhereBuilder[TEntity])

func NewWhere[TEntity any]() *WhereBuilder[TEntity] {
	return &WhereBuilder[TEntity]{builderBase: builderBase{}, PointerMap: make(WherePointerMap)}
}

func (b *WhereBuilder[TEntity]) addClause(fieldPointer any, operator WhereEnum, value any) {
	fieldName := b.builderBase.fieldPointerJSONTag(fieldPointer)
	if _, exists := b.PointerMap[fieldName]; !exists {
		b.PointerMap[fieldName] = make(map[WhereEnum]any)
	}
	b.PointerMap[fieldName][operator] = value
}

func (b *WhereBuilder[TEntity]) Equal(fieldPointer any, value any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_Equal, value)
	return b
}

func (b *WhereBuilder[TEntity]) NotEqual(fieldPointer any, value any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_NotEqual, value)
	return b
}

func (b *WhereBuilder[TEntity]) Like(fieldPointer any, value string) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_Like, value)
	return b
}

func (b *WhereBuilder[TEntity]) NotLike(fieldPointer any, value string) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_NotLike, value)
	return b
}

func (b *WhereBuilder[TEntity]) Empty(fieldPointer any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_Empty, true)
	return b
}

func (b *WhereBuilder[TEntity]) NotEmpty(fieldPointer any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_NotEmpty, true)
	return b
}

func (b *WhereBuilder[TEntity]) In(fieldPointer any, values any) *WhereBuilder[TEntity] {
	if reflect.TypeOf(values).Kind() != reflect.Slice && reflect.TypeOf(values).Kind() != reflect.Array {
		panic("'values' must be a slice or array for IN operator")
	}
	b.addClause(fieldPointer, WhereEnum_In, values)
	return b
}

func (b *WhereBuilder[TEntity]) NotIn(fieldPointer any, values any) *WhereBuilder[TEntity] {
	if reflect.TypeOf(values).Kind() != reflect.Slice && reflect.TypeOf(values).Kind() != reflect.Array {
		panic("'values' must be a slice or array for NOT IN operator")
	}
	b.addClause(fieldPointer, WhereEnum_NotIn, values)
	return b
}

func (b *WhereBuilder[TEntity]) GreaterThan(fieldPointer any, value any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_GreaterThan, value)
	return b
}

func (b *WhereBuilder[TEntity]) NotGreaterThan(fieldPointer any, value any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_NotGreaterThan, value)
	return b
}

func (b *WhereBuilder[TEntity]) GreaterEqual(fieldPointer any, value any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_GreaterEqual, value)
	return b
}

func (b *WhereBuilder[TEntity]) NotGreaterEqual(fieldPointer any, value any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_NotGreaterEqual, value)
	return b
}

func (b *WhereBuilder[TEntity]) LowerThan(fieldPointer any, value any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_LowerThan, value)
	return b
}

func (b *WhereBuilder[TEntity]) NotLowerThan(fieldPointer any, value any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_NotLowerThan, value)
	return b
}

func (b *WhereBuilder[TEntity]) LowerEqual(fieldPointer any, value any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_LowerEqual, value)
	return b
}

func (b *WhereBuilder[TEntity]) NotLowerEqual(fieldPointer any, value any) *WhereBuilder[TEntity] {
	b.addClause(fieldPointer, WhereEnum_NotLowerEqual, value)
	return b
}

func (b *WhereBuilder[TEntity]) ToJSON() *WhereBuilder[json.RawMessage] {
	if b == nil {
		return nil
	}

	return &WhereBuilder[json.RawMessage]{
		builderBase: b.builderBase,
		PointerMap:  b.PointerMap,
	}
}
