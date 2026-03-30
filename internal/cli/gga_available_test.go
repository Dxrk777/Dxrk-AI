package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dxrk/dxrk/internal/system"
)

// TestDxrkAvailableDetectsViaLookPath verifies that dxrkAvailable returns true
// when dxrk is found on PATH via cmdLookPath.
func TestDxrkAvailableDetectsViaLookPath(t *testing.T) {
	origLookPath := cmdLookPath
	cmdLookPath = func(file string) (string, error) {
		if file == "dxrk" {
			return "/usr/local/bin/dxrk", nil
		}
		return "", os.ErrNotExist
	}
	t.Cleanup(func() { cmdLookPath = origLookPath })

	if !dxrkAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("dxrkAvailable() = false, want true when dxrk is on PATH")
	}
}

// TestDxrkAvailableDetectsViaLocalBin verifies that dxrkAvailable returns true
// when dxrk exists at ~/.local/bin/dxrk (default for install.sh on Linux/macOS).
func TestDxrkAvailableDetectsViaLocalBin(t *testing.T) {
	tmpHome := t.TempDir()
	localBin := filepath.Join(tmpHome, ".local", "bin")
	if err := os.MkdirAll(localBin, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(localBin, "dxrk"), []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}

	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return tmpHome, nil }
	osStat = os.Stat
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

	if !dxrkAvailable(system.PlatformProfile{OS: "linux", PackageManager: "apt"}) {
		t.Fatal("dxrkAvailable() = false, want true when dxrk is at ~/.local/bin/dxrk")
	}
}

// TestDxrkAvailableDetectsViaHomebrewOptPrefix verifies that dxrkAvailable returns
// true when dxrk exists at /opt/homebrew/bin/dxrk (Apple Silicon Homebrew default).
func TestDxrkAvailableDetectsViaHomebrewOptPrefix(t *testing.T) {
	tmpDir := t.TempDir()
	fakeOptHomebrew := filepath.Join(tmpDir, "opt", "homebrew", "bin", "dxrk")
	if err := os.MkdirAll(filepath.Dir(fakeOptHomebrew), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fakeOptHomebrew, []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}

	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return tmpDir, nil }
	// Override osStat to redirect well-known brew paths to our temp dir.
	osStat = func(name string) (os.FileInfo, error) {
		switch name {
		case "/opt/homebrew/bin/dxrk":
			return os.Stat(fakeOptHomebrew)
		case "/usr/local/bin/dxrk":
			return nil, os.ErrNotExist
		default:
			return os.Stat(name)
		}
	}
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

	if !dxrkAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("dxrkAvailable() = false, want true when dxrk is at /opt/homebrew/bin/dxrk")
	}
}

// TestDxrkAvailableDetectsViaHomebrewUsrLocalPrefix verifies that dxrkAvailable
// returns true when dxrk exists at /usr/local/bin/dxrk (Intel Mac Homebrew default).
func TestDxrkAvailableDetectsViaHomebrewUsrLocalPrefix(t *testing.T) {
	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return t.TempDir(), nil }
	osStat = func(name string) (os.FileInfo, error) {
		switch name {
		case "/opt/homebrew/bin/dxrk":
			return nil, os.ErrNotExist
		case "/usr/local/bin/dxrk":
			// Simulate dxrk present here.
			return os.Stat(os.DevNull)
		default:
			return nil, os.ErrNotExist
		}
	}
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

	if !dxrkAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("dxrkAvailable() = false, want true when dxrk is at /usr/local/bin/dxrk")
	}
}

// TestDxrkAvailableReturnsFalseWhenNotFound verifies that dxrkAvailable returns
// false when dxrk is not found via any detection path.
func TestDxrkAvailableReturnsFalseWhenNotFound(t *testing.T) {
	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return t.TempDir(), nil }
	osStat = func(name string) (os.FileInfo, error) { return nil, os.ErrNotExist }
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

	if dxrkAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("dxrkAvailable() = true, want false when dxrk is not installed anywhere")
	}
}

// TestDxrkAvailableBrewPathsSkippedOnLinux verifies that the Homebrew-specific
// paths (/opt/homebrew/bin/dxrk, /usr/local/bin/dxrk) are NOT checked on Linux
// even if those paths happen to exist (they never exist there in practice, but
// the guard ensures no cross-platform false positives).
func TestDxrkAvailableBrewPathsSkippedOnLinux(t *testing.T) {
	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return t.TempDir(), nil }

	statCallCount := 0
	osStat = func(name string) (os.FileInfo, error) {
		if name == "/opt/homebrew/bin/dxrk" || name == "/usr/local/bin/dxrk" {
			statCallCount++
		}
		return nil, os.ErrNotExist
	}
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

	dxrkAvailable(system.PlatformProfile{OS: "linux", PackageManager: "apt"})
	if statCallCount > 0 {
		t.Fatalf("dxrkAvailable() checked Homebrew paths on Linux (%d calls), expected 0", statCallCount)
	}
}
