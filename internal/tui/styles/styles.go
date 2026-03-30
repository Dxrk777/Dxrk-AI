package styles

import "github.com/charmbracelet/lipgloss"

// HexRedPassion color palette - Morado y Rojo
var (
	ColorBase       = lipgloss.Color("#1a0a1a") // Fondo oscuro morado
	ColorSurface    = lipgloss.Color("#2d1b3d") // Superficie morado oscuro
	ColorOverlay    = lipgloss.Color("#6b3a5c") // Overlay morado medio
	ColorText       = lipgloss.Color("#f0e6f0") // Texto claro
	ColorSubtext    = lipgloss.Color("#c9a0dc") // Subtexto lila
	ColorLavender   = lipgloss.Color("#da70d6") // Lavanda/Rosa
	ColorGreen      = lipgloss.Color("#ff6b9d") // Rosa neon (verde改成rojo)
	ColorPeach      = lipgloss.Color("#ff1493") // Rosa intenso
	ColorRed        = lipgloss.Color("#ff0040") // Rojo pasión
	ColorBlue       = lipgloss.Color("#9932cc") // Púrpura oscuro
	ColorMauve      = lipgloss.Color("#ff69b4") // Rosa fuerte
	ColorYellow     = lipgloss.Color("#ff4500") // Rojo naranja
	ColorTeal       = lipgloss.Color("#c71585") // Rosa medio
	ColorHotPink    = lipgloss.Color("#ff1493") // Rosa hot
	ColorDeepPurple = lipgloss.Color("#8b008b") // Púrpura profundo
	ColorCrimson    = lipgloss.Color("#dc143c") // Carmesí
	ColorMagenta    = lipgloss.Color("#ff00ff") // Magenta brillante
)

// Cursor is the prefix used for the currently focused item.
const Cursor = "▸ "

// Tagline returns the welcome screen tagline with the given version.
func Tagline(version string) string {
	return "Dxrk Hex " + version + " — Tu compañero digital infinito 🔥"
}

// Pre-built reusable styles - ESTILO ROJO/MORADO PASSION
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorHotPink).
			Bold(true)

	HeadingStyle = lipgloss.NewStyle().
			Foreground(ColorMauve).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorHotPink).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorMauve)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorPeach)

	FrameStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorHotPink).
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
			Foreground(ColorHotPink).
			Bold(true)
)
