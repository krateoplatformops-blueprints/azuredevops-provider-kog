package repositorypermission

import (
	"path"
	"reflect"
	"strings"
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

// CreateToken creates a security token for repository permissions
// The token format is: repoV2/{projectId}/{repositoryId}
func CreateToken(projectID, repoID string) string {
	return path.Join("repoV2/", projectID, repoID)
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
		return "administerpermission"
	case GenericRead:
		return "genericread"
	case GenericContribute:
		return "genericcontribute"
	case ForcePush:
		return "forcepush"
	case CreateBranch:
		return "createbranch"
	case CreateTag:
		return "createtag"
	case ManageNote:
		return "managenote"
	case PolicyExempt:
		return "policyexempt"
	case CreateRepository:
		return "createrepository"
	case DeleteRepository:
		return "deleterepository"
	case RenameRepository:
		return "renamerepository"
	case EditPolicies:
		return "editpolicies"
	case RemoveOthersLocks:
		return "removeotherslocks"
	case ManagePermissions:
		return "managepermissions"
	case PullRequestContribute:
		return "pullrequestcontribute"
	case PullRequestBypassPolicy:
		return "pullrequestbypasspolicy"
	case ViewAdvSecAlerts:
		return "viewadvsecalerts"
	case DismissAdvSecAlerts:
		return "dismissadvsecalerts"
	case ManageAdvSecScanning:
		return "manageadvsecscanning"
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
		AdministerPermission:    comparePermissionBits(bitmask, int(Administer)),
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

	if f.AdministerPermission {
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
		key = strings.ToLower(key)

		switch key {
		case "administerpermission":
			flags.AdministerPermission = true
		case "genericread":
			flags.GenericRead = true
		case "genericcontribute":
			flags.GenericContribute = true
		case "forcepush":
			flags.ForcePush = true
		case "createbranch":
			flags.CreateBranch = true
		case "createtag":
			flags.CreateTag = true
		case "managenote":
			flags.ManageNote = true
		case "policyexempt":
			flags.PolicyExempt = true
		case "createrepository":
			flags.CreateRepository = true
		case "deleterepository":
			flags.DeleteRepository = true
		case "renamerepository":
			flags.RenameRepository = true
		case "editpolicies":
			flags.EditPolicies = true
		case "removeotherslocks":
			flags.RemoveOthersLocks = true
		case "managepermissions":
			flags.ManagePermissions = true
		case "pullrequestcontribute":
			flags.PullRequestContribute = true
		case "pullrequestbypasspolicy":
			flags.PullRequestBypassPolicy = true
		case "viewadvsecalerts":
			flags.ViewAdvSecAlerts = true
		case "dismissadvsecalerts":
			flags.DismissAdvSecAlerts = true
		case "manageadvsecscanning":
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
	AdministerPermission    bool `json:"administerpermission"`
	GenericRead             bool `json:"genericread"`
	GenericContribute       bool `json:"genericcontribute"`
	ForcePush               bool `json:"forcepush"`
	CreateBranch            bool `json:"createbranch"`
	CreateTag               bool `json:"createtag"`
	ManageNote              bool `json:"managenote"`
	PolicyExempt            bool `json:"policyexempt"`
	CreateRepository        bool `json:"createrepository"`
	DeleteRepository        bool `json:"deleterepository"`
	RenameRepository        bool `json:"renamerepository"`
	EditPolicies            bool `json:"editpolicies"`
	RemoveOthersLocks       bool `json:"removeotherslocks"`
	ManagePermissions       bool `json:"managepermissions"`
	PullRequestContribute   bool `json:"pullrequestcontribute"`
	PullRequestBypassPolicy bool `json:"pullrequestbypasspolicy"`
	ViewAdvSecAlerts        bool `json:"viewadvsecalerts"`
	DismissAdvSecAlerts     bool `json:"dismissadvsecalerts"`
	ManageAdvSecScanning    bool `json:"manageadvsecscanning"`
}

// RepositoryPermissionReq represents the request body for setting repository permissions
type RepositoryPermissionRequest struct {
	Permissions PermissionsInfo `json:"permissions" yaml:"permissions"`
}
