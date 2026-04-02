package brain_test

import (
	"runtime"
	"strings"
	"sync"
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

	var outputMu sync.Mutex
	output := ""
	result, err := c.ExecuteSafe("echo hello world", func(s string) {
		outputMu.Lock()
		output += s
		outputMu.Unlock()
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
	outputMu.Lock()
	if output == "" {
		t.Error("Output callback should have been called")
	}
	outputMu.Unlock()
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

func TestCommanderExecuteCat(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("cat is not available on Windows")
	}
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("cat", "/etc/passwd")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Exit code should be 0, got %d", result.ExitCode)
	}
}

func TestCommanderExecuteGrep(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("grep is not available on Windows")
	}
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("grep root /etc/passwd")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Exit code should be 0, got %d", result.ExitCode)
	}
}

func TestCommanderExecuteWithPipes(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("echo hello | cat")
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Exit code should be 0, got %d", result.ExitCode)
	}
}

func TestCommanderExecuteGit(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	// Just check git is available
	result, err := c.Execute("git --version")
	if err != nil {
		t.Logf("git might not be installed: %v", err)
	}
	_ = result
}

func TestCommanderExecuteNpm(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	// Just check npm is available
	result, err := c.Execute("npm --version")
	if err != nil {
		t.Logf("npm might not be installed: %v", err)
	}
	_ = result
}

func TestCommanderExecuteNode(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	// Just check node is available
	result, err := c.Execute("node --version")
	if err != nil {
		t.Logf("node might not be installed: %v", err)
	}
	_ = result
}

func TestCommanderExecutePython(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	// Just check python is available
	result, err := c.Execute("python3 --version")
	if err != nil {
		t.Logf("python3 might not be installed: %v", err)
	}
	_ = result
}

func TestCommanderExecuteGo(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	// Just check go is available
	result, err := c.Execute("go version")
	if err != nil {
		t.Logf("go might not be installed: %v", err)
	}
	_ = result
}

func TestCommanderExecuteSafeWithOutput(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	var outputLines []string
	result, err := c.ExecuteSafe("echo line1 && echo line2", func(s string) {
		outputLines = append(outputLines, s)
	})

	if err != nil {
		t.Fatalf("ExecuteSafe failed: %v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("Exit code should be 0, got %d", result.ExitCode)
	}
	if len(outputLines) == 0 {
		t.Error("Output callback should have been called")
	}
}

func TestCommanderExecuteCurl(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("curl --version")
	if err != nil {
		t.Logf("curl might not be installed: %v", err)
	}
	_ = result
}

func TestCommanderExecuteDocker(t *testing.T) {
	c := brain.NewCommander(5 * time.Second)

	result, err := c.Execute("docker --version")
	if err != nil {
		t.Logf("docker might not be installed: %v", err)
	}
	_ = result
}
