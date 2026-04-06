package screens

import (
	"strings"

	"github.com/Dxrk777/Dxrk/internal/tui/styles"
)

// VaultOptions returns the vault menu options.
func VaultOptions() []string {
	return []string{
		"Encrypt API Key",
		"Decrypt API Key",
		"Generate Secure Password",
		"Status",
		"Back",
	}
}

// VaultState holds the state for the vault screen.
type VaultState struct {
	InputText  string
	OutputText string
	Mode       string // "encrypt", "decrypt", "password", "status"
	Error      string
}

func RenderVault(cursor int, state VaultState) string {
	var b strings.Builder

	b.WriteString(styles.RenderLogo())
	b.WriteString("\n\n")
	b.WriteString(styles.HeadingStyle.Render("🔐 Vault - Secure Storage"))
	b.WriteString("\n\n")

	// Instructions based on mode
	instructions := map[string]string{
		"encrypt":  "Enter API Key to encrypt:",
		"decrypt":  "Enter encrypted text to decrypt:",
		"password": "Press Enter to generate a secure password:",
		"status":   "Vault encryption status:",
	}

	if msg, ok := instructions[state.Mode]; ok {
		b.WriteString(styles.SubtextStyle.Render(msg))
		b.WriteString("\n\n")
	}

	if state.InputText != "" {
		b.WriteString(styles.SubtextStyle.Render("Input: " + state.InputText))
		b.WriteString("\n")
	}

	if state.OutputText != "" {
		b.WriteString(styles.SuccessStyle.Render("✓ " + state.OutputText))
		b.WriteString("\n")
	}

	if state.Error != "" {
		b.WriteString(styles.WarningStyle.Render("✗ " + state.Error))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(renderOptions(VaultOptions(), cursor))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navigate • enter: select • esc: back"))

	return styles.FrameStyle.Render(b.String())
}
