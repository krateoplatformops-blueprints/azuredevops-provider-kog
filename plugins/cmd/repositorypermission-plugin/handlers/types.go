package repositorypermission

import (
	"path"
	"reflect"
)

const (
	// securityNamespaceID is the GUID for Git Repositories security namespace in Azure DevOps
	securityNamespaceID = "2e9eb7ed-3c0a-47d4-87c1-0ffdd275fd87"
)

// IdentityPermission represents the permission settings for an identity
type IdentityPermission struct {
	Descriptor string `json:"descriptor"`
	Allow      int    `json:"allow"`
	Deny       int    `json:"deny"`
}

// PermissionResponse represents the response from Azure DevOps permission API
type PermissionResponse struct {
	Count int                  `json:"count"`
	Value []IdentityPermission `json:"value"`
}

// AccessControlEntry represents a single access control entry
type AccessControlEntry struct {
	Descriptor string `json:"descriptor"`
	Allow      int    `json:"allow"`
	Deny       int    `json:"deny"`
}

// AccessControlUpdate represents the update request for access control entries
type AccessControlUpdate struct {
	Merge                bool                 `json:"merge"`
	Token                string               `json:"token"`
	AccessControlEntries []AccessControlEntry `json:"accessControlEntries"`
}

// CreateTokenObjectLevel creates a security token for repository permissions
// The token format is: repoV2/{projectId}/{repositoryId}
func CreateTokenObjectLevel(projectID, repoID string) string {
	return path.Join("repoV2/", projectID, repoID)
}

// CreateTokenProjectLevel creates a security token for project-level permissions
// The token format is: repoV2/{projectId}
func CreateTokenProjectLevel(projectID string) string {
	return path.Join("repoV2/", projectID)
}

// PermissionBit represents individual permission bits in Azure DevOps
type PermissionBit int

// Permission bits for repository access control
// Each bit represents a specific permission that can be allowed or denied
const (
	Administer              PermissionBit = 1 << iota // 1
	GenericRead                                       // 2
	GenericContribute                                 // 4
	ForcePush                                         // 8
	CreateBranch                                      // 16
	CreateTag                                         // 32
	ManageNote                                        // 64
	PolicyExempt                                      // 128
	CreateRepository                                  // 256
	DeleteRepository                                  // 512
	RenameRepository                                  // 1024
	EditPolicies                                      // 2048
	RemoveOthersLocks                                 // 4096
	ManagePermissions                                 // 8192
	PullRequestContribute                             // 16384
	PullRequestBypassPolicy                           // 32768
	ViewAdvSecAlerts                                  // 65536
	DismissAdvSecAlerts                               // 131072
	ManageAdvSecScanning                              // 262144
)

// String returns the string representation of the permission bit
func (p PermissionBit) String() string {
	switch p {
	case Administer:
		return "Administer"
	case GenericRead:
		return "GenericRead"
	case GenericContribute:
		return "GenericContribute"
	case ForcePush:
		return "ForcePush"
	case CreateBranch:
		return "CreateBranch"
	case CreateTag:
		return "CreateTag"
	case ManageNote:
		return "ManageNote"
	case PolicyExempt:
		return "PolicyExempt"
	case CreateRepository:
		return "CreateRepository"
	case DeleteRepository:
		return "DeleteRepository"
	case RenameRepository:
		return "RenameRepository"
	case EditPolicies:
		return "EditPolicies"
	case RemoveOthersLocks:
		return "RemoveOthersLocks"
	case ManagePermissions:
		return "ManagePermissions"
	case PullRequestContribute:
		return "PullRequestContribute"
	case PullRequestBypassPolicy:
		return "PullRequestBypassPolicy"
	case ViewAdvSecAlerts:
		return "ViewAdvSecAlerts"
	case DismissAdvSecAlerts:
		return "DismissAdvSecAlerts"
	case ManageAdvSecScanning:
		return "ManageAdvSecScanning"
	default:
		return ""
	}
}

// comparePermissionBits checks if a specific permission bit is set in the bitmask
func comparePermissionBits(bitmask int, check int) bool {
	return (bitmask & check) == check
}

