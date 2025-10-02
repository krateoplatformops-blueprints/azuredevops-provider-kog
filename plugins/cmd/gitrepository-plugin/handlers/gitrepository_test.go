package gitrepository

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
	fmt.Printf("MockHTTPClient: No response configured for URL: %s. Returning 404.\n", key)
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

// Test data constants
const (
	testOrg                 = "testorg"
	testProject             = "testproject"
	testRepoID              = "test-repo-id"
	testAPIVersion          = "7.2-preview.2"
	testGitPushesAPIVersion = "7.2-preview.3"
	testAuthHeader          = "Basic dGVzdDp0ZXN0"
	testUsername            = "test"
	testPassword            = "test"
)

var (
	repoCreateURL = fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories?api-version=%s", testOrg, testProject, testAPIVersion)
	repoUpdateURL = fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s?api-version=%s", testOrg, testProject, testRepoID, testAPIVersion)
	repoPushesURL = fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s/pushes?api-version=%s", testOrg, testProject, testRepoID, testGitPushesAPIVersion)
	repoRefsURL   = fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s/refs?filter=heads/main&api-version=%s", testOrg, testProject, testRepoID, testAPIVersion)

	validCreateRepoReqBody = `{
		"name": "test-repo",
		"initialize": false
	}`

	validCreateRepoReqBodyWithInit = `{
		"name": "test-repo",
		"initialize": true
	}`

	validCreateRepoReqBodyWithDefaultBranch = `{
		"name": "test-repo",
		"initialize": true,
		"defaultBranch": "refs/heads/feature"
	}`

	validCreateRepoReqBodyWithDefaultBranchNoInit = `{
		"name": "test-repo",
		"initialize": false,
		"defaultBranch": "refs/heads/feature"
	}`

	validCreateRepoReqBodyFork = `{
		"name": "fork-repo",
		"parentRepository": {
			"id": "parent-repo-id",
			"project": {
				"id": "parent-project-id"
			}
		}
	}`

	validCreateRepoReqBodyForkWithDefaultBranch = `{
		"name": "fork-repo",
		"parentRepository": {
			"id": "parent-repo-id",
			"project": {
				"id": "parent-project-id"
			}
		},
		"defaultBranch": "refs/heads/feature"
	}`

	validCreateRepoResp = `{
		"id": "test-repo-id",
		"name": "test-repo",
		"defaultBranch": "refs/heads/main",
		"project": {
			"id": "test-project-id",
			"name": "testproject"
		}
	}`

	validCreateRepoRespWithFeatureBranch = `{
		"id": "test-repo-id",
		"name": "test-repo",
		"defaultBranch": "refs/heads/feature",
		"project": {
			"id": "test-project-id",
			"name": "testproject"
		}
	}`

	validUpdateRepoResp = `{
		"id": "test-repo-id",
		"name": "test-repo",
		"defaultBranch": "refs/heads/feature",
		"project": {
			"id": "test-project-id",
			"name": "testproject"
		}
	}`

	repoNotFoundResp = `{
		"message": "Repository not found"
	}`

	unauthorizedResp = `{
		"message": "Unauthorized"
	}`

	branchExistsResp = `{
		"value": [
			{
				"name": "refs/heads/main"
			}
		]
	}`

	branchDoesNotExistResp = `{
		"value": []
	}`
)

// Test constructor functions
func TestPostGitRepository(t *testing.T) {
	t.Run("returns valid handler", func(t *testing.T) {
		client := &http.Client{}
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		opts := handlers.HandlerOptions{
			Client: client,
			Log:    &logger,
		}

		handlerInterface := PostGitRepository(opts)

		if handlerInterface == nil {
			t.Fatal("PostGitRepository should return a non-nil handler")
		}

		// Verify it implements the Handler interface
		var _ handlers.Handler = handlerInterface

		// Verify the handler has the correct type and options
		h, ok := handlerInterface.(*postHandler)
		if !ok {
			t.Fatal("PostGitRepository should return a *postHandler")
		}

		if h.Client != client {
			t.Error("Handler should have the provided client")
		}

		if h.Log != &logger {
			t.Error("Handler should have the provided logger")
		}
	})
}

