---
title: Architecture Overview
category: architecture
tags: [go, library, shared-utilities, trust-wallet]
confidence: high
source: source-analysis
updated: 2026-07-22
---

# Architecture Overview

`github.com/trustwallet/go-libs` is a **public, shared Go library** that provides reusable infrastructure utilities for Trust Wallet's Go backend services. It is a flat-package library (no internal app layers, no cmd/ entrypoints) — every top-level directory is a self-contained, importable package.

## Purpose and Role

This repo is the shared-utilities extraction point from Trust Wallet's Go microservice ecosystem. When a pattern proves stable inside a private project it is promoted here so all services can consume it without reimplementing it. The complementary primitives repo (`github.com/trustwallet/go-primitives`) holds domain-specific types (coin definitions, etc.); go-libs holds infrastructure concerns (HTTP, metrics, caching, DB, messaging, etc.).

## Package Map

| Package | Responsibility |
|---|---|
| `cache/redis` | Redis client wrapper (single + cluster, TLS, Get/Set/Del/MGet/Scan/Watch/SetNX/SetXX, JSON marshal/unmarshal) |
| `client` | HTTP + JSON-RPC client builder with metrics, caching, retry, and request pipeline |
| `database` | GORM/Postgres wrapper with read/write splitting, transaction context propagation, and migration runner |
| `metrics` | Prometheus performance metrics (PerformanceMetric, HttpServerMetric, NullablePerformanceMetric) |
| `middleware` | Gin middleware: metrics, caching, cache-control, sentry error capture, graceful shutdown |
| `logging` | Logrus-based logger with per-component entries and viper-driven config |
| `mq` | AMQP (RabbitMQ) client: connect, declare queues/exchanges, publish, consume with auto-reconnect |
| `worker` | Periodic background worker with builder API, metrics, graceful stop, and interval/one-shot modes |
| `testy` | Integration test helpers (real Postgres + Redis connections from env vars, tagged test filtering) |
| `crypto` | AES-GCM encryption/decryption and ECDSA signing utilities |
| `gin` | Gin HMAC request authentication and server setup helpers |
| `health` | HTTP health-check endpoint handler |
| `httplib` | Low-level HTTP server helper and file downloader |
| `eventer` | Simple event-client abstraction over the logger |
| `ctask` | Concurrent task runner (DoAll pattern for parallel fan-out with error aggregation) |
| `mock` | HTTP mock server for testing (serves static JSON responses) |
| `set` | Generic ordered set data structure |
| `slice` | Generic slice utilities: filter, partition, search |
| `nullable` | Nullable scalar types (NullableString, etc.) for JSON/DB serialization |
| `config/viper` | Viper-based configuration loader helper |
| `blockchain/binance` | Binance chain API and explorer client (legacy blockchain client) |
| `pkg/nullable` | Additional nullable type implementations |

## Design Principles

1. **Flat package layout** — each package is importable independently; no shared internal/ layer. A consumer imports only what it needs.
2. **Interface-first** — key entry points expose interfaces (e.g. `DBContextGetter`, `TrxContextGetter`, `PerformanceMetric`, `Worker`) so consumers can mock or replace implementations in tests.
3. **Options pattern** — configuration is passed via `Option func(...)` variadic parameters (see `cache/redis`, `database`, `mq`, `worker`) to allow backward-compatible extension.
4. **Prometheus integration** — metrics are not global singletons; callers pass a `prometheus.Registerer` to every metric constructor, enabling per-service registry isolation and testability.
5. **Context threading** — database transactions are threaded through `context.Context` (the `trxKey` context key in `database/db.go`) so business logic can stay unaware of transaction boundaries.

## See Also
- [project-structure.md](project-structure.md)
- [dependency-graph.md](dependency-graph.md)
- [layers.md](layers.md)
- [call-graph.md](call-graph.md)
- [crypto and auth](../security/crypto-and-auth.md) <!-- rel:strong -->
