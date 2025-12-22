package graphgroup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
)

const (
	requestTimeout             = 30 * time.Second
	contentTypeJSON            = "application/json"
	projectDomainPrefix        = "vstfs:///Classification/TeamProject/"
	continuationTokenHeaderKey = "X-MS-ContinuationToken" // Response header for continuation token
)

func GetGraphGroups(opts handlers.HandlerOptions) handlers.Handler {
	return &getHandler{baseHandler: newBaseHandler(opts)}
}

// Compile-time interface check
var _ handlers.Handler = (*getHandler)(nil)

// baseHandler contains common functionality shared by all handlers
type baseHandler struct {
	handlers.HandlerOptions
}

// newBaseHandler creates a new base handler with the given options
func newBaseHandler(opts handlers.HandlerOptions) *baseHandler {
	return &baseHandler{HandlerOptions: opts}
}

type getHandler struct {
	*baseHandler
}

// requestParams holds common request parameters extracted from HTTP requests
type requestParams struct {
	Organization      string
	ScopeDescriptor   string
	SubjectTypes      string
	ContinuationToken string
	APIVersion        string
	AuthHeader        string
}

// extractRequestParams extracts and validates parameters from GET requests
func (h *getHandler) extractRequestParams(r *http.Request) (*requestParams, error) {
	params := &requestParams{
		Organization:      r.PathValue("organization"),
		ScopeDescriptor:   r.URL.Query().Get("scopeDescriptor"),
		SubjectTypes:      r.URL.Query().Get("subjectTypes"),
		ContinuationToken: r.URL.Query().Get("continuationToken"),
		APIVersion:        r.URL.Query().Get("api-version"),
		AuthHeader:        r.Header.Get("Authorization"),
	}

	if err := params.validate(); err != nil {
		return nil, err
	}

	return params, nil
}

// validate checks if all required parameters are present and valid
func (p *requestParams) validate() error {
	if p.Organization == "" {
		return fmt.Errorf("organization parameter is required")
	}
	if p.APIVersion == "" {
		return fmt.Errorf("api-version is required")
	}
	if p.AuthHeader == "" {
		return fmt.Errorf("request rejected due to missing or invalid Basic authentication")
	}

	return nil
}

