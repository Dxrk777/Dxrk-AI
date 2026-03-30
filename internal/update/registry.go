package update

// Tools is the static registry of managed tools that can be checked for updates.
//
// InstallMethod controls which upgrade strategy the executor uses:
//   - InstallBrew: managed via homebrew (macOS/Linux with brew)
//   - InstallGoInstall: installed via `go install <GoImportPath>@version`
//   - InstallBinary: downloaded binary from GitHub Releases (atomic replace)
//
// For brew-managed platforms the executor picks brew regardless of the
// field here; InstallMethod represents the non-brew fallback strategy.
var Tools = []ToolInfo{
	{
		Name:          "dxrk",
		Owner:         "Dxrk",
		Repo:          "dxrk",
		DetectCmd:     nil, // version comes from build-time ldflags (app.Version)
		VersionPrefix: "v",
		// dxrk: brew on macOS, binary release download on Linux/Windows.
		InstallMethod: InstallBinary,
	},
	{
		Name:          "engram",
		Owner:         "Dxrk",
		Repo:          "engram",
		DetectCmd:     []string{"engram", "version"},
		VersionPrefix: "v",
		// engram: brew on macOS/Linux-brew, binary download elsewhere.
		InstallMethod: InstallBinary,
	},
}
