---
name: dxrk-multi-agent-coordinator
description: >
  Implement multi-agent coordination using async workers and structured message passing.
  Trigger: Multi-agent systems, parallel task execution, agent communication.
license: Apache-2.0
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use

- Building systems that need parallel agent execution
- Implementing Coordinator pattern for task orchestration
- Creating agent-to-agent communication channels
- Managing worker agent lifecycle
- Designing task delegation systems

## Coordinator Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                    COORDINATOR PATTERN                            │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│                    ┌─────────────────┐                          │
│                    │   COORDINATOR   │                          │
│                    │   (Main Agent)  │                          │
│                    └────────┬────────┘                          │
│                             │                                    │
│              ┌──────────────┼──────────────┐                    │
│              │              │              │                    │
│              ▼              ▼              ▼                    │
│        ┌─────────┐  ┌─────────┐  ┌─────────┐                   │
│        │ WORKER 1│  │ WORKER 2│  │ WORKER 3│                   │
│        │(Async)  │  │(Async)  │  │(Async)  │                   │
│        └────┬────┘  └────┬────┘  └────┬────┘                   │
│             │            │            │                         │
│             └────────────┬┴────────────┘                         │
│                          │                                      │
│                    task-notification                             │
│                    (XML/JSON injection)                         │
│                          │                                      │
│                    ┌─────┴─────┐                                │
│                    │  Mailbox  │                                │
│                    │  (Files)  │                                │
│                    └───────────┘                                │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

## Key Insight: Don't Wait, Fork and Inject

**The Problem**: LLM conversation is synchronous by nature. How do you parallelize?

**The Solution**: Fork workers immediately, return task-ids, inject results later.

```typescript
interface TaskNotification {
  task_id: string;
  status: 'pending' | 'completed' | 'failed' | 'killed';
  summary: string;
  result?: unknown;
  error?: string;
}

interface WorkerResult {
  task_id: string;
  content: string;  // Will be injected as user message
}

// Main agent spawns workers and continues
class CoordinatorAgent {
  async coordinate(task: Task): Promise<Result> {
    // Spawn workers IMMEDIATELY, don't wait
    const workers = task.subtasks.map(subtask => 
      this.spawnWorker(subtask)
    );
    
    // Continue working while workers run
    const partialResult = await this.continueMainTask(task);
    
    // Wait for workers and inject results
    const notifications = await Promise.all(
      workers.map(w => w.notification)
    );
    
    // Inject worker results as user messages
    const workerMessages = notifications.map(n => 
      this.buildTaskNotificationMessage(n)
    );
    
    // Final synthesis with all results
    return this.synthesize(task, partialResult, workerMessages);
  }
  
  private buildTaskNotificationMessage(n: TaskNotification): Message {
    return {
      role: 'user',
      content: `<task-notification>
        <task_id>${n.task_id}</task_id>
        <status>${n.status}</status>
        <summary>${n.summary}</summary>
        ${n.result ? `<result>${JSON.stringify(n.result)}</result>` : ''}
        ${n.error ? `<error>${n.error}</error>` : ''}
      </task-notification>`
    };
  }
}
```

## Mailbox Communication Pattern

For agents that need to communicate while running:

```typescript
interface MailboxMessage {
  from: string;
  to: string;
  type: 'command' | 'status' | 'result' | 'error';
  payload: unknown;
  timestamp: number;
}

class MailboxFileSystem {
  private basePath: string;
  
  constructor(basePath: string) {
    this.basePath = basePath;
  }
  
  // Write message to recipient's mailbox
  async write(to: string, message: MailboxMessage): Promise<void> {
    const mailboxPath = `${this.basePath}/${to}/inbox`;
    await fs.mkdir(mailboxPath, { recursive: true });
    
    const filename = `${message.timestamp}-${message.from}.json`;
    await fs.writeFile(
      `${mailboxPath}/${filename}`,
      JSON.stringify(message)
    );
  }
  
  // Read all messages from own mailbox
  async readMailbox(agentId: string): Promise<MailboxMessage[]> {
    const mailboxPath = `${this.basePath}/${agentId}/inbox`;
    
    if (!await fs.exists(mailboxPath)) return [];
    
    const files = await fs.readdir(mailboxPath);
    const messages = await Promise.all(
      files
        .filter(f => f.endsWith('.json'))
        .map(f => fs.readFile(`${mailboxPath}/${f}`, 'utf-8'))
    );
    
    return messages.map(m => JSON.parse(m));
  }
  
  // Clear processed messages
  async clearMailbox(agentId: string): Promise<void> {
    const mailboxPath = `${this.basePath}/${agentId}/inbox`;
    await fs.rm(mailboxPath, { recursive: true });
  }
}

// SendMessageTool for workers
class SendMessageTool {
  name = 'send_message';
  description = 'Send a message to another agent';
  
  async execute(input: { to: string; type: string; payload: unknown }): Promise<void> {
    const message: MailboxMessage = {
      from: this.currentAgentId,
      to: input.to,
      type: input.type as MailboxMessage['type'],
      payload: input.payload,
      timestamp: Date.now()
    };
    
    await mailbox.write(input.to, message);
  }
}
```

