#!/bin/bash
# =============================================================================
# Dxrk AI - Instalador de Actualización Automática
# =============================================================================
# Instala un servicio que actualiza dxrk automáticamente.
#
# Uso: ./scripts/install-auto-update.sh [--uninstall]
#
# Platforms:
#   - macOS: LaunchAgent (~/Library/LaunchAgents/)
#   - Linux: Cron job
# =============================================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

UNINSTALL=false

# Parsear argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
    --uninstall)
        UNINSTALL=true
        shift
        ;;
    -h | --help)
        echo "Uso: $0 [--uninstall]"
        exit 0
        ;;
    esac
done

echo ""
echo -e "${BLUE}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║${NC}      ${GREEN}Dxrk AI - Auto Update Installer${NC}         ${BLUE}║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════╝${NC}"
echo ""

# Detectar SO
OS=$(uname -s)

install_cron() {
    echo -e "${BLUE}[1/3]${NC} Configurando cron job..."

    CRON_CMD="@daily curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/update.sh | bash"

    # Agregar al crontab
    (
        crontab -l 2>/dev/null | grep -v "Dxrk AI"
        echo "# Dxrk AI Auto Update"
        echo "$CRON_CMD"
    ) | crontab -

    echo -e "${GREEN}✓${NC} Cron job instalado"
    echo -e "  ${DIM}Ejecutará la actualización diariamente a medianoche${NC}"
}

uninstall_cron() {
    echo -e "${YELLOW}[1/3]${NC} Removiendo cron job..."
    crontab -l 2>/dev/null | grep -v "Dxrk AI" | grep -v "@daily.*Dxrk" | crontab -
    echo -e "${GREEN}✓${NC} Cron job removido"
}

install_launchagent() {
    echo -e "${BLUE}[1/3]${NC} Configurando LaunchAgent..."

    PLIST_DIR="$HOME/Library/LaunchAgents"
    PLIST_FILE="$PLIST_DIR/com.dxrk-hex.autoupdate.plist"

    mkdir -p "$PLIST_DIR"

    # Crear plist
    cat >"$PLIST_FILE" <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.dxrk-hex.autoupdate</string>
    <key>ProgramArguments</key>
    <array>
        <string>/bin/bash</string>
        <string>-c</string>
        <string>curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/update.sh | bash</string>
    </array>
    <key>StartInterval</key>
    <integer>86400</integer>
    <key>RunAtLoad</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/dxrk-autoupdate.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/dxrk-autoupdate-error.log</string>
</dict>
</plist>
EOF

    # Cargar el servicio
    launchctl unload "$PLIST_FILE" 2>/dev/null || true
    launchctl load "$PLIST_FILE"

    echo -e "${GREEN}✓${NC} LaunchAgent instalado"
    echo -e "  ${DIM}Ejecutará la actualización cada 24 horas${NC}"
    echo -e "  ${DIM}Logs en: /tmp/dxrk-autoupdate*.log${NC}"
}

uninstall_launchagent() {
    echo -e "${YELLOW}[1/3]${NC} Removiendo LaunchAgent..."
    PLIST_FILE="$HOME/Library/LaunchAgents/com.dxrk-hex.autoupdate.plist"
    launchctl unload "$PLIST_FILE" 2>/dev/null || true
    rm -f "$PLIST_FILE"
    echo -e "${GREEN}✓${NC} LaunchAgent removido"
}

install_systemd() {
    echo -e "${BLUE}[1/3]${NC} Configurando systemd timer..."

    SYSTEMD_DIR="$HOME/.config/systemd/user"
    mkdir -p "$SYSTEMD_DIR"

    # Crear service
    cat >"$SYSTEMD_DIR/dxrk-autoupdate.service" <<'EOF'
[Unit]
Description=Dxrk AI Auto Update

[Service]
Type=oneshot
ExecStart=/bin/bash -c "curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/update.sh | bash"
EOF

    # Crear timer
    cat >"$SYSTEMD_DIR/dxrk-autoupdate.timer" <<'EOF'
[Unit]
Description=Dxrk AI Auto Update Timer

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
EOF

    # Habilitar
    systemctl --user daemon-reload
    systemctl --user enable dxrk-autoupdate.timer
    systemctl --user start dxrk-autoupdate.timer

    echo -e "${GREEN}✓${NC} Systemd timer instalado"
    echo -e "  ${DIM}Ejecutará la actualización diariamente${NC}"
}

uninstall_systemd() {
    echo -e "${YELLOW}[1/3]${NC} Removiendo systemd timer..."
    systemctl --user stop dxrk-autoupdate.timer 2>/dev/null || true
    systemctl --user disable dxrk-autoupdate.timer 2>/dev/null || true
    rm -f "$HOME/.config/systemd/user/dxrk-autoupdate."*
    echo -e "${GREEN}✓${NC} Systemd timer removido"
}

# Instalar según SO
if [[ "$UNINSTALL" == true ]]; then
    echo -e "${YELLOW}Desinstalando actualización automática...${NC}"
    echo ""

    case "$OS" in
    Darwin)
        uninstall_launchagent
        ;;
    Linux)
        if command -v systemctl &>/dev/null; then
            uninstall_systemd
        else
            uninstall_cron
        fi
        ;;
    *)
        uninstall_cron
        ;;
    esac

    echo ""
    echo -e "${GREEN}✓ Desinstalación completada${NC}"
    exit 0
fi

echo -e "${BLUE}[2/3]${NC} Verificando instalación de dxrk..."
if command -v dxrk &>/dev/null; then
    echo -e "${GREEN}✓${NC} Dxrk encontrado: $(dxrk --version)"
else
    echo -e "${YELLOW}!${NC} Dxrk no está instalado"
    echo -e "  ${DIM}¿Quieres instalarlo primero?${NC}"
    read -p "  Instalar dxrk ahora? [y/N] " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.sh | bash
    fi
fi

echo ""
echo -e "${BLUE}[3/3]${NC} Instalando actualización automática..."

case "$OS" in
Darwin)
    install_launchagent
    ;;
Linux)
    if command -v systemctl &>/dev/null; then
        install_systemd
    else
        install_cron
    fi
    ;;
*)
    install_cron
    ;;
esac

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║${NC}         ${GREEN}Instalación completada!${NC}                  ${GREEN}║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "Dxrk se actualizará automáticamente."
echo ""
echo -e "Comandos útiles:"
echo -e "  ${DIM}# Verificar estado del servicio${NC}"
if [[ "$OS" == "Darwin" ]]; then
    echo -e "  launchctl list | grep dxrk"
elif [[ "$OS" == "Linux" ]] && command -v systemctl &>/dev/null; then
    echo -e "  systemctl --user status dxrk-autoupdate.timer"
else
    echo -e "  crontab -l | grep dxrk"
fi
echo ""
echo -e "  ${DIM}# Desinstalar${NC}"
echo -e "  ./scripts/install-auto-update.sh --uninstall"
echo ""
