package swagger

// Swagger https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md
type Swagger struct {
	SwaggerVersion      string                `json:"swagger" yaml:"swagger"`
	Info                *Info                 `json:"info" yaml:"info"`
	Host                string                `json:"host,omitempty" yaml:"host,omitempty"`
	BasePath            string                `json:"basePath,omitempty" yaml:"basePath,omitempty"`
	Schemes             []string              `json:"schemes,omitempty" yaml:"schemes,omitempty"`
	Consumes            []string              `json:"consumes,omitempty" yaml:"consumes,omitempty"`
	Produces            []string              `json:"produces,omitempty" yaml:"produces,omitempty"`
	Parameters          map[string]*Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Paths               map[string]*Item      `json:"paths" yaml:"paths"`
	Definitions         map[string]*Schema    `json:"definitions,omitempty" yaml:"definitions,omitempty"`
	SecurityDefinitions map[string]*Security  `json:"securityDefinitions,omitempty" yaml:"securityDefinitions,omitempty"`
	Security            []map[string][]string `json:"security,omitempty" yaml:"security,omitempty"`
	Tags                []*Tag                `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs        *ExternalDocs         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// Info The object provides metadata about the API. The metadata can be used by the clients if needed, and can be presented in the Swagger-UI for convenience.
type Info struct {
	Title          string   `json:"title" yaml:"title"`
	Description    string   `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *License `json:"license,omitempty" yaml:"license,omitempty"`
	Version        string   `json:"version" yaml:"version"`
}

// Contact information for the exposed API.
type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	EMail string `json:"email,omitempty" yaml:"email,omitempty"`
}

// License information for the exposed API.
type License struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

// Item Describes the operations available on a single path.
// A Path Item may be empty, due to ACL constraints.
// The path itself is still exposed to the documentation viewer but they will not know which operations and parameters are available.
type Item struct {
	Ref     string     `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Get     *Operation `json:"get,omitempty" yaml:"get,omitempty"`
	Put     *Operation `json:"put,omitempty" yaml:"put,omitempty"`
	Post    *Operation `json:"post,omitempty" yaml:"post,omitempty"`
	Delete  *Operation `json:"delete,omitempty" yaml:"delete,omitempty"`
	Options *Operation `json:"options,omitempty" yaml:"options,omitempty"`
	Head    *Operation `json:"head,omitempty" yaml:"head,omitempty"`
	Patch   *Operation `json:"patch,omitempty" yaml:"patch,omitempty"`
}

// Operation Describes a single API operation on a path.
type Operation struct {
	Tags         []string              `json:"tags,omitempty" yaml:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  string                `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	OperationID  string                `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Consumes     []string              `json:"consumes,omitempty" yaml:"consumes,omitempty"`
	Produces     []string              `json:"produces,omitempty" yaml:"produces,omitempty"`
	Parameters   []*Parameter          `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Responses    map[string]*Response  `json:"responses,omitempty" yaml:"responses,omitempty"`
	Schemes      []string              `json:"schemes,omitempty" yaml:"schemes,omitempty"`
	Deprecated   bool                  `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security     []map[string][]string `json:"security,omitempty" yaml:"security,omitempty"`
}

// ExternalDocs Allows referencing an external resource for extended documentation.
// document: https://swagger.staging.wpengine.com/specification/#externalDocumentationObject
type ExternalDocs struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string `json:"url" yaml:"url"`
}

