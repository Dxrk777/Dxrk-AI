---
name: security
description: >
  Security best practices for code. Trigger: When handling authentication, encryption, secrets, or security-sensitive code.
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use
- Implementing authentication/authorization
- Handling secrets and credentials
- Encrypting/decrypting data
- Input validation and sanitization
- API security (CORS, rate limiting, CSRF)

## Critical Patterns

### Secrets management (REQUIRED)
```bash
# ❌ NEVER commit secrets
API_KEY=sk-1234567890

# ✅ Use environment variables or vault
export API_KEY=$(vault kv get -field=key secret/app)

# ✅ .gitignore
.env
*.pem
*.key
credentials.json
```

### Password hashing (REQUIRED)
```go
// ❌ Never use MD5/SHA for passwords
hash := md5.Sum([]byte(password))

// ✅ Use bcrypt or argon2
import "golang.org/x/crypto/bcrypt"
hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
```

### Input validation
```go
// ❌ SQL injection
query := fmt.Sprintf("SELECT * FROM users WHERE id='%s'", id)

// ✅ Parameterized query
query := "SELECT * FROM users WHERE id=$1"
row := db.QueryRow(query, id)
```

### JWT best practices
```go
// ✅ Short-lived access tokens (15min)
// ✅ Refresh tokens with rotation
// ✅ Validate issuer, audience, expiry
// ✅ Use RS256, not HS256 for multi-service
```

## Anti-Patterns
### Don't: Log sensitive data
```go
log.Printf("User password: %s", password)  // ❌
log.Printf("User %s logged in", username)  // ✅
```

### Don't: Use http for sensitive data
```go
http.ListenAndServe(":8080", handler)  // ❌ No TLS
http.ListenAndServeTLS(":443", "cert.pem", "key.pem", handler)  // ✅
```

## Quick Reference
| Check | Tool |
|-------|------|
| Secrets | `gitleaks detect` |
| Dependencies | `npm audit` / `go vuln check` |
| SAST | `semgrep --config=auto` |
| Container | `trivy image myapp:latest` |
