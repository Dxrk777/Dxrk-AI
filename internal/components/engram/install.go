package engram

import (
	"github.com/Dxrk777/Dxrk-AI/internal/installcmd"
	"github.com/Dxrk777/Dxrk-AI/internal/model"
	"github.com/Dxrk777/Dxrk-AI/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentEngram)
}
