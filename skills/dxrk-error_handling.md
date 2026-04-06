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
║              DXRK AI INTELLIGENCE SYSTEM — SKILL PACKAGE                      ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
-->

# Error Handling Skill

## Trigger
"error handling", "exception", "try-catch", "error recovery", "retry"

## Error Types
1. **Operational Errors**: Expected, handleable (bad input, not found)
2. **Programmer Errors**: Bugs, should crash and fix (null reference, type error)

## Never Swallow Errors
```javascript
// BAD
try {
  processData(data)
} catch (e) {
  // Silent fail - bad!
}

// GOOD
try {
  processData(data)
} catch (e) {
  logger.error('Failed to process data', { error: e, data })
  throw new AppError('PROCESSING_FAILED', e)
}
```

## Custom Error Class
```javascript
class AppError extends Error {
  constructor(code, message, statusCode = 500, details = {}) {
    super(message)
    this.code = code
    this.statusCode = statusCode
    this.details = details
  }
}

// Usage
throw new AppError('USER_NOT_FOUND', 'User not found', 404)
```

## Retry Pattern
```javascript
async function retry(fn, maxAttempts = 3, delay = 1000) {
  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      return await fn()
    } catch (error) {
      if (attempt === maxAttempts) throw error
      logger.warn(`Attempt ${attempt} failed, retrying...`, { error })
      await sleep(delay * attempt) // Exponential backoff
    }
  }
}
```

## Global Error Handler (Express)
```javascript
app.use((err, req, res, next) => {
  if (err instanceof AppError) {
    return res.status(err.statusCode).json({
      error: { code: err.code, message: err.message, details: err.details }
    })
  }
  
  logger.error('Unexpected error', { error: err })
  res.status(500).json({ error: { code: 'INTERNAL_ERROR' } })
})
```
