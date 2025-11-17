package pullrequest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockHTTPClient implements http.Client's Do method for testing.
type mockHTTPClient struct {
	responses map[string]*http.Response
	errors    map[string]error
	requests  []*http.Request
}

func newMockHTTPClient() *mockHTTPClient {
	return &mockHTTPClient{
		responses: make(map[string]*http.Response),
		errors:    make(map[string]error),
		requests:  make([]*http.Request, 0),
	}
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.requests = append(m.requests, req)
	key := req.URL.String()
	if err, exists := m.errors[key]; exists {
		return nil, err
	}
	if resp, exists := m.responses[key]; exists {
		return resp, nil
	}
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(`{"message": "Not Found"}`)),
		Header:     make(http.Header),
	}, nil
}

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

func createTestPatchHandler(mockClient *mockHTTPClient) *patchHandler {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	return &patchHandler{
		baseHandler: &baseHandler{
			HandlerOptions: handlers.HandlerOptions{
				Client: mockClient,
				Log:    &logger,
			},
		},
	}
}

const (
	testOrg        = "test-org"
	testProject    = "test-project"
	testRepoID     = "test-repo-id"
	testPullReqID  = "101"
	testAPIVersion = "7.2-preview.2"
	testAuthHeader = "Basic dXNlcjpwYXNzd29yZA=="
	testUsername   = "user"
	testPassword   = "password"
)

var (
	pullRequestCreateURL = fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s/pullrequests?api-version=%s", testOrg, testProject, testRepoID, testAPIVersion)
	pullRequestUpdateURL = fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s/pullrequests/%s?api-version=%s", testOrg, testProject, testRepoID, testPullReqID, testAPIVersion)

	// Superset request body, as sent to our plugin
	validPluginRequestBody = `{
		"title": "My Test PR",
		"description": "This is a test.",
		"sourceRefName": "refs/heads/feature/test",
		"targetRefName": "refs/heads/main",
		"status": "active",
		"isDraft": false,
		"reviewers": [{"id": "reviewer-id"}],
		"completionOptions": {
			"deleteSourceBranch": true,
			"squashMerge": true
		}
	}`

	// Request body for retargeting a PR
	retargetRequestBody = `{
		"targetRefName": "refs/heads/develop"
	}`

	// Response body from Azure DevOps API
	validAzdoResponse = `{
		"pullRequestId": 101,
		"title": "My Test PR",
		"description": "This is a test.",
		"sourceRefName": "refs/heads/feature/test",
		"targetRefName": "refs/heads/main",
		"status": "active",
		"isDraft": false
	}`

	notFoundResponse = `{"message": "Pull request not found"}`
	unauthorizedResp = `{"message": "Unauthorized"}`
	conflictResponse = `{"message": "PR already exists"}`
	genericError     = `{"message": "An error occurred"}`
)

