---
name: dxrk-prompt-registry
description: >
  Design modular system prompts using the segment registry pattern with BOUNDARY caching.
  Trigger: System prompts, prompt engineering, caching strategies, prompt composition.
license: Apache-2.0
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use

- Building system prompts for AI agents
- Organizing large prompts into maintainable segments
- Implementing prompt caching strategies
- Creating priority override chains
- Managing dynamic vs static prompt content

## Segment Registry Pattern

The key insight: **separate concerns into independent functions, then compose**.

```
┌──────────────────────────────────────────────────────────────────┐
│                   PROMPT REGISTRY ARCHITECTURE                    │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    REGISTERED SEGMENTS                      │ │
│  ├─────────────────────────────────────────────────────────────┤ │
│  │                                                             │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │ │
│  │  │   SYSTEM    │  │   AGENT     │  │   RULES     │        │ │
│  │  │  IDENTITY   │  │  PROFILE    │  │  & SAFETY   │        │ │
│  │  │  (static)  │  │  (dynamic)  │  │  (static)   │        │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘        │ │
│  │                                                             │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │ │
│  │  │   TOOLS     │  │ ENVIRONMENT │  │    MCP      │        │ │
│  │  │  CATALOG    │  │   CONTEXT   │  │ INSTRUCTIONS│        │ │
│  │  │  (dynamic) │  │  (dynamic)  │  │  (dynamic)  │        │ │
│  │  └─────────────┘  └─────────────┘  └─────────────┘        │ │
│  │                                                             │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                            │                                     │
│                            ▼                                     │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    COMPOSE + CACHE                         │ │
│  │                                                             │ │
│  │  [STATIC segments] ──► [BOUNDARY] ──► [DYNAMIC segments]  │ │
│  │  (Global Cache)      (separator)      (Per-request)        │ │
│  │                                                             │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

## Implementation Pattern

```typescript
interface PromptSegment {
  name: string;
  content: string | (() => string | Promise<string>);
  cachePolicy: 'static' | 'dynamic' | 'per-request';
  priority: number;
}

// Segment registration decorator
function registerSegment(options: {
  name: string;
  cachePolicy: 'static' | 'dynamic' | 'per-request';
  priority?: number;
}) {
  return function <T extends (...args: any[]) => string | Promise<string>>(
    target: T
  ): void {
    promptRegistry.register({
      name: options.name,
      content: target,
      cachePolicy: options.cachePolicy,
      priority: options.priority ?? 0
    });
  };
}

const promptRegistry = new PromptSegmentRegistry();

@registerSegment({ name: 'system-identity', cachePolicy: 'static', priority: 100 })
function systemIdentity(): string {
  return `You are an expert software engineer assistant.
Your goal is to help users build high-quality software.
You have deep knowledge of:
- Design patterns and architecture
- Multiple programming languages
- Testing best practices
- Security considerations`;
}

@registerSegment({ name: 'safety-rules', cachePolicy: 'static', priority: 90 })
function safetyRules(): string {
  return `
## Safety Rules
1. Always verify inputs before processing
2. Never expose sensitive data in responses
3. Confirm destructive actions before executing
4. Handle errors gracefully and informatively`;
}

@registerSegment({ name: 'tools-catalog', cachePolicy: 'dynamic', priority: 50 })
function toolsCatalog(): string {
  const tools = getAvailableTools();
  return tools.map(t => `- ${t.name}: ${t.description}`).join('\n');
}

@registerSegment({ name: 'environment-context', cachePolicy: 'per-request', priority: 40 })
function environmentContext(): string {
  return `Current directory: ${process.cwd()}
Project: ${getProjectName()}
Stack: ${getTechStack()}`;
}

@registerSegment({ name: 'mcp-instructions', cachePolicy: 'dynamic', priority: 30 })
function mcpInstructions(): string {
  const mcpTools = getMCPTools();
  if (mcpTools.length === 0) return '';
  
  return `
## MCP Tools Available
${mcpTools.map(t => `- mcp__${t.server}__${t.name}`).join('\n')}`;
}
```

## BOUNDARY Cache Strategy

The BOUNDARY marker separates cacheable static content from dynamic content:

```typescript
interface ComposedPrompt {
  static: string;      // Can be globally cached
  boundary: string;    // Cache separator
  dynamic: string;     // Per-request (no cache)
  full: string;        // Complete prompt
}

