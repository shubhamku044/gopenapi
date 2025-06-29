package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shubhamku044/gopenapi/internal/models"
)

// Test constants to avoid goconst linting issues
const (
	testAPITitle   = "Test API"
	testAPIVersion = "1.0.0"
	testTitle      = "Test"
	testModule     = "test/module"
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

		// Check that all expected files are generated in new structure
		expectedFiles := []string{
			"go.mod",
			"main.go",
			"handlers/api.go",
			"generated/api/interfaces.go",
			"generated/models/models.go",
			"generated/server/router.go",
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

	// Test case 3: Generate handler template only once
	t.Run("GenerateHandlersOnce", func(t *testing.T) {
		tempDir := t.TempDir()
		config := Config{
			OutputDir:   tempDir,
			PackageName: "testapi",
			ModuleName:  "test/module",
		}

		// Create handlers directory and api.go with custom content
		handlersDir := filepath.Join(tempDir, "handlers")
		err := os.MkdirAll(handlersDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create handlers directory: %v", err)
		}

		handlerPath := filepath.Join(handlersDir, "api.go")
		customContent := "// Custom handler implementation content"
		err = os.WriteFile(handlerPath, []byte(customContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create custom api.go: %v", err)
		}

		err = GenerateCode(spec, config)
		if err != nil {
			t.Fatalf("GenerateCode failed: %v", err)
		}

		// Check that handlers/api.go content is preserved
		content, err := os.ReadFile(handlerPath)
		if err != nil {
			t.Fatalf("Failed to read handlers/api.go: %v", err)
		}
		if string(content) != customContent {
			t.Errorf("handlers/api.go was regenerated when it should have been preserved")
		}
	})

	// Test case 4: Always regenerate files in generated/ directory
	t.Run("AlwaysRegenerateGenerated", func(t *testing.T) {
		tempDir := t.TempDir()
		config := Config{
			OutputDir:   tempDir,
			PackageName: "testapi",
			ModuleName:  "test/module",
		}

		// Create project structure first
		err := createProjectStructure(tempDir)
		if err != nil {
			t.Fatalf("Failed to create project structure: %v", err)
		}

		// Create interfaces.go with custom content
		interfacePath := filepath.Join(tempDir, "generated", "api", "interfaces.go")
		customContent := "// Custom interface content that should be overwritten"
		err = os.WriteFile(interfacePath, []byte(customContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create custom interfaces.go: %v", err)
		}

		err = GenerateCode(spec, config)
		if err != nil {
			t.Fatalf("GenerateCode failed: %v", err)
		}

		// Check that generated/api/interfaces.go was regenerated
		content, err := os.ReadFile(interfacePath)
		if err != nil {
			t.Fatalf("Failed to read interfaces.go: %v", err)
		}
		if string(content) == customContent {
			t.Errorf("generated/api/interfaces.go was not regenerated when it should have been")
		}

		// Check it contains expected generated content
		if !contains(string(content), "type APIHandlers interface") {
			t.Errorf("interfaces.go does not contain expected generated content")
		}
	})
}

func TestGenerateGoMod(t *testing.T) {
	tempDir := t.TempDir()
	moduleName := testModule

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

	expectedContent := "module " + testModule
	if !contains(string(content), expectedContent) {
		t.Errorf("go.mod does not contain expected module name")
	}

	// Check that gin dependency is included
	if !contains(string(content), "github.com/gin-gonic/gin") {
		t.Errorf("go.mod does not contain gin dependency")
	}
}

func TestCreateDirectories(t *testing.T) {
	tempDir := t.TempDir()

	err := createProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("createProjectStructure failed: %v", err)
	}

	// Check that all expected directories are created with new structure
	expectedDirs := []string{
		tempDir,
		filepath.Join(tempDir, "handlers"),
		filepath.Join(tempDir, "generated", "api"),
		filepath.Join(tempDir, "generated", "models"),
		filepath.Join(tempDir, "generated", "server"),
	}

	for _, dir := range expectedDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Expected directory %s was not created", dir)
		}
	}
}

