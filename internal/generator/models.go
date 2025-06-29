package generator

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/shubhamku044/gopenapi/internal/models"
	"github.com/shubhamku044/gopenapi/pkg/utils"
)

func GenerateModels(spec *models.OpenAPISpec, baseDir string) error {
	needsTimeImport := false
	for _, schema := range spec.Components.Schemas {
		if hasTimeFields(schema) {
			needsTimeImport = true
			break
		}
	}

	imports := ""
	if needsTimeImport {
		imports = `import (
	"time"
)`
	}

	modelsTemplate := `package models

` + imports + `

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

func hasTimeFields(schema models.Schema) bool {
	if schema.Type == "string" && (schema.Format == "date" || schema.Format == "date-time") {
		return true
	}

	for _, propSchema := range schema.Properties {
		if hasTimeFields(propSchema) {
			return true
		}
	}

	if schema.Items != nil && hasTimeFields(*schema.Items) {
		return true
	}

	return false
}
