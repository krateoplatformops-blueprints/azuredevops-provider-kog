package pullrequest

import "time"

// This file is generated from the OpenAPI specification.

// GitPullRequest defines the structure for a pull request.
type GitPullRequest struct {
	Links                 *ReferenceLinks                  `json:"_links,omitempty"`
	ArtifactID            string                           `json:"artifactId,omitempty"`
	AutoCompleteSetBy     *IdentityRef                     `json:"autoCompleteSetBy,omitempty"`
	ClosedBy              *IdentityRef                     `json:"closedBy,omitempty"`
	ClosedDate            time.Time                        `json:"closedDate,omitempty"`
	CodeReviewID          int32                            `json:"codeReviewId,omitempty"`
	Commits               []GitCommitRef                   `json:"commits,omitempty"`
	CompletionOptions     *GitPullRequestCompletionOptions `json:"completionOptions,omitempty"`
	CompletionQueueTime   time.Time                        `json:"completionQueueTime,omitempty"`
	CreatedBy             *IdentityRef                     `json:"createdBy,omitempty"`
	CreationDate          time.Time                        `json:"creationDate,omitempty"`
	Description           string                           `json:"description,omitempty"`
	ForkSource            *GitForkRef                      `json:"forkSource,omitempty"`
	IsDraft               bool                             `json:"isDraft,omitempty"`
	Labels                []WebApiTagDefinition            `json:"labels,omitempty"`
	LastMergeCommit       *GitCommitRef                    `json:"lastMergeCommit,omitempty"`
	LastMergeSourceCommit *GitCommitRef                    `json:"lastMergeSourceCommit,omitempty"`
	LastMergeTargetCommit *GitCommitRef                    `json:"lastMergeTargetCommit,omitempty"`
	MergeFailureMessage   string                           `json:"mergeFailureMessage,omitempty"`
	MergeFailureType      string                           `json:"mergeFailureType,omitempty"`
	MergeID               string                           `json:"mergeId,omitempty"`
	MergeOptions          *GitPullRequestMergeOptions      `json:"mergeOptions,omitempty"`
	MergeStatus           string                           `json:"mergeStatus,omitempty"`
	PullRequestID         int32                            `json:"pullRequestId,omitempty"`
	RemoteURL             string                           `json:"remoteUrl,omitempty"`
	Repository            *GitRepository                   `json:"repository,omitempty"`
	Reviewers             []IdentityRefWithVote            `json:"reviewers,omitempty"`
	SourceRefName         string                           `json:"sourceRefName,omitempty"`
	Status                string                           `json:"status,omitempty"`
	SupportsIterations    bool                             `json:"supportsIterations,omitempty"`
	TargetRefName         string                           `json:"targetRefName,omitempty"`
	Title                 string                           `json:"title,omitempty"`
	URL                   string                           `json:"url,omitempty"`
	WorkItemRefs          []ResourceRef                    `json:"workItemRefs,omitempty"`
}

// GitPullRequestCompletionOptions defines completion options for a pull request.
type GitPullRequestCompletionOptions struct {
	AutoCompleteIgnoreConfigIds []int32 `json:"autoCompleteIgnoreConfigIds,omitempty"`
	BypassPolicy                bool    `json:"bypassPolicy,omitempty"`
	BypassReason                string  `json:"bypassReason,omitempty"`
	DeleteSourceBranch          bool    `json:"deleteSourceBranch,omitempty"`
	MergeCommitMessage          string  `json:"mergeCommitMessage,omitempty"`
	MergeStrategy               string  `json:"mergeStrategy,omitempty"`
	SquashMerge                 bool    `json:"squashMerge,omitempty"`
	TransitionWorkItems         bool    `json:"transitionWorkItems,omitempty"`
	TriggeredByAutoComplete     bool    `json:"triggeredByAutoComplete,omitempty"`
}

// GitPullRequestMergeOptions defines merge options for a pull request.
type GitPullRequestMergeOptions struct {
	ConflictAuthorshipCommits  bool `json:"conflictAuthorshipCommits,omitempty"`
	DetectRenameFalsePositives bool `json:"detectRenameFalsePositives,omitempty"`
	DisableRenames             bool `json:"disableRenames,omitempty"`
}

