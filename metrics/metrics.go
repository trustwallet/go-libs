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
	Start(labelValues ...string) time.Time
	Duration(start time.Time, labelValues ...string)
	Success(labelValues ...string)
	Failure(labelValues ...string)
}

type performanceMetric struct {
	executionStarted         *prometheus.GaugeVec
	executionDurationSeconds *prometheus.HistogramVec
	executionSucceededTotal  *prometheus.CounterVec
	executionFailedTotal     *prometheus.CounterVec
}

type metricLabel struct {
	Key string
	Value string
}

func NewPerformanceMetric(
	namespace string,
	labelNames []string,
	reg prometheus.Registerer,
	labels ...metricLabel,
) PerformanceMetric {
	executionStarted := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      executionStartedKey,
		Help:      "Last Unix time when execution started.",
	}, labelNames)

	executionDurationSeconds := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      executionDurationSecondsKey,
		Help:      "Duration of the executions.",
	}, labelNames)

	executionSucceededTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      executionSucceededTotalKey,
		Help:      "Total number of the executions which succeeded.",
	}, labelNames)

	executionFailedTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      executionFailedTotalKey,
		Help:      "Total number of the executions which failed.",
	}, labelNames)

	staticLabels := make(map[string]string)
	for _, label := range labels {
		staticLabels[label.Key] = label.Value
	}

	Register(staticLabels, reg, executionStarted, executionDurationSeconds, executionSucceededTotal, executionFailedTotal)

	return &performanceMetric{
		executionStarted:         executionStarted,
		executionDurationSeconds: executionDurationSeconds,
		executionSucceededTotal:  executionSucceededTotal,
		executionFailedTotal:     executionFailedTotal,
	}
}

func (m *performanceMetric) Start(labelValues ...string) time.Time {
	start := time.Now()
	m.executionStarted.WithLabelValues(labelValues...).SetToCurrentTime()
	return start
}

func (m *performanceMetric) Duration(start time.Time, labelValues ...string) {
	duration := time.Since(start)
	m.executionDurationSeconds.WithLabelValues(labelValues...).Observe(duration.Seconds())
}

func (m *performanceMetric) Success(labelValues ...string) {
	m.executionSucceededTotal.WithLabelValues(labelValues...).Inc()
	m.executionFailedTotal.WithLabelValues(labelValues...).Add(0)
}

func (m *performanceMetric) Failure(labelValues ...string) {
	m.executionFailedTotal.WithLabelValues(labelValues...).Inc()
	m.executionSucceededTotal.WithLabelValues(labelValues...).Add(0)
}

type NullablePerformanceMetric struct{}

func (NullablePerformanceMetric) Start(_ ...string) time.Time {
	// NullablePerformanceMetric is a no-op, so returning empty value
	return time.Time{}
}
func (NullablePerformanceMetric) Duration(_ time.Time, _ ...string) {}
func (NullablePerformanceMetric) Success(_ ...string)               {}
func (NullablePerformanceMetric) Failure(_ ...string)               {}
