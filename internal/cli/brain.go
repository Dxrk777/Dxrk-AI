package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Dxrk777/Dxrk/internal/brain"
)

// RunBrain executes the brain command with the given arguments.
func RunBrain(args []string) error {
	// Get home directory for memory storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve home directory: %w", err)
	}

	// Create brain components
	memDir := strings.Join([]string{homeDir, ".dxrk", "memory"}, "/")

	b := brain.New(&brain.Config{
		MemoryDir:      memDir,
		CommandTimeout: 30 * time.Second,
	})

	// Initialize memory if possible
	var mem *brain.Memory
	mem, err = brain.NewMemory(memDir)
	if err != nil {
		// Memory is optional, continue without it
		mem = nil
	}

	// Create commander
	commander := brain.NewCommander(30 * time.Second)

	// Create emailer (not configured by default)
	var emailer *brain.Emailer

	// Parse and execute command
	input := strings.Join(args, " ")

	if len(args) == 0 {
		// No arguments - show help
		result := brain.Think("help", b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil
	}

	// Handle special commands
	switch strings.ToLower(args[0]) {
	case "help", "-h", "--help":
		result := brain.Think("help", b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "status":
		result := brain.Think("status", b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "history":
		result := brain.Think("history", b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "run", "execute", "cmd":
		// Extract command from args
		var cmd string
		switch args[0] {
		case "run":
			cmd = strings.Join(args[1:], " ")
		case "execute":
			cmd = strings.Join(args[1:], " ")
		case "cmd":
			cmd = strings.Join(args[1:], " ")
		}

		if cmd == "" {
			return fmt.Errorf("no command specified. Usage: dxrk brain run <command>")
		}

		result := brain.Think("run "+cmd, b, commander, emailer, mem)
		fmt.Println(result.Response)
		if result.Command != nil && result.Command.ExitCode != 0 {
			os.Exit(result.Command.ExitCode)
		}
		return nil

	case "remember", "search", "what":
		query := strings.Join(args[1:], " ")
		result := brain.Think("remember "+query, b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "email", "send":
		emailArgs := strings.Join(args[1:], " ")
		result := brain.Think("send "+emailArgs, b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "agents":
		result := brain.Think("agents", b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "version":
		result := brain.Think("version", b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "sync":
		result := brain.Think("sync", b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "update":
		result := brain.Think("update", b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "backup":
		result := brain.Think("backup", b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "install":
		pkg := strings.Join(args[1:], " ")
		result := brain.Think("install "+pkg, b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "uninstall":
		pkg := strings.Join(args[1:], " ")
		result := brain.Think("uninstall "+pkg, b, commander, emailer, mem)
		fmt.Println(result.Response)
		return nil

	case "configure", "config":
		fmt.Println("⚙️  Brain Configuration")
		fmt.Println()
		fmt.Println("Memory directory:", memDir)
		fmt.Println("Command timeout: 30s")
		fmt.Println("Email: Not configured (use dxrk brain email configure)")
		fmt.Println("Connector: Not configured (use dxrk brain connector configure)")
		return nil

	default:
		// Treat as natural language query
		result := brain.Think(input, b, commander, emailer, mem)
		fmt.Println(result.Response)
		if !result.Success {
			return fmt.Errorf("command failed")
		}
		return nil
	}
}
