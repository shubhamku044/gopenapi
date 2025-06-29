package parser

import (
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
