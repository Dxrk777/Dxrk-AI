---
name: dxrk-homebrew-release
description: >
  Release workflow for Dxrk777 homebrew-tap projects (GGA, Dxrk-AI).
  Trigger: When user asks to release, bump version, update homebrew, or publish a new version.
license: Apache-2.0
metadata:
  author: dxrk777
  version: "3.0"
---

## When to Use

- User asks to "release", "bump", "publish", or "update homebrew"
- User mentions a new version number
- User says "homebrew-tap" or "formula"
- After creating a git tag in GGA or Dxrk-AI repos

Before running any `gh release create` or `gh release edit` command with markdown content, also load `release-note-safety`.

## Supported Projects

| Project | Repo | Formula | Tag Format | Type |
|---------|------|---------|------------|------|
| GGA | `gentleman-guardian-angel` | `gga.rb` | `V{version}` (e.g., `V2.6.2`) | Tarball (builds from source) |
| Dxrk-AI | `Dxrk-AI` | `dxrk-ai.rb` | `v{version}` (e.g., `v2.5.1`) | Pre-built binaries |

---

## Dxrk-AI Release Process (Pre-built Binaries)

### Step 1: Build Binaries

```bash
cd installer
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dxrk-installer-darwin-amd64 ./cmd/dxrk-installer
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dxrk-installer-darwin-arm64 ./cmd/dxrk-installer
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dxrk-installer-linux-amd64 ./cmd/dxrk-installer
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dxrk-installer-linux-arm64 ./cmd/dxrk-installer
```

### Step 2: Commit, Tag, and Push

```bash
git add -A
git commit -m "feat/fix: description"
git push origin main
git tag v{VERSION}
git push origin v{VERSION}
```

### Step 3: Create GitHub Release

```bash
gh release create v{VERSION} \
  installer/dxrk-installer-darwin-amd64 \
  installer/dxrk-installer-darwin-arm64 \
  installer/dxrk-installer-linux-amd64 \
  installer/dxrk-installer-linux-arm64 \
  --title "v{VERSION}" \
  --notes "## Changes
- {description of changes}"
```

### Step 4: Get SHA256 of Binaries

```bash
shasum -a 256 installer/dxrk-installer-darwin-arm64 installer/dxrk-installer-darwin-amd64 installer/dxrk-installer-linux-amd64 installer/dxrk-installer-linux-arm64
```

### Step 5: Update Formula

Update `homebrew-tap/Formula/dxrk-ai.rb`:

```ruby
class DxrkAi < Formula
  desc "Interactive TUI installer for Dxrk-AI development environment"
  homepage "https://github.com/Dxrk777/Dxrk-AI"
  version "{VERSION}"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/Dxrk777/Dxrk-AI/releases/download/v#{version}/dxrk-installer-darwin-arm64"
      sha256 "{SHA256_DARWIN_ARM64}"
    end
    on_intel do
      url "https://github.com/Dxrk777/Dxrk-AI/releases/download/v#{version}/dxrk-installer-darwin-amd64"
      sha256 "{SHA256_DARWIN_AMD64}"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/Dxrk777/Dxrk-AI/releases/download/v#{version}/dxrk-installer-linux-amd64"
      sha256 "{SHA256_LINUX_AMD64}"
    end
  end

  def install
    if OS.mac? && Hardware::CPU.arm?
      bin.install "dxrk-installer-darwin-arm64" => "dxrk-ai"
    elsif OS.mac? && Hardware::CPU.intel?
      bin.install "dxrk-installer-darwin-amd64" => "dxrk-ai"
    elsif OS.linux? && Hardware::CPU.intel?
      bin.install "dxrk-installer-linux-amd64" => "dxrk-ai"
    end
  end

  test do
    system "#{bin}/dxrk-ai", "--help"
  end
end
```

### Step 6: Commit to Both Repos

```bash
# In Dxrk-AI repo
git add homebrew-tap/Formula/dxrk-ai.rb
git commit -m "chore(homebrew): bump version to v{VERSION}"
git push origin main

# In homebrew-tap repo
cd /tmp && rm -rf homebrew-tap
git clone git@github.com:Dxrk777/homebrew-tap.git
cp {path-to}/Dxrk-AI/homebrew-tap/Formula/dxrk-ai.rb /tmp/homebrew-tap/Formula/
cd /tmp/homebrew-tap
git add -A
git commit -m "chore: bump dxrk-ai to v{VERSION}"
git push origin main
```

---

## GGA Release Process (Tarball - Builds from Source)

### Step 1: Verify Tag Exists

```bash
git tag --list | tail -5
```

### Step 2: Get SHA256 of Tarball

```bash
curl -sL https://github.com/Dxrk777/gentleman-guardian-angel/archive/refs/tags/V{VERSION}.tar.gz | shasum -a 256
```

### Step 3: Update Formula

Update `homebrew-tap/Formula/gga.rb`:

```ruby
url "https://github.com/Dxrk777/gentleman-guardian-angel/archive/refs/tags/V{VERSION}.tar.gz"
sha256 "{NEW_SHA256}"
version "{VERSION}"
```

### Step 4: Commit and Push

```bash
cd ~/work/homebrew-tap
git add -A
git commit -m "chore: bump gga to V{VERSION}"
git push
```

---

## Quick Reference Commands

```bash
# Build Dxrk-AI binaries
cd installer && GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dxrk-installer-darwin-amd64 ./cmd/dxrk-installer && GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dxrk-installer-darwin-arm64 ./cmd/dxrk-installer && GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dxrk-installer-linux-amd64 ./cmd/dxrk-installer && GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dxrk-installer-linux-arm64 ./cmd/dxrk-installer

# SHA256 for binaries
shasum -a 256 installer/dxrk-installer-*

# SHA256 for GGA tarball
curl -sL https://github.com/Dxrk777/gentleman-guardian-angel/archive/refs/tags/V{VERSION}.tar.gz | shasum -a 256

# Create GitHub release with binaries
gh release create v{VERSION} installer/dxrk-installer-* --title "v{VERSION}" --notes "## Changes"
```

---

## Checklist

### Dxrk-AI
- [ ] Binaries built for all platforms (darwin-amd64, darwin-arm64, linux-amd64, linux-arm64)
- [ ] Changes committed and pushed to main
- [ ] Tag created and pushed (v{VERSION})
- [ ] GitHub release created with binaries attached
- [ ] SHA256 computed for all binaries
- [ ] Formula updated in Dxrk-AI repo
- [ ] Formula copied to homebrew-tap repo and pushed

### GGA
- [ ] Tag exists (V{VERSION})
- [ ] SHA256 computed from tarball
- [ ] Formula updated in homebrew-tap
- [ ] Committed and pushed
