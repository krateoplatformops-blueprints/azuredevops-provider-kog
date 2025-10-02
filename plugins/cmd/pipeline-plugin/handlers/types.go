package pipeline

// Pipeline represents the response from:
// GET /{organization}/{project}/_apis/pipelines/{id}
type Pipeline struct {
	Links         *ReferenceLinks        `json:"_links,omitempty"`
	Configuration *PipelineConfiguration `json:"configuration,omitempty"`
	URL           string                 `json:"url,omitempty"`
	// Embedded fields from PipelineBase
	Folder   string `json:"folder,omitempty"`
	ID       int32  `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Revision int32  `json:"revision,omitempty"`
}

// PipelineConfiguration represents the configuration of a pipeline
type PipelineConfiguration struct {
	Type       string           `json:"type,omitempty"` // enum: unknown, yaml, designerJson, justInTime, designerHyphenJson
	Path       string           `json:"path,omitempty"`
	Repository *BuildRepository `json:"repository,omitempty"` // Required
}

// ReferenceLinks represents a collection of REST reference links
type ReferenceLinks struct {
	Links map[string]interface{} `json:"links,omitempty"`
}

// ---

// CreatePipelineParametersComplete represents the request body for:
// POST /{organization}/{project}/_apis/pipelines
type CreatePipelineParameters struct {
	Configuration *PipelineConfigurationParameters `json:"configuration"` // Required
	Folder        string                           `json:"folder,omitempty"`
	Name          string                           `json:"name"` // Required
}

// CreatePipelineConfigurationParametersComplete represents the configuration for creating a pipeline
type PipelineConfigurationParameters struct {
	Type       string           `json:"type"` // Required - enum: unknown, yaml, designerJson, justInTime, designerHyphenJson
	Path       string           `json:"path,omitempty"`
	Repository *BuildRepository `json:"repository"` // Required
}

// BuildRepository represents repository information for the pipeline
type BuildRepository struct {
	ID   string `json:"id"`   // Required
	Type string `json:"type"` // Required - enum: unknown, gitHub, azureReposGit, azureReposGitHyphenated
}

// BuildDefinitionMinimal represents the request made by the plugin to Azure DevOps API when updating a pipeline
// The plugin will construct this object starting from the request body (UpdatePipelineParameters) coming from RDC
type BuildDefinitionMinimal struct {
	Name       string           `json:"name,omitempty"`
	Path       string           `json:"path,omitempty"`
	Repository *BuildRepository `json:"repository,omitempty"`
	Type       string           `json:"type,omitempty"`
	Revision   int32            `json:"revision,omitempty"`
	ID         int32            `json:"id,omitempty"`
	Process    *Process         `json:"process,omitempty"`
}

type Process struct {
	YAMLFilename string `json:"yamlFilename,omitempty"` // Required if type is yaml
}

type UpdatePipelineParameters struct {
	Configuration *PipelineConfigurationParameters `json:"configuration"`
	Folder        string                           `json:"folder,omitempty"`
	Name          string                           `json:"name"`
	ID            int32                            `json:"id"` // maybe to be removed since RDC does not include it
	Revision      int32                            `json:"revision"`
}

// ConfigurationType enum values
const (
	ConfigurationTypeUnknown            = "unknown"
	ConfigurationTypeYAML               = "yaml"
	ConfigurationTypeDesignerJSON       = "designerJson"
	ConfigurationTypeJustInTime         = "justInTime"
	ConfigurationTypeDesignerHyphenJSON = "designerHyphenJson"
)

// RepositoryType enum values
const (
	RepositoryTypeUnknown                 = "unknown"
	RepositoryTypeGitHub                  = "gitHub"
	RepositoryTypeAzureReposGit           = "azureReposGit"
	RepositoryTypeAzureReposGitHyphenated = "azureReposGitHyphenated"
)

// GetPipelineResponse represents the response for getting a single pipeline
type GetPipelineResponse Pipeline

// CreatePipelineRequest represents the request for creating a pipeline
//type CreatePipelineRequest CreatePipelineParametersComplete

// CreatePipelineResponse represents the response for creating a pipeline
//type CreatePipelineResponse Pipeline

// UpdatePipelineResponse represents the request for updating a pipeline
type UpdatePipelineRequest UpdatePipelineParameters

// UpdatePipelineResponse represents the response for updating a pipeline
type UpdatePipelineResponse Pipeline
