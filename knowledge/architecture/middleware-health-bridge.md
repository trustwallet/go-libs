---
title: Middleware ↔ Health Domain Bridge
category: architecture
tags: [middleware, health, cross-domain, bridge, gin, cache-control]
confidence: high
source: middleware/cache.go, middleware/cache_control.go, health/http.go
updated: 2026-07-22
---

# Middleware ↔ Health Domain Bridge

## Cross-Domain Connection

The relation graph identifies two cross-domain bridges between the `middleware` package and the `health` package:

1. `CacheMiddleware` → `handle` (in `health/http.go`)
2. `CacheControl` → `handle` (in `health/http.go`)

Both bridges have `pairLinks: 2`, meaning these are the only two cross-domain links between middleware and health.

## Why These Links Exist

The health check HTTP handler (`health/http.go#handle`) is an ordinary Gin handler function. When consumers register the health endpoint with the Gin router **and** also apply the cache or cache-control middleware to the router group, the middleware wraps the health endpoint. The bridge is therefore **incidental to composition, not a deliberate dependency** — the `middleware` package does not import `health`, and `health` does not import `middleware`. The connection arises at the call-site (consumer service) where both are applied to the same router.

This means the bridge is **not a real architectural coupling** — it is an artifact of the AST's heuristic call-edge inference. No code change is needed; this is an expected false positive in the relation graph.

## Health Package

`health.Handle(checks ...HealthCheck)` returns a Gin `HandlerFunc` that runs a list of health-check functions and returns HTTP 200 (all pass) or 503 (any fail). It is used by consumer services to expose `/livez`, `/readyz`, or `/health` endpoints.

## See Also
- [middleware.md](../features/middleware.md)
- [health.md](../features/health.md)
- [architecture/overview.md](overview.md)
- [crypto and auth](../security/crypto-and-auth.md) <!-- rel:strong -->
- [mq pattern](../patterns/mq-pattern.md) <!-- rel:strong -->
