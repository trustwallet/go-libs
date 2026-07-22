---
category: architecture
confidence: low
documentType: explanation
scope: repo
contentHash: e7e22f90f276
tags: [api, middleware, request]
source: architecture/call-graph.md
verified: 2026-07-22
splitPartIndex: 2
splitPartTotal: 2
canonical: false
synthetic: split-part
---

## Call Graph (Part 2)

- Call edges: **70** across **467** symbols
- Reachable from entry points (exported symbols + routes): **403**
- Dependency cycles: **0**
- Cross-domain bridges: **2**
- Dead-code candidates (non-exported, zero callers, unreachable): **62**

### Most-coupled symbols (god-node ranking)

| Symbol | Fan-in | Fan-out | Degree |
|--------|--------|---------|--------|
| `NewReqBuilder` | 7 | 0 | 7 |
| `RunMigrationsFromEnv` | 0 | 4 | 4 |
| `publish` | 3 | 1 | 4 |
| `getRuntimeTags` | 3 | 1 | 4 |
| `genID` | 3 | 0 | 3 |
| `handle` | 3 | 0 | 3 |
| `CacheMiddleware` | 0 | 3 | 3 |
| `InitClusterWithTLS` | 0 | 2 | 2 |
| `InitWithTLS` | 0 | 2 | 2 |
| `reportMonitoringMetricsIfEnabled` | 0 | 2 | 2 |
| `InitJSONClient` | 0 | 2 | 2 |
| `RpcCall` | 0 | 2 | 2 |
| `RpcCallRaw` | 0 | 2 | 2 |
| `GetRSAPrivateKey` | 2 | 0 | 2 |
| `NewMigrationRunner` | 1 | 1 | 2 |

### Cross-domain bridges

> Calls that cross a domain boundary — the integration points between subsystems. A sole link between two domains is the connection you'd least expect.

| From (domain) | To (domain) | Domain links |
|------|------|-----|
| `CacheControl` (middleware) | `handle` (health) | 2 |
| `CacheMiddleware` (middleware) | `handle` (health) | 2 |

### Possible duplicate entities (name variants)

> Symbol names that normalize identically — likely the same entity spelled inconsistently. Unify or distinguish in the docs.

- `Builder` / `builder`
- `Consumer` / `consumer`
- `Contains` / `contains`
- `ContainsAll` / `containsAll`
- `ContainsAny` / `containsAny`
- `Downloader` / `downloader`
- `Exchange` / `exchange`
- `GetBody` / `getBody`
- `HttpServerMetric` / `httpServerMetric`
- `Init` / `init`
- `Logger` / `logger`
- `PerformanceMetric` / `performanceMetric`
- `Process` / `process`
- `Publish` / `publish`
- `PublishWithConfig` / `publishWithConfig`
- `Pusher` / `pusher`
- `Queue` / `queue`
- `Reconnect` / `reconnect`
- `Server` / `server`
- `Start` / `start`
- `Worker` / `worker`

### Dead-code candidates

> Non-exported symbols with no resolved callers, unreachable from any entry point. Static analysis cannot see dynamic dispatch — verify before removing.

- `redisClient` (`cache/redis/client_interface.go`)
- `redisClientForTest` (`cache/redis/client_interface.go`)
- `reconnectCluster` (`cache/redis/redis.go`)
- `constructHttpRequest` (`client/client_execute.go`)
- `metricsEnabled` (`client/client_execute.go`)
- `reportMonitoringMetricsIfEnabled` (`client/client_execute.go`)
- `setRequestHeaders` (`client/client_execute.go`)
- `httpClientMetrics` (`client/client_metrics.go`)
- `observeDuration` (`client/client_metrics.go`)
- `observeResult` (`client/client_metrics.go`)
- `generateKey` (`client/clientcache.go`)
- `getCache` (`client/clientcache.go`)
- `init` (`client/clientcache.go`)
- `memCache` (`client/clientcache.go`)
- `setCache` (`client/clientcache.go`)
- `fillDefaultValues` (`client/jsonrpc.go`)
- `pathMetricEnabled` (`client/request.go`)
- `applyDefaultValue` (`database/config.go`)
- `transactionKey` (`database/db.go`)
- `logger` (`database/migrate.go`)
- `migrationsLogger` (`database/migrate.go`)
- `noopLogger` (`database/migrate.go`)
- `runDown` (`database/migrate.go`)
- `runForce` (`database/migrate.go`)
- `runUp` (`database/migrate.go`)
- `verifySignature` (`gin/hmac.go`)
- `server` (`health/http.go`)
- `downloader` (`httplib/downloader.go`)
- `api` (`httplib/server.go`)
- `serve` (`httplib/server.go`)
- `init` (`logging/formatter_strict_text.go`)
- `init` (`logging/logger.go`)
- `httpServerMetric` (`metrics/http_metrics.go`)
- `performanceMetric` (`metrics/metrics.go`)
- `pusher` (`metrics/pusher.go`)
- `cachedWriter` (`middleware/cache.go`)
- `cacheResponse` (`middleware/cache.go`)
- `deleteCache` (`middleware/cache.go`)
- `getCache` (`middleware/cache.go`)
- `init` (`middleware/cache.go`)
- `memCache` (`middleware/cache.go`)
- `setCache` (`middleware/cache.go`)
- `consume` (`mq/consumer.go`)
- `consumer` (`mq/consumer.go`)
- `getRemainingRetries` (`mq/consumer.go`)
- `getSanitizedPrefetchCount` (`mq/consumer.go`)
- `messageChannel` (`mq/consumer.go`)
- `process` (`mq/consumer.go`)
- `exchange` (`mq/exchange.go`)
- `initNotifyCloseListeners` (`mq/mq.go`)

## See Also
- [client](../features/client.md) <!-- rel:strong -->
- [mq](../features/mq.md) <!-- rel:related -->
- [database](../features/database.md) <!-- rel:related -->
- [middleware](../features/middleware.md) <!-- rel:related -->
- [metrics](../features/metrics.md) <!-- rel:related -->
