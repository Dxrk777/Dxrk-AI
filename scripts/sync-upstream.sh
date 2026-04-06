#!/bin/bash
# ============================================================================
# Dxrk AI - Sync Upstream Script
# ============================================================================
# Este script sincroniza Dxrk AI con Dxrk-AI (upstream).
# Se ejecuta automáticamente cuando Dxrk-AI libera una nueva versión.
#
# Uso:
#   ./sync-upstream.sh <Dxrk-AI-tag>
#
# Ejemplo:
#   ./sync-upstream.sh v1.16.0
#
# El script:
# 1. Agrega upstream como remote si no existe
# 2. Fetch del último upstream
# 3. Hace merge del upstream/main
# 4. Resuelve conflictos (prioriza cambios de Dxrk)
# 5. Crea commit de sync
# 6. Push a origin main
#
# Luego el workflow de Release crea el nuevo tag y release.
# ============================================================================

set -e # Salir en caso de error

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuración
UPSTREAM_REPO="https://github.com/Gentleman-Programming/Dxrk-AI.git"
UPSTREAM_NAME="Dxrk-AI"
BRANCH="${BRANCH:-main}"

# Archivos de estado
PERCENT_FILE=".current-percent"
LAST_TAG_FILE=".last-upstream-tag"

# ============================================================================
# Funciones de utilidad
# ============================================================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# ============================================================================
# Verificaciones iniciales
# ============================================================================

verify_git() {
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        log_error "No estamos en un repositorio git"
        exit 1
    fi
    log_info "Repositorio git verificado"
}

# ============================================================================
# Gestión del remote upstream
# ============================================================================

setup_upstream() {
    # Verificar si upstream ya existe
    if git remote get-url "$UPSTREAM_NAME" >/dev/null 2>&1; then
        CURRENT_URL=$(git remote get-url "$UPSTREAM_NAME")
        if [[ "$CURRENT_URL" != "$UPSTREAM_REPO" ]]; then
            log_warn "Upstream existe con URL diferente, actualizando..."
            git remote set-url "$UPSTREAM_NAME" "$UPSTREAM_REPO"
        else
            log_info "Upstream ya configurado: $UPSTREAM_REPO"
        fi
    else
        log_info "Agregando upstream: $UPSTREAM_REPO"
        git remote add "$UPSTREAM_NAME" "$UPSTREAM_REPO"
    fi

    # Fetch del upstream
    log_info "Fetching upstream..."
    git fetch "$UPSTREAM_NAME" --tags --quiet
}

# ============================================================================
# Merge del upstream
# ============================================================================

merge_upstream() {
    log_info "Merging ${UPSTREAM_NAME}/${BRANCH}..."

    # Configurar estrategia para resolver conflictos automáticamente
    # Prioriza los cambios locales (Dxrk) sobre los remotos (Dxrk-AI)
    git config merge.ours.driver true

    # Hacer merge (los conflictos se resuelven automáticamente a favor nuestro)
    if git merge "${UPSTREAM_NAME}/${BRANCH}" --no-edit --no-ff; then
        log_success "Merge completado sin conflictos"
    else
        log_warn "Hubo conflictos, resolviendo automáticamente a favor de Dxrk..."
        # Resolver conflictos a favor de nuestro branch
        git checkout --ours .
        git add -A
        git commit -m "chore: resolve merge conflicts favoring Dxrk changes
        
        Auto-resolved conflicts during upstream sync from Dxrk-AI
        PRIORITY: Dxrk customizations over upstream defaults"
        log_success "Conflictos resueltos"
    fi
}

# ============================================================================
# Cálculo de la nueva versión
# ============================================================================

