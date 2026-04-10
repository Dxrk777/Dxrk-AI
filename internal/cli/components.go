package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Dxrk777/Dxrk-AI/internal/agents"
	"github.com/Dxrk777/Dxrk-AI/internal/components/dxrk"
	"github.com/Dxrk777/Dxrk-AI/internal/components/dxrk-memory"
	"github.com/Dxrk777/Dxrk-AI/internal/components/mcp"
	"github.com/Dxrk777/Dxrk-AI/internal/components/permissions"
	"github.com/Dxrk777/Dxrk-AI/internal/components/persona"
	"github.com/Dxrk777/Dxrk-AI/internal/components/sdd"
	"github.com/Dxrk777/Dxrk-AI/internal/components/skills"
	"github.com/Dxrk777/Dxrk-AI/internal/components/theme"
	"github.com/Dxrk777/Dxrk-AI/internal/model"
	"github.com/Dxrk777/Dxrk-AI/internal/system"
)

// applyEngram installs the dxrk-memory binary (if needed) and injects dxrk-memory configuration
// into each agent adapter. On brew systems, it uses `brew install`; on Linux/Windows,
// it downloads a pre-built binary from GitHub releases.
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
					return fmt.Errorf("dxrk-memory setup for %q: %w", adapter.Agent(), err)
				}
			}
		}
		if _, err := engram.Inject(homeDir, adapter); err != nil {
			return fmt.Errorf("inject dxrk-memory for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// installEngramBinary checks if dxrk-memory is already available on PATH. If not, it installs
// it via brew (macOS/Linux with Homebrew) or downloads a pre-built binary from GitHub.
// Returns the path to the dxrk-memory binary.
func installEngramBinary(profile system.PlatformProfile) (string, error) {
	if _, err := cmdLookPath("dxrk-memory"); err == nil {
		return "dxrk-memory", nil
	}

	if profile.PackageManager == "brew" {
		commands, err := engram.InstallCommand(profile)
		if err != nil {
			return "", fmt.Errorf("resolve install command for component %q: %w", model.ComponentEngram, err)
		}
		if err := runCommandSequence(commands); err != nil {
			return "", err
		}
		return "dxrk-memory", nil
	}

	// Linux / Windows: download pre-built binary
	engramBinaryPath, err := engramDownloadFn(profile)
	if err != nil {
		return "", fmt.Errorf("download dxrk-memory binary: %w", err)
	}

	// Add to PATH so subsequent commands can find it
	binDir := filepath.Dir(engramBinaryPath)
	if err := system.AddToUserPath(binDir); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: could not add %s to PATH: %v\n", binDir, err)
	}

	return engramBinaryPath, nil
}

// applyMCP injects MCP (Model Context Protocol) configuration into each agent adapter.
// This includes Context7 server configuration for enhanced context capabilities.
func applyMCP(adapters []agents.Adapter, homeDir string) error {
	for _, adapter := range adapters {
		if _, err := mcp.Inject(homeDir, adapter); err != nil {
			return fmt.Errorf("inject mcp for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// applyPersona injects persona configuration (AGENTS.md system prompt) into each agent.
// Personas define the behavior and personality of the AI agent.
func applyPersona(adapters []agents.Adapter, homeDir string, personaID model.PersonaID) error {
	for _, adapter := range adapters {
		if _, err := persona.Inject(homeDir, adapter, personaID); err != nil {
			return fmt.Errorf("inject persona for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// applyPermissions injects permission configuration into each agent adapter.
// Permissions control which files and directories the agent can access.
func applyPermissions(adapters []agents.Adapter, homeDir string) error {
	for _, adapter := range adapters {
		if _, err := permissions.Inject(homeDir, adapter); err != nil {
			return fmt.Errorf("inject permissions for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// applySDD injects SDD (Spec-Driven Development) configuration into each agent adapter.
// SDD provides structured workflows for planning and implementing changes.
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

// applySkills injects agent skills into each adapter. Skills provide specialized
// capabilities like branch-pr workflow, issue creation, testing patterns, etc.
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

// applyDxrk installs the dxrk CLI tool (if needed) and injects its configuration
// into each agent. On Windows, it also ensures the PowerShell shim is in place.
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

// ensureDxrkInstalled checks if dxrk is available and installs it if not.
// Handles the case where the install script fails due to missing TTY but dxrk is actually installed.
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

// addDxrkToPath adds the dxrk bin directory to the user's PATH persistently.
// On Windows, this modifies the user registry via PowerShell.
func addDxrkToPath(homeDir string) error {
	dxrkBinDir := filepath.Join(homeDir, "bin")
	return system.AddToUserPath(dxrkBinDir)
}

// applyTheme injects visual theme configuration into each agent adapter.
// Themes define colors, fonts, and other UI preferences for the agent.
func applyTheme(adapters []agents.Adapter, homeDir string) error {
	for _, adapter := range adapters {
		if _, err := theme.Inject(homeDir, adapter); err != nil {
			return fmt.Errorf("inject theme for %q: %w", adapter.Agent(), err)
		}
	}
	return nil
}

// ComponentApplier handles applying a component to multiple agent adapters.
// It encapsulates the context needed for component installation and injection.
type ComponentApplier struct {
	adapters     []agents.Adapter
	agentIDs     []model.AgentID
	profile      system.PlatformProfile
	homeDir      string
	workspaceDir string
	selection    model.Selection
}

// NewComponentApplier creates a new ComponentApplier with the given context.
// Adapters are the agent adapters to inject into, agentIDs are the target agents,
// profile is the platform profile, homeDir is the user's home directory,
// workspaceDir is the current workspace directory, and selection is the user's choices.
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

// Apply injects the specified component into all agent adapters.
// It delegates to the appropriate applyXxx function based on the component type.
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
