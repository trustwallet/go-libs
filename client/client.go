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
	Host             string
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
	return fmt.Sprintf(
		"Failed request status %d for url: (%s), body: (%s)",
		e.StatusCode,
		e.URL.RequestURI(),
		string(e.Body),
	)
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
			if _, ok := err.(*prometheus.AlreadyRegisteredError); ok {
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
func TimeoutOption(timeout time.Duration) Option {
	return func(request *Request) error {
		httpClient, ok := request.HttpClient.(*http.Client)
		if !ok {
			return errors.New("unable to set timeout: httpclient is not *http.Client")
		}

		httpClient.Timeout = timeout
		return nil
	}
}

func ProxyOption(proxyURL string) Option {
	return func(request *Request) error {
		if proxyURL == "" {
			return nil
		}

		httpClient, ok := request.HttpClient.(*http.Client)
		if !ok {
			return errors.New("unable to set proxy: httpclient is not *http.Client")
		}

		return setHttpClientTransportProxy(httpClient, proxyURL)
	}
}

func WithHttpClient(httpClient HTTPClient) Option {
	return func(request *Request) error {
		request.HttpClient = httpClient
		return nil
	}
}

func WithExtraHeader(key, value string) Option {
	return func(request *Request) error {
		request.Headers[key] = value
		return nil
	}
}

func WithHost(host string) Option {
	return func(request *Request) error {
		request.Host = host
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

func setHttpClientTransportProxy(client *http.Client, proxyUrl string) error {
	if proxyUrl == "" {
		return errors.New("empty proxy url")
	}
	url, err := url.Parse(proxyUrl)
	if err != nil {
		return err
	}

	if client.Transport == nil {
		client.Transport = &http.Transport{Proxy: http.ProxyURL(url)}
		return nil
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		return errors.New("http client transport is not *http.Transport")
	}
	transport.Proxy = http.ProxyURL(url)
	return nil
}