calculate_new_version() {
    local upstream_tag="$1"

    # Obtener versión actual o usar default
    if [[ -f "$PERCENT_FILE" ]]; then
        CURRENT_PERCENT=$(cat "$PERCENT_FILE")
    else
        CURRENT_PERCENT="0.03"
    fi

    # Extraer tipo de release de Dxrk-AI
    # Formato: v1.15.6 -> major=1, minor=15, patch=6
    VERSION_NUM=$(echo "$upstream_tag" | sed 's/v//')
    MAJOR=$(echo "$VERSION_NUM" | cut -d. -f1)
    MINOR=$(echo "$VERSION_NUM" | cut -d. -f2)
    PATCH=$(echo "$VERSION_NUM" | cut -d. -f3)

    # Calcular incremento
    if [[ "$MINOR" == "0" && "$PATCH" == "0" ]]; then
        # Major release: +10.00%
        INCREMENT=10.00
        RELEASE_TYPE="major"
    elif [[ "$PATCH" == "0" ]]; then
        # Minor release: +0.50%
        INCREMENT=0.50
        RELEASE_TYPE="minor"
    else
        # Patch release: +0.05%
        INCREMENT=0.05
        RELEASE_TYPE="patch"
    fi

    # Calcular nuevo porcentaje (usar awk en vez de bc - disponible en GitHub Actions)
    NEW_PERCENT=$(awk "BEGIN {printf \"%.2f\", $CURRENT_PERCENT + $INCREMENT}")

    # Formatear como tag Dxrk (vXXX.XX%)
    NEW_TAG=$(printf "v%06.2f%%" "$NEW_PERCENT")

    # Guardar estado
    echo "$NEW_PERCENT" >"$PERCENT_FILE"
    echo "$upstream_tag" >"$LAST_TAG_FILE"

    # Exportar para el workflow
    echo "NEW_TAG=$NEW_TAG"
    echo "NEW_PERCENT=$NEW_PERCENT"
    echo "RELEASE_TYPE=$RELEASE_TYPE"

    log_success "Nueva versión calculada: $NEW_TAG"
    log_info "Tipo de release: $RELEASE_TYPE (+$INCREMENT%)"
}

# ============================================================================
# Commit y push
# ============================================================================

push_changes() {
    log_info "Verificando cambios para commit..."

    if git diff --quiet && git diff --cached --quiet; then
        log_info "No hay cambios para hacer commit"
        return 0
    fi

    log_info "Haciendo commit de sync..."
    git add -A

    # Mensaje de commit informativo
    COMMIT_MSG=$(
        cat <<EOF
chore(sync): merge upstream from Dxrk-AI

Sync con la última versión de Dxrk-AI.
Auto-generado por sync-upstream.sh.

Release type: ${RELEASE_TYPE:-unknown}
Upstream tag: ${1:-latest}

[Este commit mantiene Dxrk AI sincronizado con upstream]
EOF
    )

    git commit -m "$COMMIT_MSG"

    log_info "Push a origin..."
    git push origin "$BRANCH"

    log_success "Cambios empujados a origin/${BRANCH}"
}

# ============================================================================
# Función principal
# ============================================================================

main() {
    echo ""
    echo "╔═══════════════════════════════════════════════════════════╗"
    echo "║         DXRK AI - UPSTREAM SYNC                        ║"
    echo "║         Sincronizando con Dxrk-AI                       ║"
    echo "╚═══════════════════════════════════════════════════════════╝"
    echo ""

    # Verificar argumentos
    if [[ -z "$1" ]]; then
        log_error "Uso: $0 <Dxrk-AI-tag>"
        log_error "Ejemplo: $0 v1.15.7"
        exit 1
    fi

    UPSTREAM_TAG="$1"

    # Ejecutar pasos
    verify_git
    setup_upstream
    merge_upstream
    calculate_new_version "$UPSTREAM_TAG"
    push_changes "$UPSTREAM_TAG"

    echo ""
    echo "╔═══════════════════════════════════════════════════════════╗"
    echo "║         SYNC COMPLETADO                                 ║"
    echo "║                                                         ║"
    echo "║   Nuevo tag listo para release: ${NEW_TAG}               ║"
    echo "║                                                         ║"
    echo "║   Para crear el release:                                ║"
    echo "║   git tag ${NEW_TAG}                                     ║"
    echo "║   git push origin ${NEW_TAG}                             ║"
    echo "╚═══════════════════════════════════════════════════════════╝"
    echo ""

    # Output para GitHub Actions
    if [[ -n "$GITHUB_OUTPUT" ]]; then
        echo "new_tag=${NEW_TAG}" >>"$GITHUB_OUTPUT"
        echo "new_percent=${NEW_PERCENT}" >>"$GITHUB_OUTPUT"
        echo "release_type=${RELEASE_TYPE}" >>"$GITHUB_OUTPUT"
    fi
}

# Ejecutar
main "$@"