// bitmaskToPermissionFlags converts a permission bitmask to a PermissionFlags struct
func bitmaskToPermissionFlags(bitmask int) PermissionFlags {
	return PermissionFlags{
		Administer:              comparePermissionBits(bitmask, int(Administer)),
		GenericRead:             comparePermissionBits(bitmask, int(GenericRead)),
		GenericContribute:       comparePermissionBits(bitmask, int(GenericContribute)),
		ForcePush:               comparePermissionBits(bitmask, int(ForcePush)),
		CreateBranch:            comparePermissionBits(bitmask, int(CreateBranch)),
		CreateTag:               comparePermissionBits(bitmask, int(CreateTag)),
		ManageNote:              comparePermissionBits(bitmask, int(ManageNote)),
		PolicyExempt:            comparePermissionBits(bitmask, int(PolicyExempt)),
		CreateRepository:        comparePermissionBits(bitmask, int(CreateRepository)),
		DeleteRepository:        comparePermissionBits(bitmask, int(DeleteRepository)),
		RenameRepository:        comparePermissionBits(bitmask, int(RenameRepository)),
		EditPolicies:            comparePermissionBits(bitmask, int(EditPolicies)),
		RemoveOthersLocks:       comparePermissionBits(bitmask, int(RemoveOthersLocks)),
		ManagePermissions:       comparePermissionBits(bitmask, int(ManagePermissions)),
		PullRequestContribute:   comparePermissionBits(bitmask, int(PullRequestContribute)),
		PullRequestBypassPolicy: comparePermissionBits(bitmask, int(PullRequestBypassPolicy)),
		ViewAdvSecAlerts:        comparePermissionBits(bitmask, int(ViewAdvSecAlerts)),
		DismissAdvSecAlerts:     comparePermissionBits(bitmask, int(DismissAdvSecAlerts)),
		ManageAdvSecScanning:    comparePermissionBits(bitmask, int(ManageAdvSecScanning)),
	}
}

// permissionFlagsToBitmask converts a PermissionFlags struct to a bitmask
func permissionFlagsToBitmask(f PermissionFlags) int {
	bitmask := 0

	if f.Administer {
		bitmask |= int(Administer)
	}
	if f.GenericRead {
		bitmask |= int(GenericRead)
	}
	if f.GenericContribute {
		bitmask |= int(GenericContribute)
	}
	if f.ForcePush {
		bitmask |= int(ForcePush)
	}
	if f.CreateBranch {
		bitmask |= int(CreateBranch)
	}
	if f.CreateTag {
		bitmask |= int(CreateTag)
	}
	if f.ManageNote {
		bitmask |= int(ManageNote)
	}
	if f.PolicyExempt {
		bitmask |= int(PolicyExempt)
	}
	if f.CreateRepository {
		bitmask |= int(CreateRepository)
	}
	if f.DeleteRepository {
		bitmask |= int(DeleteRepository)
	}
	if f.RenameRepository {
		bitmask |= int(RenameRepository)
	}
	if f.EditPolicies {
		bitmask |= int(EditPolicies)
	}
	if f.RemoveOthersLocks {
		bitmask |= int(RemoveOthersLocks)
	}
	if f.ManagePermissions {
		bitmask |= int(ManagePermissions)
	}
	if f.PullRequestContribute {
		bitmask |= int(PullRequestContribute)
	}
	if f.PullRequestBypassPolicy {
		bitmask |= int(PullRequestBypassPolicy)
	}
	if f.ViewAdvSecAlerts {
		bitmask |= int(ViewAdvSecAlerts)
	}
	if f.DismissAdvSecAlerts {
		bitmask |= int(DismissAdvSecAlerts)
	}
	if f.ManageAdvSecScanning {
		bitmask |= int(ManageAdvSecScanning)
	}

	return bitmask
}

// filterTrueFlags returns a map of only the true permission flags
func filterTrueFlags(f PermissionFlags) map[string]bool {
	out := make(map[string]bool)
	v := reflect.ValueOf(f)
	t := reflect.TypeOf(f)

	for i := 0; i < v.NumField(); i++ {
		name := t.Field(i).Tag.Get("json")
		if name == "" {
			continue
		}
		if v.Field(i).Bool() {
			out[name] = true
		}
	}

	return out
}

func permissionFlagsToMap(f PermissionFlags) map[string]bool {
	out := make(map[string]bool)
	v := reflect.ValueOf(f)
	t := reflect.TypeOf(f)

	for i := 0; i < v.NumField(); i++ {
		name := t.Field(i).Tag.Get("json")
		if name == "" {
			continue
		}
		out[name] = v.Field(i).Bool()
	}

	return out
}

