---
name: dxrk-free-memory-stack
description: >
  Set up free persistent storage: Supabase, Redis, Neon, Turso, MongoDB Atlas.
  Trigger: Database setup, caching, vector storage, persistent memory storage.
license: Apache-2.0
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use

- Setting up free persistent databases
- Implementing session storage
- Adding vector search capabilities
- Creating cache layers
- Building memory systems for AI agents

## Free Tier Comparison

| Service | Type | Free Tier | Persistence | Best For |
|---------|------|-----------|-------------|----------|
| **Supabase** | PostgreSQL + pgvector | 500 MB | Permanent | Full-featured SQL + vectors |
| **Neon** | Serverless PostgreSQL | 512 MB | Permanent | Auto-scaling Postgres |
| **Turso** | Distributed SQLite | 500 DBs, 9 GB | Permanent | Edge databases |
| **Upstash Redis** | Redis | 10K commands/day | Permanent | Sessions, caching |
| **MongoDB Atlas** | NoSQL | 512 MB | Permanent | Flexible schemas |
| **Cloudflare KV** | Key-Value | 100K writes/day | Permanent | Global key-value |
| **PlanetScale** | MySQL | 5 GB | Permanent | Serverless MySQL |

## Supabase (RECOMMENDED)

Full PostgreSQL + pgvector for semantic search.

### Setup

```bash
# 1. Create account at https://supabase.com
# 2. Create new project
# 3. Get connection string from Settings > Database

# Using Supabase CLI (optional)
npx supabase init
npx supabase start
```

### Connection

```typescript
import { Pool } from 'pg';

const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
  ssl: { rejectUnauthorized: false }
});

// Or use Supabase client
import { createClient } from '@supabase/supabase-js';

const supabase = createClient(
  process.env.SUPABASE_URL!,
  process.env.SUPABASE_ANON_KEY!
);
```

### Schema for Agent Memory

```sql
-- Episodic memory (conversations)
CREATE TABLE memory_episodes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id TEXT NOT NULL,
  role TEXT CHECK (role IN ('user', 'assistant', 'system')),
  content TEXT NOT NULL,
  embedding VECTOR(1536),
  importance FLOAT DEFAULT 0.5,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Knowledge base (learned facts)
CREATE TABLE knowledge_base (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  topic TEXT NOT NULL,
  subtopic TEXT,
  content TEXT NOT NULL,
  source_url TEXT,
  confidence FLOAT DEFAULT 0.9,
  language TEXT,
  tags TEXT[],
  last_updated TIMESTAMPTZ DEFAULT NOW()
);

-- Error patterns and solutions
CREATE TABLE error_patterns (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  error_pattern TEXT NOT NULL,
  error_type TEXT,
  solution TEXT NOT NULL,
  language TEXT,
  framework TEXT,
  success_count INT DEFAULT 1,
  fail_count INT DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Projects and context
CREATE TABLE projects (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  description TEXT,
  stack TEXT[],
  git_url TEXT,
  notes TEXT,
  active BOOLEAN DEFAULT true,
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Vector indexes for semantic search
CREATE INDEX ON memory_episodes 
  USING ivfflat (embedding vector_cosine_ops)
  WITH (lists = 100);

CREATE INDEX ON knowledge_base 
  USING ivfflat (embedding vector_cosine_ops)
  WITH (lists = 100);
```

### Usage Examples

```typescript
// Store conversation
async function storeMessage(sessionId: string, role: string, content: string) {
  const embedding = await generateEmbedding(content);
  
  await supabase.from('memory_episodes').insert({
    session_id: sessionId,
    role,
    content,
    embedding
  });
}

// Semantic search
async function searchMemory(query: string, limit = 5) {
  const queryEmbedding = await generateEmbedding(query);
  
  const { data, error } = await supabase.rpc('match_memory_episodes', {
    query_embedding: queryEmbedding,
    match_threshold: 0.7,
    match_count: limit
  });
  
  return data;
}

// SQL function for semantic search
const matchMemorySQL = `
CREATE OR REPLACE FUNCTION match_memory_episodes(
  query_embedding VECTOR(1536),
  match_threshold FLOAT,
  match_count INT
)
RETURNS TABLE (
  id UUID,
  content TEXT,
  similarity FLOAT
)
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
  SELECT
    me.id,
    me.content,
    1 - (me.embedding <=> query_embedding) AS similarity
  FROM memory_episodes me
  WHERE 1 - (me.embedding <=> query_embedding) > match_threshold
  ORDER BY me.embedding <=> query_embedding
  LIMIT match_count;
END;
$$;
`;
```

## Upstash Redis (Caching)

Great for sessions, rate limiting, and fast caching.

