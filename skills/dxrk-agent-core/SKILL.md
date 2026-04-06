---
name: dxrk-agent-core
description: >
  DXRK Agent Architecture knowledge base: agent loops, tool control, multi-agent, prompts.
  Trigger: agent architecture, tool permissions, system prompts, multi-agent, agent design.
license: Apache-2.0
metadata:
  author: dxrk
  version: "2.0"
---

## Overview

This skill contains DXRK's core knowledge for building production agent systems. Based on battle-tested patterns from leading agent frameworks, optimized for DXRK's architecture.

## Auto-Load This Skill When

- Building agent loops or query engines
- Implementing tool permissions/security
- Designing system prompts
- Creating multi-agent systems
- Any agent architecture work

---

# PART 1: AGENT LOOP ARCHITECTURE

## Dual-Layer Generator Pattern

DXRK uses **two SEPARATE layers**:

```
┌─────────────────────────────────────────────────────────────┐
│  OUTER: QueryEngine (AsyncGenerator)                      │
│  - Session management, message persistence                  │
│  - Token tracking, cost management                         │
│  - SDK output streaming                                    │
├─────────────────────────────────────────────────────────────┤
│  INNER: queryLoop (while true)                          │
│  - "Think-Act-Observe" cycle                            │
│  - API streaming calls                                     │
│  - Tool execution                                         │
│  - Auto-compaction                                        │
└─────────────────────────────────────────────────────────────┘
```

### Core Loop Flow

```
0. Pre-iteration: AutoCompact → Snip → Token check
   ↓
1. Stream API call → tool_use blocks trigger IMMEDIATE execution
   ↓
2. Decision: tool_use present? YES → execute / NO → terminate
   ↓
3. Tool execution (partitioned by concurrency safety)
   ↓
4. Termination check: abort, max_turns, stop_hook, token_budget
   ↓
5. Update messages → continue loop
```

### Termination Reasons

| Code | Meaning |
|------|---------|
| `completed` | No tool_use in final message |
| `aborted_streaming` | User interrupt |
| `max_turns` | Turn limit reached |
| `blocking_limit` | Context too large, can't compact |
| `stop_hook_prevented` | External hook blocked |

## Streaming Tool Execution

**KEY INSIGHT**: Tools don't wait for full API response.

```typescript
interface Tool {
  name: string;
  execute: (input: unknown) => Promise<Result>;
  isConcurrencySafe: boolean;  // DECLARE THIS!
}

// Read tools → SAFE (parallel)
// Write tools → UNSAFE (serial)
```

## Token Budget System

```typescript
// Before each API call:
if (tokenCount > threshold * 0.85) {
  messages = await compact(messages);  // LLM summary
}

// Nudge pattern for long tasks
function shouldNudge(messages, targetBudget): boolean {
  return used < targetBudget && !hasRecentToolCalls(messages);
}
```

## Error Recovery Paths (3 Independent)

```
Path 1: 413 prompt_too_long → Reactive Compact → Retry
Path 2: max_output_tokens → Upgrade to 64k + nudge (3x) → Fail
Path 3: model_overloaded → Switch fallback model → Notify
```

---

# PART 2: TOOL CONTROL & SECURITY

## 12-Step Permission Pipeline

```
┌──────────────────────────────────────────────────────────────────┐
│  STEP 1a: Tool in alwaysDenyRules?          → DENY (immune)    │
│  STEP 1b: Tool in alwaysAskRules?            → ASK (immune)       │
│  STEP 1c: tool.checkPermissions()            → DENY/allow         │
│  STEP 1d: Returns DENY?                    → DENY                │
│  STEP 1e: requiresUserInteraction()?          → ASK (immune)       │
│  STEP 1f: Content matches ask rule?          → ASK (immune)       │
│  STEP 1g: Security path (.git/, .bashrc)?   → ASK (immune)       │
│  STEP 2a: bypassPermissions mode?              → ALLOW              │
│  STEP 2b: Tool in alwaysAllowRules?            → ALLOW              │
│  STEP 2c: Content matches allow rule?         → ALLOW              │
│  STEP 3:   Default                            → ASK                │
│  STEP ✓:   Show confirmation dialog            → User decides      │
└──────────────────────────────────────────────────────────────────┘
```

