# Usage

← [Back to README](../README.md)

---

## Persona Modes

| Persona | ID | Description |
|---------|-----|-------------|
| Gentleman | `gentleman` | Teaching-oriented mentor persona — pushes back on bad practices, explains the why |
| Neutral | `neutral` | Same teacher, same philosophy, no regional language — warm and professional |
| Custom | `custom` | Bring your own persona instructions |

---

## Interactive TUI

Just run it — the Bubbletea TUI guides you through agent selection, components, skills, and presets:

```bash
dxrk
```

---

## CLI Commands

### install

First-time setup — detects your tools, configures agents, injects all components:

```bash
# Full ecosystem for multiple agents
dxrk install \
  --agent claude-code,opencode,gemini-cli \
  --preset full-gentleman

# Minimal setup for Cursor
dxrk install \
  --agent cursor \
  --preset minimal

# Pick specific components and skills
dxrk install \
  --agent claude-code \
  --component engram,sdd,skills,context7,persona,permissions \
  --skill go-testing,skill-creator,branch-pr,issue-creation \
  --persona gentleman

# Dry-run first (preview plan without applying changes)
dxrk install --dry-run \
  --agent claude-code,opencode \
  --preset full-gentleman
```

### sync

Refresh managed assets to the current version. Use after `brew upgrade dxrk` or when you want your local configs aligned with the latest release. Does NOT reinstall binaries (engram, GGA) — only updates prompt content, skills, MCP configs, and SDD orchestrators.

```bash
# Sync all installed agents
dxrk sync

# Sync specific agents only
dxrk sync --agent cursor --agent windsurf

# Sync a specific component
dxrk sync --component sdd
dxrk sync --component skills
dxrk sync --component engram
```

Sync is safe and idempotent — running it twice produces no changes the second time.

### update / upgrade

Check for and install new versions of `dxrk` itself:

```bash
# Check if a newer version is available
dxrk update

# Upgrade to the latest release (downloads new binary, replaces current)
dxrk upgrade
```

After upgrading, run `dxrk sync` to refresh all managed assets to the new version's content.

### version

```bash
dxrk version
dxrk --version
dxrk -v
```

---

## CLI Flags (install)

| Flag | Description |
|------|-------------|
| `--agent`, `--agents` | Agents to configure (comma-separated) |
| `--component`, `--components` | Components to install (comma-separated) |
| `--skill`, `--skills` | Skills to install (comma-separated) |
| `--persona` | Persona mode: `gentleman`, `neutral`, `custom` |
| `--preset` | Preset: `full-gentleman`, `ecosystem-only`, `minimal`, `custom` |
| `--dry-run` | Preview the install plan without applying changes |

## CLI Flags (sync)

| Flag | Description |
|------|-------------|
| `--agent`, `--agents` | Agents to sync (defaults to all installed agents) |
| `--component` | Sync a specific component only: `sdd`, `engram`, `context7`, `skills`, `gga`, `permissions`, `theme` |
| `--include-permissions` | Include permissions sync (opt-in) |
| `--include-theme` | Include theme sync (opt-in) |

---

## Typical Workflow

```bash
# First time: install everything
brew install dxrk-team/tap/dxrk
dxrk install --agent claude-code,cursor --preset full-gentleman

# After a new release: upgrade + sync
brew upgrade dxrk
dxrk sync

# Adding a new agent later
dxrk install --agent windsurf --preset full-gentleman
```

---

## Dependency Management

`dxrk` auto-detects prerequisites before installation and provides platform-specific guidance:

- **Detected tools**: git, curl, node, npm, brew, go
- **Version checks**: validates minimum versions where applicable
- **Platform-aware hints**: suggests `brew install`, `apt install`, `pacman -S`, `dnf install`, or `winget install` depending on your OS
- **Node LTS alignment**: on apt/dnf systems, Node.js hints use NodeSource LTS bootstrap before package install
- **Dependency-first approach**: detects what's installed, calculates what's needed, shows the full dependency tree before installing anything, then verifies each dependency after installation