// Parameter Describes a single operation parameter.
type Parameter struct {
	Name        string        `json:"name,omitempty" yaml:"name,omitempty"`
	In          string        `json:"in,omitempty" yaml:"in,omitempty"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool          `json:"required,omitempty" yaml:"required,omitempty"`
	Schema      *Schema       `json:"schema,omitempty" yaml:"schema,omitempty"`
	Type        string        `json:"type,omitempty" yaml:"type,omitempty"`
	Format      string        `json:"format,omitempty" yaml:"format,omitempty"`
	Items       *Items        `json:"items,omitempty" yaml:"items,omitempty"`
	Enum        []interface{} `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     interface{}   `json:"default,omitempty" yaml:"default,omitempty"`
	Minimum     *float64      `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum     *float64      `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	MinLength   *int64        `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength   *int64        `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	Ref         string        `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// Schema The Schema Object allows the definition of input and output data types.
// These types can be objects, but also primitives and arrays.
// This object is based on the JSON Schema Specification Draft 4 and uses a predefined subset of it.
// On top of this subset, there are extensions provided by this specification to allow for more complete documentation.
type Schema struct {
	Ref         string                `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Title       string                `json:"title,omitempty" yaml:"title,omitempty"`
	Format      string                `json:"format,omitempty" yaml:"format,omitempty"`
	Description string                `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool                  `json:"required,omitempty" yaml:"required,omitempty"`
	Type        string                `json:"type,omitempty" yaml:"type,omitempty"`
	Items       *Schema               `json:"items,omitempty" yaml:"items,omitempty"`
	Properties  map[string]*Propertie `json:"properties,omitempty" yaml:"properties,omitempty"`
	Enum        []interface{}         `json:"enum,omitempty" yaml:"enum,omitempty"`
}

// Propertie properties are taken from the JSON Schema definition but their definitions were adjusted to the Swagger Specification.
// Their definition is the same as the one from JSON Schema, only where the original definition references the JSON Schema definition, the Schema Object definition is used instead.
type Propertie struct {
	Ref                  string                `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Title                string                `json:"title,omitempty" yaml:"title,omitempty"`
	Description          string                `json:"description,omitempty" yaml:"description,omitempty"`
	Default              interface{}           `json:"default,omitempty" yaml:"default,omitempty"`
	Type                 string                `json:"type,omitempty" yaml:"type,omitempty"`
	Example              string                `json:"example,omitempty" yaml:"example,omitempty"`
	Required             bool                  `json:"required,omitempty" yaml:"required,omitempty"`
	Format               string                `json:"format,omitempty" yaml:"format,omitempty"`
	ReadOnly             bool                  `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	Properties           map[string]*Propertie `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items                *Propertie            `json:"items,omitempty" yaml:"items,omitempty"`
	AdditionalProperties *Propertie            `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	Enum                 []interface{}         `json:"enum,omitempty" yaml:"enum,omitempty"`
	Minimum              *float64              `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum              *float64              `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	MinLength            *int64                `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength            *int64                `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
}

// Items A limited subset of JSON-Schema's items object. It is used by parameter definitions that are not located in "body".
type Items struct {
	Type             string   `json:"type,omitempty" yaml:"type,omitempty"`
	Format           string   `json:"format,omitempty" yaml:"format,omitempty"`
	Items            []*Items `json:"items,omitempty" yaml:"items,omitempty"` //Required if type is "array". Describes the type of items in the array.
	CollectionFormat string   `json:"collectionFormat,omitempty" yaml:"collectionFormat,omitempty"`
	Default          string   `json:"default,omitempty" yaml:"default,omitempty"`
}

// Response A container for the expected responses of an operation.
type Response struct {
	Description string             `json:"description" yaml:"description"`
	Schema      *Schema            `json:"schema,omitempty" yaml:"schema,omitempty"`
	Headers     map[string]*Header `json:"headers,omitempty" yaml:"headers,omitempty"`
	Ref         string             `json:"$ref,omitempty" yaml:"$ref,omitempty"`
}

// Header describes the type of the header.
type Header struct {
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string        `json:"type,omitempty" yaml:"type,omitempty"`
	Format      string        `json:"format,omitempty" yaml:"format,omitempty"`
	Items       *Items        `json:"items,omitempty" yaml:"items,omitempty"`
	Enum        []interface{} `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     interface{}   `json:"default,omitempty" yaml:"default,omitempty"`
	Minimum     *float64      `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Maximum     *float64      `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	MinLength   *int64        `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength   *int64        `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
}

// Security Allows the definition of a security scheme that can be used by the operations
type Security struct {
	Type             string            `json:"type,omitempty" yaml:"type,omitempty"` // Valid values are "basic", "apiKey" or "oauth2".
	Description      string            `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string            `json:"name,omitempty" yaml:"name,omitempty"`
	In               string            `json:"in,omitempty" yaml:"in,omitempty"`     // Valid values are "query" or "header".
	Flow             string            `json:"flow,omitempty" yaml:"flow,omitempty"` // Valid values are "implicit", "password", "application" or "accessCode".
	AuthorizationURL string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty" yaml:"scopes,omitempty"` // The available scopes for the OAuth2 security scheme.
}

// Tag Allows adding meta data to a single tag that is used by the Operation Object
type Tag struct {
	Name         string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}