## 5 Permission Modes

| Mode | Auto-Allow | Prompts |
|------|------------|---------|
| `default` | Nothing | All risky ops |
| `acceptEdits` | Local edits | External ops |
| `plan` | Nothing | Review plan first |
| `bypassPermissions` | Everything | Nothing (hard checks only) |
| `auto` | AI classifier decides | Depends |

## Rule Types

```typescript
// alwaysAllowRules
"Bash"                           // All Bash
"Bash(git status:*)"           // Specific command
"Edit(/src/**)"                 // Path pattern

// alwaysDenyRules
"Bash(rm -rf:*)"               // Dangerous commands
"Bash(npm publish:*)"          // Publishing

// alwaysAskRules (immune to bypass)
"Bash(git push:*)"             // All pushes need confirm
"Bash(curl * | bash:*)"        // Pipe to bash
```

## Bash Rule Matching

```typescript
// 1. stripSafeWrappers() → remove timeout/time/nice/nohup
// 2. splitCommandWithOperators() → handle &&, |, ;
// 3. Match patterns:
//    - Exact: "git status"
//    - Prefix: "git:*"
//    - Wildcard: "Bash"
```

## Denial Tracking

```typescript
// 3 consecutive denials → downgrade to interactive
// 20 cumulative denials → downgrade to interactive
// Prevents infinite loops of blocked operations
```

---

# PART 3: MULTI-AGENT COORDINATOR

## 3 Agent Execution Models

```
┌──────────────────────────────────────────────────────────────┐
│  LOCAL ASYNC (local_agent)                                  │
│  - AgentTool() spawns background worker                   │
│  - TaskID: a-{16HEX}                                     │
│  - Results via <task-notification> XML                     │
│  - SendMessageTool for continuation                        │
├──────────────────────────────────────────────────────────────┤
│  IN-PROCESS TEAMMATE (in_process_teammate)               │
│  - Team Swarm feature                                       │
│  - Identity: agentName@teamName                             │
│  - AsyncLocalStorage context isolation                      │
│  - Mailbox filesystem communication                         │
├──────────────────────────────────────────────────────────────┤
│  REMOTE (remote_agent)                                     │
│  - Deployed to remote worker cluster                       │
│  - TaskID: r-{16HEX}                                     │
│  - isolation: 'remote' parameter                           │
└──────────────────────────────────────────────────────────────┘
```

## Coordinator Pattern

```typescript
// Environment: DXRK_COORDINATOR_MODE=1

// Coordinator responsibilities:
// - DON'T execute work directly
// - Fork/parcel tasks to Workers
// - Synthesize results before responding
// - Use SendMessage to continue Workers
// - Use TaskStop to terminate Workers
```

## Task Notification Format

```xml
<task-notification>
  <task-id>a-1a2b3c4d5e6f7890</task-id>
  <status>completed</status>
  <summary>Found 3 auth-related files</summary>
  <result>The authentication flow uses...</result>
  <usage>
    <total-tokens>4521</total-tokens>
    <duration-ms>12400</duration-ms>
  </usage>
</task-notification>
```

## Mailbox Communication

```typescript
// File: ~/.claude/teams/{team}/inboxes/{agent}.json
[
  {
    "from": "leader",
    "text": "Please check the database schema",
    "timestamp": "2025-01-01T10:00:00Z",
    "read": false,
    "summary": "Check database schema"
  }
]

// Broadcast: from: "*" → all teammates
```

## Fork vs Worker

