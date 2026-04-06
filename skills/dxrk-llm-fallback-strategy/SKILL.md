---
name: dxrk-llm-fallback-strategy
description: >
  Implement multi-provider LLM strategies with automatic fallback and cost optimization.
  Trigger: LLM integration, API fallbacks, multi-model strategies, cost optimization.
license: Apache-2.0
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use

- Building production LLM integrations
- Implementing fallback systems
- Optimizing LLM costs
- Handling API rate limits
- Building multi-model pipelines

## Free LLM Providers (2024-2026)

### Tier 1: Unlimited Free (No Daily Limits)

| Provider | Models | API Endpoint | Best For |
|----------|--------|--------------|----------|
| **ZhipuAI (GLM)** | glm-4.7-flash, glm-4.5-flash | `api.mlops-community.com` | Primary model |
| **Groq** | llama-3.1-70b, mixtral-8x7b | `api.groq.com` | Fast inference |

### Tier 2: Free with Limits (200 req/day)

| Provider | Models | API Endpoint | Best For |
|----------|--------|--------------|----------|
| **OpenRouter** | moonshotai/kimi-k2.5, qwen/qwen3.6-plus | `openrouter.ai` | Model variety |
| **Together AI** | llama-3-70b, mixtral-8x7b | `api.together.xyz` | Open models |

### Tier 3: Free Credits

| Provider | Credits | Renewal |
|----------|---------|---------|
| **GroqCloud** | $0 free + usage | Permanent |
| **Cohere** | Free tier | Permanent |
| **Dxrk** | $5 free | New accounts |

## Multi-Provider Fallback Client

```typescript
interface LLMConfig {
  provider: 'zhipu' | 'groq' | 'openrouter' | 'openai' | 'anthropic';
  model: string;
  apiKey: string;
  baseUrl?: string;
  maxTokens: number;
  temperature: number;
}

interface LLMRequest {
  messages: Message[];
  systemPrompt?: string;
  temperature?: number;
  maxTokens?: number;
}

interface LLMResponse {
  content: string;
  model: string;
  usage: {
    promptTokens: number;
    completionTokens: number;
    totalTokens: number;
  };
  latency: number;
  provider: string;
}

class MultiProviderLLM {
  private providers: Map<string, LLMConfig> = new Map();
  private circuitBreakers: Map<string, CircuitBreaker> = new Map();
  
  constructor(configs: LLMConfig[]) {
    for (const config of configs) {
      this.providers.set(config.provider, config);
      this.circuitBreakers.set(
        config.provider,
        new CircuitBreaker({ failureThreshold: 5, resetTimeout: 60000 })
      );
    }
  }
  
  async complete(request: LLMRequest): Promise<LLMResponse> {
    const providerOrder = this.getProviderOrder(request);
    
    for (const provider of providerOrder) {
      const circuitBreaker = this.circuitBreakers.get(provider)!;
      
      if (circuitBreaker.isOpen()) {
        console.log(`Circuit breaker open for ${provider}, skipping`);
        continue;
      }
      
      try {
        const response = await this.callProvider(provider, request);
        circuitBreaker.recordSuccess();
        return response;
      } catch (error) {
        circuitBreaker.recordFailure();
        console.error(`Provider ${provider} failed:`, error);
        
        // Check if error is retryable
        if (!this.isRetryable(error)) {
          throw error; // Non-retryable error
        }
      }
    }
    
    throw new Error('All LLM providers failed');
  }
  
  private getProviderOrder(request: LLMRequest): string[] {
    // Strategy: Prefer free/unlimited, fallback to limited
    const priorities = {
      'zhipu': 1,      // Unlimited free
      'groq': 2,       // Free, very fast
      'openrouter': 3, // Free with limits
      'openai': 4,     // Paid
      'anthropic': 5   // Paid, best quality
    };
    
    return Array.from(this.providers.keys())
      .sort((a, b) => (priorities[a] || 99) - (priorities[b] || 99));
  }
}
```

## Circuit Breaker Pattern

