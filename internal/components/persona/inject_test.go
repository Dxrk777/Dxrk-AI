package persona

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

<<<<<<< HEAD
	"github.com/Dxrk777/Dxrk-Hex/internal/agents"
	"github.com/Dxrk777/Dxrk-Hex/internal/agents/claude"
	"github.com/Dxrk777/Dxrk-Hex/internal/agents/opencode"
	"github.com/Dxrk777/Dxrk-Hex/internal/assets"
	"github.com/Dxrk777/Dxrk-Hex/internal/model"
=======
	"github.com/gentleman-programming/gentle-ai/internal/agents"
	"github.com/gentleman-programming/gentle-ai/internal/agents/claude"
	"github.com/gentleman-programming/gentle-ai/internal/agents/opencode"
	"github.com/gentleman-programming/gentle-ai/internal/assets"
	"github.com/gentleman-programming/gentle-ai/internal/model"
>>>>>>> upstream/main
)

func claudeAdapter() agents.Adapter   { return claude.NewAdapter() }
func opencodeAdapter() agents.Adapter { return opencode.NewAdapter() }

func TestInjectClaudeDxrkWritesSectionWithRealContent(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter(), model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	path := filepath.Join(home, ".claude", "CLAUDE.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "<!-- dxrk:persona -->") {
		t.Fatal("CLAUDE.md missing open marker for persona")
	}
	if !strings.Contains(text, "<!-- /dxrk:persona -->") {
		t.Fatal("CLAUDE.md missing close marker for persona")
	}
	// Real content check — the embedded persona has these patterns.
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("CLAUDE.md missing real persona content (expected 'Senior Architect')")
	}
}

