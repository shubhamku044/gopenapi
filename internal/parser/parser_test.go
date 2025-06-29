package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/shubhamku044/gopenapi/internal/models"
	"gopkg.in/yaml.v3"
)

// ParseOpenAPISpec parses OpenAPI spec from YAML content (helper for testing)
func ParseOpenAPISpec(data []byte) (*models.OpenAPISpec, error) {
	var spec models.OpenAPISpec
	err := yaml.Unmarshal(data, &spec)
	if err != nil {
		return nil, err
	}

	ProcessSpec(&spec)
	return &spec, nil
}

func TestParseOpenAPISpec(t *testing.T) {
	// Test case 1: Parse a valid OpenAPI spec
	t.Run("ParseValidSpec", func(t *testing.T) {
		yamlContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /users:
    get:
      operationId: list_users
      summary: List users
      responses:
        '200':
          description: OK
    post:
      operationId: create_user
      summary: Create user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/User'
      responses:
        '201':
          description: Created
  /users/{id}:
    get:
      operationId: get_user
      summary: Get user by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: OK
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        email:
          type: string
          format: email
        created_at:
          type: string
          format: date-time
`

		spec, err := ParseOpenAPISpec([]byte(yamlContent))
		if err != nil {
			t.Fatalf("ParseOpenAPISpec failed: %v", err)
		}

		// Check basic info
		if spec.Info.Title != "Test API" {
			t.Errorf("Expected title 'Test API', got '%s'", spec.Info.Title)
		}
		if spec.Info.Version != "1.0.0" {
			t.Errorf("Expected version '1.0.0', got '%s'", spec.Info.Version)
		}

		// Check paths
		if len(spec.Paths) != 2 {
			t.Errorf("Expected 2 paths, got %d", len(spec.Paths))
		}

		// Check /users path
		usersPath, exists := spec.Paths["/users"]
		if !exists {
			t.Errorf("Expected /users path to exist")
		}
		if len(usersPath) != 2 {
			t.Errorf("Expected 2 operations for /users, got %d", len(usersPath))
		}

		// Check GET /users
		getUsersOp, exists := usersPath["get"]
		if !exists {
			t.Errorf("Expected GET operation for /users")
		}
		if getUsersOp.OperationID != "list_users" {
			t.Errorf("Expected operationId 'list_users', got '%s'", getUsersOp.OperationID)
		}

		// Check /users/{id} path
		userByIdPath, exists := spec.Paths["/users/{id}"]
		if !exists {
			t.Errorf("Expected /users/{id} path to exist")
		}

		// Check GET /users/{id} parameters
		getUserOp, exists := userByIdPath["get"]
		if !exists {
			t.Errorf("Expected GET operation for /users/{id}")
		}
		if len(getUserOp.Parameters) != 1 {
			t.Errorf("Expected 1 parameter for GET /users/{id}, got %d", len(getUserOp.Parameters))
		}

		param := getUserOp.Parameters[0]
		if param.Name != "id" {
			t.Errorf("Expected parameter name 'id', got '%s'", param.Name)
		}
		if param.In != "path" {
			t.Errorf("Expected parameter in 'path', got '%s'", param.In)
		}

		// Check components
		if len(spec.Components.Schemas) != 1 {
			t.Errorf("Expected 1 schema, got %d", len(spec.Components.Schemas))
		}

		userSchema, exists := spec.Components.Schemas["User"]
		if !exists {
			t.Errorf("Expected User schema to exist")
		}
		if userSchema.Type != "object" {
			t.Errorf("Expected User schema type 'object', got '%s'", userSchema.Type)
		}
		if len(userSchema.Properties) != 4 {
			t.Errorf("Expected 4 properties for User, got %d", len(userSchema.Properties))
		}

		// Check specific properties
		emailProp, exists := userSchema.Properties["email"]
		if !exists {
			t.Errorf("Expected email property to exist")
		}
		if emailProp.Format != "email" {
			t.Errorf("Expected email format 'email', got '%s'", emailProp.Format)
		}

		createdAtProp, exists := userSchema.Properties["created_at"]
		if !exists {
			t.Errorf("Expected created_at property to exist")
		}
		if createdAtProp.Format != "date-time" {
			t.Errorf("Expected created_at format 'date-time', got '%s'", createdAtProp.Format)
		}
	})

	// Test case 2: Parse invalid YAML
	t.Run("ParseInvalidYAML", func(t *testing.T) {
		invalidYAML := `
invalid yaml content
  - missing proper structure
`

		_, err := ParseOpenAPISpec([]byte(invalidYAML))
		if err == nil {
			t.Errorf("Expected error when parsing invalid YAML, but got nil")
		}
	})

	// Test case 3: Parse empty spec
	t.Run("ParseEmptySpec", func(t *testing.T) {
		emptyYAML := `
openapi: 3.0.0
info:
  title: Empty API
  version: 1.0.0
`

		spec, err := ParseOpenAPISpec([]byte(emptyYAML))
		if err != nil {
			t.Fatalf("ParseOpenAPISpec failed: %v", err)
		}

		if len(spec.Paths) != 0 {
			t.Errorf("Expected 0 paths for empty spec, got %d", len(spec.Paths))
		}
		if len(spec.Components.Schemas) != 0 {
			t.Errorf("Expected 0 schemas for empty spec, got %d", len(spec.Components.Schemas))
		}
	})

	// Test case 4: Parse spec with array schema
	t.Run("ParseArraySchema", func(t *testing.T) {
		yamlContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /tags:
    get:
      operationId: list_tags
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: array
                items:
                  type: string
components:
  schemas:
    StringArray:
      type: array
      items:
        type: string
    ObjectArray:
      type: array
      items:
        $ref: '#/components/schemas/Tag'
    Tag:
      type: object
      properties:
        name:
          type: string
`

		spec, err := ParseOpenAPISpec([]byte(yamlContent))
		if err != nil {
			t.Fatalf("ParseOpenAPISpec failed: %v", err)
		}

		// Check array schemas
		stringArraySchema, exists := spec.Components.Schemas["StringArray"]
		if !exists {
			t.Errorf("Expected StringArray schema to exist")
		}
		if stringArraySchema.Type != "array" {
			t.Errorf("Expected StringArray type 'array', got '%s'", stringArraySchema.Type)
		}
		if stringArraySchema.Items.Type != "string" {
			t.Errorf("Expected StringArray items type 'string', got '%s'", stringArraySchema.Items.Type)
		}

		objectArraySchema, exists := spec.Components.Schemas["ObjectArray"]
		if !exists {
			t.Errorf("Expected ObjectArray schema to exist")
		}
		if objectArraySchema.Type != "array" {
			t.Errorf("Expected ObjectArray type 'array', got '%s'", objectArraySchema.Type)
		}
	})
}

