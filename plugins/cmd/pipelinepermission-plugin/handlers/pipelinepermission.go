package pipelinepermission

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
)

// Handler constructors
func GetPipelinePermission(opts handlers.HandlerOptions) handlers.Handler {
	return &getHandler{baseHandler: newBaseHandler(opts)}
}

// Interface compliance verification
var _ handlers.Handler = &getHandler{}

// Base handler with common functionality
type baseHandler struct {
	handlers.HandlerOptions
}

// Constructor for the base handler
func newBaseHandler(opts handlers.HandlerOptions) *baseHandler {
	return &baseHandler{HandlerOptions: opts}
}

// Handler types embedding the base handler
type getHandler struct {
	*baseHandler
}

// Common methods, defined once on baseHandler
func (h *baseHandler) makeAzuredevopsRequest(method, url string, authHeader string, body []byte) (*http.Response, error) {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Basic Auth header (Azure DevOps requires Basic Auth for API access)
	if authHeader != "" {
		h.Log.Print("Using provided Authorization header for Basic authentication")
		req.Header.Set("Authorization", authHeader)
	} else {
		h.Log.Print("No Authorization header provided, Basic authentication required")
		return nil, fmt.Errorf("no Authorization header provided, Basic authentication required")
	}

	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

func (h *baseHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	h.Log.Print(message)
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func (h *baseHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(body)
}

// GET handler implementation
// @Summary Get the pipeline permission of a resource
// @Description Get
// @ID get-pipelinepermission
// @Param organization path string true "Organization name"
// @Param project path string true "Project name"
// @Param resourceType path string true "Resource type (e.g., pipelines, repositories)"
// @Param resourceId path string true "Resource ID (e.g., pipeline ID, repository ID)"
// @Param api-version query string true "API version (e.g., 7.2-preview.1)"
// @Param Authorization header string true "Basic Auth header (Basic <base64-encoded-username:password>)"
// @Produce json
// @Success 200 {object} ResourcePipelinePermissions "Pipeline permission details"
// @Router /api/{organization}/{project}/pipelines/pipelinepermissions/{resourceType}/{resourceId} [get]
func (h *getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	organization := r.PathValue("organization")
	project := r.PathValue("project")
	resourceType := r.PathValue("resourceType")
	resourceId := r.PathValue("resourceId")
	apiVersion := r.URL.Query().Get("api-version")
	authHeader := r.Header.Get("Authorization")

	// This single check handles missing headers, incorrect formats, and empty credentials.
	username, password, ok := r.BasicAuth()
	if !ok || username == "" || password == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Request rejected due to missing or invalid Basic authentication")
		return
	}

	h.Log.Printf("Getting pipeline permission for resource %s/%s/%s/%s", organization, project, resourceType, resourceId)

	// Get PipelinePermission
	err := h.getPipelinePermissionAndRespond(w, organization, project, resourceType, resourceId, apiVersion, authHeader)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error getting pipeline permission: %v", err))
	}
}

func (h *getHandler) getPipelinePermissionAndRespond(w http.ResponseWriter, organization, project, resourceType, resourceId, apiVersion, authHeader string) error {
	// Construct the URL for the Azure DevOps API
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/pipelines/pipelinepermissions/%s/%s?api-version=%s", organization, project, resourceType, resourceId, apiVersion)

	// Make the request to Azure DevOps API
	resp, err := h.makeAzuredevopsRequest("GET", url, authHeader, nil)
	if err != nil {
		return fmt.Errorf("failed to get pipeline permission: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body first, as it's useful for both success and error cases
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 status codes after reading the body
	if resp.StatusCode != http.StatusOK {
		h.Log.Printf("Azure DevOps API returned a non-200 status: %d. Body: %s", resp.StatusCode, string(body))
		if resp.StatusCode == http.StatusNotFound {
			h.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Pipeline permission for resource %s/%s/%s/%s not found", organization, project, resourceType, resourceId))
			return nil
		}
		// For other non-200 responses, log the error and return the raw response
		h.writeJSONResponse(w, resp.StatusCode, body)
		return nil
	}

	// Process the response body through the transformation pipeline
	processedBody, err := h.processPipelinepermissionResponse(body)
	if err != nil {
		h.Log.Printf("Failed to process response, returning original raw response from Azure DevOps API: %v", err)
		// Fallback to sending the original, valid response
		h.writeJSONResponse(w, http.StatusOK, body)
		return nil
	}

	h.writeJSONResponse(w, http.StatusOK, processedBody)
	h.Log.Printf("Successfully retrieved pipeline permission for resource %s/%s/%s/%s", organization, project, resourceType, resourceId)
	return nil
}

func (h *getHandler) processPipelinepermissionResponse(body []byte) ([]byte, error) {

	finalBody := body

	// Adjusting the response body if needed
	// Check if the field 'allPipelines' exists in the response body
	if _, err := ReadFieldFromBody(body, "allPipelines"); err != nil {

		// If the field does not exist, we will add it
		h.Log.Printf("Field 'allPipelines' not found in response body: %v", err)

		// Add the field 'allPipelines' with value 'authorized: false' (map)
		authorizedFalse := map[string]interface{}{
			"authorized": false,
		}

		finalBody, err = AddFieldToBody(body, "allPipelines", authorizedFalse)
		if err != nil {
			return nil, fmt.Errorf("failed to add message field: %w", err)
		}
		h.Log.Printf("Field 'allPipelines' added to response body with value 'false'")
		return finalBody, nil

	}

	// If the field was found, no processing is needed (happy path)
	h.Log.Printf("Field 'allPipelines' already found in AzureDevOps raw response body, body not modified")
	return finalBody, nil

	// Add message field
	// If added, types need to be adjusted accordingly
	//message := fmt.Sprintf()
	//finalBody, err := AddFieldToResponse(correctedBody, "message", message)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to add message field: %w", err)
	//}

}
