---
name: docker
description: >
  Docker containerization patterns. Trigger: When writing Dockerfiles, docker-compose, or containerizing applications.
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use
- Creating or modifying Dockerfiles
- Working with docker-compose.yml
- Containerizing applications
- Optimizing Docker builds
- Deploying containerized apps

## Critical Patterns

### Multi-stage builds (REQUIRED)
```dockerfile
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

FROM node:20-alpine AS production
WORKDIR /app
COPY --from=build /app/dist ./dist
COPY --from=build /app/node_modules ./node_modules
USER node
EXPOSE 3000
CMD ["node", "dist/main.js"]
```

### Non-root user (REQUIRED)
```dockerfile
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup
USER appuser
```

### .dockerignore (REQUIRED)
```
node_modules
.git
.env
*.log
dist
coverage
```

## Anti-Patterns
### Don't: Use latest tag in production
```dockerfile
FROM node:latest  # ❌ Non-reproducible
FROM node:20-alpine  # ✅ Pinned version
```

### Don't: Store secrets in image
```dockerfile
ENV API_KEY=secret123  # ❌ Visible in image layers
```

## Quick Reference
| Task | Pattern |
|------|---------|
| Build | `docker build -t app:latest .` |
| Run | `docker run -p 3000:3000 app:latest` |
| Compose | `docker compose up -d` |
| Logs | `docker logs -f container_name` |
| Clean | `docker system prune -af` |
