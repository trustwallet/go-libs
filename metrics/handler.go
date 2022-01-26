package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Handler(engine *gin.Engine, path string) {
	engine.GET(path, gin.WrapH(promhttp.Handler()))
}
