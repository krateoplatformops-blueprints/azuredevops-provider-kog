package graphgroup

// GraphGroupsResponse represents the enhanced response with filtering metadata
type GraphGroupsResponse struct {
	OriginalCount       int          `json:"originalCount"`       // Original count from Azure DevOps API
	Count               int          `json:"count"`               // Count after filtering (equals originalCount if no filtering)
	NumberOfFilteredOut int          `json:"numberOfFilteredOut"` // Number of groups filtered out (0 if no filtering)
	Value               []GraphGroup `json:"value"`               // Array of graph groups
}

// GraphGroup represents a single graph group from Azure DevOps
// Based on the actual API response structure from GraphMember schema
type GraphGroup struct {
	Links           map[string]interface{} `json:"_links,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Descriptor      string                 `json:"descriptor"`
	DisplayName     string                 `json:"displayName"`
	Domain          string                 `json:"domain"`
	IsCrossProject  *bool                  `json:"isCrossProject,omitempty"`
	MailAddress     *string                `json:"mailAddress"`
	Origin          string                 `json:"origin"`
	OriginId        string                 `json:"originId"`
	PrincipalName   string                 `json:"principalName,omitempty"`
	SubjectKind     string                 `json:"subjectKind"`
	URL             string                 `json:"url,omitempty"`
}

// AzureDevOpsGraphGroupsResponse represents the original response from Azure DevOps API
// This is used internally to unmarshal the response before filtering
type AzureDevOpsGraphGroupsResponse struct {
	Count int          `json:"count"`
	Value []GraphGroup `json:"value"`
}
