package screens

import (
	"fmt"
	"strings"

	"github.com/Dxrk777/Dxrk/internal/state"
	"github.com/Dxrk777/Dxrk/internal/tui/styles"
)

// MemoryOptions returns the memory menu options.
func MemoryOptions() []string {
	return []string{
		"View History",
		"View Statistics",
		"Quick Install (last setup)",
		"Clear History",
		"Back",
	}
}

// RenderMemory renders the memory/status screen.
func RenderMemory(cursor int, homeDir string) string {
	var b strings.Builder

	b.WriteString(styles.RenderLogo())
	b.WriteString("\n\n")
	b.WriteString(styles.HeadingStyle.Render("🧠 Memory System"))
	b.WriteString("\n\n")

	// Get stats
	stats, err := state.GetStats(homeDir)
	if err != nil {
		stats = map[string]any{
			"total_installs": 0,
			"total_agents":   0,
			"unique_agents":  0,
		}
	}

	// Display stats
	b.WriteString(styles.SubtextStyle.Render("Installation Statistics:"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("  📦 Total Installations: %d\n", stats["total_installs"]))
	b.WriteString(fmt.Sprintf("  🤖 Total Agents Installed: %d\n", stats["total_agents"]))
	b.WriteString(fmt.Sprintf("  🔧 Unique Agents: %d\n", stats["unique_agents"]))

	// Get last install
	last, _ := state.GetLastInstall(homeDir)
	if last != nil {
		b.WriteString("\n")
		b.WriteString(styles.SubtextStyle.Render("Last Installation:"))
		b.WriteString("\n\n")

		b.WriteString(fmt.Sprintf("  Agents: %s\n", strings.Join(last.Agents, ", ")))
		b.WriteString(fmt.Sprintf("  Preset: %s\n", last.Preset))
		if last.Success {
			b.WriteString("  Status: ✅ Success\n")
		} else {
			b.WriteString("  Status: ❌ Failed\n")
		}
	} else {
		b.WriteString("\n")
		b.WriteString(styles.WarningStyle.Render("No installations recorded yet."))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(renderOptions(MemoryOptions(), cursor))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navigate • enter: select • esc: back"))

	return styles.FrameStyle.Render(b.String())
}
