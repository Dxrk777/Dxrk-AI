# Usage

← [Back to README](../README.md)

---

## Persona Modes

| Persona | ID | Description |
|---------|-----|-------------|
<<<<<<< HEAD
| Dxrk | `dxrk` | Teaching-oriented mentor persona — pushes back on bad practices, explains the why |
=======
| Gentleman | `gentleman` | Teaching-oriented mentor persona — pushes back on bad practices, explains the why |
>>>>>>> upstream/main
| Neutral | `neutral` | Same teacher, same philosophy, no regional language — warm and professional |
| Custom | `custom` | Bring your own persona instructions |

---

## Interactive TUI

Just run it — the Bubbletea TUI guides you through agent selection, components, skills, and presets:

```bash
<<<<<<< HEAD
dxrk
=======
gentle-ai
>>>>>>> upstream/main
```

---

## CLI Commands

### install

First-time setup — detects your tools, configures agents, injects all components:

```bash
# Full ecosystem for multiple agents
<<<<<<< HEAD
dxrk install \
  --agent claude-code,opencode,gemini-cli \
  --preset full-dxrk

# Minimal setup for Cursor
dxrk install \
=======
gentle-ai install \
  --agent claude-code,opencode,gemini-cli \
  --preset full-gentleman

# Minimal setup for Cursor
gentle-ai install \
>>>>>>> upstream/main
  --agent cursor \
  --preset minimal

# Pick specific components and skills
<<<<<<< HEAD
dxrk install \
  --agent claude-code \
  --component engram,sdd,skills,context7,persona,permissions \
  --skill go-testing,skill-creator,branch-pr,issue-creation \
  --persona dxrk

# Dry-run first (preview plan without applying changes)
dxrk install --dry-run \
  --agent claude-code,opencode \
  --preset full-dxrk
=======
gentle-ai install \
  --agent claude-code \
  --component engram,sdd,skills,context7,persona,permissions \
  --skill go-testing,skill-creator,branch-pr,issue-creation \
  --persona gentleman

# Dry-run first (preview plan without applying changes)
gentle-ai install --dry-run \
  --agent claude-code,opencode \
  --preset full-gentleman
>>>>>>> upstream/main
```

### sync

<<<<<<< HEAD
Refresh managed assets to the current version. Use after `brew upgrade dxrk` or when you want your local configs aligned with the latest release. Does NOT reinstall binaries (engram, Dxrk) — only updates prompt content, skills, MCP configs, and SDD orchestrators.

```bash
# Sync all installed agents
dxrk sync

# Sync specific agents only
dxrk sync --agent cursor --agent windsurf

# Sync a specific component
dxrk sync --component sdd
dxrk sync --component skills
dxrk sync --component engram
=======
Refresh managed assets to the current version. Use after `brew upgrade gentle-ai` or when you want your local configs aligned with the latest release. Does NOT reinstall binaries (engram, GGA) — only updates prompt content, skills, MCP configs, and SDD orchestrators.

```bash
# Sync all installed agents
gentle-ai sync

# Sync specific agents only
gentle-ai sync --agent cursor --agent windsurf

# Sync a specific component
gentle-ai sync --component sdd
gentle-ai sync --component skills
gentle-ai sync --component engram
>>>>>>> upstream/main
```

Sync is safe and idempotent — running it twice produces no changes the second time.

### update / upgrade

<<<<<<< HEAD
Check for and install new versions of `dxrk` itself:

```bash
# Check if a newer version is available
dxrk update

# Upgrade to the latest release (downloads new binary, replaces current)
dxrk upgrade
```

After upgrading, run `dxrk sync` to refresh all managed assets to the new version's content.
=======
Check for and install new versions of `gentle-ai` itself:

```bash
# Check if a newer version is available
gentle-ai update

# Upgrade to the latest release (downloads new binary, replaces current)
gentle-ai upgrade
```

After upgrading, run `gentle-ai sync` to refresh all managed assets to the new version's content.
>>>>>>> upstream/main

### version

```bash
<<<<<<< HEAD
dxrk version
dxrk --version
dxrk -v
=======
gentle-ai version
gentle-ai --version
gentle-ai -v
>>>>>>> upstream/main
```

---

## CLI Flags (install)

| Flag | Description |
|------|-------------|
| `--agent`, `--agents` | Agents to configure (comma-separated) |
| `--component`, `--components` | Components to install (comma-separated) |
| `--skill`, `--skills` | Skills to install (comma-separated) |
<<<<<<< HEAD
| `--persona` | Persona mode: `dxrk`, `neutral`, `custom` |
| `--preset` | Preset: `full-dxrk`, `ecosystem-only`, `minimal`, `custom` |
=======
| `--persona` | Persona mode: `gentleman`, `neutral`, `custom` |
| `--preset` | Preset: `full-gentleman`, `ecosystem-only`, `minimal`, `custom` |
>>>>>>> upstream/main
| `--dry-run` | Preview the install plan without applying changes |

## CLI Flags (sync)

| Flag | Description |
|------|-------------|
| `--agent`, `--agents` | Agents to sync (defaults to all installed agents) |
<<<<<<< HEAD
| `--component` | Sync a specific component only: `sdd`, `engram`, `context7`, `skills`, `dxrk`, `permissions`, `theme` |
=======
| `--component` | Sync a specific component only: `sdd`, `engram`, `context7`, `skills`, `gga`, `permissions`, `theme` |
>>>>>>> upstream/main
| `--include-permissions` | Include permissions sync (opt-in) |
| `--include-theme` | Include theme sync (opt-in) |

---

## Typical Workflow

```bash
# First time: install everything
<<<<<<< HEAD
brew install dxrk-programming/tap/dxrk
dxrk install --agent claude-code,cursor --preset full-dxrk

# After a new release: upgrade + sync
brew upgrade dxrk
dxrk sync

# Adding a new agent later
dxrk install --agent windsurf --preset full-dxrk
=======
brew install gentleman-programming/tap/gentle-ai
gentle-ai install --agent claude-code,cursor --preset full-gentleman

# After a new release: upgrade + sync
brew upgrade gentle-ai
gentle-ai sync

# Adding a new agent later
gentle-ai install --agent windsurf --preset full-gentleman
>>>>>>> upstream/main
```

---

## Dependency Management

<<<<<<< HEAD
`dxrk` auto-detects prerequisites before installation and provides platform-specific guidance:
=======
`gentle-ai` auto-detects prerequisites before installation and provides platform-specific guidance:
>>>>>>> upstream/main

- **Detected tools**: git, curl, node, npm, brew, go
- **Version checks**: validates minimum versions where applicable
- **Platform-aware hints**: suggests `brew install`, `apt install`, `pacman -S`, `dnf install`, or `winget install` depending on your OS
- **Node LTS alignment**: on apt/dnf systems, Node.js hints use NodeSource LTS bootstrap before package install
- **Dependency-first approach**: detects what's installed, calculates what's needed, shows the full dependency tree before installing anything, then verifies each dependency after installation
