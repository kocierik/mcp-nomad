#!/bin/bash

set -e

chmod +x ./packages/build.sh
./packages/build.sh

# Mac
mkdir -p ./packages/npm-mcp-nomad-darwin-x64/bin
cp dist/mcp-nomad_darwin_amd64_v1/mcp-nomad ./packages/npm-mcp-nomad-darwin-x64/bin/mcp-nomad
chmod +x ./packages/npm-mcp-nomad-darwin-x64/bin/mcp-nomad
mkdir -p ./packages/npm-mcp-nomad-darwin-arm64/bin
cp dist/mcp-nomad_darwin_arm64_v8.0/mcp-nomad ./packages/npm-mcp-nomad-darwin-arm64/bin/mcp-nomad
chmod +x ./packages/npm-mcp-nomad-darwin-arm64/bin/mcp-nomad

# Linux
mkdir -p ./packages/npm-mcp-nomad-linux-x64/bin
cp dist/mcp-nomad_linux_amd64_v1/mcp-nomad ./packages/npm-mcp-nomad-linux-x64/bin/mcp-nomad
chmod +x ./packages/npm-mcp-nomad-linux-x64/bin/mcp-nomad
mkdir -p ./packages/npm-mcp-nomad-linux-arm64/bin
cp dist/mcp-nomad_linux_arm64_v8.0/mcp-nomad ./packages/npm-mcp-nomad-linux-arm64/bin/mcp-nomad
chmod +x ./packages/npm-mcp-nomad-linux-arm64/bin/mcp-nomad

# Windows
mkdir -p ./packages/npm-mcp-nomad-win32-x64/bin
cp dist/mcp-nomad_windows_amd64_v1/mcp-nomad.exe ./packages/npm-mcp-nomad-win32-x64/bin/mcp-nomad.exe
mkdir -p ./packages/npm-mcp-nomad-win32-arm64/bin
cp dist/mcp-nomad_windows_arm64_v8.0/mcp-nomad.exe ./packages/npm-mcp-nomad-win32-arm64/bin/mcp-nomad.exe

cd packages/npm-mcp-nomad-darwin-x64
npm publish --access public

cd ../npm-mcp-nomad-darwin-arm64
npm publish --access public

cd ../npm-mcp-nomad-linux-x64
npm publish --access public

npm publish --access public
cd ../npm-mcp-nomad-linux-arm64

cd ../npm-mcp-nomad-win32-x64
npm publish --access public

cd ../npm-mcp-nomad-win32-arm64
npm publish --access public

cd ../npm-mcp-nomad
npm publish --access public

cd -
