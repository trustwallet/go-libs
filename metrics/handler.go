package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/trustwallet/go-libs/httplib"
)

func InitHandler(engine *gin.Engine, path string) {
	engine.GET(path, gin.WrapH(promhttp.Handler()))
}

func NewMetricsServer(appName string, port string, path string) httplib.Server {
	router := gin.Default()

	prometheus.DefaultRegisterer.Unregister(collectors.NewGoCollector())
	prometheus.DefaultRegisterer.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	InitHandler(router, path)

	return httplib.NewHTTPServer(router, port)
}
