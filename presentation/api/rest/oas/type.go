package oas

// #region SchemaTypeEnum

type SchemaTypeEnum string

const (
	SchemaType_Object  SchemaTypeEnum = "object"
	SchemaType_String  SchemaTypeEnum = "string"
	SchemaType_Integer SchemaTypeEnum = "integer"
	SchemaType_Number  SchemaTypeEnum = "number"
	SchemaType_Boolean SchemaTypeEnum = "boolean"
	SchemaType_Array   SchemaTypeEnum = "array"
	SchemaType_Null    SchemaTypeEnum = "null"
)

// #endregion
// #region ContentTypeEnum

type ContentTypeEnum string

const (
	ContentType_ApplicationJson ContentTypeEnum = "application/json"
	ContentType_TextPlain       ContentTypeEnum = "text/plain"
	ContentType_TextHtml        ContentTypeEnum = "text/html"
	ContentType_TextXml         ContentTypeEnum = "text/xml"
	ContentType_TextCsv         ContentTypeEnum = "text/csv"
	ContentType_ImageJpeg       ContentTypeEnum = "image/jpeg"
	ContentType_ImagePng        ContentTypeEnum = "image/png"
	ContentType_ImageGif        ContentTypeEnum = "image/gif"
	ContentType_ImageSvg        ContentTypeEnum = "image/svg+xml"
	ContentType_ImageWebp       ContentTypeEnum = "image/webp"
)

// #endregion
// #region Schema

type Schema struct {
	Type                  []SchemaTypeEnum   `json:"type,omitempty"`
	Format                string             `json:"format,omitempty"`
	Description           string             `json:"description,omitempty"`
	Properties            map[string]*Schema `json:"properties,omitempty"`
	Required              []string           `json:"required,omitempty"`
	Items                 []*Schema          `json:"items,omitempty"`
	Title                 string             `json:"title,omitempty"`
	MultipleOf            *float64           `json:"multipleOf,omitempty"`
	Maximum               *float64           `json:"maximum,omitempty"`
	ExclusiveMaximum      *bool              `json:"exclusiveMaximum,omitempty"`
	Minimum               *float64           `json:"minimum,omitempty"`
	ExclusiveMinimum      *bool              `json:"exclusiveMinimum,omitempty"`
	MaxLength             *int64             `json:"maxLength,omitempty"`
	MinLength             *int64             `json:"minLength,omitempty"`
	Pattern               string             `json:"pattern,omitempty"`
	MaxItems              *int64             `json:"maxItems,omitempty"`
	MinItems              *int64             `json:"minItems,omitempty"`
	UniqueItems           *bool              `json:"uniqueItems,omitempty"`
	MaxProperties         *int64             `json:"maxProperties,omitempty"`
	MinProperties         *int64             `json:"minProperties,omitempty"`
	Enum                  []any              `json:"enum,omitempty"`
	AllOf                 []*Schema          `json:"allOf,omitempty"`
	OneOf                 []*Schema          `json:"oneOf,omitempty"`
	AnyOf                 []*Schema          `json:"anyOf,omitempty"`
	Not                   *Schema            `json:"not,omitempty"`
	AdditionalProperties  any                `json:"additionalProperties,omitempty"`
	Default               any                `json:"default,omitempty"`
	Discriminator         *Discriminator     `json:"discriminator,omitempty"`
	ReadOnly              *bool              `json:"readOnly,omitempty"`
	WriteOnly             *bool              `json:"writeOnly,omitempty"`
	XML                   *XML               `json:"xml,omitempty"`
	ExternalDocs          *ExternalDocs      `json:"externalDocs,omitempty"`
	Example               any                `json:"example,omitempty"`
	Deprecated            *bool              `json:"deprecated,omitempty"`
	DependentSchemas      map[string]*Schema `json:"dependentSchemas,omitempty"`
	UnevaluatedItems      any                `json:"unevaluatedItems,omitempty"`
	UnevaluatedProperties any                `json:"unevaluatedProperties,omitempty"`
	If                    *Schema            `json:"if,omitempty"`
	Then                  *Schema            `json:"then,omitempty"`
	Else                  *Schema            `json:"else,omitempty"`
	ContentMediaType      string             `json:"contentMediaType,omitempty"`
	ContentEncoding       string             `json:"contentEncoding,omitempty"`
	Ref                   string             `json:"$ref,omitempty"`
	Const                 any                `json:"const,omitempty"`
}

type Discriminator struct {
	PropertyName string            `json:"propertyName"`
	Mapping      map[string]string `json:"mapping,omitempty"`
}

