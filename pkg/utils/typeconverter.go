package utils

import (
	"strings"

	"github.com/shubhamku044/gopenapi/internal/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// GetGoType converts an OpenAPI schema to a Go type
func GetGoType(schema models.Schema) string {
	// Handle $ref
	if schema.Ref != "" {
		// Extract the model name from the reference
		parts := strings.Split(schema.Ref, "/")
		return parts[len(parts)-1]
	}

	// Handle different types
	switch schema.Type {
	case "integer":
		switch schema.Format {
		case "int32":
			return "int32"
		case "int64":
			return "int64"
		default:
			return "int"
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
	case "boolean":
		return "bool"
	case "string":
		switch schema.Format {
		case "byte":
			return "[]byte"
		case "binary":
			return "[]byte"
		case "date", "date-time":
			return "time.Time"
		default:
			return "string"
		}
	case "array":
		if schema.Items != nil {
			return "[]" + GetGoType(*schema.Items)
		}
		return "[]interface{}"
	case "object":
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}

// ConvertPathToGin converts an OpenAPI path to a Gin path
func ConvertPathToGin(path string) string {
	// Convert OpenAPI path params {param} to Gin format :param
	return strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", "")
}

// ToCamelCase converts a string to CamelCase
func ToCamelCase(s string) string {
	// Convert snake_case or kebab-case to CamelCase
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")

	// Use cases.Title from golang.org/x/text/cases instead of deprecated strings.Title
	title := cases.Title(language.English)
	for i := range parts {
		if len(parts[i]) > 0 {
			parts[i] = title.String(parts[i])
		}
	}
	return strings.Join(parts, "")
}
