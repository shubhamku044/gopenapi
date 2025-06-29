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
			name:     "string type",
			schema:   models.Schema{Type: "string"},
			expected: "string",
		},
		{
			name:     "integer type",
			schema:   models.Schema{Type: "integer"},
			expected: "int",
		},
		{
			name:     "boolean type",
			schema:   models.Schema{Type: "boolean"},
			expected: "bool",
		},
		{
			name:     "date-time format",
			schema:   models.Schema{Type: "string", Format: "date-time"},
			expected: "time.Time",
		},
		{
			name:     "int64 format",
			schema:   models.Schema{Type: "integer", Format: "int64"},
			expected: "int64",
		},
		{
			name:     "float64 number",
			schema:   models.Schema{Type: "number"},
			expected: "float64",
		},
		{
			name:     "array of strings",
			schema:   models.Schema{Type: "array", Items: &models.Schema{Type: "string"}},
			expected: "[]string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetGoType(tt.schema)
			if result != tt.expected {
				t.Errorf("GetGoType() = %v, want %v", result, tt.expected)
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
			name:     "snake_case",
			input:    "user_name",
			expected: "UserName",
		},
		{
			name:     "kebab-case",
			input:    "user-id",
			expected: "UserId",
		},
		{
			name:     "single word",
			input:    "user",
			expected: "User",
		},
		{
			name:     "mixed case",
			input:    "get_user_by_id",
			expected: "GetUserById",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamelCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertPathToGin(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple path",
			input:    "/users",
			expected: "/users",
		},
		{
			name:     "path with parameter",
			input:    "/users/{id}",
			expected: "/users/:id",
		},
		{
			name:     "path with multiple parameters",
			input:    "/users/{userId}/posts/{postId}",
			expected: "/users/:userId/posts/:postId",
		},
		{
			name:     "root path",
			input:    "/",
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertPathToGin(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertPathToGin() = %v, want %v", result, tt.expected)
			}
		})
	}
}
