---
name: dxrk-master-index
description: >
  Master index of all DXRK skills with COMPLETE auto-trigger mappings.
  Trigger: When needing to find the right skill, skill directory reference.
license: Apache-2.0
metadata:
  author: dxrk
  version: "3.0"
---

## 🔥 Auto-Load Triggers (Complete)

Load skills automatically based on these patterns:

### 🤖 AI Agents & Architecture (PRIORITY)

| Pattern Detected | Skill Loaded |
|-----------------|--------------|
| `agent loop`, `query engine`, `streaming tool`, `async generator`, `while true` | `dxrk-agent-core` |
| `dual layer`, `generator pattern`, `token budget`, `auto compact` | `dxrk-agent-core` |
| `tool permission`, `security`, `bash rule`, `tool control`, `deny rules` | `dxrk-agent-core` |
| `multi-agent`, `coordinator`, `worker`, `fork`, `teammate` | `dxrk-agent-core` |
| `system prompt`, `prompt registry`, `boundary`, `prompt cache`, `override` | `dxrk-agent-core` |
| `built-in agent`, `explore agent`, `plan agent`, `verification agent` | `dxrk-agent-core` |
| `mcp`, `model context protocol`, `mcp server`, `tool discovery` | `dxrk-agent-core` |
| `agent core`, `source code`, `query.ts`, `permissions.ts` | `dxrk-agent-core` |

### 🎯 Architecture Patterns

| Pattern Detected | Skill Loaded |
|-----------------|--------------|
| `permission pipeline`, `12 step`, `always allow`, `always deny` | `dxrk-tool-control-patterns` |
| `prompt registry`, `segment registry`, `prompt composition` | `dxrk-prompt-registry` |
| `self-healing`, `auto-repair`, `error recovery`, `autonomous` | `dxrk-auto-repair-protocol` |
| `llm fallback`, `multi-provider`, `circuit breaker`, `cost` | `dxrk-llm-fallback-strategy` |
| `task notification`, `mailbox`, `async worker`, `fork agent` | `dxrk-multi-agent-coordinator` |

### ☁️ Infrastructure

| Pattern Detected | Skill Loaded |
|-----------------|--------------|
| `deploy`, `fly.io`, `railway`, `render`, `docker` | `dxrk-free-deployment` |
| `supabase`, `redis`, `neon`, `turso`, `database`, `vector` | `dxrk-free-memory-stack` |
| `docker`, `container`, `dockerfile` | `dxrk-docker-development` |
| `kubernetes`, `helm`, `k8s` | `dxrk-helm-chart-builder` |

### 🎨 Frontend

| Pattern Detected | Skill Loaded |
|-----------------|--------------|
| `react`, `component`, `useState`, `useEffect` | `react-19` |
| `angular`, `standalone`, `signal`, `inject` | `scope-rule-architect-angular` |
| `next.js`, `app router`, `server component` | `nextjs-15` |
| `tailwind`, `css`, `styling` | `tailwind-4` |

### ⚙️ Backend

| Pattern Detected | Skill Loaded |
|-----------------|--------------|
| `python`, `fastapi`, `django` | `pytest` |
| `dotnet`, `csharp`, `.net` | `dotnet` |
| `api`, `rest`, `endpoint` | `django-drf` |

### 🧪 Testing

| Pattern Detected | Skill Loaded |
|-----------------|--------------|
| `go test`, `bubbletea`, `tui` | `go-testing` |
| `playwright`, `e2e`, `page object` | `playwright` |
| `pytest`, `python test`, `fixture` | `pytest` |

### 📋 Workflow

| Pattern Detected | Skill Loaded |
|-----------------|--------------|
| `pr`, `pull request`, `branch` | `branch-pr` |
| `issue`, `bug report`, `feature request` | `issue-creation` |
| `triage`, `backlog`, `issues` | `backlog-triage` |

---

## 📚 Full Skill List (123 skills)

### 🏆 DXRK Agent Core (COMPLETE)
- `dxrk-agent-core` — **COMPLETE** DXRK source analysis (ALL patterns)

### 🤖 AI Agents & Architecture
- `dxrk-agent-loop-architect` — Dual-layer AsyncGenerator
- `dxrk-tool-control-patterns` — 12-step security pipeline
- `dxrk-multi-agent-coordinator` — Async workers + mailbox
- `dxrk-prompt-registry` — Modular prompts
- `dxrk-auto-repair-protocol` — 5-phase self-healing
- `dxrk-llm-fallback-strategy` — Multi-provider LLM
- `dxrk-ai-agent-development` — Agent patterns
- `dxrk-mcp-server-builder` — MCP servers
- `dxrk-multi-agent-patterns` — Agent design

### ☁️ Infrastructure
- `dxrk-free-deployment` — Fly.io, Railway, Render
- `dxrk-free-memory-stack` — Supabase, Redis, Neon, Turso
- `dxrk-docker-development` — Docker patterns
- `dxrk-helm-chart-builder` — Kubernetes Helm
- `dxrk-observability-designer` — Monitoring/logging

### 🎨 Frontend
- `react-19` — React 19 patterns
- `dxrk-react-patterns` — React patterns
- `dxrk-react-state-management` — State patterns
- `nextjs-15` — Next.js 15
- `tailwind-4` — Tailwind CSS 4
- `zustand-5` — Zustand state

### ⚙️ Backend
- `dotnet` — .NET 9
- `django-drf` — Django REST Framework
- `ai-sdk-5` — Vercel AI SDK

### 🧪 Testing
- `go-testing` — Go testing
- `playwright` — Playwright E2E
- `pytest` — Python pytest

### 🔐 Security & DevOps
- `dxrk-api-security-testing` — API security
- `dxrk-github-actions-templates` — CI/CD

### 📋 Workflow
- `branch-pr` — PR workflow
- `issue-creation` — Issue workflow
- `backlog-triage` — Issue triage
- `repo-hardening` — Repo setup

### 🗄️ Data & Architecture
- `dxrk-database-designer` — Schema design
- `dxrk-rag-implementation` — RAG systems
- `dxrk-architecture-decision-records` — ADRs
- `dxrk-ddd-strategic-design` — DDD

---

## 📂 Source References (DXRK Agent Core)

```
/home/dxrk/Documentos/DARK GORE/claude-source-learning/
├── index.html           # Overview, architecture diagram
├── agent-loop.html      # Dual-layer AsyncGenerator, streaming tools
├── tool-control.html    # 12-step pipeline, 5 permission modes
├── multi-agent.html     # Coordinator, workers, mailbox
├── system-prompt.html  # 14 segments, BOUNDARY caching
├── built-in-agents.html # Explore, Plan, Verification agents
├── insights.html       # 10 reusable patterns
└── 核心机制.png        # Architecture diagram
```

---

## 🎯 Quick Reference

### Agent Building
```
Everything:    dxrk-agent-core
Loops:         dxrk-agent-loop-architect
Permissions:   dxrk-tool-control-patterns
Multi-Agent:   dxrk-multi-agent-coordinator
Prompts:       dxrk-prompt-registry
Errors:        dxrk-auto-repair-protocol
LLM:          dxrk-llm-fallback-strategy
```

### Free Stack
```
Deploy:        dxrk-free-deployment
Database:      dxrk-free-memory-stack
```

### Common Stacks
```
React:         react-19 + zustand-5 + tailwind-4
Angular:       scope-rule-architect-angular
Next.js:       nextjs-15
Python:        pytest + django-drf
Go:            go-testing
```
