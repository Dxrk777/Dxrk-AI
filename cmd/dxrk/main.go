// Package main is the entry point for the Dxrk Hex CLI application.
//
// Dxrk Hex es un configurador de ecosistema de IA que permite instalar y configurar
// agentes de IA como Claude Code, OpenCode, Cursor, Gemini CLI, y más.
//
// Uso:
//
//	dхrk              - Inicia la interfaz TUI interactiva
//	dхrk install      - Instala componentes seleccionados
//	dхrk update       - Verifica actualizaciones
//	dхrk upgrade      - Actualiza a una nueva versión
//	dхrk version      - Muestra la versión actual
//
// Para más información, visita: https://github.com/Dxrk777/Dxrk-Hex
package main

import (
	"fmt"
	"os"

	"github.com/Dxrk777/Dxrk-Hex/internal/app"
)

// =============================================================================
// Constantes
// =============================================================================

// version es establecida por GoReleaser via ldflags en tiempo de build.
// Ejemplo: -ldflags "-X main.version={{.Version}}"
//
// Sistema de versionado:
// - Las versiones se expresan como porcentaje (ej: 000.01%, 050.00%, 100.00%)
// - Cada release sube el porcentaje según las features agregadas
// - 000.01% = Initial Release
// - 100.00% = MVP Achieved
// Tag is the semantic version tag (v0.0.3)
// Display is what the user sees (000.03%)
var version = "000.03%"

// =============================================================================
// Función Principal
// =============================================================================

// main es el punto de entrada del binary.
//
// Flujo de ejecución:
//
// 1. ResolveVersion() determina la versión efectiva:
//   - Si se compiló con ldflags, usa esa versión
//   - Si se instaló con `go install`, usa la versión del tag
//   - Si es build local, usa "dev"
//
// 2. app.Run() inicia la aplicación:
//   - Si no hay argumentos: muestra la TUI interactiva
//   - Si hay argumentos: procesa comandos CLI
//
// 3. En caso de error:
//   - Imprime el error a stderr
//   - Sale con código de error 1
func main() {
	// Resolver la versión de la aplicación
	// Esto permite que el binary muestre la versión correcta
	// independientemente de cómo fue instalado
	app.Version = app.ResolveVersion(version)

	// Ejecutar la aplicación
	// app.Run() maneja toda la lógica de negocio:
	// - Detección del sistema operativo
	// - Verificación de plataforma soportada
	// - Auto-update check
	// - Dispatch a TUI o CLI según corresponda
	if err := app.Run(); err != nil {
		// Solo imprimimos errores a stderr
		// No usamos logging porque esto es el entry point
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// =============================================================================
// Notas de Implementación
// =============================================================================
//
// Versionado Porcentual
// --------------------
// Dxrk Hex usa un sistema de versionado basado en porcentaje en lugar de semver.
// Esto refleja el estado de desarrollo del proyecto:
//
//   000.01% - 000.99% : Inicial (Pre-alpha)
//   001.00% - 009.99% : Alpha temprana
//   010.00% - 049.99% : Alpha
//   050.00% - 079.99% : Beta
//   080.00% - 099.99% : RC (Release Candidate)
//   100.00%         : MVP (Producto Mínimo Viable)
//
// Build Flags
// -----------
// Para compilar manualmente con una versión específica:
//
//   go build -ldflags "-X main.version=0.02%" -o dxrk ./cmd/dxrk
//
// Para compilar en modo desarrollo:
//
//   go build -o dxrk ./cmd/dxrk
//   # La versión será "dev"
//
// Integración con GoReleaser
// --------------------------
// GoReleaser usa elldflags para injectar la versión en tiempo de build:
//   ldflags: -s -w -X main.version={{.Version}}
//
// El {{.Version}} viene del tag git (sin el prefijo 'v').
// =============================================================================
