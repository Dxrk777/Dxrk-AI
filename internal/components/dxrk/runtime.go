<<<<<<< HEAD:internal/components/dxrk/runtime.go
package dxrk
=======
package gga
>>>>>>> upstream/main:internal/components/gga/runtime.go

import (
	"fmt"
	"path/filepath"

<<<<<<< HEAD:internal/components/dxrk/runtime.go
	"github.com/Dxrk777/Dxrk-Hex/internal/assets"
	"github.com/Dxrk777/Dxrk-Hex/internal/components/filemerge"
)

// RuntimeLibDir returns the runtime lib path used by dxrk.
func RuntimeLibDir(homeDir string) string {
	return filepath.Join(homeDir, ".local", "share", "dxrk", "lib")
}

// RuntimeBinDir returns ~/.local/share/dxrk/bin — where Dxrk's bash script lives on Linux/Windows.
func RuntimeBinDir(homeDir string) string {
	return filepath.Join(homeDir, ".local", "share", "dxrk", "bin")
=======
	"github.com/gentleman-programming/gentle-ai/internal/assets"
	"github.com/gentleman-programming/gentle-ai/internal/components/filemerge"
)

// RuntimeLibDir returns the runtime lib path used by gga.
func RuntimeLibDir(homeDir string) string {
	return filepath.Join(homeDir, ".local", "share", "gga", "lib")
}

// RuntimeBinDir returns ~/.local/share/gga/bin — where GGA's bash script lives on Linux/Windows.
func RuntimeBinDir(homeDir string) string {
	return filepath.Join(homeDir, ".local", "share", "gga", "bin")
>>>>>>> upstream/main:internal/components/gga/runtime.go
}

// RuntimePRModePath returns the expected pr_mode.sh runtime path.
func RuntimePRModePath(homeDir string) string {
	return filepath.Join(RuntimeLibDir(homeDir), "pr_mode.sh")
}

<<<<<<< HEAD:internal/components/dxrk/runtime.go
// RuntimePS1Path returns the expected dxrk.ps1 path.
// On Windows, the shim goes to ~/bin/ (same dir as the bash dxrk script,
// already in PATH) so PowerShell finds it as a native command.
func RuntimePS1Path(homeDir string) string {
	return filepath.Join(homeDir, "bin", "dxrk.ps1")
}

// EnsureRuntimeAssets ensures critical dxrk runtime files are current.
=======
// RuntimePS1Path returns the expected gga.ps1 path.
// On Windows, the shim goes to ~/bin/ (same dir as the bash gga script,
// already in PATH) so PowerShell finds it as a native command.
func RuntimePS1Path(homeDir string) string {
	return filepath.Join(homeDir, "bin", "gga.ps1")
}

// EnsureRuntimeAssets ensures critical gga runtime files are current.
>>>>>>> upstream/main:internal/components/gga/runtime.go
//
// Behavior change from "only-if-missing" to "always-write":
// WriteFileAtomic performs a content-equality check — it is a no-op when the
// embedded asset matches the file on disk, and an atomic replace when it differs.
<<<<<<< HEAD:internal/components/dxrk/runtime.go
// This guarantees pr_mode.sh stays current after dxrk updates without
=======
// This guarantees pr_mode.sh stays current after gentle-ai updates without
>>>>>>> upstream/main:internal/components/gga/runtime.go
// touching the file on every sync when nothing has changed.
func EnsureRuntimeAssets(homeDir string) error {
	prModePath := RuntimePRModePath(homeDir)

<<<<<<< HEAD:internal/components/dxrk/runtime.go
	content, err := assets.Read("dxrk/pr_mode.sh")
	if err != nil {
		return fmt.Errorf("read embedded dxrk runtime asset pr_mode.sh: %w", err)
	}

	if _, err := filemerge.WriteFileAtomic(prModePath, []byte(content), 0o755); err != nil {
		return fmt.Errorf("write dxrk runtime file %q: %w", prModePath, err)
=======
	content, err := assets.Read("gga/pr_mode.sh")
	if err != nil {
		return fmt.Errorf("read embedded gga runtime asset pr_mode.sh: %w", err)
	}

	if _, err := filemerge.WriteFileAtomic(prModePath, []byte(content), 0o755); err != nil {
		return fmt.Errorf("write gga runtime file %q: %w", prModePath, err)
>>>>>>> upstream/main:internal/components/gga/runtime.go
	}

	return nil
}

<<<<<<< HEAD:internal/components/dxrk/runtime.go
// EnsurePowerShellShim writes dxrk.ps1 to the Dxrk bin directory.
=======
// EnsurePowerShellShim writes gga.ps1 to the GGA bin directory.
>>>>>>> upstream/main:internal/components/gga/runtime.go
// Uses WriteFileAtomic: no-op when content matches, atomic replace otherwise.
// Must only be called on Windows (caller is responsible for the OS guard).
func EnsurePowerShellShim(homeDir string) error {
	ps1Path := RuntimePS1Path(homeDir)

<<<<<<< HEAD:internal/components/dxrk/runtime.go
	content, err := assets.Read("dxrk/dxrk.ps1")
	if err != nil {
		return fmt.Errorf("read embedded dxrk runtime asset dxrk.ps1: %w", err)
	}

	if _, err := filemerge.WriteFileAtomic(ps1Path, []byte(content), 0o755); err != nil {
		return fmt.Errorf("write dxrk runtime file %q: %w", ps1Path, err)
=======
	content, err := assets.Read("gga/gga.ps1")
	if err != nil {
		return fmt.Errorf("read embedded gga runtime asset gga.ps1: %w", err)
	}

	if _, err := filemerge.WriteFileAtomic(ps1Path, []byte(content), 0o755); err != nil {
		return fmt.Errorf("write gga runtime file %q: %w", ps1Path, err)
>>>>>>> upstream/main:internal/components/gga/runtime.go
	}

	return nil
}
