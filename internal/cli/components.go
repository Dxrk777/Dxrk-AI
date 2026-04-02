package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Dxrk777/Dxrk-Hex/internal/agents"
	"github.com/Dxrk777/Dxrk-Hex/internal/components/dxrk"
	"github.com/Dxrk777/Dxrk-Hex/internal/components/engram"
	"github.com/Dxrk777/Dxrk-Hex/internal/components/mcp"
	"github.com/Dxrk777/Dxrk-Hex/internal/components/permissions"
	"github.com/Dxrk777/Dxrk-Hex/internal/components/persona"
	"github.com/Dxrk777/Dxrk-Hex/internal/components/sdd"
	"github.com/Dxrk777/Dxrk-Hex/internal/components/skills"
	"github.com/Dxrk777/Dxrk-Hex/internal/components/theme"
	"github.com/Dxrk777/Dxrk-Hex/internal/model"
	"github.com/Dxrk777/Dxrk-Hex/internal/system"
)

// applyEngram installs engram and injects it into adapters.
func applyEngram(adapters []agents.Adapter, profile system.PlatformProfile, homeDir string) error {
	engramBinaryPath, err := installEngramBinary(profile)
	if err != nil {
		return err
	}

	setupMode := engram.ParseSetupMode(os.Getenv(engram.SetupModeEnvVar))
	setupStrict := engram.ParseSetupStrict(os.Getenv(engram.SetupStrictEnvVar))

	for _, adapter := range adapters {
		if engram.ShouldAttemptSetup(setupMode, adapter.Agent()) {
			slug, _ := engram.SetupAgentSlug(adapter.Agent())
			if err := runCommand(engramBinaryPath, "setup", slug); err != nil {
				if setupStrict {
					return fmt.Errorf("engram setup for %q: %w", adapter.Agent(), err)
				}
			}
		}
		if _, err := engram.Inject(homeDir, adapter); err != nil {
			return fmt.Errorf("inject engram for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// installEngramBinary ensures engram is available and returns its path.
func installEngramBinary(profile system.PlatformProfile) (string, error) {
	if _, err := cmdLookPath("engram"); err == nil {
		return "engram", nil
	}

	if profile.PackageManager == "brew" {
		commands, err := engram.InstallCommand(profile)
		if err != nil {
			return "", fmt.Errorf("resolve install command for component %q: %w", model.ComponentEngram, err)
		}
		if err := runCommandSequence(commands); err != nil {
			return "", err
		}
		return "engram", nil
	}

	// Linux / Windows: download pre-built binary
	engramBinaryPath, err := engramDownloadFn(profile)
	if err != nil {
		return "", fmt.Errorf("download engram binary: %w", err)
	}

	// Add to PATH so subsequent commands can find it
	binDir := filepath.Dir(engramBinaryPath)
	if err := system.AddToUserPath(binDir); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: could not add %s to PATH: %v\n", binDir, err)
	}

	return engramBinaryPath, nil
}

// applyMCP injects MCP configuration (including Context7) into adapters.
func applyMCP(adapters []agents.Adapter, homeDir string) error {
	for _, adapter := range adapters {
		if _, err := mcp.Inject(homeDir, adapter); err != nil {
			return fmt.Errorf("inject mcp for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// applyPersona injects persona configuration into adapters.
func applyPersona(adapters []agents.Adapter, homeDir string, personaID model.PersonaID) error {
	for _, adapter := range adapters {
		if _, err := persona.Inject(homeDir, adapter, personaID); err != nil {
			return fmt.Errorf("inject persona for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// applyPermissions injects permissions configuration into adapters.
func applyPermissions(adapters []agents.Adapter, homeDir string) error {
	for _, adapter := range adapters {
		if _, err := permissions.Inject(homeDir, adapter); err != nil {
			return fmt.Errorf("inject permissions for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// applySDD injects SDD configuration into adapters.
func applySDD(adapters []agents.Adapter, homeDir, workspaceDir string, selection model.Selection) error {
	for _, adapter := range adapters {
		opts := sdd.InjectOptions{
			OpenCodeModelAssignments: selection.ModelAssignments,
			ClaudeModelAssignments:   selection.ClaudeModelAssignments,
			WorkspaceDir:             workspaceDir,
			StrictTDD:                selection.StrictTDD,
		}
		if _, err := sdd.Inject(homeDir, adapter, selection.SDDMode, opts); err != nil {
			return fmt.Errorf("inject sdd for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// applySkills injects skills into adapters.
func applySkills(adapters []agents.Adapter, homeDir string, skillIDs []model.SkillID) error {
	if len(skillIDs) == 0 {
		return nil
	}
	for _, adapter := range adapters {
		if _, err := skills.Inject(homeDir, adapter, skillIDs); err != nil {
			return fmt.Errorf("inject skills for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// applyDxrk installs dxrk and injects its configuration.
func applyDxrk(agentIDs []model.AgentID, adapters []agents.Adapter, profile system.PlatformProfile, homeDir string) error {
	if err := ensureDxrkInstalled(profile); err != nil {
		return err
	}

	if err := dxrk.EnsureRuntimeAssets(homeDir); err != nil {
		return fmt.Errorf("ensure dxrk runtime assets: %w", err)
	}

	if runtime.GOOS == "windows" {
		if err := dxrk.EnsurePowerShellShim(homeDir); err != nil {
			return fmt.Errorf("ensure dxrk powershell shim: %w", err)
		}
		if err := addDxrkToPath(homeDir); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: could not add dxrk to PATH: %v\n", err)
		}
	}

	if _, err := dxrk.Inject(homeDir, agentIDs); err != nil {
		return fmt.Errorf("inject dxrk config: %w", err)
	}
	return nil
}

// ensureDxrkInstalled installs dxrk if not already available.
func ensureDxrkInstalled(profile system.PlatformProfile) error {
	if dxrkAvailable(profile) {
		return nil
	}

	commands, err := dxrk.InstallCommand(profile)
	if err != nil {
		return fmt.Errorf("resolve install command for component %q: %w", model.ComponentDxrk, err)
	}

	installErr := runCommandSequence(commands)
	if installErr != nil {
		if dxrkAvailable(profile) {
			// Dxrk install uses `set -e` and `read -p` — without a TTY,
			// `read` fails but dxrk was installed successfully
			fmt.Fprintf(os.Stderr, "WARNING: dxrk install reported an error but dxrk is available — continuing. Error was: %v\n", installErr)
			return nil
		}
		return installErr
	}
	return nil
}

// addDxrkToPath adds dxrk bin directory to user PATH.
func addDxrkToPath(homeDir string) error {
	dxrkBinDir := filepath.Join(homeDir, "bin")
	return system.AddToUserPath(dxrkBinDir)
}

// applyTheme injects theme configuration into adapters.
func applyTheme(adapters []agents.Adapter, homeDir string) error {
	for _, adapter := range adapters {
		if _, err := theme.Inject(homeDir, adapter); err != nil {
			return fmt.Errorf("inject theme for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// ComponentApplier handles applying a component to adapters.
type ComponentApplier struct {
	adapters     []agents.Adapter
	agentIDs     []model.AgentID
	profile      system.PlatformProfile
	homeDir      string
	workspaceDir string
	selection    model.Selection
}

// NewComponentApplier creates a new ComponentApplier.
func NewComponentApplier(adapters []agents.Adapter, agentIDs []model.AgentID, profile system.PlatformProfile, homeDir, workspaceDir string, selection model.Selection) *ComponentApplier {
	return &ComponentApplier{
		adapters:     adapters,
		agentIDs:     agentIDs,
		profile:      profile,
		homeDir:      homeDir,
		workspaceDir: workspaceDir,
		selection:    selection,
	}
}

// Apply applies the specified component to all adapters.
func (ca *ComponentApplier) Apply(component model.ComponentID) error {
	switch component {
	case model.ComponentEngram:
		return applyEngram(ca.adapters, ca.profile, ca.homeDir)
	case model.ComponentContext7:
		return applyMCP(ca.adapters, ca.homeDir)
	case model.ComponentPersona:
		return applyPersona(ca.adapters, ca.homeDir, ca.selection.Persona)
	case model.ComponentPermission:
		return applyPermissions(ca.adapters, ca.homeDir)
	case model.ComponentSDD:
		return applySDD(ca.adapters, ca.homeDir, ca.workspaceDir, ca.selection)
	case model.ComponentSkills:
		skillIDs := selectedSkillIDs(ca.selection)
		return applySkills(ca.adapters, ca.homeDir, skillIDs)
	case model.ComponentDxrk:
		return applyDxrk(ca.agentIDs, ca.adapters, ca.profile, ca.homeDir)
	case model.ComponentTheme:
		return applyTheme(ca.adapters, ca.homeDir)
	default:
		return fmt.Errorf("component %q is not supported in install runtime", component)
	}
}
