package pipeline

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
	"github.com/rs/zerolog"
)

// mockHTTPClient implements http.Client's Do method for testing
// this is needed for external API calls in the handler
// it allows us to simulate responses and errors without making real HTTP requests (e.g., for Azure DevOps API calls).
type mockHTTPClient struct {
	responses map[string]*http.Response
	errors    map[string]error
	requests  []*http.Request
}

// newMockHTTPClient creates a new instance of mockHTTPClient
// with empty maps for responses and errors
// and an empty slice for requests.
func newMockHTTPClient() *mockHTTPClient {
	return &mockHTTPClient{
		responses: make(map[string]*http.Response),
		errors:    make(map[string]error),
		requests:  make([]*http.Request, 0),
	}
}

// Do implements the http.Client Do method for mockHTTPClient.
// It simulates sending an HTTP request and returns a response or an error
func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Store the request for verification
	m.requests = append(m.requests, req)

	key := req.URL.String()

	// Check if there's an error configured for this URL
	if err, exists := m.errors[key]; exists {
		return nil, err
	}

	// Return configured response or default 404
	if resp, exists := m.responses[key]; exists {
		return resp, nil
	}

	// Default response
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(`{"message": "Not Found"}`)),
		Header:     make(http.Header),
	}, nil
}

// setResponse allows setting a predefined response for a specific URL
func (m *mockHTTPClient) setResponse(url string, statusCode int, body string) {
	m.responses[url] = &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func (m *mockHTTPClient) setError(url string, err error) {
	m.errors[url] = err
}

func (m *mockHTTPClient) getLastRequest() *http.Request {
	if len(m.requests) == 0 {
		return nil
	}
	return m.requests[len(m.requests)-1]
}

func (m *mockHTTPClient) getRequestCount() int {
	return len(m.requests)
}

func (m *mockHTTPClient) reset() {
	m.responses = make(map[string]*http.Response)
	m.errors = make(map[string]error)
	m.requests = make([]*http.Request, 0)
}

// createTestGetHandler creates a GET handler instance for testing with a mock client
func createTestGetHandler(mockClient *mockHTTPClient) *getHandler {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	h := &getHandler{
		baseHandler: &baseHandler{
			HandlerOptions: handlers.HandlerOptions{
				Client: mockClient,
				Log:    &logger,
			},
		},
	}
	return h
}

// createTestDeleteHandler creates a DELETE handler instance for testing with a mock client
func createTestDeleteHandler(mockClient *mockHTTPClient) *deleteHandler {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	h := &deleteHandler{
		baseHandler: &baseHandler{
			HandlerOptions: handlers.HandlerOptions{
				Client: mockClient,
				Log:    &logger,
			},
		},
	}
	return h
}

// createTestPutHandler creates a PUT handler instance for testing with a mock client
func createTestPutHandler(mockClient *mockHTTPClient) *putHandler {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	h := &putHandler{
		baseHandler: &baseHandler{
			HandlerOptions: handlers.HandlerOptions{
				Client: mockClient,
				Log:    &logger,
			},
		},
	}
	return h
}

// Test data constants
const (
	testOrg         = "testorg"
	testProject     = "testproject"
	testPipelineID  = "123"
	testAPIVersion  = "7.2-preview.1"
	testAuthHeader  = "Basic dGVzdDp0ZXN0"
	testUsername    = "test"
	testPassword    = "test"
	buildAPIVersion = "7.2-preview.7"
)

var (
	pipelineGetURL    = fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/pipelines/%s?api-version=%s", testOrg, testProject, testPipelineID, testAPIVersion)
	pipelineDeleteURL = fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/build/definitions/%s?api-version=%s", testOrg, testProject, testPipelineID, buildAPIVersion)
	pipelinePutURL    = fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/build/definitions/%s?api-version=%s", testOrg, testProject, testPipelineID, buildAPIVersion)

	validPipelineResp = `{
		"id": 123,
		"name": "test-pipeline",
		"folder": "\\TestFolder",
		"url": "https://dev.azure.com/testorg/testproject/_apis/pipelines/123",
		"revision": 1,
		"configuration": {
			"type": "yaml",
			"path": "azure-pipelines.yml",
			"repository": {
				"id": "repo123",
				"type": "azureReposGit"
			}
		},
		"_links": {
			"self": {
				"href": "https://dev.azure.com/testorg/testproject/_apis/pipelines/123"
			}
		}
	}`

	validPipelineRespProcessed = `{
		"id": 123,
		"name": "test-pipeline",
		"folder": "TestFolder",
		"url": "https://dev.azure.com/testorg/testproject/_apis/pipelines/123",
		"revision": 1,
		"configuration": {
			"type": "yaml",
			"path": "azure-pipelines.yml",
			"repository": {
				"id": "repo123",
				"type": "azureReposGit"
			}
		},
		"_links": {
			"self": {
				"href": "https://dev.azure.com/testorg/testproject/_apis/pipelines/123"
			}
		}
	}`

	validBuildDefinitionResp = `{
		"id": 123,
		"name": "test-pipeline",
		"path": "\\TestFolder",
		"revision": 2,
		"type": "build",
		"process": {
			"yamlFilename": "azure-pipelines.yml"
		},
		"repository": {
			"id": "repo123",
			"type": "TfsGit"
		},
		"_links": {
			"self": {
				"href": "https://dev.azure.com/testorg/testproject/_apis/build/definitions/123"
			}
		}
	}`

	pipelineNotFoundResp = `{
		"message": "Pipeline not found"
	}`

	unauthorizedResp = `{
		"message": "Unauthorized"
	}`

	validPutRequestBody = `{
		"name": "updated-pipeline",
		"folder": "\\UpdatedFolder",
		"revision": 1,
		"configuration": {
			"type": "yaml",
			"path": "updated-pipeline.yml",
			"repository": {
				"id": "repo456",
				"type": "azureReposGit"
			}
		}
	}`
)

// Test constructor functions
func TestGetPipeline(t *testing.T) {
	t.Run("returns valid handler", func(t *testing.T) {
		client := &http.Client{}
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		opts := handlers.HandlerOptions{
			Client: client,
			Log:    &logger,
		}

		handlerInterface := GetPipeline(opts)

		if handlerInterface == nil {
			t.Fatal("GetPipeline should return a non-nil handler")
		}

		// Verify it implements the Handler interface
		var _ handlers.Handler = handlerInterface

		// Verify the handler has the correct type and options
		h, ok := handlerInterface.(*getHandler)
		if !ok {
			t.Fatal("GetPipeline should return a *getHandler")
		}

		if h.Client != client {
			t.Error("Handler should have the provided client")
		}

		if h.Log != &logger {
			t.Error("Handler should have the provided logger")
		}
	})
}

func TestDeletePipeline(t *testing.T) {
	t.Run("returns valid handler", func(t *testing.T) {
		client := &http.Client{}
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		opts := handlers.HandlerOptions{
			Client: client,
			Log:    &logger,
		}

		handlerInterface := DeletePipeline(opts)

		if handlerInterface == nil {
			t.Fatal("DeletePipeline should return a non-nil handler")
		}

		// Verify it implements the Handler interface
		var _ handlers.Handler = handlerInterface

		// Verify the handler has the correct type and options
		h, ok := handlerInterface.(*deleteHandler)
		if !ok {
			t.Fatal("DeletePipeline should return a *deleteHandler")
		}

		if h.Client != client {
			t.Error("Handler should have the provided client")
		}

		if h.Log != &logger {
			t.Error("Handler should have the provided logger")
		}
	})
}

func TestPutPipeline(t *testing.T) {
	t.Run("returns valid handler", func(t *testing.T) {
		client := &http.Client{}
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		opts := handlers.HandlerOptions{
			Client: client,
			Log:    &logger,
		}

		handlerInterface := PutPipeline(opts)

		if handlerInterface == nil {
			t.Fatal("PutPipeline should return a non-nil handler")
		}

		// Verify it implements the Handler interface
		var _ handlers.Handler = handlerInterface

		// Verify the handler has the correct type and options
		h, ok := handlerInterface.(*putHandler)
		if !ok {
			t.Fatal("PutPipeline should return a *putHandler")
		}

		if h.Client != client {
			t.Error("Handler should have the provided client")
		}

		if h.Log != &logger {
			t.Error("Handler should have the provided logger")
		}
	})
}

// Test GET handler
func TestGetHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name                 string
		organization         string
		project              string
		pipelineID           string
		apiVersion           string
		authHeader           string
		setupMock            func(*mockHTTPClient)
		expectedStatus       int
		expectedContentType  string
		expectedBodyContains string
		expectedRequestCount int
		verifyRequests       func(t *testing.T, mockClient *mockHTTPClient)
	}{
		{
			name:         "successful pipeline retrieval",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(pipelineGetURL, http.StatusOK, validPipelineResp)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: `"id":123`,
			expectedRequestCount: 1,
			verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
				if mockClient.getRequestCount() != 1 {
					t.Errorf("Expected 1 request, got %d", mockClient.getRequestCount())
				}

				req := mockClient.getLastRequest()
				if req.URL.String() != pipelineGetURL {
					t.Errorf("Request URL = %s, want %s", req.URL.String(), pipelineGetURL)
				}
				if req.Header.Get("Authorization") != testAuthHeader {
					t.Errorf("Request Authorization header = %s, want %s", req.Header.Get("Authorization"), testAuthHeader)
				}
				if req.Method != "GET" {
					t.Errorf("Request Method = %s, want GET", req.Method)
				}
			},
		},
		{
			name:         "successful pipeline retrieval with folder processing",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(pipelineGetURL, http.StatusOK, validPipelineResp)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: `"folder":"TestFolder"`, // Should be processed to remove leading backslash
			expectedRequestCount: 1,
		},
		{
			name:         "missing organization parameter",
			organization: "", // This is the key for this test
			project:      testProject,
			pipelineID:   testPipelineID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				// No setup needed as validation happens before API call
			},
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "Organization parameter is required",
			expectedRequestCount: 0,
		},
		{
			name:         "missing project parameter",
			organization: testOrg,
			project:      "", // This is the key for this test
			pipelineID:   testPipelineID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				// No setup needed as validation happens before API call
			},
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "Project parameter is required",
			expectedRequestCount: 0,
		},
		{
			name:         "missing api version parameter",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			apiVersion:   "",
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				// No setup needed as validation happens before API call
			},
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "API version is required",
			expectedRequestCount: 0,
		},
		{
			name:         "missing authorization header",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			apiVersion:   testAPIVersion,
			authHeader:   "",
			setupMock: func(mockClient *mockHTTPClient) {
				// No setup needed as validation happens before API call
			},
			expectedStatus:       http.StatusUnauthorized,
			expectedContentType:  "",
			expectedBodyContains: "Request rejected due to missing or invalid Basic authentication",
			expectedRequestCount: 0,
		},
		{
			name:         "pipeline not found",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(pipelineGetURL, http.StatusNotFound, pipelineNotFoundResp)
			},
			expectedStatus:       http.StatusNotFound,
			expectedContentType:  "",
			expectedBodyContains: fmt.Sprintf("Pipeline with ID %s not found", testPipelineID),
			expectedRequestCount: 1,
		},
		{
			name:         "unauthorized access",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(pipelineGetURL, http.StatusUnauthorized, unauthorizedResp)
			},
			expectedStatus:       http.StatusUnauthorized,
			expectedContentType:  "application/json",
			expectedBodyContains: "Unauthorized",
			expectedRequestCount: 1,
		},
		{
			name:         "network error",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setError(pipelineGetURL, fmt.Errorf("network error"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedContentType:  "",
			expectedBodyContains: "Error getting pipeline",
			expectedRequestCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := newMockHTTPClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			handler := createTestGetHandler(mockClient)

			// Create request
			// We create a dummy URL because the handler gets path values from the context, not by parsing the URL itself.
			url := fmt.Sprintf("/api/placeholder/placeholder/pipelines/placeholder?api-version=%s", tt.apiVersion)
			req := httptest.NewRequest("GET", url, nil)

			// Set path values directly on the request. This is the key to making the validation tests work.
			req.SetPathValue("organization", tt.organization)
			req.SetPathValue("project", tt.project)
			req.SetPathValue("id", tt.pipelineID)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
				// Set BasicAuth for validation check inside the handler
				req.SetBasicAuth(testUsername, testPassword)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute the handler directly
			handler.ServeHTTP(rr, req)

			// Verify status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.expectedStatus)
			}

			// Verify content type if expected
			if tt.expectedContentType != "" {
				contentType := rr.Header().Get("Content-Type")
				if !strings.HasPrefix(contentType, tt.expectedContentType) {
					t.Errorf("handler returned wrong content type: got %v want %v", contentType, tt.expectedContentType)
				}
			}

			// Verify response body contains expected content
			if tt.expectedBodyContains != "" {
				body := rr.Body.String()
				if !strings.Contains(body, tt.expectedBodyContains) {
					t.Errorf("handler response body does not contain expected content.\nGot: %s\nWant to contain: %s", body, tt.expectedBodyContains)
				}
			}

			// Verify request count
			if mockClient.getRequestCount() != tt.expectedRequestCount {
				t.Errorf("expected %d requests, got %d", tt.expectedRequestCount, mockClient.getRequestCount())
			}

			// Run custom request verification if provided
			if tt.verifyRequests != nil {
				tt.verifyRequests(t, mockClient)
			}
		})
	}
}

