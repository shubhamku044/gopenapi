package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/shubhamku044/gopenapi/internal/models"
	"gopkg.in/yaml.v3"
)

// ParseSpecFile parses an OpenAPI specification file (YAML or JSON)
func ParseSpecFile(filePath string) (*models.OpenAPISpec, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var spec models.OpenAPISpec
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".yaml" || ext == ".yml" {
		err = yaml.Unmarshal(data, &spec)
	} else if ext == ".json" {
		err = json.Unmarshal(data, &spec)
	} else {
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return nil, err
	}

	// Process the spec to add derived fields
	ProcessSpec(&spec)

	return &spec, nil
}

// ProcessSpec processes the OpenAPI spec to add derived fields
func ProcessSpec(spec *models.OpenAPISpec) {
	// Add the HTTP method to each operation
	for path, methods := range spec.Paths {
		for method, op := range methods {
			op.Method = strings.ToUpper(method)

			// Add default tags if none are provided
			if len(op.Tags) == 0 {
				op.Tags = []string{"default"}
			}

			// Handle missing operation IDs
			if op.OperationID == "" {
				// Generate a default operationID based on method and path
				pathParts := strings.Split(strings.Trim(path, "/"), "/")
				var name string
				if len(pathParts) > 0 {
					name = pathParts[len(pathParts)-1]
				} else {
					name = "root"
				}
				op.OperationID = strings.ToLower(method) + ToCamelCase(name)
			}

			// Update the operation in the map
			methods[method] = op
		}
	}
}

// ToCamelCase converts a string to CamelCase
func ToCamelCase(s string) string {
	// Convert snake_case or kebab-case to CamelCase
	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")
	for i := range parts {
		if i == 0 {
			parts[i] = strings.Title(parts[i])
		} else {
			parts[i] = strings.Title(parts[i])
		}
	}
	return strings.Join(parts, "")
}
