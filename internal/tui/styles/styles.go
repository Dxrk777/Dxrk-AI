package styles

import "github.com/charmbracelet/lipgloss"

// HexNeonCyber color palette - Neon Cyberpunk (Masculino)
var (
	ColorBase       = lipgloss.Color("#000000") // Negro total
	ColorSurface    = lipgloss.Color("#0a0a0a") // Negro profundo
	ColorOverlay    = lipgloss.Color("#00FFFF") // Cian brillante
	ColorText       = lipgloss.Color("#FFFFFF") // Blanco puro
	ColorSubtext    = lipgloss.Color("#00CED1") // Turquesa oscuro
	ColorLavender   = lipgloss.Color("#FF4500") // Naranja Rojo (OrangeRed)
	ColorGreen      = lipgloss.Color("#00FF00") // Verde neón
	ColorPeach      = lipgloss.Color("#FFD700") // Oro/Dorado
	ColorRed        = lipgloss.Color("#FF073A") // Rojo neón
	ColorBlue       = lipgloss.Color("#00BFFF") // Azul cielo brillante
	ColorMauve      = lipgloss.Color("#FF6600") // Naranja neón
	ColorYellow     = lipgloss.Color("#FFFF00") // Amarillo neón
	ColorTeal       = lipgloss.Color("#00FFFF") // Cian
	ColorHotPink    = lipgloss.Color("#FF00FF") // Magenta neón
	ColorDeepPurple = lipgloss.Color("#8B00FF") // Violeta neón
	ColorCrimson    = lipgloss.Color("#FF2400") // Escarlata
	ColorMagenta    = lipgloss.Color("#FF1493") // Rosa intenso
)

// Cursor is the prefix used for the currently focused item.
const Cursor = "▸ "

// Tagline returns the welcome screen tagline with the given version.
func Tagline(version string) string {
	return "DXRK HEX " + version + " — Tu compañero digital 🔥"
}

// Pre-built reusable styles - NEON CYBERPUNK
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorTeal).
			Bold(true)

	HeadingStyle = lipgloss.NewStyle().
			Foreground(ColorGreen).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(ColorBlue)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorYellow).
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
			BorderForeground(ColorTeal).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorOverlay).
			Padding(0, 1)

	ProgressFilled = lipgloss.NewStyle().
			Foreground(ColorGreen)

	ProgressEmpty = lipgloss.NewStyle().
			Foreground(ColorSurface)

	PercentStyle = lipgloss.NewStyle().
			Foreground(ColorLavender).
			Bold(true)
)
