package update

import (
	"fmt"
	"strings"
	"testing"
)

func TestRenderCLI_IncompleteCheckDoesNotClaimUpToDate(t *testing.T) {
	results := []UpdateResult{
		{Tool: ToolInfo{Name: "dxrk"}, InstalledVersion: "1.0.0", LatestVersion: "1.0.0", Status: UpToDate},
		{Tool: ToolInfo{Name: "dxrk-memory"}, Status: CheckFailed, Err: fmt.Errorf("timeout")},
	}

	out := RenderCLI(results)

	if strings.Contains(out, "All tools are up to date!") {
		t.Fatalf("RenderCLI must not claim all tools are up to date when checks fail:\n%s", out)
	}
	if !strings.Contains(out, "Update check incomplete") {
		t.Fatalf("RenderCLI must mention incomplete checks:\n%s", out)
	}
	if !strings.Contains(out, "check failed") {
		t.Fatalf("RenderCLI must surface failed rows:\n%s", out)
	}
}

func TestCheckFailures(t *testing.T) {
	results := []UpdateResult{
		{Tool: ToolInfo{Name: "dxrk"}, Status: UpToDate},
		{Tool: ToolInfo{Name: "dxrk-memory"}, Status: CheckFailed},
		{Tool: ToolInfo{Name: "dxrk-guardian"}, Status: CheckFailed},
	}

	failed := CheckFailures(results)
	if len(failed) != 2 {
		t.Fatalf("len(CheckFailures) = %d, want 2", len(failed))
	}
	if failed[0] != "dxrk-memory" || failed[1] != "dxrk-guardian" {
		t.Fatalf("CheckFailures() = %v, want [dxrk-memory dxrk-guardian]", failed)
	}
	if !HasCheckFailures(results) {
		t.Fatalf("HasCheckFailures() = false, want true")
	}
}
