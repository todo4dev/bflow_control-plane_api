package oas

import (
	"net/http"
	"reflect"
	"src/core/meta"
	"src/domain/exception"
	"strings"
	"time"
)

// #region BuildSchema

type BuildSchema struct {
	schema *Schema
}

func ObjectMetadata(metadata *meta.ObjectMetadata) *Schema {
	if metadata == nil || metadata.Type == nil {
		return Object()
	}

	base := Struct(reflect.New(metadata.Type).Elem().Interface())

	if metadata.Description != "" {
		base.Description = metadata.Description
	}
	if metadata.Example != nil {
		base.Example = metadata.Example
	}

	t := metadata.Type

	for fieldName, fieldMeta := range metadata.Fields {
		field, ok := t.FieldByName(fieldName)
		if !ok {
			continue
		}

		jsonName := field.Name
		if tag := field.Tag.Get("json"); tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" && parts[0] != "-" {
				jsonName = parts[0]
			}
		}

		prop, ok := base.Properties[jsonName]
		if !ok || prop == nil {
			continue
		}

		if fieldMeta.Description != "" {
			prop.Description = fieldMeta.Description
		}
		if fieldMeta.Example != nil {
			prop.Example = fieldMeta.Example
		}

		if fieldMeta.Nullable {
			hasNull := false
			for _, tpe := range prop.Type {
				if tpe == SchemaType_Null {
					hasNull = true
					break
				}
			}
			if !hasNull {
				prop.Type = append(prop.Type, SchemaType_Null)
			}
		}
	}

	return base
}
func Struct(s any) *Schema {
	schema := &Schema{
		Type:       []SchemaTypeEnum{SchemaType_Object},
		Properties: make(map[string]*Schema),
	}

	t := reflect.TypeOf(s)
	if t == nil {
		return schema
	}
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return schema
	}

	timeType := reflect.TypeOf(time.Time{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		fieldName := field.Name
		omitEmpty := false

		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] == "-" {
				// ignora campo
				continue
			}
			if parts[0] != "" {
				fieldName = parts[0]
			}
			for _, opt := range parts[1:] {
				if opt == "omitempty" {
					omitEmpty = true
					break
				}
			}
		}

		fieldType := field.Type
		isPtr := fieldType.Kind() == reflect.Pointer
		if isPtr {
			fieldType = fieldType.Elem()
		}

		var propSchema *Schema

		switch fieldType.Kind() {
		case reflect.Struct:
			if fieldType == timeType {
				propSchema = &Schema{
					Type:   []SchemaTypeEnum{SchemaType_String},
					Format: "date-time",
				}
			} else {
				propSchema = Struct(reflect.New(fieldType).Elem().Interface())
			}
		case reflect.Slice, reflect.Array:
			elemType := fieldType.Elem()

			switch elemType.Kind() {
			case reflect.Struct:
				if elemType == timeType {
					itemSchema := &Schema{
						Type:   []SchemaTypeEnum{SchemaType_String},
						Format: "date-time",
					}
					propSchema = &Schema{
						Type:  []SchemaTypeEnum{SchemaType_Array},
						Items: []*Schema{itemSchema},
					}
				} else {
					propSchema = &Schema{
						Type: []SchemaTypeEnum{SchemaType_Array},
						Items: []*Schema{
							Struct(reflect.New(elemType).Elem().Interface()),
						},
					}
				}
			default:
				itemSchema := &Schema{
					Type: []SchemaTypeEnum{goTypeToOASType(elemType)},
				}
				propSchema = &Schema{
					Type:  []SchemaTypeEnum{SchemaType_Array},
					Items: []*Schema{itemSchema},
				}
			}
		default:
			propSchema = &Schema{
				Type: []SchemaTypeEnum{goTypeToOASType(fieldType)},
			}
		}

		// required = não ponteiro e sem omitempty
		if !isPtr && !omitEmpty {
			schema.Required = append(schema.Required, fieldName)
		}

		schema.Properties[fieldName] = propSchema
	}

	return schema
}
func Object(fn ...func(*BuildSchema)) *Schema {
	builder := &BuildSchema{
		schema: &Schema{
			Type:       []SchemaTypeEnum{SchemaType_Object},
			Properties: make(map[string]*Schema),
		},
	}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}
func String(fn ...func(*BuildSchema)) *Schema {
	builder := &BuildSchema{schema: &Schema{Type: []SchemaTypeEnum{SchemaType_String}}}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}
