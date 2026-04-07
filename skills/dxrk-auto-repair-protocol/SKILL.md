---
name: dxrk-auto-repair-protocol
description: >
  Implement autonomous self-healing systems that detect, diagnose, and repair errors.
  Trigger: Error handling, self-healing systems, autonomous repair, error recovery.
license: Apache-2.0
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use

- Building autonomous error recovery systems
- Implementing self-healing applications
- Creating error pattern databases
- Designing diagnostic pipelines
- Building resilient AI agent systems

## 5-Phase Auto-Repair Protocol

```
┌──────────────────────────────────────────────────────────────────┐
│                    AUTO-REPAIR PROTOCOL                           │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  PHASE 1: DETECTION ────────────────────────────────────────    │
│  • Monitor errors in real-time                                   │
│  • Classify: syntax | logic | dependency | state | network      │
│  • Determine severity: CRITICAL / MAJOR / MINOR                  │
│                                                                  │
│  PHASE 2: DIAGNOSTIC ──────────────────────────────────────     │
│  • Reproduce error in isolated environment                       │
│  • Trace full stack                                             │
│  • Identify root cause (not symptom)                             │
│  • Map all affected components                                  │
│                                                                  │
│  PHASE 3: CORRECTION ───────────────────────────────────────     │
│  • Generate minimum 3 solution candidates                        │
│  • Evaluate: impact, risk, speed                                 │
│  • Apply optimal solution                                        │
│  • Document change with explanatory comment                     │
│                                                                  │
│  PHASE 4: VERIFICATION ────────────────────────────────────     │
│  • Run unit tests for component                                 │
│  • Run integration tests                                         │
│  • Confirm no regressions                                       │
│  • If fails → return to PHASE 2 with new hypothesis             │
│                                                                  │
│  PHASE 5: LEARNING ────────────────────────────────────────     │
│  • Store pattern → solution in knowledge base                   │
│  • Update preventive rules                                      │
│  • If critical → alert owner                                    │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

## Error Classification System

```typescript
type ErrorCategory = 
  | 'syntax'        // Compilation, parsing
  | 'logic'         // Business logic bugs
  | 'dependency'    // Missing deps, version conflicts
  | 'state'         // Race conditions, deadlocks
  | 'network'       // Connectivity, timeouts
  | 'permission'    // Auth, access control
  | 'resource'      // Memory, disk, CPU
  | 'external'      // Third-party APIs
  | 'data'          // Validation, corruption
  | 'unknown';      // Unclassified

type Severity = 'critical' | 'major' | 'minor';

interface ErrorClassification {
  category: ErrorCategory;
  severity: Severity;
  retryable: boolean;
  userFacing: boolean;
}

function classifyError(error: Error, context?: ErrorContext): ErrorClassification {
  const errorMessage = error.message.toLowerCase();
  
  // Syntax errors
  if (/syntaxerror|unexpected token|parse error/i.test(errorMessage)) {
    return { 
      category: 'syntax', 
      severity: 'critical', 
      retryable: false,
      userFacing: false 
    };
  }
  
  // Network errors
  if (/timeout|connection|econnrefused|enotfound|network/i.test(errorMessage)) {
    return { 
      category: 'network', 
      severity: 'major', 
      retryable: true,
      userFacing: true 
    };
  }
  
  // Dependency errors
  if (/cannot find module|modulenotfounderror|import.*failed/i.test(errorMessage)) {
    return { 
      category: 'dependency', 
      severity: 'critical', 
      retryable: false,
      userFacing: false 
    };
  }
  
  // Permission errors
  if (/permission denied|eacces|unauthorized|forbidden/i.test(errorMessage)) {
    return { 
      category: 'permission', 
      severity: 'critical', 
      retryable: false,
      userFacing: true 
    };
  }
  
  // Resource errors
  if (/out of memory|heap|eoverflow|enospc|cannot allocate/i.test(errorMessage)) {
    return { 
      category: 'resource', 
      severity: 'critical', 
      retryable: false,
      userFacing: true 
    };
  }
  
  // External API errors
  if (/api|rate limit|quota|4\d{2}|5\d{2}/i.test(errorMessage)) {
    return { 
      category: 'external', 
      severity: 'major', 
      retryable: true,
      userFacing: true 
    };
  }
  
  return { 
    category: 'unknown', 
    severity: 'minor', 
    retryable: false,
    userFacing: false 
  };
}
```

## 5-Attempt Recovery Strategy

```typescript
interface RecoveryAttempt {
  attempt: number;
  strategy: 'direct' | 'alternative' | 'decompose' | 'search' | 'minimal';
  result: 'success' | 'fail' | 'escalate';
  details?: string;
}