### Setup

```bash
# 1. Create account at https://upstash.com
# 2. Create new Redis database
# 3. Copy REST API URL and token
```

### TypeScript Client

```typescript
import { Redis } from '@upstash/redis';

const redis = new Redis({
  url: process.env.UPSTASH_REDIS_REST_URL!,
  token: process.env.UPSTASH_REDIS_REST_TOKEN!
});
```

### Usage Patterns

```typescript
// Session storage
async function setSession(userId: string, data: object) {
  await redis.set(`session:${userId}`, JSON.stringify(data), {
    ex: 60 * 60 * 24 * 7  // 7 days
  });
}

async function getSession(userId: string) {
  const data = await redis.get(`session:${userId}`);
  return data ? JSON.parse(data as string) : null;
}

// Rate limiting (token bucket)
async function checkRateLimit(userId: string, limit: number, window: number) {
  const key = `ratelimit:${userId}`;
  
  const count = await redis.incr(key);
  if (count === 1) {
    await redis.expire(key, window);
  }
  
  return count <= limit;
}

// Distributed locks
async function acquireLock(key: string, ttl: number): Promise<boolean> {
  const result = await redis.set(key, '1', {
    nx: true,  // Only set if not exists
    ex: ttl    // Expiration
  });
  
  return result === 'OK';
}

// Cache with invalidation
async function cacheWithInvalidation<T>(
  key: string,
  fetcher: () => Promise<T>,
  ttl: number
): Promise<T> {
  const cached = await redis.get(key);
  
  if (cached) {
    return JSON.parse(cached as string) as T;
  }
  
  const data = await fetcher();
  await redis.set(key, JSON.stringify(data), { ex: ttl });
  
  return data;
}
```

## Neon (Serverless PostgreSQL)

Auto-scales to zero, no cold starts.

```typescript
import { neon } from '@neondatabase/serverless';

const sql = neon(process.env.DATABASE_URL!);

// Use with Drizzle ORM
import { drizzle } from 'drizzle-orm/neon-http';
import { pgTable, serial, text, timestamp } from 'drizzle-orm/pg-core';

const users = pgTable('users', {
  id: serial('id').primaryKey(),
  name: text('name').notNull(),
  email: text('email').notNull(),
  createdAt: timestamp('created_at').defaultNow()
});

const db = drizzle(sql);

// Query
const allUsers = await db.select().from(users);
```

## Turso (Distributed SQLite)

Perfect for edge deployments and multi-region.

```typescript
import { createClient } from '@libsql/client';

const db = createClient({
  url: 'libsql://my-db.turso.io',
  authToken: process.env.TURSO_AUTH_TOKEN
});

// Create tables
await db.execute(`
  CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  )
`);

// Insert
await db.execute({
  sql: 'INSERT INTO messages (session_id, role, content) VALUES (?, ?, ?)',
  args: [sessionId, role, content]
});

// Query
const messages = await db.execute({
  sql: 'SELECT * FROM messages WHERE session_id = ? ORDER BY created_at',
  args: [sessionId]
});
```

## MongoDB Atlas (NoSQL)

Flexible schemas for document storage.

```typescript
import { MongoClient } from 'mongodb';

const client = new MongoClient(process.env.MONGODB_URI!);
const db = client.db('myapp');

// Collections
const sessions = db.collection('sessions');
const knowledge = db.collection('knowledge');

// Store session
await sessions.insertOne({
  sessionId,
  messages: [],
  metadata: {},
  createdAt: new Date(),
  updatedAt: new Date()
});

// Semantic search with Atlas Search
await db.collection('knowledge').createIndex({
  description: 'text',
  content: 'text'
}, { name: 'knowledge_text_index' });
```

## Keep-Alive Cron Job

Prevent free tiers from sleeping:

```bash
# Use https://cron-job.org (free)
# Set up cron job to hit your endpoint every 5 days
# Example: GET https://your-app.fly.dev/health
```

## Commands

```bash
# Supabase CLI
npx supabase init
npx supabase start
npx supabase db reset

# Neon CLI
npm install -g neonctl
neonctl auth login
neonctl connection-string --database-id <id>

# Turso CLI
brew install tursodotdev/tap/turso
turso auth login
turso db create mydb
turso db shell mydb
```

## Schema Selection Guide

| Use Case | Recommended |
|----------|-------------|
| Full agent memory with vectors | **Supabase** (pgvector) |
| Fast session/caching | **Upstash Redis** |
| Auto-scaling SQL | **Neon** |
| Edge SQLite | **Turso** |
| Flexible documents | **MongoDB Atlas** |
| Simple key-value | **Cloudflare KV** |
