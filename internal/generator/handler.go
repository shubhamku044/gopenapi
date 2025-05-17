package generator

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/shubhamku044/gopenapi/internal/parser"
)

// HandlerGenerator generates Go handler functions from OpenAPI paths
type HandlerGenerator struct {
	Config Config
	Models map[string]Model
}

// HandlerFile represents a generated handler file
type HandlerFile struct {
	Name        string
	PackageName string
	Imports     []string
	Handlers    []Handler
}

// Handler represents a handler function
type Handler struct {
	Name        string
	Description string
	Method      string
	Path        string
	OperationID string
	Parameters  []Parameter
	RequestBody *RequestBodyDef
	Responses   map[string]Response
	Tags        []string
}

// Parameter represents a handler parameter
type Parameter struct {
	Name        string
	Type        string
	Location    string // "query", "path", "header", "cookie"
	Description string
	Required    bool
}

// RequestBodyDef represents a request body definition
type RequestBodyDef struct {
	Type        string
	ContentType string
	Required    bool
}

// Response represents a handler response
type Response struct {
	Code        string
	Type        string
	ContentType string
	Description string
}

// NewHandlerGenerator creates a new handler generator
func NewHandlerGenerator(config Config, models map[string]Model) *HandlerGenerator {
	return &HandlerGenerator{
		Config: config,
		Models: models,
	}
}

// GenerateHandlers generates Go handler functions from OpenAPI paths
func (g *HandlerGenerator) GenerateHandlers(spec *parser.OpenAPISpec) ([]OutputFile, error) {
	var files []OutputFile

	// Group handlers by tag (one file per tag)
	handlersByTag := make(map[string][]Handler)

	// Process each path
	for pathURL, pathItem := range spec.Paths {
		// Process each operation (GET, POST, etc.)
		operations := map[string]*parser.Operation{
			"GET":     pathItem.Get,
			"POST":    pathItem.Post,
			"PUT":     pathItem.Put,
			"DELETE":  pathItem.Delete,
			"OPTIONS": pathItem.Options,
			"HEAD":    pathItem.Head,
			"PATCH":   pathItem.Patch,
			"TRACE":   pathItem.Trace,
		}

		for method, operation := range operations {
			if operation == nil {
				continue
			}

			handler, err := g.generateHandler(pathURL, method, operation)
			if err != nil {
				return nil, fmt.Errorf("error generating handler for %s %s: %w", method, pathURL, err)
			}

			// Determine tag for grouping
			tag := "default"
			if len(operation.Tags) > 0 {
				tag = operation.Tags[0]
			} else {
				// If no tag is specified, use the first path segment
				segments := strings.Split(strings.Trim(pathURL, "/"), "/")
				if len(segments) > 0 {
					tag = segments[0]
				}
			}

			// Group handlers by tag
			handlersByTag[tag] = append(handlersByTag[tag], handler)
		}
	}

	// Create a file for each group of handlers
	for tag, handlers := range handlersByTag {
		handlerFile := HandlerFile{
			Name:        sanitizeFileName(tag) + "_handlers.go",
			PackageName: "handlers",
			Handlers:    handlers,
		}

		// Add imports based on handler types
		handlerFile.Imports = g.determineImports(handlers)

		// Generate the file content
		content, err := g.renderHandlerFile(handlerFile)
		if err != nil {
			return nil, fmt.Errorf("error rendering handler file for tag %s: %w", tag, err)
		}

		files = append(files, OutputFile{
			Path:    "api/handlers/" + handlerFile.Name,
			Content: content,
		})
	}

	return files, nil
}

