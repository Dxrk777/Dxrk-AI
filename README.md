<div align="center">

<img width="1600" height="757" alt="Dxrk Hex Logo" src="https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/assets/logo.png" />

<h1>DXRK HEX</h1>

<p><strong>DARK HEX SYSTEM — One command. Any agent. Any OS. PROTOCOL ACTIVE.</strong></p>

<p>
<a href="https://github.com/Dxrk777/Dxrk-Hex/releases"><img src="https://img.shields.io/github/v/release/Dxrk777/Dxrk-Hex?color=5dfc8e&label=Version" alt="Release"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
<img src="https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white" alt="Go 1.25">
<img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey" alt="Platform">
<a href="https://github.com/Dxrk777/Dxrk-Hex/actions/workflows/test.yml"><img src="https://img.shields.io/github/actions/workflow/status/Dxrk777/Dxrk-Hex/test.yml?branch=main&label=Tests" alt="Tests"></a>
<a href="https://goreportcard.com/report/github.com/Dxrk777/Dxrk-Hex"><img src="https://goreportcard.com/badge/github.com/Dxrk777/Dxrk-Hex" alt="Go Report Card"></a>
<a href="https://codecov.io/gh/Dxrk777/Dxrk-Hex"><img src="https://img.shields.io/codecov/c/github/Dxrk777/Dxrk-Hex?token=placeholder" alt="Coverage"></a>
<a href="https://github.com/Dxrk777/Dxrk-Hex/releases"><img src="https://img.shields.io/badge/Versioning-Dxrk%20Index-blue" alt="Dxrk Index"></a>
</p>

> **📊 Dxrk Index System**: Dxrk Hex usa un sistema de versioning basado en porcentaje. El Dxrk Index refleja el progreso del proyecto en una escala de 0% a 100%.
> 
> | Index | Milestone | Descripción |
> |-------|-----------|-------------|
> | 000.01% | Initial Release | Primeros tests |
> | 010.00% | Skills System | Sistema de skills básico |
> | 050.00% | Multi-Platform | Soporte multiplataforma |
> | 100.00% | MVP Achieved | Producto mínimo viable |

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

## 🧠 Dxrk Hex Brain

The Brain is your **unified command center** that integrates memory, commands, email, and remote control.

### TUI (Interactive)

Launch the TUI and select **🧠 Brain** from the menu:
```bash
dxrk
# Navigate to 🧠 Brain
```

### CLI (Command Line)

Use the brain directly from terminal:
```bash
# Ask anything
dxrk brain "what agents are installed?"

# Execute shell commands
dxrk brain run "git status"
dxrk brain run "npm install"

# View history
dxrk brain history

# Check status
dxrk brain status

# Search memory
dxrk brain remember "last install"

# Send email (configure SMTP first)
dxrk brain email to user@example.com subject Hello body "Hi there!"

# List agents
dxrk brain agents

# Get help
dxrk brain help
```

### Brain Features

| Feature | Description |
|---------|-------------|
| **Memory** | Persistent command history stored in `~/.dxrk/memory` |
| **Commands** | Execute shell commands with timeout and safety |
| **Email** | Send email notifications via SMTP |
| **Remote** | Control via Telegram, Discord, WhatsApp |
| **Skills** | 60+ coding patterns integrated |

---

## Install

### Homebrew (macOS / Linux)

```bash
brew tap Dxrk777/tap
brew install dxrk
```

### Go install (any platform with Go 1.25+)

```bash
go install github.com/Dxrk777/Dxrk-Hex/cmd/dxrk@latest
```

### Windows (PowerShell)

```powershell
# Option 1: Scoop (recommended for Windows)
scoop bucket add dxrk https://github.com/Dxrk777/scoop-bucket
scoop install dxrk

# Option 2: PowerShell installer (downloads binary from GitHub Releases)
irm https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.ps1 | iex

# Option 3: Go install (requires Go 1.25+)
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