func Integer(fn ...func(*BuildSchema)) *Schema {
	builder := &BuildSchema{schema: &Schema{Type: []SchemaTypeEnum{SchemaType_Integer}}}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}
func Number(fn ...func(*BuildSchema)) *Schema {
	builder := &BuildSchema{schema: &Schema{Type: []SchemaTypeEnum{SchemaType_Number}}}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}
func Boolean(fn ...func(*BuildSchema)) *Schema {
	builder := &BuildSchema{schema: &Schema{Type: []SchemaTypeEnum{SchemaType_Boolean}}}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}
func Array(fn ...func(*BuildSchema)) *Schema {
	builder := &BuildSchema{schema: &Schema{Type: []SchemaTypeEnum{SchemaType_Array}}}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}
func Null(fn ...func(*BuildSchema)) *Schema {
	builder := &BuildSchema{schema: &Schema{Type: []SchemaTypeEnum{SchemaType_Null}}}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}
func (b *BuildSchema) Type(types ...SchemaTypeEnum) *BuildSchema {
	b.schema.Type = types
	return b
}
func (b *BuildSchema) Format(format string) *BuildSchema {
	b.schema.Format = format
	return b
}
func (b *BuildSchema) Title(title ...string) *BuildSchema {
	b.schema.Title = strings.Join(title, " ")
	return b
}
func (b *BuildSchema) Description(description ...string) *BuildSchema {
	b.schema.Description = strings.Join(description, " ")
	return b
}
func (b *BuildSchema) Property(name string, schema *Schema) *BuildSchema {
	if b.schema.Properties == nil {
		b.schema.Properties = make(map[string]*Schema)
	}
	b.schema.Properties[name] = schema
	return b
}
func (b *BuildSchema) Required(fields ...string) *BuildSchema {
	b.schema.Required = append(b.schema.Required, fields...)
	return b
}
func (b *BuildSchema) Items(schema *Schema) *BuildSchema {
	b.schema.Items = make([]*Schema, 1)
	b.schema.Items[0] = schema
	return b
}
func (b *BuildSchema) Enum(values ...any) *BuildSchema {
	b.schema.Enum = values
	return b
}
func (b *BuildSchema) Default(value any) *BuildSchema {
	b.schema.Default = value
	return b
}
func (b *BuildSchema) Example(value any) *BuildSchema {
	b.schema.Example = value
	return b
}
func (b *BuildSchema) MinLength(min int64) *BuildSchema {
	b.schema.MinLength = &min
	return b
}
func (b *BuildSchema) MaxLength(max int64) *BuildSchema {
	b.schema.MaxLength = &max
	return b
}
func (b *BuildSchema) Minimum(min float64) *BuildSchema {
	b.schema.Minimum = &min
	return b
}
func (b *BuildSchema) Maximum(max float64) *BuildSchema {
	b.schema.Maximum = &max
	return b
}
func (b *BuildSchema) Pattern(pattern string) *BuildSchema {
	b.schema.Pattern = pattern
	return b
}
func (b *BuildSchema) Deprecated(deprecated bool) *BuildSchema {
	b.schema.Deprecated = &deprecated
	return b
}
func (b *BuildSchema) ReadOnly(readOnly bool) *BuildSchema {
	b.schema.ReadOnly = &readOnly
	return b
}
func (b *BuildSchema) WriteOnly(writeOnly bool) *BuildSchema {
	b.schema.WriteOnly = &writeOnly
	return b
}
func (b *BuildSchema) Ref(ref string) *BuildSchema {
	b.schema.Ref = ref
	return b
}
func (b *BuildSchema) Nullable() *BuildSchema {
	if b.schema.Type != nil {
		b.schema.Type = append(b.schema.Type, SchemaType_Null)
	}
	return b
}
func (b *BuildSchema) Build() *Schema {
	return b.schema
}

// #endregion
// #region BuildMediaType

type BuildMediaType struct {
	mediaType *MediaType
}

func NewMediaType() *BuildMediaType {
	return &BuildMediaType{mediaType: &MediaType{}}
}
func (b *BuildMediaType) Schema(schema *Schema) *BuildMediaType {
	b.mediaType.Schema = schema
	return b
}
func (b *BuildMediaType) Example(example any) *BuildMediaType {
	b.mediaType.Example = example
	return b
}
func (b *BuildMediaType) Examples(name string, example *ExampleObject) *BuildMediaType {
	if b.mediaType.Examples == nil {
		b.mediaType.Examples = make(map[string]*ExampleObject)
	}
	b.mediaType.Examples[name] = example
	return b
}
func (b *BuildMediaType) Build() *MediaType {
	return b.mediaType
}

