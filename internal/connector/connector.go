// Package connector implements remote control for Dxrk Hex via messaging platforms.
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
//  4. Control Dxrk Hex from Telegram/Discord/WhatsApp
package connector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
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
	defer c.mu.Unlock()

	if !c.running {
		return
	}

	c.running = false
	if c.server != nil {
		c.server.Close()
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
	case "help":
		return c.buildHelpText()
	case "status":
		return c.executeStatus()
	case "list":
		return c.executeList()
	case "install":
		return fmt.Sprintf("Install command: %s (use GUI or CLI for full install)", args)
	case "backup":
		return "Backup command: Use GUI or CLI for backup operations"
	case "version":
		return "Dxrk Hex v000.13%"
	default:
		return fmt.Sprintf("Unknown command: %s. Send 'help' for available commands.", cmd)
	}
}

// buildHelpText generates the help message.
func (c *Connector) buildHelpText() string {
	return `📋 Dxrk Hex Remote Commands
─────────────────────
• help — Show this help message
• status — System status
• list — List installed agents
• install <agent> — Request installation
• backup — Backup configuration
• version — Show version

Example: Send "status" to check system`
}

// executeStatus returns system status.
func (c *Connector) executeStatus() string {
	return `✅ Dxrk Hex Status
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
