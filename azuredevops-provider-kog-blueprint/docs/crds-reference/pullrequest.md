# API Reference

Packages:

- [azuredevops.ogen.krateo.io/v1alpha1](#azuredevopsogenkrateoiov1alpha1)

# azuredevops.ogen.krateo.io/v1alpha1

Resource Types:

- [PullRequest](#pullrequest)




## PullRequest
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
      <td>PullRequest</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#pullrequestspec">spec</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatus">status</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.spec
<sup><sup>[↩ Parent](#pullrequest)</sup></sup>





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
        <td><b><a href="#pullrequestspecconfigurationref">configurationRef</a></b></td>
        <td>object</td>
        <td>
          A reference to the Configuration CR that holds all the needed configuration for this resource. OASGen Provider added this field automatically.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>organization</b></td>
        <td>string</td>
        <td>
          PARAMETER: path - The name of the Azure DevOps organization.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>project</b></td>
        <td>string</td>
        <td>
          PARAMETER: path - Project ID or project name<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>repositoryId</b></td>
        <td>string</td>
        <td>
          PARAMETER: path - The repository ID of the pull request's target branch.<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#pullrequestspecautocompletesetby">autoCompleteSetBy</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequestspeccompletionoptions">completionOptions</a></b></td>
        <td>object</td>
        <td>
          Preferences about how the pull request should be completed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          The description of the pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>ignoreTargetRefAndChooseDynamically</b></td>
        <td>boolean</td>
        <td>
          This optional parameter allows clients to use server-side dynamic choices for the target ref. Due to preexisting contracts, users _must_ specify a target ref, but this option will cause the server to ignore it and choose dynamically from the user's favorites (or the default branch).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isDraft</b></td>
        <td>boolean</td>
        <td>
          Draft / WIP pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequestspecmergeoptions">mergeOptions</a></b></td>
        <td>object</td>
        <td>
          The options which are used when a pull request merge is created.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sourceRefName</b></td>
        <td>string</td>
        <td>
          The name of the source branch of the pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>enum</td>
        <td>
          The status of the pull request.<br/>
          <br/>
            <i>Enum</i>: notSet, active, abandoned, completed, all<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>supportsIterations</b></td>
        <td>boolean</td>
        <td>
          If true, this pull request supports multiple iterations. Iteration support means individual pushes to the source branch of the pull request can be reviewed and comments left in one iteration will be tracked across future iterations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>targetRefName</b></td>
        <td>string</td>
        <td>
          The name of the target branch of the pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>title</b></td>
        <td>string</td>
        <td>
          The title of the pull request.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.spec.configurationRef
<sup><sup>[↩ Parent](#pullrequestspec)</sup></sup>



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


### PullRequest.spec.autoCompleteSetBy
<sup><sup>[↩ Parent](#pullrequestspec)</sup></sup>





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


### PullRequest.spec.completionOptions
<sup><sup>[↩ Parent](#pullrequestspec)</sup></sup>



Preferences about how the pull request should be completed.

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
        <td><b>autoCompleteIgnoreConfigIds</b></td>
        <td>[]integer</td>
        <td>
          List of any policy configuration Id's which auto-complete should not wait for. Only applies to optional policies (isBlocking == false). Auto-complete always waits for required policies (isBlocking == true).<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>bypassPolicy</b></td>
        <td>boolean</td>
        <td>
          If true, policies will be explicitly bypassed while the pull request is completed.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>bypassReason</b></td>
        <td>string</td>
        <td>
          If policies are bypassed, this reason is stored as to why bypass was used.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>deleteSourceBranch</b></td>
        <td>boolean</td>
        <td>
          If true, the source branch of the pull request will be deleted after completion.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mergeCommitMessage</b></td>
        <td>string</td>
        <td>
          If set, this will be used as the commit message of the merge commit.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mergeStrategy</b></td>
        <td>enum</td>
        <td>
          Specify the strategy used to merge the pull request during completion. If MergeStrategy is not set to any value, a no-FF merge will be created if SquashMerge == false. If MergeStrategy is not set to any value, the pull request commits will be squashed if SquashMerge == true. The SquashMerge property is deprecated. It is recommended that you explicitly set MergeStrategy in all cases. If an explicit value is provided for MergeStrategy, the SquashMerge property will be ignored.<br/>
          <br/>
            <i>Enum</i>: noFastForward, squash, rebase, rebaseMerge<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>squashMerge</b></td>
        <td>boolean</td>
        <td>
          SquashMerge is deprecated. You should explicitly set the value of MergeStrategy. If MergeStrategy is set to any value, the SquashMerge value will be ignored. If MergeStrategy is not set, the merge strategy will be no-fast-forward if this flag is false, or squash if true.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>transitionWorkItems</b></td>
        <td>boolean</td>
        <td>
          If true, we will attempt to transition any work items linked to the pull request into the next logical state (i.e. Active -> Resolved)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.spec.mergeOptions
<sup><sup>[↩ Parent](#pullrequestspec)</sup></sup>



The options which are used when a pull request merge is created.

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
        <td><b>conflictAuthorshipCommits</b></td>
        <td>boolean</td>
        <td>
          If true, conflict resolutions applied during the merge will be put in separate commits to preserve authorship info for git blame, etc.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>detectRenameFalsePositives</b></td>
        <td>boolean</td>
        <td>
          If true, renames where there is more than one valid way to map the original file locations to renamed file locations will be treated as false positives and ignored.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>disableRenames</b></td>
        <td>boolean</td>
        <td>
          If true, rename detection will not be performed during the merge.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status
<sup><sup>[↩ Parent](#pullrequest)</sup></sup>





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
        <td><b><a href="#pullrequeststatusclosedby">closedBy</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>closedDate</b></td>
        <td>string</td>
        <td>
          The date when the pull request was closed (completed, abandoned, or merged externally). (format: date-time)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatusconditionsindex">conditions</a></b></td>
        <td>[]object</td>
        <td>
          Conditions of the resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuscreatedby">createdBy</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>creationDate</b></td>
        <td>string</td>
        <td>
          The date when the pull request was created. (format: date-time)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isDraft</b></td>
        <td>boolean</td>
        <td>
          Draft / WIP pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslabelsindex">labels</a></b></td>
        <td>[]object</td>
        <td>
          The labels associated with the pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommit">lastMergeSourceCommit</a></b></td>
        <td>object</td>
        <td>
          Provides properties that describe a Git commit and associated metadata.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>mergeStatus</b></td>
        <td>enum</td>
        <td>
          The current status of the pull request merge.<br/>
          <br/>
            <i>Enum</i>: notSet, queued, conflicts, succeeded, rejectedByPolicy, failure<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pullRequestId</b></td>
        <td>integer</td>
        <td>
          The ID of the pull request. (format: int32)<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>remoteUrl</b></td>
        <td>string</td>
        <td>
          Used internally.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatusreviewersindex">reviewers</a></b></td>
        <td>[]object</td>
        <td>
          A list of reviewers on the pull request along with the state of their votes.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sourceRefName</b></td>
        <td>string</td>
        <td>
          The name of the source branch of the pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>enum</td>
        <td>
          The status of the pull request.<br/>
          <br/>
            <i>Enum</i>: notSet, active, abandoned, completed, all<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>targetRefName</b></td>
        <td>string</td>
        <td>
          The name of the target branch of the pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>title</b></td>
        <td>string</td>
        <td>
          The title of the pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>url</b></td>
        <td>string</td>
        <td>
          Used internally.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.closedBy
<sup><sup>[↩ Parent](#pullrequeststatus)</sup></sup>





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
        <td><b><a href="#pullrequeststatusclosedby_links">_links</a></b></td>
        <td>object</td>
        <td>
          The class to represent a collection of REST reference links.<br/>
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
        <td><b>directoryAlias</b></td>
        <td>string</td>
        <td>
          Deprecated - Can be retrieved by querying the Graph user referenced in the "self" entry of the IdentityRef "_links" dictionary<br/>
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
        <td><b>id</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imageUrl</b></td>
        <td>string</td>
        <td>
          Deprecated - Available in the "avatar" entry of the IdentityRef "_links" dictionary<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>inactive</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be retrieved by querying the Graph membership state referenced in the "membershipState" entry of the GraphUser "_links" dictionary<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isAadIdentity</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsAadUserType/Descriptor.IsAadGroupType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isContainer</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsGroupType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isDeletedInOrigin</b></td>
        <td>boolean</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>profileUrl</b></td>
        <td>string</td>
        <td>
          Deprecated - not in use in most preexisting implementations of ToIdentityRef<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>uniqueName</b></td>
        <td>string</td>
        <td>
          Deprecated - use Domain+PrincipalName instead<br/>
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


### PullRequest.status.closedBy._links
<sup><sup>[↩ Parent](#pullrequeststatusclosedby)</sup></sup>



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


### PullRequest.status.conditions[index]
<sup><sup>[↩ Parent](#pullrequeststatus)</sup></sup>



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


### PullRequest.status.createdBy
<sup><sup>[↩ Parent](#pullrequeststatus)</sup></sup>





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
        <td><b><a href="#pullrequeststatuscreatedby_links">_links</a></b></td>
        <td>object</td>
        <td>
          The class to represent a collection of REST reference links.<br/>
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
        <td><b>directoryAlias</b></td>
        <td>string</td>
        <td>
          Deprecated - Can be retrieved by querying the Graph user referenced in the "self" entry of the IdentityRef "_links" dictionary<br/>
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
        <td><b>id</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imageUrl</b></td>
        <td>string</td>
        <td>
          Deprecated - Available in the "avatar" entry of the IdentityRef "_links" dictionary<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>inactive</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be retrieved by querying the Graph membership state referenced in the "membershipState" entry of the GraphUser "_links" dictionary<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isAadIdentity</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsAadUserType/Descriptor.IsAadGroupType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isContainer</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsGroupType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isDeletedInOrigin</b></td>
        <td>boolean</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>profileUrl</b></td>
        <td>string</td>
        <td>
          Deprecated - not in use in most preexisting implementations of ToIdentityRef<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>uniqueName</b></td>
        <td>string</td>
        <td>
          Deprecated - use Domain+PrincipalName instead<br/>
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


### PullRequest.status.createdBy._links
<sup><sup>[↩ Parent](#pullrequeststatuscreatedby)</sup></sup>



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


### PullRequest.status.labels[index]
<sup><sup>[↩ Parent](#pullrequeststatus)</sup></sup>





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
        <td><b>active</b></td>
        <td>boolean</td>
        <td>
          Whether or not the tag definition is active.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          ID of the tag definition. (format: uuid)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          The name of the tag definition.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>url</b></td>
        <td>string</td>
        <td>
          Resource URL for the Tag Definition.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.lastMergeSourceCommit
<sup><sup>[↩ Parent](#pullrequeststatus)</sup></sup>



Provides properties that describe a Git commit and associated metadata.

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
        <td><b><a href="#pullrequeststatuslastmergesourcecommit_links">_links</a></b></td>
        <td>object</td>
        <td>
          The class to represent a collection of REST reference links.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitauthor">author</a></b></td>
        <td>object</td>
        <td>
          User info and date for Git operations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>changeCounts</b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitchangesindex">changes</a></b></td>
        <td>[]object</td>
        <td>
          An enumeration of the changes included with the commit.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>comment</b></td>
        <td>string</td>
        <td>
          Comment or message of the commit.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>commentTruncated</b></td>
        <td>boolean</td>
        <td>
          Indicates if the comment is truncated from the full Git commit comment message.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>commitId</b></td>
        <td>string</td>
        <td>
          ID (SHA-1) of the commit.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>commitTooManyChanges</b></td>
        <td>boolean</td>
        <td>
          Indicates that commit contains too many changes to be displayed<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitcommitter">committer</a></b></td>
        <td>object</td>
        <td>
          User info and date for Git operations.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>parents</b></td>
        <td>[]string</td>
        <td>
          An enumeration of the parent commit IDs for this commit.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitpush">push</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>remoteUrl</b></td>
        <td>string</td>
        <td>
          Remote URL path to the commit.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitstatusesindex">statuses</a></b></td>
        <td>[]object</td>
        <td>
          A list of status metadata from services and extensions that may associate additional information to the commit.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>url</b></td>
        <td>string</td>
        <td>
          REST URL for this resource.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitworkitemsindex">workItems</a></b></td>
        <td>[]object</td>
        <td>
          A list of workitems associated with this commit.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.lastMergeSourceCommit._links
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommit)</sup></sup>



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


### PullRequest.status.lastMergeSourceCommit.author
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommit)</sup></sup>



User info and date for Git operations.

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
        <td><b>date</b></td>
        <td>string</td>
        <td>
          Date of the Git operation. (format: date-time)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>email</b></td>
        <td>string</td>
        <td>
          Email address of the user performing the Git operation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imageUrl</b></td>
        <td>string</td>
        <td>
          Url for the user's avatar.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the user performing the Git operation.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.lastMergeSourceCommit.changes[index]
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommit)</sup></sup>





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
        <td><b>changeId</b></td>
        <td>integer</td>
        <td>
          ID of the change within the group of changes. (format: int32)<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>changeType</b></td>
        <td>enum</td>
        <td>
          The type of change that was made to the item.<br/>
          <br/>
            <i>Enum</i>: none, add, edit, encoding, rename, delete, undelete, branch, merge, lock, rollback, sourceRename, targetRename, property, all<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>item</b></td>
        <td>string</td>
        <td>
          Current version. (format: T)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitchangesindexnewcontent">newContent</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitchangesindexnewcontenttemplate">newContentTemplate</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>originalPath</b></td>
        <td>string</td>
        <td>
          Original path of item if different from current path.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>sourceServerItem</b></td>
        <td>string</td>
        <td>
          Path of the item on the server.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>url</b></td>
        <td>string</td>
        <td>
          URL to retrieve the item.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.lastMergeSourceCommit.changes[index].newContent
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommitchangesindex)</sup></sup>





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
        <td><b>content</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>contentType</b></td>
        <td>enum</td>
        <td>
          <br/>
          <br/>
            <i>Enum</i>: rawText, base64Encoded<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.lastMergeSourceCommit.changes[index].newContentTemplate
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommitchangesindex)</sup></sup>





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
          Name of the Template<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>type</b></td>
        <td>string</td>
        <td>
          Type of the Template<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.lastMergeSourceCommit.committer
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommit)</sup></sup>



User info and date for Git operations.

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
        <td><b>date</b></td>
        <td>string</td>
        <td>
          Date of the Git operation. (format: date-time)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>email</b></td>
        <td>string</td>
        <td>
          Email address of the user performing the Git operation.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imageUrl</b></td>
        <td>string</td>
        <td>
          Url for the user's avatar.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name of the user performing the Git operation.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.lastMergeSourceCommit.push
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommit)</sup></sup>





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
        <td><b><a href="#pullrequeststatuslastmergesourcecommitpush_links">_links</a></b></td>
        <td>object</td>
        <td>
          The class to represent a collection of REST reference links.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>date</b></td>
        <td>string</td>
        <td>
          (format: date-time)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>pushId</b></td>
        <td>integer</td>
        <td>
          (format: int32)<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitpushpushedby">pushedBy</a></b></td>
        <td>object</td>
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


### PullRequest.status.lastMergeSourceCommit.push._links
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommitpush)</sup></sup>



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


### PullRequest.status.lastMergeSourceCommit.push.pushedBy
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommitpush)</sup></sup>





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
        <td><b><a href="#pullrequeststatuslastmergesourcecommitpushpushedby_links">_links</a></b></td>
        <td>object</td>
        <td>
          The class to represent a collection of REST reference links.<br/>
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
        <td><b>directoryAlias</b></td>
        <td>string</td>
        <td>
          Deprecated - Can be retrieved by querying the Graph user referenced in the "self" entry of the IdentityRef "_links" dictionary<br/>
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
        <td><b>id</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imageUrl</b></td>
        <td>string</td>
        <td>
          Deprecated - Available in the "avatar" entry of the IdentityRef "_links" dictionary<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>inactive</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be retrieved by querying the Graph membership state referenced in the "membershipState" entry of the GraphUser "_links" dictionary<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isAadIdentity</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsAadUserType/Descriptor.IsAadGroupType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isContainer</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsGroupType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isDeletedInOrigin</b></td>
        <td>boolean</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>profileUrl</b></td>
        <td>string</td>
        <td>
          Deprecated - not in use in most preexisting implementations of ToIdentityRef<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>uniqueName</b></td>
        <td>string</td>
        <td>
          Deprecated - use Domain+PrincipalName instead<br/>
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


### PullRequest.status.lastMergeSourceCommit.push.pushedBy._links
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommitpushpushedby)</sup></sup>



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


### PullRequest.status.lastMergeSourceCommit.statuses[index]
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommit)</sup></sup>





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
        <td><b><a href="#pullrequeststatuslastmergesourcecommitstatusesindex_links">_links</a></b></td>
        <td>object</td>
        <td>
          The class to represent a collection of REST reference links.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitstatusesindexcontext">context</a></b></td>
        <td>object</td>
        <td>
          Status context that uniquely identifies the status.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#pullrequeststatuslastmergesourcecommitstatusesindexcreatedby">createdBy</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>creationDate</b></td>
        <td>string</td>
        <td>
          Creation date and time of the status. (format: date-time)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>description</b></td>
        <td>string</td>
        <td>
          Status description. Typically describes current state of the status.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>integer</td>
        <td>
          Status identifier. (format: int32)<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>state</b></td>
        <td>enum</td>
        <td>
          State of the status.<br/>
          <br/>
            <i>Enum</i>: notSet, pending, succeeded, failed, error, notApplicable, partiallySucceeded<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>targetUrl</b></td>
        <td>string</td>
        <td>
          URL with status details.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>updatedDate</b></td>
        <td>string</td>
        <td>
          Last update date and time of the status. (format: date-time)<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.lastMergeSourceCommit.statuses[index]._links
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommitstatusesindex)</sup></sup>



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


### PullRequest.status.lastMergeSourceCommit.statuses[index].context
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommitstatusesindex)</sup></sup>



Status context that uniquely identifies the status.

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
        <td><b>genre</b></td>
        <td>string</td>
        <td>
          Genre of the status. Typically name of the service/tool generating the status, can be empty.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          Name identifier of the status, cannot be null or empty.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.lastMergeSourceCommit.statuses[index].createdBy
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommitstatusesindex)</sup></sup>





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
        <td><b><a href="#pullrequeststatuslastmergesourcecommitstatusesindexcreatedby_links">_links</a></b></td>
        <td>object</td>
        <td>
          The class to represent a collection of REST reference links.<br/>
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
        <td><b>directoryAlias</b></td>
        <td>string</td>
        <td>
          Deprecated - Can be retrieved by querying the Graph user referenced in the "self" entry of the IdentityRef "_links" dictionary<br/>
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
        <td><b>id</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imageUrl</b></td>
        <td>string</td>
        <td>
          Deprecated - Available in the "avatar" entry of the IdentityRef "_links" dictionary<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>inactive</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be retrieved by querying the Graph membership state referenced in the "membershipState" entry of the GraphUser "_links" dictionary<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isAadIdentity</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsAadUserType/Descriptor.IsAadGroupType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isContainer</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsGroupType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isDeletedInOrigin</b></td>
        <td>boolean</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>profileUrl</b></td>
        <td>string</td>
        <td>
          Deprecated - not in use in most preexisting implementations of ToIdentityRef<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>uniqueName</b></td>
        <td>string</td>
        <td>
          Deprecated - use Domain+PrincipalName instead<br/>
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


### PullRequest.status.lastMergeSourceCommit.statuses[index].createdBy._links
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommitstatusesindexcreatedby)</sup></sup>



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


