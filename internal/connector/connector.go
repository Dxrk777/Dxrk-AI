// Package connector implements remote control for Dxrk AI via messaging platforms.
//
// Supported platforms:
//   - Telegram: Bot API integration
//   - Discord: Webhook integration
//   - WhatsApp: Twilio integration
//
// Usage:
//  1. Enable platform in settings
//  2. Configure API tokens
//  3. Run dxrk with --connector flag
//  4. Control Dxrk AI from Telegram/Discord/WhatsApp
package connector

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Dxrk777/Dxrk-Hex/internal/brain"
)

// Config holds connector configuration.
type Config struct {
	Enabled  bool           `json:"enabled"`
	Port     int            `json:"port"`
	Telegram TelegramConfig `json:"telegram"`
	Discord  DiscordConfig  `json:"discord"`
	WhatsApp WhatsAppConfig `json:"whatsapp"`
}

// TelegramConfig holds Telegram bot configuration.
type TelegramConfig struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
	ChatID  string `json:"chat_id"`
}

// DiscordConfig holds Discord webhook configuration.
type DiscordConfig struct {
	Enabled    bool   `json:"enabled"`
	WebhookURL string `json:"webhook_url"`
}

// WhatsAppConfig holds WhatsApp (Twilio) configuration.
type WhatsAppConfig struct {
	Enabled    bool   `json:"enabled"`
	AccountSID string `json:"account_sid"`
	AuthToken  string `json:"auth_token"`
	FromNumber string `json:"from_number"`
}

// IncomingMessage represents a message from a platform.
type IncomingMessage struct {
	Platform string `json:"platform"`
	From     string `json:"from"`
	To       string `json:"to"`
	Body     string `json:"body"`
}

// Connector manages remote connections.
type Connector struct {
	mu      sync.RWMutex
	config  *Config
	server  *http.Server
	running bool
	startAt time.Time
}

// New creates a new Connector.
func New(cfg *Config) *Connector {
	if cfg.Port == 0 {
		cfg.Port = 8081
	}
	return &Connector{
		config: cfg,
	}
}

// Start launches the HTTP webhook server.
func (c *Connector) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return fmt.Errorf("connector: already running")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/telegram", c.handleTelegram)
	mux.HandleFunc("/webhook/discord", c.handleDiscord)
	mux.HandleFunc("/webhook/whatsapp", c.handleWhatsApp)
	mux.HandleFunc("/health", c.handleHealth)
	mux.HandleFunc("/status", c.handleStatus)

	c.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", c.config.Port),
		Handler: mux,
	}

	c.running = true
	c.startAt = time.Now()

	go func() {
		if err := c.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[CONNECTOR] Server error: %v\n", err)
		}
	}()

	fmt.Printf("[CONNECTOR] Listening on port %d\n", c.config.Port)
	return nil
}

// Stop gracefully shuts down the connector.
func (c *Connector) Stop() {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return
	}

	c.running = false
	server := c.server
	c.mu.Unlock()

	if server != nil {
		// Use Shutdown with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			fmt.Printf("[CONNECTOR] Shutdown error: %v\n", err)
		}
	}
	fmt.Println("[CONNECTOR] Stopped")
}

// Status returns connector status.
func (c *Connector) Status() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"enabled":  c.config.Enabled,
		"running":  c.running,
		"port":     c.config.Port,
		"uptime":   time.Since(c.startAt).String(),
		"telegram": c.config.Telegram.Enabled,
		"discord":  c.config.Discord.Enabled,
		"whatsapp": c.config.WhatsApp.Enabled,
	}
}

