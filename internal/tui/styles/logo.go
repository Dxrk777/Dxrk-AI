package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// logoLines contains the braille ASCII art for the Dxrk cool anime smile logo.
var logoLines = []string{
	"                                                    ",
	"                                                    ",
	"                    ████████████████████            ",
	"                  ██░░░░░░░░░░░░░░░░░░░██          ",
	"                ██░░░░░░░░░░░░░░░░░░░░░░░██        ",
	"               █░░░░░░░░░░░░░░░░░░░░░░░░░░█       ",
	"              █░░░░░░░░░░░░░░░░░░░░░░░░░░░░█      ",
	"             █░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░█     ",
	"             █░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░█     ",
	"             █░░░░██░░░░░░░░░░░░░░░██░░░░░░░█     ",
	"             █░░░░░██░░░░░░░░░░░░░██░░░░░░░░░█     ",
	"              █░░░░░░██░░░░░░░░██░░░░░░░░░░░█      ",
	"              █░░░░░░░░██░░░░██░░░░░░░░░░░░░█      ",
	"               █░░░░░░░░░██░██░░░░░░░░░░░░░░█       ",
	"                █░░░░░░░░░░░░░░░░░░░░░░░░░░█       ",
	"                 █░░░░░░░░░░░░░░░░░░░░░░░░█        ",
	"                  █░░░░░░░░░████░░░░░░░░░░█         ",
	"                   █░░░░░░░░░████░░░░░░░░░█          ",
	"                    █░░░░░░░████████░░░░░░█           ",
	"                     █░░░░░░████████░░░░█             ",
	"                      █░░░░░████████░░░█               ",
	"                       █░░░░████████░░█               ",
	"                        █░░░████████░█                 ",
	"                         ██████████████                 ",
	"                          ████████████                  ",
	"                           ██████████                   ",
	"                                                    ",
	"                    ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓                  ",
	"                   ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓               ",
	"                  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓              ",
	"                 ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓             ",
	"                ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓            ",
	"               ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓           ",
	"              ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓          ",
	"             ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓         ",
	"            ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓        ",
	"           ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓       ",
	"                                                    ",
}

// gradientColors defines the gradient for the anime smile logo.
// Morado → Rosa → Rojo → Carmesí
var gradientColors = []lipgloss.Color{
	ColorDeepPurple, // Morado profundo (arriba)
	ColorMauve,      // Rosa fuerte
	ColorHotPink,    // Rosa hot
	ColorRed,        // Rojo pasión
	ColorCrimson,    // Carmesí (abajo)
	ColorMagenta,    // Magenta brillante
	ColorPeach,      // Rosa durazno
}

// RenderLogo returns the braille ASCII logo with a gradient.
// Gradient: Morado → Rosa → Rojo → Carmesí
func RenderLogo() string {
	total := len(logoLines)
	if total == 0 {
		return ""
	}

	bands := len(gradientColors)
	var b strings.Builder

	for i, line := range logoLines {
		bandIdx := (i * bands) / total
		if bandIdx >= bands {
			bandIdx = bands - 1
		}
		style := lipgloss.NewStyle().Foreground(gradientColors[bandIdx])
		b.WriteString(style.Render(line))
		if i < total-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}
