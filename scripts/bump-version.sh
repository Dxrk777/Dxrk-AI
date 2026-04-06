#!/bin/bash
# bump-version.sh - Bump version, update changelog, create git tag
# Usage: ./scripts/bump-version.sh [major|minor|patch] [description]
# Example: ./scripts/bump-version.sh minor "add new AI architecture skills"

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get bump type (default: minor)
BUMP_TYPE=${1:-minor}
DESCRIPTION=${2:-"updates"}

# Get current version from goreleaser.yaml or git tag
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/v//' || echo "0.1.0")

# Parse version parts
IFS='.' read -ra VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR="${VERSION_PARTS[0]}"
MINOR="${VERSION_PARTS[1]:-0}"
PATCH="${VERSION_PARTS[2]:-0}"

# Bump version
case "$BUMP_TYPE" in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo -e "${RED}Invalid bump type: $BUMP_TYPE${NC}"
        echo "Usage: $0 [major|minor|patch] [description]"
        exit 1
        ;;
esac

NEW_VERSION="$MAJOR.$MINOR.$PATCH"
DATE=$(date '+%Y-%m-%d')

echo -e "${GREEN}Bumping version: ${CURRENT_VERSION} → ${NEW_VERSION}${NC}"

# Update homebrew-dxrk.rb
if [ -f "homebrew-dxrk.rb" ]; then
    sed -i "s/version \"[0-9.]*\"/version \"$NEW_VERSION\"/" homebrew-dxrk.rb
    sed -i "s|releases/download/v[0-9.]*/|releases/download/v$NEW_VERSION/|" homebrew-dxrk.rb
    echo -e "${GREEN}✓ Updated homebrew-dxrk.rb${NC}"
fi

# Add to CHANGELOG.md
CHANGELOG_ENTRY="## [$NEW_VERSION] - $DATE

### $DESCRIPTION
"

# Find position after first "## [" and before "## [Unreleased]" if exists
if grep -q "## \[Unreleased\]" CHANGELOG.md; then
    sed -i "/## \[Unreleased\]/a\\
$CHANGELOG_ENTRY\\
---" CHANGELOG.md
else
    # Add after "## [X.Y.Z]" header if exists, otherwise prepend
    if grep -q "^## \[" CHANGELOG.md; then
        sed -i "0,/^## \[.*\]/a\\
$CHANGELOG_ENTRY\\
---" CHANGELOG.md
    else
        sed -i "1i\\
$CHANGELOG_ENTRY\\
---\\
" CHANGELOG.md
    fi
fi
echo -e "${GREEN}✓ Updated CHANGELOG.md${NC}"

# Git operations
git add homebrew-dxrk.rb CHANGELOG.md 2>/dev/null || true
git add -A 2>/dev/null || true

# Create commit and tag
git commit -m "chore(release): bump version to v$NEW_VERSION

$DESCRIPTION"

git tag -a "v$NEW_VERSION" -m "Release v$NEW_VERSION: $DESCRIPTION"

echo ""
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  Version bumped to ${NEW_VERSION}${NC}"
echo -e "${GREEN}  Commit created with tag v$NEW_VERSION${NC}"
echo -e "${YELLOW}  Run: git push && git push --tags${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
