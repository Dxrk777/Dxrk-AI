---
name: git-advanced
description: >
  Advanced Git workflows and patterns. Trigger: When doing complex git operations, rebasing, bisecting, or managing branches.
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use
- Interactive rebase and squashing commits
- Bisecting to find bugs
- Cherry-picking specific commits
- Managing complex branch workflows
- Rewriting history safely

## Critical Patterns

### Interactive rebase (REQUIRED for clean history)
```bash
# Rebase last 5 commits
git rebase -i HEAD~5
# Pick, squash, fixup, reword, drop

# Autosquash (for fixup! commits)
git commit --fixup=abc1234
git rebase -i --autosquash HEAD~5
```

### Bisect to find bugs
```bash
git bisect start
git bisect bad           # Current commit is bad
git bisect good v1.0.0   # Known good commit
# Git checks out middle commit — test it
git bisect good/bad      # Mark and continue
git bisect reset         # Return to original
```

### Stash with message
```bash
git stash push -m "WIP: feature X"
git stash list
git stash pop stash@{1}
```

### Reflog recovery
```bash
git reflog                    # See all HEAD movements
git reset --hard HEAD@{3}     # Go back 3 moves
git branch recovery HEAD@{5}  # Create branch from old state
```

## Anti-Patterns
### Don't: Force push to shared branches
```bash
git push --force origin main  # ❌ Destroys others' work
git push --force-with-lease origin feature  # ✅ Safer
```

## Quick Reference
| Task | Command |
|------|---------|
| Amend last | `git commit --amend` |
| Squash | `git rebase -i HEAD~N` |
| Bisect | `git bisect start` |
| Cherry-pick | `git cherry-pick abc1234` |
| Reflog | `git reflog` |
| Clean | `git clean -fd` |
