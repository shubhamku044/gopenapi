package generator

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/shubhamku044/gopenapi/internal/models"
	"github.com/shubhamku044/gopenapi/pkg/utils"
)

// GenerateModels generates model files from OpenAPI schemas
func GenerateModels(spec *models.OpenAPISpec, baseDir string) error {
	// For MVP, we'll create a simple models file with placeholders
	// A complete implementation would parse schemas from spec.Components.Schemas
	modelsTemplate := `package models

// This file contains the data models for the API
// In a real implementation, these would be fully generated from the schema definitions

// ResponseMessage is a simple response model
type ResponseMessage struct {
	Message string ` + "`json:\"message\"`" + `
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Code    int    ` + "`json:\"code\"`" + `
	Message string ` + "`json:\"message\"`" + `
}

{{range $name, $schema := .Schemas}}
// {{$name}} represents a {{$name}} model
type {{$name}} struct {
	{{range $propName, $propSchema := $schema.Properties}}
	{{toCamelCase $propName}} {{getGoType $propSchema}} ` + "`json:\"{{$propName}}\"`" + `
	{{end}}
}
{{end}}
`

	// Create functions for the template
	funcMap := template.FuncMap{
		"toCamelCase": utils.ToCamelCase,
		"getGoType":   func(s models.Schema) string { return utils.GetGoType(s) },
	}

	tmpl, err := template.New("models").Funcs(funcMap).Parse(modelsTemplate)
	if err != nil {
		return err
	}

	data := struct {
		Schemas map[string]models.Schema
	}{
		Schemas: spec.Components.Schemas,
	}

	f, err := os.Create(filepath.Join(baseDir, "models", "models.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}