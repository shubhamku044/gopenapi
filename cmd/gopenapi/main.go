package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// OpenAPISpec represents a simplified OpenAPI specification
type OpenAPISpec struct {
	Info struct {
		Title       string `json:"title" yaml:"title"`
		Version     string `json:"version" yaml:"version"`
		Description string `json:"description" yaml:"description"`
	} `json:"info" yaml:"info"`
	Paths map[string]map[string]Operation `json:"paths" yaml:"paths"`
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
	Type       string            `json:"type" yaml:"type"`
	Format     string            `json:"format" yaml:"format"`
	Properties map[string]Schema `json:"properties" yaml:"properties"`
	Ref        string            `json:"$ref" yaml:"$ref"`
}

func main() {
	specFile := flag.String("spec", "", "Path to OpenAPI specification file (YAML or JSON)")
	outputDir := flag.String("output", "./generated", "Output directory for generated code")
	flag.Parse()

	if *specFile == "" {
		log.Fatal("Please provide an OpenAPI specification file with --spec")
	}

	// Parse the OpenAPI specification
	spec, err := parseSpecFile(*specFile)
	if err != nil {
		log.Fatalf("Failed to parse OpenAPI specification: %v", err)
	}

	// Create the output directory structure
	err = createDirectories(*outputDir)
	if err != nil {
		log.Fatalf("Failed to create output directories: %v", err)
	}

	// Generate the code
	err = generateCode(spec, *outputDir)
	if err != nil {
		log.Fatalf("Failed to generate code: %v", err)
	}

	fmt.Printf("Code successfully generated in %s\n", *outputDir)
}

func parseSpecFile(filePath string) (*OpenAPISpec, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var spec OpenAPISpec
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

	return &spec, nil
}