func TestInjectClaudeDxrkWritesOutputStyleFile(t *testing.T) {
	home := t.TempDir()

	_, err := Inject(home, claudeAdapter(), model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	// Verify output-style file was written.
	stylePath := filepath.Join(home, ".claude", "output-styles", "Dxrk.md")
	content, err := os.ReadFile(stylePath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", stylePath, err)
	}

	text := string(content)
	if !strings.Contains(text, "name: Dxrk") {
		t.Fatal("Output style file missing YAML frontmatter 'name: Dxrk'")
	}
	if !strings.Contains(text, "keep-coding-instructions: true") {
		t.Fatal("Output style file missing 'keep-coding-instructions: true'")
	}
	if !strings.Contains(text, "Dxrk Output Style") {
		t.Fatal("Output style file missing 'Dxrk Output Style' heading")
	}
}

func TestInjectClaudeDxrkMergesOutputStyleIntoSettings(t *testing.T) {
	home := t.TempDir()

	// Pre-create a settings.json with some existing content.
	settingsDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(settingsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	existingSettings := `{"permissions": {"allow": ["Read"]}, "syntaxHighlightingDisabled": true}`
	if err := os.WriteFile(filepath.Join(settingsDir, "settings.json"), []byte(existingSettings), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := Inject(home, claudeAdapter(), model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	// Verify settings.json has outputStyle merged in.
	settingsPath := filepath.Join(home, ".claude", "settings.json")
	settingsContent, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", settingsPath, err)
	}

	var settings map[string]any
	if err := json.Unmarshal(settingsContent, &settings); err != nil {
		t.Fatalf("Unmarshal settings.json error = %v", err)
	}

	outputStyle, ok := settings["outputStyle"]
	if !ok {
		t.Fatal("settings.json missing 'outputStyle' key")
	}
	if outputStyle != "Dxrk" {
		t.Fatalf("settings.json outputStyle = %q, want %q", outputStyle, "Dxrk")
	}

	// Verify existing keys were preserved.
	if _, ok := settings["permissions"]; !ok {
		t.Fatal("settings.json lost 'permissions' key during merge")
	}
	if _, ok := settings["syntaxHighlightingDisabled"]; !ok {
		t.Fatal("settings.json lost 'syntaxHighlightingDisabled' key during merge")
	}
}

func TestInjectClaudeDxrkReturnsAllFiles(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter(), model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	// Should return 3 files: CLAUDE.md, output-style, settings.json.
	if len(result.Files) != 3 {
		t.Fatalf("Inject() returned %d files, want 3: %v", len(result.Files), result.Files)
	}

	wantSuffixes := []string{"CLAUDE.md", "Dxrk.md", "settings.json"}
	for _, suffix := range wantSuffixes {
		found := false
		for _, f := range result.Files {
			if strings.HasSuffix(f, suffix) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Inject() missing file with suffix %q in %v", suffix, result.Files)
		}
	}
}

func TestInjectClaudeNeutralWritesFullPersonaWithoutRegionalLanguage(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter(), model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	path := filepath.Join(home, ".claude", "CLAUDE.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	text := string(content)
	// Neutral persona is the same teacher — should have Senior Architect.
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("Neutral persona should contain 'Senior Architect'")
	}
<<<<<<< HEAD
	// Should NOT have Dxrk-specific regional language.
=======
	// Should NOT have gentleman-specific regional language.
>>>>>>> upstream/main
	if strings.Contains(text, "Rioplatense") {
		t.Fatal("Neutral persona should not contain Rioplatense language")
	}
}

func TestInjectClaudeNeutralDoesNotWriteOutputStyle(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter(), model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	// Should only return CLAUDE.md, no output-style file.
	if len(result.Files) != 1 {
		t.Fatalf("Neutral persona returned %d files, want 1: %v", len(result.Files), result.Files)
	}

	// Output-style file should NOT exist.
	stylePath := filepath.Join(home, ".claude", "output-styles", "Dxrk.md")
	if _, err := os.Stat(stylePath); !os.IsNotExist(err) {
		t.Fatal("Neutral persona should NOT write output-style file")
	}
}

func TestInjectCustomClaudeDoesNothing(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter(), model.PersonaCustom)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if result.Changed {
		t.Fatal("Custom persona should NOT change anything")
	}
	if len(result.Files) != 0 {
		t.Fatalf("Custom persona should return no files, got %v", result.Files)
	}

	// CLAUDE.md should NOT be created.
	claudeMD := filepath.Join(home, ".claude", "CLAUDE.md")
	if _, err := os.Stat(claudeMD); !os.IsNotExist(err) {
		t.Fatal("Custom persona should NOT create CLAUDE.md")
	}
}

func TestInjectCustomOpenCodeDoesNothing(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, opencodeAdapter(), model.PersonaCustom)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if result.Changed {
		t.Fatal("Custom persona (OpenCode) should NOT change anything")
	}
	if len(result.Files) != 0 {
		t.Fatalf("Custom persona (OpenCode) should return no files, got %v", result.Files)
	}

	// AGENTS.md should NOT be created.
	agentsMD := filepath.Join(home, ".config", "opencode", "AGENTS.md")
	if _, err := os.Stat(agentsMD); !os.IsNotExist(err) {
		t.Fatal("Custom persona (OpenCode) should NOT create AGENTS.md")
	}
}

func TestInjectOpenCodeDxrkWritesAgentsFile(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, opencodeAdapter(), model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	path := filepath.Join(home, ".config", "opencode", "AGENTS.md")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("AGENTS.md missing real persona content")
	}
}

