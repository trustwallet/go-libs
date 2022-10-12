package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricNameRequestDurationSeconds = "request_duration_seconds"
	metricNameRequestTotal           = "request_total"
	labelNameUrl                     = "url"
	labelNameMethod                  = "method"
	labelNameStatus                  = "status"

	labelValueErr = "error"

	namespaceHttpClient = "httpclient"
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
		}, []string{labelNameUrl, labelNameMethod}),
		requestTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   namespaceHttpClient,
			Name:        metricNameRequestTotal,
			Help:        "Count of total outgoing http requests, with its result status in labels",
			ConstLabels: constLabels,
		}, []string{labelNameUrl, labelNameMethod, labelNameStatus}),
	}

	return m
}

func (metric *httpClientMetrics) observeDuration(req *http.Request, startTime time.Time) {
	url := getHttpReqMetricUrl(req)
	method := req.Method

	metric.durationSeconds.WithLabelValues(url, method).Observe(time.Since(startTime).Seconds())
}

func (metric *httpClientMetrics) observeResult(req *http.Request, resp *http.Response, err error) {
	url := getHttpReqMetricUrl(req)
	method := req.Method
	status := getHttpRespMetricStatus(resp, err)

	metric.requestTotal.WithLabelValues(url, method, status).Inc()
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
	promhttp.Handler()
}

func getHttpReqMetricUrl(req *http.Request) string {
	return fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path)
}

func getHttpRespMetricStatus(resp *http.Response, err error) string {
	if err != nil {
		return labelValueErr
	}
	firstDigit := resp.StatusCode / 100
	return fmt.Sprintf("%dxx", firstDigit)
}

func WithMetricsEnabled(reg prometheus.Registerer, constLabels prometheus.Labels) Option {
	return func(request *Request) error {
		request.httpMetrics = newHttpClientMetrics(constLabels)
		request.metricRegisterer = reg
		return nil
	}
}
