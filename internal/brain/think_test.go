package brain_test

import (
	"testing"
	"time"

	"github.com/Dxrk777/Dxrk/internal/brain"
)

func TestThinkEmptyInput(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("", b, nil, nil, nil)
	if result.Success {
		t.Error("Empty input should not succeed")
	}
}

func TestThinkWhitespaceInput(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("   ", b, nil, nil, nil)
	if result.Success {
		t.Error("Whitespace-only input should not succeed")
	}
}

func TestThinkHelp(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("help", b, nil, nil, nil)
	if !result.Success {
		t.Error("Help should succeed")
	}
	if result.Action != "help" {
		t.Errorf("Action should be 'help', got %s", result.Action)
	}
	if result.Response == "" {
		t.Error("Help should return response")
	}
}

func TestThinkStatus(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("status", b, nil, nil, nil)
	if !result.Success {
		t.Error("Status should succeed")
	}
	if result.Action != "status" {
		t.Errorf("Action should be 'status', got %s", result.Action)
	}
}

func TestThinkAgents(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("agents", b, nil, nil, nil)
	if !result.Success {
		t.Error("Agents should succeed")
	}
	if result.Action != "agents" {
		t.Errorf("Action should be 'agents', got %s", result.Action)
	}
}

func TestThinkVersion(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("version", b, nil, nil, nil)
	if !result.Success {
		t.Error("Version should succeed")
	}
	if result.Action != "version" {
		t.Errorf("Action should be 'version', got %s", result.Action)
	}
}

func TestThinkSync(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("sync", b, nil, nil, nil)
	if !result.Success {
		t.Error("Sync should succeed")
	}
	if result.Action != "sync" {
		t.Errorf("Action should be 'sync', got %s", result.Action)
	}
}

func TestThinkUpdate(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("update", b, nil, nil, nil)
	if !result.Success {
		t.Error("Update should succeed")
	}
	if result.Action != "update" {
		t.Errorf("Action should be 'update', got %s", result.Action)
	}
}

func TestThinkBackup(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, _ := brain.NewMemory(t.TempDir())

	result := brain.Think("backup", b, nil, nil, mem)
	if !result.Success {
		t.Error("Backup should succeed")
	}
	if result.Action != "backup" {
		t.Errorf("Action should be 'backup', got %s", result.Action)
	}
}

func TestThinkInstall(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("install opencode", b, nil, nil, nil)
	if !result.Success {
		t.Error("Install should succeed")
	}
	if result.Action != "install" {
		t.Errorf("Action should be 'install', got %s", result.Action)
	}
}

func TestThinkInstallClaude(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("install claude", b, nil, nil, nil)
	if !result.Success {
		t.Error("Install claude should succeed")
	}
}

func TestThinkUninstall(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("uninstall claude", b, nil, nil, nil)
	if !result.Success {
		t.Error("Uninstall should succeed")
	}
	if result.Action != "uninstall" {
		t.Errorf("Action should be 'uninstall', got %s", result.Action)
	}
}

func TestThinkHistory(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, _ := brain.NewMemory(t.TempDir())

	// Add some entries
	mem.Remember(brain.MemoryEntry{Type: "command", Content: "test command"})

	result := brain.Think("history", b, nil, nil, mem)
	if !result.Success {
		t.Error("History should succeed")
	}
	if result.Action != "history" {
		t.Errorf("Action should be 'history', got %s", result.Action)
	}
	if len(result.Memory) == 0 {
		t.Error("History should return memory entries")
	}
}

func TestThinkHistoryNoMemory(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("history", b, nil, nil, nil)
	// Should return error message but not crash
	if result.Response == "" {
		t.Error("Should return response")
	}
	if result.Success {
		t.Error("History without memory should not succeed")
	}
}

