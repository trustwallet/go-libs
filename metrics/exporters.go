package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	jobStartedKey         = "job_started"
	jobDurationSecondsKey = "job_duration_seconds"
)

var (
	jobPerfExporterLocker  = sync.Mutex{}
	jobPerformanceExporter *JobPerformanceExporter
)

type JobPerformanceExporter struct {
	metrics map[string]*prometheus.GaugeVec
}

func NewJobPerformanceExporter(namespace string, labelNames ...string) *JobPerformanceExporter {
	e := &JobPerformanceExporter{
		metrics: map[string]*prometheus.GaugeVec{
			jobStartedKey: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      jobStartedKey,
				Help:      "Last Unix time when job started.",
			}, labelNames),

			jobDurationSecondsKey: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      jobDurationSecondsKey,
				Help:      "Duration of the last executed job.",
			}, labelNames),
		},
	}

	for _, c := range e.metrics {
		Register(c)
	}

	return e
}

func (s *JobPerformanceExporter) Start(labelValues ...string) (time.Time, []string) {
	start := time.Now()

	s.metrics[jobStartedKey].WithLabelValues(labelValues...).SetToCurrentTime()
	return start, labelValues
}

func (s *JobPerformanceExporter) Duration(start time.Time, labelValues []string) {
	duration := time.Since(start)
	s.metrics[jobDurationSecondsKey].WithLabelValues(labelValues...).Set(duration.Seconds())
}

func GetJobPerformanceExporter(namespace string, labelNames ...string) *JobPerformanceExporter {
	jobPerfExporterLocker.Lock()
	defer jobPerfExporterLocker.Unlock()

	if jobPerformanceExporter == nil {
		jobPerformanceExporter = NewJobPerformanceExporter(namespace, labelNames...)
	}
	return jobPerformanceExporter
}
