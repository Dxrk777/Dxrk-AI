<div align="center">

# 🔥 DXRK AI

<img width="200" height="200" src="https://avatars.githubusercontent.com/u/197162710?v=4" alt="Dxrk AI" style="border-radius: 50%; border: 4px solid #ff0040; box-shadow: 0 0 30px rgba(255,0,64,0.5);">

### Tu compañero digital 🔥

<p>

[![License](https://img.shields.io/badge/License-MIT-ff0040.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go)](https://golang.org)
[![Tests](https://img.shields.io/github/actions/workflow/status/Dxrk777/Dxrk-AI/test.yml?branch=main&label=Tests)](https://github.com/Dxrk777/Dxrk-AI/actions)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-1a1a1a)](https://github.com/Dxrk777/Dxrk-AI)

</p>

<p><strong>One command. Any agent. Any OS.</strong></p>

</div>

---

## ⚡ ¿Qué es Dxrk AI?

**Dxrk AI** transforma cualquier agente de IA en una máquina perfectamente afinada con memoria, skills, workflow SDD y una personalidad que enseña mientras programa.

| Feature | Descripción |
|---------|-------------|
| 🧠 **DxrkMemory** | Memoria persistente entre sesiones |
| 🎯 **Workflow SDD** | Spec-Driven Development real |
| 🛠️ **60+ Skills** | Patrones de código curados |
| 🔌 **MCP Servers** | Herramientas que se expanden |
| ⚡ **AI Provider Switcher** | Cambiá de proveedor en segundos |
| 🎓 **Persona Docente** | Enseña mientras programa |

---

## 🚀 Instalación

### Linux / Kali Linux (Go install) ⭐ Recomendado

```bash
# Instalar Go primero (si no lo tenés)
sudo apt install golang-go   # Debian/Ubuntu/Kali
# o
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Instalar dxrk
go install github.com/Dxrk777/Dxrk-AI/cmd/dxrk@latest

# Agregar al PATH si no está
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### Linux / macOS (Script automático)

```bash
curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-AI/main/scripts/install-dxrk.sh | bash
```

### macOS (Homebrew)

```bash
brew tap Dxrk777/tap && brew install dxrk
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/Dxrk777/Dxrk-AI/main/scripts/install-dxrk.ps1 | iex
```

### Compilar desde fuente

```bash
git clone https://github.com/Dxrk777/Dxrk-AI.git
cd Dxrk-AI
go build -o dxrk ./cmd/dxrk
sudo mv dxrk /usr/local/bin/
dxrk version
```

---

## 🤖 Agentes Soportados

| Agente | Tipo | Feature Principal |
|--------|:----:|------------------|
| **Claude Code** | Full | Sub-agentes, output styles |
| **OpenCode** | Full | Routing por fase SDD |
| **Gemini CLI** | Full | Custom agents |
| **Cursor** | Full | 11 SDD agents |
| **VS Code Copilot** | Full | Ejecución paralela |
| **Codex** | Full | CLI-native |
| **Windsurf** | Full | Plan/Code Modes |
| **Antigravity** | Full | Browser/Terminal |

---

## 🧠 Dxrk Brain

```bash
dxrk                              # TUI interactivo
dxrk install --dry-run            # Preview de instalación
dxrk install --agent claude-code  # Instalar agente específico
dxrk brain "¿qué agentes están instalados?"
dxrk sync                         # Sincronizar configuraciones
dxrk upgrade                      # Actualizar todo
dxrk version                      # Ver versión
```

---

## 📚 Docs

- [Uso](docs/usage.md) — Modos de persona, TUI, CLI
- [Agentes](docs/agents.md) — Agentes soportados
- [Componentes](docs/components.md) — Skills y presets
- [Arquitectura](docs/ARCHITECTURE.md) — Diseño interno
- [DxrkMemory](https://github.com/Dxrk777/dxrk-memory) — Motor de memoria

---

<div align="center">

**[MIT License](LICENSE)** — © 2026 **Dxrk AI**

🔥 *Tu compañero digital — 100% Dxrk*

</div>
