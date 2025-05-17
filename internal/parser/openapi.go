package parser

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

// OpenAPISpec represents a parsed OpenAPI specification
type OpenAPISpec struct {
	Version    string                // OpenAPI version (e.g., "3.0.0")
	Info       Info                  // API information
	Servers    []Server              // API servers
	Paths      map[string]PathItem   // API paths
	Components Components            // Reusable components
	Tags       []Tag                 // API tags
	Security   []map[string][]string // Security requirements
}

// Info contains API metadata
type Info struct {
	Title          string  // API title
	Description    string  // API description
	Version        string  // API version
	TermsOfService string  // Terms of service URL
	Contact        Contact // Contact information
	License        License // License information
}

// Contact information
type Contact struct {
	Name  string
	URL   string
	Email string
}

// License information
type License struct {
	Name string
	URL  string
}

// Server represents an API server
type Server struct {
	URL         string
	Description string
	Variables   map[string]ServerVariable
}

// ServerVariable is a server URL variable
type ServerVariable struct {
	Enum        []string
	Default     string
	Description string
}

// PathItem represents a path in the API
type PathItem struct {
	Summary     string
	Description string
	Get         *Operation
	Put         *Operation
	Post        *Operation
	Delete      *Operation
	Options     *Operation
	Head        *Operation
	Patch       *Operation
	Trace       *Operation
	Parameters  []Parameter
}

// Operation represents an API operation
type Operation struct {
	Tags        []string
	Summary     string
	Description string
	OperationID string
	Parameters  []Parameter
	RequestBody *RequestBody
	Responses   map[string]Response
	Security    []map[string][]string
	Deprecated  bool
}

// Parameter represents an operation parameter
type Parameter struct {
	Name            string
	In              string // "query", "header", "path", "cookie"
	Description     string
	Required        bool
	Deprecated      bool
	AllowEmptyValue bool
	Schema          Schema
}

// RequestBody represents an operation request body
type RequestBody struct {
	Description string
	Content     map[string]MediaType
	Required    bool
}

// Response represents an operation response
type Response struct {
	Description string
	Headers     map[string]Header
	Content     map[string]MediaType
}

// Header represents a response header
type Header struct {
	Description string
	Schema      Schema
}

// MediaType represents a content type
type MediaType struct {
	Schema   Schema
	Examples map[string]Example
}

// Example represents a request/response example
type Example struct {
	Summary       string
	Description   string
	Value         interface{}
	ExternalValue string
}

// Schema represents a data schema
type Schema struct {
	Type        string
	Format      string
	Items       *Schema
	Properties  map[string]Schema
	Required    []string
	Description string
	Default     interface{}
	Nullable    bool
	Ref         string // $ref reference
	AllOf       []Schema
	OneOf       []Schema
	AnyOf       []Schema
	Enum        []interface{}
	Example     interface{}
	MinLength   int
	MaxLength   int
	Minimum     float64
	Maximum     float64
}

// Components contains reusable components
type Components struct {
	Schemas         map[string]Schema
	Responses       map[string]Response
	Parameters      map[string]Parameter
	Examples        map[string]Example
	RequestBodies   map[string]RequestBody
	Headers         map[string]Header
	SecuritySchemes map[string]SecurityScheme
}

// SecurityScheme represents a security scheme
type SecurityScheme struct {
	Type             string
	Description      string
	Name             string
	In               string
	Scheme           string
	BearerFormat     string
	Flows            OAuthFlows
	OpenIDConnectURL string
}

// OAuthFlows represents OAuth flow configurations
type OAuthFlows struct {
	Implicit          OAuthFlow
	Password          OAuthFlow
	ClientCredentials OAuthFlow
	AuthorizationCode OAuthFlow
}

// OAuthFlow represents an OAuth flow
type OAuthFlow struct {
	AuthorizationURL string
	TokenURL         string
	RefreshURL       string
	Scopes           map[string]string
}

// Tag represents an API tag
type Tag struct {
	Name        string
	Description string
}

// rawOpenAPISpec is used for unmarshaling the YAML
type rawOpenAPISpec struct {
	OpenAPI    string                `yaml:"openapi"`
	Info       Info                  `yaml:"info"`
	Servers    []Server              `yaml:"servers"`
	Paths      map[string]yaml.Node  `yaml:"paths"`
	Components yaml.Node             `yaml:"components"`
	Tags       []Tag                 `yaml:"tags"`
	Security   []map[string][]string `yaml:"security"`
}

// ParseFile parses an OpenAPI YAML file
func ParseFile(filename string) (*OpenAPISpec, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return Parse(data)
}

// Parse parses OpenAPI YAML data
func Parse(data []byte) (*OpenAPISpec, error) {
	var raw rawOpenAPISpec
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	// Convert raw spec to our structured format
	spec := &OpenAPISpec{
		Version:  raw.OpenAPI,
		Info:     raw.Info,
		Servers:  raw.Servers,
		Paths:    make(map[string]PathItem),
		Tags:     raw.Tags,
		Security: raw.Security,
	}

	// Parse paths
	for path, node := range raw.Paths {
		pathItem, err := parsePathItem(node)
		if err != nil {
			return nil, fmt.Errorf("error parsing path %s: %w", path, err)
		}
		spec.Paths[path] = pathItem
	}

	// Parse components
	components, err := parseComponents(raw.Components)
	if err != nil {
		return nil, fmt.Errorf("error parsing components: %w", err)
	}
	spec.Components = components

	return spec, nil
}

// parsePathItem parses a path item node
func parsePathItem(node yaml.Node) (PathItem, error) {
	var pathItem PathItem

	// This is a placeholder for the actual parsing logic
	// In a complete implementation, you would parse each field of the path item
	// by navigating the YAML node structure

	return pathItem, nil
}

// parseComponents parses the components section
func parseComponents(node yaml.Node) (Components, error) {
	var components Components

	// This is a placeholder for the actual parsing logic
	// In a complete implementation, you would parse each component type
	// by navigating the YAML node structure

	return components, nil
}

// GetGoType converts an OpenAPI type to a Go type
func GetGoType(schema Schema) string {
	// Handle $ref references
	if schema.Ref != "" {
		// Extract the type name from the reference
		parts := strings.Split(schema.Ref, "/")
		return parts[len(parts)-1]
	}

	switch schema.Type {
	case "string":
		switch schema.Format {
		case "date-time":
			return "time.Time"
		case "date":
			return "time.Time"
		case "time":
			return "time.Time"
		case "byte":
			return "[]byte"
		case "binary":
			return "[]byte"
		default:
			return "string"
		}
	case "number":
		switch schema.Format {
		case "float":
			return "float32"
		case "double":
			return "float64"
		default:
			return "float64"
		}
	case "integer":
		switch schema.Format {
		case "int32":
			return "int32"
		case "int64":
			return "int64"
		default:
			return "int"
		}
	case "boolean":
		return "bool"
	case "array":
		if schema.Items != nil {
			itemType := GetGoType(*schema.Items)
			return "[]" + itemType
		}
		return "[]interface{}"
	case "object":
		// For objects without specific properties
		if len(schema.Properties) == 0 {
			return "map[string]interface{}"
		}
		// For objects with properties, this would typically be handled by
		// generating a struct elsewhere
		return "struct{}"
	default:
		return "interface{}"
	}
}
