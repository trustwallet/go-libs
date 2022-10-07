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

	// 2 successes, 1 errors
	_ = performRequest("GET", "/success?haha=1&hoho=2", router)
	_ = performRequest("GET", "/error?hehe=1&huhu=3", router)
	_ = performRequest("GET", "/success/hihi", router)

	metricFamilies, err := r.Gather()
	require.NoError(t, err)

	const executionFailedTotal = "execution_failed_total"
	const executionSucceededTotal = "execution_succeeded_total"

	// metricFamily.Name --> label --> counter value
	expected := map[string]map[string]int{
		executionSucceededTotal: {
			"/success":       1,
			"/success/:test": 1,
			"/error":         0,
		},
		executionFailedTotal: {
			"/success":       0,
			"/success/:test": 0,
			"/error":         1,
		},
	}

	for _, metricFamily := range metricFamilies {
		expectedLabelCounterMap, ok := expected[*metricFamily.Name]
		if !ok {
			continue
		}

		require.Len(t, metricFamily.Metric, len(expectedLabelCounterMap))
		for _, metric := range metricFamily.Metric {
			require.Len(t, metric.Label, 1)
			require.Equal(t, float64(expectedLabelCounterMap[*metric.Label[0].Value]), *metric.Counter.Value)
		}
	}
}
