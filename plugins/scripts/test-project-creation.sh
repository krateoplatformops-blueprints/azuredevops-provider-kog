#!/bin/bash

# test-project-creation.sh
#
# Creates N test projects in Azure DevOps and then lists all projects.
# This script can be used to test the project plugin endpoints.
#
# Usage:
# AZURE_DEVOPS_TOKEN="your-token" ./scripts/test-project-creation.sh [ORGANIZATION] [CLEANUP] [NUM_PROJECTS]
#
# Example:
# AZURE_DEVOPS_TOKEN="xxx" ./scripts/test-project-creation.sh "krateo-kog" true 20
#
# Environment Variables:
# AZURE_DEVOPS_TOKEN: Bearer token for Azure DevOps API (required)
#
# Arguments:
# ORGANIZATION: Azure DevOps organization name (default: krateo-kog)
# CLEANUP: Set to "true" to automatically delete test projects (default: false)
# NUM_PROJECTS: Number of projects to create (default: 10)

set -e

# Configuration - Read token from environment variable
BEARER_TOKEN="${AZURE_DEVOPS_TOKEN}"
ORGANIZATION="${1:-krateo-kog}"
CLEANUP="${2:-false}"
NUM_PROJECTS="${3:-10}"
API_VERSION="7.2-preview.4"
BASE_URL="https://dev.azure.com/${ORGANIZATION}/_apis/projects"

# Output directory for JSON files
OUTPUT_DIR="./test-results"
mkdir -p "$OUTPUT_DIR"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if jq is installed
check_jq() {
  if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is not installed. Please install jq to run this script.${NC}"
    echo "Install with: brew install jq (macOS) or apt-get install jq (Linux)"
    exit 1
  fi
}

# Function to create a project
create_project() {
  local project_name=$1
  local project_number=$2
  local total_projects=$3

  echo -e "${BLUE}[$project_number/$total_projects]${NC} Creating project: ${YELLOW}${project_name}${NC}"

  local response=$(curl -s -X POST "${BASE_URL}?api-version=${API_VERSION}" \
    -H "Authorization: Bearer ${BEARER_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"${project_name}\",
      \"description\": \"Test project ${project_number} for automated testing\",
      \"visibility\": \"private\",
      \"capabilities\": {
        \"versioncontrol\": {
          \"sourceControlType\": \"Git\"
        },
        \"processTemplate\": {
          \"templateTypeId\": \"6b724908-ef14-45cf-84f8-768b5384da45\"
        }
      }
    }")

  # Check if request was successful
  if echo "$response" | jq -e '.id' > /dev/null 2>&1; then
    local operation_id=$(echo "$response" | jq -r '.id')
    echo -e "  ${GREEN}✓${NC} Operation ID: ${operation_id}"

    # Poll operation status (simplified - just check once after 2 seconds)
    sleep 2
    local operation_url="https://dev.azure.com/${ORGANIZATION}/_apis/operations/${operation_id}?api-version=7.2-preview"
    local operation_status=$(curl -s -X GET "${operation_url}" \
      -H "Authorization: Bearer ${BEARER_TOKEN}" | jq -r '.status')

    if [ "$operation_status" == "succeeded" ]; then
      echo -e "  ${GREEN}✓${NC} Status: ${operation_status}"
    else
      echo -e "  ${YELLOW}⚠${NC} Status: ${operation_status}"
    fi
  else
    echo -e "  ${RED}✗${NC} Failed to create project"
    echo -e "  ${RED}Response:${NC} $response"
  fi

  echo ""
}

# Function to list all projects and save to JSON file
list_projects() {
  local suffix=$1
  local output_file="${OUTPUT_DIR}/projects_list_${TIMESTAMP}${suffix}.json"

  echo -e "${BLUE}========================================${NC}"
  echo -e "${BLUE}Listing all projects in organization: ${YELLOW}${ORGANIZATION}${NC}"
  echo -e "${BLUE}========================================${NC}"
  echo ""

  local response=$(curl -s -X GET "${BASE_URL}?api-version=${API_VERSION}" \
    -H "Authorization: Bearer ${BEARER_TOKEN}")

  # Check if request was successful
  if echo "$response" | jq -e '.value' > /dev/null 2>&1; then
    # Save raw response to JSON file
    echo "$response" | jq . > "$output_file"
    echo -e "${GREEN}✓ Saved project list to: ${output_file}${NC}"
    echo ""

    local project_count=$(echo "$response" | jq '.value | length')
    echo -e "${GREEN}Total projects found: ${project_count}${NC}"
    echo ""

    # Display projects in a formatted table
    echo -e "${BLUE}Projects:${NC}"
    echo "$response" | jq -r '.value[] | "  - \(.name) (ID: \(.id), State: \(.state))"'

    # Filter and display only test projects created by this script
    echo ""
    echo -e "${BLUE}Test projects (TestProject_*)${NC}"
    local test_projects=$(echo "$response" | jq -r '.value[] | select(.name | startswith("TestProject_")) | "  - \(.name) (ID: \(.id), Created: \(.lastUpdateTime))"')
    if [ -n "$test_projects" ]; then
      echo "$test_projects"

      # Save filtered test projects to separate file
      local test_projects_file="${OUTPUT_DIR}/test_projects_${TIMESTAMP}${suffix}.json"
      echo "$response" | jq '{count: ([.value[] | select(.name | startswith("TestProject_"))] | length), projects: [.value[] | select(.name | startswith("TestProject_"))]}' > "$test_projects_file"
      echo ""
      echo -e "${GREEN}✓ Saved test projects to: ${test_projects_file}${NC}"
    else
      echo -e "${YELLOW}  No test projects found${NC}"
    fi
  else
    echo -e "${RED}Failed to list projects${NC}"
    echo -e "${RED}Response:${NC} $response"
    return 1
  fi
}

