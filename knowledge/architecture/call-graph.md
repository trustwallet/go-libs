---
category: architecture
confidence: low
documentType: explanation
scope: repo
contentHash: aa410351e08b
tags: [api, route, handler, middleware, request]
source: architecture/call-graph.md
verified: 2026-07-22
splitPartIndex: 1
splitPartTotal: 2
canonical: true
---

## Call Graph
<!-- sdd-knowledge-generated -->

> Deterministic call graph extracted via tree-sitter (no LLM). Direct calls are resolved by lexical scope + imports; **dynamic dispatch and ambiguous name matches are withheld** rather than guessed. The `## Calls` table is `EXTRACTED` (a single resolved target). Member calls (`obj.method()`) whose method name resolves to exactly one definition repo-wide are recovered separately under `## Inferred calls` (`INFERRED` — receiver type unverified, but a single plausible target).

## Calls

| Caller | Callee | Example args | Sites |
|--------|--------|--------------|-------|
| `InitClusterWithTLS` | `InitCluster(ctx: context.Context, redisURL: string, opts: ClusterOption): (*Redis, error)` | `(…, url, …)` | 1 |
| `InitClusterWithTLS` | `WithClusterTLS(cfg: *tls.Config): ClusterOption` | `(…)` | 1 |
| `InitWithTLS` | `Init(ctx: context.Context, host: string, opts: Option): (*Redis, error)` | `(…, url, …)` | 1 |
| `InitWithTLS` | `WithTLS(cfg: *tls.Config): Option` | `(…)` | 1 |
| `constructHttpRequest` | `GetBody(body: interface{}): (buf io.ReadWriter, err error)` | `(req.body)` | 1 |
| `Execute` | `populateResultContainer(b: []byte, resultContainer: any): error` | `(b, req.resultContainer)` | 1 |
| `reportMonitoringMetricsIfEnabled` | `getMonitoredPathTemplateIfEnabled(req: *Req): string` | `(req)` | 1 |
| `reportMonitoringMetricsIfEnabled` | `getHttpRespMetricStatus(resp: *http.Response, err: error): string` | `(res, resErr)` | 1 |
| `GetRaw` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `GetWithContext` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `PostRaw` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `PostWithContext` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `InitJSONClient` | `InitClient(baseURL: string, errorHandler: HttpErrorHandler, options: Option): Request` | `(baseUrl, errorHandler, …)` | 1 |
| `InitJSONClient` | `WithExtraHeaders(headers: map[string]string): Option` | `(jsonHeaders)` | 1 |
| `ProxyOption` | `setHttpClientTransportProxy(client: *http.Client, proxyUrl: string): error` | `(httpClient, proxyURL)` | 1 |
| `WithMetricsEnabled` | `newHttpClientMetrics(constLabels: prometheus.Labels): *httpClientMetrics` | `(constLabels)` | 1 |
| `MakeBatchRequests` | `MakeBatches(elements: []interface{}, batchSize: int): (batches [][]interface{})` | `(elements, batchSize)` | 1 |
| `fillDefaultValues` | `genID(): int64` | `()` | 1 |
| `RpcBatchCall` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `RpcCall` | `genID(): int64` | `()` | 1 |
| `RpcCall` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `RpcCallRaw` | `genID(): int64` | `()` | 1 |
| `RpcCallRaw` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `Pathf` | `NewPath(template: string, values: []any): Path` | `(pathTemplate, values)` | 1 |
| `PathStatic` | `NewStaticPath(path: string): Path` | `(path)` | 1 |
| `Load` | `bindEnvs(v: reflect.Value, parts: string)` | `(…)` | 1 |
| `GetRSAPrivateKeyFromFile` | `GetRSAPrivateKey(reader: io.Reader): (*rsa.PrivateKey, error)` | `(file)` | 1 |
| `GetRSAPrivateKeyFromString` | `GetRSAPrivateKey(reader: io.Reader): (*rsa.PrivateKey, error)` | `(…)` | 1 |
| `NewHMACSHA256Signer` | `HMACSHA256(msg: []byte, key: string): ([]byte, error)` | `(msg, key)` | 1 |
| `NewSHA256WithRSASigner` | `SHA256WithRSA(msg: []byte, privateKey: *rsa.PrivateKey): ([]byte, error)` | `(msg, privateKey)` | 1 |
| `DoAll` | `getDoAllConfigWithOptions(opts: DoAllOpt): DoAllConfig` | `(…)` | 1 |
| `Do` | `getConfigWithOptions(opts: DoOpt): DoConfig` | `(…)` | 1 |
| `NewDBGetter` | `newLogLevelFromString(logLevel: LogLevel): (gormLogger.LogLevel, error)` | `(cfg.LogLevel)` | 1 |
| `NewMigrationRunner` | `toMigrationsLogger(logger: logger): *migrationsLogger` | `(runner.logger)` | 1 |
| `Run` | `operationFn` | `(m, operation)` | 1 |
| `RunMigrationsFromEnv` | `NewMigrationRunner(dsn: string, opts: Option): (*MigrationRunner, error)` | `(dsn, …)` | 1 |
| `RunMigrationsFromEnv` | `WithFilesDir(filesDir: string): Option` | `(filesDir)` | 1 |
| `RunMigrationsFromEnv` | `WithLogger(logger: logger): Option` | `(logger)` | 1 |
| `RunMigrationsFromEnv` | `readForceVersion(): (int, error)` | `()` | 1 |
| `Log` | `sendEvents(events: []Event)` | `(events)` | 1 |
| `StartHealthCheckServer` | `handle(handler: *http.ServeMux, route: string, handleFuncs: []CheckFunc)` | `(handler, hcServer.healthCheckRoute, hcServer.healthChecks)` | 2 |
| `GetLoggerForComponent` | `GetLogger(): *logrus.Logger` | `()` | 1 |
| `NewMetricsServer` | `InitHandler(engine: *gin.Engine, path: string)` | `(router, path)` | 1 |
| `NewHttpServerMetric` | `Register(labels: prometheus.Labels, reg: prometheus.Registerer, collectors: prometheus.Collector)` | `(staticLabels, reg, requestStarted, requestDurationSeconds, requestSucceededTotal, requestClientErrTotal, requestServerErrTotal)` | 1 |
| `NewPerformanceMetric` | `Register(labels: prometheus.Labels, reg: prometheus.Registerer, collectors: prometheus.Collector)` | `(staticLabels, reg, executionStarted, executionDurationSeconds, executionSucceededTotal, executionFailedTotal)` | 1 |
| `NewPusher` | `instanceID(): string` | `()` | 1 |
| `NewPusherWithCustomClient` | `instanceID(): string` | `()` | 1 |
| `CacheControl` | `handle(handler: *http.ServeMux, route: string, handleFuncs: []CheckFunc)` | `(c)` | 1 |
| `CacheMiddleware` | `handle(handler: *http.ServeMux, route: string, handleFuncs: []CheckFunc)` | `(c)` | 1 |
| `CacheMiddleware` | `generateKey(c: *gin.Context): string` | `(c)` | 1 |
| `CacheMiddleware` | `newCachedWriter(expire: time.Duration, writer: gin.ResponseWriter, key: string): *cachedWriter` | `(expiration, c.Writer, key)` | 1 |
| `Logger` | `LoggerFormatter(): gin.LogFormatter` | `()` | 1 |
| `GetSentryErrorHandler` | `getBody(res: *http.Response): string` | `(res)` | 1 |
| `Publish` | `publish(amqpChan: *amqp.Channel, exchange: ExchangeName, key: ExchangeKey, body: []byte): error` | `(e.client.amqpChan, e.name, "", body)` | 1 |
| `PublishWithKey` | `publish(amqpChan: *amqp.Channel, exchange: ExchangeName, key: ExchangeKey, body: []byte): error` | `(e.client.amqpChan, e.name, key, body)` | 1 |
| `publish` | `publishWithConfig(amqpChan: *amqp.Channel, exchange: ExchangeName, key: ExchangeKey, body: []byte, cfg: PublishConfig): error` | `(amqpChan, exchange, key, body, PublishConfig{})` | 1 |
| `Publish` | `ExchangeKey` | `(q.name)` | 1 |
| `Publish` | `publish(amqpChan: *amqp.Channel, exchange: ExchangeName, key: ExchangeKey, body: []byte): error` | `(q.client.amqpChan, "", …, body)` | 1 |
| `PublishWithConfig` | `ExchangeKey` | `(q.name)` | 1 |
| `PublishWithConfig` | `publishWithConfig(amqpChan: *amqp.Channel, exchange: ExchangeName, key: ExchangeKey, body: []byte, cfg: PublishConfig): error` | `(q.client.amqpChan, "", …, body, cfg)` | 1 |
| `NewFromValues` | `New(): *Set[T]` | `()` | 1 |
| `Partition` | `Min(values: T): T` | `(…, …)` | 1 |
| `GetDb` | `NewIntegrationTestDb(): (*gorm.DB, error)` | `()` | 1 |
| `GetRedis` | `NewIntegrationTestRedis(): (*redis.Redis, error)` | `()` | 1 |
| `NewIntegrationTestDb` | `MustGetTestDbDSN(): string` | `()` | 1 |
| `getRuntimeTags` | `parseTags(rawTags: string): runtimeTags` | `(…)` | 1 |
| `RequireAllTestTags` | `getRuntimeTags(): runtimeTags` | `()` | 1 |
| `RequireOneOfTestTags` | `getRuntimeTags(): runtimeTags` | `()` | 1 |
| `RequireTestTag` | `getRuntimeTags(): runtimeTags` | `()` | 1 |
| `NewWorkerBuilder` | `DefaultWorkerOptions(interval: time.Duration): *WorkerOptions` | `(…)` | 1 |

