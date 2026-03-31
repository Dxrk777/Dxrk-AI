package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

<<<<<<< HEAD
	"github.com/Dxrk777/Dxrk-Hex/internal/model"
	"github.com/Dxrk777/Dxrk-Hex/internal/system"
)

// TestDxrkFixInstallErrorWhenAlreadyAvailable tests that when Dxrk install
// command fails but Dxrk is already available on the system, the error is
// swallowed and the pipeline continues instead of failing.
// This simulates the Windows scenario where install.sh fails due to TTY
// issues but Dxrk is already present.
func TestDxrkFixInstallErrorWhenAlreadyAvailable(t *testing.T) {
=======
	"github.com/gentleman-programming/gentle-ai/internal/model"
	"github.com/gentleman-programming/gentle-ai/internal/system"
)

// TestGGAFixInstallErrorWhenAlreadyAvailable tests that when GGA install
// command fails but GGA is already available on the system, the error is
// swallowed and the pipeline continues instead of failing.
// This simulates the Windows scenario where install.sh fails due to TTY
// issues but GGA is already present.
func TestGGAFixInstallErrorWhenAlreadyAvailable(t *testing.T) {
>>>>>>> upstream/main
	home := t.TempDir()

	// Save original function references
	origHome := osUserHomeDir
	origCmdLookPath := cmdLookPath
	origRunCmd := runCommand
<<<<<<< HEAD
	origDxrkAvailableCheck := dxrkAvailableCheck
=======
	origGGAAvailableCheck := ggaAvailableCheck
>>>>>>> upstream/main

	t.Cleanup(func() {
		osUserHomeDir = origHome
		cmdLookPath = origCmdLookPath
		runCommand = origRunCmd
<<<<<<< HEAD
		dxrkAvailableCheck = origDxrkAvailableCheck
=======
		ggaAvailableCheck = origGGAAvailableCheck
>>>>>>> upstream/main
	})

	// Setup mocks
	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(string) (string, error) {
		return "", errors.New("not found")
	}

	// Track if runCommand was called and capture its error
	runCommandCalled := false
	runCommand = func(name string, args ...string) error {
		runCommandCalled = true
		// Simulate install.sh failing due to TTY issue
		return errors.New("exit status 1: read: open /dev/tty: no such device or address")
	}

<<<<<<< HEAD
	// Make dxrkAvailable return false initially (simulating install needed),
	// then return true after the "install" (simulating Dxrk was already there)
	dxrkAvailableCheck = func(profile system.PlatformProfile) bool {
		// After install command runs, simulate Dxrk being available
		// (this is the fix scenario: install failed but Dxrk is there)
=======
	// Make ggaAvailable return false initially (simulating install needed),
	// then return true after the "install" (simulating GGA was already there)
	ggaAvailableCheck = func(profile system.PlatformProfile) bool {
		// After install command runs, simulate GGA being available
		// (this is the fix scenario: install failed but GGA is there)
>>>>>>> upstream/main
		return runCommandCalled
	}

	// Create a minimal config so the test can run
	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(configPath, []byte("{}"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Create the runtime manually to test the component step
	profile := system.PlatformProfile{OS: "windows", PackageManager: "winget"}
	step := componentApplyStep{
<<<<<<< HEAD
		id:           "component:dxrk",
		component:    model.ComponentDxrk,
=======
		id:           "component:gga",
		component:    model.ComponentGGA,
>>>>>>> upstream/main
		homeDir:      home,
		workspaceDir: home,
		agents:       []model.AgentID{model.AgentOpenCode},
		selection:    model.Selection{},
		profile:      profile,
	}

	err := step.Run()

<<<<<<< HEAD
	// Verify: no error should be returned (fix: error swallowed when Dxrk available)
	if err != nil {
		t.Fatalf("Run() error = %v, want nil (error should be swallowed when Dxrk is available)", err)
=======
	// Verify: no error should be returned (fix: error swallowed when GGA available)
	if err != nil {
		t.Fatalf("Run() error = %v, want nil (error should be swallowed when GGA is available)", err)
>>>>>>> upstream/main
	}

	// Verify: runCommand was called (we attempted install)
	if !runCommandCalled {
		t.Fatal("runCommand was not called, expected install to be attempted")
	}
}

<<<<<<< HEAD
// TestDxrkFixInstallErrorWhenNotAvailable tests that when Dxrk install
// command fails and Dxrk is NOT available, the error is NOT swallowed
// and is returned to the caller. This ensures we don't mask real errors.
func TestDxrkFixInstallErrorWhenNotAvailable(t *testing.T) {
=======
// TestGGAFixInstallErrorWhenNotAvailable tests that when GGA install
// command fails and GGA is NOT available, the error is NOT swallowed
// and is returned to the caller. This ensures we don't mask real errors.
func TestGGAFixInstallErrorWhenNotAvailable(t *testing.T) {
>>>>>>> upstream/main
	home := t.TempDir()

	origHome := osUserHomeDir
	origRunCmd := runCommand
<<<<<<< HEAD
	origDxrkAvailableCheck := dxrkAvailableCheck
=======
	origGGAAvailableCheck := ggaAvailableCheck
>>>>>>> upstream/main
	origCmdLookPath := cmdLookPath

	t.Cleanup(func() {
		osUserHomeDir = origHome
		runCommand = origRunCmd
<<<<<<< HEAD
		dxrkAvailableCheck = origDxrkAvailableCheck
=======
		ggaAvailableCheck = origGGAAvailableCheck
>>>>>>> upstream/main
		cmdLookPath = origCmdLookPath
	})

	osUserHomeDir = func() (string, error) { return home, nil }
	cmdLookPath = func(string) (string, error) {
		return "", errors.New("not found")
	}

<<<<<<< HEAD
	// Make dxrkAvailable ALWAYS return false (Dxrk is not available)
	dxrkAvailableCheck = func(profile system.PlatformProfile) bool {
=======
	// Make ggaAvailable ALWAYS return false (GGA is not available)
	ggaAvailableCheck = func(profile system.PlatformProfile) bool {
>>>>>>> upstream/main
		return false
	}

	// Simulate a REAL install error (not the TTY issue)
	runCommand = func(name string, args ...string) error {
		return errors.New("network error: connection refused")
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(configPath, []byte("{}"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	profile := system.PlatformProfile{OS: "windows", PackageManager: "winget"}
	step := componentApplyStep{
<<<<<<< HEAD
		id:           "component:dxrk",
		component:    model.ComponentDxrk,
=======
		id:           "component:gga",
		component:    model.ComponentGGA,
>>>>>>> upstream/main
		homeDir:      home,
		workspaceDir: home,
		agents:       []model.AgentID{model.AgentOpenCode},
		selection:    model.Selection{},
		profile:      profile,
	}

	err := step.Run()

	// Verify: error should be returned (not swallowed)
	if err == nil {
<<<<<<< HEAD
		t.Fatal("Run() expected error when Dxrk is not available and install fails, got nil")
=======
		t.Fatal("Run() expected error when GGA is not available and install fails, got nil")
>>>>>>> upstream/main
	}

	if !strings.Contains(err.Error(), "network error") {
		t.Fatalf("Expected network error in message, got: %v", err)
	}
}
