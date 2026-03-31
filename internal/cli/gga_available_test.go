package cli

import (
	"os"
	"path/filepath"
	"testing"

<<<<<<< HEAD
	"github.com/Dxrk777/Dxrk-Hex/internal/system"
)

// TestDxrkAvailableDetectsViaLookPath verifies that dxrkAvailable returns true
// when dxrk is found on PATH via cmdLookPath.
func TestDxrkAvailableDetectsViaLookPath(t *testing.T) {
	origLookPath := cmdLookPath
	cmdLookPath = func(file string) (string, error) {
		if file == "dxrk" {
			return "/usr/local/bin/dxrk", nil
=======
	"github.com/gentleman-programming/gentle-ai/internal/system"
)

// TestGGAAvailableDetectsViaLookPath verifies that ggaAvailable returns true
// when gga is found on PATH via cmdLookPath.
func TestGGAAvailableDetectsViaLookPath(t *testing.T) {
	origLookPath := cmdLookPath
	cmdLookPath = func(file string) (string, error) {
		if file == "gga" {
			return "/usr/local/bin/gga", nil
>>>>>>> upstream/main
		}
		return "", os.ErrNotExist
	}
	t.Cleanup(func() { cmdLookPath = origLookPath })

<<<<<<< HEAD
	if !dxrkAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("dxrkAvailable() = false, want true when dxrk is on PATH")
	}
}

// TestDxrkAvailableDetectsViaLocalBin verifies that dxrkAvailable returns true
// when dxrk exists at ~/.local/bin/dxrk (default for install.sh on Linux/macOS).
func TestDxrkAvailableDetectsViaLocalBin(t *testing.T) {
=======
	if !ggaAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("ggaAvailable() = false, want true when gga is on PATH")
	}
}

// TestGGAAvailableDetectsViaLocalBin verifies that ggaAvailable returns true
// when gga exists at ~/.local/bin/gga (default for install.sh on Linux/macOS).
func TestGGAAvailableDetectsViaLocalBin(t *testing.T) {
>>>>>>> upstream/main
	tmpHome := t.TempDir()
	localBin := filepath.Join(tmpHome, ".local", "bin")
	if err := os.MkdirAll(localBin, 0o755); err != nil {
		t.Fatal(err)
	}
<<<<<<< HEAD
	if err := os.WriteFile(filepath.Join(localBin, "dxrk"), []byte("fake"), 0o755); err != nil {
=======
	if err := os.WriteFile(filepath.Join(localBin, "gga"), []byte("fake"), 0o755); err != nil {
>>>>>>> upstream/main
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

<<<<<<< HEAD
	if !dxrkAvailable(system.PlatformProfile{OS: "linux", PackageManager: "apt"}) {
		t.Fatal("dxrkAvailable() = false, want true when dxrk is at ~/.local/bin/dxrk")
	}
}

// TestDxrkAvailableDetectsViaHomebrewOptPrefix verifies that dxrkAvailable returns
// true when dxrk exists at /opt/homebrew/bin/dxrk (Apple Silicon Homebrew default).
func TestDxrkAvailableDetectsViaHomebrewOptPrefix(t *testing.T) {
	tmpDir := t.TempDir()
	fakeOptHomebrew := filepath.Join(tmpDir, "opt", "homebrew", "bin", "dxrk")
=======
	if !ggaAvailable(system.PlatformProfile{OS: "linux", PackageManager: "apt"}) {
		t.Fatal("ggaAvailable() = false, want true when gga is at ~/.local/bin/gga")
	}
}

// TestGGAAvailableDetectsViaHomebrewOptPrefix verifies that ggaAvailable returns
// true when gga exists at /opt/homebrew/bin/gga (Apple Silicon Homebrew default).
func TestGGAAvailableDetectsViaHomebrewOptPrefix(t *testing.T) {
	tmpDir := t.TempDir()
	fakeOptHomebrew := filepath.Join(tmpDir, "opt", "homebrew", "bin", "gga")
>>>>>>> upstream/main
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
<<<<<<< HEAD
		case "/opt/homebrew/bin/dxrk":
			return os.Stat(fakeOptHomebrew)
		case "/usr/local/bin/dxrk":
=======
		case "/opt/homebrew/bin/gga":
			return os.Stat(fakeOptHomebrew)
		case "/usr/local/bin/gga":
>>>>>>> upstream/main
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

<<<<<<< HEAD
	if !dxrkAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("dxrkAvailable() = false, want true when dxrk is at /opt/homebrew/bin/dxrk")
	}
}

// TestDxrkAvailableDetectsViaHomebrewUsrLocalPrefix verifies that dxrkAvailable
// returns true when dxrk exists at /usr/local/bin/dxrk (Intel Mac Homebrew default).
func TestDxrkAvailableDetectsViaHomebrewUsrLocalPrefix(t *testing.T) {
=======
	if !ggaAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("ggaAvailable() = false, want true when gga is at /opt/homebrew/bin/gga")
	}
}

