package interceptor

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/weblazy/easy/utils/etrace"
)

func Trace(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// otel trace

		traceId := etrace.ExtractTraceID(c.Request.Context())
		if c.Request.Header.Get("Traceid") != "" {
			traceId = c.Request.Header.Get("Traceid")
		}
		// 服务内部生成
		if traceId == "" {
			traceId = uuid.NewString()
		}
		ctx = context.WithValue(ctx, "trace_id", traceId)
	}

}
