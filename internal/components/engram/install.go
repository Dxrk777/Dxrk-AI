package engram

import (
	"github.com/Dxrk777/Dxrk-Hex/internal/installcmd"
	"github.com/Dxrk777/Dxrk-Hex/internal/model"
	"github.com/Dxrk777/Dxrk-Hex/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentEngram)
}
