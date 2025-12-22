package builder

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Query[TEntity any] struct {
	TextCond   *string          `json:"text,omitempty"`
	WhereCond  *WherePointerMap `json:"where,omitempty"`
	FieldCond  *Field           `json:"fields,omitempty"`
	SortCond   *SortPointerMap  `json:"sort,omitempty"`
	LimitCond  *int64           `json:"limit,omitempty"`
	OffsetCond *int64           `json:"offset,omitempty"`
}

var _ json.Marshaler = (*Query[any])(nil)
var _ json.Unmarshaler = (*Query[any])(nil)

func NewQuery[TEntity any]() *Query[TEntity] {
	return &Query[TEntity]{}
}

func (q *Query[TEntity]) normalizeWhereTypes() error {
	if q.WhereCond == nil {
		return nil
	}

	for fieldName, ops := range *q.WhereCond {
		fieldType, ok := q.fieldTypeByJSONTag(fieldName)
		if !ok || fieldType == nil {
			continue
		}

		for op, rawVal := range ops {
			switch op {
			case WhereEnum_Empty, WhereEnum_NotEmpty:
				continue
			case WhereEnum_In, WhereEnum_NotIn:
				normalized, err := q.normalizeWhereSliceValue(rawVal, fieldType)
				if err != nil {
					return fmt.Errorf("field %s operator %s: %w", fieldName, op, err)
				}
				ops[op] = normalized
			default:
				normalized, err := q.normalizeWhereScalarValue(rawVal, fieldType)
				if err != nil {
					return fmt.Errorf("field %s operator %s: %w", fieldName, op, err)
				}
				ops[op] = normalized
			}
		}
	}

	return nil
}

func (q *Query[TEntity]) fieldTypeByJSONTag(fieldName string) (reflect.Type, bool) {
	var entity TEntity
	t := reflect.TypeOf(entity)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, false
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		tag := f.Tag.Get("json")
		if tag == "" {
			continue
		}

		if idx := strings.Index(tag, ","); idx != -1 {
			tag = tag[:idx]
		}

		if tag == fieldName {
			ft := f.Type
			if ft.Kind() == reflect.Pointer {
				ft = ft.Elem()
			}
			return ft, true
		}
	}

	return nil, false
}

func (q *Query[TEntity]) normalizeWhereScalarValue(value any, targetType reflect.Type) (any, error) {
	if value == nil || targetType == nil {
		return value, nil
	}

	switch targetType.Kind() {
	case reflect.Bool:
		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return nil, err
			}
			return b, nil
		default:
			return nil, fmt.Errorf("cannot convert %T to bool", value)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var n int64

		switch v := value.(type) {
		case int:
			n = int64(v)
		case int8:
			n = int64(v)
		case int16:
			n = int64(v)
		case int32:
			n = int64(v)
		case int64:
			n = v
		case float32:
			n = int64(v)
		case float64:
			n = int64(v)
		case string:
			parsed, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, err
			}
			n = parsed
		default:
			return nil, fmt.Errorf("cannot convert %T to int", value)
		}

		out := reflect.New(targetType).Elem()
		out.SetInt(n)
		return out.Interface(), nil

	case reflect.Float32, reflect.Float64:
		var f64 float64

		switch v := value.(type) {
		case float32:
			f64 = float64(v)
		case float64:
			f64 = v
		case int:
			f64 = float64(v)
		case int64:
			f64 = float64(v)
		case string:
			parsed, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			f64 = parsed
		default:
			return nil, fmt.Errorf("cannot convert %T to float", value)
		}

		out := reflect.New(targetType).Elem()
		out.SetFloat(f64)
		return out.Interface(), nil

	case reflect.String:
		if v, ok := value.(string); ok {
			return v, nil
		}
		return fmt.Sprint(value), nil

	case reflect.Struct:
		if targetType == reflect.TypeOf(time.Time{}) {
			if v, ok := value.(string); ok {
				t, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return nil, err
				}
				return t, nil
			}
			return nil, fmt.Errorf("cannot convert %T to time.Time", value)
		}
	}

	return value, nil
}

func (q *Query[TEntity]) normalizeWhereSliceValue(value any, elemType reflect.Type) (any, error) {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, fmt.Errorf("expected slice/array for IN/NIN, got %T", value)
	}

	out := reflect.MakeSlice(reflect.SliceOf(elemType), 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		converted, err := q.normalizeWhereScalarValue(v.Index(i).Interface(), elemType)
		if err != nil {
			return nil, err
		}
		out = reflect.Append(out, reflect.ValueOf(converted))
	}

	return out.Interface(), nil
}

func (q *Query[TEntity]) Text(text string) *Query[TEntity] {
	if text == "" {
		q.TextCond = nil
	} else {
		q.TextCond = &text
	}
	return q
}

func (q *Query[TEntity]) Where(fn WhereFn[TEntity]) *Query[TEntity] {
	if fn == nil {
		q.WhereCond = nil
		return q
	}

	var entity TEntity
	whereBuilder := NewWhere[TEntity]()
	fn(&entity, whereBuilder)

	if len(whereBuilder.PointerMap) == 0 {
		q.WhereCond = nil
	} else {
		q.WhereCond = &whereBuilder.PointerMap
	}

	return q
}

func (q *Query[TEntity]) Field(fn FieldFn[TEntity]) *Query[TEntity] {
	if fn == nil {
		q.FieldCond = nil
		return q
	}

	var entity TEntity
	fieldBuilder := NewField[TEntity]()
	fn(&entity, fieldBuilder)

	if len(fieldBuilder.Field.Select) == 0 && len(fieldBuilder.Field.Remove) == 0 {
		q.FieldCond = nil
	} else {
		q.FieldCond = &fieldBuilder.Field
	}

	return q
}

func (q *Query[TEntity]) Sort(fn SortFn[TEntity]) *Query[TEntity] {
	if fn == nil {
		q.SortCond = nil
		return q
	}

	var entity TEntity
	sortBuilder := NewSort[TEntity]()
	fn(&entity, sortBuilder)

	if len(sortBuilder.PointerMap) == 0 {
		q.SortCond = nil
	} else {
		q.SortCond = &sortBuilder.PointerMap
	}

	return q
}

func (q *Query[TEntity]) Offset(offset int64) *Query[TEntity] {
	q.OffsetCond = &offset
	return q
}

func (q *Query[TEntity]) Limit(limit int64) *Query[TEntity] {
	q.LimitCond = &limit
	return q
}

func (q *Query[TEntity]) ToJSON() *Query[json.RawMessage] {
	if q == nil {
		return nil
	}

	return &Query[json.RawMessage]{
		TextCond:   q.TextCond,
		WhereCond:  q.WhereCond,
		FieldCond:  q.FieldCond,
		SortCond:   q.SortCond,
		LimitCond:  q.LimitCond,
		OffsetCond: q.OffsetCond,
	}
}

func (q *Query[TEntity]) UnmarshalJSON(data []byte) error {
	type Alias Query[TEntity]
	var aux Alias

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*q = Query[TEntity](aux)
	return q.normalizeWhereTypes()
}

func (q *Query[TEntity]) MarshalJSON() ([]byte, error) {
	type Alias Query[TEntity]
	return json.Marshal((*Alias)(q.ToJSON()))
}
