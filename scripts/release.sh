#!/usr/bin/env bash
# shellkit release script
# Usage: ./scripts/release.sh [major|minor|patch]

set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

VERSION_TYPE=${1:-patch}
VERSION_FILE="VERSION"

# Get current version
if [[ -f "$VERSION_FILE" ]]; then
    CURRENT=$(cat "$VERSION_FILE")
else
    CURRENT=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.0")
fi

# Parse version components
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT"

# Bump version
case "$VERSION_TYPE" in
    major) NEW_VERSION="$((MAJOR + 1)).0.0" ;;
    minor) NEW_VERSION="${MAJOR}.$((MINOR + 1)).0" ;;
    patch) NEW_VERSION="${MAJOR}.${MINOR}.$((PATCH + 1))" ;;
    *) echo "Usage: $0 [major|minor|patch]"; exit 1 ;;
esac

echo "==> Releasing v${NEW_VERSION} (was v${CURRENT})..."

# Update VERSION file
echo "$NEW_VERSION" > "$VERSION_FILE"
echo "    Updated VERSION"

# Update version badge in README.md
if [[ -f "README.md" ]]; then
    sed -i '' "s/version-[0-9]*\.[0-9]*\.[0-9]*/version-${NEW_VERSION}/" README.md
    echo "    Updated README.md badge"
fi

# Generate changelog
if command -v git-cliff &>/dev/null; then
    git-cliff --tag "v${NEW_VERSION}" --output CHANGELOG.md
    echo "    Generated CHANGELOG.md"
else
    echo "    Warning: git-cliff not installed, skipping changelog" >&2
fi

# Commit and tag
git add VERSION CHANGELOG.md README.md 2>/dev/null || true
git commit -m "chore(release): v${NEW_VERSION}"
git tag -a "v${NEW_VERSION}" -m "Release v${NEW_VERSION}"

# Push and create GitHub Release
echo ""
echo "==> Pushing to origin..."
git push origin main --tags

if command -v gh &>/dev/null; then
    echo "==> Creating GitHub Release..."
    gh release create "v${NEW_VERSION}" --title "v${NEW_VERSION}" --notes-file CHANGELOG.md
    echo "    Created GitHub Release"
fi

echo ""
echo "==> Release v${NEW_VERSION} complete!"
