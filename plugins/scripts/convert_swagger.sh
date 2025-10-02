#!/bin/bash

# Default values
INPUT_FILE=""
OUTPUT_BASE=""
OUTPUT_DIR="."

# Parse command line arguments
while getopts "i:o:" opt; do
    case $opt in
        i)
            INPUT_FILE="$OPTARG"
            ;;
        o)
            # Split the argument into directory and filename
            OUTPUT_DIR=$(dirname "$OPTARG")
            OUTPUT_BASE=$(basename "$OPTARG")
            ;;
        \?)
            echo "Invalid option: -$OPTARG" >&2
            exit 1
            ;;
        :)
            echo "Option -$OPTARG requires an argument." >&2
            exit 1
            ;;
    esac
done

# Check if input file is provided
if [ -z "$INPUT_FILE" ]; then
    echo "Usage: $0 -i <swagger_v2_file> [-o output_base]"
    echo "Example: $0 -i swagger.json -o openapi/v3/spec"
    echo "This will create spec.json and spec.yaml in openapi/v3/"
    exit 1
fi

# Check if file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: File $INPUT_FILE not found"
    exit 1
fi

# Set default output base if not provided
if [ -z "$OUTPUT_BASE" ]; then
    FILENAME=$(basename -- "$INPUT_FILE")
    FILENAME_NOEXT="${FILENAME%.*}"
    OUTPUT_BASE="${FILENAME_NOEXT}_openapi_v3"
fi

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Full output paths
JSON_OUTPUT="${OUTPUT_DIR%/}/${OUTPUT_BASE}.json"
YAML_OUTPUT="${OUTPUT_DIR%/}/${OUTPUT_BASE}.yaml"

# API endpoint
API_URL="https://converter.swagger.io/api/convert"

# Send the file to the converter API
echo "Converting $INPUT_FILE to OpenAPI v3..."
response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "Accept: application/json" \
    --data-binary @"$INPUT_FILE" \
    "$API_URL")

# Check if conversion was successful
if [ $? -ne 0 ]; then
    echo "Error: Failed to connect to the converter API"
    exit 1
fi

# Check if response contains error
if echo "$response" | grep -q '"error"'; then
    echo "Error in conversion:"
    echo "$response" | jq -r '.message'
    exit 1
fi

# Save JSON output
echo "$response" | jq '.' > "$JSON_OUTPUT"

# Convert to YAML and save
echo "$response" | jq '.' | yq -P  > "$YAML_OUTPUT"

if [ $? -ne 0 ]; then
    echo "Warning: Couldn't generate YAML output (is yq installed?)"
    echo "JSON output saved to $JSON_OUTPUT"
else
    echo "Conversion successful. Output files:"
    echo "- JSON: $JSON_OUTPUT"
    echo "- YAML: $YAML_OUTPUT"
fi