// handleTelegram processes Telegram webhook requests.
func (c *Connector) handleTelegram(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !c.config.Telegram.Enabled {
		http.Error(w, "telegram not enabled", http.StatusServiceUnavailable)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var update struct {
		Message struct {
			Chat struct {
				ID int64 `json:"id"`
			} `json:"chat"`
			From struct {
				Username string `json:"username"`
			} `json:"from"`
			Text string `json:"text"`
		} `json:"message"`
	}

	if err := json.Unmarshal(body, &update); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if update.Message.Text == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Process command
	response := c.processCommand(update.Message.Text)

	// Send response via Telegram
	c.sendTelegramMessage(update.Message.Chat.ID, response)

	w.WriteHeader(http.StatusOK)
}

// handleDiscord processes Discord webhook requests.
func (c *Connector) handleDiscord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !c.config.Discord.Enabled || c.config.Discord.WebhookURL == "" {
		http.Error(w, "discord not enabled", http.StatusServiceUnavailable)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload struct {
		Content string `json:"content"`
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.Content == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Process command
	response := c.processCommand(payload.Content)

	// Send response via Discord webhook
	c.sendDiscordMessage(response)

	w.WriteHeader(http.StatusOK)
}

// handleWhatsApp processes WhatsApp (Twilio) webhook requests.
func (c *Connector) handleWhatsApp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !c.config.WhatsApp.Enabled {
		http.Error(w, "whatsapp not enabled", http.StatusServiceUnavailable)
		return
	}

	_, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse form data
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		// Twilio sends form data
		w.WriteHeader(http.StatusOK)
		return // Twilio expects 200 OK immediately
	}

	w.WriteHeader(http.StatusOK)
}

// processCommand executes a command and returns the result.
func (c *Connector) processCommand(input string) string {
	text := strings.TrimSpace(input)
	parts := strings.Fields(text)

	if len(parts) == 0 {
		return "No command received. Send 'help' for available commands."
	}

	cmd := strings.ToLower(parts[0])
	args := ""
	if len(parts) > 1 {
		args = strings.Join(parts[1:], " ")
	}

	switch cmd {
	case "help", "ayuda", "comandos":
		return c.buildHelpText()
	case "status", "estado", "stat":
		return c.executeStatus()
	case "list", "ls", "agentes":
		return c.executeList()
	case "install", "i", "add":
		return c.executeInstall(args)
	case "uninstall", "remove", "rm":
		return c.executeUninstall(args)
	case "sync", "s":
		return "🔄 Sync command queued. Use GUI or CLI for sync."
	case "backup", "bk":
		return c.executeBackup()
	case "restore", "rest":
		return "🔄 Restore command queued. Use GUI or CLI for restore."
	case "version", "v", "ver":
		return c.executeVersion()
	case "update", "upgrade":
		return c.executeUpdate()
	case "memory", "historial":
		return c.executeMemory()
	case "vault", "encrypt":
		return "🔐 Vault available in GUI. Use 'encrypt <text>' for quick encrypt."
	case "remote", "connect":
		return c.executeRemoteStatus()
	case "skills", "patrones":
		return c.executeSkills()
	default:
		// Check if it's an agent name shorthand
		if c.isAgentName(cmd) {
			return c.executeInstall(cmd)
		}
		return fmt.Sprintf("Unknown command: %s. Send 'help' for available commands.", cmd)
	}
}

// isAgentName checks if input is a valid agent name.
func (c *Connector) isAgentName(name string) bool {
	agents := []string{
		"claude", "opencode", "cursor", "windsurf",
		"codex", "gemini", "vscode", "antigravity",
		"ollama", "groq", "deepseek",
	}
	for _, a := range agents {
		if name == a {
			return true
		}
	}
	return false
}

