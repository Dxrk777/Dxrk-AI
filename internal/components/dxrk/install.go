package dxrk

import (
	"github.com/dxrk/dxrk/internal/installcmd"
	"github.com/dxrk/dxrk/internal/model"
	"github.com/dxrk/dxrk/internal/system"
)

func InstallCommand(profile system.PlatformProfile) ([][]string, error) {
	return installcmd.NewResolver().ResolveComponentInstall(profile, model.ComponentDxrk)
}

func ShouldInstall(enabled bool) bool {
	return enabled
}
