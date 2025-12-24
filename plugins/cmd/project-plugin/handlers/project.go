package project

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
)

const (
	requestTimeout               = 30 * time.Second
	contentTypeJSON              = "application/json"
	maxPollingAttempts           = 10
	pollingInterval              = 1 * time.Second // Wait time between polling attempts
	authErrorMessage             = "request rejected due to missing or invalid Basic authentication"
	azureDevOpsBaseURL           = "https://dev.azure.com"
	projectsAPIPath              = "_apis/projects"
	operationsAPIPath            = "_apis/operations"
	operationsAPIVersionFallback = "7.2-preview"
)

// PostProject returns a handler for POST project creation requests
func PostProject(opts handlers.HandlerOptions) handlers.Handler {
	return &postHandler{baseHandler: newBaseHandler(opts)}
}

// PatchProject returns a handler for PATCH project update requests
func PatchProject(opts handlers.HandlerOptions) handlers.Handler {
	return &patchHandler{baseHandler: newBaseHandler(opts)}
}

// DeleteProject returns a handler for DELETE project deletion requests
func DeleteProject(opts handlers.HandlerOptions) handlers.Handler {
	return &deleteHandler{baseHandler: newBaseHandler(opts)}
}

// Compile-time interface checks
var (
	_ handlers.Handler = (*postHandler)(nil)
	_ handlers.Handler = (*patchHandler)(nil)
	_ handlers.Handler = (*deleteHandler)(nil)
)

// baseHandler contains common functionality shared by all handlers
type baseHandler struct {
	handlers.HandlerOptions
}

// newBaseHandler creates a new base handler with the given options
func newBaseHandler(opts handlers.HandlerOptions) *baseHandler {
	return &baseHandler{HandlerOptions: opts}
}

type postHandler struct {
	*baseHandler
}

type patchHandler struct {
	*baseHandler
}

type deleteHandler struct {
	*baseHandler
}

// requestParams holds common request parameters extracted from HTTP requests
type requestParams struct {
	Organization string
	ProjectID    string // Used for PATCH and DELETE
	APIVersion   string
	AuthHeader   string
}

// extractRequestParams extracts and validates parameters from HTTP requests
func (h *baseHandler) extractRequestParams(r *http.Request, needsProjectID bool) (*requestParams, error) {
	params := &requestParams{
		Organization: r.PathValue("organization"),
		ProjectID:    r.PathValue("projectId"), // If not needed, can be empty
		APIVersion:   r.URL.Query().Get("api-version"),
		AuthHeader:   r.Header.Get("Authorization"),
	}

	if err := params.validate(needsProjectID); err != nil {
		return nil, err
	}

	return params, nil
}

// validate checks if all required parameters are present and valid
func (p *requestParams) validate(needsProjectID bool) error {
	if p.Organization == "" {
		return fmt.Errorf("organization parameter is required")
	}
	if needsProjectID && p.ProjectID == "" {
		return fmt.Errorf("projectId parameter is required")
	}
	if p.APIVersion == "" {
		return fmt.Errorf("api-version is required")
	}
	if p.AuthHeader == "" {
		return fmt.Errorf(authErrorMessage)
	}

	return nil
}

// buildProjectURL constructs the Azure DevOps project API URL
func buildProjectURL(organization, projectID, apiVersion string) string {
	if projectID == "" {
		return fmt.Sprintf("%s/%s/%s?api-version=%s", azureDevOpsBaseURL, organization, projectsAPIPath, apiVersion)
	}
	return fmt.Sprintf("%s/%s/%s/%s?api-version=%s", azureDevOpsBaseURL, organization, projectsAPIPath, projectID, apiVersion)
}

// buildOperationURL constructs the Azure DevOps operation status API URL
func buildOperationURL(organization, operationID string) string {
	operationsAPIVersion := os.Getenv("OPERATIONS_API_VERSION")
	if operationsAPIVersion == "" {
		operationsAPIVersion = operationsAPIVersionFallback // Fallback Operations API version if not set
	}
	return fmt.Sprintf("%s/%s/%s/%s?api-version=%s", azureDevOpsBaseURL, organization, operationsAPIPath, operationID, operationsAPIVersion)
}

