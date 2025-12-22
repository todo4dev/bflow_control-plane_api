package core

import (
	"net/http"
	"path"
	"reflect"
	"slices"
	"strings"

	"src/core/di"
	"src/core/doc"
	"src/domain/exception"
	"src/presentation/api/rest/oas"
)

type RouteEnum string

const (
	Route_Get     RouteEnum = http.MethodGet
	Route_Post    RouteEnum = http.MethodPost
	Route_Put     RouteEnum = http.MethodPut
	Route_Patch   RouteEnum = http.MethodPatch
	Route_Delete  RouteEnum = http.MethodDelete
	Route_Options RouteEnum = http.MethodOptions
	Route_Head    RouteEnum = http.MethodHead
)

type HandlerFN func(ctx HttpContext) error

type GuardFN func(ctx HttpContext) error

type InterceptorFN func(ctx HttpContext, next HandlerFN) error

type Route struct {
	Method      RouteEnum
	Path        string
	Handler     HandlerFN
	OperationID string

	Guards       []GuardFN
	Interceptors []InterceptorFN
}

type Router struct {
	BasePath string
	Routes   []Route
}

func NewRouter() Router {
	return Router{BasePath: "/"}
}

func (r Router) PrefixPath(basePath string) Router {
	if basePath == "" {
		r.BasePath = "/"
		return r
	}

	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}

	r.BasePath = path.Clean(basePath)
	if r.BasePath == "." {
		r.BasePath = "/"
	}

	return r
}

func (r Router) Push(routeBuilder *RouteBuilder) Router {
	if routeBuilder == nil {
		return r
	}

	originalPath := routeBuilder.Route.Path

	prefixedPath := path.Clean(r.BasePath + originalPath)
	if !strings.HasPrefix(prefixedPath, "/") {
		prefixedPath = "/" + prefixedPath
	}
	if prefixedPath == "//" {
		prefixedPath = "/"
	}

	// move o path no OpenAPI de originalPath -> prefixedPath
	if prefixedPath != originalPath {
		openapi := di.Resolve[oas.OpenAPI]()

		if openapi.Paths != nil {
			if pathItem, ok := openapi.Paths[originalPath]; ok {
				// se já existir algo em prefixedPath, sobrescreve (é o mesmo endpoint lógico)
				openapi.Paths[prefixedPath] = pathItem
				delete(openapi.Paths, originalPath)
			}
		}
	}

	// atualiza o path da rota (usado pelo server para registrar no Fiber)
	routeBuilder.Route.Path = prefixedPath

	r.Routes = append(r.Routes, routeBuilder.Route)
	return r
}

type IRestController interface {
	Router() Router
}

type RouteBuilder struct {
	Route          Route
	buildOperation *oas.BuildOperation
}

// NewRoute cria um builder que orquestra tanto o Route quanto o OpenAPI.
func NewRoute(optionalPath ...string) *RouteBuilder {
	return &RouteBuilder{}
}
func (b *RouteBuilder) method(method RouteEnum, path string) *RouteBuilder {
	b.Route.Method = method
	b.Route.Path = path

	openapi := di.Resolve[oas.OpenAPI]()

	switch method {
	case Route_Get:
		openapi.Path(path).Get(func(o *oas.BuildOperation) { b.buildOperation = o })
	case Route_Post:
		openapi.Path(path).Post(func(o *oas.BuildOperation) { b.buildOperation = o })
	case Route_Put:
		openapi.Path(path).Put(func(o *oas.BuildOperation) { b.buildOperation = o })
	case Route_Delete:
		openapi.Path(path).Delete(func(o *oas.BuildOperation) { b.buildOperation = o })
	case Route_Patch:
		openapi.Path(path).Patch(func(o *oas.BuildOperation) { b.buildOperation = o })
	case Route_Options:
		openapi.Path(path).Options(func(o *oas.BuildOperation) { b.buildOperation = o })
	case Route_Head:
		openapi.Path(path).Head(func(o *oas.BuildOperation) { b.buildOperation = o })
	default:
		panic("method not supported")
	}

	return b
}