// #endregion
// #region BuildHeader

type BuildHeader struct {
	header *Header
}

func NewHeader() *BuildHeader {
	return &BuildHeader{header: &Header{}}
}
func (b *BuildHeader) Description(description ...string) *BuildHeader {
	b.header.Description = strings.Join(description, " ")
	return b
}
func (b *BuildHeader) Required(required bool) *BuildHeader {
	b.header.Required = required
	return b
}
func (b *BuildHeader) Build() *Header {
	return b.header
}

// #endregion
// #region BuildResponse

type BuildResponse struct {
	response *Response
}

func NewResponse() *BuildResponse {
	return &BuildResponse{response: &Response{}}
}

func (b *BuildResponse) Description(description ...string) *BuildResponse {
	b.response.Description = strings.Join(description, " ")
	return b
}

func (b *BuildResponse) Content(contentType ContentTypeEnum, fn func(*BuildMediaType)) *BuildResponse {
	if b.response.Content == nil {
		b.response.Content = make(map[string]*MediaType)
	}

	mtBuilder := NewMediaType()
	fn(mtBuilder)

	b.response.Content[string(contentType)] = mtBuilder.Build()
	return b
}

func (b *BuildResponse) Header(name string, fn func(*BuildHeader)) *BuildResponse {
	if b.response.Headers == nil {
		b.response.Headers = make(map[string]*Header)
	}

	hBuilder := NewHeader()
	fn(hBuilder)

	b.response.Headers[name] = hBuilder.Build()
	return b
}

func (b *BuildResponse) Build() *Response {
	return b.response
}

// ThrowsFromMetadata usa o StructMetadata (ex: de uma Query) para
// configurar esta response com base no primeiro `Throws` declarado.
//
// Exemplo de uso:
//
//	meta := doc.GetStructMetadataAs[query.HealthQuery]()
//	r.ThrowsFromMetadata(meta)
//
// Ele procura o primeiro Throws, resolve a StructMetadata do tipo de erro
// (ex: exception.InternalException) e usa isso como schema + example.
func (b *BuildResponse) ThrowsFromMetadata(metadata *meta.ObjectMetadata) *BuildResponse {
	if metadata == nil || len(metadata.Throws) == 0 {
		return b
	}

	throws := metadata.Throws[0]
	exceptionMeta := meta.GetObjectMetadataByType(throws.ErrorType)

	description := throws.Description
	if description == "" {
		description = metadata.Description
	}
	if description == "" {
		description = http.StatusText(http.StatusInternalServerError)
	}

	b.Description(description).
		Content(ContentType_ApplicationJson, func(m *BuildMediaType) {
			if exceptionMeta != nil {
				m.Schema(ObjectMetadata(exceptionMeta))
				if exceptionMeta.Example != nil {
					m.Example(exceptionMeta.Example)
				}
			} else {
				m.Schema(Object())
			}
		})

	return b
}

// #endregion
// #region BuildRequestBody

type BuildRequestBody struct {
	requestBody *RequestBody
}

func NewRequestBody() *BuildRequestBody {
	return &BuildRequestBody{requestBody: &RequestBody{Content: make(map[string]*MediaType)}}
}
func (b *BuildRequestBody) Description(description ...string) *BuildRequestBody {
	b.requestBody.Description += strings.Join(description, " ")
	return b
}
func (b *BuildRequestBody) Required(required bool) *BuildRequestBody {
	b.requestBody.Required = required
	return b
}
func (b *BuildRequestBody) Content(contentType ContentTypeEnum, fn func(*BuildMediaType)) *BuildRequestBody {
	if b.requestBody.Content == nil {
		b.requestBody.Content = make(map[string]*MediaType)
	}

	mtBuilder := NewMediaType()
	fn(mtBuilder)

	b.requestBody.Content[string(contentType)] = mtBuilder.Build()
	return b
}
func (b *BuildRequestBody) Build() *RequestBody {
	return b.requestBody
}

// #endregion
// #region BuildParameter

type BuildParameter struct {
	parameter *Parameter
}