// makeAzureDevOpsRequest performs an HTTP request to Azure DevOps API
func (h *baseHandler) makeAzureDevOpsRequest(ctx context.Context, method, url, authHeader string, body []byte) (*http.Response, error) {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	if bodyReader != nil {
		req.Header.Set("Content-Type", contentTypeJSON)
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// writeErrorResponse writes an error response to the client
func (h *baseHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	h.Log.Print(message)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

// writeJSONResponse writes a JSON response to the client
func (h *baseHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, body []byte, continuationToken string) {
	w.Header().Set("Content-Type", contentTypeJSON)
	if continuationToken != "" {
		w.Header().Set(continuationTokenHeaderKey, continuationToken)
	}
	w.WriteHeader(statusCode)
	w.Write(body)
}

// buildAzureDevOpsURL constructs the Azure DevOps API URL with query parameters
func (h *getHandler) buildAzureDevOpsURL(params *requestParams) (string, error) {
	baseURL := fmt.Sprintf("https://vssps.dev.azure.com/%s/_apis/graph/groups", params.Organization)

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("api-version", params.APIVersion)

	// Add optional parameters if present
	if params.ScopeDescriptor != "" {
		q.Set("scopeDescriptor", params.ScopeDescriptor)
	}
	if params.SubjectTypes != "" {
		q.Set("subjectTypes", params.SubjectTypes)
	}
	if params.ContinuationToken != "" {
		q.Set("continuationToken", params.ContinuationToken)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// filterOrganizationLevelGroups filters out project-level groups
// Project-level groups have domain starting with "vstfs:///Classification/TeamProject/"
func filterOrganizationLevelGroups(groups []GraphGroup) []GraphGroup {
	filtered := make([]GraphGroup, 0, len(groups))
	for _, group := range groups {
		if !strings.HasPrefix(group.Domain, projectDomainPrefix) {
			filtered = append(filtered, group)
		}
	}
	return filtered
}

// processGraphGroupsResponse processes the Azure DevOps response and applies filtering if needed
func (h *getHandler) processGraphGroupsResponse(body []byte, shouldFilter bool) ([]byte, error) {
	var azureResp AzureDevOpsGraphGroupsResponse
	if err := json.Unmarshal(body, &azureResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	originalCount := azureResp.Count
	filteredGroups := azureResp.Value

	if shouldFilter {
		h.Log.Printf("Filtering organization-level groups (removing project-level groups)")
		filteredGroups = filterOrganizationLevelGroups(azureResp.Value)
		h.Log.Printf("Filtered out %d groups from %d total groups", originalCount-len(filteredGroups), originalCount)
		h.Log.Printf("Returning %d organization-level groups", len(filteredGroups))
	} else {
		h.Log.Printf("No filtering applied (scopeDescriptor is set or filtering not needed)")
	}

	newCount := len(filteredGroups)
	numberOfFilteredOut := originalCount - newCount

	enhancedResp := GraphGroupsResponse{
		OriginalCount:       originalCount,
		Count:               newCount,
		NumberOfFilteredOut: numberOfFilteredOut,
		Value:               filteredGroups,
	}

	return json.Marshal(enhancedResp)
}

// ServeHTTP handles GET requests for graph groups with optional filtering
// @Summary Get graph groups with optional filtering
// @Description Retrieve graph groups from Azure DevOps. When scopeDescriptor is not provided, filters out project-level groups to return only organization-level groups.
// @ID get-graphgroups
// @Tags Graph Groups
// @Param organization path string true "Organization name"
// @Param scopeDescriptor query string false "Scope descriptor (e.g., project descriptor). When provided, returns all groups within that scope without filtering."
// @Param subjectTypes query string false "Comma-separated list of subject types to filter by"
// @Param continuationToken query string false "Continuation token for pagination"
// @Param api-version query string true "API version (e.g., 7.2-preview.1)"
// @Param Authorization header string true "Basic Auth header"
// @Produce json
// @Success 200 {object} GraphGroupsResponse "Retrieved graph groups with filtering metadata"
// @Header 200 {string} X-MS-ContinuationToken "Continuation token for pagination (if more results available)"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /plugin/{organization}/graph/groups [get]
func (h *getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	// Extract and validate request parameters
	params, err := h.extractRequestParams(r)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "request rejected due to missing or invalid Basic authentication" {
			statusCode = http.StatusUnauthorized
		}
		h.writeErrorResponse(w, statusCode, err.Error())
		return
	}

	h.Log.Printf("Fetching graph groups for organization: %s", params.Organization)
	if params.ScopeDescriptor != "" {
		h.Log.Printf("Scope descriptor provided: %s", params.ScopeDescriptor)
	} else {
		h.Log.Printf("No scope descriptor provided, will filter to organization-level groups only")
	}

	// Build Azure DevOps API URL
	azureURL, err := h.buildAzureDevOpsURL(params)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error building Azure DevOps URL: %s", err.Error()))
		return
	}

	h.Log.Printf("Calling Azure DevOps API: %s", azureURL)

	// Call Azure DevOps API
	resp, err := h.makeAzureDevOpsRequest(ctx, http.MethodGet, azureURL, params.AuthHeader, nil)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error calling Azure DevOps API: %s", err.Error()))
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error reading response body: %s", err.Error()))
		return
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		h.Log.Printf("Azure DevOps API returned non-200 status: %d - %s", resp.StatusCode, string(body))
		h.writeErrorResponse(w, resp.StatusCode, fmt.Sprintf("Azure DevOps API error: %s", string(body)))
		return
	}

	// Extract continuation token from response headers
	continuationToken := resp.Header.Get(continuationTokenHeaderKey)
	if continuationToken != "" {
		h.Log.Printf("Continuation token received: %s", continuationToken)
	}

	// Determine if filtering should be applied
	// Filtering is applied only when scopeDescriptor is NOT provided
	shouldFilter := params.ScopeDescriptor == ""
	h.Log.Printf("Should apply filtering: %t", shouldFilter)

	// Process and potentially filter the response
	processedBody, err := h.processGraphGroupsResponse(body, shouldFilter)
	if err != nil {
		h.Log.Printf("Error processing response: %s", err.Error())
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error processing response: %s", err.Error()))
		return
	}

	// Write enhanced JSON response with continuation token
	h.writeJSONResponse(w, http.StatusOK, processedBody, continuationToken)
	h.Log.Printf("Successfully retrieved and processed graph groups for organization: %s", params.Organization)
}