## Task State Machine

```
    ┌─────────┐
    │ PENDING │
    └────┬────┘
         │ start()
         ▼
    ┌─────────┐
    │ RUNNING │◄──────────────┐
    └────┬────┘               │
         │                   │ resume()
         │ complete()        │
         ▼                   │
    ┌──────────┐              │
    │COMPLETED │──────────────┘
    └──────────┘  restart()
         ▲
         │ fail()
         │
    ┌─────────┐
    │ FAILED  │
    └─────────┘
         ▲
         │ kill()
         │
    ┌─────────┐
    │ KILLED  │
    └─────────┘
```

## Agent Inheritance Pattern

Child agents inherit parent's capabilities:

```typescript
interface AgentConfig {
  id: string;
  name: string;
  prompt: string;
  tools: Tool[];
  
  // MCP configuration
  mcpServers?: string[];
  inheritParentMCP?: boolean;  // KEY: inherit parent's MCP clients
  
  // Coordinator settings
  isCoordinator?: boolean;
  workerLimit?: number;
}

class AgentHierarchy {
  createWorker(
    parentConfig: AgentConfig,
    task: SubTask
  ): AgentConfig {
    return {
      id: `worker-${task.id}`,
      name: task.name,
      prompt: task.prompt,
      tools: parentConfig.tools,  // Inherit tools
      
      // Inherit MCP if configured
      mcpServers: parentConfig.inheritParentMCP 
        ? parentConfig.mcpServers 
        : task.mcpServers,
      
      inheritParentMCP: parentConfig.inheritParentMCP ?? false
    };
  }
}
```

## MCP Dynamic Registration

```typescript
interface MCPTool {
  name: string;  // Will be prefixed: mcp__serverName__toolName
  description: string;
  inputSchema: JSONSchema;
}

interface MCPServerStatus {
  server: string;
  status: 'connected' | 'pending' | 'failed' | 'needs-auth' | 'disabled';
  tools: MCPTool[];
}

// MCP client implementation
class MCPClient {
  private servers: Map<string, MCPClientConnection> = new Map();
  
  async connect(serverConfig: MCPServerConfig): Promise<void> {
    // Establish connection
    const connection = new MCPClientConnection(serverConfig);
    await connection.connect();
    
    // Discover tools via ListTools RPC
    const tools = await connection.listTools();
    
    // Memoize and namespace tools
    const namespacedTools = tools.map(t => ({
      ...t,
      name: `mcp__${serverConfig.name}__${t.name}`
    }));
    
    this.servers.set(serverConfig.name, connection);
    this.registerTools(namespacedTools);
  }
  
  async disconnect(serverName: string): Promise<void> {
    const connection = this.servers.get(serverName);
    if (connection) {
      await connection.disconnect();
      this.servers.delete(serverName);
      this.unregisterTools(`mcp__${serverName}__`);
    }
  }
}
```

## Key Files (DXRK Reference)

| File | Lines | Purpose |
|------|-------|---------|
| `AgentTool.tsx` | - | Worker spawning |
| `runAgent.ts` | 200+ | Agent execution |
| `coordinatorMode.ts` | 300+ | Coordinator prompt/logic |
| `LocalAgentTask.tsx` | - | Task lifecycle |
| `teammateMailbox.ts` | 200+ | Mailbox filesystem |
| `mcp/client.ts` | 3348 | MCP discovery & registration |

## Commands

```bash
# Test multi-agent coordination
npm run test:multi-agent

# Simulate worker spawning
npm run simulate:workers

# Test mailbox communication
npm run test:mailbox
```

## Resources

- **Multi-Agent Source**: `/home/dxrk/Documentos/DARK GORE/claude-source-learning/multi-agent.html`
- **Built-in Agents**: `/home/dxrk/Documentos/DARK GORE/claude-source-learning/built-in-agents.html`
