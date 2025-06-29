package generator

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/shubhamku044/gopenapi/internal/models"
	"github.com/shubhamku044/gopenapi/pkg/utils"
)

// Config represents generator configuration
type Config struct {
	OutputDir   string
	PackageName string
	ModuleName  string
}

// GenerateCode generates all code from an OpenAPI spec
func GenerateCode(spec *models.OpenAPISpec, config Config) error {
	if config.ModuleName == "" {
		config.ModuleName = config.PackageName
	}

	// Create the output directory structure
	err := createDirectories(config.OutputDir)
	if err != nil {
		return err
	}

	// Generate go.mod file for the generated code
	err = GenerateGoMod(config.OutputDir, config.ModuleName)
	if err != nil {
		return err
	}

	// Generate main.go file only if it doesn't exist (user might have customized it)
	mainPath := filepath.Join(config.OutputDir, "main.go")
	if _, err := os.Stat(mainPath); os.IsNotExist(err) {
		err = GenerateMainFile(spec, config.OutputDir, config.ModuleName)
		if err != nil {
			return err
		}
	}

	// Always regenerate server (infrastructure code)
	err = GenerateServerFile(spec, config.OutputDir, config.PackageName, config.ModuleName)
	if err != nil {
		return err
	}

	// Always regenerate API interfaces (based on OpenAPI spec)
	err = GenerateAPIFile(spec, config.OutputDir)
	if err != nil {
		return err
	}

	// Always regenerate models (based on OpenAPI spec)
	err = GenerateModels(spec, config.OutputDir)
	if err != nil {
		return err
	}

	// Generate README.md with usage examples
	err = GenerateReadme(spec, config.OutputDir, config.PackageName)
	if err != nil {
		return err
	}

	// Generate implementation template if it doesn't exist
	err = GenerateImplementationTemplate(spec, config.OutputDir, config.ModuleName)
	if err != nil {
		return err
	}

	return nil
}

// GenerateImplementationTemplate generates a template showing users where to implement business logic
func GenerateImplementationTemplate(spec *models.OpenAPISpec, baseDir string, moduleName string) error {
	implPath := filepath.Join(baseDir, "api", "implementation.go")

	// Only generate if it doesn't exist (user might have implemented it)
	if _, err := os.Stat(implPath); !os.IsNotExist(err) {
		return nil
	}

	implTemplate := `package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// APIImplementation provides a concrete implementation of the API interface
// TODO: Implement your business logic in these methods
type APIImplementation struct {
	// Add your dependencies here (database, services, etc.)
	// db     *sql.DB
	// logger *log.Logger
}

// NewAPI creates a new API implementation
func NewAPI() API {
	return &APIImplementation{
		// Initialize your dependencies here
	}
}

{{range .Operations}}
// {{.HandlerName}} {{.Description}}
{{.Comment}}
func (impl *APIImplementation) {{.HandlerName}}(c *gin.Context{{.Params}}) {
	// TODO: Implement your business logic here
	{{.DefaultResponse}}
}

{{end}}`

	tmpl, err := template.New("implementation").Parse(implTemplate)
	if err != nil {
		return err
	}

	// Prepare operations data
	var operations []struct {
		HandlerName     string
		Description     string
		Comment         string
		Params          string
		DefaultResponse string
	}

	for path, pathOps := range spec.Paths {
		for method, op := range pathOps {
			handlerName := utils.ToCamelCase(op.OperationID)
			description := op.Summary
			if description == "" {
				description = op.Description
			}

			// Build comment
			comment := "// " + strings.ToUpper(method) + " " + path
			if description != "" {
				comment += "\n// " + description
			}

			// Build parameters
			var params []string
			for _, param := range op.Parameters {
				if param.In == "path" {
					params = append(params, param.Name+" string")
				}
			}
			paramStr := ""
			if len(params) > 0 {
				paramStr = ", " + strings.Join(params, ", ")
			}

			// Default response based on method
			defaultResp := `c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented",
		"message": "Please implement this endpoint",
	})`

			if strings.ToUpper(method) == "POST" {
				defaultResp = `c.JSON(http.StatusCreated, gin.H{
		"message": "Created successfully",
	})`
			} else if strings.ToUpper(method) == "GET" {
				defaultResp = `c.JSON(http.StatusOK, gin.H{
		"data": "TODO: Return your data here",
	})`
			}

			operations = append(operations, struct {
				HandlerName     string
				Description     string
				Comment         string
				Params          string
				DefaultResponse string
			}{
				HandlerName:     handlerName,
				Description:     description,
				Comment:         comment,
				Params:          paramStr,
				DefaultResponse: defaultResp,
			})
		}
	}

	data := struct {
		ModuleName string
		HasModels  bool
		Operations []struct {
			HandlerName     string
			Description     string
			Comment         string
			Params          string
			DefaultResponse string
		}
	}{
		ModuleName: moduleName,
		HasModels:  len(spec.Components.Schemas) > 0,
		Operations: operations,
	}

	f, err := os.Create(implPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

// GenerateGoMod generates a go.mod file for the generated code
func GenerateGoMod(baseDir string, moduleName string) error {
	goModContent := `module ` + moduleName + `

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
)
`

	f, err := os.Create(filepath.Join(baseDir, "go.mod"))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(goModContent)
	return err
}

// createDirectories creates the output directory structure
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
