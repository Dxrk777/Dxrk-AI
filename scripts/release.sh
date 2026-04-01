#!/bin/bash
# =============================================================================
# Dxrk Hex - Release Script
# =============================================================================
# Usage: ./scripts/release.sh <version>
# Example: ./scripts/release.sh 0.15%
# =============================================================================

set -e

VERSION="${1:-}"

if [ -z "$VERSION" ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 0.15%"
    echo ""
    echo "This will:"
    echo "  1. Create a git tag v<version>"
    echo "  2. Push to origin"
    echo "  3. GitHub Actions will build and release automatically"
    exit 1
fi

TAG="v${VERSION}"

echo "🚀 Creating release ${TAG}..."
echo ""

# Check for uncommitted changes
if [ -n "$(git status --porcelain)" ]; then
    echo "⚠️  Warning: You have uncommitted changes:"
    git status --short
    echo ""
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Aborted."
        exit 1
    fi
fi

# Create tag
echo "📌 Creating tag ${TAG}..."
git tag "${TAG}"

# Push tag (this triggers the release workflow)
echo "📤 Pushing tag to origin..."
git push origin "${TAG}"

echo ""
echo "✅ Release ${TAG} triggered!"
echo ""
echo "Watch the release at:"
echo "  https://github.com/Dxrk777/Dxrk-Hex/actions"
echo ""
echo "The release will include:"
echo "  - Binary for Linux (amd64, arm64)"
echo "  - Binary for macOS (amd64, arm64)"
echo "  - Binary for Windows (amd64, arm64)"
echo "  - Homebrew formula update"
