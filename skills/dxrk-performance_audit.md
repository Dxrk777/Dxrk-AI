<!--
╔═══════════════════════════════════════════════════════════════════════════════╗
║                                                                               ║
║     ██████╗ ██╗  ██╗██████╗ ██╗  ██╗     ██████╗ ███╗   ███╗███████╗       ║
║     ██╔══██╗╚██╗██╔╝██╔══██╗██║ ██╔╝    ██╔═══██╗████╗ ████║██╔════╝       ║
║     ██║  ██║ ╚███╔╝ ██████╔╝█████╔╝     ██║   ██║██╔████╔██║█████╗         ║
║     ██║  ██║ ██╔██╗ ██╔══██╗██╔═██╗     ██║   ██║██║╚██╔╝██║██╔══╝         ║
║     ██████╔╝██╔╝ ██╗██║  ██║██║  ██╗    ╚██████╔╝██║ ╚═╝ ██║███████╗       ║
║     ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝     ╚═════╝ ╚═╝     ╚═╝╚══════╝       ║
║                                                                               ║
║              DXRK HEX INTELLIGENCE SYSTEM — SKILL PACKAGE                      ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
-->

# Performance Audit Skill

## Trigger
"performance", "optimize", "slow", "bottleneck", "profiling", "latency"

## Purpose
Identify and fix performance issues in codebases.

## Analysis Areas
- **Database**: N+1 queries, missing indexes, inefficient JOINs
- **API**: Response times, pagination, caching opportunities
- **Frontend**: Bundle size, render optimization, lazy loading
- **Memory**: Leaks, allocations, garbage collection pressure
- **Network**: Request batching, compression, CDN opportunities

## Common Patterns to Find
```javascript
// BAD: N+1
users.forEach(u => u.posts = db.find('posts', {userId: u.id}))

// GOOD: Batch load
const userIds = users.map(u => u.id)
const posts = db.find('posts', {userId: { $in: userIds }})
```

## Output Format
```
## Performance Audit

### Critical Issues
| Issue | Location | Impact | Fix |

### Medium Issues
| Issue | Location | Impact | Fix |

### Quick Wins
| Change | Impact |
```
