package update

// Tools is the static registry of managed tools that can be checked for updates.
//
// InstallMethod controls which upgrade strategy the executor uses:
//   - InstallBrew:      managed via homebrew (macOS/Linux with brew)
//   - InstallGoInstall: installed via `go install <GoImportPath>@version`
//   - InstallBinary:    downloaded binary from GitHub Releases (atomic replace)
//   - InstallScript:    runs an install.sh script
var Tools = []ToolInfo{
	{
		Name:          "dxrk",
		Owner:         "Dxrk777",
		Repo:          "Dxrk-AI",
		DetectCmd:     nil, // version comes from build-time ldflags
		VersionPrefix: "v",
		InstallMethod: InstallBinary,
	},
	{
		Name:          "dxrk-memory",
		Owner:         "Dxrk777",
		Repo:          "dxrk-memory",
		DetectCmd:     []string{"dxrk-memory", "version"},
		VersionPrefix: "v",
		InstallMethod: InstallBinary,
	},
	{
		Name:          "dxrk-guardian",
		Owner:         "Dxrk777",
		Repo:          "dxrk-guardian",
		DetectCmd:     []string{"dxrk-guardian", "--version"},
		VersionPrefix: "v",
		InstallMethod: InstallScript,
	},
}
