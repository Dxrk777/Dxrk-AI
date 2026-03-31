---
name: batch
description: >
  Orchestrate large-scale changes across a codebase in parallel. Use when migrating, refactoring multiple files, or applying patterns across the project.
disable-model-invocation: true
argument-hint: "<instruction>"
allowed-tools: Read, Edit, Write, Grep, Glob, Bash
context: fork
agent: general
---

# Batch Orchestration

Orchestrate the following large-scale change across the codebase: $ARGUMENTS

## Workflow

1. **Research**: Understand the scope
   - Use Grep and Glob to find all affected files
   - Count the total number of changes needed
   - Identify any dependencies between changes

2. **Plan**: Break into independent units
   - Group files that can be changed in parallel
   - Identify shared patterns across files
   - Present a numbered plan with file counts

3. **Execute**: Apply changes systematically
   - Work through groups in order
   - Apply the same pattern to each file
   - Verify each group compiles/builds

4. **Verify**: Ensure everything works
   - Run `go build ./...` (or equivalent)
   - Run `go test ./...` (or equivalent)
   - Check for any remaining references to old patterns

## Rules
- Never change more than 5 files without pausing to verify
- Always show what you're about to change before editing
- If a change doesn't match the pattern, skip it and report
- Keep a running count of files changed
