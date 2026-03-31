package styles

import "github.com/charmbracelet/lipgloss"

// HexDarkPunk color palette - Dark Punk/Gothic
var (
	ColorBase       = lipgloss.Color("#000000") // Negro total
	ColorSurface    = lipgloss.Color("#0d0d0d") // Negro profundo
	ColorOverlay    = lipgloss.Color("#1a0a1a") // Púrpura oscuro
	ColorText       = lipgloss.Color("#e0e0e0") // Blanco hueso
	ColorSubtext    = lipgloss.Color("#909090") // Gris medio
	ColorLavender   = lipgloss.Color("#9932CC") // Púrpura oscuro (DarkOrchid)
	ColorGreen      = lipgloss.Color("#8B008B") // Púrpura oscuro (DarkMagenta)
	ColorPeach      = lipgloss.Color("#C71585") // Rojo rosado oscuro (MediumVioletRed)
	ColorRed        = lipgloss.Color("#DC143C") // Carmesí
	ColorBlue       = lipgloss.Color("#800080") // Púrpura
	ColorMauve      = lipgloss.Color("#FF1493") // Rosa oscuro (DeepPink)
	ColorYellow     = lipgloss.Color("#B22222") // Rojo fuego (FireBrick)
	ColorTeal       = lipgloss.Color("#4B0082") // Índigo
	ColorHotPink    = lipgloss.Color("#FF00FF") // Magenta
	ColorDeepPurple = lipgloss.Color("#4A0E4E") // Púrpura muy oscuro
	ColorCrimson    = lipgloss.Color("#8B0000") // Rojo sangre oscuro (DarkRed)
	ColorMagenta    = lipgloss.Color("#C000C0") // Magenta oscuro
)

// Cursor is the prefix used for the currently focused item.
const Cursor = "▸ "

// Tagline returns the welcome screen tagline with the given version.
func Tagline(version string) string {
	return "DXRK HEX " + version + " — Tu compañero digital 🔥"
}

// Pre-built reusable styles - ESTILO DARK PUNK/GOTHIC
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorMagenta).
			Bold(true)

	HeadingStyle = lipgloss.NewStyle().
			Foreground(ColorLavender).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorTeal)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(ColorMauve)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorMagenta).
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
			Foreground(ColorMagenta)

	ProgressEmpty = lipgloss.NewStyle().
			Foreground(ColorOverlay)

	PercentStyle = lipgloss.NewStyle().
			Foreground(ColorHotPink).
			Bold(true)
)
