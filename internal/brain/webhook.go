package brain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookConfig holds webhook configuration.
type WebhookConfig struct {
	// URL of the webhook
	URL string
	// Type: discord, slack, teams, generic
	Type string
	// Timeout for requests
	Timeout time.Duration
}

// DefaultWebhookConfig returns a default webhook configuration.
func DefaultWebhookConfig() *WebhookConfig {
	return &WebhookConfig{
		Timeout: 10 * time.Second,
	}
}

// WebhookMessage represents a message to send via webhook.
type WebhookMessage struct {
	// Content is the main message text
	Content string `json:"content,omitempty"`
	// Embeds for rich messages
	Embeds []WebhookEmbed `json:"embeds,omitempty"`
	// Username for the webhook bot
	Username string `json:"username,omitempty"`
	// AvatarURL for the webhook bot
	AvatarURL string `json:"avatar_url,omitempty"`
}

// WebhookEmbed represents an embedded message.
type WebhookEmbed struct {
	// Title of the embed
	Title string `json:"title,omitempty"`
	// Description of the embed
	Description string `json:"description,omitempty"`
	// Color of the embed (hex value)
	Color int `json:"color,omitempty"`
	// Fields to display
	Fields []WebhookField `json:"fields,omitempty"`
	// Footer text
	Footer WebhookFooter `json:"footer,omitempty"`
	// Timestamp
	Timestamp string `json:"timestamp,omitempty"`
}

// WebhookField represents a field in an embed.
type WebhookField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// WebhookFooter represents a footer in an embed.
type WebhookFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// WebhookClient is a client for sending webhook notifications.
type WebhookClient struct {
	config *WebhookConfig
	client *http.Client
}

// NewWebhookClient creates a new Webhook client.
func NewWebhookClient(cfg *WebhookConfig) *WebhookClient {
	if cfg == nil {
		cfg = DefaultWebhookConfig()
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	return &WebhookClient{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// Send sends a webhook message.
func (w *WebhookClient) Send(msg WebhookMessage) error {
	if w.config.URL == "" {
		return fmt.Errorf("webhook URL is required")
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", w.config.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// SendSimple sends a simple text message.
func (w *WebhookClient) SendSimple(content string) error {
	return w.Send(WebhookMessage{
		Content: content,
	})
}

// SendWithEmbed sends a rich embed message.
func (w *WebhookClient) SendWithEmbed(embed WebhookEmbed) error {
	return w.Send(WebhookMessage{
		Embeds: []WebhookEmbed{embed},
	})
}

// NotifyInstall sends a notification about an installation.
func (w *WebhookClient) NotifyInstall(agent, status string) error {
	embed := WebhookEmbed{
		Title:       "🤖 Agent Installed",
		Description: fmt.Sprintf("**%s** has been %s", agent, status),
		Color:       0x5dfc8e, // Green
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: WebhookFooter{
			Text: "Dxrk Hex",
		},
	}

	return w.SendWithEmbed(embed)
}

// NotifyError sends an error notification.
func (w *WebhookClient) NotifyError(context, message string) error {
	embed := WebhookEmbed{
		Title:       "❌ Error",
		Description: message,
		Color:       0xaa0033, // Red
		Fields: []WebhookField{
			{Name: "Context", Value: context, Inline: true},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: WebhookFooter{
			Text: "Dxrk Hex",
		},
	}

	return w.SendWithEmbed(embed)
}

// NotifyUpdate sends an update notification.
func (w *WebhookClient) NotifyUpdate(version, changelog string) error {
	embed := WebhookEmbed{
		Title:       "🔄 Update Available",
		Description: fmt.Sprintf("Version **%s** is available", version),
		Color:       0x3399ff, // Blue
		Fields: []WebhookField{
			{Name: "Changelog", Value: changelog, Inline: false},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: WebhookFooter{
			Text: "Dxrk Hex",
		},
	}

	return w.SendWithEmbed(embed)
}

// NotifySync sends a sync notification.
func (w *WebhookClient) NotifySync(filesChanged int, status string) error {
	embed := WebhookEmbed{
		Title:       "🔄 Config Synced",
		Description: fmt.Sprintf("**%d** files have been %s", filesChanged, status),
		Color:       0x5dfc8e, // Green
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: WebhookFooter{
			Text: "Dxrk Hex",
		},
	}

	return w.SendWithEmbed(embed)
}

// IsConfigured returns true if the webhook is configured.
func (w *WebhookClient) IsConfigured() bool {
	return w.config.URL != ""
}

// String returns a string representation of the Webhook client.
func (w *WebhookClient) String() string {
	return fmt.Sprintf("WebhookClient{URL: %s, Type: %s}", w.config.URL, w.config.Type)
}
