package styles

import "github.com/charmbracelet/lipgloss"

// HexCyberpunk color palette - Cyberpunk (Suave)
var (
	ColorBase       = lipgloss.Color("#0a0a0a") // Negro profundo
	ColorSurface    = lipgloss.Color("#111111") // Negro suave
	ColorOverlay    = lipgloss.Color("#1a3a5c") // Azul oscuro cyber
	ColorText       = lipgloss.Color("#c0c0c0") // Gris claro
	ColorSubtext    = lipgloss.Color("#709090") // Gris azulado
	ColorLavender   = lipgloss.Color("#ff6b35") // Naranja cyber
	ColorGreen      = lipgloss.Color("#00ff9f") // Verde cyber (mint)
	ColorPeach      = lipgloss.Color("#ffd700") // Dorado suave
	ColorRed        = lipgloss.Color("#ff3366") // Rojo cyber
	ColorBlue       = lipgloss.Color("#00d4ff") // Cyan cyber
	ColorMauve      = lipgloss.Color("#ff9500") // Naranja suave
	ColorYellow     = lipgloss.Color("#ffee00") // Amarillo cyber
	ColorTeal       = lipgloss.Color("#00ffff") // Cian cyber
	ColorHotPink    = lipgloss.Color("#ff00ff") // Magenta
	ColorDeepPurple = lipgloss.Color("#9d00ff") // Violeta cyber
	ColorCrimson    = lipgloss.Color("#ff0055") // Rosa cyber
	ColorMagenta    = lipgloss.Color("#ff66b2") // Rosa suave
)

// Cursor is the prefix used for the currently focused item.
const Cursor = "▸ "

// Tagline returns the welcome screen tagline with the given version.
func Tagline(version string) string {
	return "DXRK HEX " + version + " — Tu compañero digital 🔥"
}

// Pre-built reusable styles - CYBERPUNK SUAVE
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorTeal).
			Bold(true)

	HeadingStyle = lipgloss.NewStyle().
			Foreground(ColorBlue).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(ColorGreen)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorLavender).
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
			BorderForeground(ColorBlue).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorOverlay).
			Padding(0, 1)

	ProgressFilled = lipgloss.NewStyle().
			Foreground(ColorBlue)

	ProgressEmpty = lipgloss.NewStyle().
			Foreground(ColorSurface)

	PercentStyle = lipgloss.NewStyle().
			Foreground(ColorYellow).
			Bold(true)
)
