package brain

import (
	"testing"
)

func TestDefaultWebhookConfig(t *testing.T) {
	cfg := DefaultWebhookConfig()

	if cfg.Timeout != 10*1000000000 {
		t.Errorf("Expected timeout 10s, got %v", cfg.Timeout)
	}
}

func TestNewWebhookClient(t *testing.T) {
	cfg := &WebhookConfig{
		URL:     "https://discord.com/api/webhooks/test",
		Type:    "discord",
		Timeout: 5 * 1000000000,
	}

	client := NewWebhookClient(cfg)

	if client.config.URL != "https://discord.com/api/webhooks/test" {
		t.Errorf("Expected URL, got %s", client.config.URL)
	}
}

func TestNewWebhookClientNilConfig(t *testing.T) {
	client := NewWebhookClient(nil)

	if client.config.Timeout != 10*1000000000 {
		t.Errorf("Expected default timeout")
	}
}

func TestWebhookClientIsConfigured(t *testing.T) {
	tests := []struct {
		name     string
		config   WebhookConfig
		expected bool
	}{
		{
			name:     "configured",
			config:   WebhookConfig{URL: "https://test.com/webhook"},
			expected: true,
		},
		{
			name:     "empty URL",
			config:   WebhookConfig{URL: ""},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewWebhookClient(&tt.config)
			if client.IsConfigured() != tt.expected {
				t.Errorf("IsConfigured() = %v, want %v", client.IsConfigured(), tt.expected)
			}
		})
	}
}

func TestWebhookClientString(t *testing.T) {
	client := NewWebhookClient(&WebhookConfig{
		URL:  "https://test.com/webhook",
		Type: "discord",
	})

	str := client.String()
	if str == "" {
		t.Error("String() should return a non-empty value")
	}
}

func TestWebhookMessage(t *testing.T) {
	msg := WebhookMessage{
		Content:  "Hello, World!",
		Username: "Dxrk Hex",
	}

	if msg.Content != "Hello, World!" {
		t.Errorf("Expected content, got %s", msg.Content)
	}

	if msg.Username != "Dxrk Hex" {
		t.Errorf("Expected username, got %s", msg.Username)
	}
}

func TestWebhookEmbed(t *testing.T) {
	embed := WebhookEmbed{
		Title:       "Test Title",
		Description: "Test Description",
		Color:       0x5dfc8e,
		Fields: []WebhookField{
			{Name: "Field 1", Value: "Value 1", Inline: true},
			{Name: "Field 2", Value: "Value 2", Inline: false},
		},
		Footer: WebhookFooter{
			Text:    "Footer text",
			IconURL: "https://example.com/icon.png",
		},
	}

	if embed.Title != "Test Title" {
		t.Errorf("Expected title, got %s", embed.Title)
	}

	if len(embed.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(embed.Fields))
	}
}

func TestWebhookSendNoURL(t *testing.T) {
	client := NewWebhookClient(&WebhookConfig{
		URL: "", // No URL
	})

	err := client.SendSimple("test")
	if err == nil {
		t.Error("Expected error for empty URL")
	}
}

func TestWebhookSendSimple(t *testing.T) {
	client := NewWebhookClient(&WebhookConfig{
		URL: "http://localhost:99999", // Invalid URL
	})

	err := client.SendSimple("test")
	// Should fail because URL is not reachable
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestWebhookClientIsConfiguredTrue(t *testing.T) {
	client := NewWebhookClient(&WebhookConfig{
		URL: "https://discord.com/api/webhooks/123/abc",
	})

	if !client.IsConfigured() {
		t.Error("Expected IsConfigured() to return true")
	}
}
