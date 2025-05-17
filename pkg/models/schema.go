package models

import (
	"strings"
)

// Schema represents an OpenAPI schema object
type Schema struct {
	Type       string              `json:"type,omitempty" yaml:"type,omitempty"`
	Format     string              `json:"format,omitempty" yaml:"format,omitempty"`
	Required   []string            `json:"required,omitempty" yaml:"required,omitempty"`
	Properties map[string]*Schema  `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items      *Schema             `json:"items,omitempty" yaml:"items,omitempty"`
	Ref        string              `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	Enum       []interface{}       `json:"enum,omitempty" yaml:"enum,omitempty"`
	AllOf      []*Schema           `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	OneOf      []*Schema           `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOf      []*Schema           `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
}

// GetGoType returns the Go type for a schema
func (s *Schema) GetGoType() string {
	if s.Ref != "" {
		// Extract type name from reference
		// e.g., "#/components/schemas/Pet" -> "Pet"
		// This is simplified and would need more logic in a real implementation
		return "models." + s.Ref[strings.LastIndex(s.Ref, "/")+1:]
	}

	switch s.Type {
	case "string":
		switch s.Format {
		case "date-time":
			return "time.Time"
		case "binary":
			return "[]byte"
		default:
			return "string"
		}
	case "integer":
		switch s.Format {
		case "int64":
			return "int64"
		default:
			return "int"
		}
	case "number":
		switch s.Format {
		case "float":
			return "float32"
		case "double":
			return "float64"
		default:
			return "float64"
		}
	case "boolean":
		return "bool"
	case "array":
		if s.Items != nil {
			return "[]" + s.Items.GetGoType()
		}
		return "[]interface{}"
	case "object":
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}
