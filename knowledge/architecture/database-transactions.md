---
title: Database Layer — Context-Scoped Transactions
category: architecture
tags: [database, postgres, gorm, transactions, context, read-write-split]
confidence: high
source: database/db.go, database/migrate.go, database/migration_runner_env.go
updated: 2026-07-22
---

# Database Layer — Context-Scoped Transactions

## Overview

The `database` package wraps GORM with two design invariants that all consumers must understand:

1. **Transactions are threaded through `context.Context`** — business logic does not receive a `*gorm.DB` directly. Instead, it calls `getter.DBFrom(ctx)` which returns either the transaction or the main DB depending on whether the context carries an active transaction.
2. **Read/write splitting is transparent** — `dbresolver` routes reads to replicas and writes to the primary. Callers that must force a specific connection use `.Clauses(dbresolver.Write)` or `.Clauses(dbresolver.Read)`.

## Key Types

### `DBContextGetter` interface

```go
type DBContextGetter interface {
    DBFrom(ctx context.Context) *gorm.DB
}
```

Business logic should depend on this interface, not on `*DBGetter` directly. `DBFrom(ctx)` returns the active transaction `*gorm.DB` if one was injected by `Transaction(...)`, or the main DB otherwise.

### `TrxContextGetter` interface

```go
type TrxContextGetter interface {
    Transaction(ctx context.Context, fc func(ctx context.Context) error) error
}
```

Starts a GORM transaction and injects it into a new `ctx` via a private `trxKey`. The business-logic callback receives this enriched context — all calls to `DBFrom(enrichedCtx)` inside the callback automatically use the transaction.

### `DBGetter` struct

Implements both interfaces. Created via `NewDBGetter(cfg DBConfig)` or `NewDBGetterFromGormInstance(db)` (for testing with an existing GORM instance).

## Read/Write Splitting Caveat

When using replicas (`cfg.ReadonlyUrl` is set), a read immediately after a write may return stale data if the replica has not yet caught up with replication lag. This is a documented caveat in the source (`db.go:41-45`). Callers with strict consistency requirements must explicitly use `Clauses(dbresolver.Write)`.

## Migration Runner

`MigrationRunner` (in `database/migrate.go`) manages `golang-migrate` migrations:

- Default migration directory: `"dbmigrations"` (configurable via `WithFilesDir`)
- Supported operations: `"up"` (apply all pending), `"down"` (roll back one step), `"force"` (set version without running SQL)
- `RunMigrationsFromEnv` (in `migration_runner_env.go`) reads config from environment variables — this is the **god-node** entry point for migrations (fan-out: 4 downstream calls)

## See Also
- [database.md](../features/database.md)
- [architecture/overview.md](overview.md)
- [integration testing](../tests/integration-testing.md) <!-- rel:strong -->
- [worker pattern](../patterns/worker-pattern.md) <!-- rel:related -->
- [mq pattern](../patterns/mq-pattern.md) <!-- rel:related -->
