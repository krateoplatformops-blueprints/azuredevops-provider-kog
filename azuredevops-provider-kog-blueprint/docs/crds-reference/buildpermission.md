# API Reference

Packages:

- [azuredevops.ogen.krateo.io/v1alpha1](#azuredevopsogenkrateoiov1alpha1)

# azuredevops.ogen.krateo.io/v1alpha1

Resource Types:

- [BuildPermission](#buildpermission)




## BuildPermission
<sup><sup>[↩ Parent](#azuredevopsogenkrateoiov1alpha1 )</sup></sup>








<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>azuredevops.ogen.krateo.io/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>BuildPermission</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#buildpermissionspec">spec</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#buildpermissionstatus">status</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.spec
<sup><sup>[↩ Parent](#buildpermission)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#buildpermissionspecconfigurationref">configurationRef</a></b></td>
        <td>object</td>
        <td>
          A reference to the Configuration CR that holds all the needed configuration for this resource. OASGen Provider added this field automatically.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>organization</b></td>
        <td>string</td>
        <td>
          PARAMETER: path - Organization name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>projectId</b></td>
        <td>string</td>
        <td>
          PARAMETER: path - Project ID<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>buildDefinitionId</b></td>
        <td>string</td>
        <td>
          PARAMETER: query - Build Definition ID<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#buildpermissionspecpermissions">permissions</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>projectLevel</b></td>
        <td>boolean</td>
        <td>
          PARAMETER: query - Whether to manage permissions at the project level (true) or build level (pipeline level) (false). Default is false.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.spec.configurationRef
<sup><sup>[↩ Parent](#buildpermissionspec)</sup></sup>



A reference to the Configuration CR that holds all the needed configuration for this resource. OASGen Provider added this field automatically.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          Namespace of the referenced Configuration CR. If not provided, the same namespace will be used.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.spec.permissions
<sup><sup>[↩ Parent](#buildpermissionspec)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#buildpermissionspecpermissionsallow">allow</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#buildpermissionspecpermissionsdeny">deny</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#buildpermissionspecpermissionsidentity">identity</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.spec.permissions.allow
<sup><sup>[↩ Parent](#buildpermissionspecpermissions)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>AbandonBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Abandon builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>AdministerBuildPermissions</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Administer build permissions'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DestroyBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Destroy builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditBuildQuality</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit build quality'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditPipelineQueueConfigurationPermission</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit queue build configuration'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageBuildQualities</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage build qualities'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageBuildQueue</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage build queue'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageStageRunOrderPermission</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage stage run order'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>OverrideBuildCheckInValidation</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Override check-in validation by build'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>QueueBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Queue builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RetainIndefinitely</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Retain indefinitely'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>StopBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Stop builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>UpdateBuildInformation</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Update build information'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'View build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'View Builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.spec.permissions.deny
<sup><sup>[↩ Parent](#buildpermissionspecpermissions)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>AbandonBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Abandon builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>AdministerBuildPermissions</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Administer build permissions'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DestroyBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Destroy builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditBuildQuality</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit build quality'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditPipelineQueueConfigurationPermission</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit queue build configuration'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageBuildQualities</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage build qualities'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageBuildQueue</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage build queue'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageStageRunOrderPermission</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage stage run order'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>OverrideBuildCheckInValidation</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Override check-in validation by build'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>QueueBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Queue builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RetainIndefinitely</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Retain indefinitely'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>StopBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Stop builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>UpdateBuildInformation</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Update build information'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'View build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'View Builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.spec.permissions.identity
<sup><sup>[↩ Parent](#buildpermissionspecpermissions)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>descriptor</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.status
<sup><sup>[↩ Parent](#buildpermission)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#buildpermissionstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#buildpermissionstatuspermissions">permissions</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.status.conditions[index]
<sup><sup>[↩ Parent](#buildpermissionstatus)</sup></sup>



A Condition that may apply to a resource.

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>lastTransitionTime</b></td>
        <td>string</td>
        <td>
          LastTransitionTime is the last time this condition transitioned from one
status to another.<br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>reason</b></td>
        <td>string</td>
        <td>
          A Reason for this condition's last transition from one status to another.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          Status of this condition; is it currently True, False, or Unknown?<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of this condition. At most one of each condition type may apply to
a resource at any point in time.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          A Message containing details about this condition's last transition from
one status to another, if any.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.status.permissions
<sup><sup>[↩ Parent](#buildpermissionstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#buildpermissionstatuspermissionsallow">allow</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#buildpermissionstatuspermissionsdeny">deny</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#buildpermissionstatuspermissionsidentity">identity</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.status.permissions.allow
<sup><sup>[↩ Parent](#buildpermissionstatuspermissions)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>AbandonBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Abandon builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>AdministerBuildPermissions</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Administer build permissions'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DestroyBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Destroy builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditBuildQuality</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit build quality'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditPipelineQueueConfigurationPermission</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit queue build configuration'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageBuildQualities</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage build qualities'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageBuildQueue</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage build queue'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageStageRunOrderPermission</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage stage run order'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>OverrideBuildCheckInValidation</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Override check-in validation by build'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>QueueBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Queue builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RetainIndefinitely</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Retain indefinitely'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>StopBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Stop builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>UpdateBuildInformation</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Update build information'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'View build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'View Builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.status.permissions.deny
<sup><sup>[↩ Parent](#buildpermissionstatuspermissions)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>AbandonBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Abandon builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>AdministerBuildPermissions</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Administer build permissions'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DestroyBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Destroy builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditBuildQuality</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit build quality'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditPipelineQueueConfigurationPermission</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit queue build configuration'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageBuildQualities</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage build qualities'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageBuildQueue</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage build queue'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageStageRunOrderPermission</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage stage run order'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>OverrideBuildCheckInValidation</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Override check-in validation by build'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>QueueBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Queue builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RetainIndefinitely</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Retain indefinitely'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>StopBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Stop builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>UpdateBuildInformation</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Update build information'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewBuildDefinition</b></td>
        <td>boolean</td>
        <td>
          displayName: 'View build pipeline'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewBuilds</b></td>
        <td>boolean</td>
        <td>
          displayName: 'View Builds'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### BuildPermission.status.permissions.identity
<sup><sup>[↩ Parent](#buildpermissionstatuspermissions)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>descriptor</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>
