package pullrequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
)

func PatchPullRequest(opts handlers.HandlerOptions) handlers.Handler {
	return &patchHandler{baseHandler: newBaseHandler(opts)}
}

var _ handlers.Handler = &patchHandler{}

// Base handler with common functionality
type baseHandler struct {
	handlers.HandlerOptions
}

// Constructor for the base handler
func newBaseHandler(opts handlers.HandlerOptions) *baseHandler {
	return &baseHandler{HandlerOptions: opts}
}

type patchHandler struct {
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

	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
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

// Helper to get the current state of a pull request by ID
// For example, to determine which fields need to be updated in a PATCH operation
// since Azure DevOps throws an error if you try to set a field to its current value (`TargetRefName`)
func (h *patchHandler) getPullRequestByID(organization, project, repositoryId, pullRequestId, apiVersion, authHeader string) (*GitPullRequest, error) {
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s/pullrequests/%s?api-version=%s",
		organization, project, repositoryId, pullRequestId, apiVersion)

	resp, err := h.makeAzuredevopsRequest("GET", url, authHeader, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to azure devops: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("azure devops returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	var currentPR GitPullRequest
	if err := json.Unmarshal(body, &currentPR); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pull request response: %w", err)
	}

	return &currentPR, nil
}

// PATCH handler implementation
// @Summary Update a pull request
// @Description Update an existing pull request.
// @ID update-pullrequest
// @Param organization path string true "Organization name"
// @Param project path string true "Project name or ID"
// @Param repositoryId path string true "Repository ID"
// @Param pullRequestId path int true "Pull Request ID"
// @Param api-version query string true "API version (e.g., 7.2-preview.2)"
// @Param Authorization header string true "Basic Auth header"
// @Param pullRequest body UpdatePullRequestReq true "Pull request updates"
// @Accept json
// @Produce json
// @Success 200 {object} GitPullRequest "Updated pull request"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 409 {string} string "Conflict"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/{organization}/{project}/git/repositories/{repositoryId}/pullrequests/{pullRequestId} [patch]
func (h *patchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	organization := r.PathValue("organization")
	if organization == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "organization parameter is required")
		return
	}

	project := r.PathValue("project")
	if project == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "project parameter is required")
		return
	}

	repositoryId := r.PathValue("repositoryId")
	if repositoryId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "repositoryId parameter is required")
		return
	}

	pullRequestId := r.PathValue("pullRequestId")
	if pullRequestId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "pullRequestId parameter is required")
		return
	}

	apiVersion := r.URL.Query().Get("api-version")
	if apiVersion == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "api-version is required")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "request rejected due to missing or invalid Basic authentication")
		return
	}

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("error reading request body: %s", err.Error()))
		return
	}

	// log request body for debugging
	h.Log.Printf("PatchPullRequest received request body: %s", string(requestBody))

	var desiredPR UpdatePullRequestReq
	if err := json.Unmarshal(requestBody, &desiredPR); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("error unmarshaling request body: %s", err.Error()))
		return
	}

	h.Log.Print("Desired pull request state for update:")
	for k, v := range map[string]interface{}{
		"Status":                desiredPR.Status,
		"Title":                 desiredPR.Title,
		"Description":           desiredPR.Description,
		"CompletionOptions":     desiredPR.CompletionOptions,
		"MergeOptions":          desiredPR.MergeOptions,
		"AutoCompleteSetBy":     desiredPR.AutoCompleteSetBy,
		"TargetRefName":         desiredPR.TargetRefName,
		"LastMergeSourceCommit": desiredPR.LastMergeSourceCommit,
	} {
		h.Log.Printf(" - %s: %+v", k, v)
	}

	// Call 1: get the current state of the PR
	currentPR, err := h.getPullRequestByID(organization, project, repositoryId, pullRequestId, apiVersion, authHeader)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error getting current pull request state: %s", err.Error()))
		return
	}

	// Build a patch payload with only the fields that have changed, if any
	patchPayload := make(map[string]interface{})

	// Early exit if `status` is `completed``, we cannot update it anymore
	// {"$id":"1","innerException":null,"message":"TF401181:
	// The pull request cannot be edited due to its state.",
	// "typeName":"Microsoft.TeamFoundation.Git.Server.GitPullRequestNotEditableException,
	//  Microsoft.TeamFoundation.Git.Server","typeKey":"GitPullRequestNotEditableException",
	// "errorCode":0,"eventId":3000}
	if currentPR.Status == "completed" {
		h.writeErrorResponse(w, http.StatusConflict, "cannot update pull request: pull request is already completed. Please change the modified fields in the spec of the Custom Resource to their previous state. After some erroneous update retries, the following observe action will succeed.")
		return
	}

	if desiredPR.Status != "" && desiredPR.Status != currentPR.Status {
		patchPayload["status"] = desiredPR.Status
	}

	if desiredPR.Title != "" && desiredPR.Title != currentPR.Title {
		patchPayload["title"] = desiredPR.Title
	}

	if desiredPR.Description != "" && desiredPR.Description != currentPR.Description {
		patchPayload["description"] = desiredPR.Description
	}

	if desiredPR.CompletionOptions != nil && !reflect.DeepEqual(desiredPR.CompletionOptions, currentPR.CompletionOptions) {
		patchPayload["completionOptions"] = desiredPR.CompletionOptions
	}

	if desiredPR.MergeOptions != nil && !reflect.DeepEqual(desiredPR.MergeOptions, currentPR.MergeOptions) {
		patchPayload["mergeOptions"] = desiredPR.MergeOptions
	}

	if desiredPR.AutoCompleteSetBy != nil && desiredPR.AutoCompleteSetBy.Id != currentPR.AutoCompleteSetBy.Id {
		patchPayload["autoCompleteSetBy"] = desiredPR.AutoCompleteSetBy
	}

	if desiredPR.TargetRefName != "" && desiredPR.TargetRefName != currentPR.TargetRefName {
		patchPayload["targetRefName"] = desiredPR.TargetRefName
	}

	// LastMergeSourceCommit is required only when updating status to "completed"
	// If not set, Azure DevOps returns an error
	// If set in other cases, Azure DevOps returns an error
	//{"$id":"1","innerException":null,
	// "message":"Invalid argument value.\r\nParameter name: You must specify a valid LastMergeSourceCommit to perform this operation.",
	// "typeName":"Microsoft.TeamFoundation.SourceControl.WebServer.InvalidArgumentValueException,
	// Microsoft.TeamFoundation.SourceControl.WebServer",
	// "typeKey":"InvalidArgumentValueException","errorCode":0,"eventId":0}
	if desiredPR.LastMergeSourceCommit != nil && desiredPR.Status == "completed" {
		patchPayload["lastMergeSourceCommit"] = desiredPR.LastMergeSourceCommit
	}

	// If nothing changed, return the current PR and do nothing.
	if len(patchPayload) == 0 {
		h.Log.Print("No changes detected for pull request update, returning current state.")
		currentPRBody, _ := json.Marshal(currentPR)
		h.writeJSONResponse(w, http.StatusOK, currentPRBody)
		return
	}

	h.Log.Print("PatchPullRequest updating pull request with payload")
	for k, v := range patchPayload {
		h.Log.Printf(" - %s: %+v", k, v)
	}

	patchBody, err := json.Marshal(patchPayload)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error marshaling patch payload: %s", err.Error()))
		return
	}

	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s/pullrequests/%s?api-version=%s", organization, project, repositoryId, pullRequestId, apiVersion)

	// Call 2: make the PATCH request to update the PR
	resp, err := h.makeAzuredevopsRequest("PATCH", url, authHeader, patchBody)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error making request to azure devops: %s", err.Error()))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error reading response body: %s", err.Error()))
		return
	}

	h.writeJSONResponse(w, resp.StatusCode, body)
}
