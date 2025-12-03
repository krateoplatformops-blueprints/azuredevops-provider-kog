package repositorypermission

import (
	"fmt"
	"strings"
)

// IdentityResponse represents the response from Azure DevOps identities API
type IdentityResponse struct {
	Count int        `json:"count"`
	Value []Identity `json:"value"`
}

// Identity represents a single identity in Azure DevOps
type Identity struct {
	ID                  string   `json:"id"`
	Descriptor          string   `json:"descriptor"`
	SubjectDescriptor   string   `json:"subjectDescriptor"`
	ProviderDisplayName string   `json:"providerDisplayName"`
	CustomDisplayName   string   `json:"customDisplayName"`
	IsActive            bool     `json:"isActive"`
	Members             []string `json:"members"`
	MemberOf            []string `json:"memberOf"`
	MemberIds           []string `json:"memberIds"`
	ResourceVersion     int      `json:"resourceVersion"`
	MetaTypeID          int      `json:"metaTypeId"`
}

// UserType represents the type of identity in Azure DevOps
type UserType string

const (
	// BuildService represents the build service identity type
	BuildService UserType = "build-service"
	// AzureGroup represents the Azure DevOps group identity type
	AzureGroup UserType = "azure-group"
)

// IsValid checks if the UserType is valid
func (u UserType) IsValid() bool {
	switch u {
	case BuildService, AzureGroup:
		return true
	default:
		return false
	}
}

// ResolveIdentityDescriptorFromUserType returns the descriptor prefix for the given UserType
func (u UserType) ResolveIdentityDescriptorFromUserType() (string, error) {
	switch u {
	case BuildService:
		return "Microsoft.TeamFoundation.ServiceIdentity", nil
	case AzureGroup:
		return "Microsoft.TeamFoundation.Identity", nil
	default:
		return "", fmt.Errorf("the specified usertype is not valid: %s", u)
	}
}

// IdentityParams contains parameters needed to identify an identity
type IdentityParams struct {
	Type        UserType
	Name        string // Name is ignored if Type is BuildService
	ProjectID   string
	ProjectName string
}

// buildFilterValue constructs the filterValue query parameter for the identity API
func (i *IdentityParams) buildFilterValue() (string, error) {
	switch i.Type {
	case BuildService:
		return i.ProjectID, nil
	case AzureGroup:
		// Format: [ProjectName]\GroupName
		return fmt.Sprintf("[%s]\\%s", i.ProjectName, i.Name), nil
	default:
		return "", fmt.Errorf("unsupported identity type: %s", i.Type)
	}
}

// Validate checks if the IdentityParams are valid
func (i *IdentityParams) Validate() error {
	if !i.Type.IsValid() {
		return fmt.Errorf("invalid identity type: %s", i.Type)
	}

	if i.ProjectID == "" {
		return fmt.Errorf("project ID is required")
	}

	if i.Type == AzureGroup {
		if i.Name == "" {
			return fmt.Errorf("identity name is required for azure-group type")
		}
		if i.ProjectName == "" {
			return fmt.Errorf("project name is required for azure-group type")
		}
	}

	return nil
}

// IdentityMatch finds the matching identity from the response based on the identity parameters
func (resp *IdentityResponse) IdentityMatch(identityParams *IdentityParams) (*Identity, error) {
	if err := identityParams.Validate(); err != nil {
		return nil, fmt.Errorf("invalid identity parameters: %w", err)
	}

	resolvedDescriptor, err := identityParams.Type.ResolveIdentityDescriptorFromUserType()
	if err != nil {
		return nil, err
	}

	for i := range resp.Value {
		identity := &resp.Value[i]

		// Check descriptor contains the resolved descriptor
		if !strings.Contains(identity.Descriptor, resolvedDescriptor) {
			continue
		}

		// Match based on identity type
		switch identityParams.Type {
		case BuildService:
			// For build service, match by project ID
			if identity.ProviderDisplayName == identityParams.ProjectID {
				return identity, nil
			}
		case AzureGroup:
			// For Azure group, match by formatted name: [ProjectName]\GroupName
			expectedName := fmt.Sprintf("[%s]\\%s", identityParams.ProjectName, identityParams.Name)
			if identity.ProviderDisplayName == expectedName {
				return identity, nil
			}
		}
	}

	return nil, fmt.Errorf("identity not found matching criteria from identity parameters: type=%s, projectId=%s, name=%s",
		identityParams.Type, identityParams.ProjectID, identityParams.Name)
}