# Function to clean up test projects (optional)
cleanup_projects() {
  echo -e "${BLUE}========================================${NC}"
  echo -e "${BLUE}Cleaning up test projects...${NC}"
  echo -e "${BLUE}========================================${NC}"
  echo ""

  local response=$(curl -s -X GET "${BASE_URL}?api-version=${API_VERSION}" \
    -H "Authorization: Bearer ${BEARER_TOKEN}")

  # Get all test project IDs
  local test_project_ids=$(echo "$response" | jq -r '.value[] | select(.name | startswith("TestProject_")) | .id')

  if [ -z "$test_project_ids" ]; then
    echo -e "${YELLOW}No test projects found to clean up.${NC}"
    return
  fi

  local count=0
  for project_id in $test_project_ids; do
    count=$((count + 1))
    local project_name=$(echo "$response" | jq -r ".value[] | select(.id == \"$project_id\") | .name")
    echo -e "${BLUE}[$count]${NC} Deleting project: ${YELLOW}${project_name}${NC} (ID: ${project_id})"

    local delete_response=$(curl -s -X DELETE "${BASE_URL}/${project_id}?api-version=${API_VERSION}" \
      -H "Authorization: Bearer ${BEARER_TOKEN}")

    if echo "$delete_response" | jq -e '.id' > /dev/null 2>&1; then
      local operation_id=$(echo "$delete_response" | jq -r '.id')
      echo -e "  ${GREEN}✓${NC} Delete operation ID: ${operation_id}"
    else
      echo -e "  ${RED}✗${NC} Failed to delete project"
    fi
  done

  echo -e "${GREEN}Cleanup initiated for ${count} projects.${NC}"
  echo -e "${YELLOW}Note: Deletion is async. Wait a few seconds before listing again.${NC}"
}

# Main script
main() {
  echo -e "${BLUE}========================================${NC}"
  echo -e "${BLUE}Azure DevOps Project Creation Test${NC}"
  echo -e "${BLUE}========================================${NC}"
  echo -e "Organization: ${YELLOW}${ORGANIZATION}${NC}"
  echo -e "API Version: ${YELLOW}${API_VERSION}${NC}"
  echo -e "Projects to create: ${YELLOW}${NUM_PROJECTS}${NC}"
  echo -e "Auto-cleanup: ${YELLOW}${CLEANUP}${NC}"
  echo -e "Output directory: ${YELLOW}${OUTPUT_DIR}${NC}"
  echo ""

  # Validate BEARER_TOKEN is provided via environment variable
  if [ -z "$BEARER_TOKEN" ]; then
    echo -e "${RED}Error: AZURE_DEVOPS_TOKEN environment variable is required${NC}"
    echo -e "${YELLOW}Usage: AZURE_DEVOPS_TOKEN=\"your-token\" ./scripts/test-project-creation.sh [ORGANIZATION] [CLEANUP] [NUM_PROJECTS]${NC}"
    exit 1
  fi

  # Validate NUM_PROJECTS is a positive integer
  if ! [[ "$NUM_PROJECTS" =~ ^[0-9]+$ ]] || [ "$NUM_PROJECTS" -le 0 ]; then
    echo -e "${RED}Error: NUM_PROJECTS must be a positive integer (got: $NUM_PROJECTS)${NC}"
    exit 1
  fi

  # Check dependencies
  check_jq

  # Generate timestamp for unique project names and output files
  TIMESTAMP=$(date +%Y%m%d_%H%M%S)

  # Create N projects
  echo -e "${BLUE}Creating ${NUM_PROJECTS} test projects...${NC}"
  echo ""

  for ((i=1; i<=NUM_PROJECTS; i++)); do
    project_name="TestProject_${TIMESTAMP}_${i}"
    create_project "$project_name" "$i" "$NUM_PROJECTS"
  done

  # Wait a bit for all operations to complete
  echo -e "${YELLOW}Waiting 3 seconds for all operations to complete...${NC}"
  sleep 3
  echo ""

  # List all projects and save to JSON (before cleanup)
  list_projects "_before_cleanup"

  # Automatic cleanup based on CLEANUP flag
  echo ""
  if [[ "$CLEANUP" == "true" ]]; then
    echo -e "${YELLOW}Auto-cleanup enabled. Deleting test projects...${NC}"
    cleanup_projects
    echo ""
    echo -e "${YELLOW}Waiting 5 seconds for deletions to complete...${NC}"
    sleep 5
    echo ""
    # List projects again after cleanup
    list_projects "_after_cleanup"
  else
    echo -e "${YELLOW}Auto-cleanup disabled. Test projects will remain in the organization.${NC}"
    echo -e "${YELLOW}To enable cleanup, run with: AZURE_DEVOPS_TOKEN=\"xxx\" ./scripts/test-project-creation.sh [ORG] true [NUM_PROJECTS]${NC}"
  fi

  echo ""
  echo -e "${GREEN}✓ Script completed successfully!${NC}"
  echo -e "${GREEN}✓ Created ${NUM_PROJECTS} projects. Results saved to: ${OUTPUT_DIR}/${NC}"
}

# Run main function
main
