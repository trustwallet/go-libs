package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/trustwallet/go-libs/metrics"
)

func MetricsMiddleware(namespace string, labels prometheus.Labels, reg prometheus.Registerer) gin.HandlerFunc {
	perfMetric := metrics.NewPerformanceMetric(namespace, []string{"request_path"}, labels, reg)
	return func(c *gin.Context) {
		matchedRoute := c.FullPath()

		// if route not found call next and immediately return
		if matchedRoute == "" {
			c.Next()
			return
		}

		startTime := perfMetric.Start(matchedRoute)
		c.Next()
		perfMetric.Duration(startTime, matchedRoute)

		statusCode := c.Writer.Status()
		if successfulHttpStatusCode(statusCode) {
			perfMetric.Success(matchedRoute)
		} else {
			perfMetric.Failure(matchedRoute)
		}
	}
}

func successfulHttpStatusCode(statusCode int) bool {
	return 200 <= statusCode && statusCode < 300
}
