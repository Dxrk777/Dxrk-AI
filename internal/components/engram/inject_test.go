package dxrk-memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Dxrk777/Dxrk-AI/internal/agents"
	"github.com/Dxrk777/Dxrk-AI/internal/agents/claude"
	"github.com/Dxrk777/Dxrk-AI/internal/agents/codex"
	"github.com/Dxrk777/Dxrk-AI/internal/agents/gemini"
	"github.com/Dxrk777/Dxrk-AI/internal/agents/opencode"
	"github.com/Dxrk777/Dxrk-AI/internal/agents/vscode"
)

func claudeAdapter() agents.Adapter   { return claude.NewAdapter() }
func opencodeAdapter() agents.Adapter { return opencode.NewAdapter() }
func codexAdapter() agents.Adapter    { return codex.NewAdapter() }
func geminiAdapter() agents.Adapter   { return gemini.NewAdapter() }

func skipOnWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test requires Unix paths (homebrew) and is not applicable on Windows")
	}
}

// assertArgsHaveToolsAgent is a shared helper that validates a JSON file
// contains the MCP "dxrk-memory" entry with --tools=agent in args.
func assertArgsHaveToolsAgent(t *testing.T, path string) {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	text := string(content)
	if !strings.Contains(text, `"--tools=agent"`) {
		t.Fatalf("file %q missing --tools=agent in args; got:\n%s", path, text)
	}
}

func TestInjectClaudeWritesMCPConfig(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	// Check MCP JSON file was created.
	mcpPath := filepath.Join(home, ".claude", "mcp", "dxrk-memory.json")
	mcpContent, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(dxrk-memory.json) error = %v", err)
	}

	// Parse the JSON and validate the "command" key exists and references dxrk-memory.
	// The command may be an absolute path (if dxrk-memory is on PATH) or the relative
	// string "dxrk-memory" (if not found). Both are valid.
	var parsed map[string]any
	if err := json.Unmarshal(mcpContent, &parsed); err != nil {
		t.Fatalf("Unmarshal(dxrk-memory.json) error = %v", err)
	}
	cmd, ok := parsed["command"].(string)
	if !ok || cmd == "" {
		t.Fatalf("dxrk-memory.json missing or empty command field; got: %s", mcpContent)
	}
	// Command must either be the literal "dxrk-memory" or an absolute path ending in "dxrk-memory".
	base := filepath.Base(cmd)
	if base != "dxrk-memory" && base != "dxrk-memory.exe" {
		t.Fatalf("dxrk-memory.json command %q does not reference dxrk-memory binary; got: %s", cmd, mcpContent)
	}
	if _, ok := parsed["args"]; !ok {
		t.Fatal("dxrk-memory.json missing args field")
	}
	// RED: must include --tools=agent
	assertArgsHaveToolsAgent(t, mcpPath)
}

func TestInjectClaudeWritesProtocolSection(t *testing.T) {
	home := t.TempDir()

	_, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	claudeMDPath := filepath.Join(home, ".claude", "CLAUDE.md")
	content, err := os.ReadFile(claudeMDPath)
	if err != nil {
		t.Fatalf("ReadFile(CLAUDE.md) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "<!-- Dxrk-AI:dxrk-memory-protocol -->") {
		t.Fatal("CLAUDE.md missing open marker for dxrk-memory-protocol")
	}
	if !strings.Contains(text, "<!-- /Dxrk-AI:dxrk-memory-protocol -->") {
		t.Fatal("CLAUDE.md missing close marker for dxrk-memory-protocol")
	}
	// Real content check.
	if !strings.Contains(text, "mem_save") {
		t.Fatal("CLAUDE.md missing real dxrk-memory protocol content (expected 'mem_save')")
	}
}

func TestInjectClaudeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}
}

