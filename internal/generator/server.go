package generator

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/shubhamku044/gopenapi/internal/models"
	"github.com/shubhamku044/gopenapi/pkg/utils"
)

// PathParam represents a path parameter for route generation
type PathParam struct {
	Name string
	Type string
}

// Route represents a route for the server
type Route struct {
	Method        string
	Path          string
	HandlerName   string
	PathParams    []PathParam
	HasPathParams bool
}

// GenerateServerFile generates the server file
func GenerateServerFile(spec *models.OpenAPISpec, baseDir string, packageName string) error {
	serverTemplate := `package server

import (
	"github.com/gin-gonic/gin"
	"{{.PackageName}}/generated/api"
)

// Server represents the API server
type Server struct {
	router *gin.Engine
	api    api.API
}

// ServerOption represents a server option function
type ServerOption func(*Server)

// WithMiddleware adds middleware to the server
func WithMiddleware(middleware ...gin.HandlerFunc) ServerOption {
	return func(s *Server) {
		s.router.Use(middleware...)
	}
}

// WithMode sets the gin mode (debug, release, test)
func WithMode(mode string) ServerOption {
	return func(s *Server) {
		gin.SetMode(mode)
	}
}

// NewServer creates a new API server
func NewServer(api api.API, options ...ServerOption) *Server {
	s := &Server{
		router: gin.Default(),
		api:    api,
	}
	
	// Apply options
	for _, option := range options {
		option(s)
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

	// Prepare route data
	var routes []Route
	for path, operations := range spec.Paths {
		for method, op := range operations {
			ginPath := utils.ConvertPathToGin(path)
			handlerName := utils.ToCamelCase(op.OperationID)

			// Extract path parameters
			var pathParams []PathParam
			for _, param := range op.Parameters {
				if param.In == "path" {
					pathParams = append(pathParams, PathParam{
						Name: param.Name,
						Type: utils.GetGoType(param.Schema),
					})
				}
			}

			routes = append(routes, Route{
				Method:        strings.Title(strings.ToUpper(method)),
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