// statusCodeFromError determines the HTTP status code based on the error type
func statusCodeFromError(err error) int {
	if err.Error() == authErrorMessage {
		return http.StatusUnauthorized
	}
	return http.StatusBadRequest
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
func (h *baseHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(statusCode)
	w.Write(body)
}

// transformOperationReference transforms an OperationReference to a TransformedOperationReference
// by prefixing all field names with "operation"
func transformOperationReference(opRef *OperationReference) *TransformedOperationReference {
	return &TransformedOperationReference{
		OperationID:       opRef.ID,
		OperationPluginID: opRef.PluginID,
		OperationStatus:   opRef.Status,
		OperationURL:      opRef.URL,
	}
}

// callAzureDevOpsAndParseOperation calls Azure DevOps API and parses the operation reference from response
func (h *baseHandler) callAzureDevOpsAndParseOperation(ctx context.Context, method, azureURL string, params *requestParams, body []byte) (*OperationReference, error) {
	h.Log.Printf("Calling Azure DevOps API: %s %s", method, azureURL)

	// Call Azure DevOps API
	resp, err := h.makeAzureDevOpsRequest(ctx, method, azureURL, params.AuthHeader, body)
	if err != nil {
		return nil, fmt.Errorf("error calling Azure DevOps API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check for non-success status codes
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		h.Log.Printf("Azure DevOps API returned non-success status: %d - %s", resp.StatusCode, string(respBody))
		return nil, fmt.Errorf("Azure DevOps API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Parse the OperationReference from Azure DevOps
	var opRef OperationReference
	if err := json.Unmarshal(respBody, &opRef); err != nil {
		return nil, fmt.Errorf("failed to unmarshal operation reference: %w", err)
	}

	h.Log.Printf("Received operation reference: id=%s, status=%s", opRef.ID, opRef.Status)
	return &opRef, nil
}

// handlePostOperation handles POST operation - returns transformed operation reference (async)
func (h *baseHandler) handlePostOperation(ctx context.Context, w http.ResponseWriter, azureURL string, params *requestParams, body []byte) {
	// Call Azure DevOps and parse operation reference
	opRef, err := h.callAzureDevOpsAndParseOperation(ctx, http.MethodPost, azureURL, params, body)
	if err != nil {
		statusCode := http.StatusInternalServerError
		h.writeErrorResponse(w, statusCode, err.Error())
		return
	}

	// Transform the operation reference
	transformed := transformOperationReference(opRef)
	h.Log.Printf("Transformed to: operationId=%s, operationStatus=%s", transformed.OperationID, transformed.OperationStatus)

	// Marshal the transformed response
	transformedBody, err := json.Marshal(transformed)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal transformed response: %s", err.Error()))
		return
	}

	// Write the transformed response with 202 Accepted
	h.writeJSONResponse(w, http.StatusAccepted, transformedBody)
	h.Log.Printf("Successfully transformed POST response for organization: %s", params.Organization)
}

// pollOperation polls the operation endpoint until it succeeds or fails
func (h *baseHandler) pollOperation(ctx context.Context, operationID, organization, authHeader string) error {
	operationURL := buildOperationURL(organization, operationID)

	for attempt := 0; attempt < maxPollingAttempts; attempt++ {
		time.Sleep(pollingInterval)

		h.Log.Printf("Polling operation %s (attempt %d/%d)", operationID, attempt+1, maxPollingAttempts)

		resp, err := h.makeAzureDevOpsRequest(ctx, http.MethodGet, operationURL, authHeader, nil)
		if err != nil {
			return fmt.Errorf("failed to poll operation: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // Close immediately after reading

		if err != nil {
			return fmt.Errorf("failed to read operation response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("operation polling failed with status %d: %s", resp.StatusCode, string(body))
		}

		var opRef OperationReference
		if err := json.Unmarshal(body, &opRef); err != nil {
			return fmt.Errorf("failed to unmarshal operation status: %w", err)
		}

		h.Log.Printf("Operation status: %s", opRef.Status)

		switch opRef.Status {
		case "succeeded":
			h.Log.Printf("Operation %s succeeded", operationID)
			return nil
		case "failed", "cancelled":
			return fmt.Errorf("operation %s failed with status: %s", operationID, opRef.Status)
		}
		// Continue polling if status is notSet, queued, or inProgress
	}

	return fmt.Errorf("operation %s timed out after %d attempts", operationID, maxPollingAttempts)
}

// handlePatchOperation handles PATCH operation - polls operation until complete, then returns actual updated project (sync)
func (h *baseHandler) handlePatchOperation(ctx context.Context, w http.ResponseWriter, azureURL string, params *requestParams, body []byte) {
	// Call Azure DevOps and parse operation reference
	opRef, err := h.callAzureDevOpsAndParseOperation(ctx, http.MethodPatch, azureURL, params, body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Printf("Starting polling for operation %s", opRef.ID)

	// Poll the operation until it completes
	if err := h.pollOperation(ctx, opRef.ID, params.Organization, params.AuthHeader); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("operation polling failed: %s", err.Error()))
		return
	}

	// Operation succeeded - now get the actual updated project
	projectURL := buildProjectURL(params.Organization, params.ProjectID, params.APIVersion)

	h.Log.Printf("Fetching updated project: GET %s", projectURL)

	projectResp, err := h.makeAzureDevOpsRequest(ctx, http.MethodGet, projectURL, params.AuthHeader, nil)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error fetching updated project: %s", err.Error()))
		return
	}
	defer projectResp.Body.Close()

	projectBody, err := io.ReadAll(projectResp.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error reading project response: %s", err.Error()))
		return
	}

	if projectResp.StatusCode != http.StatusOK {
		h.Log.Printf("Failed to fetch project: %d - %s", projectResp.StatusCode, string(projectBody))
		h.writeErrorResponse(w, projectResp.StatusCode, fmt.Sprintf("error fetching project: %s", string(projectBody)))
		return
	}

	// Unmarshal the project response into ProjectResponseBody struct
	var project ProjectResponseBody
	if err := json.Unmarshal(projectBody, &project); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to unmarshal project response: %s", err.Error()))
		return
	}

	// Marshal the project back to JSON (validates structure)
	responseBody, err := json.Marshal(project)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal project response: %s", err.Error()))
		return
	}

	// Return the actual project data with 200 OK (operation completed successfully)
	h.writeJSONResponse(w, http.StatusOK, responseBody)
	h.Log.Printf("Successfully completed PATCH operation and returned updated project for organization: %s", params.Organization)
}

// handleDeleteOperation handles DELETE operation - polls operation until complete, then returns 204 No Content (sync)
func (h *baseHandler) handleDeleteOperation(ctx context.Context, w http.ResponseWriter, azureURL string, params *requestParams) {
	// Call Azure DevOps and parse operation reference
	opRef, err := h.callAzureDevOpsAndParseOperation(ctx, http.MethodDelete, azureURL, params, nil)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.Log.Printf("Starting polling for delete operation %s", opRef.ID)

	// Poll the operation until it completes
	if err := h.pollOperation(ctx, opRef.ID, params.Organization, params.AuthHeader); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("operation polling failed: %s", err.Error()))
		return
	}

	// Operation succeeded - return 204 No Content
	w.WriteHeader(http.StatusNoContent)
	h.Log.Printf("Successfully completed DELETE operation for project %s in organization: %s", params.ProjectID, params.Organization)
}