```typescript
class CircuitBreaker {
  private state: 'closed' | 'open' | 'half-open' = 'closed';
  private failures = 0;
  private lastFailureTime = 0;
  
  constructor(private options: {
    failureThreshold: number;
    resetTimeout: number;
    halfOpenSuccessThreshold?: number;
  }) {}
  
  recordSuccess(): void {
    this.failures = 0;
    this.state = 'closed';
  }
  
  recordFailure(): void {
    this.failures++;
    this.lastFailureTime = Date.now();
    
    if (this.failures >= this.options.failureThreshold) {
      this.state = 'open';
    }
  }
  
  isOpen(): boolean {
    if (this.state === 'open') {
      // Check if reset timeout has passed
      if (Date.now() - this.lastFailureTime > this.options.resetTimeout) {
        this.state = 'half-open';
        return false;
      }
      return true;
    }
    return false;
  }
}
```

## Retry with Exponential Backoff

```typescript
interface RetryConfig {
  maxRetries: number;
  baseDelay: number;
  maxDelay: number;
  backoffMultiplier: number;
  retryableErrors?: RegExp[];
}

async function withRetry<T>(
  fn: () => Promise<T>,
  config: RetryConfig
): Promise<T> {
  let lastError: Error;
  
  for (let attempt = 0; attempt <= config.maxRetries; attempt++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error as Error;
      
      if (attempt === config.maxRetries) break;
      
      // Check if error is retryable
      if (config.retryableErrors && !config.retryableErrors.some(r => r.test(lastError.message))) {
        throw lastError;
      }
      
      // Calculate delay with exponential backoff + jitter
      const delay = Math.min(
        config.baseDelay * Math.pow(config.backoffMultiplier, attempt),
        config.maxDelay
      );
      const jitter = delay * 0.1 * Math.random();
      
      console.log(`Retry ${attempt + 1}/${config.maxRetries} after ${delay}ms`);
      await sleep(delay + jitter);
    }
  }
  
  throw lastError!;
}

// Usage
const response = await withRetry(
  () => llm.complete(request),
  {
    maxRetries: 3,
    baseDelay: 1000,
    maxDelay: 30000,
    backoffMultiplier: 2,
    retryableErrors: [
      /timeout/i,
      /rate.limit/i,
      /429/i,
      /500/i,
      /502/i,
      /503/i
    ]
  }
);
```

## Token Budget Management

```typescript
interface TokenBudget {
  daily: number;
  monthly: number;
  perRequest: number;
}

class TokenBudgetManager {
  private usage = {
    daily: { amount: 0, resetAt: this.getMidnight() },
    monthly: { amount: 0, resetAt: this.getMonthEnd() }
  };
  
  constructor(private budgets: TokenBudget) {}
  
  async checkAndConsume(tokens: number): Promise<boolean> {
    this.resetIfNeeded();
    
    if (
      this.usage.daily.amount + tokens > this.budgets.daily ||
      this.usage.monthly.amount + tokens > this.budgets.monthly
    ) {
      return false; // Budget exceeded
    }
    
    this.usage.daily.amount += tokens;
    this.usage.monthly.amount += tokens;
    
    return true;
  }
  
  getRemaining(): { daily: number; monthly: number } {
    this.resetIfNeeded();
    
    return {
      daily: this.budgets.daily - this.usage.daily.amount,
      monthly: this.budgets.monthly - this.usage.monthly.amount
    };
  }
  
  private resetIfNeeded(): void {
    const now = Date.now();
    
    if (now >= this.usage.daily.resetAt) {
      this.usage.daily = { amount: 0, resetAt: this.getMidnight() };
    }
    
    if (now >= this.usage.monthly.resetAt) {
      this.usage.monthly = { amount: 0, resetAt: this.getMonthEnd() };
    }
  }
}
```

## Model Selection by Task

