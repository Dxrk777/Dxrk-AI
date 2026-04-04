package screens

import (
	"strings"

	"github.com/Dxrk777/Dxrk-Hex/internal/tui/styles"
)

// BrainState holds the state for the Brain screen.
type BrainState struct {
	// Mode: "menu", "chat", "execute", "email", "status"
	Mode string
	// Input is the current text input (for chat/execute/email modes)
	Input string
	// InputPos is the cursor position in the input
	InputPos int
	// Output is the last response from the brain
	Output string
	// Waiting indicates if we're waiting for a response
	Waiting bool
	// LastAction tracks what was last selected in menu mode
	LastAction string
	// Cursor for menu navigation
	Cursor int
}

// NewBrainState creates a new BrainState with default values.
func NewBrainState() BrainState {
	return BrainState{
		Mode:   "menu",
		Input:  "",
		Cursor: 0,
	}
}

// BrainMenuOptions returns the brain menu options.
func BrainMenuOptions() []string {
	return []string{
		"💬 Ask anything",
		"💻 Execute command",
		"📧 Send email",
		"📊 System status",
		"📜 View history",
		"⚙️  Configure",
		"🔙 Back",
	}
}

// RenderBrain renders the brain interaction screen.
func RenderBrain(state BrainState, cursor int) string {
	var b strings.Builder

	b.WriteString(styles.RenderLogo())
	b.WriteString("\n\n")
	b.WriteString(styles.HeadingStyle.Render("🧠 Dxrk Hell Brain"))
	b.WriteString("\n\n")

	switch state.Mode {
	case "menu":
		b.WriteString(styles.SubtextStyle.Render("The brain unifies memory, commands, email, and more."))
		b.WriteString("\n\n")
		b.WriteString(renderOptions(BrainMenuOptions(), cursor))
		b.WriteString("\n")
		b.WriteString(styles.HelpStyle.Render("j/k: navigate • enter: select • esc: back"))

	case "chat":
		b.WriteString(styles.SubtextStyle.Render("💬 Ask Dxrk Hell Brain"))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtextStyle.Render("Type your question or command naturally."))
		b.WriteString("\n\n")
		b.WriteString(styles.InputStyle.Render("> " + state.Input + "█"))
		b.WriteString("\n\n")
		if state.Output != "" {
			b.WriteString(styles.SubtextStyle.Render("Response:"))
			b.WriteString("\n")
			b.WriteString(state.Output)
			b.WriteString("\n\n")
		}
		if state.Waiting {
			b.WriteString(styles.HelpStyle.Render("Processing... (enter to send, esc to cancel)"))
		} else {
			b.WriteString(styles.HelpStyle.Render("enter: send • esc: cancel"))
		}

	case "execute":
		b.WriteString(styles.SubtextStyle.Render("💻 Execute Shell Command"))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtextStyle.Render("Enter a shell command to execute."))
		b.WriteString("\n\n")
		b.WriteString(styles.InputStyle.Render("> " + state.Input + "█"))
		b.WriteString("\n\n")
		if state.Output != "" {
			b.WriteString(styles.SubtextStyle.Render("Output:"))
			b.WriteString("\n")
			b.WriteString(state.Output)
			b.WriteString("\n\n")
		}
		if state.Waiting {
			b.WriteString(styles.HelpStyle.Render("Executing... (enter to send, esc to cancel)"))
		} else {
			b.WriteString(styles.HelpStyle.Render("enter: execute • esc: cancel"))
		}

	case "email":
		b.WriteString(styles.SubtextStyle.Render("📧 Send Email"))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtextStyle.Render("Format: send email to <address> subject <subject> body <message>"))
		b.WriteString("\n\n")
		b.WriteString(styles.InputStyle.Render("> " + state.Input + "█"))
		b.WriteString("\n\n")
		if state.Output != "" {
			b.WriteString(styles.SubtextStyle.Render("Result:"))
			b.WriteString("\n")
			b.WriteString(state.Output)
			b.WriteString("\n\n")
		}
		if state.Waiting {
			b.WriteString(styles.HelpStyle.Render("Sending... (enter to send, esc to cancel)"))
		} else {
			b.WriteString(styles.HelpStyle.Render("enter: send • esc: cancel"))
		}

	case "status":
		b.WriteString(styles.SubtextStyle.Render("📊 System Status"))
		b.WriteString("\n\n")
		if state.Output != "" {
			b.WriteString(state.Output)
		} else {
			b.WriteString(styles.SubtextStyle.Render("Loading status..."))
		}
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("esc: back"))

	case "history":
		b.WriteString(styles.SubtextStyle.Render("📜 Command History"))
		b.WriteString("\n\n")
		if state.Output != "" {
			b.WriteString(state.Output)
		} else {
			b.WriteString(styles.SubtextStyle.Render("No history yet."))
		}
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("esc: back"))

	case "configure":
		b.WriteString(styles.SubtextStyle.Render("⚙️  Configure Brain"))
		b.WriteString("\n\n")
		b.WriteString(styles.SubtextStyle.Render("Brain configuration options:"))
		b.WriteString("\n\n")
		b.WriteString("  • Memory Directory: ~/.dxrk/memory\n")
		b.WriteString("  • Command Timeout: 30s\n")
		b.WriteString("  • Email: Not configured\n")
		b.WriteString("  • Connector: Not configured\n")
		b.WriteString("\n")
		b.WriteString(styles.HelpStyle.Render("Configure via: dxrk brain configure"))
		b.WriteString("\n\n")
		b.WriteString(styles.HelpStyle.Render("esc: back"))
	}

	return styles.FrameStyle.Render(b.String())
}

// BrainMenuOptionCount returns the number of options in the brain menu.
func BrainMenuOptionCount() int {
	return len(BrainMenuOptions())
}
