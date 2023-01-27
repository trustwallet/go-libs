package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/trustwallet/go-libs/metrics"
)

const labelPath = "path"
const labelMethod = "method"

func MetricsMiddleware(namespace string, labels prometheus.Labels, reg prometheus.Registerer) gin.HandlerFunc {
	perfMetric := metrics.NewHttpServerMetric(namespace, []string{labelPath, labelMethod}, labels, reg)
	return func(c *gin.Context) {
		path := c.FullPath()
		method := c.Request.Method

		// route not found, call next and immediately return
		if path == "" {
			c.Next()
			return
		}

		labelValues := []string{path, method}

		startTime := perfMetric.Start(labelValues...)
		c.Next()
		perfMetric.Duration(startTime, labelValues...)

		statusCode := c.Writer.Status()
		switch {
		case 200 <= statusCode && statusCode <= 299:
			perfMetric.Success(labelValues...)
		case 500 <= statusCode && statusCode <= 599:
			perfMetric.ServerError(labelValues...)
		case 400 <= statusCode && statusCode <= 499:
			perfMetric.ClientError(labelValues...)
		}
	}
}