func TestInjectOpenCodeNeutralPreservesManagedSections(t *testing.T) {
	home := t.TempDir()

<<<<<<< HEAD
	// First install Dxrk persona + simulate SDD/engram sections
	_, err := Inject(home, opencodeAdapter(), model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject(Dxrk) error = %v", err)
=======
	// First install gentleman persona + simulate SDD/engram sections
	_, err := Inject(home, opencodeAdapter(), model.PersonaGentleman)
	if err != nil {
		t.Fatalf("Inject(gentleman) error = %v", err)
>>>>>>> upstream/main
	}

	path := filepath.Join(home, ".config", "opencode", "AGENTS.md")

	// Simulate SDD and engram sections appended by sdd.Inject and engram.Inject
	existing, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
<<<<<<< HEAD
	withSections := string(existing) + "\n\n<!-- dxrk:sdd-orchestrator -->\nSDD orchestrator content here\n<!-- /dxrk:sdd-orchestrator -->\n\n<!-- dxrk:engram-protocol -->\nEngram protocol content here\n<!-- /dxrk:engram-protocol -->\n"
=======
	withSections := string(existing) + "\n\n<!-- gentle-ai:sdd-orchestrator -->\nSDD orchestrator content here\n<!-- /gentle-ai:sdd-orchestrator -->\n\n<!-- gentle-ai:engram-protocol -->\nEngram protocol content here\n<!-- /gentle-ai:engram-protocol -->\n"
>>>>>>> upstream/main
	if err := os.WriteFile(path, []byte(withSections), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Now switch to neutral persona
	result, err := Inject(home, opencodeAdapter(), model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject(neutral) error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject(neutral) should report changed")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() after neutral error = %v", err)
	}
	text := string(content)

	// Neutral content should be present
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("AGENTS.md missing neutral persona content")
	}
	if strings.Contains(text, "Rioplatense") {
		t.Fatal("AGENTS.md has Rioplatense language in neutral persona — should be neutral tone")
	}

	// Managed sections MUST be preserved
<<<<<<< HEAD
	if !strings.Contains(text, "<!-- dxrk:sdd-orchestrator -->") {
		t.Fatal("AGENTS.md lost SDD orchestrator section after switching to neutral persona")
	}
	if !strings.Contains(text, "<!-- dxrk:engram-protocol -->") {
		t.Fatal("AGENTS.md lost engram protocol section after switching to neutral persona")
	}

	// Dxrk-specific language should be gone — neutral has the same personality but no regional language
=======
	if !strings.Contains(text, "<!-- gentle-ai:sdd-orchestrator -->") {
		t.Fatal("AGENTS.md lost SDD orchestrator section after switching to neutral persona")
	}
	if !strings.Contains(text, "<!-- gentle-ai:engram-protocol -->") {
		t.Fatal("AGENTS.md lost engram protocol section after switching to neutral persona")
	}

	// Gentleman-specific language should be gone — neutral has the same personality but no regional language
>>>>>>> upstream/main
	if strings.Contains(text, "Rioplatense") {
		t.Fatal("AGENTS.md still has Rioplatense language after switching to neutral")
	}
}

