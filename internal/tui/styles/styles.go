package styles

import "github.com/charmbracelet/lipgloss"

// HexGothic color palette - Gótico
var (
	ColorBase       = lipgloss.Color("#0a0a0a") // Negro profundo
	ColorSurface    = lipgloss.Color("#141414") // Negro gótico
	ColorOverlay    = lipgloss.Color("#2a0a2a") // Púrpura oscuro
	ColorText       = lipgloss.Color("#d0d0d0") // Gris claro
	ColorSubtext    = lipgloss.Color("#808080") // Gris medio
	ColorLavender   = lipgloss.Color("#8b0000") // Rojo sangre oscuro
	ColorGreen      = lipgloss.Color("#556b2f") // Verde oliva oscuro
	ColorPeach      = lipgloss.Color("#cd853f") // Bronce
	ColorRed        = lipgloss.Color("#b22222") // Rojo fuego
	ColorBlue       = lipgloss.Color("#4a0080") // Púrpura profundo
	ColorMauve      = lipgloss.Color("#8b008b") // Púrpura oscuro
	ColorYellow     = lipgloss.Color("#daa520") // Dorado antiguo
	ColorTeal       = lipgloss.Color("#008080") // Verde azulado oscuro
	ColorHotPink    = lipgloss.Color("#800050") // Rojo vino
	ColorDeepPurple = lipgloss.Color("#330033") // Púrpura negro
	ColorCrimson    = lipgloss.Color("#990000") // Rojo carmesí oscuro
	ColorMagenta    = lipgloss.Color("#800040") // Rojo oscuro
)

// Cursor is the prefix used for the currently focused item.
const Cursor = "▸ "

// Tagline returns the welcome screen tagline with the given version.
func Tagline(version string) string {
	return "DXRK HEX " + version + " — Tu compañero digital 🔥"
}

// Pre-built reusable styles - GÓTICO
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorMagenta).
			Bold(true)

	HeadingStyle = lipgloss.NewStyle().
			Foreground(ColorLavender).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(ColorMauve)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorCrimson).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorPeach)

	FrameStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorDeepPurple).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorOverlay).
			Padding(0, 1)

	ProgressFilled = lipgloss.NewStyle().
			Foreground(ColorLavender)

	ProgressEmpty = lipgloss.NewStyle().
			Foreground(ColorSurface)

	PercentStyle = lipgloss.NewStyle().
			Foreground(ColorCrimson).
			Bold(true)
)
