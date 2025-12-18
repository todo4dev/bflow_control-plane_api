package cqrs

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/todo4dev/bflow/control-plane/api/core/di"
)

var (
	commandExecutorsMutex  sync.RWMutex
	commandExecutorsByType = map[reflect.Type]FnCommandExecutor{}
)

type Command any

type ICommandHandler[TCommand any, TResult any] interface {
	Handle(ctx context.Context, command TCommand) (TResult, error)
}

type FnCommandExecutor func(ctx context.Context, command Command) (any, error)

func RegisterCommandHandler[TCommand any, TResult any, THandler ICommandHandler[TCommand, TResult]](factoryFunction any) {
	di.RegisterAs[THandler](factoryFunction)
	commandType := reflect.TypeOf((*TCommand)(nil)).Elem()

	commandExecutorsMutex.Lock()
	defer commandExecutorsMutex.Unlock()

	if _, exists := commandExecutorsByType[commandType]; exists {
		panic(fmt.Sprintf("cqrs: CommandHandler already registered for type %v", commandType))
	}

	executor := func(ctx context.Context, command Command) (any, error) {
		typedCommand, ok := command.(TCommand)
		if !ok {
			return nil, fmt.Errorf("cqrs: expected command of type %v, got %T", commandType, command)
		}

		handler := di.Resolve[THandler]()

		result, err := handler.Handle(ctx, typedCommand)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	commandExecutorsByType[commandType] = executor
}

func ExecuteCommand[T any](ctx context.Context, command Command) (*T, error) {
	if command == nil {
		return nil, fmt.Errorf("cqrs: nil command")
	}

	commandType := reflect.TypeOf(command)

	commandExecutorsMutex.RLock()
	executor, exists := commandExecutorsByType[commandType]
	commandExecutorsMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("cqrs: no CommandHandler registered for command type %v", commandType)
	}

	anyResult, err := executor(ctx, command)
	if err != nil {
		return nil, err
	}

	if anyResult == nil {
		return nil, nil
	}

	typedResult, ok := anyResult.(*T)
	if !ok {
		return nil, fmt.Errorf("cqrs: command handler returned unexpected type %T, expected *%s", anyResult, reflect.TypeOf((*T)(nil)).Elem())
	}

	return typedResult, nil
}

func MustExecuteCommand[T any](ctx context.Context, command Command) *T {
	result, err := ExecuteCommand[T](ctx, command)
	if err != nil {
		panic(err)
	}
	return result
}