```typescript
interface ModelSelector {
  select(task: TaskType, requirements: Requirements): LLMConfig;
}

type TaskType = 
  | 'chat' 
  | 'code' 
  | 'reasoning' 
  | 'fast' 
  | 'creative' 
  | 'analysis';

interface Requirements {
  quality: 'low' | 'medium' | 'high';
  latency: 'fast' | 'normal' | 'slow';
  cost: 'minimize' | 'balanced' | 'maximize';
}

const modelMatrix: Record<TaskType, Record<string, string>> = {
  chat: {
    fast: 'zhipu/glm-4.7-flash',
    balanced: 'groq/llama-3.1-70b-versatile',
    quality: 'openai/gpt-4o'
  },
  code: {
    fast: 'groq/llama-3.1-70b-versatile',
    balanced: 'openrouter/qwen/qwen3.6-plus',
    quality: 'anthropic/claude-3.5-sonnet'
  },
  reasoning: {
    fast: 'groq/mixtral-8x7b-32768',
    balanced: 'groq/llama-3.1-70b-versatile',
    quality: 'openai/gpt-4o'
  }
};

function selectModel(
  task: TaskType,
  requirements: Requirements
): LLMConfig {
  const key = requirements.quality === 'low' ? 'fast' 
    : requirements.quality === 'high' ? 'quality' 
    : 'balanced';
  
  const modelName = modelMatrix[task][key];
  
  // Return configured provider for this model
  return getProviderForModel(modelName);
}
```

## Cost Tracking & Alerts

```typescript
interface CostAlert {
  threshold: number;
  percentage: number;
  notified: boolean;
}

class CostTracker {
  private dailyCost = 0;
  private monthlyCost = 0;
  private alerts: CostAlert[] = [];
  
  constructor(
    private dailyBudget: number,
    private monthlyBudget: number
  ) {
    this.setupAlerts();
  }
  
  private setupAlerts(): void {
    this.alerts = [
      { threshold: this.dailyBudget * 0.5, percentage: 50, notified: false },
      { threshold: this.dailyBudget * 0.8, percentage: 80, notified: false },
      { threshold: this.dailyBudget * 0.95, percentage: 95, notified: false },
      { threshold: this.monthlyBudget * 0.5, percentage: 50, notified: false },
      { threshold: this.monthlyBudget * 0.8, percentage: 80, notified: false }
    ];
  }
  
  recordCost(cost: number): void {
    this.dailyCost += cost;
    this.monthlyCost += cost;
    this.checkAlerts();
  }
  
  private checkAlerts(): void {
    for (const alert of this.alerts) {
      if (
        !alert.notified &&
        (this.dailyCost >= alert.threshold || this.monthlyCost >= alert.threshold)
      ) {
        alert.notified = true;
        this.notify(`${alert.percentage}% budget used`);
      }
    }
  }
  
  private notify(message: string): void {
    console.warn(`[COST ALERT] ${message}`);
    // Send to alerting system (Slack, email, etc.)
  }
}
```

## Streaming Responses

```typescript
async function* streamComplete(
  request: LLMRequest
): AsyncGenerator<string> {
  const response = await fetch(CHAT_COMPLETION_URL, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${API_KEY}`
    },
    body: JSON.stringify({
      model: request.model,
      messages: request.messages,
      stream: true,
      temperature: request.temperature,
      max_tokens: request.maxTokens
    })
  });
  
  const reader = response.body?.getReader();
  if (!reader) throw new Error('No response body');
  
  const decoder = new TextDecoder();
  
  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    
    const chunk = decoder.decode(value);
    const lines = chunk.split('\n');
    
    for (const line of lines) {
      if (line.startsWith('data: ')) {
        const data = JSON.parse(line.slice(6));
        if (data.choices?.[0]?.delta?.content) {
          yield data.choices[0].delta.content;
        }
      }
    }
  }
}
```

## Commands

```bash
# Test multi-provider fallback
npm run test:llm-fallback

# Benchmark models
npm run benchmark:models

# Check cost usage
npm run llm:cost-report

# Test rate limits
npm run test:rate-limits

# Simulate provider failures
npm run simulate:llm-failure
```

## Resources

- **ZhipuAI**: https://www.zhipuai.cn/
- **Groq**: https://console.groq.com/
- **OpenRouter**: https://openrouter.ai/
- **Together AI**: https://api.together.xyz/
