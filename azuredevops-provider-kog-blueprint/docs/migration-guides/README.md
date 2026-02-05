# Azure DevOps Provider KOG - Migration Guide

This document is an index of migration guides to help you migrate your resources from the [Azure DevOps Provider (classic)](https://github.com/krateoplatformops/azuredevops-provider) to the [Azure DevOps Provider KOG](https://github.com/krateoplatformops-blueprints/azuredevops-provider-kog).

Currently, the Azure DevOps Provider KOG supports the following resources:
- `GitRepository`
- `Pipeline`
- `PipelinePermission`
- `PullRequest`
- `PolicyConfiguration`
- `RepositoryPermission`
- `BuildPermission`
- `Team`
- `GraphGroup`
- `Project`
- `User`
- `Environment`
- `VariableGroup`

Therefore, currently, Azure DevOps Provider KOG **is not a drop-in replacement** for the Azure DevOps Provider (classic), but it is a new provider that supports another set of resources.

## Summary

- [Pre-requisites](#pre-requisites)
- [Index of migration guides](#index-of-migration-guides)

## Pre-requisites

- You have a the [Azure DevOps Provider (classic)](https://github.com/krateoplatformops/azuredevops-provider) installed and properly configured in your cluster (including `ConnectorConfig`resource for authentication).
- You have the [Azure DevOps Provider KOG](https://github.com/krateoplatformops-blueprints/azuredevops-provider-kog) installed and properly configured in your cluster (including configuration resources for authentication, see [Authentication](../../README.md#authentication) and [Configuration](../../README.md#configuration)).


## Index of migration guides

- [GitRepository migration example](gitrepository-migration.md)
- [Pipeline migration example](pipeline-migration.md)
- [PipelinePermission migration example](pipelinepermission-migration.md)
- [PullRequest migration example](pullrequest-migration.md)
- [PolicyConfiguration migration example](policyconfiguration-migration.md)
- [RepositoryPermission migration example](repositorypermission-migration.md)
- [Team migration example](team-migration.md)
- [GraphGroup migration example](graphgroup-migration.md)
- [Project migration example](project-migration.md)
- [User migration example](user-migration.md)
- [Environment migration example](environment-migration.md)
- [VariableGroup migration example](variablegroup-migration.md)
