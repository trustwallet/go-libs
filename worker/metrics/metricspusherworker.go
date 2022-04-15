package metrics

import (
	"github.com/trustwallet/go-libs/metrics"
	"github.com/trustwallet/go-libs/worker"
)

func NewMetricsPusherWorker(options *worker.WorkerOptions, pusher metrics.Pusher) worker.Worker {
	return worker.NewBuilder("metrics_pusher", pusher.Push).
		WithOptions(options).
		WithStop(pusher.Close).
		Build()
}
