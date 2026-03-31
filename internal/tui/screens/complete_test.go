package screens

import (
	"strings"
	"testing"
)

<<<<<<< HEAD
func TestRenderCompleteSuccessShowsDxrkNotesWhenInstalled(t *testing.T) {
	out := RenderComplete(CompletePayload{
		ConfiguredAgents:    1,
		InstalledComponents: 1,
		DxrkInstalled:        true,
	})

	if !strings.Contains(out, "Dxrk (per project)") {
		t.Fatalf("missing Dxrk section: %q", out)
	}
	if !strings.Contains(out, "dxrk init") || !strings.Contains(out, "dxrk install") {
		t.Fatalf("missing Dxrk repo commands: %q", out)
	}
}

func TestRenderCompleteSuccessHidesDxrkNotesWhenNotInstalled(t *testing.T) {
	out := RenderComplete(CompletePayload{
		ConfiguredAgents:    1,
		InstalledComponents: 1,
		DxrkInstalled:        false,
	})

	if strings.Contains(out, "Dxrk (per project)") {
		t.Fatalf("unexpected Dxrk section: %q", out)
=======
func TestRenderCompleteSuccessShowsGGANotesWhenInstalled(t *testing.T) {
	out := RenderComplete(CompletePayload{
		ConfiguredAgents:    1,
		InstalledComponents: 1,
		GGAInstalled:        true,
	})

	if !strings.Contains(out, "GGA (per project)") {
		t.Fatalf("missing GGA section: %q", out)
	}
	if !strings.Contains(out, "gga init") || !strings.Contains(out, "gga install") {
		t.Fatalf("missing GGA repo commands: %q", out)
	}
}

func TestRenderCompleteSuccessHidesGGANotesWhenNotInstalled(t *testing.T) {
	out := RenderComplete(CompletePayload{
		ConfiguredAgents:    1,
		InstalledComponents: 1,
		GGAInstalled:        false,
	})

	if strings.Contains(out, "GGA (per project)") {
		t.Fatalf("unexpected GGA section: %q", out)
>>>>>>> upstream/main
	}
}
