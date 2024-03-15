package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Execute executes http request as described in Req.
//
// If Req.WriteTo is specified, it will also populate the resultContainer
func (r *Request) Execute(ctx context.Context, req *Req) ([]byte, error) {
	request, err := r.constructHttpRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	res, err := r.HttpClient.Do(request)
	r.reportMonitoringMetricsIfEnabled(startTime, request, req, res, err)
	if err != nil {
		return nil, err
	}

	if req.rawResponseContainer != nil && res != nil {
		*req.rawResponseContainer = *res
	}

	err = r.HttpErrorHandler(res, request.URL.String())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return nil, &HttpError{
			StatusCode: res.StatusCode,
			URL:        *request.URL,
			Body:       b,
		}
	}

	err = populateResultContainer(b, req.resultContainer)
	if err != nil {
		return b, err
	}

	return b, nil
}

// constructHttpRequest constructs a http.Request object from description in Req and common headers in r.
func (r *Request) constructHttpRequest(ctx context.Context, req *Req) (*http.Request, error) {
	body, err := GetBody(req.body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, req.method, r.GetURL(req.path.String(), req.query), body)
	if err != nil {
		return nil, err
	}

	r.setRequestHeaders(request, req)

	if r.Host != "" {
		request.Host = r.Host
	}
	return request, nil
}

func (r *Request) reportMonitoringMetricsIfEnabled(
	startTime time.Time, request *http.Request,
	req *Req, res *http.Response, resErr error,
) {
	if r.metricsEnabled() {
		url := r.GetURL(getMonitoredPathTemplateIfEnabled(req), nil)
		method := request.Method
		name := req.metricName
		status := getHttpRespMetricStatus(res, resErr)

		r.httpMetrics.observeDuration(url, method, name, startTime)
		r.httpMetrics.observeResult(url, method, name, status)
	}
}

// setRequestHeaders sets the given httpRequest with the common headers from the client, and headers specified in Req.
// If there are duplicated headers, the headers specified in Req takes precedence.
func (r *Request) setRequestHeaders(httpRequest *http.Request, req *Req) {
	headersSlice := []map[string]string{r.Headers, req.headers}
	for _, headers := range headersSlice {
		for key, value := range headers {
			httpRequest.Header.Set(key, value)
		}
	}
}

// populateResultContainer populates the given resultContainer if it's not nil
func populateResultContainer(b []byte, resultContainer any) error {
	if resultContainer != nil {
		err := json.Unmarshal(b, resultContainer)
		if err != nil {
			return err
		}
	}
	return nil
}

func getMonitoredPathTemplateIfEnabled(req *Req) string {
	if !req.pathMetricEnabled {
		return ""
	}
	return req.path.template
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
	if body == nil {
		return
	}

	var bs []byte
	bs, err = json.Marshal(body)
	if err != nil {
		return
	}

	buf = bytes.NewBuffer(bs)
	return
}
