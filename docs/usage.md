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
[Dd]xrk-ai
```

---

## CLI Commands

### install

First-time setup — detects your tools, configures agents, injects all components:

```bash
# Full ecosystem for multiple agents
[Dd]xrk-ai install \
  --agent claude-code,opencode,gemini-cli \
  --preset full-gentleman

# Minimal setup for Cursor
[Dd]xrk-ai install \
  --agent cursor \
  --preset minimal

# Pick specific components and skills
[Dd]xrk-ai install \
  --agent claude-code \
  --component DxrkMemory,sdd,skills,context7,persona,permissions \
  --skill go-testing,skill-creator,branch-pr,issue-creation \
  --persona gentleman

# Dry-run first (preview plan without applying changes)
[Dd]xrk-ai install --dry-run \
  --agent claude-code,opencode \
  --preset full-gentleman
```

### sync

Refresh managed assets to the current version. Use after `brew upgrade [Dd]xrk-ai` or when you want your local configs aligned with the latest release. Does NOT reinstall binaries (DxrkMemory, GGA) — only updates prompt content, skills, MCP configs, and SDD orchestrators.

```bash
# Sync all installed agents
[Dd]xrk-ai sync

# Sync specific agents only
[Dd]xrk-ai sync --agent cursor --agent windsurf

# Sync a specific component
[Dd]xrk-ai sync --component sdd
[Dd]xrk-ai sync --component skills
[Dd]xrk-ai sync --component DxrkMemory
```

Sync is safe and idempotent — running it twice produces no changes the second time.

### update / upgrade

Check for and install new versions of `[Dd]xrk-ai` itself:

```bash
# Check if a newer version is available
[Dd]xrk-ai update

# Upgrade to the latest release (downloads new binary, replaces current)
[Dd]xrk-ai upgrade
```

After upgrading, run `[Dd]xrk-ai sync` to refresh all managed assets to the new version's content.

### version

```bash
[Dd]xrk-ai version
[Dd]xrk-ai --version
[Dd]xrk-ai -v
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
| `--component` | Sync a specific component only: `sdd`, `DxrkMemory`, `context7`, `skills`, `gga`, `permissions`, `theme` |
| `--include-permissions` | Include permissions sync (opt-in) |
| `--include-theme` | Include theme sync (opt-in) |

---

## Typical Workflow

```bash
# First time: install everything
brew install dxrk777/tap/[Dd]xrk-ai
[Dd]xrk-ai install --agent claude-code,cursor --preset full-gentleman

# After a new release: upgrade + sync
brew upgrade [Dd]xrk-ai
[Dd]xrk-ai sync

# Adding a new agent later
[Dd]xrk-ai install --agent windsurf --preset full-gentleman
```

---

## Dependency Management

`[Dd]xrk-ai` auto-detects prerequisites before installation and provides platform-specific guidance:

- **Detected tools**: git, curl, node, npm, brew, go
- **Version checks**: validates minimum versions where applicable
- **Platform-aware hints**: suggests `brew install`, `apt install`, `pacman -S`, `dnf install`, or `winget install` depending on your OS
- **Node LTS alignment**: on apt/dnf systems, Node.js hints use NodeSource LTS bootstrap before package install
- **Dependency-first approach**: detects what's installed, calculates what's needed, shows the full dependency tree before installing anything, then verifies each dependency after installation
