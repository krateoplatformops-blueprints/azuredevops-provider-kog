# Azure DevOps Provider KOG - Plugins

This README is only related to the plugins contained in this repository. For more information about the Azure DevOps Provider KOG (Krateo Operator Generator), please refer to the [main README](../README.md).

This repository contains the source code of a set of plugins written in Go for the Azure DevOps Provider KOG.

These plugins are specialized web services that address some inconsistencies in the Azure DevOps REST API or provide a way to use Azure DevOps API seamlessly within KOG (Krateo Operator Generator). This may be due differences in request/response schemas, missing fields, or additional validations needed when managing Azure DevOps resources. 

Sometimes, since the REST API is not originally designed to be used by a Kubernetes operator, these plugins fill the gap by providing the necessary functionality to manage Azure DevOps resources effectively through KOG.
On the other hand, some missing features of Krateo Operator Generator are addressed by these plugins as well.

In addition, some resources managed by Azure DevOps Provider KOG are **abstractions** (e.g, `RepositoryPermission`, `BuildPermission`) that do not have a direct counterpart in the Azure DevOps REST API. In these cases, the plugins implement the necessary logic to map the Kubernetes Custom Resource specifications to the appropriate Azure DevOps API calls.

They are designed to work as a middleware between the [Rest Dynamic Controller](https://github.com/krateoplatformops/rest-dynamic-controller/) and the Azure DevOps REST API.

## Summary

- [Project structure](#project-structure)
- [Architecture](#architecture)
- [GitRepository plugin](#gitrepository-plugin)
  - [POST GitRepository](#post-gitrepository)
- [Pipeline plugin](#pipeline-plugin)
  - [GET Pipeline](#get-pipeline)
  - [UPDATE Pipeline](#update-pipeline)
  - [DELETE Pipeline](#delete-pipeline)
- [PipelinePermission plugin](#pipelinepermission-plugin)
  - [GET PipelinePermission](#get-pipelinepermission)
- [PullRequest plugin](#pullrequest-plugin)
  - [PATCH PullRequest](#patch-pullrequest)
- [Azure DevOps API Reference](#azuredevops-api-reference)
- [Authentication](#authentication)
- [API documentation](#api-documentation)
- [Development and Testing Guide](#development-and-testing-guide)

## Project structure

The project is organized as a Go workspace-based monorepo containing multiple, independent plugins. Each plugin resides in its own subdirectory under the `cmd/` directory, while shared code is located in the `pkg/` directory.
Each plugin is a standalone Go module, allowing for independent building and versioning.
Therefore, each plugin can be built into its own container image. Please refer to the [release process section of the main readme](../README.md#release-process) for more details.

## Architecture

The diagram below illustrates how the Azure DevOps Plugin interacts within the Krateo ecosystem, including the `rest-dynamic-controller` and the `azuredevops-provider-kog` as well as the Azure DevOps REST API.
Note that this is only a **partial representation** showing only 3 managed resources as an example (`Pipeline`, `PipelinePermission`, and `GitRepository`).

```mermaid
graph TD
    %% Define node order early to influence layout
    KOG -- Instanciates --> RDC1
    KOG -- Instanciates --> RDC2
    KOG -- Instanciates --> RDC3

    RDC1 --> Pipeline
    RDC1 -- Direct calls --> ADO_API

    RDC2 --> PipelinePermission
    RDC2 -- Direct calls --> ADO_API
    
    
    RDC3 -- Direct calls --> ADO_API
    RDC3 --> GitRepository
    
    Pipeline -- "Fixes folder path,<br>consistent update/delete" --> ADO_API
    PipelinePermission -- "Ensures allPipelines<br>property is set" --> ADO_API
    GitRepository -- "Supports defaultBranch,<br>initialization, fork validation" --> ADO_API

    subgraph Krateo Ecosystem
        KOG[azuredevops-provider-kog <br> Helm Chart]
        RDC1[rest-dynamic-controller-1 <br> Pipeline]
        RDC2[rest-dynamic-controller-2 <br> PipelinePermission]
        RDC3[rest-dynamic-controller-3 <br> GitRepository]

    end

    Pipeline[Pipeline<br>Plugin]
    PipelinePermission[PipelinePermission<br>Plugin]
    GitRepository[GitRepository<br>Plugin]
    
    subgraph Azure DevOps
        ADO_API[Azure DevOps REST API]
    end

    %% Styling
    style RDC1 fill:#e0f2f7,stroke:#333,stroke-width:1px
    style RDC2 fill:#e0f2f7,stroke:#333,stroke-width:1px
    style RDC3 fill:#e0f2f7,stroke:#333,stroke-width:1px
    style KOG fill:#e0f2f7,stroke:#333,stroke-width:1px
    style ADO_API fill:#f0f0f0,stroke:#333,stroke-width:1px
    style Pipeline fill:#d4f8d4,stroke:#333,stroke-width:1px
    style PipelinePermission fill:#d4f8d4,stroke:#333,stroke-width:1px
    style GitRepository fill:#d4f8d4,stroke:#333,stroke-width:1px

```

## GitRepository plugin

### POST GitRepository

**Description**:
This endpoint creates a new GitRepository in the specified Azure DevOps project.
It allows you to specify the `initialize` field to indicate whether the repository should be initialized with a first commit. (Note: you cannot initialize a repository with a first commit if you are forking a repository).
It allows you to specify the `defaultBranch` field to set the default branch of the repository.

<details>
<summary><b>Why This Endpoint Exists</b></summary>
<br/>

- The standard Azure DevOps REST API has two different request body schemas for creating (`POST`) and updating (`PATCH`) Git repositories. In particular, the field `defaultBranch` is only available in the `PATCH` request body.
- This endpoint allows you to create a Git repository with the `defaultBranch` field, which is not supported in the standard Azure DevOps REST API for the `POST` request body. Practically performing a `PATCH` operation on the repository immediately after creation.
- Moreover, it allows you to initialize the repository with a first commit by setting the `initialize` field to `true`.
- In addition, it performs additional validations related to branch existence (for forks) and repository initialization.
- Another additional validation is that it checks if the `sourceRef` branch exists in the parent repository when forking a repository. If it does not exist, it returns a `400 Bad Request` error.

</details>

<details>
<summary><b>Request</b></summary>
<br/>

```http
POST /api/{organization}/{projectId}/git/repositories
```

**Path parameters**:
- `organization` (string, required): The name of the Azure DevOps organization.
- `projectId` (string, required): The ID or name of the Azure DevOps project.

**Query parameters**:
- `api-version` (string, required): The version of the Azure DevOps REST API to use. For example, `7.2-preview.2`.
- `sourceRef` (string, optional): The source reference for the repository. This is typically a branch name (e.g., `refs/heads/main`).

**Request body example**:
```json
{
  "name": "string",
  "defaultBranch": "string",    // Adjusted field
  "initialize": true,           // Adjusted field

  // From here, optional, fork-related fields:
  "parentRepository": {
    "id": "4b8c6f64-5717-4562-b3fc-2c963f66afa6",
    "project": {
      "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    }
  },
  "project": {
    "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  }
}
```

> The field `projectId` (path parameter) can be either the project ID or the project name. The fields `project.id` and `parentRepository.project.id` in the request body must be the project ID (not the project name) and are required when forking a repository. If you are not forking a repository, you have to omit these fields.

</details>

<details>
<summary><b>Response</b></summary>
<br/>

**Response status codes**:
- `201 Created`: The GitRrepository was successfully created.
- `202 Accepted`: The GitRrepository was successfully created but `defaultBranch` specified in the request body does not exist in the repository.
- `400 Bad Request`: The request body is invalid, the `sourceRef` branch does not exist in the parent repository or other validation errors occurred.
- `401 Unauthorized`: The request is not authorized. Ensure that the `Authorization` header is set correctly.
- `500 Internal Server Error`: An unexpected error occurred while processing the request.

**Response body example**:
```json
{
  "_links": {
    "links": {
      "additionalProp1": {},
      "additionalProp2": {},
      "additionalProp3": {}
    }
  },
  "creationDate": "2025-07-06T12:28:03.454Z",
  "defaultBranch": "string",
  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "isDisabled": true,
  "isFork": true,
  "isInMaintenance": true,
  "name": "string",
  "parentRepository": {
    "collection": {
      "avatarUrl": "string",
      "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      "name": "string",
      "url": "string"
    },
    "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "isFork": true,
    "name": "string",
    "project": {
      "abbreviation": "string",
      "defaultTeamImageUrl": "string",
      "description": "string",
      "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
      "lastUpdateTime": "2025-07-06T12:28:03.454Z",
      "name": "string",
      "revision": 0,
      "state": "deleting",
      "url": "string",
      "visibility": "private"
    },
    "remoteUrl": "string",
    "sshUrl": "string",
    "url": "string"
  },
  "project": {
    "abbreviation": "string",
    "defaultTeamImageUrl": "string",
    "description": "string",
    "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "lastUpdateTime": "2025-07-06T12:28:03.454Z",
    "name": "string",
    "revision": 0,
    "state": "deleting",
    "url": "string",
    "visibility": "private"
  },
  "remoteUrl": "string",
  "size": 0,
  "sshUrl": "string",
  "url": "string",
  "validRemoteUrls": [
    "string"
  ],
  "webUrl": "string"
}
```

</details>

## Pipeline plugin

### GET Pipeline

**Description**:
This endpoint retrieves a specific pipeline by its ID in the specified Azure DevOps project.
It returns the pipeline details, including its ID, name, and other metadata.

<details>
<summary><b>Why This Endpoint Exists</b></summary>
<br/>

- The standard Azure DevOps REST API return the `folder` field with an "escaped backslash" as prefix like `"folder":"\\test-folder"`.
- This endpoint returns the `folder` field without the "escaped backslash" prefix, allowing a correct comparison with the `folder` field set in the `spec` of the `Pipeline` resource. Otherwise, the reconciliation loop in KOG would always detect a difference, for example, between `"\\test-folder"` (from Azure DevOps REST API) and `"test-folder"` (from `spec`), leading to infinite drift detections and useless updates.

</details>

<details>
<summary><b>Request</b></summary>
<br/>

```http
GET /api/{organization}/{project}/pipelines/{id}
```

**Path parameters**:
- `organization` (string, required): The name of the Azure DevOps organization.
- `project` (string, required): The name of the Azure DevOps project.
- `id` (string, required): The ID of the pipeline to retrieve.

**Query parameters**:
- `api-version` (string, required): The version of the Azure DevOps REST API to use. For example, `7.2-preview.1`.

</details>

<details>
<summary><b>Response</b></summary>
<br/>

**Response status codes**:
- `200 OK`: The request was successful and the pipeline details are returned.
- `400 Bad Request`: The request is invalid. Ensure that the `organization`, `project`, and `id` parameters are correct.
- `401 Unauthorized`: The request is not authorized. Ensure that the `Authorization` header is set correctly.
- `500 Internal Server Error`: An unexpected error occurred while processing the request.

**Response body example**:
```json
{
  "_links":{
    "self":{
      "href":"string"
      },
    "web":{
      "href":"string"
    }
  },
  "configuration":{
    "path":"pipelines/test_inner_pipeline.yml",
    "repository":{
      "id":"string",
      "type":"azureReposGit"
    },
    "type":"yaml"
  },
  "folder":"test-folder-kog", // Adjusted field, from "\\test-folder-kog" to "test-folder-kog"
  "id":49,
  "name":"test-pipeline-kog-1",
  "revision":1,
  "url":"string"
}
```

</details>

### UPDATE Pipeline

**Description**:
This endpoint updates an existing pipeline in the specified Azure DevOps project.
In particular, it allows you to change the pipeline's name, folder, and configuration details such as the path to the configuration file.

<details>
<summary><b>Why This Endpoint Exists</b></summary>
<br/>

- The standard Azure DevOps REST API does not have a `/pipelines/{id}` endpoint for updating pipelines.
- In order to update a pipeline, you need to use the `/build/definitions/{id}` endpoint, which is not consistent with the `/pipelines/{id}` endpoint used for retrieving pipelines.
- This endpoint provides a consistent way to update pipelines using the `/pipelines/{id}` endpoint and the same request body schema as the `POST /pipelines` endpoint of Azure DevOps REST API.
- In particular, the plugin creates a `BuildDefinitionMinimal` object starting from the request body and then performs a `PUT` request to the `/build/definitions/{id}` endpoint of Azure DevOps REST API.
- A needed adjustement related to the repository type is performed, as the Azure DevOps REST API returns different values for the `repository.type` field depending on the endpoint used to retrieve the pipeline. For instance, even if a pipeline is linked to a `azureReposGit` repository, the `/build/definitions/{id}` endpoint returns `repository.type` as `TfsGit`, while the `/pipelines/{id}` endpoint returns `repository.type` as `azureReposGit`.
- Moreover, since this endpoint under the hood uses the `/build/definitions/{id}` Azure DevOps endpoint, the plugin set the correct `api-version` parameter needed to update a pipeline using the `/build/definitions/{id}` endpoint (`7.2-preview.7`).

> Currently, the `api-version` parameter is passed as an environment variable to the plugin by the related Helm chart.

</details>

<details><summary><b>Request</b></summary>
<br/>

```http
PUT /api/{organization}/{project}/pipelines/{id}
```

**Path parameters**:
- `organization` (string, required): The name of the Azure DevOps organization.
- `project` (string, required): The name of the Azure DevOps project.
- `id` (string, required): The ID of the pipeline to update.

**Request body example**:
```json
{
  "configuration":{
    "path":"pipelines/inner_folder/another_config.yml",
    "repository":{
      "id":"string",
      "type":"azureReposGit"
    },
    "type":"yaml"
  },
  "folder":"test-folder-kog",
  "name":"test-pipeline-kog-1-v2",
  "revision":"3"
}
```

</details>

<details><summary><b>Response</b></summary>
<br/>

**Response status codes**:
- `200 OK`: The pipeline was successfully updated.
- `400 Bad Request`: The request body is invalid or the pipeline ID does not exist.
- `401 Unauthorized`: The request is not authorized. Ensure that the `Authorization` header is set correctly.
- `404 Not Found`: The specified pipeline does not exist in the project.
- `500 Internal Server Error`: An unexpected error occurred while processing the request.

**Response body example**:
```json
{
  "_links":{
    "self":{
      "href":"string"
      },
    "web":{
      "href":"string"
    }
  },
  "configuration":{
    "path":"pipelines/test_inner_pipeline.yml",
    "repository":{
      "id":"string",
      "type":"azureReposGit" // Adjusted field
    },
    "type":"yaml"
  },
  "folder":"test-folder-kog", // Adjusted field
  "id":49,
  "name":"test-pipeline-kog-1",
  "revision":1,
  "url":"string"
}
```

</details>

### DELETE Pipeline

**Description**:
This endpoint deletes a specific pipeline by its ID in the specified Azure DevOps project.

<details>
<summary><b>Why This Endpoint Exists</b></summary>
<br/>

- The standard Azure DevOps REST API does not have a `/pipelines/{id}` endpoint for deleting pipelines.
- In order to delete a pipeline, you need to use the `/build/definitions/{id}` endpoint, which currently support a different `api-version` parameter when compared to the `/pipelines/{id}` endpoint used for retrieving pipelines.
- This endpoint sets the correct `api-version` parameter needed to delete a pipeline using the `/build/definitions/{id}` endpoint (`7.2-preview.7`).

> Currently, the `api-version` parameter for BuildDefinitions is passed as an environment variable to the plugin by the related Helm chart.

</details>

<details><summary><b>Request</b></summary>
<br/>

```http
DELETE /api/{organization}/{project}/pipelines/{id}
```

**Path parameters**:
- `organization` (string, required): The name of the Azure DevOps organization.
- `project` (string, required): The name of the Azure DevOps project.
- `id` (string, required): The ID of the pipeline to delete.

</details>

<details><summary><b>Response</b></summary>
<br/>

**Response status codes**:
- `204 No Content`: The pipeline was successfully deleted.
- `400 Bad Request`: The request is invalid or the pipeline ID does not exist.
- `401 Unauthorized`: The request is not authorized. Ensure that the `Authorization` header is set correctly.
- `404 Not Found`: The specified pipeline does not exist in the project.
- `500 Internal Server Error`: An unexpected error occurred while processing the request.

</details>

---

## PipelinePermission plugin

### GET PipelinePermission

**Description**: 
Given a `ResourceType` and `ResourceId`, it returns authorized definitions for that resource.
More precisely, it returns the list of `pipelines` that have permissions to access the specified resource and the fact whether `allPipelines` have access to it.

<details>
<summary><b>Why This Endpoint Exists</b></summary>
<br/>

- The standard Azure DevOps REST API **does not return** the `allPipelines` property when said property is set to `authorized: false` on Azure DevOps (default behavior).
- This endpoint checks if the response from the Azure DevOps REST API contains the `allPipelines` property and, if not, it adds it with a value of `authorized: false`. This is necessary to have a correct reconciliation loop in KOG, as the absence of the `allPipelines` property in the response would determine an incomplete comparison with the `spec` of the `PipelinePermission` resource.

</details>

<details>
<summary><b>Request</b></summary>
<br/>

```http
GET /api/{organization}/{project}/pipelines/pipelinepermissions/{resourceType}/{resourceId}
```

**Path parameters**:
- `organization` (string, required): The name of the Azure DevOps organization.
- `project` (string, required): The name of the Azure DevOps project.
- `resourceType` (string, required): The type of resource for which permissions are being requested (e.g., `repository`, `environment`, `queue`).
- `resourceId` (string, required): The ID of the resource for which permissions are being requested.

**Query parameters**:
- `api-version` (string, required): The version of the Azure DevOps REST API to use. For example, `7.2-preview.2`.
</details>

<details>
<summary><b>Response</b></summary>
<br/>

**Response status codes**:
- `200 OK`: The request was successful.
- `401 Unauthorized`: The request is not authorized. Ensure that the `Authorization` header is set correctly.
- `500 Internal Server Error`: An unexpected error occurred while processing the request.

**Response body example**:
```json
{
  "resource": {
    "type":"environment",
    "id":"7"
  },
  "allPipelines":{
    "authorized":false // Adjusted field, added if missing in the response from Azure DevOps REST API
  },
  "pipelines": [
    {
      "id":14,
      "authorized":true,
      "authorizedBy": {
        "displayName":"<REDACTED>",
        "id":"<REDACTED>",
        "uniqueName":"<REDACTED>",
        "descriptor":"<REDACTED>"
      },
      "authorizedOn":"2025-06-30T14:33:02.06Z"
    },
    {
      "id":15,
      "authorized":true,
      "authorizedBy": {
        "displayName":"<REDACTED>",
        "id":"<REDACTED>",
        "uniqueName":"<REDACTED>",
        "descriptor":"<REDACTED>"
      },
      "authorizedOn":"2025-06-30T14:33:02.06Z"
    }
  ]
}
```
</details>

---

## PullRequest plugin

### PATCH PullRequest

**Description**:
This endpoint updates an existing pull request in the specified Azure DevOps project.

<details>
<summary><b>Why This Endpoint Exists</b></summary>
<br/>

- The endpoint calculates only the modified fields between the current state of the pull request in Azure DevOps and the desired state specified in the request body (coming from the Kubernetes CR), and it sends only those modified fields in the `PATCH` request to the Azure DevOps REST API. Indeed, Azure DevOps throws an error if you try to set a field to its current value, for instance in the case of `TargetRefName` field.
- In addition, it manages `LastMergeSourceCommit` field which is required only when updating status to "completed": if not set, Azure DevOps returns an error. If set in other cases, Azure DevOps returns an error.
- It also performs additional validations related to pull request state transitions: an already completed pull request cannot be updated.

</details>

<details>
<summary><b>Request</b></summary>
<br/>

```http
PATCH /api/{organization}/{project}/git/repositories/{repositoryId}/pullrequests/{pullRequestId}
```

**Path parameters**:
- `organization` (string, required):T he name of the Azure DevOps organization.
- `project` (string, required): Project ID or project name.
- `repositoryId` (string, required):  The repository ID of the pull request's target branch.
- `pullRequestId` (string, required): The ID of the pull request to update

**Query parameters**:
- `api-version` (string, required): The version of the Azure DevOps REST API to use. For example, `7.2-preview.2`.

**Request body example**:
```json
{
  "completionOptions": {
    "deleteSourceBranch": false,
    "mergeCommitMessage": "Squash merge for feature X, test from Krateo",
    "mergeStrategy": "squash",
    "transitionWorkItems": true
  },
  "description": "Test PR description from Krateo",
  "lastMergeSourceCommit": {
    "commitId": "<REDACTED>",
    "url": "<REDACTED>"
  },
  "mergeOptions": {
    "conflictAuthorshipCommits": true,
    "detectRenameFalsePositives": false,
    "disableRenames": false
  },
  "status": "active",
  "targetRefName": "refs/heads/main",
  "title": "Test PR from Krateo test-branch"
}
```

</details>

<details>
<summary><b>Response</b></summary>
<br/>

**Response status codes**:
- `200 OK`: The pull request was successfully updated.
- `400 Bad Request`: The request body is invalid or the pull request ID does not exist.
- `401 Unauthorized`: The request is not authorized. Ensure that the `Authorization` header is set correctly.
- `404 Not Found`: The specified pull request does not exist in the repository.
- `409 Conflict`: The pull request cannot be updated due to a conflict (e.g., trying to update a completed pull request).
- `500 Internal Server Error`: An unexpected error occurred while processing the request.

**Response body example**:
```json
{
  "repository": {
    "id": "5605b0ba-e2fa-4aab-af0b-0888321b3a08",
    "name": "repo-1-classic",
    "url": "<REDACTED>",
    "project": {
      "id": "99837031-4e4e-4753-9a47-73fcc4cba766",
      "name": "project-1-classic",
      "description": "test",
      "url": "<REDACTED>",
      "state": "wellFormed",
      "revision": 626,
      "visibility": "private",
      "lastUpdateTime": "2025-12-03T08:38:51.013Z"
    },
    "size": 10828,
    "remoteUrl": "<REDACTED>",
    "sshUrl": "<REDACTED>",
    "webUrl": "<REDACTED>",
    "isDisabled": false,
    "isInMaintenance": false
  },
  "pullRequestId": 402,
  "codeReviewId": 402,
  "status": "active",
  "createdBy": {
    "displayName": "Leonardo Vicentini",
    "url": "<REDACTED>",
    "_links": {
      "avatar": {
        "href": "<REDACTED>"
      }
    },
    "id": "<REDACTED>",
    "uniqueName": "leonardo.vicentini@krateo.io",
    "imageUrl": "<REDACTED>",
    "descriptor": "<REDACTED>"
  },
  "creationDate": "2025-12-03T14:44:56.2184364Z",
  "title": "Test PR from Krateo test-branch",
  "description": "Test PR description from Krateo",
  "sourceRefName": "refs/heads/test-3-12-25",
  "targetRefName": "refs/heads/main",
  "mergeStatus": "succeeded",
  "isDraft": false,
  "mergeId": "<REDACTED>",
  "lastMergeSourceCommit": {
    "commitId": "<REDACTED>",
    "url": "<REDACTED>"
  },
  "lastMergeTargetCommit": {
    "commitId": "<REDACTED>",
    "url": "<REDACTED>"
  },
  "lastMergeCommit": {
    "commitId": "<REDACTED>",
    "author": {
      "name": "Leonardo Vicentini",
      "email": "leonardo.vicentini@krateo.io",
      "date": "2025-11-17T15:08:26Z"
    },
    "committer": {
      "name": "Leonardo Vicentini",
      "email": "leonardo.vicentini@krateo.io",
      "date": "2025-11-17T15:08:26Z"
    },
    "comment": "Squash merge for feature X, test from Krateo",
    "url": "<REDACTED>"
  },
  "reviewers": [],
  "url": "<REDACTED>",
  "_links": {
    "self": {
      "href": "<REDACTED>"
    },
    "repository": {
      "href": "<REDACTED>"
    },
    "workItems": {
      "href": "<REDACTED>"
    },
    "sourceBranch": {
      "href": "<REDACTED>"
    },
    "targetBranch": {
      "href": "<REDACTED>"
    },
    "statuses": {
      "href": "<REDACTED>"
    },
    "sourceCommit": {
      "href": "<REDACTED>"
    },
    "targetCommit": {
      "href": "<REDACTED>"
    },
    "createdBy": {
      "href": "<REDACTED>"
    },
    "iterations": {
      "href": "<REDACTED>"
    }
  },
  "completionOptions": {
    "mergeCommitMessage": "Squash merge for feature X, test from Krateo",
    "squashMerge": true,
    "mergeStrategy": "squash",
    "transitionWorkItems": true
  },
  "mergeOptions": {
    "disableRenames": false,
    "detectRenameFalsePositives": false,
    "conflictAuthorshipCommits": true
  },
  "supportsIterations": true,
  "artifactId": "<REDACTED>"
}
```

</details>

## Azure DevOps API Reference

For complete Azure DevOps REST API documentation, visit: [Azure DevOps REST API docs](https://learn.microsoft.com/en-us/rest/api/azure/devops/) and [API Specifications](https://github.com/MicrosoftDocs/vsts-rest-api-specs/tree/master).

## Authentication

The plugin will forward the `Authorization` header passed in the request to this plugin to the Azure DevOps REST API.
In particular, it supports the Basic Authentication scheme, which is the default for Azure DevOps REST API.
How it works:
- You can generate a Personal Access Token (PAT) in Azure DevOps.
- Use the PAT as the password in the Basic Authentication header.
- The username can be any string (e.g., `user`), as Azure DevOps does not require a specific username for PAT authentication.

You can get more information in the README of the [Azure DevOps Provider KOG](https://github.com/krateoplatformops-blueprints/azuredevops-provider-kog?tab=readme-ov-file#authentication).

## API documentation

Each plugin generates and holds its own OpenAPI specification. The documentation is generated using the `swag` tool and is stored within each plugin's directory (e.g., `cmd/collaborator-plugin/docs`).

To generate or update the documentation for a specific plugin, run the `swag-init.sh` script from this `plugins` directory, passing the plugin's name as an argument.

**Example: Generate docs for `pipeline-plugin`**
```sh
./scripts/swag-init.sh pipeline-plugin
```

This will generate the necessary `swagger.json`, `swagger.yaml`, and OpenAPI v3 files in the `cmd/pipeline-plugin/docs/` directory.

You can then access the Swagger UI for each plugin at `/swagger/index.html` endpoint of the respective plugin's service.

## Development and Testing Guide

For detailed instructions on building and testing the plugins, please refer to the [development and testing guide](./docs/dev_and_testing.md).