// IdentityRef defines a reference to an identity.
type IdentityRef struct {
	GraphSubjectBase
	DirectoryAlias    string `json:"directoryAlias,omitempty"`
	Id                string `json:"id,omitempty"`
	ImageURL          string `json:"imageUrl,omitempty"`
	Inactive          bool   `json:"inactive,omitempty"`
	IsAadIdentity     bool   `json:"isAadIdentity,omitempty"`
	IsContainer       bool   `json:"isContainer,omitempty"`
	IsDeletedInOrigin bool   `json:"isDeletedInOrigin,omitempty"`
	ProfileURL        string `json:"profileUrl,omitempty"`
	UniqueName        string `json:"uniqueName,omitempty"`
}

// IdentityRefWithVote includes a vote on a pull request.
type IdentityRefWithVote struct {
	IdentityRef
	HasDeclined bool   `json:"hasDeclined,omitempty"`
	IsFlagged   bool   `json:"isFlagged,omitempty"`
	IsReapprove bool   `json:"isReapprove,omitempty"`
	IsRequired  bool   `json:"isRequired,omitempty"`
	ReviewerURL string `json:"reviewerUrl,omitempty"`
	Vote        int16  `json:"vote,omitempty"`
}

// GraphSubjectBase is a base type for identity references.
type GraphSubjectBase struct {
	Links       *ReferenceLinks `json:"_links,omitempty"`
	Descriptor  string          `json:"descriptor,omitempty"`
	DisplayName string          `json:"displayName,omitempty"`
	URL         string          `json:"url,omitempty"`
}

// ReferenceLinks represents a collection of REST reference links.
type ReferenceLinks struct {
	Links map[string]interface{} `json:"links,omitempty"`
}

// GitCommitRef provides properties that describe a Git commit.
type GitCommitRef struct {
	Links            *ReferenceLinks  `json:"_links,omitempty"`
	Author           *GitUserDate     `json:"author,omitempty"`
	ChangeCounts     map[string]int32 `json:"changeCounts,omitempty"`
	Changes          []GitChange      `json:"changes,omitempty"`
	Comment          string           `json:"comment,omitempty"`
	CommentTruncated bool             `json:"commentTruncated,omitempty"`
	CommitID         string           `json:"commitId,omitempty"`
	Committer        *GitUserDate     `json:"committer,omitempty"`
	Parents          []string         `json:"parents,omitempty"`
	RemoteURL        string           `json:"remoteUrl,omitempty"`
	Statuses         []GitStatus      `json:"statuses,omitempty"`
	URL              string           `json:"url,omitempty"`
	WorkItems        []ResourceRef    `json:"workItems,omitempty"`
}

// GitUserDate contains user info and date for Git operations.
type GitUserDate struct {
	Date     time.Time `json:"date,omitempty"`
	Email    string    `json:"email,omitempty"`
	ImageURL string    `json:"imageUrl,omitempty"`
	Name     string    `json:"name,omitempty"`
}

// GitChange represents a change in a commit.
type GitChange struct {
	Change
	ChangeID           int32        `json:"changeId,omitempty"`
	NewContentTemplate *GitTemplate `json:"newContentTemplate,omitempty"`
	OriginalPath       string       `json:"originalPath,omitempty"`
}

// Change represents a version control change.
type Change struct {
	ChangeType       string       `json:"changeType,omitempty"`
	Item             string       `json:"item,omitempty"`
	NewContent       *ItemContent `json:"newContent,omitempty"`
	SourceServerItem string       `json:"sourceServerItem,omitempty"`
	URL              string       `json:"url,omitempty"`
}

// ItemContent represents the content of an item.
type ItemContent struct {
	Content     string `json:"content,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}

// GitTemplate defines a Git template.
type GitTemplate struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// GitStatus represents a Git status.
type GitStatus struct {
	Links        *ReferenceLinks   `json:"_links,omitempty"`
	Context      *GitStatusContext `json:"context,omitempty"`
	CreatedBy    *IdentityRef      `json:"createdBy,omitempty"`
	CreationDate time.Time         `json:"creationDate,omitempty"`
	Description  string            `json:"description,omitempty"`
	ID           int32             `json:"id,omitempty"`
	State        string            `json:"state,omitempty"`
	TargetURL    string            `json:"targetUrl,omitempty"`
	UpdatedDate  time.Time         `json:"updatedDate,omitempty"`
}

// GitStatusContext uniquely identifies a status.
type GitStatusContext struct {
	Genre string `json:"genre,omitempty"`
	Name  string `json:"name,omitempty"`
}

// ResourceRef is a reference to a resource.
type ResourceRef struct {
	ID  string `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}

// GitForkRef contains information about a fork ref.
type GitForkRef struct {
	GitRef
	Repository *GitRepository `json:"repository,omitempty"`
}

