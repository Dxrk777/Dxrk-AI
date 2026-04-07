// scripts/calculate-version.go
// ============================================================================
// Dxrk - Version Calculator
// ============================================================================
// Este script calcula el nuevo porcentaje de Dxrk basándose en el
// último release de Dxrk-AI (upstream).
//
// Sistema de versioning:
//   Dxrk-AI:     v1.15.6  →  v1.15.7  →  v1.16.0  →  v2.0.0
//   Dxrk:           045.27%  →  045.32%  →  045.82%  →  055.82%
//
// Incrementos:
//   - Major (1.x → 2.x): +10.00%
//   - Minor (1.15 → 1.16): +0.50%
//   - Patch (1.15.5 → 1.15.6): +0.05%
// ============================================================================

package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	// Archivo donde se guarda el último tag procesado
	lastTagFile = ".last-upstream-tag"

	// Archivo con el porcentaje actual de Dxrk
	currentPercentFile = ".current-percent"

	// Valores iniciales si no existen
	defaultStartPercent = 0.03
)

// Incrementos por tipo de release
const (
	incrementMajor = 10.00
	incrementMinor = 0.50
	incrementPatch = 0.05
)

func main() {
	// Obtener argumentos
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run calculate-version.go <Dxrk-AI-tag>")
		fmt.Println("Ejemplo: go run calculate-version.go v1.15.7")
		os.Exit(1)
	}

	upstreamTag := os.Args[1]

	// Extraer versión de Dxrk-AI
	version, releaseType, err := parseVersion(upstreamTag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parseando versión: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📦 Dxrk-AI detectado: %s (tipo: %s)\n", version, releaseType)

	// Obtener el último tag de Dxrk o usar default
	currentPercent := getCurrentPercent()
	fmt.Printf("📊 Porcentaje actual de Dxrk: %.2f%%\n", currentPercent)

	// Calcular incremento
	increment := getIncrement(releaseType)
	fmt.Printf("📈 Incremento: +%.2f%%\n", increment)

	// Calcular nuevo porcentaje
	newPercent := currentPercent + increment

	// Formatear como tag de Dxrk
	newTag := formatDXRKTag(newPercent)
	fmt.Printf("🎯 Nuevo tag para Dxrk: %s\n", newTag)

	// Guardar el estado para el workflow
	saveState(upstreamTag, newPercent)

	// Output para el workflow (usar para siguiente paso)
	fmt.Printf("\n::set-output name=new_tag::%s\n", newTag)
	fmt.Printf("::set-output name=new_percent::%.2f\n", newPercent)
}

func parseVersion(tag string) (string, string, error) {
	// Limpiar prefijo v si existe
	tag = strings.TrimPrefix(tag, "v")

	// Parsear versión tipo 1.15.6
	parts := strings.Split(tag, ".")
	if len(parts) < 3 {
		return "", "", fmt.Errorf("versión inválida: %s (esperado formato X.Y.Z)", tag)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", "", fmt.Errorf("major inválido: %s", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", "", fmt.Errorf("minor inválido: %s", parts[1])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", "", fmt.Errorf("patch inválido: %s", parts[2])
	}

	version := fmt.Sprintf("v%d.%d.%d", major, minor, patch)

	// Determinar tipo de release
	var releaseType string
	if patch > 0 {
		releaseType = "patch"
	} else if minor > 0 || patch > 0 {
		releaseType = "minor"
	} else {
		releaseType = "major"
	}

	return version, releaseType, nil
}

func getCurrentPercent() float64 {
	// Leer último tag de Dxrk
	data, err := os.ReadFile(currentPercentFile)
	if err != nil {
		// Intentar con git
		return getGitTagPercent()
	}

	percent, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return getGitTagPercent()
	}

	return percent
}

func getGitTagPercent() float64 {
	// Ejecutar git para obtener último tag
	// Este código está preparado para el workflow

	// Intentar leer del environment (set by workflow)
	if envVal := os.Getenv("DXRK_CURRENT_PERCENT"); envVal != "" {
		if percent, err := strconv.ParseFloat(envVal, 64); err == nil {
			return percent
		}
	}

	return defaultStartPercent
}

func getIncrement(releaseType string) float64 {
	switch releaseType {
	case "major":
		return incrementMajor
	case "minor":
		return incrementMinor
	case "patch":
		return incrementPatch
	default:
		return incrementPatch
	}
}

func formatDXRKTag(percent float64) string {
	// Formato: vXXX.XX%
	// Ejemplo: v045.27%

	// Asegurar que no pase de 999.99%
	if percent > 999.99 {
		percent = 999.99
	}

	// Formatear con 2 decimales
	formatted := fmt.Sprintf("v%06.2f%%", percent)
	return formatted
}

func saveState(upstreamTag string, percent float64) {
	// Guardar último tag de upstream
	os.WriteFile(lastTagFile, []byte(upstreamTag), 0644)

	// Guardar porcentaje actual
	os.WriteFile(currentPercentFile, []byte(fmt.Sprintf("%.2f", percent)), 0644)
}

// ============================================================================
// Funciones para usar desde otros lugares
// ============================================================================

// GetLatestDxrkAIRelease obtiene el último release de Dxrk-AI desde GitHub
func GetLatestDxrkAIRelease() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/Dxrk777/Dxrk/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parsear JSON manualmente (simplificado)
	// En producción usar github.com/google/go-github/v64

	var tagName string
	// Aquí iría el parsing del JSON...
	// Por ahora retornamos error para usar en workflow

	return tagName, fmt.Errorf("usar en workflow con curl")
}

// CompareVersions compara dos versiones de Dxrk-AI
// Retorna: "major", "minor", "patch", o "" si son iguales
func CompareVersions(old, new string) string {
	oldVer, oldType, _ := parseVersion(old)
	newVer, _, _ := parseVersion(new)

	if oldVer == newVer {
		return ""
	}

	// Si changed major, es major
	oldParts := strings.Split(strings.TrimPrefix(old, "v"), ".")
	newParts := strings.Split(strings.TrimPrefix(new, "v"), ".")

	if newParts[0] != oldParts[0] {
		return "major"
	}
	if newParts[1] != oldParts[1] {
		return "minor"
	}

	return oldType
}

// IsNewVersion verifica si newTag es más nuevo que oldTag
func IsNewVersion(oldTag, newTag string) bool {
	_, _, errOld := parseVersion(oldTag)
	_, _, errNew := parseVersion(newTag)

	if errOld != nil || errNew != nil {
		return false
	}

	oldParts := strings.Split(strings.TrimPrefix(oldTag, "v"), ".")
	newParts := strings.Split(strings.TrimPrefix(newTag, "v"), ".")

	oldMajor, _ := strconv.Atoi(oldParts[0])
	newMajor, _ := strconv.Atoi(newParts[0])
	if newMajor > oldMajor {
		return true
	}

	oldMinor, _ := strconv.Atoi(oldParts[1])
	newMinor, _ := strconv.Atoi(newParts[1])
	if newMinor > oldMinor {
		return true
	}

	oldPatch, _ := strconv.Atoi(oldParts[2])
	newPatch, _ := strconv.Atoi(newParts[2])

	return newPatch > oldPatch
}