func TestGenerateHandlerTemplates(t *testing.T) {
	spec := &models.OpenAPISpec{
		Paths: map[string]map[string]models.Operation{
			"/users": {
				"get": models.Operation{
					OperationID: "list_users",
					Summary:     "List all users",
				},
				"post": models.Operation{
					OperationID: "create_user",
					Summary:     "Create a new user",
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
						"id":   {Type: "string"},
						"name": {Type: "string"},
					},
				},
			},
		},
	}

	tempDir := t.TempDir()
	moduleName := testModule

	// Create project structure first
	err := createProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	err = GenerateHandlerTemplates(spec, tempDir, moduleName)
	if err != nil {
		t.Fatalf("GenerateHandlerTemplates failed: %v", err)
	}

	// Check that handlers/api.go was created
	handlerPath := filepath.Join(tempDir, "handlers", "api.go")
	if _, err := os.Stat(handlerPath); os.IsNotExist(err) {
		t.Errorf("handlers/api.go file was not created")
	}

	// Check handlers/api.go content
	content, err := os.ReadFile(handlerPath)
	if err != nil {
		t.Fatalf("Failed to read handlers/api.go: %v", err)
	}

	expectedContent := []string{
		"package handlers",
		"APIHandlers",
		"ListUsers",
		"CreateUser",
		"TODO: Implement your business logic",
		moduleName + "/generated/api",
		moduleName + "/generated/models",
	}

	for _, expected := range expectedContent {
		if !contains(string(content), expected) {
			t.Errorf("handlers/api.go does not contain expected content: %s", expected)
		}
	}
}

