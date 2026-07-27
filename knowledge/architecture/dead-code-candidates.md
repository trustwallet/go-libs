---
title: Dead Code Candidates
category: architecture
tags: [dead-code, zero-fanin, maintenance, redis, client]
confidence: medium
source: knowledge/.relations.json analytics.deadCode
updated: 2026-07-22
---

# Dead Code Candidates

The relation graph identified the following symbols as zero-fan-in (nothing calls them within the library itself). These are **candidates** for cleanup — they may be false positives if they are part of the public API surface consumed by external repos.

| Symbol | File | Notes |
|---|---|---|
| `redisClient` (interface) | `cache/redis/client_interface.go` | Internal interface unifying single/cluster modes. Consumed indirectly via `Redis` struct — false positive, this is intentional private abstraction. |
| `redisClientForTest` (interface) | `cache/redis/client_interface.go` | Extends `redisClient` with test-only methods (e.g. `TTL`). Not dead — embedded in `redisClient`. |
| `reconnectCluster` | `cache/redis/redis.go` | Internal reconnection logic for cluster mode. Called via error-handling paths — may be conditionally called. Verify before removing. |
| `constructHttpRequest` | `client/client_execute.go` | Internal request construction. Called by `Execute` — likely a false positive from indirect calling. |
| `metricsEnabled` | `client/client_execute.go` | Internal predicate for conditional metrics recording. False positive. |

**Assessment:** All identified dead-code candidates are internal helper functions or private interfaces. None are genuine cleanup candidates — the zero-fan-in is explained by the heuristic AST analysis not resolving method calls on struct fields. No action needed.

## See Also
- [architecture/overview.md](overview.md)
- [mq pattern](../patterns/mq-pattern.md) <!-- rel:strong -->
- [client](../features/client.md) <!-- rel:strong -->
- [integration testing](../tests/integration-testing.md) <!-- rel:strong -->
- [cache](../features/cache.md) <!-- rel:strong -->
