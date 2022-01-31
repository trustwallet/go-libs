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
	Start(labelValues []string) (time.Time, []string)
	Duration(start time.Time, labelValues []string)
	Success(labelValues []string)
	Failure(labelValues []string)
}

type performanceMetric struct {
	collectors Collectors
}

func (c Collectors) GaugeVec(key string) *prometheus.GaugeVec {
	if gauge, ok := (map[string]prometheus.Collector)(c)[key]; ok {
		return gauge.(*prometheus.GaugeVec)
	}
	return nil
}

func (c Collectors) CounterVec(key string) *prometheus.CounterVec {
	if gauge, ok := (map[string]prometheus.Collector)(c)[key]; ok {
		return gauge.(*prometheus.CounterVec)
	}
	return nil
}

func NewConsumerPerformanceMetric(namespace string) PerformanceMetric {
	return NewPerformanceMetric(namespace, "consumer_queue")
}

func NewWorkerPerformanceMetric(namespace string) PerformanceMetric {
	return NewPerformanceMetric(namespace, "worker")
}

func NewPerformanceMetric(namespace string, labelNames ...string) PerformanceMetric {
	e := &performanceMetric{
		collectors: map[string]prometheus.Collector{
			executionStartedKey: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      executionStartedKey,
				Help:      "Last Unix time when execution started.",
			}, labelNames),

			executionDurationSecondsKey: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      executionDurationSecondsKey,
				Help:      "Duration of the last execution.",
			}, labelNames),

			executionSucceededTotalKey: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: namespace,
				Name:      executionSucceededTotalKey,
				Help:      "Total number of the executions wich succeeded.",
			}, labelNames),

			executionFailedTotalKey: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: namespace,
				Name:      executionFailedTotalKey,
				Help:      "Total number of the executions wich failed.",
			}, labelNames),
		},
	}

	for _, c := range e.collectors {
		Register(c)
	}

	return e
}

func (s *performanceMetric) Start(lvs []string) (time.Time, []string) {
	start := time.Now()

	s.collectors.GaugeVec(executionStartedKey).WithLabelValues(lvs...).SetToCurrentTime()
	return start, lvs
}

func (s *performanceMetric) Duration(start time.Time, lvs []string) {
	duration := time.Since(start)
	s.collectors.GaugeVec(executionDurationSecondsKey).WithLabelValues(lvs...).Set(duration.Seconds())
}

func (s *performanceMetric) Success(lvs []string) {
	s.collectors.CounterVec(executionSucceededTotalKey).WithLabelValues(lvs...).Inc()
}

func (s *performanceMetric) Failure(lvs []string) {
	s.collectors.CounterVec(executionFailedTotalKey).WithLabelValues(lvs...).Inc()
}

type NullablePerformanceMetric struct{}

func (NullablePerformanceMetric) Start(lvs []string) (time.Time, []string) {
	return time.Now(), nil
}
func (NullablePerformanceMetric) Duration(start time.Time, lvs []string) {}
func (NullablePerformanceMetric) Success(lvs []string)                   {}
func (NullablePerformanceMetric) Failure(lvs []string)                   {}
