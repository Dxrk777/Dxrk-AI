<div align="center">

<img width="1600" height="757" alt="Dxrk Hex Logo" src="https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/assets/logo.png" />

<h1>DXRK HEX</h1>

<p><strong>DARK HEX SYSTEM — One command. Any agent. Any OS. PROTOCOL ACTIVE.</strong></p>

<p>
<a href="https://github.com/Dxrk777/Dxrk-Hex/releases"><img src="https://img.shields.io/github/v/release/Dxrk777/Dxrk-Hex?color=5dfc8e&label=Version" alt="Release"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
<img src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white" alt="Go 1.24+">
<img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey" alt="Platform">
<img src="https://img.shields.io/badge/Status-DXRK%20INDEX-5dfc8e" alt="Status">
<a href="https://github.com/Gentleman-Programming/gentle-ai/releases"><img src="https://img.shields.io/badge/Upstream-gentle--ai-blue" alt="gentle-ai upstream"></a>
</p>

> **📊 Dxrk Index System**: Dxrk Hex usa un sistema de versioning basado en porcentaje que sigue a [gentle-ai](https://github.com/Gentleman-Programming/gentle-ai) automáticamente. Cada release de gentle-ai incrementa el Dxrk Index.
> 
> | gentle-ai | Dxrk Hex | Incremento |
> |------------|----------|------------|
> | v1.15.6 | 045.27% | +0.05% (patch) |
> | v1.16.0 | 045.82% | +0.50% (minor) |
> | v2.0.0 | 055.82% | +10.00% (major) |

</div>

---

## What It Does

This is NOT an AI agent installer. Most agents are easy to install. This is an **ecosystem configurator** -- it takes whatever AI coding agent(s) you use and supercharges them with the Dxrk Hex stack: persistent memory, Spec-Driven Development workflow, curated coding skills, MCP servers, an AI provider switcher, a teaching-oriented persona with security-first permissions, and per-phase model assignment so each SDD step can run on a different model.

**Before**: "I installed Claude Code / OpenCode / Cursor, but it's just a chatbot that writes code."

**After**: Your agent now has memory, skills, workflow, MCP tools, and a persona that actually teaches you.

### 8 Supported Agents

| Agent | Delegation Model | Key Feature |
|-------|:---:|---|
| **Claude Code** | Full (Task tool) | Sub-agents, output styles |
| **OpenCode** | Full (multi-mode overlay) | Per-phase model routing |
| **Gemini CLI** | Full (experimental) | Custom agents in `~/.gemini/agents/` |
| **Cursor** | Full (native subagents) | 9 SDD agents in `~/.cursor/agents/` |
| **VS Code Copilot** | Full (runSubagent) | Parallel execution |
| **Codex** | Solo-agent | CLI-native, TOML config |
| **Windsurf** | Solo-agent | Plan Mode, Code Mode, native workflows |
| **Antigravity** | Solo-agent + Mission Control | Built-in Browser/Terminal sub-agents |

---

## Quick Start

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.ps1 | iex
```

This downloads the latest release for your platform and launches the interactive TUI. No Go toolchain required.

---

## Install

### Homebrew (macOS / Linux)

```bash
brew tap Dxrk777/tap
brew install dxrk
```

### Go install (any platform with Go 1.24+)

```bash
go install github.com/Dxrk777/Dxrk-Hex/cmd/dxrk@latest
```

### Windows (PowerShell)

```powershell
# Option 1: PowerShell installer (downloads binary from GitHub Releases)
irm https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.ps1 | iex

# Option 2: Go install (requires Go 1.24+)
go install github.com/Dxrk777/Dxrk-Hex/cmd/dxrk@latest
```

### From releases

Download the binary for your platform from [GitHub Releases](https://github.com/Dxrk777/Dxrk-Hex/releases).

---

## Documentation

| Topic | Description |
|-------|-------------|
| [Intended Usage](docs/intended-usage.md) | How dxrk is meant to be used — the mental model |
| [Agents](docs/agents.md) | Supported agents, feature matrix, config paths, and per-agent notes |
| [Components, Skills & Presets](docs/components.md) | All components, Dxrk behavior, skill catalog, and preset definitions |
| [Usage](docs/usage.md) | Persona modes, interactive TUI, CLI flags, and dependency management |
| [Platforms](docs/platforms.md) | Supported platforms, Windows notes, security verification, config paths |
| [Architecture & Development](docs/architecture.md) | Codebase layout, testing, and relationship to Dxrk.Dots |

---

## Version System

Dxrk Hex uses a **percentage-based version system**:

| Version | Status |
|---------|--------|
| `000.01%` | Initial Release |
| `001.00%` | Core Installer |
| `010.00%` | Skills System |
| `050.00%` | Multi-Platform |
| `100.00%` | MVP Achieved |

---

<div align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
<p><strong>// DARK HEX SYSTEM — PROTOCOL ACTIVE //</strong></p>
</div>
