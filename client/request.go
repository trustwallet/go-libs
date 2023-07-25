package client

import (
	"net/http"
	"net/url"
)

// Req defines a http request. It is named `Req` instead of `Request` because the http client is named `Request`
// Consider renaming the client to other name.
//
// To build this struct, use NewReqBuilder.
type Req struct {
	headers              map[string]string
	resultContainer      any
	method               string
	path                 Path
	query                url.Values
	body                 any
	rawResponseContainer *http.Response

	metricName        string
	pathMetricEnabled bool
}

type ReqBuilder struct {
	req *Req
}

func NewReqBuilder() *ReqBuilder {
	return &ReqBuilder{
		req: &Req{
			headers:           map[string]string{},
			pathMetricEnabled: true,
		},
	}
}

// Headers sets the headers of the http request. Headers will be overwritten in case of duplicates
func (builder *ReqBuilder) Headers(headers map[string]string) *ReqBuilder {
	for k, v := range headers {
		builder.req.headers[k] = v
	}
	return builder
}

func (builder *ReqBuilder) WriteTo(resultContainer any) *ReqBuilder {
	builder.req.resultContainer = resultContainer
	return builder
}

func (builder *ReqBuilder) WriteRawResponseTo(resp *http.Response) *ReqBuilder {
	builder.req.rawResponseContainer = resp
	return builder
}

func (builder *ReqBuilder) Method(method string) *ReqBuilder {
	builder.req.method = method
	return builder
}

// PathStatic sets the path for the request.
// Use PathStatic ONLY if your path doesn't contain any parameters. Otherwise, use Pathf instead
func (builder *ReqBuilder) PathStatic(path string) *ReqBuilder {
	builder.req.path = NewStaticPath(path)
	return builder
}

func (builder *ReqBuilder) Pathf(pathTemplate string, values ...any) *ReqBuilder {
	builder.req.path = NewPath(pathTemplate, values)
	return builder
}

func (builder *ReqBuilder) Query(query url.Values) *ReqBuilder {
	builder.req.query = query
	return builder
}

func (builder *ReqBuilder) Body(body any) *ReqBuilder {
	builder.req.body = body
	return builder
}

func (builder *ReqBuilder) MetricName(name string) *ReqBuilder {
	builder.req.metricName = name
	return builder
}

// pathMetricEnabled is only for internal use, where it is set to false
// in deprecated wrapper functions such as Get, GetWithContext, Post, PostRaw
func (builder *ReqBuilder) pathMetricEnabled(enabled bool) *ReqBuilder {
	builder.req.pathMetricEnabled = enabled
	return builder
}

func (builder *ReqBuilder) Build() *Req {
	copiedReq := *builder.req
	return &copiedReq
}
