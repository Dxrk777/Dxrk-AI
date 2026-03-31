package connector

import (
	"testing"
)

func TestNewConnector(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		Port:    8081,
	}

	c := New(cfg)
	if c == nil {
		t.Fatal("New() returned nil")
	}
	if c.config.Port != 8081 {
		t.Errorf("Expected port 8081, got %d", c.config.Port)
	}
}

func TestConnectorStatus(t *testing.T) {
	cfg := &Config{
		Enabled: true,
		Port:    8081,
		Telegram: TelegramConfig{
			Enabled: true,
			Token:   "test-token",
		},
		Discord: DiscordConfig{
			Enabled:    true,
			WebhookURL: "https://discord.com/webhook",
		},
		WhatsApp: WhatsAppConfig{
			Enabled: true,
		},
	}

	c := New(cfg)
	status := c.Status()

	if status["enabled"] != true {
		t.Error("Expected enabled=true")
	}
	if status["port"].(int) != 8081 {
		t.Errorf("Expected port 8081, got %v", status["port"])
	}
	if status["telegram"] != true {
		t.Error("Expected telegram=true")
	}
	if status["discord"] != true {
		t.Error("Expected discord=true")
	}
	if status["whatsapp"] != true {
		t.Error("Expected whatsapp=true")
	}
}

func TestProcessCommand_Help(t *testing.T) {
	c := &Connector{}
	response := c.processCommand("help")

	if response == "" {
		t.Error("Expected help text, got empty")
	}
	if len(response) < 50 {
		t.Error("Help text too short")
	}
}

func TestProcessCommand_Status(t *testing.T) {
	c := &Connector{}
	response := c.processCommand("status")

	if response == "" {
		t.Error("Expected status response, got empty")
	}
}

func TestProcessCommand_List(t *testing.T) {
	c := &Connector{}
	response := c.processCommand("list")

	// Response may be empty or error depending on dxrk installation
	// Just verify it doesn't panic
	_ = response
}

func TestProcessCommand_Unknown(t *testing.T) {
	c := &Connector{}
	response := c.processCommand("unknowncommand")

	if response == "" {
		t.Error("Expected error message for unknown command")
	}
}

func TestProcessCommand_Empty(t *testing.T) {
	c := &Connector{}
	response := c.processCommand("")

	if response == "" {
		t.Error("Expected message for empty command")
	}
}

func TestProcessCommand_WithArgs(t *testing.T) {
	c := &Connector{}
	response := c.processCommand("install claude opencode")

	if response == "" {
		t.Error("Expected response for install command")
	}
}

func TestBuildHelpText(t *testing.T) {
	c := &Connector{}
	help := c.buildHelpText()

	expected := []string{
		"help",
		"status",
		"list",
		"install",
		"backup",
	}

	for _, cmd := range expected {
		if !contains(help, cmd) {
			t.Errorf("Help text should contain '%s'", cmd)
		}
	}
}

func TestExecuteStatus(t *testing.T) {
	c := &Connector{}
	status := c.executeStatus()

	if status == "" {
		t.Error("Expected status text, got empty")
	}
}

func TestConfigDefaults(t *testing.T) {
	cfg := &Config{}
	c := New(cfg)

	// Should default to port 8081
	if c.config.Port != 8081 {
		t.Errorf("Expected default port 8081, got %d", c.config.Port)
	}
}

func TestConnectorNotRunning(t *testing.T) {
	c := New(&Config{})

	// Should not be running initially
	status := c.Status()
	if status["running"] != false {
		t.Error("Expected running=false initially")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
