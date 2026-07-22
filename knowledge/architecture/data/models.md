---
category: architecture
subcategory: data
confidence: low
documentType: explanation
scope: repo
contentHash: b2c19411247b
tags: [api, handler, middleware, request]
source: architecture/data/models.md
verified: 2026-07-22
splitPartIndex: 1
splitPartTotal: 5
canonical: true
---

## Data Models & Schemas
<!-- sdd-knowledge-generated -->

> Field-level shape of data models extracted via tree-sitter: TS interfaces / type aliases, Zod `z.object` schemas, Go/Rust/Swift structs, Kotlin data classes, and Python dataclasses. Scoped to domain data — UI views/props, view-models, design tokens (theme/style/colors), and constant/identifier namespaces are excluded. Deterministic, no LLM.

## AccountMeta

_struct · `blockchain/binance/model.go`:61_

| Field | Type | Optional |
|-------|------|----------|
| `Balances` | `[]TokenBalance` | no |

## api

_struct · `httplib/server.go`:19_

| Field | Type | Optional |
|-------|------|----------|
| `router` | `http.Handler` | no |
| `port` | `string` | no |
| `h2c` | `bool` | no |

## AssetInfoResp

_struct · `client/api/backend/model.go`:4_

| Field | Type | Optional |
|-------|------|----------|
| `Name` | `string` | no |
| `Symbol` | `string` | no |
| `Type` | `string` | no |
| `Decimals` | `int` | no |
| `AssetID` | `string` | no |

## Bep2Asset

_struct · `blockchain/binance/explorer/model.go`:4_

| Field | Type | Optional |
|-------|------|----------|
| `Asset` | `string` | no |
| `Name` | `string` | no |
| `AssetImg` | `string` | no |
| `MappedAsset` | `string` | no |
| `Decimals` | `int` | no |

## Bep2Assets

_struct · `blockchain/binance/explorer/model.go`:12_

| Field | Type | Optional |
|-------|------|----------|
| `AssetInfoList` | `[]Bep2Asset` | no |

## builder

_struct · `worker/worker.go`:19_

| Field | Type | Optional |
|-------|------|----------|
| `worker` | `*worker` | no |

## cachedWriter

_struct · `middleware/cache.go`:37_

| Field | Type | Optional |
|-------|------|----------|
| `status` | `int` | no |
| `written` | `bool` | no |
| `expire` | `time.Duration` | no |
| `key` | `string` | no |

## cacheResponse

_struct · `middleware/cache.go`:31_

| Field | Type | Optional |
|-------|------|----------|
| `Status` | `int` | no |
| `Header` | `http.Header` | no |
| `Data` | `[]byte` | no |

## Client

_struct · `blockchain/binance/api/client.go`:14_

| Field | Type | Optional |
|-------|------|----------|
| `req` | `client.Request` | no |

## Client

_struct · `blockchain/binance/client.go`:14_

| Field | Type | Optional |
|-------|------|----------|
| `req` | `client.Request` | no |

## Client

_struct · `blockchain/binance/explorer/client.go`:13_

| Field | Type | Optional |
|-------|------|----------|
| `req` | `client.Request` | no |

## Client

_struct · `client/api/backend/client.go`:10_

| Field | Type | Optional |
|-------|------|----------|
| `req` | `client.Request` | no |

## Client

_struct · `mq/mq.go`:26_

| Field | Type | Optional |
|-------|------|----------|
| `url` | `string` | no |
| `conn` | `*amqp.Connection` | no |
| `amqpChan` | `*amqp.Channel` | no |
| `connClients` | `[]ConnectionClient` | no |
| `connCheckTimeout` | `time.Duration` | no |

## consumer

_struct · `mq/consumer.go`:17_

| Field | Type | Optional |
|-------|------|----------|
| `client` | `*Client` | no |
| `queue` | `Queue` | no |
| `messageProcessor` | `MessageProcessor` | no |
| `options` | `*ConsumerOptions` | no |
| `messages` | `<-chan amqp.Delivery` | no |
| `stopChan` | `chan struct{}` | no |

## ConsumerOptions

_struct · `mq/options.go`:9_

| Field | Type | Optional |
|-------|------|----------|
| `Workers` | `int` | no |
| `Prefetch` | `int` | no |
| `RetryOnError` | `bool` | no |
| `RetryDelay` | `time.Duration` | no |
| `PerformanceMetric` | `metrics.PerformanceMetric` | no |
| `MaxRetries` | `int` | no |

## See Also
- [blockchain](../../features/blockchain.md) <!-- rel:strong -->
- [mq](../../features/mq.md) <!-- rel:strong -->
- [middleware](../../features/middleware.md) <!-- rel:strong -->
- [worker pattern](../../patterns/worker-pattern.md) <!-- rel:strong -->
- [common mistakes and anti patterns](../../guides/troubleshooting/common-mistakes-and-anti-patterns.md) <!-- rel:strong -->