func TestGenerateInterfaces(t *testing.T) {
	spec := &models.OpenAPISpec{
		Paths: map[string]map[string]models.Operation{
			"/health": {
				"get": models.Operation{
					OperationID: "health_check",
					Summary:     "Health check",
				},
			},
		},
	}

	tempDir := t.TempDir()
	moduleName := testModule

	// Create project structure first
	err := createProjectStructure(tempDir)
	if err != nil {
		t.Fatalf("Failed to create project structure: %v", err)
	}

	err = GenerateInterfaces(spec, tempDir, moduleName)
	if err != nil {
		t.Fatalf("GenerateInterfaces failed: %v", err)
	}

	// Check that generated/api/interfaces.go was created
	interfacePath := filepath.Join(tempDir, "generated", "api", "interfaces.go")
	if _, err := os.Stat(interfacePath); os.IsNotExist(err) {
		t.Errorf("generated/api/interfaces.go file was not created")
	}

	// Check content
	content, err := os.ReadFile(interfacePath)
	if err != nil {
		t.Fatalf("Failed to read interfaces.go: %v", err)
	}

	expectedContent := []string{
		"package api",
		"type APIHandlers interface",
		"HealthCheck",
		"github.com/gin-gonic/gin",
	}

	for _, expected := range expectedContent {
		if !contains(string(content), expected) {
			t.Errorf("interfaces.go does not contain expected content: %s", expected)
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

// Test functions with 0% coverage
func TestGenerateAPIFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create the api directory that GenerateAPIFile expects
	apiDir := filepath.Join(tempDir, "api")
	err := os.MkdirAll(apiDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create api directory: %v", err)
	}

	spec := &models.OpenAPISpec{}
	spec.Info.Title = testAPITitle
	spec.Info.Version = testAPIVersion
	spec.Paths = map[string]map[string]models.Operation{
		"/users": {
			"get": {
				OperationID: "list_users",
				Summary:     "List users",
			},
			"post": {
				OperationID: "create_user",
				Summary:     "Create user",
			},
		},
	}

	err = GenerateAPIFile(spec, tempDir)
	if err != nil {
		t.Fatalf("GenerateAPIFile failed: %v", err)
	}

	// Verify file was created at correct path
	apiFile := filepath.Join(tempDir, "api", "api.go")
	if _, err := os.Stat(apiFile); os.IsNotExist(err) {
		t.Errorf("Expected API file to be created at %s", apiFile)
	}

	// Verify content
	content, err := os.ReadFile(apiFile)
	if err != nil {
		t.Fatalf("Failed to read API file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "ListUsers") {
		t.Errorf("Expected API file to contain ListUsers method")
	}
	if !strings.Contains(contentStr, "CreateUser") {
		t.Errorf("Expected API file to contain CreateUser method")
	}
}

func TestGenerateMainFile(t *testing.T) {
	tempDir := t.TempDir()

	spec := &models.OpenAPISpec{}
	spec.Info.Title = testAPITitle
	spec.Info.Version = testAPIVersion

	err := GenerateMainFile(spec, tempDir, testModule)
	if err != nil {
		t.Fatalf("GenerateMainFile failed: %v", err)
	}

	// Verify file was created
	mainFile := filepath.Join(tempDir, "main.go")
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		t.Errorf("Expected main.go to be created at %s", mainFile)
	}

	// Verify content
	content, err := os.ReadFile(mainFile)
	if err != nil {
		t.Fatalf("Failed to read main.go: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "package main") {
		t.Errorf("Expected main.go to contain package main")
	}
	if !strings.Contains(contentStr, testModule) {
		t.Errorf("Expected main.go to contain module name")
	}
}

func TestGenerateServerFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create the server directory that GenerateServerFile expects
	serverDir := filepath.Join(tempDir, "server")
	err := os.MkdirAll(serverDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create server directory: %v", err)
	}

	spec := &models.OpenAPISpec{}
	spec.Info.Title = testAPITitle
	spec.Info.Version = testAPIVersion
	spec.Paths = map[string]map[string]models.Operation{
		"/health": {
			"get": {
				OperationID: "health_check",
			},
		},
	}

	err = GenerateServerFile(spec, tempDir, "testpkg", testModule)
	if err != nil {
		t.Fatalf("GenerateServerFile failed: %v", err)
	}

	// Verify file was created at correct path
	serverFile := filepath.Join(tempDir, "server", "server.go")
	if _, err := os.Stat(serverFile); os.IsNotExist(err) {
		t.Errorf("Expected server file to be created at %s", serverFile)
	}

	// Verify content
	content, err := os.ReadFile(serverFile)
	if err != nil {
		t.Fatalf("Failed to read server file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "package server") {
		t.Errorf("Expected server file to contain package server")
	}
}

// Test edge cases for better coverage
func TestGenerateCodeErrorCases(t *testing.T) {
	t.Run("InvalidOutputDirectory", func(t *testing.T) {
		// Try to write to a directory that doesn't exist and can't be created
		invalidDir := "/invalid/path/that/cannot/be/created"

		spec := &models.OpenAPISpec{}
		spec.Info.Title = testTitle

		config := Config{
			OutputDir:   invalidDir,
			PackageName: testTitle,
			ModuleName:  testTitle,
		}

		err := GenerateCode(spec, config)
		if err == nil {
			t.Errorf("Expected error when using invalid output directory")
		}
	})
}

func TestCreateProjectStructureEdgeCases(t *testing.T) {
	t.Run("CreateDirectoriesInExistingPath", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create some directories first
		handlerDir := filepath.Join(tempDir, "handlers")
		err := os.MkdirAll(handlerDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create handlers directory: %v", err)
		}

		err = createProjectStructure(tempDir)
		if err != nil {
			t.Fatalf("createProjectStructure failed: %v", err)
		}

		// Verify all directories exist
		dirs := []string{
			"handlers",
			"generated/api",
			"generated/models",
			"generated/server",
		}

		for _, dir := range dirs {
			dirPath := filepath.Join(tempDir, dir)
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				t.Errorf("Expected directory %s to exist", dir)
			}
		}
	})
}

func TestGenerateGoModIfNotExistsEdgeCases(t *testing.T) {
	t.Run("WithExistingGoMod", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create existing go.mod
		existingContent := "module existing/module\n\ngo 1.21"
		goModPath := filepath.Join(tempDir, "go.mod")
		err := os.WriteFile(goModPath, []byte(existingContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create existing go.mod: %v", err)
		}

		err = GenerateGoModIfNotExists(tempDir, "newmodule")
		if err != nil {
			t.Fatalf("GenerateGoModIfNotExists failed: %v", err)
		}

		// Verify original content is preserved
		content, err := os.ReadFile(goModPath)
		if err != nil {
			t.Fatalf("Failed to read go.mod: %v", err)
		}

		if !strings.Contains(string(content), "existing/module") {
			t.Errorf("Expected existing go.mod content to be preserved")
		}
	})
}

func TestGenerateUserMainIfNotExistsEdgeCases(t *testing.T) {
	t.Run("WithExistingMainGo", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create existing main.go
		existingContent := "package main\n\n// existing main"
		mainPath := filepath.Join(tempDir, "main.go")
		err := os.WriteFile(mainPath, []byte(existingContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create existing main.go: %v", err)
		}

		spec := &models.OpenAPISpec{}
		spec.Info.Title = testTitle
		err = GenerateUserMainIfNotExists(spec, tempDir, testModule)
		if err != nil {
			t.Fatalf("GenerateUserMainIfNotExists failed: %v", err)
		}

		// Verify original content is preserved
		content, err := os.ReadFile(mainPath)
		if err != nil {
			t.Fatalf("Failed to read main.go: %v", err)
		}

		if !strings.Contains(string(content), "existing main") {
			t.Errorf("Expected existing main.go content to be preserved")
		}
	})
}

func TestGenerateHandlerTemplatesEdgeCases(t *testing.T) {
	t.Run("WithExistingHandlers", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create handlers directory with existing file
		handlerDir := filepath.Join(tempDir, "handlers")
		err := os.MkdirAll(handlerDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create handlers directory: %v", err)
		}

		existingFile := filepath.Join(handlerDir, "existing.go")
		err = os.WriteFile(existingFile, []byte("// existing handler"), 0600)
		if err != nil {
			t.Fatalf("Failed to create existing handler: %v", err)
		}

		spec := &models.OpenAPISpec{
			Paths: map[string]map[string]models.Operation{
				"/test": {
					"get": {OperationID: "test_get"},
				},
			},
		}

		err = GenerateHandlerTemplates(spec, tempDir, testModule)
		if err != nil {
			t.Fatalf("GenerateHandlerTemplates failed: %v", err)
		}

		// Should not create api.go since handlers directory has files
		apiFile := filepath.Join(handlerDir, "api.go")
		if _, err := os.Stat(apiFile); err == nil {
			t.Errorf("Should not create api.go when handlers directory has existing files")
		}
	})

	t.Run("WithEmptyHandlersDirectory", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create empty handlers directory
		handlerDir := filepath.Join(tempDir, "handlers")
		err := os.MkdirAll(handlerDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create handlers directory: %v", err)
		}

		spec := &models.OpenAPISpec{
			Paths: map[string]map[string]models.Operation{
				"/users": {
					"get":  {OperationID: "list_users"},
					"post": {OperationID: "create_user"},
				},
				"/users/{id}": {
					"get":    {OperationID: "get_user"},
					"put":    {OperationID: "update_user"},
					"delete": {OperationID: "delete_user"},
				},
			},
		}

		err = GenerateHandlerTemplates(spec, tempDir, testModule)
		if err != nil {
			t.Fatalf("GenerateHandlerTemplates failed: %v", err)
		}

		// Should create api.go
		apiFile := filepath.Join(handlerDir, "api.go")
		if _, err := os.Stat(apiFile); os.IsNotExist(err) {
			t.Errorf("Expected api.go to be created")
		}

		// Verify content includes all methods
		content, err := os.ReadFile(apiFile)
		if err != nil {
			t.Fatalf("Failed to read api.go: %v", err)
		}

		contentStr := string(content)
		methods := []string{"ListUsers", "CreateUser", "GetUser", "UpdateUser", "DeleteUser"}
		for _, method := range methods {
			if !strings.Contains(contentStr, method) {
				t.Errorf("Expected api.go to contain method %s", method)
			}
		}
	})
}

