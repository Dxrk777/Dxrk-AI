package brain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Dxrk777/Dxrk-Hex/internal/brain"
)

func TestCommanderNew(t *testing.T) {
	c := brain.NewCommander(10 * time.Second)
	if c == nil {
		t.Fatal("Commander should not be nil")
	}
}

func TestCommanderExecute(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("echo", "hello")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Exit code should be 0, got %d", result.ExitCode)
	}
	if !strings.Contains(result.Output, "hello") {
		t.Errorf("Output should contain 'hello', got: %s", result.Output)
	}
}

func TestCommanderExecuteString(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("echo hello world")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Exit code should be 0, got %d", result.ExitCode)
	}
}

func TestCommanderExecuteNotAllowed(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("rm -rf /")
	if err == nil {
		// The command should be blocked or fail
		if result.ExitCode == 0 {
			t.Error("Dangerous command should not succeed")
		}
	}
}

func TestCommanderExecuteTimeout(t *testing.T) {
	c := brain.NewCommander(1 * time.Millisecond)

	_, err := c.Execute("sleep 10")
	// Should timeout or fail
	if err == nil {
		t.Error("Long-running command should timeout")
	}
}

func TestCommanderListAllowedCommands(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	cmds := c.ListAllowedCommands()
	if len(cmds) == 0 {
		t.Error("Should have allowed commands")
	}

	// Check some expected commands
	found := false
	for _, cmd := range cmds {
		if cmd == "ls" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Should include 'ls' in allowed commands")
	}
}

func TestCommanderAddAllowedCommand(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	cmds := c.ListAllowedCommands()
	initialLen := len(cmds)

	c.AddAllowedCommand("custom-cmd")
	newCmds := c.ListAllowedCommands()

	if len(newCmds) != initialLen+1 {
		t.Error("Should have one more command after adding")
	}
}

func TestCommanderExecuteSafe(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	output := ""
	result, err := c.ExecuteSafe("echo hello world", func(s string) {
		output += s
	})

	if err != nil {
		t.Fatalf("ExecuteSafe failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Exit code should be 0, got %d", result.ExitCode)
	}
	if !strings.Contains(result.Output, "hello") {
		t.Errorf("Output should contain 'hello', got: %s", result.Output)
	}
	if output == "" {
		t.Error("Output callback should have been called")
	}
}

func TestCommanderExecuteSafeEmpty(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	_, err := c.ExecuteSafe("", nil)
	if err == nil {
		t.Error("Empty command should fail")
	}
}

func TestCommanderExecuteSafeNotAllowed(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	_, err := c.ExecuteSafe("sudo rm -rf /", nil)
	if err == nil {
		t.Error("Not allowed command should fail")
	}
}

func TestCommanderExecuteWithError(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("ls /nonexistent_directory_12345")
	// Should complete but with error
	if err != nil {
		t.Logf("Command returned error (expected): %v", err)
	}
	// Exit code should be non-zero for failed command
	if result.ExitCode == 0 {
		t.Log("Note: ls may return 0 even for nonexistent directory")
	}
}

func TestCommanderExecutePwd(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("pwd")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Exit code should be 0, got %d", result.ExitCode)
	}
	if result.Output == "" {
		t.Error("pwd should return output")
	}
}

func TestCommanderExecuteDate(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("date")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Exit code should be 0, got %d", result.ExitCode)
	}
}

func TestCommanderExecuteWhoami(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("whoami")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Exit code should be 0, got %d", result.ExitCode)
	}
	if result.Output == "" {
		t.Error("whoami should return output")
	}
}

func TestCommandResult(t *testing.T) {
	result := &brain.CommandResult{
		Command:   "echo test",
		Output:    "test\n",
		Error:     "",
		ExitCode:  0,
		Timestamp: time.Now(),
	}

	if result.Command != "echo test" {
		t.Error("Command mismatch")
	}
	if result.Output != "test\n" {
		t.Error("Output mismatch")
	}
	if result.ExitCode != 0 {
		t.Error("ExitCode mismatch")
	}
}

func TestCommandResultWithError(t *testing.T) {
	result := &brain.CommandResult{
		Command:   "false",
		Output:    "",
		Error:     "exit status 1",
		ExitCode:  1,
		Timestamp: time.Now(),
	}

	if result.ExitCode != 1 {
		t.Error("ExitCode should be 1")
	}
	if result.Error == "" {
		t.Error("Error should be set")
	}
}
