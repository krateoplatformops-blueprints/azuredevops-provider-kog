# Membership Samples

This folder contains sample YAML files for managing Azure DevOps Graph memberships using the Krateo OASGen Provider.

## Overview

A Membership represents a relationship between a **member** (user, group, or team) and a **container** (group or team). This resource allows you to manage who belongs to which groups or teams in Azure DevOps.

## Membership Types

The same Membership resource can represent different types of relationships:

| Member Type | Container Type | Description |
|-------------|----------------|-------------|
| User | Group | Add a user to a group |
| User | Team | Add a user to a team |
| Group | Group | Nested groups (add a group as member of another group) |
| Team | Group | Add a team as member of a group |

## Prerequisites

1. Deploy the Membership RestDefinition using the `azuredevops-provider-kog-membership` Helm chart
2. Create a MembershipConfiguration CR with authentication details
3. Obtain the descriptors for the member and container:
   - User descriptors: `GET /{organization}/_apis/graph/users`
   - Group descriptors: `GET /{organization}/_apis/graph/groups`
   - Team descriptors: Teams are represented as groups with specific characteristics

## Sample Files

- `membership_1_user_to_group.yaml` - Add a user to a group
- `membership_2_group_to_group.yaml` - Add a group as member of another group (nested groups)

## Getting Descriptors

### Get User Descriptor
```bash
curl -u ":<PAT_TOKEN>" \
  "https://vssps.dev.azure.com/{organization}/_apis/graph/users?subjectTypes=aad&api-version=7.2-preview.1"
```

### Get Group Descriptor
```bash
curl -u ":<PAT_TOKEN>" \
  "https://vssps.dev.azure.com/{organization}/_apis/graph/groups?api-version=7.2-preview.1"
```

## Important Notes

- **No update operation**: Memberships cannot be updated. To change a membership, delete and recreate it.
- **Idempotent creation**: Creating a membership that already exists will succeed (returns 201).
- **Both descriptors required**: You must specify both `subjectDescriptor` (member) and `containerDescriptor` (container).
