# Krateo Azure DevOps Provider KOG - Troubleshooting Guide

This document serves as a troubleshooting guide for the Krateo Azure DevOps Provider KOG.

## Summary

- [Summary](#summary)
- [GitRepository](#gitrepository)
  - [1. New Repository Cases](#1-new-repository-cases)
    - [Case 1.1: New Repo Without Initialization](#case-11-new-repo-without-initialization)
    - [Case 1.2: New Repo With Initialization (Default Branch Omitted)](#case-12-new-repo-with-initialization-default-branch-omitted)
    - [Case 1.3: New Repo With Initialization and Custom Default Branch](#case-13-new-repo-with-initialization-and-custom-default-branch)
    - [Case 1.4: New Repo With Custom Default Branch (Without Initialization Flag set to `true`)](#case-14-new-repo-with-custom-default-branch-without-initialization-flag-set-to-true)
  - [2. Forked Repository Cases (Parent Repository Initialized)](#2-forked-repository-cases-parent-repository-initialized)
    - [Case 2.1: Forked Repo With Nonexistent SourceRef](#case-21-forked-repo-with-nonexistent-sourceref)
    - [Case 2.2: Forked Repo With Existing Default Branch in Parent](#case-22-forked-repo-with-existing-default-branch-in-parent)
    - [Case 2.3: Forked Repo With Nonexistent Default Branch (Fallback to Parent Default)](#case-23-forked-repo-with-nonexistent-default-branch-fallback-to-parent-default)
  - [3. Forked Repository Cases (Parent Repository Not Initialized)](#3-forked-repository-cases-parent-repository-not-initialized)
    - [Case 3.1: Fork With Invalid `sourceRef`](#case-31-fork-with-invalid-sourceref)
    - [Case 3.2: Fork With No Branch Configuration](#case-32-fork-with-no-branch-configuration)
    - [Case 3.3: Fork With Nonexistent Default Branch](#case-33-fork-with-nonexistent-default-branch)
- [PolicyConfiguration](#policyconfiguration)
- [PipelinePermission](#pipelinepermission)
  - [PipelinePermission for Agent Pools](#pipelinepermission-for-agent-pools)
- [GraphGroup](#graphgroup)
  - [Adding non-VS groups with custom scope (project level)](#adding-non-vs-groups-with-custom-scope-project-level)
  - [Setting multiple `groupDescriptors` for nested groups](#setting-multiple-groupdescriptors-for-nested-groups)
  - [AAD Groups not listed in Graph Groups API if not nested](#aad-groups-not-listed-in-graph-groups-api-if-not-nested) d

## GitRepository 

This section outlines the expected behaviors of GitRepositories under different configurations during creation and forking.
Each case states the input configuration and resulting behavior.

Some sample Custom Resources are provided in the [`/samples/gitrepository/` directory](../chart/samples/gitrepository/) of the chart. 
Files are named after some of the cases described below, e.g., [`gitrepository_2.3.yaml`](../chart/samples/gitrepository/gitrepository_2.3.yaml) corresponds to Case 2.3.

### 1. New Repository Cases

#### Case 1.1: New Repo Without Initialization
- **Input**:
  - `initialize`: *omitted*
  - `defaultBranch`: *omitted*
- **Result**: 
  - Repository is created **empty**, with **no branches**.

---

#### Case 1.2: New Repo With Initialization (Default Branch Omitted)
- **Input**:
  - `initialize`: `true`
  - `defaultBranch`: *omitted*
- **Result**:
  - Repository is initialized with a **first commit** on the `main` branch (which is set as `defaultBranch` as it is the default branch in Azure DevOps).

---

#### Case 1.3: New Repo With Initialization and Custom Default Branch
- **Input**:
  - `initialize`: `true`
  - `defaultBranch`: `test-branch`
- **Result**:
  - Repository is initialized with a **first commit** on the `test-branch` branch.

---

#### Case 1.4: New Repo With Custom Default Branch (Without Initialization Flag set to `true`)
- **Input**:
  - `initialize`: `false` or *omitted*
  - `defaultBranch`: `test-branch`
- **Result**:
  - Error 400: `When specifying a 'defaultBranch' for a new repository, 'initialize' must be set to true`

---

### 2. Forked Repository Cases (Parent Repository Initialized)

> In fork scenarios, the `initialize` field is **ignored**.

#### Case 2.1: Forked Repo With Nonexistent SourceRef
- **Preconditions**:
  - Parent repository is **initialized** but does **not** have a branch named `test-branch`.
- **Input**:
  - `defaultBranch`: anything or omitted
  - `sourceRef`: `refs/heads/test-branch-not-existing`
- **Result**:
  - **Error 400**: `SourceRef 'refs/heads/test-branch-not-existing'` does not exist in parent repository.

---

#### Case 2.2: Forked Repo With Existing Default Branch in Parent
- **Preconditions**:
  - Parent repository is **initialized** and has a branch named `test-branch` (either as default or non-default).
- **Input**:
  - `defaultBranch`: `test-branch` (which exists in parent repo)
  - `sourceRef`: *includes* `test-branch` or *omitted*
- **Result**:
  - Repository is forked with `test-branch` as the **default branch** (git history is copied). Note that in the parent repo, `test-branch` can be the default branch or any other branch.

---

#### Case 2.3: Forked Repo With Nonexistent Default Branch set (Fallback to Parent Default)
- **Preconditions**:
  - Parent repository is **initialized** and has a default branch (e.g., `main`) but does **not** have a branch named `test-branch`.
- **Input**:
  - `defaultBranch`: `test-branch` (does **not** exist in parent)
  - `sourceRef`: *omitted* or references an **existing** branch in the parent
- **Result**:
  - Repository is forked with the **default branch of the parent** repository (e.g., `main`) as the default branch (git history is copied).
  - Status code: **202 Accepted**
  - Notes: In the context of the `gitrepository-controller`, the controller awaits the manual creation of the `test-branch` by the user. 
  The Kubernetes resource will be in a `Ready: False` state until the user creates the `test-branch` on Azure DevOps.
  After the user manually creates the `test-branch` on the GitRepository (not in the parent), the `gitrepository-controller` will update the repository to set `test-branch` as the default branch. Finally, the GitRepository CR will become `Ready: True`.

---

### 3. Forked Repository Cases (Parent Repository Not Initialized)

Note: it is not recommended to fork from uninitialized repositories.

> In fork scenarios, the `initialize` field is **ignored**.

#### Case 3.1: Fork With Invalid `sourceRef`
- **Preconditions**:
  - Parent repository is **not initialized** (no branches exist).
- **Input**:
  - Parent repo is **not initialized**
  - `defaultBranch`: *omitted*
  - `sourceRef`: `test-branch-not-existing`
- **Result**:
  - **Error 400**: `SourceRef 'refs/heads/test-branch-not-existing'` does not exist in parent repository.

---

#### Case 3.2: Fork With No Branch Configuration
- **Input**:
  - Parent repo is **not initialized**
  - `defaultBranch`: *omitted*
  - `sourceRef`: *omitted*
- **Result**:
  - Repository is forked in an **uninitialized** state.
  - No branches are created.

---

#### Case 3.3: Fork With Nonexistent Default Branch
- **Input**:
  - Parent repo is **not initialized**
  - `defaultBranch`: `test-branch`
  - `sourceRef`: *omitted*
- **Result**:
  - Repository is forked in an **uninitialized** state.
  - No branches are created.
  - Status code: **202 Accepted**
  - Notes: In the context of the `gitrepository-controller`, the controller awaits the manual creation of the `test-branch` by the user. 
  The Kubernetes resource will be in a `Ready: False` state until the user creates the `test-branch` on Azure DevOps.
  After the user manually creates the `test-branch` on the GitRepository (not in the parent), the `gitrepository-controller` will update the repository to set `test-branch` as the default branch. Finally, the GitRepository CR will become `Ready: True`.

## PolicyConfiguration

To list all available policy types in your Azure DevOps organization, you can use the following `curl` command. Make sure to replace `<ORG>`, `<PROJECT>`, and `TOKEN`:
```sh
curl -X GET "https://dev.azure.com/<ORG>/<PROJECT>/_apis/policy/types?api-version=7.2-preview.1" \
-H "Authorization: Basic TOKEN"
```

```json
{
  "count": 17,
  "value": [
    {
      "description": "GitRepositorySettingsPolicyName",
      "id": "0517f88d-4ec5-4343-9d26-9930ebd53069",
      "displayName": "GitRepositorySettingsPolicyName"
    },
    {
      "description": "This policy will reject pushes to a repository for files that contain credentials or secrets.",
      "id": "ec003f37-8db0-4e10-992a-a2895045752c",
      "displayName": "Secrets scanning restriction"
    },
    {
      "description": "This policy will reject pushes to a repository for files that contain credentials or secrets.",
      "id": "90f9629b-664b-4804-a560-dd79b0c628f8",
      "displayName": "Secrets scanning restriction"
    },
    {
      "description": "This policy will reject pushes to a repository for paths which exceed the specified length.",
      "id": "001a79cf-fda1-4c4e-9e7c-bac40ee5ead8",
      "displayName": "Path Length restriction"
    },
    {
      "description": "This policy will require a successful proof of presence for the PR.",
      "id": "67ed70bd-2a6b-4006-af44-be590463f46d",
      "displayName": "Proof Of Presence"
    },
    {
      "description": "This policy will reject pushes to a repository for names which aren't valid on all supported client OSes.",
      "id": "db2b9b4c-180d-4529-9701-01541d19f36b",
      "displayName": "Reserved names restriction"
    },
    {
      "description": "This policy ensures that pull requests use a consistent merge strategy.",
      "id": "fa4e907d-c16b-4a4c-9dfa-4916e5d171ab",
      "displayName": "Require a merge strategy"
    },
    {
      "description": "Check if the pull request has any active comments",
      "id": "c6a1889d-b943-4856-b76f-9e46bb6b0df2",
      "displayName": "Comment requirements"
    },
    {
      "description": "This policy will require a successful status to be posted before updating protected refs.",
      "id": "cbdc66da-9728-4af8-aada-9a5a32e4a226",
      "displayName": "Status"
    },
    {
      "description": "Git repository settings",
      "id": "7ed39669-655c-494e-b4a0-a08b4da0fcce",
      "displayName": "Git repository settings"
    },
    {
      "description": "This policy will require a successful build has been performed before updating protected refs.",
      "id": "0609b952-1397-4640-95ec-e00a01b2c241",
      "displayName": "Build"
    },
    {
      "description": "This policy will reject pushes to a repository for files which exceed the specified size.",
      "id": "2e26e725-8201-4edd-8bf5-978563c34a80",
      "displayName": "File size restriction"
    },
    {
      "description": "This policy will reject pushes to a repository which add file paths that match the specified patterns.",
      "id": "51c78909-e838-41a2-9496-c647091e3c61",
      "displayName": "File name restriction"
    },
    {
      "description": "This policy will block pushes from including commits where the author email does not match the specified patterns.",
      "id": "77ed4bd3-b063-4689-934a-175e4d0a78d7",
      "displayName": "Commit author email validation"
    },
    {
      "description": "This policy will ensure that required reviewers are added for modified files matching specified patterns.",
      "id": "fd2167ab-b0be-447a-8ec8-39368250530e",
      "displayName": "Required reviewers"
    },
    {
      "description": "This policy will ensure that a minimum number of reviewers have approved a pull request before completion.",
      "id": "fa4e907d-c16b-4a4c-9dfa-4906e5d171dd",
      "displayName": "Minimum number of reviewers"
    },
    {
      "description": "This policy encourages developers to link commits to work items.",
      "id": "40e92b44-2fe1-4dd6-b3d8-74a9c21d0c6e",
      "displayName": "Work item linking"
    }
  ]
}
```

## Team

Unlike other Azure DevOps resources, the field `projectId` in the Team resource can only accept the **Project ID** (a UUID) and not the Project Name.
This is due to the fact that the Azure DevOps REST API for Teams returns always the Project ID in the response, even if the Team was created using the Project Name. This difference would lead to infinite reconciliation loops in the controller if the Project Name was used.
For instance, if you create a Team with `projectId: MyProject`, the API will return `projectId: 850e8400-e29b-41d4-a716-446655440000` (the actual Project ID), causing the controller to think that the resource is out of sync and attempt to update it again with `projectId: MyProject`, leading to an infinite loop.

## PipelinePermission

### PipelinePermission for Agent Pools

In the case of managing a `PipelinePermission` for an Agent Pool, the `resourceType` should be set to `queue`. [Reference on StackOverflow](https://stackoverflow.com/questions/77258168/how-to-update-azure-pipeline-permissions-for-resource-using-cli#comment139630674_77258662).

## GraphGroup

Plugin to be implemented:

if scopeDescriptor is not set
it means the group searched is at organization level
so we want to filter out all project level groups
so if they have
"domain": "vstfs:///Classification/TeamProject/
then the proxied response should filter them out

if instead, scopeDescriptor is set, we are already scoping to a project (native behavior of the AzureDevops API)
and so we don't need to filter anything

note also that project level groups cannot be AAD groups (materialized by mail or originId), so no need to worry about that case

### Adding non-VS groups with custom scope (project level)

When trying to add a non-Visual Studio created group (e.g., an AAD group) to a custom scope (e.g., a project level scope) in Azure DevOps, an error is encountered.

Reference error:
```json
{
  "$id": "1",
  "errorCode": 0,
  "eventId": 3000,
  "innerException": null,
  "message": "VS860004: Only Visual Studio created groups can be added to a custom scope.",
  "typeKey": "GraphBadRequestException",
  "typeName": "Microsoft.VisualStudio.Services.Graph.GraphBadRequestException, Microsoft.VisualStudio.Services.WebApi"
}
```

### Setting multiple `groupDescriptors` for nested groups

An error could be experienced when trying to set more than one `groupDescriptor` for nested groups in Azure DevOps.
Reference error:
```json
{
  "$id": "1",
  "errorCode": 0,
  "eventId": 0,
  "innerException": null,
  "message": "Given identifier length is out of range of valid values.\r\nParameter name: identifier\r\nActual value was vssgp.Uy0xLTktMTU1MTM3NDI0NS0yMzgyNjQzNzU1LTExODE2NjQzMjMtMzEwNTYzNjM2Ni00MjA0NjczODI2LTEtNjU2MTUwMTQzLTM3NTk4OTk3MTItMjIxMjU1NjYyMi0yMzI3NjU1NTkz,vssgp.Uy0xLTktMTU1MTM3NDI0NS0yMzgyNjQzNzU1LTExODE2NjQzMjMtMzEwNTYzNjM2Ni00MjA0NjczODI2LTEtMTAwOTI3ODYzMi0zMDcyMjM3MTI2LTMxNDYxODc2MDEtMjMxOTM2MjE4.",
  "typeKey": "ArgumentOutOfRangeException",
  "typeName": "System.ArgumentOutOfRangeException, mscorlib"
}
```

It is not known if this behavior can be experienced only in certain scenarios or every time multiple `groupDescriptors` are set.
Currently, no solution has been found to this issue.

### AAD Groups not listed in Graph Groups API if not nested

If a AAD group is materialized (by mail or originId) into Azure DevOps, **without being nested under another group** using the `groupDescriptors` field, it will not be listed in the response of the API for listing groups.

An example scenario:
- `GET /{organization}/_apis/graph/groups` does not list the group
- but `GET /{organization}/_apis/graph/groups/{groupDescriptor}` can fetch it
- and it can be deleted too using `DELETE /{organization}/_apis/graph/groups/{groupDescriptor}`

This can also be verified in the Azure DevOps portal when navigating to the Organization Settings -> Permissions -> Groups page: https://dev.azure.com/krateo-kog/_settings/groups.

However, the group can still be accessed directly if you know its descriptor, for example: https://dev.azure.com/krateo-kog/_settings/groups?subjectDescriptor=aadgp.AAABBBCCC

Currently, no solution has been found to this issue.


{
  "$id": "1",
  "errorCode": 0,
  "eventId": 4207,
  "innerException": null,
  "message": "TF50258: An error occurred finding the group. There is no group with the security identifier (SID) S-1-9-1551374245-1204400969-2402986413-2179408616-3-1968947946-954710338-3207406156-2488723893.",
  "typeKey": "FindGroupSidDoesNotExistException",
  "typeName": "Microsoft.VisualStudio.Services.Identity.FindGroupSidDoesNotExistException, Microsoft.VisualStudio.Services.WebApi"
}