#!/usr/bin/env bash

INSPECTOR_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$INSPECTOR_DIR"

pnpm install

pnpm exec mcp-inspector go run ../../main.go "$@"