// Test DELETE handler
func TestDeleteHandler_ServeHTTP(t *testing.T) {
	// Set environment variable for testing
	os.Setenv("BUILD_DEFINITIONS_API_VERSION", buildAPIVersion)
	defer os.Unsetenv("BUILD_DEFINITIONS_API_VERSION")

	tests := []struct {
		name                 string
		organization         string
		project              string
		pipelineID           string
		authHeader           string
		setupMock            func(*mockHTTPClient)
		expectedStatus       int
		expectedContentType  string
		expectedBodyContains string
		expectedRequestCount int
		verifyRequests       func(t *testing.T, mockClient *mockHTTPClient)
	}{
		{
			name:         "successful pipeline deletion",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(pipelineDeleteURL, http.StatusNoContent, "")
			},
			expectedStatus:       http.StatusNoContent,
			expectedContentType:  "",
			expectedBodyContains: "",
			expectedRequestCount: 1,
			verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
				if mockClient.getRequestCount() != 1 {
					t.Errorf("Expected 1 request, got %d", mockClient.getRequestCount())
				}

				req := mockClient.getLastRequest()
				if req.URL.String() != pipelineDeleteURL {
					t.Errorf("Request URL = %s, want %s", req.URL.String(), pipelineDeleteURL)
				}
				if req.Header.Get("Authorization") != testAuthHeader {
					t.Errorf("Request Authorization header = %s, want %s", req.Header.Get("Authorization"), testAuthHeader)
				}
				if req.Method != "DELETE" {
					t.Errorf("Request Method = %s, want DELETE", req.Method)
				}
			},
		},
		{
			name:         "pipeline not found for deletion",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(pipelineDeleteURL, http.StatusNotFound, pipelineNotFoundResp)
			},
			expectedStatus:       http.StatusNotFound,
			expectedContentType:  "",
			expectedBodyContains: fmt.Sprintf("Pipeline with ID %s not found", testPipelineID),
			expectedRequestCount: 1,
		},
		{
			name:         "unauthorized deletion",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(pipelineDeleteURL, http.StatusUnauthorized, unauthorizedResp)
			},
			expectedStatus:       http.StatusUnauthorized,
			expectedContentType:  "application/json",
			expectedBodyContains: "Unauthorized",
			expectedRequestCount: 1,
		},
		{
			name:         "missing authorization header for deletion",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			authHeader:   "",
			setupMock: func(mockClient *mockHTTPClient) {
				// No setup needed as validation happens before API call
			},
			expectedStatus:       http.StatusUnauthorized,
			expectedContentType:  "",
			expectedBodyContains: "Request rejected due to missing or invalid Basic authentication",
			expectedRequestCount: 0,
		},
		{
			name:         "network error during deletion",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setError(pipelineDeleteURL, fmt.Errorf("network error"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedContentType:  "",
			expectedBodyContains: "Error deleting pipeline",
			expectedRequestCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := newMockHTTPClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			handler := createTestDeleteHandler(mockClient)

			// Create request
			req := httptest.NewRequest("DELETE", "/api/placeholder", nil)

			// Set path values
			req.SetPathValue("organization", tt.organization)
			req.SetPathValue("project", tt.project)
			req.SetPathValue("id", tt.pipelineID)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
				req.SetBasicAuth(testUsername, testPassword)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute
			handler.ServeHTTP(rr, req)

			// Verify status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.expectedStatus)
			}

			// Verify content type if expected
			if tt.expectedContentType != "" {
				contentType := rr.Header().Get("Content-Type")
				if !strings.HasPrefix(contentType, tt.expectedContentType) {
					t.Errorf("handler returned wrong content type: got %v want %v", contentType, tt.expectedContentType)
				}
			}

			// Verify response body contains expected content
			if tt.expectedBodyContains != "" {
				body := rr.Body.String()
				if !strings.Contains(body, tt.expectedBodyContains) {
					t.Errorf("handler response body does not contain expected content.\nGot: %s\nWant to contain: %s", body, tt.expectedBodyContains)
				}
			}

			// Verify request count
			if mockClient.getRequestCount() != tt.expectedRequestCount {
				t.Errorf("expected %d requests, got %d", tt.expectedRequestCount, mockClient.getRequestCount())
			}

			// Run custom request verification if provided
			if tt.verifyRequests != nil {
				tt.verifyRequests(t, mockClient)
			}
		})
	}
}