func TestInjectOpenCodeMergesEngramToSettings(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false")
	}

	// Should include opencode.json and AGENTS.md (fallback protocol injection).
	if len(result.Files) != 2 {
		t.Fatalf("Inject() files = %v, want exactly 2 (opencode.json + AGENTS.md)", result.Files)
	}

	configPath := filepath.Join(home, ".config", "opencode", "opencode.json")
	config, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(opencode.json) error = %v", err)
	}

	text := string(config)
	if !strings.Contains(text, `"dxrk-memory"`) {
		t.Fatal("opencode.json missing dxrk-memory server entry")
	}
	if !strings.Contains(text, `"mcp"`) {
		t.Fatal("opencode.json missing mcp key")
	}
	if strings.Contains(text, `"mcpServers"`) {
		t.Fatal("opencode.json should use 'mcp' key, not 'mcpServers'")
	}
	if !strings.Contains(text, `"type": "local"`) {
		t.Fatal("opencode.json dxrk-memory missing type: local")
	}
	// OpenCode 1.3.3+: command must be an array, no separate "args" field.
	if !strings.Contains(text, `"--tools=agent"`) {
		t.Fatal("opencode.json missing --tools=agent in command array")
	}
	if strings.Contains(text, `"args"`) {
		t.Fatal("opencode.json must NOT have a separate args field — command must be an array")
	}

	// Verify NO plugin files or plugin arrays exist.
	pluginPath := filepath.Join(home, ".config", "opencode", "plugins", "dxrk-memory.ts")
	if _, err := os.Stat(pluginPath); err == nil {
		t.Fatal("plugin file should NOT exist — old approach removed")
	}
	if strings.Contains(text, `"plugins"`) {
		t.Fatal("opencode.json should NOT contain plugins key")
	}

	agentsPath := filepath.Join(home, ".config", "opencode", "AGENTS.md")
	agentsContent, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("ReadFile(AGENTS.md) error = %v", err)
	}
	agentsText := string(agentsContent)
	if !strings.Contains(agentsText, "<!-- Dxrk-AI:dxrk-memory-protocol -->") {
		t.Fatal("AGENTS.md missing dxrk-memory protocol section marker")
	}
	if !strings.Contains(agentsText, "mem_save") {
		t.Fatal("AGENTS.md missing dxrk-memory protocol content (expected 'mem_save')")
	}
}

func TestInjectOpenCodeIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject() first changed = false")
	}

	second, err := Inject(home, opencodeAdapter())
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true")
	}
}

// TestInjectOpenCodeMigratesFromOldFormat verifies that when a user's
// opencode.json contains the old v1.11.3 format (separate "args" key),
// Inject() replaces mcp.dxrk-memory atomically so that "args" is absent and
// "command" is an array — the format required by OpenCode 1.3.3+.
func TestInjectOpenCodeMigratesFromOldFormat(t *testing.T) {
	home := t.TempDir()

	mockEngramLookPath(t, "/opt/homebrew/bin/dxrk-memory", "")

	adapter := opencodeAdapter()
	configPath := adapter.SettingsPath(home)
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}

	// Pre-seed with the old v1.11.3 format.
	oldFormat := `{"mcp": {"dxrk-memory": {"command": "/opt/homebrew/bin/dxrk-memory", "args": ["mcp","--tools=agent"], "type": "local"}}}`
	if err := os.WriteFile(configPath, []byte(oldFormat), 0o644); err != nil {
		t.Fatalf("WriteFile(opencode.json) error = %v", err)
	}

	result, err := Inject(home, adapter)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject() changed = false; expected migration to produce a change")
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(opencode.json) error = %v", err)
	}

	// (1) "args" key must be absent from mcp.dxrk-memory.
	if strings.Contains(string(content), `"args"`) {
		t.Fatalf("mcp.dxrk-memory still contains 'args' key after migration; got:\n%s", content)
	}

	// (2) command must be a []any containing the dxrk-memory binary.
	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("Unmarshal(opencode.json) error = %v", err)
	}
	mcpMap, _ := parsed["mcp"].(map[string]any)
	engramMap, _ := mcpMap["dxrk-memory"].(map[string]any)
	cmdRaw, ok := engramMap["command"]
	if !ok {
		t.Fatalf("mcp.dxrk-memory missing command key; got:\n%s", content)
	}
	cmdArr, ok := cmdRaw.([]any)
	if !ok {
		t.Fatalf("mcp.dxrk-memory.command must be []any after migration, got %T; got:\n%s", cmdRaw, content)
	}
	if len(cmdArr) == 0 {
		t.Fatalf("mcp.dxrk-memory.command array is empty; got:\n%s", content)
	}
	firstElem, _ := cmdArr[0].(string)
	if firstElem == "" {
		t.Fatalf("mcp.dxrk-memory.command[0] is empty or not a string; got:\n%s", content)
	}
	// Must end with "dxrk-memory".
	if filepath.Base(firstElem) != "dxrk-memory" {
		t.Fatalf("mcp.dxrk-memory.command[0] = %q does not end with 'dxrk-memory'; got:\n%s", firstElem, content)
	}

	// (3) Second Inject() call must be idempotent (changed=false).
	second, err := Inject(home, adapter)
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true; expected idempotent (no change)")
	}
}

