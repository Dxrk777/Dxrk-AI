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

# API Design Skill

## Trigger
"design API", "REST", "GraphQL", "endpoint", "API design", "create endpoint"

## Principles
1. **Consistency**: Uniform patterns across all endpoints
2. **Simplicity**: Complex operations through resource composition
3. **Clarity**: Self-documenting through naming and structure
4. **Safety**: Authentication, authorization, rate limiting

## REST Best Practices
```
GET    /users          # List users
GET    /users/:id      # Get single user
POST   /users          # Create user
PUT    /users/:id      # Full update
PATCH  /users/:id      # Partial update
DELETE /users/:id      # Delete user

# Nested resources (max 2 levels)
GET    /users/:id/posts
GET    /users/:id/posts/:postId/comments
```

## Response Format
```json
{
  "data": { },
  "meta": { "page": 1, "total": 100 },
  "errors": [ ]
}
```

## Status Codes
- 200: Success
- 201: Created
- 204: No Content (DELETE)
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 422: Validation Error
- 429: Rate Limited
- 500: Server Error
