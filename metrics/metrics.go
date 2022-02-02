package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	executionStartedKey         = "execution_started"
	executionDurationSecondsKey = "execution_duration_seconds"
	executionSucceededTotalKey  = "execution_succeeded_total"
	executionFailedTotalKey     = "execution_failed_total"
)

type Collectors map[string]prometheus.Collector

type PerformanceMetric interface {
	Start() time.Time
	Duration(start time.Time)
	Success()
	Failure()
}

type performanceMetric struct {
	executionStarted         *prometheus.GaugeVec
	executionDurationSeconds *prometheus.GaugeVec
	executionSucceededTotal  *prometheus.CounterVec
	executionFailedTotal     *prometheus.CounterVec
}

func NewPerformanceMetric(namespace string, labels prometheus.Labels, reg prometheus.Registerer) PerformanceMetric {
	executionStarted := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      executionStartedKey,
		Help:      "Last Unix time when execution started.",
	}, nil)

	executionDurationSeconds := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      executionDurationSecondsKey,
		Help:      "Duration of the last execution.",
	}, nil)

	executionSucceededTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      executionSucceededTotalKey,
		Help:      "Total number of the executions wich succeeded.",
	}, nil)

	executionFailedTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      executionFailedTotalKey,
		Help:      "Total number of the executions wich failed.",
	}, nil)

	Register(labels, reg, executionStarted, executionDurationSeconds, executionSucceededTotal, executionFailedTotal)

	return &performanceMetric{
		executionStarted:         executionStarted,
		executionDurationSeconds: executionDurationSeconds,
		executionSucceededTotal:  executionSucceededTotal,
		executionFailedTotal:     executionFailedTotal,
	}
}

func (s *performanceMetric) Start() time.Time {
	start := time.Now()
	s.executionStarted.WithLabelValues().SetToCurrentTime()
	return start
}

func (s *performanceMetric) Duration(start time.Time) {
	duration := time.Since(start)
	s.executionDurationSeconds.WithLabelValues().Set(duration.Seconds())
}

func (s *performanceMetric) Success() {
	s.executionSucceededTotal.WithLabelValues().Inc()
	s.executionFailedTotal.WithLabelValues().Add(0)
}

func (s *performanceMetric) Failure() {
	s.executionFailedTotal.WithLabelValues().Inc()
	s.executionSucceededTotal.WithLabelValues().Add(0)
}

type NullablePerformanceMetric struct{}

func (NullablePerformanceMetric) Start() time.Time {
	return time.Now()
}
func (NullablePerformanceMetric) Duration(start time.Time) {}
func (NullablePerformanceMetric) Success()                 {}
func (NullablePerformanceMetric) Failure()                 {}
