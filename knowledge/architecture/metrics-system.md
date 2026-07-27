---
title: Metrics System
category: architecture
tags: [prometheus, metrics, performance, http, observability]
confidence: high
source: metrics/metrics.go, metrics/http_metrics.go, metrics/register.go, worker/metrics/
updated: 2026-07-22
---

# Metrics System

## Overview

The `metrics` package provides a **Prometheus-based performance metrics abstraction** that is used across the library (middleware, worker, HTTP client). All metric constructors accept a `prometheus.Registerer` to avoid global registry conflicts in multi-service deployments.

## Key Types

### `PerformanceMetric` interface

```go
type PerformanceMetric interface {
    Start(labelValues ...string) time.Time
    Duration(start time.Time, labelValues ...string)
    Success(labelValues ...string)
    Failure(labelValues ...string)
}
```

Tracks four counters/gauges per labeled operation:
- `execution_started` — gauge: last start time (Unix)
- `execution_duration_seconds` — histogram: latency
- `execution_succeeded_total` — counter
- `execution_failed_total` — counter

### `NullablePerformanceMetric`

A no-op implementation of `PerformanceMetric`. Use when metrics are disabled or in tests where you do not want Prometheus counters incremented.

### `HttpServerMetric`

Specialization for HTTP server handlers — adds `client_error_total` (4xx) alongside the standard success/failure counters. Used by `middleware.MetricsMiddleware`.

## Usage Pattern

```go
reg := prometheus.NewRegistry()
metric := metrics.NewPerformanceMetric("my_service", []string{"operation"}, nil, reg)

start := metric.Start("import_assets")
defer func() {
    metric.Duration(start, "import_assets")
    if err != nil {
        metric.Failure("import_assets")
    } else {
        metric.Success("import_assets")
    }
}()
```

## Worker Metrics

`worker/metrics/` contains the worker-specific Prometheus metric wrappers. The `worker` package uses these to track worker execution counts and durations, separate from the general `PerformanceMetric` (which is for any timed operation).

## See Also
- [metrics.md](../features/metrics.md)
- [middleware.md](../features/middleware.md)
- [worker.md](../features/worker.md)
- [worker pattern](../patterns/worker-pattern.md) <!-- rel:strong -->
- [integration testing](../tests/integration-testing.md) <!-- rel:related -->
