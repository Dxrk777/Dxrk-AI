package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/Dxrk777/Dxrk-AI/internal/components/engram"
	"github.com/Dxrk777/Dxrk-AI/internal/model"
	"github.com/Dxrk777/Dxrk-AI/internal/planner"
	"github.com/Dxrk777/Dxrk-AI/internal/verify"
)

// runPostApplyVerification runs all post-installation verification checks.
func runPostApplyVerification(homeDir string, selection model.Selection, resolved planner.ResolvedPlan) verify.Report {
	checks := make([]verify.Check, 0)
	adapters := resolveAdapters(resolved.Agents)

	for _, component := range resolved.OrderedComponents {
		for _, path := range componentPaths(homeDir, selection, adapters, component) {
			currentPath := path
			checks = append(checks, verify.Check{
				ID:          "verify:file:" + currentPath,
				Description: "required file exists",
				Run: func(context.Context) error {
					if _, err := os.Stat(currentPath); err != nil {
						return err
					}
					return nil
				},
			})
		}
	}

	if hasComponent(resolved.OrderedComponents, model.ComponentEngram) {
		checks = append(checks, engramHealthChecks()...)
	}
	checks = append(checks, antigravityCollisionCheck(resolved.Agents)...)

	return verify.BuildReport(verify.RunChecks(context.Background(), checks))
}

// engramHealthChecks returns health checks for dxrk-memory binary.
func engramHealthChecks() []verify.Check {
	return []verify.Check{
		{
			ID:          "verify:engram:binary",
			Description: "dxrk-memory binary on PATH (restart shell if missing)",
			Soft:        true,
			Run: func(context.Context) error {
				if err := engram.VerifyInstalled(); err != nil {
					return fmt.Errorf("%w\nIf engram was installed via `go install`, add it to PATH:\n  %s", err, engramPathGuidance(os.Getenv("SHELL")))
				}
				return nil
			},
		},
		{
			ID:          "verify:engram:version",
			Description: "engram version returns valid output",
			Soft:        true,
			Run: func(context.Context) error {
				if err := engram.VerifyInstalled(); err != nil {
					// Binary not on PATH — skip version check gracefully.
					return nil
				}
				_, err := engram.VerifyVersion()
				return err
			},
		},
	}
}

// antigravityCollisionCheck returns a soft verify check that warns the user
// when both Antigravity and Gemini CLI are selected. Both agents write to
// ~/.gemini/GEMINI.md — content is merged (not overwritten) but the user
// should be aware.
func antigravityCollisionCheck(agents []model.AgentID) []verify.Check {
	hasAntigravity := false
	hasGemini := false
	for _, id := range agents {
		if id == model.AgentAntigravity {
			hasAntigravity = true
		}
		if id == model.AgentGeminiCLI {
			hasGemini = true
		}
	}
	if !hasAntigravity || !hasGemini {
		return nil
	}
	return []verify.Check{
		{
			ID:          "verify:antigravity:rules-collision",
			Description: "Antigravity and Gemini CLI share ~/.gemini/GEMINI.md",
			Soft:        true,
			Run: func(context.Context) error {
				return fmt.Errorf(
					"both Antigravity and Gemini CLI write rules to ~/.gemini/GEMINI.md\n" +
						"Content is merged, not overwritten — rules from both agents coexist in the same file.\n" +
						"This is expected behavior. No action required unless you want to separate them manually.",
				)
			},
		},
	}
}
