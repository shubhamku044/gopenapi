package generator

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/shubhamku044/gopenapi/internal/models"
	"github.com/shubhamku044/gopenapi/pkg/utils"
)

const (
	pathParameterType = "path"
)

// APIMethod represents an API method for generation
type APIMethod struct {
	Name        string
	HandlerName string
	Comment     string
	Parameters  string
}

// GenerateAPIFile generates the API interface file
func GenerateAPIFile(spec *models.OpenAPISpec, baseDir string) error {
	apiTemplate := `package api

import (
	"github.com/gin-gonic/gin"
)

// API defines the interface for API operations
type API interface {
{{range .Methods}}
	{{.Comment}}
	{{.HandlerName}}(c *gin.Context{{.Parameters}})
{{end}}
}
`

	tmpl, err := template.New("api").Parse(apiTemplate)
	if err != nil {
		return err
	}

	// Generate methods from OpenAPI spec
	var methods []APIMethod
	for path, operations := range spec.Paths {
		for method, op := range operations {
			handlerName := utils.ToCamelCase(op.OperationID)

			// Build comment
			comment := "// " + handlerName + " handles " + strings.ToUpper(method) + " " + path
			if op.Summary != "" {
				comment += "\n\t// " + op.Summary
			}
			if op.Description != "" {
				comment += "\n\t// " + op.Description
			}

			// Build parameters
			var params []string
			for _, param := range op.Parameters {
				if param.In == pathParameterType {
					paramType := utils.GetGoType(param.Schema)
					params = append(params, param.Name+" "+paramType)
				}
			}

			paramStr := ""
			if len(params) > 0 {
				paramStr = ", " + strings.Join(params, ", ")
			}

			methods = append(methods, APIMethod{
				Name:        strings.ToUpper(method) + " " + path,
				HandlerName: handlerName,
				Comment:     comment,
				Parameters:  paramStr,
			})
		}
	}

	data := struct {
		Methods []APIMethod
	}{
		Methods: methods,
	}

	f, err := os.Create(filepath.Join(baseDir, "api", "api.go"))
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}