func TestInjectCursorMergesEngramToSettings(t *testing.T) {
	home := t.TempDir()

	cursorAdapter, err := agents.NewAdapter("cursor")
	if err != nil {
		t.Fatalf("NewAdapter(cursor) error = %v", err)
	}

	result, injectErr := Inject(home, cursorAdapter)
	if injectErr != nil {
		t.Fatalf("Inject(cursor) error = %v", injectErr)
	}

	// Cursor uses MCPConfigFile strategy — dxrk-memory gets merged into mcp.json.
	if !result.Changed {
		t.Fatalf("Inject(cursor) changed = false")
	}
}

func TestInjectCursorWithMalformedMCPJsonRecovery(t *testing.T) {
	// Real Windows users may have a ~/.cursor/mcp.json that starts with non-JSON
	// content (e.g. "allow: all" or just "a"). The installer should recover by
	// treating the broken file as {} and proceeding with the overlay merge.
	home := t.TempDir()

	cursorAdapter, err := agents.NewAdapter("cursor")
	if err != nil {
		t.Fatalf("NewAdapter(cursor) error = %v", err)
	}

	// Pre-create ~/.cursor/mcp.json with invalid (non-JSON) content.
	mcpPath := cursorAdapter.MCPConfigPath(home, "dxrk-memory")
	if err := os.MkdirAll(filepath.Dir(mcpPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	if err := os.WriteFile(mcpPath, []byte("allow: all"), 0o644); err != nil {
		t.Fatalf("WriteFile(malformed mcp.json) error = %v", err)
	}

	result, injectErr := Inject(home, cursorAdapter)
	if injectErr != nil {
		t.Fatalf("Inject(cursor) with malformed mcp.json error = %v; want nil (should recover)", injectErr)
	}
	if !result.Changed {
		t.Fatalf("Inject(cursor) changed = false; want true")
	}

	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(mcp.json) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, `"mcpServers"`) {
		t.Fatalf("mcp.json missing mcpServers key after recovery; got:\n%s", text)
	}
	if !strings.Contains(text, `"dxrk-memory"`) {
		t.Fatalf("mcp.json missing dxrk-memory server after recovery; got:\n%s", text)
	}
}

func TestInjectVSCodeMergesEngramToMCPConfigFile(t *testing.T) {
	home := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	adapter := vscode.NewAdapter()

	result, err := Inject(home, adapter)
	if err != nil {
		t.Fatalf("Inject(vscode) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(vscode) changed = false")
	}

	mcpPath := adapter.MCPConfigPath(home, "dxrk-memory")
	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(mcp.json) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, `"servers"`) {
		t.Fatal("mcp.json missing servers key")
	}
	if !strings.Contains(text, `"dxrk-memory"`) {
		t.Fatal("mcp.json missing dxrk-memory server")
	}
	if !strings.Contains(text, `"mcp"`) {
		t.Fatal("mcp.json missing dxrk-memory args mcp")
	}
	if strings.Contains(text, `"mcpServers"`) {
		t.Fatal("mcp.json should use 'servers' key, not 'mcpServers'")
	}
	// RED: VS Code overlay must include --tools=agent
	assertArgsHaveToolsAgent(t, mcpPath)
}

// ─── Gemini tests ─────────────────────────────────────────────────────────────

func TestInjectGeminiToolsFlagPresent(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, geminiAdapter())
	if err != nil {
		t.Fatalf("Inject(gemini) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(gemini) changed = false")
	}

	settingsPath := filepath.Join(home, ".gemini", "settings.json")
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile(settings.json) error = %v", err)
	}
	text := string(content)
	if !strings.Contains(text, `"mcpServers"`) {
		t.Fatal("settings.json missing mcpServers key")
	}
	if !strings.Contains(text, `"dxrk-memory"`) {
		t.Fatal("settings.json missing dxrk-memory entry")
	}
	// RED: Gemini overlay must use --tools=agent
	if !strings.Contains(text, `"--tools=agent"`) {
		t.Fatal("settings.json missing --tools=agent in args")
	}
}

