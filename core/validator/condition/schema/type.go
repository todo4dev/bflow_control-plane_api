package schema

import "src/core/validator/_base"

type ConditionMatchMode int

const (
	// ConditionMatchAny: passa se pelo menos 1 schema passar (default).
	ConditionMatchAny ConditionMatchMode = iota

	// ConditionMatchOne: passa se exatamente 1 schema passar.
	ConditionMatchOne

	// ConditionMatchAll: passa se todos os schemas passarem.
	ConditionMatchAll
)

type ConditionSchema struct {
	alternatives []_base.IValidable
	matchMode    ConditionMatchMode
}

func Condition(schemas ..._base.IValidable) *ConditionSchema {
	return &ConditionSchema{
		alternatives: append([]_base.IValidable(nil), schemas...),
		matchMode:    ConditionMatchAny,
	}
}

func (builder *ConditionSchema) Try(schemas ..._base.IValidable) *ConditionSchema {
	if len(schemas) == 0 {
		return builder
	}

	builder.alternatives = append(builder.alternatives, schemas...)
	return builder
}

func (builder *ConditionSchema) Match(mode ConditionMatchMode) *ConditionSchema {
	builder.matchMode = mode
	return builder
}

type alternativeSchema struct {
	predicate       func() bool
	thenSchema      _base.IValidable
	otherwiseSchema _base.IValidable
}

func (s alternativeSchema) Validate() error {
	if s.predicate == nil {
		return nil
	}

	if s.predicate() {
		if s.thenSchema == nil {
			return nil
		}
		return s.thenSchema.Validate()
	}

	if s.otherwiseSchema != nil {
		return s.otherwiseSchema.Validate()
	}

	return nil
}

func (builder *ConditionSchema) Condition(
	predicate func() bool,
	thenSchema _base.IValidable,
	otherwiseSchema _base.IValidable,
) *ConditionSchema {
	builder.alternatives = append(builder.alternatives, alternativeSchema{
		predicate:       predicate,
		thenSchema:      thenSchema,
		otherwiseSchema: otherwiseSchema,
	})
	return builder
}
func (builder *ConditionSchema) Validate() error {
	if len(builder.alternatives) == 0 {
		return nil
	}

	switch builder.matchMode {
	case ConditionMatchAll:
		return builder.validateAll()
	case ConditionMatchOne:
		return builder.validateOne()
	default:
		return builder.validateAny()
	}
}
func (builder *ConditionSchema) validateAny() error {
	var collected []_base.ValidationError

	for _, schema := range builder.alternatives {
		if schema == nil {
			continue
		}

		err := schema.Validate()
		if err == nil {
			return nil
		}

		if validationError, ok := err.(*_base.Error); ok {
			collected = append(collected, validationError.Errors...)
		} else {
			collected = append(collected, _base.ValidationError{
				Code:    "condition.altError",
				Message: err.Error(),
			})
		}
	}

	if len(collected) == 0 {
		collected = append(collected, _base.ValidationError{
			Code:    "condition.noMatch",
			Message: "no condition matched",
		})
	}

	return &_base.Error{Errors: collected}
}
func (builder *ConditionSchema) validateOne() error {
	var collected []_base.ValidationError
	successCount := 0

	for _, schema := range builder.alternatives {
		if schema == nil {
			continue
		}

		err := schema.Validate()
		if err == nil {
			successCount++
			continue
		}

		if validationError, ok := err.(*_base.Error); ok {
			collected = append(collected, validationError.Errors...)
		} else {
			collected = append(collected, _base.ValidationError{
				Code:    "condition.altError",
				Message: err.Error(),
			})
		}
	}

	if successCount == 1 {
		return nil
	}

	if successCount == 0 {
		if len(collected) == 0 {
			collected = append(collected, _base.ValidationError{
				Code:    "condition.noMatch",
				Message: "no condition matched",
			})
		}
		return &_base.Error{Errors: collected}
	}

	return &_base.Error{
		Errors: []_base.ValidationError{
			{
				Code:    "condition.multipleMatches",
				Message: "value matches more than one condition",
			},
		},
	}
}
func (builder *ConditionSchema) validateAll() error {
	var collected []_base.ValidationError

	for _, schema := range builder.alternatives {
		if schema == nil {
			continue
		}

		err := schema.Validate()
		if err == nil {
			continue
		}

		if validationError, ok := err.(*_base.Error); ok {
			collected = append(collected, validationError.Errors...)
		} else {
			collected = append(collected, _base.ValidationError{
				Code:    "condition.altError",
				Message: err.Error(),
			})
		}
	}

	if len(collected) == 0 {
		return nil
	}

	return &_base.Error{Errors: collected}
}
