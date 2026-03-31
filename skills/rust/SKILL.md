---
name: rust
description: >
  Rust language patterns and best practices. Trigger: When writing Rust code, working with Cargo, or Rust projects.
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use
- Writing Rust source files (.rs)
- Configuring Cargo.toml
- Working with ownership, borrowing, lifetimes
- Error handling patterns
- Async Rust with tokio

## Critical Patterns

### Error handling with Result (REQUIRED)
```rust
use thiserror::Error;

#[derive(Error, Debug)]
enum AppError {
    #[error("database error: {0}")]
    Database(#[from] sqlx::Error),
    #[error("not found: {0}")]
    NotFound(String),
}

fn get_user(id: &str) -> Result<User, AppError> {
    // use ? operator
}
```

### Ownership patterns (REQUIRED)
```rust
// Prefer borrowing over ownership transfer
fn process(data: &str) -> String {  // ✅ borrows
    data.to_uppercase()
}

// Use Arc<Mutex<T>> for shared mutable state
use std::sync::{Arc, Mutex};
let shared = Arc::new(Mutex::new(Vec::new()));
```

### Builder pattern
```rust
#[derive(Default)]
struct Config {
    host: String,
    port: u16,
}

impl Config {
    fn builder() -> Self { Self::default() }
    fn host(mut self, host: &str) -> Self {
        self.host = host.to_string();
        self
    }
    fn port(mut self, port: u16) -> Self {
        self.port = port;
        self
    }
}
```

## Anti-Patterns
### Don't: Use unwrap() in production
```rust
let val = map.get("key").unwrap();  // ❌ Panics if missing
let val = map.get("key").ok_or("missing key")?;  // ✅ Returns error
```

## Quick Reference
| Task | Command |
|------|---------|
| Build | `cargo build --release` |
| Test | `cargo test` |
| Lint | `cargo clippy` |
| Format | `cargo fmt` |
| Doc | `cargo doc --open` |
