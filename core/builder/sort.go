package builder

type SortEnum int

const (
	SortEnum_Desc SortEnum = -1
	SortEnum_Asc  SortEnum = 1
)

type SortPointerMap map[string]SortEnum

type SortBuilder[TEntity any] struct {
	builderBase
	PointerMap SortPointerMap
}

type SortFn[TEntity any] func(entity *TEntity, sortBuilder *SortBuilder[TEntity])

func NewSort[TEntity any]() *SortBuilder[TEntity] {
	return &SortBuilder[TEntity]{builderBase: builderBase{}, PointerMap: make(SortPointerMap)}
}

func (s *SortBuilder[TEntity]) Desc(fieldPointer any) {
	fieldName := s.builderBase.fieldPointerJSONTag(fieldPointer)
	s.PointerMap[fieldName] = SortEnum_Desc
}

func (s *SortBuilder[TEntity]) Asc(fieldPointer any) {
	fieldName := s.builderBase.fieldPointerJSONTag(fieldPointer)
	s.PointerMap[fieldName] = SortEnum_Asc
}
