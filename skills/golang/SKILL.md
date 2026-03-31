---
name: golang
description: >
  Go language best practices and patterns. Trigger: When writing Go code, working with modules, or Go projects.
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use
- Writing Go source files (.go)
- Working with Go modules (go.mod)
- Implementing interfaces
- Error handling and testing
- Concurrency patterns (goroutines, channels)

## Critical Patterns

### Error handling (REQUIRED)
```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to read config: %w", err)
}

// Sentinel errors
var ErrNotFound = errors.New("not found")
```

### Table-driven tests (REQUIRED)
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

### Context for cancellation
```go
func fetchWithTimeout(ctx context.Context, url string) ([]byte, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := http.DefaultClient.Do(req)
    // ...
}
```

## Anti-Patterns
### Don't: Ignore errors
```go
data, _ := os.ReadFile("file.txt")  // ❌ Silent failure
data, err := os.ReadFile("file.txt")  // ✅ Handle it
```

### Don't: Use panic for normal errors
```go
panic("file not found")  // ❌ Use error return
return fmt.Errorf("file not found: %w", err)  // ✅
```

## Quick Reference
| Task | Command |
|------|---------|
| Build | `go build ./...` |
| Test | `go test ./...` |
| Lint | `golangci-lint run` |
| Format | `go fmt ./...` |
| Vet | `go vet ./...` |
| Mod tidy | `go mod tidy` |
