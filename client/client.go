package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

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
			Timeout: time.Second * 15,
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
	client := InitClient(baseUrl, errorHandler, options...)
	client.Headers = map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	return client
}

var DefaultErrorHandler = func(res *http.Response, uri string) error {
	return nil
}

// TimeoutOption is an option to set timeout for the http client calls
// value unit is nanoseconds
func TimeoutOption(timeout time.Duration) Option {
	return func(request *Request) error {
		request.SetTimeout(timeout)

		return nil
	}
}

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

func (r *Request) SetTimeout(timeout time.Duration) {
	r.HttpClient.(*http.Client).Timeout = timeout
}

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

func (r *Request) AddHeader(key, value string) {
	r.Headers[key] = value
}

func (r *Request) GetWithContext(ctx context.Context, result interface{}, path Path, query url.Values) error {
	uri := r.GetURL(path.String(), query)
	return r.Execute(ctx, "GET", uri, path.template, nil, result)
}

func (r *Request) Get(result interface{}, path Path, query url.Values) error {
	return r.GetWithContext(context.Background(), result, path, query)
}

func (r *Request) Post(result interface{}, path Path, body interface{}) error {
	return r.PostWithContext(context.Background(), result, path, body)
}

func (r *Request) GetRaw(path Path, query url.Values) ([]byte, error) {
	uri := r.GetURL(path.String(), query)
	return r.ExecuteRaw(context.Background(), "GET", uri, path.template, nil)
}

func (r *Request) PostRaw(path Path, body interface{}) ([]byte, error) {
	buf, err := GetBody(body)
	if err != nil {
		return nil, err
	}
	uri := r.GetBase(path.String())

	return r.ExecuteRaw(context.Background(), "POST", uri, path.template, buf)
}

func (r *Request) PostWithContext(ctx context.Context, result interface{}, path Path, body interface{}) error {
	buf, err := GetBody(body)
	if err != nil {
		return err
	}
	uri := r.GetBase(path.String())
	return r.Execute(ctx, "POST", uri, path.template, buf, result)
}

func (r *Request) Execute(ctx context.Context, method string, url, pathTemplate string, body io.Reader, result interface{}) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	b, err := r.execute(ctx, req, pathTemplate)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, result)
	if err != nil {
		return err
	}

	return nil
}

func (r *Request) ExecuteRaw(ctx context.Context, method string, url, pathTemplate string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	return r.execute(ctx, req, pathTemplate)
}

func (r *Request) execute(ctx context.Context, req *http.Request, reqPathTemplate string) ([]byte, error) {
	c := r.HttpClient

	startTime := time.Now()
	res, err := c.Do(req.WithContext(ctx))

	if r.metricsEnabled() {
		r.httpMetrics.observeDuration(req, reqPathTemplate, startTime)
		r.httpMetrics.observeResult(req, reqPathTemplate, res, err)
	}

	if err != nil {
		return nil, err
	}

	err = r.HttpErrorHandler(res, req.URL.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, &HttpError{
			StatusCode: res.StatusCode,
			URL:        *res.Request.URL,
			Body:       body,
		}
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *Request) GetBase(path string) string {
	baseURL := strings.TrimRight(r.BaseURL, "/")
	if path == "" {
		return baseURL
	}
	path = strings.TrimLeft(path, "/")
	return fmt.Sprintf("%s/%s", baseURL, path)
}

func (r *Request) GetURL(path string, query url.Values) string {
	baseURL := r.GetBase(path)
	if query == nil {
		return baseURL
	}
	queryStr := query.Encode()
	return fmt.Sprintf("%s?%s", baseURL, queryStr)
}

func (r *Request) metricsEnabled() bool {
	return r.httpMetrics != nil
}

func GetBody(body interface{}) (buf io.ReadWriter, err error) {
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
	}
	return
}
