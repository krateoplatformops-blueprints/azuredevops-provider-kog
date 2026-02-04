package variablegroup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/krateoplatformops/azuredevops-provider-kog/pkg/handlers"
)

const (
	requestTimeout  = 30 * time.Second
	contentTypeJSON = "application/json"
)

// --- Constructors ---

// GetVariableGroup returns a handler for GET /plugin/{organization}/{project}/variablegroups/{groupId}.
// Proxies to the ADO individual GET endpoint and patches the response:
//   - normalises boolean fields (isReadOnly / isSecret) that ADO omits when false
//   - fills secret variable values from the CR spec body (ADO always returns null for secrets)
func GetVariableGroup(opts handlers.HandlerOptions) handlers.Handler {
	return &getHandler{baseHandler: newBaseHandler(opts)}
}

// FindByVariableGroup returns a handler for GET /plugin/{organization}/{project}/variablegroups.
// Proxies to the ADO list endpoint and, for the item whose name matches the CR spec:
//   - normalises boolean fields
//   - fills secret variable values from the CR spec body
//   - populates variableGroupProjectReferences (null in list responses) from the CR spec body
func FindByVariableGroup(opts handlers.HandlerOptions) handlers.Handler {
	return &findByHandler{baseHandler: newBaseHandler(opts)}
}

// DeleteVariableGroup returns a handler for DELETE /plugin/{organization}/{project}/variablegroups/{groupId}.
// The ADO delete endpoint is organisation-scoped and requires projectIds as a query param;
// this handler extracts the project from the path and forwards it correctly.
func DeleteVariableGroup(opts handlers.HandlerOptions) handlers.Handler {
	return &deleteHandler{baseHandler: newBaseHandler(opts)}
}

// Compile-time interface checks
var (
	_ handlers.Handler = (*getHandler)(nil)
	_ handlers.Handler = (*findByHandler)(nil)
	_ handlers.Handler = (*deleteHandler)(nil)
)

// --- Base handler ---

type baseHandler struct {
	handlers.HandlerOptions
}

func newBaseHandler(opts handlers.HandlerOptions) *baseHandler {
	return &baseHandler{HandlerOptions: opts}
}

