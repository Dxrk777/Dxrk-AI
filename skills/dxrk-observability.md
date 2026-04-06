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

# Observability Skill

## Trigger
"observability", "monitoring", "logging", "tracing", "metrics", "alerting", "SRE"

## Three Pillars

### 1. Logs
- Structured JSON format
- Correlation IDs
- Log levels: ERROR, WARN, INFO, DEBUG
- Include context: userId, requestId, sessionId

```json
{
  "timestamp": "2026-04-05T12:00:00Z",
  "level": "error",
  "message": "Payment failed",
  "context": {
    "userId": "123",
    "orderId": "456",
    "errorCode": "DECLINED"
  }
}
```

### 2. Metrics
- Request rate (RPS)
- Error rate (4xx, 5xx)
- Latency (p50, p95, p99)
- Saturation (CPU, memory, connections)

### 3. Traces
- Distributed tracing
- Span hierarchy
- Service mesh integration

## Health Check Pattern
```javascript
app.get('/health', async (req, res) => {
  const checks = {
    db: await checkDatabase(),
    cache: await checkCache(),
    external: await checkExternal()
  }
  
  const healthy = Object.values(checks).every(c => c.healthy)
  
  res.status(healthy ? 200 : 503).json({
    status: healthy ? 'healthy' : 'unhealthy',
    checks,
    timestamp: new Date().toISOString()
  })
})
```

## Alerting Rules
- Error rate > 1% for 5 minutes
- P99 latency > 2s for 5 minutes
- Availability < 99.9%
