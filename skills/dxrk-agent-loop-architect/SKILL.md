---
name: dxrk-agent-loop-architect
description: >
  Architect agent systems using the DXRK dual-layer generator pattern.
  Trigger: Building agent loops, query engines, streaming tool execution.
license: Apache-2.0
metadata:
  author: dxrk
  version: "2.0"
---

## When to Use

- Building custom agent loops from scratch
- Implementing streaming tool execution
- Designing session management vs behavior separation
- Creating token budget and auto-compaction systems
- Implementing error recovery paths

## Dual-Layer Architecture (DXRK Pattern)

### The Core Insight

DXRK uses **two layers** that must be SEPARATED:

```
┌─────────────────────────────────────────────────────────────┐
│  OUTER LAYER: QueryEngine (AsyncGenerator)                 │
│  - Manages cross-turn message history                      │
│  - Handles SDK output streaming                             │
│  - Tracks token usage & costs                               │
│  - Message persistence (sessionStorage)                     │
│  - USD budget checks                                        │
│  - Permission denial tracking                               │
├─────────────────────────────────────────────────────────────┤
│  INNER LAYER: queryLoop (while(true))                     │
│  - Single request's complete "think-act-observe" cycle     │
│  - Streams API calls                                        │
│  - Extracts and executes tool_use blocks                    │
│  - Auto-compaction / Context Collapse                      │
│  - Error recovery (prompt_too_long, max_tokens)           │
│  - Token budget nudge mechanism                             │
└─────────────────────────────────────────────────────────────┘
```

### Implementation Pattern

```typescript
// OUTER: Session Management
class QueryEngine {
  private mutableMessages: Message[] = [];
  
  async *submitMessage(userInput: string): AsyncGenerator<SDKMessage> {
    for await (const message of this.queryLoop(userInput)) {
      this.mutableMessages.push(message);  // Persist across turns
      yield this.buildSDKMessage(message);   // Stream to caller
    }
  }
}

// INNER: Behavior Loop
async function* queryLoop(input: string): AsyncGenerator<Message> {
  let messages = buildContext(input);
  
  while (true) {
    // 1. Auto-compaction check BEFORE API call
    if (await shouldCompact(messages)) {
      messages = await compact(messages);
    }
    
    // 2. Stream API call - tools can start BEFORE stream ends
    for await (const chunk of api.stream(messages)) {
      if (chunk.type === 'tool_use') {
        executor.addTool(chunk.tool); // Start execution immediately!
      }
      yield chunk;
    }
    
    // 3. Check if more tools needed
    if (!hasToolCalls(yieldedMessages)) break;
    
    // 4. Execute tools
    const results = await executor.runTools();
    messages = [...messages, ...results];
  }
}
```

## Streaming Tool Executor Pattern

**KEY INSIGHT**: Tools don't wait for the full API response.

```typescript
interface Tool {
  name: string;
  execute: (input: unknown) => Promise<Result>;
  isConcurrencySafe: boolean; // CRITICAL: declare this!
}

class StreamingToolExecutor {
  private pending: Map<string, Promise<Result>> = new Map();
  
  addTool(toolCall: ToolCall): void {
    // Start execution IMMEDIATELY when tool_use block arrives
    const promise = tool.execute(toolCall.input);
    this.pending.set(toolCall.id, promise);
  }
  
  async runTools(): Promise<ToolResult[]> {
    const results: ToolResult[] = [];
    
    // Partition by concurrency safety
    const [safe, unsafe] = partition(this.pending, t => t.isConcurrencySafe);
    
    // Safe tools: run in PARALLEL (up to 10)
    const safeResults = await Promise.all(safe.map(t => t.promise));
    results.push(...safeResults);
    
    // Unsafe tools: run SERIALLY (writes, bash)
    for (const t of unsafe) {
      results.push(await t.promise);
    }
    
    return results;
  }
}
```

### Concurrency Safety Rules

| Tool Type | Concurrency Safe | Reason |
|-----------|-----------------|--------|
| Read files | ✅ YES | Read-only, no side effects |
| Glob/Search | ✅ YES | Read-only |
| MCP tools (read) | ✅ YES | Depends on MCP server |
| Write files | ❌ NO | Potential race conditions |
| Bash execute | ❌ NO | Side effects, ordering matters |
| External API (POST) | ❌ NO | Non-idempotent |

