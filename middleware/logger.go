package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Logger(skipPaths ...string) gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: LoggerFormatter(),
		SkipPaths: skipPaths,
	})
}

func LoggerFormatter() gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}
}