## Callers (reverse)

| Symbol | Called by |
|--------|-----------|
| `bindEnvs` | `Load` |
| `DefaultWorkerOptions` | `NewWorkerBuilder` |
| `ExchangeKey` | `Publish`, `PublishWithConfig` |
| `generateKey` | `CacheMiddleware` |
| `genID` | `RpcCall`, `RpcCallRaw`, `fillDefaultValues` |
| `getBody` | `GetSentryErrorHandler` |
| `GetBody` | `constructHttpRequest` |
| `getConfigWithOptions` | `Do` |
| `getDoAllConfigWithOptions` | `DoAll` |
| `getHttpRespMetricStatus` | `reportMonitoringMetricsIfEnabled` |
| `GetLogger` | `GetLoggerForComponent` |
| `getMonitoredPathTemplateIfEnabled` | `reportMonitoringMetricsIfEnabled` |
| `GetRSAPrivateKey` | `GetRSAPrivateKeyFromFile`, `GetRSAPrivateKeyFromString` |
| `getRuntimeTags` | `RequireAllTestTags`, `RequireOneOfTestTags`, `RequireTestTag` |
| `handle` | `CacheControl`, `CacheMiddleware`, `StartHealthCheckServer` |
| `HMACSHA256` | `NewHMACSHA256Signer` |
| `Init` | `InitWithTLS` |
| `InitClient` | `InitJSONClient` |
| `InitCluster` | `InitClusterWithTLS` |
| `InitHandler` | `NewMetricsServer` |
| `instanceID` | `NewPusher`, `NewPusherWithCustomClient` |
| `LoggerFormatter` | `Logger` |
| `MakeBatches` | `MakeBatchRequests` |
| `Min` | `Partition` |
| `MustGetTestDbDSN` | `NewIntegrationTestDb` |
| `New` | `NewFromValues` |
| `newCachedWriter` | `CacheMiddleware` |
| `newHttpClientMetrics` | `WithMetricsEnabled` |
| `NewIntegrationTestDb` | `GetDb` |
| `NewIntegrationTestRedis` | `GetRedis` |
| `newLogLevelFromString` | `NewDBGetter` |
| `NewMigrationRunner` | `RunMigrationsFromEnv` |
| `NewPath` | `Pathf` |
| `NewReqBuilder` | `GetRaw`, `GetWithContext`, `PostRaw`, `PostWithContext`, `RpcBatchCall`, `RpcCall`, `RpcCallRaw` |
| `NewStaticPath` | `PathStatic` |
| `operationFn` | `Run` |
| `parseTags` | `getRuntimeTags` |
| `populateResultContainer` | `Execute` |
| `publish` | `Publish`, `PublishWithKey` |
| `publishWithConfig` | `PublishWithConfig`, `publish` |
| `readForceVersion` | `RunMigrationsFromEnv` |
| `Register` | `NewHttpServerMetric`, `NewPerformanceMetric` |
| `sendEvents` | `Log` |
| `setHttpClientTransportProxy` | `ProxyOption` |
| `SHA256WithRSA` | `NewSHA256WithRSASigner` |
| `toMigrationsLogger` | `NewMigrationRunner` |
| `WithClusterTLS` | `InitClusterWithTLS` |
| `WithExtraHeaders` | `InitJSONClient` |
| `WithFilesDir` | `RunMigrationsFromEnv` |
| `WithLogger` | `RunMigrationsFromEnv` |
| `WithTLS` | `InitWithTLS` |