| Aspect | Fork | Worker |
|--------|-------|--------|
| Context | Inherits full parent | Fresh start |
| System prompt | Identical | Independent |
| Prompt cache | High hit rate | Independent |
| Recursion | Blocked by boilerplate | Blocked by default |

---

# PART 4: SYSTEM PROMPT ARCHITECTURE

## 14 Prompt Segments

```
┌─────────────────────────────────────────────────────────────┐
│  STATIC (70% of prompt, Global Cache)                     │
├─────────────────────────────────────────────────────────────┤
│  1. IntroSection → Identity, URL generation limits        │
│  2. SystemSection → Tool output rendering, hooks          │
│  3. DoingTasksSection → No over-design, no unnecessary     │
│  4. ActionsSection → Blast radius, reversible vs risky      │
│  5. UsingToolsSection → Dedicated tools > Bash             │
│  6. ToneAndStyleSection → No emoji unless asked          │
│  7. OutputEfficiencySection → Brief, direct, no filler    │
├═════════════════════════════════════════════════════════════┤
│  ════════════ BOUNDARY ════════════                       │
├═════════════════════════════════════════════════════════════┤
│  DYNAMIC (Session/User cached)                            │
├─────────────────────────────────────────────────────────────┤
│  8.  SessionGuidance → Agent tool usage                   │
│  9.  Memory → CLAUDE.md files (up to 5 levels up)      │
│  10. EnvironmentInfo → cwd, git, platform, model, date    │
│  11. Language → Forced response language                  │
│  12. MCP Instructions → Server-specific (NO CACHE)        │
│  13. Scratchpad → Session temp directory                  │
│  14. TokenBudget → Target token goal                      │
└─────────────────────────────────────────────────────────────┘
```

## 4-Layer Priority Override

```
┌─────────────────────────────────────────────────────────────┐
│  Priority 0 (HIGHEST): overrideSystemPrompt               │
│  Priority 1: Coordinator System Prompt                    │
│  Priority 2: agentSystemPrompt (PROACTIVE: append)       │
│  Priority 3: customSystemPrompt                            │
│  Priority 4 (LOWEST): defaultSystemPrompt                │
│  + ALWAYS: appendSystemPrompt (appended regardless)       │
└─────────────────────────────────────────────────────────────┘
```

## Key Prompt Insights

```markdown
# DOING TASKS (Immutable Rules)
- Don't add features beyond what was asked
- Don't create abstractions for one-time operations
- Three similar lines > one helper with one call
- Default to NO comments (only if WHY is non-obvious)
- Be honest about test failures

# TONE AND STYLE
- Only emoji if user explicitly asks
- Short, concise responses
- Code reference: file:line (e.g., src/main.ts:42)
- GitHub issue: owner/repo#123
- No colon before tool calls

# EXECUTING ACTIONS WITH CARE
- Local/reversible → freely take (edits, tests, servers)
- Risky → warrant confirmation (destructive, irreversible, affecting others)
- Do NOT use destructive as shortcut
```

## Memory/CLAUDE.md Loading

```typescript
// Scans up to 5 parent directories for CLAUDE.md
// Markdown format (human-friendly, GitHub previewable)
// Loaded into system prompt for project conventions
```

---

# PART 5: BUILT-IN AGENTS

## 5 Built-in Agents

| Agent | Model | Tools | Role |
|-------|-------|-------|------|
| **Explore** | haiku/inherit | Read-only | Fast codebase search |
| **Plan** | inherit | Read-only | Architecture design |
| **Verification** | inherit | Read + /tmp | Adversarial testing |
| **Guide** | haiku | Search + web | DXRK usage help |
| **General** | default | All | Fallback |

## Explore Agent Prompt

```markdown
=== CRITICAL: READ-ONLY MODE - NO FILE MODIFICATIONS ===

STRICTLY PROHIBITED:
- Creating files (no Write, touch)
- Modifying files (no Edit)
- Deleting files (no rm)
- Creating temp files anywhere including /tmp

STRENGTHS:
- Rapid file finding with glob
- Powerful regex search with Grep
- File reading and analysis

NOTE: Be FAST. Spawn multiple parallel tool calls.
```

