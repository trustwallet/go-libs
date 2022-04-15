package metrics

import (
	"github.com/trustwallet/go-libs/metrics"
	"github.com/trustwallet/go-libs/worker"
)

func NewMetricsPusherWorker(options *worker.WorkerOptions, pusher metrics.Pusher) *worker.Worker {
	return worker.InitWorker("metrics_pusher", options, pusher.Push).WithStop(pusher.Close)
}
