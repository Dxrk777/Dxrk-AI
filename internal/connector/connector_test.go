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
		"Remote Commands",
		"status",
		"list",
		"install",
		"agents",
		"backup",
		"sync",
		"update",
		"version",
		"memory",
		"vault",
		"skills",
	}

	for _, cmd := range expected {
		if !contains(help, cmd) {
			t.Errorf("Help text should contain '%s'", cmd)
		}
	}
}

func TestListAvailableAgents(t *testing.T) {
	c := &Connector{}
	agents := c.listAvailableAgents()

	expected := []string{
		"claude",
		"opencode",
		"cursor",
		"windsurf",
		"ollama",
		"groq",
		"deepseek",
	}

	for _, agent := range expected {
		if !contains(agents, agent) {
			t.Errorf("Agent list should contain '%s'", agent)
		}
	}
}

func TestIsAgentName(t *testing.T) {
	c := &Connector{}

	// Valid agents
	valid := []string{"claude", "opencode", "cursor", "windsurf", "ollama", "groq"}
	for _, agent := range valid {
		if !c.isAgentName(agent) {
			t.Errorf("Expected '%s' to be valid agent name", agent)
		}
	}

	// Invalid agents
	invalid := []string{"invalid", "unknown", "random", ""}
	for _, agent := range invalid {
		if c.isAgentName(agent) {
			t.Errorf("Expected '%s' to be invalid agent name", agent)
		}
	}
}

func TestExecuteInstall(t *testing.T) {
	c := &Connector{}

	// Test with specific agent
	resp := c.executeInstall("claude")
	if !contains(resp, "Claude") {
		t.Error("Should mention Claude")
	}

	// Test with unknown agent
	resp = c.executeInstall("unknownagent")
	if !contains(resp, "Unknown agent") {
		t.Error("Should mention unknown agent")
	}

	// Test with no args
	resp = c.executeInstall("")
	if !contains(resp, "Available Agents") {
		t.Error("Should list available agents")
	}
}

func TestProcessCommand_AgentShortcut(t *testing.T) {
	c := &Connector{}

	// Test direct agent name
	resp := c.processCommand("claude")
	if !contains(resp, "Claude") {
		t.Errorf("Expected Claude install message, got: %s", resp)
	}

	resp = c.processCommand("ollama")
	if !contains(resp, "Ollama") {
		t.Errorf("Expected Ollama install message, got: %s", resp)
	}

	resp = c.processCommand("opencode")
	if !contains(resp, "OpenCode") {
		t.Errorf("Expected OpenCode install message, got: %s", resp)
	}
}