func TestHasTimeFieldsFunction(t *testing.T) {
	t.Run("SchemaWithTimeFields", func(t *testing.T) {
		schema := models.Schema{
			Type:   "string",
			Format: "date-time",
		}

		result := hasTimeFields(schema)
		if !result {
			t.Errorf("Expected hasTimeFields to return true for schema with date-time format")
		}
	})

	t.Run("SchemaWithDateFields", func(t *testing.T) {
		schema := models.Schema{
			Type:   "string",
			Format: "date",
		}

		result := hasTimeFields(schema)
		if !result {
			t.Errorf("Expected hasTimeFields to return true for schema with date format")
		}
	})

	t.Run("SchemaWithTimeFieldsInProperties", func(t *testing.T) {
		schema := models.Schema{
			Type: "object",
			Properties: map[string]models.Schema{
				"created_at": {Type: "string", Format: "date-time"},
				"name":       {Type: "string"},
			},
		}

		result := hasTimeFields(schema)
		if !result {
			t.Errorf("Expected hasTimeFields to return true for schema with date-time fields in properties")
		}
	})

	t.Run("SchemaWithoutTimeFields", func(t *testing.T) {
		schema := models.Schema{
			Type: "object",
			Properties: map[string]models.Schema{
				"name":  {Type: "string"},
				"email": {Type: "string"},
			},
		}

		result := hasTimeFields(schema)
		if result {
			t.Errorf("Expected hasTimeFields to return false for schema without date/date-time fields")
		}
	})

	t.Run("EmptySchema", func(t *testing.T) {
		schema := models.Schema{}

		result := hasTimeFields(schema)
		if result {
			t.Errorf("Expected hasTimeFields to return false for empty schema")
		}
	})
}

