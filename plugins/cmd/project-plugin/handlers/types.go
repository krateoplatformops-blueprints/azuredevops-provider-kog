package project

import "fmt"

// OperationReference represents the Azure DevOps operation reference response
// returned by create and update operations
type OperationReference struct {
	ID       string `json:"id,omitempty"`
	PluginID string `json:"pluginId,omitempty"`
	Status   string `json:"status,omitempty"`
	URL      string `json:"url,omitempty"`
}

// TransformedOperationReference represents the transformed operation reference
// with fields prefixed to avoid confusion with project fields
type TransformedOperationReference struct {
	OperationID       string `json:"operationId,omitempty"`
	OperationPluginID string `json:"operationPluginId,omitempty"`
	OperationStatus   string `json:"operationStatus,omitempty"`
	OperationURL      string `json:"operationUrl,omitempty"`
}

// =======================================================================================================

// ProjectLinks represents the _links object in Azure DevOps project response
type ProjectLinks struct {
	Self       *Link `json:"self,omitempty"`
	Collection *Link `json:"collection,omitempty"`
	Web        *Link `json:"web,omitempty"`
}

// Link represents a hyperlink reference
type Link struct {
	Href string `json:"href,omitempty"`
}

// DefaultTeam represents the default team information in a project
type DefaultTeam struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// ProjectState represents the possible states of a project
type ProjectState string

const (
	ProjectStateDeleting      ProjectState = "deleting"
	ProjectStateNew           ProjectState = "new"
	ProjectStateWellFormed    ProjectState = "wellFormed"
	ProjectStateCreatePending ProjectState = "createPending"
	ProjectStateAll           ProjectState = "all"
	ProjectStateUnchanged     ProjectState = "unchanged"
	ProjectStateDeleted       ProjectState = "deleted"
)

// ProjectVisibility represents the visibility level of a project
type ProjectVisibility string

const (
	ProjectVisibilityPrivate ProjectVisibility = "private"
	ProjectVisibilityPublic  ProjectVisibility = "public"
)

type Capabilities struct {
	VersionControl  VersionControl  `json:"versioncontrol,omitempty"`
	ProcessTemplate ProcessTemplate `json:"processTemplate,omitempty"`
}

type VersionControl struct {
	SourceControlType string `json:"sourceControlType,omitempty"`
}

type ProcessTemplate struct {
	TemplateTypeID string `json:"templateTypeId,omitempty"`
}

// ProjectResponseBody represents an Azure DevOps project
// @Description Azure DevOps project object in response body
type ProjectResponseBody struct {
	Abbreviation        string        `json:"abbreviation,omitempty"`
	DefaultTeamImageUrl string        `json:"defaultTeamImageUrl,omitempty"`
	Description         string        `json:"description,omitempty"`
	ID                  string        `json:"id,omitempty"`
	LastUpdateTime      string        `json:"lastUpdateTime,omitempty"`
	Name                string        `json:"name,omitempty"`
	Revision            int           `json:"revision,omitempty"`
	State               string        `json:"state,omitempty" enums:"deleting,new,wellFormed,createPending,all,unchanged,deleted"`
	URL                 string        `json:"url,omitempty"`
	Visibility          string        `json:"visibility,omitempty" enums:"private,public"`
	Links               *ProjectLinks `json:"_links,omitempty"`
	Capabilities        *Capabilities `json:"capabilities,omitempty"`
	DefaultTeam         *DefaultTeam  `json:"defaultTeam,omitempty"`
}

// Project represents an Azure DevOps project
// @Description Azure DevOps project object in request body for creation
type ProjectRequestBodyCreate struct {
	Description  string        `json:"description,omitempty"`
	Name         string        `json:"name,omitempty"`
	Visibility   string        `json:"visibility,omitempty" enums:"private,public"`
	Capabilities *Capabilities `json:"capabilities,omitempty"`
}

// ValidateVisibility validates the project visibility against allowed values
func (p *ProjectRequestBodyCreate) ValidateVisibility() error {
	if p.Visibility == "" {
		return nil // Visibility is optional
	}

	validVisibilities := map[string]bool{
		string(ProjectVisibilityPrivate): true,
		string(ProjectVisibilityPublic):  true,
	}

	if !validVisibilities[p.Visibility] {
		return fmt.Errorf("invalid visibility '%s': must be one of [private, public]", p.Visibility)
	}
	return nil
}

// Validate performs comprehensive validation of the ProjectRequestBody
func (p *ProjectRequestBodyCreate) Validate() error {
	// In future, add more validation methods here if needed

	if p.Name == "" {
		return fmt.Errorf("project name is required")
	}

	if err := p.ValidateVisibility(); err != nil {
		return err
	}
	return nil
}

// Project represents an Azure DevOps project
// @Description Azure DevOps project object in request body for update
type ProjectRequestBodyUpdate struct {
	Abbreviation string `json:"abbreviation,omitempty"` // Only updatable field, not allowed in create
	Description  string `json:"description,omitempty"`
}

// This is needed to merge both create and update fields since oasgen will create the spec of the CRD starting from the request body
// of the create operation, so we need to include the updatable fields from the update operation as well.
// and then add logic in the handler to separate them accordingly.
type ProjectRequestBodyMerged struct {
	Abbreviation string        `json:"abbreviation,omitempty"` // Only updatable field, not allowed in create
	Description  string        `json:"description,omitempty"`
	Name         string        `json:"name,omitempty"`
	Visibility   string        `json:"visibility,omitempty" enums:"private,public"`
	Capabilities *Capabilities `json:"capabilities,omitempty"`
}
