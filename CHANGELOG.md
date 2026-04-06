# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.2] - 2026-04-05

### Style

- gofmt formatting on go-rod example files


### Fixed

- go-rod example files excluded from Go build

## [1.1.0] - 2026-04-05

### 🚀 Major Update: DXRK Branding + AI Architecture Skills

Complete ecosystem overhaul with full DXRK branding and AI architecture skills.

### Added

#### 🧠 DXRK Agent Core
- Complete agent architecture knowledge base (renamed from dxrk-claude-code-internals)
- Dual-layer AsyncGenerator pattern
- 12-step permission pipeline
- Multi-agent coordination with async workers and mailbox
- Modular prompt system with BOUNDARY cache separator
- 5-phase self-healing auto-repair protocol
- Multi-provider LLM fallback strategy

#### 📚 AI Architecture Skills (10 new)
- `dxrk-agent-core` - Core architecture patterns
- `dxrk-agent-loop-architect` - Generator pattern
- `dxrk-tool-control-patterns` - Security pipeline
- `dxrk-multi-agent-coordinator` - Worker patterns
- `dxrk-prompt-registry` - Prompt composition
- `dxrk-auto-repair-protocol` - Error recovery
- `dxrk-llm-fallback-strategy` - LLM routing
- `dxrk-free-deployment` - Fly.io, Railway, Render
- `dxrk-free-memory-stack` - Supabase, Redis, Neon, Turso
- `dxrk-master-index` - Complete skill catalog

#### 🔄 Automated Orchestration
- Full auto-trigger system in AGENTS.md
- 123+ skills with automatic loading based on context
- Engram persistent memory integration
- SDD (Spec-Driven Development) workflow ready

### Changed

- All "Claude Code" and "Anthropic" references → DXRK branding
- AGENTS.md with comprehensive auto-load triggers
- Skills now properly orchestrated from ~/.claude/skills/

### Fixed

- Version bump script created: `scripts/bump-version.sh`

## [1.0.0] - 2026-04-01

### 🎉 Major Release

This is the first major release of Dxrk Hex with the Brain module, CLI integration, and comprehensive testing.

### Added

#### 🧠 Brain Module
- **Memory System**: Persistent command history stored in `~/.dxrk/memory`
- **Commander**: Secure shell command execution with timeout and safety checks
- **Email Integration**: SMTP email notifications
- **Webhook Support**: Discord, Slack, and Teams notifications
- **Engram Client**: HTTP client for Engram MCP server integration

#### 💻 CLI Commands
- `dxrk brain` - Unified command center
  - `dxrk brain help` - Show help
  - `dxrk brain status` - Check system status
  - `dxrk brain history` - View command history
  - `dxrk brain run <command>` - Execute shell commands
  - `dxrk brain agents` - List available agents
  - `dxrk brain remember <query>` - Search memory
  - `dxrk brain email to <addr> subject <sub> body <msg>` - Send emails

#### 🎨 TUI Integration
- Brain screen integrated in TUI menu (option 9)
- Interactive modes: chat, execute, email, status, history, configure
- Memory persistence across sessions

#### 🔗 Connector Module
- Remote control via Telegram, Discord, WhatsApp
- HTTP webhook handlers with proper error handling
- 62% test coverage

#### 📚 Documentation
- `docs/ARCHITECTURE.md` - Project structure and design patterns
- `docs/TROUBLESHOOTING.md` - Common issues and solutions
- `README.md` updated with Brain documentation
- Badges for tests, coverage, and Go Report Card

#### 🧪 Testing
- Brain module: 73% test coverage
- CLI brain tests: 25 test cases
- HTTP handler tests for connector
- Webhook and Engram client tests

### Changed

- **README.md**: Added badges, Brain documentation section
- **GoReleaser config**: Fixed homepage URL
- **Release workflow**: Ready for automated releases

### Fixed

- Brain TUI integration with real brain module
- Test assertions for brain CLI commands
- HTTP handler error responses

## [0.0.3] - Previous Release

- See [gentle-ai releases](https://github.com/Gentleman-Programming/gentle-ai/releases) for upstream changes

## Upstream Sync

Dxrk Hex es un fork de [gentle-ai](https://github.com/Gentleman-Programming/gentle-ai) y sincroniza automáticamente los cambios upstream diariamente via GitHub Actions.

### Installation

```bash
# Homebrew
brew install Dxrk777/tap/dxrk

# Go install
go install github.com/Dxrk777/Dxrk-Hex/cmd/dxrk@latest

# Script
curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.sh | bash
```

### Usage

```bash
# Interactive TUI
dxrk

# CLI install
dxrk install opencode

# Brain commands
dxrk brain status
dxrk brain run "git status"
dxrk brain help
```
