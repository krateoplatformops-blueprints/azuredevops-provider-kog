package project

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
	"github.com/rs/zerolog"
)

// mockHTTPClient implements http.Client's Do method for testing.
// It stores mock responses and allows multiple requests to the same URL,
// which is critical for testing polling operations.
type mockHTTPClient struct {
	responses map[string]mockResponse
	errors    map[string]error
	requests  []*http.Request
}

// mockResponse stores HTTP response data as plain values instead of using *http.Response directly.
// This is necessary because http.Response.Body is an io.ReadCloser that can only be read once (consumed on first read).
// PATCH and DELETE handlers poll operation URLs multiple times (up to 10 attempts),
// making the same GET request repeatedly to check operation status.
//
// mockResponse:
//   - Stores body as a string (immutable data, not a stream)
//   - Creates a fresh io.NopCloser(strings.NewReader(body)) on each Do() call
//   - Each request gets a new reader, allowing unlimited re-reads of the same data
type mockResponse struct {
	statusCode int
	body       string // Stored as string to allow multiple reads
	header     http.Header
}

// newMockHTTPClient creates a new instance of mockHTTPClient
func newMockHTTPClient() *mockHTTPClient {
	return &mockHTTPClient{
		responses: make(map[string]mockResponse),
		errors:    make(map[string]error),
		requests:  make([]*http.Request, 0),
	}
}

// Do implements the http.Client Do method for mockHTTPClient
func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.requests = append(m.requests, req)

	// Include HTTP method in key to differentiate PATCH vs GET on same URL
	key := req.Method + " " + req.URL.String()

	if err, exists := m.errors[key]; exists {
		return nil, err
	}

	if mockResp, exists := m.responses[key]; exists {
		// Create a fresh reader each time to allow multiple reads
		header := mockResp.header
		if header == nil {
			header = make(http.Header)
		}
		return &http.Response{
			StatusCode: mockResp.statusCode,
			Body:       io.NopCloser(strings.NewReader(mockResp.body)),
			Header:     header,
		}, nil
	}

	fmt.Printf("MockHTTPClient: No response configured for %s %s. Returning 404.\n", req.Method, req.URL.String())
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(`{"message": "Not Found"}`)),
		Header:     make(http.Header),
	}, nil
}

