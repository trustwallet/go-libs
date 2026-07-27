---
title: Worker Pattern
category: patterns
tags: [worker, background, periodic, builder, graceful-shutdown, context]
confidence: high
source: worker/worker.go, worker/options.go
updated: 2026-07-22
---

# Worker Pattern

## Overview

The `worker` package provides a reusable **periodic background task** abstraction. It is the canonical way to implement scheduled work in Trust Wallet Go services.

## Builder API

```go
w := worker.NewWorkerBuilder("my_worker", func() error {
    // business logic
    return doWork()
}).
    WithOptions(worker.DefaultWorkerOptions(5 * time.Minute)).
    WithStop(func() error {
        // cleanup on shutdown
        return nil
    }).
    Build()

var wg sync.WaitGroup
w.Start(ctx, &wg)
wg.Wait()
```

## `WorkerOptions`

| Field | Default | Purpose |
|---|---|---|
| `Interval` | 1 min (configurable) | How often the worker function runs |
| `RunImmediately` | depends | Run once immediately before the first tick |

Setting `Interval = -1` puts the worker in "hold" mode — it blocks until the context is cancelled (for one-shot or event-driven workers).

## Graceful Shutdown

Workers respect context cancellation. When `ctx` is cancelled:
1. The current `workerFn()` is allowed to finish.
2. The optional `stopFn` is called.
3. `wg.Done()` is signalled.

## Metrics Integration

Worker execution is automatically tracked via the `worker/metrics` package — each `Start` call registers a Prometheus counter and histogram scoped to the worker name.

## See Also
- [worker.md](../features/worker.md)
- [metrics-system.md](../architecture/metrics-system.md)
- [overview](../architecture/overview.md) <!-- rel:strong -->
- [models](../architecture/data/models.md) <!-- rel:strong -->
- [database transactions](../architecture/database-transactions.md) <!-- rel:related -->
