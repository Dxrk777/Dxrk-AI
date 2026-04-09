package skills

import (
	"io/fs"
	"strings"

	"github.com/Dxrk777/Dxrk-AI/internal/assets"
	"github.com/Dxrk777/Dxrk-AI/internal/model"
)

// sddSkills are the SDD orchestrator skills — always included.
// These are also managed by the SDD component, so the skills component
// skips them to prevent duplicate writes.
var sddSkills = []model.SkillID{
	model.SkillSDDInit,
	model.SkillSDDExplore,
	model.SkillSDDPropose,
	model.SkillSDDSpec,
	model.SkillSDDDesign,
	model.SkillSDDTasks,
	model.SkillSDDApply,
	model.SkillSDDVerify,
	model.SkillSDDArchive,
	model.SkillSDDOnboard,
	model.SkillJudgmentDay,
}

// foundationSkills are baseline skills for the "recommended" tier.
// FIX: SkillSkillRegistry was previously missing from this list.
var foundationSkills = []model.SkillID{
	model.SkillGoTesting,
	model.SkillCreator,
	model.SkillBranchPR,
	model.SkillIssueCreation,
	model.SkillSkillRegistry, // was missing — now included
}

// SkillsForPreset returns which skills should be installed for a given preset.
//
//   - "minimal"        → SDD skills only
//   - "ecosystem-only" → SDD + foundation skills
//   - "full-dxrk"      → ALL skills discovered from embedded assets (dynamic)
//   - "custom"         → empty (caller provides explicit list)
func SkillsForPreset(preset model.PresetID) []model.SkillID {
	switch preset {
	case model.PresetMinimal:
		return copySkills(sddSkills)
	case model.PresetEcosystemOnly:
		return copySkills(append(sddSkills, foundationSkills...))
	case model.PresetFullDxrk:
		// FIX: previously returned a hardcoded list that was identical to
		// ecosystem-only, ignoring all other embedded skills. Now we
		// discover every skill present in the embedded asset FS so that
		// new skills added to the repo are automatically activated.
		return AllSkillIDs()
	case model.PresetCustom:
		return nil
	default:
		// Unknown preset — default to full.
		return AllSkillIDs()
	}
}

// AllSkillIDs returns every skill ID discovered from the embedded assets FS.
//
// FIX: the previous implementation returned a hardcoded list that missed
// SkillSkillRegistry and any future skills added to the repo. This version
// walks assets.FS/skills and builds the list dynamically so all embedded
// skills are activated automatically without requiring code changes.
func AllSkillIDs() []model.SkillID {
	var ids []model.SkillID

	entries, err := fs.ReadDir(assets.FS, "skills")
	if err != nil {
		// Fallback to static list if FS walk fails (should never happen).
		return fallbackAllSkillIDs()
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		// Skip the shared helpers directory.
		if name == "_shared" || strings.HasPrefix(name, "_") {
			continue
		}
		ids = append(ids, model.SkillID(name))
	}

	if len(ids) == 0 {
		return fallbackAllSkillIDs()
	}

	return ids
}

// fallbackAllSkillIDs returns the statically known skill list.
// Used only when the embedded FS walk fails unexpectedly.
func fallbackAllSkillIDs() []model.SkillID {
	all := make([]model.SkillID, 0, len(sddSkills)+len(foundationSkills))
	all = append(all, sddSkills...)
	all = append(all, foundationSkills...)
	return all
}

func copySkills(src []model.SkillID) []model.SkillID {
	dst := make([]model.SkillID, len(src))
	copy(dst, src)
	return dst
}
