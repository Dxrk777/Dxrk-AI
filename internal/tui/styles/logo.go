package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// logoLines contains the braille ASCII art for the Dxrk Dark Skull punk logo.
var logoLines = []string{
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠑⠀⠈⣴⣿⡄⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠿⠟⠛⠛⠩⠙⠀⠠⣤⣖⡂⠁⠀⠀⢠⠀⠀⡅⠄⡅⠀⢼⠀⠀⠅⠠⣸⠀⢀⣿⣾⢟⣿",
	"⣿⣿⣿⣿⣿⣿⣹⣿⣿⣇⠀⢀⣾⠟⣉⠀⠀⣭⣩⣿⡛⢿⣿⣿⣿⣿⡋⣄⣦⠟⠡⠴⠾⢿⣿⣷⡎⢀⠄⠀⠀⠀⠀⠘⠀⢀⡏⠀⠀⠀⡏⡄⠀⠀⠀⣿⠀⢸⣻⢰⣿⣿",
	"⣾⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⡿⠋⠰⠁⠀⠀⣉⣙⠛⠻⣿⣿⣿⣿⣿⠟⢉⣤⡶⢖⣚⣶⣆⠈⠻⣿⣄⠄⠀⢀⡀⠀⠀⠀⠘⠇⠀⠀⠀⠘⡠⠀⠀⠀⡏⠀⠈⣶⢸⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⡟⢀⠸⠁⠀⠁⠀⢁⡇⠛⡻⠶⠤⠈⠻⠿⠿⠟⢀⣿⠿⠉⠁⣰⠀⠤⠀⠄⢀⣤⡆⠀⠀⠀⢠⣄⠄⠔⡠⠀⠀⠈⠐⠂⠀⠁⠀⠃⠒⣶⣽⢸⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⠟⣡⡅⠁⠀⢂⠀⠀⡀⠀⡀⡄⣂⣠⣄⠀⢰⣆⡄⠐⢱⣷⡄⣠⣶⣖⢀⣀⡾⡸⢛⠁⠀⠈⠒⠀⠀⠀⠀⡅⠀⢀⢲⣿⠆⠀⠘⢰⣿⡗⡷⣿⣸⣿⣿",
	"⣿⣿⣿⣿⣿⡟⣱⡆⣷⠀⠀⠠⡏⠀⠀⡇⠀⢀⢹⣿⣿⡟⡀⠞⡼⣿⡄⢼⣿⣿⣾⣿⣿⣿⡿⢣⡄⠐⠛⠀⠐⠀⠀⢀⠀⠀⡀⢀⡑⠸⠀⠀⠄⠀⠿⢃⢳⣿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣇⣿⡀⢇⠘⠏⠀⠰⢥⣃⠘⠿⢿⣛⣔⣽⠽⠯⢱⣿⡢⠨⣟⣛⣛⣛⢫⠶⡇⠀⠀⠀⠀⡠⠂⣠⣾⠀⠀⡁⠈⠀⣰⠋⠀⣀⣖⡵⣎⢽⣿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⠈⡀⠀⠀⠀⢼⢉⡈⣿⣿⣿⣿⡹⣷⣞⠐⢨⡿⢸⣿⣿⣿⣿⣾⡤⡧⢘⢀⠀⠀⠠⣤⠿⠿⠀⠀⠃⠐⠁⠀⡄⢠⣭⡩⣾⣿⣿⣿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⠀⠐⠀⠀⠈⢑⠃⠘⣿⣿⣿⣿⢟⣾⣰⣤⣶⣿⣿⣿⣿⣿⣿⠁⠇⠀⠈⢀⣴⣶⣿⡆⠰⠄⢰⢀⠄⣴⠇⢠⢲⣻⣷⡹⣿⣿⣿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⡢⠀⠀⢠⠀⡇⡠⢹⣿⠟⣵⣶⣾⣦⣝⡿⣿⣿⣿⣿⣿⣇⠀⠀⠀⠀⠘⠛⠌⠀⣺⡗⠂⠈⢈⠀⣴⡆⠸⣸⢸⣿⣷⣹⣿⣿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣇⡖⠄⡙⣼⡇⡇⣇⠹⡨⠻⠭⢭⡁⠙⠛⠠⠙⢻⣿⣿⣷⢠⡀⠉⢁⠀⢠⣬⡥⠄⠀⠀⠠⢻⢀⢸⠇⠄⢟⣼⣿⢟⣹⣿⣿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡤⠀⢸⡜⣷⣿⣿⡇⠁⠚⢿⣿⣿⣭⠤⣀⣴⠞⠙⠛⡏⢠⢄⠀⠠⣄⢠⣤⠴⠢⡂⠁⣄⠈⠀⠀⠀⠄⠹⣿⣷⣿⣿⣿⣿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡗⡆⢸⣿⣹⣿⣿⢡⠀⠑⣶⣶⡶⠒⠞⠛⠃⠀⠀⢰⣧⢸⣿⠀⠴⠟⡘⢁⢰⣿⡃⠀⠀⠀⠄⠀⠀⠀⢆⢹⢹⣿⣿⣿⣿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣟⢃⢸⣿⣿⣿⣿⠘⠀⣏⡜⢏⡆⠠⠤⣤⠀⠆⠀⠸⠇⢈⣉⠀⡶⠋⣠⣼⣾⣉⠉⢀⠈⠀⠀⠀⠀⠘⣼⣻⢸⣿⣿⣿⣿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⠋⣢⢻⣿⣿⣿⣿⠰⣤⣿⣿⠈⡑⠐⠁⠈⢀⣠⡄⠀⠸⣿⣿⠘⢳⣇⣿⣿⣿⣿⢦⡄⠀⠀⠀⠀⠀⠀⢘⣿⡸⣿⣿⣿⣿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣮⣿⡗⣸⣿⣿⣿⢰⡇⣼⣿⣶⡇⠠⣷⢸⣏⣿⡗⡆⠰⠄⣮⡰⣼⣧⣿⣿⠀⡛⠸⡊⠀⠀⠀⠀⠀⠀⠈⢿⣧⡇⣿⣿⢿⣿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣾⣿⣿⣿⣿⠘⣧⣿⣿⣿⣿⡇⡟⢸⣿⣿⣷⣷⠸⡆⢹⠎⢙⠛⢸⣿⠀⠐⢈⣧⠄⠀⠀⡄⠀⠀⠀⠀⠙⠇⠿⠿⡶⡿⣿⣿⣿",
	"⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⣿⣿⣿⣿⣿⣿⡧⣿⣿⣿⠿⠛⠀⠃⢸⢀⣙⡷⠀⡃⢀⣤⡧⣀⣀⣀⠀⠀⠀⠀⠀⠀⠁⠀⠀⠀⠁⠃⢰⣧⣿",
}

// gradientColors defines the gradient for the skull logo.
var gradientColors = []lipgloss.Color{
	ColorBase,       // Negro profundo
	ColorDeepPurple, // Rojo muy oscuro
	ColorRed,        // Rojo oscuro
	ColorCrimson,    // Carmesí
	ColorRed,        // Rojo
	ColorDeepPurple, // Rojo oscuro
	ColorBase,       // Negro profundo
}

// RenderLogo returns the braille ASCII logo with a gradient.
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
