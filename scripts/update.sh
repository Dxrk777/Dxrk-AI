#!/bin/bash
# =============================================================================
# Dxrk AI - Auto Update Script
# =============================================================================
# Detecta cómo está instalado dxrk y actualiza automáticamente.
#
# Uso: ./scripts/update.sh [--check]
#
# Opciones:
#   --check    Solo verifica si hay actualización disponible (no actualiza)
#   --force    Forzar actualización aunque esté actualizado
# =============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BOLD='\033[1m'
DIM='\033[2m'

CHECK_ONLY=false
FORCE_UPDATE=false

# Parsear argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
    --check)
        CHECK_ONLY=true
        shift
        ;;
    --force)
        FORCE_UPDATE=true
        shift
        ;;
    -h | --help)
        echo "Uso: $0 [--check] [--force]"
        echo ""
        echo "Opciones:"
        echo "  --check    Solo verificar si hay actualización (no actualiza)"
        echo "  --force    Forzar actualización aunque ya esté actualizado"
        exit 0
        ;;
    *)
        echo "Opción desconocida: $1"
        exit 1
        ;;
    esac
done

# Funciones de utilidad
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Banner
echo ""
echo -e "${BOLD}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${BOLD}║${NC}          ${GREEN}Dxrk AI Auto-Updater${NC}              ${BOLD}║${NC}"
echo -e "${BOLD}╚══════════════════════════════════════════════════╝${NC}"
echo ""

# Detectar si dxrk está instalado
detect_installation() {
    log_info "Detectando instalación de dxrk..."

    if command -v dxrk &>/dev/null; then
        DXRK_PATH=$(command -v dxrk)
        log_success "Dxrk encontrado en: $DXRK_PATH"

        # Detectar método de instalación
        if [[ "$DXRK_PATH" == "/usr/local/bin/dxrk" ]] || [[ "$DXRK_PATH" == "/opt/homebrew/bin/dxrk" ]]; then
            INSTALL_METHOD="homebrew"
            log_info "Método de instalación: Homebrew"
        elif [[ "$DXRK_PATH" == *"go/bin"* ]]; then
            INSTALL_METHOD="go"
            log_info "Método de instalación: Go Install"
        else
            INSTALL_METHOD="binary"
            log_info "Método de instalación: Binario manual"
        fi

        # Obtener versión actual
        CURRENT_VERSION=$(dxrk --version 2>/dev/null || echo "unknown")
        log_info "Versión actual: $CURRENT_VERSION"

        return 0
    else
        log_warn "Dxrk no está instalado"
        INSTALL_METHOD="none"
        return 1
    fi
}