type XML struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute *bool  `json:"attribute,omitempty"`
	Wrapped   *bool  `json:"wrapped,omitempty"`
}

type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url"`
}

// #endregion
// #region Parameter

type Parameter struct {
	Name            string  `json:"name"`
	In              string  `json:"in"`
	Description     string  `json:"description,omitempty"`
	Required        bool    `json:"required,omitempty"`
	Schema          *Schema `json:"schema,omitempty"`
	Example         any     `json:"example,omitempty"`
	Deprecated      bool    `json:"deprecated,omitempty"`
	AllowEmptyValue bool    `json:"allowEmptyValue,omitempty"`
	AllowReserved   bool    `json:"allowReserved,omitempty"`
}

// #endregion
// #region Operation

type Operation struct {
	Tags         []string               `json:"tags,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	Description  string                 `json:"description,omitempty"`
	OperationId  string                 `json:"operationId,omitempty"`
	Parameters   []*Parameter           `json:"parameters,omitempty"`
	RequestBody  *RequestBody           `json:"requestBody,omitempty"`
	Responses    map[int]*Response      `json:"responses,omitempty"`
	Deprecated   bool                   `json:"deprecated,omitempty"`
	ExternalDocs *ExternalDocs          `json:"externalDocs,omitempty"`
	Security     []*SecurityRequirement `json:"security,omitempty"`
	Servers      []*Server              `json:"servers,omitempty"`
}

// #endregion
// #region Response

type Response struct {
	Description string                `json:"description"`
	Content     map[string]*MediaType `json:"content,omitempty"`
	Headers     map[string]*Header    `json:"headers,omitempty"`
	Links       map[string]*Link      `json:"links,omitempty"`
}

// #endregion
// #region MediaType

type MediaType struct {
	Schema   *Schema                   `json:"schema,omitempty"`
	Example  any                       `json:"example,omitempty"`
	Examples map[string]*ExampleObject `json:"examples,omitempty"`
	Encoding map[string]*Encoding      `json:"encoding,omitempty"`
}

// #endregion
// #region ExampleObject

type ExampleObject struct {
	Summary       string `json:"summary,omitempty"`
	Description   string `json:"description,omitempty"`
	Value         any    `json:"value,omitempty"`
	ExternalValue string `json:"externalValue,omitempty"`
}

// #endregion
// #region Encoding

type Encoding struct {
	ContentType   string             `json:"contentType,omitempty"`
	Headers       map[string]*Header `json:"headers,omitempty"`
	Style         string             `json:"style,omitempty"`
	Explode       *bool              `json:"explode,omitempty"`
	AllowReserved *bool              `json:"allowReserved,omitempty"`
}

// #endregion
// #region Header

type Header struct {
	Description     string  `json:"description,omitempty"`
	Required        bool    `json:"required,omitempty"`
	Schema          *Schema `json:"schema,omitempty"`
	Example         any     `json:"example,omitempty"`
	Deprecated      bool    `json:"deprecated,omitempty"`
	AllowEmptyValue bool    `json:"allowEmptyValue,omitempty"`
}

// #endregion
// #region Link

type Link struct {
	OperationRef string         `json:"operationRef,omitempty"`
	OperationId  string         `json:"operationId,omitempty"`
	Parameters   map[string]any `json:"parameters,omitempty"`
	RequestBody  any            `json:"requestBody,omitempty"`
	Description  string         `json:"description,omitempty"`
	Server       *Server        `json:"server,omitempty"`
}

// #endregion
// #region Server

type Server struct {
	URL         string                     `json:"url"`
	Description string                     `json:"description,omitempty"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty"`
}

// #endregion
// #region ServerVariable

type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

// #endregion
// #region SecurityRequirement

type SecurityRequirement map[string][]string

// #endregion
// #region SecurityScheme

type SecurityScheme struct {
	Type             string      `json:"type"`
	Description      string      `json:"description,omitempty"`
	Name             string      `json:"name,omitempty"`
	In               string      `json:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty"`
	OpenIdConnectURL string      `json:"openIdConnectUrl,omitempty"`
}

// #endregion
// #region OAuthFlows

type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

// #endregion
// #region OAuthFlow

type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes"`
}

// #endregion
// #region Tag

