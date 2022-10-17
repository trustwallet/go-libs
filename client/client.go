package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

const defaultTimeout = 5 * time.Second

type Request struct {
	BaseURL          string
	Headers          map[string]string
	HttpClient       HTTPClient
	HttpErrorHandler HttpErrorHandler

	// Monitoring
	metricRegisterer prometheus.Registerer
	httpMetrics      *httpClientMetrics
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpError struct {
	StatusCode int
	URL        url.URL
	Body       []byte
}

func (e *HttpError) Error() string {
	return fmt.Sprintf("Failed request status %d for url: (%s)", e.StatusCode, e.URL.RequestURI())
}

type HttpErrorHandler func(res *http.Response, uri string) error

type Option func(request *Request) error

func InitClient(baseURL string, errorHandler HttpErrorHandler, options ...Option) Request {
	if errorHandler == nil {
		errorHandler = DefaultErrorHandler
	}

	client := Request{
		Headers: make(map[string]string),
		HttpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		HttpErrorHandler: errorHandler,
		BaseURL:          baseURL,
	}

	for _, option := range options {
		err := option(&client)
		if err != nil {
			log.Fatal("Could not initialize http client", err)
		}
	}

	if client.metricsEnabled() {
		err := client.metricRegisterer.Register(client.httpMetrics)
		if err != nil {
			if _, ok := err.(*prometheus.AlreadyRegisteredError); !ok {
				log.WithError(err).Warn("metric already registered")
			} else {
				log.WithError(err).Error("could not initialize http client metrics")
			}
		}
	}

	return client
}

func InitJSONClient(baseUrl string, errorHandler HttpErrorHandler, options ...Option) Request {
	jsonHeaders := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	client := InitClient(
		baseUrl,
		errorHandler,
		append(options, WithExtraHeaders(jsonHeaders))...)
	return client
}

var DefaultErrorHandler = func(res *http.Response, uri string) error {
	return nil
}

// TimeoutOption is an option to set timeout for the http client calls
// value unit is nanoseconds
//
// Deprecated: Internal http.Client shouldn't be modified after construction. Use WithHttpClient instead
func TimeoutOption(timeout time.Duration) Option {
	return func(request *Request) error {
		request.SetTimeout(timeout)

		return nil
	}
}

// Deprecated: Internal http.Client shouldn't be modified after construction. Use WithHttpClient instead
func ProxyOption(proxyURL string) Option {
	return func(request *Request) error {
		if proxyURL == "" {
			return nil
		}

		err := request.SetProxy(proxyURL)
		if err != nil {
			return err
		}

		return nil
	}
}

func WithHttpClient(httpClient HTTPClient) Option {
	return func(request *Request) error {
		request.HttpClient = httpClient
		return nil
	}
}

func WithExtraHeaders(headers map[string]string) Option {
	return func(request *Request) error {
		for k, v := range headers {
			request.Headers[k] = v
		}
		return nil
	}
}

func WithMetricsEnabled(reg prometheus.Registerer, constLabels prometheus.Labels) Option {
	return func(request *Request) error {
		request.httpMetrics = newHttpClientMetrics(constLabels)
		request.metricRegisterer = reg
		return nil
	}
}

// Deprecated: Internal http.Client shouldn't be modified after construction. Use WithHttpClient instead
func (r *Request) SetTimeout(timeout time.Duration) {
	r.HttpClient.(*http.Client).Timeout = timeout
}

// Deprecated: Internal http.Client shouldn't be modified after construction. Use WithHttpClient instead
func (r *Request) SetProxy(proxyUrl string) error {
	if proxyUrl == "" {
		return errors.New("empty proxy url")
	}
	url, err := url.Parse(proxyUrl)
	if err != nil {
		return err
	}
	r.HttpClient.(*http.Client).Transport = &http.Transport{Proxy: http.ProxyURL(url)}
	return nil
}

// Deprecated: Headers shouldn't be modified after construction. Use WithExtraHeaders instead
func (r *Request) AddHeader(key, value string) {
	r.Headers[key] = value
}
