---
name: simplify
description: >
  Review recently changed files for code quality, reuse, and efficiency. Use when the user says "review", "clean up", "simplify", or "optimize".
argument-hint: "[focus area]"
allowed-tools: Read, Edit, Write, Grep, Glob, Bash
context: fork
agent: general
---

# Code Simplification

## Focus
$ARGUMENTS

## Review Process

### 1. Find Recent Changes
```bash
git diff --name-only HEAD~5
```

### 2. Analyze Each File
For each changed file:
- **Duplication**: Is the same pattern repeated? Extract to a function
- **Complexity**: Can nested logic be simplified? Use early returns
- **Naming**: Are names clear and consistent?
- **Error handling**: Are errors properly handled?
- **Dead code**: Is there unused code? Remove it

### 3. Apply Fixes
- Group related simplifications
- Make one logical change at a time
- Verify after each change

### 4. Run Quality Checks
```bash
go vet ./...
go test ./...
golangci-lint run
```

## Rules
- Don't change behavior — only simplify implementation
- If you're unsure, leave it as is and note it
- Prefer readability over cleverness
