package gitrepository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
)

// Handler constructors
func PostGitRepository(opts handlers.HandlerOptions) handlers.Handler {
	return &postHandler{baseHandler: newBaseHandler(opts)}
}

// Interface compliance verification
var _ handlers.Handler = &postHandler{}

// Base handler with common functionality
type baseHandler struct {
	handlers.HandlerOptions
}

// Constructor for the base handler
func newBaseHandler(opts handlers.HandlerOptions) *baseHandler {
	return &baseHandler{HandlerOptions: opts}
}

// Handler types embedding the base handler
type postHandler struct {
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

// POST handler implementation
// @Summary Create a new GitRepository on Azure DevOps
// @Description Create a new GitRepository on Azure DevOps using the provided organization, project, and repository details.
// @ID post-gitrepository
// @Param organization path string true "Organization name"
// @Param projectId path string true "Project ID or name"
// @Param api-version query string true "API version (e.g., 7.2-preview.2)"
// @Param sourceRef query string false "Specify the source refs to use while creating a fork repo"
// @Param Authorization header string true "Basic Auth header (Basic <base64-encoded-username:password>)"
// @Param gitrepositoryCreate body CreateRepositoryRequest true "GitRepository creation request body (with additional fields handled by the plugin)"
// @Accept json
// @Produce json
// @Success 201 {object} CreateRepositoryResponse "GitRepository details"
// @Success 202 {object} CreateRepositoryResponse "GitRepository details (repo created but creation of branch deisgnated as default branch is pending, user must create it, then the gitrepository-controller will update the default branch later)"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Router /api/{organization}/{projectId}/git/repositories [post]
func (h *postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	organization := r.PathValue("organization")
	projectId := r.PathValue("projectId")
	apiVersion := r.URL.Query().Get("api-version")
	sourceRef := r.URL.Query().Get("sourceRef") // Optional sourceRef parameter
	authHeader := r.Header.Get("Authorization")

	// Validate required parameters
	if organization == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Organization parameter is required")
		return
	}
	if projectId == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Project ID parameter is required")
		return
	}
	if apiVersion == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "API version parameter is required")
		return
	}

	// This single check handles missing headers, incorrect formats, and empty credentials.
	username, password, ok := r.BasicAuth()
	if !ok || username == "" || password == "" {
		h.writeErrorResponse(w, http.StatusUnauthorized, "Request rejected due to missing or invalid Basic authentication")
		return
	}

	// Read and parse the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var createRequest CreateRepositoryRequest
	if err := json.Unmarshal(body, &createRequest); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON in request body")
		return
	}

	// Validate required fields
	if createRequest.Name == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Repository name is required")
		return
	}

	requestedDefaultBranch := createRequest.DefaultBranch

	// Decide if a defaultBranch update is needed
	needsDefaultBranchUpdate := requestedDefaultBranch != ""

	if requestedDefaultBranch == "" {
		h.Log.Print("No default branch specified in request")
	} else {
		h.Log.Printf("Custom default branch '%s' specified", requestedDefaultBranch)
	}

	// Validate sourceRef if provided
	valid, err := h.validateSourceRef(organization, projectId, &createRequest, sourceRef, apiVersion, authHeader, w)
	if !valid || err != nil {
		return // Stop execution if this validation failed
	}

	// For new repositories (not forks), if a default branch is specified, initialization must be explicitly enabled.
	// Since Azure DevOps does not allow setting a default branch on an uninitialized repository (empty repository)
	// Therefore, we return a 400 Bad Request if this condition is not met.
	if createRequest.ParentRepository == nil && needsDefaultBranchUpdate && !createRequest.Initialize {
		h.writeErrorResponse(w, http.StatusBadRequest, "When specifying a 'defaultBranch' for a new repository, 'initialize' must be set to true")
		return
	}

	// Create the repository request body for Azure DevOps (without defaultBranch)
	azureDevOpsRequest := GitRepositoryCreateOptionsMinimal{
		Name:             createRequest.Name,
		ParentRepository: createRequest.ParentRepository,
		Project:          createRequest.Project,
		// Note: defaultBranch is intentionally omitted for the initial POST request as POST does not support it
	}

	// 1. Create repository (fork or new)
	createdRepo, err := h.createGitRepository(organization, projectId, apiVersion, authHeader, sourceRef, azureDevOpsRequest)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create repository: %v", err))
		return
	}

	defaultBranchCreationPending := false

	// 2. Handle initialization vs fork logic
	if createRequest.ParentRepository != nil {
		// This is a fork - probably some branches already exist from parent
		// we could check if the requested default branch exists in the parent repository
		h.Log.Printf("Created fork repository '%s' from parent repository '%s'", createdRepo.Name, createRequest.ParentRepository.ID)

		// For forks, check if the desired default branch exists thanks to sourceRef from parent
		if needsDefaultBranchUpdate {
			exists, err := h.branchExists(organization, projectId, createdRepo.ID, requestedDefaultBranch, apiVersion, authHeader)
			if err != nil {
				h.Log.Printf("Error checking if branch '%s' exists in fork: %v", requestedDefaultBranch, err)
				h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to check if branch '%s' exists in fork: %v", requestedDefaultBranch, err))
				return
			}
			if !exists {
				// we will return 202 accepted in this case to indicate that the branch creation is pending
				h.Log.Printf("Branch '%s' does not exist in fork repository '%s'. Cannot set it as default. However, repository is created successfully", requestedDefaultBranch, createdRepo.Name)
				defaultBranchCreationPending = true
				needsDefaultBranchUpdate = false // No need to set default branch now (it will raise an error), it needs to be created later by the user
			}
			h.Log.Printf("Branch '%s' exists in fork repository '%s'", requestedDefaultBranch, createdRepo.Name)
			// If the branch exists, we can proceed to set it as default branch later
		}
	} else {
		// This is a new repository - handle initialization and branch creation
		if createRequest.Initialize {
			// Repository is to be initialized
			h.Log.Printf("Repository '%s' is not a fork, proceeding with initialization", createdRepo.Name)

			// Initialize with the requested branch (or default if none specified)
			initBranch := requestedDefaultBranch
			if initBranch == "" {
				initBranch = "refs/heads/main"
			}

			h.Log.Printf("Repository '%s' will be initialized with an initial commit on branch '%s'", createdRepo.Name, initBranch)

			if err := h.initializeRepository(organization, projectId, createdRepo.ID, authHeader, initBranch); err != nil {
				h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to initialize repository '%s': %v", createdRepo.Name, err))
				return
			}
			h.Log.Printf("Successfully initialized repository '%s' with an initial commit", createdRepo.Name)
			// After initialization, the branch exists, so we can and go directly to setting it as default if needed
		} else {
			// Repository is not initialized and is not a fork
			h.Log.Printf("Repository '%s' was created without initialization. No branches exist yet.", createdRepo.Name)

			// If a default branch was requested, this is an error since we cannot set a default branch on an uninitialized repo
			if needsDefaultBranchUpdate {
				h.Log.Printf("Cannot set default branch '%s' on uninitialized repository '%s' - no branches exist", requestedDefaultBranch, createdRepo.Name)
				h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Cannot set default branch '%s' on uninitialized repository '%s' - no branches exist", requestedDefaultBranch, createdRepo.Name))
				return
			}
		}
	}

	// 3. Update default branch if needed
	// At this point, we know branches exist (either from fork or initialization)
	if needsDefaultBranchUpdate {
		h.Log.Printf("Setting default branch to '%s' for repository '%s'", requestedDefaultBranch, createdRepo.Name)

		// At this point, we are sure the branch exists (either it was created or it was already there)
		h.Log.Printf("Updating repository '%s' to set default branch to '%s'", createdRepo.Name, requestedDefaultBranch)
		updatedRepo, err := h.updateRepositoryDefaultBranch(organization, projectId, createdRepo.ID, requestedDefaultBranch, apiVersion, authHeader)
		if err != nil {
			h.Log.Printf("Failed to set default branch '%s': %v", requestedDefaultBranch, err)
			h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to set default branch '%s': %v", requestedDefaultBranch, err))
			return
		}
		// Use the updated repository information
		createdRepo = updatedRepo
		h.Log.Printf("Successfully set default branch to '%s'", requestedDefaultBranch)

	}

	// Convert the response to the expected format which is GitRepository
	response := CreateRepositoryResponse(*createdRepo)

	responseBody, err := json.Marshal(response)
	if err != nil {
		h.Log.Printf("Failed to marshal response: %v", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to marshal response")
		return
	}

	// Forked repo case:
	// If the branch creation is pending (currently in the newly forked repo there is not that branch, since is not inherited from the parent),
	// we return 202 Accepted (GitRepository created but the requested default branch is not yet created)
	// this is to indicate that the user must create the branch first in the newly forked repository,
	// then the gitrepository-controller will update the default branch later
	// But for now, the default branch is the one of the parent repository
	if defaultBranchCreationPending {
		h.Log.Printf("Branch '%s' creation is pending in fork repository '%s'. Returning 202 Accepted", requestedDefaultBranch, createdRepo.Name)
		h.writeJSONResponse(w, http.StatusAccepted, responseBody)
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, responseBody)
	h.Log.Printf("Successfully created GitRepository '%s' in organization '%s', project '%s'", createdRepo.Name, organization, projectId)
}