// GitRef represents a Git reference.
type GitRef struct {
	Links          *ReferenceLinks `json:"_links,omitempty"`
	Creator        *IdentityRef    `json:"creator,omitempty"`
	IsLocked       bool            `json:"isLocked,omitempty"`
	IsLockedBy     *IdentityRef    `json:"isLockedBy,omitempty"`
	Name           string          `json:"name,omitempty"`
	ObjectID       string          `json:"objectId,omitempty"`
	PeeledObjectID string          `json:"peeledObjectId,omitempty"`
	Statuses       []GitStatus     `json:"statuses,omitempty"`
	URL            string          `json:"url,omitempty"`
}

// WebApiTagDefinition represents a tag definition.
type WebApiTagDefinition struct {
	Active bool   `json:"active,omitempty"`
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	URL    string `json:"url,omitempty"`
}

// GitRepository represents a Git repository.
type GitRepository struct {
	Links            *ReferenceLinks       `json:"_links,omitempty"`
	DefaultBranch    string                `json:"defaultBranch,omitempty"`
	ID               string                `json:"id,omitempty"`
	IsFork           bool                  `json:"isFork,omitempty"`
	Name             string                `json:"name,omitempty"`
	ParentRepository *GitRepositoryRef     `json:"parentRepository,omitempty"`
	Project          *TeamProjectReference `json:"project,omitempty"`
	RemoteURL        string                `json:"remoteUrl,omitempty"`
	Size             int64                 `json:"size,omitempty"`
	SshURL           string                `json:"sshUrl,omitempty"`
	URL              string                `json:"url,omitempty"`
	WebURL           string                `json:"webUrl,omitempty"`
}

// GitRepositoryRef is a shallow reference to a Git repository.
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

// TeamProjectReference is a shallow reference to a TeamProject.
type TeamProjectReference struct {
	Abbreviation        string    `json:"abbreviation,omitempty"`
	DefaultTeamImageURL string    `json:"defaultTeamImageUrl,omitempty"`
	Description         string    `json:"description,omitempty"`
	ID                  string    `json:"id,omitempty"`
	LastUpdateTime      time.Time `json:"lastUpdateTime,omitempty"`
	Name                string    `json:"name,omitempty"`
	Revision            int64     `json:"revision,omitempty"`
	State               string    `json:"state,omitempty"`
	URL                 string    `json:"url,omitempty"`
	Visibility          string    `json:"visibility,omitempty"`
}

// TeamProjectCollectionReference is a reference to a TeamProjectCollection.
type TeamProjectCollectionReference struct {
	AvatarURL string `json:"avatarUrl,omitempty"`
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
}

// CreatePullRequestReq defines the request body for creating a pull request.
//type CreatePullRequestReq struct {
//	CompletionOptions  *GitPullRequestCompletionOptions `json:"completionOptions,omitempty"`
//	Description        string                           `json:"description,omitempty"`
//	IsDraft            bool                             `json:"isDraft,omitempty"`
//	MergeOptions       *GitPullRequestMergeOptions      `json:"mergeOptions,omitempty"`
//	Reviewers          []IdentityRef                    `json:"reviewers,omitempty"`
//	SourceRefName      string                           `json:"sourceRefName,omitempty"`
//	Status             string                           `json:"status,omitempty"`
//	SupportsIterations bool                             `json:"supportsIterations,omitempty"`
//	TargetRefName      string                           `json:"targetRefName,omitempty"`
//	Title              string                           `json:"title,omitempty"`
//}

type AutoCompleteSetByOnlyId struct {
	Id string `json:"id,omitempty"`
}

// UpdatePullRequestReq defines the request body for updating a pull request.
type UpdatePullRequestReq struct {
	Status                string                           `json:"status,omitempty"`
	Title                 string                           `json:"title,omitempty"`
	Description           string                           `json:"description,omitempty"`
	CompletionOptions     *GitPullRequestCompletionOptions `json:"completionOptions,omitempty"`
	MergeOptions          *GitPullRequestMergeOptions      `json:"mergeOptions,omitempty"`
	AutoCompleteSetBy     *AutoCompleteSetByOnlyId         `json:"autoCompleteSetBy,omitempty"`
	TargetRefName         string                           `json:"targetRefName,omitempty"`
	LastMergeSourceCommit *GitCommitRef                    `json:"lastMergeSourceCommit,omitempty"` // Needed to update status to "completed"
}
