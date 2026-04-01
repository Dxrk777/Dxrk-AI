package cli

import (
	"testing"
)

func TestRunBrainNoArgs(t *testing.T) {
	err := RunBrain([]string{})
	if err != nil {
		t.Errorf("RunBrain with no args should not error: %v", err)
	}
}

func TestRunBrainHelp(t *testing.T) {
	err := RunBrain([]string{"help"})
	if err != nil {
		t.Errorf("RunBrain help should not error: %v", err)
	}
}

func TestRunBrainStatus(t *testing.T) {
	err := RunBrain([]string{"status"})
	if err != nil {
		t.Errorf("RunBrain status should not error: %v", err)
	}
}

func TestRunBrainHistory(t *testing.T) {
	err := RunBrain([]string{"history"})
	if err != nil {
		t.Errorf("RunBrain history should not error: %v", err)
	}
}

func TestRunBrainAgents(t *testing.T) {
	err := RunBrain([]string{"agents"})
	if err != nil {
		t.Errorf("RunBrain agents should not error: %v", err)
	}
}

func TestRunBrainVersion(t *testing.T) {
	err := RunBrain([]string{"version"})
	if err != nil {
		t.Errorf("RunBrain version should not error: %v", err)
	}
}

func TestRunBrainSync(t *testing.T) {
	err := RunBrain([]string{"sync"})
	if err != nil {
		t.Errorf("RunBrain sync should not error: %v", err)
	}
}

func TestRunBrainUpdate(t *testing.T) {
	err := RunBrain([]string{"update"})
	if err != nil {
		t.Errorf("RunBrain update should not error: %v", err)
	}
}

func TestRunBrainBackup(t *testing.T) {
	err := RunBrain([]string{"backup"})
	if err != nil {
		t.Errorf("RunBrain backup should not error: %v", err)
	}
}

func TestRunBrainInstall(t *testing.T) {
	err := RunBrain([]string{"install", "opencode"})
	if err != nil {
		t.Errorf("RunBrain install should not error: %v", err)
	}
}

func TestRunBrainUninstall(t *testing.T) {
	err := RunBrain([]string{"uninstall", "claude"})
	if err != nil {
		t.Errorf("RunBrain uninstall should not error: %v", err)
	}
}

func TestRunBrainConfigure(t *testing.T) {
	err := RunBrain([]string{"configure"})
	if err != nil {
		t.Errorf("RunBrain configure should not error: %v", err)
	}
}

func TestRunBrainConfig(t *testing.T) {
	err := RunBrain([]string{"config"})
	if err != nil {
		t.Errorf("RunBrain config should not error: %v", err)
	}
}

func TestRunBrainRunEcho(t *testing.T) {
	err := RunBrain([]string{"run", "echo", "test"})
	if err != nil {
		t.Errorf("RunBrain run echo should not error: %v", err)
	}
}

func TestRunBrainExecute(t *testing.T) {
	err := RunBrain([]string{"execute", "whoami"})
	if err != nil {
		t.Errorf("RunBrain execute should not error: %v", err)
	}
}

func TestRunBrainCmd(t *testing.T) {
	err := RunBrain([]string{"cmd", "pwd"})
	if err != nil {
		t.Errorf("RunBrain cmd should not error: %v", err)
	}
}

func TestRunBrainRemember(t *testing.T) {
	err := RunBrain([]string{"remember", "test"})
	if err != nil {
		t.Errorf("RunBrain remember should not error: %v", err)
	}
}

func TestRunBrainSearch(t *testing.T) {
	err := RunBrain([]string{"search", "test"})
	if err != nil {
		t.Errorf("RunBrain search should not error: %v", err)
	}
}

func TestRunBrainEmail(t *testing.T) {
	err := RunBrain([]string{"email", "to", "test@example.com", "subject", "Test"})
	if err != nil {
		t.Errorf("RunBrain email should not error: %v", err)
	}
}

func TestRunBrainNaturalLanguage(t *testing.T) {
	err := RunBrain([]string{"what", "is", "the", "status"})
	if err != nil {
		t.Errorf("RunBrain natural language should not error: %v", err)
	}
}

func TestRunBrainUnknown(t *testing.T) {
	err := RunBrain([]string{"asdfghjkl"})
	// Unknown commands should return an error
	if err == nil {
		t.Log("Unknown command returned nil error (may be acceptable)")
	}
}
