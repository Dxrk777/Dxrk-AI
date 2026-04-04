package screens

import (
	"testing"
)

func TestBrainMenuOptions(t *testing.T) {
	options := BrainMenuOptions()

	// Verify we have exactly 7 options
	if len(options) != 7 {
		t.Errorf("BrainMenuOptions() returned %d options, want 7", len(options))
	}

	// Verify the expected options exist
	expected := []string{
		"💬 Ask anything",
		"💻 Execute command",
		"📧 Send email",
		"📊 System status",
		"📜 View history",
		"⚙️  Configure",
		"🔙 Back",
	}

	for i, exp := range expected {
		if i < len(options) && options[i] != exp {
			t.Errorf("BrainMenuOptions()[%d] = %q, want %q", i, options[i], exp)
		}
	}
}

func TestBrainMenuOptionCount(t *testing.T) {
	count := BrainMenuOptionCount()
	expected := 7

	if count != expected {
		t.Errorf("BrainMenuOptionCount() = %d, want %d", count, expected)
	}
}

func TestNewBrainState(t *testing.T) {
	state := NewBrainState()

	if state.Mode != "menu" {
		t.Errorf("NewBrainState().Mode = %q, want %q", state.Mode, "menu")
	}

	if state.Input != "" {
		t.Errorf("NewBrainState().Input = %q, want empty", state.Input)
	}

	if state.Cursor != 0 {
		t.Errorf("NewBrainState().Cursor = %d, want 0", state.Cursor)
	}

	if state.Waiting != false {
		t.Errorf("NewBrainState().Waiting = %v, want false", state.Waiting)
	}
}

func TestRenderBrain_MenuMode(t *testing.T) {
	state := NewBrainState()
	state.Mode = "menu"
	state.Cursor = 0

	output := RenderBrain(state, 0)

	// Verify the output contains expected text
	expectedTexts := []string{
		"Dxrk Hex Brain",
		"Ask anything",
		"Execute command",
		"navigate",
	}

	for _, text := range expectedTexts {
		if !contains(output, text) {
			t.Errorf("RenderBrain() output should contain %q", text)
		}
	}
}

func TestRenderBrain_ChatMode(t *testing.T) {
	state := NewBrainState()
	state.Mode = "chat"
	state.Input = "hello"
	state.Cursor = 0

	output := RenderBrain(state, 0)

	// Verify the output contains expected text
	expectedTexts := []string{
		"Ask Dxrk Hex Brain",
		"hello",
	}

	for _, text := range expectedTexts {
		if !contains(output, text) {
			t.Errorf("RenderBrain() output should contain %q", text)
		}
	}
}

func TestRenderBrain_ExecuteMode(t *testing.T) {
	state := NewBrainState()
	state.Mode = "execute"
	state.Input = "ls -la"
	state.Cursor = 0

	output := RenderBrain(state, 0)

	// Verify the output contains expected text
	expectedTexts := []string{
		"Execute Shell Command",
		"ls -la",
	}

	for _, text := range expectedTexts {
		if !contains(output, text) {
			t.Errorf("RenderBrain() output should contain %q", text)
		}
	}
}

func TestRenderBrain_EmailMode(t *testing.T) {
	state := NewBrainState()
	state.Mode = "email"
	state.Input = "to user@test.com subject Hello"
	state.Cursor = 0

	output := RenderBrain(state, 0)

	// Verify the output contains expected text
	expectedTexts := []string{
		"Send Email",
		"to user@test.com",
	}

	for _, text := range expectedTexts {
		if !contains(output, text) {
			t.Errorf("RenderBrain() output should contain %q", text)
		}
	}
}

func TestRenderBrain_StatusMode(t *testing.T) {
	state := NewBrainState()
	state.Mode = "status"
	state.Output = "System Status: Running"
	state.Cursor = 0

	output := RenderBrain(state, 0)

	// Verify the output contains expected text
	if !contains(output, "System Status") {
		t.Errorf("RenderBrain() output should contain %q", "System Status")
	}

	if !contains(output, "System Status: Running") {
		t.Errorf("RenderBrain() output should contain %q", "System Status: Running")
	}
}

func TestRenderBrain_HistoryMode(t *testing.T) {
	state := NewBrainState()
	state.Mode = "history"
	state.Output = "Command History\n\ncmd1\ncmd2"
	state.Cursor = 0

	output := RenderBrain(state, 0)

	// Verify the output contains expected text
	if !contains(output, "Command History") {
		t.Errorf("RenderBrain() output should contain %q", "Command History")
	}
}

func TestRenderBrain_ConfigureMode(t *testing.T) {
	state := NewBrainState()
	state.Mode = "configure"
	state.Output = "Configuration options here"
	state.Cursor = 0

	output := RenderBrain(state, 0)

	// Verify the output contains expected text
	if !contains(output, "Configure") {
		t.Errorf("RenderBrain() output should contain %q", "Configure")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
