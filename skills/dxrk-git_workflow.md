<!--
╔═══════════════════════════════════════════════════════════════════════════════╗
║                                                                               ║
║     ██████╗ ██╗  ██╗██████╗ ██╗  ██╗     ██████╗ ███╗   ███╗███████╗       ║
║     ██╔══██╗╚██╗██╔╝██╔══██╗██║ ██╔╝    ██╔═══██╗████╗ ████║██╔════╝       ║
║     ██║  ██║ ╚███╔╝ ██████╔╝█████╔╝     ██║   ██║██╔████╔██║█████╗         ║
║     ██║  ██║ ██╔██╗ ██╔══██╗██╔═██╗     ██║   ██║██║╚██╔╝██║██╔══╝         ║
║     ██████╔╝██╔╝ ██╗██║  ██║██║  ██╗    ╚██████╔╝██║ ╚═╝ ██║███████╗       ║
║     ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝     ╚═════╝ ╚═╝     ╚═╝╚══════╝       ║
║                                                                               ║
║              DXRK AI INTELLIGENCE SYSTEM — SKILL PACKAGE                      ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
-->

# Git Workflow Skill

## Trigger
"git workflow", "branch strategy", "conventional commits", "PR", "merge strategy"

## Branch Naming
```
feature/TICKET-description
fix/TICKET-description
hotfix/TICKET-description
chore/task-description
```

## Conventional Commits
```
feat: add user authentication
fix: resolve login timeout issue
refactor: extract validation logic
docs: update API documentation
test: add integration tests
chore: update dependencies
```

## PR Process
1. Fork from main
2. Create feature branch
3. Make atomic commits
4. Write descriptive commit messages
5. Open PR with template
6. Address review feedback
7. Squash and merge

## PR Template
```markdown
## Summary
<What does this PR do?>

## Changes
- <list of changes>

## Testing
- <how was this tested?>

## Screenshots (if UI)
```

## Git Commands Reference
```bash
# Clean up merged branches
git branch --merged | grep -v main | xargs git branch -d

# Interactive rebase
git rebase -i HEAD~3

# Stash with message
git stash push -m "WIP: feature"

# Amend without changing message
git commit --amend --no-edit
```
