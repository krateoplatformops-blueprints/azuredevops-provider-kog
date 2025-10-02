#!/bin/bash

# swag-init.sh
#
# Generates OpenAPI documentation for a specific plugin.
# This script must be run from the root of the 'plugins' directory.
#
# Usage:
# ./scripts/swag-init.sh <plugin-name>
#
# Example:
# ./scripts/swag-init.sh collaborator-plugin

set -e

# Check for the plugin name argument
if [ -z "$1" ]; then
  echo "Error: Missing plugin name."
  echo "Usage: $0 <plugin-name>"
  echo "Example: $0 collaborator-plugin"
  exit 1
fi

PLUGIN_NAME=$1
PLUGIN_DIR="cmd/${PLUGIN_NAME}"

# Verify that the target plugin directory exists
if [ ! -d "$PLUGIN_DIR" ]; then
  echo "Error: Plugin directory not found at ${PLUGIN_DIR}"
  exit 1
fi

echo "Generating OpenAPI documentation for ${PLUGIN_NAME}..."

# Define paths relative to the plugin's directory
MAIN_GO_PATH="${PLUGIN_DIR}/main.go"
OUTPUT_DIR="${PLUGIN_DIR}/docs"

# Generate the Swagger v2 JSON file
# The --output flag sets the destination directory for the 'docs' folder
swag init --parseDependency -g "${MAIN_GO_PATH}" --output "${OUTPUT_DIR}"

# Define the input and output paths for the conversion script
SWAGGER_V2_JSON="${OUTPUT_DIR}/swagger.json"
OPENAPI_V3_OUTPUT="${OUTPUT_DIR}/openapi3"

echo "Converting Swagger v2 to OpenAPI v3..."
# Call the conversion script, which will create openapi3.json and openapi3.yaml
./scripts/convert_swagger.sh -i "${SWAGGER_V2_JSON}" -o "${OPENAPI_V3_OUTPUT}"

echo "Documentation generation complete for ${PLUGIN_NAME}."
echo "Output is in ${OUTPUT_DIR}/"