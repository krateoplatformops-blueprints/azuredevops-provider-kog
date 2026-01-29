#!/bin/bash
set -eo pipefail

#
# Script to generate documentation for Azure DevOps Provider KOG blueprints.
#
# Usage: ./generate_docs.sh [--all | --name <resource_name> | -n <resource_name>]
#   --all: Process all blueprint directories.
#   --name <resource_name>, -n <resource_name>: Process only the specified blueprint resource (e.g., 'gitrepository').
#
# For each selected directory matching '*-blueprint', it performs the following steps:
# 1. Cleans up any previously generated documentation files.
# 2. Sets a dummy version in the Chart.yaml file.
# 3. Installs the Helm chart into a temporary namespace.
# 4. Waits for the corresponding 'restdefinition' Kubernetes resource to become available.
# 5. Extracts the main CRD name
# 6. Saves the main CRD and any associated configuration CRD
# 7. Uses 'crdoc' to generate markdown documentation from the main CRD.
# 8. Uninstalls the Helm chart and cleans up the temporary namespace.
#
# Requirements:
# - 'kubectl' context must be configured to point to a Kubernetes cluster.
# - 'helm' CLI must be installed.
# - 'crdoc' CLI must be installed.
#

# --- Configuration ---
# Namespace for temporary Helm releases.
RELEASE_NAMESPACE="docgen-ns"
# Timeout for waiting on Kubernetes resources.
K8S_WAIT_TIMEOUT="7m"

# --- Location Independent Setup ---
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
printf "Script Directory: %s\n" "$SCRIPT_DIR"
PROJECT_ROOT=$( cd -- "$SCRIPT_DIR/.." &> /dev/null && pwd )
printf "Project Root Directory: %s\n" "$PROJECT_ROOT"
PROVIDER_NAME="azuredevops-provider-kog"
PROVIDER_RESOURCE_PREFIX="azuredevops"


# Change to the project root directory to ensure all paths are relative to it.
cd "$PROJECT_ROOT"

# --- Usage Function ---
usage() {
  echo "Usage: $0 [--all | --name <resource_name> | -n <resource_name>]" >&2
  echo "  --all: Process all blueprint directories." >&2
  echo "  --name <resource_name>, -n <resource_name>: Process only the specified blueprint resource (e.g., 'gitrepository')." >&2
  exit 1
}

# --- Argument Parsing ---
process_all=false
target_resource_name=""

while [[ "$#" -gt 0 ]]; do
  case "$1" in
    --all)
      process_all=true
      shift
      ;;
    --name|-n)
      if [[ -n "$2" ]] && [[ "$2" != --* ]]; then
        target_resource_name="$2"
        shift 2
      else
        echo "ERROR: --name|-n requires a resource name." >&2
        usage
      fi
      ;;
    *)
      echo "ERROR: Unknown option: $1" >&2
      usage
      ;;
  esac
done

# Validate arguments
if ! "$process_all" && [[ -z "$target_resource_name" ]]; then
  echo "ERROR: Please specify either --all or --name <resource_name>." >&2
  usage
fi

if "$process_all" && [[ -n "$target_resource_name" ]]; then
  echo "ERROR: Cannot use --all with --name/-n simultaneously." >&2
  usage
fi

# --- Body ---

# Ensure the script is run from the correct directory.
if [[ ! -f "README.md" ]] || [[ ! -d "$PROVIDER_NAME-blueprint" ]]; then
  echo "ERROR: Please run this script from the root of the '$PROVIDER_NAME' directory." >&2
  exit 1
fi

# Check for required command-line tools.
for cmd in "helm" "kubectl" "crdoc"; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "ERROR: Required command '$cmd' is not installed. Please install it and try again." >&2
    exit 1
  fi
done

# check oasgen-provider is installed

# Clean up the release namespace at the start, in case of a previous failed run.
echo "Cleaning up release namespace '$RELEASE_NAMESPACE' if it exists..."
kubectl delete namespace "$RELEASE_NAMESPACE" --ignore-not-found=true
echo "---"

UMBRELLA_DIR="$PROVIDER_NAME-blueprint"
echo "Umbrella Directory: $UMBRELLA_DIR"

# Determine which blueprint directories to process.
BLUEPRINT_DIRS=()
if "$process_all"; then
  
  BLUEPRINT_DIRS=( *-blueprint )

  # exclude umbrella blueprint directory of shape <provider-name>-blueprint
  for i in "${!BLUEPRINT_DIRS[@]}"; do
    if [[ "${BLUEPRINT_DIRS[i]}" == "$PROVIDER_NAME-blueprint" ]]; then
      unset 'BLUEPRINT_DIRS[i]'
    fi
  done
  # Re-index the array to remove any gaps created by unset
  BLUEPRINT_DIRS=( "${BLUEPRINT_DIRS[@]}" )
  
  printf "Found %d blueprint directories to process:\n" "${#BLUEPRINT_DIRS[@]}"
  for dir in "${BLUEPRINT_DIRS[@]}"; do
    printf " - %s\n" "$dir"
  done