### PullRequest.status.lastMergeSourceCommit.workItems[index]
<sup><sup>[↩ Parent](#pullrequeststatuslastmergesourcecommit)</sup></sup>





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
        <td><b>url</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.reviewers[index]
<sup><sup>[↩ Parent](#pullrequeststatus)</sup></sup>





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
        <td><b><a href="#pullrequeststatusreviewersindex_links">_links</a></b></td>
        <td>object</td>
        <td>
          The class to represent a collection of REST reference links.<br/>
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
        <td><b>directoryAlias</b></td>
        <td>string</td>
        <td>
          Deprecated - Can be retrieved by querying the Graph user referenced in the "self" entry of the IdentityRef "_links" dictionary<br/>
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
        <td><b>hasDeclined</b></td>
        <td>boolean</td>
        <td>
          Indicates if this reviewer has declined to review this pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>id</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>imageUrl</b></td>
        <td>string</td>
        <td>
          Deprecated - Available in the "avatar" entry of the IdentityRef "_links" dictionary<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>inactive</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be retrieved by querying the Graph membership state referenced in the "membershipState" entry of the GraphUser "_links" dictionary<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isAadIdentity</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsAadUserType/Descriptor.IsAadGroupType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isContainer</b></td>
        <td>boolean</td>
        <td>
          Deprecated - Can be inferred from the subject type of the descriptor (Descriptor.IsGroupType)<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isDeletedInOrigin</b></td>
        <td>boolean</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isFlagged</b></td>
        <td>boolean</td>
        <td>
          Indicates if this reviewer is flagged for attention on this pull request.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isReapprove</b></td>
        <td>boolean</td>
        <td>
          Indicates if this approve vote should still be handled even though vote didn't change.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>isRequired</b></td>
        <td>boolean</td>
        <td>
          Indicates if this is a required reviewer for this pull request. <br /> Branches can have policies that require particular reviewers are required for pull requests.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>profileUrl</b></td>
        <td>string</td>
        <td>
          Deprecated - not in use in most preexisting implementations of ToIdentityRef<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>reviewerUrl</b></td>
        <td>string</td>
        <td>
          URL to retrieve information about this identity<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>uniqueName</b></td>
        <td>string</td>
        <td>
          Deprecated - use Domain+PrincipalName instead<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>url</b></td>
        <td>string</td>
        <td>
          This url is the full route to the source resource of this graph subject.<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>vote</b></td>
        <td>integer</td>
        <td>
          Vote on a pull request:<br /> 10 - approved 5 - approved with suggestions 0 - no vote -5 - waiting for author -10 - rejected (format: int16)<br/>
          <br/>
            <i>Format</i>: int32<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>votedFor</b></td>
        <td>[]object</td>
        <td>
          Groups or teams that this reviewer contributed to. <br /> Groups and teams can be reviewers on pull requests but can not vote directly.  When a member of the group or team votes, that vote is rolled up into the group or team vote.  VotedFor is a list of such votes.<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### PullRequest.status.reviewers[index]._links
<sup><sup>[↩ Parent](#pullrequeststatusreviewersindex)</sup></sup>



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