func TestGenerateInterfacesEdgeCases(t *testing.T) {
	t.Run("SpecWithComplexOperations", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create the directory structure that GenerateInterfaces expects
		err := createProjectStructure(tempDir)
		if err != nil {
			t.Fatalf("Failed to create project structure: %v", err)
		}

		spec := &models.OpenAPISpec{
			Paths: map[string]map[string]models.Operation{
				"/api/v1/users/{user_id}/posts/{post_id}": {
					"get": {
						OperationID: "get_user_post_by_id",
						Parameters: []models.Parameter{
							{Name: "user_id", In: "path"},
							{Name: "post_id", In: "path"},
							{Name: "include", In: "query"},
						},
					},
				},
			},
		}

		err = GenerateInterfaces(spec, tempDir, "complex/module")
		if err != nil {
			t.Fatalf("GenerateInterfaces failed: %v", err)
		}

		// Verify file creation
		interfaceFile := filepath.Join(tempDir, "generated", "api", "interfaces.go")
		if _, err := os.Stat(interfaceFile); os.IsNotExist(err) {
			t.Errorf("Expected interfaces.go to be created")
		}
	})
}

func TestGenerateRouterEdgeCases(t *testing.T) {
	t.Run("SpecWithVariousHTTPMethods", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create the directory structure that GenerateRouter expects
		err := createProjectStructure(tempDir)
		if err != nil {
			t.Fatalf("Failed to create project structure: %v", err)
		}

		spec := &models.OpenAPISpec{
			Paths: map[string]map[string]models.Operation{
				"/api/resource": {
					"get":     {OperationID: "get_resource"},
					"post":    {OperationID: "create_resource"},
					"put":     {OperationID: "update_resource"},
					"delete":  {OperationID: "delete_resource"},
					"patch":   {OperationID: "patch_resource"},
					"head":    {OperationID: "head_resource"},
					"options": {OperationID: "options_resource"},
				},
			},
		}

		err = GenerateRouter(spec, tempDir, "testmodule")
		if err != nil {
			t.Fatalf("GenerateRouter failed: %v", err)
		}

		// Verify file creation
		routerFile := filepath.Join(tempDir, "generated", "server", "router.go")
		if _, err := os.Stat(routerFile); os.IsNotExist(err) {
			t.Errorf("Expected router.go to be created")
		}

		// Verify content includes all HTTP methods
		content, err := os.ReadFile(routerFile)
		if err != nil {
			t.Fatalf("Failed to read router.go: %v", err)
		}

		contentStr := string(content)
		methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
		for _, method := range methods {
			if !strings.Contains(contentStr, method) {
				t.Errorf("Expected router.go to contain HTTP method %s", method)
			}
		}
	})
}

