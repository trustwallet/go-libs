package worker

import (
	"time"

	"github.com/trustwallet/go-libs/metrics"
)

type WorkerOptions struct {
	Interval          time.Duration
	RunImmediately    bool
	RunConsequently   bool
	PerformanceMetric metrics.PerformanceMetric
}

func DefaultWorkerOptions(interval time.Duration) WorkerOptions {
	return WorkerOptions{
		Interval:          interval,
		RunImmediately:    true,
		RunConsequently:   false,
		PerformanceMetric: &metrics.NullablePerformanceMetric{},
	}
}

func (o WorkerOptions) WithPerformanceMetric(metric metrics.PerformanceMetric) WorkerOptions {
	o.PerformanceMetric = metric
	return o
}

func (o WorkerOptions) ShouldFinishBeforeNextStart() WorkerOptions {
	o.RunConsequently = true
	return o
}
