package buildpermission

import (
	"reflect"
)

const (
	// securityNamespaceID is the GUID for Build security namespace in Azure DevOps
	securityNamespaceID = "33344d9c-fc72-4d6f-aba5-fa317101a7e9"
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

// CreateTokenObjectLevel creates a security token for object-level permissions
// The token format is: {projectID}/{buildDefinitionID} for object-level permissions
func CreateTokenObjectLevel(projectID, buildDefinitionID string) string {
	return projectID + "/" + buildDefinitionID
}

// PermissionBit represents individual permission bits in Azure DevOps
type PermissionBit int

// Permission bits for repository access control
// Each bit represents a specific permission that can be allowed or denied
const (
	ViewBuilds                               PermissionBit = 1 << iota // 1
	EditBuildQuality                                                   // 2
	RetainIndefinitely                                                 // 4
	DeleteBuilds                                                       // 8
	ManageBuildQualities                                               // 16
	DestroyBuilds                                                      // 32
	UpdateBuildInformation                                             // 64
	QueueBuilds                                                        // 128
	ManageBuildQueue                                                   // 256
	StopBuilds                                                         // 512
	ViewBuildDefinition                                                // 1024
	EditBuildDefinition                                                // 2048
	DeleteBuildDefinition                                              // 4096
	OverrideBuildCheckInValidation                                     // 8192
	AdministerBuildPermissions                                         // 16384
	CreateBuildDefinition                                              // 32768
	EditPipelineQueueConfigurationPermission                           // 65536
	ManageStageRunOrderPermission                                      // 131072
	AbandonBuilds                                                      // 262144
)

// String returns the string representation of the permission bit
func (p PermissionBit) String() string {
	switch p {
	case ViewBuilds:
		return "ViewBuilds"
	case EditBuildQuality:
		return "EditBuildQuality"
	case RetainIndefinitely:
		return "RetainIndefinitely"
	case DeleteBuilds:
		return "DeleteBuilds"
	case ManageBuildQualities:
		return "ManageBuildQualities"
	case DestroyBuilds:
		return "DestroyBuilds"
	case UpdateBuildInformation:
		return "UpdateBuildInformation"
	case QueueBuilds:
		return "QueueBuilds"
	case ManageBuildQueue:
		return "ManageBuildQueue"
	case StopBuilds:
		return "StopBuilds"
	case ViewBuildDefinition:
		return "ViewBuildDefinition"
	case EditBuildDefinition:
		return "EditBuildDefinition"
	case DeleteBuildDefinition:
		return "DeleteBuildDefinition"
	case OverrideBuildCheckInValidation:
		return "OverrideBuildCheckInValidation"
	case AdministerBuildPermissions:
		return "AdministerBuildPermissions"
	case CreateBuildDefinition:
		return "CreateBuildDefinition"
	case EditPipelineQueueConfigurationPermission:
		return "EditPipelineQueueConfigurationPermission"
	case ManageStageRunOrderPermission:
		return "ManageStageRunOrderPermission"
	case AbandonBuilds:
		return "AbandonBuilds"
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
		ViewBuilds:                               comparePermissionBits(bitmask, int(ViewBuilds)),
		EditBuildQuality:                         comparePermissionBits(bitmask, int(EditBuildQuality)),
		RetainIndefinitely:                       comparePermissionBits(bitmask, int(RetainIndefinitely)),
		DeleteBuilds:                             comparePermissionBits(bitmask, int(DeleteBuilds)),
		ManageBuildQualities:                     comparePermissionBits(bitmask, int(ManageBuildQualities)),
		DestroyBuilds:                            comparePermissionBits(bitmask, int(DestroyBuilds)),
		UpdateBuildInformation:                   comparePermissionBits(bitmask, int(UpdateBuildInformation)),
		QueueBuilds:                              comparePermissionBits(bitmask, int(QueueBuilds)),
		ManageBuildQueue:                         comparePermissionBits(bitmask, int(ManageBuildQueue)),
		StopBuilds:                               comparePermissionBits(bitmask, int(StopBuilds)),
		ViewBuildDefinition:                      comparePermissionBits(bitmask, int(ViewBuildDefinition)),
		EditBuildDefinition:                      comparePermissionBits(bitmask, int(EditBuildDefinition)),
		DeleteBuildDefinition:                    comparePermissionBits(bitmask, int(DeleteBuildDefinition)),
		OverrideBuildCheckInValidation:           comparePermissionBits(bitmask, int(OverrideBuildCheckInValidation)),
		AdministerBuildPermissions:               comparePermissionBits(bitmask, int(AdministerBuildPermissions)),
		CreateBuildDefinition:                    comparePermissionBits(bitmask, int(CreateBuildDefinition)),
		EditPipelineQueueConfigurationPermission: comparePermissionBits(bitmask, int(EditPipelineQueueConfigurationPermission)),
		ManageStageRunOrderPermission:            comparePermissionBits(bitmask, int(ManageStageRunOrderPermission)),
		AbandonBuilds:                            comparePermissionBits(bitmask, int(AbandonBuilds)),
	}
}

// permissionFlagsToBitmask converts a PermissionFlags struct to a bitmask
func permissionFlagsToBitmask(f PermissionFlags) int {
	bitmask := 0

	if f.ViewBuilds {
		bitmask |= int(ViewBuilds)
	}
	if f.EditBuildQuality {
		bitmask |= int(EditBuildQuality)
	}
	if f.RetainIndefinitely {
		bitmask |= int(RetainIndefinitely)
	}
	if f.DeleteBuilds {
		bitmask |= int(DeleteBuilds)
	}
	if f.ManageBuildQualities {
		bitmask |= int(ManageBuildQualities)
	}
	if f.DestroyBuilds {
		bitmask |= int(DestroyBuilds)
	}
	if f.UpdateBuildInformation {
		bitmask |= int(UpdateBuildInformation)
	}
	if f.QueueBuilds {
		bitmask |= int(QueueBuilds)
	}
	if f.ManageBuildQueue {
		bitmask |= int(ManageBuildQueue)
	}
	if f.StopBuilds {
		bitmask |= int(StopBuilds)
	}
	if f.ViewBuildDefinition {
		bitmask |= int(ViewBuildDefinition)
	}
	if f.EditBuildDefinition {
		bitmask |= int(EditBuildDefinition)
	}
	if f.DeleteBuildDefinition {
		bitmask |= int(DeleteBuildDefinition)
	}
	if f.OverrideBuildCheckInValidation {
		bitmask |= int(OverrideBuildCheckInValidation)
	}
	if f.AdministerBuildPermissions {
		bitmask |= int(AdministerBuildPermissions)
	}
	if f.CreateBuildDefinition {
		bitmask |= int(CreateBuildDefinition)
	}
	if f.EditPipelineQueueConfigurationPermission {
		bitmask |= int(EditPipelineQueueConfigurationPermission)
	}
	if f.ManageStageRunOrderPermission {
		bitmask |= int(ManageStageRunOrderPermission)
	}
	if f.AbandonBuilds {
		bitmask |= int(AbandonBuilds)
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
		case "ViewBuilds":
			flags.ViewBuilds = true
		case "EditBuildQuality":
			flags.EditBuildQuality = true
		case "RetainIndefinitely":
			flags.RetainIndefinitely = true
		case "DeleteBuilds":
			flags.DeleteBuilds = true
		case "ManageBuildQualities":
			flags.ManageBuildQualities = true
		case "DestroyBuilds":
			flags.DestroyBuilds = true
		case "UpdateBuildInformation":
			flags.UpdateBuildInformation = true
		case "QueueBuilds":
			flags.QueueBuilds = true
		case "ManageBuildQueue":
			flags.ManageBuildQueue = true
		case "StopBuilds":
			flags.StopBuilds = true
		case "ViewBuildDefinition":
			flags.ViewBuildDefinition = true
		case "EditBuildDefinition":
			flags.EditBuildDefinition = true
		case "DeleteBuildDefinition":
			flags.DeleteBuildDefinition = true
		case "OverrideBuildCheckInValidation":
			flags.OverrideBuildCheckInValidation = true
		case "AdministerBuildPermissions":
			flags.AdministerBuildPermissions = true
		case "CreateBuildDefinition":
			flags.CreateBuildDefinition = true
		case "EditPipelineQueueConfigurationPermission":
			flags.EditPipelineQueueConfigurationPermission = true
		case "ManageStageRunOrderPermission":
			flags.ManageStageRunOrderPermission = true
		case "AbandonBuilds":
			flags.AbandonBuilds = true
		}
	}

	return flags
}