// ServeHTTP handles POST requests for project creation
// @Summary Create a project
// @Description Create a new project in Azure DevOps. Returns a transformed operation reference.
// @ID Projects_Create
// @Tags Projects
// @Accept json
// @Produce json
// @Param organization path string true "The name of the Azure DevOps organization."
// @Param api-version query string true "API version (e.g., 7.2-preview.4)"
// @Param Authorization header string true "Basic Auth header"
// @Param project body ProjectRequestBodyMerged true "The project to create."
// @Success 202 {object} TransformedOperationReference "Transformed operation reference with prefixed fields"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /plugin/{organization}/projects [post]
func (h *postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	// Extract and validate request parameters (POST doesn't need projectId)
	params, err := h.extractRequestParams(r, false)
	if err != nil {
		h.writeErrorResponse(w, statusCodeFromError(err), err.Error())
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to read request body: %s", err.Error()))
		return
	}
	defer r.Body.Close()

	// Unmarshal and validate request body
	// In the annotations, we use ProjectRequestBodyMerged to include both create and update fields,
	// but here we only need to use and validate the create fields.
	var projectReq ProjectRequestBodyCreate
	if err := json.Unmarshal(body, &projectReq); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON in request body: %s", err.Error()))
		return
	}

	// Validate project request body
	if err := projectReq.Validate(); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Marshal back to JSON for Azure DevOps API call
	validatedBody, err := json.Marshal(projectReq)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal project request: %s", err.Error()))
		return
	}

	// Build Azure DevOps API URL
	azureURL := buildProjectURL(params.Organization, "", params.APIVersion)

	// Handle POST operation (async - returns transformed operation reference)
	h.handlePostOperation(ctx, w, azureURL, params, validatedBody)
}

