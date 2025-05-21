package models

// OpenAPISpec represents a simplified OpenAPI specification
type OpenAPISpec struct {
	Info struct {
		Title       string `json:"title" yaml:"title"`
		Version     string `json:"version" yaml:"version"`
		Description string `json:"description" yaml:"description"`
	} `json:"info" yaml:"info"`
	Paths      map[string]map[string]Operation `json:"paths" yaml:"paths"`
	Components struct {
		Schemas map[string]Schema `json:"schemas" yaml:"schemas"`
	} `json:"components" yaml:"components"`
}

// Operation represents an API operation
type Operation struct {
	Method      string              // HTTP method (GET, POST, etc.) - populated during processing
	OperationID string              `json:"operationId" yaml:"operationId"`
	Summary     string              `json:"summary" yaml:"summary"`
	Description string              `json:"description" yaml:"description"`
	Parameters  []Parameter         `json:"parameters" yaml:"parameters"`
	RequestBody *RequestBody        `json:"requestBody" yaml:"requestBody"`
	Responses   map[string]Response `json:"responses" yaml:"responses"`
	Tags        []string            `json:"tags" yaml:"tags"`
}

// Parameter represents an API parameter
type Parameter struct {
	Name        string `json:"name" yaml:"name"`
	In          string `json:"in" yaml:"in"`
	Required    bool   `json:"required" yaml:"required"`
	Description string `json:"description" yaml:"description"`
	Schema      Schema `json:"schema" yaml:"schema"`
}

// RequestBody represents an API request body
type RequestBody struct {
	Required bool `json:"required" yaml:"required"`
	Content  map[string]struct {
		Schema Schema `json:"schema" yaml:"schema"`
	} `json:"content" yaml:"content"`
}

// Response represents an API response
type Response struct {
	Description string `json:"description" yaml:"description"`
	Content     map[string]struct {
		Schema Schema `json:"schema" yaml:"schema"`
	} `json:"content" yaml:"content"`
}

// Schema represents a data schema
type Schema struct {
	Type                 string            `json:"type" yaml:"type"`
	Format               string            `json:"format" yaml:"format"`
	Properties           map[string]Schema `json:"properties" yaml:"properties"`
	Items                *Schema           `json:"items" yaml:"items"`
	Ref                  string            `json:"$ref" yaml:"$ref"`
	Required             []string          `json:"required" yaml:"required"`
	Description          string            `json:"description" yaml:"description"`
	Enum                 []interface{}     `json:"enum" yaml:"enum"`
	AllOf                []Schema          `json:"allOf" yaml:"allOf"`
	OneOf                []Schema          `json:"oneOf" yaml:"oneOf"`
	AnyOf                []Schema          `json:"anyOf" yaml:"anyOf"`
	Not                  *Schema           `json:"not" yaml:"not"`
	AdditionalProperties *bool             `json:"additionalProperties" yaml:"additionalProperties"`
}
