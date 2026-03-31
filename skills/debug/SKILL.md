---
name: debug
description: >
  Enable debug logging and troubleshoot issues systematically. Use when something isn't working, when investigating errors, or when the user says "debug" or "fix".
argument-hint: "[description]"
allowed-tools: Read, Grep, Bash, Edit
context: fork
agent: Explore
---

# Debug Mode

## Investigation

$ARGUMENTS

## Systematic Debugging Process

1. **Reproduce**: Understand the exact error
   - Check error messages and stack traces
   - Identify when it started happening
   - Note any recent changes

2. **Isolate**: Narrow down the problem
   - Use Grep to find error-related code
   - Check logs and output
   - Run affected code paths

3. **Root Cause**: Find the actual problem
   - Don't fix symptoms — find the root cause
   - Check edge cases and null values
   - Verify assumptions about data types

4. **Fix**: Apply the minimal fix
   - Make the smallest possible change
   - Don't refactor while debugging
   - Add a test that reproduces the bug

5. **Verify**: Confirm the fix works
   - Run relevant tests
   - Check for regressions
   - Document the root cause

## Debug Commands
- `go test -v -run TestName ./package/` — Run specific test
- `go test -race ./...` — Check for race conditions
- `go vet ./...` — Static analysis
- `grep -rn "error\|Error\|panic" .` — Find error handling

## Rules
- Never guess — always verify with evidence
- Fix the root cause, not the symptom
- Add a regression test for every bug fix
