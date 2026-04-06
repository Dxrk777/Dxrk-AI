# Architecture

## Overview

Dxrk AI is a self-hosted AI coding agent with custom features.

## Project Structure

```
Dxrk/
├── cmd/dxrk/              # Main entry point
├── internal/
│   ├── agents/            # Agent definitions (Claude, OpenCode, etc.)
│   ├── app/               # Application orchestration
│   ├── backup/            # Backup/restore functionality
│   ├── brain/             # 🆕 Brain module (memory, commands, email)
│   ├── catalog/           # Agent catalog
│   ├── cli/               # CLI commands
│   ├── components/         # Ecosystem components
│   │   ├── dxrk/         # Dxrk runtime
│   │   ├── engram/       # Engram memory system
│   │   ├── mcp/          # MCP server
│   │   ├── sdd/          # Spec-Driven Development
│   │   ├── skills/       # Coding skills
│   │   └── theme/        # Terminal theme
│   ├── connector/        # Remote control (Telegram, Discord, WhatsApp)
│   ├── model/             # Domain models
│   ├── opencode/         # OpenCode configuration
│   ├── pipeline/         # Installation pipeline
│   ├── planner/          # Dependency planning
│   ├── state/             # Persistent state
│   ├── system/            # System detection
│   ├── tui/               # Terminal UI
│   ├── update/            # Update checking
│   ├── vault/             # Encryption module
│   └── verify/            # Installation verification
├── assets/                # Static assets
├── docs/                  # Documentation
├── scripts/               # Installation scripts
└── skills/               # AI coding skills
```

## Core Components

### Brain Module

The Brain is the central orchestrator for Dxrk AI features:

```
┌─────────────────────────────────────────────────────┐
│                    Dxrk AI Brain                    │
├─────────────────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────┐  ┌─────────┐            │
│  │ Memory  │  │Commands │  │  Email  │            │
│  └────┬────┘  └────┬────┘  └────┬────┘            │
│       └────────────┼────────────┘                  │
│                    ▼                               │
│              ┌─────────┐                            │
│              │  Think  │                            │
│              └────┬────┘                            │
│                   │                                 │
├───────────────────┼─────────────────────────────────┤
│                   ▼                                 │
│  ┌─────────────────────────────────────────┐       │
│  │      BrainIntegration (Connector)         │       │
│  └────────────────────┬────────────────────┘       │
│                       │                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐         │
│  │ Telegram │  │ Discord  │  │ WhatsApp │         │
│  └──────────┘  └──────────┘  └──────────┘         │
└─────────────────────────────────────────────────────┘
```

### Installation Pipeline

```
Detection → Agents → Persona → Preset → Components → Review → Installing → Complete
```

### TUI Architecture

The TUI follows the Bubble Tea pattern:

```go
type Model struct {
    Screen Screen
    // ...
}

func (m Model) Init() tea.Cmd     // Initialize
func (m Model) Update(msg)       // Handle events
func (m Model) View() string     // Render
```

## Design Patterns

### Container-Presentational

Components are split into:
- **Presentational**: Render UI, handle user input
- **Container**: Handle data fetching, state management

### Dependency Injection

Testable code uses interface injection:

```go
var (
    osStatModelCache = os.Stat  // Package-level for testing
)
```

## Data Flow

```
User Input → TUI Model → Brain → Memory/Commander/Emailer → Response
                    ↓
              Connector → Telegram/Discord/WhatsApp
```

## Testing Strategy

- Unit tests for pure functions
- Integration tests for I/O operations
- Mock external services for reproducibility

## Key Files

| File | Purpose |
|------|---------|
| `internal/app/app.go` | Main entry, CLI/TUI dispatch |
| `internal/brain/think.go` | Natural language processing |
| `internal/tui/model.go` | TUI state machine |
| `internal/pipeline/` | Installation orchestration |

## Extension Points

### Adding a New Agent

1. Add to `internal/agents/`
2. Add to catalog in `internal/catalog/`
3. Add to model in `internal/model/`
4. Add tests

### Adding a New Brain Command

1. Add handler in `think.go`
2. Add tests
3. Update CLI if needed

### Adding a New Skill

1. Add to `internal/assets/skills/`
2. Document in skill registry
