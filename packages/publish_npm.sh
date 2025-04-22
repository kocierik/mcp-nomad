#!/bin/bash

set -e

# Mac
mkdir -p ./packages/npm-mcp-nomad-darwin-x64/bin
cp dist/mcp-nomad-go_darwin_amd64_v1/mcp-nomad-go ./packages/npm-mcp-nomad-darwin-x64/bin/mcp-nomad-go
chmod +x ./packages/npm-mcp-nomad-darwin-x64/bin/mcp-nomad-go
mkdir -p ./packages/npm-mcp-nomad-darwin-arm64/bin
cp dist/mcp-nomad-go_darwin_arm64_v8.0/mcp-nomad-go ./packages/npm-mcp-nomad-darwin-arm64/bin/mcp-nomad-go
chmod +x ./packages/npm-mcp-nomad-darwin-arm64/bin/mcp-nomad-go

# Linux
mkdir -p ./packages/npm-mcp-nomad-linux-x64/bin
cp dist/mcp-nomad-go_linux_amd64_v1/mcp-nomad-go ./packages/npm-mcp-nomad-linux-x64/bin/mcp-nomad-go
chmod +x ./packages/npm-mcp-nomad-linux-x64/bin/mcp-nomad-go
mkdir -p ./packages/npm-mcp-nomad-linux-arm64/bin
cp dist/mcp-nomad-go_linux_arm64_v8.0/mcp-nomad-go ./packages/npm-mcp-nomad-linux-arm64/bin/mcp-nomad-go
chmod +x ./packages/npm-mcp-nomad-linux-arm64/bin/mcp-nomad-go

# Windows
mkdir -p ./packages/npm-mcp-nomad-win32-x64/bin
cp dist/mcp-nomad-go_windows_amd64_v1/mcp-nomad-go.exe ./packages/npm-mcp-nomad-win32-x64/bin/mcp-nomad-go.exe
mkdir -p ./packages/npm-mcp-nomad-win32-arm64/bin
cp dist/mcp-nomad-go_windows_arm64_v8.0/mcp-nomad-go.exe ./packages/npm-mcp-nomad-win32-arm64/bin/mcp-nomad-go.exe

cd packages/npm-mcp-nomad-darwin-x64
npm publish --access public

cd ../npm-mcp-nomad-darwin-arm64
npm publish --access public

cd ../npm-mcp-nomad-linux-x64
npm publish --access public

cd ../npm-mcp-nomad-linux-arm64
npm publish --access public

cd ../npm-mcp-nomad-win32-x64
npm publish --access public

cd ../npm-mcp-nomad-win32-arm64
npm publish --access public

cd ../npm-mcp-nomad
npm publish --access public

cd -
