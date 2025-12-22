package cqrs

import "context"

type Command any

type ICommandHandler[TCommand any, TResult any] interface {
	Handle(ctx context.Context, command TCommand) (TResult, error)
}

var commandRegistry = newHandlerRegistry("command")

func RegisterCommandHandler[TCommand any, TResult any, THandler ICommandHandler[TCommand, TResult]](factoryFunction any) {
	register[TCommand, TResult, THandler](commandRegistry, factoryFunction)
}

func ExecuteCommand[T any](ctx context.Context, command Command) (*T, error) {
	return execute[T](commandRegistry, ctx, command)
}

func MustExecuteCommand[T any](ctx context.Context, command Command) *T {
	result, err := ExecuteCommand[T](ctx, command)
	if err != nil {
		panic(err)
	}
	return result
}
