package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Generator handles the parsing of gopenapi specs and code generation
type Generator struct {
	SpecPath   string
	OutputPath string
	Spec       map[string]interface{}
}

// NewGenerator creates a new generator instance
func NewGenerator(specPath, outputPath string) (*Generator, error) {
	return &Generator{
		SpecPath:   specPath,
		OutputPath: outputPath,
	}, nil
}

// LoadSpec loads and parses the gopenapi spec file
func (g *Generator) LoadSpec() error {
	data, err := os.ReadFile(g.SpecPath)
	if err != nil {
		return fmt.Errorf("failed to read spec file: %w", err)
	}

	// Determine if the file is JSON or YAML based on extension
	ext := filepath.Ext(g.SpecPath)
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &g.Spec); err != nil {
			return fmt.Errorf("failed to parse JSON spec: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &g.Spec); err != nil {
			return fmt.Errorf("failed to parse YAML spec: %w", err)
		}
	default:
		return fmt.Errorf("unsupported spec file format: %s", ext)
	}

	// Check if this is a gopenapi spec or an OpenAPI spec
	_, isGopenapi := g.Spec["gopenapi"]
	_, isOpenapi := g.Spec["openapi"]

	// If it's a gopenapi spec, validate required sections
	if isGopenapi {
		_, hasModels := g.Spec["models"]
		if !hasModels {
			return fmt.Errorf("invalid gopenapi spec: missing 'models' section")
		}
	} else if isOpenapi {
		// It's an OpenAPI spec, check for components
		components, hasComponents := g.Spec["components"].(map[string]interface{})
		if !hasComponents {
			return fmt.Errorf("invalid OpenAPI spec: missing 'components' section")
		}
		
		_, hasSchemas := components["schemas"].(map[string]interface{})
		if !hasSchemas {
			return fmt.Errorf("invalid OpenAPI spec: missing 'schemas' in components")
		}
	} else {
		return fmt.Errorf("invalid spec: neither gopenapi nor openapi version field found")
	}

	return nil
}

// ModelTemplate is the template for generating model structs
const ModelTemplate = `package {{.PackageName}}

{{range .Imports}}import "{{.}}"
{{end}}

{{range .Models}}
// {{.Name}} {{.Description}}
type {{.Name}} struct {
{{range .Properties}}    {{.Name}} {{.Type}} ` + "`json:\"{{.JSONName}},omitempty\" yaml:\"{{.JSONName}},omitempty\"{{if .Required}} validate:\"required\"{{end}}`" + `
{{end}}
}
{{end}}
`

// ModelData represents the data for the model template
type ModelData struct {
	PackageName string
	Imports     []string
	Models      []ModelInfo
}

// ModelInfo represents a model to be generated
type ModelInfo struct {
	Name        string
	Description string
	Properties  []PropertyInfo
}

// PropertyInfo represents a property of a model
type PropertyInfo struct {
	Name     string
	Type     string
	JSONName string
	Required bool
}

