# metrics package

Add dependency to the project

```sh
go get github.com/trustwallet/go-libs/metrics
```

## Features

* The `handler.go` contains very simple method to register Prometheus middleware with `gin-gonic` engine.
* The `register.go` contains another simple method to register Prometheus collectors
  (metrics) with the Default Registerer (which also includes golang specific metrics by default).
* The `exporters.go` is a place for the generic metrics exporters.
  * `JobPerformanceExporter` allows to track generic job performance by tracking its start time and duration.
* The `pusher.go` configures Prometheus Pushgateway client.


## How it works?

The application might have various exporters which are collecting system metrics as a background processes (in-process worker), standalone services or metric values updated by some code logic. 

<details> 
<summary>Example of the API metrics exporter (standalone service)</summary><p> 


```go
type APIMetricsExporter struct {
	tickersCache cache.Data
	metrics      map[string]*prometheus.GaugeVec
}

func NewAPIMetricsExporter(tickersCache cache.Data) *APIMetricsExporter {
	s := &APIMetricsExporter{
		tickersCache: tickersCache,
		metrics: map[string]*prometheus.GaugeVec{
			"tickers_cached_total": prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: "market_api",
				Name:      "tickers_cached_total",
				Help:      "Total number of tickers cached",
			}, nil),
		},
	}

	for _, c := range s.metrics {
		metrics.Register(c)
	}

	return s
}

func (s *APIMetricsExporter) Export() error {
	log.Info("export market api metrics")

	tCached, err := s.tickersCache.GetAllTickers()
	if err != nil {
		return errors.Wrap(err, "failed to get all cached tickers")
	}

	// market_api_tickers_cached_total
	s.metrics[tickersCachedTotalKey].
		WithLabelValues().Set(float64(len(tCached)))

	return nil
}
```
</p></details><br/>
 

The exporters should only care about registering the collectors that they manage internally with Prometheus registerer.

```go
for _, c := range s.metrics {
	metrics.Register(c)
}
```

### Simple Worker

The example `APIMetricsExporter` (see above) is designed as standalone service,
thus it should be periodically invoked by some in-process worker.

The _go-libs/worker_ package `SimpleWorker` was introduced which allows
to simply configure interval and the calling method.

```go
exporter := NewAPIMetricsExporter(tickersCache)
exportWorker := worker.NewSimpleWorker(config.Default.Metrics.UpdateTime, exporter.Export)

exportWorker.Start(ctx, waitGroup)
```


### Scrape Mode

The `web` applications that are hosted to serve incoming requests usually have the `/metrics` or similar endpoint exposed.

It's only one extra line in the main application logic to expose the registered collectors for Prometheus scrapper.

```go
engine := gin.New()
// ...

metrics.Handler(engine, config.Default.Metrics.Path)
```

### Push Mode

The `worker` applications are launched without a capability to serve the incoming requests, thus `/metrics` endpoint (Prometheus handler) cannot be utilized there. Instead, the Prometheus Pushgateway should be used, an intermediary service which allows to push metrics from jobs which cannot be scraped directly.

Assuming, the example `APIMetricsExporter` (see above) is launched as part of the worker application.

It has been integrated with `SimpleWorker` already to export data from some underlying services (collect system metrics).

Next, to make these collectors push their values to Prometheus Pushgateway server,
need to initialize the Pushgateway client and set up another worker which pushes 
registered collectors' values.

```go
pusher := metrics.NewPusher(config.Default.Metrics.PushgatewayURL, "metrics_worker")
pusherWorker := worker.NewSimpleWorker(config.Default.Metrics.UpdateTime, pusher.Push)

pusherWorker.Start(ctx, waitGroup)
```

📎 When pushing the collectors' values to Pushgateway the `instance` label is 
automatically set from `DYNO` (set by Heroku) or `INSTANCE_ID` (generic variable that can be set easily) environment variables; otherwise instance is `local`.

### Job Performance Exporter

The _go-libs/metrics_ package contains one very simple, but useful `JobPerformanceExporter` service to collect metrics about any background task (short-lived or long-running, doesn't matter).

Its initialization function accepts `namespace` as a first parameter, and any number of the optional `labelNames` parameters (later when executing exporter methods the same amount of the
label values has to be passed, see example below).

In the following example the metrics will be prefixed with `market_worker_`
(e.g. `market_worker_job_started`) and have `worker` label (e.g. `worker="tickers_cache"`)

```go
exporter := metrics.NewJobPerformanceExporter("market_worker", "worker")
```

The exporter is designed as a service that should be called from the code on job started and finished, and provides two relevant functions `Start(...)` and `Duration(...)`.

The elegant one-liner to track job performance. The trick here that arguments passed 
to deferred code are resolved instantly, and then duration is calculated before
returning from the function.

```go
func (j *job) invokeJob() {
	defer j.exporter.Duration(Start("tickers_cache"))

    // The rest of the job logic goes here
}
```

📎 In some cases application runs multiple in-process workers that should share the 
same `JobPerformanceExporter`; otherwise Prometheus will return an error while 
attempting to register several exporters with the same collectors (by key). 
For this reason there is a `GetJobPerformanceExporter(...)` function with the 
same interface as `NewJobPerformanceExporter(...)` which acts as a Singleton
providing pointer to the same exporter instance.

## Useful Readings

- [Metric and label naming best practices](https://prometheus.io/docs/practices/naming/)
- [Things to watch out for when instrumenting Prometheus collectors](https://prometheus.io/docs/practices/instrumentation/#things-to-watch-out-for) 
- [When to use the Pushgateway](https://prometheus.io/docs/practices/pushing/)
