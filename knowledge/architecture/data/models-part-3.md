---
category: architecture
subcategory: data
confidence: low
documentType: explanation
scope: org
contentHash: dbbcc4be8d30
tags: [metrics, prometheus]
source: architecture/data/models.md
verified: 2026-07-22
splitPartIndex: 3
splitPartTotal: 5
canonical: false
synthetic: split-part
---

## Data Models & Schemas (Part 3)

_struct · `client/clientcache.go`:22_

| Field | Type | Optional |
|-------|------|----------|
| `cache` | `*cache.Cache` | no |

## memCache

_struct · `middleware/cache.go`:26_

| Field | Type | Optional |
|-------|------|----------|
| `cache` | `*cache.Cache` | no |

## MetricsPusherClient

_struct · `metrics/pusher.go`:13_

| Field | Type | Optional |
|-------|------|----------|
| `client` | `client.Request` | no |

## MigrationRunner

_struct · `database/migrate.go`:37_

| Field | Type | Optional |
|-------|------|----------|
| `mgr` | `*migrate.Migrate` | no |
| `filesDir` | `string` | no |
| `logger` | `logger` | no |

## migrationsLogger

_struct · `database/migrate.go`:162_

| Field | Type | Optional |
|-------|------|----------|
| `logger` | `logger` | no |

## MockDBContextGetter

_struct · `database/mock_db.go`:16_

| Field | Type | Optional |
|-------|------|----------|
| `ctrl` | `*gomock.Controller` | no |
| `recorder` | `*MockDBContextGetterMockRecorder` | no |

## MockDBContextGetterMockRecorder

_struct · `database/mock_db.go`:22_

| Field | Type | Optional |
|-------|------|----------|
| `mock` | `*MockDBContextGetter` | no |

## MockTrxContextGetter

_struct · `database/mock_db.go`:53_

| Field | Type | Optional |
|-------|------|----------|
| `ctrl` | `*gomock.Controller` | no |
| `recorder` | `*MockTrxContextGetterMockRecorder` | no |

## MockTrxContextGetterMockRecorder

_struct · `database/mock_db.go`:59_

| Field | Type | Optional |
|-------|------|----------|
| `mock` | `*MockTrxContextGetter` | no |

## NodeInfoResponse

_struct · `blockchain/binance/model.go`:6_

| Field | Type | Optional |
|-------|------|----------|
| `SyncInfo` | `struct { LatestBlockHeight int `json:"latest_block_height"` }` | no |

## OperationData

_struct · `database/migrate.go`:31_

| Field | Type | Optional |
|-------|------|----------|
| `ID` | `string` | no |
| `ForceVersion` | `int` | no |

## OrderedSet

_struct · `set/ordered.go`:3_

| Field | Type | Optional |
|-------|------|----------|
| `valuesSet` | `map[T]struct{}` | no |
| `values` | `[]T` | no |

## Path

_struct · `client/path.go`:5_

| Field | Type | Optional |
|-------|------|----------|
| `template` | `string` | no |
| `values` | `[]any` | no |

## performanceMetric

_struct · `metrics/metrics.go`:25_

| Field | Type | Optional |
|-------|------|----------|
| `executionStarted` | `*prometheus.GaugeVec` | no |
| `executionDurationSeconds` | `*prometheus.HistogramVec` | no |
| `executionSucceededTotal` | `*prometheus.CounterVec` | no |
| `executionFailedTotal` | `*prometheus.CounterVec` | no |

## PublishConfig

_struct · `mq/queue.go`:69_

| Field | Type | Optional |
|-------|------|----------|
| `MaxRetries` | `*int` | no |
| `DeliveryMode` | `DeliveryMode` | no |

## pusher

_struct · `metrics/pusher.go`:37_

| Field | Type | Optional |
|-------|------|----------|
| `pusher` | `*push.Pusher` | no |

## queue

_struct · `mq/queue.go`:5_

| Field | Type | Optional |
|-------|------|----------|
| `name` | `QueueName` | no |
| `client` | `*Client` | no |
