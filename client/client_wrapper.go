package client

import (
	"context"
	"net/http"
	"net/url"
)

// Deprecated: Left as backwards-compatibility. Use Execute(NewReqBuilder()) for better APIs and monitoring
func (r *Request) GetWithContext(ctx context.Context, result interface{}, path string, query url.Values) error {
	_, err := r.Execute(ctx, NewReqBuilder().
		Method(http.MethodGet).
		PathStatic(path).
		Query(query).
		WriteTo(result).
		pathMetricEnabled(false).
		Build())
	return err
}

// Deprecated: Left as backwards-compatibility. Use Execute(NewReqBuilder()) for better APIs and monitoring
func (r *Request) Get(result interface{}, path string, query url.Values) error {
	return r.GetWithContext(context.Background(), result, path, query)
}

// Deprecated: Left as backwards-compatibility. Use Execute(NewReqBuilder()) for better APIs and monitoring
func (r *Request) Post(result interface{}, path string, body interface{}) error {
	return r.PostWithContext(context.Background(), result, path, body)
}

// Deprecated: Left as backwards-compatibility. Use Execute(NewReqBuilder()) for better APIs and monitoring
func (r *Request) GetRaw(path string, query url.Values) ([]byte, error) {
	return r.Execute(context.Background(), NewReqBuilder().
		Method(http.MethodGet).
		PathStatic(path).
		Query(query).
		pathMetricEnabled(false).
		Build())
}

// Deprecated: Left as backwards-compatibility. Use Execute(NewReqBuilder()) for better APIs and monitoring
func (r *Request) PostRaw(path string, body interface{}) ([]byte, error) {
	return r.Execute(context.Background(), NewReqBuilder().
		Method(http.MethodPost).
		PathStatic(path).
		Body(body).
		pathMetricEnabled(false).
		Build())
}

// Deprecated: Left as backwards-compatibility. Use Execute(NewReqBuilder()) for better APIs and monitoring
func (r *Request) PostWithContext(ctx context.Context, result interface{}, path string, body interface{}) error {
	_, err := r.Execute(ctx, NewReqBuilder().
		Method(http.MethodPost).
		PathStatic(path).
		Body(body).
		WriteTo(result).
		pathMetricEnabled(false).
		Build())
	return err
}