func PathParameter() *BuildParameter {
	return &BuildParameter{parameter: &Parameter{In: "path", Required: true}}
}
func QueryParameter() *BuildParameter {
	return &BuildParameter{parameter: &Parameter{In: "query"}}
}
func HeaderParameter() *BuildParameter {
	return &BuildParameter{parameter: &Parameter{In: "header"}}
}
func CookieParameter() *BuildParameter {
	return &BuildParameter{parameter: &Parameter{In: "cookie"}}
}
func (b *BuildParameter) Name(name string) *BuildParameter {
	b.parameter.Name = name
	return b
}
func (b *BuildParameter) In(in string) *BuildParameter {
	b.parameter.In = in
	return b
}
func (b *BuildParameter) Description(description ...string) *BuildParameter {
	b.parameter.Description += strings.Join(description, " ")
	return b
}
func (b *BuildParameter) Required(required bool) *BuildParameter {
	b.parameter.Required = required
	return b
}
func (b *BuildParameter) Schema(schema *Schema) *BuildParameter {
	b.parameter.Schema = schema
	return b
}
func (b *BuildParameter) Example(example any) *BuildParameter {
	b.parameter.Example = example
	return b
}
func (b *BuildParameter) Deprecated(deprecated bool) *BuildParameter {
	b.parameter.Deprecated = deprecated
	return b
}
func (b *BuildParameter) Build() *Parameter {
	return b.parameter
}

// #endregion
// #region BuildOperation

type BuildOperation struct {
	operation *Operation
}

func NewOperation() *BuildOperation {
	return &BuildOperation{operation: &Operation{}}
}
func (b *BuildOperation) Tags(tags ...string) *BuildOperation {
	b.operation.Tags = tags
	return b
}
func (b *BuildOperation) Summary(summary ...string) *BuildOperation {
	b.operation.Summary += strings.Join(summary, " ")
	return b
}
func (b *BuildOperation) Description(description ...string) *BuildOperation {
	b.operation.Description += strings.Join(description, " ")
	return b
}
func (b *BuildOperation) OperationId(operationId string) *BuildOperation {
	b.operation.OperationId = operationId
	return b
}
func (b *BuildOperation) PathParameter(fn func(*BuildParameter)) *BuildOperation {
	builder := PathParameter()
	fn(builder)
	b.operation.Parameters = append(b.operation.Parameters, builder.Build())
	return b
}
func (b *BuildOperation) QueryParameter(fn func(*BuildParameter)) *BuildOperation {
	builder := QueryParameter()
	fn(builder)
	b.operation.Parameters = append(b.operation.Parameters, builder.Build())
	return b
}
func (b *BuildOperation) HeaderParameter(fn func(*BuildParameter)) *BuildOperation {
	builder := HeaderParameter()
	fn(builder)
	b.operation.Parameters = append(b.operation.Parameters, builder.Build())
	return b
}
func (b *BuildOperation) CookieParameter(fn func(*BuildParameter)) *BuildOperation {
	builder := CookieParameter()
	fn(builder)
	b.operation.Parameters = append(b.operation.Parameters, builder.Build())
	return b
}
func (b *BuildOperation) RequestBody(fn func(*BuildRequestBody)) *BuildOperation {
	builder := NewRequestBody()
	fn(builder)
	b.operation.RequestBody = builder.Build()
	return b
}
func (b *BuildOperation) Response(statusCode int, fn func(*BuildResponse)) *BuildOperation {
	builder := NewResponse()
	fn(builder)
	if b.operation.Responses == nil {
		b.operation.Responses = make(map[int]*Response)
	}

	b.operation.Responses[statusCode] = builder.Build()

	// 500
	internalMetadata := meta.GetObjectMetadataAs[exception.Internal]()
	internalSchema := Struct(reflect.New(reflect.TypeOf(exception.Internal{})).Elem().Interface())
	b.operation.Responses[http.StatusInternalServerError] = NewResponse().
		Description(internalMetadata.Description).
		Content(ContentType_ApplicationJson, func(m *BuildMediaType) { m.Schema(internalSchema).Example(internalMetadata.Example) }).
		Build()

	return b
}
func (b *BuildOperation) Deprecated(deprecated bool) *BuildOperation {
	b.operation.Deprecated = deprecated
	return b
}
func (b *BuildOperation) Security(requirement *SecurityRequirement) *BuildOperation {
	b.operation.Security = append(b.operation.Security, requirement)
	return b
}
func (b *BuildOperation) Build() *Operation {
	return b.operation
}

