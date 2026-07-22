---
category: architecture
subcategory: data
confidence: low
documentType: explanation
scope: repo
contentHash: d4bd08b21830
tags: [request, response]
source: architecture/data/models.md
verified: 2026-07-22
splitPartIndex: 4
splitPartTotal: 5
canonical: false
synthetic: split-part
---

## Data Models & Schemas (Part 4)

_struct · `cache/redis/redis.go`:18_

| Field | Type | Optional |
|-------|------|----------|
| `client` | `redisClient` | no |
| `isCluster` | `bool` | no |

## Req

_struct · `client/request.go`:12_

| Field | Type | Optional |
|-------|------|----------|
| `headers` | `map[string]string` | no |
| `resultContainer` | `any` | no |
| `method` | `string` | no |
| `path` | `Path` | no |
| `query` | `url.Values` | no |
| `body` | `any` | no |
| `rawResponseContainer` | `*http.Response` | no |
| `metricName` | `string` | no |
| `pathMetricEnabled` | `bool` | no |

## ReqBuilder

_struct · `client/request.go`:25_

| Field | Type | Optional |
|-------|------|----------|
| `req` | `*Req` | no |

## Request

_struct · `client/client.go`:17_

| Field | Type | Optional |
|-------|------|----------|
| `BaseURL` | `string` | no |
| `Headers` | `map[string]string` | no |
| `Host` | `string` | no |
| `HttpClient` | `HTTPClient` | no |
| `HttpErrorHandler` | `HttpErrorHandler` | no |
| `metricRegisterer` | `prometheus.Registerer` | no |
| `httpMetrics` | `*httpClientMetrics` | no |

## RpcError

_struct · `client/jsonrpc.go`:38_

| Field | Type | Optional |
|-------|------|----------|
| `Code` | `int` | no |
| `Message` | `string` | no |
| `Data` | `string` | no |

## RpcRequest

_struct · `client/jsonrpc.go`:17_

| Field | Type | Optional |
|-------|------|----------|
| `JsonRpc` | `string` | no |
| `Method` | `string` | no |
| `Params` | `interface{}` | no |
| `Id` | `int64` | no |

## RpcResponse

_struct · `client/jsonrpc.go`:24_

| Field | Type | Optional |
|-------|------|----------|
| `JsonRpc` | `string` | no |
| `Error` | `*RpcError` | no |
| `Result` | `interface{}` | no |
| `Id` | `int64` | no |

## RpcResponseRaw

_struct · `client/jsonrpc.go`:31_

| Field | Type | Optional |
|-------|------|----------|
| `JsonRpc` | `string` | no |
| `Error` | `*RpcError` | no |
| `Result` | `json.RawMessage` | no |
| `Id` | `int64` | no |

## server

_struct · `health/http.go`:21_

| Field | Type | Optional |
|-------|------|----------|
| `healthCheckRoute` | `string` | no |
| `readinessCheckRoute` | `string` | no |
| `port` | `int` | no |
| `healthChecks` | `[]CheckFunc` | no |
| `readinessChecks` | `[]CheckFunc` | no |

## Set

_struct · `set/set.go`:7_

| Field | Type | Optional |
|-------|------|----------|
| `values` | `map[T]struct{}` | no |

## Status

_struct · `eventer/client.go`:15_

| Field | Type | Optional |
|-------|------|----------|
| `Status` | `bool` | no |

## SubTransactions

_struct · `blockchain/binance/model.go`:50_

| Field | Type | Optional |
|-------|------|----------|
| `TxHash` | `string` | no |
| `BlockHeight` | `int` | no |
| `TxType` | `string` | no |
| `FromAddr` | `string` | no |
| `ToAddr` | `string` | no |
| `TxAsset` | `string` | no |
| `TxFee` | `string` | no |
| `Value` | `string` | no |

## TextFormatterConfig

_struct · `logging/formatter_strict_text.go`:8_

| Field | Type | Optional |
|-------|------|----------|
| `ForceColors` | `bool` | no |
| `DisableColors` | `bool` | no |
| `DisableTimestamp` | `bool` | no |
| `FullTimestamp` | `bool` | no |
| `TimestampFormat` | `string` | no |
| `DisableSorting` | `bool` | no |

## Token

_struct · `blockchain/binance/model.go`:74_
