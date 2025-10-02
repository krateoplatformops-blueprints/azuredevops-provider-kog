package pipeline

import (
	"encoding/json"
	"fmt"
	"strings"
)

// function to add a field to the response body
func AddFieldToBody(body []byte, fieldName string, fieldValue interface{}) ([]byte, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	data[fieldName] = fieldValue

	return json.Marshal(data)
}

// function to read the field from a body and return the value (generic function)
func ReadFieldFromBody(body []byte, fieldName string) (interface{}, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check if the field exists in the data map
	value, exists := data[fieldName]
	if !exists {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	return value, nil
}

func normalizeRepoTypeFromAzure(typ string) string {
	switch typ {
	case "TfsGit":
		return RepositoryTypeAzureReposGit
	case "GitHub":
		return RepositoryTypeGitHub
	default:
		return strings.TrimSpace(strings.ToLower(typ))
	}
}

func normalizeRepoTypeToAzure(typ string) string {
	switch typ {
	case RepositoryTypeAzureReposGit:
		return "TfsGit"
	case RepositoryTypeGitHub:
		return "GitHub"
	default:
		return strings.TrimSpace(strings.ToLower(typ))
	}
}
