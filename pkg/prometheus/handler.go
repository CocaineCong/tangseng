package prometheus

import (
	"net/http"
	"strings"

	"github.com/CocaineCong/tangseng/config"
	"github.com/CocaineCong/tangseng/consts"
	log "github.com/CocaineCong/tangseng/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// GatewayHandler for expose gateway metrics
func GatewayHandler() gin.HandlerFunc {
	EtcdRegister(config.Conf.Server.Metrics, consts.GatewayJobForPrometheus)
	handler := promhttp.Handler()
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// RpcHandler is for launch a http server
// to expose metrics
func RpcHandler(addr string) {
	port := strings.Split(addr, ":")[1]
	http.Handle("/metrics", promhttp.Handler())
	log.LogrusObj.Panic(http.ListenAndServe(":"+port, nil))
}