// ─── Codex tests ──────────────────────────────────────────────────────────────

func TestInjectCodexWritesTOMLMCP(t *testing.T) {
	home := t.TempDir()

	result, err := Inject(home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(codex) changed = false")
	}

	configPath := filepath.Join(home, ".codex", "config.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(config.toml) error = %v", err)
	}
	text := string(content)
	if !strings.Contains(text, "[mcp_servers.dxrk-memory]") {
		t.Fatalf("config.toml missing [mcp_servers.dxrk-memory] block; got:\n%s", text)
	}
	// command must reference the dxrk-memory binary — either relative ("dxrk-memory") or an
	// absolute path (when dxrk-memory is on PATH). Both are valid.
	if !strings.Contains(text, "command = ") {
		t.Fatalf("config.toml missing command field; got:\n%s", text)
	}
	cmdLine := ""
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "command = ") {
			cmdLine = strings.TrimSpace(line)
			break
		}
	}
	if cmdLine == "" {
		t.Fatalf("config.toml missing command line; got:\n%s", text)
	}
	// The command value must end with "dxrk-memory" or "dxrk-memory.exe".
	cmdVal := strings.TrimPrefix(cmdLine, "command = ")
	cmdVal = strings.Trim(cmdVal, `"`)
	base := filepath.Base(cmdVal)
	if base != "dxrk-memory" && base != "dxrk-memory.exe" {
		t.Fatalf("config.toml command %q does not reference dxrk-memory binary; got:\n%s", cmdVal, text)
	}
	if !strings.Contains(text, `"--tools=agent"`) {
		t.Fatalf("config.toml missing --tools=agent; got:\n%s", text)
	}
}

func TestInjectCodexWritesInstructionFiles(t *testing.T) {
	home := t.TempDir()

	_, err := Inject(home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) error = %v", err)
	}

	instructionsPath := filepath.Join(home, ".codex", "dxrk-memory-instructions.md")
	content, err := os.ReadFile(instructionsPath)
	if err != nil {
		t.Fatalf("ReadFile(dxrk-memory-instructions.md) error = %v", err)
	}
	if !strings.Contains(string(content), "mem_save") {
		t.Fatal("dxrk-memory-instructions.md missing expected content (mem_save)")
	}

	compactPath := filepath.Join(home, ".codex", "dxrk-memory-compact-prompt.md")
	compactContent, err := os.ReadFile(compactPath)
	if err != nil {
		t.Fatalf("ReadFile(dxrk-memory-compact-prompt.md) error = %v", err)
	}
	if !strings.Contains(string(compactContent), "FIRST ACTION REQUIRED") {
		t.Fatal("dxrk-memory-compact-prompt.md missing expected content (FIRST ACTION REQUIRED)")
	}
}

func TestInjectCodexInjectsTOMLKeys(t *testing.T) {
	skipOnWindows(t)
	home := t.TempDir()

	_, err := Inject(home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) error = %v", err)
	}

	configPath := filepath.Join(home, ".codex", "config.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(config.toml) error = %v", err)
	}
	text := string(content)

	instructionsPath := filepath.Join(home, ".codex", "dxrk-memory-instructions.md")
	if !strings.Contains(text, `model_instructions_file`) {
		t.Fatalf("config.toml missing model_instructions_file key; got:\n%s", text)
	}
	if !strings.Contains(text, instructionsPath) {
		t.Fatalf("config.toml model_instructions_file does not reference %q; got:\n%s", instructionsPath, text)
	}

	compactPath := filepath.Join(home, ".codex", "dxrk-memory-compact-prompt.md")
	if !strings.Contains(text, `experimental_compact_prompt_file`) {
		t.Fatalf("config.toml missing experimental_compact_prompt_file key; got:\n%s", text)
	}
	if !strings.Contains(text, compactPath) {
		t.Fatalf("config.toml experimental_compact_prompt_file does not reference %q; got:\n%s", compactPath, text)
	}
}

// ─── Engram setup absolute path preservation tests ────────────────────────────

