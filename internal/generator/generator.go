package generator

import (
	"os"
	"path/filepath"

	"github.com/shubhamku044/gopenapi/internal/models"
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

	err = GenerateGoMod(config.OutputDir, config.ModuleName)
	if err != nil {
		return err
	}

	// Generate server
	err = GenerateServerFile(spec, config.OutputDir, config.PackageName, config.ModuleName)
	if err != nil {
		return err
	}

	// Generate API interfaces
	err = GenerateAPIFile(spec, config.OutputDir)
	if err != nil {
		return err
	}

	// Generate models
	err = GenerateModels(spec, config.OutputDir)
	if err != nil {
		return err
	}

	// Generate README.md with usage examples
	err = GenerateReadme(spec, config.OutputDir, config.PackageName)
	if err != nil {
		return err
	}

	return nil
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
