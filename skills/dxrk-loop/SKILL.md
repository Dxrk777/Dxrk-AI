---
name: loop
description: >
  Run a prompt repeatedly on a schedule. Use for polling, monitoring, or periodic checks. Example: /loop 5m check if tests pass
disable-model-invocation: true
argument-hint: "<interval> <prompt>"
allowed-tools: Read, Bash, Grep
---

# Scheduled Loop

Running: $ARGUMENTS

## Parse Arguments
- First argument is the interval (e.g., 5m, 30s, 1h)
- Remaining arguments are the prompt to execute

## Execution Pattern
1. Parse the interval from $ARGUMENTS[0]
2. Execute the remaining prompt
3. Report results
4. Wait for interval
5. Repeat until stopped

## Supported Intervals
- `30s` — 30 seconds
- `5m` — 5 minutes
- `1h` — 1 hour

## Rules
- Always show a timestamp with each result
- If the command fails 3 consecutive times, stop and report
- Summarize patterns across runs
