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

# Code Review Skill

## Trigger
"review code", "PR review", "analyze changes", "check quality", "code review"

## Purpose
Perform comprehensive code review covering correctness, security, performance, and maintainability.

## Review Areas
- **Correctness**: Logic errors, edge cases, null handling
- **Security**: Vulnerabilities, injection risks, auth issues
- **Performance**: N+1 queries, memory leaks, inefficient algorithms
- **Maintainability**: Code smells, duplication, naming conventions
- **Testing**: Coverage, test quality, edge cases covered

## Process
1. Understand the change context and purpose
2. Analyze code structure and architecture fit
3. Check for common anti-patterns
4. Verify tests are adequate
5. Provide actionable feedback

## Output Format
```
## Code Review: <PR/Change Title>

### Summary
<Overall assessment>

### Must Fix
- <Critical issues>

### Should Fix
- <Important improvements>

### Suggestions
- <Nice-to-have improvements>
```
