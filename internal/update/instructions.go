package update

import (
	"github.com/Dxrk777/Dxrk-AI/internal/system"
)

// updateHint returns a platform-specific instruction string for updating the given tool.
func updateHint(tool ToolInfo, profile system.PlatformProfile) string {
	switch tool.Name {
	case "dxrk":
		return dxrkHint(profile)
	case "dxrk-memory":
		return engramHint(profile)
	case "dxrk-guardian":
		return ggaHint(profile)
	default:
		return ""
	}
}

func dxrkHint(profile system.PlatformProfile) string {
	switch profile.OS {
	case "darwin":
		return "brew upgrade dxrk"
	case "linux":
		return "curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-AI/main/scripts/install.sh | bash"
	case "windows":
		return "irm https://raw.githubusercontent.com/Dxrk777/Dxrk-AI/main/scripts/install.ps1 | iex"
	default:
		return ""
	}
}

func engramHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade dxrk-memory"
	default:
		return "dxrk upgrade (downloads pre-built binary)"
	}
}

func ggaHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade gga"
	default:
		return "See https://github.com/Dxrk777/Dxrk-AI"
	}
}
