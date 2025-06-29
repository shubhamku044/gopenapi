package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shubhamku044/gopenapi/internal/models"
)

func TestGenerateCode(t *testing.T) {
	// Create a test OpenAPI spec
	spec := &models.OpenAPISpec{
		Paths: map[string]map[string]models.Operation{
			"/users": {
				"get": models.Operation{
					OperationID: "list_users",
					Summary:     "List all users",
					Description: "Returns a list of users",
				},
				"post": models.Operation{
					OperationID: "create_user",
					Summary:     "Create a new user",
					Description: "Creates a new user",
				},
			},
			"/users/{id}": {
				"get": models.Operation{
					OperationID: "get_user",
					Summary:     "Get user by ID",
					Parameters: []models.Parameter{
						{
							Name: "id",
							In:   "path",
							Schema: models.Schema{
								Type: "string",
							},
						},
					},
				},
			},
		},
		Components: struct {
			Schemas map[string]models.Schema `json:"schemas" yaml:"schemas"`
		}{
			Schemas: map[string]models.Schema{
				"User": {
					Type: "object",
					Properties: map[string]models.Schema{
						"id": {
							Type:   "string",
							Format: "uuid",
						},
						"name": {
							Type: "string",
						},
						"email": {
							Type:   "string",
							Format: "email",
						},
					},
				},
			},
		},
	}

	// Test case 1: Generate code with default config
	t.Run("GenerateCodeWithDefaults", func(t *testing.T) {
		tempDir := t.TempDir()
		config := Config{
			OutputDir:   tempDir,
			PackageName: "testapi",
			ModuleName:  "test/module",
		}

		err := GenerateCode(spec, config)
		if err != nil {
			t.Fatalf("GenerateCode failed: %v", err)
		}

		// Check that all expected files are generated
		expectedFiles := []string{
			"go.mod",
			"main.go",
			"api/api.go",
			"api/implementation.go",
			"models/models.go",
			"server/server.go",
			"README.md",
		}

		for _, file := range expectedFiles {
			filePath := filepath.Join(tempDir, file)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("Expected file %s was not generated", file)
			}
		}
	})

	// Test case 2: Don't regenerate main.go if it exists
	t.Run("DontRegenerateMainGo", func(t *testing.T) {
		tempDir := t.TempDir()
		config := Config{
			OutputDir:   tempDir,
			PackageName: "testapi",
			ModuleName:  "test/module",
		}

		// Create main.go with custom content
		mainPath := filepath.Join(tempDir, "main.go")
		customContent := "// Custom main.go content"
		err := os.WriteFile(mainPath, []byte(customContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create custom main.go: %v", err)
		}

		err = GenerateCode(spec, config)
		if err != nil {
			t.Fatalf("GenerateCode failed: %v", err)
		}

		// Check that main.go content is preserved
		content, err := os.ReadFile(mainPath)
		if err != nil {
			t.Fatalf("Failed to read main.go: %v", err)
		}
		if string(content) != customContent {
			t.Errorf("main.go was regenerated when it should have been preserved")
		}
	})

	// Test case 3: Generate implementation template only once
	t.Run("GenerateImplementationOnce", func(t *testing.T) {
		tempDir := t.TempDir()
		config := Config{
			OutputDir:   tempDir,
			PackageName: "testapi",
			ModuleName:  "test/module",
		}

		// Create api directory
		apiDir := filepath.Join(tempDir, "api")
		err := os.MkdirAll(apiDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create api directory: %v", err)
		}

		// Create implementation.go with custom content
		implPath := filepath.Join(apiDir, "implementation.go")
		customContent := "// Custom implementation content"
		err = os.WriteFile(implPath, []byte(customContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create custom implementation.go: %v", err)
		}

		err = GenerateCode(spec, config)
		if err != nil {
			t.Fatalf("GenerateCode failed: %v", err)
		}

		// Check that implementation.go content is preserved
		content, err := os.ReadFile(implPath)
		if err != nil {
			t.Fatalf("Failed to read implementation.go: %v", err)
		}
		if string(content) != customContent {
			t.Errorf("implementation.go was regenerated when it should have been preserved")
		}
	})
}

func TestGenerateGoMod(t *testing.T) {
	tempDir := t.TempDir()
	moduleName := "test/module"

	err := GenerateGoModIfNotExists(tempDir, moduleName)
	if err != nil {
		t.Fatalf("GenerateGoMod failed: %v", err)
	}

	// Check that go.mod file was created
	goModPath := filepath.Join(tempDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Errorf("go.mod file was not created")
	}

	// Check go.mod content
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	expectedContent := "module test/module"
	if !contains(string(content), expectedContent) {
		t.Errorf("go.mod does not contain expected module name")
	}
}

func TestCreateDirectories(t *testing.T) {
	tempDir := t.TempDir()

	err := createProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("createDirectories failed: %v", err)
	}

	// Check that all expected directories are created
	expectedDirs := []string{
		tempDir,
		filepath.Join(tempDir, "api"),
		filepath.Join(tempDir, "models"),
		filepath.Join(tempDir, "server"),
	}

	for _, dir := range expectedDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Expected directory %s was not created", dir)
		}
	}
}

func TestGenerateImplementationTemplate(t *testing.T) {
	spec := &models.OpenAPISpec{
		Paths: map[string]map[string]models.Operation{
			"/users": {
				"get": models.Operation{
					OperationID: "list_users",
					Summary:     "List all users",
				},
			},
		},
	}

	tempDir := t.TempDir()
	moduleName := "test/module"

	// Create api directory
	apiDir := filepath.Join(tempDir, "api")
	err := os.MkdirAll(apiDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create api directory: %v", err)
	}

	err = GenerateHandlerTemplates(spec, tempDir, moduleName)
	if err != nil {
		t.Fatalf("GenerateImplementationTemplate failed: %v", err)
	}

	// Check that implementation.go was created
	implPath := filepath.Join(apiDir, "implementation.go")
	if _, err := os.Stat(implPath); os.IsNotExist(err) {
		t.Errorf("implementation.go file was not created")
	}

	// Check implementation.go content
	content, err := os.ReadFile(implPath)
	if err != nil {
		t.Fatalf("Failed to read implementation.go: %v", err)
	}

	expectedContent := []string{
		"package api",
		"APIImplementation",
		"ListUsers",
		"TODO: Implement your business logic",
	}

	for _, expected := range expectedContent {
		if !contains(string(content), expected) {
			t.Errorf("implementation.go does not contain expected content: %s", expected)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
