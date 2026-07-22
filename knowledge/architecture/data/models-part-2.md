---
category: architecture
subcategory: data
confidence: low
documentType: explanation
scope: repo
contentHash: f033d6f6e5d6
tags: [metrics, prometheus]
source: architecture/data/models.md
verified: 2026-07-22
splitPartIndex: 2
splitPartTotal: 5
canonical: false
synthetic: split-part
---

## Data Models & Schemas (Part 2)

_struct · `database/config.go`:42_

| Field | Type | Optional |
|-------|------|----------|
| `Url` | `string` | no |
| `ReadonlyUrl` | `*string` | no |
| `LogLevel` | `LogLevel` | no |
| `ConnPool` | `*DBConnPool` | no |

## DBConnPool

_struct · `database/config.go`:34_

| Field | Type | Optional |
|-------|------|----------|
| `MaxIdleConns` | `int` | no |
| `ConnMaxIdleTime` | `time.Duration` | no |
| `MaxOpenConns` | `int` | no |
| `ConnMaxLifetime` | `time.Duration` | no |

## DBGetter

_struct · `database/db.go`:31_

| Field | Type | Optional |
|-------|------|----------|
| `db` | `*gorm.DB` | no |

## DeclareConfig

_struct · `mq/queue.go`:54_

| Field | Type | Optional |
|-------|------|----------|
| `Durable` | `bool` | no |
| `AutoDelete` | `bool` | no |
| `Exclusive` | `bool` | no |
| `NoWait` | `bool` | no |
| `Args` | `map[string]interface{}` | no |

## DoAllConfig

_struct · `ctask/do_all.go`:11_

| Field | Type | Optional |
|-------|------|----------|
| `WorkerNum` | `int` | no |

## DoAllResp

_struct · `ctask/do_all.go`:15_

| Field | Type | Optional |
|-------|------|----------|
| `Result` | `R` | no |
| `Error` | `error` | no |

## DoConfig

_struct · `ctask/doer.go`:11_

| Field | Type | Optional |
|-------|------|----------|
| `WorkerNum` | `int` | no |

## downloader

_struct · `httplib/downloader.go`:15_

| Field | Type | Optional |
|-------|------|----------|
| `client` | `http.Client` | no |
| `bytesSizeLimit` | `int64` | no |

## Event

_struct · `eventer/client.go`:19_

| Field | Type | Optional |
|-------|------|----------|
| `Name` | `string` | no |
| `CreatedAt` | `int64` | no |
| `Params` | `map[string]string` | no |

## exchange

_struct · `mq/exchange.go`:5_

| Field | Type | Optional |
|-------|------|----------|
| `name` | `ExchangeName` | no |
| `client` | `*Client` | no |

## HmacVerifier

_struct · `gin/hmac.go`:20_

| Field | Type | Optional |
|-------|------|----------|
| `keys` | `[][]byte` | no |
| `sigFN` | `StrFromCtx` | no |
| `sigEncoder` | `func([]byte) string` | no |

## httpClientMetrics

_struct · `client/client_metrics.go`:25_

| Field | Type | Optional |
|-------|------|----------|
| `durationSeconds` | `*prometheus.HistogramVec` | no |
| `requestTotal` | `*prometheus.CounterVec` | no |

## HttpError

_struct · `client/client.go`:33_

| Field | Type | Optional |
|-------|------|----------|
| `StatusCode` | `int` | no |
| `URL` | `url.URL` | no |
| `Body` | `[]byte` | no |

## httpServerMetric

_struct · `metrics/http_metrics.go`:24_

| Field | Type | Optional |
|-------|------|----------|
| `requestStarted` | `*prometheus.GaugeVec` | no |
| `requestDurationSeconds` | `*prometheus.HistogramVec` | no |
| `requestSucceededTotal` | `*prometheus.CounterVec` | no |
| `requestClientErrTotal` | `*prometheus.CounterVec` | no |
| `requestServerErrTotal` | `*prometheus.CounterVec` | no |

## IntegrationTestSuite

_struct · `testy/integration_test_suite.go`:30_

| Field | Type | Optional |
|-------|------|----------|
| `db` | `*gorm.DB` | no |
| `redis` | `*redis.Redis` | no |

## MarketPair

_struct · `blockchain/binance/model.go`:83_

| Field | Type | Optional |
|-------|------|----------|
| `BaseAssetSymbol` | `string` | no |
| `LotSize` | `string` | no |
| `QuoteAssetSymbol` | `string` | no |
| `TickSize` | `string` | no |
