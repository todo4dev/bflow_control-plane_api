package cqrs

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/todo4dev/bflow_control-plane_api/core/di"
)

var (
	queryExecutorsMutex  sync.RWMutex
	queryExecutorsByType = map[reflect.Type]FnQueryExecutor{}
)

type Query any

type IQueryHandler[TQuery any, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}

type FnQueryExecutor func(ctx context.Context, query Query) (any, error)

func RegisterQueryHandler[TQuery any, TResult any, THandler IQueryHandler[TQuery, TResult]](factoryFunction any) {
	di.RegisterAs[THandler](factoryFunction)

	queryType := reflect.TypeOf((*TQuery)(nil)).Elem()

	queryExecutorsMutex.Lock()
	defer queryExecutorsMutex.Unlock()

	if _, exists := queryExecutorsByType[queryType]; exists {
		panic(fmt.Sprintf("cqrs: QueryHandler already registered for type %v", queryType))
	}

	executor := func(ctx context.Context, query Query) (any, error) {
		typedQuery, ok := query.(TQuery)
		if !ok {
			return nil, fmt.Errorf("cqrs: expected query of type %v, got %T", queryType, query)
		}

		handler := di.Resolve[THandler]()

		result, err := handler.Handle(ctx, typedQuery)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	queryExecutorsByType[queryType] = executor
}

func ExecuteQuery[T any](ctx context.Context, query Query) (*T, error) {
	if query == nil {
		return nil, fmt.Errorf("cqrs: nil query")
	}

	queryType := reflect.TypeOf(query)

	queryExecutorsMutex.RLock()
	executor, exists := queryExecutorsByType[queryType]
	queryExecutorsMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("cqrs: no QueryHandler registered for query type %v", queryType)
	}

	anyResult, err := executor(ctx, query)
	if err != nil {
		return nil, err
	}

	typedResult, ok := anyResult.(*T)
	if !ok {
		return nil, fmt.Errorf("cqrs: expected result of type %v, got %T", reflect.TypeOf(typedResult), typedResult)
	}
	return typedResult, nil
}

func MustExecuteQuery[T any](ctx context.Context, query Query) *T {
	result, err := ExecuteQuery[T](ctx, query)
	if err != nil {
		panic(err)
	}
	return result
}