## Inferred calls (member, unique name)

> `INFERRED` (confidence 0.9): `obj.method()` where `method` resolves to exactly one definition repo-wide. Receiver type is not resolved; common method names (init/update/onCreate) remain withheld as ambiguous. Use SCIP (`--scip`) for type-resolved member calls at full certainty.

| Caller | Callee | Example args | Sites |
|--------|--------|--------------|-------|
| `GetTransactionsByAddress` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `GetTransactionsByAddress` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `GetTransactionsByAddress` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `GetTransactionsByAddress` | `PathStatic(path: string): *ReqBuilder` | `("bc/api/v1/txs")` | 1 |
| `GetTransactionsByAddress` | `Query(query: url.Values): *ReqBuilder` | `(params)` | 1 |
| `GetTransactionsByAddress` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `InitClient` | `InitJSONClient(baseUrl: string, errorHandler: HttpErrorHandler, options: Option): Request` | `(url, errorHandler)` | 1 |
| `FetchAccountMeta` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `FetchAccountMeta` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `FetchAccountMeta` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `FetchAccountMeta` | `Pathf(pathTemplate: string, values: any): *ReqBuilder` | `("/api/v1/account/%s", address)` | 1 |
| `FetchAccountMeta` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `FetchMarketPairs` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `FetchMarketPairs` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `FetchMarketPairs` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `FetchMarketPairs` | `PathStatic(path: string): *ReqBuilder` | `("/api/v1/markets")` | 1 |
| `FetchMarketPairs` | `Query(query: url.Values): *ReqBuilder` | `(params)` | 1 |
| `FetchMarketPairs` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `FetchNodeInfo` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `FetchNodeInfo` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `FetchNodeInfo` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `FetchNodeInfo` | `PathStatic(path: string): *ReqBuilder` | `("/api/v1/node-info")` | 1 |
| `FetchNodeInfo` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `FetchTokens` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `FetchTokens` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `FetchTokens` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `FetchTokens` | `PathStatic(path: string): *ReqBuilder` | `("/api/v1/tokens")` | 1 |
| `FetchTokens` | `Query(query: url.Values): *ReqBuilder` | `(params)` | 1 |
| `FetchTokens` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `FetchTransactionsByAddressAndTokenID` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `FetchTransactionsByAddressAndTokenID` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `FetchTransactionsByAddressAndTokenID` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `FetchTransactionsByAddressAndTokenID` | `PathStatic(path: string): *ReqBuilder` | `("/api/v1/transactions")` | 1 |
| `FetchTransactionsByAddressAndTokenID` | `Query(query: url.Values): *ReqBuilder` | `(params)` | 1 |
| `FetchTransactionsByAddressAndTokenID` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `FetchTransactionsInBlock` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `FetchTransactionsInBlock` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `FetchTransactionsInBlock` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `FetchTransactionsInBlock` | `Pathf(pathTemplate: string, values: any): *ReqBuilder` | `("api/v2/transactions-in-block/%d", blockNumber)` | 1 |
| `FetchTransactionsInBlock` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `InitClient` | `InitJSONClient(baseUrl: string, errorHandler: HttpErrorHandler, options: Option): Request` | `(url, errorHandler, …)` | 1 |
| `InitClient` | `WithExtraHeader(key: string): Option` | `("apikey", apiKey)` | 1 |
| `FetchBep2Assets` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `FetchBep2Assets` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `FetchBep2Assets` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `FetchBep2Assets` | `PathStatic(path: string): *ReqBuilder` | `("/api/v1/assets")` | 1 |
| `FetchBep2Assets` | `Query(query: url.Values): *ReqBuilder` | `(params)` | 1 |
| `FetchBep2Assets` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `InitClient` | `InitJSONClient(baseUrl: string, errorHandler: HttpErrorHandler, options: Option): Request` | `(url, errorHandler)` | 1 |
| `HealthCheck` | `IsAvailable(ctx: context.Context): bool` | `(ctx)` | 1 |
| `HealthCheck` | `New(): *Set[T]` | `("redis is not available")` | 1 |
| `Reconnect` | `reconnectCluster(ctx: context.Context, redisURL: string): error` | `(ctx, host)` | 1 |
| `GetAssetInfo` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `GetAssetInfo` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `GetAssetInfo` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `GetAssetInfo` | `Pathf(pathTemplate: string, values: any): *ReqBuilder` | `("/v1/assets/%s", assetID)` | 1 |
| `GetAssetInfo` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `InitClient` | `InitJSONClient(baseUrl: string, errorHandler: HttpErrorHandler, options: Option): Request` | `(url, errorHandler)` | 1 |
| `constructHttpRequest` | `GetURL(path: string, query: url.Values): string` | `(…, req.query)` | 1 |
| `constructHttpRequest` | `setRequestHeaders(httpRequest: *http.Request, req: *Req)` | `(request, req)` | 1 |
| `Execute` | `constructHttpRequest(ctx: context.Context, req: *Req): (*http.Request, error)` | `(ctx, req)` | 1 |
| `Execute` | `reportMonitoringMetricsIfEnabled(startTime: time.Time, request: *http.Request, req: *Req, res: *http.Response, resErr: error)` | `(startTime, request, req, res, err)` | 1 |
| `Execute` | `HttpErrorHandler` | `(res, …)` | 1 |
| `GetURL` | `GetBase(path: string): string` | `(path)` | 1 |
| `reportMonitoringMetricsIfEnabled` | `GetURL(path: string, query: url.Values): string` | `(…, nil)` | 1 |
| `reportMonitoringMetricsIfEnabled` | `metricsEnabled(): bool` | `()` | 1 |
| `reportMonitoringMetricsIfEnabled` | `observeDuration(url: string, startTime: time.Time)` | `(url, method, name, startTime)` | 1 |
| `reportMonitoringMetricsIfEnabled` | `observeResult(url: string)` | `(url, method, name, status)` | 1 |
| `Get` | `GetWithContext(ctx: context.Context, result: interface{}, path: string, query: url.Values): error` | `(…, result, path, query)` | 1 |
| `GetRaw` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `GetRaw` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `GetRaw` | `pathMetricEnabled(enabled: bool): *ReqBuilder` | `(false)` | 1 |
| `GetRaw` | `PathStatic(path: string): *ReqBuilder` | `(path)` | 1 |
| `GetRaw` | `Query(query: url.Values): *ReqBuilder` | `(query)` | 1 |
| `GetWithContext` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(ctx, …)` | 1 |
| `GetWithContext` | `Method(method: string): *ReqBuilder` | `(http.MethodGet)` | 1 |
| `GetWithContext` | `pathMetricEnabled(enabled: bool): *ReqBuilder` | `(false)` | 1 |
| `GetWithContext` | `PathStatic(path: string): *ReqBuilder` | `(path)` | 1 |
| `GetWithContext` | `Query(query: url.Values): *ReqBuilder` | `(query)` | 1 |
| `GetWithContext` | `WriteTo(resultContainer: any): *ReqBuilder` | `(result)` | 1 |
| `Post` | `PostWithContext(ctx: context.Context, result: interface{}, path: string, body: interface{}): error` | `(…, result, path, body)` | 1 |
| `PostRaw` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `PostRaw` | `Body(body: any): *ReqBuilder` | `(body)` | 1 |
| `PostRaw` | `Method(method: string): *ReqBuilder` | `(http.MethodPost)` | 1 |
| `PostRaw` | `pathMetricEnabled(enabled: bool): *ReqBuilder` | `(false)` | 1 |
| `PostRaw` | `PathStatic(path: string): *ReqBuilder` | `(path)` | 1 |
| `PostWithContext` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(ctx, …)` | 1 |
| `PostWithContext` | `Body(body: any): *ReqBuilder` | `(body)` | 1 |
| `PostWithContext` | `Method(method: string): *ReqBuilder` | `(http.MethodPost)` | 1 |
| `PostWithContext` | `pathMetricEnabled(enabled: bool): *ReqBuilder` | `(false)` | 1 |
| `PostWithContext` | `PathStatic(path: string): *ReqBuilder` | `(path)` | 1 |
| `PostWithContext` | `WriteTo(resultContainer: any): *ReqBuilder` | `(result)` | 1 |
| `InitClient` | `metricsEnabled(): bool` | `()` | 1 |
| `InitClient` | `Register(labels: prometheus.Labels, reg: prometheus.Registerer, collectors: prometheus.Collector)` | `(client.httpMetrics)` | 1 |
| `ProxyOption` | `New(): *Set[T]` | `("unable to set proxy: httpclient is not )` | 1 |
| `setHttpClientTransportProxy` | `New(): *Set[T]` | `("empty proxy url")` | 2 |
| `SetProxy` | `New(): *Set[T]` | `("empty proxy url")` | 1 |
| `TimeoutOption` | `New(): *Set[T]` | `("unable to set timeout: httpclient is no)` | 1 |
| `generateKey` | `GetBase(path: string): string` | `(path)` | 1 |
| `getCache` | `New(): *Set[T]` | `("validator cache: invalid cache key")` | 3 |
| `GetWithCache` | `GetWithCacheAndContext(ctx: context.Context, result: interface{}, path: string, query: url.Values, cache: time.Duration): error` | `(…, result, path, query, cache)` | 1 |
| `GetWithCacheAndContext` | `GetWithContext(ctx: context.Context, result: interface{}, path: string, query: url.Values): error` | `(ctx, result, path, query)` | 1 |
| `init` | `New(): *Set[T]` | `(…, …)` | 1 |
| `PostWithCache` | `PostWithCacheAndContext(ctx: context.Context, result: interface{}, path: string, body: interface{}, cache: time.Duration): error` | `(…, result, path, body, cache)` | 1 |
| `PostWithCacheAndContext` | `PostWithContext(ctx: context.Context, result: interface{}, path: string, body: interface{}): error` | `(ctx, result, path, body)` | 1 |
| `setCache` | `New(): *Set[T]` | `(…)` | 1 |
| `RpcBatchCall` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `RpcBatchCall` | `fillDefaultValues(): RpcRequests` | `()` | 1 |
| `RpcBatchCall` | `Body(body: any): *ReqBuilder` | `(…)` | 1 |
| `RpcBatchCall` | `Method(method: string): *ReqBuilder` | `(http.MethodPost)` | 1 |
| `RpcBatchCall` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `RpcCall` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `RpcCall` | `GetObject(toType: interface{}): error` | `(result)` | 1 |
| `RpcCall` | `Body(body: any): *ReqBuilder` | `(req)` | 1 |
| `RpcCall` | `Method(method: string): *ReqBuilder` | `(http.MethodPost)` | 1 |
| `RpcCall` | `MetricName(name: string): *ReqBuilder` | `(method)` | 1 |
| `RpcCall` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `RpcCallRaw` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `RpcCallRaw` | `Body(body: any): *ReqBuilder` | `(req)` | 1 |
| `RpcCallRaw` | `Method(method: string): *ReqBuilder` | `(http.MethodPost)` | 1 |
| `RpcCallRaw` | `MetricName(name: string): *ReqBuilder` | `(method)` | 1 |
| `RpcCallRaw` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `bindEnvs` | `Type` | `()` | 1 |
| `GetRSAPrivateKey` | `New(): *Set[T]` | `("decoded key is empty")` | 3 |
| `HMACSHA256` | `Write(data: []byte): (int, error)` | `(msg)` | 1 |
| `HMACSHA256` | `New(): *Set[T]` | `(sha256.New, …)` | 1 |
| `SHA256WithRSA` | `Write(data: []byte): (int, error)` | `(msg)` | 1 |
| `SHA256WithRSA` | `New(): *Set[T]` | `("private key is empty")` | 2 |
| `NewDBGetter` | `applyDefaultValue()` | `()` | 1 |
| `NewDBGetter` | `Register(labels: prometheus.Labels, reg: prometheus.Registerer, collectors: prometheus.Collector)` | `(dbresolver.Config{ Sources: )` | 1 |
| `NewDBGetter` | `New(): *Set[T]` | `(…, gormLogger.Config{ SlowThreshold: )` | 2 |
| `NewMigrationRunner` | `New(): *Set[T]` | `(…, dsn)` | 1 |
| `runDown` | `Info(null: interface{})` | `(…)` | 1 |
| `runForce` | `Info(null: interface{})` | `(…)` | 1 |
| `runUp` | `Info(null: interface{})` | `("running migrate UP")` | 2 |
| `RunMigrationsFromEnv` | `Info(null: interface{})` | `(…)` | 3 |
| `RunMigrationsFromEnv` | `Version(): (version uint, dirty bool, err error)` | `()` | 2 |
| `Init` | `InitJSONClient(baseUrl: string, errorHandler: HttpErrorHandler, options: Option): Request` | `(url, middleware.SentryErrorHandler)` | 1 |
| `SendBatch` | `Execute(ctx: context.Context, req: *Req): ([]byte, error)` | `(…, …)` | 1 |
| `SendBatch` | `Body(body: any): *ReqBuilder` | `(events)` | 1 |
| `SendBatch` | `Method(method: string): *ReqBuilder` | `(http.MethodPost)` | 1 |
| `SendBatch` | `NewReqBuilder(): *ReqBuilder` | `()` | 1 |
| `SendBatch` | `PathStatic(path: string): *ReqBuilder` | `("")` | 1 |
| `SendBatch` | `WriteTo(resultContainer: any): *ReqBuilder` | `(…)` | 1 |
| `sendEvents` | `SendBatch(events: []Event): (status Status, err error)` | `(events)` | 1 |
| `SignedHandler` | `verifySignature(msg: []byte, sig: string): error` | `(…, sig)` | 1 |
| `SignedHandler` | `New(): *Set[T]` | `("cannot extract plaintext")` | 3 |
| `verifySignature` | `Write(data: []byte): (int, error)` | `(msg)` | 1 |
| `verifySignature` | `New(): *Set[T]` | `(sha256.New, signatureKey)` | 1 |
| `SetupGracefulServeWithUnixFile` | `Info(null: interface{})` | `("Server Shutdown: ", err)` | 2 |
| `SetupGracefulServeWithUnixFile` | `Remove(val: T)` | `(unixFile)` | 1 |
| `SetupGracefulShutdown` | `Info(null: interface{})` | `("Server Shutdown: ", err)` | 5 |
| `handle` | `WriteHeader(code: int)` | `(http.StatusOK)` | 1 |
| `StartHealthCheckServer` | `Info(null: interface{})` | `("server shutdown: ", err)` | 1 |
| `Run` | `serve(ctx: context.Context, wg: *sync.WaitGroup)` | `(ctx, wg)` | 1 |
| `serve` | `Info(null: interface{})` | `("Starting the API server")` | 3 |
| `init` | `New(): *Set[T]` | `()` | 1 |
| `NewMetricsServer` | `NewHTTPServer(router: http.Handler, port: string): Server` | `(router, port)` | 1 |
| `Close` | `Delete(ctx: context.Context, key: string): error` | `()` | 1 |
| `NewMetricsPusherClient` | `WithExtraHeader(key: string): Option` | `("X-API-Key", key)` | 1 |
| `NewPusher` | `New(): *Set[T]` | `(pushgatewayURL, jobName)` | 1 |
| `NewPusherWithCustomClient` | `New(): *Set[T]` | `(pushgatewayURL, …)` | 1 |
| `Register` | `GetLogger(): *logrus.Logger` | `()` | 1 |
| `CacheMiddleware` | `deleteCache(key: string)` | `(key)` | 2 |
| `CacheMiddleware` | `Write(data: []byte): (int, error)` | `(mc.Data)` | 1 |
| `CacheMiddleware` | `WriteHeader(code: int)` | `(mc.Status)` | 1 |
| `deleteCache` | `Delete(ctx: context.Context, key: string): error` | `(key)` | 1 |
| `getCache` | `New(): *Set[T]` | `("validator cache: failed to cast cache t)` | 2 |
| `init` | `New(): *Set[T]` | `(…, …)` | 1 |
| `setCache` | `New(): *Set[T]` | `(…)` | 1 |
| `Write` | `New(): *Set[T]` | `("validator cache: failed to marshal cach)` | 1 |
| `WriteString` | `New(): *Set[T]` | `(…)` | 3 |
| `MetricsMiddleware` | `ClientError(labelValues: string)` | `(…)` | 1 |
| `MetricsMiddleware` | `NewHttpServerMetric(namespace: string, labelNames: []string, staticLabels: prometheus.Labels, reg: prometheus.Registerer): HttpServerMetric` | `(namespace, []string{labelPath, labelMethod, labelSt, labels, reg)` | 1 |
| `MetricsMiddleware` | `ServerError(labelValues: string)` | `(…)` | 1 |
| `SetupGracefulShutdown` | `Info(null: interface{})` | `("Shutdown timeout: ...", timeout)` | 2 |
| `consume` | `getRemainingRetries(delivery: amqp.Delivery): int32` | `(msg)` | 1 |
| `consume` | `process(queueName: string, body: []byte): error` | `(queueName, msg.Body)` | 1 |
| `consume` | `PublishWithConfig(body: []byte, cfg: PublishConfig): error` | `(msg.Body, PublishConfig{ MaxRetries: nullabl)` | 1 |
| `consume` | `Int(i: int): *int` | `(…)` | 1 |
| `messageChannel` | `getSanitizedPrefetchCount(): int` | `()` | 1 |
| `process` | `Failure(labelValues: string)` | `()` | 1 |
| `process` | `Process(m: Message): error` | `(body)` | 1 |
| `Start` | `consume(ctx: context.Context)` | `(ctx)` | 1 |
| `Start` | `messageChannel(): (<-chan amqp.Delivery, error)` | `()` | 1 |
| `HealthCheck` | `New(): *Set[T]` | `("connection is closed")` | 1 |
| `InitConsumer` | `InitQueue(name: QueueName): Queue` | `(queueName)` | 1 |
| `ListenConnection` | `Info(null: interface{})` | `("start listen connection")` | 3 |
| `ListenConnection` | `initNotifyCloseListeners(): (<-chan *amqp.Error, <-chan *amqp.Error)` | `()` | 2 |
| `ListenConnection` | `reconnectWithRetry(ctx: context.Context): error` | `(ctx)` | 1 |
| `ListenConnectionAsync` | `ListenConnection(ctx: context.Context): error` | `(ctx)` | 1 |
| `reconnectWithRetry` | `Info(null: interface{})` | `("Connecting to MQ... Attempt ", …)` | 2 |
| `reconnectWithRetry` | `reconnect(): error` | `()` | 1 |
| `StartConsumers` | `AddConnectionClient(connClient: ConnectionClient)` | `(consumer)` | 1 |
| `Declare` | `DeclareWithConfig(cfg: DeclareConfig): error` | `(DeclareConfig{Durable: true})` | 1 |
| `MarshalJSON` | `ToSlice(): []T` | `()` | 1 |
| `UnmarshalJSON` | `Clear()` | `()` | 1 |
| `containsAll` | `contains(targetTag: string): bool` | `(targetTag)` | 1 |
| `containsAny` | `contains(targetTag: string): bool` | `(targetTag)` | 1 |
| `RequireAllTestTags` | `containsAll(targetTags: string): bool` | `(…)` | 1 |
| `RequireOneOfTestTags` | `containsAny(targetTags: string): bool` | `(…)` | 1 |
| `RequireTestTag` | `contains(targetTag: string): bool` | `(testTag)` | 1 |
| `NewMetricsPusherWorker` | `NewWorkerBuilder(name: string, workerFn: func() error): Builder` | `("metrics_pusher", pusher.Push)` | 1 |
| `NewMetricsPusherWorker` | `WithOptions(options: *WorkerOptions): Builder` | `(options)` | 1 |
| `NewMetricsPusherWorker` | `WithStop(stopFn: func() error): Builder` | `(pusher.Close)` | 1 |
| `hold` | `Info(null: interface{})` | `("worker started, but won't be executed")` | 3 |
| `invoke` | `Failure(labelValues: string)` | `()` | 1 |
| `start` | `Info(null: interface{})` | `("run immediately")` | 4 |
| `Start` | `hold(ctx: context.Context, wg: *sync.WaitGroup)` | `(ctx, wg)` | 1 |
| `start` | `invoke()` | `()` | 2 |
| `Start` | `start(ctx: context.Context, wg: *sync.WaitGroup)` | `(ctx, wg)` | 1 |

## See Also
- [client](../features/client.md) <!-- rel:strong -->
- [mq](../features/mq.md) <!-- rel:related -->
- [database](../features/database.md) <!-- rel:related -->
- [metrics](../features/metrics.md) <!-- rel:related -->
- [middleware](../features/middleware.md) <!-- rel:related -->