func TestInjectVSCodeNeutralPreservesManagedSections(t *testing.T) {
	home := t.TempDir()

	vscodeAdapter, err := agents.NewAdapter("vscode-copilot")
	if err != nil {
		t.Fatalf("NewAdapter(vscode-copilot) error = %v", err)
	}

<<<<<<< HEAD
	_, err = Inject(home, vscodeAdapter, model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject(Dxrk) error = %v", err)
=======
	_, err = Inject(home, vscodeAdapter, model.PersonaGentleman)
	if err != nil {
		t.Fatalf("Inject(gentleman) error = %v", err)
>>>>>>> upstream/main
	}

	path := vscodeAdapter.SystemPromptFile(home)

	existing, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
<<<<<<< HEAD
	withSections := string(existing) + "\n\n<!-- dxrk:sdd-orchestrator -->\nSDD content\n<!-- /dxrk:sdd-orchestrator -->\n"
=======
	withSections := string(existing) + "\n\n<!-- gentle-ai:sdd-orchestrator -->\nSDD content\n<!-- /gentle-ai:sdd-orchestrator -->\n"
>>>>>>> upstream/main
	if err := os.WriteFile(path, []byte(withSections), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err = Inject(home, vscodeAdapter, model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject(neutral) error = %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() after neutral error = %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("instructions file missing neutral persona content")
	}
	if strings.Contains(text, "Rioplatense") {
		t.Fatal("instructions file has Rioplatense language in neutral persona")
	}
<<<<<<< HEAD
	if !strings.Contains(text, "<!-- dxrk:sdd-orchestrator -->") {
=======
	if !strings.Contains(text, "<!-- gentle-ai:sdd-orchestrator -->") {
>>>>>>> upstream/main
		t.Fatal("instructions file lost SDD section after switching to neutral persona")
	}
	if !strings.Contains(text, "---\nname:") {
		t.Fatal("instructions file lost YAML frontmatter")
	}
}

func TestInjectNeutralPreservesWhenMarkerAtByteZero(t *testing.T) {
	home := t.TempDir()

	opencodeAdapter, err := agents.NewAdapter("opencode")
	if err != nil {
		t.Fatalf("NewAdapter(opencode) error = %v", err)
	}

	promptPath := opencodeAdapter.SystemPromptFile(home)
	if err := os.MkdirAll(filepath.Dir(promptPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	// File starts DIRECTLY with a managed marker at byte 0 — no persona preamble.
<<<<<<< HEAD
	markerOnly := "<!-- dxrk:sdd-orchestrator -->\nSDD content\n<!-- /dxrk:sdd-orchestrator -->\n"
=======
	markerOnly := "<!-- gentle-ai:sdd-orchestrator -->\nSDD content\n<!-- /gentle-ai:sdd-orchestrator -->\n"
>>>>>>> upstream/main
	if err := os.WriteFile(promptPath, []byte(markerOnly), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err = Inject(home, opencodeAdapter, model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject(neutral) error = %v", err)
	}

	content, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("missing neutral persona content")
	}
<<<<<<< HEAD
	if !strings.Contains(text, "<!-- dxrk:sdd-orchestrator -->") {
=======
	if !strings.Contains(text, "<!-- gentle-ai:sdd-orchestrator -->") {
>>>>>>> upstream/main
		t.Fatal("SDD section destroyed when marker was at byte 0")
	}
}

func TestInjectNeutralIdempotentWithManagedSections(t *testing.T) {
	home := t.TempDir()

	opencodeAdapter, err := agents.NewAdapter("opencode")
	if err != nil {
		t.Fatalf("NewAdapter(opencode) error = %v", err)
	}

	promptPath := opencodeAdapter.SystemPromptFile(home)
	if err := os.MkdirAll(filepath.Dir(promptPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	// Set up: neutral + managed sections
	// Simulate a file with neutral persona + managed sections.
	// Use a fingerprint from the real neutral asset so the test is realistic.
	neutralContent := assets.MustRead("generic/persona-neutral.md")
<<<<<<< HEAD
	initial := neutralContent + "\n\n<!-- dxrk:sdd-orchestrator -->\nSDD content\n<!-- /dxrk:sdd-orchestrator -->\n\n<!-- dxrk:engram-protocol -->\nEngram content\n<!-- /dxrk:engram-protocol -->\n"
=======
	initial := neutralContent + "\n\n<!-- gentle-ai:sdd-orchestrator -->\nSDD content\n<!-- /gentle-ai:sdd-orchestrator -->\n\n<!-- gentle-ai:engram-protocol -->\nEngram content\n<!-- /gentle-ai:engram-protocol -->\n"
>>>>>>> upstream/main
	if err := os.WriteFile(promptPath, []byte(initial), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// First neutral inject
	result1, err := Inject(home, opencodeAdapter, model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject(neutral) first error = %v", err)
	}

	// Second neutral inject — should be idempotent
	result2, err := Inject(home, opencodeAdapter, model.PersonaNeutral)
	if err != nil {
		t.Fatalf("Inject(neutral) second error = %v", err)
	}

	if result2.Changed && !result1.Changed {
		t.Fatal("second neutral inject should not report changed when first didn't")
	}

	content, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(content)

	// Verify no duplication
<<<<<<< HEAD
	if strings.Count(text, "<!-- dxrk:sdd-orchestrator -->") != 1 {
=======
	if strings.Count(text, "<!-- gentle-ai:sdd-orchestrator -->") != 1 {
>>>>>>> upstream/main
		t.Fatal("SDD section duplicated after idempotent neutral inject")
	}
	if strings.Count(text, "## Rules") != 1 {
		t.Fatal("neutral persona duplicated after idempotent inject")
	}
<<<<<<< HEAD
	if strings.Count(text, "<!-- dxrk:engram-protocol -->") != 1 {
=======
	if strings.Count(text, "<!-- gentle-ai:engram-protocol -->") != 1 {
>>>>>>> upstream/main
		t.Fatal("engram section duplicated after idempotent neutral inject")
	}
}

func TestInjectClaudeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, claudeAdapter(), model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, claudeAdapter(), model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}
}

func TestInjectOpenCodeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, opencodeAdapter(), model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, opencodeAdapter(), model.PersonaDxrk)
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}
}

func TestInjectWindsurfIsIdempotent(t *testing.T) {
	home := t.TempDir()

	windsurfAdapter, err := agents.NewAdapter("windsurf")
	if err != nil {
		t.Fatalf("NewAdapter(windsurf) error = %v", err)
	}

<<<<<<< HEAD
	first, err := Inject(home, windsurfAdapter, model.PersonaDxrk)
=======
	first, err := Inject(home, windsurfAdapter, model.PersonaGentleman)
>>>>>>> upstream/main
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	promptPath := windsurfAdapter.SystemPromptFile(home)
	contentAfterFirst, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("ReadFile() after first inject error = %v", err)
	}

<<<<<<< HEAD
	second, err := Inject(home, windsurfAdapter, model.PersonaDxrk)
=======
	second, err := Inject(home, windsurfAdapter, model.PersonaGentleman)
>>>>>>> upstream/main
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true — persona was duplicated in global_rules.md")
	}

	contentAfterSecond, err := os.ReadFile(promptPath)
	if err != nil {
		t.Fatalf("ReadFile() after second inject error = %v", err)
	}

	if string(contentAfterFirst) != string(contentAfterSecond) {
		t.Fatal("global_rules.md content changed on second inject — persona was duplicated")
	}
}

<<<<<<< HEAD
func TestInjectCursorDxrkWritesRulesFileWithRealContent(t *testing.T) {
=======
func TestInjectCursorGentlemanWritesRulesFileWithRealContent(t *testing.T) {
>>>>>>> upstream/main
	home := t.TempDir()

	cursorAdapter, err := agents.NewAdapter("cursor")
	if err != nil {
		t.Fatalf("NewAdapter(cursor) error = %v", err)
	}

	result, injectErr := Inject(home, cursorAdapter, model.PersonaDxrk)
	if injectErr != nil {
		t.Fatalf("Inject(cursor) error = %v", injectErr)
	}

	if !result.Changed {
		t.Fatalf("Inject(cursor, Dxrk) changed = false")
	}

	// Verify the generic persona content was used — not just neutral one-liner.
	path := filepath.Join(home, ".cursor", "rules", "dxrk.mdc")
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, readErr)
	}

	text := string(content)
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("Cursor persona missing 'Senior Architect' — got neutral fallback instead of generic persona")
	}
	if !strings.Contains(text, "Skills") {
		t.Fatal("Cursor persona missing skills section")
	}
}

func TestInjectGeminiDxrkWritesSystemPromptWithRealContent(t *testing.T) {
	home := t.TempDir()

	geminiAdapter, err := agents.NewAdapter("gemini-cli")
	if err != nil {
		t.Fatalf("NewAdapter(gemini-cli) error = %v", err)
	}

	result, injectErr := Inject(home, geminiAdapter, model.PersonaDxrk)
	if injectErr != nil {
		t.Fatalf("Inject(gemini) error = %v", injectErr)
	}

	if !result.Changed {
		t.Fatal("Inject(gemini, Dxrk) changed = false")
	}

	path := filepath.Join(home, ".gemini", "GEMINI.md")
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, readErr)
	}

	text := string(content)
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("Gemini persona missing 'Senior Architect'")
	}
}

