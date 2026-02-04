package variablegroup

import "fmt"

// boolPtr returns a pointer to the given bool value.  Used during
// normalisation to fill in the boolean fields that ADO omits when false.
func boolPtr(b bool) *bool { return &b }

// -------------------------------------------------------------------
// Azure DevOps API wire types – derived from the VariableGroups OAS
// (TaskAgent 7.2-preview)
// -------------------------------------------------------------------

// VariableValue represents a single variable entry inside a variable
// group.  ADO omits isReadOnly / isSecret when their value is false;
// the pointer types let us detect that during normalisation and still
// marshal the field as explicit false (not null / absent).
type VariableValue struct {
	Value      string `json:"value"`
	IsReadOnly *bool  `json:"isReadOnly"`
	IsSecret   *bool  `json:"isSecret"`
}

// ProjectReference identifies an Azure DevOps project by ID and/or name.
type ProjectReference struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// VariableGroupProjectReference associates a variable group with a
// single project.  Matches the VariableGroupProjectReference schema in
// the OAS.
type VariableGroupProjectReference struct {
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	ProjectReference *ProjectReference `json:"projectReference,omitempty"`
}

// IdentityRef is a minimal representation of the Azure DevOps identity
// that appears in the read-only createdBy / modifiedBy response fields.
// Only the most commonly consumed fields are retained; deprecated ADO
// fields (imageUrl, profileUrl, directoryAlias …) are intentionally
// omitted.
type IdentityRef struct {
	Descriptor  string `json:"descriptor,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	URL         string `json:"url,omitempty"`
	ID          string `json:"id,omitempty"`
	UniqueName  string `json:"uniqueName,omitempty"`
}

// VariableGroup is the full response object returned by the Azure
// DevOps VariableGroups API for a single variable group.  Collection
// fields (Variables, VariableGroupProjectReferences) deliberately omit
// the omitempty tag so that a nil / null value round-trips faithfully
// instead of being silently dropped.
type VariableGroup struct {
	ID                             int                             `json:"id"`
	Name                           string                          `json:"name"`
	Description                    string                          `json:"description,omitempty"`
	Type                           string                          `json:"type,omitempty"`
	IsShared                       *bool                           `json:"isShared,omitempty"`
	CreatedBy                      *IdentityRef                    `json:"createdBy,omitempty"`
	CreatedOn                      string                          `json:"createdOn,omitempty"`
	ModifiedBy                     *IdentityRef                    `json:"modifiedBy,omitempty"`
	ModifiedOn                     string                          `json:"modifiedOn,omitempty"`
	Variables                      map[string]VariableValue        `json:"variables"`
	VariableGroupProjectReferences []VariableGroupProjectReference `json:"variableGroupProjectReferences"`
	ProviderData                   map[string]interface{}                 `json:"providerData,omitempty"`
}

// VariableGroupListResponse is the {count, value} envelope returned by
// the Azure DevOps list endpoint for variable groups.
type VariableGroupListResponse struct {
	Count int             `json:"count"`
	Value []VariableGroup `json:"value"`
}

// VariableGroupParameters mirrors the VariableGroupParameters OAS
// schema – the request body used for create / update.  It is also the
// shape of the CR spec body that RDC forwards to the plugin on GET and
// FINDBY requests.
type VariableGroupParameters struct {
	Name                           string                          `json:"name"`
	Description                    string                          `json:"description,omitempty"`
	Type                           string                          `json:"type,omitempty"`
	Variables                      map[string]VariableValue        `json:"variables"`
	VariableGroupProjectReferences []VariableGroupProjectReference `json:"variableGroupProjectReferences,omitempty"`
	ProviderData                   map[string]interface{}                 `json:"providerData,omitempty"`
}

// -------------------------------------------------------------------
// Request extraction
// -------------------------------------------------------------------

// requestParams holds the path / query / header values extracted from
// an incoming plugin request.
type requestParams struct {
	Organization string
	Project      string
	GroupID      string
	APIVersion   string
	AuthHeader   string
}

// validateBase checks the fields common to every route
// (organisation, project, api-version, authorisation).
func (p *requestParams) validateBase() error {
	if p.Organization == "" {
		return fmt.Errorf("organization parameter is required")
	}
	if p.Project == "" {
		return fmt.Errorf("project parameter is required")
	}
	if p.APIVersion == "" {
		return fmt.Errorf("api-version is required")
	}
	if p.AuthHeader == "" {
		return fmt.Errorf("request rejected due to missing or invalid authentication")
	}
	return nil
}

// validate checks all fields including GroupID.  Use for routes that
// include {groupId} in the path (get-by-ID, delete).
func (p *requestParams) validate() error {
	if err := p.validateBase(); err != nil {
		return err
	}
	if p.GroupID == "" {
		return fmt.Errorf("groupId parameter is required")
	}
	return nil
}
