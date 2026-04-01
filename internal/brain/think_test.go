package brain_test

import (
	"testing"
	"time"

	"github.com/Dxrk777/Dxrk-Hex/internal/brain"
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
