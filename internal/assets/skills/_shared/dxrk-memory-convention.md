# dxrk-memory Convention

This file provides reference documentation for dxrk-memory (formerly DxrkMemory/Engram) protocol.

## Overview

dxrk-memory is a persistent memory system that survives across sessions and compactions.

## Core Commands

### mem_save — Save an observation

```
mem_save(
  title: "short descriptive title",
  type: "bugfix|decision|architecture|discovery|pattern|config|preference",
  scope: "project|personal",
  topic_key: "optional stable key for updates",
  content: "what, why, where, learned"
)
```

### mem_search — Search memory

```
mem_search(query: "search keywords", project: "optional project filter")
```

### mem_get_observation — Get full content

```
mem_get_observation(id: "observation-id")
```

### mem_context — Quick session check

```
mem_context()
```

### mem_session_summary — End of session

```
mem_session_summary(
  goal: "what we worked on",
  instructions: "user preferences discovered",
  discoveries: "technical findings",
  accomplished: "completed items",
  next_steps: "remaining work",
  relevant_files: "important paths"
)
```

## Usage Guidelines

1. **Save proactively** — After architecture decisions, bug fixes, pattern establishment
2. **Search before starting** — Check memory when user mentions past work
3. **Session close** — Always call mem_session_summary before ending

## For SDD Workflow

See `engram-convention.md` for SDD-specific persistence patterns.