// Test constructor functions
//
//	func TestPostPullRequest(t *testing.T) {
//		t.Run("returns valid handler", func(t *testing.T) {
//			client := &http.Client{}
//			logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
//			opts := handlers.HandlerOptions{
//				Client: client,
//				Log:    &logger,
//			}
//
//			handlerInterface := PostPullRequest(opts)
//
//			if handlerInterface == nil {
//				t.Fatal("PostPullRequest should return a non-nil handler")
//			}
//
//			var _ handlers.Handler = handlerInterface
//
//			h, ok := handlerInterface.(*postHandler)
//			if !ok {
//				t.Fatal("PostPullRequest should return a *postHandler")
//			}
//
//			if h.Client != client {
//				t.Error("Handler should have the provided client")
//			}
//
//			if h.Log != &logger {
//				t.Error("Handler should have the provided logger")
//			}
//		})
//	}
//
//	func TestPatchPullRequest(t *testing.T) {
//		t.Run("returns valid handler", func(t *testing.T) {
//			client := &http.Client{}
//			logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
//			opts := handlers.HandlerOptions{
//				Client: client,
//				Log:    &logger,
//			}
//
//			handlerInterface := PatchPullRequest(opts)
//
//			if handlerInterface == nil {
//				t.Fatal("PatchPullRequest should return a non-nil handler")
//			}
//
//			var _ handlers.Handler = handlerInterface
//
//			h, ok := handlerInterface.(*patchHandler)
//			if !ok {
//				t.Fatal("PatchPullRequest should return a *patchHandler")
//			}
//
//			if h.Client != client {
//				t.Error("Handler should have the provided client")
//			}
//
//			if h.Log != &logger {
//				t.Error("Handler should have the provided logger")
//			}
//		})
//	}
//
//	func TestPostHandler_ServeHTTP(t *testing.T) {
//		tests := []struct {
//			name                 string
//			organization         string
//			project              string
//			repositoryId         string
//			apiVersion           string
//			authHeader           string
//			requestBody          string
//			setupMock            func(*mockHTTPClient)
//			expectedStatus       int
//			expectedBody         string
//			isJSONResponse       bool
//			expectedRequestCount int
//			verifyRequests       func(t *testing.T, mockClient *mockHTTPClient)
//		}{
//			{
//				name:         "successful pull request creation",
//				organization: testOrg,
//				project:      testProject,
//				repositoryId: testRepoID,
//				apiVersion:   testAPIVersion,
//				authHeader:   testAuthHeader,
//				requestBody:  validPluginRequestBody,
//				setupMock: func(mc *mockHTTPClient) {
//					mc.setResponse(pullRequestCreateURL, http.StatusCreated, validAzdoResponse)
//				},
//				expectedStatus:       http.StatusCreated,
//				expectedBody:         validAzdoResponse,
//				isJSONResponse:       true,
//				expectedRequestCount: 1,
//				verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
//					req := mockClient.getLastRequest()
//					assert.Equal(t, pullRequestCreateURL, req.URL.String())
//					assert.Equal(t, testAuthHeader, req.Header.Get("Authorization"))
//					assert.Equal(t, "POST", req.Method)
//					body, _ := io.ReadAll(req.Body)
//					assert.NotContains(t, string(body), "status")
//				},
//			},
//			{
//				name:                 "missing organization",
//				organization:         "",
//				project:              testProject,
//				repositoryId:         testRepoID,
//				apiVersion:           testAPIVersion,
//				authHeader:           testAuthHeader,
//				requestBody:          validPluginRequestBody,
//				setupMock:            nil,
//				expectedStatus:       http.StatusBadRequest,
//				expectedBody:         "organization parameter is required",
//				expectedRequestCount: 0,
//			},
//			{
//				name:                 "missing project",
//				organization:         testOrg,
//				project:              "",
//				repositoryId:         testRepoID,
//				apiVersion:           testAPIVersion,
//				authHeader:           testAuthHeader,
//				requestBody:          validPluginRequestBody,
//				setupMock:            nil,
//				expectedStatus:       http.StatusBadRequest,
//				expectedBody:         "project parameter is required",
//				expectedRequestCount: 0,
//			},
//			{
//				name:                 "missing repositoryId",
//				organization:         testOrg,
//				project:              testProject,
//				repositoryId:         "",
//				apiVersion:           testAPIVersion,
//				authHeader:           testAuthHeader,
//				requestBody:          validPluginRequestBody,
//				setupMock:            nil,
//				expectedStatus:       http.StatusBadRequest,
//				expectedBody:         "repositoryId parameter is required",
//				expectedRequestCount: 0,
//			},
//			{
//				name:                 "missing api-version",
//				organization:         testOrg,
//				project:              testProject,
//				repositoryId:         testRepoID,
//				apiVersion:           "",
//				authHeader:           testAuthHeader,
//				requestBody:          validPluginRequestBody,
//				setupMock:            nil,
//				expectedStatus:       http.StatusBadRequest,
//				expectedBody:         "api-version is required",
//				expectedRequestCount: 0,
//			},
//			{
//				name:                 "missing authorization header",
//				organization:         testOrg,
//				project:              testProject,
//				repositoryId:         testRepoID,
//				apiVersion:           testAPIVersion,
//				authHeader:           "",
//				requestBody:          validPluginRequestBody,
//				setupMock:            nil,
//				expectedStatus:       http.StatusUnauthorized,
//				expectedBody:         "request rejected due to missing or invalid Basic authentication",
//				expectedRequestCount: 0,
//			},
//			{
//				name:                 "bad request - invalid json",
//				organization:         testOrg,
//				project:              testProject,
//				repositoryId:         testRepoID,
//				apiVersion:           testAPIVersion,
//				authHeader:           testAuthHeader,
//				requestBody:          `{"title": "test"`,
//				setupMock:            nil,
//				expectedStatus:       http.StatusBadRequest,
//				expectedBody:         "error unmarshaling request body",
//				expectedRequestCount: 0,
//			},
//			{
//				name:         "azure devops returns conflict",
//				organization: testOrg,
//				project:      testProject,
//				repositoryId: testRepoID,
//				apiVersion:   testAPIVersion,
//				authHeader:   testAuthHeader,
//				requestBody:  validPluginRequestBody,
//				setupMock: func(mc *mockHTTPClient) {
//					mc.setResponse(pullRequestCreateURL, http.StatusConflict, conflictResponse)
//				},
//				expectedStatus:       http.StatusConflict,
//				expectedBody:         conflictResponse,
//				isJSONResponse:       true,
//				expectedRequestCount: 1,
//			},
//			{
//				name:         "azure devops returns unauthorized",
//				organization: testOrg,
//				project:      testProject,
//				repositoryId: testRepoID,
//				apiVersion:   testAPIVersion,
//				authHeader:   testAuthHeader,
//				requestBody:  validPluginRequestBody,
//				setupMock: func(mc *mockHTTPClient) {
//					mc.setResponse(pullRequestCreateURL, http.StatusUnauthorized, unauthorizedResp)
//				},
//				expectedStatus:       http.StatusUnauthorized,
//				expectedBody:         unauthorizedResp,
//				isJSONResponse:       true,
//				expectedRequestCount: 1,
//			},
//			{
//				name:         "network error on create",
//				organization: testOrg,
//				project:      testProject,
//				repositoryId: testRepoID,
//				apiVersion:   testAPIVersion,
//				authHeader:   testAuthHeader,
//				requestBody:  validPluginRequestBody,
//				setupMock: func(mc *mockHTTPClient) {
//					mc.setError(pullRequestCreateURL, fmt.Errorf("network error"))
//				},
//				expectedStatus:       http.StatusInternalServerError,
//				expectedBody:         "error making request to azure devops",
//				expectedRequestCount: 1,
//			},
//		}
//
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				mockClient := newMockHTTPClient()
//				if tt.setupMock != nil {
//					tt.setupMock(mockClient)
//				}
//
//				handler := createTestPostHandler(mockClient)
//				req := httptest.NewRequest("POST", "/", strings.NewReader(tt.requestBody))
//				req.SetPathValue("organization", tt.organization)
//				req.SetPathValue("project", tt.project)
//				req.SetPathValue("repositoryId", tt.repositoryId)
//				req.URL.RawQuery = "api-version=" + tt.apiVersion
//				if tt.authHeader != "" {
//					req.Header.Set("Authorization", tt.authHeader)
//					req.SetBasicAuth(testUsername, testPassword)
//				}
//
//				rr := httptest.NewRecorder()
//				handler.ServeHTTP(rr, req)
//
//				assert.Equal(t, tt.expectedStatus, rr.Code, "handler returned wrong status code")
//
//				if tt.expectedBody != "" {
//					if tt.isJSONResponse {
//						assert.JSONEq(t, tt.expectedBody, rr.Body.String(), "handler returned wrong JSON body")
//					} else {
//						assert.Contains(t, rr.Body.String(), tt.expectedBody, "handler returned wrong body content")
//					}
//				}
//
//				assert.Equal(t, tt.expectedRequestCount, mockClient.getRequestCount(), "unexpected number of requests")
//
//				if tt.verifyRequests != nil {
//					tt.verifyRequests(t, mockClient)
//				}
//			})
//		}
//	}

