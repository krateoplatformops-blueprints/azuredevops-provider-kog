# Membership Samples

This folder contains sample YAML files for managing Azure DevOps Graph memberships using the OASGen Provider.

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

1. Membership RestDefinition is installed in the cluster
2. Create a MembershipConfiguration CR with authentication details
3. Obtain the descriptors for the member and container

## Sample Files

- `membership_1_user_to_group.yaml` - Add a user to a group
- `membership_2_group_to_team.yaml` - Add a group to a team
- `membership_3_group_to_group_project_level.yaml` - Add a group to another group at the project level
- `membership_4_group_to_group_organization_level.yaml` - Add a group to another group at the organization level
- `membership_5_team_to_group.yaml` - Add a team to a group