// GenerateModels generates model structs from model definitions
func (g *Generator) GenerateModels() error {
	// Ensure output directory exists
	modelsDir := filepath.Join(g.OutputPath, "models")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	// Get models from the gopenapi spec
	models, err := g.getSchemas()
	if err != nil {
		return err
	}

	// Generate model data
	modelData := ModelData{
		PackageName: "models",
		Imports:     []string{"time"},
		Models:      []ModelInfo{},
	}

	// Process each model
	for _, modelInterface := range models {
		modelMap, ok := modelInterface.(map[string]interface{})
		if !ok {
			return fmt.Errorf("model is not a valid object")
		}
		
		name, ok := modelMap["name"].(string)
		if !ok {
			return fmt.Errorf("model name is not a string")
		}
		
		model, err := g.processSchema(name, modelMap)
		if err != nil {
			return err
		}
		modelData.Models = append(modelData.Models, model)
	}

	// Create the template
	tmpl, err := template.New("model").Parse(ModelTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Create the output file
	outputFile := filepath.Join(modelsDir, "models.go")
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Execute the template
	if err := tmpl.Execute(file, modelData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// getSchemas extracts the models from the gopenapi spec
func (g *Generator) getSchemas() ([]interface{}, error) {
	// Extract models from the gopenapi spec
	models, ok := g.Spec["models"].([]interface{})
	if !ok {
		// Try the old OpenAPI format as a fallback
		components, ok := g.Spec["components"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("no models found in spec")
		}
		
		schemas, ok := components["schemas"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("no schemas found in components")
		}
		
		// Convert map to array for compatibility
		modelArray := make([]interface{}, 0, len(schemas))
		for name, schema := range schemas {
			schemaMap, ok := schema.(map[string]interface{})
			if !ok {
				continue
			}
			
			// Create a model in our custom format
			model := map[string]interface{}{
				"name": name,
				"description": "Converted from OpenAPI schema",
			}
			
			// Copy required properties
			if required, ok := schemaMap["required"].([]interface{}); ok {
				model["required"] = required
			}
			
			// Convert properties
			if props, ok := schemaMap["properties"].(map[string]interface{}); ok {
				properties := make([]interface{}, 0, len(props))
				for propName, propSchema := range props {
					propMap, ok := propSchema.(map[string]interface{})
					if !ok {
						continue
					}
					
					// Create property in our custom format
					property := map[string]interface{}{
						"name": propName,
					}
					
					// Copy type and format
					if propType, ok := propMap["type"].(string); ok {
						property["type"] = propType
					}
					if format, ok := propMap["format"].(string); ok {
						property["format"] = format
					}
					if enum, ok := propMap["enum"].([]interface{}); ok {
						property["enum"] = enum
					}
					
					properties = append(properties, property)
				}
				model["properties"] = properties
			}
			
			modelArray = append(modelArray, model)
		}
		
		return modelArray, nil
	}

	return models, nil
}

// processSchema converts a gopenapi model to a ModelInfo
func (g *Generator) processSchema(name string, model map[string]interface{}) (ModelInfo, error) {
	description := "Generated from gopenapi model"
	if desc, ok := model["description"].(string); ok {
		description = desc
	}

	modelInfo := ModelInfo{
		Name:        name,
		Description: description,
		Properties:  []PropertyInfo{},
	}

	// Get required properties
	requiredProps := []string{}
	if required, ok := model["required"].([]interface{}); ok {
		for _, req := range required {
			if reqStr, ok := req.(string); ok {
				requiredProps = append(requiredProps, reqStr)
			}
		}
	}

	// Process properties
	if properties, ok := model["properties"].([]interface{}); ok {
		for _, propInterface := range properties {
			propMap, ok := propInterface.(map[string]interface{})
			if !ok {
				continue
			}

			propName, ok := propMap["name"].(string)
			if !ok {
				continue
			}

			// Check if property is required
			isRequired := false
			for _, req := range requiredProps {
				if req == propName {
					isRequired = true
					break
				}
			}

			// Convert property name to Go style (camelCase)
			goName := toCamelCase(propName)

			// Get property type
			goType := g.getGoTypeFromProperty(propMap)

			// Add property to model
			modelInfo.Properties = append(modelInfo.Properties, PropertyInfo{
				Name:     goName,
				Type:     goType,
				JSONName: propName,
				Required: isRequired,
			})
		}
	}

	return modelInfo, nil
}

// getGoTypeFromProperty converts a gopenapi property type to a Go type
func (g *Generator) getGoTypeFromProperty(property map[string]interface{}) string {
	// Get type
	typeStr, _ := property["type"].(string)
	format, _ := property["format"].(string)
	
	// Check for reference to another model
	if typeStr == "" {
		if items, ok := property["items"].(string); ok {
			return "[]" + items
		}
	}

	switch typeStr {
	case "string":
		switch format {
		case "date-time":
			return "time.Time"
		case "binary":
			return "[]byte"
		default:
			return "string"
		}
	case "integer":
		switch format {
		case "int64":
			return "int64"
		default:
			return "int"
		}
	case "number":
		switch format {
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
		// Check for array type with items
		if items, ok := property["items"].(string); ok {
			return "[]" + items
		}
		return "[]interface{}"
	case "object":
		return "map[string]interface{}"
	default:
		return "interface{}"
	}
}

// toCamelCase converts a string to camelCase
func toCamelCase(s string) string {
	// Split the string by non-alphanumeric characters
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", r)
	})

	// Capitalize each part except the first one
	for i := 0; i < len(parts); i++ {
		if i == 0 {
			parts[i] = strings.ToLower(parts[i])
		} else {
			parts[i] = strings.Title(parts[i])
		}
	}

	return strings.Join(parts, "")
}

// GenerateServer generates server handlers and routing
func (g *Generator) GenerateServer() error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.OutputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get endpoints from the gopenapi spec
	endpoints, ok := g.Spec["endpoints"].([]interface{})
	if !ok {
		// For now, just return nil since server generation is optional
		return nil
	}

	// TODO: Implement server code generation with endpoints
	_ = endpoints // Use endpoints to avoid unused variable warning
	return nil
}

// GenerateClient generates HTTP client code
func (g *Generator) GenerateClient() error {
	// Ensure output directory exists
	clientDir := filepath.Join(g.OutputPath, "client")
	if err := os.MkdirAll(clientDir, 0755); err != nil {
		return fmt.Errorf("failed to create client directory: %w", err)
	}

	// Get endpoints from the gopenapi spec
	endpoints, ok := g.Spec["endpoints"].([]interface{})
	if !ok {
		// For now, just return nil since client generation is optional
		return nil
	}

	// TODO: Implement client code generation with endpoints
	_ = endpoints // Use endpoints to avoid unused variable warning
	return nil
}

// GenerateAll generates all code components
func (g *Generator) GenerateAll() error {
	if err := g.GenerateModels(); err != nil {
		return err
	}

	if err := g.GenerateServer(); err != nil {
		return err
	}

	if err := g.GenerateClient(); err != nil {
		return err
	}

	return nil
}
