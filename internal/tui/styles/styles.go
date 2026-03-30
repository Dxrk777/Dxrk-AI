package styles

import "github.com/charmbracelet/lipgloss"

// HexRedPassion color palette - Dark Punk/Gothic
var (
	ColorBase       = lipgloss.Color("#0a0a0a") // Negro profundo
	ColorSurface    = lipgloss.Color("#1a0a0a") // Negro con tinte rojo
	ColorOverlay    = lipgloss.Color("#2a0a0a") // Rojo oscuro
	ColorText       = lipgloss.Color("#f0f0f0") // Blanco hueso
	ColorSubtext    = lipgloss.Color("#c0c0c0") // Gris claro
	ColorLavender   = lipgloss.Color("#ff3333") // Rojo brillante
	ColorGreen      = lipgloss.Color("#ff0000") // Rojo intenso
	ColorPeach      = lipgloss.Color("#ff4444") // Rojo claro
	ColorRed        = lipgloss.Color("#cc0000") // Rojo oscuro
	ColorBlue       = lipgloss.Color("#8b0000") // Rojo sangre
	ColorMauve      = lipgloss.Color("#ff6666") // Rojo rosado
	ColorYellow     = lipgloss.Color("#cc3333") // Rojo marrón
	ColorTeal       = lipgloss.Color("#aa0000") // Rojo medio
	ColorHotPink    = lipgloss.Color("#ff0000") // Rojo brillante
	ColorDeepPurple = lipgloss.Color("#660000") // Rojo muy oscuro
	ColorCrimson    = lipgloss.Color("#dc143c") // Carmesí
	ColorMagenta    = lipgloss.Color("#ff0033") // Rojo magenta
)

// Cursor is the prefix used for the currently focused item.
const Cursor = "▸ "

// Tagline returns the welcome screen tagline with the given version.
func Tagline(version string) string {
	return "DXRK HEX " + version + " — Tu compañero digital oscuro 🔥"
}

// Pre-built reusable styles - ESTILO DARK PUNK/GOTHIC
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	HeadingStyle = lipgloss.NewStyle().
			Foreground(ColorCrimson).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorRed).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorCrimson)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorPeach)

	FrameStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorRed).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorOverlay).
			Padding(0, 1)

	ProgressFilled = lipgloss.NewStyle().
			Foreground(ColorRed)

	ProgressEmpty = lipgloss.NewStyle().
			Foreground(ColorOverlay)

	PercentStyle = lipgloss.NewStyle().
			Foreground(ColorCrimson).
			Bold(true)
)