// createGitRepository performs the actual repository creation via Azure DevOps API
func (h *postHandler) createGitRepository(organization, projectId, apiVersion, authHeader, sourceRef string, request GitRepositoryCreateOptionsMinimal) (*GitRepository, error) {
	// Construct the URL for the Azure DevOps API
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories?api-version=%s", organization, projectId, apiVersion)

	// check if sourceRef is provided and append it to the URL
	if sourceRef != "" {
		url += fmt.Sprintf("&sourceRef=%s", sourceRef)
	}

	// Marshal the request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal create request: %w", err)
	}

	h.Log.Printf("Creating repository with request body: %s", string(requestBody))

	// Make the POST request to Azure DevOps API
	resp, err := h.makeAzuredevopsRequest("POST", url, authHeader, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to make create repository request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read create repository response: %w", err)
	}

	// Check for non-201 status codes
	if resp.StatusCode != http.StatusCreated {
		h.Log.Printf("Azure DevOps API returned non-201 status for repository creation: %d. Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("azure DevOps API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var createdRepo GitRepository
	if err := json.Unmarshal(body, &createdRepo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create repository response: %w", err)
	}

	return &createdRepo, nil
}

// updateRepositoryDefaultBranch updates the default branch of an existing repository
func (h *postHandler) updateRepositoryDefaultBranch(organization, projectId, repositoryId, defaultBranch, apiVersion, authHeader string) (*GitRepository, error) {
	// Construct the URL for the Azure DevOps API
	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s?api-version=%s", organization, projectId, repositoryId, apiVersion)

	// Create the update request
	updateRequest := GitRepositoryUpdateOptions{
		DefaultBranch: defaultBranch,
	}

	h.Log.Printf("Repository update request body: %+v", updateRequest)

	// Marshal the request body
	requestBody, err := json.Marshal(updateRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update request: %w", err)
	}

	h.Log.Printf("Updating repository default branch with request body: %s", string(requestBody))

	// Make the PATCH request to Azure DevOps API
	resp, err := h.makeAzuredevopsRequest("PATCH", url, authHeader, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to make update repository request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read update repository response: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		h.Log.Printf("Azure DevOps API returned non-200 status for repository update: %d. Body: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("azure DevOps API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var updatedRepo GitRepository
	if err := json.Unmarshal(body, &updatedRepo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update repository response: %w", err)
	}

	return &updatedRepo, nil
}

// Helper functions in the postHandler for:
// - repository initialization
// - branch existence check
// - validating sourceRef

func (h *postHandler) initializeRepository(organization, projectId, repositoryId, authHeader, branchToInit string) error {
	// Ensure branch name has proper format
	if !strings.HasPrefix(branchToInit, "refs/heads/") {
		branchToInit = "refs/heads/" + branchToInit
	}

	// Git Pushes may require a different API version, we pass it via env variable GIT_PUSHES_API_VERSION

	apiVersion := os.Getenv("GIT_PUSHES_API_VERSION")
	if apiVersion == "" {
		h.Log.Print("GIT_PUSHES_API_VERSION environment variable not set, using default API version")
		apiVersion = "7.2-preview.3" // Default Git Pushes API version if not set
	}

	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s/pushes?api-version=%s", organization, projectId, repositoryId, apiVersion)

	requestBody := map[string]interface{}{
		"refUpdates": []map[string]string{
			{
				"name":        branchToInit,
				"oldObjectId": "0000000000000000000000000000000000000000",
			},
		},
		"commits": []map[string]interface{}{
			{
				"comment": "Initial commit",
				"changes": []map[string]interface{}{
					{
						"changeType": "add",
						"item": map[string]string{
							"path": "/README.md",
						},
						"newContent": map[string]string{
							"content":     "# New Repository\n\nThis repository was initialized automatically.",
							"contentType": "rawtext",
						},
					},
				},
			},
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal initial commit request: %w", err)
	}

	h.Log.Printf("Initializing repository with request body: %s", string(body))

	resp, err := h.makeAzuredevopsRequest("POST", url, authHeader, body)
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to initialize repository, status: %d, body: %s", resp.StatusCode, string(responseBody))
	}

	h.Log.Printf("Successfully initialized repository '%s' with branch '%s'", repositoryId, branchToInit)
	return nil
}

func (h *postHandler) branchExists(organization, projectId, repositoryId, branchName, apiVersion, authHeader string) (bool, error) {
	// Remove the 'refs/heads/' prefix if present for the `refs` API endpoint
	branchNameForAPI := strings.TrimPrefix(branchName, "refs/heads/")

	url := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/git/repositories/%s/refs?filter=heads/%s&api-version=%s", organization, projectId, repositoryId, branchNameForAPI, apiVersion)

	h.Log.Printf("Checking if branch '%s' exists in repository '%s'", branchNameForAPI, repositoryId)
	h.Log.Printf("Branch existence check URL: %s", url)

	// Sleep to allow any potential delays in branch creation
	// Without this, we might check too early after creation and it has been observed that Azure DevOps API can take a moment to reflect new branches
	// This is a workaround for the observed delay in branch creation visibility
	// Minimum wait time tested: 1 second
	h.Log.Printf("Waiting some seconds before checking branch existence...")
	time.Sleep(1 * time.Second)

	resp, err := h.makeAzuredevopsRequest("GET", url, authHeader, nil)
	if err != nil {
		return false, fmt.Errorf("failed to check branch existence: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to check branch existence, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read branch check response: %w", err)
	}

	var refsResponse struct {
		Value []interface{} `json:"value"`
	}

	if err := json.Unmarshal(body, &refsResponse); err != nil {
		return false, fmt.Errorf("failed to unmarshal refs response: %w", err)
	}

	// log value for debugging
	h.Log.Printf("Branch existence check response: %s", string(body))

	return len(refsResponse.Value) > 0, nil
}

func (h *postHandler) validateSourceRef(organization string, projectId string, createRequest *CreateRepositoryRequest, sourceRef, apiVersion, authHeader string, w http.ResponseWriter) (bool, error) {

	if createRequest.ParentRepository != nil {

		// check if sourceRef is provided and if it is a valid branch
		if sourceRef != "" {
			h.Log.Printf("Original SourceRef provided: %s", sourceRef)

			// decode sourceRef if it is URL encoded
			sourceRef, err := url.QueryUnescape(sourceRef)
			if err != nil {
				h.Log.Printf("Error decoding sourceRef '%s': %v", sourceRef, err)
				h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid sourceRef '%s': %v", sourceRef, err))
				return false, err
			}
			h.Log.Printf("Decoded SourceRef: %s", sourceRef)

			if !strings.HasPrefix(sourceRef, "refs/heads/") {
				h.writeErrorResponse(w, http.StatusBadRequest, "sourceRef must start with 'refs/heads/'")
				return false, fmt.Errorf("sourceRef must start with 'refs/heads/'")
			}
			exists, err := h.branchExists(organization, projectId, createRequest.ParentRepository.ID, sourceRef, apiVersion, authHeader)
			if err != nil {
				h.Log.Printf("Error checking if sourceRef '%s' exists in parent repository '%s': %v", sourceRef, createRequest.ParentRepository.ID, err)
				h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to check if sourceRef '%s' exists in parent repository '%s': %v", sourceRef, createRequest.ParentRepository.ID, err))
				return false, err
			}
			if exists {
				h.Log.Printf("SourceRef '%s' exists in parent repository '%s'", sourceRef, createRequest.ParentRepository.ID)
			} else {
				h.Log.Printf("SourceRef '%s' does not exist in parent repository '%s'", sourceRef, createRequest.ParentRepository.ID)
				h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("SourceRef '%s' does not exist in parent repository '%s'", sourceRef, createRequest.ParentRepository.ID))
				return false, fmt.Errorf("sourceRef '%s' does not exist in parent repository '%s'", sourceRef, createRequest.ParentRepository.ID)
			}
		}
		h.Log.Printf("SourceRef validation complete for repository '%s': SourceRef='%s'", createRequest.Name, sourceRef)
	} else {
		// No parent repository, so no sourceRef validation needed
		h.Log.Printf("No parent repository, skipping sourceRef validation for repository '%s'", createRequest.Name)
		if sourceRef != "" {
			h.Log.Printf("Warning: sourceRef '%s' provided for a new repository without a parent repository. This will be ignored.", sourceRef)
		}
	}
	return true, nil
}
