# Data Models & Schemas

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

## DBConfig

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

## memCache

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

## Redis

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

| Field | Type | Optional |
|-------|------|----------|
| `ContractAddress` | `string` | no |
| `Name` | `string` | no |
| `OriginalSymbol` | `string` | no |
| `Owner` | `string` | no |
| `Symbol` | `string` | no |
| `TotalSupply` | `string` | no |

## TokenBalance

_struct · `blockchain/binance/model.go`:65_

| Field | Type | Optional |
|-------|------|----------|
| `Free` | `string` | no |
| `Frozen` | `string` | no |
| `Locked` | `string` | no |
| `Symbol` | `string` | no |

## TransactionData

_struct · `blockchain/binance/model.go`:38_

| Field | Type | Optional |
|-------|------|----------|
| `OrderData` | `struct { Symbol string `json:"symbol"` OrderType string `json:"orderType"` Side string `json:"side"` Price string `json:"price"` Quantity string `json:"quantity"` TimeInForce string `json:"timeInForce"` OrderID string `json:"orderId"` }` | no |

## TransactionsInBlockResponse

_struct · `blockchain/binance/model.go`:12_

| Field | Type | Optional |
|-------|------|----------|
| `BlockHeight` | `int` | no |
| `Tx` | `[]Tx` | no |

## TransactionsResponse

_struct · `blockchain/binance/api/model.go`:4_

| Field | Type | Optional |
|-------|------|----------|
| `Total` | `int` | no |
| `Tx` | `[]Tx` | no |

## Tx

_struct · `blockchain/binance/api/model.go`:11_

| Field | Type | Optional |
|-------|------|----------|
| `Hash` | `string` | no |
| `BlockHeight` | `int` | no |
| `BlockTime` | `int64` | no |
| `Type` | `Type` | no |
| `Fee` | `int` | no |
| `Code` | `int` | no |
| `Source` | `int` | no |
| `Sequence` | `int` | no |
| `Memo` | `string` | no |
| `Log` | `string` | no |
| `Data` | `string` | no |
| `Asset` | `string` | no |
| `Amount` | `float64` | no |
| `FromAddr` | `string` | no |
| `ToAddr` | `string` | no |

## Tx

_struct · `blockchain/binance/model.go`:19_

| Field | Type | Optional |
|-------|------|----------|
| `TxHash` | `string` | no |
| `BlockHeight` | `int` | no |
| `TxType` | `TxType` | no |
| `TimeStamp` | `time.Time` | no |
| `FromAddr` | `interface{}` | no |
| `ToAddr` | `interface{}` | no |
| `Value` | `string` | no |
| `TxAsset` | `string` | no |
| `TxFee` | `string` | no |
| `OrderID` | `string` | no |
| `Code` | `int` | no |
| `Data` | `string` | no |
| `Memo` | `string` | no |
| `Source` | `int` | no |
| `SubTransactions` | `[]SubTransactions` | no |
| `Sequence` | `int` | no |

## worker

_struct · `worker/worker.go`:54_

| Field | Type | Optional |
|-------|------|----------|
| `name` | `string` | no |
| `workerFn` | `func() error` | no |
| `stopFn` | `func() error` | no |
| `options` | `*WorkerOptions` | no |

## WorkerOptions

_struct · `worker/options.go`:9_

| Field | Type | Optional |
|-------|------|----------|
| `Interval` | `time.Duration` | no |
| `RunImmediately` | `bool` | no |
| `RunConsequently` | `bool` | no |
| `PerformanceMetric` | `metrics.PerformanceMetric` | no |

