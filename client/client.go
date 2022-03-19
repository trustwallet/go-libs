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
	HttpClient       *http.Client
	HttpErrorHandler HttpErrorHandler
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

func TimeoutOption(seconds time.Duration) Option {
	return func(request *Request) error {
		request.SetTimeout(seconds)

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

func (r *Request) SetTimeout(seconds time.Duration) {
	r.HttpClient.Timeout = time.Second * seconds
}

func (r *Request) SetProxy(proxyUrl string) error {
	if proxyUrl == "" {
		return errors.New("empty proxy url")
	}
	url, err := url.Parse(proxyUrl)
	if err != nil {
		return err
	}
	r.HttpClient.Transport = &http.Transport{Proxy: http.ProxyURL(url)}
	return nil
}

func (r *Request) AddHeader(key, value string) {
	r.Headers[key] = value
}

func (r *Request) GetWithContext(result interface{}, path string, query url.Values, ctx context.Context) error {
	uri := r.GetURL(path, query)
	return r.Execute("GET", uri, nil, result, ctx)
}

func (r *Request) Get(result interface{}, path string, query url.Values) error {
	uri := r.GetURL(path, query)
	return r.Execute("GET", uri, nil, result, context.Background())
}

func (r *Request) Post(result interface{}, path string, body interface{}) error {
	buf, err := GetBody(body)
	if err != nil {
		return err
	}
	uri := r.GetBase(path)
	return r.Execute("POST", uri, buf, result, context.Background())
}

func (r *Request) GetRaw(path string, query url.Values) ([]byte, error) {
	var result interface{}

	err := r.Get(&result, path, query)
	if err != nil {
		return nil, err
	}

	return ExtractBody(result)
}

func (r *Request) PostRaw(path string, body interface{}) ([]byte, error) {
	var result interface{}

	err := r.Post(&result, path, body)
	if err != nil {
		return nil, err
	}

	return ExtractBody(result)
}

func (r *Request) PostWithContext(result interface{}, path string, body interface{}, ctx context.Context) error {
	buf, err := GetBody(body)
	if err != nil {
		return err
	}
	uri := r.GetBase(path)
	return r.Execute("POST", uri, buf, result, ctx)
}

func (r *Request) Execute(method string, url string, body io.Reader, result interface{}, ctx context.Context) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	c := r.HttpClient

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	err = r.HttpErrorHandler(res, url)
	if err != nil {
		return err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return &HttpError{
			StatusCode: res.StatusCode,
			URL:        *res.Request.URL,
			Body:       body,
		}
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, result)
	if err != nil {
		return err
	}
	return err
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

func ExtractBody(body interface{}) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(buf)
}
