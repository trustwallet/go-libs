# metrics package

Add dependency to the project

```sh
go get github.com/trustwallet/go-libs/metrics
```

## Features

* The `handler.go` contains very simple method to register Prometheus middleware with `gin-gonic` engine.
* The `metrics.go` is a place for the generic metric exporters.
  * `PerformanceMetric` allows to track generic job performance, start time,
     duration, success and failed executions.
* The `register.go` contains another simple method to register Prometheus collectors 
  with the Default Registerer (which also includes golang specific metrics by default).
* The `pusher.go` configures Prometheus Pushgateway client (see Push mode below).
* The `go-libs/worker` and `go-libs/mq` packages integration for automatic job performance tracking.

## How it works?

The application might have various exporters which are collecting system metrics
as a background processes (in-process worker), standalone services or metric values
updated directly from code logic.

<details> 
<summary>Example of the API metrics (standalone service)</summary><p> 


```go
type APIMetrics struct {
	tickersCache cache.Data
	collectors   map[string]*prometheus.GaugeVec
}

func NewAPIMetrics(tickersCache cache.Data) *APIMetrics {
	s := &APIMetrics{
		tickersCache: tickersCache,
		collectors: map[string]*prometheus.GaugeVec{
			"tickers_cached_total": prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: "market_api",
				Name:      "tickers_cached_total",
				Help:      "Total number of tickers cached",
			}, nil),
		},
	}

	for _, c := range s.collectors {
		metrics.Register(c)
	}

	return s
}

func (s *APIMetrics) Export() error {
	log.Info("export market api metrics")

	tCached, err := s.tickersCache.GetAllTickers()
	if err != nil {
		return errors.Wrap(err, "failed to get all cached tickers")
	}

	// market_api_tickers_cached_total
	s.collectors["tickers_cached_total"].
		WithLabelValues().Set(float64(len(tCached)))

	return nil
}
```
</p></details><br/>
 

The metric exporters should only care about registering the collectors
that they manage internally with Prometheus registerer.

```go
...
for _, c := range s.collectors {
	metrics.Register(c)
}
...
```

### Simple Worker

The example `APIMetrics` (see above) is designed as standalone service,
thus it should be periodically invoked by some in-process worker to collect metric values.

New `go-libs/worker` design allows to define a simple worker for this purpose with 
just a few lines.

```go
metrics := api.NewAPIMetrics(tickersCache)
exportWorker := worker.InitWorker("metrics_exporter",
			worker.DefaultWorkerOptions(config.Default.Metrics.UpdateTime),
			metrics.Export)
...
exportWorker.Start(ctx, waitGroup)
```

### Scrape Mode

The `web` applications that are hosted to serve incoming requests usually
have the `/metrics` or similar endpoint exposed.

One extra line in the main application logic to expose the registered
collectors for Prometheus scrapper.

```go
engine := gin.New()
// ...

metrics.InitHandler(engine, config.Default.Metrics.Path)
```

### Push Mode

The `worker` applications are launched without a capability to serve the incoming requests,
thus `/metrics` endpoint (Prometheus handler) cannot be utilized there. Instead,
the Prometheus Pushgateway should be used, an intermediary service which allows to push metrics
from jobs which cannot be scraped directly.

Assuming, the example `APIMetrics` (see above) is launched as part of the worker app.

It has been integrated with `metrics_exporter` worker already to periodically collect metrics
data from some underlying services.

Next, to make the metric collectors push the values to Prometheus Pushgateway server,
need to initialize the Pushgateway client and set up another worker which pushes 
registered collectors' values.

```go
jobName := "market_worker"
pusher := metrics.NewPusher(config.Default.Metrics.PushgatewayURL, jobName)
pusherWorker := worker.InitWorker(
	"metrics_pusher",
	worker.DefaultWorkerOptions(config.Default.Metrics.UpdateTime),
	pusher.Push)

...
pusherWorker.Start(ctx, wg)
```

📎 When pushing the collectors' values to Pushgateway the `instance` label is 
automatically set from `DYNO` (set by Heroku) or `INSTANCE_ID` (generic variable that
can be set easily) environment variables; otherwise instance is `local`.

### Job Performance Metrics

The _go-libs/metrics_ package contains one very simple, but useful `PerformanceMetric` 
service to collect metrics about any task (short-lived or long-running, doesn't matter).

Its initialization function accepts `namespace` as a first parameter, and any number of the
optional `labelNames` parameters (later when executing exporter methods the same amount of the
label values has to be passed, see example below).

In the following example the metrics will be prefixed with `market_worker_`
(e.g. `market_worker_job_started`) and have `worker` label (e.g. `worker="tickers_cache"`)

```go
metric := metrics.NewPerformanceMetric("market_worker", "worker")
// or shorter...
metric := metrics.NewWorkerPerformanceMetric("market_worker")
```

This metric is designed as a service that should be called from the code on
job started and finished (`Start` and `Duration` functions), and
on success or failure job execution (`Success` and `Failure` functions).

### go-libs/worker integration

The latest `worker` package has already integrated with `metrics` package 
and tracks worker function performance automatically.

Collected metrics from workers will have the `worker` label
set automatically.

The following code is already part of the `worker` package:
```go
func (w *worker) invoke() {
	metric := w.options.PerformanceMetric
	if metric == nil { 			// 1. dummy perf metric
		metric = &metrics.NullablePerformanceMetric{}
	}

	lvs := []string{w.name} 	// 2. worker name as label
	t, _ := metric.Start(lvs) 	// 3. collect worker start time
	err := w.workerFn()  		// 4. invoke worker function
	metric.Duration(t, lvs) 	// 5. collect worker execution duration

	if err != nil {
		metric.Failure(lvs) 	// 6. increment failure counter on error
		log.WithField("worker", w.name).Error(err)
	} else {
		metric.Success(lvs) 	// 7. increment success counter
	}
}
```

As you see, if the `PerformanceMetric` option wasn't initialized it does no-op.

To initialize the worker performance tracking use the following code:

```go
workers := make([]worker.Worker, 0)

// init all workers 
workers = append(workers, worker.InitWorker(...))


// enable performance metric
if config.Default.Metrics.Enabled {
	metric := metrics.NewConsumerPerformanceMetric("market_worker")
	for _, w := range workers {
		w.Options().WithPerformanceMetric(metric)
	}
}

// start all workers
for _, w := range workers {
	w.Start(ctx, waitGroup)
}
```

For all workers for which `WithPerformanceMetric(metric)` was executed 
the performance metrics will be collected (and pushed if `metrics_pusher` is configured).

### go-libs/mq integration

The latest `mq` package has already integrated with `metrics` package 
and tracks functions processing incoming messages performance automatically.

Collected metrics from consumers will have the `consumer_queue` label
set automatically.

It has similar (as worker above) code to wrap
message processing function into performance metrics collector.

The registration is also very similar:

```go
consumers := make([]mq.Consumer, 0)

// init all consumers
consumers = append(consumers, mqClient.InitConsumer(...))

// enable performance metric
if config.Default.Metrics.Enabled {
	metric := metrics.NewConsumerPerformanceMetric("market_consumer")
	for _, c := range consumers {
		c.Options().WithPerformanceMetric(metric)
	}
}
```

If the `PerformanceMetric` option wasn't initialized it does no-op.

## Useful Readings

- [Metric and label naming best practices](https://prometheus.io/docs/practices/naming/)
- [Things to watch out for when instrumenting Prometheus collectors](https://prometheus.io/docs/practices/instrumentation/#things-to-watch-out-for) 
- [When to use the Pushgateway](https://prometheus.io/docs/practices/pushing/)