async function autoRecover(
  error: Error,
  context: ExecutionContext,
  onAlert?: (alert: Alert) => void
): Promise<RecoveryResult> {
  const attempts: RecoveryAttempt[] = [];
  
  // Attempt 1: Direct fix
  attempts.push(await attemptDirectFix(error, context));
  if (attempts[0].result === 'success') {
    return { success: true, attempts, solution: attempts[0].details };
  }
  
  // Attempt 2: Alternative paradigm
  attempts.push(await attemptAlternative(error, context));
  if (attempts[1].result === 'success') {
    return { success: true, attempts, solution: attempts[1].details };
  }
  
  // Attempt 3: Decompose into sub-problems
  attempts.push(await attemptDecomposition(error, context));
  if (attempts[2].result === 'success') {
    return { success: true, attempts, solution: attempts[2].details };
  }
  
  // Attempt 4: Search for documented solution
  attempts.push(await attemptSearch(error, context));
  if (attempts[3].result === 'success') {
    return { success: true, attempts, solution: attempts[3].details };
  }
  
  // Attempt 5: Minimal viable solution + plan
  attempts.push(await attemptMinimal(error, context));
  if (attempts[4].result === 'success') {
    return { success: true, attempts, solution: attempts[4].details };
  }
  
  // ALL FAILED: Escalate to owner
  await escalateToOwner(error, context, attempts, onAlert);
  
  return { 
    success: false, 
    attempts, 
    escalation: true,
    message: 'All recovery attempts failed. Owner has been notified.'
  };
}

async function attemptDirectFix(error: Error, context: ExecutionContext): Promise<RecoveryAttempt> {
  // Try the most straightforward fix
  // For example: add missing semicolon, fix typo, etc.
  try {
    const fix = analyzeAndFix(error);
    await applyFix(fix);
    return { attempt: 1, strategy: 'direct', result: 'success', details: fix.description };
  } catch (e) {
    return { attempt: 1, strategy: 'direct', result: 'fail' };
  }
}

async function attemptAlternative(error: Error, context: ExecutionContext): Promise<RecoveryAttempt> {
  // Try a different paradigm
  // For example: sync → async, pull → push, etc.
  try {
    const alternative = findAlternativeApproach(error, context);
    await applyAlternative(alternative);
    return { attempt: 2, strategy: 'alternative', result: 'success', details: alternative.description };
  } catch (e) {
    return { attempt: 2, strategy: 'alternative', result: 'fail' };
  }
}
```

## Error Pattern Knowledge Base

```typescript
interface ErrorPattern {
  pattern: RegExp | string;
  category: ErrorCategory;
  solutions: Solution[];
  successRate: number;
  lastUsed: Date;
}

interface Solution {
  description: string;
  code?: string;
  steps: string[];
  estimatedFixTime: number;
}

const errorKnowledgeBase: ErrorPattern[] = [
  {
    pattern: /connection refused|ECONNREFUSED/i,
    category: 'network',
    solutions: [
      {
        description: 'Retry with exponential backoff',
        steps: ['Wait 1s', 'Retry', 'Wait 2s', 'Retry', 'Wait 4s', 'Retry'],
        estimatedFixTime: 30
      },
      {
        description: 'Check if service is running',
        steps: ['Check port', 'Start service if needed', 'Verify firewall'],
        estimatedFixTime: 60
      }
    ],
    successRate: 0.85,
    lastUsed: new Date()
  },
  {
    pattern: /cannot find module '(.*)'|Modulenotfounderror: no module named '(.*)'/i,
    category: 'dependency',
    solutions: [
      {
        description: 'Install missing dependency',
        code: 'npm install {module}',
        steps: ['Identify module', 'Install', 'Verify'],
        estimatedFixTime: 30
      },
      {
        description: 'Check import path',
        steps: ['Verify path', 'Fix relative imports', 'Check package.json exports'],
        estimatedFixTime: 60
      }
    ],
    successRate: 0.92,
    lastUsed: new Date()
  }
];