func TestProcessCommand_SpanishCommands(t *testing.T) {
	c := &Connector{}

	// Spanish commands
	resp := c.processCommand("ayuda")
	if !contains(resp, "Remote Commands") {
		t.Errorf("Expected help text, got: %s", resp)
	}

	resp = c.processCommand("estado")
	if !contains(resp, "Online") {
		t.Errorf("Expected status, got: %s", resp)
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

func TestProcessCommand_Sync(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("sync")
	if resp == "" {
		t.Error("Expected sync response")
	}
}

func TestProcessCommand_Update(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("update")
	if resp == "" {
		t.Error("Expected update response")
	}
}

func TestProcessCommand_Memory(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("memory")
	if resp == "" {
		t.Error("Expected memory response")
	}
}

func TestProcessCommand_Vault(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("vault")
	if resp == "" {
		t.Error("Expected vault response")
	}
}

func TestProcessCommand_Remote(t *testing.T) {
	c := New(&Config{})
	resp := c.processCommand("remote")
	if resp == "" {
		t.Error("Expected remote response")
	}
}

func TestProcessCommand_Skills(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("skills")
	if resp == "" {
		t.Error("Expected skills response")
	}
}

func TestProcessCommand_Version(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("version")
	if resp == "" {
		t.Error("Expected version response")
	}
}

func TestProcessCommand_VersionShort(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("v")
	if resp == "" {
		t.Error("Expected version short response")
	}
}

func TestProcessCommand_VersionAlt(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("ver")
	if resp == "" {
		t.Error("Expected version alt response")
	}
}

func TestProcessCommand_Backup(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("backup")
	if resp == "" {
		t.Error("Expected backup response")
	}
}

func TestProcessCommand_Restore(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("restore")
	if resp == "" {
		t.Error("Expected restore response")
	}
}

func TestProcessCommand_Agents(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("agentes")
	if resp == "" {
		t.Error("Expected agentes response")
	}
}

func TestProcessCommand_Comandos(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("comandos")
	if resp == "" {
		t.Error("Expected comandos response")
	}
}

func TestProcessCommand_StatusAlt(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("stat")
	if resp == "" {
		t.Error("Expected stat response")
	}
}

func TestExecuteBackup(t *testing.T) {
	c := &Connector{}
	resp := c.executeBackup()
	if !contains(resp, "Backup") {
		t.Error("Expected backup message")
	}
}

func TestExecuteVersion(t *testing.T) {
	c := &Connector{}
	resp := c.executeVersion()
	if !contains(resp, "Dxrk") {
		t.Error("Expected version message")
	}
	if !contains(resp, "000") {
		t.Error("Expected version number")
	}
}

func TestExecuteUpdate(t *testing.T) {
	c := &Connector{}
	resp := c.executeUpdate()
	if !contains(resp, "Update") {
		t.Error("Expected update message")
	}
}

func TestExecuteMemory(t *testing.T) {
	c := &Connector{}
	resp := c.executeMemory()
	if !contains(resp, "Memory") {
		t.Error("Expected memory message")
	}
}

func TestExecuteRemoteStatus(t *testing.T) {
	c := New(&Config{Port: 8081})
	resp := c.executeRemoteStatus()
	if !contains(resp, "Remote") {
		t.Error("Expected remote status message")
	}
}

func TestExecuteSkills(t *testing.T) {
	c := &Connector{}
	resp := c.executeSkills()
	if !contains(resp, "Skills") {
		t.Error("Expected skills message")
	}
	if !contains(resp, "60") {
		t.Error("Expected 60+ skills mention")
	}
}

func TestProcessCommand_Uninstall(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("uninstall claude")
	if !contains(resp, "Uninstall") {
		t.Error("Expected uninstall response")
	}
}

func TestProcessCommand_Remove(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("remove ollama")
	if !contains(resp, "Uninstall") {
		t.Error("Expected remove response")
	}
}

func TestProcessCommand_Rm(t *testing.T) {
	c := &Connector{}
	resp := c.processCommand("rm cursor")
	if !contains(resp, "Uninstall") {
		t.Error("Expected rm response")
	}
}

func TestExecuteUninstall(t *testing.T) {
	c := &Connector{}
	resp := c.executeUninstall("claude")
	if !contains(resp, "claude") {
		t.Error("Expected uninstall for claude")
	}
}

func TestExecuteUninstallEmpty(t *testing.T) {
	c := &Connector{}
	resp := c.executeUninstall("")
	if !contains(resp, "Usage") {
		t.Error("Expected usage message for empty args")
	}
}

func TestExecuteInstall_Groq(t *testing.T) {
	c := &Connector{}
	resp := c.executeInstall("groq")
	if !contains(resp, "GROQ") {
		t.Error("Expected groq message")
	}
	if !contains(resp, "gratis") {
		t.Error("Expected free mention")
	}
}

func TestExecuteInstall_Deepseek(t *testing.T) {
	c := &Connector{}
	resp := c.executeInstall("deepseek")
	if !contains(resp, "DEEPSEEK") {
		t.Error("Expected deepseek message")
	}
}

func TestExecuteInstall_Cursor(t *testing.T) {
	c := &Connector{}
	resp := c.executeInstall("cursor")
	if !contains(resp, "CURSOR") {
		t.Error("Expected cursor message")
	}
}

func TestExecuteInstall_Windsurf(t *testing.T) {
	c := &Connector{}
	resp := c.executeInstall("windsurf")
	if !contains(resp, "WIND") {
		t.Error("Expected windsurf message")
	}
}

func TestExecuteInstall_Gemini(t *testing.T) {
	c := &Connector{}
	resp := c.executeInstall("gemini")
	if !contains(resp, "GEMINI") {
		t.Error("Expected gemini message")
	}
}

func TestExecuteInstall_VSCode(t *testing.T) {
	c := &Connector{}
	resp := c.executeInstall("vscode")
	if !contains(resp, "VS CODE") {
		t.Error("Expected vscode message")
	}
}

func TestExecuteInstall_Antigravity(t *testing.T) {
	c := &Connector{}
	resp := c.executeInstall("antigravity")
	if !contains(resp, "ANTIGRAVITY") {
		t.Error("Expected antigravity message")
	}
}

func TestExecuteInstall_Codex(t *testing.T) {
	c := &Connector{}
	resp := c.executeInstall("codex")
	if !contains(resp, "CODEX") {
		t.Error("Expected codex message")
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
