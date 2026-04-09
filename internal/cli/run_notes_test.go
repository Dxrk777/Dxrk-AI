package cli

import (
	"runtime"
	"strings"
	"testing"

	"github.com/Dxrk777/Dxrk-AI/internal/model"
	"github.com/Dxrk777/Dxrk-AI/internal/planner"
	"github.com/Dxrk777/Dxrk-AI/internal/verify"
)

func TestWithPostInstallNotesAddsDxrkNextSteps(t *testing.T) {
	report := verify.Report{Ready: true, FinalNote: "You're ready."}
	resolved := planner.ResolvedPlan{OrderedComponents: []model.ComponentID{model.ComponentDxrk}}

	updated := withPostInstallNotes(report, resolved)
	if !strings.Contains(updated.FinalNote, "Dxrk is now installed globally") {
		t.Fatalf("FinalNote missing Dxrk global install note: %q", updated.FinalNote)
	}
	if !strings.Contains(updated.FinalNote, "dxrk init") || !strings.Contains(updated.FinalNote, "dxrk install") {
		t.Fatalf("FinalNote missing Dxrk repo setup steps: %q", updated.FinalNote)
	}
}

func TestWithPostInstallNotesDoesNotChangeNonDxrk(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("GOBIN path not applicable on Windows")
	}
	// Set GOBIN to a directory already in PATH so that withGoInstallPathNote
	// does not append a PATH guidance note for the Engram component.
	t.Setenv("GOBIN", "/usr/local/bin")

	report := verify.Report{Ready: true, FinalNote: "You're ready."}
	resolved := planner.ResolvedPlan{OrderedComponents: []model.ComponentID{model.ComponentEngram}}

	updated := withPostInstallNotes(report, resolved)
	if updated.FinalNote != report.FinalNote {
		t.Fatalf("FinalNote changed unexpectedly: %q", updated.FinalNote)
	}
}
