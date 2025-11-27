# API Reference

Packages:

- [azuredevops.ogen.krateo.io/v1alpha1](#azuredevopsogenkrateoiov1alpha1)

# azuredevops.ogen.krateo.io/v1alpha1

Resource Types:

- [RepositoryPermission](#repositorypermission)




## RepositoryPermission
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
      <td>RepositoryPermission</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#repositorypermissionspec">spec</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#repositorypermissionstatus">status</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RepositoryPermission.spec
<sup><sup>[↩ Parent](#repositorypermission)</sup></sup>





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
        <td><b><a href="#repositorypermissionspecconfigurationref">configurationRef</a></b></td>
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
        <td><b><a href="#repositorypermissionspecpermissions">permissions</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>projectLevel</b></td>
        <td>boolean</td>
        <td>
          PARAMETER: query - Whether to manage permissions at the project level (true) or repository level (false). Default is false.<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>repositoryId</b></td>
        <td>string</td>
        <td>
          PARAMETER: query - Repository ID<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RepositoryPermission.spec.configurationRef
<sup><sup>[↩ Parent](#repositorypermissionspec)</sup></sup>



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


### RepositoryPermission.spec.permissions
<sup><sup>[↩ Parent](#repositorypermissionspec)</sup></sup>





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
        <td><b><a href="#repositorypermissionspecpermissionsallow">allow</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#repositorypermissionspecpermissionsdeny">deny</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#repositorypermissionspecpermissionsidentity">identity</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RepositoryPermission.spec.permissions.allow
<sup><sup>[↩ Parent](#repositorypermissionspecpermissions)</sup></sup>





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
        <td><b>Administer</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Administer'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateBranch</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create branch'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateTag</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create tag'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete or disable repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DismissAdvSecAlerts</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: manage and dismiss alerts'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditPolicies</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit policies'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ForcePush</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Force push (rewrite history, delete branches and tags)'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>GenericContribute</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Contribute'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>GenericRead</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Read'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageAdvSecScanning</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: manage settings'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageNote</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage notes'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManagePermissions</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage permissions'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PolicyExempt</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Bypass policies when pushing'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PullRequestBypassPolicy</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Bypass policies when completing pull requests'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PullRequestContribute</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Contribute to pull requests'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RemoveOthersLocks</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Remove others' locks'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RenameRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Rename repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewAdvSecAlerts</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: view alerts'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RepositoryPermission.spec.permissions.deny
<sup><sup>[↩ Parent](#repositorypermissionspecpermissions)</sup></sup>





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
        <td><b>Administer</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Administer'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateBranch</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create branch'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateTag</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create tag'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete or disable repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DismissAdvSecAlerts</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: manage and dismiss alerts'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditPolicies</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit policies'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ForcePush</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Force push (rewrite history, delete branches and tags)'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>GenericContribute</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Contribute'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>GenericRead</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Read'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageAdvSecScanning</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: manage settings'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageNote</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage notes'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManagePermissions</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage permissions'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PolicyExempt</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Bypass policies when pushing'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PullRequestBypassPolicy</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Bypass policies when completing pull requests'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PullRequestContribute</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Contribute to pull requests'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RemoveOthersLocks</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Remove others' locks'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RenameRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Rename repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewAdvSecAlerts</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: view alerts'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RepositoryPermission.spec.permissions.identity
<sup><sup>[↩ Parent](#repositorypermissionspecpermissions)</sup></sup>





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


### RepositoryPermission.status
<sup><sup>[↩ Parent](#repositorypermission)</sup></sup>





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
        <td><b><a href="#repositorypermissionstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#repositorypermissionstatuspermissions">permissions</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RepositoryPermission.status.conditions[index]
<sup><sup>[↩ Parent](#repositorypermissionstatus)</sup></sup>



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


### RepositoryPermission.status.permissions
<sup><sup>[↩ Parent](#repositorypermissionstatus)</sup></sup>





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
        <td><b><a href="#repositorypermissionstatuspermissionsallow">allow</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#repositorypermissionstatuspermissionsdeny">deny</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#repositorypermissionstatuspermissionsidentity">identity</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RepositoryPermission.status.permissions.allow
<sup><sup>[↩ Parent](#repositorypermissionstatuspermissions)</sup></sup>





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
        <td><b>Administer</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Administer'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateBranch</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create branch'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateTag</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create tag'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete or disable repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DismissAdvSecAlerts</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: manage and dismiss alerts'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditPolicies</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit policies'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ForcePush</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Force push (rewrite history, delete branches and tags)'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>GenericContribute</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Contribute'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>GenericRead</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Read'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageAdvSecScanning</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: manage settings'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageNote</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage notes'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManagePermissions</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage permissions'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PolicyExempt</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Bypass policies when pushing'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PullRequestBypassPolicy</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Bypass policies when completing pull requests'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PullRequestContribute</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Contribute to pull requests'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RemoveOthersLocks</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Remove others' locks'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RenameRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Rename repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewAdvSecAlerts</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: view alerts'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RepositoryPermission.status.permissions.deny
<sup><sup>[↩ Parent](#repositorypermissionstatuspermissions)</sup></sup>





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
        <td><b>Administer</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Administer'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateBranch</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create branch'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>CreateTag</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Create tag'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DeleteRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Delete or disable repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>DismissAdvSecAlerts</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: manage and dismiss alerts'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>EditPolicies</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Edit policies'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ForcePush</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Force push (rewrite history, delete branches and tags)'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>GenericContribute</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Contribute'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>GenericRead</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Read'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageAdvSecScanning</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: manage settings'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManageNote</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage notes'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ManagePermissions</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Manage permissions'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PolicyExempt</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Bypass policies when pushing'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PullRequestBypassPolicy</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Bypass policies when completing pull requests'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>PullRequestContribute</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Contribute to pull requests'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RemoveOthersLocks</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Remove others' locks'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>RenameRepository</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Rename repository'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ViewAdvSecAlerts</b></td>
        <td>boolean</td>
        <td>
          displayName: 'Advanced Security: view alerts'<br/>
          <br/>
            <i>Default</i>: false<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### RepositoryPermission.status.permissions.identity
<sup><sup>[↩ Parent](#repositorypermissionstatuspermissions)</sup></sup>





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
