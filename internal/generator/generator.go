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
}

// GenerateCode generates all code from an OpenAPI spec
func GenerateCode(spec *models.OpenAPISpec, config Config) error {
	// Create the output directory structure
	err := createDirectories(config.OutputDir)
	if err != nil {
		return err
	}

	// Generate server
	err = GenerateServerFile(spec, config.OutputDir, config.PackageName)
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