func TestThinkRemember(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, err := brain.NewMemory(t.TempDir())
	if err != nil {
		t.Fatalf("NewMemory failed: %v", err)
	}

	// Add entries
	mem.Remember(brain.MemoryEntry{Type: "install", Content: "installed opencode"})
	mem.Remember(brain.MemoryEntry{Type: "command", Content: "ran git status"})

	result := brain.Think("remember opencode", b, nil, nil, mem)
	if !result.Success {
		t.Error("Remember should succeed")
	}
	if result.Action != "query" {
		t.Errorf("Action should be 'query', got %s", result.Action)
	}
}

func TestThinkRememberNoQuery(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, _ := brain.NewMemory(t.TempDir())

	mem.Remember(brain.MemoryEntry{Type: "test", Content: "entry 1"})

	result := brain.Think("remember", b, nil, nil, mem)
	if !result.Success {
		t.Error("Remember without query should succeed")
	}
}

func TestThinkSearch(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, _ := brain.NewMemory(t.TempDir())

	mem.Remember(brain.MemoryEntry{Type: "install", Content: "installed claude"})

	result := brain.Think("search claude", b, nil, nil, mem)
	if !result.Success {
		t.Error("Search should succeed")
	}
}

func TestThinkWhat(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, _ := brain.NewMemory(t.TempDir())

	mem.Remember(brain.MemoryEntry{Type: "test", Content: "something"})

	result := brain.Think("what is this", b, nil, nil, mem)
	if !result.Success {
		t.Error("What query should succeed")
	}
}

func TestThinkUnknown(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("asdfghjkl", b, nil, nil, nil)
	// Unknown commands should still return a response (not crash)
	if result.Response == "" {
		t.Error("Unknown command should return response")
	}
}

func TestThinkEmailSend(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("send email to test@example.com subject Hello body Test", b, nil, nil, nil)
	// Should handle the command (even if email not configured)
	if result.Action != "email" {
		t.Errorf("Action should be 'email', got %s", result.Action)
	}
}

func TestThinkEmailNotConfigured(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("send email to test@example.com", b, nil, nil, nil)
	// Should return message about email not configured
	if result.Response == "" {
		t.Error("Should return message about email config")
	}
}

func TestThinkExecuteCommand(t *testing.T) {
	b := brain.New(&brain.Config{})
	c := brain.NewCommander(5 * time.Second)

	result := brain.Think("run echo hello", b, c, nil, nil)
	if !result.Success {
		t.Error("Execute should succeed")
	}
	if result.Action != "execute" {
		t.Errorf("Action should be 'execute', got %s", result.Action)
	}
}

func TestThinkExecuteWithArgs(t *testing.T) {
	b := brain.New(&brain.Config{})
	c := brain.NewCommander(5 * time.Second)

	result := brain.Think("execute whoami", b, c, nil, nil)
	if !result.Success {
		t.Error("Execute should succeed")
	}
}

func TestThinkCmd(t *testing.T) {
	b := brain.New(&brain.Config{})
	c := brain.NewCommander(5 * time.Second)

	result := brain.Think("cmd date", b, c, nil, nil)
	if !result.Success {
		t.Error("Cmd should succeed")
	}
}

func TestThinkEmailTest(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("email test", b, nil, nil, nil)
	// Should handle the command
	if result.Action != "email_test" {
		t.Errorf("Action should be 'email_test', got %s", result.Action)
	}
}

func TestThinkCaseInsensitive(t *testing.T) {
	b := brain.New(&brain.Config{})

	// Test uppercase
	result := brain.Think("HELP", b, nil, nil, nil)
	if !result.Success {
		t.Error("HELP should work")
	}

	// Test mixed case
	result = brain.Think("Status", b, nil, nil, nil)
	if !result.Success {
		t.Error("Status should work")
	}
}

func TestThinkWithLeadingTrailingSpaces(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("  status  ", b, nil, nil, nil)
	if !result.Success {
		t.Error("Status with spaces should work")
	}
}

