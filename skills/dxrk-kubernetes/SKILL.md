---
name: kubernetes
description: >
  Kubernetes orchestration patterns. Trigger: When writing K8s manifests, Helm charts, or deploying to clusters.
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use
- Writing Kubernetes YAML manifests
- Creating Helm charts
- Configuring deployments, services, ingress
- Managing secrets and configmaps
- Setting up horizontal pod autoscaling

## Critical Patterns

### Deployment with health checks (REQUIRED)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    spec:
      containers:
      - name: app
        image: myapp:latest
        ports:
        - containerPort: 3000
        livenessProbe:
          httpGet:
            path: /healthz
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /ready
            port: 3000
          initialDelaySeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
```

### Secret management (REQUIRED)
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
stringData:
  DB_PASSWORD: "use-external-secrets-operator-in-prod"
```

## Anti-Patterns
### Don't: Use :latest in production
### Don't: Run as root
### Don't: Skip resource limits

## Quick Reference
| Task | Command |
|------|---------|
| Apply | `kubectl apply -f manifest.yaml` |
| Logs | `kubectl logs -f deployment/app` |
| Scale | `kubectl scale deployment app --replicas=5` |
| Debug | `kubectl exec -it pod-name -- /bin/sh` |
| Port-forward | `kubectl port-forward svc/app 3000:3000` |
