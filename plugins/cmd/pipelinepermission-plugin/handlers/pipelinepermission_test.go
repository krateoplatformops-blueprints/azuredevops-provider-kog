package pipelinepermission

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
		},
		nil
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

// Test data constants
const (
	testOrg          = "testorg"
	testProject      = "testproject"
	testResourceType = "pipelines"
	testResourceID   = "123"
	testAPIVersion   = "7.2-preview.1"
	testAuthHeader   = "Basic dGVzdDp0ZXN0"
	testUsername     = "test"
	testPassword     = "test"
)

var (
	permissionGetURL = fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/pipelines/pipelinepermissions/%s/%s?api-version=%s", testOrg, testProject, testResourceType, testResourceID, testAPIVersion)

	validPermissionRespWithAllPipelines    = `{"resource":{"type":"pipeline","id":"123"},"pipelines":[],"allPipelines":{"authorized":true}}`
	validPermissionRespWithoutAllPipelines = `{"resource":{"type":"pipeline","id":"123"},"pipelines":[]}`
	permissionNotFoundResp                 = `{"message": "Resource not found"}`
	unauthorizedResp                       = `{"message": "Unauthorized"}`
)

// Test constructor functions
func TestGetPipelinePermission(t *testing.T) {
	t.Run("returns valid handler", func(t *testing.T) {
		client := &http.Client{}
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		opts := handlers.HandlerOptions{
			Client: client,
			Log:    &logger,
		}

		handlerInterface := GetPipelinePermission(opts)

		if handlerInterface == nil {
			t.Fatal("GetPipelinePermission should return a non-nil handler")
		}

		// Verify it implements the Handler interface
		var _ handlers.Handler = handlerInterface

		// Verify the handler has the correct type and options
		h, ok := handlerInterface.(*getHandler)
		if !ok {
			t.Fatal("GetPipelinePermission should return a *getHandler")
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
		resourceType         string
		resourceID           string
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
			name:         "successful retrieval with allPipelines field present",
			organization: testOrg,
			project:      testProject,
			resourceType: testResourceType,
			resourceID:   testResourceID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(permissionGetURL, http.StatusOK, validPermissionRespWithAllPipelines)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: `"allPipelines":{"authorized":true}`,
			expectedRequestCount: 1,
			verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
				if mockClient.getRequestCount() != 1 {
					t.Errorf("Expected 1 request, got %d", mockClient.getRequestCount())
				}
				req := mockClient.getLastRequest()
				if req.URL.String() != permissionGetURL {
					t.Errorf("Request URL = %s, want %s", req.URL.String(), permissionGetURL)
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
			name:         "successful retrieval without allPipelines field present",
			organization: testOrg,
			project:      testProject,
			resourceType: testResourceType,
			resourceID:   testResourceID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(permissionGetURL, http.StatusOK, validPermissionRespWithoutAllPipelines)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: `"allPipelines":{"authorized":false}`,
			expectedRequestCount: 1,
		},
		{
			name:                 "missing authorization header",
			organization:         testOrg,
			project:              testProject,
			resourceType:         testResourceType,
			resourceID:           testResourceID,
			apiVersion:           testAPIVersion,
			authHeader:           "",
			setupMock:            nil,
			expectedStatus:       http.StatusUnauthorized,
			expectedBodyContains: "Request rejected due to missing or invalid Basic authentication",
			expectedRequestCount: 0,
		},
		{
			name:         "resource not found",
			organization: testOrg,
			project:      testProject,
			resourceType: testResourceType,
			resourceID:   testResourceID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(permissionGetURL, http.StatusNotFound, permissionNotFoundResp)
			},
			expectedStatus:       http.StatusNotFound,
			expectedContentType:  "",
			expectedBodyContains: fmt.Sprintf("Pipeline permission for resource %s/%s/%s/%s not found", testOrg, testProject, testResourceType, testResourceID),
			expectedRequestCount: 1,
		},
		{
			name:         "unauthorized access",
			organization: testOrg,
			project:      testProject,
			resourceType: testResourceType,
			resourceID:   testResourceID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(permissionGetURL, http.StatusUnauthorized, unauthorizedResp)
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
			resourceType: testResourceType,
			resourceID:   testResourceID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setError(permissionGetURL, fmt.Errorf("network error"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: "Error getting pipeline permission",
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
			url := fmt.Sprintf("/api/placeholder/placeholder/pipelines/pipelinepermissions/placeholder/placeholder?api-version=%s", tt.apiVersion)
			req := httptest.NewRequest("GET", url, nil)

			// Set path values
			req.SetPathValue("organization", tt.organization)
			req.SetPathValue("project", tt.project)
			req.SetPathValue("resourceType", tt.resourceType)
			req.SetPathValue("resourceId", tt.resourceID)

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
