package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

const (
	requestStartedKey         = "request_started"
	requestDurationSecondsKey = "request_duration_seconds"
	requestSucceededTotalKey  = "request_succeeded_total"
	requestClientErrTotalKey  = "request_client_error_total"
	requestServerErrTotalKey  = "request_server_error_total"
)

type HttpServerMetric interface {
	Start(labelValues ...string) time.Time
	Duration(start time.Time, labelValues ...string)
	Success(labelValues ...string)
	ServerError(labelValues ...string)
	ClientError(labelValues ...string)
}

type httpServerMetric struct {
	requestStarted         *prometheus.GaugeVec
	requestDurationSeconds *prometheus.HistogramVec
	requestSucceededTotal  *prometheus.CounterVec
	requestClientErrTotal  *prometheus.CounterVec
	requestServerErrTotal  *prometheus.CounterVec
}

func NewHttpServerMetric(
	namespace string,
	labelNames []string,
	reg prometheus.Registerer,
	labels ...Label,
) HttpServerMetric {
	requestStarted := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      requestStartedKey,
		Help:      "Last Unix time when request started.",
	}, labelNames)

	requestDurationSeconds := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      requestDurationSecondsKey,
		Help:      "Duration of the executions.",
	}, labelNames)

	requestSucceededTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      requestSucceededTotalKey,
		Help:      "Total number of the 2xx requests which succeeded.",
	}, labelNames)

	requestClientErrTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      requestClientErrTotalKey,
		Help:      "Total number of the 4xx requests.",
	}, labelNames)

	requestServerErrTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      requestServerErrTotalKey,
		Help:      "Total number of the 5xx requests.",
	}, labelNames)

	staticLabels := make(map[string]string)
	for _, label := range labels {
		staticLabels[label.Key] = label.Value
	}

	Register(staticLabels, reg, requestStarted, requestDurationSeconds, requestSucceededTotal, requestClientErrTotal, requestServerErrTotal)

	return &httpServerMetric{
		requestStarted:         requestStarted,
		requestDurationSeconds: requestDurationSeconds,
		requestSucceededTotal:  requestSucceededTotal,
		requestClientErrTotal:  requestClientErrTotal,
		requestServerErrTotal:  requestServerErrTotal,
	}
}

func (m *httpServerMetric) Start(labelValues ...string) time.Time {
	start := time.Now()
	m.requestStarted.WithLabelValues(labelValues...).SetToCurrentTime()
	return start
}

func (m *httpServerMetric) Duration(start time.Time, labelValues ...string) {
	duration := time.Since(start)
	m.requestDurationSeconds.WithLabelValues(labelValues...).Observe(duration.Seconds())
}

func (m *httpServerMetric) Success(labelValues ...string) {
	m.requestSucceededTotal.WithLabelValues(labelValues...).Inc()
	m.requestServerErrTotal.WithLabelValues(labelValues...).Add(0)
	m.requestClientErrTotal.WithLabelValues(labelValues...).Add(0)
}

func (m *httpServerMetric) ServerError(labelValues ...string) {
	m.requestSucceededTotal.WithLabelValues(labelValues...).Add(0)
	m.requestServerErrTotal.WithLabelValues(labelValues...).Inc()
	m.requestClientErrTotal.WithLabelValues(labelValues...).Add(0)
}

func (m *httpServerMetric) ClientError(labelValues ...string) {
	m.requestSucceededTotal.WithLabelValues(labelValues...).Add(0)
	m.requestServerErrTotal.WithLabelValues(labelValues...).Add(0)
	m.requestClientErrTotal.WithLabelValues(labelValues...).Inc()
}
