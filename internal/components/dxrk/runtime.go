package dxrk

import (
	"fmt"
	"path/filepath"

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
}

// RuntimePRModePath returns the expected pr_mode.sh runtime path.
func RuntimePRModePath(homeDir string) string {
	return filepath.Join(RuntimeLibDir(homeDir), "pr_mode.sh")
}

// RuntimePS1Path returns the expected dxrk.ps1 path.
// On Windows, the shim goes to ~/bin/ (same dir as the bash dxrk script,
// already in PATH) so PowerShell finds it as a native command.
func RuntimePS1Path(homeDir string) string {
	return filepath.Join(homeDir, "bin", "dxrk.ps1")
}

// EnsureRuntimeAssets ensures critical dxrk runtime files are current.
//
// Behavior change from "only-if-missing" to "always-write":
// WriteFileAtomic performs a content-equality check — it is a no-op when the
// embedded asset matches the file on disk, and an atomic replace when it differs.
// This guarantees pr_mode.sh stays current after dxrk updates without
// touching the file on every sync when nothing has changed.
func EnsureRuntimeAssets(homeDir string) error {
	prModePath := RuntimePRModePath(homeDir)

	content, err := assets.Read("dxrk/pr_mode.sh")
	if err != nil {
		return fmt.Errorf("read embedded dxrk runtime asset pr_mode.sh: %w", err)
	}

	if _, err := filemerge.WriteFileAtomic(prModePath, []byte(content), 0o755); err != nil {
		return fmt.Errorf("write dxrk runtime file %q: %w", prModePath, err)
	}

	return nil
}

// EnsurePowerShellShim writes dxrk.ps1 to the Dxrk bin directory.
// Uses WriteFileAtomic: no-op when content matches, atomic replace otherwise.
// Must only be called on Windows (caller is responsible for the OS guard).
func EnsurePowerShellShim(homeDir string) error {
	ps1Path := RuntimePS1Path(homeDir)

	content, err := assets.Read("dxrk/dxrk.ps1")
	if err != nil {
		return fmt.Errorf("read embedded dxrk runtime asset dxrk.ps1: %w", err)
	}

	if _, err := filemerge.WriteFileAtomic(ps1Path, []byte(content), 0o755); err != nil {
		return fmt.Errorf("write dxrk runtime file %q: %w", ps1Path, err)
	}

	return nil
}
