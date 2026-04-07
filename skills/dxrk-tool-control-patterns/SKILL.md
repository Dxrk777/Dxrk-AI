---
name: dxrk-tool-control-patterns
description: >
  Implement layered tool permission systems with immune-to-bypass security.
  Trigger: Tool permissions, bash rules, security controls, permission modes.
license: Apache-2.0
metadata:
  author: dxrk
  version: "2.0"
---

## When to Use

- Building tool permission systems
- Implementing bash command allow/deny rules
- Designing multi-mode security (dev/prod)
- Creating immune-to-bypass security layers
- Implementing content-level tool filtering

## 12-Step Permission Decision Pipeline

DXRK's security model uses a **sequential decision pipeline** where more dangerous checks come FIRST:

```
┌──────────────────────────────────────────────────────────────────┐
│                  12-STEP PERMISSION PIPELINE                      │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Step 1-2:  HARD-CODED SAFETY CHECKS (immune to bypass)         │
│             - Deny list patterns                                  │
│             - Dangerous command blocking                          │
│                                                                  │
│  Step 3-4:  BYPASS MODE CHECK                                    │
│             - If bypassPermissions=true, skip user prompts       │
│             - But HARD checks still run!                          │
│                                                                  │
│  Step 5-6:  MODE-SPECIFIC RULES                                  │
│             - Development vs Production rules differ             │
│                                                                  │
│  Step 7-8:  TOOL-SPECIFIC PERMISSIONS                           │
│             - Read tools have different rules than write tools    │
│                                                                  │
│  Step 9-10: USER CONFIGURATION                                   │
│             - ~/.claude/settings.json                            │
│             - Project CLAUDE.md overrides                         │
│                                                                  │
│  Step 11:   CONTENT-LEVEL FILTERING                              │
│             - Regex patterns on command content                  │
│                                                                  │
│  Step 12:   FINAL DECISION                                       │
│             - Allow / Deny / Ask User                           │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

## Implementation Pattern

```typescript
interface PermissionResult {
  decision: 'allow' | 'deny' | 'ask';
  reason: string;
  bypassImmune: boolean;
}

interface Tool {
  name: string;
  isDangerous: boolean;
  contentFilter?: RegExp[];
}

class PermissionEngine {
  private hardcodedDenials: Pattern[] = [
    /rm\s+-rf\s+\/(?!.*--no-preserve-root)/,
    /drop\s+database/i,
    /format\s+.*drive/i,
  ];
  
  private immuneChecks = [
    // These ALWAYS run, even in bypass mode
    (tool: Tool, input: unknown) => this.checkHardcodedDenials(tool, input),
    (tool: Tool, input: unknown) => this.checkDestructivePatterns(tool, input),
  ];
  
  async hasPermission(
    tool: Tool,
    input: unknown,
    context: PermissionContext
  ): Promise<PermissionResult> {
    // PHASE 1: Immune checks (NEVER bypassed)
    for (const check of this.immuneChecks) {
      const result = check(tool, input);
      if (result.deny) {
        return { 
          decision: 'deny', 
          reason: result.reason,
          bypassImmune: true
        };
      }
    }
    
    // PHASE 2: Bypass mode check
    if (context.bypassPermissions) {
      return { decision: 'allow', reason: 'bypass_mode', bypassImmune: false };
    }
    
    // PHASE 3: Mode-specific rules
    const modeResult = await this.checkModeRules(tool, input, context.mode);
    if (modeResult.decision !== 'ask') {
      return modeResult;
    }
    
    // PHASE 4: Tool-specific permissions
    const toolResult = await this.checkToolRules(tool, input);
    if (toolResult.decision !== 'ask') {
      return toolResult;
    }
    
    // PHASE 5: User configuration
    const userResult = await this.checkUserConfig(tool, input, context.userConfig);
    if (userResult.decision !== 'ask') {
      return userResult;
    }
    
    // PHASE 6: Content-level filtering
    const contentResult = await this.checkContent(tool, input);
    if (contentResult.deny) {
      return { decision: 'deny', reason: contentResult.reason, bypassImmune: false };
    }
    
    return { decision: 'ask', reason: 'user_confirmation_required', bypassImmune: false };
  }
}
```

## 5 Permission Modes

| Mode | Description | Auto-Allow | User Prompts |
|------|-------------|------------|-------------|
| `development` | Local dev, fast iteration | Read + safe write | Destructive ops |
| `production` | Live systems, strict | Nothing | Everything |
| `restricted` | Learning/sandbox | Nothing | Everything |
| `gentle` | Beginners, maximum safety | Nothing | All writes |
| `bypass` | CI/CD, automation | Everything | Nothing (hard checks only) |

## Bash Rule Pattern

```typescript
interface BashRule {
  pattern: RegExp | string;
  decision: 'allow' | 'deny' | 'ask';
  reason?: string;
}