async function findKnownSolution(error: Error): Promise<Solution | null> {
  for (const pattern of errorKnowledgeBase) {
    if (typeof pattern.pattern === 'string') {
      if (error.message.includes(pattern.pattern)) {
        // Update usage stats
        pattern.lastUsed = new Date();
        return pattern.solutions[0];
      }
    } else {
      if (pattern.pattern.test(error.message)) {
        pattern.lastUsed = new Date();
        return pattern.solutions[0];
      }
    }
  }
  
  return null;
}

async function learnNewPattern(error: Error, solution: Solution): Promise<void> {
  const pattern: ErrorPattern = {
    pattern: extractPattern(error.message),
    category: classifyError(error).category,
    solutions: [solution],
    successRate: 1.0,
    lastUsed: new Date()
  };
  
  errorKnowledgeBase.push(pattern);
  await saveToKnowledgeBase(pattern);
}
```

## Monitoring & Alerting

```typescript
interface Alert {
  severity: 'info' | 'warning' | 'critical';
  title: string;
  message: string;
  context: Record<string, unknown>;
  timestamp: Date;
  retryable: boolean;
}

class AutoRepairMonitor {
  private metrics = {
    totalErrors: 0,
    recovered: 0,
    escalated: 0,
    byCategory: new Map<ErrorCategory, { total: number; recovered: number }>()
  };
  
  async recordError(error: Error, classification: ErrorClassification, result: RecoveryResult): Promise<void> {
    this.metrics.totalErrors++;
    
    const categoryMetrics = this.metrics.byCategory.get(classification.category) 
      || { total: 0, recovered: 0 };
    categoryMetrics.total++;
    this.metrics.byCategory.set(classification.category, categoryMetrics);
    
    if (result.success) {
      this.metrics.recovered++;
      categoryMetrics.recovered++;
    } else {
      this.metrics.escalated++;
    }
    
    await this.persistMetrics();
  }
  
  getHealthScore(): number {
    if (this.metrics.totalErrors === 0) return 100;
    return Math.round((this.metrics.recovered / this.metrics.totalErrors) * 100);
  }
  
  getReport(): string {
    return `
Auto-Repair Health Report
========================
Total Errors: ${this.metrics.totalErrors}
Recovered: ${this.metrics.recovered}
Escalated: ${this.metrics.escalated}
Health Score: ${this.getHealthScore()}%

By Category:
${Array.from(this.metrics.byCategory.entries())
  .map(([cat, m]) => `  ${cat}: ${m.recovered}/${m.total}`)
  .join('\n')}
    `.trim();
  }
}
```

## Commands

```bash
# Run self-healing test suite
npm run test:auto-repair

# Simulate error scenarios
npm run simulate:errors

# View repair history
npm run repair:history

# Reset knowledge base
npm run repair:reset-kb

# Generate health report
npm run repair:report
```

## Integration with Engram Memory

Store repair patterns in persistent memory:

```typescript
import { engram } from '@/lib/memory';

async function saveRepairPattern(
  error: Error,
  solution: Solution,
  success: boolean
): Promise<void> {
  await engram.save({
    type: 'repair',
    title: `Fixed: ${error.message.slice(0, 50)}`,
    content: {
      errorPattern: error.message,
      category: classifyError(error).category,
      solution: solution,
      success: success,
      timestamp: new Date().toISOString()
    }
  });
}

async function loadRepairPatterns(category?: ErrorCategory): Promise<ErrorPattern[]> {
  const patterns = await engram.search({
    type: 'repair',
    category: category
  });
  
  return patterns.map(p => ({
    pattern: p.content.errorPattern,
    category: p.content.category,
    solution: p.content.solution
  }));
}
```

## Resources

- **Inspired by**: DXRK's 3 error recovery paths
- **Pattern Base**: `/home/dxrk/Documentos/DARK GORE/Dxrk.md` (Auto-Repair Module)