// TestInjectClaudePreservesAbsoluteCommandFromEngramSetup verifies that when
// `dxrk-memory setup claude-code` has already written an absolute-path command to
// ~/.claude", "mcp", "dxrk-memory.json (Engram v1.10.3+ behaviour), a subsequent call to
// Inject() does NOT overwrite the absolute path with the relative "dxrk-memory".
func TestInjectClaudePreservesAbsoluteCommandFromEngramSetup(t *testing.T) {
	skipOnWindows(t)
	home := t.TempDir()

	// Simulate what `dxrk-memory setup claude-code` writes on v1.10.3+:
	// an absolute path as the command value.
	absPath := "/opt/homebrew/bin/dxrk-memory"
	mcpPath := filepath.Join(home, ".claude", "mcp", "dxrk-memory.json")
	if err := os.MkdirAll(filepath.Dir(mcpPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	setupContent := []byte(`{
  "command": "/opt/homebrew/bin/dxrk-memory",
  "args": ["mcp", "--tools=agent"]
}
`)
	if err := os.WriteFile(mcpPath, setupContent, 0o644); err != nil {
		t.Fatalf("WriteFile(dxrk-memory.json) error = %v", err)
	}

	// Now run Inject — should NOT overwrite the absolute command.
	_, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(dxrk-memory.json) error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, absPath) {
		t.Fatalf("Inject() overwrote absolute command path; want %q preserved, got:\n%s", absPath, text)
	}
	// Still must have --tools=agent.
	assertArgsHaveToolsAgent(t, mcpPath)
}

// TestInjectClaudePreservesAbsoluteCommandIsIdempotent verifies that calling
// Inject() twice when an absolute-path dxrk-memory.json already exists does not
// cause repeated writes (idempotency).
func TestInjectClaudePreservesAbsoluteCommandIsIdempotent(t *testing.T) {
	skipOnWindows(t)
	home := t.TempDir()

	absPath := "/usr/local/bin/dxrk-memory"
	mcpPath := filepath.Join(home, ".claude", "mcp", "dxrk-memory.json")
	if err := os.MkdirAll(filepath.Dir(mcpPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	setupContent := []byte(`{
  "command": "/usr/local/bin/dxrk-memory",
  "args": ["mcp", "--tools=agent"]
}
`)
	if err := os.WriteFile(mcpPath, setupContent, 0o644); err != nil {
		t.Fatalf("WriteFile(dxrk-memory.json) error = %v", err)
	}

	first, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() first error = %v", err)
	}

	second, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject() second changed = true after absolute-path setup; want idempotent (no change)")
	}

	// Absolute path must still be present.
	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(dxrk-memory.json) error = %v", err)
	}
	if !strings.Contains(string(content), absPath) {
		t.Fatalf("absolute command path %q was lost after second Inject(); got:\n%s", absPath, string(content))
	}
	_ = first // first result not the focus of this test
}

// TestInjectClaudeAddsToolsAgentWhenSetupWritesBareArgs verifies that if
// `dxrk-memory setup` wrote an absolute command but with bare args (no --tools=agent),
// Inject() adds --tools=agent while preserving the absolute path.
func TestInjectClaudeAddsToolsAgentWhenSetupWritesBareArgs(t *testing.T) {
	skipOnWindows(t)
	home := t.TempDir()

	absPath := "/home/user/go/bin/dxrk-memory"
	mcpPath := filepath.Join(home, ".claude", "mcp", "dxrk-memory.json")
	if err := os.MkdirAll(filepath.Dir(mcpPath), 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	// Bare mcp arg without --tools=agent — older dxrk-memory setup format.
	setupContent := []byte(`{
  "command": "/home/user/go/bin/dxrk-memory",
  "args": ["mcp"]
}
`)
	if err := os.WriteFile(mcpPath, setupContent, 0o644); err != nil {
		t.Fatalf("WriteFile(dxrk-memory.json) error = %v", err)
	}

	_, err := Inject(home, claudeAdapter())
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	content, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("ReadFile(dxrk-memory.json) error = %v", err)
	}
	text := string(content)

	// Absolute path must be preserved.
	if !strings.Contains(text, absPath) {
		t.Fatalf("absolute path %q was lost; got:\n%s", absPath, text)
	}
	// --tools=agent must be added.
	assertArgsHaveToolsAgent(t, mcpPath)
}

