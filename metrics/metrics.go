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
	executionDurationSeconds *prometheus.HistogramVec
	executionSucceededTotal  *prometheus.CounterVec
	executionFailedTotal     *prometheus.CounterVec
}

func NewPerformanceMetric(namespace string, labels prometheus.Labels, reg prometheus.Registerer) PerformanceMetric {
	executionStarted := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      executionStartedKey,
		Help:      "Last Unix time when execution started.",
	}, nil)

	executionDurationSeconds := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      executionDurationSecondsKey,
		Help:      "Duration of the executions.",
	}, nil)

	executionSucceededTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      executionSucceededTotalKey,
		Help:      "Total number of the executions which succeeded.",
	}, nil)

	executionFailedTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      executionFailedTotalKey,
		Help:      "Total number of the executions which failed.",
	}, nil)

	Register(labels, reg, executionStarted, executionDurationSeconds, executionSucceededTotal, executionFailedTotal)

	return &performanceMetric{
		executionStarted:         executionStarted,
		executionDurationSeconds: executionDurationSeconds,
		executionSucceededTotal:  executionSucceededTotal,
		executionFailedTotal:     executionFailedTotal,
	}
}

func (m *performanceMetric) Start() time.Time {
	start := time.Now()
	m.executionStarted.WithLabelValues().SetToCurrentTime()
	return start
}

func (m *performanceMetric) Duration(start time.Time) {
	duration := time.Since(start)
	m.executionDurationSeconds.WithLabelValues().Observe(duration.Seconds())
}

func (m *performanceMetric) Success() {
	m.executionSucceededTotal.WithLabelValues().Inc()
	m.executionFailedTotal.WithLabelValues().Add(0)
}

func (m *performanceMetric) Failure() {
	m.executionFailedTotal.WithLabelValues().Inc()
	m.executionSucceededTotal.WithLabelValues().Add(0)
}

type NullablePerformanceMetric struct{}

func (NullablePerformanceMetric) Start() time.Time {
	// NullablePerformanceMetric is a no-op, so returning empty value
	return time.Time{}
}
func (NullablePerformanceMetric) Duration(_ time.Time) {}
func (NullablePerformanceMetric) Success()             {}
func (NullablePerformanceMetric) Failure()             {}