func createDirectories(baseDir string) error {
	dirs := []string{
		baseDir,
		filepath.Join(baseDir, "api"),
		filepath.Join(baseDir, "models"),
		filepath.Join(baseDir, "server"),
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateCode(spec *OpenAPISpec, baseDir string) error {
	// Generate server
	err := generateServerFile(spec, baseDir)
	if err != nil {
		return err
	}

	// Generate API interfaces
	err = generateAPIFile(spec, baseDir)
	if err != nil {
		return err
	}

	// Generate models
	err = generateModels(spec, baseDir)
	if err != nil {
		return err
	}

	// Generate README.md with usage examples
	err = generateReadme(spec, baseDir)
	if err != nil {
		return err
	}

	return nil
}

func generateServerFile(spec *OpenAPISpec, baseDir string) error {
	serverTemplate := `package server

import (
	"github.com/gin-gonic/gin"
	"{{.PackageName}}/api"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	api    api.API
}

// NewServer creates a new API server
func NewServer(api api.API) *Server {
	s := &Server{
		router: gin.Default(),
		api:    api,
	}
	s.setupRoutes()
	return s
}

// Start starts the server on the specified address
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

// GetRouter returns the Gin router instance
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
{{range .Routes}}
	s.router.{{.Method}}("{{.Path}}", func(c *gin.Context) {
		{{if .HasPathParams}}
		{{range .PathParams}}
		{{.Name}} := c.Param("{{.Name}}")
		{{end}}
		{{end}}
		s.api.{{.HandlerName}}(c{{if .HasPathParams}}, {{range $index, $param := .PathParams}}{{if $index}}, {{end}}{{.Name}}{{end}}{{end}})
	})
{{end}}
}
`

	tmpl, err := template.New("server").Parse(serverTemplate)
	if err != nil {
		return err
	}

	// Use a proper import path for the generated code
	packageName := "github.com/shubhamku044/gopenapi/generated"

	// Prepare route data
	type PathParam struct {
		Name string
		Type string
	}

	type Route struct {
		Method        string
		Path          string
		HandlerName   string
		PathParams    []PathParam
		HasPathParams bool
	}

	var routes []Route
	for path, operations := range spec.Paths {
		for method, op := range operations {
			// Set the HTTP method in the Operation struct
			op.Method = strings.Title(strings.ToLower(method))
			
			ginPath := convertPathToGin(path)
			handlerName := toCamelCase(op.OperationID)

			// Extract path parameters
			var pathParams []PathParam
			for _, param := range op.Parameters {
				if param.In == "path" {
					pathParams = append(pathParams, PathParam{
						Name: param.Name,
						Type: getGoType(param.Schema),
					})
				}
			}

			routes = append(routes, Route{
				Method:        strings.ToUpper(method),
				Path:          ginPath,
				HandlerName:   handlerName,
				PathParams:    pathParams,
				HasPathParams: len(pathParams) > 0,
			})
		}
	}

	data := struct {
		PackageName string
		Routes      []Route
	}{
		PackageName: packageName,
		Routes:      routes,
	}

	f, err := os.Create(filepath.Join(baseDir, "server", "server.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

func generateAPIFile(spec *OpenAPISpec, baseDir string) error {
	apiTemplate := `package api

import (
	"github.com/gin-gonic/gin"
)

// API defines the interface for all API operations
type API interface {
	{{range .Handlers}}
	// {{.HandlerName}} handles {{.Method}} {{.Path}}
	// {{.Summary}}
	{{.HandlerName}}(c *gin.Context{{if .HasPathParams}}, {{range $index, $param := .PathParams}}{{if $index}}, {{end}}{{.Name}} {{.Type}}{{end}}{{end}})
	{{end}}
}

// DefaultAPI provides a default implementation of the API interface
type DefaultAPI struct {}

// NewAPI creates a new API handler
func NewAPI() API {
	return &DefaultAPI{}
}

{{range .Handlers}}
// {{.HandlerName}} handles {{.Method}} {{.Path}}
func (a *DefaultAPI) {{.HandlerName}}(c *gin.Context{{if .HasPathParams}}, {{range $index, $param := .PathParams}}{{if $index}}, {{end}}{{.Name}} {{.Type}}{{end}}{{end}}) {
	// TODO: Implement me
	c.JSON(200, gin.H{"message": "Not implemented"})
}
{{end}}
`

	tmpl, err := template.New("api").Parse(apiTemplate)
	if err != nil {
		return err
	}

	// Transform the operations for template rendering
	type PathParam struct {
		Name string
		Type string
	}

	type Handler struct {
		Method        string
		Path          string
		HandlerName   string
		Summary       string
		PathParams    []PathParam
		HasPathParams bool
	}

	var handlers []Handler
	for path, operations := range spec.Paths {
		for method, op := range operations {
			// Extract path parameters
			var pathParams []PathParam
			for _, param := range op.Parameters {
				if param.In == "path" {
					pathParams = append(pathParams, PathParam{
						Name: param.Name,
						Type: getGoType(param.Schema),
					})
				}
			}

			handlers = append(handlers, Handler{
				Method:        method,
				Path:          path,
				HandlerName:   toCamelCase(op.OperationID),
				Summary:       op.Summary,
				PathParams:    pathParams,
				HasPathParams: len(pathParams) > 0,
			})
		}
	}

	data := struct {
		Handlers []Handler
	}{
		Handlers: handlers,
	}

	f, err := os.Create(filepath.Join(baseDir, "api", "api.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

func generateModels(spec *OpenAPISpec, baseDir string) error {
	// In a real implementation, you would extract models from spec.Components.Schemas
	// For this MVP, we'll just create a placeholder models file
	modelsContent := `package models

// This file contains the data models for the API
// In a real implementation, these would be generated based on the schema definitions

// ResponseMessage is a simple response model
type ResponseMessage struct {
	Message string ` + "`json:\"message\"`" + `
}
`

	return ioutil.WriteFile(filepath.Join(baseDir, "models", "models.go"), []byte(modelsContent), 0644)
}

func generateReadme(spec *OpenAPISpec, baseDir string) error {
	readmeTemplate := `# Generated API for {{.Title}}

This code was generated using GopenAPI for the {{.Title}} API (version {{.Version}}).

## Usage

To use this generated API in your project:

### 1. Implement the API interface

Create a custom implementation of the API interface:

` + "```go" + `
package main

import (
	"github.com/gin-gonic/gin"
	
	"yourmodule/generated/api"
	"yourmodule/generated/server"
)

// CustomAPI implements the generated API interface
type CustomAPI struct {
	// Add your dependencies here (database connections, etc.)
}

// Ensure CustomAPI implements the API interface
var _ api.API = (*CustomAPI)(nil)

// Example implementation of a generated method
func (a *CustomAPI) ListUsers(c *gin.Context) {
	// Your implementation here
	users := []map[string]interface{}{
		{"id": "1", "name": "John Doe"},
		{"id": "2", "name": "Jane Smith"},
	}
	c.JSON(200, users)
}

// Implement all other methods defined in the API interface...
` + "```" + `

### 2. Create and run the server

` + "```go" + `
package main

import (
	"log"
	
	"yourmodule/generated/server"
)

func main() {
	// Create your API implementation
	customAPI := &CustomAPI{}
	
	// Create the server with your API implementation
	srv := server.NewServer(customAPI)
	
	// Start the server
	log.Println("Starting server on :8080")
	if err := srv.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
` + "```" + `

### 3. Or use with an existing Gin application

` + "```go" + `
package main

import (
	"github.com/gin-gonic/gin"
	
	"yourmodule/generated/server"
)

func main() {
	// Create your existing Gin router
	router := gin.Default()
	
	// Add some custom middleware or routes
	router.Use(customMiddleware())
	router.GET("/health", healthCheckHandler)
	
	// Create your API implementation
	customAPI := &CustomAPI{}
	
	// Create the server with your API implementation
	srv := server.NewServer(customAPI)
	
	// Get the Gin router from the server and add it to your existing router
	apiRouter := srv.GetRouter()
	
	// Use the API routes in your application
	router.Any("/api/*path", func(c *gin.Context) {
		apiRouter.HandleContext(c)
	})
	
	// Start your server
	router.Run(":8080")
}
` + "```" + `
`

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return err
	}

	data := struct {
		Title   string
		Version string
	}{
		Title:   spec.Info.Title,
		Version: spec.Info.Version,
	}

	f, err := os.Create(filepath.Join(baseDir, "README.md"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

// Helper functions

func convertPathToGin(path string) string {
	// Convert OpenAPI path params {param} to Gin format :param
	return strings.ReplaceAll(strings.ReplaceAll(path, "{", ":"), "}", "")
}

func toCamelCase(s string) string {
	// Simple implementation - in a real generator, you'd want more robust handling
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}
	return strings.Join(parts, "")
}

func getGoType(schema Schema) string {
	// This is a simplified type mapping
	// In a real implementation, you'd handle more complex types and references
	switch schema.Type {
	case "integer":
		return "int"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	default:
		return "string"
	}
}
