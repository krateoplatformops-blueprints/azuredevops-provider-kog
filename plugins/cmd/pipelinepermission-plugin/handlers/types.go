package pipelinepermission

import (
	"time"
)

// ResourcePipelinePermissions represents the main response structure for pipeline permissions
type ResourcePipelinePermissions struct {
	AllPipelines *Permission          `json:"allPipelines,omitempty"`
	Pipelines    []PipelinePermission `json:"pipelines,omitempty"`
	Resource     *Resource            `json:"resource,omitempty"`
}

// Permission represents basic permission information
type Permission struct {
	Authorized   bool         `json:"authorized"`
	AuthorizedBy *IdentityRef `json:"authorizedBy,omitempty"`
	AuthorizedOn *time.Time   `json:"authorizedOn,omitempty"`
}

// PipelinePermission extends Permission with an ID field
type PipelinePermission struct {
	Permission
	ID int32 `json:"id"`
}

// Resource represents a resource with basic identification
type Resource struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

// IdentityRef represents an identity reference extending GraphSubjectBase
type IdentityRef struct {
	GraphSubjectBase
	DirectoryAlias    string `json:"directoryAlias,omitempty"`
	ID                string `json:"id,omitempty"`
	ImageURL          string `json:"imageUrl,omitempty"`
	Inactive          bool   `json:"inactive"`
	IsAadIdentity     bool   `json:"isAadIdentity"`
	IsContainer       bool   `json:"isContainer"`
	IsDeletedInOrigin bool   `json:"isDeletedInOrigin"`
	ProfileURL        string `json:"profileUrl,omitempty"`
	UniqueName        string `json:"uniqueName,omitempty"`
}

// GraphSubjectBase represents the base structure for graph subjects
type GraphSubjectBase struct {
	Links       *ReferenceLinks `json:"_links,omitempty"`
	Descriptor  string          `json:"descriptor,omitempty"`
	DisplayName string          `json:"displayName,omitempty"`
	URL         string          `json:"url,omitempty"`
}

// ReferenceLinks represents a collection of REST reference links
type ReferenceLinks struct {
	Links map[string]interface{} `json:"links,omitempty"`
}
