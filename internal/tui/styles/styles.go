package styles

import "github.com/charmbracelet/lipgloss"

// HexGothic color palette - GOTHIC BLOOD
var (
	ColorBase    = lipgloss.Color("#0a0a0a") // Negro profundo
	ColorSurface = lipgloss.Color("#141414") // Negro gótico
	ColorOverlay = lipgloss.Color("#2a0a2a") // Púrpura oscuro
	ColorText    = lipgloss.Color("#ffffff") // Blanco puro
	ColorSubtext = lipgloss.Color("#cccccc") // Gris claro

	// Gothic Blood Colors - Visibles y Variados
	ColorDeepPurple  = lipgloss.Color("#330066") // Púrpura oscuro
	ColorForestGreen = lipgloss.Color("#006633") // Verde bosque
	ColorDarkTeal    = lipgloss.Color("#005555") // Verde azulado oscuro
	ColorBurgundy    = lipgloss.Color("#800040") // Burdeos
	ColorBlood       = lipgloss.Color("#aa0033") // Rojo sangre
	ColorDarkRose    = lipgloss.Color("#cc2255") // Rosa oscuro
	ColorCrimson     = lipgloss.Color("#ff3366") // Carmesí brillante

	// Aliases para compatibilidad
	ColorRed     = lipgloss.Color("#aa0033") // Rojo sangre
	ColorMagenta = lipgloss.Color("#ff3366") // Carmesí
)

// Cursor is the prefix used for the currently focused item.
const Cursor = "▸ "

// Tagline returns the welcome screen tagline with the given version.
func Tagline(version string) string {
	return "DXRK HEX " + version + " — Tu compañero digital 🔥"
}

// Pre-built reusable styles - GOTHIC BLOOD
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorMagenta).
			Bold(true)

	HeadingStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(ColorBlood)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorCrimson).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorDarkTeal)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorBurgundy)

	FrameStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorMagenta).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorOverlay).
			Padding(0, 1)

	ProgressFilled = lipgloss.NewStyle().
			Foreground(ColorForestGreen)

	ProgressEmpty = lipgloss.NewStyle().
			Foreground(ColorSurface)

	PercentStyle = lipgloss.NewStyle().
			Foreground(ColorCrimson).
			Bold(true)
)