func TestInjectVSCodeDxrkWritesInstructionsFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	vscodeAdapter, err := agents.NewAdapter("vscode-copilot")
	if err != nil {
		t.Fatalf("NewAdapter(vscode-copilot) error = %v", err)
	}

	result, injectErr := Inject(home, vscodeAdapter, model.PersonaDxrk)
	if injectErr != nil {
		t.Fatalf("Inject(vscode) error = %v", injectErr)
	}

	if !result.Changed {
		t.Fatal("Inject(vscode, Dxrk) changed = false")
	}

	path := vscodeAdapter.SystemPromptFile(home)
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, readErr)
	}

	text := string(content)
	if !strings.Contains(text, "applyTo: \"**\"") {
		t.Fatal("VS Code instructions file missing YAML frontmatter applyTo pattern")
	}
	if !strings.Contains(text, "Senior Architect") {
		t.Fatal("VS Code persona missing 'Senior Architect'")
	}
}

// --- Auto-heal tests: Claude Code stale free-text persona ---

<<<<<<< HEAD
// legacyClaudePersonaBlock simulates a Dxrk persona block that was written
=======
// legacyClaudePersonaBlock simulates a Gentleman persona block that was written
>>>>>>> upstream/main
// directly (without markers) by an old installer or manually by the user.
const legacyClaudePersonaBlock = `## Rules

- NEVER add "Co-Authored-By" or any AI attribution to commits. Use conventional commits format only.
- Never build after changes.

## Personality

Senior Architect, 15+ years experience, GDE & MVP.

## Language

- Spanish input → Rioplatense Spanish.

## Behavior

- Push back when user asks for code without context.

`

