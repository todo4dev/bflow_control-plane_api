package di

import (
	"fmt"
	"reflect"
	"sync"
)

var (
	registryMutex    sync.RWMutex
	providerRegistry = map[reflect.Type][]*Provider{}
)

type Provider struct {
	FactoryFunction reflect.Value
	OutputType      reflect.Type
	IsSingleton     bool
	CachedInstance  reflect.Value
}

func typeOf[T any]() reflect.Type {
	var zeroValue *T
	return reflect.TypeOf(zeroValue).Elem()
}

func registerProvider(factoryFunction any, isSingleton bool, asType reflect.Type) {
	if factoryFunction == nil {
		panic("di: nil factory function provided")
	}

	factoryValue := reflect.ValueOf(factoryFunction)
	factoryType := factoryValue.Type()

	if factoryType.Kind() != reflect.Func {
		panic("di: factory must be a function")
	}
	if factoryType.NumOut() != 1 {
		panic("di: factory function must return exactly one value")
	}

	outputType := factoryType.Out(0)

	if asType != nil {
		if !outputType.AssignableTo(asType) {
			panic(fmt.Sprintf(
				"di: factory return type %v is not assignable to %v",
				outputType,
				asType,
			))
		}
		outputType = asType
	}

	providerInstance := &Provider{
		FactoryFunction: factoryValue,
		OutputType:      outputType,
		IsSingleton:     isSingleton,
	}

	registryMutex.Lock()
	providerRegistry[outputType] = append(providerRegistry[outputType], providerInstance)
	registryMutex.Unlock()
}

func resolveByType(targetType reflect.Type) reflect.Value {
	registryMutex.RLock()
	providers := providerRegistry[targetType]
	registryMutex.RUnlock()

	if len(providers) == 0 {
		panic(fmt.Sprintf("di: no provider registered for type %v", targetType))
	}

	return buildInstance(providers[0])
}

func buildInstance(providerInstance *Provider) reflect.Value {
	if providerInstance.IsSingleton {
		registryMutex.Lock()
		if providerInstance.CachedInstance.IsValid() {
			instance := providerInstance.CachedInstance
			registryMutex.Unlock()
			return instance
		}

		instance := callFactoryWithDependencies(providerInstance)
		providerInstance.CachedInstance = instance
		registryMutex.Unlock()

		return instance
	}

	return callFactoryWithDependencies(providerInstance)
}

func callFactoryWithDependencies(providerInstance *Provider) reflect.Value {
	factoryType := providerInstance.FactoryFunction.Type()
	numberOfInputs := factoryType.NumIn()

	arguments := make([]reflect.Value, numberOfInputs)

	for index := 0; index < numberOfInputs; index++ {
		dependencyType := factoryType.In(index)
		arguments[index] = resolveByType(dependencyType)
	}

	outputValues := providerInstance.FactoryFunction.Call(arguments)
	return outputValues[0]
}

func Register(factoryFunction any) {
	registerProvider(factoryFunction, false, nil)
}

func RegisterAs[TType any](factoryFunction any) {
	targetType := typeOf[TType]()
	registerProvider(factoryFunction, false, targetType)
}

func Singleton(factoryFunction any) {
	registerProvider(factoryFunction, true, nil)
}

func SingletonAs[TType any](factoryFunction any) {
	targetType := typeOf[TType]()
	registerProvider(factoryFunction, true, targetType)
}

func Resolve[T any]() T {
	targetType := typeOf[T]()
	value := resolveByType(targetType)
	return value.Interface().(T)
}

func ResolveAll[T any]() []T {
	targetType := typeOf[T]()

	registryMutex.RLock()
	providers := providerRegistry[targetType]
	registryMutex.RUnlock()

	if len(providers) == 0 {
		return nil
	}

	results := make([]T, 0, len(providers))

	for _, providerInstance := range providers {
		value := buildInstance(providerInstance)
		results = append(results, value.Interface().(T))
	}

	return results
}
