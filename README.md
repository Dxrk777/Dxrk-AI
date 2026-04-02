<div align="center">

<img width="800" height="400" alt="Dxrk Hex Logo" src="https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/assets/logo.png" />

# DXRK HEX

<p>
<a href="https://github.com/Dxrk777/Dxrk-Hex/releases"><img src="https://img.shields.io/github/v/release/Dxrk777/Dxrk-Hex?color=ff0040&label=Version" alt="Release"></a>
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-red.svg" alt="License: MIT"></a>
<img src="https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white" alt="Go 1.25">
<a href="https://github.com/Dxrk777/Dxrk-Hex/actions/workflows/test.yml"><img src="https://img.shields.io/github/actions/workflow/status/Dxrk777/Dxrk-Hex/test.yml?branch=main&label=Tests" alt="Tests"></a>
</p>

<p><strong>Tu compañero digital 🔥 — One command. Any agent. Any OS.</strong></p>

</div>

---

## ⚡ Qué es Dxrk Hex?

**Dxrk Hex NO es un instalador de agentes de IA.** 

La mayoría de los agentes son fáciles de instalar. Esto es un **configurador de ecosistema** — transforma cualquier agente de IA que uses en una máquina perfectamente afinada:

- 🧠 **Memoria persistente** — tu agente recuerda todo
- 🎯 **Workflow SDD** — Spec-Driven Development de verdad
- 🛠️ **60+ skills** — patrones de código curados
- 🔌 **MCP servers** — herramientas que se expanden
- ⚡ **AI Provider Switcher** — cambia de proveedor en segundos
- 🎓 **Persona docente** — enseña mientras programa
- 🎛️ **Modelo por fase** — cada paso SDD usa el mejor modelo

### 💡 Antes vs Después

| Antes | Después |
|-------|---------|
| "Instalé Claude Code pero es solo un chatbot" | Tu agente tiene memoria, skills, workflow y enseña |

---

## 🚀 Quick Start

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.ps1 | iex
```

Sin Go, sin complicaciones. Descarga, instala, usa.

---

## 🤖 8 Agentes Soportados

| Agente | Modelo | Feature Principal |
|--------|:------:|------------------|
| **Claude Code** | Full (Task tool) | Sub-agentes, output styles |
| **OpenCode** | Full (multi-mode) | Routing por fase SDD |
| **Gemini CLI** | Full (experimental) | Custom agents en `~/.gemini/agents/` |
| **Cursor** | Full (native subagents) | 11 SDD agents en `~/.cursor/agents/` |
| **VS Code Copilot** | Full (runSubagent) | Ejecución paralela |
| **Codex** | Solo-agent | CLI-native, TOML config |
| **Windsurf** | Solo-agent | Plan Mode, Code Mode |
| **Antigravity** | Solo-agent + Mission | Browser/Terminal built-in |

---

## 🧠 Dxrk Hex Brain

Tu **centro de comando unificado** que integra memoria, comandos, email y control remoto.

### TUI Interactivo

```bash
dxrk
# Navega a 🧠 Brain
```

### CLI Directo

```bash
# Pregunta lo que quieras
dxrk brain "¿qué agentes están instalados?"

# Ejecuta comandos shell
dxrk brain run "git status"
dxrk brain run "npm install"

# Ver historial
dxrk brain history

# Buscar en memoria
dxrk brain remember "última instalación"

# Control remoto
dxrk brain telegram "reiniciar servidor"
```

---

## 📦 Instalación

### Homebrew (macOS / Linux)

```bash
brew tap Dxrk777/tap
brew install dxrk
```

### Go install (cualquier plataforma con Go 1.25+)

```bash
go install github.com/Dxrk777/Dxrk-Hex/cmd/dxrk@latest
```

### Windows

```powershell
# Scoop (recomendado)
scoop bucket add dxrk https://github.com/Dxrk777/scoop-bucket
scoop install dxrk

# PowerShell installer
irm https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.ps1 | iex
```

### Releases

Descarga el binary desde [GitHub Releases](https://github.com/Dxrk777/Dxrk-Hex/releases).

---

## 📚 Documentación

| Tema | Descripción |
|------|-------------|
| [Intended Usage](docs/intended-usage.md) | Cómo usar Dxrk Hex — el modelo mental |
| [Agents](docs/agents.md) | Agentes soportados, matrix de features, paths |
| [Components & Skills](docs/components.md) | Componentes, catálogo de skills, presets |
| [Usage](docs/usage.md) | Modos de persona, TUI interactivo, CLI |
| [Platforms](docs/platforms.md) | Plataformas soportadas, notas de Windows |

---

## 🎯 Sistema de Versiones

Dxrk Hex usa un **sistema de versionado porcentual**:

| Index | Status | Descripción |
|-------|--------|-------------|
| `000.01%` | Initial Release | Primeros tests |
| `010.00%` | Skills System | Sistema de skills básico |
| `050.00%` | Multi-Platform | Soporte multiplataforma |
| `100.00%` | MVP Achieved | Producto mínimo viable |

---

<div align="center">

**[MIT License](LICENSE)** — © 2026 Dxrk Hex

🔥 *Tu compañero digital*

</div>
