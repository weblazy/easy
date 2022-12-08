package interceptor

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/weblazy/easy/etrace"
	"github.com/weblazy/easy/transport"
)

func Trace(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set(transport.PrefixPass+"traceid", etrace.ExtractTraceID(c.Request.Context()))
	}
}
