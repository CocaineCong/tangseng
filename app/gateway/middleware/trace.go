package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"sync/atomic"
	"time"
)

// TODO 若是想用x-request-id 方式连通链路, 也可看middleware下的trace, 按需求来选, 推荐官方的open-telemetry的方式 因jaeger-client已被舍弃
var (
	version string
	incrNum uint64
	pid     = os.Getpid()
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-Id")
		// use pid as traceId
		if traceID == "" {
			traceID = NewTraceID()
		}

		ctx := context.WithValue(c, "trace-id", traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Trace-Id", traceID)

		c.Next()
	}
}

// NewTraceID New trace id
func NewTraceID() string {
	return fmt.Sprintf("trace-id-%d-%s-%d",
		os.Getpid(),
		time.Now().Format("2006.01.02.15.04.05.999"),
		atomic.AddUint64(&incrNum, 1))
}
