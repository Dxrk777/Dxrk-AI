package update

import (
<<<<<<< HEAD
	"github.com/Dxrk777/Dxrk-Hex/internal/system"
=======
	"github.com/gentleman-programming/gentle-ai/internal/system"
>>>>>>> upstream/main
)

// updateHint returns a platform-specific instruction string for updating the given tool.
func updateHint(tool ToolInfo, profile system.PlatformProfile) string {
	switch tool.Name {
<<<<<<< HEAD
	case "dxrk":
		return dxrkHint(profile)
	case "engram":
		return engramHint(profile)
=======
	case "gentle-ai":
		return gentleAIHint(profile)
	case "engram":
		return engramHint(profile)
	case "gga":
		return ggaHint(profile)
>>>>>>> upstream/main
	default:
		return ""
	}
}

<<<<<<< HEAD
func dxrkHint(profile system.PlatformProfile) string {
	switch profile.OS {
	case "darwin":
		return "brew upgrade dxrk"
	case "linux":
		return "curl -fsSL https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.sh | bash"
	case "windows":
		return "irm https://raw.githubusercontent.com/Dxrk777/Dxrk-Hex/main/scripts/install-dxrk.ps1 | iex"
=======
func gentleAIHint(profile system.PlatformProfile) string {
	switch profile.OS {
	case "darwin":
		return "brew upgrade gentle-ai"
	case "linux":
		return "curl -fsSL https://raw.githubusercontent.com/Gentleman-Programming/gentle-ai/main/scripts/install.sh | bash"
	case "windows":
		return "irm https://raw.githubusercontent.com/Gentleman-Programming/gentle-ai/main/scripts/install.ps1 | iex"
>>>>>>> upstream/main
	default:
		return ""
	}
}

func engramHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade engram"
	default:
<<<<<<< HEAD
		return "dxrk upgrade (downloads pre-built binary)"
=======
		return "gentle-ai upgrade (downloads pre-built binary)"
	}
}

func ggaHint(profile system.PlatformProfile) string {
	switch profile.PackageManager {
	case "brew":
		return "brew upgrade gga"
	default:
		return "See https://github.com/Gentleman-Programming/gentleman-guardian-angel"
>>>>>>> upstream/main
	}
}
