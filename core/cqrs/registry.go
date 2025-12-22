package cqrs

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"src/core/di"
)

type handlerRegistry struct {
	mutex     sync.RWMutex
	executors map[reflect.Type]func(ctx context.Context, message any) (any, error)
	kindName  string
}

func newHandlerRegistry(kindName string) *handlerRegistry {
	return &handlerRegistry{
		executors: map[reflect.Type]func(context.Context, any) (any, error){},
		kindName:  kindName,
	}
}

func register[TMessage any, TResult any, THandler interface {
	Handle(ctx context.Context, message TMessage) (TResult, error)
}](registry *handlerRegistry, factoryFunction any) {
	di.RegisterAs[THandler](factoryFunction)

	messageKey := normalizeType(reflect.TypeFor[TMessage]())

	registry.mutex.Lock()
	defer registry.mutex.Unlock()

	if _, exists := registry.executors[messageKey]; exists {
		panic(fmt.Sprintf("cqrs: %s handler already registered for type %v (normalized)", registry.kindName, messageKey))
	}

	registry.executors[messageKey] = func(ctx context.Context, message any) (any, error) {
		typedMessage, err := coerce[TMessage](message, registry.kindName)
		if err != nil {
			return nil, err
		}

		handler := di.Resolve[THandler]()
		return handler.Handle(ctx, typedMessage)
	}
}

func execute[TResult any](registry *handlerRegistry, ctx context.Context, message any) (*TResult, error) {
	messageKey, err := normalizedTypeKeyOfValue(message, registry.kindName)
	if err != nil {
		return nil, err
	}

	registry.mutex.RLock()
	executor, exists := registry.executors[messageKey]
	registry.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("cqrs: no %s handler registered for type %v (normalized)", registry.kindName, messageKey)
	}

	anyResult, err := executor(ctx, message)
	if err != nil {
		return nil, err
	}

	if anyResult == nil {
		return nil, nil
	}

	typedResult, err := coerce[*TResult](anyResult, "result")
	if err != nil {
		return nil, fmt.Errorf("cqrs: %s handler returned unexpected result: %w", registry.kindName, err)
	}

	return typedResult, nil
}