func TestInjectCodexIsIdempotent(t *testing.T) {
	home := t.TempDir()

	first, err := Inject(home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) first error = %v", err)
	}
	if !first.Changed {
		t.Fatalf("Inject(codex) first changed = false")
	}

	second, err := Inject(home, codexAdapter())
	if err != nil {
		t.Fatalf("Inject(codex) second error = %v", err)
	}
	if second.Changed {
		t.Fatalf("Inject(codex) second changed = true (should be idempotent)")
	}

	// Verify only one [mcp_servers.dxrk-memory] block.
	configPath := filepath.Join(home, ".codex", "config.toml")
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("ReadFile(config.toml) error = %v", err)
	}
	count := strings.Count(string(content), "[mcp_servers.dxrk-memory]")
	if count != 1 {
		t.Fatalf("config.toml has %d [mcp_servers.dxrk-memory] blocks, want exactly 1; got:\n%s", count, string(content))
	}
}

// ─── Absolute path resolution tests ──────────────────────────────────────────

// mockEngramLookPath sets EngramLookPath to a mock and restores it after the test.
func mockEngramLookPath(t *testing.T, result string, errMsg string) {
	t.Helper()
	orig := EngramLookPath
	EngramLookPath = func(string) (string, error) {
		if errMsg != "" {
			return "", fmt.Errorf("%s", errMsg)
		}
		return result, nil
	}
	t.Cleanup(func() { EngramLookPath = orig })
}

// TestEngramInjectUsesAbsolutePathWhenAvailable verifies that when dxrk-memory is
// resolvable on PATH, its absolute path is written into the MCP config file
// for agents that use StrategyMCPConfigFile (e.g. Windsurf).
func TestEngramInjectUsesAbsolutePathWhenAvailable(t *testing.T) {
	home := t.TempDir()

	absPath := "/usr/local/bin/dxrk-memory"
	mockEngramLookPath(t, absPath, "")

	windsurfAdapter, err := agents.NewAdapter("windsurf")
	if err != nil {
		t.Fatalf("NewAdapter(windsurf) error = %v", err)
	}

	result, injectErr := Inject(home, windsurfAdapter)
	if injectErr != nil {
		t.Fatalf("Inject(windsurf) error = %v", injectErr)
	}
	if !result.Changed {
		t.Fatalf("Inject(windsurf) changed = false")
	}

	mcpPath := windsurfAdapter.MCPConfigPath(home, "dxrk-memory")
	content, readErr := os.ReadFile(mcpPath)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", mcpPath, readErr)
	}

	// Parse and validate the command field contains the absolute path.
	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("Unmarshal(%q) error = %v", mcpPath, err)
	}

	mcpServersRaw, ok := parsed["mcpServers"]
	if !ok {
		t.Fatalf("mcp_config.json missing mcpServers key; got:\n%s", content)
	}
	mcpServers, ok := mcpServersRaw.(map[string]any)
	if !ok {
		t.Fatalf("mcpServers has unexpected type: %T", mcpServersRaw)
	}
	engramServerRaw, ok := mcpServers["dxrk-memory"]
	if !ok {
		t.Fatalf("mcpServers missing dxrk-memory entry; got:\n%s", content)
	}
	engramServer, ok := engramServerRaw.(map[string]any)
	if !ok {
		t.Fatalf("dxrk-memory server has unexpected type: %T", engramServerRaw)
	}

	cmd, _ := engramServer["command"].(string)
	if cmd != absPath {
		t.Fatalf("mcp_config.json command = %q, want absolute path %q", cmd, absPath)
	}
}

// TestEngramInjectFallsBackToRelativeWhenNotFound verifies that when dxrk-memory
// cannot be resolved on PATH, the config falls back to the relative "dxrk-memory"
// command string.
func TestEngramInjectFallsBackToRelativeWhenNotFound(t *testing.T) {
	home := t.TempDir()

	mockEngramLookPath(t, "", "not found")

	windsurfAdapter, err := agents.NewAdapter("windsurf")
	if err != nil {
		t.Fatalf("NewAdapter(windsurf) error = %v", err)
	}

	result, injectErr := Inject(home, windsurfAdapter)
	if injectErr != nil {
		t.Fatalf("Inject(windsurf) error = %v", injectErr)
	}
	if !result.Changed {
		t.Fatalf("Inject(windsurf) changed = false")
	}

	mcpPath := windsurfAdapter.MCPConfigPath(home, "dxrk-memory")
	content, readErr := os.ReadFile(mcpPath)
	if readErr != nil {
		t.Fatalf("ReadFile(%q) error = %v", mcpPath, readErr)
	}

	text := string(content)
	if !strings.Contains(text, `"command": "dxrk-memory"`) {
		t.Fatalf("mcp_config.json should use relative fallback 'dxrk-memory'; got:\n%s", text)
	}
}