// makeAzureDevOpsRequest performs an HTTP request to Azure DevOps API.
func (h *baseHandler) makeAzureDevOpsRequest(ctx context.Context, method, reqURL, authHeader string, body []byte) (*http.Response, error) {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
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

func (h *baseHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	h.Log.Print(message)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func (h *baseHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(statusCode)
	w.Write(body)
}

// extractParams populates a requestParams from the incoming request.
// GroupID will be empty for routes that do not include {groupId} in the
// path; callers choose validate() or validateBase() accordingly.
func (h *baseHandler) extractParams(r *http.Request) *requestParams {
	return &requestParams{
		Organization: r.PathValue("organization"),
		Project:      r.PathValue("project"),
		GroupID:      r.PathValue("groupId"),
		APIVersion:   r.URL.Query().Get("api-version"),
		AuthHeader:   r.Header.Get("Authorization"),
	}
}

// --- GET handler ---

type getHandler struct {
	*baseHandler
}

// @Summary Get a variable group (plugin)
// @Description Retrieves a variable group by ID from Azure DevOps, then patches secret
// variable values and normalises boolean fields using values from the CR spec body.
// @Tags VariableGroup
// @Accept json
// @Produce json
// @Param organization path string true "Organization name"
// @Param project path string true "Project ID (project name is not supported in this field)"
// @Param groupId path string true "Variable group ID"
// @Param api-version query string true "API version"
// @Param body body VariableGroupParameters true "CR spec forwarded by RDC; used to restore secret variable values"
// @Success 200 {object} VariableGroup "Variable group with normalised booleans and patched secrets"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /plugin/{organization}/{project}/variablegroups/{groupId} [post]
func (h *getHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	params := h.extractParams(r)
	if err := params.validate(); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid request parameters: %v", err))
		return
	}

	// CR spec body sent by RDC – used to patch secret values
	crSpec, err := parseCRSpec(r)
	if err != nil {
		h.Log.Printf("GET: no CR spec in body, secret patching will be skipped: %v", err)
	}
	h.Log.Printf("GET: parsed CR spec: %+v", crSpec)

	// Call ADO individual GET endpoint.
	// This endpoint correctly populates variableGroupProjectReferences.
	adoURL := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/distributedtask/variablegroups/%s",
		params.Organization, params.Project, params.GroupID)
	u, err := url.Parse(adoURL)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse ADO URL: %v", err))
		return
	}
	q := u.Query()
	q.Set("api-version", params.APIVersion)
	u.RawQuery = q.Encode()

	h.Log.Printf("GET: calling ADO API: %s", u.String())

	resp, err := h.makeAzureDevOpsRequest(ctx, http.MethodGet, u.String(), params.AuthHeader, nil)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to call ADO API: %v", err))
		return
	}
	defer resp.Body.Close()

	adoBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to read ADO response: %v", err))
		return
	}

	h.Log.Printf("================== GET: ADO response received =========================")
	h.Log.Printf("GET: ADO response status: %d", resp.StatusCode)
	h.Log.Printf("GET: ADO response body: %s", string(adoBody))
	h.Log.Printf("=======================================================================")

	// Handle "null" body with 200 status as 404 Not Found
	if isNullResponseBody(resp.StatusCode, adoBody) {
		h.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("variable group %s not found", params.GroupID))
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		h.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("variable group %s not found", params.GroupID))
		return
	}
	if resp.StatusCode != http.StatusOK {
		h.writeErrorResponse(w, resp.StatusCode, fmt.Sprintf("ADO returned status %d: %s", resp.StatusCode, string(adoBody)))
		return
	}

	var vg VariableGroup
	if err := json.Unmarshal(adoBody, &vg); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse ADO response: %v", err))
		return
	}

	// Normalise booleans and patch secret values from CR spec
	normalizeVariableGroup(&vg, crSpec)

	// Extract project reference to top-level field for RDC
	extractProjectReference(&vg)

	result, err := json.Marshal(vg)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal response: %v", err))
		return
	}

	h.Log.Printf("GET: returning normalised response for variable group id=%d", vg.ID)
	h.writeJSONResponse(w, http.StatusOK, result)
}

// --- FINDBY handler ---

type findByHandler struct {
	*baseHandler
}

// @Summary List variable groups (plugin)
// @Description Retrieves all variable groups for a project from Azure DevOps.
// For the item matching the CR name: patches secret values, normalises booleans,
// and populates variableGroupProjectReferences from the CR spec body.
// Non-matching items are only normalised (booleans).
// @Tags VariableGroup
// @Accept json
// @Produce json
// @Param organization path string true "Organization name"
// @Param project path string true "Project ID (project name is not supported in this field)"
// @Param api-version query string true "API version"
// @Param body body VariableGroupParameters true "CR spec forwarded by RDC; used to patch the matching item"
// @Success 200 {object} VariableGroupListResponse "List response with normalised items"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /plugin/{organization}/{project}/variablegroups [post]
func (h *findByHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	params := h.extractParams(r)
	if err := params.validateBase(); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid request parameters: %v", err))
		return
	}

	// CR spec body sent by RDC – used to patch the matching item
	crSpec, err := parseCRSpec(r)
	if err != nil {
		h.Log.Printf("FINDBY: no CR spec in body, secret/ref patching will be skipped: %v", err)
	}
	h.Log.Printf("FINDBY: parsed CR spec: %+v", crSpec)

	// Call ADO list endpoint
	adoURL := fmt.Sprintf("https://dev.azure.com/%s/%s/_apis/distributedtask/variablegroups",
		params.Organization, params.Project)
	u, err := url.Parse(adoURL)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse ADO URL: %v", err))
		return
	}
	q := u.Query()
	q.Set("api-version", params.APIVersion)
	u.RawQuery = q.Encode()

	h.Log.Printf("FINDBY: calling ADO list API: %s", u.String())

	resp, err := h.makeAzureDevOpsRequest(ctx, http.MethodGet, u.String(), params.AuthHeader, nil)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to call ADO API: %v", err))
		return
	}
	defer resp.Body.Close()

	adoBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to read ADO response: %v", err))
		return
	}

	if resp.StatusCode != http.StatusOK {
		h.writeErrorResponse(w, resp.StatusCode, fmt.Sprintf("ADO returned status %d: %s", resp.StatusCode, string(adoBody)))
		return
	}

	var listResp VariableGroupListResponse
	if err := json.Unmarshal(adoBody, &listResp); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to parse ADO list response: %v", err))
		return
	}

	// Walk items: full patch on the matching item, boolean-only normalisation on others
	var crName string
	if crSpec != nil {
		crName = crSpec.Name
	}

	for i := range listResp.Value {
		if listResp.Value[i].Name == crName && crSpec != nil {
			normalizeVariableGroup(&listResp.Value[i], crSpec)
			patchProjectReferences(&listResp.Value[i], crSpec)
		} else {
			normalizeVariableGroup(&listResp.Value[i], nil)
		}
		extractProjectReference(&listResp.Value[i])
	}

	result, err := json.Marshal(listResp)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to marshal response: %v", err))
		return
	}

	h.Log.Printf("FINDBY: returning normalised list (matching name=%q)", crName)
	h.writeJSONResponse(w, http.StatusOK, result)
}

