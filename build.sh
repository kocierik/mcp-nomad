#!/bin/bash

# Exit on error
set -e

VERSION="1.1.0"
BINARY_NAME="mcp-nomad-go"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

# Clean and create dist directory
echo "Cleaning dist directory..."
rm -rf dist
mkdir -p dist

# Build targets
declare -A TARGETS=(
  ["darwin_amd64"]="1"
  ["darwin_arm64"]="8.0"
  ["linux_amd64"]="1"
  ["linux_arm64"]="8.0"
  ["windows_amd64"]="1"
  ["windows_arm64"]="8.0"
)

# Build binaries
echo "Building binaries..."
for target in "${!TARGETS[@]}"; do
    GOOS="${target%%_*}"
    GOARCH="${target##*_}"
    VERSION_SUFFIX="${TARGETS[$target]}"
    OUTPUT_DIR="dist/${BINARY_NAME}_${GOOS}_${GOARCH}_v${VERSION_SUFFIX}"

    echo "Building for $GOOS/$GOARCH..."
    mkdir -p "$OUTPUT_DIR"
    if [[ "$GOOS" == "windows" ]]; then
        go build -o "${OUTPUT_DIR}/${BINARY_NAME}.exe" .
    else
        go build -o "${OUTPUT_DIR}/${BINARY_NAME}" .
    fi
done

# Compress built binaries
echo "Compressing binaries..."
cd dist || exit 1

for dir in ${BINARY_NAME}_*; do
    if [[ "$dir" == *"windows"* ]]; then
        echo "Zipping $dir..."
        zip -r "${dir}.zip" "$dir"
    else
        echo "Creating tar.gz for $dir..."
        tar czf "${dir}.tar.gz" "$dir"
    fi
done

# Show contents
echo "Final contents of dist/:"
ls -la | cat

echo "✅ All builds and compressions completed successfully!"

