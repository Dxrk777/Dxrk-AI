package dxrk

import (
	"github.com/Dxrk777/Dxrk-AI/internal/installcmd"
	"github.com/Dxrk777/Dxrk-AI/internal/model"
	"github.com/Dxrk777/Dxrk-AI/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentDxrk)
}

func ShouldInstall(enabled bool) bool {
	return enabled
}