// #endregion
// #region BuildPathItem

type BuildPathItem struct {
	pathItem *PathItem
}

func NewPathItem() *BuildPathItem {
	return &BuildPathItem{pathItem: &PathItem{}}
}
func (b *BuildPathItem) Summary(summary ...string) *BuildPathItem {
	b.pathItem.Summary += strings.Join(summary, " ")
	return b
}
func (b *BuildPathItem) Description(description ...string) *BuildPathItem {
	b.pathItem.Description += strings.Join(description, " ")
	return b
}
func (b *BuildPathItem) Get(fn func(*BuildOperation)) *BuildPathItem {
	builder := NewOperation()
	fn(builder)
	b.pathItem.GetOperation = builder.Build()
	return b
}
func (b *BuildPathItem) Post(fn func(*BuildOperation)) *BuildPathItem {
	builder := NewOperation()
	fn(builder)
	b.pathItem.PostOperation = builder.Build()
	return b
}
func (b *BuildPathItem) Put(fn func(*BuildOperation)) *BuildPathItem {
	builder := NewOperation()
	fn(builder)
	b.pathItem.PutOperation = builder.Build()
	return b
}
func (b *BuildPathItem) Delete(fn func(*BuildOperation)) *BuildPathItem {
	builder := NewOperation()
	fn(builder)
	b.pathItem.DeleteOperation = builder.Build()
	return b
}
func (b *BuildPathItem) Patch(fn func(*BuildOperation)) *BuildPathItem {
	builder := NewOperation()
	fn(builder)
	b.pathItem.PatchOperation = builder.Build()
	return b
}
func (b *BuildPathItem) Parameter(fn func(*BuildParameter)) *BuildPathItem {
	builder := PathParameter()
	fn(builder)
	b.pathItem.Parameters = append(b.pathItem.Parameters, builder.Build())
	return b
}
func (b *BuildPathItem) Build() *PathItem {
	return b.pathItem
}
func (b *BuildPathItem) GetBuilder(op *Operation) *BuildPathItem {
	b.pathItem.GetOperation = op
	return b
}
func (b *BuildPathItem) PostBuilder(op *Operation) *BuildPathItem {
	b.pathItem.PostOperation = op
	return b
}
func (b *BuildPathItem) PutBuilder(op *Operation) *BuildPathItem {
	b.pathItem.PutOperation = op
	return b
}
func (b *BuildPathItem) DeleteBuilder(op *Operation) *BuildPathItem {
	b.pathItem.DeleteOperation = op
	return b
}
func (b *BuildPathItem) PatchBuilder(op *Operation) *BuildPathItem {
	b.pathItem.PatchOperation = op
	return b
}

// #endregion
// #region BuildInfo

type BuildInfo struct {
	info *Info
}

func NewInfo(fn ...func(*BuildInfo)) *Info {
	builder := &BuildInfo{info: &Info{}}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}

func (b *BuildInfo) Summary(summary ...string) *BuildInfo {
	b.info.Summary += strings.Join(summary, " ")
	return b
}

func (b *BuildInfo) Description(description ...string) *BuildInfo {
	b.info.Description += strings.Join(description, " ")
	return b
}

func (b *BuildInfo) TermsOfService(termsOfService ...string) *BuildInfo {
	b.info.TermsOfService += strings.Join(termsOfService, " ")
	return b
}

func (b *BuildInfo) Contact(name, url, email string) *BuildInfo {
	b.info.Contact = &Contact{Name: name, URL: url, Email: email}
	return b
}

func (b *BuildInfo) License(name, url string) *BuildInfo {
	b.info.License = &License{Name: name, URL: url}
	return b
}

func (b *BuildInfo) Build() *Info {
	return b.info
}

// #endregion
// #region BuildSecurityScheme

type BuildSecurityScheme struct {
	securityScheme *SecurityScheme
}

func NewSecurityScheme(fn ...func(*BuildSecurityScheme)) *SecurityScheme {
	builder := &BuildSecurityScheme{securityScheme: &SecurityScheme{}}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}
func (b *BuildSecurityScheme) Type(t string) *BuildSecurityScheme {
	b.securityScheme.Type = t
	return b
}
func (b *BuildSecurityScheme) Name(name string) *BuildSecurityScheme {
	b.securityScheme.Name = name
	return b
}