// ServeHTTP handles PATCH requests for project updates
// @Summary Update a project
// @Description Update an existing project in Azure DevOps. Polls the operation until complete and returns the updated project.
// @ID Projects_Update
// @Tags Projects
// @Accept json
// @Produce json
// @Param organization path string true "The name of the Azure DevOps organization."
// @Param projectId path string true "Project ID"
// @Param api-version query string true "API version (e.g., 7.2-preview.4)"
// @Param Authorization header string true "Basic Auth header"
// @Param project body ProjectRequestBodyUpdate true "Project update request"
// @Success 200 {object} ProjectResponseBody "Updated project object"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /plugin/{organization}/projects/{projectId} [patch]
func (h *patchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	// Extract and validate request parameters (PATCH needs projectId)
	params, err := h.extractRequestParams(r, true)
	if err != nil {
		h.writeErrorResponse(w, statusCodeFromError(err), err.Error())
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to read request body: %s", err.Error()))
		return
	}
	defer r.Body.Close()

	// Unmarshal and validate request body
	var projectReq ProjectRequestBodyUpdate
	if err := json.Unmarshal(body, &projectReq); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON in request body: %s", err.Error()))
		return
	}

	// Marshal back to JSON for Azure DevOps API call
	validatedBody, err := json.Marshal(projectReq)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal project request: %s", err.Error()))
		return
	}

	// Build Azure DevOps API URL
	azureURL := buildProjectURL(params.Organization, params.ProjectID, params.APIVersion)

	// Handle PATCH operation (sync - polls operation and returns updated project)
	h.handlePatchOperation(ctx, w, azureURL, params, validatedBody)
}

// ServeHTTP handles DELETE requests for project deletion
// @Summary Delete a project
// @Description Delete an existing project in Azure DevOps. Polls the operation until complete and returns 204 No Content.
// @ID Projects_Delete
// @Tags Projects
// @Produce json
// @Param organization path string true "The name of the Azure DevOps organization."
// @Param projectId path string true "Project ID"
// @Param api-version query string true "API version (e.g., 7.2-preview.4)"
// @Param Authorization header string true "Basic Auth header"
// @Success 204 "No Content - project deleted successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /plugin/{organization}/projects/{projectId} [delete]
func (h *deleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	// Extract and validate request parameters (DELETE needs projectId)
	params, err := h.extractRequestParams(r, true)
	if err != nil {
		h.writeErrorResponse(w, statusCodeFromError(err), err.Error())
		return
	}

	// Build Azure DevOps API URL
	azureURL := buildProjectURL(params.Organization, params.ProjectID, params.APIVersion)

	// Handle DELETE operation (sync - polls operation and returns 204 No Content)
	h.handleDeleteOperation(ctx, w, azureURL, params)
}
