package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// Test data structures for consistent testing
var (
	testJSONData = map[string]interface{}{
		"id":   123,
		"name": "John Doe",
		"user": map[string]interface{}{
			"id":          456,
			"username":    "johndoe",
			"html_url":    "https://example.com/johndoe",
			"permissions": []string{"read", "write"},
			"profile": map[string]interface{}{
				"email":    "john@example.com",
				"verified": true,
				"settings": map[string]interface{}{
					"theme":         "dark",
					"notifications": true,
				},
			},
		},
		"metadata": map[string]interface{}{
			"created_at": "2023-01-01",
			"updated_at": "2023-12-01",
		},
	}

	testJSONBytes, _ = json.Marshal(testJSONData)
)

// Mock HTTP Response helper
func createMockResponse(data []byte, statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(data)),
		Header:     make(http.Header),
	}
}

// TestResponseFlattener_FlattenResponse tests the main FlattenResponse method
func TestResponseFlattener_FlattenResponse(t *testing.T) {
	tests := []struct {
		name         string
		mappings     []FieldMapping
		responseData []byte
		statusCode   int
		wantErr      bool
		validate     func(t *testing.T, result []byte)
	}{
		{
			name: "successful flattening with single mapping",
			mappings: []FieldMapping{
				{SourcePath: "user.html_url", TargetKey: "html_url"},
			},
			responseData: testJSONBytes,
			statusCode:   200,
			wantErr:      false,
			validate: func(t *testing.T, result []byte) {
				var flattened map[string]interface{}
				if err := json.Unmarshal(result, &flattened); err != nil {
					t.Fatalf("Failed to unmarshal result: %v", err)
				}

				// Check that original data is preserved
				if flattened["id"] != float64(123) {
					t.Errorf("Original field 'id' not preserved") // something went wrong
				}

				// Check that flattened field exists
				if flattened["html_url"] != "https://example.com/johndoe" {
					t.Errorf("Flattened field 'html_url' incorrect: got %v", flattened["html_url"])
				}
			},
		},
		{
			name: "multiple mappings",
			mappings: []FieldMapping{
				{SourcePath: "user.html_url", TargetKey: "html_url"},
				{SourcePath: "user.permissions", TargetKey: "permissions"},
				{SourcePath: "user.profile.email", TargetKey: "email"},
			},
			responseData: testJSONBytes,
			statusCode:   200,
			wantErr:      false,
			validate: func(t *testing.T, result []byte) {
				var flattened map[string]interface{}
				if err := json.Unmarshal(result, &flattened); err != nil {
					t.Fatalf("Failed to unmarshal result: %v", err)
				}

				// Check that original data is preserved
				if flattened["id"] != float64(123) {
					t.Errorf("Original field 'id' not preserved") // something went wrong
				}

				expectedFields := map[string]interface{}{
					"html_url":    "https://example.com/johndoe",
					"permissions": []interface{}{"read", "write"},
					"email":       "john@example.com",
				}

				for key, expected := range expectedFields {
					if !reflect.DeepEqual(flattened[key], expected) {
						t.Errorf("Field %s: expected %v, got %v", key, expected, flattened[key])
					}
				}
			},
		},
		{
			name:         "empty response body should fail",
			mappings:     []FieldMapping{},
			responseData: []byte(""),
			statusCode:   200,
			wantErr:      true,
		},
		{
			name:         "invalid JSON should fail",
			mappings:     []FieldMapping{},
			responseData: []byte("{invalid json}"),
			statusCode:   200,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rf := &ResponseFlattener{Mappings: tt.mappings}
			resp := createMockResponse(tt.responseData, tt.statusCode)

			result, err := rf.FlattenResponse(resp)

			if (err != nil) != tt.wantErr {
				t.Errorf("FlattenResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestResponseFlattener_FlattenBytes tests the FlattenBytes method directly
func TestResponseFlattener_FlattenBytes(t *testing.T) {
	tests := []struct {
		name     string
		mappings []FieldMapping
		input    []byte
		wantErr  bool
		validate func(t *testing.T, result []byte)
	}{
		{
			name: "basic nested field extraction",
			mappings: []FieldMapping{
				{SourcePath: "user.id", TargetKey: "user_id"},
			},
			input:   testJSONBytes,
			wantErr: false,
			validate: func(t *testing.T, result []byte) {
				var flattened map[string]interface{}
				json.Unmarshal(result, &flattened)

				// Check that original data is preserved
				if flattened["id"] != float64(123) {
					t.Errorf("Original field 'id' not preserved") // something went wrong
				}

				if flattened["user_id"] != float64(456) {
					t.Errorf("Expected user_id to be 456, got %v", flattened["user_id"])
				}
			},
		},
		{
			name: "deep nested field extraction",
			mappings: []FieldMapping{
				{SourcePath: "user.profile.settings.theme", TargetKey: "theme"},
			},
			input:   testJSONBytes,
			wantErr: false,
			validate: func(t *testing.T, result []byte) {
				var flattened map[string]interface{}
				json.Unmarshal(result, &flattened)

				// Check that original data is preserved
				if flattened["id"] != float64(123) {
					t.Errorf("Original field 'id' not preserved") // something went wrong
				}

				if flattened["theme"] != "dark" {
					t.Errorf("Expected theme to be 'dark', got %v", flattened["theme"])
				}
			},
		},
		{
			name: "nonexistent field should error",
			mappings: []FieldMapping{
				{SourcePath: "user.nonexistent", TargetKey: "missing"},
			},
			input:   testJSONBytes,
			wantErr: true,
		},
		{
			name: "nonexistent intermediate field should error",
			mappings: []FieldMapping{
				{SourcePath: "user.nonexistent.field", TargetKey: "missing"},
			},
			input:   testJSONBytes,
			wantErr: true,
		},
		{
			name: "nonexistent root in nested path should error",
			mappings: []FieldMapping{
				{SourcePath: "missing.field.deep", TargetKey: "missing"},
			},
			input:   testJSONBytes,
			wantErr: true,
		},
		{
			name: "invalid path should error",
			mappings: []FieldMapping{
				{SourcePath: "user.id.invalid", TargetKey: "invalid"},
			},
			input:   testJSONBytes,
			wantErr: true,
		},
		{
			name:     "empty mappings should preserve original",
			mappings: []FieldMapping{},
			input:    testJSONBytes,
			wantErr:  false,
			validate: func(t *testing.T, result []byte) {
				var original, flattened map[string]interface{}
				json.Unmarshal(testJSONBytes, &original)
				json.Unmarshal(result, &flattened)

				if !reflect.DeepEqual(original, flattened) {
					t.Error("Empty mappings should preserve original structure")
				}
			},
		},
		{
			name: "overwrite existing field",
			mappings: []FieldMapping{
				{SourcePath: "user.id", TargetKey: "id"}, // This should overwrite the root "id"
			},
			input:   testJSONBytes,
			wantErr: false,
			validate: func(t *testing.T, result []byte) {
				var flattened map[string]interface{}
				json.Unmarshal(result, &flattened)
				if flattened["id"] != float64(456) { // Should be user.id, not root id
					t.Errorf("Expected id to be overwritten with user.id (456), got %v", flattened["id"])
				}
			},
		},
		{
			name:     "invalid JSON should error",
			mappings: []FieldMapping{},
			input:    []byte("{invalid}"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rf := &ResponseFlattener{Mappings: tt.mappings}

			result, err := rf.FlattenBytes(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("FlattenBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestResponseFlattener_extractValue tests the private extractValue method
func TestResponseFlattener_extractValue(t *testing.T) {
	rf := &ResponseFlattener{}

	tests := []struct {
		name     string
		data     map[string]interface{}
		path     string
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "simple field extraction",
			data:     testJSONData,
			path:     "name",
			expected: "John Doe",
			wantErr:  false,
		},
		{
			name:     "nested field extraction",
			data:     testJSONData,
			path:     "user.username",
			expected: "johndoe",
			wantErr:  false,
		},
		{
			name:     "deep nested field extraction",
			data:     testJSONData,
			path:     "user.profile.settings.theme",
			expected: "dark",
			wantErr:  false,
		},
		{
			name:     "array field extraction",
			data:     testJSONData,
			path:     "user.permissions",
			expected: []string{"read", "write"},
			wantErr:  false,
		},
		{
			name:     "boolean field extraction",
			data:     testJSONData,
			path:     "user.profile.verified",
			expected: true,
			wantErr:  false,
		},
		{
			name:    "nonexistent root field",
			data:    testJSONData,
			path:    "nonexistent",
			wantErr: true,
		},
		{
			name:    "nonexistent nested field",
			data:    testJSONData,
			path:    "user.nonexistent",
			wantErr: true,
		},
		{
			name:    "nonexistent intermediate field in path",
			data:    testJSONData,
			path:    "user.nonexistent.field",
			wantErr: true,
		},
		{
			name:    "nonexistent root field in nested path",
			data:    testJSONData,
			path:    "nonexistent.field.deep",
			wantErr: true,
		},
		{
			name:    "invalid path - trying to access field on non-object",
			data:    testJSONData,
			path:    "name.invalid",
			wantErr: true,
		},
		{
			name:    "empty path",
			data:    testJSONData,
			path:    "",
			wantErr: true,
		},
		{
			name:    "path with empty parts",
			data:    testJSONData,
			path:    "user..profile",
			wantErr: true,
		},
		{
			name:    "single dot path",
			data:    testJSONData,
			path:    ".",
			wantErr: true,
		},
		{
			name:    "path with only dots",
			data:    testJSONData,
			path:    "...",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := rf.extractValue(tt.data, tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("extractValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("extractValue() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestFieldMapping tests the FieldMapping struct
func TestFieldMapping(t *testing.T) {
	mapping := FieldMapping{
		SourcePath: "user.profile.email",
		TargetKey:  "user_email",
	}

	if mapping.SourcePath != "user.profile.email" {
		t.Errorf("SourcePath not set correctly")
	}

	if mapping.TargetKey != "user_email" {
		t.Errorf("TargetKey not set correctly")
	}
}

// TestResponseFlattener_EdgeCases tests various edge cases
func TestResponseFlattener_EdgeCases(t *testing.T) {
	t.Run("empty response body", func(t *testing.T) {
		rf := &ResponseFlattener{}
		resp := &http.Response{
			Body: io.NopCloser(strings.NewReader("")), // Empty body
		}

		_, err := rf.FlattenResponse(resp)
		if err == nil {
			t.Error("Expected error when response body is empty")
		}
	})

	t.Run("nil response body causes panic - caught by defer", func(t *testing.T) {
		// This test demonstrates that a nil body will cause a panic
		// In production, we should ensure resp.Body is never nil
		rf := &ResponseFlattener{}
		resp := &http.Response{
			Body: nil, // This will cause panic in io.ReadAll
		}

		// Use a deferred function to catch the panic
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when response body is nil, but no panic occurred")
			}
		}()

		// This should panic
		_, _ = rf.FlattenResponse(resp)
	})

	t.Run("empty JSON object", func(t *testing.T) {
		rf := &ResponseFlattener{
			Mappings: []FieldMapping{
				{SourcePath: "nonexistent", TargetKey: "missing"},
			},
		}

		emptyJSON := []byte("{}")
		_, err := rf.FlattenBytes(emptyJSON)
		if err == nil {
			t.Error("Expected error when trying to extract from empty object")
		}
	})

	t.Run("complex data types preservation", func(t *testing.T) {
		complexData := map[string]interface{}{
			"array": []interface{}{1, 2, 3, "string", true},
			"nested": map[string]interface{}{
				"number": 42.5,
				"null":   nil,
			},
		}

		complexJSON, _ := json.Marshal(complexData)

		rf := &ResponseFlattener{
			Mappings: []FieldMapping{
				{SourcePath: "nested.number", TargetKey: "extracted_number"},
			},
		}

		result, err := rf.FlattenBytes(complexJSON)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		var flattened map[string]interface{}
		json.Unmarshal(result, &flattened)

		// Check array preservation
		if !reflect.DeepEqual(flattened["array"], []interface{}{1.0, 2.0, 3.0, "string", true}) {
			t.Error("Array not preserved correctly")
		}

		// Check extracted number
		if flattened["extracted_number"] != 42.5 {
			t.Errorf("Expected extracted_number to be 42.5, got %v", flattened["extracted_number"])
		}
	})
}

// Benchmark tests
func BenchmarkResponseFlattener_FlattenBytes(b *testing.B) {
	rf := &ResponseFlattener{
		Mappings: []FieldMapping{
			{SourcePath: "user.html_url", TargetKey: "html_url"},
			{SourcePath: "user.permissions", TargetKey: "permissions"},
			{SourcePath: "user.profile.email", TargetKey: "email"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := rf.FlattenBytes(testJSONBytes)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkResponseFlattener_extractValue(b *testing.B) {
	rf := &ResponseFlattener{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := rf.extractValue(testJSONData, "user.profile.settings.theme")
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

// Example tests for documentation
func ExampleResponseFlattener_FlattenResponse() {
	// Create a flattener with field mappings
	rf := &ResponseFlattener{
		Mappings: []FieldMapping{
			{SourcePath: "user.html_url", TargetKey: "html_url"},
			{SourcePath: "user.permissions", TargetKey: "permissions"},
		},
	}

	// Mock response data
	responseData := `{
		"id": 123,
		"user": {
			"html_url": "https://example.com/user",
			"permissions": ["read", "write"]
		}
	}`

	resp := createMockResponse([]byte(responseData), 200)

	result, err := rf.FlattenResponse(resp)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Flattened: %s\n", result)
	// Output will include both original fields and flattened fields at root level
}

func ExampleResponseFlattener_FlattenBytes() {
	rf := &ResponseFlattener{
		Mappings: []FieldMapping{
			{SourcePath: "user.email", TargetKey: "email"},
		},
	}

	jsonData := []byte(`{"user": {"email": "test@example.com"}}`)

	result, err := rf.FlattenBytes(jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Result: %s\n", result)
	// Output: Result: {"email":"test@example.com","user":{"email":"test@example.com"}}
}