elif [[ -n "$target_resource_name" ]]; then
  SPECIFIC_BLUEPRINT_DIR="$PROVIDER_NAME-${target_resource_name}-blueprint"
  if [[ ! -d "$SPECIFIC_BLUEPRINT_DIR" ]]; then
    echo "ERROR: Blueprint directory '$SPECIFIC_BLUEPRINT_DIR' not found." >&2
    exit 1
  fi
  BLUEPRINT_DIRS=( "$SPECIFIC_BLUEPRINT_DIR" )
  printf "Processing specific blueprint directory: %s\n" "$SPECIFIC_BLUEPRINT_DIR"
fi

# Loop through each selected blueprint directory.
for BLUEPRINT_DIR in "${BLUEPRINT_DIRS[@]}"; do
  if [[ -d "$BLUEPRINT_DIR" ]]; then
    # Run each blueprint processing in a subshell to isolate changes and ensure cleanup.
    (
      set -eo pipefail
      
      echo "Processing blueprint: $BLUEPRINT_DIR"
      cd "$BLUEPRINT_DIR"

      # Ensure Chart.yaml is restored even if the script fails.
      trap 'if [[ -f "Chart.yaml.original" ]]; then echo "Restoring original Chart.yaml in $(pwd)..."; mv -f "Chart.yaml.original" "Chart.yaml"; fi' EXIT ERR INT
      cp Chart.yaml Chart.yaml.original

      # Extract the resource name from the directory name (e.g., 'gitrepository').
      RESOURCE_NAME=$(echo "$BLUEPRINT_DIR" | sed "s/$PROVIDER_NAME-//" | sed 's/-blueprint//')
      printf "Resource Name: %s\n" "$RESOURCE_NAME"
      RELEASE_NAME="ado-$RESOURCE_NAME"
      printf "Release Name: %s\n" "$RELEASE_NAME"

      # Define paths for documentation output.
      CRD_DIR="$PROJECT_ROOT/$UMBRELLA_DIR/docs/crds"
      CONFIG_CRD_DIR="$CRD_DIR/configuration-crds"
      CRDOC_DIR="$PROJECT_ROOT/$UMBRELLA_DIR/docs/crds-reference"
      echo "CRD Directory: $CRD_DIR"
      echo "Configuration CRD Directory: $CONFIG_CRD_DIR"
      echo "CRDoc Output Directory: $CRDOC_DIR"

      # Ensure documentation directories exist
      echo "Ensuring documentation directories exist..."
      mkdir -p "$CONFIG_CRD_DIR"
      mkdir -p "$CRDOC_DIR"

      # Temporarily modify Chart.yaml to set a dummy version for installation.
      sed -i.bak '/^version:/ { s/: .*$/: 0.0.0-dev/ ; }' Chart.yaml
      rm -f Chart.yaml.bak

      # Install the Helm chart.
      echo "Installing Helm chart '$RELEASE_NAME'..."
      helm install "$RELEASE_NAME" . --namespace "$RELEASE_NAMESPACE" --create-namespace
      echo "Helm chart '$RELEASE_NAME' installed."

      # print all restdefinitions in the namespace for debugging
      echo "Current RestDefinitions in namespace '$RELEASE_NAMESPACE':"
      kubectl get restdefinition.ogen.krateo.io -n "$RELEASE_NAMESPACE"

      # The RestDefinition name follows the pattern: $RELEASE_NAME-$RESOURCE_NAME
      # TODO: maybe a more robust way to determine the RD name?
      #REST_DEF_NAME="${RELEASE_NAME}-${RESOURCE_NAME}"
      REST_DEF_NAME="azuredevops-provider-kog-${RESOURCE_NAME}"
      
      echo "Waiting for RestDefinition '$REST_DEF_NAME' to be ready..."
      if ! kubectl wait restdefinition.ogen.krateo.io "$REST_DEF_NAME" \
        --for=condition=Ready \
        --namespace "$RELEASE_NAMESPACE" \
        --timeout="$K8S_WAIT_TIMEOUT"; then
        echo "ERROR: RestDefinition '$REST_DEF_NAME' did not become ready in time." >&2
        echo "--- Debug Info ---"
        kubectl get restdefinition.ogen.krateo.io "$REST_DEF_NAME" -n "$RELEASE_NAMESPACE" -o yaml
        kubectl get pods -n "$RELEASE_NAMESPACE"
        echo "------------------"
        helm uninstall "$RELEASE_NAME" --namespace "$RELEASE_NAMESPACE"
        continue
      fi

      # Find main CRD by checking for 'restresources' category in the CRD spec and matching the resource name within all CRDs in the cluster.
      # Like
      #categories:
      # #- team
      #- restconfigs

      #Loop through all CRDs in the cluster and find the main CRD and config CRD by looking at categories
      # the main CRD will have a category matching the resource name and also the category 'restresources'
      # the config CRD will have a category matching the resource name and also the category 'restconfigs'
      MAIN_CRD_NAME=""
      CONFIG_CRD_NAME=""
      CONFIG_CRD_KIND=""
      echo "Identifying main and configuration CRDs..."
      ALL_CRDS=$(kubectl get crds -o name | grep "ogen.krateo.io" | sed 's,customresourcedefinition.apiextensions.k8s.io/,,g')
      echo "Total CRDs found in cluster: $(echo "$ALL_CRDS" | wc -l)"
      echo "CRDs List:"
      echo "$ALL_CRDS"
      for CRD in $ALL_CRDS; do
        CATEGORIES=$(kubectl get crd "$CRD" -o json | jq -r '.spec.names.categories[]')
        echo "Checking CRD: $CRD with categories: $CATEGORIES"
        for CATEGORY in $CATEGORIES; do
          echo "Inspecting category: $CATEGORY"
          if [[ "$CATEGORY" == "$RESOURCE_NAME" ]]; then
            echo "Resource name match found in CRD: $CRD"
            if echo "$CATEGORIES" | grep -qw "restresources"; then
              echo "Found 'restresources' category in CRD: $CRD"
              MAIN_CRD_NAME="$CRD"
              echo "Found main CRD: $MAIN_CRD_NAME"
            elif echo "$CATEGORIES" | grep -qw "restconfigs"; then
              echo "Found 'restconfigs' category in CRD: $CRD"
              CONFIG_CRD_NAME="$CRD"
              CONFIG_CRD_KIND=$(kubectl get crd "$CRD" -o jsonpath='{.spec.names.kind}')
              echo "Found configuration CRD: $CONFIG_CRD_NAME"
            fi
          fi
        done
      done

      # Save the main CRD definition to a file with the naming convention <resource_name>-crd.yaml.
      MAIN_CRD_FILE_PATH="$CRD_DIR/$RESOURCE_NAME-crd.yaml"
      echo "Saving main CRD: '$MAIN_CRD_NAME'"
      kubectl get crd "$MAIN_CRD_NAME" -o yaml >"$MAIN_CRD_FILE_PATH"
      echo "Main CRD saved to '$MAIN_CRD_FILE_PATH'"

      #Remove from file the lines in metadata that are not needed
      #echo "Cleaning up main CRD file metadata..."
      #yq eval 'del(.metadata.creationTimestamp, .metadata.generation, .metadata.resourceVersion, .metadata.uid)' -i "$MAIN_CRD_FILE_PATH"
      #echo "Main CRD file metadata cleaned."

      # Remove status section if exists
      #echo "Removing status section from main CRD file if it exists..."
      #yq eval 'del(.status)' -i "$MAIN_CRD_FILE_PATH"
      #echo "Status section removed from main CRD file."

      # Save configuration CRD to a file if it exists.
      if [[ -n "$CONFIG_CRD_NAME" ]]; then
        echo "Configuration CRD kind: $CONFIG_CRD_KIND"

        #kind lowercase for filename
        CONFIG_CRD_FILE_KIND_LOWER=$(echo "$CONFIG_CRD_KIND" | tr '[:upper:]' '[:lower:]')
        echo "Configuration CRD file kind lowercase: $CONFIG_CRD_FILE_KIND_LOWER"

        CONFIG_CRD_FILE_PATH="$CONFIG_CRD_DIR/$CONFIG_CRD_FILE_KIND_LOWER-crd.yaml"
        echo "Saving configuration CRD: '$CONFIG_CRD_NAME'"
        kubectl get crd "$CONFIG_CRD_NAME" -o yaml >"$CONFIG_CRD_FILE_PATH"
        echo "Configuration CRD saved to '$CONFIG_CRD_FILE_PATH'"
      else
        echo "No configuration CRD found for resource '$RESOURCE_NAME'."
      fi

      # Generate markdown documentation using crdoc.
      echo "Generating documentation for '$RESOURCE_NAME' with crdoc..."
      crdoc --resources "$MAIN_CRD_FILE_PATH" --output "$CRDOC_DIR/$RESOURCE_NAME.md"
      echo "Documentation generated at '$CRDOC_DIR/$RESOURCE_NAME.md'"

      # Uninstall the Helm chart.
      echo "Uninstalling Helm chart '$RELEASE_NAME'..."
      helm uninstall "$RELEASE_NAME" --namespace "$RELEASE_NAMESPACE"
    )
    echo "Finished processing '$BLUEPRINT_DIR'."
    echo "---"
  fi
done

# Final cleanup of the namespace.
echo "Final cleanup of namespace '$RELEASE_NAMESPACE'..."
kubectl delete namespace "$RELEASE_NAMESPACE" --ignore-not-found=true

echo "Documentation generation complete."

# todo: clean metatdata from generated crds
# todo: remove status from generated crds