# Obtener última versión desde GitHub
get_latest_version() {
    log_info "Verificando última versión en GitHub..."

    # Intentar con gh CLI
    if command -v gh &>/dev/null; then
        LATEST_TAG=$(gh release view --repo Dxrk777/Dxrk --json tagName --jq '.tagName' 2>/dev/null || echo "")
    fi

    # Fallback a curl
    if [[ -z "$LATEST_TAG" ]]; then
        LATEST_TAG=$(curl -s "https://api.github.com/repos/Dxrk777/Dxrk-AI/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "")
    fi

    if [[ -n "$LATEST_TAG" ]]; then
        log_success "Última versión: $LATEST_TAG"
        return 0
    else
        log_error "No se pudo obtener la última versión"
        return 1
    fi
}

# Comparar versiones
needs_update() {
    if [[ "$CURRENT_VERSION" == "unknown" ]] || [[ "$CURRENT_VERSION" == "$LATEST_TAG" ]]; then
        return 1 # No necesita actualización
    fi
    return 0 # Necesita actualización
}

# Actualizar con Homebrew
update_homebrew() {
    log_info "Actualizando con Homebrew..."

    if [[ "$CHECK_ONLY" == true ]]; then
        log_info "Modo verificación - no se actualiza"
        return 0
    fi

    echo ""
    log_info "Ejecutando: brew upgrade dxrk..."
    echo ""

    if brew upgrade dxrk; then
        log_success "Actualización completada!"
        echo ""
        log_info "Nueva versión: $(dxrk --version)"
        return 0
    else
        log_error "Error al actualizar con Homebrew"
        return 1
    fi
}

# Actualizar con Go
update_go() {
    log_info "Actualizando con Go..."

    if [[ "$CHECK_ONLY" == true ]]; then
        log_info "Modo verificación - no se actualiza"
        return 0
    fi

    echo ""
    log_info "Ejecutando: go install github.com/Dxrk777/Dxrk-AI/cmd/dxrk@latest"
    echo ""

    if go install github.com/Dxrk777/Dxrk-AI/cmd/dxrk@latest; then
        log_success "Actualización completada!"
        echo ""
        log_info "Nueva versión: $(dxrk --version)"
        return 0
    else
        log_error "Error al actualizar con Go"
        return 1
    fi
}

# Actualizar binario manual
update_binary() {
    log_info "Actualizando binario..."

    if [[ "$CHECK_ONLY" == true ]]; then
        log_info "Modo verificación - no se actualiza"
        return 0
    fi

    # Detectar SO y arquitectura
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    # Normalizar arquitectura
    case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64 | aarch64)
        ARCH="arm64"
        ;;
    esac

    # Normalizar SO
    case $OS in
    darwin)
        OS="darwin"
        ;;
    linux)
        OS="linux"
        ;;
    mingw* | cygwin* | msys*)
        OS="windows"
        ;;
    *)
        log_error "SO no soportado: $OS"
        return 1
        ;;
    esac

    # Determinar extensión
    if [[ "$OS" == "windows" ]]; then
        EXT="zip"
    else
        EXT="tar.gz"
    fi

    FILENAME="dxrk_${LATEST_TAG#v}_${OS}_${ARCH}.${EXT}"
    DOWNLOAD_URL="https://github.com/Dxrk777/Dxrk-AI/releases/download/${LATEST_TAG}/${FILENAME}"

    log_info "Descargando: $FILENAME"

    # Crear directorio temporal
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"

    # Descargar
    if command -v curl &>/dev/null; then
        curl -sL "$DOWNLOAD_URL" -o "dxrk.${EXT}"
    elif command -v wget &>/dev/null; then
        wget -q "$DOWNLOAD_URL" -O "dxrk.${EXT}"
    else
        log_error "Necesitas curl o wget para descargar"
        return 1
    fi

    # Extraer
    if [[ "$EXT" == "tar.gz" ]]; then
        tar -xzf "dxrk.${EXT}"
    else
        unzip -q "dxrk.${EXT}"
    fi

    # Instalar
    log_info "Instalando en: $DXRK_PATH"

    if [[ -w "$DXRK_PATH" ]] || sudo -v; then
        if [[ -w "$DXRK_PATH" ]]; then
            cp dxrk "$DXRK_PATH"
        else
            sudo cp dxrk "$DXRK_PATH"
            sudo chmod +x "$DXRK_PATH"
        fi

        log_success "Actualización completada!"
        echo ""
        log_info "Nueva versión: $(dxrk --version)"

        # Limpiar
        cd /
        rm -rf "$TMP_DIR"

        return 0
    else
        log_error "No se pudo escribir en $DXRK_PATH"
        log_info "Puedes actualizar manualmente con:"
        log_info "  sudo cp dxrk $DXRK_PATH"
        cd /
        rm -rf "$TMP_DIR"
        return 1
    fi
}

# Instalar desde cero si no está instalado
install_fresh() {
    log_warn "Dxrk no está instalado"
    echo ""
    log_info "Ejecutando instalador..."
    echo ""

    if [[ "$CHECK_ONLY" == true ]]; then
        log_info "Modo verificación - no se instala"
        return 0
    fi

    # Detectar SO
    OS=$(uname -s)

    if [[ "$OS" == "Darwin" ]]; then
        echo "curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-AI/main/scripts/install-dxrk.sh | bash"
        curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-AI/main/scripts/install-dxrk.sh | bash
    elif [[ "$OS" == "Linux" ]]; then
        curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-AI/main/scripts/install-dxrk.sh | bash
    else
        log_info "Descarga manual desde:"
        log_info "  https://github.com/Dxrk777/Dxrk-AI/releases"
    fi
}

# Main
main() {
    # Detectar instalación
    if ! detect_installation; then
        install_fresh
        exit 0
    fi

    # Obtener última versión
    if ! get_latest_version; then
        log_error "No se pudo verificar la versión"
        exit 1
    fi

    echo ""

    # Verificar si necesita actualización
    if needs_update; then
        log_warn "Hay una nueva versión disponible!"
        echo ""
        echo -e "  ${DIM}Actual:${NC} $CURRENT_VERSION"
        echo -e "  ${GREEN}Nueva:${NC}  $LATEST_TAG"
        echo ""

        if [[ "$CHECK_ONLY" == true ]]; then
            log_info "Usa $0 --force para actualizar"
            exit 0
        fi

        if [[ "$FORCE_UPDATE" == false ]]; then
            read -p "¿Actualizar? [y/N] " -n 1 -r
            echo ""
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                log_info "Actualización cancelada"
                exit 0
            fi
        fi

        echo ""

        # Actualizar según método
        case $INSTALL_METHOD in
        homebrew)
            update_homebrew
            ;;
        go)
            update_go
            ;;
        binary)
            update_binary
            ;;
        *)
            log_error "Método de instalación desconocido"
            exit 1
            ;;
        esac
    else
        log_success "Ya tienes la última versión!"
        echo ""
        log_info "Versión: $CURRENT_VERSION"
        exit 0
    fi
}

main
