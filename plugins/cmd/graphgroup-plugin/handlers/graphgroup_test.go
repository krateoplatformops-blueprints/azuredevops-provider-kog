package graphgroup

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
type mockHTTPClient struct {
	responses map[string]*http.Response
	errors    map[string]error
	requests  []*http.Request
}

// newMockHTTPClient creates a new instance of mockHTTPClient
func newMockHTTPClient() *mockHTTPClient {
	return &mockHTTPClient{
		responses: make(map[string]*http.Response),
		errors:    make(map[string]error),
		requests:  make([]*http.Request, 0),
	}
}

// Do implements the http.Client Do method for mockHTTPClient
func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	m.requests = append(m.requests, req)

	key := req.URL.String()

	if err, exists := m.errors[key]; exists {
		return nil, err
	}

	if resp, exists := m.responses[key]; exists {
		return resp, nil
	}

	fmt.Printf("MockHTTPClient: No response configured for URL: %s. Returning 404.\n", key)
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(strings.NewReader(`{"message": "Not Found"}`)),
		Header:     make(http.Header),
	}, nil
}

// setResponse allows setting a predefined response for a specific URL
func (m *mockHTTPClient) setResponse(url string, statusCode int, body string, headers map[string]string) {
	header := make(http.Header)
	for k, v := range headers {
		header.Set(k, v)
	}
	m.responses[url] = &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     header,
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

// Test data constants
const (
	testOrg        = "testorg"
	testAPIVersion = "7.2-preview.1"
	testAuthHeader = "Basic dGVzdDp0ZXN0"
	testUsername   = "test"
	testPassword   = "test"
)

var (
	baseGraphGroupsURL = fmt.Sprintf("https://vssps.dev.azure.com/%s/_apis/graph/groups?api-version=%s", testOrg, testAPIVersion)

	// Response with mixed organization and project-level groups
	mixedGroupsResp = `{
		"count": 5,
		"value": [
			{
				"descriptor": "org-group-1",
				"displayName": "Project Collection Admins",
				"domain": "vstfs:///Framework/IdentityDomain/00000000-0000-0000-0000-000000000001",
				"origin": "vsts",
				"originId": "org-id-1",
				"principalName": "[Org]\\Admins",
				"subjectKind": "group"
			},
			{
				"descriptor": "project-group-1",
				"displayName": "Build Administrators",
				"domain": "vstfs:///Classification/TeamProject/project-uuid-1",
				"origin": "vsts",
				"originId": "proj-id-1",
				"principalName": "[Project1]\\Build Admins",
				"subjectKind": "group"
			},
			{
				"descriptor": "org-group-2",
				"displayName": "Project Collection Service Accounts",
				"domain": "vstfs:///Framework/IdentityDomain/00000000-0000-0000-0000-000000000002",
				"origin": "vsts",
				"originId": "org-id-2",
				"principalName": "[Org]\\Service",
				"subjectKind": "group"
			},
			{
				"descriptor": "project-group-2",
				"displayName": "Contributors",
				"domain": "vstfs:///Classification/TeamProject/project-uuid-2",
				"origin": "vsts",
				"originId": "proj-id-2",
				"principalName": "[Project2]\\Contributors",
				"subjectKind": "group"
			},
			{
				"descriptor": "project-group-3",
				"displayName": "Readers",
				"domain": "vstfs:///Classification/TeamProject/project-uuid-3",
				"origin": "vsts",
				"originId": "proj-id-3",
				"principalName": "[Project3]\\Readers",
				"subjectKind": "group"
			}
		]
	}`

	// Response with only organization-level groups
	orgOnlyGroupsResp = `{
		"count": 2,
		"value": [
			{
				"descriptor": "org-group-1",
				"displayName": "Project Collection Admins",
				"domain": "vstfs:///Framework/IdentityDomain/00000000-0000-0000-0000-000000000001",
				"origin": "vsts",
				"originId": "org-id-1",
				"principalName": "[Org]\\Admins",
				"subjectKind": "group"
			},
			{
				"descriptor": "org-group-2",
				"displayName": "Project Collection Service Accounts",
				"domain": "vstfs:///Framework/IdentityDomain/00000000-0000-0000-0000-000000000002",
				"origin": "vsts",
				"originId": "org-id-2",
				"principalName": "[Org]\\Service",
				"subjectKind": "group"
			}
		]
	}`

	// Response with only project-level groups
	projectOnlyGroupsResp = `{
		"count": 2,
		"value": [
			{
				"descriptor": "project-group-1",
				"displayName": "Build Administrators",
				"domain": "vstfs:///Classification/TeamProject/project-uuid-1",
				"origin": "vsts",
				"originId": "proj-id-1",
				"principalName": "[Project1]\\Build Admins",
				"subjectKind": "group"
			},
			{
				"descriptor": "project-group-2",
				"displayName": "Contributors",
				"domain": "vstfs:///Classification/TeamProject/project-uuid-2",
				"origin": "vsts",
				"originId": "proj-id-2",
				"principalName": "[Project2]\\Contributors",
				"subjectKind": "group"
			}
		]
	}`

	emptyGroupsResp = `{
		"count": 0,
		"value": []
	}`

	unauthorizedResp = `{
		"message": "Unauthorized"
	}`

	badRequestResp = `{
		"message": "Bad Request"
	}`
)

// Test constructor function
func TestGetGraphGroups(t *testing.T) {
	t.Run("returns valid handler", func(t *testing.T) {
		client := &http.Client{}
		logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
		opts := handlers.HandlerOptions{
			Client: client,
			Log:    &logger,
		}

		handlerInterface := GetGraphGroups(opts)

		if handlerInterface == nil {
			t.Fatal("GetGraphGroups should return a non-nil handler")
		}

		var _ handlers.Handler = handlerInterface

		h, ok := handlerInterface.(*getHandler)
		if !ok {
			t.Fatal("GetGraphGroups should return a *getHandler")
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
		scopeDescriptor      string
		subjectTypes         string
		continuationToken    string
		apiVersion           string
		authHeader           string
		setupMock            func(*mockHTTPClient)
		expectedStatus       int
		expectedContentType  string
		expectedBodyContains []string
		expectedRequestCount int
		verifyResponse       func(t *testing.T, body string)
	}{
		{
			name:              "filters project-level groups when no scopeDescriptor",
			organization:      testOrg,
			apiVersion:        testAPIVersion,
			authHeader:        testAuthHeader,
			scopeDescriptor:   "", // No scope descriptor - should filter
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(baseGraphGroupsURL, http.StatusOK, mixedGroupsResp, nil)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: []string{`"originalCount":5`, `"count":2`, `"numberOfFilteredOut":3`, `"Project Collection Admins"`, `"Project Collection Service Accounts"`},
			expectedRequestCount: 1,
			verifyResponse: func(t *testing.T, body string) {
				// Should NOT contain project-level groups
				if strings.Contains(body, "Build Administrators") {
					t.Error("Response should not contain project-level groups")
				}
				if strings.Contains(body, "Contributors") {
					t.Error("Response should not contain project-level groups")
				}
				if strings.Contains(body, "vstfs:///Classification/TeamProject/") {
					t.Error("Response should not contain project domain prefix")
				}
			},
		},
		{
			name:              "no filtering when scopeDescriptor is provided",
			organization:      testOrg,
			apiVersion:        testAPIVersion,
			authHeader:        testAuthHeader,
			scopeDescriptor:   "scp.PROJECT_SCOPE_DESCRIPTOR",
			setupMock: func(mockClient *mockHTTPClient) {
				url := baseGraphGroupsURL + "&scopeDescriptor=scp.PROJECT_SCOPE_DESCRIPTOR"
				mockClient.setResponse(url, http.StatusOK, projectOnlyGroupsResp, nil)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: []string{`"originalCount":2`, `"count":2`, `"numberOfFilteredOut":0`, `"Build Administrators"`, `"Contributors"`},
			expectedRequestCount: 1,
			verifyResponse: func(t *testing.T, body string) {
				// Should contain all project-level groups (no filtering)
				if !strings.Contains(body, "Build Administrators") {
					t.Error("Response should contain all groups when scopeDescriptor is set")
				}
			},
		},
		{
			name:              "returns all org-level groups when only org groups exist",
			organization:      testOrg,
			apiVersion:        testAPIVersion,
			authHeader:        testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(baseGraphGroupsURL, http.StatusOK, orgOnlyGroupsResp, nil)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: []string{`"originalCount":2`, `"count":2`, `"numberOfFilteredOut":0`},
			expectedRequestCount: 1,
		},
		{
			name:              "filters all groups when only project groups exist",
			organization:      testOrg,
			apiVersion:        testAPIVersion,
			authHeader:        testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(baseGraphGroupsURL, http.StatusOK, projectOnlyGroupsResp, nil)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: []string{`"originalCount":2`, `"count":0`, `"numberOfFilteredOut":2`, `"value":[]`},
			expectedRequestCount: 1,
			verifyResponse: func(t *testing.T, body string) {
				if strings.Contains(body, "Build Administrators") {
					t.Error("All project-level groups should be filtered out")
				}
			},
		},
		{
			name:              "handles empty response from Azure DevOps",
			organization:      testOrg,
			apiVersion:        testAPIVersion,
			authHeader:        testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(baseGraphGroupsURL, http.StatusOK, emptyGroupsResp, nil)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: []string{`"originalCount":0`, `"count":0`, `"numberOfFilteredOut":0`, `"value":[]`},
			expectedRequestCount: 1,
		},
		{
			name:              "preserves continuation token in response header",
			organization:      testOrg,
			apiVersion:        testAPIVersion,
			authHeader:        testAuthHeader,
			continuationToken: "",
			setupMock: func(mockClient *mockHTTPClient) {
				headers := map[string]string{
					"X-MS-ContinuationToken": "next-page-token-123",
				}
				mockClient.setResponse(baseGraphGroupsURL, http.StatusOK, mixedGroupsResp, headers)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: []string{`"originalCount":5`, `"count":2`},
			expectedRequestCount: 1,
			verifyResponse: func(t *testing.T, body string) {
				// Continuation token is checked in the response headers, not body
			},
		},
		{
			name:              "forwards continuation token to Azure DevOps",
			organization:      testOrg,
			apiVersion:        testAPIVersion,
			authHeader:        testAuthHeader,
			continuationToken: "page-2-token",
			setupMock: func(mockClient *mockHTTPClient) {
				url := baseGraphGroupsURL + "&continuationToken=page-2-token"
				mockClient.setResponse(url, http.StatusOK, orgOnlyGroupsResp, nil)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: []string{`"originalCount":2`, `"count":2`},
			expectedRequestCount: 1,
		},
		{
			name:              "forwards subjectTypes parameter",
			organization:      testOrg,
			apiVersion:        testAPIVersion,
			authHeader:        testAuthHeader,
			subjectTypes:      "Microsoft.IdentityModel.Claims.ClaimsIdentity",
			setupMock: func(mockClient *mockHTTPClient) {
				url := baseGraphGroupsURL + "&subjectTypes=Microsoft.IdentityModel.Claims.ClaimsIdentity"
				mockClient.setResponse(url, http.StatusOK, orgOnlyGroupsResp, nil)
			},
			expectedStatus:       http.StatusOK,
			expectedContentType:  "application/json",
			expectedBodyContains: []string{`"originalCount":2`, `"count":2`},
			expectedRequestCount: 1,
		},
		{
			name:                 "returns 400 when organization is missing",
			organization:         "",
			apiVersion:           testAPIVersion,
			authHeader:           testAuthHeader,
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: []string{"organization parameter is required"},
			expectedRequestCount: 0,
		},
		{
			name:                 "returns 400 when api-version is missing",
			organization:         testOrg,
			apiVersion:           "",
			authHeader:           testAuthHeader,
			setupMock:            nil,
			expectedStatus:       http.StatusBadRequest,
			expectedContentType:  "",
			expectedBodyContains: []string{"api-version is required"},
			expectedRequestCount: 0,
		},
		{
			name:                 "returns 401 when authorization header is missing",
			organization:         testOrg,
			apiVersion:           testAPIVersion,
			authHeader:           "",
			setupMock:            nil,
			expectedStatus:       http.StatusUnauthorized,
			expectedContentType:  "",
			expectedBodyContains: []string{"request rejected due to missing or invalid Basic authentication"},
			expectedRequestCount: 0,
		},
		{
			name:         "forwards Azure DevOps API errors",
			organization: testOrg,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setResponse(baseGraphGroupsURL, http.StatusUnauthorized, unauthorizedResp, nil)
			},
			expectedStatus:       http.StatusUnauthorized,
			expectedContentType:  "",
			expectedBodyContains: []string{"Azure DevOps API error"},
			expectedRequestCount: 1,
		},
		{
			name:         "handles network errors",
			organization: testOrg,
			apiVersion:   testAPIVersion,
			authHeader:   testAuthHeader,
			setupMock: func(mockClient *mockHTTPClient) {
				mockClient.setError(baseGraphGroupsURL, fmt.Errorf("network timeout"))
			},
			expectedStatus:       http.StatusInternalServerError,
			expectedContentType:  "",
			expectedBodyContains: []string{"error calling Azure DevOps API"},
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

			// Build URL with query parameters
			url := fmt.Sprintf("/plugin/%s/graph/groups?api-version=%s", tt.organization, tt.apiVersion)
			if tt.scopeDescriptor != "" {
				url += fmt.Sprintf("&scopeDescriptor=%s", tt.scopeDescriptor)
			}
			if tt.subjectTypes != "" {
				url += fmt.Sprintf("&subjectTypes=%s", tt.subjectTypes)
			}
			if tt.continuationToken != "" {
				url += fmt.Sprintf("&continuationToken=%s", tt.continuationToken)
			}

			req := httptest.NewRequest("GET", url, nil)

			// Set path values
			req.SetPathValue("organization", tt.organization)

			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
				req.SetBasicAuth(testUsername, testPassword)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute handler
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
			body := rr.Body.String()
			for _, expected := range tt.expectedBodyContains {
				if !strings.Contains(body, expected) {
					t.Errorf("handler response body does not contain expected content.\nGot: %s\nWant to contain: %s", body, expected)
				}
			}

			// Verify request count
			if mockClient.getRequestCount() != tt.expectedRequestCount {
				t.Errorf("expected %d requests, got %d", tt.expectedRequestCount, mockClient.getRequestCount())
			}

			// Run custom verification if provided
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, body)
			}

			// Verify continuation token in response headers if applicable
			if tt.name == "preserves continuation token in response header" {
				token := rr.Header().Get("X-MS-ContinuationToken")
				if token != "next-page-token-123" {
					t.Errorf("Expected continuation token in response header, got: %s", token)
				}
			}
		})
	}
}