func (b *RouteBuilder) Get(path string) *RouteBuilder {
	return b.method(Route_Get, path)
}

func (b *RouteBuilder) Post(path string) *RouteBuilder {
	return b.method(Route_Post, path)
}

func (b *RouteBuilder) Put(path string) *RouteBuilder {
	return b.method(Route_Put, path)
}

func (b *RouteBuilder) Delete(path string) *RouteBuilder {
	return b.method(Route_Delete, path)
}

func (b *RouteBuilder) Patch(path string) *RouteBuilder {
	return b.method(Route_Patch, path)
}

func (b *RouteBuilder) Options(path string) *RouteBuilder {
	return b.method(Route_Options, path)
}

func (b *RouteBuilder) Head(path string) *RouteBuilder {
	return b.method(Route_Head, path)
}

func (b *RouteBuilder) OperationId(operationID string) *RouteBuilder {
	b.Route.OperationID = operationID
	if b.buildOperation != nil {
		b.buildOperation.OperationId(operationID)
	}
	return b
}

func (b *RouteBuilder) Summary(summary string) *RouteBuilder {
	if b.buildOperation != nil {
		b.buildOperation.Summary(summary)
	}
	return b
}

func (b *RouteBuilder) Description(description string) *RouteBuilder {
	if b.buildOperation != nil {
		b.buildOperation.Description(description)
	}
	return b
}

func (b *RouteBuilder) Tags(tags ...string) *RouteBuilder {
	if b.buildOperation != nil {
		b.buildOperation.Tags(tags...)
	}
	return b
}

func (b *RouteBuilder) RequestBody(fn func(*oas.BuildRequestBody)) *RouteBuilder {
	if b.buildOperation != nil {
		b.buildOperation.RequestBody(fn)
	}
	return b
}

func (b *RouteBuilder) Response(statusCode int, fn func(*oas.BuildResponse)) *RouteBuilder {
	if b.buildOperation != nil {
		b.buildOperation.Response(statusCode, fn)
	}
	return b
}

func (b *RouteBuilder) responseException(statusCode int, ex any) *RouteBuilder {
	metadata := doc.GetObjectMetadataOf(ex)
	if metadata == nil {
		return b.Response(statusCode, func(r *oas.BuildResponse) {
			r.Description(http.StatusText(statusCode)).
				Content(oas.ContentType_ApplicationJson, func(m *oas.BuildMediaType) {
					m.Schema(oas.Object())
				})
		})
	}

	exType := reflect.TypeOf(ex)
	if exType.Kind() == reflect.Pointer {
		exType = exType.Elem()
	}

	payloadType := exType

	if exType.Kind() == reflect.Struct {
		for i := 0; i < exType.NumField(); i++ {
			field := exType.Field(i)
			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				payloadType = field.Type
				break
			}
		}
	}

	payloadZero := reflect.New(payloadType).Elem().Interface()
	payloadSchema := oas.Struct(payloadZero)

	return b.Response(statusCode, func(r *oas.BuildResponse) {
		description := metadata.Description
		if description == "" {
			description = http.StatusText(statusCode)
		}
		r.Description(description).Content(oas.ContentType_ApplicationJson, func(m *oas.BuildMediaType) {
			m.Schema(payloadSchema).Example(metadata.Example)
		})
	})
}

