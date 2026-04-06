package brain

import (
	"fmt"
	"strings"
	"time"
)

// ThinkResult holds the result of a Think operation.
type ThinkResult struct {
	Response  string         `json:"response"`
	Command   *CommandResult `json:"command,omitempty"`
	Email     *EmailResult   `json:"email,omitempty"`
	Memory    []MemoryEntry  `json:"memory,omitempty"`
	Action    string         `json:"action"`
	Success   bool           `json:"success"`
	Timestamp time.Time      `json:"timestamp"`
}

// EmailResult holds the result of an email operation.
type EmailResult struct {
	Sent    bool     `json:"sent"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Error   string   `json:"error,omitempty"`
}

// Think processes user input and coordinates all brain modules.
// It handles:
// - Natural language commands
// - Shell command execution
// - Email sending
// - Memory queries
// - Status checks
func Think(input string, brain *Brain, commander *Commander, emailer *Emailer, memory *Memory) *ThinkResult {
	result := &ThinkResult{
		Timestamp: time.Now(),
		Success:   true,
	}

	input = strings.TrimSpace(input)
	if input == "" {
		result.Response = "I didn't receive any input. What would you like me to do?"
		result.Success = false
		return result
	}

	// Parse input
	lower := strings.ToLower(input)

	// Check for commands
	if strings.HasPrefix(lower, "run ") || strings.HasPrefix(lower, "execute ") || strings.HasPrefix(lower, "cmd ") {
		cmd := strings.TrimPrefix(input, "run ")
		cmd = strings.TrimPrefix(cmd, "execute ")
		cmd = strings.TrimPrefix(cmd, "cmd ")
		return executeCommand(cmd, commander, memory)
	}

	// Check for email
	if strings.HasPrefix(lower, "send ") || strings.HasPrefix(lower, "correo ") {
		return handleEmail(input, emailer, memory)
	}

	// Check for email test BEFORE generic "email " to avoid matching "email test"
	if strings.HasPrefix(lower, "email test") {
		return handleEmailTest(emailer)
	}

	if strings.HasPrefix(lower, "email ") {
		return handleEmail(input, emailer, memory)
	}

	// Check for memory query
	if strings.HasPrefix(lower, "remember ") || strings.HasPrefix(lower, "what ") || strings.HasPrefix(lower, "search ") {
		return handleMemoryQuery(input, memory)
	}

	// Handle specific commands
	switch {
	case strings.HasPrefix(lower, "help"):
		return handleHelp()

	case strings.HasPrefix(lower, "status"):
		return handleStatus(brain)

	case strings.HasPrefix(lower, "history"):
		return handleHistory(memory)

	case strings.HasPrefix(lower, "agents"):
		return handleAgents()

	case strings.HasPrefix(lower, "install "):
		return handleInstall(input)

	case strings.HasPrefix(lower, "uninstall "):
		return handleUninstall(input)

	case strings.HasPrefix(lower, "backup"):
		return handleBackup(memory)

	case strings.HasPrefix(lower, "version"):
		return handleVersion()

	case strings.HasPrefix(lower, "sync"):
		return handleSync()

	case strings.HasPrefix(lower, "update"):
		return handleUpdate()

	case strings.HasPrefix(lower, "email test"):
		return handleEmailTest(emailer)

	default:
		result.Response = fmt.Sprintf("I understand you want: '%s'\n\nI'm not sure how to help with that yet. Try:\n- 'help' - Show available commands\n- 'status' - Check system status\n- 'run <command>' - Execute a shell command\n- 'send email to <address> subject <subject>' - Send an email", input)
		result.Action = "unknown"
		return result
	}
}

// executeCommand runs a shell command.
func executeCommand(cmd string, commander *Commander, memory *Memory) *ThinkResult {
	result := &ThinkResult{
		Timestamp: time.Now(),
		Success:   true,
		Action:    "execute",
	}

	cmdResult, err := commander.Execute(cmd)
	if err != nil {
		result.Response = fmt.Sprintf("Failed to execute command: %v", err)
		result.Success = false
		result.Command = cmdResult
		return result
	}

	result.Command = cmdResult
	if cmdResult.Error != "" {
		result.Response = fmt.Sprintf("Command executed with errors:\n%s\n%s", cmdResult.Output, cmdResult.Error)
		result.Success = cmdResult.ExitCode == 0
	} else {
		result.Response = cmdResult.Output
		if result.Response == "" {
			result.Response = "Command executed successfully (no output)"
		}
	}

	// Save to memory
	if memory != nil {
		memory.Remember(MemoryEntry{
			Type:    "command",
			Content: fmt.Sprintf("Executed: %s", cmd),
			Metadata: map[string]interface{}{
				"exit_code": cmdResult.ExitCode,
				"duration":  cmdResult.Duration.String(),
			},
		})
	}

	return result
}

// handleEmail processes email commands.
func handleEmail(input string, emailer *Emailer, memory *Memory) *ThinkResult {
	result := &ThinkResult{
		Timestamp: time.Now(),
		Action:    "email",
	}

	if emailer == nil || !emailer.IsConfigured() {
		result.Response = "Email is not configured. Set up SMTP settings first."
		result.Success = false
		return result
	}

	// Parse: "send email to user@example.com subject Hello body Message"
	parts := strings.Fields(input)
	var to, subject, body string
	parsing := ""

	for i, part := range parts {
		lower := strings.ToLower(part)
		switch {
		case lower == "to":
			if i+1 < len(parts) {
				to = parts[i+1]
			}
		case lower == "subject":
			parsing = "subject"
		case lower == "body":
			parsing = "body"
		case lower == "message":
			parsing = "body"
		default:
			if parsing == "subject" {
				subject += " " + part
			} else if parsing == "body" {
				body += " " + part
			}
		}
	}

	subject = strings.TrimSpace(subject)
	body = strings.TrimSpace(body)

	if to == "" {
		result.Response = "Please specify a recipient: send email to <address> subject <subject> body <message>"
		result.Success = false
		return result
	}

	email := Email{
		To:      []string{to},
		Subject: subject,
		Body:    body,
	}

	if err := emailer.Send(email); err != nil {
		result.Response = fmt.Sprintf("Failed to send email: %v", err)
		result.Success = false
		result.Email = &EmailResult{To: []string{to}, Subject: subject, Error: err.Error()}
		return result
	}

	result.Response = fmt.Sprintf("✅ Email sent to %s", to)
	result.Success = true
	result.Email = &EmailResult{Sent: true, To: []string{to}, Subject: subject}

	// Save to memory
	if memory != nil {
		memory.Remember(MemoryEntry{
			Type:    "email",
			Content: fmt.Sprintf("Sent email to %s: %s", to, subject),
		})
	}

	return result
}

// handleMemoryQuery processes memory queries.
func handleMemoryQuery(input string, memory *Memory) *ThinkResult {
	result := &ThinkResult{
		Timestamp: time.Now(),
		Action:    "query",
	}

	if memory == nil {
		result.Response = "Memory system not available"
		result.Success = false
		return result
	}

	query := strings.TrimPrefix(strings.ToLower(input), "remember ")
	query = strings.TrimPrefix(query, "what ")
	query = strings.TrimPrefix(query, "search ")
	query = strings.TrimSpace(query)

	if query == "" {
		result.Memory = memory.History()
		result.Response = formatHistory(result.Memory)
		return result
	}

	entries := memory.Query(query)
	if len(entries) == 0 {
		result.Response = fmt.Sprintf("No results found for: '%s'", query)
		result.Success = true
		return result
	}

	result.Success = true

	result.Memory = entries
	result.Response = fmt.Sprintf("Found %d results for '%s':\n%s", len(entries), query, formatHistory(entries))
	return result
}

// handleHelp shows available commands.
func handleHelp() *ThinkResult {
	return &ThinkResult{
		Response: `📋 Dxrk AI Commands

🔧 SYSTEM:
  status       - Check system status
  version      - Show version info
  history      - Show command history
  sync         - Sync with upstream
  update       - Check for updates
  backup       - Create backup

🤖 AGENTS:
  agents       - List available agents
  install <name> - Install an agent
  uninstall <name> - Remove an agent

💻 SHELL:
  run <command> - Execute shell command
  Examples:
    run ls -la
    run git status
    run npm install

📧 EMAIL:
  send email to <address> subject <subject> body <message>
  Examples:
    send email to user@example.com subject Hello body Hi there!
    email to admin@company.com subject Report body Here's the report

🔍 MEMORY:
  remember <text> - Search memory
  history         - Show recent history
  search <query>  - Search past commands

💡 TIP: Just type your request naturally!`,
		Action:    "help",
		Success:   true,
		Timestamp: time.Now(),
	}
}

// handleStatus shows system status.
func handleStatus(brain *Brain) *ThinkResult {
	status := brain.Status()
	return &ThinkResult{
		Response: fmt.Sprintf(`✅ Dxrk AI Status

🟢 Running: %v
⏱️  Uptime: %s
💾 Memory: %v
🔐 Vault: %v
📧 Email: %v`,
			status["running"],
			status["uptime"],
			status["memory"],
			status["vault"],
			status["email"],
		),
		Action:    "status",
		Success:   true,
		Timestamp: time.Now(),
	}
}

// handleHistory shows command history.
func handleHistory(memory *Memory) *ThinkResult {
	result := &ThinkResult{
		Timestamp: time.Now(),
		Action:    "history",
	}

	if memory == nil {
		result.Response = "Memory system not available"
		result.Success = false
		return result
	}

	entries := memory.Recent(20)
	result.Memory = entries
	result.Response = formatHistory(entries)
	result.Success = true
	return result
}

// handleAgents lists available agents.
func handleAgents() *ThinkResult {
	return &ThinkResult{
		Response: `🤖 Available Agents

🎯 PROFESSIONAL:
  opencode    - OpenCode PRO ★ (Recommended)
  claude      - Claude Code (Anthropic)

📦 CLI TOOLS:
  cursor      - Cursor AI
  windsurf    - Windsurf AI
  codex       - OpenAI Codex
  gemini      - Google Gemini CLI
  vscode      - VS Code + AI
  antigravity - AGENT Framework

🆓 FREE LLM:
  ollama     - Local LLM (offline)
  groq       - Free API (fast)
  deepseek   - Free API (smart)

Install: dxrk install <agent>`,
		Action:    "agents",
		Success:   true,
		Timestamp: time.Now(),
	}
}

// handleInstall processes install command.
func handleInstall(input string) *ThinkResult {
	agent := strings.TrimPrefix(input, "install ")
	agent = strings.TrimSpace(agent)

	return &ThinkResult{
		Response: fmt.Sprintf(`📦 Install Request: %s

To install this agent, use:
  dxrk install %s

Or use the GUI: Menu → Start installation → %s`, agent, agent, agent),
		Action:    "install",
		Success:   true,
		Timestamp: time.Now(),
	}
}

// handleUninstall processes uninstall command.
func handleUninstall(input string) *ThinkResult {
	agent := strings.TrimPrefix(input, "uninstall ")
	agent = strings.TrimSpace(agent)

	return &ThinkResult{
		Response: fmt.Sprintf(`🗑️ Uninstall Request: %s

To uninstall this agent, use:
  dxrk uninstall %s`, agent, agent),
		Action:    "uninstall",
		Success:   true,
		Timestamp: time.Now(),
	}
}

// handleBackup creates a backup.
func handleBackup(memory *Memory) *ThinkResult {
	if memory != nil {
		memory.Remember(MemoryEntry{
			Type:    "backup",
			Content: "Backup requested",
		})
	}

	return &ThinkResult{
		Response: `💾 Backup

To create a backup:
  dxrk backup

Backups are stored in:
  ~/.dxrk/backups/`,
		Action:    "backup",
		Success:   true,
		Timestamp: time.Now(),
	}
}

// handleVersion shows version info.
func handleVersion() *ThinkResult {
	return &ThinkResult{
		Response: `📦 Dxrk AI
Version: 000.13%
Build: Latest

GitHub: github.com/Dxrk777/Dxrk-Hex
Upstream: github.com/Gentleman-Programming/gentle-ai`,
		Action:    "version",
		Success:   true,
		Timestamp: time.Now(),
	}
}

// handleSync syncs with upstream.
func handleSync() *ThinkResult {
	return &ThinkResult{
		Response: `🔄 Sync

Automatic sync is enabled (daily).
To force sync:
  dxrk sync

Or wait for the next automatic sync.`,
		Action:    "sync",
		Success:   true,
		Timestamp: time.Now(),
	}
}

// handleUpdate checks for updates.
func handleUpdate() *ThinkResult {
	return &ThinkResult{
		Response: `🔄 Update Check

To check for updates:
  dxrk upgrade

Dxrk AI checks for updates hourly.`,
		Action:    "update",
		Success:   true,
		Timestamp: time.Now(),
	}
}

// handleEmailTest tests email configuration.
func handleEmailTest(emailer *Emailer) *ThinkResult {
	result := &ThinkResult{
		Timestamp: time.Now(),
		Action:    "email_test",
	}

	if emailer == nil || !emailer.IsConfigured() {
		result.Response = "Email is not configured"
		result.Success = false
		return result
	}

	if err := emailer.TestConnection(); err != nil {
		result.Response = fmt.Sprintf("❌ Email test failed: %v", err)
		result.Success = false
		return result
	}

	result.Response = "✅ Email test sent successfully!"
	result.Success = true
	return result
}

// formatHistory formats memory entries for display.
func formatHistory(entries []MemoryEntry) string {
	if len(entries) == 0 {
		return "No history found."
	}

	var sb strings.Builder
	sb.WriteString("📜 Recent History:\n\n")

	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]
		sb.WriteString(fmt.Sprintf("[%s] %s - %s\n",
			entry.Timestamp.Format("15:04"),
			entry.Type,
			entry.Content,
		))
	}

	return sb.String()
}
