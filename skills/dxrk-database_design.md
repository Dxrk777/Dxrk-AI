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

# Database Design Skill

## Trigger
"database design", "schema", "migration", "SQL", "PostgreSQL", "MongoDB"

## Normalization Levels
1. **1NF**: Atomic values, no repeating groups
2. **2NF**: No partial dependencies (1NF + no composite keys)
3. **3NF**: No transitive dependencies
4. **BCNF**: Every determinant is a candidate key

## PostgreSQL Patterns
```sql
-- Index strategy
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
CREATE INDEX idx_orders_user_id ON orders(user_id, created_at DESC);

-- Soft delete
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP;
CREATE INDEX idx_users_deleted ON users(deleted_at) WHERE deleted_at IS NOT NULL;

-- Audit trail
CREATE TABLE audit_log (
  id SERIAL PRIMARY KEY,
  table_name TEXT NOT NULL,
  record_id INTEGER NOT NULL,
  action TEXT NOT NULL,
  old_data JSONB,
  new_data JSONB,
  changed_by UUID REFERENCES users(id),
  changed_at TIMESTAMP DEFAULT NOW()
);
```

## MongoDB Patterns
```javascript
// Document design
{
  _id: ObjectId,
  userId: ObjectId,        // Reference
  items: [{                 // Embedded
    productId: ObjectId,
    quantity: Number,
    price: Number
  }],
  total: Number,
  createdAt: Date
}

// Indexes
db.orders.createIndex({ userId: 1, createdAt: -1 })
db.orders.createIndex({ "items.productId": 1 })
```
