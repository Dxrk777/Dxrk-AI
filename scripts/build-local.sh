#!/usr/bin/env bash
# Build dxrk from source — works on any Linux/macOS with Go installed
set -euo pipefail

echo "🔥 Building Dxrk AI from source..."

# Check Go
if ! command -v go &>/dev/null; then
    echo "❌ Go not found. Install with: sudo apt install golang-go"
    exit 1
fi

GO_VERSION=$(go version | grep -oP '\d+\.\d+(\.\d+)?' | head -1)
echo "✅ Go $GO_VERSION found"

# Run from the repo root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Tidy deps
echo "📦 Downloading dependencies..."
go mod tidy

# Build
echo "🔨 Building..."
CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=dev" -o dxrk ./cmd/dxrk

echo "✅ Binary built: $(pwd)/dxrk"
echo ""
echo "To install system-wide:"
echo "  sudo mv dxrk /usr/local/bin/"
echo ""
echo "Or add to your PATH:"
echo "  export PATH=\$PATH:$(pwd)"
echo ""
echo "Test it:"
echo "  ./dxrk version"
echo "  ./dxrk help"
