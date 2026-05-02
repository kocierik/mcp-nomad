#!/bin/bash

set -e
REPO_ROOT=$(cd "$(dirname "$0")/.." && pwd)
cd "$REPO_ROOT"

# CI passes NPM_TOKEN via env but .npmrc with _authToken lives under packages/*/.
# npm whoami must see a registry token from the repo root, so write a disposable userconfig.
NPM_CI_AUTH_RC=""
EXIT_cleanup() {
  restore_npmrc
  [ -n "$NPM_CI_AUTH_RC" ] && rm -f "$NPM_CI_AUTH_RC"
}
trap EXIT_cleanup EXIT

NPM_PACKAGES=(
  npm-mcp-nomad-darwin-x64
  npm-mcp-nomad-darwin-arm64
  npm-mcp-nomad-linux-x64
  npm-mcp-nomad-linux-arm64
  npm-mcp-nomad-win32-x64
  npm-mcp-nomad-win32-arm64
  npm-mcp-nomad
)

NPMRC_BACKUP_DIR="$REPO_ROOT/packages/.npmrc-backups"
restore_npmrc() {
  if [ "${RESTORE_NPMRC:-0}" = "1" ] && [ -d "$NPMRC_BACKUP_DIR" ]; then
    for pkg in "${NPM_PACKAGES[@]}"; do
      [ -f "$NPMRC_BACKUP_DIR/$pkg" ] && mv "$NPMRC_BACKUP_DIR/$pkg" "$REPO_ROOT/packages/$pkg/.npmrc"
    done
    rmdir "$NPMRC_BACKUP_DIR" 2>/dev/null || true
  fi
}
if [ -n "${NPM_TOKEN:-}" ]; then
  NPM_CI_AUTH_RC=$(mktemp)
  chmod 600 "$NPM_CI_AUTH_RC"
  printf '//registry.npmjs.org/:_authToken=%s\n' "$NPM_TOKEN" >"$NPM_CI_AUTH_RC"
  export NPM_CONFIG_USERCONFIG="$NPM_CI_AUTH_RC"
fi

# Fail early if not logged in to npm (avoids confusing 404 after build)
if ! npm whoami --registry=https://registry.npmjs.org &>/dev/null; then
  echo "Error: Not logged in to npm. Run 'npm login' or set a valid NPM_TOKEN secret (valid publish token)."
  exit 1
fi

# Package dirs have .npmrc with _authToken=${NPM_TOKEN}. When NPM_TOKEN is not set (local publish),
# move .npmrc out of package dirs so npm uses ~/.npmrc and the backup is not included in the tarball.
# Note: npm may require 2FA or a granular token with "Bypass 2FA" to publish; use NPM_TOKEN in that case.
if [ -z "${NPM_TOKEN:-}" ]; then
  echo "NPM_TOKEN not set: using npm login credentials (hiding package .npmrc for publish)."
  RESTORE_NPMRC=1
  mkdir -p "$NPMRC_BACKUP_DIR"
  for pkg in "${NPM_PACKAGES[@]}"; do
    [ -f "packages/$pkg/.npmrc" ] && mv "packages/$pkg/.npmrc" "$NPMRC_BACKUP_DIR/$pkg"
  done
fi

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

for pkg in "${NPM_PACKAGES[@]}"; do
  cd "$REPO_ROOT/packages/$pkg"
  npm publish --access public
done