// Test POST handler
func TestPostHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name                 string
		organization         string
		project              string
		apiVersion           string
		sourceRef            string
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
			name:         "successful new repository creation (no init, no default branch)",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBody,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(repoCreateURL, http.StatusCreated, validCreateRepoResp)
			},
			expectedStatus:       http.StatusCreated,
			expectedContentType:  "application/json",
			expectedBodyContains: `"id":"test-repo-id"`,
			expectedRequestCount: 1,
			verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
				req := mockClient.getLastRequest()
				if req.URL.String() != repoCreateURL {
					t.Errorf("Request URL = %s, want %s", req.URL.String(), repoCreateURL)
				}
				if req.Method != "POST" {
					t.Errorf("Request Method = %s, want POST", req.Method)
				}
			},
		},
		{
			name:         "successful new repository creation with initialization",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBodyWithInit,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(repoCreateURL, http.StatusCreated, validCreateRepoResp)
				mockClient.setResponse(repoPushesURL, http.StatusCreated, `{}`)
			},
			expectedStatus:       http.StatusCreated,
			expectedContentType:  "application/json",
			expectedBodyContains: `"id":"test-repo-id"`,
			expectedRequestCount: 2, // Create + Initialize
			verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
				if mockClient.getRequestCount() != 2 {
					t.Errorf("Expected 2 requests, got %d", mockClient.getRequestCount())
				}
				// Verify create request
				createReq := mockClient.requests[0]
				if createReq.URL.String() != repoCreateURL {
					t.Errorf("Create Request URL = %s, want %s", createReq.URL.String(), repoCreateURL)
				}
				// Verify initialize request
				initReq := mockClient.requests[1]
				if initReq.URL.String() != repoPushesURL {
					t.Errorf("Initialize Request URL = %s, want %s", initReq.URL.String(), repoPushesURL)
				}
				if initReq.Method != "POST" {
					t.Errorf("Initialize Request Method = %s, want POST", initReq.Method)
				}
			},
		},
		{
			name:         "successful new repository creation with initialization and custom default branch",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBodyWithDefaultBranch,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(repoCreateURL, http.StatusCreated, validCreateRepoResp)
				mockClient.setResponse(repoPushesURL, http.StatusCreated, `{}`)
				mockClient.setResponse(repoUpdateURL, http.StatusOK, validUpdateRepoResp)
			},
			expectedStatus:       http.StatusCreated,
			expectedContentType:  "application/json",
			expectedBodyContains: `"defaultBranch":"refs/heads/feature"`,
			expectedRequestCount: 3, // Create + Initialize + Update Default Branch
			verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
				if mockClient.getRequestCount() != 3 {
					t.Errorf("Expected 3 requests, got %d", mockClient.getRequestCount())
				}
				// Verify create request
				createReq := mockClient.requests[0]
				if createReq.URL.String() != repoCreateURL {
					t.Errorf("Create Request URL = %s, want %s", createReq.URL.String(), repoCreateURL)
				}
				// Verify initialize request
				initReq := mockClient.requests[1]
				if initReq.URL.String() != repoPushesURL {
					t.Errorf("Initialize Request URL = %s, want %s", initReq.URL.String(), repoPushesURL)
				}
				// Verify update default branch request
				updateReq := mockClient.requests[2]
				if updateReq.URL.String() != repoUpdateURL {
					t.Errorf("Update Request URL = %s, want %s", updateReq.URL.String(), repoUpdateURL)
				}
				if updateReq.Method != "PATCH" {
					t.Errorf("Request Method = %s, want PATCH", updateReq.Method)
				}
			},
		},
		{
			name:         "new repository creation fails if defaultBranch specified but initialize is false",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBodyWithDefaultBranchNoInit,
			setupMock: func(mockClient *mockHTTPClient) {
				// No mock setup needed as validation happens before API call
			},
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "When specifying a 'defaultBranch' for a new repository, 'initialize' must be set to true",
			expectedRequestCount: 0,
		},
		{
			name:         "successful fork creation",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBodyFork,
			sourceRef:    "refs/heads/main", // Example sourceRef
			setupMock: func(mockClient *mockHTTPClient) {
				// Mock branch existence check for parent repo
				mockClient.setResponse(fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/parent-repo-id/refs?filter=heads/main&api-version=%s", testOrg, testProject, testAPIVersion), http.StatusOK, branchExistsResp)
				mockClient.setResponse(repoCreateURL+"&sourceRef=refs/heads/main", http.StatusCreated, validCreateRepoResp)
			},
			expectedStatus:       http.StatusCreated,
			expectedContentType:  "application/json",
			expectedBodyContains: `"id":"test-repo-id"`,
			expectedRequestCount: 2, // Branch check + Create
			verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
				if mockClient.getRequestCount() != 2 {
					t.Errorf("Expected 2 requests, got %d", mockClient.getRequestCount())
				}
				// Verify branch check request
				branchCheckReq := mockClient.requests[0]
				if !strings.Contains(branchCheckReq.URL.String(), "parent-repo-id/refs?filter=heads/main") {
					t.Errorf("Branch check URL mismatch: got %s", branchCheckReq.URL.String())
				}
				// Verify create request
				createReq := mockClient.requests[1]
				if !strings.Contains(createReq.URL.String(), "sourceRef=refs/heads/main") {
					t.Errorf("Create Request URL missing sourceRef: got %s", createReq.URL.String())
				}
			},
		},
		{
			name:         "successful fork creation with default branch (branch exists in fork)",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBodyForkWithDefaultBranch,
			sourceRef:    "refs/heads/feature",
			setupMock: func(mockClient *mockHTTPClient) {
				// Mock branch existence check for parent repo
				mockClient.setResponse(fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/parent-repo-id/refs?filter=heads/feature&api-version=%s", testOrg, testProject, testAPIVersion), http.StatusOK, branchExistsResp)
				// Mock create repo
				mockClient.setResponse(repoCreateURL+"&sourceRef=refs/heads/feature", http.StatusCreated, validCreateRepoResp)
				// Mock branch existence check for newly created fork
				mockClient.setResponse(fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s/refs?filter=heads/feature&api-version=%s", testOrg, testProject, testRepoID, testAPIVersion), http.StatusOK, branchExistsResp)
				// Mock update default branch
				mockClient.setResponse(repoUpdateURL, http.StatusOK, validUpdateRepoResp)
			},
			expectedStatus:       http.StatusCreated,
			expectedContentType:  "application/json",
			expectedBodyContains: `"defaultBranch":"refs/heads/feature"`,
			expectedRequestCount: 4, // Parent branch check + Create + Fork branch check + Update Default Branch
		},
		{
			name:         "fork creation with pending default branch (branch does not exist in fork)",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBodyForkWithDefaultBranch,
			sourceRef:    "refs/heads/new-feature", // Assume this branch doesn't exist in parent or fork initially
			setupMock: func(mockClient *mockHTTPClient) {
				// Mock branch existence check for parent repo (assume it exists for sourceRef validation)
				mockClient.setResponse(fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/parent-repo-id/refs?filter=heads/new-feature&api-version=%s", testOrg, testProject, testAPIVersion), http.StatusOK, branchExistsResp)
				// Mock create repo
				mockClient.setResponse(repoCreateURL+"&sourceRef=refs/heads/new-feature", http.StatusCreated, validCreateRepoResp)
				// Mock branch existence check for newly created fork (does NOT exist)
				mockClient.setResponse(fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s/refs?filter=heads/feature&api-version=%s", testOrg, testProject, testRepoID, testAPIVersion), http.StatusOK, branchDoesNotExistResp)
			},
			expectedStatus:       http.StatusAccepted, // 202 Accepted
			expectedContentType:  "application/json",
			expectedBodyContains: `"id":"test-repo-id"`, // Should return the created repo, but with original default branch
			expectedRequestCount: 3,                     // Parent branch check + Create + Fork branch check
		},
		{
			name:                 "missing organization parameter",
			organization:         "",
			project:              testProject,
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			requestBody:          validCreateRepoReqBody,
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "Organization parameter is required",
			expectedRequestCount: 0,
		},
		{
			name:                 "missing project parameter",
			organization:         testOrg,
			project:              "",
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			requestBody:          validCreateRepoReqBody,
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "Project ID parameter is required",
			expectedRequestCount: 0,
		},
		{
			name:                 "missing api version parameter",
			organization:         testOrg,
			project:              testProject,
			apiVersion:           "",
			authHeader:           testAuthHeader,
			requestBody:          validCreateRepoReqBody,
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "API version parameter is required",
			expectedRequestCount: 0,
		},
		{
			name:                 "missing authorization header",
			organization:         testOrg,
			project:              testProject,
			apiVersion:           testAPIVersion,
			authHeader:           "",
			requestBody:          validCreateRepoReqBody,
			setupMock:            nil,
			expectedStatus:       http.StatusUnauthorized,
			expectedContentType:  "",
			expectedBodyContains: "Request rejected due to missing or invalid Basic authentication",
			expectedRequestCount: 0,
		},
		{
			name:                 "invalid request body - bad json",
			organization:         testOrg,
			project:              testProject,
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			requestBody:          `{"name": "test"`, // Malformed JSON
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "Invalid JSON in request body",
			expectedRequestCount: 0,
		},
		{
			name:                 "missing repository name in request body",
			organization:         testOrg,
			project:              testProject,
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			requestBody:          `{"initialize": true}`, // Missing name
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "Repository name is required",
			expectedRequestCount: 0,
		},
		{
			name:         "network error during initial repository creation",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBody,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setError(repoCreateURL, fmt.Errorf("network error"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedContentType:  "",
			expectedBodyContains: "Failed to create repository",
			expectedRequestCount: 1,
		},
		{
			name:         "azure devops returns non-201 for create",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBody,
			setupMock: func(mockHTTPClient *mockHTTPClient) {
				mockHTTPClient.setResponse(repoCreateURL, http.StatusBadRequest, `{"message": "Bad Request"}`)
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedContentType:  "",
			expectedBodyContains: "Failed to create repository",
			expectedRequestCount: 1,
		},
		{
			name:         "network error during initialization",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBodyWithInit,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(repoCreateURL, http.StatusCreated, validCreateRepoResp)
				mockClient.setError(repoPushesURL, fmt.Errorf("init network error"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedContentType:  "",
			expectedBodyContains: "Failed to initialize repository",
			expectedRequestCount: 2,
		},
		{
			name:         "network error during default branch update",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBodyWithDefaultBranch,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(repoCreateURL, http.StatusCreated, validCreateRepoResp)
				mockClient.setResponse(repoPushesURL, http.StatusCreated, `{}`)
				mockClient.setError(repoUpdateURL, fmt.Errorf("update network error"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedContentType:  "",
			expectedBodyContains: "Failed to set default branch",
			expectedRequestCount: 3,
		},
		{
			name:                 "invalid sourceRef format",
			organization:         testOrg,
			project:              testProject,
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			requestBody:          validCreateRepoReqBodyFork,
			sourceRef:            "main", // Invalid format
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "sourceRef must start with 'refs/heads/'",
			expectedRequestCount: 0,
		},
		{
			name:         "sourceRef branch not found in parent repository",
			organization: testOrg,
			project:      testProject,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			requestBody:  validCreateRepoReqBodyFork,
			sourceRef:    "refs/heads/non-existent-branch",
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/parent-repo-id/refs?filter=heads/non-existent-branch&api-version=%s", testOrg, testProject, testAPIVersion), http.StatusOK, branchDoesNotExistResp)
			},
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: "SourceRef 'refs/heads/non-existent-branch' does not exist in parent repository 'parent-repo-id'",
			expectedRequestCount: 1, // Only branch existence check
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := newMockHTTPClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			handler := createTestPostHandler(mockClient)

			// Create request
			url := fmt.Sprintf("/api/%s/%s/git/repositories?api-version=%s", tt.organization, tt.project, tt.apiVersion)
			if tt.sourceRef != "" {
				url += fmt.Sprintf("&sourceRef=%s", tt.sourceRef)
			}
			req := httptest.NewRequest("POST", url, strings.NewReader(tt.requestBody))

			// Set path values directly on the request.
			req.SetPathValue("organization", tt.organization)
			req.SetPathValue("projectId", tt.project)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
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