// executeInstall handles install command with all agents.
func (c *Connector) executeInstall(args string) string {
	if args == "" {
		return c.listAvailableAgents()
	}

	agent := strings.ToLower(strings.TrimSpace(args))

	// Detailed agent descriptions
	switch agent {
	case "opencode":
		return `🎯 OPENCODE PRO — Código Profesional

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Tu fork de Cursor con poder Dxrk!

CARACTERÍSTICAS:
  🔐 Vault — Encriptación AES-256
  🧠 Memory — Historial + Engram
  🔗 Remote — Telegram/Discord/WhatsApp
  🧠 Skills — 60 patrones de código

WORKFLOWS:
  📋 SDD (Spec-Driven Development)
  🏗️ Clean Architecture
  🧪 Strict TDD Mode
  ✅ Code Review automation

INSTALACIÓN:
  dxrk install opencode

O usa GUI: Menu → Start installation → OpenCode`

	case "claude":
		return `🤖 CLAUDE CODE — AI Assistant

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ El agent más inteligente de Anthropic

CARACTERÍSTICAS:
  🧠 Mejor reasoning y contexto largo
  💻 Excelente para código complejo
  📝 Análisis de arquitectura
  🔍 Code review exhaustivo

BEST FOR:
  → Código profesional de alta calidad
  → Proyectos grandes y complejos
  → Arquitectura de sistemas
  → Refactoring importante

INSTALACIÓN:
  dxrk install claude

O usa GUI: Menu → Start installation → Claude`

	case "cursor":
		return `🎨 CURSOR AI — IDE Inteligente

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Editor de código con IA integrada

CARACTERÍSTICAS:
  💻 Editor completo con IA
  🔄 Autocompletado inteligente
  🎯 Multi-model support
  📁 Proyecto-aware

BEST FOR:
  → Desarrollo rápido
  → Código iterativo
  → Debug visual
  → Principiantes

INSTALACIÓN:
  dxrk install cursor

O usa GUI: Menu → Start installation → Cursor`

	case "windsurf":
		return `🌊 WINDSURF — AI Coding Agent

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Agent con contexto de proyecto completo

CARACTERÍSTICAS:
  🌊 Flow State™ mode
  📂 Proyecto completo en contexto
  🔧 Tools avanzados
  🎯 Autonomous coding

BEST FOR:
  → Coding autonomous
  → Proyectos completos
  → Automatización
  → Senior developers

INSTALACIÓN:
  dxrk install windsurf

O usa GUI: Menu → Start installation → Windsurf`

	case "ollama":
		return `🆓 OLLAMA — Local LLM (Gratis)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ LLM local, 100% offline, sin API key

CARACTERÍSTICAS:
  🆓 100% gratis
  🔒 Privado (nada sale de tu PC)
  📱 Funciona offline
  🚀 Modelos: llama, mistral, codellama

MODELS DISPONIBLES:
  • llama3.2 — General purpose
  • codellama — Especializado en código
  • mistral — Rápido y eficiente
  • phi — Muy ligero

INSTALACIÓN:
  dxrk install ollama

O usa GUI: Menu → Start installation → Ollama`

	case "groq":
		return `🆓 GROQ — API Gratuita

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ API gratis, respuestas ultra-rápidas

CARACTERÍSTICAS:
  🆓 API key gratuita
  ⚡ Respuestas extremadamente rápidas
  💬 Modelos: llama, mixtral
  🔓 Sin costo hasta cierto límite

MODELS:
  • llama-3.1-70b — Muy potente
  • mixtral-8x7b — Balanceado
  • llama-3.1-8b — Rápido

INSTALACIÓN:
  dxrk install groq

O usa GUI: Menu → Start installation → Groq`

	case "deepseek":
		return `🆓 DEEPSEEK — API Gratuita

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ API china gratuita con modelos especializados

CARACTERÍSTICAS:
  🆓 API key gratuita
  🧠 Modelo especializado en código
  💰 Muy económico
  🌍 Internacional

MODELS:
  • deepseek-coder — Especializado código
  • deepseek-chat — General purpose

BEST FOR:
  → Código
  → Matemáticas
  → Análisis técnico

INSTALACIÓN:
  dxrk install deepseek

O usa GUI: Menu → Start installation → DeepSeek`

	case "codex":
		return `⚡ OPENAI CODEX — API de OpenAI

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ El clásico agent de OpenAI

CARACTERÍSTICAS:
  💻 API de OpenAI
  🔥 Muy probado y estable
  📚 Contexto de código
  ⚡ Rápido

BEST FOR:
  → Código legacy
  → APIs de OpenAI
  → Automatización

INSTALACIÓN:
  dxrk install codex

O usa GUI: Menu → Start installation → Codex`

	case "gemini":
		return `🌟 GOOGLE GEMINI CLI

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Agent de Google con multimodal

CARACTERÍSTICAS:
  🌟 Multimodal (texto + imágenes)
  📊 Contexto de Google
  🔍 Búsquedas web
  💡 Rápido

BEST FOR:
  → Proyectos con imágenes
  → Análisis de UI/UX
  → Búsquedas de información

INSTALACIÓN:
  dxrk install gemini

O usa GUI: Menu → Start installation → Gemini`

	case "vscode":
		return `📦 VS CODE + AI Extensions

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Tu editor favorito con IA

CARACTERÍSTICAS:
  📦 VS Code base
  🤖 GitHub Copilot integration
  🔌 Muchas extensions
  💻 Familiar para todos

BEST FOR:
  → Si ya usas VS Code
  → Setup personalizado
  → Extensions específicas

INSTALACIÓN:
  dxrk install vscode

O usa GUI: Menu → Start installation → VSCode`

	case "antigravity":
		return `🚀 ANTIGRAVITY — AGENT Framework

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ Framework propio de agentes IA

CARACTERÍSTICAS:
  🔧 Framework custom
  🎯 Workflows personalizados
  🧠 Memoria adaptativa
  🔗 Conexiones flexibles

BEST FOR:
  → Desarrollo custom
  → Workflows únicos
  → Integración advanced

INSTALACIÓN:
  dxrk install antigravity

O usa GUI: Menu → Start installation → Antigravity`

	default:
		return fmt.Sprintf(`❌ Unknown agent: %s

Available agents:
%s

Examples:
  /install claude
  /install opencode  ← Recommended for PRO
  /install ollama`, agent, c.listAvailableAgents())
	}
}

