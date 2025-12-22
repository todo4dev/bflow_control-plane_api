package cqrs

import "context"

type Query any

type IQueryHandler[TQuery any, TResult any] interface {
	Handle(ctx context.Context, query TQuery) (TResult, error)
}

var queryRegistry = newHandlerRegistry("query")

func RegisterQueryHandler[TQuery any, TResult any, THandler IQueryHandler[TQuery, TResult]](factoryFunction any) {
	register[TQuery, TResult, THandler](queryRegistry, factoryFunction)
}

func ExecuteQuery[T any](ctx context.Context, query Query) (*T, error) {
	return execute[T](queryRegistry, ctx, query)
}

func MustExecuteQuery[T any](ctx context.Context, query Query) *T {
	result, err := ExecuteQuery[T](ctx, query)
	if err != nil {
		panic(err)
	}
	return result
}