## Token Budget & Auto-Compaction

### Pre-call Check Pattern

```typescript
async function shouldCompact(messages: Message[]): Promise<boolean> {
  const tokenCount = await countTokens(messages);
  const threshold = getContextLimit() * 0.85; // 85% threshold
  
  return tokenCount > threshold;
}

async function compact(messages: Message[]): Promise<Message[]> {
  // Option 1: LLM Summary (best quality)
  const summary = await llm.summarize(messages);
  return [messages[0], createSummaryMessage(summary), ...messages.at(-1)];
  
  // Option 2: Sliding Window (faster)
  const windowSize = Math.floor(getContextLimit() * 0.5);
  return messages.slice(-windowSize);
}
```

### Nudge Pattern for Long Tasks

```typescript
// When token budget is specified, inject nudge messages
function shouldNudge(messages: Message[], targetBudget: number): boolean {
  const used = countTokens(messages);
  return used < targetBudget && !hasRecentToolCalls(messages);
}

function injectNudge(): Message {
  return {
    role: 'user',
    content: "Continue working on the task. You've used X tokens but should aim for Y. Keep going."
  };
}
```

## Error Recovery Paths (3 Independent Paths)

```
┌──────────────────────────────────────────────────────────────────┐
│                     ERROR RECOVERY PATHS                          │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Path 1: prompt_too_long (413)                                  │
│    → Trigger Reactive Compact                                    │
│    → Auto-summarize history                                      │
│    → Retry with compressed context                              │
│                                                                  │
│  Path 2: max_output_tokens                                      │
│    → Upgrade to max_tokens=64k                                   │
│    → Inject nudge message                                        │
│    → Retry up to 3 times                                        │
│                                                                  │
│  Path 3: model_overloaded                                       │
│    → Switch to fallback model                                    │
│    → Clear orphan messages                                       │
│    → Notify user via system message                              │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### Implementation

```typescript
async function handleAPIError(error: APIError, context: LoopContext): Promise<RecoveryAction> {
  switch (error.type) {
    case 'PROMPT_TOO_LONG':
      const compressed = await compactMessages(context.messages);
      return { action: 'RETRY', messages: compressed };
    
    case 'MAX_TOKETS':
      if (context.retryCount < 3) {
        return { 
          action: 'RETRY_NUDGE', 
          maxTokens: 65536,
          nudge: "Continue your previous thought..."
        };
      }
      return { action: 'FAIL', reason: 'max_retries_exceeded' };
    
    case 'MODEL_UNAVAILABLE':
      return { 
        action: 'FALLBACK_MODEL',
        model: getFallbackModel()
      };
    
    default:
      return { action: 'FAIL', error };
  }
}
```

## Termination Conditions

| Condition | Reason Code | Trigger |
|-----------|-------------|---------|
| Normal completion | `completed` | No tool_use in final assistant message |
| User interrupt | `aborted_streaming` | AbortController.signal.aborted |
| Turn limit | `max_turns` | maxTurns config reached |
| Context blocked | `blocking_limit` | Token > threshold, can't compact |
| Stop hook | `stop_hook_prevented` | Hook returns block signal |

## Key Files (DXRK Reference)

| File | Lines | Purpose |
|------|-------|---------|
| `query.ts` | - | Main loop (while true) |
| `QueryEngine.ts` | - | Session layer |
| `StreamingToolExecutor.ts` | - | Concurrent tool execution |
| `toolOrchestration.ts` | - | Tool partitioning |
| `autoCompact.ts` | - | Token budget & compression |
| `tokenBudget.ts` | - | Budget management |

## Commands

```bash
# Test agent loop locally
npm run test:agent-loop

# Benchmark token consumption
npm run benchmark:tokens

# Run streaming performance test
npm run test:streaming
```

## Resources

- **DXRK Agent Core**: `dxrk-agent-core` skill (comprehensive guide)
- **Patterns Source**: `/home/dxrk/Documentos/DARK GORE/claude-source-learning/`
