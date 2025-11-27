package buildpermission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
)

const (
	requestTimeout  = 30 * time.Second
	contentTypeJSON = "application/json"
)

func GetBuildPermission(opts handlers.HandlerOptions) handlers.Handler {
	return &getHandler{baseHandler: newBaseHandler(opts)}
}

func PostBuildPermission(opts handlers.HandlerOptions) handlers.Handler {
	return &postHandler{baseHandler: newBaseHandler(opts)}
}

// Compile-time interface checks
var (
	_ handlers.Handler = (*getHandler)(nil)
	_ handlers.Handler = (*postHandler)(nil)
)

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

type postHandler struct {
	*baseHandler
}

// requestParams holds common request parameters extracted from HTTP requests
type requestParams struct {
	Organization       string
	ProjectID          string
	BuildDefinitionID  string
	IdentityDescriptor string
	IdentityType       string
	IdentityName       string
	ProjectLevel       bool
	APIVersion         string
	AuthHeader         string
}

// extractRequestParamsGET extracts and validates parameters from GET requests (uses query params)
func (h *getHandler) extractRequestParamsGET(r *http.Request) (*requestParams, error) {
	params := &requestParams{
		Organization:       r.PathValue("organization"),
		ProjectID:          r.PathValue("projectId"),
		BuildDefinitionID:  r.URL.Query().Get("buildDefinitionId"),
		IdentityDescriptor: r.URL.Query().Get("identityDescriptor"),
		IdentityType:       r.URL.Query().Get("identityType"),
		IdentityName:       r.URL.Query().Get("identityName"),
		ProjectLevel:       r.URL.Query().Get("projectLevel") == "true", // to cast to bool from string, if malformed will be false
		APIVersion:         r.URL.Query().Get("api-version"),
		AuthHeader:         r.Header.Get("Authorization"),
	}

	if err := params.validate(); err != nil {
		return nil, err
	}

	return params, nil
}