// generateHandler generates a handler from an OpenAPI operation
func (g *HandlerGenerator) generateHandler(pathURL, method string, operation *parser.Operation) (Handler, error) {
	// Generate a handler name from the operation ID or path
	name := operation.OperationID
	if name == "" {
		// Generate a name from the method and path
		name = method + toCamelCase(strings.ReplaceAll(pathURL, "/", "_"))
	}

	handler := Handler{
		Name:        formatHandlerName(name),
		Description: operation.Description,
		Method:      method,
		Path:        pathURL,
		OperationID: operation.OperationID,
		Parameters:  []Parameter{},
		Responses:   make(map[string]Response),
		Tags:        operation.Tags,
	}

	// Process parameters
	for _, param := range operation.Parameters {
		parameter := Parameter{
			Name:        param.Name,
			Type:        parser.GetGoType(param.Schema),
			Location:    param.In,
			Description: param.Description,
			Required:    param.Required,
		}
		handler.Parameters = append(handler.Parameters, parameter)
	}

	// Process request body
	if operation.RequestBody != nil {
		for contentType, mediaType := range operation.RequestBody.Content {
			// Use the first content type (usually application/json)
			handler.RequestBody = &RequestBodyDef{
				Type:        parser.GetGoType(mediaType.Schema),
				ContentType: contentType,
				Required:    operation.RequestBody.Required,
			}
			break
		}
	}

	// Process responses
	for code, response := range operation.Responses {
		for contentType, mediaType := range response.Content {
			// Use the first content type (usually application/json)
			handler.Responses[code] = Response{
				Code:        code,
				Type:        parser.GetGoType(mediaType.Schema),
				ContentType: contentType,
				Description: response.Description,
			}
			break
		}
	}

	return handler, nil
}

// determineImports determines necessary imports for a handler file
func (g *HandlerGenerator) determineImports(handlers []Handler) []string {
	importsMap := make(map[string]bool)

	// Standard imports
	importsMap["net/http"] = true
	importsMap["encoding/json"] = true

	// Add models package if using any models
	hasModels := false
	for _, handler := range handlers {
		if handler.RequestBody != nil && !isBuiltinType(handler.RequestBody.Type) {
			hasModels = true
		}
		for _, response := range handler.Responses {
			if !isBuiltinType(response.Type) {
				hasModels = true
			}
		}
	}

	if hasModels {
		importsMap[g.Config.PackageName+"/api/models"] = true
	}

	// Convert map to slice
	var imports []string
	for imp := range importsMap {
		imports = append(imports, imp)
	}

	return imports
}

// renderHandlerFile renders a handler file using templates
func (g *HandlerGenerator) renderHandlerFile(handlerFile HandlerFile) (string, error) {
	// Use a simple template for now
	tmplText := `// Code generated by gopenapi; DO NOT EDIT.
package {{.PackageName}}

import (
{{- range .Imports}}
	"{{.}}"
{{- end}}
)

{{range .Handlers}}
// {{.Name}} {{.Description}}
// {{.Method}} {{.Path}}
func {{.Name}}(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement handler logic
	{{if .RequestBody}}
	// Parse request body
	var requestBody {{.RequestBody.Type}}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	{{end}}

	{{if index .Responses "200"}}
	// Return successful response
	w.Header().Set("Content-Type", "{{(index .Responses "200").ContentType}}")
	w.WriteHeader(http.StatusOK)
	{{if ne (index .Responses "200").Type "string"}}
	// TODO: Populate response data
	response := {{(index .Responses "200").Type}}{}
	json.NewEncoder(w).Encode(response)
	{{else}}
	w.Write([]byte("Success"))
	{{end}}
	{{else}}
	w.WriteHeader(http.StatusNoContent)
	{{end}}
}
{{end}}
`

	tmpl, err := template.New("handler").Parse(tmplText)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, handlerFile); err != nil {
		return "", fmt.Errorf("error executing template: %w", err)
	}

	return buf.String(), nil
}

// formatHandlerName formats an operation name as a Go handler function name
func formatHandlerName(name string) string {
	// Ensure it starts with a capital letter for export
	if len(name) > 0 {
		return strings.ToUpper(name[0:1]) + name[1:]
	}
	return "Handler"
}

// sanitizeFileName sanitizes a tag name for use as a filename
func sanitizeFileName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}

// isBuiltinType checks if a type is a Go built-in type
func isBuiltinType(typeName string) bool {
	builtins := map[string]bool{
		"string":                 true,
		"int":                    true,
		"int32":                  true,
		"int64":                  true,
		"float32":                true,
		"float64":                true,
		"bool":                   true,
		"interface{}":            true,
		"map[string]interface{}": true,
	}

	// Also check for arrays of built-in types
	if strings.HasPrefix(typeName, "[]") {
		return isBuiltinType(typeName[2:])
	}

	return builtins[typeName]
}