type Tag struct {
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

// #endregion
// #region RequestBody

type RequestBody struct {
	Description string                `json:"description,omitempty"`
	Content     map[string]*MediaType `json:"content"`
	Required    bool                  `json:"required,omitempty"`
}

// #endregion
// #region Components

type Components struct {
	Schemas         map[string]*Schema         `json:"schemas,omitempty"`
	Responses       map[string]*Response       `json:"responses,omitempty"`
	Parameters      map[string]*Parameter      `json:"parameters,omitempty"`
	Examples        map[string]*ExampleObject  `json:"examples,omitempty"`
	RequestBodies   map[string]*RequestBody    `json:"requestBodies,omitempty"`
	Headers         map[string]*Header         `json:"headers,omitempty"`
	SecuritySchemes map[string]*SecurityScheme `json:"securitySchemes,omitempty"`
	Links           map[string]*Link           `json:"links,omitempty"`
	Callbacks       map[string]*Callback       `json:"callbacks,omitempty"`
}

// #endregion
// #region Callback

type Callback map[string]*PathItem

// #endregion
// #region Info

type Info struct {
	Title          string   `json:"title"`
	Summary        string   `json:"summary,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
	Version        string   `json:"version"`
}

// #endregion
// #region Contact

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// #endregion
// #region License

type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

// #endregion
// #region PathItem

type PathItem struct {
	Ref              string       `json:"$ref,omitempty"`
	Summary          string       `json:"summary,omitempty"`
	Description      string       `json:"description,omitempty"`
	GetOperation     *Operation   `json:"get,omitempty"`
	PostOperation    *Operation   `json:"post,omitempty"`
	PutOperation     *Operation   `json:"put,omitempty"`
	DeleteOperation  *Operation   `json:"delete,omitempty"`
	OptionsOperation *Operation   `json:"options,omitempty"`
	HeadOperation    *Operation   `json:"head,omitempty"`
	PatchOperation   *Operation   `json:"patch,omitempty"`
	TraceOperation   *Operation   `json:"trace,omitempty"`
	Servers          []*Server    `json:"servers,omitempty"`
	Parameters       []*Parameter `json:"parameters,omitempty"`
}

func (p *PathItem) Get(fn func(*BuildOperation)) *PathItem {
	builder := NewOperation()
	fn(builder)
	p.GetOperation = builder.Build()
	return p
}

func (p *PathItem) Post(fn func(*BuildOperation)) *PathItem {
	builder := NewOperation()
	fn(builder)
	p.PostOperation = builder.Build()
	return p
}

func (p *PathItem) Put(fn func(*BuildOperation)) *PathItem {
	builder := NewOperation()
	fn(builder)
	p.PutOperation = builder.Build()
	return p
}

func (p *PathItem) Delete(fn func(*BuildOperation)) *PathItem {
	builder := NewOperation()
	fn(builder)
	p.DeleteOperation = builder.Build()
	return p
}

func (p *PathItem) Options(fn func(*BuildOperation)) *PathItem {
	builder := NewOperation()
	fn(builder)
	p.OptionsOperation = builder.Build()
	return p
}

func (p *PathItem) Head(fn func(*BuildOperation)) *PathItem {
	builder := NewOperation()
	fn(builder)
	p.HeadOperation = builder.Build()
	return p
}

func (p *PathItem) Patch(fn func(*BuildOperation)) *PathItem {
	builder := NewOperation()
	fn(builder)
	p.PatchOperation = builder.Build()
	return p
}

func (p *PathItem) Trace(fn func(*BuildOperation)) *PathItem {
	builder := NewOperation()
	fn(builder)
	p.TraceOperation = builder.Build()
	return p
}

// #endregion
// #region Paths

type Paths map[string]*PathItem

// #endregion
// #region Webhooks

type Webhooks map[string]*PathItem

// #endregion
// #region OpenAPI

type OpenAPI struct {
	OpenAPI           string                 `json:"openapi"`
	Info              *Info                  `json:"info"`
	JsonSchemaDialect string                 `json:"jsonSchemaDialect,omitempty"`
	Servers           []*Server              `json:"servers,omitempty"`
	Paths             Paths                  `json:"paths,omitempty"`
	Webhooks          Webhooks               `json:"webhooks,omitempty"`
	Components        *Components            `json:"components,omitempty"`
	Security          []*SecurityRequirement `json:"security,omitempty"`
	Tags              []*Tag                 `json:"tags,omitempty"`
	ExternalDocs      *ExternalDocs          `json:"externalDocs,omitempty"`
}

func (o *OpenAPI) Path(path string) *PathItem {
	if o.Paths == nil {
		o.Paths = make(Paths)
	}

	pathItem := o.Paths[path]
	if pathItem == nil {
		pathItem = &PathItem{}
		o.Paths[path] = pathItem
	}
	return pathItem
}

// #endregion