func (b *BuildSecurityScheme) Scheme(scheme string) *BuildSecurityScheme {
	b.securityScheme.Scheme = scheme
	return b
}

func (b *BuildSecurityScheme) BearerFormat(bearerFormat string) *BuildSecurityScheme {
	b.securityScheme.BearerFormat = bearerFormat
	return b
}

func (b *BuildSecurityScheme) Description(description ...string) *BuildSecurityScheme {
	b.securityScheme.Description = strings.Join(description, " ")
	return b
}

func (b *BuildSecurityScheme) Build() *SecurityScheme {
	return b.securityScheme
}

// #endregion
// #region BuildComponents

type BuildComponents struct {
	components *Components
}

func NewComponents(fn ...func(*BuildComponents)) *Components {
	builder := &BuildComponents{components: &Components{}}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}
func (b *BuildComponents) SecurityScheme(name string, fn func(*BuildSecurityScheme)) *BuildComponents {
	if b.components.SecuritySchemes == nil {
		b.components.SecuritySchemes = make(map[string]*SecurityScheme)
	}

	scheme, ok := b.components.SecuritySchemes[name]
	if !ok || scheme == nil {
		scheme = &SecurityScheme{}
		b.components.SecuritySchemes[name] = scheme
	}

	fn(&BuildSecurityScheme{securityScheme: scheme})
	return b
}
func (b *BuildComponents) Build() *Components {
	return b.components
}

// #endregion
// #region BuildOpenAPI

type BuildOpenAPI struct {
	openapi OpenAPI
}

func NewOpenAPI(fn ...func(*BuildOpenAPI)) OpenAPI {
	builder := &BuildOpenAPI{openapi: OpenAPI{
		OpenAPI: "3.1.0",
		Paths:   make(Paths),
	}}
	if len(fn) != 0 {
		fn[0](builder)
	}
	return builder.Build()
}
func (b *BuildOpenAPI) Info(title, version string, fn func(*BuildInfo)) *BuildOpenAPI {
	if b.openapi.Info == nil {
		b.openapi.Info = &Info{}
	}
	b.openapi.Info.Title = title
	b.openapi.Info.Version = version
	fn(&BuildInfo{info: b.openapi.Info})
	return b
}
func (b *BuildOpenAPI) JsonSchemaDialect(dialect string) *BuildOpenAPI {
	b.openapi.JsonSchemaDialect = dialect
	return b
}
func (b *BuildOpenAPI) Server(url, description string) *BuildOpenAPI {
	b.openapi.Servers = append(b.openapi.Servers, &Server{URL: url, Description: description})
	return b
}
func (b *BuildOpenAPI) Path(path string, fn func(*BuildPathItem)) *BuildOpenAPI {
	builder := NewPathItem()
	fn(builder)
	if b.openapi.Paths == nil {
		b.openapi.Paths = make(Paths)
	}
	b.openapi.Paths[path] = builder.Build()
	return b
}
func (b *BuildOpenAPI) Components(fn func(*BuildComponents)) *BuildOpenAPI {
	if b.openapi.Components == nil {
		b.openapi.Components = &Components{}
	}
	fn(&BuildComponents{components: b.openapi.Components})
	return b
}
func (b *BuildOpenAPI) SecurityScheme(name string, fn func(*BuildSecurityScheme)) *BuildOpenAPI {
	if b.openapi.Components == nil {
		b.openapi.Components = &Components{}
	}
	fn(&BuildSecurityScheme{securityScheme: b.openapi.Components.SecuritySchemes[name]})
	return b
}
func (b *BuildOpenAPI) Tag(name, description string) *BuildOpenAPI {
	b.openapi.Tags = append(b.openapi.Tags, &Tag{Name: name, Description: description})
	return b
}
func (b *BuildOpenAPI) Build() OpenAPI {
	return b.openapi
}

// #endregion

func goTypeToOASType(t reflect.Type) SchemaTypeEnum {
	switch t.Kind() {
	case reflect.Bool:
		return SchemaType_Boolean
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return SchemaType_Integer
	case reflect.Float32, reflect.Float64:
		return SchemaType_Number
	case reflect.String:
		return SchemaType_String
	case reflect.Slice, reflect.Array:
		return SchemaType_Array
	case reflect.Map:
		return SchemaType_Object
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			return SchemaType_String
		}
		return SchemaType_Object
	default:
		return SchemaType_String
	}
}
