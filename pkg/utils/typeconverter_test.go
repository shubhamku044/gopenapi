package utils

import (
	"testing"

	"github.com/shubhamku044/gopenapi/internal/models"
)

func TestGetGoType(t *testing.T) {
	tests := []struct {
		name     string
		schema   models.Schema
		expected string
	}{
		{
			name:     "String type",
			schema:   models.Schema{Type: "string"},
			expected: "string",
		},
		{
			name:     "String with date-time format",
			schema:   models.Schema{Type: "string", Format: "date-time"},
			expected: "time.Time",
		},
		{
			name:     "String with date format",
			schema:   models.Schema{Type: "string", Format: "date"},
			expected: "time.Time",
		},
		{
			name:     "String with byte format",
			schema:   models.Schema{Type: "string", Format: "byte"},
			expected: "[]byte",
		},
		{
			name:     "String with binary format",
			schema:   models.Schema{Type: "string", Format: "binary"},
			expected: "[]byte",
		},
		{
			name:     "String with email format",
			schema:   models.Schema{Type: "string", Format: "email"},
			expected: "string",
		},
		{
			name:     "String with uuid format",
			schema:   models.Schema{Type: "string", Format: "uuid"},
			expected: "string",
		},
		{
			name:     "Integer type",
			schema:   models.Schema{Type: "integer"},
			expected: "int",
		},
		{
			name:     "Integer with int32 format",
			schema:   models.Schema{Type: "integer", Format: "int32"},
			expected: "int32",
		},
		{
			name:     "Integer with int64 format",
			schema:   models.Schema{Type: "integer", Format: "int64"},
			expected: "int64",
		},
		{
			name:     "Number type",
			schema:   models.Schema{Type: "number"},
			expected: "float64",
		},
		{
			name:     "Number with float format",
			schema:   models.Schema{Type: "number", Format: "float"},
			expected: "float32",
		},
		{
			name:     "Number with double format",
			schema:   models.Schema{Type: "number", Format: "double"},
			expected: "float64",
		},
		{
			name:     "Boolean type",
			schema:   models.Schema{Type: "boolean"},
			expected: "bool",
		},
		{
			name:     "Array of strings",
			schema:   models.Schema{Type: "array", Items: &models.Schema{Type: "string"}},
			expected: "[]string",
		},
		{
			name:     "Array of integers",
			schema:   models.Schema{Type: "array", Items: &models.Schema{Type: "integer"}},
			expected: "[]int",
		},
		{
			name:     "Array of objects (interface{})",
			schema:   models.Schema{Type: "array", Items: &models.Schema{Type: "object"}},
			expected: "[]map[string]interface{}",
		},
		{
			name:     "Array with reference",
			schema:   models.Schema{Type: "array", Items: &models.Schema{Ref: "#/components/schemas/User"}},
			expected: "[]User",
		},
		{
			name:     "Array without items (fallback)",
			schema:   models.Schema{Type: "array"},
			expected: "[]interface{}",
		},
		{
			name:     "Object type",
			schema:   models.Schema{Type: "object"},
			expected: "map[string]interface{}",
		},
		{
			name:     "Reference to User",
			schema:   models.Schema{Ref: "#/components/schemas/User"},
			expected: "User",
		},
		{
			name:     "Reference to nested schema",
			schema:   models.Schema{Ref: "#/components/schemas/api/v1/UserProfile"},
			expected: "UserProfile",
		},
		{
			name:     "Empty reference",
			schema:   models.Schema{Ref: ""},
			expected: "interface{}",
		},
		{
			name:     "Unknown type",
			schema:   models.Schema{Type: "unknown"},
			expected: "interface{}",
		},
		{
			name:     "Empty schema",
			schema:   models.Schema{},
			expected: "interface{}",
		},
		{
			name:     "Nested array of references",
			schema:   models.Schema{Type: "array", Items: &models.Schema{Type: "array", Items: &models.Schema{Ref: "#/components/schemas/Tag"}}},
			expected: "[][]Tag",
		},
		{
			name:     "Array of date-time strings",
			schema:   models.Schema{Type: "array", Items: &models.Schema{Type: "string", Format: "date-time"}},
			expected: "[]time.Time",
		},
		{
			name:     "Array of byte arrays",
			schema:   models.Schema{Type: "array", Items: &models.Schema{Type: "string", Format: "byte"}},
			expected: "[][]byte",
		},
		{
			name:     "Complex reference path",
			schema:   models.Schema{Ref: "#/components/schemas/api/v2/user/Profile"},
			expected: "Profile",
		},
		{
			name:     "Reference with special characters",
			schema:   models.Schema{Ref: "#/components/schemas/User-Profile"},
			expected: "User-Profile",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GetGoType(test.schema)
			if result != test.expected {
				t.Errorf("GetGoType(%+v) = %q, expected %q", test.schema, result, test.expected)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Snake case",
			input:    "hello_world",
			expected: "HelloWorld",
		},
		{
			name:     "Kebab case",
			input:    "hello-world",
			expected: "HelloWorld",
		},
		{
			name:     "Mixed separators",
			input:    "hello_world-test",
			expected: "HelloWorldTest",
		},
		{
			name:     "Single word",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Already camelCase",
			input:    "helloWorld",
			expected: "Helloworld",
		},
		{
			name:     "Already PascalCase",
			input:    "HelloWorld",
			expected: "Helloworld",
		},
		{
			name:     "Numbers in string",
			input:    "user_id_v2",
			expected: "UserIdV2",
		},
		{
			name:     "Special characters",
			input:    "api_v1_user",
			expected: "ApiV1User",
		},
		{
			name:     "Multiple underscores",
			input:    "hello__world",
			expected: "HelloWorld",
		},
		{
			name:     "Multiple hyphens",
			input:    "hello--world",
			expected: "HelloWorld",
		},
		{
			name:     "Starting with separator",
			input:    "_hello_world",
			expected: "HelloWorld",
		},
		{
			name:     "Ending with separator",
			input:    "hello_world_",
			expected: "HelloWorld",
		},
		{
			name:     "Only separators",
			input:    "___",
			expected: "",
		},
		{
			name:     "Single character",
			input:    "a",
			expected: "A",
		},
		{
			name:     "Uppercase input",
			input:    "HELLO_WORLD",
			expected: "HelloWorld",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ToCamelCase(test.input)
			if result != test.expected {
				t.Errorf("ToCamelCase(%q) = %q, expected %q", test.input, result, test.expected)
			}
		})
	}
}

func TestConvertPathToGin(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "Simple path parameter",
			path:     "/users/{id}",
			expected: "/users/:id",
		},
		{
			name:     "Multiple path parameters",
			path:     "/users/{userId}/posts/{postId}",
			expected: "/users/:userId/posts/:postId",
		},
		{
			name:     "No path parameters",
			path:     "/users",
			expected: "/users",
		},
		{
			name:     "Path with query-like syntax (should not change)",
			path:     "/users?active=true",
			expected: "/users?active=true",
		},
		{
			name:     "Complex path",
			path:     "/api/v1/users/{userId}/posts/{postId}/comments/{commentId}",
			expected: "/api/v1/users/:userId/posts/:postId/comments/:commentId",
		},
		{
			name:     "Empty path",
			path:     "",
			expected: "",
		},
		{
			name:     "Root path",
			path:     "/",
			expected: "/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ConvertPathToGin(test.path)
			if result != test.expected {
				t.Errorf("ConvertPathToGin(%q) = %q, expected %q", test.path, result, test.expected)
			}
		})
	}
}
