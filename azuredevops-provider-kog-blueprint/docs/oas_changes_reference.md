# Azure DevOps Provider KOG - OpenAPI Specification (OAS) Changes reference

This documents serves as a reference for changes made to the OpenAPI Specification (OAS) of the resources managed by the Azure DevOps Provider KOG.
Note that the changes are made to **comply with some requirements** of the OASGen Provider and/or Rest Dynamic Controller or to **fix issues or inconsistencies** within the Azure DevOps REST API.

## Summary

- [GitRepository](#gitrepository)
- [Pipeline](#pipeline)
- [PipelinePermission](#pipelinepermission)
- [PullRequest](#pullrequest)
- [PolicyConfiguration](#policyconfiguration)
- [RepositoryPermission](#repositorypermission)
- [BuildPermission](#buildpermission)
- [Team](#team)

## `GitRepository`

**Version**: 7.2-preview.2
**Original specification file**:
- Link: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/master/specification/git/7.2/git.json
- Permalink: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/ed36f85e1796b78f1c88961a45396e31c6618000/specification/git/7.2/git.json

**Transformations**: 
- Original file is converted from JSON to YAML and from Swagger 2.0 to OpenAPI 3.0.1 with the following tool: https://editor.swagger.io/
- File is pruned to only include the endpoints and schemas relevant to GitRepositories.

**List of change made to the OpenAPI Specification (OAS)**:
- The path parameter `project` has been changed to `projectId` on every endpoint that requires it.
Otherwise there would be a potential clash between `project` field in path and `project` field in the request/response body (which is a "fork-related" field, used when forking a repository) that would cause issues with Rest Dynamic Controller operations due to ambiguity.
Note that `projectId` could be either a project name or a project ID even if the field is named `projectId`.
- The path parameter `repositoryId` has been changed to `id` on every endpoint that requires it. This is done to be aligned with the naming convention used in response bodies by the Azure DevOps REST API.
- In the `delete` endpoint the response status code has been changed from `200` to `204` as the Azure DevOps REST API actually returns a `204 No Content` status code when a repository is deleted successfully.
- Commented out the query parameters: `includeLinks`, `includeAllUrls`, `includeHidden`.
- Some addtional schemas are added, such as `GitRepositoryUpdateOptions`. This schema is used in the `update` operation to allow updating only the name and default branch of a GitRepository.
- Commented out the `oauth2` security scheme since it not supported by OASGen Provider and Rest Dynamic Controller while authentication is handled at the root level of the file and with basic auth scheme.
- Commented out `x-ms-*` items as they are specific to Microsoft internal tooling.
- Commented out the `components.parameters` section at root level as it just to inform about api versions.

## `Pipeline`

**Version**: 7.2-preview.1
**Original specification file**:
- Link: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/master/specification/pipelines/7.2/pipelines.json
- Permalink: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/a69e0a84db58a99dc4243957289b6d825dcb2af0/specification/pipelines/7.2/pipelines.json

**Transformations**: 
- Original file is converted from JSON to YAML and from Swagger 2.0 to OpenAPI 3.0.1 with the following tool: https://editor.swagger.io/
- File is pruned to only include the endpoints and schemas relevant to Pipelines.

**List of change made to the OpenAPI Specification (OAS)**:
- The schema for the request body of the `create` operation has been modified to include additional fields not documented in the original OpenAPI Specification (OAS) but required for a successful operation (`configuration.repository` and `configuration.path`).
- Note: Build Definitions APIs are used, via a plugin, in order to perform `update` and `delete` operations which are not available for Pipelines. Build Definitions are an older version of Pipelines and have a superset of the features of Pipelines and it is currently used with version 7.2-preview.7.
- Commented out `name` field in the `BuildRepository` schema to allow only setting the `id` field when specifying the repository of a Pipeline.
- Commented out the `oauth2` security scheme since it not supported by OASGen Provider and Rest Dynamic Controller while authentication is handled at the root level of the file and with basic auth scheme.
- Commented out `x-ms-*` items as they are specific to Microsoft internal tooling.
- Commented out the `components.parameters` section at root level as it just to inform about api versions.

## `PipelinePermission`

**Version**: 7.2-preview.1
**Original specification file**:
- Link: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/master/specification/approvalsAndChecks/7.2/pipelinePermissions.json
- Permalink: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/a69e0a84db58a99dc4243957289b6d825dcb2af0/specification/approvalsAndChecks/7.2/pipelinePermissions.json

**Transformations**: 
- Original file is converted from JSON to YAML and from Swagger 2.0 to OpenAPI 3.0.1 with the following tool: https://editor.swagger.io/
- File is pruned to only include the endpoints and schemas relevant to PipelinePermissions.

**List of change made to the OpenAPI Specification (OAS)**:
- The `RequestBody` schema is defined for the request body of the `PATCH` operation.
This is done to reduce the fields of the request body.
The PATCH operation is used in the RestDefinition `pipelinepermission` for the `create` and `update` operations.
- An `enum` has been added to the `resourceType` field for both the POST and PATCH operations, to restrict its possible values to the supported resource types (repository, environment, queue, teamproject, endpoint, variablegroup, securefile).
- Commented out the `oauth2` security scheme since it not supported by OASGen Provider and Rest Dynamic Controller while authentication is handled at the root level of the file and with basic auth scheme.
- Commented out `x-ms-*` items as they are specific to Microsoft internal tooling.
- Commented out the `components.parameters` section at root level as it just to inform about api versions.

## `PullRequest`

**Version**: `7.2-preview.2`
**Original specification file**:
- Link: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/master/specification/git/7.2/git.json
- Permalink: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/03cce9dd0355aa5a5fa808dcc0bd15aadda383b9/specification/git/7.2/git.json

**Transformations**: 
- Original file is converted from JSON to YAML and from Swagger 2.0 to OpenAPI 3.0.1 with the following tool: https://editor.swagger.io/
- File is pruned to only include the endpoints and schemas relevant to PullRequests.

**List of change made to the OpenAPI Specification (OAS)**:
- Added the following plugin endpoint:
  - `PATCH /api/{organization}/{project}/git/repositories/{repositoryId}/pullrequests/{pullRequestId}`
- Commented out the following endpoint that is not used in favor of the new plugin endpoint:
  - `PATCH /{organization}/{project}/_apis/git/repositories/{repositoryId}/pullrequests/{pullRequestId}`
- Modified status code for the `POST /{organization}/{project}/_apis/git/repositories/{repositoryId}/pullrequests` endpoint from `200` to `201` since the API actually returns `201 Created` on successful creation and not `200 OK`.
- Created new schemas for request bodies (`CreatePullRequestReq` and `UpdatePullRequestReq`), to ensure that only the relevant fields are included when creating or updating a Pull Request while the original `GitPullRequest` schema includes many read-only fields that should not be part of the request body.
- Commented out the `oauth2` security scheme since it not supported by OASGen Provider and Rest Dynamic Controller while authentication is handled at the root level of the file and with basic auth scheme.
- Commented out `x-ms-*` items as they are specific to Microsoft internal tooling.
- Commented out the `components.parameters` section at root level as it just to inform about api versions.

## `PolicyConfiguration`

**Version**: `7.2-preview.1`
**Original specification file**:
- Link: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/master/specification/policy/7.2/policy.json
- Permalink: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/a69e0a84db58a99dc4243957289b6d825dcb2af0/specification/policy/7.2/policy.json

**Transformations**: 
- Original file is converted from JSON to YAML and from Swagger 2.0 to OpenAPI 3.0.1 with the following tool: https://editor.swagger.io/
- File is pruned to only include the endpoints and schemas relevant to PolicyConfigurations.

**List of change made to the OpenAPI Specification (OAS)**:
- Created the following schemas required for proper functioning of the PolicyConfiguration resource and Rest Dynamic Controller operations: `PolicyReq`, `PolicyConfigurationSettings`, `PolicyScope`. These schemas were inferred from the examples provided at: https://github.com/MicrosoftDocs/vsts-rest-api-specs/tree/master/specification/policy/7.2/httpExamples.
`PolicyReq` is used as request body schema for `create` and `update` operations instead of the original `PolicyConfiguration` schema which contains read-only fields that should not be part of the request body but are included in the response body.
- Commented out the `oauth2` security scheme since it not supported by OASGen Provider and Rest Dynamic Controller while authentication is handled at the root level of the file and with basic auth scheme.
- Commented out `x-ms-*` items as they are specific to Microsoft internal tooling.
- Commented out the `components.parameters` section at root level as it just to inform about api versions.

## `RepositoryPermission`

Note: `RepositoryPermission` resource is an abstraction over the `AccessControlEntries` resource provided by the Azure DevOps REST API.
The OAS file is obtained by using thw `swag` tool to generate an OpenAPI Specification (OAS) from Go code. In this case, the Go code of the plugin that implements the `RepositoryPermission` resource over the `AccessControlEntries` resource.

**Version**: `7.2-preview.1` (AccessControlEntries)
**Original specification file** (AccessControlEntries):
- Link: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/master/specification/security/7.2/security.json
- Permalink: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/a69e0a84db58a99dc4243957289b6d825dcb2af0/specification/security/7.2/security.json

**List of change made to the OpenAPI Specification (OAS)** (generated entirely from the plugin code):
- Manually commented out the `additionalProperties` fields in the `allow` and `deny` properties to provide a fixed set of permissions that can be set on a RepositoryPermission resource. This is necessary to obtain a correct reconciliation of the resource by the Rest Dynamic Controller.
For example, part of the the modified `PermissionInfo` schema is as follows:
```yaml
allow:
  type: object
  # FIXED
  #additionalProperties:
  #  type: boolean
  properties:
    Administer:
      type: boolean
      default: false
      description: "displayName: 'Administer' "
    GenericRead:
      type: boolean
      default: false
      description: "displayName: 'Read' "
    GenericContribute:
      type: boolean
      default: false
      description: "displayName: 'Contribute' "
    ForcePush:
      type: boolean
      default: false
      description: "displayName: 'Force push (rewrite history, delete branches and tags)' "
    CreateBranch:
      type: boolean
      default: false
      description: "displayName: 'Create branch' "
...
```

## `BuildPermission`

Note: `BuildPermission` resource is an abstraction over the `AccessControlEntries` resource provided by the Azure DevOps REST API.
The OAS file is obtained by using thw `swag` tool to generate an OpenAPI Specification (OAS) from Go code. In this case, the Go code of the plugin that implements the `BuildPermission` resource over the `AccessControlEntries` resource.

**Version**: `7.2-preview.1` (AccessControlEntries)
**Original specification file** (AccessControlEntries):
- Link: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/master/specification/security/7.2/security.json
- Permalink: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/a69e0a84db58a99dc4243957289b6d825dcb2af0/specification/security/7.2/security.json

**List of change made to the OpenAPI Specification (OAS)** (generated entirely from the plugin code):
- Manually commented out the `additionalProperties` fields in the `allow` and `deny` properties to provide a fixed set of permissions that can be set on a BuildPermission resource. This is necessary to obtain a correct reconciliation of the resource by the Rest Dynamic Controller.
For example, part of the the modified `PermissionInfo` schema is as follows:
```yaml
allow:
  type: object
  # FIXED
  #additionalProperties:
  #  type: boolean
  properties:
    ViewBuilds:
      type: boolean
      default: false
      description: "displayName: 'View Builds' "
    EditBuildQuality:
      type: boolean
      default: false
      description: "displayName: 'Edit build quality' "
    RetainIndefinitely:
      type: boolean
      default: false
      description: "displayName: 'Retain indefinitely' "
    DeleteBuilds:
      type: boolean
      default: false
      description: "displayName: 'Delete builds' "
...
```

## `Team`

**Version**: `7.2-preview.3`
**Original specification file**:
- Link: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/master/specification/core/7.2/core.json
- Permalink: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/a69e0a84db58a99dc4243957289b6d825dcb2af0/specification/core/7.2/core.json

**Transformations**: 
- Original file is converted from JSON to YAML and from Swagger 2.0 to OpenAPI 3.0.1 with the following tool: https://editor.swagger.io/
- File is pruned to only include the endpoints and schemas relevant to Team resource.

**List of change made to the OpenAPI Specification (OAS)**:
- Changed the status code for the `create` operation from `200` to `201` since the Azure DevOps REST API actually returns a `201 Created` status code when a Team is created successfully.
- Changed the status code for the `delete` operation from `200` to `204` since the Azure DevOps REST API actually returns a `204 No Content` status code when a Team is deleted successfully.
- Commented out the `oauth2` security scheme since it not supported by OASGen Provider and Rest Dynamic Controller while authentication is handled at the root level of the file and with basic auth scheme.
- Commented out `x-ms-*` items as they are specific to Microsoft internal tooling.
- Commented out the `components.parameters` section at root level as it just to inform about api versions.

## `Project`

NOT YET AVAILABLE ON THE CHART
WORK IN PROGRESS

"7.2-preview.4"
https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/master/specification/core/7.2/core.json
https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/a69e0a84db58a99dc4243957289b6d825dcb2af0/specification/core/7.2/core.json

- Commented out the `components.parameters` section at root level as it just to inform about api versions.

200 to 202 in create (it is actually returning 202)

- Commented out the `oauth2` security scheme since it not supported by OASGen Provider and Rest Dynamic Controller while authentication is handled at the root level of the file and with basic auth scheme.
- Commented out `x-ms-*` items as they are specific to Microsoft internal tooling.
- Commented out the `components.parameters` section at root level as it just to inform about api versions.


## GraphGroup

**Version**: 7.2-preview.1
**Original specification file**:
- Link: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/master/specification/graph/7.2/graph.json
- Permalink: https://github.com/MicrosoftDocs/vsts-rest-api-specs/blob/a69e0a84db58a99dc4243957289b6d825dcb2af0/specification/graph/7.2/graph.json



- created new schema `GraphGroupCreationContextUnified` for the request body of the `create` operation.
- changed response code from 200 to 201 in create operation




## User

TODO