func TestThinkMemoryQueryNoResults(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, _ := brain.NewMemory(t.TempDir())

	result := brain.Think("remember nonexistent", b, nil, nil, mem)
	if !result.Success {
		t.Error("Query with no results should succeed")
	}
}

func TestThinkEmailMultipleWordsSubject(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("email to test@example.com subject This is a long subject body message", b, nil, nil, nil)
	if result.Action != "email" {
		t.Errorf("Action should be 'email', got %s", result.Action)
	}
}

func TestThinkResult(t *testing.T) {
	result := &brain.ThinkResult{
		Response:  "test response",
		Action:    "test",
		Success:   true,
		Timestamp: time.Now(),
	}

	if result.Response != "test response" {
		t.Error("Response mismatch")
	}
	if result.Action != "test" {
		t.Error("Action mismatch")
	}
	if !result.Success {
		t.Error("Success should be true")
	}
}

func TestThinkResultWithEmail(t *testing.T) {
	result := &brain.ThinkResult{
		Response: "email sent",
		Action:   "email",
		Success:  true,
		Email: &brain.EmailResult{
			Sent:    true,
			To:      []string{"test@example.com"},
			Subject: "Test",
		},
		Timestamp: time.Now(),
	}

	if result.Email == nil {
		t.Error("Email result should not be nil")
	}
	if !result.Email.Sent {
		t.Error("Email should be marked as sent")
	}
}

func TestThinkResultWithCommand(t *testing.T) {
	result := &brain.ThinkResult{
		Response: "command output",
		Action:   "execute",
		Success:  true,
		Command: &brain.CommandResult{
			Command:   "echo test",
			Output:    "test",
			ExitCode:  0,
			Timestamp: time.Now(),
		},
	}

	if result.Command == nil {
		t.Error("Command result should not be nil")
	}
	if result.Command.ExitCode != 0 {
		t.Error("Exit code should be 0")
	}
}

// ─── Additional executeCommand tests ──────────────────────────────────────────

func TestThinkExecuteWithMemory(t *testing.T) {
	b := brain.New(&brain.Config{})
	c := brain.NewCommander(5 * time.Second)
	mem, _ := brain.NewMemory(t.TempDir())

	result := brain.Think("run echo memory test", b, c, nil, mem)
	if !result.Success {
		t.Error("Execute with memory should succeed")
	}
}

func TestThinkExecuteEmptyOutput(t *testing.T) {
	b := brain.New(&brain.Config{})
	c := brain.NewCommander(5 * time.Second)

	// "true" command succeeds with no output
	result := brain.Think("run true", b, c, nil, nil)
	if result.Action != "execute" {
		t.Errorf("Action should be 'execute', got %s", result.Action)
	}
}

func TestThinkExecuteNoArgs(t *testing.T) {
	b := brain.New(&brain.Config{})
	c := brain.NewCommander(5 * time.Second)

	result := brain.Think("run", b, c, nil, nil)
	// Should not crash, returns error response
	if result.Response == "" {
		t.Error("Should return error response")
	}
}

func TestThinkExecuteCommandNotFound(t *testing.T) {
	b := brain.New(&brain.Config{})
	c := brain.NewCommander(5 * time.Second)

	result := brain.Think("run nonexistentcommand12345xyz", b, c, nil, nil)
	// Should complete but report error
	if result.Response == "" {
		t.Error("Should return response for command not found")
	}
}

// ─── Additional handleEmail tests ───────────────────────────────────────────

func TestThinkEmailSpanish(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("correo a test@example.com", b, nil, nil, nil)
	if result.Action == "email" {
		t.Log("Email command detected (Spanish)")
	}
}

func TestThinkEmailMissingTo(t *testing.T) {
	b := brain.New(&brain.Config{})
	e := brain.NewEmailer(brain.EmailConfig{
		Host:     "smtp.test.com",
		Port:     587,
		User:     "test",
		Password: "test",
		From:     "test@test.com",
	})

	result := brain.Think("send email subject Hello body Test", b, nil, e, nil)
	if result.Success {
		t.Error("Email without recipient should not succeed")
	}
}

