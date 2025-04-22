#!/bin/bash

# Exit on error
set -e

VERSION="1.0.0"
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

# Build for Darwin (macOS)
echo "Building for darwin/amd64..."
mkdir -p "dist/${BINARY_NAME}_darwin_amd64_v1"
GOOS=darwin GOARCH=amd64 go build -o "dist/${BINARY_NAME}_darwin_amd64_v1/${BINARY_NAME}" .

echo "Building for darwin/arm64..."
mkdir -p "dist/${BINARY_NAME}_darwin_arm64_v8.0"
GOOS=darwin GOARCH=arm64 go build -o "dist/${BINARY_NAME}_darwin_arm64_v8.0/${BINARY_NAME}" .

# Build for Linux
echo "Building for linux/amd64..."
mkdir -p "dist/${BINARY_NAME}_linux_amd64_v1"
GOOS=linux GOARCH=amd64 go build -o "dist/${BINARY_NAME}_linux_amd64_v1/${BINARY_NAME}" .

echo "Building for linux/arm64..."
mkdir -p "dist/${BINARY_NAME}_linux_arm64_v8.0"
GOOS=linux GOARCH=arm64 go build -o "dist/${BINARY_NAME}_linux_arm64_v8.0/${BINARY_NAME}" .

# Build for Windows
echo "Building for windows/amd64..."
mkdir -p "dist/${BINARY_NAME}_windows_amd64_v1"
GOOS=windows GOARCH=amd64 go build -o "dist/${BINARY_NAME}_windows_amd64_v1/${BINARY_NAME}.exe" .

echo "Building for windows/arm64..."
mkdir -p "dist/${BINARY_NAME}_windows_arm64_v8.0"
GOOS=windows GOARCH=arm64 go build -o "dist/${BINARY_NAME}_windows_arm64_v8.0/${BINARY_NAME}.exe" .

echo "All builds completed successfully!"
echo "Binaries are available in the dist/ directory" 