#!/bin/bash

# Get the current version from the main package.json file
previous_version=$(grep -Po '"version": "\K[^"]*' ./packages/npm-mcp-nomad/package.json)
new_version="${1}"

if [ -z "$new_version" ]; then
  echo "Usage: $0 <new_version>"
  exit 1
fi

if [ -z "$previous_version" ]; then
  echo "Error: Could not find current version in package.json"
  exit 1
fi

echo "Updating version from $previous_version to $new_version"

# replace previous version with new version in all .json files in ./packages folder 
find ./packages -type f -name '*.json' -exec sed -i "s/${previous_version}/${new_version}/g" {} \;

echo "Version update complete"