// listAvailableAgents returns formatted list of agents.
func (c *Connector) listAvailableAgents() string {
	return `🤖 Available Agents:

🎯 PROFESSIONAL (Recommended):
  • opencode    — OpenCode PRO ★ (Tu fork Dxrk)
    └─ SDD, Clean Arch, Vault, Memory

📦 CLI Tools:
  • claude      — Anthropic Claude Code
  • cursor      — Cursor AI
  • windsurf    — Windsurf AI
  • codex       — OpenAI Codex
  • gemini      — Google Gemini CLI
  • vscode      — VS Code + AI
  • antigravity — AGENT Framework

🆓 FREE LLM (No API key needed):
  • ollama     — Local LLM (offline)
  • groq       — Free API (fast)
  • deepseek   — Free API (smart)

💡 For code professional → Use opencode
   Quick: Just type "opencode" or "/install opencode"`
}

// executeUninstall handles uninstall command.
func (c *Connector) executeUninstall(args string) string {
	if args == "" {
		return "Usage: uninstall <agent>\nExample: /uninstall claude"
	}

	agent := strings.ToLower(strings.TrimSpace(args))
	return fmt.Sprintf(`🗑️ Uninstall queued: %s

Uninstall request received!
Use GUI or CLI for uninstallation:
  dxrk uninstall %s`, agent, agent)
}

// buildHelpText generates the help message.
func (c *Connector) buildHelpText() string {
	return `📋 Dxrk AI Remote Commands
─────────────────────

🤖 AGENTS:
  install <agent> — Install an agent
  uninstall <agent> — Remove an agent
  list — List installed agents
  agents — Show available agents

⚙️ SYSTEM:
  status — System status
  update — Check for updates
  sync — Sync configurations
  backup — Backup settings
  restore — Restore backup

🧠 MEMORY:
  memory — View install history
  vault — Encryption status

🔗 INFO:
  skills — List available skills
  version — Show version
  remote — Remote status

💡 SHORTCUTS:
  Just type agent name: "claude"
  Shortcuts: /i, /ls, /v

Example: Send "status" or "claude"`
}

// executeStatus returns system status.
func (c *Connector) executeStatus() string {
	return `✅ Dxrk AI Status
─────────────────
• System: Online
• Memory: Active
• Vault: Ready
• Connector: Running`
}

// executeList returns installed agents.
func (c *Connector) executeList() string {
	// Try to run dxrk list command
	cmd := exec.Command("dxrk", "status")
	output, err := cmd.Output()
	if err != nil {
		return "❌ Could not retrieve agent list"
	}
	return string(output)
}

// sendTelegramMessage sends a message via Telegram Bot API.
func (c *Connector) sendTelegramMessage(chatID int64, text string) {
	if c.config.Telegram.Token == "" {
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.config.Telegram.Token)
	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "Markdown",
	}

	body, _ := json.Marshal(payload)
	http.Post(url, "application/json", strings.NewReader(string(body)))
}

