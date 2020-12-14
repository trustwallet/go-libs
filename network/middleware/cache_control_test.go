package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCacheControl(t *testing.T) {
	router := gin.New()
	router.GET("/cache_ping_control", CacheControl(time.Second*30, func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong "+fmt.Sprint(time.Now().UnixNano()))
	}))

	w1 := performRequest("GET", "/cache_ping_control", router)
	w1CacheControl := w1.Header().Get("Cache-Control")
	assert.NotEqual(t, "no-cache", w1CacheControl)
	time.Sleep(time.Second * 1)
	w2 := performRequest("GET", "/cache_ping_control", router)
	w2CacheControl := w2.Header().Get("Cache-Control")

	assert.Equal(t, w1CacheControl, w2CacheControl)
	assert.Equal(t, "max-age=30", w2CacheControl)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, http.StatusOK, w2.Code)

}