function composePrompt(context: PromptContext): ComposedPrompt {
  const segments = promptRegistry.getAll();
  
  // Sort by priority (higher = appears first)
  segments.sort((a, b) => b.priority - a.priority);
  
  const staticSegments: string[] = [];
  const dynamicSegments: string[] = [];
  
  for (const segment of segments) {
    const content = typeof segment.content === 'function' 
      ? segment.content() 
      : segment.content;
    
    if (segment.cachePolicy === 'static') {
      staticSegments.push(content);
    } else {
      dynamicSegments.push(content);
    }
  }
  
  return {
    static: staticSegments.join('\n\n'),
    boundary: '\n\n<!-- BOUNDARY -->\n\n',
    dynamic: dynamicSegments.join('\n\n'),
    full: staticSegments.join('\n\n') + 
          '\n\n<!-- BOUNDARY -->\n\n' + 
          dynamicSegments.join('\n\n')
  };
}
```

## 4-Layer Priority Override Chain

DXRK uses this priority chain:

```
┌──────────────────────────────────────────────────────────────────┐
│                    PROMPT OVERRIDE PRIORITY                       │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Layer 1 (HIGHEST): overridePrompt    ──► COMPLETE override     │
│           │                                                         │
│           ▼                                                         │
│  Layer 2: Coordinator Prompt          ──► Coordinator additions     │
│           │                                                         │
│           ▼                                                         │
│  Layer 3: Agent Definition Prompt    ──► Agent-specific rules    │
│           │                                                         │
│           ▼                                                         │
│  Layer 4: User Custom Prompt        ──► User preferences         │
│           │                                                         │
│           ▼                                                         │
│  Layer 5 (LOWEST): Default Prompt   ──► Framework defaults      │
│                                                                  │
│  + appendSystemPrompt()            ──► ALWAYS appended to end    │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### Implementation

```typescript
interface PromptOverride {
  layer: 1 | 2 | 3 | 4 | 5;
  type: 'replace' | 'append';
  content: string;
  protected?: boolean;  // Immune to being overridden
}

class PromptBuilder {
  private overrides: PromptOverride[] = [];
  
  // Add override (respects priority)
  addOverride(override: PromptOverride): void {
    // Don't add if lower-priority override exists
    const existing = this.overrides.find(
      o => o.layer > override.layer && !o.protected
    );
    
    if (!existing) {
      this.overrides.push(override);
      this.overrides.sort((a, b) => a.layer - b.layer);
    }
  }
  
  // Append always goes to end (safe operation)
  append(content: string): void {
    this.overrides.push({
      layer: 0,  // Special: always last
      type: 'append',
      content
    });
  }
  
  build(): string {
    // Start with default
    let prompt = DEFAULT_PROMPT;
    
    // Apply overrides in priority order
    for (const override of this.overrides) {
      if (override.type === 'replace') {
        prompt = override.content;
      } else {
        prompt += '\n\n' + override.content;
      }
    }
    
    return prompt;
  }
}
```

## Dynamic Segment Examples

### Session/Project Context

```typescript
@registerSegment({ name: 'project-context', cachePolicy: 'per-request' })
function projectContext(): string {
  const project = getCurrentProject();
  
  return `
## Project Context
Name: ${project.name}
Language: ${project.language}
Framework: ${project.framework || 'None'}
Key Files:
${project.keyFiles.map(f => `- ${f}`).join('\n')}`;
}
```

### Working Memory

```typescript
@registerSegment({ name: 'working-memory', cachePolicy: 'per-request' })
function workingMemory(): string {
  const memory = getSessionMemory();
  
  if (memory.length === 0) return '';
  
  return `
## Session Memory
${memory.map(m => `- ${m}`).join('\n')}`;
}
```

### Recent Actions

```typescript
@registerSegment({ name: 'recent-actions', cachePolicy: 'per-request' })
function recentActions(): string {
  const recent = getRecentToolCalls().slice(-5);
  
  if (recent.length === 0) return '';
  
  return `
## Recent Actions
${recent.map(a => `[${a.timestamp}] ${a.tool}: ${a.summary}`).join('\n')}`;
}
```

## Key Files (DXRK Reference)

| File | Lines | Purpose |
|------|-------|---------|
| `prompts.ts` | 914 | System prompt construction |
| `systemPrompt.ts` | 123 | Override priority chain |
| `systemPromptSections.ts` | - | Section registry |

## Commands

```bash
# Test prompt composition
npm run test:prompt-compose

# Benchmark caching
npm run benchmark:prompt-cache

# Validate prompt segments
npm run validate:prompts
```

## Resources

- **System Prompt Source**: `/home/dxrk/Documentos/DARK GORE/claude-source-learning/system-prompt.html`
