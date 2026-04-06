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

# Testing Strategy Skill

## Trigger
"testing strategy", "test pyramid", "coverage", "TDD", "unit test", "integration test"

## Test Pyramid
```
        /\
       /  \
      / E2E\        Few, slow, high confidence
     /------\
    /Integr. \      Some, medium speed
   /----------\
  /  Unit Tests \   Many, fast, low confidence
 /______________\
```

## Unit Test Pattern
```javascript
describe('calculateTotal', () => {
  it('sums item prices correctly', () => {
    const items = [
      { price: 10, quantity: 2 },
      { price: 5, quantity: 1 }
    ]
    expect(calculateTotal(items)).toBe(25)
  })
  
  it('handles empty array', () => {
    expect(calculateTotal([])).toBe(0)
  })
  
  it('throws for negative prices', () => {
    expect(() => calculateTotal([{ price: -1 }]))
      .toThrow('Price cannot be negative')
  })
})
```

## Integration Test Pattern
```javascript
describe('POST /users', () => {
  it('creates user and returns 201', async () => {
    const res = await request(app)
      .post('/users')
      .send({ email: 'test@example.com', name: 'Test' })
      .expect(201)
    
    expect(res.body).toMatchObject({
      id: expect.any(String),
      email: 'test@example.com'
    })
    
    // Verify in database
    const user = await db.users.findById(res.body.id)
    expect(user).toBeDefined()
  })
})
```

## Coverage Targets
- Unit: 80%+ line coverage
- Integration: Critical paths covered
- E2E: Happy paths + critical user flows
