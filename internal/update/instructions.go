package update

import (
	"github.com/Dxrk777/Dxrk-Hex/internal/system"
)

// updateHint returns a platform-specific instruction string for updating the given tool.
func updateHint(tool ToolInfo, profile system.PlatformProfile) string {
	switch tool.Name {
	case "dxrk":
		return dxrkHint(profile)
	case "engram":
		return engramHint(profile)
	default:
		return ""
	}
}

func dxrkHint(profile system.PlatformProfile) string {
	switch profile.OS {
	case "darwin":
		return "brew upgrade dxrk"
	case "linux":
		return "curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.sh | bash"
	case "windows":
		return "irm https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.ps1 | iex"
	default:
		return ""
	}
}

func engramHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade engram"
	default:
		return "dxrk upgrade (downloads pre-built binary)"
	}
}