// Additional comprehensive tests to improve coverage
func TestGenerateModels(t *testing.T) {
	tempDir := t.TempDir()

	spec := &models.OpenAPISpec{}
	spec.Info.Title = testAPITitle
	spec.Info.Version = testAPIVersion
	spec.Components.Schemas = map[string]models.Schema{
		"User": {
			Type: "object",
			Properties: map[string]models.Schema{
				"id":   {Type: "integer"},
				"name": {Type: "string"},
			},
		},
	}

	err := GenerateModels(spec, tempDir)
	if err != nil {
		t.Fatalf("GenerateModels failed: %v", err)
	}

	// Verify file was created
	modelsFile := filepath.Join(tempDir, "models", "models.go")
	if _, err := os.Stat(modelsFile); os.IsNotExist(err) {
		t.Errorf("Expected models file to be created at %s", modelsFile)
	}

	// Verify content
	content, err := os.ReadFile(modelsFile)
	if err != nil {
		t.Fatalf("Failed to read models file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "type User struct") {
		t.Errorf("Expected models file to contain User struct")
	}
}

func TestGenerateReadme(t *testing.T) {
	tempDir := t.TempDir()

	spec := &models.OpenAPISpec{}
	spec.Info.Title = testAPITitle
	spec.Info.Version = testAPIVersion
	spec.Info.Description = "A test API for demonstration"

	err := GenerateReadme(spec, tempDir, testModule)
	if err != nil {
		t.Fatalf("GenerateReadme failed: %v", err)
	}

	// Verify file was created
	readmeFile := filepath.Join(tempDir, "README.md")
	if _, err := os.Stat(readmeFile); os.IsNotExist(err) {
		t.Errorf("Expected README.md to be created at %s", readmeFile)
	}

	// Verify content
	content, err := os.ReadFile(readmeFile)
	if err != nil {
		t.Fatalf("Failed to read README.md: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, testAPITitle) {
		t.Errorf("Expected README.md to contain API title")
	}
}

func TestGenerateGoModIfNotExists(t *testing.T) {
	t.Run("WithoutExistingGoMod", func(t *testing.T) {
		tempDir := t.TempDir()

		err := GenerateGoModIfNotExists(tempDir, testModule)
		if err != nil {
			t.Fatalf("GenerateGoModIfNotExists failed: %v", err)
		}

		// Verify file was created
		goModFile := filepath.Join(tempDir, "go.mod")
		if _, err := os.Stat(goModFile); os.IsNotExist(err) {
			t.Errorf("Expected go.mod to be created at %s", goModFile)
		}

		// Verify content
		content, err := os.ReadFile(goModFile)
		if err != nil {
			t.Fatalf("Failed to read go.mod: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, testModule) {
			t.Errorf("Expected go.mod to contain module name")
		}
	})

	t.Run("WithExistingGoMod", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create existing go.mod
		existingContent := "module existing/module\n\ngo 1.21"
		goModPath := filepath.Join(tempDir, "go.mod")
		err := os.WriteFile(goModPath, []byte(existingContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create existing go.mod: %v", err)
		}

		err = GenerateGoModIfNotExists(tempDir, "newmodule")
		if err != nil {
			t.Fatalf("GenerateGoModIfNotExists failed: %v", err)
		}

		// Verify original content is preserved
		content, err := os.ReadFile(goModPath)
		if err != nil {
			t.Fatalf("Failed to read go.mod: %v", err)
		}

		if !strings.Contains(string(content), "existing/module") {
			t.Errorf("Expected existing go.mod content to be preserved")
		}
	})
}