// Additional tests for missing coverage
func TestParseSpecFile(t *testing.T) {
	t.Run("ValidYAMLFile", func(t *testing.T) {
		tempDir := t.TempDir()
		yamlContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /users:
    get:
      operationId: list_users
      summary: List users
      responses:
        '200':
          description: OK
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: string
        name:
          type: string
`

		specFile := filepath.Join(tempDir, "test.yaml")
		err := os.WriteFile(specFile, []byte(yamlContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		spec, err := ParseSpecFile(specFile)
		if err != nil {
			t.Fatalf("ParseSpecFile failed: %v", err)
		}

		if spec.Info.Title != "Test API" {
			t.Errorf("Expected title 'Test API', got '%s'", spec.Info.Title)
		}
	})

	t.Run("ValidJSONFile", func(t *testing.T) {
		tempDir := t.TempDir()
		jsonContent := `{
			"openapi": "3.0.0",
			"info": {
				"title": "JSON API",
				"version": "1.0.0"
			},
			"paths": {
				"/test": {
					"get": {
						"operationId": "test_get",
						"responses": {
							"200": {
								"description": "OK"
							}
						}
					}
				}
			}
		}`

		specFile := filepath.Join(tempDir, "test.json")
		err := os.WriteFile(specFile, []byte(jsonContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		spec, err := ParseSpecFile(specFile)
		if err != nil {
			t.Fatalf("ParseSpecFile failed: %v", err)
		}

		if spec.Info.Title != "JSON API" {
			t.Errorf("Expected title 'JSON API', got '%s'", spec.Info.Title)
		}
	})

	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := ParseSpecFile("nonexistent.yaml")
		if err == nil {
			t.Errorf("Expected error for non-existent file")
		}
	})

	t.Run("InvalidYAMLFile", func(t *testing.T) {
		tempDir := t.TempDir()
		invalidContent := `invalid: yaml: content: [unclosed`

		specFile := filepath.Join(tempDir, "invalid.yaml")
		err := os.WriteFile(specFile, []byte(invalidContent), 0600)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		_, err = ParseSpecFile(specFile)
		if err == nil {
			t.Errorf("Expected error for invalid YAML")
		}
	})

	t.Run("UnsupportedFileExtension", func(t *testing.T) {
		tempDir := t.TempDir()
		content := `openapi: 3.0.0`

		specFile := filepath.Join(tempDir, "test.txt")
		err := os.WriteFile(specFile, []byte(content), 0600)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		_, err = ParseSpecFile(specFile)
		if err == nil {
			t.Errorf("Expected error for unsupported file extension")
		}
	})
}

func TestToCamelCaseParser(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"test-case", "TestCase"},
		{"simple", "Simple"},
		{"multi_word_test", "MultiWordTest"},
		{"", ""},
		{"single", "Single"},
		{"kebab-case-example", "KebabCaseExample"},
		{"mixed_case-example", "MixedCaseExample"},
		{"UPPER_CASE", "UpperCase"},
		{"camelCase", "Camelcase"},
		{"PascalCase", "Pascalcase"},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := ToCamelCase(test.input)
			if result != test.expected {
				t.Errorf("ToCamelCase(%q) = %q, expected %q", test.input, result, test.expected)
			}
		})
	}
}

// Test ProcessSpec with more edge cases for better coverage
func TestProcessSpecExtended(t *testing.T) {
	t.Run("ProcessSpecWithComplexStructure", func(t *testing.T) {
		spec := &models.OpenAPISpec{
			Paths: make(map[string]map[string]models.Operation),
		}

		// Add path with multiple operations
		spec.Paths["/api/users/{user_id}/posts/{post_id}"] = map[string]models.Operation{
			"get": {
				OperationID: "get_user_post",
			},
			"put": {
				OperationID: "update_user_post",
			},
			"delete": {
				OperationID: "delete_user_post",
			},
		}

		// Add another path
		spec.Paths["/api/health-check"] = map[string]models.Operation{
			"get": {
				OperationID: "health_check",
			},
		}

		ProcessSpec(spec)

		// Verify processing worked
		if len(spec.Paths) != 2 {
			t.Errorf("Expected 2 paths after processing, got %d", len(spec.Paths))
		}

		// Check that all operations are still present
		userPostOps := spec.Paths["/api/users/{user_id}/posts/{post_id}"]
		if len(userPostOps) != 3 {
			t.Errorf("Expected 3 operations for user post path, got %d", len(userPostOps))
		}

		healthOps := spec.Paths["/api/health-check"]
		if len(healthOps) != 1 {
			t.Errorf("Expected 1 operation for health check path, got %d", len(healthOps))
		}

		// Verify HTTP methods are added
		for method, op := range userPostOps {
			expectedMethod := strings.ToUpper(method)
			if op.Method != expectedMethod {
				t.Errorf("Expected method %s, got %s", expectedMethod, op.Method)
			}
		}
	})

	t.Run("ProcessSpecWithEmptyPaths", func(t *testing.T) {
		spec := &models.OpenAPISpec{
			Paths: make(map[string]map[string]models.Operation),
		}

		ProcessSpec(spec)

		if len(spec.Paths) != 0 {
			t.Errorf("Expected 0 paths after processing empty spec, got %d", len(spec.Paths))
		}
	})

	t.Run("ProcessSpecWithMissingOperationIDs", func(t *testing.T) {
		spec := &models.OpenAPISpec{
			Paths: map[string]map[string]models.Operation{
				"/users": {
					"get": {
						// No OperationID provided
					},
					"post": {
						// No OperationID provided
					},
				},
				"/": {
					"get": {
						// Root path test
					},
				},
				"/complex/path/with/segments": {
					"delete": {
						// Complex path test
					},
				},
			},
		}

		ProcessSpec(spec)

		// Check generated operation IDs
		usersOps := spec.Paths["/users"]
		if usersOps["get"].OperationID != "getUsers" {
			t.Errorf("Expected operation ID 'getUsers', got '%s'", usersOps["get"].OperationID)
		}
		if usersOps["post"].OperationID != "postUsers" {
			t.Errorf("Expected operation ID 'postUsers', got '%s'", usersOps["post"].OperationID)
		}

		rootOps := spec.Paths["/"]
		// For root path "/", splitting by "/" gives empty parts, so name becomes empty,
		// which results in just the method name
		if rootOps["get"].OperationID != "get" {
			t.Errorf("Expected operation ID 'get', got '%s'", rootOps["get"].OperationID)
		}

		complexOps := spec.Paths["/complex/path/with/segments"]
		if complexOps["delete"].OperationID != "deleteSegments" {
			t.Errorf("Expected operation ID 'deleteSegments', got '%s'", complexOps["delete"].OperationID)
		}
	})

	t.Run("ProcessSpecWithExistingOperationIDs", func(t *testing.T) {
		spec := &models.OpenAPISpec{
			Paths: map[string]map[string]models.Operation{
				"/users": {
					"get": {
						OperationID: "listAllUsers",
					},
				},
			},
		}

		ProcessSpec(spec)

		// Should preserve existing operation ID
		usersOps := spec.Paths["/users"]
		if usersOps["get"].OperationID != "listAllUsers" {
			t.Errorf("Expected operation ID to be preserved as 'listAllUsers', got '%s'", usersOps["get"].OperationID)
		}
	})

	t.Run("ProcessSpecWithTags", func(t *testing.T) {
		spec := &models.OpenAPISpec{
			Paths: map[string]map[string]models.Operation{
				"/users": {
					"get": {
						OperationID: "list_users",
						Tags:        []string{"users", "admin"},
					},
					"post": {
						OperationID: "create_user",
						// No tags - should get default
					},
				},
			},
		}

		ProcessSpec(spec)

		usersOps := spec.Paths["/users"]

		// Check that existing tags are preserved
		getTags := usersOps["get"].Tags
		if len(getTags) != 2 || getTags[0] != "users" || getTags[1] != "admin" {
			t.Errorf("Expected tags [users, admin], got %v", getTags)
		}

		// Check that default tags are added when missing
		postTags := usersOps["post"].Tags
		if len(postTags) != 1 || postTags[0] != "default" {
			t.Errorf("Expected tags [default], got %v", postTags)
		}
	})
}