class BashPermissionEngine {
  private rules: BashRule[] = [
    // IMMEDIATE DENY (immune to bypass for safety)
    { pattern: /^rm\s+-rf\s+\//, decision: 'deny', reason: 'Root deletion blocked' },
    { pattern: /^dd\s+.*of=\/(?!.*conv=notrunc)/, decision: 'deny', reason: 'Disk write blocked' },
    { pattern: /^mkfs/, decision: 'deny', reason: 'Format blocked' },
    { pattern: /^:()\s*{\s*:\|:&\s*};/, decision: 'deny', reason: 'Fork bomb detected' },
    
    // GIT RULES (configurable)
    { pattern: /^git\s+status/, decision: 'allow', reason: 'Read-only' },
    { pattern: /^git\s+log/, decision: 'allow', reason: 'Read-only' },
    { pattern: /^git\s+diff/, decision: 'allow', reason: 'Read-only' },
    { pattern: /^git\s+add/, decision: 'ask', reason: 'Staging changes' },
    { pattern: /^git\s+commit/, decision: 'ask', reason: 'Creating commit' },
    { pattern: /^git\s+push/, decision: 'ask', reason: 'Pushing to remote' },
    { pattern: /^git\s+push\s+--force/, decision: 'deny', reason: 'Force push blocked by default' },
    
    // NPM/YARN RULES
    { pattern: /^npm\s+(install|ci)/, decision: 'allow', reason: 'Installing dependencies' },
    { pattern: /^npm\s+run/, decision: 'allow', reason: 'Running scripts' },
    { pattern: /^npm\s+publish/, decision: 'ask', reason: 'Publishing package' },
    
    // DOCKER RULES
    { pattern: /^docker\s+(ps|images|logs)/, decision: 'allow', reason: 'Read operations' },
    { pattern: /^docker\s+run/, decision: 'ask', reason: 'Starting container' },
    { pattern: /^docker\s+rm/, decision: 'ask', reason: 'Removing container' },
    { pattern: /^docker\s+system\s+prune/, decision: 'deny', reason: 'Destructive cleanup blocked' },
  ];
  
  evaluate(command: string): PermissionResult {
    for (const rule of this.rules) {
      if (this.matches(command, rule.pattern)) {
        return {
          decision: rule.decision,
          reason: rule.reason || `Matched rule: ${rule.pattern}`,
          bypassImmune: rule.decision === 'deny' && this.isHardcoded(rule)
        };
      }
    }
    
    // Default: ask for unknown commands
    return { decision: 'ask', reason: 'no_matching_rule' };
  }
}
```

## Hooks System Pattern

```typescript
// Four hook types for tool lifecycle
interface Hooks {
  'pre-tool-use': PreToolHook[];
  'post-tool-use': PostToolHook[];
  'user-prompt-submit': UserPromptHook[];
  'stop': StopHook[];
}

interface PreToolHook {
  name: string;
  command: string;  // Shell command, interacts via stdin/stdout
  timeout: number;
}

// Hook stdin format
interface PreToolHookInput {
  tool_name: string;
  tool_input: unknown;
  session_id: string;
}

// Hook stdout format (must be valid JSON)
interface PreToolHookOutput {
  action: 'allow' | 'deny' | 'modify';
  modified_input?: unknown;  // If action is 'modify'
  reason?: string;
}

// Configuration in settings.json
const hookConfig = {
  hooks: {
    'pre-tool-use': [
      { 
        name: 'security-scan',
        command: '/path/to/security-hook.sh',
        timeout: 5000
      }
    ],
    'post-tool-use': [
      {
        name: 'audit-log',
        command: '/path/to/audit-hook.sh',
        timeout: 2000
      }
    ]
  }
};
```

## Key Files (DXRK Reference)

| File | Lines | Purpose |
|------|-------|---------|
| `permissions.ts` | - | Main permission pipeline |
| `bashPermissions.ts` | - | Bash rule matching |
| `useCanUseTool.tsx` | - | UI for permissions |
| `hooks.ts` | - | Hook system |

## Commands

```bash
# Test permission engine
npm run test:permissions

# Simulate bash rule evaluation
./scripts/test-bash-rules.sh "git push --force"

# Run hook integration test
npm run test:hooks
```

## Resources

- **DXRK Agent Core**: `dxrk-agent-core` skill (comprehensive guide)
- **Patterns Source**: `/home/dxrk/Documentos/DARK GORE/claude-source-learning/`