// setResponse allows setting a predefined response for a specific HTTP method and URL
func (m *mockHTTPClient) setResponse(method, url string, statusCode int, body string) {
	key := method + " " + url
	m.responses[key] = mockResponse{
		statusCode: statusCode,
		body:       body,
		header:     make(http.Header),
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

func (m *mockHTTPClient) reset() {
	m.responses = make(map[string]mockResponse)
	m.errors = make(map[string]error)
	m.requests = make([]*http.Request, 0)
}

// createTestPostHandler creates a POST handler instance for testing with a mock client
func createTestPostHandler(mockClient *mockHTTPClient) *postHandler {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	h := &postHandler{
		baseHandler: &baseHandler{
			HandlerOptions: handlers.HandlerOptions{
				Client: mockClient,
				Log:    &logger,
			},
		},
	}
	return h
}

// createTestPatchHandler creates a PATCH handler instance for testing with a mock client
func createTestPatchHandler(mockClient *mockHTTPClient) *patchHandler {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	h := &patchHandler{
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

// Test data constants
const (
	testOrg        = "test-org"
	testProjectID  = "proj-123"
	testAPIVersion = "7.2-preview.4"
	testAuthHeader = "Basic dGVzdDp0ZXN0"
	testOpID       = "a1b2c3d4-5678-90ab-cdef-1234567890ab"
)

var (
	// Common request bodies
	validProjectBody = `{"name": "test-project", "description": "test"}`
	validUpdateBody  = `{"name": "updated-project", "description": "updated"}`

	// Common response bodies
	operationReference = `{"id":"` + testOpID + `","status":"queued","url":"https://dev.azure.com/` + testOrg + `/_apis/operations/` + testOpID + `","pluginId":"12345678-1234-1234-1234-123456789abc"}`
	operationSucceeded = `{"id":"` + testOpID + `","status":"succeeded","url":"https://dev.azure.com/` + testOrg + `/_apis/operations/` + testOpID + `"}`
	operationFailed    = `{"id":"op-fail-123","status":"failed","url":"https://dev.azure.com/` + testOrg + `/_apis/operations/op-fail-123"}`
	projectResponse    = `{"id":"` + testProjectID + `","name":"updated-project","description":"updated","state":"wellFormed"}`

	// Common URLs (using buildProjectURL and buildOperationURL from production code)
	testPostURL  = buildProjectURL(testOrg, "", testAPIVersion)
	testPatchURL = buildProjectURL(testOrg, testProjectID, testAPIVersion)
	testOpURL    = buildOperationURL(testOrg, testOpID)
)

// TestPostProject_ServeHTTP tests the POST handler using table-driven tests
func TestPostProject_ServeHTTP(t *testing.T) {
	tests := []struct {
		name                 string
		organization         string
		apiVersion           string
		authHeader           string
		requestBody          string
		setupMock            func(*mockHTTPClient)
		expectedStatus       int
		expectedBodyContains []string
		verifyResponse       func(t *testing.T, body string)
	}{
		{
			name:         "successful project creation",
			organization: testOrg,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validProjectBody,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(http.MethodPost, testPostURL, http.StatusOK, operationReference)
			},
			expectedStatus:       http.StatusAccepted,
			expectedBodyContains: []string{"operationId", "operationStatus", "operationUrl", "operationPluginId"},
			verifyResponse: func(t *testing.T, body string) {
				// Verify original fields NOT present
				if strings.Contains(body, `"id"`) && !strings.Contains(body, `"operationId"`) {
					t.Errorf("Response should NOT contain original 'id' field, got: %s", body)
				}
				if strings.Contains(body, `"status"`) && !strings.Contains(body, `"operationStatus"`) {
					t.Errorf("Response should NOT contain original 'status' field (only operationStatus), got: %s", body)
				}
			},
		},
		{
			name:                 "missing organization",
			organization:         "",
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			requestBody:          "{}",
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: []string{"organization parameter is required"},
		},
		{
			name:                 "missing api version",
			organization:         testOrg,
			apiVersion:           "",
			authHeader:           testAuthHeader,
			requestBody:          "{}",
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: []string{"api-version is required"},
		},
		{
			name:                 "missing authorization",
			organization:         testOrg,
			apiVersion:           testAPIVersion,
			authHeader:           "",
			requestBody:          "{}",
			setupMock:            nil,
			expectedStatus:       http.StatusUnauthorized,
			expectedBodyContains: []string{"authentication"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockClient := newMockHTTPClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}
			handler := createTestPostHandler(mockClient)

			// Create request
			url := fmt.Sprintf("/plugin/project/%s?api-version=%s", tt.organization, tt.apiVersion)
			req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(tt.requestBody))
			req.SetPathValue("organization", tt.organization)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Execute
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Verify status
			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d. Body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}

			// Verify body contains expected strings
			body := w.Body.String()
			for _, expected := range tt.expectedBodyContains {
				if !strings.Contains(body, expected) {
					t.Errorf("body missing %q, got: %s", expected, body)
				}
			}

			// Custom verification
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, body)
			}
		})
	}
}

