package gitrepository

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Custom time type to handle Azure DevOps timestamp format
type AzureDevOpsTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for AzureDevOpsTime
func (t *AzureDevOpsTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from the JSON string
	timeStr := strings.Trim(string(data), "\"")

	// List of possible time formats used by Azure DevOps
	formats := []string{
		"2006-01-02T15:04:05.999999999", // With nanoseconds
		"2006-01-02T15:04:05.999999",    // With microseconds
		"2006-01-02T15:04:05.999",       // With milliseconds
		"2006-01-02T15:04:05",           // Without fractional seconds
		time.RFC3339,                    // Standard RFC3339
		time.RFC3339Nano,                // RFC3339 with nanoseconds
	}

	for _, format := range formats {
		if parsedTime, err := time.Parse(format, timeStr); err == nil {
			t.Time = parsedTime
			return nil
		}
	}

	return fmt.Errorf("unable to parse time: %s", timeStr)
}

// MarshalJSON implements json.Marshaler for AzureDevOpsTime
func (t AzureDevOpsTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Format(time.RFC3339))
}

// ProjectState represents the state of a project
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

// ProjectVisibility represents the visibility of a project
type ProjectVisibility string

const (
	ProjectVisibilityPrivate ProjectVisibility = "private"
	ProjectVisibilityPublic  ProjectVisibility = "public"
)

// ReferenceLinks represents a collection of REST reference links
type ReferenceLinks struct {
	Links map[string]interface{} `json:"links,omitempty"`
}

// TeamProjectCollectionReference represents a reference to a TeamProjectCollection
type TeamProjectCollectionReference struct {
	AvatarURL string `json:"avatarUrl,omitempty"`
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
}

// TeamProjectReference represents a shallow reference to a TeamProject
type TeamProjectReference struct {
	Abbreviation        string            `json:"abbreviation,omitempty"`
	DefaultTeamImageURL string            `json:"defaultTeamImageUrl,omitempty"`
	Description         string            `json:"description,omitempty"`
	ID                  string            `json:"id,omitempty"`
	LastUpdateTime      *AzureDevOpsTime  `json:"lastUpdateTime,omitempty"`
	Name                string            `json:"name,omitempty"`
	Revision            int64             `json:"revision,omitempty"`
	State               ProjectState      `json:"state,omitempty"`
	URL                 string            `json:"url,omitempty"`
	Visibility          ProjectVisibility `json:"visibility,omitempty"`
}

// TeamProjectReferenceMinimal represents a minimal reference to a TeamProject
type TeamProjectReferenceMinimal struct {
	ID string `json:"id,omitempty"`
}

// GitRepositoryRef represents a reference to a Git repository
type GitRepositoryRef struct {
	Collection *TeamProjectCollectionReference `json:"collection,omitempty"`
	ID         string                          `json:"id,omitempty"`
	IsFork     bool                            `json:"isFork,omitempty"`
	Name       string                          `json:"name,omitempty"`
	Project    *TeamProjectReference           `json:"project,omitempty"`
	RemoteURL  string                          `json:"remoteUrl,omitempty"`
	SshURL     string                          `json:"sshUrl,omitempty"`
	URL        string                          `json:"url,omitempty"`
}

// GitRepositoryRefMinimal represents a minimal reference to a Git repository
type GitRepositoryRefMinimal struct {
	ID      string                       `json:"id,omitempty"`
	Project *TeamProjectReferenceMinimal `json:"project,omitempty"`
}

// GitRepository represents a Git repository (response body)
type GitRepository struct {
	Links            *ReferenceLinks       `json:"_links,omitempty"`
	CreationDate     *AzureDevOpsTime      `json:"creationDate,omitempty"`
	DefaultBranch    string                `json:"defaultBranch,omitempty"`
	ID               string                `json:"id,omitempty"`
	IsDisabled       bool                  `json:"isDisabled,omitempty"`
	IsFork           bool                  `json:"isFork,omitempty"`
	IsInMaintenance  bool                  `json:"isInMaintenance,omitempty"`
	Name             string                `json:"name,omitempty"`
	ParentRepository *GitRepositoryRef     `json:"parentRepository,omitempty"`
	Project          *TeamProjectReference `json:"project,omitempty"`
	RemoteURL        string                `json:"remoteUrl,omitempty"`
	Size             int64                 `json:"size,omitempty"`
	SshURL           string                `json:"sshUrl,omitempty"`
	URL              string                `json:"url,omitempty"`
	ValidRemoteURLs  []string              `json:"validRemoteUrls,omitempty"`
	WebURL           string                `json:"webUrl,omitempty"`
}

// GitRepositoryCreateOptionsMinimal represents the minimal options for creating a Git repository (request body for POST)
// Additional field from UPDATE added and handled with business logic in this plugin: DefaultBranch
type GitRepositoryCreateOptionsMinimalPlugin struct {
	Name             string                       `json:"name,omitempty"`
	ParentRepository *GitRepositoryRefMinimal     `json:"parentRepository,omitempty"`
	Project          *TeamProjectReferenceMinimal `json:"project,omitempty"`
	DefaultBranch    string                       `json:"defaultBranch,omitempty"`
	Initialize       bool                         `json:"initialize,omitempty"` // Indicates if the repository should be initialized with an initial commit
}

type GitRepositoryCreateOptionsMinimal struct {
	Name             string                       `json:"name,omitempty"`
	ParentRepository *GitRepositoryRefMinimal     `json:"parentRepository,omitempty"`
	Project          *TeamProjectReferenceMinimal `json:"project,omitempty"`
	DefaultBranch    string                       `json:"defaultBranch,omitempty"`
}

// GitRepositoryCreateOptions represents the full options for creating a Git repository
// Currently not used, instead we use GitRepositoryCreateOptionsMinimal
//type GitRepositoryCreateOptions struct {
//	Name             string                `json:"name,omitempty"`
//	ParentRepository *GitRepositoryRef     `json:"parentRepository,omitempty"`
//	Project          *TeamProjectReference `json:"project,omitempty"`
//}

// GitRepositoryUpdateOptions represents the options for updating a Git repository (request body for PATCH)
type GitRepositoryUpdateOptions struct {
	Name          string `json:"name,omitempty"`
	DefaultBranch string `json:"defaultBranch,omitempty"`
}

// Request/Response wrapper types for API operations

// ListRepositoriesResponse represents the response for listing repositories
//type ListRepositoriesResponse []GitRepository

// GetRepositoryResponse represents the response for getting a single repository
type GetRepositoryResponse GitRepository

// CreateRepositoryRequest represents the request for creating a repository
type CreateRepositoryRequest GitRepositoryCreateOptionsMinimalPlugin

type CreateRepositoryResponse GitRepository

// UpdateRepositoryRequest represents the request for updating a repository
type UpdateRepositoryRequest GitRepositoryUpdateOptions

// UpdateRepositoryResponse represents the response for updating a repository
type UpdateRepositoryResponse GitRepository