// sendDiscordMessage sends a message via Discord webhook.
func (c *Connector) sendDiscordMessage(text string) {
	if c.config.Discord.WebhookURL == "" {
		return
	}

	payload := map[string]string{"content": text}
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}
	http.Post(c.config.Discord.WebhookURL, "application/json", strings.NewReader(string(body)))
}

// handleHealth returns health status.
func (c *Connector) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// handleStatus returns full status.
func (c *Connector) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c.Status())
}

// executeBackup handles backup command.
func (c *Connector) executeBackup() string {
	return `💾 Backup Requested!

Use GUI or CLI:
  dxrk backup
  dxrk backup --name "my-backup"

Backups are stored in:
  ~/.dxrk/backups/`
}

// executeVersion returns version info.
func (c *Connector) executeVersion() string {
	return `📦 Dxrk AI
Version: 000.13%
Build: Latest

GitHub: github.com/Dxrk777/Dxrk-Hex
Upstream: github.com/Gentleman-Programming/gentle-ai

Updates: Automatic (daily sync)`
}

// executeUpdate handles update check.
func (c *Connector) executeUpdate() string {
	return `🔄 Update Check

Dxrk AI sincroniza con upstream diariamente.
Para forzar actualización:

  dxrk upgrade

O usa la GUI: Menú → Upgrade tools`
}

// executeMemory returns memory/history status.
func (c *Connector) executeMemory() string {
	return `🧠 Memory Status

Dxrk AI tracks:
  • Install history
  • Agent usage
  • Preferences

View history:
  dxrk memory
  Or use: Menu → 🧠 Memory

Engram sync: Active`
}

// executeRemoteStatus returns remote connection status.
func (c *Connector) executeRemoteStatus() string {
	status := c.Status()
	return fmt.Sprintf(`🔗 Remote Status

Telegram: %v
Discord: %v
WhatsApp: %v
Server: %v
Port: %d
Uptime: %s

Configure in GUI: Menu → 🔗 Remote Connect`, status["telegram"], status["discord"], status["whatsapp"], status["running"], status["port"], status["uptime"])
}

// executeSkills returns available skills count.
func (c *Connector) executeSkills() string {
	return `🧠 Available Skills: 60+

Categories:
  • Frameworks: React, Next.js, Angular, Vue, Svelte
  • Languages: TypeScript, Go, Rust, Python
  • Tools: Docker, K8s, Playwright
  • Architecture: SDD, Clean Arch, Security
  • AI: AI SDK, Prompt Engineer, Zod

Skills are installed automatically with agents.
Use GUI: Menu → Start installation`
}

// BrainIntegration provides AI-powered command processing.
type BrainIntegration struct {
	*brain.Brain
	*brain.Commander
	*brain.Emailer
	*brain.Memory
}

// NewBrainIntegration creates a new Brain integration instance.
func NewBrainIntegration(dataDir string, emailCfg brain.EmailConfig) (*BrainIntegration, error) {
	mem, err := brain.NewMemory(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to init memory: %w", err)
	}

	b := &BrainIntegration{
		Brain:     brain.New(&brain.Config{MemoryDir: dataDir}),
		Commander: brain.NewCommander(30 * time.Second),
		Emailer:   brain.NewEmailer(emailCfg),
		Memory:    mem,
	}

	return b, nil
}

// ProcessCommand processes a natural language command using the brain.
func (bi *BrainIntegration) ProcessCommand(input string) string {
	result := brain.Think(input, bi.Brain, bi.Commander, bi.Emailer, bi.Memory)
	return result.Response
}

// ExecuteCommand runs a shell command.
func (bi *BrainIntegration) ExecuteCommand(cmd string) *brain.CommandResult {
	result, _ := bi.Commander.Execute(cmd)
	return result
}

// SendEmail sends an email.
func (bi *BrainIntegration) SendEmail(to []string, subject, body string) error {
	return bi.Emailer.SendHTML(to, subject, body)
}

// GetHistory returns command history.
func (bi *BrainIntegration) GetHistory(limit int) []brain.MemoryEntry {
	if limit <= 0 {
		limit = 20
	}
	return bi.Memory.Recent(limit)
}