// extractRequestParamsPOST extracts and validates parameters from POST requests (uses body)
func (h *postHandler) extractRequestParamsPOST(r *http.Request, permReq *BuildPermissionRequest) (*requestParams, error) {
	params := &requestParams{
		Organization:       r.PathValue("organization"),
		ProjectID:          r.PathValue("projectId"),
		BuildDefinitionID:  r.URL.Query().Get("buildDefinitionId"),
		IdentityDescriptor: permReq.Permissions.Identity.Descriptor,
		IdentityType:       permReq.Permissions.Identity.Type,
		IdentityName:       permReq.Permissions.Identity.Name,
		ProjectLevel:       r.URL.Query().Get("projectLevel") == "true", // to cast to bool from string, if malformed will be false
		APIVersion:         r.URL.Query().Get("api-version"),
		AuthHeader:         r.Header.Get("Authorization"),
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
	if p.ProjectID == "" {
		return fmt.Errorf("projectId parameter is required")
	}
	if !p.ProjectLevel && p.BuildDefinitionID == "" {
		return fmt.Errorf("buildDefinitionId parameter is required when projectLevel is false")
	}
	if p.APIVersion == "" {
		return fmt.Errorf("api-version is required")
	}
	if p.AuthHeader == "" {
		return fmt.Errorf("request rejected due to missing or invalid Basic authentication")
	}

	return p.validateIdentityParams()
}

// validateIdentityParams validates identity-related parameters
func (p *requestParams) validateIdentityParams() error {
	if p.IdentityDescriptor != "" {
		return nil // If descriptor is provided, we don't need other params
	}

	if p.IdentityType == "" {
		return fmt.Errorf("either 1) identityDescriptor or 2) identityType and identityName or 3) only identityType if identityType is 'build-service' must be provided as query parameters")
	}

	if p.IdentityType != string(BuildService) && p.IdentityName == "" {
		return fmt.Errorf("either 1) identityDescriptor or 2) identityType and identityName or 3) only identityType if identityType is 'build-service' must be provided as query parameters")
	}

	if p.IdentityType != string(BuildService) && p.IdentityType != string(AzureGroup) {
		return fmt.Errorf("the specified identityType is not valid or not supported, must be either 'build-service' or 'azure-group'")
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
func (h *baseHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(statusCode)
	w.Write(body)
}

// getIdentity retrieves identity information from Azure DevOps
func (h *baseHandler) getIdentity(ctx context.Context, organization, apiVersion, authHeader string, identityParams IdentityParams) (*IdentityResponse, error) {
	baseURL := fmt.Sprintf("https://vssps.dev.azure.com/%s/_apis/identities", organization)

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("api-version", apiVersion)
	q.Set("searchFilter", "General")
	q.Set("queryMembership", "None")

	filterValue, err := identityParams.buildFilterValue()
	if err != nil {
		return nil, err
	}
	q.Set("filterValue", filterValue)
	u.RawQuery = q.Encode()

	resp, err := h.makeAzureDevOpsRequest(ctx, http.MethodGet, u.String(), authHeader, nil)
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

	var identityResp IdentityResponse
	if err := json.Unmarshal(body, &identityResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal identity response: %w", err)
	}

	return &identityResp, nil
}

// extractDescriptor extracts the descriptor from identity response based on identity parameters
func (h *baseHandler) extractDescriptor(identityParams *IdentityParams, identityResponse *IdentityResponse) (string, error) {
	matchedIdentity, err := identityResponse.IdentityMatch(identityParams)
	if err != nil {
		return "", fmt.Errorf("failed to match identity: %w", err)
	}
	if matchedIdentity == nil {
		return "", fmt.Errorf("no matching identity found")
	}

	h.Log.Printf("Matched identity descriptor: %s", matchedIdentity.Descriptor)
	return matchedIdentity.Descriptor, nil
}

// getProjectName retrieves the project name from project ID
func (h *baseHandler) getProjectName(ctx context.Context, organization, projectID, apiVersion, authHeader string) (string, error) {
	baseURL := fmt.Sprintf("https://dev.azure.com/%s/_apis/projects/%s", organization, projectID)

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("api-version", apiVersion)
	u.RawQuery = q.Encode()

	resp, err := h.makeAzureDevOpsRequest(ctx, http.MethodGet, u.String(), authHeader, nil)
	if err != nil {
		return "", fmt.Errorf("failed to make GET request to azure devops: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("azure devops returned non-200 status: %d - %s", resp.StatusCode, string(body))
	}

	var projectData struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &projectData); err != nil {
		return "", fmt.Errorf("failed to unmarshal project response: %w", err)
	}

	return projectData.Name, nil
}

// resolveIdentityDescriptor resolves the identity descriptor either from the params or by fetching it
func (h *baseHandler) resolveIdentityDescriptor(ctx context.Context, params *requestParams, projectName string) (string, error) {
	if params.IdentityDescriptor != "" {
		return params.IdentityDescriptor, nil
	}

	identityParams := &IdentityParams{
		Type:        UserType(params.IdentityType),
		Name:        params.IdentityName,
		ProjectID:   params.ProjectID,
		ProjectName: projectName,
	}

	identityResp, err := h.getIdentity(ctx, params.Organization, params.APIVersion, params.AuthHeader, *identityParams)
	if err != nil {
		return "", fmt.Errorf("error getting identity: %w", err)
	}

	descriptor, err := h.extractDescriptor(identityParams, identityResp)
	if err != nil {
		return "", fmt.Errorf("error extracting identity descriptor: %w", err)
	}

	return descriptor, nil
}

// getBuildPermissions retrieves build permissions by making a POST request with empty permissions bitmask
func (h *getHandler) getBuildPermissions(ctx context.Context, organization, apiVersion, authHeader, projectID, buildDefinitionID, descriptor string, projectLevel bool) (*PermissionResponse, error) {
	baseURL := fmt.Sprintf("https://dev.azure.com/%s/_apis/accesscontrolentries/%s", organization, securityNamespaceID)

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("api-version", apiVersion)
	u.RawQuery = q.Encode()

	token := ""
	if projectLevel {
		token = projectID
	} else {
		token = CreateTokenObjectLevel(projectID, buildDefinitionID)
	}

	// Create AccessControlUpdate with empty allow/deny to fetch current permissions
	accessControlUpdate := AccessControlUpdate{
		Merge: true, // Merge 0 bitmasks with existing permissions instead of overwriting
		Token: token,
		AccessControlEntries: []AccessControlEntry{
			{
				Descriptor: descriptor,
				Allow:      0,
				Deny:       0,
			},
		},
	}

	bodyBytes, err := json.Marshal(accessControlUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access control update: %w", err)
	}

	resp, err := h.makeAzureDevOpsRequest(ctx, http.MethodPost, u.String(), authHeader, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request to azure devops: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("azure devops returned non-200 status: %d - %s", resp.StatusCode, string(respBody))
	}

	var permissionResp PermissionResponse
	if err := json.Unmarshal(respBody, &permissionResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal permission response: %w", err)
	}

	// Log retrieved permissions for debugging
	//for _, perm := range permissionResp.Value {
	//	h.Log.Printf("Descriptor: %s, Allow: %d, Deny: %d", perm.Descriptor, perm.Allow, perm.Deny)
	//}

	return &permissionResp, nil
}

// setBuildPermissions sets build permissions by making a POST request
func (h *postHandler) setBuildPermissions(ctx context.Context, organization, apiVersion, authHeader, projectID, buildDefinitionID, descriptor string, projectLevel bool, allow, deny int) (*PermissionResponse, error) {
	baseURL := fmt.Sprintf("https://dev.azure.com/%s/_apis/accesscontrolentries/%s", organization, securityNamespaceID)

	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("api-version", apiVersion)
	u.RawQuery = q.Encode()

	token := ""
	if projectLevel {
		token = projectID
	} else {
		token = CreateTokenObjectLevel(projectID, buildDefinitionID)
	}

	accessControlUpdate := AccessControlUpdate{
		Merge: false, // Don't merge, replace the permissions
		Token: token,
		AccessControlEntries: []AccessControlEntry{
			{
				Descriptor: descriptor,
				Allow:      allow,
				Deny:       deny,
			},
		},
	}

	bodyBytes, err := json.Marshal(accessControlUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal access control update: %w", err)
	}

	//h.Log.Printf("Setting permissions - Allow bitmask: %d, Deny bitmask: %d", allow, deny)

	resp, err := h.makeAzureDevOpsRequest(ctx, http.MethodPost, u.String(), authHeader, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request to azure devops: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("azure devops returned non-200 status: %d - %s", resp.StatusCode, string(respBody))
	}

	var permissionResp PermissionResponse
	if err := json.Unmarshal(respBody, &permissionResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal permission response: %w", err)
	}

	return &permissionResp, nil
}

// buildResponseFromPermissions constructs the final response object (BuildPermissionResponse) from permission data
func (h *baseHandler) buildResponseFromPermissions(params *requestParams, descriptor string, permissionResp *PermissionResponse) (*BuildPermissionResponse, error) {
	if len(permissionResp.Value) == 0 {
		return nil, fmt.Errorf("no permission data returned")
	}

	allowBits := permissionResp.Value[0].Allow
	denyBits := permissionResp.Value[0].Deny
	//h.Log.Printf("Retrieved permission bitmask - Allow: %d, Deny: %d", allowBits, denyBits)

	// Decode bitmask to PermissionFlags structs
	allowFlags := bitmaskToPermissionFlags(allowBits)
	denyFlags := bitmaskToPermissionFlags(denyBits)

	//h.Log.Printf("Decoded allowed permissions: %+v", allowFlags)
	//h.Log.Printf("Decoded denied permissions: %+v", denyFlags)

	//allow := filterTrueFlags(allowFlags)
	//deny := filterTrueFlags(denyFlags)
	allow := permissionFlagsToMap(allowFlags)
	deny := permissionFlagsToMap(denyFlags)

	//h.Log.Printf("Allowed permissions map: %+v", allow)
	//h.Log.Printf("Denied permissions map: %+v", deny)

	return &BuildPermissionResponse{
		Organization:      params.Organization,
		ProjectID:         params.ProjectID,
		BuildDefinitionID: params.BuildDefinitionID,
		ProjectLevel:      params.ProjectLevel,
		PermissionsInfo: PermissionsInfo{
			Identity: BuildPermissionIdentity{
				Descriptor: descriptor,
				Type:       params.IdentityType,
				Name:       params.IdentityName,
			},
			Allow: allow,
			Deny:  deny,
		},
	}, nil
}

// ServeHTTP handles GET requests for build permissions
// @Summary Get build permissions
// @Description Retrieve the permissions for a specific build definition in Azure DevOps
// @ID get-buildpermission
// @Tags Build Permissions
// @Param organization path string true "Organization name"
// @Param projectId path string true "Project ID"
// @Param buildDefinitionId query string false "Build Definition ID"
// @Param identityType query string false "Type of identity (e.g., 'azure-group', 'build-service')"
// @Param identityName query string false "Name of the identity (e.g., Contributors), not used if identityType is 'build-service' or if identityDescriptor is provided"
// @Param identityDescriptor query string false "Descriptor of the identity, has priority over identityType and identityName"
// @Param projectLevel query bool false "Whether to manage permissions at the project level (true) or build level (pipeline level) (false). Default is false." Default(false)
// @Param api-version query string true "API version (e.g., 7.2-preview.2)"
// @Param Authorization header string true "Basic Auth header"
// @Produce json
// @Success 200 {object} BuildPermissionResponse "Retrieved build permission"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /plugin/buildpermission/{organization}/{projectId} [get]
func (h *getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	// Extract and validate request parameters
	params, err := h.extractRequestParamsGET(r)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "request rejected due to missing or invalid Basic authentication" {
			statusCode = http.StatusUnauthorized
		}
		h.writeErrorResponse(w, statusCode, err.Error())
		return
	}

	// Get project name from project ID
	projectName, err := h.getProjectName(ctx, params.Organization, params.ProjectID, params.APIVersion, params.AuthHeader)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error getting project name: %s", err.Error()))
		return
	}

	// Resolve identity descriptor
	descriptor, err := h.resolveIdentityDescriptor(ctx, params, projectName)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get build permissions
	permissionResp, err := h.getBuildPermissions(ctx, params.Organization, params.APIVersion, params.AuthHeader, params.ProjectID, params.BuildDefinitionID, descriptor, params.ProjectLevel)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error getting build permissions: %s", err.Error()))
		return
	}

	// Build final response
	finalResponse, err := h.buildResponseFromPermissions(params, descriptor, permissionResp)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error building response: %s", err.Error()))
		return
	}

	finalRespBytes, err := json.Marshal(finalResponse)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error marshaling final response: %s", err.Error()))
		return
	}

	h.writeJSONResponse(w, http.StatusOK, finalRespBytes)
}

