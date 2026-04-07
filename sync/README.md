# DXRK Config

Repositorio centralizado para sincronizar skills y configuraciones de **Dxrk** a todos los agentes de IA.

## Repos related

- **Dxrk-Hex**: https://github.com/Dxrk777/Dxrk-Hex (configs, skills, sync system)
- **Dxrk-Config**: https://github.com/Dxrk777/Dxrk-OpenCode-Config (sync target repo)

## Estructura

```
.dxrk-config/
├── sync/
│   └── sync.sh              ← Script de sincronización
├── skills/                  ← 137 skills centralizadas
│   ├── dxrk-golang/
│   ├── dxrk-react-patterns/
│   ├── dxrk-ui-ux-pro-max/
│   └── ...
└── personas/               ← Configs por agente
    ├── opencode/
    │   └── AGENTS.md
    ├── cursor/
    │   ├── rules-dxrk.mdc
    │   └── gentle-ai.mdc
    ├── vscode/
    │   ├── dxrk.instructions.md
    │   └── Dxrk-AI.instructions.md
    └── claude/
        └── CLAUDE.md
```

## Uso

1. **Clonar** este repo en `~/.dxrk-config/`:
   ```bash
   git clone https://github.com/Dxrk777/Dxrk-Config.git ~/.dxrk-config
   ```

2. **Sincronizar** a todos los agentes:
   ```bash
   cd ~/.dxrk-config/sync
   ./sync.sh --all
   ```

3. **Sincronizar** solo un agente específico:
   ```bash
   ./sync.sh --cursor   # Solo Cursor
   ./sync.sh --opencode # Solo OpenCode
   ./sync.sh --claude   # Solo Claude Code
   ```

## Agentes Soportados

- **OpenCode** → `~/.config/opencode/`
- **Cursor** → `~/.cursor/`
- **Claude Code** → `~/.claude/`
- **VSCode** → `~/.config/Code/User/`
- **Windsurf** → `~/.windsurf/`

## Agregar Nuevo Agente

1. Agregar la ruta en `sync.sh` en el array `AGENT_DIRS`
2. Agregar la lógica de sync en la función `sync_agent()`
3. Ejecutar `./sync.sh --all`

## Workflow

1. Hacer cambios en `~/.dxrk-config/`
2. Commit y push a GitHub
3. En otra máquina: `git pull` + `./sync.sh --all`

---

**Dxrk** — *"Tu mentor digital. No te doy la respuesta — te ayudo a encontrarla."*