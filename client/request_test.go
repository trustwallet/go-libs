package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequest_WriteRawResponseTo(t *testing.T) {
	const (
		pathOk  = "/ok"
		path5xx = "/5xx"
	)

	tests := []struct {
		name       string
		path       string
		statusCode int
		headers    http.Header
	}{
		{
			name:       "Test write raw response with statusOK",
			path:       pathOk,
			statusCode: http.StatusOK,
			headers: http.Header{
				"Content-Type":         []string{"application/json"},
				"x-aptos-block-height": []string{"73287085"},
			},
		},
		{
			name:       "Test write raw response with status5xx",
			path:       path5xx,
			statusCode: http.StatusInternalServerError,
			headers: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := InitClient("http://www.example.com", nil,
				WithHttpClient(&http.Client{
					Transport: RoundTripperFunc(func(request *http.Request) (*http.Response, error) {
						switch request.URL.Path {
						case pathOk:
							return &http.Response{
								StatusCode: http.StatusOK,
								Body:       io.NopCloser(strings.NewReader(`{"Data": "ok"}`)),
								Header:     tt.headers,
							}, nil
						case path5xx:
							return &http.Response{
								StatusCode: http.StatusInternalServerError,
								Request:    request,
								Body:       io.NopCloser(strings.NewReader(`{"Data": "5xx"}`)),
								Header:     tt.headers,
							}, nil
						default:
							return nil, nil
						}
					}),
				}),
			)
			var resp http.Response
			_, _ = client.Execute(context.Background(), NewReqBuilder().Method(http.MethodGet).PathStatic(tt.path).WriteRawResponseTo(&resp).Build())
			require.Equal(t, tt.headers, resp.Header)
			require.Equal(t, tt.statusCode, resp.StatusCode)
		})
	}
}