// mapToPermissionFlags converts a map of permission names to a PermissionFlags struct
func mapToPermissionFlags(perms map[string]bool) PermissionFlags {
	flags := PermissionFlags{}

	for key, enabled := range perms {
		if !enabled {
			continue
		}

		// Normalize the key to lowercase for case-insensitive matching
		//key = strings.ToLower(key)

		switch key {
		case "Administer":
			flags.Administer = true
		case "GenericRead":
			flags.GenericRead = true
		case "GenericContribute":
			flags.GenericContribute = true
		case "ForcePush":
			flags.ForcePush = true
		case "CreateBranch":
			flags.CreateBranch = true
		case "CreateTag":
			flags.CreateTag = true
		case "ManageNote":
			flags.ManageNote = true
		case "PolicyExempt":
			flags.PolicyExempt = true
		case "CreateRepository":
			flags.CreateRepository = true
		case "DeleteRepository":
			flags.DeleteRepository = true
		case "RenameRepository":
			flags.RenameRepository = true
		case "EditPolicies":
			flags.EditPolicies = true
		case "RemoveOthersLocks":
			flags.RemoveOthersLocks = true
		case "ManagePermissions":
			flags.ManagePermissions = true
		case "PullRequestContribute":
			flags.PullRequestContribute = true
		case "PullRequestBypassPolicy":
			flags.PullRequestBypassPolicy = true
		case "ViewAdvSecAlerts":
			flags.ViewAdvSecAlerts = true
		case "DismissAdvSecAlerts":
			flags.DismissAdvSecAlerts = true
		case "ManageAdvSecScanning":
			flags.ManageAdvSecScanning = true
		}
	}

	return flags
}

// RepositoryPermissionResponse represents the complete response for repository permissions
type RepositoryPermissionResponse struct {
	Organization    string          `json:"organization" yaml:"organization"`
	ProjectID       string          `json:"projectId" yaml:"projectId"`
	RepositoryID    string          `json:"repositoryId" yaml:"repositoryId"`
	ProjectLevel    bool            `json:"projectLevel" yaml:"projectLevel"`
	PermissionsInfo PermissionsInfo `json:"permissions" yaml:"permissions"`
}

// PermissionsInfo contains the detailed permission information
type PermissionsInfo struct {
	Identity RepositoryPermissionIdentity `json:"identity" yaml:"identity"`
	Allow    map[string]bool              `json:"allow" yaml:"allow"`
	Deny     map[string]bool              `json:"deny" yaml:"deny"`
}

// RepositoryPermissionIdentity represents the identity associated with permissions
type RepositoryPermissionIdentity struct {
	Type       string `json:"type" yaml:"type"`
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Descriptor string `json:"descriptor,omitempty" yaml:"descriptor,omitempty"`
}

// PermissionFlags represents all available repository permissions as boolean flags
type PermissionFlags struct {
	Administer              bool `json:"Administer"`
	GenericRead             bool `json:"GenericRead"`
	GenericContribute       bool `json:"GenericContribute"`
	ForcePush               bool `json:"ForcePush"`
	CreateBranch            bool `json:"CreateBranch"`
	CreateTag               bool `json:"CreateTag"`
	ManageNote              bool `json:"ManageNote"`
	PolicyExempt            bool `json:"PolicyExempt"`
	CreateRepository        bool `json:"CreateRepository"`
	DeleteRepository        bool `json:"DeleteRepository"`
	RenameRepository        bool `json:"RenameRepository"`
	EditPolicies            bool `json:"EditPolicies"`
	RemoveOthersLocks       bool `json:"RemoveOthersLocks"`
	ManagePermissions       bool `json:"ManagePermissions"`
	PullRequestContribute   bool `json:"PullRequestContribute"`
	PullRequestBypassPolicy bool `json:"PullRequestBypassPolicy"`
	ViewAdvSecAlerts        bool `json:"ViewAdvSecAlerts"`
	DismissAdvSecAlerts     bool `json:"DismissAdvSecAlerts"`
	ManageAdvSecScanning    bool `json:"ManageAdvSecScanning"`
}

// RepositoryPermissionRequest represents the request body for setting repository permissions
type RepositoryPermissionRequest struct {
	Permissions PermissionsInfo `json:"permissions" yaml:"permissions"`
}