// Test PUT handler
func TestPutHandler_ServeHTTP(t *testing.T) {
	// Set environment variable for testing
	os.Setenv("BUILD_DEFINITIONS_API_VERSION", buildAPIVersion)
	defer os.Unsetenv("BUILD_DEFINITIONS_API_VERSION")

	tests := []struct {
		name                 string
		organization         string
		project              string
		pipelineID           string
		authHeader           string
		requestBody          string
		setupMock            func(*mockHTTPClient)
		expectedStatus       int
		expectedContentType  string
		expectedBodyContains string
		expectedRequestCount int
		verifyRequests       func(t *testing.T, mockClient *mockHTTPClient)
	}{
		{
			name:         "successful pipeline update",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			authHeader:   testAuthHeader,
			requestBody:  validPutRequestBody,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(pipelinePutURL, http.StatusOK, validBuildDefinitionResp)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: `"id":123`,
			expectedRequestCount: 1,
			verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
				req := mockClient.getLastRequest()
				if req.URL.String() != pipelinePutURL {
					t.Errorf("Request URL = %s, want %s", req.URL.String(), pipelinePutURL)
				}
				if req.Method != "PUT" {
					t.Errorf("Request Method = %s, want PUT", req.Method)
				}
				body, _ := io.ReadAll(req.Body)
				// Check for a key part of the transformed request
				if !strings.Contains(string(body), `"name":"updated-pipeline"`) {
					t.Errorf("Request body does not contain expected name. Got: %s", string(body))
				}
				if !strings.Contains(string(body), `"id":123`) {
					t.Errorf("Request body does not contain expected id. Got: %s", string(body))
				}
			},
		},
		{
			name:         "successful update with folder processing",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			authHeader:   testAuthHeader,
			requestBody:  validPutRequestBody,
			setupMock: func(mockClient *mockHTTPClient) {
				// The response from the API has a path with a leading backslash
				mockClient.setResponse(pipelinePutURL, http.StatusOK, validBuildDefinitionResp)
			},
			expectedStatus:      http.StatusOK,
			expectedContentType: "application/json",
			// The final response to the client should have the backslash removed.
			expectedBodyContains: `"folder":"TestFolder"`,
			expectedRequestCount: 1,
		},
		{
			name:                 "missing authorization header",
			organization:         testOrg,
			project:              testProject,
			pipelineID:           testPipelineID,
			authHeader:           "",
			requestBody:          validPutRequestBody,
			setupMock:            nil,
			expectedStatus:       http.StatusUnauthorized,
			expectedBodyContains: "Request rejected due to missing or invalid Basic authentication",
			expectedRequestCount: 0,
		},
		{
			name:                 "invalid request body - bad json",
			organization:         testOrg,
			project:              testProject,
			pipelineID:           testPipelineID,
			authHeader:           testAuthHeader,
			requestBody:          `{"name": "test"`,
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "Invalid JSON in request body",
			expectedRequestCount: 0,
		},
		{
			name:                 "invalid pipeline id - not a number",
			organization:         testOrg,
			project:              testProject,
			pipelineID:           "not-a-number",
			authHeader:           testAuthHeader,
			requestBody:          validPutRequestBody,
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: "Invalid pipeline ID: not-a-number",
			expectedRequestCount: 0,
		},
		{
			name:         "azure devops returns 404",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			authHeader:   testAuthHeader,
			requestBody:  validPutRequestBody,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(pipelinePutURL, http.StatusNotFound, pipelineNotFoundResp)
			},
			expectedStatus:       http.StatusNotFound,
			expectedContentType:  "",
			expectedBodyContains: fmt.Sprintf("Pipeline with ID %s not found", testPipelineID),
			expectedRequestCount: 1,
		},
		{
			name:         "network error",
			organization: testOrg,
			project:      testProject,
			pipelineID:   testPipelineID,
			authHeader:   testAuthHeader,
			requestBody:  validPutRequestBody,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setError(pipelinePutURL, fmt.Errorf("network error"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: "Failed to update pipeline",
			expectedRequestCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := newMockHTTPClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			handler := createTestPutHandler(mockClient)

			// Create request
			req := httptest.NewRequest("PUT", "/api/placeholder", strings.NewReader(tt.requestBody))

			// Set path values
			req.SetPathValue("organization", tt.organization)
			req.SetPathValue("project", tt.project)
			req.SetPathValue("id", tt.pipelineID)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
				req.SetBasicAuth(testUsername, testPassword)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute
			handler.ServeHTTP(rr, req)

			// Verify status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.expectedStatus)
			}

			// Verify content type if expected
			if tt.expectedContentType != "" {
				contentType := rr.Header().Get("Content-Type")
				if !strings.HasPrefix(contentType, tt.expectedContentType) {
					t.Errorf("handler returned wrong content type: got %v want %v", contentType, tt.expectedContentType)
				}
			}

			// Verify response body
			if tt.expectedBodyContains != "" {
				body := rr.Body.String()
				if !strings.Contains(body, tt.expectedBodyContains) {
					t.Errorf("handler response body does not contain expected content.\nGot: %s\nWant to contain: %s", body, tt.expectedBodyContains)
				}
			}

			// Verify request count
			if mockClient.getRequestCount() != tt.expectedRequestCount {
				t.Errorf("expected %d requests, got %d", tt.expectedRequestCount, mockClient.getRequestCount())
			}

			// Run custom request verification
			if tt.verifyRequests != nil {
				tt.verifyRequests(t, mockClient)
			}
		})
	}
}
