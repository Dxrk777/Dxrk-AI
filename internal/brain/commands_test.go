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