// ServeHTTP handles POST requests for setting build permissions
// @Summary Set build permissions
// @Description Set or update the permissions for a specific build definition in Azure DevOps. Identity information comes from the request body.
// @ID post-buildpermission
// @Tags Build Permissions
// @Accept json
// @Produce json
// @Param organization path string true "Organization name"
// @Param projectId path string true "Project ID"
// @Param buildDefinitionId query string false "Build Definition ID"
// @Param projectLevel query bool false "Whether to manage permissions at the project level (true) or build level (pipeline level) (false). Default is false." Default(false)
// @Param api-version query string true "API version (e.g., 7.2-preview.2)"
// @Param Authorization header string true "Basic Auth header"
// @Param body body BuildPermissionRequest true "Permission settings to apply. Identity info (type, name, descriptor) is specified in permissions.identity."
// @Success 200 {object} BuildPermissionResponse "Updated build permission"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /plugin/buildpermission/{organization}/{projectId} [post]
func (h *postHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	// Parse request body first
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to read request body: %s", err.Error()))
		return
	}
	defer r.Body.Close()

	var permReq BuildPermissionRequest
	if err := json.Unmarshal(body, &permReq); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to unmarshal request body: %s", err.Error()))
		return
	}

	// log the received permission request for debugging
	//h.Log.Printf("Received permission request: %+v", permReq)

	// Extract and validate request parameters (identity info comes from body)
	params, err := h.extractRequestParamsPOST(r, &permReq)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "request rejected due to missing or invalid Basic authentication" {
			statusCode = http.StatusUnauthorized
		}
		h.writeErrorResponse(w, statusCode, err.Error())
		return
	}

	// Get project name from project ID
	projectName, err := h.getProjectName(ctx, params.Organization, params.ProjectID, params.APIVersion, params.AuthHeader)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error getting project name: %s", err.Error()))
		return
	}

	// Resolve identity descriptor
	descriptor, err := h.resolveIdentityDescriptor(ctx, params, projectName)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert permission maps to bitmasks
	allowFlags := mapToPermissionFlags(permReq.Permissions.Allow)
	denyFlags := mapToPermissionFlags(permReq.Permissions.Deny)

	allowBitmask := permissionFlagsToBitmask(allowFlags)
	denyBitmask := permissionFlagsToBitmask(denyFlags)

	//h.Log.Printf("[POST] Setting permissions - Allow flags: %+v (bitmask: %d)", allowFlags, allowBitmask)
	//h.Log.Printf("[POST] Setting permissions - Deny flags: %+v (bitmask: %d)", denyFlags, denyBitmask)

	// Set build permissions
	permissionResp, err := h.setBuildPermissions(ctx, params.Organization, params.APIVersion, params.AuthHeader, params.ProjectID, params.BuildDefinitionID, descriptor, params.ProjectLevel, allowBitmask, denyBitmask)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error setting build permissions: %s", err.Error()))
		return
	}

	// Build final response
	finalResponse, err := h.buildResponseFromPermissions(params, descriptor, permissionResp)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error building response: %s", err.Error()))
		return
	}

	finalRespBytes, err := json.Marshal(finalResponse)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("error marshaling final response: %s", err.Error()))
		return
	}

	h.writeJSONResponse(w, http.StatusOK, finalRespBytes)
}
