package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func TestClientMetrics(t *testing.T) {
	const (
		pathOk  = "/ok"
		path5xx = "/5xx"
		pathErr = "/err"

		basePath = "/api"
	)

	reg := prometheus.NewPedanticRegistry()

	client := InitClient("http://www.example.com"+basePath, nil,
		WithMetricsEnabled(reg, prometheus.Labels{"app": "test"}),
		WithHttpClient(&http.Client{
			Transport: RoundTripperFunc(func(request *http.Request) (*http.Response, error) {
				switch request.URL.Path {
				case basePath + pathOk:
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"Data": "ok"}`)),
					}, nil
				case basePath + path5xx:
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Request:    request,
						Body:       io.NopCloser(strings.NewReader(`{"Data": "5xx"}`)),
					}, nil
				case basePath + pathErr:
					return nil, errors.New("oops")
				default:
					return nil, nil
				}
			}),
		}))

	_, _ = client.Execute(context.Background(), NewReqBuilder().Method(http.MethodGet).PathStatic(pathOk).Build())
	_, _ = client.Execute(context.Background(), NewReqBuilder().Method(http.MethodGet).PathStatic(path5xx).Build())
	_, _ = client.Execute(context.Background(), NewReqBuilder().Method(http.MethodGet).PathStatic(pathErr).Build())
	_, _ = client.Execute(context.Background(), NewReqBuilder().Method(http.MethodGet).PathStatic(pathErr).Build())

	_, _ = client.Execute(context.Background(), NewReqBuilder().Method(http.MethodPost).PathStatic(path5xx).Build())
	_, _ = client.Execute(context.Background(), NewReqBuilder().Method(http.MethodPost).PathStatic(pathErr).Build())

	_, _ = client.Execute(context.Background(), NewReqBuilder().
		Method(http.MethodPost).
		PathStatic(pathErr).
		MetricName("postError").
		Build())

	// using Pathf
	_, _ = client.Execute(context.Background(), NewReqBuilder().
		Method(http.MethodGet).Pathf("%s", pathErr).Build())
	_, _ = client.Execute(context.Background(), NewReqBuilder().
		Method(http.MethodGet).Pathf("%s", pathOk).Build())
	_, _ = client.Execute(context.Background(), NewReqBuilder().
		Method(http.MethodGet).Pathf("%s", pathOk).Build())
	_, _ = client.Execute(context.Background(), NewReqBuilder().
		Method(http.MethodGet).Pathf("%s", path5xx).Build())

	type Resp struct {
		Data string
	}
	var resp Resp
	ctx := context.Background()
	_, _ = client.Execute(ctx, NewReqBuilder().
		Method(http.MethodGet).
		WriteTo(&resp).
		PathStatic(path5xx).
		Build())
	_, _ = client.Execute(ctx, NewReqBuilder().
		Method(http.MethodGet).
		WriteTo(&resp).
		PathStatic(path5xx).
		Build())
	_, _ = client.Execute(ctx, NewReqBuilder().
		Method(http.MethodGet).
		WriteTo(&resp).
		PathStatic(pathOk).
		Build())
	_, _ = client.Execute(ctx, NewReqBuilder().
		Method(http.MethodGet).
		WriteTo(&resp).
		PathStatic(pathErr).
		Build())

	_, _ = client.Execute(ctx, NewReqBuilder().
		Method(http.MethodPost).
		WriteTo(&resp).
		PathStatic(path5xx).
		Build())
	_, _ = client.Execute(ctx, NewReqBuilder().
		Method(http.MethodPost).
		WriteTo(&resp).
		PathStatic(path5xx).
		Build())
	_, _ = client.Execute(ctx, NewReqBuilder().
		Method(http.MethodPost).
		WriteTo(&resp).
		PathStatic(pathOk).
		Build())
	_, _ = client.Execute(ctx, NewReqBuilder().
		Method(http.MethodPost).
		WriteTo(&resp).
		PathStatic(pathErr).
		Build())

	mfs, err := reg.Gather()
	require.NoError(t, err)
	require.NotNil(t, mfs)

	// metricFamily.Name --> Concat(label_name=label_value) --> counter value
	expected := map[string]map[string]int{
		namespaceHttpClient + "_" + metricNameRequestTotal: {
			"app=test method=GET name= status=2xx url=http://www.example.com/api/ok":    1,
			"app=test method=GET name= status=5xx url=http://www.example.com/api/5xx":   1,
			"app=test method=GET name= status=error url=http://www.example.com/api/err": 2,

			"app=test method=POST name= status=5xx url=http://www.example.com/api/5xx":   1,
			"app=test method=POST name= status=error url=http://www.example.com/api/err": 1,

			"app=test method=POST name=postError status=error url=http://www.example.com/api/err": 1,

			"app=test method=GET name= status=2xx url=http://www.example.com/api":   1,
			"app=test method=GET name= status=5xx url=http://www.example.com/api":   2,
			"app=test method=GET name= status=error url=http://www.example.com/api": 1,

			"app=test method=POST name= status=2xx url=http://www.example.com/api":   1,
			"app=test method=POST name= status=5xx url=http://www.example.com/api":   2,
			"app=test method=POST name= status=error url=http://www.example.com/api": 1,

			// Pathf metrics
			"app=test method=GET name= status=2xx url=http://www.example.com/api/%s":   2,
			"app=test method=GET name= status=5xx url=http://www.example.com/api/%s":   1,
			"app=test method=GET name= status=error url=http://www.example.com/api/%s": 1,
		},
	}

	testedMetricCount := 0
	for _, mf := range mfs {
		expectedLabelCounterMap, ok := expected[*mf.Name]
		if !ok {
			continue
		}
		testedMetricCount++

		require.Len(t, mf.Metric, len(expectedLabelCounterMap))
		for _, metric := range mf.Metric {
			labelNameValues := make([]string, len(metric.Label))
			for idx, label := range metric.Label {
				labelNameValues[idx] = fmt.Sprintf("%s=%s", *label.Name, *label.Value)
			}

			joinedLabels := strings.Join(labelNameValues, " ")
			expectedCounter := float64(expectedLabelCounterMap[joinedLabels])
			require.Equal(t, expectedCounter, *metric.Counter.Value)
		}
	}
	require.Equal(t, len(expected), testedMetricCount, "makes sure all expected metrics are tested")
}

func Test_getHttpRespMetricStatus(t *testing.T) {
	type args struct {
		resp *http.Response
		err  error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "err is not nil, return labelValueErr",
			args: args{
				resp: nil,
				err:  errors.New("oops"),
			},
			want: labelValueErr,
		},
		{
			name: "err is nil, 1xx",
			args: args{
				resp: &http.Response{StatusCode: http.StatusContinue},
				err:  nil,
			},
			want: "1xx",
		},
		{
			name: "err is nil, 2xx",
			args: args{
				resp: &http.Response{StatusCode: http.StatusOK},
				err:  nil,
			},
			want: "2xx",
		},
		{
			name: "err is nil, 3xx",
			args: args{
				resp: &http.Response{StatusCode: http.StatusMovedPermanently},
				err:  nil,
			},
			want: "3xx",
		},
		{
			name: "err is nil, 4xx",
			args: args{
				resp: &http.Response{StatusCode: http.StatusBadRequest},
				err:  nil,
			},
			want: "4xx",
		},
		{
			name: "err is nil, 5xx",
			args: args{
				resp: &http.Response{StatusCode: http.StatusHTTPVersionNotSupported},
				err:  nil,
			},
			want: "5xx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getHttpRespMetricStatus(tt.args.resp, tt.args.err), "getHttpRespMetricStatus(%v, %v)", tt.args.resp, tt.args.err)
		})
	}
}
