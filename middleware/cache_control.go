package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func CacheControl(duration time.Duration, handle gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()
		cacheControlValue := uint(duration.Seconds())
		c.Writer.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", cacheControlValue))
		handle(c)
	}
}
