#!/usr/bin/env bash
set -e

VERSION="v0.1.1"
OS=$(uname -s)
ARCH=$(uname -m)

# Map arch names
if [ "$ARCH" = "x86_64" ]; then
  ARCH="x86_64"
elif [ "$ARCH" = "arm64" ]; then
  ARCH="arm64"
fi

# Check for GitHub token
if [ -z "$GITHUB_TOKEN" ]; then
  echo "Error: GITHUB_TOKEN environment variable is required for private repository"
  exit 1
fi

ASSET_NAME="neko-cli_${OS}_${ARCH}.tar.gz"
REPO="nekoman-hq/neko-cli"

# Get asset ID from GitHub API using jq for reliable parsing
ASSET_ID=$(curl -s -H "Authorization: token $GITHUB_TOKEN" \
  "https://api.github.com/repos/$REPO/releases/tags/$VERSION" \
  | grep -B 2 "\"name\": \"$ASSET_NAME\"" \
  | grep '"id"' \
  | head -n 1 \
  | grep -o '[0-9]\+')

if [ -z "$ASSET_ID" ]; then
  echo "Error: Asset $ASSET_NAME not found in release $VERSION"
  exit 1
fi

echo "Downloading $ASSET_NAME..."

# Download using API with Accept header for binary content
TMP=$(mktemp -d)
curl -L -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/octet-stream" \
  "https://api.github.com/repos/$REPO/releases/assets/$ASSET_ID" \
  | tar -xz -C "$TMP"

chmod +x "$TMP/neko-cli"
sudo mv "$TMP/neko-cli" /usr/local/bin/neko
rm -rf "$TMP"

echo "neko-cli $VERSION installed successfully!"