func TestInjectClaudeAutoHealsStaleFreeTextPersona(t *testing.T) {
	home := t.TempDir()

	// Pre-populate CLAUDE.md with legacy persona content (no markers) followed
	// by a properly-marked section from a previous installer run.
	claudeMD := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.MkdirAll(filepath.Dir(claudeMD), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// Simulate a stale install: free-text persona block at top, then a different
	// marked section below (e.g., from a previous SDD install).
<<<<<<< HEAD
	stalePreamble := legacyClaudePersonaBlock + "\n<!-- dxrk:sdd -->\nOld SDD content.\n<!-- /dxrk:sdd -->\n"
=======
	stalePreamble := legacyClaudePersonaBlock + "\n<!-- gentle-ai:sdd -->\nOld SDD content.\n<!-- /gentle-ai:sdd -->\n"
>>>>>>> upstream/main
	if err := os.WriteFile(claudeMD, []byte(stalePreamble), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

<<<<<<< HEAD
	result, err := Inject(home, claudeAdapter(), model.PersonaDxrk)
=======
	result, err := Inject(home, claudeAdapter(), model.PersonaGentleman)
>>>>>>> upstream/main
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject() should have changed the file to remove the legacy block")
	}

	content, err := os.ReadFile(claudeMD)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(content)

	// The file should now have the persona inside markers, not as free text.
<<<<<<< HEAD
	if !strings.Contains(text, "<!-- dxrk:persona -->") {
		t.Fatal("CLAUDE.md missing persona marker after heal")
	}
	if !strings.Contains(text, "<!-- /dxrk:persona -->") {
=======
	if !strings.Contains(text, "<!-- gentle-ai:persona -->") {
		t.Fatal("CLAUDE.md missing persona marker after heal")
	}
	if !strings.Contains(text, "<!-- /gentle-ai:persona -->") {
>>>>>>> upstream/main
		t.Fatal("CLAUDE.md missing persona close marker after heal")
	}

	// The existing SDD section must be preserved.
<<<<<<< HEAD
	if !strings.Contains(text, "<!-- dxrk:sdd -->") {
=======
	if !strings.Contains(text, "<!-- gentle-ai:sdd -->") {
>>>>>>> upstream/main
		t.Fatal("CLAUDE.md lost the sdd section during heal")
	}
	if !strings.Contains(text, "Old SDD content.") {
		t.Fatal("CLAUDE.md lost the sdd section content during heal")
	}

	// The persona content must NOT appear twice (no duplicate blocks).
	firstPersonaIdx := strings.Index(text, "Senior Architect")
	if firstPersonaIdx < 0 {
		t.Fatal("CLAUDE.md missing 'Senior Architect' persona content")
	}
	// Verify there's no second occurrence outside the markers.
	lastPersonaIdx := strings.LastIndex(text, "Senior Architect")
	if firstPersonaIdx != lastPersonaIdx {
		// It's OK if the same string appears inside the single persona marker block
		// multiple times (e.g., content + newlines), but there must not be a
		// separate free-text block also containing it.
		// Check: everything before the open marker should NOT contain "Senior Architect".
<<<<<<< HEAD
		openMarkerIdx := strings.Index(text, "<!-- dxrk:persona -->")
=======
		openMarkerIdx := strings.Index(text, "<!-- gentle-ai:persona -->")
>>>>>>> upstream/main
		if openMarkerIdx >= 0 && strings.Contains(text[:openMarkerIdx], "Senior Architect") {
			t.Fatal("CLAUDE.md still has 'Senior Architect' before the persona marker — legacy block not fully stripped")
		}
	}
}

func TestInjectClaudeAutoHealStalePersonaOnlyFile(t *testing.T) {
	home := t.TempDir()

	// CLAUDE.md contains ONLY the legacy persona block (no markers at all).
	claudeMD := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.MkdirAll(filepath.Dir(claudeMD), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	if err := os.WriteFile(claudeMD, []byte(legacyClaudePersonaBlock), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

<<<<<<< HEAD
	result, err := Inject(home, claudeAdapter(), model.PersonaDxrk)
=======
	result, err := Inject(home, claudeAdapter(), model.PersonaGentleman)
>>>>>>> upstream/main
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject() should have changed the file")
	}

	content, err := os.ReadFile(claudeMD)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(content)

	// Must have markers now.
<<<<<<< HEAD
	if !strings.Contains(text, "<!-- dxrk:persona -->") {
=======
	if !strings.Contains(text, "<!-- gentle-ai:persona -->") {
>>>>>>> upstream/main
		t.Fatal("CLAUDE.md missing persona marker")
	}

	// Must NOT have the legacy free-text block before markers.
<<<<<<< HEAD
	openMarkerIdx := strings.Index(text, "<!-- dxrk:persona -->")
=======
	openMarkerIdx := strings.Index(text, "<!-- gentle-ai:persona -->")
>>>>>>> upstream/main
	if openMarkerIdx >= 0 {
		before := text[:openMarkerIdx]
		if strings.Contains(before, "## Rules") {
			t.Fatal("legacy '## Rules' block still present before persona marker")
		}
	}
}

func TestInjectClaudeHealDoesNotTouchNonPersonaContent(t *testing.T) {
	home := t.TempDir()

	// CLAUDE.md has user content that does NOT match persona fingerprints.
	claudeMD := filepath.Join(home, ".claude", "CLAUDE.md")
	if err := os.MkdirAll(filepath.Dir(claudeMD), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	userContent := "# My custom config\n\nI like turtles.\n"
	if err := os.WriteFile(claudeMD, []byte(userContent), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

<<<<<<< HEAD
	result, err := Inject(home, claudeAdapter(), model.PersonaDxrk)
=======
	result, err := Inject(home, claudeAdapter(), model.PersonaGentleman)
>>>>>>> upstream/main
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatal("Inject() should write persona section")
	}

	content, err := os.ReadFile(claudeMD)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	text := string(content)

	// User content must be preserved.
	if !strings.Contains(text, "I like turtles.") {
		t.Fatal("user content was erased — heal was too aggressive")
	}
	// Persona section must be appended.
<<<<<<< HEAD
	if !strings.Contains(text, "<!-- dxrk:persona -->") {
=======
	if !strings.Contains(text, "<!-- gentle-ai:persona -->") {
>>>>>>> upstream/main
		t.Fatal("persona section not appended")
	}
}

// --- Auto-heal tests: VSCode stale legacy path cleanup ---

func TestInjectVSCodeCleansLegacyGitHubPersonaFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

<<<<<<< HEAD
	// Plant an old-style Dxrk persona file at the legacy path.
=======
	// Plant an old-style Gentleman persona file at the legacy path.
>>>>>>> upstream/main
	legacyPath := filepath.Join(home, ".github", "copilot-instructions.md")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	// Old installer wrote raw persona content without YAML frontmatter.
	oldContent := "## Personality\n\nSenior Architect, 15+ years experience.\n"
	if err := os.WriteFile(legacyPath, []byte(oldContent), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	vscodeAdapter, err := agents.NewAdapter("vscode-copilot")
	if err != nil {
		t.Fatalf("NewAdapter(vscode-copilot) error = %v", err)
	}

<<<<<<< HEAD
	result, injectErr := Inject(home, vscodeAdapter, model.PersonaDxrk)
=======
	result, injectErr := Inject(home, vscodeAdapter, model.PersonaGentleman)
>>>>>>> upstream/main
	if injectErr != nil {
		t.Fatalf("Inject(vscode) error = %v", injectErr)
	}
	if !result.Changed {
		t.Fatal("Inject(vscode) should report changed (legacy cleanup + new file write)")
	}

	// Legacy file must be gone.
	if _, statErr := os.Stat(legacyPath); !os.IsNotExist(statErr) {
		t.Fatal("legacy ~/.github/copilot-instructions.md was NOT removed by auto-heal")
	}

	// New file must exist at the current path.
	newPath := vscodeAdapter.SystemPromptFile(home)
	content, readErr := os.ReadFile(newPath)
	if readErr != nil {
		t.Fatalf("ReadFile new path %q error = %v", newPath, readErr)
	}
	if !strings.Contains(string(content), "applyTo: \"**\"") {
		t.Fatal("new VSCode instructions file missing YAML frontmatter")
	}
}

func TestInjectVSCodePreservesNonPersonaGitHubFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	// Plant a .github/copilot-instructions.md that has user content (not a
<<<<<<< HEAD
	// Dxrk persona) — it must NOT be deleted.
=======
	// Gentleman persona) — it must NOT be deleted.
>>>>>>> upstream/main
	legacyPath := filepath.Join(home, ".github", "copilot-instructions.md")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	userContent := "# My custom Copilot instructions\n\nAlways be concise.\n"
	if err := os.WriteFile(legacyPath, []byte(userContent), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	vscodeAdapter, err := agents.NewAdapter("vscode-copilot")
	if err != nil {
		t.Fatalf("NewAdapter(vscode-copilot) error = %v", err)
	}

<<<<<<< HEAD
	_, injectErr := Inject(home, vscodeAdapter, model.PersonaDxrk)
=======
	_, injectErr := Inject(home, vscodeAdapter, model.PersonaGentleman)
>>>>>>> upstream/main
	if injectErr != nil {
		t.Fatalf("Inject(vscode) error = %v", injectErr)
	}

	// User's file must still exist.
	remaining, readErr := os.ReadFile(legacyPath)
	if readErr != nil {
		t.Fatalf("legacy user file was deleted: ReadFile error = %v", readErr)
	}
	if string(remaining) != userContent {
		t.Fatalf("user file content was modified: got %q", string(remaining))
	}
}

<<<<<<< HEAD
func TestNeutralAndDxrkToneSectionsMatch(t *testing.T) {
	neutral := assets.MustRead("generic/persona-neutral.md")
	Dxrk := assets.MustRead("generic/persona-Dxrk.md")
=======
func TestNeutralAndGentlemanToneSectionsMatch(t *testing.T) {
	neutral := assets.MustRead("generic/persona-neutral.md")
	gentleman := assets.MustRead("generic/persona-gentleman.md")
>>>>>>> upstream/main

	extractSection := func(content, section string) string {
		idx := strings.Index(content, "## "+section)
		if idx < 0 {
			return ""
		}
		rest := content[idx:]
		nextIdx := strings.Index(rest[1:], "\n## ")
		if nextIdx < 0 {
			return rest
		}
		return rest[:nextIdx+1]
	}

	neutralTone := extractSection(neutral, "Tone")
<<<<<<< HEAD
	DxrkTone := extractSection(Dxrk, "Tone")

	if neutralTone != DxrkTone {
		t.Fatalf("## Tone sections diverged:\nneutral:\n%s\nDxrk:\n%s", neutralTone, DxrkTone)
=======
	gentlemanTone := extractSection(gentleman, "Tone")

	if neutralTone != gentlemanTone {
		t.Fatalf("## Tone sections diverged:\nneutral:\n%s\ngentleman:\n%s", neutralTone, gentlemanTone)
>>>>>>> upstream/main
	}
}

func TestInjectVSCodeIdempotentAfterHeal(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))

	// Plant legacy file and run inject twice — second run should be idempotent.
	legacyPath := filepath.Join(home, ".github", "copilot-instructions.md")
	if err := os.MkdirAll(filepath.Dir(legacyPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	if err := os.WriteFile(legacyPath, []byte("## Personality\n\nSenior Architect.\n"), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	vscodeAdapter, err := agents.NewAdapter("vscode-copilot")
	if err != nil {
		t.Fatalf("NewAdapter(vscode-copilot) error = %v", err)
	}

<<<<<<< HEAD
	first, err := Inject(home, vscodeAdapter, model.PersonaDxrk)
=======
	first, err := Inject(home, vscodeAdapter, model.PersonaGentleman)
>>>>>>> upstream/main
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatal("first inject should have changed")
	}

<<<<<<< HEAD
	second, err := Inject(home, vscodeAdapter, model.PersonaDxrk)
=======
	second, err := Inject(home, vscodeAdapter, model.PersonaGentleman)
>>>>>>> upstream/main
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("second inject should be idempotent (changed = false), but changed = true")
	}
}
