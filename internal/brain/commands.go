package brain

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// CommandResult holds the result of a command execution.
type CommandResult struct {
	Command   string        `json:"command"`
	Output    string        `json:"output"`
	Error     string        `json:"error,omitempty"`
	ExitCode  int           `json:"exit_code"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// Commander handles command execution with timeout and safety checks.
type Commander struct {
	timeout     time.Duration
	allowedCmds map[string]bool // whitelist of allowed commands
}

// NewCommander creates a new Commander instance.
func NewCommander(timeout time.Duration) *Commander {
	return &Commander{
		timeout: timeout,
		allowedCmds: map[string]bool{
			"ls":       true,
			"pwd":      true,
			"date":     true,
			"whoami":   true,
			"echo":     true,
			"cat":      true,
			"head":     true,
			"tail":     true,
			"grep":     true,
			"find":     true,
			"ps":       true,
			"df":       true,
			"free":     true,
			"uptime":   true,
			"curl":     true,
			"wget":     true,
			"git":      true,
			"npm":      true,
			"pnpm":     true,
			"yarn":     true,
			"docker":   true,
			"kubectl":  true,
			"go":       true,
			"python":   true,
			"python3":  true,
			"node":     true,
			"rustc":    true,
			"cargo":    true,
			"dxrk":     true,
			"claude":   true,
			"opencode": true,
		},
	}
}

// Execute runs a command and returns the result.
func (c *Commander) Execute(cmd string, args ...string) (*CommandResult, error) {
	if len(args) == 0 {
		// Parse command string
		parts := strings.Fields(cmd)
		if len(parts) == 0 {
			return nil, fmt.Errorf("empty command")
		}
		cmd = parts[0]
		args = parts[1:]
	}

	// Security check: verify command is allowed
	if !c.isAllowed(cmd) {
		return &CommandResult{
			Command:   cmd,
			Error:     fmt.Sprintf("command '%s' is not allowed", cmd),
			ExitCode:  1,
			Timestamp: time.Now(),
		}, fmt.Errorf("command not allowed: %s", cmd)
	}

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Create command
	execCmd := exec.CommandContext(ctx, cmd, args...)
	execCmd.Dir = "" // Use current directory

	// Capture output
	var stdout, stderr bytes.Buffer
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr

	// Execute
	start := time.Now()
	err := execCmd.Run()
	duration := time.Since(start)

	// Build result
	result := &CommandResult{
		Command:   cmd,
		Output:    stdout.String(),
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		result.Error = stderr.String()
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = "command timed out"
			result.ExitCode = 124 // standard timeout exit code
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = 1
		}
	}

	return result, nil
}

// ExecuteSafe runs a command with output callback (for streaming).
func (c *Commander) ExecuteSafe(cmd string, outputFunc func(string)) (*CommandResult, error) {
	// Parse command string
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	cmdName := parts[0]
	args := parts[1:]

	// Security check
	if !c.isAllowed(cmdName) {
		return nil, fmt.Errorf("command not allowed: %s", cmdName)
	}

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// Create command
	execCmd := exec.CommandContext(ctx, cmdName, args...)

	// Set up pipe for streaming output
	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if _, err = execCmd.StderrPipe(); err != nil {
		return nil, err
	}

	// Start command
	start := time.Now()
	if err = execCmd.Start(); err != nil {
		return &CommandResult{
			Command:   cmdName,
			Error:     err.Error(),
			ExitCode:  1,
			Timestamp: time.Now(),
		}, err
	}

	// Read output (use sync read to avoid race conditions)
	output, err := io.ReadAll(stdout)
	if err != nil {
		output = []byte{}
	}
	outputStr := string(output)

	// Call output callback if provided (for streaming simulation)
	if outputFunc != nil && outputStr != "" {
		outputFunc(outputStr)
	}

	// Wait for completion
	err = execCmd.Wait()
	duration := time.Since(start)

	result := &CommandResult{
		Command:   cmdName,
		Output:    outputStr,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			result.Error = "command timed out"
			result.ExitCode = 124
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
			result.Error = string(result.Error)
		} else {
			result.ExitCode = 1
			result.Error = err.Error()
		}
	}

	return result, nil
}

// isAllowed checks if a command is in the whitelist.
func (c *Commander) isAllowed(cmd string) bool {
	// Allow any dxrk commands
	if strings.HasPrefix(cmd, "dxrk") {
		return true
	}
	return c.allowedCmds[cmd]
}

// AddAllowedCommand adds a command to the whitelist.
func (c *Commander) AddAllowedCommand(cmd string) {
	c.allowedCmds[cmd] = true
}

// ListAllowedCommands returns all allowed commands.
func (c *Commander) ListAllowedCommands() []string {
	cmds := make([]string, 0, len(c.allowedCmds))
	for cmd := range c.allowedCmds {
		cmds = append(cmds, cmd)
	}
	return cmds
}
