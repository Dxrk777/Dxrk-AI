#!/bin/bash
# Dxrk Config Sync - Sincroniza skills y configs a todos los agentes
# Uso: ./sync.sh [--all | --opencode | --cursor | --claude | --windsurf]

set -e

DXRK_CONFIG_DIR="$HOME/.dxrk-config"
SKILLS_SOURCE="$DXRK_CONFIG_DIR/skills"
PERSONAS_SOURCE="$DXRK_CONFIG_DIR/personas"

# Directorios de destino por agente
declare -A AGENT_DIRS=(
  ["opencode"]="$HOME/.config/opencode"
  ["cursor"]="$HOME/.cursor"
  ["claude"]="$HOME/.claude"
  ["vscode"]="$HOME/.config/Code/User"
  ["windsurf"]="$HOME/.windsurf"
)

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Verificar que existe la config central
check_source() {
  if [ ! -d "$SKILLS_SOURCE" ]; then
    log_error "No se encontró $SKILLS_SOURCE"
    log_info "Clona el repo: git clone https://github.com/Dxrk777/Dxrk-Config.git $DXRK_CONFIG_DIR"
    exit 1
  fi
}

# Sincronizar skills
sync_skills() {
  local agent=$1
  local dest_dir="$AGENT_DIRS[$agent]/skills"
  
  if [ -z "$dest_dir" ]; then
    log_error "Agent desconocido: $agent"
    return 1
  fi
  
  # Crear directorio si no existe
  mkdir -p "$dest_dir"
  
  # Copiar skills (excluyendo archivos no deseados)
  rsync -av --exclude='__pycache__' --exclude='*.pyc' \
    "$SKILLS_SOURCE/" "$dest_dir/" 2>/dev/null || \
    cp -r "$SKILLS_SOURCE"/* "$dest_dir/" 2>/dev/null || true
  
  log_info "Skills sincronizadas para $agent"
}

# Sincronizar agentes específicos
sync_agent() {
  local agent=$1
  local base_dir="$AGENT_DIRS[$agent]"
  
  if [ -z "$base_dir" ]; then
    log_error "Agent desconocido: $agent"
    return 1
  fi
  
  # Sync skills
  sync_skills "$agent"
  
  # Sync personas/config según el agente
  case $agent in
    opencode)
      if [ -f "$PERSONAS_SOURCE/opencode/AGENTS.md" ]; then
        cp "$PERSONAS_SOURCE/opencode/AGENTS.md" "$base_dir/"
      fi
      ;;
    cursor)
      if [ -d "$PERSONAS_SOURCE/cursor/rules" ]; then
        mkdir -p "$base_dir/rules"
        cp -r "$PERSONAS_SOURCE/cursor/rules/"* "$base_dir/rules/"
      fi
      ;;
    claude)
      if [ -f "$PERSONAS_SOURCE/claude/CLAUDE.md" ]; then
        cp "$PERSONAS_SOURCE/claude/CLAUDE.md" "$base_dir/"
      fi
      ;;
    vscode)
      if [ -d "$PERSONAS_SOURCE/vscode/prompts" ]; then
        mkdir -p "$base_dir/prompts"
        cp -r "$PERSONAS_SOURCE/vscode/prompts/"* "$base_dir/prompts/"
      fi
      ;;
  esac
  
  log_info "Sincronización completa para $agent"
}

# Sincronizar todos los agentes
sync_all() {
  log_info "Sincronizando todos los agentes..."
  for agent in "${!AGENT_DIRS[@]}"; do
    sync_agent "$agent"
  done
  log_info "✓ Todos los agentes sincronizados"
}

# Main
main() {
  check_source
  
  case "${1:-all}" in
    --all|-a)
      sync_all
      ;;
    --opencode|-o)
      sync_agent "opencode"
      ;;
    --cursor|-c)
      sync_agent "cursor"
      ;;
    --claude|-l)
      sync_agent "claude"
      ;;
    --vscode|-v)
      sync_agent "vscode"
      ;;
    --windsurf|-w)
      sync_agent "windsurf"
      ;;
    --help|-h)
      echo "Dxrk Config Sync"
      echo "Uso: $0 [opción]"
      echo ""
      echo "Opciones:"
      echo "  --all, -a      Sincronizar todos los agentes (default)"
      echo "  --opencode, -o Sincronizar solo OpenCode"
      echo "  --cursor, -c   Sincronizar solo Cursor"
      echo "  --claude, -l   Sincronizar solo Claude Code"
      echo "  --vscode, -v   Sincronizar solo VSCode"
      echo "  --windsurf, -w Sincronizar solo Windsurf"
      echo "  --help, -h     Mostrar esta ayuda"
      ;;
    *)
      log_error "Opción desconocida: $1"
      echo "Usa --help para ver las opciones disponibles"
      exit 1
      ;;
  esac
}

main "$@"