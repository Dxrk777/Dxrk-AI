---
name: graphql
description: >
  GraphQL API design patterns. Trigger: When writing GraphQL schemas, resolvers, or working with GraphQL APIs.
metadata:
  author: dxrk
  version: "1.0"
---

## When to Use
- Writing GraphQL schemas (.graphql)
- Implementing resolvers (Go, Node.js, Python)
- Querying/mutating GraphQL APIs
- Setting up subscriptions
- Optimizing N+1 queries with dataloaders

## Critical Patterns

### Schema-first design (REQUIRED)
```graphql
type User {
  id: ID!
  name: String!
  email: String!
  posts: [Post!]!
}

type Query {
  user(id: ID!): User
  users(limit: Int = 10, offset: Int = 0): [User!]!
}

type Mutation {
  createUser(input: CreateUserInput!): User!
}

input CreateUserInput {
  name: String!
  email: String!
}
```

### Dataloader pattern (REQUIRED for Go)
```go
type Loaders struct {
    UserByID *dataloader.Loader[string, *model.User]
}

func NewLoaders(repo Repository) *Loaders {
    return &Loaders{
        UserByID: dataloader.NewBatchedLoader(
            func(ctx context.Context, ids []string) []*dataloader.Result[*model.User] {
                users, err := repo.GetUsersByIDs(ctx, ids)
                // batch and return
            },
        ),
    }
}
```

## Anti-Patterns
### Don't: Return unbounded lists
```graphql
type Query {
  users: [User!]!  # ❌ Could return millions
  users(limit: Int!, offset: Int!): [User!]!  # ✅ Paginated
}
```

## Quick Reference
| Task | Tool |
|------|------|
| Generate Go | `go run github.com/99designs/gqlgen generate` |
| Playground | Apollo Studio, GraphiQL |
| Client | urql, Apollo Client, graphql-request |