// ResponseThrowsFromMetadata adiciona automaticamente todas as exceções declaradas
// via doc.Throws[...] usando o StructMetadata já resolvido.
//
// Uso recomendado:
//
//	cmdMeta := doc.GetStructMetadataAs[command.LoginProviderCallbackCommand]()
//	.ResponseThrowsFromMetadata(cmdMeta)
//
// Ou com alias curto:
//
//	.ThrowsFrom(cmdMeta)
func (b *RouteBuilder) ResponseThrowsFromMetadata(metadata *doc.ObjectMetadata, fallback ...any) *RouteBuilder {
	if metadata == nil && len(fallback) > 0 {
		metadata = doc.GetObjectMetadataOf(fallback[0])
	}
	if metadata == nil || len(metadata.Throws) == 0 {
		return b
	}

	descriptionsByStatusCode := make(map[int][]string)
	processedStatusCodes := make(map[int]struct{})

	for _, throw := range metadata.Throws {
		errorInstance := reflect.New(throw.ErrorType).Interface()
		statusCode := GetHTTPStatus(errorInstance)
		errorMeta := doc.GetObjectMetadataByType(throw.ErrorType)

		desc := throw.Description
		if desc == "" && errorMeta != nil && errorMeta.Description != "" {
			desc = errorMeta.Description
		}
		if desc == "" {
			desc = http.StatusText(statusCode)
		}

		found := slices.Contains(descriptionsByStatusCode[statusCode], desc)
		if !found {
			descriptionsByStatusCode[statusCode] = append(descriptionsByStatusCode[statusCode], desc)
		}

		if _, ok := processedStatusCodes[statusCode]; !ok {
			processedStatusCodes[statusCode] = struct{}{}
			b.Response(statusCode, func(r *oas.BuildResponse) {
				r.Description(desc)
				if errorMeta != nil {
					r.Content(oas.ContentType_ApplicationJson, func(m *oas.BuildMediaType) {
						m.Schema(oas.ObjectMetadata(errorMeta))
						if errorMeta.Example != nil {
							m.Example(errorMeta.Example)
						}
					})
				}
			})
		}
	}

	for statusCode := range processedStatusCodes {
		if descriptions, ok := descriptionsByStatusCode[statusCode]; ok {
			finalDesc := strings.Join(descriptions, " | ")
			b.Response(statusCode, func(r *oas.BuildResponse) {
				r.Description(finalDesc)
			})
		}
	}

	return b
}

func (b *RouteBuilder) ResponseNotfoundException() *RouteBuilder {
	return b.responseException(http.StatusNotFound, exception.NotFoundException{})
}

func (b *RouteBuilder) ResponsePreconditionFailedException() *RouteBuilder {
	return b.responseException(http.StatusPreconditionFailed, exception.PreconditionFailedException{})
}

func (b *RouteBuilder) ResponseUnauthorizedException() *RouteBuilder {
	return b.responseException(http.StatusUnauthorized, exception.UnauthorizedException{})
}

func (b *RouteBuilder) ResponseUnprocessableEntityException() *RouteBuilder {
	return b.responseException(http.StatusUnprocessableEntity, exception.UnprocessableEntityException{})
}

func (b *RouteBuilder) ResponseValidationException() *RouteBuilder {
	return b.responseException(http.StatusBadRequest, exception.ValidationException{})
}

func (b *RouteBuilder) ResponseConflictException() *RouteBuilder {
	return b.responseException(http.StatusConflict, exception.ConflictException{})
}

func (b *RouteBuilder) ResponseForbiddenException() *RouteBuilder {
	return b.responseException(http.StatusForbidden, exception.ForbiddenException{})
}

func (b *RouteBuilder) ResponseMethodNotAllowedException() *RouteBuilder {
	return b.responseException(http.StatusMethodNotAllowed, exception.MethodNotAllowedException{})
}

func (b *RouteBuilder) ResponseNotAcceptableException() *RouteBuilder {
	return b.responseException(http.StatusNotAcceptable, exception.NotAcceptableException{})
}

func (b *RouteBuilder) Deprecated(deprecated bool) *RouteBuilder {
	if b.buildOperation != nil {
		b.buildOperation.Deprecated(deprecated)
	}
	return b
}

func (b *RouteBuilder) Security(requirement *oas.SecurityRequirement) *RouteBuilder {
	if b.buildOperation != nil {
		b.buildOperation.Security(requirement)
	}
	return b
}

func (b *RouteBuilder) Handler(handler HandlerFN) *RouteBuilder {
	b.Route.Handler = handler
	return b
}

func (b *RouteBuilder) UseGuards(guards ...GuardFN) *RouteBuilder {
	b.Route.Guards = append(b.Route.Guards, guards...)
	return b
}

func (b *RouteBuilder) UseInterceptors(interceptors ...InterceptorFN) *RouteBuilder {
	b.Route.Interceptors = append(b.Route.Interceptors, interceptors...)
	return b
}
