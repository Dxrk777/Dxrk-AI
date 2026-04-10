# Components, Skills & Presets

← [Back to README](../README.md)

---

## Components

| Component | ID | Description |
|-----------|-----|-------------|
| DxrkMemory | `DxrkMemory` | Persistent cross-session memory — managed automatically by the agent, no manual interaction needed |
| SDD | `sdd` | Spec-Driven Development workflow (9 phases) — the agent handles SDD organically when the task warrants it, or when you ask; you don't need to learn the commands |
| Skills | `skills` | Curated coding skill library |
| Context7 | `context7` | MCP server for live framework/library documentation |
| Persona | `persona` | Dxrk, neutral, or custom behavior mode |
| Permissions | `permissions` | Security-first defaults and guardrails |
| Dxrk | `dxrk` | Dxrk Guardian Angel — AI provider switcher |
| Theme | `theme` | Dxrk Kanagawa theme overlay |

## Dxrk Behavior

`dxrk --component dxrk` installs/provisions the `dxrk` binary globally on your machine.

It does **not** run project-level hook setup automatically (`dxrk init` / `dxrk install`) because that should be an explicit decision per repository.

After global install, enable Dxrk per project with:

```bash
dxrk init
dxrk install
```

---

## Skills

### Included Skills (installed by dxrk)

14 skill files organized by category, embedded in the binary and injected into your agent's configuration:

#### SDD (Spec-Driven Development)

| Skill | ID | Description |
|-------|-----|-------------|
| SDD Init | `sdd-init` | Bootstrap SDD context in a project |
| SDD Explore | `sdd-explore` | Investigate codebase before committing to a change |
| SDD Propose | `sdd-propose` | Create change proposal with intent, scope, approach |
| SDD Spec | `sdd-spec` | Write specifications with requirements and scenarios |
| SDD Design | `sdd-design` | Technical design with architecture decisions |
| SDD Tasks | `sdd-tasks` | Break down a change into implementation tasks |
| SDD Apply | `sdd-apply` | Implement tasks following specs and design |
| SDD Verify | `sdd-verify` | Validate implementation matches specs |
| SDD Archive | `sdd-archive` | Sync delta specs to main specs and archive |
| Judgment Day | `judgment-day` | Parallel adversarial review — two independent judges review the same target |

#### Foundation

| Skill | ID | Description |
|-------|-----|-------------|
| Go Testing | `go-testing` | Go testing patterns including Bubbletea TUI testing |
| Skill Creator | `skill-creator` | Create new AI agent skills following the Agent Skills spec |
| Branch & PR | `branch-pr` | PR creation workflow with conventional commits, branch naming, and issue-first enforcement |
| Issue Creation | `issue-creation` | Issue filing workflow with bug report and feature request templates |

These foundation skills are installed by default with both `full-dxrk` and `ecosystem-only` presets.

### Coding Skills (separate repository)

For framework-specific skills (React 19, Angular, TypeScript, Tailwind 4, Zod 4, Playwright, etc.), see [Dxrk/Dxrk-Skills](https://github.com/Dxrk777/Dxrk-Skills). These are maintained by the community and installed separately by cloning the repo and copying skills to your agent's skills directory.

---

## Presets

| Preset | ID | What's Included |
|--------|-----|-----------------|
| Full Dxrk | `full-dxrk` | All components (DxrkMemory + SDD + Skills + Context7 + Dxrk + Persona + Permissions + Theme) + all skills + dxrk persona |
| Ecosystem Only | `ecosystem-only` | Core components (DxrkMemory + SDD + Skills + Context7 + Dxrk) + all skills + dxrk persona |
| Minimal | `minimal` | DxrkMemory + SDD skills only |
| Custom | `custom` | You pick components, skills, and persona individually |