// --- DELETE handler ---

type deleteHandler struct {
	*baseHandler
}

// @Summary Delete a variable group (plugin)
// @Description Deletes a variable group from Azure DevOps by calling the organisation-level
// delete endpoint with projectIds extracted from the URL path.
// @Tags VariableGroup
// @Produce json
// @Param organization path string true "Organization name"
// @Param project path string true "Project ID (project name is not supported in this field)"
// @Param groupId path string true "Variable group ID"
// @Param api-version query string true "API version"
// @Success 204 "No Content – variable group deleted"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /plugin/{organization}/{project}/variablegroups/{groupId} [delete]
func (h *deleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	params := h.extractParams(r)
	if err := params.validate(); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid request parameters: %v", err))
		return
	}

	h.Log.Printf("Deleting variable group: organization=%s, groupId=%s, projectId=%s",
		params.Organization, params.GroupID, params.Project)

	if err := h.deleteVariableGroupAndRespond(ctx, w, params); err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting variable group: %v", err))
	}
}

func (h *deleteHandler) deleteVariableGroupAndRespond(ctx context.Context, w http.ResponseWriter, params *requestParams) error {
	// ADO delete is organisation-scoped; project is passed as projectIds query param
	baseURL := fmt.Sprintf("https://dev.azure.com/%s/_apis/distributedtask/variablegroups/%s",
		params.Organization, params.GroupID)

	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("api-version", params.APIVersion)
	q.Set("projectIds", params.Project) // set just the one project ID we have; ADO expects a comma-separated list
	u.RawQuery = q.Encode()

	h.Log.Printf("Calling Azure DevOps DELETE API: %s", u.String())

	resp, err := h.makeAzureDevOpsRequest(ctx, http.MethodDelete, u.String(), params.AuthHeader, nil)
	if err != nil {
		return fmt.Errorf("failed to make DELETE request to azure devops: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	h.Log.Printf("================== FINBY: ADO response received =========================")
	h.Log.Printf("FINDBY: ADO response status: %d", resp.StatusCode)
	h.Log.Printf("FINDBY: ADO response body: %s", string(body))
	h.Log.Printf("=======================================================================")

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		h.Log.Printf("Successfully deleted variable group: groupId=%s", params.GroupID)
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	if resp.StatusCode == http.StatusNotFound {
		h.Log.Printf("Variable group not found: groupId=%s", params.GroupID)
		h.writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Variable group with ID %s not found", params.GroupID))
		return nil
	}

	h.Log.Printf("Azure DevOps returned non-success status: %d - %s", resp.StatusCode, string(body))
	return fmt.Errorf("azure devops returned non-success status: %d - %s", resp.StatusCode, string(body))
}

// --- Shared helpers ---

// parseCRSpec reads and unmarshals the request body as the CR spec
// forwarded by RDC.  Returns nil, error when the body is absent or
// unparseable – callers treat this as non-fatal and skip patching.
func parseCRSpec(r *http.Request) (*VariableGroupParameters, error) {
	if r.Body == nil {
		return nil, fmt.Errorf("no request body")
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	if len(bodyBytes) == 0 {
		return nil, fmt.Errorf("empty request body")
	}

	var spec VariableGroupParameters
	if err := json.Unmarshal(bodyBytes, &spec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return &spec, nil
}

// normalizeVariableGroup patches a single VariableGroup in place:
//   - Sets IsReadOnly and IsSecret to false when ADO has omitted them.
//     ADO omits these fields entirely when the value is false, which
//     causes a spurious diff against the CR spec.
//   - When crSpec is non-nil, copies the value of every secret variable
//     from crSpec.  ADO never returns secret values – even with
//     loadSecrets=true – so the CR spec is the only source of truth.
func normalizeVariableGroup(vg *VariableGroup, crSpec *VariableGroupParameters) {
	if vg == nil {
		return
	}

	for name, v := range vg.Variables {
		if v.IsReadOnly == nil {
			v.IsReadOnly = boolPtr(false)
			log.Println("Added isReadOnly=false for variable:", name)
		}
		if v.IsSecret == nil {
			v.IsSecret = boolPtr(false)
			log.Println("Added isSecret=false for variable:", name)
		}

		// Patch secret values from CR spec
		if *v.IsSecret && crSpec != nil {
			if specVar, ok := crSpec.Variables[name]; ok {
				v.Value = specVar.Value
				log.Printf("Patched secret value for variable: %s", name)
			}
		}

		vg.Variables[name] = v // map values are not addressable; write back
	}
}

// patchProjectReferences populates VariableGroupProjectReferences from
// the CR spec when the ADO list endpoint has left it nil.  The
// individual GET endpoint returns the field correctly, so this patch is
// only needed for findby responses.
func patchProjectReferences(vg *VariableGroup, crSpec *VariableGroupParameters) {
	if vg == nil || crSpec == nil {
		return
	}

	if vg.VariableGroupProjectReferences == nil {
		vg.VariableGroupProjectReferences = crSpec.VariableGroupProjectReferences
		log.Println("Patched variableGroupProjectReferences from CR spec")
	}
}

// extractProjectReference lifts projectReference out of the
// variableGroupProjectReferences array into a top-level field.
// The array is assumed to have at most one entry (a variable group
// belongs to one project in practice); if it has more, the first entry
// is used and a warning is logged.
func extractProjectReference(vg *VariableGroup) {
	if vg == nil {
		return
	}
	if len(vg.VariableGroupProjectReferences) == 0 {
		log.Println("extractProjectReference: variableGroupProjectReferences is empty, projectReference will be absent")
		return
	}
	if len(vg.VariableGroupProjectReferences) > 1 {
		log.Printf("extractProjectReference: WARNING variableGroupProjectReferences has %d entries, using first one", len(vg.VariableGroupProjectReferences))
	}
	vg.ProjectReference = vg.VariableGroupProjectReferences[0].ProjectReference
	if vg.ProjectReference != nil {
		log.Printf("extractProjectReference: set projectReference id=%s name=%s", vg.ProjectReference.ID, vg.ProjectReference.Name)
	} else {
		log.Println("extractProjectReference: variableGroupProjectReferences[0].projectReference is nil")
	}
}

// helper to check if 200 and null response body, this is in reality a 404
func isNullResponseBody(statusCode int, body []byte) bool {
	return statusCode == http.StatusOK && bytes.Equal(body, []byte("null"))
}
