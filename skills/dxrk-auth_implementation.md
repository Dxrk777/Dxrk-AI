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

# Authentication & Authorization Skill

## Trigger
"authentication", "authorization", "JWT", "OAuth", "permissions", "RBAC"

## JWT Best Practices
```javascript
// Token structure
{
  sub: "user_id",
  email: "user@example.com",
  roles: ["admin", "user"],
  iat: 1234567890,
  exp: 1234571490 // 1 hour
}

// Verification
const decoded = jwt.verify(token, process.env.JWT_SECRET, {
  algorithms: ['HS256'], // Explicit algorithm
  issuer: 'my-app',     // Issuer validation
  audience: 'my-api'     // Audience validation
})

// NEVER decode without verify
// const decoded = jwt.decode(token) // ❌ Unsafe
```

## Password Security
```javascript
// Hashing with bcrypt
const hash = await bcrypt.hash(password, 12)

// Verification
const match = await bcrypt.compare(password, hash)

// NEVER store plain text passwords
```

## RBAC Pattern
```javascript
// Middleware
const requireRole = (role) => (req, res, next) => {
  if (!req.user?.roles?.includes(role)) {
    return res.status(403).json({ error: 'Forbidden' })
  }
  next()
}

// Usage
app.delete('/users/:id', requireAuth, requireRole('admin'), deleteUser)

// Permission check
const can = (user, action, resource) => {
  const permissions = {
    admin: ['create', 'read', 'update', 'delete'],
    user: ['read', 'update_own'],
  }
  return permissions[user.role]?.includes(action)
}
```

## OAuth 2.0 Flow
1. User clicks "Login with Provider"
2. Redirect to provider's auth endpoint
3. User grants permissions
4. Provider redirects back with code
5. Exchange code for tokens
6. Store tokens securely (httpOnly cookie)
