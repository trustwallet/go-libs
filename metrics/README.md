# metrics package

Add dependency to the project

```sh
go get github.com/trustwallet/go-libs/metrics
```

## Features

* The `handler.go` contains very simple method to register Prometheus middleware with `gin-gonic` engine.
* The `metrics.go` is a place for the generic metrics services.
  * `PerformanceMetric` allows to track generic job performance, start time,
     duration, success and failed executions.
* The `register.go` contains function which allows registering collectors with custom scope (prometheus labels or `nil`)
  and target Prometheus Registerer instance (could be a Default Registerer).
* The `pusher.go` configures Prometheus Pushgateway client (see Push mode below).
* The `go-libs/worker` and `go-libs/mq` packages integration for automatic job performance tracking.

## How it works?

When we think about collecting system metrics there are several approaches available:

* Case 1 - There a system (service) that might have been built without Prometheus in mind
* Case 2 - There a custom Prometheus Collector
* Case 3 - There is a service that is Prometheus aware, but manages inner collectors

### Case 1

For this case the custom Prometheus Collector is recommended which wraps 
existing service and invokes its functions on `Collect()` method being called.

See the [example](https://github.com/prometheus/client_golang/blob/main/prometheus/example_clustermanager_test.go) in the [prometheus/client_golang](https://github.com/prometheus/client_golang) repo.

The `go-libs/metrics` package offers the `metrics.Register(namespace, labels, registerer)` method
for easier registration of custom collectors.

### Case 2

This case is pretty much the same as Case 1, but the collector itself is aware of the 
logic required to collect and set collectors' values.

<details> 
<summary>Example of the APIMetricsCollector</summary><p> 


```go
// Descriptor used by the APIMetricsCollector below.
var tickersCachedTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "market",
		Subsystem: "api",
		Name:      "tickers_cached_total",
		Help:      "Total number of tickers cached",
	}, nil)

type APIMetricsCollector struct {
	tickersCache cache.Data
}

func NewAPIMetricsCollector(tickersCache cache.Data) *APIMetricsCollector {
	return &APIMetrics{
		tickersCache: tickersCache,
	}
}

// Describe is implemented with DescribeByCollect. That's possible because the
// Collect method will always return the same metrics with the same descriptors.
func (c APIMetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect first triggers the internal collect() to fetch data from the underlying services
// and update set collectors' values.
// Then it simply executes collector.Collect().
func (c APIMetricsCollector) Collect(ch chan<- prometheus.Metric) {
	c.log.Info("collect market api metrics")
	err := c.collect()
	if err != nil {
		log.Error(err)
	}

	tickersCachedTotal.Collect(ch)
}

func (c *APIMetricsCollector) collect() error {
	log.Info("export market api metrics")

	tCached, err := s.tickersCache.GetAllTickers()
	if err != nil {
		return errors.Wrap(err, "failed to get all cached tickers")
	}

	tickersCachedTotal.WithLabelValues().Set(float64(len(tCached)))
	return nil
}
```
</p></details><br/>
 

The instance of this collector has to be registered with Prometheus client.

```go
// main.go

func initMetrics(tickersCache *memory.DataInstance) {
	// disable default go collector which produces a lot of noise
	prometheus.Unregister(collectors.NewGoCollector())

	// register prometheus http handler
	metrics.InitHandler(engine, "/metrics")

	// register collector
	metrics.Register(nil, prometheus.DefaultRegisterer, api.NewAPIMetricsCollector(tickersCache))
}
```

### Case 3 

The example of the 3rd case is the `PerformanceMetric` which is delivered as part
of this package. 

Internally it manages several collectors and registers itself with passed Prometheus Registerer. The passed `prometheus.labels` allow to initialize
multiple instances of the service to track different target jobs execution
without a collision (from Prometheus perspective).


The initialization of the metric
```go
metric := metrics.NewPerformanceMetric(
	"market_api",
	prometheus.Labels{"module": "tickers"},
	prometheus.DefaultRegisterer,
)
```

The usage of the metric service
```go
func (j *job) SomeJob() error {
	defer j.metric.Duration(j.metric.Start())

	err := doSomeWork()
	if err != nil {
		j.metric.Failure()
		return err
	}
	j.metric.Success()
}
```

## Scrape Mode

The `web` applications that are hosted to serve incoming requests usually
have the `/metrics` or similar endpoint exposed.

One extra line in the main application logic to expose the registered
collectors for Prometheus scrapper.

```go
engine := gin.New()

// register prometheus http handler
metrics.InitHandler(engine, "/metrics")
```

## Push Mode

The `worker` applications are launched without a capability to serve the incoming requests,
thus `/metrics` endpoint (Prometheus handler) cannot be utilized there. Instead,
the Prometheus Pushgateway should be used, an intermediary service which allows to push metrics
from jobs which cannot be scraped directly.

Assuming, the `PerformanceMetric` (see above) is configured to capture
worker performance.

To make the metric collectors push the values to Prometheus Pushgateway server,
need to initialize the Pushgateway client and set up the worker which pushes 
registered collectors' values.

```go
func initMetrics() (worker.Worker, error) {
	pusher := metrics.NewPusher(pushgatewayURL, "market_worker")

	// check connection to pusher
	err := pusher.Push()
	if err != nil {
		log.WithError(err).
			Error("cannot connect to pushgateway, metrics won't be pushed")
		return err, nil
	}

	return worker.InitWorker(
		"metrics_pusher",
		worker.DefaultWorkerOptions(pushInterval),
		pusher.Push,
	), nil
}

// start metrics_pusher worker
initMetrics().Start(ctx, wg)
```

📎 When pushing the collectors' values to Pushgateway the `instance` label is 
automatically set from `DYNO` (set by Heroku) or `INSTANCE_ID` (generic variable that
can be set easily) environment variables; otherwise instance is `local`.

## go-libs/worker integration

The latest `worker` package has already integrated with `metrics` package 
and tracks worker function performance automatically.

The following code is already part of the `worker` package:
```go
func (w *worker) invoke() {
	metric := w.options.PerformanceMetric
	
	// no-op perf metric
	if metric == nil {
		metric = &metrics.NullablePerformanceMetric{}
	}

	// collect worker start time and duration on func return
	defer metric.Duration(metric.Start()) 
	
	// invoke worker function
	err := w.workerFn()  		

	if err != nil {
		// increment failure counter on error
		metric.Failure()
		log.WithField("worker", w.name).Error(err)
	} else {
		// increment success counter
		metric.Success(lvs)
	}
}
```

If the `PerformanceMetric` option wasn't initialized it does no-op.

To initialize the worker performance tracking use the following code:

```go
options := worker.DefaultWorkerOptions(interval)

if metricsEnabled {
	options.PerformanceMetric = metrics.NewPerformanceMetric(
		"market_worker",
		prometheus.Labels{"worker": workerName},
		prometheus.DefaultRegisterer,
	)
}

worker := worker.InitWorker(WorkerName,	options, service.DoWork)

ctx, cancel := context.WithCancel(context.Background())
wg := &sync.WaitGroup{}
worker.Start(ctx, waitGroup)
```

This example assumes either Prometheus handler was initialized for scrape mode
or Pushgateway `pusher` was started with `metrics_pusher` worker (see above).

## go-libs/mq integration

The latest `mq` package has already integrated with `metrics` package 
and tracks functions processing incoming messages performance automatically.

The following code is already part of the `mq` package:

```go
func (c *consumer) process(queueName string, body []byte) error {
	metric := c.options.PerformanceMetric
	
	// no-op perf metric
	if metric == nil {
		metric = &metrics.NullablePerformanceMetric{}
	}

	// collect worker start time and duration on func return
	defer metric.Duration(metric.Start())

	// invoke message process function
	err := c.fn(body)

	if err != nil {
		// increment failure counter on error
		metric.Failure()
	} else {
		// increment success counter
		metric.Success()
	}

	return err
}
```

The registration is also very similar:

```go
options := mq.DefaultConsumerOptions(workersCount)

if maxRetries > 0 {
	options.MaxRetries = maxRetries
}

if metricsEnabled {
	options.PerformanceMetric = metrics.NewPerformanceMetric(
		"market_consumer",
		prometheus.Labels{"queue_name": string(queueName)},
		prometheus.DefaultRegisterer,
	)
}

mqClient, err := mq.Connect(rabbitmqURL)
if err != nil {
	log.WithError(err).Fatal("failed to init Rabbit MQ client")
}

consumer := mqClient.InitConsumer(queueName, options, service.DoMessageProcessing)

ctx, cancel := context.WithCancel(context.Background())

if err := mqClient.StartConsumers(ctx, consumer); err != nil {
	log.WithError(err).Fatal("failed to start Rabbit MQ consumers")
}
```

## Useful Readings

- [Metric and label naming best practices](https://prometheus.io/docs/practices/naming/)
- [Things to watch out for when instrumenting Prometheus collectors](https://prometheus.io/docs/practices/instrumentation/#things-to-watch-out-for) 
- [When to use the Pushgateway](https://prometheus.io/docs/practices/pushing/)
