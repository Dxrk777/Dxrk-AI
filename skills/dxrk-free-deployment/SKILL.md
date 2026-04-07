---
name: dxrk-free-deployment
description: >
  Deploy applications to free-tier cloud platforms: Fly.io, Railway, Render, Railway.
  Trigger: Deploying apps, free hosting, Docker deployment, server setup.
license: Apache-2.0
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use

- Deploying applications for free
- Setting up Docker containers for cloud
- Configuring Fly.io, Railway, or Render deployments
- Setting up SSL, domains, and environment variables
- Creating production-ready Dockerfiles

## Free Tier Comparison

| Platform | Free Tier | Persistent | SSL | Custom Domain | Best For |
|----------|-----------|------------|-----|---------------|----------|
| **Fly.io** | 3 shared VMs | ✅ Yes | ✅ Free | ✅ Free | Go, Docker-native |
| **Railway** | $5 credit/month | ✅ Yes | ✅ Auto | ✅ Free | Quick deploys |
| **Render** | 750h/month | ⚠️ Sleeps | ✅ Free | ✅ Free | Web services |
| **Neon** | 512 MB SQL | ✅ Yes | - | - | Serverless Postgres |
| **Turso** | 500 DBs, 9 GB | ✅ Yes | - | - | Distributed SQLite |

## Fly.io (RECOMMENDED for Go)

Best for Go applications - native support, no Dockerfile needed.

### Quick Deploy

```bash
# Install Fly CLI
curl -L https://fly.io/install.sh | sh

# Authenticate
fly auth login

# Launch app (auto-detects Go)
fly launch --name my-app

# Deploy
fly deploy

# Check status
fly status

# View logs
fly logs

# SSH into app
fly ssh console
```

### Go Dockerfile (if needed)

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Minimal runtime
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

### fly.toml Configuration

```toml
app = "my-app"
primary_region = "iad"

[build]
builder = "herokuish"

[deploy]
  release_command = "go run migrations.go"

[env]
  PORT = "8080"
  ENVIRONMENT = "production"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = false
  auto_start_machines = true

[[vm]]
  memory = "256mb"
  cpu_kind = "shared"
  cpus = 1
```

## Railway (Easiest Setup)

Connect GitHub repo, Railway auto-detects stack.

### railway.toml

```toml
[build]
builder = "nixpacks"

[deploy]
  numReplicas = 1
  restartPolicyType = "ON_FAILURE"
  restartPolicyMaxRetries = 10

[environment]
  PORT = "8080"
  NODE_ENV = "production"
```

### Environment Variables (Railway Dashboard)

```
# Database
DATABASE_URL=postgres://user:pass@host:5432/db

# Auth
JWT_SECRET=your-secret-key

# Feature Flags
FEATURE_NEW_UI=true
```

## Render (24/7 Free with Sleep)

### render.yaml

```yaml
services:
  - type: web
    name: api
    runtime: docker
    plan: free  # 750h/month
    region: oregon
    
    dockerfilePath: ./Dockerfile
    
    healthCheckPath: /health
    
    envVars:
      - key: PORT
        value: 8080
      - key: DATABASE_URL
        sync: false  # Set in dashboard
      - key: NODE_ENV
        value: production

cronJobs:
  - name: cleanup
    command: npm run cleanup
    schedule: "0 * * * *"  # Every hour
    region: oregon
```

## Docker Compose for Local Dev

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://postgres:postgres@db:5432/myapp
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=myapp
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

## SSL & Domain Setup

### Fly.io Custom Domain

```bash
# Add custom domain
fly certs create mydomain.com

# Check certificate status
fly certs show mydomain.com

# Add www subdomain
fly certs create www.mydomain.com
```

### Let's Encrypt (Manual)

```bash
# Install certbot
sudo apt install certbot

# Generate certificate
sudo certbot certonly --standalone -d mydomain.com -d www.mydomain.com

# Certificate location
# /etc/letsencrypt/live/mydomain.com/fullchain.pem
# /etc/letsencrypt/live/mydomain.com/privkey.pem
```

## Health Check Endpoints

Always implement health checks for cloud platforms:

```go
// Go health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    health := map[string]interface{}{
        "status": "healthy",
        "version": os.Getenv("VERSION"),
        "timestamp": time.Now().UTC(),
    }
    
    // Check dependencies
    if err := db.Ping(); err != nil {
        health["status"] = "unhealthy"
        health["db"] = "disconnected"
        w.WriteHeader(503)
    }
    
    json.NewEncoder(w).Encode(health)
}
```

```typescript
// Node.js health check
app.get('/health', async (req, res) => {
  const health = {
    status: 'ok',
    uptime: process.uptime(),
    database: await checkDatabase()
  };
  
  res.status(health.status === 'ok' ? 200 : 503).json(health);
});
```

## Commands

```bash
# Fly.io
fly launch
fly deploy
fly status
fly logs
fly secrets set KEY=value
fly secrets unset KEY

# Railway
railway login
railway init
railway up
railway open

# Render
render deploy
render logs --service=api

# Docker
docker build -t myapp .
docker run -p 8080:8080 myapp
docker-compose up -d
```

## Cost Optimization Tips

1. **Use Fly.io for Go** - Native build support, no Dockerfile needed
2. **Use Railway for quick prototypes** - $5 free credit goes far
3. **Use Render for web services** - 750h free, auto-sleeps when unused
4. **Use Neon for Postgres** - Serverless, scales to zero
5. **Use Turso for SQLite** - Free distributed database

## Resources

- **Fly.io Docs**: https://fly.io/docs/
- **Railway Docs**: https://docs.railway.app/
- **Render Docs**: https://render.com/docs