## Verification Agent Philosophy

```markdown
Your job is NOT to confirm — it's to BREAK.

Two failure patterns to AVOID:
1. Verification avoidance: "I read the code, it looks correct"
2. First 80% seduction: Polished UI → pass without checking

YOUR VALUE IS IN THE LAST 20%.
```

### Verification Strategies by Change Type

**Frontend:**
- Start dev server
- Use browser automation MCP tools
- Click, screenshot, check console
- Don't say "needs real browser" — check for MCP tools first

**Backend:**
- Start service → curl endpoints
- Verify response shape (not just status code)
- Test error handling and edge cases
- Check edge inputs (null, overflow, unicode)

**Bug Fix:**
- REPRODUCE the original bug FIRST
- Verify fix works
- Run regression tests
- Check adjacent functionality

**Refactoring:**
- Tests MUST pass (behavior unchanged)
- Diff public API surface
- Same input → same output
- No performance regression

### Verification Output Format

```markdown
### Check: POST /api/register rejects short password
**Command run:**
  curl -s -X POST localhost:8000/api/register \
    -H 'Content-Type: application/json' \
    -d '{"email":"t@t.co","password":"short"}'
**Output observed:**
  {"error": "password must be at least 8 characters"} (HTTP 400)
**Result: PASS**

---
VERDICT: PASS | FAIL | PARTIAL
```

## Plan Agent

```markdown
Your role: Software architect and planning specialist

Process:
1. Understand requirements
2. Explore thoroughly (read provided files, find patterns)
3. Design solution (consider trade-offs)
4. Detail the plan (step-by-step, dependencies)

Required output:
### Critical Files for Implementation
List 3-5 most critical files:
- path/to/file1.ts
- path/to/file2.ts
```

---

# PART 6: MCP INTEGRATION

## MCP Tool Registration

```
1. Connect MCP Server (stdio/SSE/HTTP/WebSocket)
2. ListTools RPC → discover all tools
3. Normalize names: my-tool → mcp__server-name__my_tool
4. Inject into tool pool
5. Sub-agents inherit parent MCP clients
```

## MCP Server States

| State | Meaning |
|-------|---------|
| `connected` | Tools available, instructions injected |
| `needs-auth` | OAuth/API key required |
| `pending` | Reconnecting |
| `failed` | Connection failed |
| `disabled` | User disabled |

## MCP Configuration Scopes

```yaml
local:    ~/.project/.claude/mcp.json
user:     ~/.claude/mcp.json
project:  project root
managed:   enterprise push
```

---

# QUICK REFERENCE

## Key Architecture Patterns

### Dual-Layer Pattern
```
OUTER: Session management, persistence, SDK
INNER: while(true) loop, think-act-observe
```

### Streaming Tools
```
Tool execution starts BEFORE API stream completes
isConcurrencySafe: read=true, write=false
```

### Error Recovery
```
413 → Compact → Retry
max_tokens → Upgrade + Nudge (3x)
overload → Fallback model
```

### Multi-Agent
```
Coordinator → spawns Workers
Workers → complete tasks
Results → task-notification XML injection
```

### Permission Pipeline
```
12 steps, high-priority checks first
immune-to-bypass: safety paths, ask rules
```

### Prompt Registry
```
14 segments, BOUNDARY separates static/dynamic
4-layer priority override
```

---

## Source References

Patterns extracted and adapted for DXRK from leading agent frameworks:

- `/home/dxrk/Documentos/DARK GORE/claude-source-learning/`
- DXRK internal patterns and production experience

---

## Commands

```bash
# Test agent loop patterns
npm run test:agent-loop

# Benchmark permissions
npm run test:permissions

# Test multi-agent
npm run test:multi-agent
```

---

*🔥 DXRK Agent Architecture — Built for production*