// TestPatchProject_ServeHTTP tests the PATCH handler using table-driven tests
func TestPatchProject_ServeHTTP(t *testing.T) {
	tests := []struct {
		name                 string
		organization         string
		projectID            string
		apiVersion           string
		authHeader           string
		requestBody          string
		setupMock            func(*mockHTTPClient)
		expectedStatus       int
		expectedBodyContains []string
	}{
		{
			name:         "successful update with polling",
			organization: testOrg,
			projectID:    testProjectID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validUpdateBody,
			setupMock: func(mockClient *mockHTTPClient) {
				// Step 1: PATCH returns operation reference
				mockClient.setResponse(http.MethodPatch, testPatchURL, http.StatusAccepted, operationReference)
				// Step 2: Operation polling succeeds
				mockClient.setResponse(http.MethodGet, testOpURL, http.StatusOK, operationSucceeded)
				// Step 3: Final project GET
				mockClient.setResponse(http.MethodGet, testPatchURL, http.StatusOK, projectResponse)
			},
			expectedStatus:       http.StatusOK,
			expectedBodyContains: []string{`"id":"proj-123"`, `"name":"updated-project"`, `"state":"wellFormed"`},
		},
		{
			name:         "operation polling failed",
			organization: testOrg,
			projectID:    testProjectID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validUpdateBody,
			setupMock: func(mockClient *mockHTTPClient) {
				// Step 1: PATCH returns operation reference
				failOpRef := `{"id":"op-fail-123","status":"queued","url":"https://dev.azure.com/` + testOrg + `/_apis/operations/op-fail-123"}`
				mockClient.setResponse(http.MethodPatch, testPatchURL, http.StatusAccepted, failOpRef)
				// Step 2: Operation polling returns failed status
				failOpURL := buildOperationURL(testOrg, "op-fail-123")
				mockClient.setResponse(http.MethodGet, failOpURL, http.StatusOK, operationFailed)
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: []string{"operation polling failed"},
		},
		{
			name:                 "missing projectId",
			organization:         testOrg,
			projectID:            "",
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			requestBody:          "{}",
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: []string{"projectId parameter is required"},
		},
		{
			name:                 "missing organization",
			organization:         "",
			projectID:            testProjectID,
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			requestBody:          "{}",
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: []string{"organization parameter is required"},
		},
		{
			name:                 "missing api version",
			organization:         testOrg,
			projectID:            testProjectID,
			apiVersion:           "",
			authHeader:           testAuthHeader,
			requestBody:          "{}",
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: []string{"api-version is required"},
		},
		{
			name:                 "missing authorization",
			organization:         testOrg,
			projectID:            testProjectID,
			apiVersion:           testAPIVersion,
			authHeader:           "",
			requestBody:          "{}",
			setupMock:            nil,
			expectedStatus:       http.StatusUnauthorized,
			expectedBodyContains: []string{"authentication"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockClient := newMockHTTPClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}
			handler := createTestPatchHandler(mockClient)

			// Create request
			url := fmt.Sprintf("/plugin/project/%s/%s?api-version=%s", tt.organization, tt.projectID, tt.apiVersion)
			req := httptest.NewRequest(http.MethodPatch, url, strings.NewReader(tt.requestBody))
			req.SetPathValue("organization", tt.organization)
			req.SetPathValue("projectId", tt.projectID)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Execute
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Verify status
			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d. Body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}

			// Verify body contains expected strings
			body := w.Body.String()
			for _, expected := range tt.expectedBodyContains {
				if !strings.Contains(body, expected) {
					t.Errorf("body missing %q, got: %s", expected, body)
				}
			}
		})
	}
}

