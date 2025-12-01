# API Reference

Packages:

- [azuredevops.ogen.krateo.io/v1alpha1](#azuredevopsogenkrateoiov1alpha1)

# azuredevops.ogen.krateo.io/v1alpha1

Resource Types:

- [GitRepository](#gitrepository)




## GitRepository
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
      <td>GitRepository</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#gitrepositoryspec">spec</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gitrepositorystatus">status</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GitRepository.spec
<sup><sup>[↩ Parent](#gitrepository)</sup></sup>





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
        <td><b><a href="#gitrepositoryspecconfigurationref">configurationRef</a></b></td>
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
          PARAMETER: path - Project ID or name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>defaultBranch</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>initialize</b></td>
        <td>boolean</td>
        <td>
          Indicates if the repository should be initialized with an initial commit<br/>
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
        <td><b><a href="#gitrepositoryspecparentrepository">parentRepository</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gitrepositoryspecproject">project</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sourceRef</b></td>
        <td>string</td>
        <td>
          PARAMETER: query - Specify the source refs to use while creating a fork repo<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GitRepository.spec.configurationRef
<sup><sup>[↩ Parent](#gitrepositoryspec)</sup></sup>



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


### GitRepository.spec.parentRepository
<sup><sup>[↩ Parent](#gitrepositoryspec)</sup></sup>





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
        <td><b>id</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#gitrepositoryspecparentrepositoryproject">project</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GitRepository.spec.parentRepository.project
<sup><sup>[↩ Parent](#gitrepositoryspecparentrepository)</sup></sup>





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
        <td><b>id</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GitRepository.spec.project
<sup><sup>[↩ Parent](#gitrepositoryspec)</sup></sup>





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
        <td><b>id</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GitRepository.status
<sup><sup>[↩ Parent](#gitrepository)</sup></sup>





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
        <td><b><a href="#gitrepositorystatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>defaultBranch</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          (format: uuid)<br/>
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
        <td><b>remoteUrl</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sshUrl</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>url</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### GitRepository.status.conditions[index]
<sup><sup>[↩ Parent](#gitrepositorystatus)</sup></sup>



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
