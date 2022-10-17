package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespaceHttpClient = "httpclient"

	metricNameRequestDurationSeconds = "request_duration_seconds"
	metricNameRequestTotal           = "request_total"

	labelUrl    = "url"
	labelMethod = "method"
	labelStatus = "status"
	labelName   = "name"

	labelValueErr = "error"
)

type httpClientMetrics struct {
	durationSeconds *prometheus.HistogramVec
	requestTotal    *prometheus.CounterVec
}

func newHttpClientMetrics(constLabels prometheus.Labels) *httpClientMetrics {
	m := &httpClientMetrics{
		durationSeconds: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace:   namespaceHttpClient,
			Name:        metricNameRequestDurationSeconds,
			Help:        "Histogram of duration of outgoing http requests",
			ConstLabels: constLabels,
		}, []string{labelUrl, labelMethod, labelName}),
		requestTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   namespaceHttpClient,
			Name:        metricNameRequestTotal,
			Help:        "Count of total outgoing http requests, with its result status in labels",
			ConstLabels: constLabels,
		}, []string{labelUrl, labelMethod, labelName, labelStatus}),
	}

	return m
}

func (metric *httpClientMetrics) observeDuration(url, method, name string, startTime time.Time) {
	metric.durationSeconds.WithLabelValues(url, method, name).Observe(time.Since(startTime).Seconds())
}

func (metric *httpClientMetrics) observeResult(url, method, name, status string) {
	metric.requestTotal.WithLabelValues(url, method, name, status).Inc()
}

// Describe implements prometheus.Collector interface
func (metric *httpClientMetrics) Describe(descs chan<- *prometheus.Desc) {
	metric.durationSeconds.Describe(descs)
	metric.requestTotal.Describe(descs)
}

// Collect implements prometheus.Collector interface
func (metric *httpClientMetrics) Collect(metrics chan<- prometheus.Metric) {
	metric.durationSeconds.Collect(metrics)
	metric.requestTotal.Collect(metrics)
}

func getHttpReqMetricUrl(req *http.Request, pathTemplate string) string {
	return fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, pathTemplate)
}

func getHttpRespMetricStatus(resp *http.Response, err error) string {
	if err != nil {
		return labelValueErr
	}
	firstDigit := resp.StatusCode / 100
	return fmt.Sprintf("%dxx", firstDigit)
}