// TestGGAAvailableDetectsViaHomebrewUsrLocalPrefix verifies that ggaAvailable
// returns true when gga exists at /usr/local/bin/gga (Intel Mac Homebrew default).
func TestGGAAvailableDetectsViaHomebrewUsrLocalPrefix(t *testing.T) {
>>>>>>> upstream/main
	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return t.TempDir(), nil }
	osStat = func(name string) (os.FileInfo, error) {
		switch name {
<<<<<<< HEAD
		case "/opt/homebrew/bin/dxrk":
			return nil, os.ErrNotExist
		case "/usr/local/bin/dxrk":
			// Simulate dxrk present here.
=======
		case "/opt/homebrew/bin/gga":
			return nil, os.ErrNotExist
		case "/usr/local/bin/gga":
			// Simulate gga present here.
>>>>>>> upstream/main
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

<<<<<<< HEAD
	if !dxrkAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("dxrkAvailable() = false, want true when dxrk is at /usr/local/bin/dxrk")
	}
}

// TestDxrkAvailableReturnsFalseWhenNotFound verifies that dxrkAvailable returns
// false when dxrk is not found via any detection path.
func TestDxrkAvailableReturnsFalseWhenNotFound(t *testing.T) {
=======
	if !ggaAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("ggaAvailable() = false, want true when gga is at /usr/local/bin/gga")
	}
}

// TestGGAAvailableReturnsFalseWhenNotFound verifies that ggaAvailable returns
// false when gga is not found via any detection path.
func TestGGAAvailableReturnsFalseWhenNotFound(t *testing.T) {
>>>>>>> upstream/main
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

<<<<<<< HEAD
	if dxrkAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("dxrkAvailable() = true, want false when dxrk is not installed anywhere")
	}
}

// TestDxrkAvailableBrewPathsSkippedOnLinux verifies that the Homebrew-specific
// paths (/opt/homebrew/bin/dxrk, /usr/local/bin/dxrk) are NOT checked on Linux
// even if those paths happen to exist (they never exist there in practice, but
// the guard ensures no cross-platform false positives).
func TestDxrkAvailableBrewPathsSkippedOnLinux(t *testing.T) {
=======
	if ggaAvailable(system.PlatformProfile{OS: "darwin", PackageManager: "brew"}) {
		t.Fatal("ggaAvailable() = true, want false when gga is not installed anywhere")
	}
}

// TestGGAAvailableBrewPathsSkippedOnLinux verifies that the Homebrew-specific
// paths (/opt/homebrew/bin/gga, /usr/local/bin/gga) are NOT checked on Linux
// even if those paths happen to exist (they never exist there in practice, but
// the guard ensures no cross-platform false positives).
func TestGGAAvailableBrewPathsSkippedOnLinux(t *testing.T) {
>>>>>>> upstream/main
	origLookPath := cmdLookPath
	origHomeDir := osUserHomeDir
	origStat := osStat
	cmdLookPath = func(file string) (string, error) { return "", os.ErrNotExist }
	osUserHomeDir = func() (string, error) { return t.TempDir(), nil }

	statCallCount := 0
	osStat = func(name string) (os.FileInfo, error) {
<<<<<<< HEAD
		if name == "/opt/homebrew/bin/dxrk" || name == "/usr/local/bin/dxrk" {
=======
		if name == "/opt/homebrew/bin/gga" || name == "/usr/local/bin/gga" {
>>>>>>> upstream/main
			statCallCount++
		}
		return nil, os.ErrNotExist
	}
	t.Cleanup(func() {
		cmdLookPath = origLookPath
		osUserHomeDir = origHomeDir
		osStat = origStat
	})

<<<<<<< HEAD
	dxrkAvailable(system.PlatformProfile{OS: "linux", PackageManager: "apt"})
	if statCallCount > 0 {
		t.Fatalf("dxrkAvailable() checked Homebrew paths on Linux (%d calls), expected 0", statCallCount)
=======
	ggaAvailable(system.PlatformProfile{OS: "linux", PackageManager: "apt"})
	if statCallCount > 0 {
		t.Fatalf("ggaAvailable() checked Homebrew paths on Linux (%d calls), expected 0", statCallCount)
>>>>>>> upstream/main
	}
}
