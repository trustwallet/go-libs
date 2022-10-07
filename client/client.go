package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Request struct {
	BaseURL          string
	Headers          map[string]string
	HttpClient       HTTPClient
	HttpErrorHandler HttpErrorHandler
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

func (r *Request) GetWithContext(ctx context.Context, result interface{}, path string, query url.Values) error {
	uri := r.GetURL(path, query)
	return r.Execute(ctx, "GET", uri, nil, result)
}

func (r *Request) Get(result interface{}, path string, query url.Values) error {
	return r.GetWithContext(context.Background(), result, path, query)
}

func (r *Request) Post(result interface{}, path string, body interface{}) error {
	return r.PostWithContext(context.Background(), result, path, body)
}

func (r *Request) GetRaw(path string, query url.Values) ([]byte, error) {
	uri := r.GetURL(path, query)
	return r.ExecuteRaw(context.Background(), "GET", uri, nil)
}

func (r *Request) PostRaw(path string, body interface{}) ([]byte, error) {
	buf, err := GetBody(body)
	if err != nil {
		return nil, err
	}
	uri := r.GetBase(path)

	return r.ExecuteRaw(context.Background(), "POST", uri, buf)
}

func (r *Request) PostWithContext(ctx context.Context, result interface{}, path string, body interface{}) error {
	buf, err := GetBody(body)
	if err != nil {
		return err
	}
	uri := r.GetBase(path)
	return r.Execute(ctx, "POST", uri, buf, result)
}

func (r *Request) Execute(ctx context.Context, method string, url string, body io.Reader, result interface{}) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	b, err := r.execute(ctx, req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, result)
	if err != nil {
		return err
	}

	return nil
}

func (r *Request) ExecuteRaw(ctx context.Context, method string, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	return r.execute(ctx, req)
}

func (r *Request) execute(ctx context.Context, req *http.Request) ([]byte, error) {
	c := r.HttpClient

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	err = r.HttpErrorHandler(res, req.URL.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
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
	b, err := ioutil.ReadAll(res.Body)
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

func GetBody(body interface{}) (buf io.ReadWriter, err error) {
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
	}
	return
}