// TestDeleteProject_ServeHTTP tests the DELETE handler using table-driven tests
func TestDeleteProject_ServeHTTP(t *testing.T) {
	tests := []struct {
		name                 string
		organization         string
		projectID            string
		apiVersion           string
		authHeader           string
		setupMock            func(*mockHTTPClient)
		expectedStatus       int
		expectedBodyContains []string
		verifyResponse       func(t *testing.T, body string)
	}{
		{
			name:         "successful deletion with polling",
			organization: testOrg,
			projectID:    testProjectID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				// Step 1: DELETE returns operation reference
				delOpRef := `{"id":"del-op-123","status":"queued","url":"https://dev.azure.com/` + testOrg + `/_apis/operations/del-op-123"}`
				mockClient.setResponse(http.MethodDelete, testPatchURL, http.StatusAccepted, delOpRef)
				// Step 2: Operation polling succeeds
				delOpURL := buildOperationURL(testOrg, "del-op-123")
				delOpResp := `{"id":"del-op-123","status":"succeeded","url":"https://dev.azure.com/` + testOrg + `/_apis/operations/del-op-123"}`
				mockClient.setResponse(http.MethodGet, delOpURL, http.StatusOK, delOpResp)
			},
			expectedStatus: http.StatusNoContent,
			verifyResponse: func(t *testing.T, body string) {
				if len(body) != 0 {
					t.Errorf("Expected empty response body, got: %s", body)
				}
			},
		},
		{
			name:         "operation polling failed",
			organization: testOrg,
			projectID:    testProjectID,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				// Step 1: DELETE returns operation reference
				delFailOpRef := `{"id":"del-fail-123","status":"queued","url":"https://dev.azure.com/` + testOrg + `/_apis/operations/del-fail-123"}`
				mockClient.setResponse(http.MethodDelete, testPatchURL, http.StatusAccepted, delFailOpRef)
				// Step 2: Operation polling returns failed status
				delFailOpURL := buildOperationURL(testOrg, "del-fail-123")
				delFailOpResp := `{"id":"del-fail-123","status":"failed","url":"https://dev.azure.com/` + testOrg + `/_apis/operations/del-fail-123"}`
				mockClient.setResponse(http.MethodGet, delFailOpURL, http.StatusOK, delFailOpResp)
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedBodyContains: []string{"operation polling failed"},
		},
		{
			name:                 "missing projectId",
			organization:         testOrg,
			projectID:            "",
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: []string{"projectId parameter is required"},
		},
		{
			name:                 "missing organization",
			organization:         "",
			projectID:            testProjectID,
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: []string{"organization parameter is required"},
		},
		{
			name:                 "missing api version",
			organization:         testOrg,
			projectID:            testProjectID,
			apiVersion:           "",
			authHeader:           testAuthHeader,
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedBodyContains: []string{"api-version is required"},
		},
		{
			name:                 "missing authorization",
			organization:         testOrg,
			projectID:            testProjectID,
			apiVersion:           testAPIVersion,
			authHeader:           "",
			setupMock:            nil,
			expectedStatus:       http.StatusUnauthorized,
			expectedBodyContains: []string{"authentication"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			mockClient := newMockHTTPClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}
			handler := createTestDeleteHandler(mockClient)

			// Create request
			url := fmt.Sprintf("/plugin/project/%s/%s?api-version=%s", tt.organization, tt.projectID, tt.apiVersion)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.SetPathValue("organization", tt.organization)
			req.SetPathValue("projectId", tt.projectID)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Execute
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			// Verify status
			if w.Code != tt.expectedStatus {
				t.Errorf("status = %d, want %d. Body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}

			// Verify body contains expected strings
			body := w.Body.String()
			for _, expected := range tt.expectedBodyContains {
				if !strings.Contains(body, expected) {
					t.Errorf("body missing %q, got: %s", expected, body)
				}
			}

			// Custom verification
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, body)
			}
		})
	}
}

func TestTransformOperationReference(t *testing.T) {
	opRef := &OperationReference{
		ID:       "test-id-123",
		PluginID: "plugin-456",
		Status:   "queued",
		URL:      "https://example.com/operation/test-id-123",
	}

	transformed := transformOperationReference(opRef)

	if transformed.OperationID != opRef.ID {
		t.Errorf("Expected OperationID %s, got %s", opRef.ID, transformed.OperationID)
	}
	if transformed.OperationPluginID != opRef.PluginID {
		t.Errorf("Expected OperationPluginID %s, got %s", opRef.PluginID, transformed.OperationPluginID)
	}
	if transformed.OperationStatus != opRef.Status {
		t.Errorf("Expected OperationStatus %s, got %s", opRef.Status, transformed.OperationStatus)
	}
	if transformed.OperationURL != opRef.URL {
		t.Errorf("Expected OperationURL %s, got %s", opRef.URL, transformed.OperationURL)
	}
}
