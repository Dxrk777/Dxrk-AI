package model_test

import (
	"testing"

	"github.com/Dxrk777/Dxrk/internal/model"
)

func TestClaudeModelAliasValid(t *testing.T) {
	tests := []struct {
		name     string
		alias    model.ClaudeModelAlias
		expected bool
	}{
		{"Opus is valid", model.ClaudeModelOpus, true},
		{"Sonnet is valid", model.ClaudeModelSonnet, true},
		{"Haiku is valid", model.ClaudeModelHaiku, true},
		{"Empty is invalid", model.ClaudeModelAlias(""), false},
		{"Random is invalid", model.ClaudeModelAlias("random"), false},
		{"Case sensitive opus", model.ClaudeModelAlias("Opus"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.alias.Valid(); got != tt.expected {
				t.Errorf("Valid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestClaudeModelAliasString(t *testing.T) {
	if model.ClaudeModelOpus.String() != "opus" {
		t.Errorf("Opus.String() = %v, want opus", model.ClaudeModelOpus.String())
	}
	if model.ClaudeModelSonnet.String() != "sonnet" {
		t.Errorf("Sonnet.String() = %v, want sonnet", model.ClaudeModelSonnet.String())
	}
	if model.ClaudeModelHaiku.String() != "haiku" {
		t.Errorf("Haiku.String() = %v, want haiku", model.ClaudeModelHaiku.String())
	}
}

func TestClaudeModelPresetBalanced(t *testing.T) {
	preset := model.ClaudeModelPresetBalanced()

	if preset["orchestrator"] != model.ClaudeModelOpus {
		t.Error("orchestrator should be Opus in balanced preset")
	}
	if preset["sdd-archive"] != model.ClaudeModelHaiku {
		t.Error("sdd-archive should be Haiku in balanced preset")
	}
	if preset["default"] != model.ClaudeModelSonnet {
		t.Error("default should be Sonnet in balanced preset")
	}
	if len(preset) != 10 {
		t.Errorf("balanced preset should have 10 entries, got %d", len(preset))
	}
}

func TestClaudeModelPresetPerformance(t *testing.T) {
	preset := model.ClaudeModelPresetPerformance()

	if preset["sdd-design"] != model.ClaudeModelOpus {
		t.Error("sdd-design should be Opus in performance preset")
	}
	if preset["sdd-verify"] != model.ClaudeModelOpus {
		t.Error("sdd-verify should be Opus in performance preset")
	}
	if preset["sdd-archive"] != model.ClaudeModelHaiku {
		t.Error("sdd-archive should be Haiku in performance preset")
	}
}

func TestClaudeModelPresetEconomy(t *testing.T) {
	preset := model.ClaudeModelPresetEconomy()

	// Economy should use Sonnet for everything except archive
	for phase, alias := range preset {
		if phase == "sdd-archive" {
			if alias != model.ClaudeModelHaiku {
				t.Errorf("sdd-archive should be Haiku, got %v", alias)
			}
		} else if alias != model.ClaudeModelSonnet {
			t.Errorf("phase %s should be Sonnet in economy preset, got %v", phase, alias)
		}
	}
}

func TestSelectionHasAgent(t *testing.T) {
	sel := model.Selection{
		Agents: []model.AgentID{model.AgentClaudeCode, model.AgentOpenCode},
	}

	if !sel.HasAgent(model.AgentClaudeCode) {
		t.Error("Selection should have AgentClaudeCode")
	}
	if !sel.HasAgent(model.AgentOpenCode) {
		t.Error("Selection should have AgentOpenCode")
	}
	if sel.HasAgent(model.AgentCursor) {
		t.Error("Selection should NOT have AgentCursor")
	}
}

func TestSelectionHasComponent(t *testing.T) {
	sel := model.Selection{
		Components: []model.ComponentID{model.ComponentEngram, model.ComponentSDD},
	}

	if !sel.HasComponent(model.ComponentEngram) {
		t.Error("Selection should have ComponentEngram")
	}
	if sel.HasComponent(model.ComponentPersona) {
		t.Error("Selection should NOT have ComponentPersona")
	}
}

func TestSyncOverrides(t *testing.T) {
	// Test nil means no override
	var overrides *model.SyncOverrides
	if overrides != nil {
		t.Error("SyncOverrides should be nil by default")
	}

	// Test with overrides
	strictTDD := true
	overrides = &model.SyncOverrides{
		StrictTDD: &strictTDD,
	}

	if overrides.StrictTDD == nil || !*overrides.StrictTDD {
		t.Error("StrictTDD should be true")
	}
}
