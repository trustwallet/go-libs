# Project Structure

<!-- sdd-knowledge-generated -->

## Languages

| Language | Files |
|----------|-------|
| Go | 95 |
| **Total** | **95** |

## Directory Layout

```
├── blockchain/
│   └── binance/
│       ├── api
│       ├── client.go
│       ├── explorer
│       └── model.go
├── cache/
│   └── redis/
│       ├── client_interface.go
│       ├── redis.go
│       └── redis_test.go
├── client/
│   ├── api/
│   │   └── backend
│   ├── client.go
│   ├── client_execute.go
│   ├── client_metrics.go
│   ├── client_metrics_test.go
│   ├── client_test.go
│   ├── client_wrapper.go
│   ├── client_wrapper_test.go
│   ├── clientcache.go
│   ├── clientcache_test.go
│   ├── jsonrpc.go
│   ├── jsonrpc_batch.go
│   ├── jsonrpc_batch_test.go
│   ├── jsonrpc_test.go
│   ├── path.go
│   ├── path_test.go
│   ├── request.go
│   └── request_test.go
├── config/
│   └── viper/
│       └── viper.go
├── crypto/
│   ├── aes.go
│   ├── aes_test.go
│   ├── sign.go
│   └── sign_test.go
├── ctask/
│   ├── do_all.go
│   ├── do_all_test.go
│   ├── doer.go
│   └── doer_test.go
├── database/
│   ├── config.go
│   ├── db.go
│   ├── migrate.go
│   ├── migration_runner_env.go
│   └── mock_db.go
├── eventer/
│   ├── client.go
│   └── log.go
├── gin/
│   ├── hmac.go
│   ├── hmac_test.go
│   └── setup.go
├── health/
│   ├── http.go
│   └── http_test.go
├── httplib/
│   ├── downloader.go
│   └── server.go
├── logging/
│   ├── formatter_strict_text.go
│   ├── logger.go
│   └── logger_test.go
├── metrics/
│   ├── handler.go
│   ├── http_metrics.go
│   ├── metrics.go
│   ├── pusher.go
│   └── register.go
├── middleware/
│   ├── cache.go
│   ├── cache_control.go
│   ├── cache_control_test.go
│   ├── cache_test.go
│   ├── logger.go
│   ├── metrics.go
│   ├── metrics_test.go
│   ├── sentry.go
│   ├── sentry_test.go
│   └── shutdown.go
├── mock/
│   ├── mock.go
│   └── mock_test.go
├── mq/
│   ├── consumer.go
│   ├── exchange.go
│   ├── mq.go
│   ├── options.go
│   └── queue.go
├── pkg/
│   └── nullable/
│       ├── primitives.go
│       └── time.go
├── set/
│   ├── ordered.go
│   ├── ordered_test.go
│   ├── set.go
│   └── set_test.go
├── slice/
│   ├── filter.go
│   ├── filter_test.go
│   ├── partition.go
│   ├── partition_test.go
│   ├── search.go
│   └── search_test.go
├── testy/
│   ├── integration_test_suite.go
│   ├── tagged.go
│   └── tagged_test.go
└── worker/
    ├── metrics/
    │   └── metricspusherworker.go
    ├── options.go
    ├── worker.go
    └── worker_test.go
```

## Features / Modules

| Feature | Files | Classes | Interfaces | Functions | Endpoints |
|---------|-------|---------|------------|-----------|----------|
| blockchain | 6 | 16 | 0 | 3 | 0 |
| cache | 3 | 1 | 2 | 6 | 0 |
| client | 19 | 13 | 1 | 24 | 0 |
| config | 1 | 0 | 0 | 2 | 0 |
| crypto | 4 | 0 | 1 | 9 | 0 |
| ctask | 4 | 3 | 0 | 6 | 0 |
| database | 5 | 12 | 3 | 14 | 0 |
| eventer | 2 | 3 | 0 | 3 | 0 |
| gin | 3 | 1 | 0 | 6 | 0 |
| health | 2 | 1 | 0 | 7 | 0 |
| httplib | 2 | 2 | 2 | 5 | 0 |
| logging | 3 | 1 | 0 | 7 | 0 |
| metrics | 5 | 5 | 3 | 9 | 0 |
| middleware | 10 | 3 | 0 | 22 | 0 |
| mock | 2 | 0 | 0 | 3 | 0 |
| mq | 5 | 7 | 5 | 6 | 0 |
| nullable | 2 | 0 | 0 | 16 | 0 |
| set | 4 | 2 | 0 | 3 | 0 |
| slice | 6 | 0 | 0 | 5 | 0 |
| testy | 3 | 1 | 0 | 8 | 0 |
| worker | 4 | 3 | 2 | 3 | 0 |

## Entry Points

- `httplib/server.go` (Go)

## See Also
- [crypto and auth](../security/crypto-and-auth.md) <!-- rel:strong -->
- [constitution](../constitution.md) <!-- rel:strong -->
- [httplib](../features/httplib.md) <!-- rel:strong -->
- [slice](../features/slice.md) <!-- rel:strong -->
- [health](../features/health.md) <!-- rel:strong -->
