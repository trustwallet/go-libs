package middleware

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestMetricsMiddleware(t *testing.T) {
	r := prometheus.NewRegistry()
	router := gin.New()
	router.Use(MetricsMiddleware("", nil, r))

	successGroup := router.Group("/success")
	successGroup.GET("/:test", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct{}{})
	})

	successGroup.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct{}{})
	})

	router.GET("/error", func(c *gin.Context) {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("oops error"))
	})
	router.GET("/404", func(c *gin.Context) {
		_ = c.AbortWithError(http.StatusNotFound, errors.New("404"))
	})

	// 2 successes, 1 errors
	_ = performRequest("GET", "/success?haha=1&hoho=2", router)
	_ = performRequest("GET", "/error?hehe=1&huhu=3", router)
	_ = performRequest("GET", "/success/hihi", router)
	_ = performRequest("GET", "/404", router)

	metricFamilies, err := r.Gather()
	require.NoError(t, err)
	const (
		requestSucceededTotalKey = "request_succeeded_total"
		requestClientErrTotalKey = "request_client_error_total"
		requestServerErrTotalKey = "request_server_error_total"
	)
	// metricFamily.Name --> label --> counter value
	expected := map[string]map[string]int{
		requestSucceededTotalKey: {
			"/success":       1,
			"/success/:test": 1,
			"/error":         0,
			"/404":           0,
		},
		requestServerErrTotalKey: {
			"/success":       0,
			"/success/:test": 0,
			"/error":         1,
			"/404":           0,
		},
		requestClientErrTotalKey: {
			"/success":       0,
			"/success/:test": 0,
			"/error":         0,
			"/404":           1,
		},
	}
	for _, metricFamily := range metricFamilies {
		expectedLabelCounterMap, ok := expected[*metricFamily.Name]
		if !ok {
			continue
		}
		require.Len(t, metricFamily.Metric, len(expectedLabelCounterMap))
		for _, metric := range metricFamily.Metric {
			require.Len(t, metric.Label, 3)
			labelIndexes := map[string]int{
				labelMethod: -1,
				labelPath:   -1,
				labelStatus: -1,
			}
			for idx, label := range metric.Label {
				labelIndexes[*label.Name] = idx
			}
			require.Equal(t, len(labelIndexes), 3)
			for _, labelIdx := range labelIndexes {
				require.NotEqual(t, -1, labelIdx)
			}
			pathIdx := labelIndexes[labelPath]
			path := *metric.Label[pathIdx].Value
			expectedPathMetric := float64(expectedLabelCounterMap[path])
			require.Equal(t, expectedPathMetric, *metric.Counter.Value)
		}
	}
}
