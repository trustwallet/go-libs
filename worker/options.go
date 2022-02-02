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

func DefaultWorkerOptions(interval time.Duration) *WorkerOptions {
	return &WorkerOptions{
		Interval:          interval,
		RunImmediately:    true,
		RunConsequently:   false,
		PerformanceMetric: &metrics.NullablePerformanceMetric{},
	}
}