// TestEngramInjectAbsolutePathForOpenCodeMergeStrategy verifies that the
// absolute path is used when the StrategyMergeIntoSettings strategy is
// applied for OpenCode.
func TestEngramInjectAbsolutePathForOpenCodeMergeStrategy(t *testing.T) {
	home := t.TempDir()

	absPath := "/usr/local/bin/dxrk-memory"
	mockEngramLookPath(t, absPath, "")

	adapter := opencodeAdapter()
	settingsDir := filepath.Dir(adapter.SettingsPath(home))
	os.MkdirAll(settingsDir, 0o755)
	os.WriteFile(adapter.SettingsPath(home), []byte("{}"), 0o644)

	_, err := Inject(home, adapter)
	if err != nil {
		t.Fatalf("Inject() error = %v", err)
	}

	content, err := os.ReadFile(adapter.SettingsPath(home))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	text := string(content)
	if !strings.Contains(text, absPath) {
		t.Fatalf("OpenCode settings missing absolute dxrk-memory path, got: %s", text)
	}
	// OpenCode 1.3.3+: command must be an array, no separate "args" field.
	if strings.Contains(text, `"args"`) {
		t.Fatalf("OpenCode settings must NOT have a separate args field; got: %s", text)
	}

	// Structurally verify command is a []any containing the absolute path.
	var parsed map[string]any
	if err := json.Unmarshal(content, &parsed); err != nil {
		t.Fatalf("Unmarshal(opencode.json) error = %v", err)
	}
	mcpRaw, ok := parsed["mcp"]
	if !ok {
		t.Fatalf("opencode.json missing mcp key; got:\n%s", text)
	}
	mcpMap, ok := mcpRaw.(map[string]any)
	if !ok {
		t.Fatalf("mcp key has unexpected type %T; got:\n%s", mcpRaw, text)
	}
	engramRaw, ok := mcpMap["dxrk-memory"]
	if !ok {
		t.Fatalf("mcp missing dxrk-memory key; got:\n%s", text)
	}
	engramMap, ok := engramRaw.(map[string]any)
	if !ok {
		t.Fatalf("mcp.dxrk-memory has unexpected type %T; got:\n%s", engramRaw, text)
	}
	cmdRaw, ok := engramMap["command"]
	if !ok {
		t.Fatalf("mcp.dxrk-memory missing command key; got:\n%s", text)
	}
	cmdArr, ok := cmdRaw.([]any)
	if !ok {
		t.Fatalf("mcp.dxrk-memory.command must be an array, got %T; value:\n%s", cmdRaw, text)
	}
	if len(cmdArr) == 0 {
		t.Fatalf("mcp.dxrk-memory.command array is empty; got:\n%s", text)
	}
	firstElem, ok := cmdArr[0].(string)
	if !ok || firstElem != absPath {
		t.Fatalf("mcp.dxrk-memory.command[0] = %v, want %q; got:\n%s", cmdArr[0], absPath, text)
	}
}

// TestEngramInjectAbsolutePathForGeminiMergeStrategy verifies that the
// absolute path is also used when the StrategyMergeIntoSettings strategy is
// applied (e.g. Gemini CLI).
func TestEngramInjectAbsolutePathForGeminiMergeStrategy(t *testing.T) {
	home := t.TempDir()

	absPath := "/opt/homebrew/bin/dxrk-memory"
	mockEngramLookPath(t, absPath, "")

	result, err := Inject(home, geminiAdapter())
	if err != nil {
		t.Fatalf("Inject(gemini) error = %v", err)
	}
	if !result.Changed {
		t.Fatalf("Inject(gemini) changed = false")
	}

	settingsPath := filepath.Join(home, ".gemini", "settings.json")
	content, readErr := os.ReadFile(settingsPath)
	if readErr != nil {
		t.Fatalf("ReadFile(settings.json) error = %v", readErr)
	}

	text := string(content)
	if !strings.Contains(text, absPath) {
		t.Fatalf("settings.json missing absolute path %q; got:\n%s", absPath, text)
	}
}
