package variablegroup

import "fmt"

// boolPtr returns a pointer to the given bool value.  Used during
// normalisation to fill in the boolean fields that ADO omits when false.
func boolPtr(b bool) *bool { return &b }

type VariableValue struct {
	Value      string `json:"value"`
	IsReadOnly *bool  `json:"isReadOnly"`
	IsSecret   *bool  `json:"isSecret"`
}

type ProjectReference struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type VariableGroupProjectReference struct {
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	ProjectReference *ProjectReference `json:"projectReference,omitempty"`
}

type IdentityRef struct {
	Descriptor  string `json:"descriptor,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	URL         string `json:"url,omitempty"`
	ID          string `json:"id,omitempty"`
	UniqueName  string `json:"uniqueName,omitempty"`
}

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
	// ProjectReference is lifted from VariableGroupProjectReferences[0] by the plugin.
	// It is not part of the ADO wire format; omitempty keeps it absent until populated.
	ProjectReference *ProjectReference      `json:"projectReference,omitempty"`
	ProviderData     map[string]interface{} `json:"providerData,omitempty"`
}

type VariableGroupListResponse struct {
	Count int             `json:"count"`
	Value []VariableGroup `json:"value"`
}

type VariableGroupParameters struct {
	Name                           string                          `json:"name"`
	Description                    string                          `json:"description,omitempty"`
	Type                           string                          `json:"type,omitempty"`
	Variables                      map[string]VariableValue        `json:"variables"`
	VariableGroupProjectReferences []VariableGroupProjectReference `json:"variableGroupProjectReferences,omitempty"`
	ProviderData                   map[string]interface{}          `json:"providerData,omitempty"`
}

type requestParams struct {
	Organization string
	Project      string
	GroupID      string
	APIVersion   string
	AuthHeader   string
}

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

func (p *requestParams) validate() error {
	if err := p.validateBase(); err != nil {
		return err
	}
	if p.GroupID == "" {
		return fmt.Errorf("groupId parameter is required")
	}
	return nil
}
