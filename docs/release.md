# Release Process Guide

This document outlines the process for releasing new versions of the components within this repository. The project uses a **component-based tagging strategy**, where each independent part (plugin, resource blueprint, main blueprint) is versioned and released separately.

## Standard Release Workflow

The standard workflow is:
1.  Merge your completed feature or bugfix into the `main` branch via a pull request.
2.  Pull the latest changes to your local `main` branch.
3.  Create and push the appropriate tag from your local `main` branch.

## Tagging Convention

Tags are composed of two parts: `component-name/version`. The versioning scheme is Semantic Versioning (`X.Y.Z`).

| Component Type          | Tag Format Example                  | Workflow Triggered                 |
| ----------------------- | ----------------------------------- | ---------------------------------- |
| Plugin Source           | `collaborator-plugin/1.0.3`         | `plugins-source-tag.yaml`          |
| Resource Blueprint      | `collaborator-blueprint/1.0.3`      | `resource-blueprint-tag.yaml`      |
| Main (Umbrella) Blueprint | `1.0.3`                             | `main-blueprint-tag.yaml`          |

---

## Comprehensive Release Scenario

This scenario walks through a complete release cycle, covering all component types.

**Goal:** Release a new feature that requires changes to the `collaborator` plugin, its blueprint, and the main umbrella chart.

**Prerequisites:**
- All code changes have been merged into the `main` branch.
- Your local `main` branch is up-to-date (`git pull origin main`).

### Step 1: Release the Plugin

First, release the container image for the plugin.

- **Action:** Tag and push the new version for the `collaborator-plugin`.
- **Command:**
  ```sh
  git tag collaborator-plugin/1.2.0
  git push origin tag collaborator-plugin/1.2.0
  ```
- **Result:** The `plugins-source-tag` workflow runs, publishing the `collaborator-plugin:1.2.0` container image to the GitHub Container Registry.
- **Verification:** Wait for the workflow to complete successfully before proceeding.

### Step 2: Release the Resource Blueprint

Now that the new plugin image is available, release the Helm chart that uses it.

- **Action:** Tag and push the new version for the `collaborator-blueprint`.
- **Command:**
  ```sh
  git tag collaborator-blueprint/1.2.0
  git push origin tag collaborator-blueprint/1.2.0
  ```
- **Result:** The `resource-blueprint-tag` workflow runs. It will:
  1.  Find the `collaborator-plugin:1.2.0` image you just published (latest tag).
  2.  Set the chart's `version` to `1.2.0` and its `appVersion` to `1.2.0` (appVersion reflects the plugin image version).
  3.  Publish the new `github-provider-kog-collaborator` chart to the Helm repository.
- **Verification:** Wait for the workflow to complete successfully.

### Step 3: Release the Main Blueprint

Finally, release the main umbrella chart to make the new version of the collaborator blueprint available to end-users.

- **Action:** Tag and push the new version for the main blueprint.
- **Command:**
  ```sh
  git tag 2.0.0
  git push origin tag 2.0.0
  ```
- **Result:** The `main-blueprint-tag` workflow runs. It will:
  1.  Find the `github-provider-kog-collaborator:1.2.0` chart you just published.
  2.  Update its `Chart.yaml` to use this new version as a dependency (field `version` under the element in dependencies).
  3.  Set its own chart `version` to `2.0.0`.
  4.  Publish the new `github-provider-kog` chart to the Helm repository.

This completes the end-to-end release process.
