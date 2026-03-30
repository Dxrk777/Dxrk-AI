<div align="center">

<img width="3276" height="1280" alt="image" src="https://github.com/user-attachments/assets/3a3e4ae1-b9f4-4ce9-8fd0-3833812beb99" />

<h1>Dxrk</h1>

<p><strong>One command. Any agent. Any OS. The Dxrk ecosystem -- configured and ready.</strong></p>

<p>
<a href="https://github.com/Dxrk777777/Dxrk-Hex/releases"><img src="https://img.shields.io/github/v/release/dxrk777/Dxrk-Hex" alt="Release"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
<img src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white" alt="Go 1.24+">
<img src="https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-lightgrey" alt="Platform">
</p>

</div>

---

## What It Does

This is NOT an AI agent installer. Most agents are easy to install. This is an **ecosystem configurator** -- it takes whatever AI coding agent(s) you use and supercharges them with the Dxrk stack: persistent memory, Spec-Driven Development workflow, curated coding skills, MCP servers, an AI provider switcher, a teaching-oriented persona with security-first permissions, and per-phase model assignment so each SDD step can run on a different model.

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
curl -fsSL https://raw.githubusercontent.com/dxrk777/Dxrk-Hex/main/scripts/install-dxrk.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/dxrk777/Dxrk-Hex/main/scripts/install-dxrk.ps1 | iex
```

This downloads the latest release for your platform and launches the interactive TUI. No Go toolchain required.

---

## Install

### Homebrew (macOS / Linux)

```bash
brew tap dxrk/tap
brew install dxrk
```

### Go install (any platform with Go 1.24+)

```bash
go install github.com/Dxrk777777/Dxrk-Hex/cmd/dxrk@latest
```

### Windows (PowerShell)

```powershell
# Option 1: PowerShell installer (downloads binary from GitHub Releases)
irm https://raw.githubusercontent.com/dxrk777/Dxrk-Hex/main/scripts/install-dxrk.ps1 | iex

# Option 2: Go install (requires Go 1.24+)
go install github.com/Dxrk777777/Dxrk-Hex/cmd/dxrk@latest
```

### From releases

Download the binary for your platform from [GitHub Releases](https://github.com/Dxrk777777/Dxrk-Hex/releases).

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

<div align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
</div>