func TestPatchHandler_ServeHTTP(t *testing.T) {
	// Current state of the PR in Azure DevOps
	currentPRState := `
	{
		"pullRequestId": 101,
		"title": "Original Title",
		"description": "Original Description",
		"status": "active"
	}`

	tests := []struct {
		name                 string
		requestBody          string
		setupMock            func(*mockHTTPClient)
		expectedStatus       int
		expectedRequestCount int
		verifyRequests       func(t *testing.T, mockClient *mockHTTPClient)
	}{
		{
			name: "success - changes detected",
			requestBody: `{
				"title": "Updated Title",
				"description": "Updated Description"
			}`,
			setupMock: func(mc *mockHTTPClient) {
				// This single response will be used for both the internal GET and the PATCH call.
				mc.setResponse(pullRequestUpdateURL, http.StatusOK, currentPRState)
			},
			expectedStatus:       http.StatusOK,
			expectedRequestCount: 2, // 1 GET, 1 PATCH
			verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
				require.Len(t, mockClient.requests, 2, "Expected two requests to be made")

				// Verify GET request
				getReq := mockClient.requests[0]
				assert.Equal(t, "GET", getReq.Method)
				assert.Equal(t, pullRequestUpdateURL, getReq.URL.String())

				// Verify PATCH request
				patchReq := mockClient.requests[1]
				assert.Equal(t, "PATCH", patchReq.Method)
				body, _ := io.ReadAll(patchReq.Body)

				var payload map[string]interface{}
				err := json.Unmarshal(body, &payload)
				require.NoError(t, err)

				assert.Equal(t, "Updated Title", payload["title"])
				assert.Equal(t, "Updated Description", payload["description"])
				assert.Len(t, payload, 2, "Payload should only contain changed fields")
			},
		},
		{
			name: "success - no changes detected",
			requestBody: `{
				"title": "Original Title"
			}`,
			setupMock: func(mc *mockHTTPClient) {
				mc.setResponse(pullRequestUpdateURL, http.StatusOK, currentPRState)
			},
			expectedStatus:       http.StatusOK,
			expectedRequestCount: 1, // Only the GET call
			verifyRequests: func(t *testing.T, mockClient *mockHTTPClient) {
				require.Len(t, mockClient.requests, 1, "Expected only one request to be made")
				getReq := mockClient.requests[0]
				assert.Equal(t, "GET", getReq.Method)
			},
		},
		{
			name:        "error - initial GET fails",
			requestBody: `{"title": "Updated Title"}`,
			setupMock: func(mc *mockHTTPClient) {
				mc.setError(pullRequestUpdateURL, fmt.Errorf("network error"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedRequestCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := newMockHTTPClient()
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}

			handler := createTestPatchHandler(mockClient)
			req := httptest.NewRequest("PATCH", "/", strings.NewReader(tt.requestBody))
			req.SetPathValue("organization", testOrg)
			req.SetPathValue("project", testProject)
			req.SetPathValue("repositoryId", testRepoID)
			req.SetPathValue("pullRequestId", testPullReqID)
			req.URL.RawQuery = "api-version=" + testAPIVersion
			req.Header.Set("Authorization", testAuthHeader)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedRequestCount, mockClient.getRequestCount())

			if tt.verifyRequests != nil {
				tt.verifyRequests(t, mockClient)
			}
		})
	}
}
