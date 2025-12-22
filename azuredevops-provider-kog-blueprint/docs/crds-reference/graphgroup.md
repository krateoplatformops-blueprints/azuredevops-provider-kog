# API Reference

Packages:

- [azuredevops.ogen.krateo.io/v1alpha1](#azuredevopsogenkrateoiov1alpha1)

# azuredevops.ogen.krateo.io/v1alpha1

Resource Types:

- [GraphGroup](#graphgroup)




## GraphGroup
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
      <td>GraphGroup</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#graphgroupspec">spec</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#graphgroupstatus">status</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GraphGroup.spec
<sup><sup>[↩ Parent](#graphgroup)</sup></sup>





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
        <td><b><a href="#graphgroupspecconfigurationref">configurationRef</a></b></td>
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
        <td><b>description</b></td>
        <td>string</td>
        <td>
          Used by VSTS groups; if set this will be the group description, otherwise ignored<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>descriptor</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>displayName</b></td>
        <td>string</td>
        <td>
          Used by VSTS groups; if set this will be the group DisplayName, otherwise ignored<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>groupDescriptors</b></td>
        <td>string</td>
        <td>
          PARAMETER: query - A comma separated list of descriptors referencing groups you want the graph group to join<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mailAddress</b></td>
        <td>string</td>
        <td>
          This should be the mail address or the group in the source AD or AAD provider. Example: jamal@contoso.com Team Services will communicate with the source provider to fill all other fields on creation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>originId</b></td>
        <td>string</td>
        <td>
          This should be the object id or sid of the group from the source AD or AAD provider. Example: d47d025a-ce2f-4a79-8618-e8862ade30dd Team Services will communicate with the source provider to fill all other fields on creation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>scopeDescriptor</b></td>
        <td>string</td>
        <td>
          PARAMETER: query - Scope descriptor (e.g., project descriptor). When provided, returns all groups within that scope without filtering.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>storageKey</b></td>
        <td>string</td>
        <td>
          Optional: If provided, we will use this identifier for the storage key of the created group (format: uuid)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GraphGroup.spec.configurationRef
<sup><sup>[↩ Parent](#graphgroupspec)</sup></sup>



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


### GraphGroup.status
<sup><sup>[↩ Parent](#graphgroup)</sup></sup>





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
        <td><b><a href="#graphgroupstatus_links">_links</a></b></td>
        <td>object</td>
        <td>
          The class to represent a collection of REST reference links.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#graphgroupstatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>descriptor</b></td>
        <td>string</td>
        <td>
          The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>displayName</b></td>
        <td>string</td>
        <td>
          This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>domain</b></td>
        <td>string</td>
        <td>
          This represents the name of the container of origin for a graph member. (For MSA this is "Windows Live ID", for AD the name of the domain, for AAD the tenantID of the directory, for VSTS groups the ScopeId, etc)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>legacyDescriptor</b></td>
        <td>string</td>
        <td>
          [Internal Use Only] The legacy descriptor is here in case you need to access old version IMS using identity descriptor.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mailAddress</b></td>
        <td>string</td>
        <td>
          The email address of record for a given graph member. This may be different than the principal name.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>origin</b></td>
        <td>string</td>
        <td>
          The type of source provider for the origin identifier (ex:AD, AAD, MSA)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>originId</b></td>
        <td>string</td>
        <td>
          The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>principalName</b></td>
        <td>string</td>
        <td>
          This is the PrincipalName of this graph member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the graph member by VSTS.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>subjectKind</b></td>
        <td>string</td>
        <td>
          This field identifies the type of the graph subject (ex: Group, Scope, User).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>url</b></td>
        <td>string</td>
        <td>
          This url is the full route to the source resource of this graph subject.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GraphGroup.status._links
<sup><sup>[↩ Parent](#graphgroupstatus)</sup></sup>



The class to represent a collection of REST reference links.

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
        <td><b>links</b></td>
        <td>object</td>
        <td>
          The readonly view of the links.  Because Reference links are readonly, we only want to expose them as read only.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GraphGroup.status.conditions[index]
<sup><sup>[↩ Parent](#graphgroupstatus)</sup></sup>



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