func TestThinkEmailWithMessageKeyword(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("email to user@test.com subject Test message This is the body", b, nil, nil, nil)
	if result.Action != "email" {
		t.Errorf("Action should be 'email', got %s", result.Action)
	}
}

func TestThinkEmailFullCommand(t *testing.T) {
	b := brain.New(&brain.Config{})

	result := brain.Think("send email to user@example.com subject 'My Subject' body 'My Message'", b, nil, nil, nil)
	if result.Action != "email" {
		t.Errorf("Action should be 'email', got %s", result.Action)
	}
}

// ─── Additional memory query tests ──────────────────────────────────────────

func TestThinkRememberEmptyQuery(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, _ := brain.NewMemory(t.TempDir())

	mem.Remember(brain.MemoryEntry{Type: "test", Content: "entry1"})
	mem.Remember(brain.MemoryEntry{Type: "test", Content: "entry2"})

	result := brain.Think("remember", b, nil, nil, mem)
	if !result.Success {
		t.Error("Remember with empty query should succeed")
	}
	// Empty query should return all history
	if len(result.Memory) == 0 && result.Response == "" {
		t.Error("Should return history or response")
	}
}

func TestThinkQueryMultipleMatches(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, _ := brain.NewMemory(t.TempDir())

	mem.Remember(brain.MemoryEntry{Type: "install", Content: "installed opencode"})
	mem.Remember(brain.MemoryEntry{Type: "install", Content: "installed claude"})
	mem.Remember(brain.MemoryEntry{Type: "command", Content: "ran git status"})

	result := brain.Think("remember installed", b, nil, nil, mem)
	if len(result.Memory) != 2 {
		t.Errorf("Should return 2 matches, got %d", len(result.Memory))
	}
}

// ─── Additional email test tests ────────────────────────────────────────────

func TestThinkEmailTestConfigured(t *testing.T) {
	b := brain.New(&brain.Config{})
	e := brain.NewEmailer(brain.EmailConfig{
		Host:     "smtp.test.com",
		Port:     587,
		User:     "test",
		Password: "test",
		From:     "test@test.com",
	})

	result := brain.Think("email test", b, nil, e, nil)
	// Should attempt to test connection
	if result.Action != "email_test" {
		t.Errorf("Action should be 'email_test', got %s", result.Action)
	}
}

// ─── formatHistory tests ────────────────────────────────────────────────────

func TestFormatHistoryEmpty(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, _ := brain.NewMemory(t.TempDir())

	result := brain.Think("history", b, nil, nil, mem)
	if result.Response == "" {
		t.Error("Should return formatted history")
	}
}

func TestFormatHistorySingleEntry(t *testing.T) {
	b := brain.New(&brain.Config{})
	mem, _ := brain.NewMemory(t.TempDir())

	mem.Remember(brain.MemoryEntry{Type: "command", Content: "echo test"})

	result := brain.Think("history", b, nil, nil, mem)
	if !contains(result.Response, "command") {
		t.Error("Response should contain entry type")
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

// ─── executeCommand edge cases ───────────────────────────────────────────────

func TestExecuteCommandWithPipes(t *testing.T) {
	b := brain.New(&brain.Config{})
	c := brain.NewCommander(5 * time.Second)

	result := brain.Think("run echo hello | cat", b, c, nil, nil)
	if result.Success {
		t.Log("Pipeline command executed")
	}
}

func TestExecuteCommandWithRedirect(t *testing.T) {
	b := brain.New(&brain.Config{})
	c := brain.NewCommander(5 * time.Second)

	// This should be blocked by security
	result := brain.Think("run echo test > /tmp/test.txt", b, c, nil, nil)
	t.Logf("Redirect result: %s", result.Response)
}
