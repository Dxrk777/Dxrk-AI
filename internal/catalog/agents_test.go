package catalog

import (
	"testing"

	"github.com/Dxrk777/Dxrk-Hex/internal/model"
)

func TestAllAgents(t *testing.T) {
	agents := AllAgents()

	if len(agents) == 0 {
		t.Error("AllAgents should return at least one agent")
	}

	// Check all agents have required fields
	for _, agent := range agents {
		if agent.ID == "" {
			t.Error("Agent ID should not be empty")
		}
		if agent.Name == "" {
			t.Error("Agent Name should not be empty")
		}
		if agent.Tier != model.TierFull {
			t.Errorf("Agent %s should have TierFull, got %v", agent.Name, agent.Tier)
		}
	}
}

func TestMVPAgents(t *testing.T) {
	agents := MVPAgents()

	if len(agents) != 2 {
		t.Errorf("MVPAgents should return 2 agents, got %d", len(agents))
	}

	// MVP should be Claude Code and OpenCode
	if agents[0].ID != model.AgentClaudeCode {
		t.Errorf("First MVP agent should be ClaudeCode, got %v", agents[0].ID)
	}
	if agents[1].ID != model.AgentOpenCode {
		t.Errorf("Second MVP agent should be OpenCode, got %v", agents[1].ID)
	}
}

func TestIsMVPAgent(t *testing.T) {
	if !IsMVPAgent(model.AgentClaudeCode) {
		t.Error("ClaudeCode should be MVP agent")
	}
	if !IsMVPAgent(model.AgentOpenCode) {
		t.Error("OpenCode should be MVP agent")
	}
	if IsMVPAgent(model.AgentCursor) {
		t.Error("Cursor should NOT be MVP agent")
	}
}

func TestIsSupportedAgent(t *testing.T) {
	// All known agents should be supported
	supported := []model.AgentID{
		model.AgentClaudeCode,
		model.AgentOpenCode,
		model.AgentGeminiCLI,
		model.AgentCodex,
		model.AgentCursor,
		model.AgentVSCodeCopilot,
		model.AgentAntigravity,
		model.AgentWindsurf,
	}

	for _, id := range supported {
		if !IsSupportedAgent(id) {
			t.Errorf("Agent %v should be supported", id)
		}
	}

	// Unknown agents should not be supported
	if IsSupportedAgent(model.AgentID("unknown-agent")) {
		t.Error("Unknown agent should not be supported")
	}
}

func TestAgentCopiesAreIndependent(t *testing.T) {
	// Modifying returned agents should not affect the original
	agents1 := AllAgents()
	agents2 := AllAgents()

	if len(agents1) != len(agents2) {
		t.Error("Copies should have same length")
	}

	agents1[0].Name = "Modified Name"
	if agents2[0].Name == "Modified Name" {
		t.Error("Modifying copy should not affect original")
	}
}

func TestSkillStruct(t *testing.T) {
	skill := Skill{
		ID:       model.SkillSDDInit,
		Name:     "sdd-init",
		Category: "sdd",
		Priority: "p0",
	}

	if skill.ID != model.SkillSDDInit {
		t.Errorf("Skill ID mismatch: got %v, want %v", skill.ID, model.SkillSDDInit)
	}
	if skill.Name != "sdd-init" {
		t.Errorf("Skill Name mismatch: got %v, want sdd-init", skill.Name)
	}
}

func TestAgentStruct(t *testing.T) {
	agent := Agent{
		ID:         model.AgentClaudeCode,
		Name:       "Claude Code",
		Tier:       model.TierFull,
		ConfigPath: "~/.claude",
	}

	if agent.ID != model.AgentClaudeCode {
		t.Errorf("Agent ID mismatch: got %v, want %v", agent.ID, model.AgentClaudeCode)
	}
	if agent.ConfigPath != "~/.claude" {
		t.Errorf("Agent ConfigPath mismatch: got %v, want ~/.claude", agent.ConfigPath)
	}
}
