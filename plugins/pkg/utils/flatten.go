package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// FieldMapping defines how to extract and rename fields
type FieldMapping struct {
	SourcePath string // e.g., "user.permissions", "user.html_url", "user.id"
	TargetKey  string // e.g., "permissions", "html_url", "id"
}

// ResponseFlattener handles flattening of HTTP response bodies
type ResponseFlattener struct {
	Mappings []FieldMapping
}

// FlattenResponse reads and flattens an HTTP response body
func (rf *ResponseFlattener) FlattenResponse(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return rf.FlattenBytes(body)
}

// FlattenBytes flattens a JSON byte array
func (rf *ResponseFlattener) FlattenBytes(body []byte) ([]byte, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Start with the complete original data to preserve everything
	flattened := make(map[string]interface{})

	// Copy the entire original response
	for key, value := range data {
		flattened[key] = value
	}

	// Then, add the flattened nested fields to root level
	for _, mapping := range rf.Mappings {
		value, err := rf.extractValue(data, mapping.SourcePath)
		if err != nil {
			return nil, fmt.Errorf("failed to extract %s: %w", mapping.SourcePath, err)
		}
		// Add/override at root level
		flattened[mapping.TargetKey] = value
	}

	return json.Marshal(flattened)
}

// extractValue extracts a value from nested map using dot notation path
func (rf *ResponseFlattener) extractValue(data map[string]interface{}, path string) (interface{}, error) {

	if path == "" {
		return nil, fmt.Errorf("empty path not allowed")
	}

	parts := strings.Split(path, ".")
	current := data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - return the value
			value, exists := current[part]
			if !exists {
				return nil, fmt.Errorf("field %s not found", part)
			}
			return value, nil
		}

		// Intermediate part - navigate deeper
		next, exists := current[part]
		if !exists {
			return nil, fmt.Errorf("field %s not found in path %s", part, path)
		}

		nextMap, ok := next.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("field %s is not an object in path %s", part, path)
		}
		current = nextMap
	}

	return nil, fmt.Errorf("invalid path %s", path)
}