// BuildPermissionResponse represents the complete response for build permissions
type BuildPermissionResponse struct {
	Organization      string          `json:"organization" yaml:"organization"`
	ProjectID         string          `json:"projectId" yaml:"projectId"`
	BuildDefinitionID string          `json:"buildDefinitionId" yaml:"buildDefinitionId"`
	ProjectLevel      bool            `json:"projectLevel" yaml:"projectLevel"`
	PermissionsInfo   PermissionsInfo `json:"permissions" yaml:"permissions"`
}

// PermissionsInfo contains the detailed permission information
type PermissionsInfo struct {
	Identity BuildPermissionIdentity `json:"identity" yaml:"identity"`
	Allow    map[string]bool         `json:"allow" yaml:"allow"`
	Deny     map[string]bool         `json:"deny" yaml:"deny"`
}

// BuildPermissionIdentity represents the identity associated with permissions
type BuildPermissionIdentity struct {
	Type       string `json:"type" yaml:"type"`
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
	Descriptor string `json:"descriptor,omitempty" yaml:"descriptor,omitempty"`
}

// PermissionFlags represents all available repository permissions as boolean flags
type PermissionFlags struct {
	ViewBuilds                               bool `json:"ViewBuilds"`
	EditBuildQuality                         bool `json:"EditBuildQuality"`
	RetainIndefinitely                       bool `json:"RetainIndefinitely"`
	DeleteBuilds                             bool `json:"DeleteBuilds"`
	ManageBuildQualities                     bool `json:"ManageBuildQualities"`
	DestroyBuilds                            bool `json:"DestroyBuilds"`
	UpdateBuildInformation                   bool `json:"UpdateBuildInformation"`
	QueueBuilds                              bool `json:"QueueBuilds"`
	ManageBuildQueue                         bool `json:"ManageBuildQueue"`
	StopBuilds                               bool `json:"StopBuilds"`
	ViewBuildDefinition                      bool `json:"ViewBuildDefinition"`
	EditBuildDefinition                      bool `json:"EditBuildDefinition"`
	DeleteBuildDefinition                    bool `json:"DeleteBuildDefinition"`
	OverrideBuildCheckInValidation           bool `json:"OverrideBuildCheckInValidation"`
	AdministerBuildPermissions               bool `json:"AdministerBuildPermissions"`
	CreateBuildDefinition                    bool `json:"CreateBuildDefinition"`
	EditPipelineQueueConfigurationPermission bool `json:"EditPipelineQueueConfigurationPermission"`
	ManageStageRunOrderPermission            bool `json:"ManageStageRunOrderPermission"`
	AbandonBuilds                            bool `json:"AbandonBuilds"`
}

// BuildPermissionRequest represents the request body for setting build permissions
type BuildPermissionRequest struct {
	Permissions PermissionsInfo `json:"permissions" yaml:"permissions"`
}
