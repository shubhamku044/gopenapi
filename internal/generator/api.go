package generator

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/shubhamku044/gopenapi/internal/models"
	"github.com/shubhamku044/gopenapi/pkg/utils"
)

type Handler struct {
	Method        string
	Path          string
	HandlerName   string
	Summary       string
	Description   string
	PathParams    []PathParam
	HasPathParams bool
}

func GenerateAPIFile(spec *models.OpenAPISpec, baseDir string) error {
	apiTemplate := `package api

import (
	"github.com/gin-gonic/gin"
)

// API defines the interface for all API operations
type API interface {
	{{range .Handlers}}
	// {{.HandlerName}} handles {{.Method}} {{.Path}}
	{{if .Summary}}// {{.Summary}}{{end}}
	{{if .Description}}// {{.Description}}{{end}}
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

	var handlers []Handler
	for path, operations := range spec.Paths {
		for method, op := range operations {
			var pathParams []PathParam
			for _, param := range op.Parameters {
				if param.In == "path" {
					pathParams = append(pathParams, PathParam{
						Name: param.Name,
						Type: utils.GetGoType(param.Schema),
					})
				}
			}

			handlers = append(handlers, Handler{
				Method:        strings.ToUpper(method),
				Path:          path,
				HandlerName:   utils.ToCamelCase(op.OperationID),
				Summary:       op.Summary,
				Description:   op.Description,
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
